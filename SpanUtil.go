package st

import (
	"cmp"
	"errors"
)

// Core of the span utilties: Provides methos for processing ranges.
type SpanUtil[E any, T any] struct {
	Cmp         func(a, b E) int
	Validate    bool
	TagRequired bool
}

func (s *SpanUtil[E, T]) Check(next, current SpanBoundry[E, T]) error {

	if s.Cmp(next.GetBegin(), next.GetEnd()) > 0 {
		return errors.New("GetBegin must be less than or equal to GetEnd")
	}

	if current != nil {

		if s.Compare(current, next) > 0 {

			return errors.New("SpanBoundry out of sequence")
		}
	}
	return nil
}

// Wrapper function to return a pointer to the value passed in.
func (s *SpanUtil[E, T]) GetP(x E) *E {
	return &x
}

// Creates a instance of *SpanUtil[E cmp.Ordered,T], this can be used to process most span data sets.
func NewOrderedSpanUtil[E cmp.Ordered, T any]() *SpanUtil[E, T] {
	return NewSpanUtil[E, T](cmp.Compare)
}

// Creates an instance of *SpanUtil[E,T], the value of cmp is expected to be able to compare the Span.Begin and Span.End values.
// See: [cmp.Compare] for more info.
//
// The default SpanFormat is set to: "Span: [%s -> %s], Tag: %s"
//
// [cmp.Compare]: https://pkg.go.dev/github.com/google/go-cmp/cmp#Comparer
func NewSpanUtil[E any, T any](cmp func(a, b E) int) *SpanUtil[E, T] {
	return &SpanUtil[E, T]{Cmp: cmp}
}

// This method is used to sort slice of spans in the accumulation order.
// For more details see: [slices.SortFunc].
//
// [slices.SortFunc]: https://pkg.go.dev/slices#SortedFunc
func (s *SpanUtil[E, T]) Compare(a, b SpanBoundry[E, T]) int {
	var diff int = s.Cmp(a.GetBegin(), b.GetBegin())
	if diff == 0 {
		return s.Cmp(b.GetEnd(), a.GetEnd())
	}
	return diff
}

// Returns true if a contains b.
func (s *SpanUtil[E, T]) Contains(a SpanBoundry[E, T], b E) bool {
	return s.Cmp(a.GetBegin(), b) < 1 && s.Cmp(a.GetEnd(), b) > -1
}

// Returns true if a overlaps with b or if be overlaps with a.
func (s *SpanUtil[E, T]) Overlap(a, b SpanBoundry[E, T]) bool {
	return s.Contains(a, b.GetBegin()) || s.Contains(a, b.GetEnd()) || s.Contains(b, a.GetEnd()) || s.Contains(b, a.GetEnd())
}

// This method is used to determin the outer bounds of ranges a and b.
// The first int represents comparing a.Begin to b.Begin and the second int represents comparing a.End to b.End.
func (s *SpanUtil[E, T]) ContainedBy(a, b SpanBoundry[E, T]) (int, int) {
	return s.Cmp(a.GetBegin(), b.GetBegin()), s.Cmp(a.GetEnd(), b.GetEnd())
}

// Creates a new span, error is nill unless a is greater than b.
func (s *SpanUtil[E, T]) NewSpan(a, b E, tag *T) (*Span[E, T], error) {
	if s.Cmp(a, b) > 0 {
		return nil, errors.New("Value a is greater than value b")
	}
	return &Span[E, T]{Begin: a, End: b, Tag: tag}, nil
}

// This method returns the first smallest span from the slice of Span[E,T].
func (s *SpanUtil[E, T]) FirstSpan(list *[]SpanBoundry[E, T]) *Span[E, T] {
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
func (s *SpanUtil[E, T]) NextSpan(start SpanBoundry[E, T], list *[]SpanBoundry[E, T]) *Span[E, T] {
	var begin *E = nil
	var end *E = nil

	for _, check := range *list {
		if begin == nil {
			if s.Cmp(check.GetBegin(), start.GetEnd()) > 0 {
				begin = s.GetP(check.GetBegin())
				end = s.GetP(check.GetEnd())
			} else if s.Cmp(check.GetEnd(), start.GetEnd()) > 0 {
				begin = s.GetP(check.GetEnd())
				end = s.GetP(check.GetEnd())
			}
		} else {
			if s.Cmp(check.GetBegin(), start.GetEnd()) > 0 && s.Cmp(*begin, check.GetBegin()) > 0 {
				begin = s.GetP(check.GetBegin())
			}
			if s.Cmp(*begin, check.GetEnd()) < 1 && s.Cmp(check.GetEnd(), start.GetEnd()) > 0 && s.Cmp(*end, check.GetEnd()) > 0 {
				end = s.GetP(check.GetEnd())
			}
		}
	}
	if begin != nil {
		return &Span[E, T]{Begin: *begin, End: *end}
	}
	return nil
}

// Factory interface for the creation of SpanOverlapAccumulator[E,T].
func (s *SpanUtil[E, T]) NewSpanOverlapAccumulator() *SpanOverlapAccumulator[E, T] {
	return &SpanOverlapAccumulator[E, T]{
		Validate: s.Validate,
		SpanUtil: s,
		Rss:      &OverlappingSpanSets[E, T]{Contains: nil, Span: nil},
		Pos:      -1,
	}
}

