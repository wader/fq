package recoverfn

import (
	"fmt"
	"io"
	"os"
	"runtime"
)

const stackSizeLimit = 256

type Raw struct {
	RecoverV  interface{}
	RecoverPC uintptr
	PCs       []uintptr
}

// Run runs fn and return Raw{}, true on no-panic
// on panic it recovers and return a raw stacktrace and panic value to inspect
func Run(fn func()) (Raw, bool) {
	// TODO: once?
	var recoverPC [1]uintptr
	runtime.Callers(1, recoverPC[:])

	pc, v := func() (pcs []uintptr, v interface{}) {
		defer func() {
			if recoverErr := recover(); recoverErr != nil {
				pcs = make([]uintptr, stackSizeLimit)
				pcs = pcs[0:runtime.Callers(0, pcs)]
				v = recoverErr
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

func (r Raw) RePanic() {
	var o io.Writer
	o = os.Stderr
	if p := os.Getenv("REPANIC_LOG"); p != "" {
		if f, err := os.Create(p); err == nil {
			o = f
			defer f.Close()
		}
	}

	fmt.Fprintf(o, "repanic: %v\n", r.RecoverV)
	for _, f := range r.frames(0, 0, 0) {
		fmt.Fprintf(o, "%s\n", f.Function)
		fmt.Fprintf(o, "\t%s:%d\n", f.File, f.Line)
	}
	panic(r.RecoverV)
}
