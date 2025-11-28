// Package span implements utilties for handling spans of values.
//
// The package provides a data span/range intersection library that is algorithmically implmented using generics and should work with any
// data set that has a comparable Begin and End value.
//
// # How this package treats spans.
//
// Spans in this package are expected to contain a Begin and End value. The Begin and End values should be 
// comparable with a cmp function.  The Begin value is expected to be less than or equal to the End value.
package spans

import (
	"cmp"
	"errors"
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
  Span SpanBounds[E, T]

  // When nil, Span is the only value representing this Span.
  // When not nill, contains all the Spans acumulated to create this instance.
  Contains *[]SpanBounds[E, T]
}

type SpanOverlapBounds[E any, T any] interface {
  SpanBounds[E,T]
  // Returns a pointer to overlapping spans, nill if thiis is a unique span.
  GetContains() *[]SpanBounds[E,T]
  // If this is a uniqe span returns true, otherwise returns false.
  IsUnique() bool
}

func (s *OverlappingSpanSets[E,T]) IsUnique() bool {
  return s.Contains==nil
}

func (s *OverlappingSpanSets[E,T]) GetContains() *[]SpanBounds[E,T] {
  return s.Contains
}

func (s *OverlappingSpanSets[E,T]) GetTag() *T {
  return s.Span.GetTag();
}

func (s *OverlappingSpanSets[E,T]) GetBegin() E {
  return s.Span.GetBegin()
}

func (s *OverlappingSpanSets[E,T]) GetBeginP() *E {
  return s.Span.GetBeginP();
}
func (s *OverlappingSpanSets[E,T]) GetEndP() *E {
  return s.Span.GetEndP();
}

func (s *OverlappingSpanSets[E,T]) GetEnd() E {
  return s.Span.GetEnd()
}

// Core of the span utilties: Provides methos for processing ranges.
type SpanUtil[E any, T any] struct {
  Cmp func(a, b E) int
}

// Creates a instance of *SpanUtil[E cmp.Ordered,T], this can be used to process most span data sets.
func NewOrderedSpanUtil[E cmp.Ordered, T any]() *SpanUtil[E, T] {
  return NewSpanUtil[E, T](cmp.Compare)
}

type SpanBounds[E any,T any] interface {
  // Returns the Begin value.
  GetBegin() E
  // Returns the pointer to the Begin value.
  GetBeginP() *E
  // Returns the End value.
  GetEnd() E
  // Returns the pointer to the End value.
  GetEndP() *E
  // Returns the pointer to the Tag value.
  GetTag() *T
}

func (s *Span[E,T]) GetTag() *T {
  return s.Tag;
}

func (s *Span[E,T]) GetBegin() E {
  return s.Begin
}

func (s *Span[E,T]) GetBeginP() *E {
  return &s.Begin;
}
func (s *Span[E,T]) GetEndP() *E {
  return &s.End;
}

func (s *Span[E,T]) GetEnd() E {
  return s.End
}

// Creates an instance of *SpanUtil[E,T], the value of cmp is expected to be able to compare the Span.Begin and Span.End values.
// See [[cmp.Compare]] for more info.
//
// [cmp.Compare]: https://pkg.go.dev/github.com/google/go-cmp/cmp#Comparer
func NewSpanUtil[E any, T any](cmp func(a, b E) int) *SpanUtil[E, T] {
  return &SpanUtil[E, T]{Cmp: cmp}
}

// This method is used to sort slice of spans in the accumulation order. 
// For more details see: [[slices.SortFunc]]
// 
// [slices.SortFunc]: https://pkg.go.dev/slices#SortedFunc
func (s *SpanUtil[E, T]) Compare(a, b SpanBounds[E, T]) int {
  var diff int = s.Cmp(a.GetBegin(), b.GetBegin())
  if diff == 0 {
    return s.Cmp(b.GetEnd(), a.GetEnd())
  }
  return diff
}

// Returns true if a contains b.
func (s *SpanUtil[E, T]) Contains(a SpanBounds[E, T], b E) bool {
  return s.Cmp(a.GetBegin(), b) < 1 && s.Cmp(a.GetEnd(), b) > -1
}

// Returns true if a overlaps with b or if be overlaps with a.
func (s *SpanUtil[E, T]) Overlap(a, b SpanBounds[E, T]) bool {
  return s.Contains(a, b.GetBegin()) || s.Contains(a, b.GetEnd()) || s.Contains(b, a.GetEnd()) || s.Contains(b, a.GetEnd())
}

// This method is used to determin the outer bounds of ranges a and b.  
// The first int represents comparing a.Begin to b.Begin and the second int represents comparing a.End to b.End.
func (s *SpanUtil[E, T]) ContainedBy(a, b SpanBounds[E, T]) (int, int) {
  return s.Cmp(a.GetBegin(), b.GetBegin()), s.Cmp(a.GetEnd(), b.GetEnd())
}

// Creates a new span, error is nill unless a is greater than b.
func (s *SpanUtil[E, T]) NewSpan(a,b E, tag *T) (*Span[E,T],error) {
  if(s.Cmp(a,b)>0) {
    return nil,errors.New("Value a is greater than value b")
  }
  return &Span[E,T]{Begin: a,End: b,Tag: tag},nil
}

// This method returns the first smallest span from the slice of Span[E,T].
func (s *SpanUtil[E, T]) FirstSpan(list *[]SpanBounds[E, T]) *Span[E, T] {
  var span = &Span[E, T]{Begin: (*list)[0].GetBegin(), End: (*list)[0].GetEnd()}
  var last = len(*list)
  for i := 1; i < last; i++ {
    var check = (*list)[i]
    if s.Cmp(check.GetBegin(), span.GetEnd()) == -1 {
      span.Begin = check.GetBegin()
    }
    if s.Cmp(check.GetEnd(), span.End) == -1 {
      span.End = check.GetEnd()
    }
  }
  return span
}

// This method acts as a stateless iterator that, 
// returns the next overlapping Span[E,T] or nill based on the start Span[E,T] and the slice of spans. 
// If all valid Span[E,T] values have been exausted, nil is returned.
func (s *SpanUtil[E, T]) NextSpan(start SpanBounds[E, T], list *[]SpanBounds[E, T]) *Span[E, T] {
  var begin *E = nil
  var end *E = nil
  for _, check := range *list {
    if begin == nil {
      if s.Cmp(check.GetBegin(), start.GetEnd()) > 0 {
        begin = check.GetBeginP()
        end = check.GetEndP()
      } else if s.Cmp(check.GetEnd(), start.GetEnd()) > 0 {
        begin = check.GetEndP()
        end = check.GetEndP()
      }
    } else {
      if s.Cmp(check.GetBegin(), start.GetEnd()) > 0 && s.Cmp(*begin, check.GetBegin()) > 0 {
        begin = check.GetBeginP()
      }
      if s.Cmp(*begin, check.GetEnd()) < 1 && s.Cmp(check.GetEnd(), start.GetEnd()) > 0 && s.Cmp(*end, check.GetEnd()) > 0 {
        end = check.GetEndP()
      }
    }
  }
  if begin != nil {
    return &Span[E, T]{Begin: *begin, End: *end}
  }
  return nil
}

// This is a stater structure, used to driver the creation of new OverlappingSpanSets.
type SpanOverlapAccumulator[E any, T any] struct {
  Rss *OverlappingSpanSets[E, T]
  *SpanUtil[E, T]
}

// Factory interface for the creation of SpanOverlapAccumulator[E,T].
func (s *SpanUtil[E, T]) NewSpanOverlapAccumulator() *SpanOverlapAccumulator[E, T] {
  return &SpanOverlapAccumulator[E, T]{SpanUtil: s, Rss: &OverlappingSpanSets[E, T]{Contains: nil, Span: nil}}
}

// The Accumulate method.  
// 
// For a given Span[E,T] provided:
// When the span overlaps with the current Span[E,T], the OverlappingSpanSets is expanded and the span is appened to the Contains slice.
// When the span is outside of the current Span[E,T], then a new OverlappingSpanSets is created with this span as its current span.
func (s *SpanOverlapAccumulator[E, T]) Accumulate(span SpanBounds[E, T]) *OverlappingSpanSets[E, T] {
  if s.Rss.Span == nil {
    s.Rss.Span = span
    return s.Rss
  }

  a := s.Rss.Span
  if s.Cmp(a.GetEnd(), span.GetBegin()) < 0 {
    s.Rss = &OverlappingSpanSets[E, T]{Span: span, Contains: nil}
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
      s.Rss.Contains = &[]SpanBounds[E, T]{a, span}
    } else {
      *s.Rss.Contains = append(*s.Rss.Contains, span)
    }
  }
  return s.Rss
}
