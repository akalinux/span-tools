package main

import (
	"cmp"
	"fmt"
	"github.com/akalinux/span-tools"
)

func main() {
	var u = st.NewSpanUtil(
		// use the standard Compare function
		cmp.Compare,
		// Define our Next function
		func(e int) int { return e + 1 },
	)
	var list = &[]st.SpanBoundry[int]{
		u.Ns(1, 2),
		u.Ns(2, 7),
		u.Ns(5, 11),
	}

	// Create our initial span
	var span, ok = u.FirstSpan(list)

	// Denote which set we are on
	var count = 0

	for ok {
		// Find the indexes of our input set
		var sources = u.GetOverlapIndexes(span, list)

		fmt.Printf("Overlap Set: %d, Span: %v, Columns: %v\n", count, span, sources)

		count++
		span, ok = u.NextSpan(span, list)
	}

}
