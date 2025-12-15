package st

import (
	"iter"
	"slices"
)

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

// Generates a iter.Seq2 iterator, for a chanel.
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

// This method takes a iter.Seq2 iterator of OverlappingSpanSets and initalizes the SpanOverlapAccumulator struct.
//
// # Warning
//
// This methos creates an [iter.Pull2] and exposes the resulting functions in the returned reference. If you are using this methos outside of the normal
// operations, you should a setup a defer call to  SpanOverlapAccumulator[E, T].Close() method to clean the instance up in order to prevent memory leaks or undefined behavior.
//
// [iter.Pull2]: https://pkg.go.dev/iter#hdr-Pulling_Values
func (s *SpanOverlapAccumulator[E, T]) ColumnOverlapFactory(driver iter.Seq2[int, *OverlappingSpanSets[E, T]]) *SpanOverlapColumnAccumulator[E, T] {
	var next, stop = iter.Pull2(driver)
	return s.Init(next, stop)
}

// This method takes the next and stop functions and creates a new fully initalized instance of SpanOverlapColumnAccumulator[E, T].
func (s *SpanOverlapAccumulator[E, T]) Init(next func() (int, *OverlappingSpanSets[E, T], bool), stop func()) *SpanOverlapColumnAccumulator[E, T] {
	var res = &SpanOverlapColumnAccumulator[E, T]{}
	res.ItrStop = stop
	res.ItrGetNext = next
	res.Util = s.SpanUtil
	var id, current, _ = res.ItrGetNext()
	res.SrcPos = id
	res.Next = current
	return res
}

// This is a convenience method for initalizing the iter.Seq2 stater internals based on a slice of SpanBoundry.
func (s *SpanOverlapAccumulator[E, T]) ColumnOverlapSliceFactory(list *[]SpanBoundry[E, T]) *SpanOverlapColumnAccumulator[E, T] {
	return s.ColumnOverlapFactory(s.SliceIterFactory(list))
}

// Contains the current iterator control functions and represents the column position in the iterator process.
type SpanOverlapColumnAccumulator[E any, T any] struct {
	// Representation of the data that intersected with an SpanBoundry passed to GetNext.
	// A value of nil means no data overlaps.
	Overlaps *[]*OverlappingSpanSets[E, T]

	// Where Overlaps begins relative to our OverlappingSpanSets[E,T] iteration.
	// A value of -1 means there was no overlap with the last SpanBoundry[E,T] compared.
	SrcStart int
	// Where Overlaps ends relative to our OverlappingSpanSets[E,T] iteration.
	// A value of -1 means there was no overlap with the last SpanBoundry[E,T] compared.
	SrcEnd int

	// Span utility instance
	Util *SpanUtil[E, T]

	// The iter.Pull2 "next" method generated from the iter.Seq2 instance.
	ItrGetNext func() (int, *OverlappingSpanSets[E, T], bool)
	// The iter.Pull2 "stop" method generated from the iter.Seq2 instance.
	ItrStop func()

	// The next set to operate on, when nil.
	Next *OverlappingSpanSets[E, T]

	// Denotes where we are in the orginal OverlappingSpanSets[E,T] instance.
	SrcPos int
}

// Returns true if there are more elements in this column.
func (s *SpanOverlapColumnAccumulator[E, T]) HasNext() bool {
	return s.Next != nil
}

// This method is used to call the stop method of the iter.Pull2 iterator method.
// If you are managing an instance of SpanOverlapColumnAccumulator[E,T] on your own, make sure
// to setup a defer SpanOverlapColumnAccumulator[E,T].Close() to ensure your code does not leak memory
// or run into undefined behaviors.
func (s *SpanOverlapColumnAccumulator[E, T]) Close() {
	if s.ItrStop != nil {
		s.ItrStop()
	}
}

// This method updates the currrent SpanOverlapColumnAccumulator[E, T] to represent what data intersects with the overlap SpanBoundry[E,T].
func (s *SpanOverlapColumnAccumulator[E, T]) GetNext(overlap SpanBoundry[E, T]) {
	s.Overlaps = &[]*OverlappingSpanSets[E, T]{}
	var id = s.SrcPos
	var current = s.Next
	var hasnext = current != nil
	var u = *s.Util
	s.SrcPos = -1
	s.SrcStart = -1

	for hasnext {
		s.SrcPos = id
		s.Next = current
		if u.Overlap(overlap, current) {
			if s.SrcStart == -1 {
				s.SrcStart = id
			}
			s.SrcEnd = id
			*s.Overlaps = append(*s.Overlaps, current)
			if u.Cmp(current.GetEnd(), overlap.GetEnd()) > 0 {
				return
			}
		} else if u.Cmp(current.GetBegin(), overlap.GetEnd()) > 0 {
			// current is after next, then we are done!
			return
		}
		id, current, hasnext = s.ItrGetNext()

		if !hasnext {
			if s.SrcStart == -1 {
				s.Next = nil
				s.SrcPos = -1
			}
			return
		}
	}
}
