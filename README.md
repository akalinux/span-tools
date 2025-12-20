# Span-Tools

Implements the universal span intersection algorithm. 
The algorithm represents a unified way to find intersections 
and overlaps of "one dimensional spans" of any data type.

The package is built around the SpanUtil[E] struct,
The struct requires 2 methods be set in order to implement the algorithm:
	- A "Compare" function see: [cmp.Compare] for more details.
	- A "Next" function, takes a given value and get a next value.
	 The next value must be greater than the input value

The algorithm is primarily implemented by 3 methods of the SpanUtil[E] struct:
- FirstSpan, finds the initial data span intersection.
- NextSpan, finds all subsequent data span intersections.
- CreateOverlapSpan, finds the most common intersection of all overlapping spans.

Other features of this package provide ways to consolidate overlaps and data set
iteration from various data sources.


## Basic Example

In this example we will find the intersections of 3 sets of integers.

Example Sets:

		(1,2)
		(2,7)
		(5,11)
		
Example Code:

	package main
	import (
		"github.com/akalinux/span-tools"
		"fmt"
		"cmp"
	)
		
	func main() {
		var u=st.NewSpanUtil(
			// use the standard Compare function
			cmp.Compare,
			// Define our Next function
			func(e int) int { return e+1},
		)
		var list=&[]st.SpanBoundry[int]{
			u.Ns(1,2),
			u.Ns(2,7),
			u.Ns(5,11),
		}
		var count=0
		var span,ok=u.FirstSpan(list)
		for ok {
			var sources=u.GetOverlapIndexes(span,list)
			fmt.Printf("Overlap Set: %d, Span: %v, Columns: %v\n",count,span,sources)
			count++
			span,ok=u.NextSpan(span,list)
		}
	}
    