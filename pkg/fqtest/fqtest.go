package fqtest

import (
	"bytes"
	"encoding/hex"
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
	"fq/internal/shquote"
	"fq/pkg/bitio"
	"fq/pkg/cli"
	"fq/pkg/decode"
)

type testCaseRun struct {
	lineNr          int
	testCase        *testCase
	args            string
	expectedStdout  string
	actualStdoutBuf *bytes.Buffer
	actualStderrBuf *bytes.Buffer
}

func (tcr *testCaseRun) Stdin() io.Reader  { return nil } // TOOD: special file?
func (tcr *testCaseRun) Stdout() io.Writer { return tcr.actualStdoutBuf }
func (tcr *testCaseRun) Stderr() io.Writer { return tcr.actualStderrBuf }
func (tcr *testCaseRun) Environ() []string { return nil }
func (tcr *testCaseRun) Args() []string {
	return append([]string{"fq"}, shquote.Split(tcr.args)...)
}
func (tcr *testCaseRun) Open(name string) (io.ReadSeeker, error) {
	for _, p := range tcr.testCase.parts {
		f, ok := p.(*testCaseFile)
		if ok && f.name == name {
			// if no data assume it's a real file
			if len(f.data) == 0 {
				return os.Open(filepath.Join(filepath.Dir(tcr.testCase.path), name))
			}
			return io.NewSectionReader(bytes.NewReader(f.data), 0, int64(len(f.data))), nil
		}
	}
	return nil, fmt.Errorf("%s: file not found", name)
}

type testCaseFile struct {
	name string
	data []byte
}

type testCaseComment struct {
	comment string
}

type testCase struct {
	lineNr int
	path   string
	parts  []interface{}
}

func (tc *testCase) ToActual() string {
	sb := &strings.Builder{}
	for _, p := range tc.parts {

		switch p := p.(type) {
		case *testCaseComment:
			fmt.Fprintf(sb, "#%s\n", p.comment)
		case *testCaseRun:
			fmt.Fprintf(sb, ">%s\n", p.args)
			fmt.Fprint(sb, p.actualStdoutBuf.String())
		case *testCaseFile:
			fmt.Fprintf(sb, "/%s:\n", p.name)
			sb.Write(p.data)
		default:
			panic("unreachable")
		}
	}
	return sb.String()
}

type Section struct {
	LineNr int
	Name   string
	Value  string
}

var unescapeRe = regexp.MustCompile(`\\(?:b[01]+|x[0-f]+)`)

func Unescape(s string) string {
	return unescapeRe.ReplaceAllStringFunc(s, func(r string) string {
		log.Printf("r: %s\n", r)
		switch {
		case r[1] == 'b':
			b, _ := bitio.BytesFromBitString(r[2:])
			return string(b)
		case r[1] == 'x':
			b, _ := hex.DecodeString(r[2:])
			return string(b)
		default:
			return r
		}
	})
}

func SectionParser(re *regexp.Regexp, s string) []Section {
	var sections []Section

	firstMatch := func(ss []string, fn func(s string) bool) string {
		for _, s := range ss {
			if fn(s) {
				return s
			}
		}
		return ""
	}

	const lineDelim = "\n"
	var cs *Section
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
			sections = append(sections, Section{})
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

func parseTestCases(s string) *testCase {
	te := &testCase{}
	te.parts = []interface{}{}

	for _, section := range SectionParser(regexp.MustCompile(`^#.*$|^/.*:|^>.*$`), s) {
		n, v := section.Name, section.Value

		switch {
		case strings.HasPrefix(n, "#"):
			comment := n[1:]
			te.parts = append(te.parts, &testCaseComment{comment: comment})
		case strings.HasPrefix(n, "/"):
			name := n[1 : len(n)-1]
			te.parts = append(te.parts, &testCaseFile{name: name, data: []byte(v)})
		case strings.HasPrefix(n, ">"):
			te.parts = append(te.parts, &testCaseRun{
				lineNr:          section.LineNr,
				testCase:        te,
				args:            strings.TrimPrefix(n, ">"),
				expectedStdout:  v,
				actualStdoutBuf: &bytes.Buffer{},
				actualStderrBuf: &bytes.Buffer{},
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

			for _, p := range tc.parts {
				tcr, ok := p.(*testCaseRun)
				if !ok {
					continue
				}

				t.Run(strconv.Itoa(tcr.lineNr)+":"+tcr.args, func(t *testing.T) {
					testDecodedTestCaseRun(t, registry, tcr)
				})
			}
		})

		return nil
	})

	if v := os.Getenv("WRITE_ACTUAL"); v != "" {
		for _, tc := range tcs {
			if err := ioutil.WriteFile(tc.path, []byte(tc.ToActual()), 0644); err != nil {
				t.Error(err)
			}
		}
	}
}
