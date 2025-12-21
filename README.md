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

The full example can be found: [here](https://github.com/akalinux/span-tools/blob/main/examples/example01/example01.go).

__Setup the package and imports:__

We will need to import our "st" package along with the "fmt" and "cmp" packages in order to process
the example data sets.

	import (
		"github.com/akalinux/span-tools"
		"fmt"
		"cmp"
	)

__Create our SpanUtil[E] instance:__

We will use the factory interface NewSpanUtil to generate our SpanUtil[int] instance for these examples.

	var u=st.NewSpanUtil(
		// use the standard Compare function
		cmp.Compare,
		// Define our Next function
		func(e int) int { return e+1},
	)

__Find our the initial SpanBoundry intersection:__

We need to find the initial intersection, before we can iterate through of these data sets.
The initial SpanBoundry is found by making a call to u.FirstSapn(list).

	// Create our initial span 
	var span,ok=u.FirstSpan(list)
		
	// Denote our overlap set position
	var count=0

__Iterate through all of our SpanBoundry intersections:__

We can now step through each data intersection point and output the results.
Each subsequent intersection is found by making a call to u.NextSpan(span,list).

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

__Resulting output:__

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

The full source code can be found: [here](https://github.com/akalinux/span-tools/blob/main/examples/beyondbasics/multicolumns.go)

__Create a ColumnSets[E] instance:__

The ColumnSets instance is created by a factory interface of SpanUtil.
For each instance of ColumnSets, a properly scoped call to "defer i.Close()" will require being made.

	// Build our column accumulator
	ac := u.NewColumnSets()
	
	// Always make sure a defer to close is scoped correctly!
	defer ac.Close()

__Adding each data set to our ColumnSets:__

Each data set will need to be added to the ColumnSets instance. 
The internals refer to each column as a source.
Every source added receives an id starting from 0, so we know in advance
what the id of each source is, but all AddCoulumnXXX methods of ColumnSets returns the index
of the column/source added.

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

__Iterate through the results:__

Finally we want to iterate through the resulting overlaps and intersections found in our different data sets:

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

__The resulting output:__

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

# SpanBoundry Consolidation of Duplicate(s) and Overlap(s)

In the real world data sets are often messy, out of order and contain duplicates and overlaps.
The internals of the "st" package expect SpanBoundry instance to be provided in a specific order. If data is not provided in the correct order it cannot
be processed correctly.

The expected order is as follows (i=SpanBoundry):
 - i.GetBegin() ascending order
 - i.GetEnd() in descending order
 
 This is an example unordered data set
 
	(7,11),
	(20,21),
	(2,11),
	(2,12),
	(5,19),

This is the same data ordered for consumption by the "st" package:

	(2,12),
	(2,11),
	(5,19),
	(7,11),
	(20,21),

__Enable Sorting of data sets__


The full source code can be found: [here](https://github.com/akalinux/span-tools/blob/main/examples/ConsolidateOverlaps/main.go).

The SpanUtil[E] struct has a "Sort" (default false ) flag, when set to true, all instances of
SpanOverlapAccumulator[E] created with the u.

	// Turn sorting on
	u.Sort=true

__Creating our SpanOverlapAccumulator__

The SpanUtil[E] instance provides a factory interface for the creation
of SpanOverlapAccumulator instances, the method is u.NewSpanOverlapAccumulator().

	ac :=u.NewSpanOverlapAccumulator()

__Sorting and Consolidation__

Now we need to step through the resulting sorting and consolidation
results.  The ac.SliceIterFactory(*list) method provides an iter.Seq2
factory interface that can be used to driver our for loop for us.

	// this slice will end up being sorted by the "st" internals
	unsorted :=&[]st.SpanBoundry[int]{
		// Raw       // Will be sorted to
		u.Ns(7,11),  // Row: 3
		u.Ns(20,21), // Row: 4
		u.Ns(2,11),  // Row: 1
		u.Ns(2,12),  // Row: 0
		u.Ns(5,19),  // Row: 2
	}
	
	for id,span := range ac.SliceIterFactory(unsorted) {
		fmt.Printf("OverlappingSpanSets: %d SpanBoundry (%d,%d)\n ",id,span.GetBegin(),span.GetEnd())
		fmt.Print(" Original Span values:\n")
		for _,src :=range *span.GetSources() {
			fmt.Printf("    Row: %d span: %v\n",src.SrcId,src.SpanBoundry)
		}
	}

__Resulting output:__

	OverlappingSpanSets: 0 SpanBoundry (2,19)
	  Original Span values:
	    Row: 0 span: &{2 12}
	    Row: 1 span: &{2 11}
	    Row: 2 span: &{5 19}
	    Row: 3 span: &{7 11}
	OverlappingSpanSets: 1 SpanBoundry (20,21)
	  Original Span values:
	    Row: 4 span: &{20 21}

# More Examples

For more examples see the Examples folder [examples](https://github.com/akalinux/span-tools/tree/main/examples)
