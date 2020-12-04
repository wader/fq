package deepequal

import (
	"fmt"
	"reflect"
)

type tf interface {
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
}

func testDeepEqual(fn func(format string, args ...interface{}), name string, expected interface{}, actual interface{}) {
	expectedStr := fmt.Sprintf("%#v", expected)
	actualStr := fmt.Sprintf("%#v", actual)
	if !reflect.DeepEqual(expected, actual) {
		diff := ""
		for i := len(diff); i < len(expectedStr) && i < len(actualStr); i++ {
			if expectedStr[i] != actualStr[i] {
				diff += "^"
			} else {
				diff += " "
			}
		}
		fn(`
%s
expected: %s
  actual: %s
    diff: %s`[1:],
			name, expectedStr, actualStr, diff)
	}
}

func Error(t tf, name string, expected interface{}, actual interface{}) {
	testDeepEqual(t.Errorf, name, expected, actual)
}

func Fatal(t tf, name string, expected interface{}, actual interface{}) {
	testDeepEqual(t.Fatalf, name, expected, actual)
}
