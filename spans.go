package spans

import (
	"cmp"
)

type Span[E any, T any] struct {
	Begin E
	End   E
	Tag   *T
}

type AccumulatedSpanSet[E any, T any] struct {
	*Span[E, T]
	Contains *[]*Span[E, T]
}

type SpanTests[E any, T any] struct {
	Compare         func(a, b Span[E, T]) int
	ContainedBy     func(a, b *Span[E, T]) (int, int)
	SpanAccumulator func() func(b *Span[E, T]) *AccumulatedSpanSet[E, T]
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
		var rss = AccumulatedSpanSet[E, T]{}
		rss.Span = nil
		rss.Contains = nil

		return func(b *Span[E, T]) *AccumulatedSpanSet[E, T] {
			a := rss.Span
			if rss.Span == nil {
				rss.Span = b
			} else if cmp(a.End, b.Begin) < 0 {
				// reset our accumulator
				rss = AccumulatedSpanSet[E, T]{}
				rss.Span = b;
        rss.Contains=nil
        return &rss;
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
			return &rss
		}
	}
}

func CreateCompare[E any, T any](cmp func(E, E) int) SpanTests[E, T] {
	var Compare = BuildCompare[E, T](cmp)
	var ContainedBy = BuildContainedBy[E, T](cmp)
	var SpanAccumulator = BuidAccumulator(ContainedBy, cmp)
	var ops = SpanTests[E, T]{
		Compare:         Compare,
		ContainedBy:     ContainedBy,
		SpanAccumulator: SpanAccumulator,
	}
	return ops
}

func OrderedCreateCompare[E cmp.Ordered, T any]() SpanTests[E, T] {
	return CreateCompare[E, T](cmp.Compare)
}
