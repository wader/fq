package fqtest

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"testing"

	"fq/format/registry"
	"fq/internal/deepequal"
	"fq/internal/shquote"
	"fq/pkg/bitio"
	"fq/pkg/interp"
)

type testCaseReadline struct {
	input          string
	expectedPrompt string
	expectedStdout string
}

type testCaseRunOutput struct {
	io.Writer
}

func (o testCaseRunOutput) Size() (int, int) { return 130, 25 }
func (o testCaseRunOutput) IsTerminal() bool { return true }

type testCaseRun struct {
	lineNr           int
	testCase         *testCase
	args             string
	stdin            string
	expectedStdout   string
	expectedStderr   string
	expectedExitCode int
	actualStdoutBuf  *bytes.Buffer
	actualStderrBuf  *bytes.Buffer
	actualExitCode   int
	readlines        []testCaseReadline
	readlinesPos     int
}

func (tcr *testCaseRun) Line() int { return tcr.lineNr }

func (tcr *testCaseRun) Stdin() fs.File {
	return interp.FileReader{R: bytes.NewBufferString(tcr.stdin)}
}

func (tcr *testCaseRun) Stdout() interp.Output { return testCaseRunOutput{tcr.actualStdoutBuf} }

func (tcr *testCaseRun) Stderr() io.Writer { return tcr.actualStderrBuf }

func (tcr *testCaseRun) Interrupt() chan struct{} { return nil }

func (tcr *testCaseRun) Environ() []string {
	return []string{
		"NODECODEPROGRESS=1",
	}
}

func (tcr *testCaseRun) Args() []string { return shquote.Split(tcr.args) }

func (tcr *testCaseRun) ConfigDir() (string, error) { return "/config", nil }

func (tcr *testCaseRun) FS() fs.FS { return tcr.testCase }

func (tcr *testCaseRun) Readline(prompt string, complete func(line string, pos int) (newLine []string, shared int)) (string, error) {
	tcr.actualStdoutBuf.WriteString(prompt)
	if tcr.readlinesPos >= len(tcr.readlines) {
		return "", io.EOF
	}

	lineRaw := tcr.readlines[tcr.readlinesPos].input
	line := Unescape(lineRaw)
	tcr.readlinesPos++

	if strings.HasSuffix(line, "\t") {
		tcr.actualStdoutBuf.WriteString(lineRaw + "\n")

		l := len(line) - 1
		newLine, shared := complete(line[0:l], l)
		// TODO: shared
		_ = shared
		for _, nl := range newLine {
			tcr.actualStdoutBuf.WriteString(nl + "\n")
		}

		return "", nil
	}

	tcr.actualStdoutBuf.WriteString(lineRaw + "\n")

	if line == "^D" {
		return "", io.EOF
	}

	return line, nil
}
func (tcr *testCaseRun) History() ([]string, error) { return nil, nil }

func (tcr *testCaseRun) ToExpectedStdout() string {
	sb := &strings.Builder{}

	if len(tcr.readlines) == 0 {
		fmt.Fprint(sb, tcr.expectedStdout)
	} else {
		for _, rl := range tcr.readlines {
			fmt.Fprintf(sb, "%s%s\n", rl.expectedPrompt, rl.input)
			if rl.expectedStdout != "" {
				fmt.Fprint(sb, rl.expectedStdout)
			}
		}
	}

	return sb.String()
}

func (tcr *testCaseRun) ToExpectedStderr() string {
	return tcr.expectedStderr
}

type part interface {
	Line() int
}

type testCaseFile struct {
	lineNr int
	name   string
	data   []byte
}

func (tcf *testCaseFile) Line() int { return tcf.lineNr }

type testCaseComment struct {
	lineNr  int
	comment string
}

func (tcr *testCaseComment) Line() int { return tcr.lineNr }

type testCase struct {
	path      string
	parts     []part
	wasTested bool
}

func (tc *testCase) ToActual() string {
	var partsLineSorted []part
	partsLineSorted = append(partsLineSorted, tc.parts...)
	sort.Slice(partsLineSorted, func(i, j int) bool {
		return partsLineSorted[i].Line() < partsLineSorted[j].Line()
	})

	sb := &strings.Builder{}
	for _, p := range partsLineSorted {
		switch p := p.(type) {
		case *testCaseComment:
			fmt.Fprintf(sb, "#%s\n", p.comment)
		case *testCaseRun:
			fmt.Fprintf(sb, "$%s\n", p.args)
			fmt.Fprint(sb, p.actualStdoutBuf.String())
			if p.actualExitCode != 0 {
				fmt.Fprintf(sb, "exitcode: %d\n", p.actualExitCode)
			}
			if p.stdin != "" {
				fmt.Fprint(sb, "stdin:\n")
				fmt.Fprint(sb, p.stdin)
			}
			if p.actualStderrBuf.Len() > 0 {
				fmt.Fprint(sb, "stderr:\n")
				fmt.Fprint(sb, p.actualStderrBuf.String())
			}
		case *testCaseFile:
			fmt.Fprintf(sb, "%s:\n", p.name)
			sb.Write(p.data)
		default:
			panic("unreachable")
		}
	}

	return sb.String()
}

func (tc *testCase) Open(name string) (fs.File, error) {
	for _, p := range tc.parts {
		f, ok := p.(*testCaseFile)
		if ok && f.name == name {
			// if no data assume it's a real file
			if len(f.data) == 0 {
				return os.Open(filepath.Join(filepath.Dir(tc.path), name))
			}
			return interp.FileReader{
				R: io.NewSectionReader(bytes.NewReader(f.data), 0, int64(len(f.data))),
				FileInfo: interp.FixedFileInfo{
					FName: filepath.Base(name),
					FSize: int64(len(f.data)),
				},
			}, nil
		}
	}
	return nil, fmt.Errorf("%s: file not found", name)
}

type Section struct {
	LineNr int
	Name   string
	Value  string
}

var unescapeRe = regexp.MustCompile(`\\(?:t|b|n|r|0(?:b[01]{8}|x[0-f]{2}))`)

func Unescape(s string) string {
	return unescapeRe.ReplaceAllStringFunc(s, func(r string) string {
		switch {
		case r == `\n`:
			return "\n"
		case r == `\r`:
			return "\r"
		case r == `\t`:
			return "\t"
		case r == `\b`:
			return "\b"
		case strings.HasPrefix(r, `\0b`):
			b, _ := bitio.BytesFromBitString(r[3:])
			return string(b)
		case strings.HasPrefix(r, `\0x`):
			b, _ := hex.DecodeString(r[3:])
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
	te.parts = []part{}
	var currentTestRun *testCaseRun
	const promptEnd = "> "
	replDepth := 0

	// TODO: better section splitter, too much heuristics now
	for _, section := range SectionParser(regexp.MustCompile(`^\$ .*$|^stdin:$|^stderr:$|^exitcode:.*$|^#.*$|^/.*:|^(?:>* )?[a-z\d]+> .*$`), s) {
		n, v := section.Name, section.Value

		switch {
		case strings.HasPrefix(n, "#"):
			comment := n[1:]
			te.parts = append(te.parts, &testCaseComment{lineNr: section.LineNr, comment: comment})
		case strings.HasPrefix(n, "/"):
			name := n[0 : len(n)-1]
			te.parts = append(te.parts, &testCaseFile{lineNr: section.LineNr, name: name, data: []byte(v)})
		case strings.HasPrefix(n, "$"):
			replDepth++

			if currentTestRun != nil {
				te.parts = append(te.parts, currentTestRun)
			}

			currentTestRun = &testCaseRun{
				lineNr:          section.LineNr,
				testCase:        te,
				args:            strings.TrimPrefix(n, "$"),
				expectedStdout:  v,
				actualStdoutBuf: &bytes.Buffer{},
				actualStderrBuf: &bytes.Buffer{},
			}
		case strings.HasPrefix(n, "exitcode:"):
			currentTestRun.expectedExitCode, _ = strconv.Atoi(strings.TrimSpace(strings.TrimPrefix(n, "exitcode:")))
		case strings.HasPrefix(n, "stdin"):
			currentTestRun.stdin = v
		case strings.HasPrefix(n, "stderr"):
			currentTestRun.expectedStderr = v
		case strings.Contains(n, promptEnd): // TODO: better
			i := strings.LastIndex(n, promptEnd)

			prompt := n[0:i] + promptEnd
			input := n[i+2:]

			currentTestRun.readlines = append(currentTestRun.readlines, testCaseReadline{
				input:          input,
				expectedPrompt: prompt,
				expectedStdout: v,
			})

			// TODO: hack
			if strings.Contains(input, "| repl") {
				replDepth++
			}
			if input == "^D" {
				replDepth--
			}

			if replDepth == 0 {
				te.parts = append(te.parts, currentTestRun)
				currentTestRun = nil
			}

		default:
			panic(fmt.Sprintf("%d: unexpected section %q %q", section.LineNr, n, v))
		}
	}

	if currentTestRun != nil {
		te.parts = append(te.parts, currentTestRun)
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

func testDecodedTestCaseRun(t *testing.T, registry *registry.Registry, tcr *testCaseRun) {
	q, err := interp.New(tcr, registry)
	if err != nil {
		t.Fatal(err)
	}

	err = q.Main(context.Background(), tcr.Stdout(), "dev")
	if err != nil {
		if ex, ok := err.(interp.Exiter); ok { //nolint:errorlint
			tcr.actualExitCode = ex.ExitCode()
		}
	}

	// cli.Command{Version: "test", OS: &te}.Run()
	// deepequal.Error(t, "files", te.ExpectedFiles, te.ActualFiles)
	deepequal.Error(t, "exitcode", tcr.expectedExitCode, tcr.actualExitCode)
	deepequal.Error(t, "stdout", tcr.ToExpectedStdout(), tcr.actualStdoutBuf.String())
	deepequal.Error(t, "stderr", tcr.ToExpectedStderr(), tcr.actualStderrBuf.String())
}

func TestPath(t *testing.T, registry *registry.Registry) {
	tcs := []*testCase{}

	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
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
					tc.wasTested = true
				})
			}
		})

		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	if v := os.Getenv("WRITE_ACTUAL"); v != "" {
		for _, tc := range tcs {
			if !tc.wasTested {
				continue
			}
			if err := ioutil.WriteFile(tc.path, []byte(tc.ToActual()), 0644); err != nil { //nolint:gosec
				t.Error(err)
			}
		}
	}
}
