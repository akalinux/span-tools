package main
import (
	"github.com/akalinux/span-tools"
	"fmt"
	"cmp"
)

func main() {
	var u=st.NewSpanUtil(
		// use the standard Compare function
		cmp.Compare,
		// Define our Next function
		func(e int) int { return e+1},
	 )
	var list=&[]st.SpanBoundry[int]{
		u.Ns(1,2),
		u.Ns(2,7),
		u.Ns(5,11),
	}
	var count=0
	var span,ok=u.FirstSpan(list)
	for ok {
  	var sources=u.GetOverlapIndexes(span,list)
  	fmt.Printf("Overlap Set: %d, Span: %v, Columns: %v\n",count,span,sources)
		count++
	  span,ok=u.NextSpan(span,list)
	}
	
}
