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

type StandardOS struct{}

func (StandardOS) Stdin() io.Reader                        { return os.Stdin }
func (StandardOS) Stdout() io.Writer                       { return os.Stdout }
func (StandardOS) Stderr() io.Writer                       { return os.Stderr }
func (StandardOS) Environ() []string                       { return os.Environ() }
func (StandardOS) Args() []string                          { return os.Args }
func (StandardOS) Open(name string) (io.ReadSeeker, error) { return os.Open(name) }
func (o StandardOS) Readline(prompt string, complete func(line string, pos int) (newLine []string, shared int)) (string, error) {
	// TODO: refactor, shared?
	historyFile := ""
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	historyFile = filepath.Join(cacheDir, "fq/history")
	_ = os.MkdirAll(filepath.Dir(historyFile), 0700)

	var autoComplete readline.AutoCompleter
	if complete != nil {
		autoComplete = autoCompleterFn(func(line []rune, pos int) (newLine [][]rune, length int) {
			// completeCtx, completeCtxCancelFn := context.WithTimeout(ctx, 1*time.Second)
			// defer completeCtxCancelFn()

			// // TODO: err
			// names, shared, _ := completeTrampoline(completeCtx, completeFn, c, q, string(line), pos)

			names, shared := complete(string(line), pos)

			var runeNames [][]rune
			for _, name := range names {
				runeNames = append(runeNames, []rune(name[shared:]))
			}

			return runeNames, shared
		})
	}

	l, err := readline.NewEx(&readline.Config{
		Stdin:        ioutil.NopCloser(os.Stdin),
		Stdout:       os.Stdin,
		Stderr:       os.Stderr, // TODO: ??
		HistoryFile:  historyFile,
		AutoComplete: autoComplete,
		// InterruptPrompt: "^C",
		// EOFPrompt:       "exit",

		HistorySearchFold: true,
		// FuncFilterInputRune: filterInput,

		// FuncFilterInputRune: func(r rune) (rune, bool) {
		// 	log.Printf("r: %#+v\n", r)
		// 	return r, true
		// },

		// Listener: listenerFn(func(line []rune, pos int, key rune) (newLine []rune, newPos int, ok bool) {
		// 	log.Printf("line: %#+v pos=%v key=%d\n", line, pos, key)
		// 	return line, pos, false
		// }),
	})
	if err != nil {
		return "", err
	}

	l.SetPrompt(prompt)
	src, err := l.Readline()
	if err != nil {
		return "", err
	}

	return src, nil
}

func StandardOSMain(r *decode.Registry) {
	if err := (Main{
		OS:       StandardOS{},
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
