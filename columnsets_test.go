package st

import (
	"errors"
	"testing"
)

func TestEmptyColumnSet(t *testing.T) {

	var cs = testDriver.NewColumnSets()
	if cs.OverlapCount() != -1 {
		t.Errorf("Expected the uninitalized set to be -1")
	}
	cs.columns = &[]*ColumnOverlapAccumulator[int]{}
	cs.setNext()
	cs.init()
	cs.Close()
	cs.Close()
	cs.AddColumn(nil)
	// if we get here its ok!
}
func TestInitColumSet(t *testing.T) {
	var cs = testDriver.NewColumnSets()

	var id, _ = cs.AddColumnFromSpanSlice(&[]SpanBoundry[int]{&Span[int]{Begin: 1, End: 1}})
	if id != 0 {
		t.Errorf("Expected id: 0, got %d", id)
		return
	}

	id, _ = cs.AddColumnFromSpanSlice(&[]SpanBoundry[int]{&Span[int]{Begin: 2, End: 2}})
	if id != 1 {
		t.Errorf("Expected id: 1, got %d", id)
		return
	}

	id, _ = cs.AddColumnFromSpanSlice(&[]SpanBoundry[int]{&Span[int]{Begin: 3, End: 3}})
	if id != 2 {
		t.Errorf("Expected id: 1, got %d", id)
		return
	}
	cs.init()

	if cs.OverlapCount() != 1 {
		t.Errorf("Expected OverlapCount of: 1, got: %d", cs.OverlapCount())
		return
	}
	var col = (*cs.current)[0]
	if col.ColumnId != 0 || col.GetSrcId() != 0 || col.GetEndId() != 0 {
		t.Errorf("Bad data position, expected 0, got: ColumnId: %d, src: %d, end: %d",
			col.ColumnId,
			col.GetSrcId(),
			col.GetEndId(),
		)
		return
	}
	if col.GetBegin() != 1 || col.GetEnd() != 1 {
		t.Errorf("Expected Begin: 1, End: 1, got Begin: %d, End %d", col.GetBegin(), col.GetEnd())
		return
	}

	defer cs.Close()

}

func TestFullIter(t *testing.T) {
	var cs = testDriver.NewColumnSets()

	cs.AddColumnFromSpanSlice(&[]SpanBoundry[int]{&Span[int]{Begin: 1, End: 1}})
	cs.AddColumnFromSpanSlice(&[]SpanBoundry[int]{&Span[int]{Begin: 2, End: 2}})

	cs.init()
	defer cs.Close()
	cs.setNext()

	if cs.OverlapCount() != 1 {
		t.Errorf("Expected OverlapCount of: 1, got: %d", cs.OverlapCount())
		return
	}
	if cs.overlap.GetBegin() != 2 || cs.overlap.GetEnd() != 2 {
		t.Errorf("Should have gotten our next span of 2->2, but got: %v", cs.overlap)
		return
	}
	cs.setNext()
	if cs.pos != -1 {
		t.Errorf("Should be done!, got %v", cs.overlap)
	}
}

func TestColumSetIter(t *testing.T) {

	var cs = testDriver.NewColumnSets()
	cs.AddColumnFromSpanSlice(&[]SpanBoundry[int]{
		&Span[int]{Begin: 1, End: 2},
		&Span[int]{Begin: 1, End: 1},
		&Span[int]{Begin: 3, End: 3},
		&Span[int]{Begin: 4, End: 4},
	})
	cs.AddColumnFromSpanSlice(&[]SpanBoundry[int]{
		&Span[int]{Begin: 2, End: 3},
	})
	cs.AddColumnFromSpanSlice(&[]SpanBoundry[int]{
		&Span[int]{Begin: 3, End: 5},
	})
	var expected = []SpanBoundry[int]{
		&Span[int]{Begin: 1, End: 1},
		&Span[int]{Begin: 2, End: 2},
		&Span[int]{Begin: 3, End: 3},
		&Span[int]{Begin: 4, End: 4},
		&Span[int]{Begin: 5, End: 5},
	}
	for id, res := range cs.Iter() {
		var cmp = expected[id]
		if res.GetBegin() != cmp.GetBegin() || res.GetEnd() != cmp.GetEnd() {
			t.Errorf("Expected: %v, Got: %v", cmp, res.GetSpan())
		}
		// force code to be tested
		res.GetColumns()
		res.GetSpan()
	}
}

func TestColumSetIterBreakTest(t *testing.T) {
	var cs = testDriver.NewColumnSets()
	cs.AddColumnFromSpanSlice(&[]SpanBoundry[int]{
		&Span[int]{Begin: 1, End: 2},
		&Span[int]{Begin: 1, End: 1},
		&Span[int]{Begin: 3, End: 3},
		&Span[int]{Begin: 4, End: 4},
	})
	cs.AddColumnFromSpanSlice(&[]SpanBoundry[int]{
		&Span[int]{Begin: 2, End: 3},
	})
	cs.AddColumnFromSpanSlice(&[]SpanBoundry[int]{
		&Span[int]{Begin: 3, End: 5},
	})
	for range cs.Iter() {
		break
	}
	if !cs.closed {
		t.Error("Should now be closed")
	}
	var check = cs.Iter()
	if check != nil {
		t.Error("Should not be able to create anotehr instance!")
	}
}

func TestColumSetsInitError(t *testing.T) {
	col := &[]*OverlappingSpanSets[int]{
		{
			Err: errors.New("Force init error proof"),
		},
	}
	cols := testDriver.NewColumnSets()
	cols.AddColumnFromOverlappingSpanSets(col)
	count := 0
	for range cols.Iter() {
		count++
	}
	if count != 0 {
		t.Error("Should have 0 rows!")
	}
	if nil == cols.Err {
		t.Error("We should be in an error state at this point!")
	}
}

func TestColumSetsSeqError(t *testing.T) {
	u := NewSpanUtil(
		testDriver.Cmp,
		testDriver.Next,
	)
	u.Validate = true
	u.Sort = false
	ac := u.NewColumnSets()


	ac.AddColumnFromSpanSlice(&[]SpanBoundry[int]{
		u.Ns(3, 3),
		u.Ns(5, 11),
	})
	ac.AddColumnFromSpanSlice(&[]SpanBoundry[int]{
		u.Ns(1, 7),
		u.Ns(8, 11),
	})
	ac.AddColumnFromSpanSlice(&[]SpanBoundry[int]{
		// set 0
		u.Ns(1, 2),
		u.Ns(3, 7), // will consolidate to 3-11
		u.Ns(5, 11), // will consolidate to 3-11
		// set 1
		u.Ns(12, 13), // will consolidate to 3-11
		// set 3
		u.Ns(13, 13), // will consolidate to 3-11
		u.Ns(13, 14), // will consolidate to 3-11
		u.Ns(6, 6), // will consolidate to 3-11

	})
	count :=0
	for range ac.Iter() {
		
		count++
	}
	if ac.Err==nil {
		t.Error("Should be in an error state")
	}
	if ac.ErrCol!=2 {
		t.Errorf("Expected error columnId to be 2, got %d, Error was: %v",ac.ErrCol,ac.Err);
	}
}

func TestColumnSetFactoryChanIter(t *testing.T) {

	var cs = testDriver.NewColumnSets()

	colId:=-1	
	var Add=func(list *[]SpanBoundry[int]) {
		colId++
	  s :=testDriver.NewSpanOverlapAccumulator().NewOlssChanStater()
		go func() {
			defer s.Final()
			end :=len(*list) -1
			id :=0
			span :=(*list)[0]
			for s.CanAccumulate(span) {
				id++
				if id > end {
					return
				}
				span=(*list)[id]
			}
		}()
		cs.AddColumnFromNewOlssChanStater(s)
	}
  Add(&[]SpanBoundry[int]{
		&Span[int]{Begin: 1, End: 2},
		&Span[int]{Begin: 1, End: 1},
		&Span[int]{Begin: 3, End: 3},
		&Span[int]{Begin: 4, End: 4},
	})
  Add(&[]SpanBoundry[int]{
		&Span[int]{Begin: 2, End: 3},
	})
	Add(&[]SpanBoundry[int]{
		&Span[int]{Begin: 3, End: 5},
	})
	var expected = []SpanBoundry[int]{
		&Span[int]{Begin: 1, End: 1},
		&Span[int]{Begin: 2, End: 2},
		&Span[int]{Begin: 3, End: 3},
		&Span[int]{Begin: 4, End: 4},
		&Span[int]{Begin: 5, End: 5},
	}
	for id, res := range cs.Iter() {
		var cmp = expected[id]
		if res.GetBegin() != cmp.GetBegin() || res.GetEnd() != cmp.GetEnd() {
			t.Errorf("Expected: %v, Got: %v", cmp, res.GetSpan())
		}
	}
}
func TestColumnSetFactoryChanIterBreak(t *testing.T) {

	var cs = testDriver.NewColumnSets()

	colId:=-1	
	var Add=func(list *[]SpanBoundry[int]) *OlssChanStater[int] {
		colId++
	  s :=testDriver.NewSpanOverlapAccumulator().NewOlssChanStater()
		go func() {
			defer s.Final()
			end :=len(*list) -1
			id :=0
			span :=(*list)[0]
			for s.CanAccumulate(span) {
				id++
				if id > end {
					return
				}
				span=(*list)[id]
			}
		}()
		cs.AddColumnFromNewOlssChanStater(s)
		return s
	}
  s:=Add(&[]SpanBoundry[int]{
		&Span[int]{Begin: 1, End: 2},
		&Span[int]{Begin: 1, End: 1},
		&Span[int]{Begin: 3, End: 3},
		&Span[int]{Begin: 4, End: 4},
	})
  Add(&[]SpanBoundry[int]{
		&Span[int]{Begin: 2, End: 3},
	})
	Add(&[]SpanBoundry[int]{
		&Span[int]{Begin: 3, End: 5},
	})
	var expected = []SpanBoundry[int]{
		&Span[int]{Begin: 1, End: 1},
		&Span[int]{Begin: 2, End: 2},
		&Span[int]{Begin: 3, End: 3},
		&Span[int]{Begin: 4, End: 4},
		&Span[int]{Begin: 5, End: 5},
	}
	for id, res := range cs.Iter() {
		var cmp = expected[id]
		if res.GetBegin() != cmp.GetBegin() || res.GetEnd() != cmp.GetEnd() {
			t.Errorf("Expected: %v, Got: %v", cmp, res.GetSpan())
		}
		break
	}
	// should do nothing.. if the code is broken.. this will cause a panic!
	if !s.IsShutDown {
		t.Error("Should have all ready been shut down!")
		return
	}
	
	// need to flush the channel, to test if it was actually closed
	_,ok := <- s.Chan
	for ok {
	  _,ok = <- s.Chan
	}
	if s.CanAccumulate(&Span[int]{Begin: 0,End: 1}) {
		t.Error("Should no longer be able to accumulte!")
		return
	}
	if s.Push(&OverlappingSpanSets[int]{}) {
		t.Error("Should no longer be able to push")
	}
	// if something is broken, this will cause a panic
	s.Shutdown()
}



