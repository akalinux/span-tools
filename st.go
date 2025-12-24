// Implements the universal span intersection algorithm. The algorithm represents a unified way to find intersections 
// and overlaps of "one dimensional spans" of any data type.  The package is built around the SpanUtil[E any] struct, and
// the manipulation of the SpanBoundry[E any] interface.
//
// For examples and extended documentation please see the [Project] page.
//
// The "Universal Span Intersection Algorithm" is implemented by breaking operations down into their constituent parts.
// The process of finding the overlaps in data sets is in no way constrained by the types of data.  We simply need
// a way to define our spans, compare values, and create a next value.
// 
// # The parts can be described as follows:
// 
// The SpanUtils[E any] struct requires 2 methods be passed to the constructor in order to implement the algorithm:
// - A "Compare" function see: [cmp.Compare] for more details.
// - A "Next" function, takes a given value and returns next value.
//   The next value must be greater than the input value
//
// Example of creating a new instance of SpanUtil[int]:
//
//  var u = st.NewSpanUtil(
//    //use the standard Compare function
//    cmp.Compare,
//    //Define our Next function
//    func(e int) int { return e + 1 },
//  )
// 
// The algorithm is primarily implemented by 2 methods of the SpanUtil[E any] struct:
//  - FirstSpan, finds the initial data span intersection.
//  - NextSpan, finds all subsequent data span intersections.
// 
// # Basic example
//
// Our Example Set data:
//
//  var list = &[]st.SpanBoundry[int]{
//    u.Ns(1, 2),
//    u.Ns(2, 7),
//    u.Ns(5, 11),
//  }
// 
// The act of "Finding the first span" is performed  in 3 stages:
//  1. The first stage requires finding the smallest begin and end value of all of our spans.
//  2. If a begin value in our sets is both greater than the smallest begin value and less than or equal to smallest end value,
//    then the initial end value must set to the smallest begin value, else we use the smallest end value.
//  3. We will use the begin value from stage 1 and the end value from stage 2 as our "first span"
// 
// Example code snipped:
//
//   var span, ok = u.FirstSpan(list)
//
// The act of "Finding the next span" is performed in 4 stages:
//  1. The first stage we create a "new next value" that is greater than our last span end value. 
//   This "new next value" will be used as the "begin" value for step 3.
//  2. We look for the next smallest begin or end value in our sets that are, greater than or equal to our "new next value".
//   The value from this process will be used as our next "end" value for step 3.
//  3. For each set that overlaps with the span defined by the "begin" from step 1 and the "end" from step 2:
//   We need to look for the largest begin value and the smallest end value.
//  4. We will used the largest begin value and smallest end value as our "next span"
//  
// We can repeat the "Finding the next set" until step 1 yields a value beyond any end value in our sets:
//
//  // Denote which set we are on
//  var count = 0
//  
//  for ok {
//    // Find the indexes of our input set
//    var sources = u.GetOverlapIndexes(span, list)
//    fmt.Printf("Overlap Set: %d, Span: %v, Columns: %v\n", count, span, sources)
//    count++
//    span, ok = u.NextSpan(span, list)
//  }
// 
// # Beyond the basics
// 
// Finding overlaps between lists of lists takes a bit more work, but is greatly simplified by this package.
// In this example we will create an instance of st.ColumnSets[int].  The ColumnSets instance is created by a factory interface of SpanUtil.
// 
// Build our column accumulator:
//
//  ac := u.NewColumnSets()
//  // Always make sure a defer to close is scoped correctly!
//  defer ac.Close()
//
// In order to compare a list of list of spans we will add the following constraints:
//  1. Each list of list of spans must be presented in a specific order. 
//   The order is defined as: begin value in ascending order, end value in descending order.
//  2. Each time a overlapping value is encountered a new Larger span consisting of the smallest begin value and largest
//   end value must be created.  As a side effect of this the original spans that caused this overlapping set should be
//   retained to explain where this new larger span came from. 
//  3. To find the end of an overlapping set we must continue until we find the next span that does not overlap with our current
//    span or we run out of spans to process.
// 
// Example data set:
//
//  // We will map our ColumnId to our Set Name
//  m := make(map[int]string)
//
//  var seta = &[]st.SpanBoundry[int]{
//    u.Ns(1, 2),
//    u.Ns(3, 7),  // will consolidate to 3-11
//    u.Ns(5, 11), // will consolidate to 3-11
//  }
//  m[0] = "SetA"
//  
//  var setb = &[]st.SpanBoundry[int]{
//    u.Ns(3, 3),
//    u.Ns(5, 11),
//  }
//  m[1] = "SetB"
//  var setc = &[]st.SpanBoundry[int]{
//    u.Ns(1, 7),
//    u.Ns(8, 11),
//  }
//  m[2] = "SetC"
// 
//  // The internals sort slices and applys the constrains automatically
//  ac.AddColumnFromSpanSlice(seta)
//  ac.AddColumnFromSpanSlice(setb)
//  ac.AddColumnFromSpanSlice(setc)
//
// With our constraints applied, we can now consolidate list of list of spans, however we will need to add some enhancements to 
// the process.
// 
// The enhancements are as follows:
//  1. Now our our "Finding the first span" list of spans are pulled the first member of each constrained and consolidated list
//    list of spans.  We will refer the "first span" as the "current span".
//  2. For each constrained list of spans: iterate through each constrained span and save all the constrained spans that overlap with
//    our current span.  When we find a span with an end or begin value beyond the current span end value we or we have exhausted this list of
//    spans, this list iteration is complete.
//  3. Now we apply "Finding the next span" to our current list of constrained spans.  We will refer to the "next span" as the "current span"
//  4. Repeat steps 2 and 3 until we have exhausted all lists of constrained spans.
//
// Example Iterator loop:
//
//  header := "+-----+--------------------+------------------------------------+\n"
//  fmt.Print(header)
//  fmt.Print("| Seq | Begin and End      | Set Name:(Row,Row)                 |\n")
//  for pos, res := range ac.Iter() {
//  // check if there were errors
//  if ac.Err != nil {
//    fmt.Printf("Error: on Column: %s, error was: %v\n",m[ac.ErrCol],ac.Err)
//    return
//  }
//  cols := res.GetColumns()
//  names := []string{}
//  for _, column := range *cols {
//    str :=fmt.Sprintf("%s:(%d-%d)",m[column.ColumnId],column.GetSrcId(),column.GetEndId())
//    names = append(names, str)
//    }
//    fmt.Print(header)
//    fmt.Printf("| %- 3d | Begin:% 3d, End:% 3d | %- 34s |\n",
//      pos,
//      res.GetBegin(),
//      res.GetEnd(),
//      strings.Join(names, ", "),
//    )
//  }
//  fmt.Print(header)
//
// # Integrating go routines and streaming data sets
//
// The internals of the st package, can be used to create context instances to mange communication between go routines for us.
// If the spans are recived out of order or fail to pass error checking constraints, then the main iterator loop will be halted.
// In this example we will simulate streaming the same data set via go routines.
//
// First we create an "Add" function to create our go routien and push data into our column channel based iterators:
//
//  var Add = func(list *[]st.SpanBoundry[int]) int {
//    s := u.NewSpanOverlapAccumulator().NewOlssChanStater()
//  
//    // Please note, we must start the go routine before we add
//    // the accumulator to the ColumnSets instance.  If not we
//    // will run into a race condition.
//    go func() {
//      // Scope our cleanup code to the go routine
//      defer s.Final()
//  
//      end := len(*list) - 1
//      id := 0
//      span := (*list)[id]
//      for s.CanAccumulate(span) {
//        id++
//        if id > end {
//          return
//        }
//        span = (*list)[id]
//      }
//      }()
//  
//    // Adding the st.OlssChanStater instance to the ColumnSets
//    return ac.AddColumnFromNewOlssChanStater(s)
//  }
//
// We can now add our data sets:
//
//  Add(&[]st.SpanBoundry[int]{
//    u.Ns(1, 2),
//    u.Ns(3, 7),  // will consolidate to 3-11
//    u.Ns(5, 11), // will consolidate to 3-11
//  })
//  Add(&[]st.SpanBoundry[int]{
//    u.Ns(3, 3),
//    u.Ns(5, 11),
//  })
//  Add(&[]st.SpanBoundry[int]{
//    u.Ns(1, 7),
//    u.Ns(8, 11),
//  })
//
// The for loop and map remain unchanged from our previous example.  The only differnce is the internals have no way 
// to sort the data before it is consolidated.
//
// [Project]: https://github.com/akalinux/span-tools
// [cmp.Compare]: https://pkg.go.dev/cmp#Compare
package st
