package st

// Represents the stater for going to the: next span.
type SpanIterSeq2Stater[E any] struct {
	Current *OverlappingSpanSets[E]
	Next    *OverlappingSpanSets[E]
	Sa      *SpanOverlapAccumulator[E]
	Id      int
}

// Returns true if this SpanBoundry[E] created a new data intersrection.
// If the value is true you must make a call to s.GetNext() instance method before calling this method again!
func (s *SpanIterSeq2Stater[E]) SetNext(span SpanBoundry[E]) bool {
	cmp,_ := s.Sa.Accumulate(span)
	if s.Current == nil {
		s.Current = cmp
		return false
	}
	if s.Current == cmp {
		return false
	}
	s.Next = cmp
	return true
}

// Returns true if we have any more intersections.
func (s *SpanIterSeq2Stater[E]) HasNext() bool {
	return s.Current != nil
}

// Returns the current index and intersection if any.
func (s *SpanIterSeq2Stater[E]) GetNext() (int, *OverlappingSpanSets[E]) {
	var next = s.Current
	s.Current = s.Next
	s.Next = nil
	s.Id++

	return s.Id, next
}

