package st

import "testing"

func TestEmptyColumnSet(t *testing.T) {

	var cs = testDriver.NewColumnSets()
	if(cs.OverlapCount()!=-1) {
		t.Errorf("Expected the uninitalized set to be -1")
	}
	cs.columns=&[]*ColumnOverlapAccumulator[int]{}
	cs.setNext()
	cs.init()
	cs.Close()
	cs.Close()
	cs.AddColumn(nil)
	// if we get here its ok!
}
func TestInitColumSet(t *testing.T) {
	var cs = testDriver.NewColumnSets()

	var id,_ = cs.AddColumnFromSpanSlice(&[]SpanBoundry[int]{&Span[int]{Begin: 1, End: 1}})
	if id != 0 {
		t.Errorf("Expected id: 0, got %d", id)
		return
	}

	id,_ = cs.AddColumnFromSpanSlice(&[]SpanBoundry[int]{&Span[int]{Begin: 2, End: 2}})
	if id != 1 {
		t.Errorf("Expected id: 1, got %d", id)
		return
	}

	id,_ = cs.AddColumnFromSpanSlice(&[]SpanBoundry[int]{&Span[int]{Begin: 3, End: 3}})
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
	if col.GetBegin()!=1 || col.GetEnd()!=1 {
		t.Errorf("Expected Begin: 1, End: 1, got Begin: %d, End %d",col.GetBegin(),col.GetEnd())
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
	if(cs.overlap.GetBegin()!=2 || cs.overlap.GetEnd()!=2) {
		t.Errorf("Should have gotten our next span of 2->2, but got: %v",cs.overlap)
		return
	}
	cs.setNext()
	if(cs.pos!=-1) {
		t.Errorf("Should be done!, got %v",cs.overlap)
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
	var expected=[]SpanBoundry[int]{
		&Span[int]{Begin: 1, End: 1},
		&Span[int]{Begin: 2, End: 2},
		&Span[int]{Begin: 3, End: 3},
		&Span[int]{Begin: 4, End: 4},
		&Span[int]{Begin: 5, End: 5},
	}
	for id,res := range cs.Iter() {
		var cmp=expected[id];
		if(res.GetBegin()!=cmp.GetBegin() || res.GetEnd()!=cmp.GetEnd()) {
			t.Errorf("Expected: %v, Got: %v",cmp,res.GetSpan())
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
	for  range cs.Iter() {
		break
	}
	if(!cs.closed) {
		t.Error("Should now be closed")
	}
	var check=cs.Iter()
	if(check!=nil) {
		t.Error("Should not be able to create anotehr instance!")
	} 
}


