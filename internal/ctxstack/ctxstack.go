// ctxstack manages a stack of contexts. When triggerFn returns and closeCh is not closed
// it will cancel the top context. Stack is popped when returned cancel funcition is called.
// Cancel functions need to be cancelled in reverse they were pushed.
// This can be used to keep track of contexts for nested REPL:s
package ctxstack

import (
	"context"
)

type Stack struct {
	cancelFns []func()
	closeCh   chan struct{}
}

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

func (s *Stack) Close() {
	close(s.closeCh)
}

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
