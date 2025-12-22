package st

// A representation of accumulated Span from a given source.
// The *Span[E,T] represents the span that contains all SpanBoundry[E,T] in *Contains.
// When *Contains is nil, then this struct only has 1 SpanBoundry[E,T].
// When *Contains is not nil, the original *SpanBoundry[E,T] values are contained within.
type OverlappingSpanSets[E any] struct {

	// The Span that contains all Spans in this instance.
	Span SpanBoundry[E]

	// When nil, Span is the only value representing this Span.
	// When not nil, contains all the Spans accumulated to create this instance.
	Contains *[]SpanBoundry[E]

	// Starting position in the original data set
	SrcBegin int

	// Ending position in the original data set
	SrcEnd int
	
	Err error
}

func (s *OverlappingSpanSets[E]) GetSrcId() int {
	return s.SrcBegin
}


func (s *OverlappingSpanSets[E]) GetEndId() int {
	return s.SrcEnd
}

func (s *OverlappingSpanSets[E]) GetOverlaps() *[]*OverlappingSpanSets[E] {
	return &[]*OverlappingSpanSets[E]{s}
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

func (s *OverlappingSpanSets[E]) GetFirstSpan() (int, SpanBoundry[E]) {
	if s.IsUnique() {
		return s.SrcBegin, s.Span
	}
	return s.SrcBegin, (*s.Contains)[0]
}

func (s *OverlappingSpanSets[E]) GetLastSpan() (int, SpanBoundry[E]) {
	if s.IsUnique() {
		return s.SrcEnd, s.Span
	}
	return s.SrcEnd, (*s.Contains)[len(*s.Contains)-1]
}

type OvelapSources[E any] struct {
	SpanBoundry[E]
	SrcId int
}

func (s *OverlappingSpanSets[E]) GetSources() *[]*OvelapSources[E] {
	res := &[]*OvelapSources[E]{}
	if s.IsUnique() {
		*res = append(*res, &OvelapSources[E]{
			SpanBoundry: s.Span,
			SrcId:       s.SrcBegin,
		})
	} else {
		for id, span := range *s.Contains {
			*res = append(*res, &OvelapSources[E]{
				SpanBoundry: span,
				SrcId:       s.SrcBegin + id,
			})
		}
	}
	return res
}
