package cli

import (
	"bufio"
	"context"
	"fmt"
	"fq/format/registry"
	"fq/internal/profile"
	"fq/pkg/interp"
	"io"
	"log"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/chzyer/readline"
)

func MaybeProfile() func() {
	return profile.Start(os.Getenv("CPUPROFILE"), os.Getenv("MEMPROFILE"))
}

func MaybeLogFile() {
	if lf := os.Getenv("LOGFILE"); lf != "" {
		log.SetOutput(func() io.Writer { f, _ := os.Create(lf); return f }())
	}
}

type Exiter interface {
	ExitCode() int
}

type autoCompleterFn func(line []rune, pos int) (newLine [][]rune, length int)

func (a autoCompleterFn) Do(line []rune, pos int) (newLine [][]rune, length int) {
	return a(line, pos)
}

type standardOS struct {
	rl                  *readline.Instance
	output              interp.Output
	interruptSignalChan chan os.Signal
	interruptChan       chan struct{}
}

func newStandardOS() (*standardOS, error) {
	historyFile := ""
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return nil, err
	}
	historyFile = filepath.Join(cacheDir, "fq/history")
	_ = os.MkdirAll(filepath.Dir(historyFile), 0700)

	rl, err := readline.NewEx(&readline.Config{
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
		output:              standardOsOutput{rl: rl},
	}, nil
}

func (o standardOS) Close() error {
	// TODO: only close if is terminal somehow? otherwise reset will write
	// to stdout and mess up raw output
	// o.rl.Close()
	close(o.interruptSignalChan)
	return nil
}

type standardOsOutput struct {
	rl *readline.Instance
}

func (o standardOsOutput) Write(p []byte) (n int, err error) {
	return o.rl.Write(p)
}

func (o standardOsOutput) Size() (int, int) {
	w, h, _ := readline.GetSize(int(os.Stdout.Fd()))
	return w, h
}

func (o standardOsOutput) IsTerminal() bool {
	return readline.IsTerminal(int(os.Stdout.Fd()))
}

func (o *standardOS) Stdin() io.Reader         { return o.rl.Config.Stdin }
func (o *standardOS) Stdout() interp.Output    { return o.output }
func (o *standardOS) Stderr() io.Writer        { return o.rl.Stderr() }
func (o *standardOS) Interrupt() chan struct{} { return o.interruptChan }
func (*standardOS) Args() []string             { return os.Args }
func (*standardOS) Environ() []string          { return os.Environ() }
func (*standardOS) ConfigDir() (string, error) {
	p, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(p, "fq"), nil
}
func (*standardOS) Open(name string) (io.ReadSeeker, error) { return os.Open(name) }
func (o *standardOS) Readline(prompt string, complete func(line string, pos int) (newLine []string, shared int)) (string, error) {
	if complete != nil {
		o.rl.Config.AutoComplete = autoCompleterFn(func(line []rune, pos int) (newLine [][]rune, length int) {
			names, shared := complete(string(line), pos)
			var runeNames [][]rune
			for _, name := range names {
				runeNames = append(runeNames, []rune(name[shared:]))
			}

			return runeNames, shared
		})
	}

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
func (o *standardOS) History() ([]string, error) {
	// TODO: refactor history handling to use internal fs?
	r, err := os.Open(o.rl.Config.HistoryFile)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	var hs []string
	lineScanner := bufio.NewScanner(r)
	for lineScanner.Scan() {
		hs = append(hs, lineScanner.Text())
	}
	if err := lineScanner.Err(); err != nil {
		return nil, err
	}
	return hs, nil
}

func Main(r *registry.Registry, version string) {
	sos, err := newStandardOS()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer sos.Close()
	i, err := interp.New(sos, r)
	defer i.Stop()
	if err != nil {
		fmt.Fprintln(sos.Stderr(), err)
		os.Exit(1)
	}

	if err := i.Main(context.Background(), sos.Stdout(), version); err != nil {
		fmt.Fprintln(sos.Stderr(), err)
		if ex, ok := err.(Exiter); ok {
			os.Exit(ex.ExitCode())
		}
		os.Exit(1)
	}
}
