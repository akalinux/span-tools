package spans

import (
	"cmp"
	// "fmt"
)

type Span[E any, T any] struct {
  // Start of the Span.
	Begin E
  // End of the Span.
	End   E
  // Pointer to data set used to identify this Span[E,T]
	Tag   *T
}

type AccumulatedSpanSet[E any, T any] struct {
  
  // The Span that contains all Spans in this instance.
	*Span[E, T]
  
  // When nil, Span is the only value representing this Span.
  // When not nill, contains all the Spans acumulated to create this instance. 
	Contains *[]*Span[E, T]
}

type SpanTests[E any, T any] struct {
  // Sort method used by slices.SortFunc to sort slices in order of accumulation.  
  // Order is represent as the smallest Begin value and largest End value.
	Compare         func(a, b Span[E, T]) int
  
  // Compuates the smallest Begin value and largest End value given the input a,b *Span[E,T].  
  // Return values: First int represents the smallest Span[E,T] Begin, Second value represents the largest Span[E,T] End.
	ContainedBy     func(a, b *Span[E, T]) (int, int)
  
  // Returns true if the a *Span[E,T] contains the value b E.
	Contains        func(a *Span[E, T], b E) bool
  
  // Creates a Closure used to accumualte new Span overlaps.
  // Each time a new pointer is returned, a new Span overlap has been found.
  // Values passed to the accumuator(s *Span[E,T]) function must be in the oder produced by slices.SortFunc(SpanTests[E,T].Compare)
	SpanAccumulator func() func(b *Span[E, T]) *AccumulatedSpanSet[E, T]
  // Comapres a to b: returns -1,0,1 based on the 2 values compared
	Cmp             func(a,b E) int
  // Used to compare 2 spans, returns true if they overlap
  Overlap         func(a,b *Span[E,T])  bool
}

func BuildOverlap[E any ,T any](contains func (a *Span[E,T], b E) bool) func(a,b *Span[E,T]) bool {
  return func (a,b *Span[E,T]) bool {
    return contains(a,b.Begin) || contains(a,b.End) || (contains(b,a.Begin) || contains(b,a.End))
  }
}

func BuildCompare[E any, T any](cmp func(E, E) int) func(a, b Span[E, T]) int {
	return func(a, b Span[E, T]) int {
		var diff int = cmp(a.Begin, b.Begin)
		if diff == 0 {
			return cmp(b.End, a.End)
		}
		return diff
	}
}

func BuildContainedBy[E any, T any](cmp func(E, E) int) func(a, b *Span[E, T]) (int, int) {
	return func(a, b *Span[E, T]) (int, int) {
		return cmp(a.Begin, b.Begin), cmp(a.End, b.End)
	}
}

func BuidAccumulator[E any, T any](
	ContainedBy func(a, b *Span[E, T]) (int, int),
	cmp func(a, b E) int,
) func() func(b *Span[E, T]) *AccumulatedSpanSet[E, T] {
	return func() func(b *Span[E, T]) *AccumulatedSpanSet[E, T] {
		var rss = &AccumulatedSpanSet[E, T]{}
		rss.Span = nil
		rss.Contains = nil

		return func(b *Span[E, T]) *AccumulatedSpanSet[E, T] {
			a := rss.Span
			if rss.Span == nil {
				rss.Span = b
			} else if cmp(a.End, b.Begin) < 0 {
				// reset our accumulator
				//fmt.Println("New Range")
				rss = nil
				rss = &AccumulatedSpanSet[E, T]{}
				rss.Span = b
				rss.Contains = nil
				return rss
			} else {
				x, y := ContainedBy(a, b)
				s := rss.Contains
				if x|y == 0 {
					rss.Span = a
				} else {
					var r = Span[E, T]{}
					if x < 0 {
						r.Begin = a.Begin
					} else {
						r.Begin = b.Begin
					}
					if y > 0 {
						r.End = a.End
					} else {
						r.End = b.End
					}
					rss.Span = &r

				}
				if s == nil {
					var c = []*Span[E, T]{a, b}
					rss.Contains = &c
				} else {
					*rss.Contains = append(*rss.Contains, b)
				}
			}
			return rss
		}
	}
}

func BuildContains[E any, T any](cmp func(a,b E) int) func (a *Span[E,T],b E) bool {
  return func(a *Span[E,T], b E) bool {
    return cmp(a.Begin,b)<1 && cmp(a.End,b) >-1
  }
}

func CreateCompare[E any, T any](cmp func(E, E) int) SpanTests[E, T] {
	var Compare = BuildCompare[E, T](cmp)
	var ContainedBy = BuildContainedBy[E, T](cmp)
	var SpanAccumulator = BuidAccumulator(ContainedBy, cmp)
  var Contains = BuildContains[E,T](cmp);
  var Overlap = BuildOverlap(Contains);
	var ops = SpanTests[E, T]{
		Compare:         Compare,
		ContainedBy:     ContainedBy,
		SpanAccumulator: SpanAccumulator,
		Cmp:             cmp,
    Contains:        Contains,
    Overlap:         Overlap,
	}
	return ops
}

func OrderedCreateCompare[E cmp.Ordered, T any]() SpanTests[E, T] {
	return CreateCompare[E, T](cmp.Compare)
}
