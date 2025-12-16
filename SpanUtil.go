package st

import (
	"cmp"
	"errors"
	"iter"
)

// Core of the span utilties: Provides methos for processing ranges.
// Its recommended that an instance of this structure be created via the constructor util methods such as NewSpanUtil(Cmp) or NewOrderedSpanUtil(),
type SpanUtil[E any, T any] struct {

	// Compare function.  This function should be atomic and be able t compare the E type by return -1,0,1.
	Cmp func(a, b E) int

	// Turns validation on for new child objects created.
	Validate bool

	// Denots if a tag is required, when true tag values cannot be nil.
	TagRequired bool
}

// This method is used to verify the sanyty of the next and current value.
// The comparison operation is performed in 2 stages:
// 1. next.GetBegin() must be less than or equal to next.GetEnd().
// 2. When the current value is not nil, then next must come after current.
// Returns nil when checks pass, the error is not nill when checks fail.
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

// Factory interface for the creation of SpanOverlapAccumulator[E,T].
func (s *SpanUtil[E, T]) NewSpanOverlapAccumulator() *SpanOverlapAccumulator[E, T] {
	return &SpanOverlapAccumulator[E, T]{
		Validate: s.Validate,
		SpanUtil: s,
		Rss:      &OverlappingSpanSets[E, T]{Contains: nil, Span: nil},
		Pos:      -1,
	}
}

// This method takes a iter.Seq2 iterator of OverlappingSpanSets and initalizes the ColumnOverlapAccumulator struct.
//
// # Warning
//
// This methos creates an [iter.Pull2] and exposes the resulting functions in the returned struct pointer. If you are using this method outside of the normal
// operations, you should a setup a defer call to  ColumnOverlapAccumulator[E, T].Close() method to clean the instance up in order to prevent memory leaks or undefined behavior.
//
// [iter.Pull2]: https://pkg.go.dev/iter#hdr-Pulling_Values
func (s *SpanUtil[E, T]) ColumnOverlapFactory(driver iter.Seq2[int, *OverlappingSpanSets[E, T]]) *ColumnOverlapAccumulator[E, T] {
	var next, stop = iter.Pull2(driver)
	return s.ColumnOverlapFactoryBuilder(next, stop)
}

// This method takes the next and stop functions and creates a new fully initalized instance of ColumnOverlapAccumulator[E, T].
func (s *SpanUtil[E, T]) ColumnOverlapFactoryBuilder(next func() (int, *OverlappingSpanSets[E, T], bool), stop func()) *ColumnOverlapAccumulator[E, T] {
	var res = &ColumnOverlapAccumulator[E, T]{}
	res.ItrStop = stop
	res.ItrGetNext = next
	res.Util = s
	var id, current, ok = res.ItrGetNext()
	if ok {
		res.SrcPos = id
		res.Next = current
	}
	return res
}

func (s *SpanUtil[E, T]) NewColumnSets() *ColumnSets[E, T] {
	return &ColumnSets[E, T]{
		Columns: &[]*ColumnOverlapAccumulator[E, T]{},
		Active:  &[]bool{},
		Util:    s,
	}
}

func (s *SpanUtil[E, T]) GetNextBegin(current E, list *[]SpanBoundry[E, T]) *E {
	var next *E = nil
	for _, span := range *list {
		var cmp = s.Cmp(span.GetBegin(), current)
		if next == nil {
			if cmp > 0 {
				next = s.GetP(span.GetBegin())
			}
		} else if cmp > 0 && s.Cmp(span.GetBegin(), *next) < 0 {
			next = s.GetP(span.GetBegin())
		}
	}
	return next
}

func (s *SpanUtil[E, T]) GetNextEnd(current E, list *[]SpanBoundry[E, T]) *E {
	var next *E = nil
	for _, span := range *list {
		var cmp = s.Cmp(span.GetEnd(), current)
		if next == nil {
			if cmp > 0 {
				next = s.GetP(span.GetEnd())
			}
		} else if cmp > 0 && s.Cmp(span.GetEnd(), *next) < 0 {
			next = s.GetP(span.GetEnd())
		}
	}
	return next
}

// This method acts as a stateless iterator that,
// returns the next overlapping Span[E,T] or nill based on the start SpanBoundry[E,T] and the slice of spans.
// If all valid SpanBoundry[E,T] values have been exausted, nil is returned.
func (s *SpanUtil[E, T]) NextSpan(start SpanBoundry[E, T], list *[]SpanBoundry[E, T]) *Span[E, T] {
	var begin *E = s.GetNextBegin(start.GetEnd(), list)
	var end *E = s.GetNextEnd(start.GetEnd(), list)

	if begin != nil {
		var nextEnd *E = s.GetNextBegin(*begin, list)
		if nextEnd != nil {
			if s.Cmp(*nextEnd, *end) < 0 {
				end = nextEnd
			}
		}

		return &Span[E, T]{Begin: *begin, End: *end}
	} else if end != nil {
		return &Span[E, T]{Begin: start.GetEnd(), End: *end}
	}
	return nil
}
