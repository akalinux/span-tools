package st

// Representation of a Span/Range of values in a generic context.
// The assumption is that Begin is less than or equal to the End value.
type Span[E any] struct {
	// Start of the Span.
	Begin E
	// End of the Span.
	End E
}

// Pure pointer Span struct
type SpanRef[E any] struct {
	// Start of the Span.
	Begin *E
	// End of the Span.
	End *E
}

func (s *SpanRef[E]) GetBegin() E {
	return *s.Begin
}

func (s *SpanRef[E]) GetEnd() E {
	return *s.End
}

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

