package st


type ColumnSets[E any,T any] struct {
	Uti *SpanUtil[E,T]
	Current *[]*ColumnOverlapAccumulator[E,T]
	Active  *[]bool
}

