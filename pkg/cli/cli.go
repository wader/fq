package cli

import (
	"context"
	"flag"
	"fmt"
	"fq"
	"fq/pkg/decode"
	"fq/pkg/osenv"
	"fq/pkg/query"
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

type StandardOS struct {
	rl *readline.Instance
}

func newStandardOS() (*StandardOS, error) {
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

	return &StandardOS{rl: rl}, nil
}

func (*StandardOS) Stdin() io.Reader                        { return os.Stdin }
func (*StandardOS) Stdout() io.Writer                       { return os.Stdout }
func (*StandardOS) Stderr() io.Writer                       { return os.Stderr }
func (*StandardOS) Environ() []string                       { return os.Environ() }
func (*StandardOS) Args() []string                          { return os.Args }
func (*StandardOS) Open(name string) (io.ReadSeeker, error) { return os.Open(name) }
func (o *StandardOS) Readline(prompt string, complete func(line string, pos int) (newLine []string, shared int)) (string, error) {
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

func StandardOSMain(r *decode.Registry) {
	o, err := newStandardOS()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if err := (Main{
		OS:       o,
		Registry: r,
	}).Run(); err != nil {
		os.Exit(1)
	}
}

type Main struct {
	OS       osenv.OS
	Registry *decode.Registry
}

// Run cli main
func (m Main) Run() error {
	err := m.run()
	if err != nil && err != flag.ErrHelp {
		fmt.Fprintf(m.OS.Stderr(), "%s\n", err)
	}
	return err
}

func (m Main) run() error {
	// TODO: pass with some kind of env?
	filename := "asd"

	var args []interface{}
	for _, a := range m.OS.Args() {
		args = append(args, a)
	}

	q := query.NewQuery(query.QueryOptions{
		Variables: map[string]interface{}{
			"$FILENAME": filename,
			"$VERSION":  fq.Version,
			"$ARGS":     args,
		},
		Registry: m.Registry,

		Environ:  m.OS.Environ, // TODO: func?
		Stdin:    m.OS.Stdin(),
		Open:     m.OS.Open,
		Readline: m.OS.Readline,
	})

	runMode := query.ScriptMode

	i, err := q.Eval(context.Background(), runMode, nil, "main", query.WriterOutput{Ctx: context.Background(), W: m.OS.Stdout()}, nil)
	if err != nil {
		return err
	}
	for {
		v, ok := i.Next()
		if !ok {
			break
		} else if err, ok := v.(error); ok {
			fmt.Fprintln(m.OS.Stderr(), err)
			break
		} else if d, ok := v.([2]interface{}); ok {
			fmt.Fprintf(m.OS.Stdout(), "%s: %v\n", d[0], d[1])
		}
	}

	return nil
}
