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

	// this slice will end up being sorted by the "st" internals
	unsorted := &[]st.SpanBoundry[int]{
		// Raw       // Will be sorted to
		u.Ns(7, 11),  // Row: 3
		u.Ns(20, 21), // Row: 4
		u.Ns(2, 11),  // Row: 1
		u.Ns(2, 12),  // Row: 0
		u.Ns(5, 19),  // Row: 2
	}

	// represents our current value
	var current st.SpanBoundry[int]
	for id, next := range *unsorted {
		err := u.Check(next, current)
		if err == nil {
			fmt.Printf("id: %d SpanBoundry: %v, OK\n", id, next)
			current=next
		} else {
			fmt.Printf("id: %d SpanBoundry: %v, Not Ok, error was: %v\n", id, next,err)
			break
		}
	}

}
