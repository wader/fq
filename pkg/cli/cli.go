package cli

import (
	"fmt"
	"fq/pkg/decode"
	"fq/pkg/interp"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/chzyer/readline"
)

type autoCompleterFn func(line []rune, pos int) (newLine [][]rune, length int)

func (a autoCompleterFn) Do(line []rune, pos int) (newLine [][]rune, length int) {
	return a(line, pos)
}

type standardOS struct {
	rl *readline.Instance
}

func newStandardOS() (*standardOS, error) {
	// TODO: refactor, shared?
	historyFile := ""
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return nil, err
	}
	historyFile = filepath.Join(cacheDir, "fq/history")
	_ = os.MkdirAll(filepath.Dir(historyFile), 0700)

	rl, err := readline.NewEx(&readline.Config{
		Stdin:             ioutil.NopCloser(os.Stdin),
		Stdout:            os.Stdin,
		Stderr:            os.Stderr,
		HistoryFile:       historyFile,
		HistorySearchFold: true,
	})
	if err != nil {
		return nil, err
	}

	return &standardOS{rl: rl}, nil
}

type standardOsOutput struct{}

func (o standardOsOutput) Write(p []byte) (n int, err error) {
	return os.Stdout.Write(p)
}

func (o standardOsOutput) Size() (int, int) {
	w, h, _ := readline.GetSize(int(os.Stdout.Fd()))
	return w, h
}

func (o standardOsOutput) IsTerminal() bool {
	return readline.IsTerminal(int(os.Stdout.Fd()))
}

func (*standardOS) Stdin() io.Reader                        { return os.Stdin }
func (*standardOS) Stdout() interp.Output                   { return standardOsOutput{} }
func (*standardOS) Stderr() io.Writer                       { return os.Stderr }
func (*standardOS) Environ() []string                       { return os.Environ() }
func (*standardOS) Args() []string                          { return os.Args }
func (*standardOS) Open(name string) (io.ReadSeeker, error) { return os.Open(name) }
func (o *standardOS) Readline(prompt string, complete func(line string, pos int) (newLine []string, shared int)) (string, error) {
	var autoComplete readline.AutoCompleter
	if complete != nil {
		autoComplete = autoCompleterFn(func(line []rune, pos int) (newLine [][]rune, length int) {
			names, shared := complete(string(line), pos)
			var runeNames [][]rune
			for _, name := range names {
				runeNames = append(runeNames, []rune(name[shared:]))
			}

			return runeNames, shared
		})
	}

	o.rl.Config.AutoComplete = autoComplete
	o.rl.SetPrompt(prompt)
	src, err := o.rl.Readline()
	if err != nil {
		return "", err
	}

	return src, nil
}

func Main(r *decode.Registry) {
	o, err := newStandardOS()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	q, err := interp.New(interp.InterpOptions{
		Registry: r,
		OS:       o,
	})
	if err != nil {
		fmt.Fprintln(o.Stderr(), err)
		os.Exit(1)
	}
	if err := q.Main(o.Stdout()); err != nil {
		os.Exit(1)
	}
}
