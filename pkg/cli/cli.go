package cli

import (
	"context"
	"fmt"
	"fq/pkg/decode"
	"fq/pkg/interp"
	"io"
	"io/ioutil"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/goinsane/readline"
)

type autoCompleterFn func(line []rune, pos int) (newLine [][]rune, length int)

func (a autoCompleterFn) Do(line []rune, pos int) (newLine [][]rune, length int) {
	return a(line, pos)
}

type standardOS struct {
	rl                  *readline.Instance
	interruptSignalChan chan os.Signal
	interruptChan       chan struct{}
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

	interruptChan := make(chan struct{}, 1)
	interruptSignalChan := make(chan os.Signal, 1)
	signal.Notify(interruptSignalChan, os.Interrupt)
	go func() {
		defer signal.Stop(interruptSignalChan)
		for range interruptSignalChan {
			select {
			case interruptChan <- struct{}{}:
			default:
			}
		}
	}()

	return &standardOS{
		rl:                  rl,
		interruptSignalChan: interruptSignalChan,
		interruptChan:       interruptChan,
	}, nil
}

func (o standardOS) Close() error {
	close(o.interruptSignalChan)
	return nil
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
func (o *standardOS) Interrupt() chan struct{}              { return o.interruptChan }
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

	// _ = autoComplete

	// s := bufio.NewScanner(os.Stdin)
	// if ok := s.Scan(); !ok {
	// 	fmt.Fprintln(os.Stdout)
	// 	return "", io.EOF
	// }
	// line := s.Text()
	// if err := s.Err(); err != nil {
	// 	return "", err
	// }
	// return line, nil

	o.rl.Config.AutoComplete = autoComplete
	o.rl.SetPrompt(prompt)
	line, err := o.rl.Readline()
	if err == readline.ErrInterrupt {
		return "", interp.ErrInterrupt
	} else if err == io.EOF {
		return "", interp.ErrEOF
	} else if err != nil {
		return "", err
	}

	return line, nil
}

func Main(r *decode.Registry) {
	o, err := newStandardOS()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer o.Close()
	i, err := interp.New(interp.InterpOptions{
		Registry: r,
		OS:       o,
	})
	if err != nil {
		fmt.Fprintln(o.Stderr(), err)
		os.Exit(1)
	}
	if err := i.Main(context.Background(), o.Stdout()); err != nil {
		fmt.Fprintln(o.Stderr(), err)
		os.Exit(1)
	}
}
