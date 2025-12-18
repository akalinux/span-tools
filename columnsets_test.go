package st

import "testing"

func TestEmptyColumnSet(t *testing.T) {

	var cs = testDriver.NewColumnSets()
	if(cs.OverlapCount()!=-1) {
		t.Errorf("Expected the uninitalized set to be -1")
	}
	cs.Columns=&[]*ColumnOverlapAccumulator[int]{}
	cs.SetNext()
	cs.Init()
	cs.Close()
	cs.Close()
	cs.AddColumn(nil)
	// if we get here its ok!
}
func TestInitColumSet(t *testing.T) {
	var cs = testDriver.NewColumnSets()

	var id = cs.AddColumnFromSpanSlice(&[]SpanBoundry[int]{&Span[int]{Begin: 1, End: 1}})
	if id != 0 {
		t.Errorf("Expected id: 0, got %d", id)
		return
	}

	id = cs.AddColumnFromSpanSlice(&[]SpanBoundry[int]{&Span[int]{Begin: 2, End: 2}})
	if id != 1 {
		t.Errorf("Expected id: 1, got %d", id)
		return
	}

	id = cs.AddColumnFromSpanSlice(&[]SpanBoundry[int]{&Span[int]{Begin: 3, End: 3}})
	if id != 2 {
		t.Errorf("Expected id: 1, got %d", id)
		return
	}
	cs.Init()

	if cs.OverlapCount() != 1 {
		t.Errorf("Expected OverlapCount of: 1, got: %d", cs.OverlapCount())
		return
	}
	var col = (*cs.Current)[0]
	if col.ColumnId != 0 || col.GetSrcStart() != 0 || col.GetSrcEnd() != 0 {
		t.Errorf("Bad data position, expected 0, got: ColumnId: %d, src: %d, end: %d",
			col.ColumnId,
			col.GetSrcStart(),
			col.GetSrcEnd(),
		)
		return
	}
	if col.GetBegin()!=1 || col.GetEnd()!=1 {
		t.Errorf("Expected Begin: 1, End: 1, got Begin: %d, End %d",col.GetBegin(),col.GetEnd())
		return
	}

	defer cs.Close()

}

func TestManualIterBasicColumSet(t *testing.T) {
	var cs = testDriver.NewColumnSets()

  cs.AddColumnFromSpanSlice(&[]SpanBoundry[int]{&Span[int]{Begin: 1, End: 1}})
	cs.AddColumnFromSpanSlice(&[]SpanBoundry[int]{&Span[int]{Begin: 2, End: 2}})
	cs.AddColumnFromSpanSlice(&[]SpanBoundry[int]{&Span[int]{Begin: 3, End: 3}})
	
	cs.Init()

  cs.SetNext()
	if cs.OverlapCount() != 1 {
		t.Errorf("Expected OverlapCount of: 1, got: %d", cs.OverlapCount())
		return
	}
	var col = (*cs.Current)[0]
	if col.ColumnId != 1 || col.GetSrcStart() != 0 || col.GetSrcEnd() != 0 {
		t.Errorf("Bad data position, got: ColumnId: %d, src: %d, end: %d",
			col.ColumnId,
			col.GetSrcStart(),
			col.GetSrcEnd(),
		)
		return
	}
	if col.GetBegin()!=2 || col.GetEnd()!=2 {
		t.Errorf("Expected Begin: 2, End: 2, got Begin: %d, End %d",col.GetBegin(),col.GetEnd())
		return
	}
  cs.SetNext()
	col = (*cs.Current)[0]
	if(col.GetOverlaps()==nil) {
		t.Error("Overlaps should not be nil in this case!");
		return
	}
	if cs.OverlapCount() != 1 {
		t.Errorf("Expected OverlapCount of: 1, got: %d", cs.OverlapCount())
		return
	}
	if col.ColumnId != 2 || col.GetSrcStart() != 0 || col.GetSrcEnd() != 0 {
		t.Errorf("Bad data position, got: ColumnId: %d, src: %d, end: %d",
			col.ColumnId,
			col.GetSrcStart(),
			col.GetSrcEnd(),
		)
		return
	}
	if col.GetBegin()!=3 || col.GetEnd()!=3 {
		t.Errorf("Expected Begin: 3, End: 3, got Begin: %d, End %d",col.GetBegin(),col.GetEnd())
		return
	}
	
	cs.SetNext()
	if(cs.Pos!=-1) {
		t.Errorf("Our last state should be -1")
		return
	}

	defer cs.Close()

}
