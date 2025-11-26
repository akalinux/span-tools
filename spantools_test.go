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
			if res.Begin != -2 || res.End != -0 || res.Contains == nil || lastRes != res {
				t.Errorf("Bad Range on element 1")
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

func TestFirstRange(t *testing.T) {
  var src *[]*Span[int,string]=&[]*Span[int, string]{
    {Begin:2,End:2},
    {Begin:0,End:1},
  }
  var span=driver.FirstSpan(src);
  if(span.Begin!=0 || span.End!=1) {
    t.Errorf("Invalid start range")
  }
}

func TestNextRange(t *testing.T) {
  var src *[]*Span[int,string]=&[]*Span[int, string]{
    {Begin:3,End:4}, // 3,4, last valid range, should get nil after this
    {Begin:2,End:2}, // 2,2, first range
    {Begin:0,End:1}, // should ignore
  }
  var first=&Span[int,string]{Begin: 0,End:1};
  var span=driver.NextSpan(first,src);
  if(span==nil) {
    t.Errorf("Should not have reached our end yet!")
    return;
  }
  if(span.Begin!=2 || span.End!=2) {
    t.Errorf("Invalid range, expected: 2->2, got %d->%d",span.Begin,span.End);
  }
  span=driver.NextSpan(span,src);
  if(span==nil) {
    t.Errorf("Should not have reached our end yet!")
    return;
  }
  if(span.Begin!=3 || span.End!=4) {
    t.Errorf("Invalid range, expected: 3->4, got %d->%d",span.Begin,span.End);
  }
  span=driver.NextSpan(span,src);
  if(span!=nil) {
    t.Errorf("End expected!")
    return;
  }
}
