package main

import (
	"cmp"
	"fmt"
	"github.com/akalinux/span-tools"
	"slices"
)

var u = st.NewSpanUtil(
	// use the standard Compare function
	cmp.Compare,
	// Define our Next function
	func(e int) int { return e + 1 },
)
func main() {
	
	// turn validation on
	u.Validate=true
	// this slice will end up being sorted by the "st" internals
	unsorted := &[]st.SpanBoundry[int]{
		// Raw       // Will be sorted to
		u.Ns(7, 11),  // Row: 3
		u.Ns(20, 21), // Row: 4
		u.Ns(2, 11),  // Row: 1
		u.Ns(2, 12),  // Row: 0
		u.Ns(5, 19),  // Row: 2
	}

	// This pass will error out
	fmt.Print("Processing our data with an invalid order\n")
	AccumulateSet(unsorted)

	
	// Once the data is sorted consolidation will work correctly
	slices.SortFunc(*unsorted, u.Compare)
	fmt.Print("\nProcessing post sort\n")
	AccumulateSet(unsorted)

}

func AccumulateSet(list *[]st.SpanBoundry[int]) {

	// Create our accumulator
	ac := u.NewSpanOverlapAccumulator()

	id := 1
	span := (*list)[0]
	ol, err := ac.Accumulate(span)
	max := len(*list) 

	var last *st.OverlappingSpanSets[int]
	for ; id < max; id++ {
		if err != nil {
			fmt.Printf("  Failed to accumulate: %v, error was: %v\n", span, err)
			return
		} else if last == ol {
			fmt.Printf("  %v has beeen absorbed int OverlappingSpanSets: (%d,%d)\n", span, ol.GetBegin(), ol.GetEnd())
		} else {
			fmt.Printf("  %v has spawned an new OverlappingSpanSets: (%d,%d)\n", span, ol.GetBegin(), ol.GetEnd())
		}
		last = ol
		span = (*list)[id]
		ol, err = ac.Accumulate(span)
	}
	if err != nil {
		fmt.Printf("  Failed to accumulate: %v, error was: %v\n", span, err)
	} else if last == ol {
		fmt.Printf("  %v has been absorbed into OverlappingSpanSets: (%d,%d)\n", span, ol.GetBegin(), ol.GetEnd())
	} else {
		fmt.Printf("  %v has spawned an new OverlappingSpanSets: (%d,%d)\n", span, ol.GetBegin(), ol.GetEnd())
	}
}
