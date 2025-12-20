package st

import (
	"errors"
	"iter"
)

// Core of the span utilties: Provides methos for processing ranges.
// Its recommended that an instance of this structure be created via the constructor util methods such as NewSpanUtil(Cmp) or NewOrderedSpanUtil(),
type SpanUtil[E any] struct {

	// Compare function.  This function should be atomic and be able t compare the E type by return -1,0,1.
	Cmp func(a, b E) int

	// Turns validation on for new child objects created.
	Validate bool

	// Denots if a tag is required, when true tag values cannot be nil.
	TagRequired bool

	// Next value function, should return the next E.
	// The new E value must always be greater than the argument passed in
	Next func(e E) E
	
	
	// Global default flag denoting if overlaps that are ajacent should be consolidated.
	// when truee: 1,2 and 2,3 
}

// This method is used to verify the sanyty of the next and current value.
// The comparison operation is performed in 2 stages:
// 1. next.GetBegin() must be less than or equal to next.GetEnd().
// 2. When the current value is not nil, then next must come after current.
// Returns nil when checks pass, the error is not nill when checks fail.
func (s *SpanUtil[E]) Check(next, current SpanBoundry[E]) error {

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
func (s *SpanUtil[E]) GetP(x E) *E {
	return &x
}

// Creates an instance of *SpanUtil[E], the value of cmp is expected to be able to compare the Span.Begin and Span.End values.
// See: [cmp.Compare] for more info.
//
// The default SpanFormat is set to: "Span: [%s -> %s], Tag: %s"
//
// [cmp.Compare]: https://pkg.go.dev/github.com/google/go-cmp/cmp#Comparer
func NewSpanUtil[E any](cmp func(a, b E) int, next func(e E) E) *SpanUtil[E] {
	return &SpanUtil[E]{
		Cmp:  cmp,
		Next: next,
	}
}

// This method is used to sort slice of spans in the accumulation order.
// For more details see: [slices.SortFunc].
//
// [slices.SortFunc]: https://pkg.go.dev/slices#SortedFunc
func (s *SpanUtil[E]) Compare(a, b SpanBoundry[E]) int {
	var diff int = s.Cmp(a.GetBegin(), b.GetBegin())
	if diff == 0 {
		return s.Cmp(b.GetEnd(), a.GetEnd())
	}
	return diff
}

// Returns true if a contains b.
func (s *SpanUtil[E]) Contains(a SpanBoundry[E], b E) bool {
	return s.Cmp(a.GetBegin(), b) < 1 && s.Cmp(a.GetEnd(), b) > -1
}

// Returns true if a overlaps with b or if be overlaps with a.
func (s *SpanUtil[E]) Overlap(a, b SpanBoundry[E]) bool {
	return s.Contains(a, b.GetBegin()) || s.Contains(a, b.GetEnd()) || s.Contains(b, a.GetEnd()) || s.Contains(b, a.GetEnd())
}

// This method is used to determin the outer bounds of ranges a and b.
// The first int represents comparing a.Begin to b.Begin and the second int represents comparing a.End to b.End.
func (s *SpanUtil[E]) ContainedBy(a, b SpanBoundry[E]) (int, int) {
	return s.Cmp(a.GetBegin(), b.GetBegin()), s.Cmp(a.GetEnd(), b.GetEnd())
}

// Creates a new span, error is nill unless a is greater than b.
func (s *SpanUtil[E]) NewSpan(a, b E) (*Span[E], error) {
	if s.Cmp(a, b) > 0 {
		return nil, errors.New("Value a is greater than value b")
	}
	return &Span[E]{Begin: a, End: b}, nil
}

// This method returns the first smallest span from list.
// The new SpanBoundry[E] will contian the smallest GetBegin() and the smallest GetEnd().
// Usage(FirstSpan): Is in conjuntion with NextSpan(span,list) to iterate through a set
func (s *SpanUtil[E]) FirstSpan(list *[]SpanBoundry[E]) SpanBoundry[E] {
	var span = &Span[E]{Begin: (*list)[0].GetBegin(), End: (*list)[0].GetEnd()}
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

// Factory interface for the creation of SpanOverlapAccumulator[E].
func (s *SpanUtil[E]) NewSpanOverlapAccumulator() *SpanOverlapAccumulator[E] {
	return &SpanOverlapAccumulator[E]{
		Validate: s.Validate,
		SpanUtil: s,
		Rss:      &OverlappingSpanSets[E]{Contains: nil, Span: nil},
		Pos:      -1,
	}
}

// This method takes a iter.Seq2 iterator of OverlappingSpanSets and initalizes the ColumnOverlapAccumulator struct.
//
// # Warning
//
// This methos creates an [iter.Pull2] and exposes the resulting functions in the returned struct pointer. If you are using this method outside of the normal
// operations, you should a setup a defer call to  ColumnOverlapAccumulator[E].Close() method to clean the instance up in order to prevent memory leaks or undefined behavior.
//
// [iter.Pull2]: https://pkg.go.dev/iter#hdr-Pulling_Values
func (s *SpanUtil[E]) ColumnOverlapFactory(driver iter.Seq2[int, *OverlappingSpanSets[E]]) *ColumnOverlapAccumulator[E] {
	var next, stop = iter.Pull2(driver)
	return s.ColumnOverlapFactoryBuilder(next, stop)
}

// This method takes the next and stop functions and creates a new fully initalized instance of ColumnOverlapAccumulator[E].
func (s *SpanUtil[E]) ColumnOverlapFactoryBuilder(next func() (int, *OverlappingSpanSets[E], bool), stop func()) *ColumnOverlapAccumulator[E] {
	var res = &ColumnOverlapAccumulator[E]{}
	res.ItrStop = stop
	res.ItrGetNext = next
	res.Util = s
	var _, current, ok = res.ItrGetNext()
	if ok {
		res.Next = current
	}
	return res
}

func (s *SpanUtil[E]) NewColumnSets() *ColumnSets[E] {
	return &ColumnSets[E]{
		Util: s,
	}
}

func (s *SpanUtil[E]) CreateOverlapSpan(list *[]SpanBoundry[E]) SpanBoundry[E] {

	var begin *E
	var end *E
	var Cmp = s.Cmp
	var res SpanBoundry[E] = &Span[E]{}
	for _, span := range *list {
		if begin == nil {
			begin = s.GetP(span.GetBegin())
			end = s.GetP(span.GetEnd())
			continue
		}

		if Cmp(span.GetBegin(), *begin) > 0 {
			begin = s.GetP(span.GetBegin())
		}
		if Cmp(*end, span.GetEnd()) > 0 {
			end = s.GetP(span.GetEnd())
		}
	}
	if Cmp(*begin, *end) < 1 {
		res = &Span[E]{Begin: *begin, End: *end}
	}
	return res
}

// Finds the next common overlapping SpanBoundry[E] span in list after start, 
// or nil if no next overlap is found after start in list.
func (s *SpanUtil[E]) NextSpan(start SpanBoundry[E], list *[]SpanBoundry[E]) (SpanBoundry[E],bool) {
	var min = s.Next(start.GetEnd())
	var end *E
	var res SpanBoundry[E] 
	var Cmp = s.Cmp
	var ok=false
	for _, span := range *list {
		var cmp = Cmp(span.GetEnd(), min)
		if end == nil {
			if cmp > -1 {
				end = s.GetP(span.GetEnd())
				cmp = Cmp(span.GetBegin(), min)
				if cmp > 0 && Cmp(span.GetBegin(), *end) < 0 {
					end = s.GetP(span.GetBegin())
				}

			}
			continue
		} else if cmp > -1 && Cmp(*end, span.GetEnd()) > 0 {
			end = s.GetP(span.GetEnd())
		}
		cmp = Cmp(span.GetBegin(), min)
		if cmp > 0 && Cmp(span.GetBegin(), *end) < 0 {
			end = s.GetP(span.GetBegin())
		}
	}
	if end != nil {
		var tmp = &Span[E]{Begin: min, End: *end}
		var ol = []SpanBoundry[E]{}
		copy(ol, *list)
		ol = append(ol, tmp)
		res = s.CreateOverlapSpan(&ol)
		ok=true
	}

	return res,ok
}
