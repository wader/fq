// Package ctxstack manages a stack of contexts. When triggerFn returns and closeCh is not closed
// it will cancel the top context. Stack is popped when returned cancel funcition is called.
// Cancel functions need to be cancelled in reverse they were pushed.
// This can be used to keep track of contexts for nested REPL:s were you only want to cancel
// the current active "top" REPL.
package ctxstack

import (
	"context"
)

// Stack is a context stack
type Stack struct {
	cancelFns []func()
	closeCh   chan struct{}
}

// New ctxstack.Stack
func New(triggerCh func(closeCh chan struct{})) *Stack {
	closeCh := make(chan struct{})
	s := &Stack{closeCh: closeCh}

	go func() {
		for {
			triggerCh(closeCh)
			select {
			case <-closeCh:
			default:
				s.cancelFns[len(s.cancelFns)-1]()
				continue
			}
			break
		}
	}()

	return s
}

// Stop context stack
func (s *Stack) Stop() {
	close(s.closeCh)
}

// Push a new context and return it. Cancel to pop it.
func (s *Stack) Push(parent context.Context) (context.Context, func()) {
	stackCtx, stackCtxCancel := context.WithCancel(parent)
	stackIdx := len(s.cancelFns)
	s.cancelFns = append(s.cancelFns, stackCtxCancel)

	return stackCtx, func() {
		if stackIdx != len(s.cancelFns)-1 {
			panic("cancelled in wrong order")
		}
		s.cancelFns = s.cancelFns[0 : len(s.cancelFns)-1]
		stackCtxCancel()
	}
}
