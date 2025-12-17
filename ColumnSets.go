package st


type ColumnSets[E any] struct {
	Util *SpanUtil[E]
	Columns *[]*ColumnOverlapAccumulator[E]
	Active  *[]bool
	Overlap SpanBoundry[E]
	Closed bool
}

func(s *ColumnSets[E]) Close() {
	if(s.Closed) {
		return
	}
	s.Closed=true
	if(s.Columns!=nil) {
		for _,col :=range *s.Columns {
			col.Close()
		}
	}
}