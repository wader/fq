package script

import (
	"bytes"
	"cmp"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/wader/fq/internal/shquote"
	"github.com/wader/fq/pkg/bitio"
	"github.com/wader/fq/pkg/interp"
)

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

var escapeRe = regexp.MustCompile(`[^[:print:][:space:]]`)

func Escape(s string) string {
	return string(escapeRe.ReplaceAllFunc([]byte(s), func(r []byte) []byte {
		return []byte(fmt.Sprintf(`\0x%.2x`, r[0]))
	}))
}

type CaseReadline struct {
	expr           string
	env            []string
	input          string
	expectedPrompt string
	expectedStdout string
}

type CaseRunInput struct {
	interp.FileReader
	isTerminal bool
	width      int
	height     int
}

func (i CaseRunInput) Size() (int, int) { return i.width, i.height }
func (i CaseRunInput) IsTerminal() bool { return i.isTerminal }

type CaseRunOutput struct {
	io.Writer
	Terminal bool
	Width    int
	Height   int
}

func (o CaseRunOutput) Size() (int, int) { return o.Width, o.Height }
func (o CaseRunOutput) IsTerminal() bool { return o.Terminal }

type CaseRun struct {
	LineNr           int
	Case             *Case
	Command          string
	Env              []string
	args             []string
	StdinInitial     string
	ExpectedStdout   string
	ExpectedStderr   string
	ExpectedExitCode int
	ActualStdoutBuf  *bytes.Buffer
	ActualStderrBuf  *bytes.Buffer
	ActualExitCode   int
	Readlines        []CaseReadline
	ReadlinesPos     int
	ReadlineEnv      []string
	WasRun           bool
}

func (cr *CaseRun) Line() int { return cr.LineNr }

func (cr *CaseRun) getEnv(name string) string {
	for _, kv := range cr.Environ() {
		if strings.HasPrefix(kv, name+"=") {
			return kv[len(name)+1:]
		}
	}
	return ""
}

func (cr *CaseRun) getEnvInt(name string) int {
	n, _ := strconv.Atoi(cr.getEnv(name))
	return n
}

func (cr *CaseRun) Platform() interp.Platform {
	return interp.Platform{
		OS:        "testos",
		Arch:      "testarch",
		GoVersion: "testgo_version",
	}
}

func (cr *CaseRun) Stdin() interp.Input {
	return CaseRunInput{
		FileReader: interp.FileReader{
			R: bytes.NewBufferString(cr.StdinInitial),
		},
		isTerminal: cr.StdinInitial == "" || cr.getEnvInt("_STDIN_IS_TERMINAL") != 0,
		width:      cr.getEnvInt("_STDIN_WIDTH"),
		height:     cr.getEnvInt("_STDIN_HEIGHT"),
	}
}

func (cr *CaseRun) Stdout() interp.Output {
	var w io.Writer = cr.ActualStdoutBuf
	if cr.getEnvInt("_STDOUT_HEX") != 0 {
		w = hex.NewEncoder(cr.ActualStdoutBuf)
	}

	return CaseRunOutput{
		Writer:   w,
		Terminal: cr.getEnvInt("_STDOUT_IS_TERMINAL") != 0,
		Width:    cr.getEnvInt("_STDOUT_WIDTH"),
		Height:   cr.getEnvInt("_STDOUT_HEIGHT"),
	}
}

func (cr *CaseRun) Stderr() interp.Output {
	return CaseRunOutput{Writer: cr.ActualStderrBuf}
}

func (cr *CaseRun) InterruptChan() chan struct{} { return nil }

func (cr *CaseRun) Environ() []string {
	env := []string{
		"_STDIN_WIDTH=135",
		"_STDIN_HEIGHT=25",
		"_STDOUT_WIDTH=135",
		"_STDOUT_HEIGHT=25",
		"_STDOUT_IS_TERMINAL=1",
		"_STDIN_IS_TERMINAL=0",
		"NO_COLOR=1",
		"NO_DECODE_PROGRESS=1",
		"COMPLETION_TIMEOUT=10", // increase to make -race work better
	}
	env = append(env, cr.Env...)
	env = append(env, cr.ReadlineEnv...)

	envm := make(map[string]string)
	for _, kv := range env {
		if i := strings.IndexByte(kv, '='); i > 0 {
			envm[kv[:i]] = kv[i+1:]
		}
	}

	env = []string{}
	for k, v := range envm {
		env = append(env, k+"="+v)
	}

	return env
}

func (cr *CaseRun) Args() []string { return cr.args }

func (cr *CaseRun) ConfigDir() (string, error) { return "/config", nil }

func (cr *CaseRun) FS() fs.FS { return cr.Case }

func (cr *CaseRun) Readline(opts interp.ReadlineOpts) (string, error) {
	cr.ActualStdoutBuf.WriteString(opts.Prompt)
	if cr.ReadlinesPos >= len(cr.Readlines) {
		return "", io.EOF
	}

	expr := cr.Readlines[cr.ReadlinesPos].expr
	lineRaw := cr.Readlines[cr.ReadlinesPos].input
	line := Unescape(lineRaw)
	cr.ReadlineEnv = cr.Readlines[cr.ReadlinesPos].env
	cr.ReadlinesPos++

	if strings.HasSuffix(line, "\t") {
		cr.ActualStdoutBuf.WriteString(lineRaw + "\n")

		l := len(line) - 1
		newLine, shared := opts.CompleteFn(line[0:l], l)
		// TODO: shared
		_ = shared
		for _, nl := range newLine {
			cr.ActualStdoutBuf.WriteString(nl + "\n")
		}

		return "", nil
	}

	cr.ActualStdoutBuf.WriteString(expr + "\n")

	if line == "^D" {
		return "", io.EOF
	}

	return line, nil
}
func (cr *CaseRun) History() ([]string, error) { return nil, nil }

func (cr *CaseRun) ToExpectedStdout() string {
	sb := &strings.Builder{}

	if len(cr.Readlines) == 0 {
		fmt.Fprint(sb, cr.ExpectedStdout)
	} else {
		for _, rl := range cr.Readlines {
			fmt.Fprintf(sb, "%s%s\n", rl.expectedPrompt, rl.expr)
			if rl.expectedStdout != "" {
				fmt.Fprint(sb, rl.expectedStdout)
			}
		}
	}

	return sb.String()
}

func (cr *CaseRun) ToExpectedStderr() string {
	return cr.ExpectedStderr
}

type part interface {
	Line() int
}

type caseFile struct {
	lineNr int
	name   string
	data   []byte
}

func (cf *caseFile) Line() int { return cf.lineNr }

type caseComment struct {
	lineNr  int
	comment string
}

func (cc *caseComment) Line() int { return cc.lineNr }

type Case struct {
	Path   string
	Parts  []part
	WasRun bool
}

func (c *Case) ToActual() string {
	var partsLineSorted []part
	partsLineSorted = append(partsLineSorted, c.Parts...)
	slices.SortFunc(partsLineSorted, func(a, b part) int { return cmp.Compare(a.Line(), b.Line()) })

	sb := &strings.Builder{}
	for _, p := range partsLineSorted {
		switch p := p.(type) {
		case *caseComment:
			fmt.Fprintf(sb, "#%s\n", p.comment)
		case *CaseRun:
			fmt.Fprintf(sb, "$%s\n", p.Command)
			var s string
			if p.WasRun {
				s = p.ActualStdoutBuf.String()
			} else {
				s = p.ToExpectedStdout()
			}
			if s != "" {
				fmt.Fprint(sb, s)
				if !strings.HasSuffix(s, "\n") {
					fmt.Fprint(sb, "\\\n")
				}
			}
			if p.WasRun {
				if p.ActualExitCode != 0 {
					fmt.Fprintf(sb, "exitcode: %d\n", p.ActualExitCode)
				}
			} else {
				if p.ExpectedExitCode != 0 {
					fmt.Fprintf(sb, "exitcode: %d\n", p.ExpectedExitCode)
				}
			}
			if p.StdinInitial != "" {
				fmt.Fprint(sb, "stdin:\n")
				fmt.Fprint(sb, p.StdinInitial)
			}
			if p.WasRun {
				if p.ActualStderrBuf.Len() > 0 {
					fmt.Fprint(sb, "stderr:\n")
					fmt.Fprint(sb, p.ActualStderrBuf.String())
				}
			} else {
				if p.ExpectedStderr != "" {
					fmt.Fprint(sb, "stderr:\n")
					fmt.Fprint(sb, p.ExpectedStderr)
				}
			}
		case *caseFile:
			fmt.Fprintf(sb, "%s:\n", p.name)
			sb.Write(p.data)
		default:
			panic("unreachable")
		}
	}

	return sb.String()
}

func normalizeOSError(err error) error {
	var pe *os.PathError
	if errors.As(err, &pe) {
		pe.Err = errors.New("no such file or directory")
		pe.Path = filepath.ToSlash(pe.Path)
	}
	return err
}

func (c *Case) Open(name string) (fs.File, error) {
	const testData = "testdata"
	testDataIndex := strings.Index(c.Path, testData)
	// cwd is directory where current script file is
	testRoot := c.Path[0 : testDataIndex+len(testData)]
	testCwd := filepath.Dir(c.Path[testDataIndex+len(testData):])
	testAbsPath := filepath.Join(testCwd, name)
	fsPath := filepath.Join(testRoot, testAbsPath)

	for _, p := range c.Parts {
		f, ok := p.(*caseFile)
		if !ok {
			continue
		}
		if f.name == filepath.ToSlash(testAbsPath) {
			return interp.FileReader{
				R: io.NewSectionReader(bytes.NewReader(f.data), 0, int64(len(f.data))),
				FileInfo: interp.FixedFileInfo{
					FName: filepath.Base(name),
					FSize: int64(len(f.data)),
				},
			}, nil
		}
	}
	f, err := os.Open(fsPath)
	// normalizeOSError is used to normalize OS specific path and messages into the ones unix uses
	// this needed to make difftest work
	return f, normalizeOSError(err)
}

type Section struct {
	LineNr int
	Name   string
	Value  string

	valueSB strings.Builder
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
			cs.valueSB.WriteString(l)
			cs.valueSB.WriteString(lineDelim)
		}
	}

	for i := range sections {
		cs := &sections[i]
		cs.Value = cs.valueSB.String()
	}

	return sections
}

var kvRe = regexp.MustCompile(`^[A-Z_]+=`)

func ParseCommand(s string) (env []string, args []string) {
	parts := shquote.Split(s)
	for i, p := range parts {
		if kvRe.MatchString(p) {
			env = append(env, p)
			continue
		}
		args = parts[i:]
		break
	}

	return env, args
}

func ParseInput(s string) (env []string, input string) {
	tokens := shquote.Parse(s)
	l := 0
	for _, t := range tokens {
		if t.Separator {
			continue
		}
		if kvRe.MatchString(t.Str) {
			env = append(env, t.Str)
			l = t.End
			continue
		}
		break
	}
	return env, s[l:]
}

func ParseCases(s string) *Case {
	te := &Case{}
	te.Parts = []part{}
	var currentRun *CaseRun
	const promptEnd = ">"
	replDepth := 0

	// TODO: better section splitter, too much heuristics now
	for _, section := range SectionParser(regexp.MustCompile(
		`^\$ .*$|^stdin:$|^stderr:$|^exitcode:.*$|^#.*$|^/.*:|^[^<|"]+>.*$`,
	), s) {
		n, v := section.Name, section.Value

		switch {
		case strings.HasPrefix(n, "#"):
			comment := n[1:]
			te.Parts = append(te.Parts, &caseComment{lineNr: section.LineNr, comment: comment})
		case strings.HasPrefix(n, "/"):
			name := n[0 : len(n)-1]
			te.Parts = append(te.Parts, &caseFile{lineNr: section.LineNr, name: name, data: []byte(v)})
		case strings.HasPrefix(n, "$"):
			replDepth++

			if currentRun != nil {
				te.Parts = append(te.Parts, currentRun)
			}

			// escaped newline
			v = strings.TrimSuffix(v, "\\\n")
			command := strings.TrimPrefix(n, "$")
			env, args := ParseCommand(command)

			currentRun = &CaseRun{
				LineNr:          section.LineNr,
				Case:            te,
				Command:         command,
				Env:             env,
				args:            args,
				ExpectedStdout:  v,
				ActualStdoutBuf: &bytes.Buffer{},
				ActualStderrBuf: &bytes.Buffer{},
			}
		case strings.HasPrefix(n, "exitcode:"):
			currentRun.ExpectedExitCode, _ = strconv.Atoi(strings.TrimSpace(strings.TrimPrefix(n, "exitcode:")))
		case strings.HasPrefix(n, "stdin"):
			currentRun.StdinInitial = v
		case strings.HasPrefix(n, "stderr"):
			currentRun.ExpectedStderr = v
		case strings.Contains(n, promptEnd+" ") || strings.HasSuffix(n, promptEnd): // TODO: better
			i := strings.LastIndex(n, promptEnd+" ")
			if strings.HasSuffix(n, promptEnd) {
				i = len(n) - 1
			}

			prompt := n[0:i] + promptEnd + " "
			expr := strings.TrimSpace(n[i+1:])
			env, input := ParseInput(expr)

			currentRun.Readlines = append(currentRun.Readlines, CaseReadline{
				expr:           expr,
				env:            env,
				input:          input,
				expectedPrompt: prompt,
				expectedStdout: v,
			})

			// TODO: hack
			if strings.Contains(expr, "| repl") {
				replDepth++
			}
			if expr == "^D" {
				replDepth--
			}

		default:
			panic(fmt.Sprintf("%d: unexpected section %q %q", section.LineNr, n, v))
		}
	}

	if currentRun != nil {
		te.Parts = append(te.Parts, currentRun)
	}

	return te
}
