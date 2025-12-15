// Package st implements utilties for finding intersection of overlaping values in related data sets.
//
// The package provides a data span/range intersection library that is algorithmically implmented using generics and should work with any
// data set that has a comparable Begin and End value.   If you can compare a begin and end value of and resolve to a -1,0,1, then this library is
// able to find how that data intersects.
//
// # How this package works
//
// Spans in this package are expected to contain a Begin and End value. The Begin and End values should be
// comparable with a cmp function, see the go standard [cmp.Compare] function for more details.  The only constraint on the data sets are that the Begin value is required to be less than or equal to the End value.
//
// [cmp.Compare]: https://pkg.go.dev/cmp#Compare
package st