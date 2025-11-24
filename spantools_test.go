package spantools

import (
	"slices"
	"testing"
)

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
		{
			// sorted
			{Begin: -2, End: 2},
			{Begin: -1, End: 0},
			{Begin: -1, End: 0},
			{Begin: 0, End: 1},
			{Begin: 0, End: 1},
		},
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
    {
			{Begin: -1, End: 0},
			{Begin: 2, End: 2},
			{Begin: 2, End: 2},
			{Begin: 5, End: 6},
			{Begin: 9, End: 11},
    },
  },
}

type cdTests[E any, T any] struct {
  ResolvedSpanSet[T,E]
  src []Span[E,T]
} 

func TestOneContainerForAllSort(t *testing.T) {
	var driver = OrderedCreateCompare[int, string]()
	for setId, testSet := range testSets {
		var unsorted = make([]Span[int,string],len(testSet[0]))
    copy(unsorted,testSet[0]);
		slices.SortFunc(unsorted, driver.Compare)

		var sorted = testSet[1]
		for idx, span := range unsorted {
			var expected = sorted[idx]
			if span.Begin != expected.Begin || span.End != expected.End {
				t.Errorf("Error comparing test sort data set: %d, row: %d, Expected: %d,%d, Got: %d,%d", setId, idx,expected.Begin,expected.End,span.Begin,span.End)
			}
		}
	}

}
