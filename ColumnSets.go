package st


type ColumnSets[E any,T any] struct {
	Util *SpanUtil[E,T]
	Columns *[]*ColumnOverlapAccumulator[E,T]
	Active  *[]bool
	Overlap SpanBoundry[E,T]
	Closed bool
}

func(s *ColumnSets[E, T]) Close() {
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