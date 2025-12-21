package st

// A representation of accumulated Span from a given source.
// The *Span[E,T] represents the span that contains all SpanBoundry[E,T] in *Contains.
// When *Contains is nil, then this struct only has 1 SpanBoundry[E,T].
// When *Contains is not nil, the orignal *SpanBoundry[E,T] values are contained within.
type OverlappingSpanSets[E any] struct {

	// The Span that contains all Spans in this instance.
	Span SpanBoundry[E]

	// When nil, Span is the only value representing this Span.
	// When not nill, contains all the Spans acumulated to create this instance.
	Contains *[]SpanBoundry[E]

	// Starting position in the original data set
	SrcBegin int

	// Ending position in the original data set
	SrcEnd int
}


func (s *OverlappingSpanSets[E]) IsUnique() bool {
	return s.Contains == nil
}

func (s *OverlappingSpanSets[E]) GetContains() *[]SpanBoundry[E] {
	return s.Contains
}

func (s *OverlappingSpanSets[E]) GetBegin() E {
	return s.Span.GetBegin()
}

func (s *OverlappingSpanSets[E]) GetEnd() E {
	return s.Span.GetEnd()
}



