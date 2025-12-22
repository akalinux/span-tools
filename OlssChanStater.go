package st

import (
	"context"
)

// Context aware goroutine channel accumulator instance.
// This struct is meant to be instantiated factory interfaces.
type OlssChanStater[E any] struct {
	Chan       chan *OverlappingSpanSets[E]
	Closed     bool
	Stater     SpanIterSeq2Stater[E]
	Ctx        context.Context
	Cancle     func()
	IsShutDown bool
}

// Acts as a control method in for loops, handles pushing data to the channel from
// within a goroutine.
//
// Example:
//
//  defer s.Final()
//  for s.CanAccumulate(span) {
//    // get your next SpanBoundry
//  }
func (s *OlssChanStater[E]) CanAccumulate(span SpanBoundry[E]) bool {
	if s.Closed {
		return false
	}
	if !s.Stater.SetNext(span) {
		return true
	}
	_, next := s.Stater.GetNext()
	return s.Push(next)
}

// Attempts to push the next value to the channel, if this instance is Closed or
// if the context has been cancled, then the method returns false.
func (s *OlssChanStater[E]) Push(next *OverlappingSpanSets[E]) bool {
	if s.Closed {
		return false
	}
	select {
	case <-s.Ctx.Done():
		if !s.Closed {
			s.Closed=true
			close(s.Chan)
		}
		return false
	case s.Chan <- next:
		return true
	}
}

// Shuts down the context from the Column accumulator thread.  Do not call this outside
// of the thread running the ColumnOverlapAccumulator instance or you will get undefined
// behavior.
func (s *OlssChanStater[E]) Shutdown()  {
	if s.IsShutDown {
		return
	}
	s.IsShutDown = true
	s.Cancle()
}

// Call this method with a derfer statement in your goroutine when you are done processing
// SpanBoundry instances.
//
// Example:
//
//  defer s.Final()
func (s *OlssChanStater[E]) Final() bool {
	if s.Closed {
		return false
	}

	res :=false
	if s.Stater.HasNext() {
		_, ol := s.Stater.GetNext()
		res= s.Push(ol)
	}

	close(s.Chan)
	s.Closed=true
	return res
}
