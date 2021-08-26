package cli

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/wader/fq/format/registry"
	"github.com/wader/fq/internal/profile"
	"github.com/wader/fq/pkg/interp"

	"github.com/wader/readline"
)

func MaybeProfile() func() {
	return profile.Start(os.Getenv("CPUPROFILE"), os.Getenv("MEMPROFILE"))
}

func MaybeLogFile() {
	// used during dev to redirect log to file, useful when debugging repl etc
	if lf := os.Getenv("LOGFILE"); lf != "" {
		if f, err := os.Create(lf); err == nil {
			log.SetOutput(f)
		}
	}
}

type autoCompleterFn func(line []rune, pos int) (newLine [][]rune, length int)

func (a autoCompleterFn) Do(line []rune, pos int) (newLine [][]rune, length int) {
	return a(line, pos)
}

type standardOS struct {
	rl                  *readline.Instance
	interruptSignalChan chan os.Signal
	interruptChan       chan struct{}
}

func newStandardOS() *standardOS {
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
		interruptSignalChan: interruptSignalChan,
		interruptChan:       interruptChan,
	}
}

func (o *standardOS) Stdin() fs.File {
	return interp.FileReader{
		R: os.Stdin,
		FileInfo: interp.FixedFileInfo{
			FName: "stdin",
			FMode: fs.ModeIrregular,
		},
	}
}

type standardOsOutput struct {
	os *standardOS
}

func (o standardOsOutput) Write(p []byte) (n int, err error) {
	if o.os.rl != nil {
		return o.os.rl.Write(p)
	}
	return os.Stdout.Write(p)
}

func (o standardOsOutput) Size() (int, int) {
	w, h, _ := readline.GetSize(int(os.Stdout.Fd()))
	return w, h
}

func (o standardOsOutput) IsTerminal() bool {
	return readline.IsTerminal(int(os.Stdout.Fd()))
}

func (o *standardOS) Stdout() interp.Output { return standardOsOutput{os: o} }

func (o *standardOS) Stderr() io.Writer { return os.Stderr }

func (o *standardOS) Interrupt() chan struct{} { return o.interruptChan }

func (*standardOS) Args() []string { return os.Args }

func (*standardOS) Environ() []string { return os.Environ() }

func (*standardOS) ConfigDir() (string, error) {
	p, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(p, "fq"), nil
}

type standardOSFS struct{}

func (standardOSFS) Open(name string) (fs.File, error) { return os.Open(name) }

func (*standardOS) FS() fs.FS { return standardOSFS{} }

func (o *standardOS) Readline(prompt string, complete func(line string, pos int) (newLine []string, shared int)) (string, error) {
	if o.rl == nil {
		var err error

		var historyFile string
		cacheDir, err := os.UserCacheDir()
		if err != nil {
			return "", err
		}
		historyFile = filepath.Join(cacheDir, "fq/history")
		_ = os.MkdirAll(filepath.Dir(historyFile), 0700)

		o.rl, err = readline.NewEx(&readline.Config{
			HistoryFile:       historyFile,
			HistorySearchFold: true,
		})
		if err != nil {
			return "", err
		}
	}

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
	if errors.Is(err, readline.ErrInterrupt) {
		return "", interp.ErrInterrupt
	} else if errors.Is(err, io.EOF) {
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

func (o *standardOS) Close() error {
	// only close if is terminal otherwise ansi reset will write
	// to stdout and mess up raw output
	if o.rl != nil {
		o.rl.Close()
	}
	close(o.interruptSignalChan)
	return nil
}

func Main(r *registry.Registry, version string) {
	os.Exit(func() int {
		sos := newStandardOS()
		defer sos.Close()
		i, err := interp.New(sos, r)
		defer i.Stop()
		if err != nil {
			fmt.Fprintln(sos.Stderr(), err)
			return 1
		}

		if err := i.Main(context.Background(), sos.Stdout(), version); err != nil {
			if ex, ok := err.(interp.Exiter); ok { //nolint:errorlint
				return ex.ExitCode()
			}
			return 1
		}

		return 0
	}())
}
