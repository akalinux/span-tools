package st

import("testing"
"iter")

func TestOverlapsLoopBreak(t *testing.T) {
	var count=0
	for range testDriver.NewSpanOverlapAccumulator().SliceIterFactory(&MultiSet) {
	  count++
		break
	}
	
	if(count!=1) {
		t.Errorf("Failed To break!");
	}
	
}
func TestMultSetDataOverlaps(t *testing.T) {
	var ac = testDriver.NewSpanOverlapAccumulator()
	var expected =[][]int{
		{0,0},
		{1,2},
		{3,3},
		{4,5},
		{6,6},
	}
	var count=0
	for idx,ol := range ac.SliceIterFactory(&MultMultiiSet) {
		count++
		if(expected[idx][0]!=ol.SrcBegin) {
			t.Errorf("Invalid SrcBegin, expected: %d, got %d for position: %d",expected[idx][0],ol.SrcBegin,idx)
			return
		}
		if(expected[idx][1]!=ol.SrcEnd) {
			t.Errorf("Invalid SrcEnd, expected: %d, got %d for position: %d",expected[idx][1],ol.SrcEnd,idx)
			return
		}
	}
	if(count!=len(expected)) {
		t.Errorf("Iterator count missmatch!, expected %d, got %d",len(expected),count)
	}
}

func TestPull2MultSetDataOverlaps(t *testing.T) {
	var next,stop=iter.Pull2(testDriver.NewSpanOverlapAccumulator().SliceIterFactory(&MultMultiiSet))
	defer stop()
	var count=0
	var expected =[][]int{
		{0,0},
		{1,2},
		{3,3},
		{4,5},
		{6,6},
	}
	var idx,ol,ok =next()
	for ok {
		count++
		if(expected[idx][0]!=ol.SrcBegin) {
			t.Errorf("Invalid SrcBegin, expected: %d, got %d for position: %d",expected[idx][0],ol.SrcBegin,idx)
			return
		}
		if(expected[idx][1]!=ol.SrcEnd) {
			t.Errorf("Invalid SrcEnd, expected: %d, got %d for position: %d",expected[idx][1],ol.SrcEnd,idx)
			return
		}
		if !ol.IsUnique() && ol.GetContains()==nil {
			t.Errorf("Contains should not be empty if the object is uniqe")
			return
		}
	  idx,ol,ok =next()
	}
	if(count!=len(expected)) {
		t.Errorf("Iterator count missmatch!, expected %d, got %d",len(expected),count)
	}
}