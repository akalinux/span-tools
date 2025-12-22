package st

import (
	"testing"
)

var AllSet = []SpanBoundry[int]{
	// sorted
	&Span[int]{Begin: -2, End: 2},
	&Span[int]{Begin: -1, End: 0},
	&Span[int]{Begin: -1, End: 0},
	&Span[int]{Begin: 0, End: 1},
	&Span[int]{Begin: 0, End: 1},
}

var MultiSet = []SpanBoundry[int]{
	&Span[int]{Begin: -1, End: 0},            // 0
	&Span[int]{Begin: 2, End: 2}, //1
	&Span[int]{Begin: 2, End: 2}, //1
	&Span[int]{Begin: 5, End: 6},             // 2
	&Span[int]{Begin: 9, End: 11},            // 3
	// -1
}

func TestAccumulateIter(t *testing.T) {
	var sa = testDriver.NewSpanOverlapAccumulator()
	sa.Sort = true

	var exp = []*OverlappingSpanSets[int]{
		{
			Span:     &Span[int]{Begin: -1, End: 0},
			SrcBegin: 0,
			SrcEnd:   0,
			Contains: nil,
		},
		{
			SrcBegin: 1,
			SrcEnd:   2,
			Span:     &Span[int]{Begin: 2, End: 2},
			Contains: &[]SpanBoundry[int]{
				MultiSet[1],
				MultiSet[2],
			},
		},
		{
			Span:     &Span[int]{Begin: 5, End: 6},
			SrcBegin: 3,
			SrcEnd:   3,
			Contains: nil,
		},
		{
			Span:     &Span[int]{Begin: 9, End: 11},
			SrcBegin: 4,
			SrcEnd:   4,
			Contains: nil,
		},
	}
	for idx, res := range sa.NewOverlappingSpanSetsIterSeq2FromSpanBoundrySlice(&MultiSet) {
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
	var first,_ = acc.Accumulate(MultMultiiSet[0])
	var next,_ = acc.Accumulate(MultMultiiSet[1])

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
	next,_ = acc.Accumulate(MultMultiiSet[2])
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
	next,_ = acc.Accumulate(MultMultiiSet[3])
	if next.GetBegin() != 5 || next.GetEnd() != 6 {
		t.Errorf("Expected 5,6 got: %d,%d", next.GetBegin(), next.GetEnd())
		return
	}
	if first == next {
		t.Errorf("First and next should not be the same!")
		return
	}
	first = next
	next,_ = acc.Accumulate(MultMultiiSet[4])
	if next.GetBegin() != 9 || next.GetEnd() != 11 {
		t.Errorf("Expected 9,11 got: %d,%d", next.GetBegin(), next.GetEnd())
		return
	}
	if first == next {
		t.Errorf("First and next should not be the same!")
		return
	}

	first = next
	next,_ = acc.Accumulate(MultMultiiSet[5])
	if next.GetBegin() != 9 || next.GetEnd() != 11 {
		t.Errorf("Expected 9,11 got: %d,%d", next.GetBegin(), next.GetEnd())
		return
	}
	if first != next {
		t.Errorf("First and next should not be the same!")
		return
	}

	first = next
	next,_ = acc.Accumulate(MultMultiiSet[6])
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

	var c = make(chan SpanBoundry[int], len(MultMultiiSet))
	for _, span := range MultMultiiSet {
		c <- span
	}
	close(c)
	var count = 0
	for idx, span := range testDriver.NewSpanOverlapAccumulator().NewOverlappingSpanSetsFromSpanBoundryChan(c) {
		count++
		if idx == 3 {
			if span.SrcBegin != 4 || span.SrcEnd != 5 {
				t.Errorf("Expected index point: 4,5, got %d,%d", span.SrcBegin, span.SrcEnd)
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
	for idx, span := range testDriver.NewSpanOverlapAccumulator().NewOverlappingSpanSetsIterSeq2FromSpanBoundrySlice(&MultMultiiSet) {
		count++
		if idx == 3 {
			if span.SrcBegin != 4 || span.SrcEnd != 5 {
				t.Errorf("Expected index point: 4,5, got %d,%d", span.SrcBegin, span.SrcEnd)
				return
			}
		}
	}
	if count != 5 {
		t.Errorf("Force Final iterator block to exersize for channel test failed? Expected a total 5 for got: %d", count)
	}
}

func TestSpanIterFactory(t *testing.T) {
	var sa = testDriver.NewSpanOverlapAccumulator().NewSpanIterSeq2Stater()
	if sa.HasNext() {
		t.Errorf("We should not have next")
		return
	}

	var check = sa.SetNext(MultMultiiSet[0])
	var showSa = func() {
		t.Errorf(
			"Invalid state: \n  SetNext: %v\n  HasNext: %v\n  Next: %+v\n  Current: %+v\n  Id: %d\n",
			check,
			sa.HasNext(),
			sa.Next,
			sa.Current,
			sa.Id,
		)
	}
	if check || !sa.HasNext() {
		showSa()
		return
	}
	check = sa.SetNext(MultMultiiSet[1])
	if !check {
		showSa()
		return
	}
	var id, span = sa.GetNext()
	var showSpan = func() {
		t.Errorf("Span: %+v, id: %d", span, id)
	}
	if id != 0 || span == nil {
		showSpan()
		return
	}
	check = sa.SetNext(MultMultiiSet[2])
	if check {
		showSa()
	}
	check = sa.SetNext(MultMultiiSet[3])
	if !check {
		showSa()
		return
	}
	id, span = sa.GetNext()
	if !sa.HasNext() {
		showSa()
		return
	}
	check = sa.SetNext(MultMultiiSet[4])
	if !check {
		showSa()
		return
	}
	id, span = sa.GetNext()
	check = sa.SetNext(MultMultiiSet[5])
	if check {
		showSa()
		return
	}
	check = sa.SetNext(MultMultiiSet[6])
	if !check {
		showSa()
		return
	}
	id, span = sa.GetNext()
	if !sa.HasNext() {
		showSa()
		return
	}
	id, span = sa.GetNext()
	if sa.HasNext() {
		showSa()
		return
	}
}

