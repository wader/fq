package interp

import (
	"fmt"
	"fq/pkg/decode"
	"io"

	"github.com/itchyny/gojq"
)

var _ gojq.JQValue = (*decodeError)(nil)
var _ Display = (*decodeError)(nil)

type decodeError struct {
	v decode.DecodeFormatsError
}

func (de decodeError) Display(w io.Writer, opts Options) error {
	if !opts.Verbose {
		fmt.Fprintf(w, "Failed to decode, try -d <format> to force format\n")
		return nil
	}

	for _, err := range de.v.Errs {
		fmt.Fprintf(w, "%s: %s\n", err.Format.Name, err.Err.Error())
		for _, f := range err.Stacktrace.Frames() {
			fmt.Fprintf(w, "%s\n", f.Function)
			fmt.Fprintf(w, "  %s:%d\n", f.File, f.Line)
		}
	}
	return nil
}

func (de decodeError) JQValueLength() interface{} {
	return len(de.v.Errs)
}
func (de decodeError) JQValueIndex(index int) interface{} {
	if index < 0 || index >= len(de.v.Errs) {
		return nil
	}
	return formatError{de.v.Errs[index]}
}

func (de decodeError) JQValueSlice(start int, end int) interface{} {
	return nil
}
func (de decodeError) JQValueKey(name string) interface{} {
	return fmt.Errorf("can't index array with string")
}
func (de decodeError) JQValueEach() interface{} {
	var props []gojq.PathValue
	for i, e := range de.v.Errs {
		props = append(props, gojq.PathValue{Path: i, Value: formatError{e}})
	}
	return props
}
func (de decodeError) JQValueType() string {
	return "array"
}

func (de decodeError) JQValueKeys() interface{} {
	var is []interface{}
	for i := range de.v.Errs {
		is = append(is, i)
	}
	return is
}

func (de decodeError) JQValueHasKey(key interface{}) interface{} {
	i, ok := key.(int)
	if !ok {
		return fmt.Errorf("cannot index array with %v", key)
	}
	return i < 0 || i >= len(de.v.Errs)
}

func (de decodeError) JQValueToNumber() interface{} {
	return funcTypeError{name: "tonumber", typ: "decodeError"}
}

func (de decodeError) JQValueToString() interface{} {
	return funcTypeError{name: "tostring", typ: "decodeError"}
}

func (de decodeError) JQValue() interface{} {
	var ea []interface{}
	for _, e := range de.v.Errs {
		ea = append(ea, formatError{e}.JQValue())
	}

	return ea
}

var _ gojq.JQValue = (*formatError)(nil)

type formatError struct {
	v decode.FormatError
}

func (fe formatError) JQValueLength() interface{} {
	return 3
}
func (fe formatError) JQValueIndex(index int) interface{} {
	return fmt.Errorf("can't index object")
}
func (fe formatError) JQValueSlice(start int, end int) interface{} {
	return fmt.Errorf("can't slice object")
}
func (fe formatError) JQValueKey(name string) interface{} {
	switch name {
	case "format":
		return fe.v.Format.Name
	case "error":
		return fe.v.Err.Error()
	case "stacktrace":
		var st []interface{}
		for _, f := range fe.v.Stacktrace.Frames() {
			st = append(st, f.Function)
		}
		return st
	}
	return nil
}
func (fe formatError) JQValueEach() interface{} {
	return []gojq.PathValue{
		{Path: "format", Value: fe.JQValueKey("format")},
		{Path: "error", Value: fe.JQValueKey("error")},
		{Path: "stacktrace", Value: fe.JQValueKey("stacktrace")},
	}
}
func (fe formatError) JQValueType() string {
	return "object"
}

func (fe formatError) JQValueKeys() interface{} {
	return []interface{}{"format", "error", "stackrace"}
}

func (fe formatError) JQValueHasKey(key interface{}) interface{} {
	s, ok := key.(string)
	if !ok {
		return fmt.Errorf("cannot index object with %v", key)
	}
	for _, e := range []string{"format", "error", "stackrace"} {
		if s == e {
			return true
		}
	}
	return false
}

func (fe formatError) JQValueToNumber() interface{} {
	return funcTypeError{name: "tonumber", typ: "formatError"}
}

func (fe formatError) JQValueToString() interface{} {
	return funcTypeError{name: "tostring", typ: "formatError"}
}

func (fe formatError) JQValue() interface{} {
	return map[string]interface{}{
		"format":     fe.JQValueKey("format"),
		"error":      fe.JQValueKey("error"),
		"stacktrace": fe.JQValueKey("stacktrace"),
	}
}
