package st

// Contains the current iterator control functions and represents the column position in the iterator process.
type ColumnOverlapAccumulator[E any, T any] struct {
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
func (s *ColumnOverlapAccumulator[E, T]) HasNext() bool {
	return s.Next != nil
}

// This method is used to call the stop method of the iter.Pull2 iterator method.
// If you are managing an instance of ColumnOverlapAccumulator[E,T] on your own, make sure
// to setup a defer SpanOverlapColumnAccumulator[E,T].Close() to ensure your code does not leak memory
// or run into undefined behaviors.
func (s *ColumnOverlapAccumulator[E, T]) Close() {
	if s.ItrStop != nil {
		s.ItrStop()
	}
}

// This method updates the currrent instance  to represent what data intersects with the overlap SpanBoundry[E,T].
func (s *ColumnOverlapAccumulator[E, T]) SetNext(overlap SpanBoundry[E, T]) {
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


