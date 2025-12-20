package st

import "testing"

// Validates the inital range of a list of ranges
func TestFirstRange(t *testing.T) {
	var src *[]SpanBoundry[int] = &[]SpanBoundry[int]{
		&Span[int]{Begin: 2, End: 2},
		&Span[int]{Begin: 0, End: 1},
	}
	var span,_ = testDriver.FirstSpan(src)
	if span.GetBegin() != 0 || span.GetEnd() != 1 {
		t.Errorf("Invalid start range")
	}
}


func CommonNextSpan(src *[]SpanBoundry[int],expected *[]SpanBoundry[int],t *testing.T) {

	var check,ok = testDriver.NextSpan(&Span[int]{Begin: -1, End: -1}, src)
	for id,span := range *expected {
		if !ok {
			t.Errorf("Should have gotten as our next return value for id: %d",id)
			return
		}
		if span.GetBegin() != check.GetBegin() {
			t.Errorf("Invalid begin, expected: %v, got: %v for set: %d", span, check,id)
			return
		}
		if span.GetEnd() != check.GetEnd() {
			t.Errorf("Invalid end, expected: %v, got: %v for set: %d", span, check,id)
			return
		}
		check,ok = testDriver.NextSpan(check, src)
	}
	if ok {
		t.Errorf("Expected last call to NextSpan to return nil, got: %v",check)
		return
	}	
}

func TestGetNextSpan(t *testing.T) {
	var src *[]SpanBoundry[int] = &[]SpanBoundry[int]{
		&Span[int]{Begin: 3, End: 3},
		&Span[int]{Begin: 2, End: 2},
		&Span[int]{Begin: 0, End: 1},
	}
	var expected *[]SpanBoundry[int] = &[]SpanBoundry[int]{
		&Span[int]{Begin: 0, End: 1},
		&Span[int]{Begin: 2, End: 2},
		&Span[int]{Begin: 3, End: 3},
	}
	CommonNextSpan(src,expected,t)
}

func TestGetNextSpanReducedColumns(t *testing.T) {
	var src *[]SpanBoundry[int] = &[]SpanBoundry[int]{
		&Span[int]{Begin: 0, End: 1},
		&Span[int]{Begin: 3, End: 6},
		&Span[int]{Begin: 2, End: 5},
	}
	var expected *[]SpanBoundry[int] = &[]SpanBoundry[int]{
		&Span[int]{Begin: 0, End: 1},
		&Span[int]{Begin: 2, End: 3},
		&Span[int]{Begin: 4, End: 5},
		&Span[int]{Begin: 6, End: 6},
	}
	CommonNextSpan(src,expected,t)
}

func TestGetNextSpanAllColumnsSetA(t *testing.T) {
	var src *[]SpanBoundry[int] = &[]SpanBoundry[int]{
		&Span[int]{Begin: 0, End: 1},
		&Span[int]{Begin: 2, End: 5},
		&Span[int]{Begin: 3, End: 4},
		&Span[int]{Begin: 3, End: 6},
	}
	var expected *[]SpanBoundry[int] = &[]SpanBoundry[int]{
		&Span[int]{Begin: 0, End: 1},
		&Span[int]{Begin: 2, End: 3},
		&Span[int]{Begin: 4, End: 4},
		&Span[int]{Begin: 5, End: 5},
		&Span[int]{Begin: 6, End: 6},
	}

	CommonNextSpan(src,expected,t)
}

func TestGetNextSpanAllColumnsSetB(t *testing.T) {
	var src *[]SpanBoundry[int] = &[]SpanBoundry[int]{
		&Span[int]{Begin: 0, End: 1},
		&Span[int]{Begin: 2, End: 5},
		&Span[int]{Begin: 3, End: 4},
		&Span[int]{Begin: 4, End: 6},
	}
	var expected *[]SpanBoundry[int] = &[]SpanBoundry[int]{
		&Span[int]{Begin: 0, End: 1},
		&Span[int]{Begin: 2, End: 3},
		&Span[int]{Begin: 4, End: 4},
		&Span[int]{Begin: 5, End: 5},
		&Span[int]{Begin: 6, End: 6},
	}

	CommonNextSpan(src,expected,t)
}
