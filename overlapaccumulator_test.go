package st

import (
	"iter"
	"testing"
)

func TestOverlapsLoopBreak(t *testing.T) {
	var count = 0
	for range testDriver.NewSpanOverlapAccumulator().NewOlssSeq2FromSbSlice(&MultiSet) {
		count++
		break
	}

	if count != 1 {
		t.Errorf("Failed To break!")
	}

}
func TestMultSetDataOverlaps(t *testing.T) {
	var ac = testDriver.NewSpanOverlapAccumulator()
	var expected = [][]int{
		{0, 0},
		{1, 2},
		{3, 3},
		{4, 5},
		{6, 6},
	}
	var count = 0
	for idx, ol := range ac.NewOlssSeq2FromSbSlice(&MultMultiiSet) {
		count++
		if expected[idx][0] != ol.SrcBegin {
			t.Errorf("Invalid SrcBegin, expected: %d, got %d for position: %d", expected[idx][0], ol.SrcBegin, idx)
			return
		}
		if expected[idx][1] != ol.SrcEnd {
			t.Errorf("Invalid SrcEnd, expected: %d, got %d for position: %d", expected[idx][1], ol.SrcEnd, idx)
			return
		}
		_, span := ol.GetFirstSpan()
		if span == nil {
			t.Errorf("Expected non nil span, got nil?")
			return
		}
		_, span = ol.GetLastSpan()
		if span == nil {
			t.Errorf("Expected non nil span, got nil?")
			return
		}
		raw := ol.GetSources()
		if len(*raw) == 0 {
			t.Errorf("Should get at least one span, got %d on set: %d", len(*raw),idx)
			return
		}
	  sets :=ol.GetOverlaps()
		if len(*sets)==0 {
			t.Errorf("Should never get 0 overlaps!")
			return
		}
		if(ol.GetSrcId()!=ol.SrcBegin || ol.GetEndId()!=ol.SrcEnd) {
			t.Error("Interface methouds for index begin and end should match the internal structure states")
			return
		} 
	}
	if count != len(expected) {
		t.Errorf("Iterator count missmatch!, expected %d, got %d", len(expected), count)
	}
}

func TestPull2MultSetDataOverlaps(t *testing.T) {
	var next, stop = iter.Pull2(testDriver.NewSpanOverlapAccumulator().NewOlssSeq2FromSbSlice(&MultMultiiSet))
	defer stop()
	var count = 0
	var expected = [][]int{
		{0, 0},
		{1, 2},
		{3, 3},
		{4, 5},
		{6, 6},
	}
	var idx, ol, ok = next()
	for ok {
		count++
		if expected[idx][0] != ol.SrcBegin {
			t.Errorf("Invalid SrcBegin, expected: %d, got %d for position: %d", expected[idx][0], ol.SrcBegin, idx)
			return
		}
		if expected[idx][1] != ol.SrcEnd {
			t.Errorf("Invalid SrcEnd, expected: %d, got %d for position: %d", expected[idx][1], ol.SrcEnd, idx)
			return
		}
		if !ol.IsUnique() && ol.GetContains() == nil {
			t.Errorf("Contains should not be empty if the object is uniqe")
			return
		}
		idx, ol, ok = next()
	}
	if count != len(expected) {
		t.Errorf("Iterator count missmatch!, expected %d, got %d", len(expected), count)
	}
}

func TestOverlapAdjacentConsolidate(t *testing.T) {
	var ac = testDriver.NewSpanOverlapAccumulator()
	ac.Consolidate = true
	var expected = [][]int{
		{0, 3},
		{5, 9},
	}
	var count = -1
	for idx, ol := range ac.NewOlssSeq2FromSbSlice(&[]SpanBoundry[int]{
		&Span[int]{Begin: 0, End: 2},
		&Span[int]{Begin: 3, End: 3},
		&Span[int]{Begin: 5, End: 6},
		&Span[int]{Begin: 7, End: 9},
	}) {
		count++
		if expected[count][0] != ol.GetBegin() {
			t.Errorf("Invalid Begin, expected: %d, got %v for position: %d", expected[idx][0], ol.Span, idx)
			return
		}
		if expected[count][1] != ol.GetEnd() {
			t.Errorf("Invalid End, expected: %d, got %v for position: %d", expected[idx][1], ol.Span, idx)
			return
		}
	}
	if count != len(expected)-1 {
		t.Errorf("Iterator count missmatch!, expected %d, got %d", len(expected), count)
	}
}
