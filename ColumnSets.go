package st

type CurrentColumn[E any] struct {
	ColumnOverlap[E]
	ColumnId int
}
type ColumnSets[E any] struct {
	Util    *SpanUtil[E]
	Columns *[]*ColumnOverlapAccumulator[E]
	Active  *[]int
	Overlap SpanBoundry[E]
	Closed  bool
	Pos     int
	Current *[]*CurrentColumn[E]
}

func (s *ColumnSets[E]) Close() {
	if s.Closed {
		return
	}
	s.Closed = true
	if s.Columns != nil {
		for _, col := range *s.Columns {
			col.Close()
		}
	}
}

func (s *ColumnSets[E]) OverlapCount() int {
	if s.Current == nil {
		return -1
	}
	return len(*s.Current)
}

// Appends an a column accumulator to the current column set.
// Returns the id of the column, if the instance is closed returns -1.
func (s *ColumnSets[E]) AddColumn(c *ColumnOverlapAccumulator[E]) int {
	if s.Closed {
		return -1
	}
	if s.Columns == nil {
		s.Columns = &[]*ColumnOverlapAccumulator[E]{}
	}
	*s.Columns = append(*s.Columns, c)
	return len(*s.Columns) - 1
}

func (s *ColumnSets[E]) AddColumnFromSpanSlice(list *[]SpanBoundry[E]) int {
	return s.AddColumn(s.Util.NewSpanOverlapAccumulator().ColumnOverlapSliceFactory(list))
}

func (s *ColumnSets[E]) Init() {
	var check = []int{}
	var test = &[]SpanBoundry[E]{}
	for i, span := range *s.Columns {
		if span.HasNext() {
			check = append(check, i)
			*test = append(*test, span)
		}
	}
	if len(check) == 0 {
		s.Pos = -1
		return
	}
	s.Pos = 0
	s.Overlap = s.Util.FirstSpan(test)
	s.Active = &check
	s.SetCurrent()

}
func (s *ColumnSets[E]) SetCurrent() {
	s.Current = &[]*CurrentColumn[E]{}
	for _, i := range *s.Active {
		var col = (*s.Columns)[i]
		col.SetNext(s.Overlap)
		if col.InOverlap() {
			var res = &CurrentColumn[E]{
				ColumnId:      i,
				ColumnOverlap: col,
			}
			*s.Current = append(*s.Current, res)
		}
	}
}

func (s *ColumnSets[E]) SetNext() {
	var check = []int{}
	var test = &[]SpanBoundry[E]{}
	for i, span := range *s.Columns {
		if span.HasNext() {
			check = append(check, i)
			*test = append(*test, span)
		}
	}
	if len(check) == 0 {
		s.Pos = -1
		return
	}

	s.Overlap = s.Util.NextSpan(s.Overlap, test)
	if s.Overlap == nil {
		s.Pos = -1
		return
	}
	s.Pos++

	s.Active = &check
	s.SetCurrent()
}
