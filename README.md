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

In this example we create a ColumnSets instance from a SpanUtil instance.  The instance of ColumnSets
will be used to find the intersection of 3 data sets.  To add a bit of flavor, two of the data sets will contain
overlapping values.  Please note when a source is processed as a column, the overlapping data sets are consolidated 
together.  Once all columns have been added, we iterate over the result set to see how our data intersects.
 
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
For each instance of ColumnSets, a properly scoped call to "defer ac.Close()" will require being made.

	// Build our column accumulator
	ac := u.NewColumnSets()
	
	// Always make sure a defer to close is scoped correctly!
	defer ac.Close()

__Adding each data set to our ColumnSets:__

Each data set will need to be added to the ColumnSets instance. The internals refer to each column as a source.
Every source added receives an id starting from 0 for each new column/source the id is incremented by 1.  As a
note all AddCoulumnXXX methods of ColumnSets returns the index of the column/source that was added.

Note: once the iteration over the data sets begins, it is no longer possible to add additional columns to the
ColumnSets instance.  The iterators provided by the "st" class in general are considered one shot iterators.

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
		// check if there were errors
		if ac.Err != nil {
			fmt.Printf("Error: on Column: %s, error was: %v\n",m[ac.ErrCol],ac.Err)
			return
		}
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

## Intersections and Go Routines

The previous example works, but has some significant scaling limitations.
For one, the implementation is limited to pre-loaded slices in memory and
the code doesn't really take advantage of go routines.

In the real world we would expect to be able to query a system in one
go routine ,and process the results in another while the data is streamed 
back to us.The ColumnSets instance supports the use of go routines via
a chan based iterator of OverlappingSpanSets.

This example will use the same data set but simulate calling multiple systems
via go routines.  We will also manually drive the accumulation of the data
by iterating through each slice one element at a time, as if it was a
result from an external system.

The full source code can be found: [here](https://github.com/akalinux/span-tools/blob/main/examples/multichan/main.go)

__Creating an instance of ColumnSets__

The creation of the ColumnSets instance remains unchanged from our previous examples.
No special changes or initialization settings are required to enable the use of chan
based communications via go routines.  The top level ColumnSets.AddColumnXX methods
are agnostic to how the data is accumulated, so you can mix and match ass many
variations as needed.

__Creation of our go routines__

The creation of our go routines is done via a declaration of a closure.
Each call to our closure will spawn an new go routine and begin pushing data to our
initialized chan instance under the hood for us.

The OlssChanStater[E] struct when instantiated by a factory interface will be initialized
with 3 important instances:
 - The context.Context used to shut down the should the iterator fail be stopped
 - The chan of OverlappingSpanSets[E] required to communicate back to the thread running the ColumnSets instance
 - The Stater instance used to manage accumulation of SpanBoundry instances and push them to our chan

Here is the "Add" closure added to the "main" function:

	var Add = func(list *[]st.SpanBoundry[int]) int {
		s := u.NewSpanOverlapAccumulator().NewOlssChanStater()

		// Please note, we must start the go routine before we add
		// the accumulator to the ColumnSets instance.  If not we
		// will run into a race condition.
		go func() {
			// Scope our cleanup code to the go routine
			defer s.Final()
	
			end := len(*list) - 1
			id := 0
			span := (*list)[id]
			for s.CanAccumulate(span) {
				id++
				if id > end {
					return
				}
				span = (*list)[id]
			}
		}()

		// Adding the st.OlssChanStater instance to the ColumnSets
		return ac.AddColumnFromNewOlssChanStater(s)
	}


__Adding our data sets__

Since the "Add" closure will manage the creation of our go routines and
provides the columnId for us, adding our data sets will look a little
different.

	// We will map our ColumnId to our Set Name
	m := make(map[int]string)

	m[Add(
		&[]st.SpanBoundry[int]{
			u.Ns(1, 2),
			u.Ns(3, 7),  // will consolidate to 3-11
			u.Ns(5, 11), // will consolidate to 3-11
		},
	)] = "SetA"

	m[Add(&[]st.SpanBoundry[int]{
		u.Ns(3, 3),
		u.Ns(5, 11),
	})] = "SetB"

	m[Add(&[]st.SpanBoundry[int]{
		u.Ns(1, 7),
		u.Ns(8, 11),
	})] = "SetC"

__Iteration through the result set__

The rest of our code, the for loop, error checking and final output statements
remain completely unchanged.

__The resulting output__

Please note, because our data set is completely unchanged, our output is also going
to be exactly the same.  The only difference is how we accumulated the data.
Because the ColumnSets instance iterator manages the context.Context instances,
the calls to s.CanAccumulate(span) will only return true if the iterator is actively
engaged in the loop.  If we call break, or the internals run into an error the 
required call to the context cancel function is made, which will cleanup our go routine.

# SpanBoundry Consolidation of Duplicates and Overlaps

In the real world data sets are often messy, out of order, and contain duplicates/overlaps.
The internals of the "st" package expect SpanBoundry instances to be provided in a specific order. 
If data is not provided in the correct order it cannot be processed correctly.

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

__Sorting of data sets__

The full source code can be found: [here](https://github.com/akalinux/span-tools/blob/main/examples/ConsolidateOverlaps/main.go).

The SpanUtil[E] struct has a "Sort" flag, when set to true( the default ), all instances of
SpanOverlapAccumulator[E] created with the factory interface u.NewSpanOverlapAccumulator() will have
the Sort flag set to true by default.  If you wish to disable sorting you can either set the sort flag
to false on the SpanUtils.Sort=false instance before creating other object or by configuring the flag
on the instance of the SpanOverlapAccululator.Sort=false.

__Creating our SpanOverlapAccumulator__

The SpanUtil[E] instance provides a factory interface for the creation
of SpanOverlapAccumulator instances, the method is u.NewSpanOverlapAccumulator().

	ac :=u.NewSpanOverlapAccumulator()

__Sorting and Consolidation__

Now we need to step through the resulting sorted and consolidated
results.  The ac.NewOlssSeq2FromSbSlice(list) 
method provides an iter.Seq2 factory interface that can be used to driver 
our for loop for us.

	// this slice will end up being sorted by the "st" internals
	unsorted :=&[]st.SpanBoundry[int]{
		// Raw       // Will be sorted to
		u.Ns(7,11),  // Row: 3
		u.Ns(20,21), // Row: 4
		u.Ns(2,11),  // Row: 1
		u.Ns(2,12),  // Row: 0
		u.Ns(5,19),  // Row: 2
	}
	
	for id,span := range ac.NewOlssSeq2FromSbSlice(unsorted) {
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

## Manual Consolidation and Error Checking

Data integrity is very important: the internals of the "st" package check for errors
by default.  Error checking can be disabled on the SpanUltil instance by setting
u.Valudate=false.  The iterators of the "st" package stop progressing if an error
is encountered.  Typically error checking is done in an instance of SpanOverlapAccumulator.
This is generally a good place to stop the iteration process.
The SpanOverlapAccumulator instance provides a method called s.Accumulate(SpanBoundry).
This method returns both the OverlappingSpanSets instance and a pointer to error.
If the error instance is not nil then the SpanOverlapAccumulator has encountered an error.

The SpanUtil[E] instance provides a check method for validating both SpanBondry instances and
validating SpanBoundry instances in sequence.  The name of the method is Check.

__Error checking Example:__

This example checks each element of a slice of SpanBoundry instances to see if they are both valid
and in the correct order.

Example code can be found: [here](https://github.com/akalinux/span-tools/blob/main/examples/ErrorExample/main.go)

__To check if a SpanBoundry instance is valid:__

In this case if err is  not nil, then the span is valid.  Validity is defined as
span.GetBegin() is less than or equal to span.GetEnd().

	err :=u.Check(span,nil)
	
	if err!=nil {
	  // invalid span
	}	

__To check if the next SpanBoundry should be after the current SpanBoundry:__

This method performs 2 checks
 - First next is checked for validity
 - Checks if next comes after current or is equal to current

Note: current is not checked for validity.

	err :=u.Check(next,current)
	
	if err!=nil {
	  // next is out of order in relation to current
	}	

__Manual Consolidation with Error checking enabled:__

As noted, error checking is enabled by default. In this example we will iterate 
through the SpanBoundry slice twice.  In the first pass we will provide an unsorted
list that will error out during the consolidation process.  The 2nd pass we will
first sort our list and then enter the consolidation process.

The source code for this example can be found: [here](https://github.com/akalinux/span-tools/blob/main/examples/ManualConsolidation/main.go).

__First we need to turn validation on__

We will be using the same data set as our previous example, the main differences come in 3 parts.
  - The import of the "slices" package for sorting
  - Turning validation on
  - The introduction of an additional function into the "main" package

__Our updated imports for this example:__

	import (
		"cmp"
		"fmt"
		"github.com/akalinux/span-tools"
		"slices"
	)

In our main package we define a function called AccumulateSet, and it handles processing each
manual accumulation pass. Please see the source code 
[here](https://github.com/akalinux/span-tools/blob/main/examples/ManualConsolidation/main.go) 
for more details

__Example1, the expected error pass:__

	// This pass will error out
	fmt.Print("Processing our data with an invalid order\n")
	AccumulateSet(unsorted)

__Output from this section:__

Notice that we run into an error when we get to SpanBoundry (2,11).  The error is caused
by the detection of a sequencing inconsistency.

	Processing our data with an invalid order
	  &{7 11} has spawned an new OverlappingSpanSets: (7,11)
	  &{20 21} has spawned an new OverlappingSpanSets: (20,21)
	  Failed to accumulate: &{2 11}, error was: SpanBoundry out of sequence

__Example 2, the expected success pass:__

	// Once the data is sorted consolidation will work correctly
	slices.SortFunc(*unsorted, u.Compare)
	fmt.Print("\nProcessing post sort\n")
	AccumulateSet(unsorted)

__Output from this section:__

In the output data set, take note that the OverlappingSpanSets expands from (2,12)
to encompass (2,19) and a new OverlappingSpanSets is only created when a non overlapping
SpanBoundry is introduced to the Accumulator method.

	Processing post sort
	  &{2 12} has spawned an new OverlappingSpanSets: (2,12)
	  &{2 11} has been absorbed into OverlappingSpanSets: (2,12)
	  &{5 19} has been absorbed into OverlappingSpanSets: (2,19)
	  &{7 11} has been absorbed into OverlappingSpanSets: (2,19)
	  &{20 21} has spawned an new OverlappingSpanSets: (20,21)

# More Examples

For more examples see the Examples folder [examples](https://github.com/akalinux/span-tools/tree/main/examples)
