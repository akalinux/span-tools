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

Other features of this package:
- Provide ways to consolidate overlaps.
- Iterate through intersections of multiple data sets.

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

# Beyond The Basics

The basic example works, but its not very useful.  In the real world
we generally have multiple data sources.  Usually we want to find the intersections
of between those different data sources.

In this example we create a ColumnSets instance from a SpanUtil instance and add each
data set as a column.  Once all columns have been added, we iterate over the result
set which contains the data and how it intersects.
 
In this example we will use 3 data sets, one of which contains overlapping values.
Please note when a source is processed as a column, the overlapping data sets are consolidated together.

Example Data sets:

	SetA:
			(1, 2),
			(3, 7),  // will consolidate to 3-11
			(5, 11), // will consolidate to 3-11
	
	SetB:
			(3, 3),
			(5, 11),
	
	SetC:
			(1, 7),  // will consolidate to 1-11
			(8, 11), // will consolidate to 1-11


Example Code:

	package main
	
	import (
		"cmp"
		"fmt"
		"strings"
		"github.com/akalinux/span-tools"
	)
	
	func main() {
		u := st.NewSpanUtil(
			// use the standard Compare function
			cmp.Compare,
				
			// Define our Next function
			func(e int) int { return e + 1 },
		)
		// Build our column accumulator
		ac := u.NewColumnSets()
		
		// Always make sure a defer to Close is scoped correctly!
		defer ac.Close()
		
		// We will map our ColumnId to our Set Name
		m := make(map[int]string)
		
		var seta = &[]st.SpanBoundry[int]{
			u.Ns(1, 2),
			u.Ns(3, 7),  // will consolidate to 3-11
			u.Ns(5, 11), // will consolidate to 3-11
		}
		ac.AddColumnFromSpanSlice(seta)
		m[0] = "SetA"
		
		var setb = &[]st.SpanBoundry[int]{
			u.Ns(3, 3),
			u.Ns(5, 11),
		}
		ac.AddColumnFromSpanSlice(setb)
		m[1] = "SetB"
		
		var setc = &[]st.SpanBoundry[int]{
			u.Ns(1, 7),
			u.Ns(8, 11),
		}
		ac.AddColumnFromSpanSlice(setc)
		m[2] = "SetC"
		
		header := "+-----+--------------------+------------------------------------+\n"
		fmt.Print(header)
		fmt.Print("| Seq | Begin and End      | Set Name:(Row,Row)                 |\n")
		for pos, res := range ac.Iter() {
			cols := res.GetColumns()
			names := []string{}
			for _, column := range *cols {
				str :=fmt.Sprintf("%s:(%d-%d)",m[column.ColumnId],column.GetSrcId(),column.GetEndId())
				names = append(names, str)
			}
			fmt.Print(header)
			fmt.Printf("| %- 3d | Begin:% 3d, End:% 3d | %- 34s |\n",
				pos,
				res.GetBegin(),
				res.GetEnd(),
				strings.Join(names, ", "),
			)
		}
		fmt.Print(header)	
	}

The resulting output would be:

	+-----+--------------------+------------------------------------+
	| Seq | Begin and End      | Set Name:(Row,Row)                 |
	+-----+--------------------+------------------------------------+
	|  0  | Begin:  1, End:  2 | SetA:(0-0), SetC:(0-0)             |
	+-----+--------------------+------------------------------------+
	|  1  | Begin:  3, End:  3 | SetA:(1-2), SetB:(0-0), SetC:(0-0) |
	+-----+--------------------+------------------------------------+
	|  2  | Begin:  4, End:  5 | SetA:(1-2), SetB:(1-1), SetC:(0-0) |
	+-----+--------------------+------------------------------------+
	|  3  | Begin:  6, End:  7 | SetA:(1-2), SetB:(1-1), SetC:(0-0) |
	+-----+--------------------+------------------------------------+
	|  4  | Begin:  8, End: 11 | SetA:(1-2), SetB:(1-1), SetC:(1-1) |
	+-----+--------------------+------------------------------------+

# More Examples

For more examples see the Examples folder [examples](./examples)
