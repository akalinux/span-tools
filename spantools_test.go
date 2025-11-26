package spans

import (
	"slices"
	"testing"
)

var AllSet = []Span[int, string]{
	// sorted
	{Begin: -2, End: 2},
	{Begin: -1, End: 0},
	{Begin: -1, End: 0},
	{Begin: 0, End: 1},
	{Begin: 0, End: 1},
}

var tagA = "a"
var tagB = "b"
var MultiSet = []Span[int, string]{
	{Begin: -1, End: 0},
	{Begin: 2, End: 2, Tag: &tagA},
	{Begin: 2, End: 2, Tag: &tagB},
	{Begin: 5, End: 6},
	{Begin: 9, End: 11},
}

// Data sets used to verify the range sort method works as expected
var testSets = [][][]Span[int, string]{
	{
		// test set 0, All consumed by 1 range
		{
			// unsorted
			{Begin: -1, End: 0},
			{Begin: 0, End: 1},
			{Begin: -1, End: 0},
			{Begin: 0, End: 1},
			{Begin: -2, End: 2},
		},
		AllSet,
	},
	// test set 1, Seperate blocks with only 1 overlap
	{
		{
			{Begin: 2, End: 2},
			{Begin: 9, End: 11},
			{Begin: 5, End: 6},
			{Begin: 2, End: 2},
			{Begin: -1, End: 0},
		},
		MultiSet,
	},
}

var driver = OrderedCreateCompare[int, string]()

// Validates sort operation, by sorting slices and compairing the the sorted elements to a manually sorted array.
func TestOneContainerForAllSort(t *testing.T) {
	for setId, testSet := range testSets {
		var unsorted = make([]Span[int, string], len(testSet[0]))
		copy(unsorted, testSet[0])
		slices.SortFunc(unsorted, driver.Compare)

		var sorted = testSet[1]
		for idx, span := range unsorted {
			var expected = sorted[idx]
			if span.Begin != expected.Begin || span.End != expected.End {
				t.Errorf("Error comparing test sort data set: %d, row: %d, Expected: %d,%d, Got: %d,%d", setId, idx, expected.Begin, expected.End, span.Begin, span.End)
			}
		}
	}

}

// Test the Accumulator function and validates that the overlaps are generated correctly.
func TestConsolidate(t *testing.T) {
	var container = AllSet[0]
	var accumulator = driver.SpanAccumulator()
	for idx, span := range AllSet {
		var res = accumulator(&span)
		if container.Begin != res.Begin || container.End != res.End {
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
	var accumulator = driver.SpanAccumulator()
	for idx, span := range MultiSet {
		var res = accumulator(&span)
		if span.Begin != res.Begin || span.End != res.End {
			t.Errorf("Range missmatch, expected: %d->%d, got: %d->%d", span.Begin, span.End, res.Begin, res.End)
		}
		switch idx {
		case 0, 1, 3, 4:
			if res.Contains != nil {
				t.Errorf("First should be nil")
			}
		case 2:
			if res.Contains == nil || len(*res.Contains) != 2 {
				t.Errorf("Container should have 2 elements")
			}
			var list = *res.Contains
			var check = *list[0].Tag + *list[1].Tag
			if check != "ab" {
				t.Errorf("tag validation failed")
			}
		}
	}
}

//  Tests the Accumulator function to make sure growth works as expected for accumulated spans.
func TestGrowth(t *testing.T) {
	src := []Span[int, string]{
		// sorted
		{Begin: -2, End: -1},
		{Begin: -1, End: 0},
		{Begin: 1, End: 1},
	}
	var accumulator = driver.SpanAccumulator()
	var lastRes *AccumulatedSpanSet[int, string] = nil
	for idx, span := range src {
		res := accumulator(&span)
		switch idx {
		case 0:
			if res.Begin != -2 || res.End != -1 || res.Contains != nil {
				t.Errorf("Bad Range on element 0")
			}
		case 1:
			if res.Begin != -2 || res.End != -0 {
				t.Errorf("Bad Range on element 1")
			}
      if(res.Contains == nil) {
				t.Errorf("Expected contains for element 1")
      } 
      if lastRes != res {
        t.Errorf("Did not expect new range!")
        return;
      }
		case 2:
			if res.Begin != 1 || res.End != 1 {
				t.Errorf("Bad Range")
			}
			if res == lastRes {
				t.Errorf("Bad result")
			}
			if nil != res.Contains {
				t.Errorf("Invalid Contains")
			}
		}
		lastRes = res
	}
}

// Validates the inital range of a list of ranges
func TestFirstRange(t *testing.T) {
	var src *[]*Span[int, string] = &[]*Span[int, string]{
		{Begin: 2, End: 2},
		{Begin: 0, End: 1},
	}
	var span = driver.FirstSpan(src)
	if span.Begin != 0 || span.End != 1 {
		t.Errorf("Invalid start range")
	}
}

// Validates the creation of the next span based on the current span.
func TestNextRange(t *testing.T) {
	var src *[]*Span[int, string] = &[]*Span[int, string]{
		{Begin: 3, End: 4}, // 3,4, last valid range, should get nil after this
		{Begin: 2, End: 2}, // 2,2, first range
		{Begin: 0, End: 1}, // should ignore
	}
	var first = &Span[int, string]{Begin: 0, End: 1}
	var span = driver.NextSpan(first, src)
	if span == nil {
		t.Errorf("Should not have reached our end yet!")
		return
	}
	if span.Begin != 2 || span.End != 2 {
		t.Errorf("Invalid range, expected: 2->2, got %d->%d", span.Begin, span.End)
	}
	span = driver.NextSpan(span, src)
	if span == nil {
		t.Errorf("Should not have reached our end yet!")
		return
	}
	if span.Begin != 3 || span.End != 4 {
		t.Errorf("Invalid range, expected: 3->4, got %d->%d", span.Begin, span.End)
	}
	span = driver.NextSpan(span, src)
	if span != nil {
		t.Errorf("End expected!")
		return
	}
}

// Validates the creation of the next span when there are gaps.
func TestNextRangeGap(t *testing.T) {
  var src *[]*Span[int, string] = &[]*Span[int, string]{
    {Begin: 4, End: 5}, // 3,4, last valid range, should get nil after this
    {Begin: 2, End: 2}, // 2,2, first range
    {Begin: 0, End: 1}, // should ignore
  }
  var first = &Span[int, string]{Begin: 0, End: 0}
  var span = driver.NextSpan(first, src)
  if span == nil {
    t.Errorf("Should not have reached our end yet!")
    return
  }
  if span.Begin != 2 || span.End != 2 {
    t.Errorf("Invalid range, expected: 2->2, got %d->%d", span.Begin, span.End)
    return;
  }
  span = driver.NextSpan(span, src)
  if span == nil {
    t.Errorf("Should not have reached our end yet!")
    return
  }
  if span.Begin != 4 || span.End != 5 {
    t.Errorf("Invalid range, expected: 4->5, got %d->%d", span.Begin, span.End)
    return;
  }
  span = driver.NextSpan(span, src)
  if span != nil {
    t.Errorf("End expected!")
    return
  }
}

// Used to test overlaping span generation.
func testNextOverlaps(t *testing.T,src *[]*Span[int,string]) {
  
  var first = &Span[int, string]{Begin: -1, End: 0}
  var span = driver.NextSpan(first, src)
  if span == nil {
    t.Errorf("Should not have reached our end yet!")
    return
  }
  if span.Begin != 1 || span.End != 1 {
    t.Errorf("Invalid range, expected: 1->1, got %d->%d", span.Begin, span.End)
    return;
  }
  span = driver.NextSpan(span, src)
  if span == nil {
    t.Errorf("Should not have reached our end yet!")
    return
  }
  if span.Begin != 2 || span.End != 3 {
    t.Errorf("Invalid range, expected: 2->3, got %d->%d", span.Begin, span.End)
  }
  span = driver.NextSpan(span, src)
  if span != nil {
    t.Errorf("End expected!")
    return
  }
}
// Validates the creation of the next span with overlaps.
func TestNextRangeOverlaps(t *testing.T) {
  var src *[]*Span[int, string] = &[]*Span[int, string]{
    {Begin: 2, End: 3}, // 2,3, next range is nil
    {Begin: 1, End: 3}, // overlaps with 0 and 2
    {Begin: 0, End: 1}, // 1,1, first range
  }
  testNextOverlaps(t,src);
}

// Validates the creation of the next span with a different data set.
func TestNextRangeOverlapsReverseOrder(t *testing.T) {
  var src *[]*Span[int, string] = &[]*Span[int, string]{
    {Begin: 0, End: 1}, // 1,1, first range
    {Begin: 1, End: 3}, // overlaps with 0 and 2
    {Begin: 2, End: 3}, // 2,3, next range is nil
  }
  testNextOverlaps(t,src);
}

func TestNextRangeOverlapsMixedOrder(t *testing.T) {
  var src *[]*Span[int, string] = &[]*Span[int, string]{
    {Begin: 1, End: 3}, // overlaps with 1 and 2
    {Begin: 0, End: 1}, // 1,1, first range
    {Begin: 2, End: 3}, // 2,3, next range is nil
  }
  testNextOverlaps(t,src);
}

// Negative and positive overlap span testing.
func TestOverlaps(t *testing.T) {
  var a =&Span[int,string]{Begin: 0,End: 1}
  var b =&Span[int,string]{Begin: 1,End: 2}
  if(!driver.Overlap(a,b)) {
    t.Errorf("Expected a and b to overlap");
  }
  if(!driver.Overlap(b,a)) {
    t.Errorf("Expected a and b to overlap");
  }
  a =&Span[int,string]{Begin: 0,End: 1}
  b =&Span[int,string]{Begin: 2,End: 2}
  if(driver.Overlap(a,b)) {
    t.Errorf("Invalid overlap of a and b ");
  }
  if(driver.Overlap(a,b)) {
    t.Errorf("Invalid overlap of b and a ");
  }
}
