package main

import (
	"cmp"
	"fmt"
	"strings"
	"github.com/akalinux/span-tools"
)

func main() {
	u := st.NewSpanUtil(
		// use the standard Compare function
		cmp.Compare,
		
		// Define our Next function
		func(e int) int { return e + 1 },
	)
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
	ac.AddColumnFromSpanSlice(seta)
	m[0] = "SetA"

	var setb = &[]st.SpanBoundry[int]{
		u.Ns(3, 3),
		u.Ns(5, 11),
	}
	ac.AddColumnFromSpanSlice(setb)
	m[1] = "SetB"

	var setc = &[]st.SpanBoundry[int]{
		u.Ns(1, 7),
		u.Ns(8, 11),
	}
	ac.AddColumnFromSpanSlice(setc)
	m[2] = "SetC"


	header := "+-----+--------------------+------------------------------------+\n"
	fmt.Print(header)
	fmt.Print("| Seq | Begin and End      | Set Name:(Row,Row)                 |\n")
	for pos, res := range ac.Iter() {
		// check if there were errors
		if ac.Err != nil {
			fmt.Printf("Error: on Column: %s, error was: %v\n",m[ac.ErrCol],ac.Err)
			return
		}
		cols := res.GetColumns()
		names := []string{}
		for _, column := range *cols {
			str :=fmt.Sprintf("%s:(%d-%d)",m[column.ColumnId],column.GetSrcId(),column.GetEndId())
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
