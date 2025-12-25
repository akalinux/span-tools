package main

import (
	"cmp"
	"fmt"
	"github.com/akalinux/span-tools"
	"strings"
)

var u = st.NewSpanUtil(
	// use the standard Compare function
	cmp.Compare,

	// Define our Next function
	func(e int) int { return e + 1 },
)

func main() {

	// Build our column accumulator
	ac := u.NewColumnSets()
	// Always make sure a defer to close is scoped correctly!
	defer ac.Close()

	var Add = func(list *[]st.SpanBoundry[int]) int {
		s := u.NewSpanOverlapAccumulator().NewOlssChanStater()

		// Please note, we must start the goroutine before we add
		// the accumulator to the ColumnSets instance.  If not we
		// will run into a race condition.
		go func() {
			// Scope our cleanup code to the goroutine
			defer s.Final()

			end := len(*list) - 1
			id := 0
			span := (*list)[id]
			for s.CanAccumulate(span) {
				id++
				if id > end {
					return
				}
				span = (*list)[id]
			}
		}()

		// Adding the st.OlssChanStater instance to the ColumnSets
		return ac.AddColumnFromNewOlssChanStater(s)
	}

	// We will map our ColumnId to our Set Name
	m := make(map[int]string)

	m[Add(
		&[]st.SpanBoundry[int]{
			u.Ns(1, 2),
			u.Ns(3, 7),  // will consolidate to 3-11
			u.Ns(5, 11), // will consolidate to 3-11
		},
	)] = "SetA"

	m[Add(&[]st.SpanBoundry[int]{
		u.Ns(3, 3),
		u.Ns(5, 11),
	})] = "SetB"

	m[Add(&[]st.SpanBoundry[int]{
		u.Ns(1, 7),
		u.Ns(8, 11),
	})] = "SetC"

	header := "+-----+--------------------+------------------------------------+\n"
	fmt.Print(header)
	fmt.Print("| Seq | Begin and End      | Set Name:(Row,Row)                 |\n")
	for pos, res := range ac.Iter() {

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
	// check if there were errors
	if ac.Err != nil {
		fmt.Printf("Error: on Column: %s, error was: %v\n", m[ac.ErrCol], ac.Err)
		return
	}
	fmt.Print(header)

}
