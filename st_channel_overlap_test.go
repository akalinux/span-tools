package st

import (
	"testing"
)

func TestOverlapChannel(t *testing.T) {
	var ac = testDriver.NewSpanOverlapAccumulator()
	var list = []*OverlappingSpanSets[int, string]{}
	for _, ol := range ac.SliceIterFactory(&MultMultiiSet) {
		list = append(list, ol)
	}
	ch := make(chan *OverlappingSpanSets[int, string], len(list))
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
	var itb=testDriver.NewSpanOverlapAccumulator().SliceIterFactoryOverlaps(&list);
  count=0
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
