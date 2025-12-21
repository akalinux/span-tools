package st

// Represents an intersection of data.  Intersections can represent data that overlaps from
// a given source, or intersections between multiple sources.
//
// Representation as a SpanBoundry is as follows:
//  - GetBegin() returns the smallest most intersection boundry.
//  - GetEnd()   returns the largest most intersection boundry.
type ColumnOverlap[E any] interface {
	SpanBoundry[E]
	// Returns the first index point from the soruce data set
	GetSrcId() int
	// Returns the last index point from the soruce data set
	GetEndId() int
	GetOverlaps() *[]*OverlappingSpanSets[E]
	// Returns both the index point of the first data soruce, and the original SpanBoundry.
	GetFirstSpan() (int,SpanBoundry[E])
	// Returns both the index point of the last data soruce, and the original SpanBoundry.
	GetLastSpan() (int,SpanBoundry[E])
	// Returns all intersecting SpanBoundry instances and their indexes
	GetSources() *[]*OvelapSources[E]
}