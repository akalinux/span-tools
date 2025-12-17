package st

import (
	"testing"
)

func MakeOverlapTestList() *[]*OverlappingSpanSets[int] {
	var ac = testDriver.NewSpanOverlapAccumulator()
	var list = []*OverlappingSpanSets[int]{}
	for _, ol := range ac.SliceIterFactory(&MultMultiiSet) {
		list = append(list, ol)
	}
	return &list
}

func TestOverlapChannel(t *testing.T) {
	var ac = testDriver.NewSpanOverlapAccumulator()
	var list =*MakeOverlapTestList() 
	ch := make(chan *OverlappingSpanSets[int], len(list))
	for _, ol := range list {
		ch <- ol
	}
	close(ch)
	var count = 0
	for id, value := range ac.ChanIterFactoryOverlaps(ch) {
		count++
		if list[id] != value {
			t.Errorf("Error, wrong object ref in chan iter??")
			return
		}
	}
	if count != len(list) {
		t.Errorf("Iterator count mismatch in chan iter??")
	}
}
func TestBreakLoopOverlapChannel(t *testing.T) {
	var ac = testDriver.NewSpanOverlapAccumulator()
	var list =*MakeOverlapTestList() 
	ch := make(chan *OverlappingSpanSets[int], len(list))
	for _, ol := range list {
		ch <- ol
	}
	close(ch)
	var count = 0
	for range ac.ChanIterFactoryOverlaps(ch) {
		count++
		break
	}
	if count != 1 {
		t.Errorf("Iterator count mismatch in chan iter??")
	}
}
func TestNilOverlapChannel(t *testing.T) {
	var itb=testDriver.NewSpanOverlapAccumulator().ChanIterFactoryOverlaps(nil);
	var count=0
	for  range itb  {
	  count++
	}
	if(count!=0) {
		t.Error("Should Not get any elements in our loop")
	}
}
func TestNilOverlapSlice(t *testing.T) {
	var itb=testDriver.NewSpanOverlapAccumulator().SliceIterFactoryOverlaps(nil);

  var count=0
	for  range itb  {
	  count++
  }
	if(count!=0) {
		t.Error("Should Not get any elements in our loop")
	}
}

func TestBreakLoopOverlapSlice(t *testing.T) {
	var list =*MakeOverlapTestList() 
	var itb=testDriver.NewSpanOverlapAccumulator().SliceIterFactoryOverlaps(&list);
  var count=0
	for range itb  {
		count++
		break
	}
	if count != 1 {
		t.Errorf("Iterator count mismatch in slice iter??")
	}
}

func TestOverlapSlice(t *testing.T) {
	var list =*MakeOverlapTestList() 
	var itb=testDriver.NewSpanOverlapAccumulator().SliceIterFactoryOverlaps(&list);
  var count=0
	for id, value := range itb  {
		count++
		if list[id] != value {
			t.Errorf("Error, wrong object ref in slice iter??")
			return
		}
	}
	if count != len(list) {
		t.Errorf("Iterator count mismatch in slice iter??")
	}
}
