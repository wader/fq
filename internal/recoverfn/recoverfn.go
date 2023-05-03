package recoverfn

import (
	"runtime"
)

const stackSizeLimit = 256

type Raw struct {
	RecoverV  any
	RecoverPC uintptr
	PCs       []uintptr
}

type RecoverableErrorer interface {
	IsRecoverableError() bool
}

// Run runs fn and return Raw{}, true for no panic
// If panic is recoverable raw stacktrace and panic value is returned
// If panic is not recoverable we just panic again
func Run(fn func()) (Raw, bool) {
	// TODO: once?
	var recoverPC [1]uintptr
	runtime.Callers(1, recoverPC[:])

	pc, v := func() (pcs []uintptr, v any) {
		defer func() {
			if recoverV := recover(); recoverV != nil {
				if re, ok := recoverV.(RecoverableErrorer); ok && re.IsRecoverableError() {
					pcs = make([]uintptr, stackSizeLimit)
					pcs = pcs[0:runtime.Callers(0, pcs)]
					v = recoverV
					return
				}
				panic(recoverV)
			}
		}()

		fn()

		return nil, nil
	}()

	if v == nil {
		return Raw{}, true
	}

	return Raw{
		RecoverV:  v,
		RecoverPC: recoverPC[0],
		PCs:       pc,
	}, false
}

func (r Raw) frames(startSkip int, bottomSkip int, bottomPC uintptr) []runtime.Frame {
	var bottomFrame runtime.Frame
	bottomIndex := -1
	if bottomPC != 0 {
		bottomPCs := [1]uintptr{bottomPC}
		bottomFrame, _ = runtime.CallersFrames(bottomPCs[:]).Next()
	}

	fs := make([]runtime.Frame, len(r.PCs))
	frames := runtime.CallersFrames(r.PCs)
	for i := 0; ; i++ {
		f, more := frames.Next()
		if !more {
			break
		}
		if bottomPC != 0 && f.Function == bottomFrame.Function {
			bottomIndex = i
		}
		fs[i] = f
	}

	endIndex := len(fs) - 1
	if bottomIndex != -1 {
		endIndex = bottomIndex - bottomSkip
	}

	return fs[startSkip:endIndex]
}

func (r Raw) Frames() []runtime.Frame {
	// 3 to skip runtime.Callers, Recover help function and runtime.gopanic
	// 1 to skip Recover defer recover() function
	return r.frames(3, 1, r.RecoverPC)
}
