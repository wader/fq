package test

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
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

type StandardOS struct{}

func (StandardOS) Stdin() io.Reader                        { return os.Stdin }
func (StandardOS) Stdout() io.Writer                       { return os.Stdout }
func (StandardOS) Stderr() io.Writer                       { return os.Stderr }
func (StandardOS) Args() []string                          { return os.Args }
func (StandardOS) Open(name string) (io.ReadSeeker, error) { return os.Open(name) }

type testCaseRun struct {
	LineNr          int
	testCase        *testCase
	_Args           []string
	ExpectedStdout  string
	ActualStdoutBuf *bytes.Buffer
	ActualStderrBuf *bytes.Buffer
}

func (tcr *testCaseRun) Args() []string    { return tcr._Args }
func (tcr *testCaseRun) Stdin() io.Reader  { return nil } // TOOD: special file?
func (tcr *testCaseRun) Stdout() io.Writer { return tcr.ActualStdoutBuf }
func (tcr *testCaseRun) Stderr() io.Writer { return tcr.ActualStderrBuf }
func (tcr *testCaseRun) Open(name string) (io.ReadSeeker, error) {
	data, _ := tcr.testCase.Files[name]
	if len(data) == 0 {
		var err error
		log.Printf("%s filepath.Join(tcr.testCase.Path, name): %#+v\n", tcr.testCase.Path, filepath.Join(tcr.testCase.Path, name))
		f, err := os.Open(filepath.Join(tcr.testCase.Path, name))
		return f, err
	}

	return io.NewSectionReader(bytes.NewReader(data), 0, int64(len(data))), nil
}

type testCase struct {
	LineNr int

	Path string

	Files map[string][]byte
	Runs  []*testCaseRun

	ExpectedStdout string
	ExpectedStderr string
}

type section struct {
	LineNr int
	Name   string
	Value  string
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

			cs.LineNr = lineNr
			cs.Name = firstMatch(sm, func(s string) bool { return len(s) != 0 })
		} else {
			// TODO: use builder somehow if performance is needed
			cs.Value += l + lineDelim
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
		{LineNr: 1, Name: "a:", Value: "c\nc\n"},
		{LineNr: 4, Name: "b:", Value: ""},
		{LineNr: 5, Name: "a:", Value: "c\n"},
		{LineNr: 7, Name: "a:", Value: ""},
	}

	deepequal.Error(t, "sections", expectedSections, actualSections)
}

func parseTestCases(s string) *testCase {
	te := &testCase{}
	te.Files = map[string][]byte{}

	// match "name:" or ">args" sections
	seenRun := false
	for _, section := range sectionParser(regexp.MustCompile(`^#.*$|^/.*:|^>.*$`), s) {
		n, v := section.Name, section.Value

		switch {
		case !seenRun && strings.HasPrefix(n, "/"):
			name := n[1 : len(n)-1]
			te.Files[name] = []byte(v)
		case strings.HasPrefix(n, ">"):
			seenRun = true
			args := append([]string{"fq"}, strings.Fields(strings.TrimPrefix(n, ">"))...)
			te.Runs = append(te.Runs, &testCaseRun{
				LineNr:          section.LineNr,
				testCase:        te,
				_Args:           args,
				ExpectedStdout:  v,
				ActualStdoutBuf: &bytes.Buffer{},
				ActualStderrBuf: &bytes.Buffer{},
			})
		default:
			panic(fmt.Sprintf("%d: unexpected section %q %q", section.LineNr, n, v))
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
	log.Printf("err: %#+v\n", err)

	//log.Printf("tcr.ActualStdoutBuf.String(): %#+v\n", tcr.ActualStdoutBuf.String())

	// cli.Command{Version: "test", OS: &te}.Run()
	// deepequal.Error(t, "files", te.ExpectedFiles, te.ActualFiles)
	deepequal.Error(t, "stdout", tcr.ExpectedStdout, tcr.ActualStdoutBuf.String())
	//deepequal.Error(t, "stderr", te.ExpectedStderr, te.ActualStderrBuf.String())
}

func TestPath(t *testing.T, registry *decode.Registry) {
	const testDataDir = "testdata"

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
			tc.Path = filepath.Dir(path) // TODO: move?
			log.Printf("tc.Path: %#+v\n", tc.Path)

			for _, tcr := range tc.Runs {
				t.Run(strconv.Itoa(tcr.LineNr), func(t *testing.T) {
					testDecodedTestCaseRun(t, registry, tcr)
				})
			}
		})

		return nil
	})
}
