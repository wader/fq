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
	"runtime"

	"github.com/wader/fq/pkg/interp"
	"golang.org/x/term"

	"github.com/ergochat/readline"
)

func maybeLogFile() {
	// used during dev to redirect log to file, useful when debugging repl etc
	if lf := os.Getenv("LOGFILE"); lf != "" {
		if f, err := os.Create(lf); err == nil {
			log.SetOutput(f)
		}
	}
}

// function implementing readline.AutoComplete interface
type autoCompleterFn func(line []rune, pos int) (newLine [][]rune, length int)

func (a autoCompleterFn) Do(line []rune, pos int) (newLine [][]rune, length int) {
	return a(line, pos)
}

type stdOS struct {
	rl            *readline.Instance
	historyFile   string
	closeChan     chan struct{}
	interruptChan chan struct{}
	completerFn   interp.CompleteFn
}

func newStandardOS() *stdOS {
	closeChan := make(chan struct{})
	interruptChan := make(chan struct{}, 1)

	// this more or less converts a os signal chan to just a struct{} chan that
	// ignores signals if forwarding it would block, also this makes sure interp
	// does not know about os.
	go func() {
		interruptSignalChan := make(chan os.Signal, 1)
		signal.Notify(interruptSignalChan, os.Interrupt)
		defer func() {
			signal.Stop(interruptSignalChan)
			close(interruptSignalChan)
			close(interruptChan)
		}()

		for {
			select {
			case <-interruptSignalChan:
				// ignore if interruptChan is full
				select {
				case interruptChan <- struct{}{}:
				default:
				}
			case <-closeChan:
				return
			}
		}
	}()

	return &stdOS{
		closeChan:     closeChan,
		interruptChan: interruptChan,
	}
}

func (stdOS) Platform() interp.Platform {
	return interp.Platform{
		OS:        runtime.GOOS,
		Arch:      runtime.GOARCH,
		GoVersion: runtime.Version(),
	}
}

type fdTerminal uintptr

func (fd fdTerminal) Size() (int, int) {
	w, h, _ := term.GetSize(int(fd))
	// TODO: old version return 0 on no terminal
	w, h = max(0, w), max(0, h)
	return w, h
}
func (fd fdTerminal) IsTerminal() bool {
	return term.IsTerminal(int(fd))
}

type stdinInput struct {
	fdTerminal
	fs.File
}

func (o *stdOS) Stdin() interp.Input {
	return stdinInput{
		fdTerminal: fdTerminal(os.Stdin.Fd()),
		File: interp.FileReader{
			R: os.Stdin,
			FileInfo: interp.FixedFileInfo{
				FName: "stdin",
				FMode: fs.ModeIrregular,
			},
		},
	}
}

type stdoutOutput struct {
	fdTerminal
	os *stdOS
}

func (o stdoutOutput) Write(p []byte) (n int, err error) {
	// Let write go thru readline if it has been used. This to have ansi color emulation
	// on windows thru readlins:s stdout rewriter
	// TODO: check if tty instead? else only color when repl
	if o.os.rl != nil {
		return o.os.rl.Write(p)
	}
	return os.Stdout.Write(p)
}

func (o *stdOS) Stdout() interp.Output {
	return stdoutOutput{fdTerminal: fdTerminal(os.Stdout.Fd()), os: o}
}

type stderrOutput struct {
	fdTerminal
}

func (o stderrOutput) Write(p []byte) (n int, err error) { return os.Stderr.Write(p) }

func (o *stdOS) Stderr() interp.Output { return stderrOutput{fdTerminal: fdTerminal(os.Stderr.Fd())} }

func (o *stdOS) InterruptChan() chan struct{} { return o.interruptChan }

func (*stdOS) Args() []string { return os.Args }

func (*stdOS) Environ() []string { return os.Environ() }

func (*stdOS) ConfigDir() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	fqDir := filepath.Join(configDir, "fq")

	if runtime.GOOS != "darwin" {
		return fqDir, nil
	}

	// this is to support fallback to ~/.config on macOS/darwin
	if _, err := os.Stat(fqDir); err == nil {
		return fqDir, nil
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(homeDir, ".config", "fq"), nil
}

type stdOSFS struct{}

func (stdOSFS) Open(name string) (fs.File, error) { return os.Open(name) }

func (*stdOS) FS() fs.FS { return stdOSFS{} }

func (o *stdOS) Readline(opts interp.ReadlineOpts) (string, error) {
	if o.rl == nil {
		var err error

		var historyFile string
		cacheDir, err := os.UserCacheDir()
		if err != nil {
			return "", err
		}
		historyFile = filepath.Join(cacheDir, "fq/history")
		_ = os.MkdirAll(filepath.Dir(historyFile), 0700)

		cfg := &readline.Config{
			HistoryFile:       historyFile,
			HistorySearchFold: true,
			Undo:              true,
		}
		cfg.AutoComplete = autoCompleterFn(func(line []rune, pos int) (newLine [][]rune, length int) {
			if o.completerFn == nil {
				return nil, 0
			}

			names, shared := o.completerFn(string(line), pos)
			var runeNames [][]rune
			for _, name := range names {
				runeNames = append(runeNames, []rune(name[shared:]))
			}

			return runeNames, shared
		})
		o.rl, err = readline.NewEx(cfg)
		if err != nil {
			return "", err
		}
		o.historyFile = historyFile
	}

	// inject completer to autocompleter
	o.completerFn = opts.CompleteFn

	o.rl.SetPrompt(opts.Prompt)
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

func (o *stdOS) History() ([]string, error) {
	// TODO: refactor history handling to use internal fs?
	r, err := os.Open(o.historyFile)
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

func (o *stdOS) Close() error {
	// only close if is terminal otherwise ansi reset will write
	// to stdout and mess up raw output
	if o.rl != nil {
		o.rl.Close()
	}
	close(o.closeChan)
	return nil
}

func Main(r *interp.Registry, version string) {
	os.Exit(func() int {
		defer maybeProfile()()
		maybeLogFile()

		sos := newStandardOS()
		defer sos.Close()
		i, err := interp.New(sos, r)
		defer i.Stop()
		if err != nil {
			fmt.Fprintln(sos.Stderr(), err)
			return 1
		}

		if err := i.Main(context.Background(), sos.Stdout(), version); err != nil {
			if ex, ok := err.(interp.Exiter); ok {
				return ex.ExitCode()
			}
			return 1
		}

		return 0
	}())
}
