

package st

import (
	"slices"
	"testing"
)



var MultMultiiSet = []SpanBoundry[int, string]{
  &Span[int, string]{Begin: -1, End: 0},            // 0
  &Span[int, string]{Begin: 2, End: 2, Tag: &tagA}, //1
  &Span[int, string]{Begin: 2, End: 2, Tag: &tagB}, //1
  &Span[int, string]{Begin: 5, End: 6},             // 2
  &Span[int, string]{Begin: 9, End: 11},            // 3
  &Span[int, string]{Begin: 9, End: 11},            // 3
	&Span[int, string]{Begin: 12, End: 12},            // 4
}

// Data sets used to verify the range sort method works as expected
var testSets = [][][]SpanBoundry[int, string]{
	{
		// test set 0, All consumed by 1 range
		{
			// unsorted
			&Span[int, string]{Begin: -1, End: 0},
			&Span[int, string]{Begin: 0, End: 1},
			&Span[int, string]{Begin: -1, End: 0},
			&Span[int, string]{Begin: 0, End: 1},
			&Span[int, string]{Begin: -2, End: 2},
		},
		AllSet,
	},
	// test set 1, Seperate blocks with only 1 overlap
	{
		{
			&Span[int, string]{Begin: 2, End: 2},
			&Span[int, string]{Begin: 9, End: 11},
			&Span[int, string]{Begin: 5, End: 6},
			&Span[int, string]{Begin: 2, End: 2},
			&Span[int, string]{Begin: -1, End: 0},
		},
		MultiSet,
	},
}

var testDriver = NewOrderedSpanUtil[int, string]()

type SpanInt struct {
	SpanRef[int, string]
}

func TestSaneSpan(t *testing.T) {
	var begin = 2
	var end = 1
	var next = &SpanInt{
		SpanRef[int, string]{
			Begin: &begin,
			End:   &end,
		},
	}

	if next.GetTag() != nil {
		t.Errorf("Should get nil for our default tag")
		return
	}
	if nil == testDriver.Check(next, nil) {
		t.Errorf("Should get an error if begin is greater than end")
		return
	}
	begin = 0
	if nil != testDriver.Check(next, nil) {
		t.Errorf("Should get not get an error")
		return
	}

	var current = &Span[int, string]{Begin: 3, End: 3}
	if nil == testDriver.Check(next, current) {
		t.Errorf("current: %d->%d, should be before next: %d->%d", current.GetBegin(), current.GetEnd(), next.GetBegin(), next.GetEnd())
	}

}

func TestNewSpan(t *testing.T) {
	var span, err = testDriver.NewSpan(1, 2, &tagA)
	if err != nil {
		t.Errorf("Creation of new valid span failed")
		return
	}
	if span == nil {
		t.Errorf("Invalid return span value")
		return
	}
	if span.GetTag() != &tagA {
		t.Errorf("Invalid return span tag pointer")
		return
	}
	if span.Begin != 1 || span.End != 2 || span.Tag != &tagA {
		t.Errorf("Invalid return span content")
		return
	}
	span, err = testDriver.NewSpan(2, 1, &tagA)
	if err == nil {
		t.Errorf("Should have an error here")
	}
}

// Validates sort operation, by sorting slices and compairing the the sorted elements to a manually sorted array.
func TestOneContainerForAllSort(t *testing.T) {
	for setId, testSet := range testSets {
		var unsorted = make([]SpanBoundry[int, string], len(testSet[0]))
		copy(unsorted, testSet[0])
		slices.SortFunc(unsorted, testDriver.Compare)

		var sorted = testSet[1]
		for idx, span := range unsorted {
			var expected = sorted[idx]
			if span.GetBegin() != expected.GetBegin() || span.GetEnd() != expected.GetEnd() {
				t.Errorf("Error comparing test sort data set: %d, row: %d, Expected: %d,%d, Got: %d,%d", setId, idx, expected.GetBegin(), expected.GetEnd(), span.GetBegin(), span.GetEnd())
			}
		}
	}

}

// Test the Accumulator function and validates that the overlaps are generated correctly.
func TestConsolidate(t *testing.T) {
	var container = AllSet[0]
	var s = testDriver.NewSpanOverlapAccumulator()
	s.Validate = false
	for idx, span := range AllSet {
		var res = s.Accumulate(span)
		if container.GetBegin() != res.GetBegin() || container.GetEnd() != res.GetEnd() {
			t.Errorf("Container out of bounds in element: %d", idx)
		}
		if idx == 0 {
			if res.Contains != nil {
				t.Errorf("First element is always natural")
			}
		} else {
			if res.Contains == nil {
				t.Errorf("Must never be nil when beyond first element")
			}
			if len(*res.Contains) != idx+1 {
				t.Errorf("Range acumulation invalid at element: %d, size was: %d", idx, len(*res.Contains))
			}
		}
	}
}

// Tests the Accumulator function with multiple both overlapping and non overlapping Spans.
func TestMergeMultiple(t *testing.T) {
	var accumulator = testDriver.NewSpanOverlapAccumulator()
	for idx, span := range MultiSet {
		var res = accumulator.Accumulate(span)
		if span.GetBegin() != res.GetBegin() || span.GetEnd() != res.GetEnd() {
			t.Errorf("Range missmatch, expected: %d->%d, got: %d->%d", span.GetBegin(), span.GetEnd(), res.GetBegin(), res.GetEnd())
		}
		switch idx {
		case 0, 1, 3, 4:
			if res.Contains != nil {
				t.Errorf("First should be nil")
				return
			}
		case 2:
			if res.Contains == nil {
				t.Errorf("Container should not be nil")
				return
			}
			if len(*res.Contains) != 2 {
				t.Errorf("Container should have 2 elements")
				return
			}
			var list = *res.Contains
			var check = *list[0].GetTag() + *list[1].GetTag()
			if check != "ab" {
				t.Errorf("tag validation failed")
			}
		}
	}
}

// Tests the Accumulator function to make sure growth works as expected for accumulated spans.
func TestGrowth(t *testing.T) {
	src := []Span[int, string]{
		// sorted
		{Begin: -2, End: -1},
		{Begin: -1, End: 0},
		{Begin: 1, End: 1},
	}
	var s = testDriver.NewSpanOverlapAccumulator()
	var lastRes *OverlappingSpanSets[int, string] = nil
	for idx, span := range src {
		res := s.Accumulate(&span)
		if res.GetTag() != nil {
			t.Errorf("Tag should be nil")
		}
		switch idx {
		case 0:
			if res.GetBegin() != -2 || res.GetEnd() != -1 || res.Contains != nil {
				t.Errorf("Bad Range on element 0")
			}
		case 1:
			if res.GetBegin() != -2 || res.GetEnd() != -0 {
				t.Errorf("Bad Range on element 1")
			}
			if res.IsUnique() {
				t.Errorf("Expected to be stand alone, but contians multiple elements")
			}
			if *lastRes != *res {
				t.Errorf("Did not expect new result!")
				return
			}
		case 2:
			if res.GetBegin() != 1 || res.GetEnd() != 1 {
				t.Errorf("Bad Range")
			}
			if *res == *lastRes {
				t.Errorf("Expected new result")
			}
			if nil != res.GetContains() {
				t.Errorf("Invalid Contains")
			}
		}
		lastRes = res
	}
}



// Negative and positive overlap span testing.
func TestOverlaps(t *testing.T) {
	var a = &Span[int, string]{Begin: 0, End: 1}
	var b = &Span[int, string]{Begin: 1, End: 2}
	if !testDriver.Overlap(a, b) {
		t.Errorf("Expected a and b to overlap")
	}
	if !testDriver.Overlap(b, a) {
		t.Errorf("Expected a and b to overlap")
	}
	a = &Span[int, string]{Begin: 0, End: 1}
	b = &Span[int, string]{Begin: 2, End: 2}
	if testDriver.Overlap(a, b) {
		t.Errorf("Invalid overlap of a and b ")
	}
	if testDriver.Overlap(a, b) {
		t.Errorf("Invalid overlap of b and a ")
	}
}

func TestMultiAccumulateSet(t *testing.T) {
	var acc = testDriver.NewSpanOverlapAccumulator()
	var first = acc.Accumulate(MultiSet[0])
	var next = acc.Accumulate(MultiSet[1])

	if first.GetBegin() != -1 || first.GetEnd() != 0 {
		t.Errorf("Expected -1,0 got: %d,%d", first.GetBegin(), first.GetEnd())
		return
	}

	if first == next {
		t.Errorf("First and next should not be the same!")
		return
	}
	if next.GetBegin() != 2 || next.GetEnd() != 2 {
		t.Errorf("Expected 2,2 got: %d,%d", next.GetBegin(), next.GetEnd())
		return
	}
	first = next
	next = acc.Accumulate(MultiSet[2])
	if next.GetBegin() != 2 || next.GetEnd() != 2 {
		t.Errorf("Expected 2,2 got: %d,%d", next.GetBegin(), next.GetEnd())
		return
	}
	if next.SrcBegin != 1 || next.SrcEnd != 2 {
		t.Errorf("Bad source index points! Expected: 1,2 got: %d,%d", next.SrcBegin, next.SrcEnd)
		return
	}
	if first != next {
		t.Errorf("First and next must be the same!")
		return
	}
	first = next
	next = acc.Accumulate(MultiSet[3])
	if next.GetBegin() != 5 || next.GetEnd() != 6 {
		t.Errorf("Expected 5,6 got: %d,%d", next.GetBegin(), next.GetEnd())
		return
	}
  if(next.SrcBegin!=3 || next.SrcEnd!=3) {
		t.Errorf("Expected 3,3 got: %d,%d", next.SrcBegin, next.SrcEnd)
    panic("STOP TESTING")
  }
	if first == next {
		t.Errorf("First and next should not be the same!")
		return
	}
	first = next
	next = acc.Accumulate(MultiSet[4])
	if next.GetBegin() != 9 || next.GetEnd() != 11 {
		t.Errorf("Expected 9,11 got: %d,%d", next.GetBegin(), next.GetEnd())
		return
	}
	if first == next {
		t.Errorf("First and next should not be the same!")
		return
	}
}






func TestChanAccumulatro(t *testing.T) {
	var c = make(chan SpanBoundry[int, string], len(MultiSet))
	for _, span := range MultiSet {
		c <- span
	}
	close(c)
	var count = 0
	for range testDriver.NewSpanOverlapAccumulator().ChanIterFactory(c) {
		count++
	}
	if count != 4 {
		t.Errorf("Expected a total 4 for got: %d", count)
	}
	c = make(chan SpanBoundry[int, string], len(AllSet))
	for _, span := range AllSet {
		c <- span
	}
	close(c)
	count = 0
	for range testDriver.NewSpanOverlapAccumulator().ChanIterFactory(c) {
		count++
	}
	if count != 1 {
		t.Errorf("Expected a total 1 for got: %d", count)
	}

	count = 0
	for range testDriver.NewSpanOverlapAccumulator().ChanIterFactory(nil) {
		count++
	}
	if count != 0 {
		t.Errorf("Expected a total 0 for got: %d", count)
	}
	c = make(chan SpanBoundry[int, string], 1)
	c <- &Span[int, string]{Begin: 11, End: 5}
	close(c)
	count = 0
	for range testDriver.NewSpanOverlapAccumulator().ChanIterFactory(c) {
		count++
	}
	if count != 0 {
		t.Errorf("Exersize Error, failed? Expected a total 1 for got: %d", count)
	}

	c = make(chan SpanBoundry[int, string], len(MultiSet))
	for _, span := range MultiSet {
		c <- span
	}
	close(c)
	count = 0
	for range testDriver.NewSpanOverlapAccumulator().ChanIterFactory(c) {
		count++
		break
	}
	if count != 1 {
		t.Errorf("Force, yeild test coverage... Expected a total 1 for got: %d", count)
	}
  

}





