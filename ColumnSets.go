package st

import (
	"iter"
	"slices"
)

// Represents a source data set of the culumn consolidation process.
type CurrentColumn[E any] struct {
	ColumnOverlap[E]
	ColumnId int
}

// This struct acts as the majordomo of the constrained span intersection iteration process.
//
// For every instance created make sure to scope the proper defer call:
//
//   defer i.Close()
type ColumnSets[E any] struct {
	Util    *SpanUtil[E]
	columns *[]*ColumnOverlapAccumulator[E]
	active  *[]int
	overlap SpanBoundry[E]
	closed  bool
	pos     int
	current *[]*CurrentColumn[E]
	itr     bool

	// The last error, nil if there were no errors
	Err error
	// ColumnId the our last error came from
	ErrCol  int
	OnClose *[]func()
}

type ColumnResults[E any] interface {
	// Returns the current columns
	GetColumns() *[]*CurrentColumn[E]

	// Denotes how many columns overlap with the current span.
	OverlapCount() int
	// Returns the SpanBoundry representing the current position in our data set.
	GetSpan() SpanBoundry[E]

	SpanBoundry[E]
}

// Returns the SpanBoundry instance that represents the intersection of our current column state.
func (s *ColumnSets[E]) GetSpan() SpanBoundry[E] {
	return s.overlap
}

// This is a wrapper for s.GetSpan.GetBegin().
func (s *ColumnSets[E]) GetBegin() E {
	return s.overlap.GetBegin()
}

// This is a wrapper for s.GetSpan.GetEnd().
func (s *ColumnSets[E]) GetEnd() E {
	return s.overlap.GetEnd()
}

func (s *ColumnSets[E]) GetColumns() *[]*CurrentColumn[E] {
	return s.current
}

// Shuts down and cleans up the instance, and any go routines that were registered.
func (s *ColumnSets[E]) Close() {
	if s.closed {
		return
	}
	s.closed = true
	if s.OnClose != nil {
		for _, todo := range *s.OnClose {
			todo()
		}
	}
	if s.columns != nil {
		for _, col := range *s.columns {
			col.Close()
		}
	}
}

// Adds a function to call when the close operation is called or the iterator optations have been completed.
func (s *ColumnSets[E]) AddOnClose(todo func()) {
	if s.OnClose == nil {
		s.OnClose = &[]func(){todo}
	} else {
		*s.OnClose = append(*s.OnClose, todo)
	}
}

// Denotes how many columns overlap with this span, if the set of current colums is empty
// then the returned value will be -1.
func (s *ColumnSets[E]) OverlapCount() int {
	if s.current == nil {
		return -1
	}
	return len(*s.current)
}

// Appends an a column accumulator to the current column set.
// Returns the id of the column, if the instance is closed returns -1.
func (s *ColumnSets[E]) AddColumn(c *ColumnOverlapAccumulator[E]) int {
	if s.closed {
		return -1
	}
	if s.columns == nil {
		s.columns = &[]*ColumnOverlapAccumulator[E]{}
	}
	*s.columns = append(*s.columns, c)
	return len(*s.columns) - 1
}

// This is a helper method that constructs an SpanOverlapAccumulator and then produces
// an iterator from the SpanOverlapAccumulator based on list.
func (s *ColumnSets[E]) AddColumnFromSpanSlice(list *[]SpanBoundry[E]) (int, *SpanOverlapAccumulator[E]) {
	var ac = s.Util.NewSpanOverlapAccumulator()
	var res = s.AddColumn(ac.NewCoaFromSbSlice(list))
	return res, ac
}

// Adds list as a column to the internals.
func (s *ColumnSets[E]) AddColumnFromOverlappingSpanSets(list *[]*OverlappingSpanSets[E]) int {
	return s.AddColumn(
		s.Util.NewColumnOverlapAccumulator(
			iter.Pull2(
				slices.All(*list),
			),
		),
	)
}

// Adds a context aware channel based ColumnOverlapAccumulator.
//
// # Warning
//
// If you don't start the go routine that appends to the channel before calling this method, 
// it will cause a race condition that will prevent the ColumnSets instancce from working 
// correctly.
func (s *ColumnSets[E]) AddColumnFromNewOlssChanStater(sa *OlssChanStater[E]) int {
	id := s.AddColumn(
		s.Util.NewColumnOverlapAccumulator(
			iter.Pull2(
				s.Util.NewOlssSeq2FromOlssChan(sa.Chan),
			),
		),
	)
	s.AddOnClose(sa.Shutdown)
	return id
}

func (s *ColumnSets[E]) init() {
	var check = []int{}
	var test = &[]SpanBoundry[E]{}

	for i, span := range *s.columns {
		if span.Err != nil {
			s.Err = span.Err
			s.ErrCol = i
			s.pos = -1
			return
		}
		if span.HasNext() {
			check = append(check, i)
			*test = append(*test, span)
		}
	}
	var init, ok = s.Util.FirstSpan(test)
	if !ok {
		s.pos = -1
		return
	}
	s.pos = 0
	s.overlap = init
	s.active = &check
	s.setCurrent()
}

func (s *ColumnSets[E]) setCurrent() {
	s.current = &[]*CurrentColumn[E]{}
	for _, i := range *s.active {
		var col = (*s.columns)[i]
		col.SetNext(s.overlap)
		if col.InOverlap() {
			var res = &CurrentColumn[E]{
				ColumnId:      i,
				ColumnOverlap: col,
			}
			*s.current = append(*s.current, res)
		}
	}
}

func (s *ColumnSets[E]) setNext() {
	var check = []int{}
	var test = &[]SpanBoundry[E]{}

	for i, span := range *s.columns {
		if span.Err != nil {
			s.Err = span.Err
			s.ErrCol = i
			s.pos = -1
			return
		}
		if span.HasNext() {
			check = append(check, i)
			*test = append(*test, span)
		}
	}
	if len(check) == 0 {
		s.pos = -1
		return
	}

	var ok bool
	s.overlap, ok = s.Util.NextSpan(s.overlap, test)
	if !ok {
		s.pos = -1
		return
	}
	s.pos++

	s.active = &check
	s.setCurrent()
}

// Creates an iteraotr to walk all added columns and find the overlaps.
func (s *ColumnSets[E]) Iter() iter.Seq2[int, ColumnResults[E]] {
	if s.itr {
		return nil
	}
	s.itr = true
	s.init()

	return func(yeild func(int, ColumnResults[E]) bool) {
		defer s.Close()
		for s.pos != -1 {
			if !yeild(s.pos, s) {
				return
			}
			s.setNext()
		}
	}
}
