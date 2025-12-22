package st

import (
	"iter"
	"slices"
)

type CurrentColumn[E any] struct {
	ColumnOverlap[E]
	ColumnId int
}

type ColumnSets[E any] struct {
	Util    *SpanUtil[E]
	columns *[]*ColumnOverlapAccumulator[E]
	active  *[]int
	overlap SpanBoundry[E]
	closed  bool
	pos     int
	current *[]*CurrentColumn[E]
	itr     bool
	Err     error
	ErrCol  int
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

func (s *ColumnSets[E]) Close() {
	if s.closed {
		return
	}
	s.closed = true
	if s.columns != nil {
		for _, col := range *s.columns {
			col.Close()
		}
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
	var res = s.AddColumn(ac.NewColumnOverlapAccumulatorFromSpanBoundrySlice(list))
	return res, ac
}

func (s *ColumnSets[E]) AddColumnFromOverlappingSpanSets(list *[]*OverlappingSpanSets[E]) int {
	return s.AddColumn(
		s.Util.NewColumnOverlapAccumulator(
			iter.Pull2(
				slices.All(*list),
			),
		),
	)
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
