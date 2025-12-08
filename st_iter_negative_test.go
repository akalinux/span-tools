package st

import "testing"
func TestBadOrder(t *testing.T) {
  var list = &[]SpanBoundry[int, string]{
    &Span[int, string]{Begin: 9, End: 11},
    &Span[int, string]{Begin: 2, End: 2},
  }
  testDriver.Validate = true
  for id, span := range testDriver.NewSpanOverlapAccumulator().SliceIterFactory(list) {
    if id > 0 {
      t.Errorf("Should stop at 0")
      return
    }
    if span.GetBegin() != 9 {
      t.Error("Should have span 0, got span 1")
    }
  }
}

func TestNilSliceIter(t *testing.T) {
  for range testDriver.NewSpanOverlapAccumulator().SliceIterFactory(nil) {
     t.Errorf("Should have gotten no ranges!")
     return
   }
}

func TestBadInitValue(t *testing.T) {
  var list = &[]SpanBoundry[int, string]{
    &Span[int, string]{Begin: 13, End: 11},
  }
  testDriver.Validate = true
  for range testDriver.NewSpanOverlapAccumulator().SliceIterFactory(list) {
    t.Errorf("Should have gotten no ranges!")
    return
  }

}