package st

import (
	"iter"
	"slices"
	"testing"
)

var AllSet = []SpanBoundry[int, string]{
	// sorted
	&Span[int, string]{Begin: -2, End: 2},
	&Span[int, string]{Begin: -1, End: 0},
	&Span[int, string]{Begin: -1, End: 0},
	&Span[int, string]{Begin: 0, End: 1},
	&Span[int, string]{Begin: 0, End: 1},
}

var tagA = "a"
var tagB = "b"
var MultiSet = []SpanBoundry[int, string]{
	&Span[int, string]{Begin: -1, End: 0},            // 0
	&Span[int, string]{Begin: 2, End: 2, Tag: &tagA}, //1
	&Span[int, string]{Begin: 2, End: 2, Tag: &tagB}, //1
	&Span[int, string]{Begin: 5, End: 6},             // 2
	&Span[int, string]{Begin: 9, End: 11},            // 3
	// -1
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
	Begin *int
	End   *int
	Tag   *string
}

func (s *SpanInt) GetTag() *string {
	return s.Tag
}

func (s *SpanInt) GetBegin() int {
	return *s.Begin
}

func (s *SpanInt) GetEnd() int {
	return *s.End
}

func TestSaneSpan(t *testing.T) {
	var begin = 2
	var end = 1
	var next = &SpanInt{Begin: &begin, End: &end}

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
  s.Validate=false;
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

// Validates the inital range of a list of ranges
func TestFirstRange(t *testing.T) {
	var src *[]SpanBoundry[int, string] = &[]SpanBoundry[int, string]{
		&Span[int, string]{Begin: 2, End: 2},
		&Span[int, string]{Begin: 0, End: 1},
	}
	var span = testDriver.FirstSpan(src)
	if span.Begin != 0 || span.End != 1 {
		t.Errorf("Invalid start range")
	}
}

// Validates the creation of the next span based on the current span.
func TestNextRange(t *testing.T) {
	var src = &[]SpanBoundry[int, string]{
		&Span[int, string]{Begin: 3, End: 4}, // 3,4, last valid range, should get nil after this
		&Span[int, string]{Begin: 2, End: 2}, // 2,2, first range
		&Span[int, string]{Begin: 0, End: 1}, // should ignore
	}
	var first = &Span[int, string]{Begin: 0, End: 1}
	var span = testDriver.NextSpan(first, src)
	if span == nil {
		t.Errorf("Should not have reached our end yet!")
		return
	}
	if span.Begin != 2 || span.End != 2 {
		t.Errorf("Invalid range, expected: 2->2, got %d->%d", span.Begin, span.End)
	}
	span = testDriver.NextSpan(span, src)
	if span == nil {
		t.Errorf("Should not have reached our end yet!")
		return
	}
	if span.Begin != 3 || span.End != 4 {
		t.Errorf("Invalid range, expected: 3->4, got %d->%d", span.Begin, span.End)
	}
	span = testDriver.NextSpan(span, src)
	if span != nil {
		t.Errorf("End expected!")
		return
	}
}

// Validates the creation of the next span when there are gaps.
func TestNextRangeGap(t *testing.T) {
	var src *[]SpanBoundry[int, string] = &[]SpanBoundry[int, string]{
		&Span[int, string]{Begin: 4, End: 5}, // 3,4, last valid range, should get nil after this
		&Span[int, string]{Begin: 2, End: 2}, // 2,2, first range
		&Span[int, string]{Begin: 0, End: 1}, // should ignore
	}
	var first = &Span[int, string]{Begin: 0, End: 0}
	var span = testDriver.NextSpan(first, src)
	if span == nil {
		t.Errorf("Should not have reached our end yet!")
		return
	}
	if span.Begin != 2 || span.End != 2 {
		t.Errorf("Invalid range, expected: 2->2, got %d->%d", span.Begin, span.End)
		return
	}
	span = testDriver.NextSpan(span, src)
	if span == nil {
		t.Errorf("Should not have reached our end yet!")
		return
	}
	if span.Begin != 4 || span.End != 5 {
		t.Errorf("Invalid range, expected: 4->5, got %d->%d", span.Begin, span.End)
		return
	}
	span = testDriver.NextSpan(span, src)
	if span != nil {
		t.Errorf("End expected!")
		return
	}
}

// Used to test overlaping span generation.
func testNextOverlaps(t *testing.T, src *[]SpanBoundry[int, string]) {

	var first = &Span[int, string]{Begin: -1, End: 0}
	var span = testDriver.NextSpan(first, src)
	if span == nil {
		t.Errorf("Should not have reached our end yet!")
		return
	}
	if span.Begin != 1 || span.End != 1 {
		t.Errorf("Invalid range, expected: 1->1, got %d->%d", span.Begin, span.End)
		return
	}
	span = testDriver.NextSpan(span, src)
	if span == nil {
		t.Errorf("Should not have reached our end yet!")
		return
	}
	if span.Begin != 2 || span.End != 3 {
		t.Errorf("Invalid range, expected: 2->3, got %d->%d", span.Begin, span.End)
	}
	span = testDriver.NextSpan(span, src)
	if span != nil {
		t.Errorf("End expected!")
		return
	}
}

// Validates the creation of the next span with overlaps.
func TestNextRangeOverlaps(t *testing.T) {
	var src *[]SpanBoundry[int, string] = &[]SpanBoundry[int, string]{
		&Span[int, string]{Begin: 2, End: 3}, // 2,3, next range is nil
		&Span[int, string]{Begin: 1, End: 3}, // overlaps with 0 and 2
		&Span[int, string]{Begin: 0, End: 1}, // 1,1, first range
	}
	testNextOverlaps(t, src)
}

// Validates the creation of the next span with a different data set.
func TestNextRangeOverlapsReverseOrder(t *testing.T) {
	var src *[]SpanBoundry[int, string] = &[]SpanBoundry[int, string]{
		&Span[int, string]{Begin: 0, End: 1}, // 1,1, first range
		&Span[int, string]{Begin: 1, End: 3}, // overlaps with 0 and 2
		&Span[int, string]{Begin: 2, End: 3}, // 2,3, next range is nil
	}
	testNextOverlaps(t, src)
}

func TestNextRangeOverlapsMixedOrder(t *testing.T) {
	var src *[]SpanBoundry[int, string] = &[]SpanBoundry[int, string]{
		&Span[int, string]{Begin: 1, End: 3}, // overlaps with 1 and 2
		&Span[int, string]{Begin: 0, End: 1}, // 1,1, first range
		&Span[int, string]{Begin: 2, End: 3}, // 2,3, next range is nil
	}
	testNextOverlaps(t, src)
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

func TestAccumulateIter(t *testing.T) {
	var id = -1
	var sa = testDriver.NewSpanOverlapAccumulator()
	sa.Sort = true

	for idx, res := range sa.SliceIterFactory(&MultiSet) {
		id++
		if id != idx {
			t.Errorf("Failed, expected idx: %d, got %d", id, idx)
			return
		}
		switch idx {
		case 0:
			{
				if res.GetBegin() != -1 || res.GetEnd() != 0 {
					t.Errorf("Invalid Range on set 0, expected: -1,0, got %d,%d", res.GetBegin(), res.GetEnd())
					return
				}
				if res.Contains != nil {
					t.Errorf("Exepcted Empty contains")
					return
				}
				if id != idx {
					t.Errorf("Expected to be at idx: 0, but we are at idx %d", idx)
				}
			}
		case 1:
			{
				if res.GetBegin() != 2 || res.GetEnd() != 2 {
					t.Errorf("Invalid Range on set 1, expected: 2,2, got %d,%d", res.GetBegin(), res.GetEnd())
					return
				}
				if res.GetContains() == nil {
					t.Errorf("Exepcted Non-Empty contains")
					return
				}
			}
		case 2:
			{
				if res.GetBegin() != 5 || res.GetEnd() != 6 {
					t.Errorf("Invalid Range on set 0, expected: 5,6... got %d,%d", res.GetBegin(), res.GetEnd())
					return
				}
				if res.Contains != nil {
					t.Errorf("Exepcted Empty contains")
					return
				}
			}
		case 3:
			{
				if res.GetBegin() != 9 || res.GetEnd() != 11 {
					t.Errorf("Invalid Range on set 0, expected: 9,11... got %d,%d", res.GetBegin(), res.GetEnd())
					return
				}
				if res.Contains != nil {
					t.Errorf("Exepcted Empty contains")
					return
				}
			}
		default:
			{
				t.Errorf("Got a range beyond 3, expected set to end at the, end is at: %d", idx)
				return
			}
		}

	}
	if id != 3 {
		t.Errorf("Expected final id to be 3, got %d", id)
		return
	}
	for idx := range testDriver.NewSpanOverlapAccumulator().SliceIterFactory(&MultiSet) {
		if idx == 0 {
			break
		}
		if idx != 0 {
			t.Errorf("Someting went wrong and our iterator broke?")
			return
		}
	}
	var count = 0
	for range testDriver.NewSpanOverlapAccumulator().SliceIterFactory(nil) {
		count++
	}
	if count != 0 {
		t.Errorf("Should have not gotten any iterator passes when our slice is nil")
		return
	}
}

func TestAccumulateIterPull2(t *testing.T) {

	var next, stop = iter.Pull2(testDriver.NewSpanOverlapAccumulator().SliceIterFactory(&MultiSet))
	defer stop()
	var idx, res, hasnext = next()
	var count = -1
	for hasnext {
		count++
		if idx != count {
			t.Errorf("Failed, expected %d, got %d", count, idx)
			return
		}
		switch idx {
		case 0:
			{
				if res.GetBegin() != -1 || res.GetEnd() != 0 {
					t.Errorf("Invalid Range on set 0, expected: -1,0, got %d,%d", res.GetBegin(), res.GetEnd())
					return
				}
				if res.Contains != nil {
					t.Errorf("Exepcted Empty contains")
					return
				}
			}
		case 1:
			{
				if res.GetBegin() != 2 || res.GetEnd() != 2 {
					t.Errorf("Invalid Range on set 1, expected: 2,2, got %d,%d", res.GetBegin(), res.GetEnd())
					return
				}
				if res.GetContains() == nil {
					t.Errorf("Exepcted Non-Empty contains")
					return
				}
			}
		case 2:
			{
				if res.GetBegin() != 5 || res.GetEnd() != 6 {
					t.Errorf("Invalid Range on set 0, expected: 5,6... got %d,%d", res.GetBegin(), res.GetEnd())
					return
				}
				if res.Contains != nil {
					t.Errorf("Exepcted Empty contains")
					return
				}
			}
		case 3:
			{
				if res.GetBegin() != 9 || res.GetEnd() != 11 {
					t.Errorf("Invalid Range on set 0, expected: 9,11... got %d,%d", res.GetBegin(), res.GetEnd())
					return
				}
				if res.Contains != nil {
					t.Errorf("Exepcted Empty contains")
					return
				}
			}
		default:
			{
				t.Errorf("Got a range beyond 3, expected set to end at the, end is at: %d", idx)
				return
			}
		}
		idx, res, hasnext = next()
	}
	if count != 3 {
		t.Errorf("Expected 3 rows, got %d", count)
	}
}
func TestColumnConsolidateIter(t *testing.T) {
	var res = testDriver.NewSpanOverlapAccumulator().ColumnOverlapSliceFactory(&MultiSet)
	res.GetNext(&Span[int, string]{Begin: -1, End: 0})
	defer res.Close()
	if !res.HasNext() {
		t.Errorf("Should Has Next")
		return
	}
	if len(*res.Overlaps) != 1 {
		t.Errorf("Should have 1 element in our slice")
		return
	}
	if (*res.Overlaps)[0].GetBegin() != -1 || (*res.Overlaps)[0].GetEnd() != 0 {
		t.Errorf("Invalid first element")
		return
	}
	if res.SrcPos != 1 {
		t.Errorf("Expected SrcPos: 1, got SrcPos: %d", res.SrcPos)
		return
	}
	res.GetNext(&Span[int, string]{Begin: 1, End: 1})
	if !res.HasNext() {
		t.Errorf("Expected a next")
		return
	}
	if res.SrcPos != 1 {
		t.Errorf("Expected SrcPos: 1, got SrcPos: %d", res.SrcPos)
		return
	}
	res.GetNext(&Span[int, string]{Begin: 2, End: 5})
	if !res.HasNext() {
		t.Errorf("Expected a next")
		return
	}

	if res.SrcPos != 2 {
		t.Errorf("Expected SrcPos: 3, got SrcPos: %d", res.SrcPos)
		return
	}
	if res.SrcStart != 1 {
		t.Errorf("Expected SrcStart: 1, got: %d", res.SrcStart)
		return
	}
	if res.SrcEnd != 2 {
		t.Errorf("Expected SrcEnd: 2, got: %d", res.SrcEnd)
		return
	}
	res.GetNext(&Span[int, string]{Begin: 6, End: 11})
	if res.SrcPos != 3 {
		t.Errorf("Expected SrcPos: 3, got SrcPos: %d", res.SrcPos)
		return
	}
	if !res.HasNext() {
		t.Errorf("Expected a next")
		return
	}
	if res.SrcStart != 2 {
		t.Errorf("Expected SrcStart: 2, got: %d", res.SrcStart)
		return
	}
	if res.SrcEnd != 3 {
		t.Errorf("Expected SrcEnd: 3, got: %d", res.SrcEnd)
		return
	}

	res.GetNext(&Span[int, string]{Begin: 12, End: 12})
	if res.HasNext() {
		t.Errorf("Expected to not have next")
		return
	}

	// Make sure we close our pull iter
	res.Close()
	res = testDriver.NewSpanOverlapAccumulator().ColumnOverlapSliceFactory(&MultiSet)
	res.GetNext(&Span[int, string]{Begin: -2, End: -2})
	if res.SrcPos != 0 {
		t.Errorf("Make sure our first span is 0, got %d", res.SrcPos)
		return
	}

	res.Close()
	res = testDriver.NewSpanOverlapAccumulator().ColumnOverlapSliceFactory(&MultiSet)
	res.GetNext(&Span[int, string]{Begin: 2, End: 2})

	if res.SrcStart != 1 {
		t.Errorf("Validate we got the correct start postion, expected 1, got %d", res.SrcPos)
		return
	}
	if res.SrcEnd != 1 {
		t.Errorf("Validate we got the correct end postion, expected 1, got %d", res.SrcEnd)
		return
	}
	if res.SrcPos != 2 {
		t.Errorf("Make sure our span id is 3, got %d", res.SrcPos)
		return
	}
	if res.Next.GetBegin() != 5 || res.Next.GetEnd() != 6 {
		t.Errorf("Make sure our next range is 5->6, got %d->%d", res.Next.GetBegin(), res.Next.GetEnd())
		return
	}
	res.Close()
	res = testDriver.NewSpanOverlapAccumulator().ColumnOverlapSliceFactory(&MultiSet)
	res.GetNext(&Span[int, string]{Begin: 20, End: 20})
	if res.SrcStart != -1 {
		t.Error("Should not have a next!")
		return
	}

}

func TestBadOrder(t *testing.T) {
	var list = &[]SpanBoundry[int, string]{
		&Span[int, string]{Begin: 9, End: 11},
		&Span[int, string]{Begin: 2, End: 2},
	}
  testDriver.Validate=true;
  for id,span := range testDriver.NewSpanOverlapAccumulator().SliceIterFactory(list) {
    if(id>0) {
      t.Errorf("Should stop at 0");
      return;
    }
    if(span.GetBegin()!=9) {
      t.Error("Should have span 0, got span 1")
    } 
  }
}

func TestBadInitValue(t *testing.T) {
  var list = &[]SpanBoundry[int, string]{
    &Span[int, string]{Begin: 13, End: 11},
  }
  testDriver.Validate=true;
  for  range testDriver.NewSpanOverlapAccumulator().SliceIterFactory(list) {
    t.Errorf("Should have gotten no ranges!");
    return;
  }

}
