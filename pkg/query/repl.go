// +build !js,!wasm

// TODO: break out os specific readline into interface?

package query

import (
	"context"
	"fmt"
	"fq/internal/ioextra"
	"fq/pkg/decode"
	"io"
	"io/ioutil"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"time"

	"github.com/chzyer/readline"
)

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
			return autoComplete(completeCtx, q, line, pos)
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
		runCtx, runCtxCancelFn := context.WithCancel(ctx)
		_ = runCtxCancelFn
		go func() {
			select {
			case <-interruptChan:
				runCtxCancelFn()
			case <-ctx.Done():
				// nop
			}
		}()

		var v []interface{}
		if len(q.inputStack) > 0 {
			v = q.inputStack[len(q.inputStack)-1]
		}
		var inputSummary []string
		if len(v) > 0 {
			first := v[0]
			if vv, ok := first.(*decode.Value); ok {
				inputSummary = append(inputSummary, vv.Path())
			} else if t, ok := valueToTypeString(first); ok {
				inputSummary = append(inputSummary, t)
			} else {
				inputSummary = append(inputSummary, "?")
			}
		}
		if len(v) > 1 {
			inputSummary = append(inputSummary, "...")
		}
		prompt := fmt.Sprintf("inputs[%d] [%s]> ", len(q.inputStack), strings.Join(inputSummary, ","))

		l.SetPrompt(prompt)

		src, err := l.Readline()
		if err == readline.ErrInterrupt {
			continue
		} else if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		if _, err := q.Run(runCtx, REPLMode, src, ioextra.ContextWriter{W: q.opts.OS.Stdout(), C: runCtx}); err != nil {
			if err != context.Canceled {
				fmt.Fprintf(q.opts.OS.Stdout(), "error: %s\n", err)
			}
		}
		runCtxCancelFn()
	}

	return nil
}
