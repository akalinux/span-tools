package st

import (
	"iter"
	"slices"
	"context"
)

// This is a stater structure, used to drive the creation of new OverlappingSpanSets.
// Each source of spans should have their own instance of SpanOverlapAccumulator.
type SpanOverlapAccumulator[E any] struct {
	Rss *OverlappingSpanSets[E]
	*SpanUtil[E]
	// When true slices passed in will be sorted.
	Sort bool

	// When not nil, this object has encounter an error
	Err error

	// Sequence counter
	Pos int

	// Turns validation on/off, default false or off
	Validate bool

	// Turns consolidation of adjacent spans on or off, default false or off
	Consolidate bool
}

// The Accumulate method.
//
// For a given span provided: When the span overlaps with the current internal span,
// the OverlappingSpanSets is expanded and the span is append to the Contains slice.
// When the span is outside of the current internal span,
// then a new OverlappingSpanSets is created with this span as its current span.
// The error value is nil, by default, when an error has happend it is no longer nil.
func (s *SpanOverlapAccumulator[E]) Accumulate(span SpanBoundry[E]) (*OverlappingSpanSets[E], error) {
	s.Pos++
	if s.Validate && s.Err==nil {
		s.Err = s.Check(span, s.Rss.Span)
	}

	if s.Err != nil {
		s.Rss.Err=s.Err
	}

	if s.Rss.Span == nil {
		s.Rss.Span = span
		return s.Rss, s.Err
	}

	a := s.Rss.Span
	if s.Cmp(a.GetEnd(), span.GetBegin()) < 0 {
		var joined = false
		if s.Consolidate {
			var next = s.Next(a.GetEnd())
			if s.Cmp(next, span.GetBegin()) == 0 {
				s.Rss.Span = s.Ns(a.GetBegin(),span.GetEnd())
				joined = true
			}
		}
		if !joined {

			s.Rss = &OverlappingSpanSets[E]{
				Span:     span,
				Contains: nil,
				SrcBegin: s.Pos,
				SrcEnd:   s.Pos,
				Err: s.Err,
			}
		}
	} else {
		x, y := s.ContainedBy(a, span)
		if x|y != 0 {
			var begin,end E;
			if x < 0 {
				begin = a.GetBegin()
			} else {
				begin = span.GetBegin()
			}
			if y > 0 {
				end = a.GetEnd()
			} else {
				end = span.GetEnd()
			}
			s.Rss.Span = s.Ns(begin,end)
		}

		if s.Rss.Contains == nil {
			s.Rss.Contains = &[]SpanBoundry[E]{a, span}
		} else {
			*s.Rss.Contains = append(*s.Rss.Contains, span)
		}
		s.Rss.SrcEnd = s.Pos
	}
	return s.Rss, s.Err
}



// Helper function to create an overlap iterator from a slice of list.
func (s *SpanOverlapAccumulator[E]) NewOlssSeq2FromOlssSlice(list *[]*OverlappingSpanSets[E]) iter.Seq2[int, *OverlappingSpanSets[E]] {
  return slices.All(*list)
}

// Generates a iter.Seq2 iterator, for a channel of SpanBoundry instances.
func (s *SpanOverlapAccumulator[E]) NewOlssSeq2FromSbChan(c <-chan SpanBoundry[E]) iter.Seq2[int, *OverlappingSpanSets[E]] {
	var sa = s.NewSpanIterSeq2Stater()
	if c != nil {
		var span, ok = <-c
		for ok {
			if sa.SetNext(span) {
				break
			}
			span, ok = <-c
		}
	}
	return func(yeild func(int, *OverlappingSpanSets[E]) bool) {
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

func (s *SpanOverlapAccumulator[E]) NewSpanIterSeq2Stater() *SpanIterSeq2Stater[E] {
	var si = &SpanIterSeq2Stater[E]{
		Sa:      s,
		Current: nil,
		Next:    nil,
		Id:      -1,
	}
	return si
}

// Factory interface for converting slices of SpanBoundaries instances into iterator sequences of OverlappingSpanSets.
func (s *SpanOverlapAccumulator[E]) NewOlssSeq2FromSbSlice(list *[]SpanBoundry[E]) iter.Seq2[int, *OverlappingSpanSets[E]] {
	var end = -1
	var pos = 0
	var au = s.NewSpanIterSeq2Stater()
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

	return func(yeild func(int, *OverlappingSpanSets[E]) bool) {
		// no list stop here
		if end == -1 {
			return
		}

		for {
			if au.HasNext() {
				var id, current = au.GetNext()

				if !yeild(id, current) {
					return
				}
				for pos < end {
					if au.SetNext((*list)[pos]) {
						pos++
						break
					}
					pos++
				}
			} else {
				return
			}

		}
	}
}

// This is a convenience method for initializing the iter.Seq2 stater internals based on a slice of SpanBoundry.
func (s *SpanOverlapAccumulator[E]) NewCoaFromSbSlice(list *[]SpanBoundry[E]) *ColumnOverlapAccumulator[E] {
	return s.NewCoaFromOlssSeq2(s.NewOlssSeq2FromSbSlice(list))
}

func (s *SpanOverlapAccumulator[E]) NewCoaFromOlssChan(c <-chan *OverlappingSpanSets[E]) *ColumnOverlapAccumulator[E] {
	return s.SpanUtil.NewCoaFromOlssSeq2(s.NewOlssSeq2FromOlssChan(c))
}


func (s *SpanOverlapAccumulator[E]) NewOlssChanStater() *OlssChanStater[E] {
	ctx,cancle :=context.WithCancel(context.Background())
	return &OlssChanStater[E]{
		Stater: *s.NewSpanIterSeq2Stater(),
		Chan: make(chan *OverlappingSpanSets[E]),
		Ctx: ctx,
		Cancel: cancle,
	}
}
