// Package ctxstack manages a stack of contexts. When triggerFn returns and closeCh is not closed
// it will cancel the top context. Stack is popped in top first order when returned cancel funcition
// is called.
// This can be used to keep track of contexts for nested REPL:s were you only want to cancel
// the current active "top" REPL.
// TODO: should New take a parent context?
package ctxstack

import (
	"context"
)

// Stack is a context stack
type Stack struct {
	cancelFns []func()
	stopCh    chan struct{}
}

// New context stack
func New(triggerCh func(stopCh chan struct{})) *Stack {
	stopCh := make(chan struct{})
	s := &Stack{stopCh: stopCh}

	go func() {
		for {
			triggerCh(stopCh)
			select {
			case <-stopCh:
				// stop if stopCh closed
			default:
				// ignore if triggered before any context pushed
				if len(s.cancelFns) > 0 {
					s.cancelFns[len(s.cancelFns)-1]()
				}
				continue
			}
			break
		}
	}()

	return s
}

// Stop context stack
func (s *Stack) Stop() {
	for i := len(s.cancelFns) - 1; i >= 0; i-- {
		s.cancelFns[i]()
	}
	close(s.stopCh)
}

// Push creates, pushes and returns new context. Cancel pops it.
func (s *Stack) Push(parent context.Context) (context.Context, func()) {
	stackCtx, stackCtxCancel := context.WithCancel(parent)
	stackIdx := len(s.cancelFns)

	s.cancelFns = append(s.cancelFns, stackCtxCancel)
	cancelled := false

	return stackCtx, func() {
		if cancelled {
			return
		}
		cancelled = true

		for i := len(s.cancelFns) - 1; i >= stackIdx; i-- {
			s.cancelFns[i]()
		}
		s.cancelFns = s.cancelFns[0:stackIdx]

		stackCtxCancel()
	}
}
