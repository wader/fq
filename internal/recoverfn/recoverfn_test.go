package recoverfn_test

import (
	"reflect"
	"testing"

	"github.com/wader/fq/internal/recoverfn"
)

func test1() {
	panic(testError(true))
}

func test2() {
	test1()
}

type testError bool

func (t testError) IsRecoverableError() bool { return bool(t) }

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

	if _, ok := r.RecoverV.(testError); !ok {
		t.Errorf("expected v %v, got %v", testError(true), r.RecoverV)
	}

	frames := r.Frames()

	expectedFramesLen := 2
	actualFramesLen := len(frames)
	if !reflect.DeepEqual(expectedFramesLen, actualFramesLen) {
		t.Errorf("expected len(frames) %v, got %v", actualFramesLen, actualFramesLen)
	}

	expectedFrame0Function := "github.com/wader/fq/internal/recoverfn_test.test1"
	actualFrame0Function := frames[0].Function
	if !reflect.DeepEqual(expectedFrame0Function, actualFrame0Function) {
		t.Errorf("expected frames[0].Function %v, got %v", expectedFrame0Function, actualFrame0Function)
	}

	expectedFrame1Function := "github.com/wader/fq/internal/recoverfn_test.test2"
	actualFrame1Function := frames[1].Function
	if !reflect.DeepEqual(expectedFrame1Function, actualFrame1Function) {
		t.Errorf("expected frames[1].Function %v, got %v", expectedFrame1Function, actualFrame1Function)
	}
}
