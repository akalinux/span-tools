package st

import "testing"

// Validates the inital range of a list of ranges
func TestFirstRange(t *testing.T) {
	var src *[]SpanBoundry[int] = &[]SpanBoundry[int]{
		&Span[int]{Begin: 2, End: 2},
		&Span[int]{Begin: 0, End: 1},
	}
	var span = testDriver.FirstSpan(src)
	if span.GetBegin() != 0 || span.GetEnd() != 1 {
		t.Errorf("Invalid start range")
	}
}

func TestGetNextBegin(t *testing.T) {
	var src *[]SpanBoundry[int] = &[]SpanBoundry[int]{
		&Span[int]{Begin: 3, End: 3},
		&Span[int]{Begin: 2, End: 2},
		&Span[int]{Begin: 0, End: 1},
	}
	var check = testDriver.GetNextBegin(1, src)
	if check == nil {
		t.Errorf("Should have gotten 2 as our next return value")
		return
	}
	if *check != 2 {
		t.Errorf("Expected 2, got %d", *check)
		return
	}
	check = testDriver.GetNextBegin(3, src)
	if check != nil {
		t.Errorf("Should have gotten nil")
		return
	}
}

func TestGetNextEnd(t *testing.T) {
	var src *[]SpanBoundry[int] = &[]SpanBoundry[int]{
		&Span[int]{Begin: 3, End: 3},
		&Span[int]{Begin: 2, End: 2},
		&Span[int]{Begin: 0, End: 1},
	}
	var check = testDriver.GetNextEnd(1, src)
	if check == nil {
		t.Errorf("Should have gotten 2 as our next return value")
		return
	}
	if *check != 2 {
		t.Errorf("Expected 2, got %d", *check)
		return
	}
	check = testDriver.GetNextEnd(3, src)
	if check != nil {
		t.Errorf("Should have gotten nil")
		return
	}
}

func TestGetNextSpan(t *testing.T) {
	var src *[]SpanBoundry[int] = &[]SpanBoundry[int]{
		&Span[int]{Begin: 0, End: 1},
		&Span[int]{Begin: 2, End: 2},
		&Span[int]{Begin: 3, End: 3},
	}
	var check = testDriver.NextSpan(&Span[int]{Begin: -1, End: 0}, src)
	var expected *[]SpanBoundry[int] = &[]SpanBoundry[int]{
		&Span[int]{Begin: 0, End: 1},
		&Span[int]{Begin: 1, End: 2},
		&Span[int]{Begin: 2, End: 3},
		&Span[int]{Begin: 3, End: 3},
	}

	for id,span := range *expected {
		if check == nil {
			t.Errorf("Should have gotten as our next return value")
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
		check = testDriver.NextSpan(check, src)
	}
	var nilCheck SpanBoundry[int]=nil
	if check == nilCheck {
		t.Errorf("Expected last call to NextSpan to return nil, got: %v",check)
		return
	}
}

func TestGetNextSpanReducedColumns(t *testing.T) {
	var src *[]SpanBoundry[int] = &[]SpanBoundry[int]{
		&Span[int]{Begin: 0, End: 1},
		&Span[int]{Begin: 2, End: 5},
		&Span[int]{Begin: 3, End: 6},
	}
	var expected []SpanBoundry[int] = []SpanBoundry[int]{
		&Span[int]{Begin: 0, End: 1},
		&Span[int]{Begin: 1, End: 2},
		&Span[int]{Begin: 2, End: 3},
		&Span[int]{Begin: 3, End: 5},
		&Span[int]{Begin: 5, End: 6},
	}
	var check = testDriver.NextSpan(&Span[int]{Begin: -1, End: 0}, src)
	for id, span := range expected {
		if check == nil {
			t.Errorf("check should not be nil for row %d", id)
			return
		}
		if span.GetBegin() != check.GetBegin() {
			t.Errorf("Invalid begin, expected: %d, got: %d on row: %d", span.GetBegin(), check.GetBegin(), id)
			return
		}
		if span.GetEnd() != check.GetEnd() {
			t.Errorf("Invalid end, expected: %d, got: %d on row: %d", span.GetEnd(), check.GetEnd(), id)
			return
		}
		check = testDriver.NextSpan(check, src)
	}
	if check == nil {
		t.Error("Expected last call to NextSpan to return nil")
	}
}

func TestGetNextSpanAllColumnsSetA(t *testing.T) {
	var src *[]SpanBoundry[int] = &[]SpanBoundry[int]{
		&Span[int]{Begin: 0, End: 7},
		&Span[int]{Begin: 2, End: 5},
		&Span[int]{Begin: 3, End: 6},
	}
	var expected []SpanBoundry[int] = []SpanBoundry[int]{
		&Span[int]{Begin: 0, End: 2},
		&Span[int]{Begin: 3, End: 5},
		&Span[int]{Begin: 5, End: 6},
		&Span[int]{Begin: 6, End: 7},
	}
	var check = testDriver.NextSpan(&Span[int]{Begin: -1, End: -1}, src)
	for id, span := range expected {
		if check == nil {
			t.Errorf("check should not be nil for row %d", id)
			return
		}
		if span.GetBegin() != check.GetBegin() {
			t.Errorf("Invalid begin, expected: %d->%d, got: %d->%d on row: %d",
				span.GetBegin(),
				span.GetEnd(), 
				check.GetBegin(),
				check.GetEnd(),
				id,
			)
			return
		}
		if span.GetEnd() != check.GetEnd() {
			t.Errorf("Invalid end, expected: %d, got: %d on row: %d", span.GetEnd(), check.GetEnd(), id)
			return
		}
		check = testDriver.NextSpan(check, src)
	}
	if check != nil {
		t.Error("Expected last call to NextSpan to return nil")
	}
}
