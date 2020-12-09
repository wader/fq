package test

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"fq/internal/deepequal"
	"fq/pkg/cli"
	"fq/pkg/decode"
)

type testCaseRun struct {
	lineNr          int
	testCase        *testCase
	args            []string
	expectedStdout  string
	actualStdoutBuf *bytes.Buffer
	actualStderrBuf *bytes.Buffer
}

func (tcr *testCaseRun) Args() []string    { return append([]string{"fq"}, tcr.args...) }
func (tcr *testCaseRun) Stdin() io.Reader  { return nil } // TOOD: special file?
func (tcr *testCaseRun) Stdout() io.Writer { return tcr.actualStdoutBuf }
func (tcr *testCaseRun) Stderr() io.Writer { return tcr.actualStderrBuf }
func (tcr *testCaseRun) Open(name string) (io.ReadSeeker, error) {
	for _, f := range tcr.testCase.files {
		if f.name == name {
			// if no data assume it's a real file
			if len(f.data) == 0 {
				return os.Open(filepath.Join(filepath.Dir(tcr.testCase.path), name))
			}
			return io.NewSectionReader(bytes.NewReader(f.data), 0, int64(len(f.data))), nil
		}
	}
	return nil, fmt.Errorf("%s: file not found", name)
}

type testFile struct {
	name string
	data []byte
}

type testCase struct {
	lineNr         int
	path           string
	files          []testFile
	runs           []*testCaseRun
	expectedStdout string
	expectedStderr string
}

func (tc *testCase) ToActual() string {
	sb := &strings.Builder{}

	for _, f := range tc.files {
		fmt.Fprintf(sb, "/%s:\n", f.name)
		sb.Write(f.data)
	}
	for _, r := range tc.runs {
		fmt.Fprintf(sb, "> %s\n", strings.Join(r.args, " "))
		fmt.Fprintf(sb, r.actualStdoutBuf.String())
	}

	return sb.String()
}

type section struct {
	lineNr int
	name   string
	value  string
}

func sectionParser(re *regexp.Regexp, s string) []section {
	var sections []section

	firstMatch := func(ss []string, fn func(s string) bool) string {
		for _, s := range ss {
			if fn(s) {
				return s
			}
		}
		return ""
	}

	const lineDelim = "\n"
	var cs *section
	lineNr := 0
	lines := strings.Split(s, lineDelim)
	// skip last if empty because of how split works "a\n" -> ["a", ""]
	if lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	for _, l := range lines {
		lineNr++

		sm := re.FindStringSubmatch(l)
		if cs == nil || len(sm) > 0 {
			sections = append(sections, section{})
			cs = &sections[len(sections)-1]

			cs.lineNr = lineNr
			cs.name = firstMatch(sm, func(s string) bool { return len(s) != 0 })
		} else {
			// TODO: use builder somehow if performance is needed
			cs.value += l + lineDelim
		}

	}

	return sections
}

func TestSectionParser(t *testing.T) {
	actualSections := sectionParser(
		regexp.MustCompile(`^(?:(a:)|(b:))$`),
		`
a:
c
c
b:
a:
c
a:
`[1:])

	expectedSections := []section{
		{lineNr: 1, name: "a:", value: "c\nc\n"},
		{lineNr: 4, name: "b:", value: ""},
		{lineNr: 5, name: "a:", value: "c\n"},
		{lineNr: 7, name: "a:", value: ""},
	}

	deepequal.Error(t, "sections", expectedSections, actualSections)
}

func parseTestCases(s string) *testCase {
	te := &testCase{}
	te.files = []testFile{}

	// match "name:" or ">args" sections
	seenRun := false
	for _, section := range sectionParser(regexp.MustCompile(`^#.*$|^/.*:|^>.*$`), s) {
		n, v := section.name, section.value

		switch {
		case !seenRun && strings.HasPrefix(n, "/"):
			name := n[1 : len(n)-1]
			te.files = append(te.files, testFile{name: name, data: []byte(v)})
		case strings.HasPrefix(n, ">"):
			seenRun = true
			args := strings.Fields(strings.TrimPrefix(n, ">"))
			te.runs = append(te.runs, &testCaseRun{
				lineNr:          section.lineNr,
				testCase:        te,
				args:            args,
				expectedStdout:  v,
				actualStdoutBuf: &bytes.Buffer{},
				actualStderrBuf: &bytes.Buffer{},
			})
		default:
			panic(fmt.Sprintf("%d: unexpected section %q %q", section.lineNr, n, v))
		}
	}

	return te
}

// func TestParseTestCase(t *testing.T) {
// 	actualTestCase := parseTestCases(`
// /a:
// input content a
// $ a b
// /a:
// expected content a
// >stdout:
// expected stdout
// >stderr:
// expected stderr
// ---
// /a2:
// input content a2
// $ a2 b2
// /a2:
// expected content a2
// >stdout:
// expected stdout2
// >stderr:
// expected stderr2
// `[1:])

// 	expectedTestCase := []testCase{
// 		{
// 			RunArgs:         []string{"a", "b"},
// 			Files:           map[string]string{"a": "input content a\n"},
// 			ExpectedFiles:   map[string]string{"a": "expected content a\n"},
// 			ExpectedStdout:  "expected stdout\n",
// 			ExpectedStderr:  "expected stderr\n",
// 			ActualStdoutBuf: &bytes.Buffer{},
// 			ActualStderrBuf: &bytes.Buffer{},
// 			ActualFiles:     map[string]string{},
// 		},
// 		{
// 			RunArgs:         []string{"a2", "b2"},
// 			Files:           map[string]string{"a2": "input content a2\n"},
// 			ExpectedFiles:   map[string]string{"a2": "expected content a2\n"},
// 			ExpectedStdout:  "expected stdout2\n",
// 			ExpectedStderr:  "expected stderr2\n",
// 			ActualStdoutBuf: &bytes.Buffer{},
// 			ActualStderrBuf: &bytes.Buffer{},
// 			ActualFiles:     map[string]string{},
// 		},
// 	}

// 	deepequal.Error(t, "testcase", expectedTestCase, actualTestCase)
// }

func testDecodedTestCaseRun(t *testing.T, registry *decode.Registry, tcr *testCaseRun) {
	//log.Printf("tcr: %#+v\n", tcr)

	m := cli.Main{
		OS:       tcr,
		Registry: registry,
	}
	err := m.Run()
	if err != nil {
		// TODO: expect error
		t.Fatal(err)
	}

	//log.Printf("tcr.ActualStdoutBuf.String(): %#+v\n", tcr.ActualStdoutBuf.String())

	// cli.Command{Version: "test", OS: &te}.Run()
	// deepequal.Error(t, "files", te.ExpectedFiles, te.ActualFiles)
	deepequal.Error(t, "stdout", tcr.expectedStdout, tcr.actualStdoutBuf.String())
	//deepequal.Error(t, "stderr", te.ExpectedStderr, te.ActualStderrBuf.String())
}

func TestPath(t *testing.T, registry *decode.Registry) {
	const testDataDir = "testdata"

	tcs := []*testCase{}

	filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) != ".fqtest" {
			return nil
		}

		t.Run(path, func(t *testing.T) {
			b, err := ioutil.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			tc := parseTestCases(string(b))
			tcs = append(tcs, tc)
			tc.path = path

			for _, tcr := range tc.runs {
				t.Run(strconv.Itoa(tcr.lineNr), func(t *testing.T) {
					testDecodedTestCaseRun(t, registry, tcr)
				})
			}
		})

		return nil
	})

	for _, tc := range tcs {
		if err := ioutil.WriteFile(tc.path, []byte(tc.ToActual()), 0644); err != nil {
			t.Error(err)
		}
	}
}
