// +build !js,!wasm

// TODO: break out os specific readline into interface?

package query

import (
	"context"
	"fmt"
	"fq/pkg/decode"
	"io"
	"io/ioutil"
	"os"
	"os/signal"
	"path/filepath"
	"time"

	"github.com/chzyer/readline"
)

// TODO: move to OS? Input/Outout interfaces?
type Output interface {
	io.Writer
	Size() (int, int)
	IsTerminal() bool
}

// TODO: move
type DiscardOutput struct{}

func (o DiscardOutput) Write(p []byte) (n int, err error) {
	return n, nil
}
func (o DiscardOutput) Size() (int, int) { return 0, 0 }
func (o DiscardOutput) IsTerminal() bool { return false }

type WriterOutput struct {
	Ctx context.Context
	W   io.Writer
}

// TODO: move
func (o WriterOutput) Write(p []byte) (n int, err error) {
	if err := o.Ctx.Err(); err != nil {
		return 0, err
	}
	return o.W.Write(p)
}

func (o WriterOutput) Size() (int, int) {
	f, ok := o.W.(interface{ Fd() uintptr })
	if ok {
		w, h, _ := readline.GetSize(int(f.Fd()))
		return w, h
	}
	return 0, 0
}

func (o WriterOutput) IsTerminal() bool {
	if f, ok := o.W.(interface{ Fd() uintptr }); ok {
		return readline.IsTerminal(int(f.Fd()))
	}
	return false
}

// REPL read-eval-print-loop
func (q *Query) REPL(ctx context.Context) error {
	// TODO: refactor
	historyFile := ""
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return err
	}
	historyFile = filepath.Join(cacheDir, "fq/history")
	_ = os.MkdirAll(filepath.Dir(historyFile), 0700)

	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt)

	l, err := readline.NewEx(&readline.Config{
		Stdin:       ioutil.NopCloser(q.opts.OS.Stdin()),
		Stdout:      q.opts.OS.Stdout(),
		Stderr:      q.opts.OS.Stderr(),
		HistoryFile: historyFile,
		AutoComplete: autoCompleterFn(func(line []rune, pos int) (newLine [][]rune, length int) {
			completeCtx, completeCtxCancelFn := context.WithTimeout(ctx, 1*time.Second)
			defer completeCtxCancelFn()
			return autoComplete(completeCtx, nil, q, line, pos)
		}),
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
		return err
	}

	for {
		if ok, err := func() (bool, error) {
			var v []interface{}
			stackLenStr := ""
			if len(q.inputStack) > 0 {
				v = q.inputStack[len(q.inputStack)-1]
			}
			if len(q.inputStack) > 1 {
				stackLenStr = fmt.Sprintf("[%d]", len(q.inputStack))
			}
			inputSummary := ""
			if len(v) > 0 {
				first := v[0]
				if vv, ok := first.(*decode.Value); ok {
					inputSummary = valuePath(vv)
				} else if t, ok := valueToTypeString(first); ok {
					inputSummary = t
				} else {
					inputSummary = "?"
				}
			}
			if len(v) > 1 {
				inputSummary = "(" + inputSummary + ",...)"
			}
			prompt := fmt.Sprintf("%s%s> ", stackLenStr, inputSummary)

			l.SetPrompt(prompt)

			src, err := l.Readline()
			if err == readline.ErrInterrupt {
				return true, nil
			} else if err == io.EOF {
				return false, nil
			}

			if err != nil {
				return false, err
			}

			interruptCtx, interruptCtxCancelFn := context.WithCancel(ctx)
			defer interruptCtxCancelFn()
			go func() {
				select {
				case <-interruptChan:
					interruptCtxCancelFn()
				case <-interruptCtx.Done():
					// nop
				}
			}()

			output := WriterOutput{
				Ctx: interruptCtx,
				W:   q.opts.OS.Stdout(),
			}

			if _, err := q.Run(interruptCtx, REPLMode, src, output); err != nil {
				if err != context.Canceled {
					fmt.Fprintf(q.opts.OS.Stdout(), "error: %s\n", err)
				}
			}

			return true, nil
		}(); !ok {
			return err
		}
	}
}
