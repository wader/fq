package decode

import (
	"fmt"
	"strings"

	"github.com/wader/fq/internal/mathx"
	"github.com/wader/fq/internal/recoverfn"
)

type RecoverableErrorer interface {
	IsRecoverableError() bool
}

type FormatError struct {
	Err        error
	Format     *Format
	Stacktrace recoverfn.Raw
}

type FormatsError struct {
	Errs []FormatError
}

func (fe FormatsError) Error() string {
	var errs []string
	for _, err := range fe.Errs {
		errs = append(errs, err.Error())
	}
	return strings.Join(errs, ", ")
}

func (fe FormatError) Error() string {
	// var fns []string
	// for _, f := range fe.Stacktrace.Frames() {
	// 	fns = append(fns, fmt.Sprintf("%s:%d:%s", f.File, f.Line, f.Function))
	// }

	return fe.Err.Error()
}

func (fe FormatError) Value() any {
	var st []any
	for _, f := range fe.Stacktrace.Frames() {
		st = append(st, f.Function)
	}

	return map[string]any{
		"format":     fe.Format.Name,
		"error":      fe.Err.Error(),
		"stacktrace": st,
	}
}

func (FormatsError) IsRecoverableError() bool { return true }

type IOError struct {
	Err      error
	Name     string
	Op       string
	ReadSize int64
	SeekPos  int64
	Pos      int64
}

func (e IOError) Error() string {
	var prefix string
	if e.Name != "" {
		prefix = e.Op + "(" + e.Name + ")"
	} else {
		prefix = e.Op
	}

	return fmt.Sprintf("%s: failed at position %s (read size %s seek pos %s): %s",
		prefix, mathx.Bits(e.Pos).StringByteBits(10), mathx.Bits(e.ReadSize).StringByteBits(10), mathx.Bits(e.SeekPos).StringByteBits(10), e.Err)
}
func (e IOError) Unwrap() error { return e.Err }

func (IOError) IsRecoverableError() bool { return true }

type DecoderError struct {
	Reason string
	Pos    int64
}

func (e DecoderError) Error() string {
	return fmt.Sprintf("error at position %s: %s", mathx.Bits(e.Pos).StringByteBits(16), e.Reason)
}

func (DecoderError) IsRecoverableError() bool { return true }
