package st

// A representation of accumulated Span from a given source.
// The *Span[E,T] represents the span that contains all SpanBoundry[E,T] in *Contains.
// When *Contains is nil, then this struct only has 1 SpanBoundry[E,T].
// When *Contains is not nil, the orignal *SpanBoundry[E,T] values are contained within.
type OverlappingSpanSets[E any, T any] struct {

	// The Span that contains all Spans in this instance.
	Span SpanBoundry[E, T]

	// When nil, Span is the only value representing this Span.
	// When not nill, contains all the Spans acumulated to create this instance.
	Contains *[]SpanBoundry[E, T]

	// Starting position in the original data set
	SrcBegin int

	// Ending position in the original data set
	SrcEnd int
}


func (s *OverlappingSpanSets[E, T]) IsUnique() bool {
	return s.Contains == nil
}

func (s *OverlappingSpanSets[E, T]) GetContains() *[]SpanBoundry[E, T] {
	return s.Contains
}

func (s *OverlappingSpanSets[E, T]) GetTag() *T {
	return s.Span.GetTag()
}

func (s *OverlappingSpanSets[E, T]) GetBegin() E {
	return s.Span.GetBegin()
}

func (s *OverlappingSpanSets[E, T]) GetEnd() E {
	return s.Span.GetEnd()
}



