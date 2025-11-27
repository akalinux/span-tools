package spans

import (
  "cmp"
)

// Representation of a Span/Range of values in a generic context.
// The assumption is that Begin is less than or equal to the End value.
type Span[E any, T any] struct {
  // Start of the Span.
  Begin E
  // End of the Span.
  End E
  // Pointer to data set used to identify this Span[E,T]
  Tag *T
}

// A representation of accumulated Span from a given source.
// The *Span[E,T] represents the span that contains all Spans[E,T] in *Contains.
// When *Contains is nil, then this struct only has 1 Span[E,T].
// When *Contains is not nil, the orignal *Span[E,T] values are contained within.
type OverlappingSpanSets[E any, T any] struct {

  // The Span that contains all Spans in this instance.
  *Span[E, T]

  // When nil, Span is the only value representing this Span.
  // When not nill, contains all the Spans acumulated to create this instance.
  Contains *[]*Span[E, T]
}

// Core of the span utilties: Provides methos for processing ranges.
type SpanUtil[E any, T any] struct {
  Cmp func(a, b E) int
}

// Creates a instance of *SpanUtil[E cmp.Ordered,T], this can be used to process most span data sets.
func NewOrderedSpanUtil[E cmp.Ordered, T any]() *SpanUtil[E, T] {
  return NewSpanUtil[E, T](cmp.Compare)
}

func NewSpanUtil[E any, T any](cmp func(a, b E) int) *SpanUtil[E, T] {
  return &SpanUtil[E, T]{Cmp: cmp}
}

// This method is used to sort slice of spans in the accumulation order. 
// For more details see: [[slices.SortFunc]]
// 
// [slices.SortFunc]: https://pkg.go.dev/slices#SortedFunc
func (s *SpanUtil[E, T]) Compare(a, b Span[E, T]) int {
  var diff int = s.Cmp(a.Begin, b.Begin)
  if diff == 0 {
    return s.Cmp(b.End, a.End)
  }
  return diff
}

// Returns true if a contains b.
func (s *SpanUtil[E, T]) Contains(a *Span[E, T], b E) bool {
  return s.Cmp(a.Begin, b) < 1 && s.Cmp(a.End, b) > -1
}

// Returns true if a overlaps with b or if be overlaps with a.
func (s *SpanUtil[E, T]) Overlap(a, b *Span[E, T]) bool {
  return s.Contains(a, b.Begin) || s.Contains(a, b.End) || s.Contains(b, a.Begin) || s.Contains(b, a.End)
}

func (s *SpanUtil[E, T]) ContainedBy(a, b *Span[E, T]) (int, int) {
  return s.Cmp(a.Begin, b.Begin), s.Cmp(a.End, b.End)
}

func (s *SpanUtil[E, T]) FirstSpan(list *[]*Span[E, T]) *Span[E, T] {
  var span = &Span[E, T]{Begin: (*list)[0].Begin, End: (*list)[0].End}
  var last = len(*list)
  for i := 1; i < last; i++ {
    var check = (*list)[i]
    if s.Cmp(check.Begin, span.Begin) == -1 {
      span.Begin = check.Begin
    }
    if s.Cmp(check.End, span.End) == -1 {
      span.End = check.End
    }
  }
  return span
}

// This method acts as a stateless iterator that, 
// returns the next overlapping Span[E,T] or nill based on the start Span[E,T] and the slice of spans. 
// If all valid Span[E,T] values have been exausted, nil is returned.
func (s *SpanUtil[E, T]) NextSpan(start *Span[E, T], list *[]*Span[E, T]) *Span[E, T] {
  var begin *E = nil
  var end *E = nil
  for _, check := range *list {
    if begin == nil {
      if s.Cmp(check.Begin, start.End) > 0 {
        begin = &check.Begin
        end = &check.End
      } else if s.Cmp(check.End, start.End) > 0 {
        begin = &check.End
        end = &check.End
      }
    } else {
      if s.Cmp(check.Begin, start.End) > 0 && s.Cmp(*begin, check.Begin) > 0 {
        begin = &check.Begin
      }
      if s.Cmp(*begin, check.End) < 1 && s.Cmp(check.End, start.End) > 0 && s.Cmp(*end, check.End) > 0 {
        end = &check.End
      }
    }
  }
  if begin != nil {
    return &Span[E, T]{Begin: *begin, End: *end}
  }
  return nil
}

type SpanOverlapAccumulator[E any, T any] struct {
  Rss *OverlappingSpanSets[E, T]
  *SpanUtil[E, T]
}

func (s *SpanUtil[E, T]) NewSpanOverlapAccumulator() *SpanOverlapAccumulator[E, T] {
  return &SpanOverlapAccumulator[E, T]{SpanUtil: s, Rss: &OverlappingSpanSets[E, T]{Contains: nil, Span: nil}}
}

func (s *SpanOverlapAccumulator[E, T]) Accumulate(b *Span[E, T]) *OverlappingSpanSets[E, T] {
  if s.Rss.Span == nil {
    s.Rss.Span = b
    return s.Rss
  }

  a := s.Rss.Span
  if s.Cmp(a.End, b.Begin) < 0 {
    s.Rss = &OverlappingSpanSets[E, T]{Span: b, Contains: nil}
  } else {
    x, y := s.ContainedBy(a, b)
    if x|y != 0 {
      var r = Span[E, T]{}
      if x < 0 {
        r.Begin = a.Begin
      }
      if y > 0 {
        r.End = a.End
      }
      s.Rss.Span = &r
    }
    if s.Rss.Contains == nil {
      s.Rss.Contains = &[]*Span[E, T]{a, b}
    } else {
      *s.Rss.Contains = append(*s.Rss.Contains, b)
    }
  }
  return s.Rss
}
