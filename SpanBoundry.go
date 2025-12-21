package st

// Representation of a Span/Range of values in a generic context.
// The assumption is that Begin is less than or equal to the End value.
type Span[E any] struct {
	// Start of the Span.
	Begin E
	// End of the Span.
	End E
}


// This interface acts as the core representation of spans for the "st" pacakge.
// Spans are represent by 2 values a "Begin" value and an "End" value.
// The Begin value should be returned by GetBegin and should be greater than
// or equal to the End value returned by GetEnd.
type SpanBoundry[E any] interface {
	// Returns the Begin value.
	GetBegin() E

	// Returns the End value.
	GetEnd() E
}

// Returns the Begin value
func (s *Span[E]) GetBegin() E {
	return s.Begin
}

// Returns the End value
func (s *Span[E]) GetEnd() E {
	return s.End
}

