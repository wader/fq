// Package difftest implement test based on serialized string output
//
// User provides a function that gets a input string. It returns a output string
// based on the input somehow and a output path to file with content to compare it
// to or to write to if in write mode.
// If there is a difference test will fail with a diff.
//
// Test inputs are read from files matching Pattern from Path.
//
// Note that output path can be the same as input which useful if the function
// implements some kind of transcript that includes both input and output.
package difftest

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pmezard/go-difflib/difflib"
)

type tf interface {
	Helper()
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
}

const green = "\x1b[32m"
const red = "\x1b[31m"
const reset = "\x1b[0m"

type Fn func(t *testing.T, path string, input string) (string, string, error)

type Options struct {
	Path        string
	Pattern     string
	ColorDiff   bool
	WriteOutput bool
	Fn          Fn
}

func testDeepEqual(t tf, color bool, printfFn func(format string, args ...interface{}), expected string, actual string) {
	t.Helper()

	diff := difflib.UnifiedDiff{
		A:        difflib.SplitLines(expected),
		B:        difflib.SplitLines(actual),
		FromFile: "expected",
		ToFile:   "actual",
		Context:  3,
	}
	uDiff, err := difflib.GetUnifiedDiffString(diff)
	if err != nil {
		t.Fatalf("%s", err)
	}
	if uDiff == "" {
		return
	}

	if color {
		lines := strings.Split(uDiff, "\n")
		var coloredLines []string
		// diff looks like this:
		// --- expected
		// +++ actual
		// @@ -5,7 +5,7 @@
		// -            a
		// +            b
		seenAt := false
		for _, l := range lines {
			if len(l) == 0 {
				continue
			}
			switch {
			case seenAt && l[0] == '+':
				coloredLines = append(coloredLines, green+l+reset)
			case seenAt && l[0] == '-':
				coloredLines = append(coloredLines, red+l+reset)
			default:
				if l[0] == '@' {
					seenAt = true
				}
				coloredLines = append(coloredLines, l)
			}
		}
		uDiff = strings.Join(coloredLines, "\n")
	}

	printfFn("%s", "\n"+uDiff)
}

func ErrorEx(t tf, color bool, expected string, actual string) {
	t.Helper()
	testDeepEqual(t, color, t.Errorf, expected, actual)
}

func Error(t tf, expected string, actual string) {
	t.Helper()
	testDeepEqual(t, false, t.Errorf, expected, actual)
}

func FatalEx(t tf, color bool, expected string, actual string) {
	t.Helper()
	testDeepEqual(t, color, t.Fatalf, expected, actual)
}

func Fatal(t tf, expected string, actual string) {
	t.Helper()
	testDeepEqual(t, false, t.Fatalf, expected, actual)
}

func TestWithOptions(t *testing.T, opts Options) {
	t.Helper()

	func() {
		t.Helper()

		// done in two steps as it seems hard to mark some functions inside filepath.Walk as test helpers
		var paths []string
		if err := filepath.Walk(opts.Path, func(path string, info os.FileInfo, err error) error {
			t.Helper()

			if err != nil {
				return err
			}
			match, err := filepath.Match(filepath.Join(filepath.Dir(path), opts.Pattern), path)
			if err != nil {
				return err
			} else if !match {
				return nil
			}

			paths = append(paths, path)
			return nil
		}); err != nil {
			t.Fatal(err)
		}

		for _, p := range paths {
			t.Run(p, func(t *testing.T) {
				t.Helper()
				input, err := ioutil.ReadFile(p)
				if err != nil {
					t.Fatal(err)
				}

				outputPath, output, err := opts.Fn(t, p, string(input))
				if err != nil {
					t.Fatal(err)
				}

				expectedOutput, expectedOutputErr := ioutil.ReadFile(outputPath)
				if opts.WriteOutput {
					if expectedOutputErr == nil && string(expectedOutput) == output {
						return
					}

					if err := ioutil.WriteFile(outputPath, []byte(output), 0644); err != nil { //nolint:gosec
						t.Fatal(err)
					}
					return
				}

				if expectedOutputErr != nil {
					t.Fatal(expectedOutputErr)
				}

				ErrorEx(t, opts.ColorDiff, string(expectedOutput), output)
			})
		}
	}()
}

func Test(t *testing.T, pattern string, fn Fn) {
	t.Helper()

	TestWithOptions(t, Options{
		Path:    "testdata",
		Pattern: pattern,
		Fn:      fn,
	})
}
