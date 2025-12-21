package main

import (
	"github.com/akalinux/span-tools"
	"fmt"
	"cmp"
	"slices"
)

func main() {
	var u = st.NewSpanUtil(
		// use the standard Compare function
		cmp.Compare,
		// Define our Next function
		func(e int) int { return e + 1 },
	)
	u.Sort=true
	
	unsorted :=&[]st.SpanBoundry[int]{
		u.Ns(7,11),
		u.Ns(20,21),
		u.Ns(2,11),
		u.Ns(2,12),
		u.Ns(5,19),
	}
	
}