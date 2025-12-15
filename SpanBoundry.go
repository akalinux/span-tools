package st

// Representation of a Span/Range of values in a generic context.
// The assumption is that Begin is less than or equal to the End value.
type Span[E any, T any] struct {
	// Start of the Span.
	Begin E
	// End of the Span.
	End E
	// Pointer to data set used to identify this Span[E,T]
	Tag *T
}

// Pure pointer Span struct
type SpanRef[E any, T any] struct {
	// Start of the Span.
	Begin *E
	// End of the Span.
	End *E
	// Pointer to data set used to identify this SpanRef[E,T]
	Tag *T
}

func (s *SpanRef[E, T]) GetTag() *T {
	return s.Tag
}

func (s *SpanRef[E, T]) GetBegin() E {
	return *s.Begin
}

func (s *SpanRef[E, T]) GetEnd() E {
	return *s.End
}

type SpanBoundry[E any, T any] interface {
	// Returns the Begin value.
	GetBegin() E

	// Returns the End value.
	GetEnd() E

	// Returns the pointer to the Tag value.
	GetTag() *T
}

// Returns a pointer to the current tag value
func (s *Span[E, T]) GetTag() *T {
	return s.Tag
}

// Returns the Begin value
func (s *Span[E, T]) GetBegin() E {
	return s.Begin
}

// Returns the End value
func (s *Span[E, T]) GetEnd() E {
	return s.End
}