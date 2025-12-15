package st

import (
	"testing"
)

var AllSet = []SpanBoundry[int, string]{
	// sorted
	&Span[int, string]{Begin: -2, End: 2},
	&Span[int, string]{Begin: -1, End: 0},
	&Span[int, string]{Begin: -1, End: 0},
	&Span[int, string]{Begin: 0, End: 1},
	&Span[int, string]{Begin: 0, End: 1},
}
var tagA = "a"
var tagB = "b"
var MultiSet = []SpanBoundry[int, string]{
	&Span[int, string]{Begin: -1, End: 0},            // 0
	&Span[int, string]{Begin: 2, End: 2, Tag: &tagA}, //1
	&Span[int, string]{Begin: 2, End: 2, Tag: &tagB}, //1
	&Span[int, string]{Begin: 5, End: 6},             // 2
	&Span[int, string]{Begin: 9, End: 11},            // 3
	// -1
}

func TestAccumulateIter(t *testing.T) {
	var sa = testDriver.NewSpanOverlapAccumulator()
	sa.Sort = true

	var exp = []*OverlappingSpanSets[int, string]{
		{
			Span:     &Span[int, string]{Begin: -1, End: 0},
			SrcBegin: 0,
			SrcEnd:   0,
			Contains: nil,
		},
		{
			SrcBegin: 1,
			SrcEnd:   2,
			Span:     &Span[int, string]{Begin: 2, End: 2},
			Contains: &[]SpanBoundry[int, string]{
				MultiSet[1],
				MultiSet[2],
			},
		},
		{
			Span:     &Span[int, string]{Begin: 5, End: 6},
			SrcBegin: 3,
			SrcEnd:   3,
			Contains: nil,
		},
    {
      Span:     &Span[int, string]{Begin: 9, End: 11},
      SrcBegin: 4,
      SrcEnd:   4,
      Contains: nil,
    },
	}
	for idx, res := range sa.SliceIterFactory(&MultiSet) {
		var cmp = exp[idx]

		if cmp.SrcBegin != res.SrcBegin {
			t.Errorf("SrcBegin Expected: %d, got: %d", cmp.SrcBegin, res.SrcBegin)
			return
		}
		if cmp.SrcEnd != res.SrcEnd {
			t.Errorf("SrcBegin Expected: %d, got: %d", cmp.SrcEnd, res.SrcEnd)
			return
		}

	}
}

func TestMultiMultiAccumulateSet(t *testing.T) {
  var acc = testDriver.NewSpanOverlapAccumulator()
  var first = acc.Accumulate(MultMultiiSet[0])
  var next = acc.Accumulate(MultMultiiSet[1])

  if first.GetBegin() != -1 || first.GetEnd() != 0 {
    t.Errorf("Expected -1,0 got: %d,%d", first.GetBegin(), first.GetEnd())
    return
  }

  if first == next {
    t.Errorf("First and next should not be the same!")
    return
  }
  if next.GetBegin() != 2 || next.GetEnd() != 2 {
    t.Errorf("Expected 2,2 got: %d,%d", next.GetBegin(), next.GetEnd())
    return
  }
  first = next
  next = acc.Accumulate(MultMultiiSet[2])
  if next.GetBegin() != 2 || next.GetEnd() != 2 {
    t.Errorf("Expected 2,2 got: %d,%d", next.GetBegin(), next.GetEnd())
    return
  }
  if next.SrcBegin != 1 || next.SrcEnd != 2 {
    t.Errorf("Bad source index points! Expected: 1,2 got: %d,%d", next.SrcBegin, next.SrcEnd)
    return
  }
  if first != next {
    t.Errorf("First and next must be the same!")
    return
  }
  first = next
  next = acc.Accumulate(MultMultiiSet[3])
  if next.GetBegin() != 5 || next.GetEnd() != 6 {
    t.Errorf("Expected 5,6 got: %d,%d", next.GetBegin(), next.GetEnd())
    return
  }
  if first == next {
    t.Errorf("First and next should not be the same!")
    return
  }
  first = next
  next = acc.Accumulate(MultMultiiSet[4])
  if next.GetBegin() != 9 || next.GetEnd() != 11 {
    t.Errorf("Expected 9,11 got: %d,%d", next.GetBegin(), next.GetEnd())
    return
  }
  if first == next {
    t.Errorf("First and next should not be the same!")
    return
  }

  first = next
  next = acc.Accumulate(MultMultiiSet[5])
  if next.GetBegin() != 9 || next.GetEnd() != 11 {
    t.Errorf("Expected 9,11 got: %d,%d", next.GetBegin(), next.GetEnd())
    return
  }
  if first != next {
    t.Errorf("First and next should not be the same!")
    return
  }
  
  first = next
  next = acc.Accumulate(MultMultiiSet[6])
  if next.GetBegin() != 12 || next.GetEnd() != 12 {
    t.Errorf("Expected 12,12 got: %d,%d", next.GetBegin(), next.GetEnd())
    return
  }
  if first == next {
    t.Errorf("First and next should be the same!")
    return
  }

}


func TestExersizeSubIterator(t *testing.T) {

  var c = make(chan SpanBoundry[int, string], len(MultMultiiSet))
  for _, span := range MultMultiiSet {
    c <- span
  }
  close(c)
  var count = 0
  for idx,span :=range testDriver.NewSpanOverlapAccumulator().ChanIterFactory(c) {
    count++
    if(idx==3) {
      if(span.SrcBegin!=4 || span.SrcEnd!=5) {
        t.Errorf("Expected index point: 4,5, got %d,%d", span.SrcBegin,span.SrcEnd);
      }
    }
  }
  if count != 5 {
    t.Errorf("Force Final iterator block to exersize for channel test failed? Expected a total 5 for got: %d", count)
  }
}

func TestExersizeSubIteratorSlice(t *testing.T) {
  var count = 0
  //fmt.Printf("Starting Slice MultiMulti iter test testing\n")
  for idx,span :=range testDriver.NewSpanOverlapAccumulator().SliceIterFactory(&MultMultiiSet) {
    count++
    if(idx==3) {
      if(span.SrcBegin!=4 || span.SrcEnd!=5) {
        t.Errorf("Expected index point: 4,5, got %d,%d", span.SrcBegin,span.SrcEnd);
        return
      }
    }
  }
  if count != 5 {
    t.Errorf("Force Final iterator block to exersize for channel test failed? Expected a total 5 for got: %d", count)
  }
}

func TestSpanIterFactory(t *testing.T) {
  var sa=testDriver.NewSpanOverlapAccumulator().SpanStatefulAccumulator()
  if(sa.HasNext()) {
    t.Errorf("We should not have next");
    return
  }
 
  var check= sa.SetNext(MultMultiiSet[0])
  var showSa=func() {
    t.Errorf(
      "Invalid state: \n  SetNext: %v\n  HasNext: %v\n  Next: %+v\n  Current: %+v\n  Id: %d\n",
      check,
      sa.HasNext(),
      sa.Next,
      sa.Current,
      sa.Id,
      );
  }
  if(check || !sa.HasNext() ) {
    showSa()
    return
  }
  check=sa.SetNext(MultMultiiSet[1])
  if(!check) {
    showSa()
    return;
  } 
  var id,span=sa.GetNext()
  var showSpan=func() {
    t.Errorf("Span: %+v, id: %d",span,id)
  }
  if(id!=0 || span==nil) {
    showSpan()
    return
  }
  check=sa.SetNext(MultMultiiSet[2])
  if(check) {
    showSa()
  }
  check=sa.SetNext(MultMultiiSet[3])
  if(!check) {
    showSa();return
  }
  id,span=sa.GetNext()
  if(!sa.HasNext()) {
    showSa()
    return
  }
  check=sa.SetNext(MultMultiiSet[4])
  if(!check) {
    showSa();
    return;
  }
  id,span=sa.GetNext()
  check=sa.SetNext(MultMultiiSet[5])
  if(check) {
    showSa()
    return
  }
  check=sa.SetNext(MultMultiiSet[6])
  if(!check) {
    showSa()
    return
  }
  id,span=sa.GetNext()
  if(!sa.HasNext()) {
    showSa()
    return
  }
  id,span=sa.GetNext()
  if(sa.HasNext()) {
    showSa()
    return;
  }
}


func TestColumnConsolidateIter(t *testing.T) {
  var res = testDriver.NewSpanOverlapAccumulator().ColumnOverlapSliceFactory(&MultiSet)
  res.SetNext(&Span[int, string]{Begin: -1, End: 0})
  defer res.Close()
  if !res.HasNext() {
    t.Errorf("Should Has Next")
    return
  }
  if len(*res.Overlaps) != 1 {
    t.Errorf("Should have 1 element in our slice")
    return
  }
  if (*res.Overlaps)[0].GetBegin() != -1 || (*res.Overlaps)[0].GetEnd() != 0 {
    t.Errorf("Invalid first element")
    return
  }
  if res.SrcPos != 1 {
    t.Errorf("Expected SrcPos: 1, got SrcPos: %d", res.SrcPos)
    return
  }
  res.SetNext(&Span[int, string]{Begin: 1, End: 1})
  if !res.HasNext() {
    t.Errorf("Expected a next")
    return
  }
  if res.SrcPos != 1 {
    t.Errorf("Expected SrcPos: 1, got SrcPos: %d", res.SrcPos)
    return
  }
  res.SetNext(&Span[int, string]{Begin: 2, End: 5})
  if !res.HasNext() {
    t.Errorf("Expected a next")
    return
  }

  if res.SrcPos != 2 {
    t.Errorf("Expected SrcPos: 3, got SrcPos: %d", res.SrcPos)
    return
  }
  if res.SrcStart != 1 {
    t.Errorf("Expected SrcStart: 1, got: %d", res.SrcStart)
    return
  }
  if res.SrcEnd != 2 {
    t.Errorf("Expected SrcEnd: 2, got: %d", res.SrcEnd)
    return
  }
  res.SetNext(&Span[int, string]{Begin: 6, End: 11})
  if res.SrcPos != 3 {
    t.Errorf("Expected SrcPos: 3, got SrcPos: %d", res.SrcPos)
    return
  }
  if !res.HasNext() {
    t.Errorf("Expected a next")
    return
  }
  if res.SrcStart != 2 {
    t.Errorf("Expected SrcStart: 2, got: %d", res.SrcStart)
    return
  }
  if res.SrcEnd != 3 {
    t.Errorf("Expected SrcEnd: 3, got: %d", res.SrcEnd)
    return
  }

  res.SetNext(&Span[int, string]{Begin: 12, End: 12})
  if res.HasNext() {
    t.Errorf("Expected to not have next")
    return
  }

  // Make sure we close our pull iter
  res.Close()
  res = testDriver.NewSpanOverlapAccumulator().ColumnOverlapSliceFactory(&MultiSet)
  res.SetNext(&Span[int, string]{Begin: -2, End: -2})
  if res.SrcPos != 0 {
    t.Errorf("Make sure our first span is 0, got %d", res.SrcPos)
    return
  }

  res.Close()
  res = testDriver.NewSpanOverlapAccumulator().ColumnOverlapSliceFactory(&MultiSet)
  res.SetNext(&Span[int, string]{Begin: 2, End: 2})

  if res.SrcStart != 1 {
    t.Errorf("Validate we got the correct start postion, expected 1, got %d", res.SrcPos)
    return
  }
  if res.SrcEnd != 1 {
    t.Errorf("Validate we got the correct end postion, expected 1, got %d", res.SrcEnd)
    return
  }
  if res.SrcPos != 2 {
    t.Errorf("Make sure our span id is 3, got %d", res.SrcPos)
    return
  }
  if res.Next.GetBegin() != 5 || res.Next.GetEnd() != 6 {
    t.Errorf("Make sure our next range is 5->6, got %d->%d", res.Next.GetBegin(), res.Next.GetEnd())
    return
  }
  res.Close()
  res = testDriver.NewSpanOverlapAccumulator().ColumnOverlapSliceFactory(&MultiSet)
  res.SetNext(&Span[int, string]{Begin: 20, End: 20})
  if res.SrcStart != -1 {
    t.Error("Should not have a next!")
    return
  }

}


