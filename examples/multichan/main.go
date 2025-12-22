package main

import (
	"cmp"
	"context"
	"fmt"
	"github.com/akalinux/span-tools"
	"iter"
	"strings"
)

var u = st.NewSpanUtil(
	// use the standard Compare function
	cmp.Compare,

	// Define our Next function
	func(e int) int { return e + 1 },
)

func buildChanColumnAccumulator(ctx context.Context, list *[]st.SpanBoundry[int]) (*st.ColumnOverlapAccumulator[int]) {
	col := make(chan *st.OverlappingSpanSets[int])
	na :=u.NewSpanOverlapAccumulator()
	go func() {

		next, stop := iter.Pull2(
			na.NewOverlappingSpanSetsIterSeq2FromSpanBoundrySlice(
				list,
			),
		)
		defer stop()
		_, res, ok := next()
		if !ok {
			return
		}
		for {
			select {
			case <-ctx.Done():
				close(col)
				return
			case col <- res:
				_, res, ok = next()
				if !ok {
					close(col)
					return
				}
			}
		}
	}()
	return na.NewColumnOverlapAccumulatorFromOverlappingSpanSetsChan(col)
}

func main() {

	// Build our column accumulator
	ac := u.NewColumnSets()

	// Always make sure a defer to close is scoped correctly!
	defer ac.Close()
	// We will map our ColumnId to our Set Name
	m := make(map[int]string)

	var seta = &[]st.SpanBoundry[int]{
		u.Ns(1, 2),
		u.Ns(3, 7),  // will consolidate to 3-11
		u.Ns(5, 11), // will consolidate to 3-11
	}
	m[0] = "SetA"
	ctxA, cancelA := context.WithCancel(context.Background())
	defer cancelA()
	ca :=buildChanColumnAccumulator(ctxA,seta)
	ac.AddColumn(ca)

	var setb = &[]st.SpanBoundry[int]{
		u.Ns(3, 3),
		u.Ns(5, 11),
	}
	ctxB, cancelB := context.WithCancel(context.Background())
	defer cancelB()
	cb :=buildChanColumnAccumulator(ctxB,setb)
	ac.AddColumn(cb)
	m[1] = "SetB"

	var setc = &[]st.SpanBoundry[int]{
		u.Ns(1, 7),
		u.Ns(8, 11),
	}
	ctxC, cancelC := context.WithCancel(context.Background())
	defer cancelC()
	cc :=buildChanColumnAccumulator(ctxC,setc)
	ac.AddColumn(cc)
	m[2] = "SetC"

	header := "+-----+--------------------+------------------------------------+\n"
	fmt.Print(header)
	fmt.Print("| Seq | Begin and End      | Set Name:(Row,Row)                 |\n")
	for pos, res := range ac.Iter() {
		// check if there were errors
		if ac.Err != nil {
			fmt.Printf("Error: on Column: %s, error was: %v\n", m[ac.ErrCol], ac.Err)
			return
		}
		cols := res.GetColumns()
		names := []string{}
		for _, column := range *cols {
			str := fmt.Sprintf("%s:(%d-%d)", m[column.ColumnId], column.GetSrcId(), column.GetEndId())
			names = append(names, str)
		}
		fmt.Print(header)
		fmt.Printf("| %- 3d | Begin:% 3d, End:% 3d | %- 34s |\n",
			pos,
			res.GetBegin(),
			res.GetEnd(),
			strings.Join(names, ", "),
		)
	}
	fmt.Print(header)

}
