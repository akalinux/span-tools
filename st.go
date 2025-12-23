
// Implements the universal span intersection algorithm. The algorithm represents a unified way to find intersections 
// and overlaps of "one dimensional spans" of any data type.  The package is built around the SpanUtil[E any] struct, and
// the manipulation of the SpanBoundry[E any] interface.
//
// For examples documentation please see the [Tutorials] page.
// 
// The SpanUtils[E any] struct requires 2 methods be passed to the constructor in order to implement the algorithm:
// - A "Compare" function see: [cmp.Compare] for more details.
// - A "Next" function, takes a given value and returns next value.
//   The next value must be greater than the input value
// 
// The algorithm is primarily implemented by 3 methods of the SpanUtil[E any] struct:
//  - FirstSpan, finds the initial data span intersection.
//  - NextSpan, finds all subsequent data span intersections.
//  - CreateOverlapSpan, finds the most common intersection of all overlapping spans.
// 
// Other features of this package:
//  - Provide ways to consolidate overlaps.
//  - Iterate through intersections of multiple data sets.
//
// [Tutorials]: https://github.com/akalinux/span-tools
// [cmp.Compare]: https://pkg.go.dev/cmp#Compare
package st