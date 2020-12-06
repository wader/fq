package deepequal

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"reflect"
	"strconv"
	"strings"
)

type tf interface {
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
}

func diff(a, b string, nLineContext int) string {
	aFile, _ := ioutil.TempFile("", "naivediff")
	defer os.Remove(aFile.Name())
	_, _ = io.Copy(aFile, bytes.NewBufferString(a))
	aFile.Close()
	bFile, _ := ioutil.TempFile("", "naivediff")
	defer os.Remove(bFile.Name())
	_, _ = io.Copy(bFile, bytes.NewBufferString(b))
	bFile.Close()
	c := exec.Command("diff", "-U", strconv.Itoa(nLineContext), aFile.Name(), bFile.Name())
	diffBuf, _ := c.Output()
	realDiff := strings.Join(strings.Split(string(diffBuf), "\n")[2:], "\n")
	return realDiff
}

func testDeepEqual(fn func(format string, args ...interface{}), name string, expected interface{}, actual interface{}) {
	expectedStr := fmt.Sprintf("%s", expected)
	actualStr := fmt.Sprintf("%s", actual)

	if !reflect.DeepEqual(expected, actual) {
		fn(`
%s diff:
%s
`,
			name, diff(expectedStr, actualStr, 5))
	}
}

func Error(t tf, name string, expected interface{}, actual interface{}) {
	testDeepEqual(t.Errorf, name, expected, actual)
}

func Fatal(t tf, name string, expected interface{}, actual interface{}) {
	testDeepEqual(t.Fatalf, name, expected, actual)
}
