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
	Compare func(a, b Span[E, T]) int

	// Compuates the smallest Begin value and largest End value given the input a,b *Span[E,T].
	// Return values: First int represents the smallest Span[E,T] Begin, Second value represents the largest Span[E,T] End.
	ContainedBy func(a, b *Span[E, T]) (int, int)

	// Returns true if the a *Span[E,T] contains the value b E.
	Contains func(a *Span[E, T], b E) bool

	// Creates a closure used to accumualte new Span overlaps.
	// Each time a new pointer is returned, a new Span overlap has been found.
	// Values passed to the accumuator(s *Span[E,T]) function must be in the oder produced by slices.SortFunc(SpanTests[E,T].Compare)
	SpanAccumulator func() func(b *Span[E, T]) *AccumulatedSpanSet[E, T]

	// Comapres a to b: returns -1,0,1 based on the 2 values compared
	Cmp func(a, b E) int

	// Used to compare 2 spans, returns true if they overlap
	Overlap func(a, b *Span[E, T]) bool

	// Returns the smallest start Span[E,T]. Assumes Compare order.
	FirstSpan func(list *[]*Span[E, T]) *Span[E, T]

	// Returns the next *Span[E,T] or nil if no remaining *Span[E,T] overlaps.  Assumes Compare order.
	NextSpan func(first *Span[E, T], list *[]*Span[E, T]) *Span[E, T]
}

func buildOverlap[E any, T any](contains func(a *Span[E, T], b E) bool) func(a, b *Span[E, T]) bool {
	return func(a, b *Span[E, T]) bool {
		return contains(a, b.Begin) || contains(a, b.End) || (contains(b, a.Begin) || contains(b, a.End))
	}
}

func buildCompare[E any, T any](cmp func(E, E) int) func(a, b Span[E, T]) int {
	return func(a, b Span[E, T]) int {
		var diff int = cmp(a.Begin, b.Begin)
		if diff == 0 {
			return cmp(b.End, a.End)
		}
		return diff
	}
}

func buildContainedBy[E any, T any](cmp func(E, E) int) func(a, b *Span[E, T]) (int, int) {
	return func(a, b *Span[E, T]) (int, int) {
		return cmp(a.Begin, b.Begin), cmp(a.End, b.End)
	}
}

func buidAccumulator[E any, T any](
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
					}
					if y > 0 {
						r.End = a.End
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

func buildContains[E any, T any](cmp func(a, b E) int) func(a *Span[E, T], b E) bool {
	return func(a *Span[E, T], b E) bool {
		return cmp(a.Begin, b) < 1 && cmp(a.End, b) > -1
	}
}

func buildFirstSpan[E any, T any](cmp func(a, b E) int) func(list *[]*Span[E, T]) *Span[E, T] {
	return func(list *[]*Span[E, T]) *Span[E, T] {
		var span = Span[E, T]{Begin: (*list)[0].Begin, End: (*list)[0].End}
		var last = len(*list)
		for i := 1; i < last; i++ {
			var check = (*list)[i]
			if cmp(check.Begin, span.Begin) == -1 {
				span.Begin = check.Begin
			}
			if cmp(check.End, span.End) == -1 {
				span.End = check.End
			}
		}
		return &span
	}
}

func buildNextSpan[E any, T any](cmp func(a, b E) int) func(start *Span[E, T], list *[]*Span[E, T]) *Span[E, T] {
	return func(start *Span[E, T], list *[]*Span[E, T]) *Span[E, T] {
		var begin *E = nil
		var end *E = nil
		for _, check := range *list {
			if begin == nil {
				if cmp(check.Begin, start.End) > 0 {
					begin = &check.Begin
					end = &check.End
				} else if cmp(check.End, start.End) > 0 {
					begin = &check.End
					end = &check.End
				} 
			} else {
				if cmp(check.Begin, start.End) > 0 {
					if cmp(*begin, check.Begin) > 0 {
						begin = &check.Begin
					}
				}
				if cmp(*begin,check.End) < 1 && cmp(check.End, start.End) > 0 {
					if cmp(*end, check.End) > 0 {
						end = &check.End
					} 
				}
			}
		}
		if begin != nil {
			return &Span[E, T]{Begin: *begin, End: *end}
		}
		return nil
	}
}

// Factory for the creation of SpanTest[E,T], the cmp func(E,E) int, is used to compare the Begin and End values of a given Span[E,T].
func CreateCompare[E any, T any](cmp func(E, E) int) SpanTests[E, T] {
	var Compare = buildCompare[E, T](cmp)
	var ContainedBy = buildContainedBy[E, T](cmp)
	var SpanAccumulator = buidAccumulator(ContainedBy, cmp)
	var Contains = buildContains[E, T](cmp)
	var Overlap = buildOverlap(Contains)
	var FirstSpan = buildFirstSpan[E, T](cmp)
	var NextSpan = buildNextSpan[E, T](cmp)
	var ops = SpanTests[E, T]{
		Compare:         Compare,
		ContainedBy:     ContainedBy,
		SpanAccumulator: SpanAccumulator,
		Cmp:             cmp,
		Contains:        Contains,
		Overlap:         Overlap,
		FirstSpan:       FirstSpan,
		NextSpan:        NextSpan,
	}
	return ops
}

// Factory for the creation of SpanTest[E,T], where the type of E is represented by any data type supported by E cmp.Ordered.
func OrderedCreateCompare[E cmp.Ordered, T any]() SpanTests[E, T] {
	return CreateCompare[E, T](cmp.Compare)
}
