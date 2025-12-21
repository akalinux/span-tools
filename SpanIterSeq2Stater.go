package st


type SpanIterSeq2Stater[E any] struct {
	Current *OverlappingSpanSets[E]
	Next    *OverlappingSpanSets[E]
	Sa      *SpanOverlapAccumulator[E]
	Id      int
}

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

func (s *SpanIterSeq2Stater[E]) HasNext() bool {
	return s.Current != nil
}

func (s *SpanIterSeq2Stater[E]) GetNext() (int, *OverlappingSpanSets[E]) {
	var next = s.Current
	s.Current = s.Next
	s.Next = nil
	s.Id++

	return s.Id, next
}

