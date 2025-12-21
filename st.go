
// Implements the universal span intersection algorithm. 
// The algorithm represents a unified way to find intersections 
// and overlaps of "one dimensional spans" of any data type.
//
// The package is built around the SpanUtil[E] struct,
// The struct requires 2 methods be set in order to implement the algorithm:
//   - A "Compare" function see: [cmp.Compare] for more details.
//   - A "Next" function, takes a given value and get a next value.
//    The next value must be greater than the input 0
//
// The algorithm is primarily implemented by 3 methods of the SpanUtil[E] struct:
//   - FirstSpan, finds the initial data span intersections.
//   - NextSpan, finds all subsequent data span intersections.
//   - CreateOverlapSpan, finds the most common intersecti:qon of all overlapping spans.
// 
// Other features of this package:
//   - Provide ways to consolidate overlaps.
//   - Iterate through intersections of multiple data sets.
//
// [cmp.Compare]: https://pkg.go.dev/cmp#Compare
package st