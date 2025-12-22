package st

import "testing"
func TestBadOrder(t *testing.T) {
  var list = &[]SpanBoundry[int]{
    &Span[int]{Begin: 9, End: 11},
    &Span[int]{Begin: 2, End: 2},
  }
	testDriver :=NewSpanUtil(testDriver.Cmp,testDriver.Next)
  testDriver.Validate = true
  testDriver.Sort= false
	
  for id, span := range testDriver.NewSpanOverlapAccumulator().NewOlssSeq2FromSbSlice(list) {
    if id > 0 {
      t.Errorf("Should stop at 0")
      return
    }
    if span.GetBegin() != 2 {
      t.Error("Should have span 1, got span 0")
    }
  }
}

func TestNilSliceIter(t *testing.T) {
  for range testDriver.NewSpanOverlapAccumulator().NewOlssSeq2FromSbSlice(nil) {
     t.Errorf("Should have gotten no ranges!")
     return
   }
}

func TestBadInitValue(t *testing.T) {
			
  var list = &[]SpanBoundry[int]{
    &Span[int]{Begin: 13, End: 11},
  }
  testDriver.Validate = true
  for _,span :=range testDriver.NewSpanOverlapAccumulator().NewOlssSeq2FromSbSlice(list) {
		if span.Err ==nil {
      t.Errorf("Should have gotten no valid ranges!, but we got: %v, %v",span.Span,span.Err)
      return
		}
  }

}