package st


// This structure represents a "data source" and how it intersects with an external SpanBoundry.
// Contains the current iterator control functions and represents the column position in the iterator process.
type ColumnOverlapAccumulator[E any] struct {
	// Representation of the data that intersected with an SpanBoundry passed to GetNext.
	// A value of nil means no data overlaps.
	Overlaps *[]*OverlappingSpanSets[E]

	// Where Overlaps begins relative to our OverlappingSpanSets[E,T] iteration.
	// A value of -1 means there was no overlap with the last SpanBoundry[E,T] compared.
	SrcStart int
	// Where Overlaps ends relative to our OverlappingSpanSets[E,T] iteration.
	// A value of -1 means there was no overlap with the last SpanBoundry[E,T] compared.
	SrcEnd int

	// Span utility instance
	Util *SpanUtil[E]

	// The iter.Pull2 "next" method generated from the iter.Seq2 instance.
	ItrGetNext func() (int, *OverlappingSpanSets[E], bool)
	// The iter.Pull2 "stop" method generated from the iter.Seq2 instance.
	ItrStop func()

	// The next set to operate on, when nil.
	Next *OverlappingSpanSets[E]

	// Denotes if the object is closed
	Closed bool
}

func (s *ColumnOverlapAccumulator[E]) GetBegin() E {
	return s.Next.GetBegin()
}

func (s *ColumnOverlapAccumulator[E]) GetEnd() E {
	return s.Next.GetEnd()
}

// Returns the first span that intersecs with this set.
func (s *ColumnOverlapAccumulator[E]) GetFirstSpan() (int,SpanBoundry[E]) {
	return (*s.Overlaps)[0].GetFirstSpan()
}

// Returns the last span that intersecs with this set.
func (s *ColumnOverlapAccumulator[E]) GetLastSpan() (int,SpanBoundry[E]) {
	return (*s.Overlaps)[len(*s.Overlaps)-1].GetLastSpan()
}

// Returns all spans with the sequence id from the orginal data source.
func (s *ColumnOverlapAccumulator[E]) GetSources() (*[]*OvelapSources[E]) {
	list :=[]*OvelapSources[E]{}
	for _,ol :=range *s.Overlaps {
		src :=ol.GetSources()
		list=append(list,(*src)...)
	}
	return &list
}

// Returns the first positional id from the orignal data set for this column.
func (s *ColumnOverlapAccumulator[E]) GetSrcId() int {
	return s.SrcStart
}

// Returns the last positional id from the orignal data set for this column.
func (s *ColumnOverlapAccumulator[E]) GetEndId() int {
	return s.SrcEnd
}

// Returns the overlap sets.
func (s *ColumnOverlapAccumulator[E]) GetOverlaps() *[]*OverlappingSpanSets[E] {
	return s.Overlaps
}

// Returns true if there are more elements in this column.
func (s *ColumnOverlapAccumulator[E]) HasNext() bool {
	return !s.Closed && s.Next != nil
}

// This method is used to call the stop method of the iter.Pull2 iterator method.
// If you are managing an instance of ColumnOverlapAccumulator[E,T] on your own, make sure
// to setup a defer SpanOverlapColumnAccumulator[E,T].Close() to ensure your code does not leak memory
// or run into undefined behaviors.
func (s *ColumnOverlapAccumulator[E]) Close() {
	if s.Closed {
		return
	}
	if s.ItrStop != nil {
		s.ItrStop()
		s.Closed = true
	}
}

// When true this instance contains elements in "Overlaps" that intersect with 
// the last value passed to SetNext.
func (s *ColumnOverlapAccumulator[E]) InOverlap() bool {
	return s.HasNext() && s.SrcStart != -1
}

// This method updates the state of the instance instance in relation to overlap.
// The overlap is considered an external point for comparison, and the internal data sets are
// updated to reflect the current intersection points if any.
func (s *ColumnOverlapAccumulator[E]) SetNext(overlap SpanBoundry[E]) {
	var current = s.Next
	var hasnext = current != nil
	var u = *s.Util
	if hasnext && s.Overlaps != nil && len(*s.Overlaps) != 1 && u.Overlap(overlap, current) {
		var ol = &[]*OverlappingSpanSets[E]{}
		for _, span := range *s.Overlaps {
			if span == current {
				break
			}
			if u.Overlap(overlap, span) {

				*ol = append(*ol, span)
			}
		}
		if len(*ol) != 0 {
			s.SrcStart = (*ol)[0].SrcBegin
			s.Overlaps = ol
		} else {
			s.SrcStart = -1
			s.SrcEnd = -1
			s.Overlaps = &[]*OverlappingSpanSets[E]{}
		}
	} else {
		s.SrcStart = -1
		s.SrcEnd = -1
		s.Overlaps = &[]*OverlappingSpanSets[E]{}
	}

	for hasnext {
		s.Next = current

		if u.Overlap(overlap, current) {
			if s.SrcStart == -1 {
				s.SrcStart = current.SrcBegin
			}
			s.SrcEnd = current.SrcEnd
			*s.Overlaps = append(*s.Overlaps, current)
			if u.Cmp(current.GetEnd(), overlap.GetEnd()) > 0 {
				return
			}
		} else if u.Cmp(current.GetBegin(), overlap.GetEnd()) > 0 {
			// current is after next, then we are done!
			return
		}
		_, current, hasnext = s.ItrGetNext()

		if !hasnext {
			if s.SrcStart == -1 {
				s.Next = nil
			}
			return
		}
	}
}
