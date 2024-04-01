// Package colorjson is gojq:s cli/encoder.go extract to be reusable and have non-global color config
// TODO: possible gojq can export it?
//
// The MIT License (MIT)
// Copyright (c) 2019-2022 itchyny
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
package colorjson

import (
	"bytes"
	"cmp"
	"fmt"
	"io"
	"math"
	"math/big"
	"slices"
	"strconv"
	"unicode/utf8"
)

type Colors struct {
	Reset     []byte
	Null      []byte
	False     []byte
	True      []byte
	Number    []byte
	String    []byte
	ObjectKey []byte
	Array     []byte
	Object    []byte
}

type Options struct {
	Color   bool
	Tab     bool
	Indent  int
	ValueFn func(v any) (any, error)
	Colors  Colors
}

type Encoder struct {
	out   io.Writer
	w     *bytes.Buffer
	depth int
	buf   [64]byte

	opts Options
}

func NewEncoder(opts Options) *Encoder {
	// reuse the buffer in multiple calls of marshal
	return &Encoder{
		w:    new(bytes.Buffer),
		opts: opts,
	}
}

func (e *Encoder) flush() error {
	_, err := e.out.Write(e.w.Bytes())
	e.w.Reset()
	return err
}

func (e *Encoder) Marshal(v any, w io.Writer) error {
	e.out = w
	err := e.encode(v)
	if ferr := e.flush(); ferr != nil && err == nil {
		err = ferr
	}
	return err
}

func (e *Encoder) encode(v any) error {
	switch v := v.(type) {
	case nil:
		e.write([]byte("null"), e.opts.Colors.Null)
	case bool:
		if v {
			e.write([]byte("true"), e.opts.Colors.True)
		} else {
			e.write([]byte("false"), e.opts.Colors.False)
		}
	case int:
		e.write(strconv.AppendInt(e.buf[:0], int64(v), 10), e.opts.Colors.Number)
	case float64:
		e.encodeFloat64(v)
	case *big.Int:
		e.write(v.Append(e.buf[:0], 10), e.opts.Colors.Number)
	case string:
		e.encodeString(v, e.opts.Colors.String)
	case []any:
		if err := e.encodeArray(v); err != nil {
			return err
		}
	case map[string]any:
		if err := e.encodeMap(v); err != nil {
			return err
		}
	case error:
		// value we're trying to encode is an error
		// this can happen if ValueFn is used and it reads from reader that gets cancelled etc
		return v
	default:
		if e.opts.ValueFn == nil {
			panic(fmt.Sprintf("unknown type and to ValueFn set: %[1]T (%[1]v)", v))
		}
		vv, err := e.opts.ValueFn(v)
		if err != nil {
			return err
		}
		return e.encode(vv)
	}
	if e.w.Len() > 8*1024 {
		return e.flush()
	}
	return nil
}

// ref: floatEncoder in encoding/json
func (e *Encoder) encodeFloat64(f float64) {
	if math.IsNaN(f) {
		e.write([]byte("null"), e.opts.Colors.Null)
		return
	}
	if f >= math.MaxFloat64 {
		f = math.MaxFloat64
	} else if f <= -math.MaxFloat64 {
		f = -math.MaxFloat64
	}
	format := byte('f')
	if x := math.Abs(f); x != 0 && x < 1e-6 || x >= 1e21 {
		format = 'e'
	}
	buf := strconv.AppendFloat(e.buf[:0], f, format, -1, 64)
	if format == 'e' {
		// clean up e-09 to e-9
		if n := len(buf); n >= 4 && buf[n-4] == 'e' && buf[n-3] == '-' && buf[n-2] == '0' {
			buf[n-2] = buf[n-1]
			buf = buf[:n-1]
		}
	}
	e.write(buf, e.opts.Colors.Number)
}

// ref: encodeState#string in encoding/json
func (e *Encoder) encodeString(s string, color []byte) {
	if color != nil {
		e.setColor(e.w, color)
	}
	e.w.WriteByte('"')
	start := 0
	for i := 0; i < len(s); {
		if b := s[i]; b < utf8.RuneSelf {
			if ' ' <= b && b <= '~' && b != '"' && b != '\\' {
				i++
				continue
			}
			if start < i {
				e.w.WriteString(s[start:i])
			}
			switch b {
			case '"':
				e.w.WriteString(`\"`)
			case '\\':
				e.w.WriteString(`\\`)
			case '\b':
				e.w.WriteString(`\b`)
			case '\f':
				e.w.WriteString(`\f`)
			case '\n':
				e.w.WriteString(`\n`)
			case '\r':
				e.w.WriteString(`\r`)
			case '\t':
				e.w.WriteString(`\t`)
			default:
				const hex = "0123456789abcdef"
				e.w.WriteString(`\u00`)
				e.w.WriteByte(hex[b>>4])
				e.w.WriteByte(hex[b&0xF])
			}
			i++
			start = i
			continue
		}
		c, size := utf8.DecodeRuneInString(s[i:])
		if c == utf8.RuneError && size == 1 {
			if start < i {
				e.w.WriteString(s[start:i])
			}
			e.w.WriteString(`\ufffd`)
			i += size
			start = i
			continue
		}
		i += size
	}
	if start < len(s) {
		e.w.WriteString(s[start:])
	}
	e.w.WriteByte('"')
	if color != nil {
		e.setColor(e.w, e.opts.Colors.Reset)
	}
}

func (e *Encoder) encodeArray(vs []any) error {
	e.writeByte('[', e.opts.Colors.Array)
	e.depth += e.opts.Indent
	for i, v := range vs {
		if i > 0 {
			e.writeByte(',', e.opts.Colors.Array)
		}
		if e.opts.Indent != 0 {
			e.writeIndent()
		}
		if err := e.encode(v); err != nil {
			return err
		}
	}
	e.depth -= e.opts.Indent
	if len(vs) > 0 && e.opts.Indent != 0 {
		e.writeIndent()
	}
	e.writeByte(']', e.opts.Colors.Array)
	return nil
}

func (e *Encoder) encodeMap(vs map[string]any) error {
	e.writeByte('{', e.opts.Colors.Object)
	e.depth += e.opts.Indent
	type keyVal struct {
		key string
		val any
	}
	kvs := make([]keyVal, len(vs))
	var i int
	for k, v := range vs {
		kvs[i] = keyVal{k, v}
		i++
	}
	slices.SortFunc(kvs, func(a, b keyVal) int {
		return cmp.Compare(a.key, b.key)
	})
	for i, kv := range kvs {
		if i > 0 {
			e.writeByte(',', e.opts.Colors.Object)
		}
		if e.opts.Indent != 0 {
			e.writeIndent()
		}
		e.encodeString(kv.key, e.opts.Colors.ObjectKey)
		e.writeByte(':', e.opts.Colors.Object)
		if e.opts.Indent != 0 {
			e.w.WriteByte(' ')
		}
		if err := e.encode(kv.val); err != nil {
			return err
		}
	}
	e.depth -= e.opts.Indent
	if len(vs) > 0 && e.opts.Indent != 0 {
		e.writeIndent()
	}
	e.writeByte('}', e.opts.Colors.Object)
	return nil
}

func (e *Encoder) writeIndent() {
	e.w.WriteByte('\n')
	if n := e.depth; n > 0 {
		if e.opts.Tab {
			e.writeIndentInternal(n, "\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t")
		} else {
			e.writeIndentInternal(n, "                                ")
		}
	}
}

func (e *Encoder) writeIndentInternal(n int, spaces string) {
	if l := len(spaces); n <= l {
		e.w.WriteString(spaces[:n])
	} else {
		e.w.WriteString(spaces)
		for n -= l; n > 0; n, l = n-l, l*2 {
			if n < l {
				l = n
			}
			e.w.Write(e.w.Bytes()[e.w.Len()-l:])
		}
	}
}

func (e *Encoder) writeByte(b byte, color []byte) {
	if color == nil {
		e.w.WriteByte(b)
	} else {
		e.setColor(e.w, color)
		e.w.WriteByte(b)
		e.setColor(e.w, e.opts.Colors.Reset)
	}
}

func (e *Encoder) write(bs []byte, color []byte) {
	if color == nil {
		e.w.Write(bs)
	} else {
		e.setColor(e.w, color)
		e.w.Write(bs)
		e.setColor(e.w, e.opts.Colors.Reset)
	}
}

func (e *Encoder) setColor(buf *bytes.Buffer, color []byte) {
	if e.opts.Color {
		buf.Write(color)
	}
}
