package spantools 

import (
	"cmp"
	"fmt"
)

type Span[E any, T any] struct {
	Begin E
	End   E
	Tag   *T
}

type SpanSlices[E any, T any] []Span[E, T]
type ResolvedSpanSet[E any, T any] struct {
	*Span[E, T]
	Contains *SpanSlices[E, T]
}

func BuildCreateSpan[E any, T any](cmp func(E, E) int) func(E, E, *T) (Span[E, T], error) {
	var createSpan = func(begin, end E, tag *T) (Span[E, T], error) {
		if cmp(end, begin) == 1 {
			return Span[E, T]{}, fmt.Errorf("Invalid range, Begin must be less than or equal to End")
		}
		newSpan := Span[E, T]{Begin: begin, End: end, Tag: tag}
		return newSpan, nil
	}
	return createSpan
}

type SpanTests[E any, T any] struct {
	Compare          func(a, b Span[E, T]) int
	ContainedBy      func(a, b Span[E, T]) (int, int)
	ResolveContainer func(a, b Span[E, T], s SpanSlices[E, T]) ResolvedSpanSet[E, T]
	CreateSpan      func(a, b E, tag *T) (Span[E, T], error)
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

func BuildContainedBy[E any, T any](cmp func(E, E) int) func(a, b Span[E, T]) (int, int) {
	return func(a, b Span[E, T]) (int, int) {
		return cmp(a.Begin, b.Begin), cmp(a.End, b.End)
	}
}

func BuildResolveContainer[E any, T any](ContainedBy func(a, b Span[E, T]) (int, int)) func(a, b Span[E, T], s SpanSlices[E, T]) ResolvedSpanSet[E, T] {
	return func(a, b Span[E, T], s SpanSlices[E, T]) ResolvedSpanSet[E, T] {
		var rs ResolvedSpanSet[E, T] = ResolvedSpanSet[E, T]{}
		x, y := ContainedBy(a, b)
		if x|y == 0 {
			rs.Span = &a
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
			rs.Span = &r

		}
		if s == nil {
			var c = SpanSlices[E, T]{a, b}
			rs.Contains = &c
		} else {
			var size = len(s) + 2
			var c = make(SpanSlices[E, T], size)
			copy(c, s)
			c[size-2] = a
			c[size-1] = b
			rs.Contains = &c
		}
		return rs
	}
}

func CreateCompare[E any, T any](cmp func(E, E) int) SpanTests[E, T] {
	var Compare = BuildCompare[E, T](cmp)
	var ContainedBy = BuildContainedBy[E, T](cmp)
	var ResolveContainer = BuildResolveContainer(ContainedBy)
	var CreateSpan = BuildCreateSpan[E, T](cmp)
	var ops = SpanTests[E, T]{
		Compare:          Compare,
		ContainedBy:      ContainedBy,
		ResolveContainer: ResolveContainer,
		CreateSpan:      CreateSpan,
	}
	return ops
}

func OrderedCreateCompare[E cmp.Ordered, T any]() SpanTests[E, T] {
  return CreateCompare[E,T](cmp.Compare);
}


