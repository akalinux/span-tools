package main

import (
	"github.com/akalinux/span-tools"
	"fmt"
	"cmp"
)

func main() {
	var u = st.NewSpanUtil(
		// use the standard Compare function
		cmp.Compare,
		// Define our Next function
		func(e int) int { return e + 1 },
	)
	// Turn sorting on
	u.Sort=true
	
	// Create our accumulator
	ac :=u.NewSpanOverlapAccumulator()
	
	// this slice will end up being sorted by the "st" internals
	unsorted :=&[]st.SpanBoundry[int]{
		// Raw       // Will be sorted to
		u.Ns(7,11),  // Row: 3
		u.Ns(20,21), // Row: 4
		u.Ns(2,11),  // Row: 1
		u.Ns(2,12),  // Row: 0
		u.Ns(5,19),  // Row: 2
	}
	
	for id,span := range ac.NewOverlappingSpanSetsIterSeq2FromSpanBoundrySlice(unsorted) {
		fmt.Printf("OverlappingSpanSets: %d SpanBoundry (%d,%d)\n ",id,span.GetBegin(),span.GetEnd())
		fmt.Print(" Original Span values:\n")
		for _,src :=range *span.GetSources() {
			fmt.Printf("    Row: %d span: %v\n",src.SrcId,src.SpanBoundry)
		}
	}
	
}