# Span-Tools

Implements the universal span intersection algorithm. 
The algorithm represents a unified way to find intersections 
and overlaps of "one dimensional spans" of any data type.

The package is built around the SpanUtil[E] struct,
The struct requires 2 methods be set in order to implement the algorithm:
- A "Compare" function see: [cmp.Compare](https://pkg.go.dev/cmp#Compare) for more details.
- A "Next" function, takes a given value and returns next value.
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
		
		// Create our initial span 
		var span,ok=u.FirstSpan(list)
		
		// Denote our overlap set position
		var count=0
		for ok {
		  // Get the indexes of the columns this overlap relates to
			var sources=u.GetOverlapIndexes(span,list)
			
			// output our intersection data
			fmt.Printf("Overlap Set: %d, Span: %v, Columns: %v\n",count,span,sources)
			
			// update our overlap set
			count++
			
			// get our next set
			span,ok=u.NextSpan(span,list)
		}
	}

Resulting output:

    Overlap Set: 0, Span: &{1 1}, Columns: &[0]
    Overlap Set: 1, Span: &{2 2}, Columns: &[0 1]
    Overlap Set: 2, Span: &{3 5}, Columns: &[1 2]
    Overlap Set: 3, Span: &{6 7}, Columns: &[1 2]
    Overlap Set: 4, Span: &{8 11}, Columns: &[2]

# More Examples

For more examples see the Examples folder [examples](./examples)
    
