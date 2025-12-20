
// Implements the universal span intersection algorithm. 
// The algorithm represents a unified way to find intersections 
// and overlaps of 1d ranges of any data type.
//
// The construction of a SpanUtil[E] instance with the following arguments passed to the constructor:
//   - A "Compare"" function see: [cmp.Compare] for more details.
//   - A "Next"" function a way take in a given value and get a next value.
//    The next value must be greater than the input vallue
//
// The SpanUtil[E] instnace provides methods and factory interfaces for processing setns of Spans represented by objects implementing the SpanBoundry[E] interface.
//
// [cmp.Compare]: https://pkg.go.dev/cmp#Compare
package st