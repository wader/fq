package apple

import "github.com/wader/fq/internal/mathx"

// PosLoopDetector is used for detecting loops when writing decoders, and can
// short-circuit infinite recursion that can cause stack overflows.
type PosLoopDetector[T mathx.Integer] []T

// Push adds the current offset to the stack and executes the supplied
// detection function
func (pld *PosLoopDetector[T]) Push(offset T, detect func()) {
	for _, o := range *pld {
		if offset == o {
			detect()
		}
	}
	*pld = append(*pld, offset)
}

// Pop removes the most recently added offset from the stack.
func (pld *PosLoopDetector[T]) Pop() {
	*pld = (*pld)[:len(*pld)-1]
}

// PushAndPop adds the current offset to the stack, executes the supplied
// detection function, and returns the Pop method. A good usage of this is to
// pair this method call with a defer statement.
func (pld *PosLoopDetector[T]) PushAndPop(offset T, detect func()) func() {
	pld.Push(offset, detect)
	return pld.Pop
}
