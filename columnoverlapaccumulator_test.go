package st

import (
	"testing"
)

type IterValidate struct {
	Next     *Span[int, string]
	SrcStart int
	SrcEnd   int
}

func TestColumConsolidateChannelOverlapAccumulator(t *testing.T) {
  var list=[]*OverlappingSpanSets[int,string]{}
	
	
	for _,span := range testDriver.NewSpanOverlapAccumulator().SliceIterFactory(&MultMultiiSet) {
		list=append(list,span)
	}
	var ts=make(chan *OverlappingSpanSets[int,string], len(list))
	for _,span := range list {
		ts <- span
	}
	close(ts);
	var ca=testDriver.NewSpanOverlapAccumulator().ColumnChanOverlapSpanSetsFactory(ts)
	ca.SetNext(MultMultiiSet[len(MultMultiiSet)-1])
	
	var _,ok = <-ts
	if(ok) {
		t.Error("Channel should be depleted")
		
	}

}

func testOverlapStruct(expected []IterValidate,t *testing.T,src *[]SpanBoundry[int,string]) {
	var res = testDriver.NewSpanOverlapAccumulator().ColumnOverlapSliceFactory(src)
	defer res.Close()
	for pos, conf := range expected {

		res.SetNext(conf.Next)

		if conf.SrcStart != res.SrcStart {
			t.Errorf("Bad position on result: %d, expected SrcStart: %d, Got: %d", pos, conf.SrcStart, res.SrcStart)
			return
		}
		if conf.SrcEnd != res.SrcEnd {
			t.Errorf("Bad position on result: %d, expected SrcEnd: %d, Got: %d", pos, conf.SrcEnd, res.SrcEnd)
			return
		}
	}

	res.Close()
	if res.HasNext() {
		t.Errorf("Iterature struct shold no longer have next!")
	}
}

func TestColumConsolidateLookBack(t *testing.T) {
	var expected = []IterValidate{
		{
			Next:     &Span[int, string]{Begin: -1, End:5},
			SrcStart: 0,
			SrcEnd:   3,
		},
		{
			Next:     &Span[int, string]{Begin: 2, End:6},
			SrcStart: 1,
			SrcEnd:   3,
		},
		{
			Next:     &Span[int, string]{Begin: 13, End:13},
			SrcStart: -1,
			SrcEnd:   -1,
		},
 }
  testOverlapStruct(expected,t,&MultMultiiSet)
}
func TestColumnConsolidateIter(t *testing.T) {

	var expected = []IterValidate{
		{
			Next:     &Span[int, string]{Begin: -1, End: -1},
			SrcStart: 0,
			SrcEnd:   0,
		},
		{
			Next:     &Span[int, string]{Begin: 1, End: 1},
			SrcStart: -1,
			SrcEnd:   -1,
		},
		{
			Next:     &Span[int, string]{Begin: 2, End: 5},
			SrcStart: 1,
			SrcEnd:   3,
		},
		{
			Next:     &Span[int, string]{Begin: 6, End: 11},
			SrcStart: 3,
			SrcEnd:   5,
		},
		{
			Next:     &Span[int, string]{Begin: 12, End: 12},
			SrcStart: 6,
			SrcEnd:   6,
		},
	}
  testOverlapStruct(expected,t,&MultMultiiSet)
	
}
