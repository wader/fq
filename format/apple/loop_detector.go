package apple

// PosLoopDetector is used for detecting loops when writing decoders, and can
// short-circuit infinite recursion that can cause stack overflows.
type PosLoopDetector []int64

// Push adds the current offset to the stack and executes the supplied
// detection function
func (pld *PosLoopDetector) Push(offset int64, detect func()) {
	for _, o := range *pld {
		if offset == o {
			detect()
		}
	}
	*pld = append(*pld, offset)
}

// Pop removes the most recently added offset from the stack.
func (pld *PosLoopDetector) Pop() {
	*pld = (*pld)[:len(*pld)-1]
}

// PushAndPop adds the current offset to the stack, executes the supplied
// detection function, and returns the Pop method. A good usage of this is to
// pair this method call with a defer statement.
func (pld *PosLoopDetector) PushAndPop(i int64, detect func()) func() {
	pld.Push(i, detect)
	return pld.Pop
}
