package rangecmp; 

import (
	"cmp"
	"fmt"
)

type Range[E any, T any] struct {
	Begin E
	End   E
	Tag   *T
}

type RangeSlices[E any, T any] []Range[E, T]
type ResolvedRangeSet[E any, T any] struct {
	*Range[E, T]
	Contains *RangeSlices[E, T]
}

func BuildCreateRange[E any, T any](cmp func(E, E) int) func(E, E, *T) (Range[E, T], error) {
	var createRange = func(begin, end E, tag *T) (Range[E, T], error) {
		if cmp(end, begin) == 1 {
			return Range[E, T]{}, fmt.Errorf("Invalid range, Begin must be less than or equal to End")
		}
		newRange := Range[E, T]{Begin: begin, End: end, Tag: tag}
		return newRange, nil
	}
	return createRange
}

type RangeTests[E any, T any] struct {
	Compare          func(a, b Range[E, T]) int
	ContainedBy      func(a, b Range[E, T]) (int, int)
	ResolveContainer func(a, b Range[E, T], s RangeSlices[E, T]) ResolvedRangeSet[E, T]
	CreateRange      func(a, b E, tag *T) (Range[E, T], error)
}

func BuildCompare[E any, T any](cmp func(E, E) int) func(a, b Range[E, T]) int {
	return func(a, b Range[E, T]) int {
		var diff int = cmp(a.Begin, b.Begin)
		if diff == 0 {
			return cmp(b.End, a.End)
		}
		return diff
	}
}

func BuildContainedBy[E any, T any](cmp func(E, E) int) func(a, b Range[E, T]) (int, int) {
	return func(a, b Range[E, T]) (int, int) {
		return cmp(a.Begin, b.Begin), cmp(a.End, b.End)
	}
}

func BuildResolveContainer[E any, T any](ContainedBy func(a, b Range[E, T]) (int, int)) func(a, b Range[E, T], s RangeSlices[E, T]) ResolvedRangeSet[E, T] {
	return func(a, b Range[E, T], s RangeSlices[E, T]) ResolvedRangeSet[E, T] {
		var rs ResolvedRangeSet[E, T] = ResolvedRangeSet[E, T]{}
		x, y := ContainedBy(a, b)
		if x|y == 0 {
			rs.Range = &a
		} else {
			var r = Range[E, T]{}
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
			rs.Range = &r

		}
		if s == nil {
			var c = RangeSlices[E, T]{a, b}
			rs.Contains = &c
		} else {
			var size = len(s) + 2
			var c = make(RangeSlices[E, T], size)
			copy(c, s)
			c[size-2] = a
			c[size-1] = b
			rs.Contains = &c
		}
		return rs
	}
}

func CreateCompare[E any, T any](cmp func(E, E) int) RangeTests[E, T] {
	var Compare = BuildCompare[E, T](cmp)
	var ContainedBy = BuildContainedBy[E, T](cmp)
	var ResolveContainer = BuildResolveContainer(ContainedBy)
	var CreateRange = BuildCreateRange[E, T](cmp)
	var ops = RangeTests[E, T]{
		Compare:          Compare,
		ContainedBy:      ContainedBy,
		ResolveContainer: ResolveContainer,
		CreateRange:      CreateRange,
	}
	return ops
}

func OrderedCreateCompare[E cmp.Ordered, T any]() RangeTests[E, T] {
  return CreateCompare[E,T](cmp.Compare);
}


