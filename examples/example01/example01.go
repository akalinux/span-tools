package main

import (
	"cmp"
	"fmt"
	"github.com/akalinux/span-tools"
)

var u = st.NewSpanUtil(
	// use the standard Compare function
	cmp.Compare,
	// Define our Next function
	func(e int) int { return e + 1 },
)

type MySpan struct {
	a int
	b int 
}

func (s *MySpan) GetBegin() int {
	return s.a
}
func (s *MySpan) GetEnd() int {
	return s.b
}

func init() {
	// overload the default span
	u.SpanFactory=func (a,b int) st.SpanBoundry[int] {
		return &MySpan{a,b}
	}
}

func main() {

	var list = &[]st.SpanBoundry[int]{
		u.Ns(1, 2),
		u.Ns(2, 7),
		u.Ns(5, 11),
	}

	// Create our initial span
	var span, ok = u.FirstSpan(list)

	// Denote which set we are on
	var count = 0

	for ok {
		// Find the indexes of our input set
		var sources = u.GetOverlapIndexes(span, list)

		fmt.Printf("Overlap Set: %d, Span: %v, Columns: %v\n", count, span, sources)

		count++
		span, ok = u.NextSpan(span, list)
	}

}
