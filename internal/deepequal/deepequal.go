package deepequal

import (
	"fmt"
	"reflect"

	"github.com/pmezard/go-difflib/difflib"
)

type tf interface {
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
}

func testDeepEqual(fn func(format string, args ...interface{}), name string, expected interface{}, actual interface{}) {
	expectedStr := fmt.Sprintf("%v", expected)
	actualStr := fmt.Sprintf("%v", actual)

	if !reflect.DeepEqual(expected, actual) {
		diff := difflib.UnifiedDiff{
			A:        difflib.SplitLines(expectedStr),
			B:        difflib.SplitLines(actualStr),
			FromFile: "expected",
			ToFile:   "actual",
			Context:  3,
		}
		uDiff, err := difflib.GetUnifiedDiffString(diff)
		if err != nil {
			panic(err)
		}
		fn(`
%s diff:
%s
`,
			name, uDiff)
	}
}

func Error(t tf, name string, expected interface{}, actual interface{}) {
	testDeepEqual(t.Errorf, name, expected, actual)
}

func Fatal(t tf, name string, expected interface{}, actual interface{}) {
	testDeepEqual(t.Fatalf, name, expected, actual)
}
