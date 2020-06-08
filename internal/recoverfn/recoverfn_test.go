package recoverfn_test

import (
	"fq/internal/recoverfn"
	"reflect"
	"testing"
)

func test1() {
	panic("hello")
}

func test2() {
	test1()
}

func TestNormal(t *testing.T) {
	_, rOk := recoverfn.Run(func() {})
	expectedROK := true
	if !reflect.DeepEqual(expectedROK, rOk) {
		t.Errorf("expected v %v, got %v", expectedROK, rOk)
	}
}

func TestPanic(t *testing.T) {
	r, rOk := recoverfn.Run(test2)

	expectedROK := false
	if !reflect.DeepEqual(expectedROK, rOk) {
		t.Errorf("expected v %v, got %v", expectedROK, rOk)
	}

	expectedV := "hello"
	if !reflect.DeepEqual(expectedV, r.RecoverV) {
		t.Errorf("expected v %v, got %v", expectedV, r.RecoverV)
	}

	frames := r.Frames()

	expectedFramesLen := 2
	actualFramesLen := len(frames)
	if !reflect.DeepEqual(expectedFramesLen, actualFramesLen) {
		t.Errorf("expected len(frames) %v, got %v", actualFramesLen, actualFramesLen)
	}

	expectedFrame0Function := "fq/internal/recoverfn_test.test1"
	actualFrame0Function := frames[0].Function
	if !reflect.DeepEqual(expectedFrame0Function, actualFrame0Function) {
		t.Errorf("expected frames[0].Function %v, got %v", expectedFrame0Function, actualFrame0Function)
	}

	expectedFrame1Function := "fq/internal/recoverfn_test.test2"
	actualFrame1Function := frames[1].Function
	if !reflect.DeepEqual(expectedFrame1Function, actualFrame1Function) {
		t.Errorf("expected frames[1].Function %v, got %v", expectedFrame1Function, actualFrame1Function)
	}
}
