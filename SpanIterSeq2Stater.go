package st

type SpanIterSeq2Stater[E any, T any] struct {
	Current *OverlappingSpanSets[E, T]
	Next    *OverlappingSpanSets[E, T]
	Sa      *SpanOverlapAccumulator[E, T]
	Id      int
}

func (s *SpanIterSeq2Stater[E, T]) SetNext(span SpanBoundry[E, T]) bool {
	var cmp = s.Sa.Accumulate(span)
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

func (s *SpanIterSeq2Stater[E, T]) HasNext() bool {
	return s.Current != nil
}

func (s *SpanIterSeq2Stater[E, T]) GetNext() (int, *OverlappingSpanSets[E, T]) {
	var next = s.Current
	s.Current = s.Next
	s.Next = nil
	s.Id++

	return s.Id, next
}

