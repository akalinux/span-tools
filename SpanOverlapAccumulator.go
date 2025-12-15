package st

import (
	"iter"
	"slices"
)
// This is a stater structure, used to driverthe creation of new OverlappingSpanSets.
type SpanOverlapAccumulator[E any, T any] struct {
	Rss *OverlappingSpanSets[E, T]
	*SpanUtil[E, T]
	// When true slices passed in will be sorted.
	Sort     bool
	Err      error
	Pos      int
	Validate bool
}



// The Accumulate method.
//
// For a given Span[E,T] provided:
// When the span overlaps with the current Span[E,T], the OverlappingSpanSets is expanded and the span is appened to the Contains slice.
// When the span is outside of the current Span[E,T], then a new OverlappingSpanSets is created with this span as its current span.
func (s *SpanOverlapAccumulator[E, T]) Accumulate(span SpanBoundry[E, T]) *OverlappingSpanSets[E, T] {
	s.Pos++
	if s.Validate {
		s.Err = s.Check(span, s.Rss.Span)
	}

	if s.Rss.Span == nil {
		s.Rss.Span = span
		return s.Rss
	}

	a := s.Rss.Span
	if s.Cmp(a.GetEnd(), span.GetBegin()) < 0 {
		s.Rss = &OverlappingSpanSets[E, T]{
			Span:     span,
			Contains: nil,
			SrcBegin: s.Pos,
			SrcEnd:   s.Pos,
		}
	} else {
		x, y := s.ContainedBy(a, span)
		if x|y != 0 {
			var r = Span[E, T]{}
			if x < 0 {
				r.Begin = a.GetBegin()
			}
			if y > 0 {
				r.End = a.GetEnd()
			}
			s.Rss.Span = &r
		}

		if s.Rss.Contains == nil {
			s.Rss.Contains = &[]SpanBoundry[E, T]{a, span}
		} else {
			*s.Rss.Contains = append(*s.Rss.Contains, span)
		}
		s.Rss.SrcEnd = s.Pos
	}
	return s.Rss
}

// Creates a channel iteraotr for channels of OverlappingSpanSets.
func (s *SpanOverlapAccumulator[E, T]) ChanIterFactoryOverlaps(c <-chan *OverlappingSpanSets[E, T]) iter.Seq2[int, *OverlappingSpanSets[E, T]] {

	if c == nil {
		return func(yeild func(int, *OverlappingSpanSets[E, T]) bool) {
		}
	}
	var i = 0
	return func(yeild func(int, *OverlappingSpanSets[E, T]) bool) {
		var ol, ok = <-c
		for ok {

			if !yeild(i, ol) {
				return
			}
			i++
		  ol, ok = <-c
		}

	}
}

// Generates a iter.Seq2 iterator, for a channel of SpanBoundry instances.
func (s *SpanOverlapAccumulator[E, T]) ChanIterFactory(c <-chan SpanBoundry[E, T]) iter.Seq2[int, *OverlappingSpanSets[E, T]] {
	var sa = s.SpanStatefulAccumulator()
	if c != nil {
		var span, ok = <-c
		for ok {
			if sa.SetNext(span) {
				break
			}
			span, ok = <-c
		}
	}
	return func(yeild func(int, *OverlappingSpanSets[E, T]) bool) {
		// no chan??? stop here
		if !sa.HasNext() {
			return
		}

		for {
			if s.Err != nil {
				return
			}
			if sa.HasNext() {
				var id, current = sa.GetNext()
				if !yeild(id, current) {
					return
				}
				var span, ok = <-c
				for ok {
					if sa.SetNext(span) {
						break
					}
					span, ok = <-c
				}
			} else {
				return
			}
		}
	}
}

func (s *SpanOverlapAccumulator[E, T]) SpanStatefulAccumulator() *SpanIterSeq2Stater[E, T] {
	var si = &SpanIterSeq2Stater[E, T]{
		Sa:      s,
		Current: nil,
		Next:    nil,
		Id:      -1,
	}
	return si
}

// Factory interface for converting slices of SpanBoundaries instances into iterator sequences of OverlappingSpanSets.
func (s *SpanOverlapAccumulator[E, T]) SliceIterFactory(list *[]SpanBoundry[E, T]) iter.Seq2[int, *OverlappingSpanSets[E, T]] {
	var end = -1
	var pos = 0
	var au = s.SpanStatefulAccumulator()
	if list != nil {
		if s.Sort {
			slices.SortFunc(*list, s.Compare)
		}
		end = len(*list)
		for pos < end {
			if au.SetNext((*list)[pos]) {
				pos++
				break
			}
			pos++
		}
	}

	return func(yeild func(int, *OverlappingSpanSets[E, T]) bool) {
		// no list stop here
		if end == -1 {
			return
		}

		for {
			if s.Err != nil {
				return
			}
			if au.HasNext() {
				var id, current = au.GetNext()
				if !yeild(id, current) {
					return
				}
			} else {
				return
			}

			for pos < end {
				if au.SetNext((*list)[pos]) {
					pos++
					break
				}
				pos++
			}
		}
	}
}




