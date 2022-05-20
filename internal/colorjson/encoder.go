// This is gojq:s cli/encoder.go extract to be reusable and non-global color config
// TODO: possible gojq can export it?
//
// The MIT License (MIT)
// Copyright (c) 2019-2021 itchyny
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
//

// skip errcheck to keep code similar to gojq version
//nolint:errcheck
package colorjson

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"math/big"
	"sort"
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

type Encoder struct {
	w       *bufio.Writer
	wErr    error
	color   bool
	tab     bool
	indent  int
	depth   int
	buf     [64]byte
	valueFn func(v any) any
	colors  Colors
}

func NewEncoder(color bool, tab bool, indent int, valueFn func(v any) any, colors Colors) *Encoder {
	// reuse the buffer in multiple calls of marshal
	return &Encoder{
		color:   color,
		tab:     tab,
		indent:  indent,
		valueFn: valueFn,
		colors:  colors,
	}
}

func (e *Encoder) Marshal(v any, w io.Writer) error {
	e.w = bufio.NewWriter(w)
	e.encode(v)
	if e.wErr != nil {
		return e.wErr
	}
	return e.w.Flush()
}

func (e *Encoder) encode(v any) {
	switch v := v.(type) {
	case nil:
		e.write([]byte("null"), e.colors.Null)
	case bool:
		if v {
			e.write([]byte("true"), e.colors.True)
		} else {
			e.write([]byte("false"), e.colors.False)
		}
	case int:
		e.write(strconv.AppendInt(e.buf[:0], int64(v), 10), e.colors.Number)
	case float64:
		e.encodeFloat64(v)
	case *big.Int:
		e.write(v.Append(e.buf[:0], 10), e.colors.Number)
	case string:
		e.encodeString(v, e.colors.String)
	case []any:
		e.encodeArray(v)
	case map[string]any:
		e.encodeMap(v)
	default:
		if e.valueFn != nil {
			v = e.valueFn(v)
		} else {
			panic(fmt.Sprintf("invalid value: %#+v", v))
		}
		e.encode(v)
	}
}

// ref: floatEncoder in encoding/json
func (e *Encoder) encodeFloat64(f float64) {
	if math.IsNaN(f) {
		e.write([]byte("null"), e.colors.Null)
		return
	}
	if f >= math.MaxFloat64 {
		f = math.MaxFloat64
	} else if f <= -math.MaxFloat64 {
		f = -math.MaxFloat64
	}
	fmt := byte('f')
	if x := math.Abs(f); x != 0 && x < 1e-6 || x >= 1e21 {
		fmt = 'e'
	}
	buf := strconv.AppendFloat(e.buf[:0], f, fmt, -1, 64)
	if fmt == 'e' {
		// clean up e-09 to e-9
		if n := len(buf); n >= 4 && buf[n-4] == 'e' && buf[n-3] == '-' && buf[n-2] == '0' {
			buf[n-2] = buf[n-1]
			buf = buf[:n-1]
		}
	}
	e.write(buf, e.colors.Number)
}

// ref: encodeState#string in encoding/json
func (e *Encoder) encodeString(s string, color []byte) {
	if e.color {
		e.w.Write(color)
	}
	e.w.WriteByte('"')
	start := 0
	for i := 0; i < len(s); {
		if b := s[i]; b < utf8.RuneSelf {
			if ']' <= b && b <= '~' || '#' <= b && b <= '[' || b == ' ' || b == '!' {
				i++
				continue
			}
			if start < i {
				e.w.WriteString(s[start:i])
			}
			e.w.WriteByte('\\')
			switch b {
			case '\\', '"':
				e.w.WriteByte(b)
			case '\b':
				e.w.WriteByte('b')
			case '\f':
				e.w.WriteByte('f')
			case '\n':
				e.w.WriteByte('n')
			case '\r':
				e.w.WriteByte('r')
			case '\t':
				e.w.WriteByte('t')
			default:
				const hex = "0123456789abcdef"
				e.w.WriteString("u00")
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
	if e.color {
		e.w.Write(e.colors.Reset)
	}
}

func (e *Encoder) encodeArray(vs []any) {
	e.writeByte('[', e.colors.Array)
	e.depth += e.indent
	for i, v := range vs {
		if e.wErr != nil {
			return
		}

		if i > 0 {
			e.writeByte(',', e.colors.Array)
		}
		if e.indent != 0 {
			e.writeIndent()
		}
		e.encode(v)
	}
	e.depth -= e.indent
	if len(vs) > 0 && e.indent != 0 {
		e.writeIndent()
	}
	e.writeByte(']', e.colors.Array)
}

func (e *Encoder) encodeMap(vs map[string]any) {
	e.writeByte('{', e.colors.Object)
	e.depth += e.indent
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
	sort.Slice(kvs, func(i, j int) bool {
		return kvs[i].key < kvs[j].key
	})
	for i, kv := range kvs {
		if e.wErr != nil {
			return
		}

		if i > 0 {
			e.writeByte(',', e.colors.Object)
		}
		if e.indent != 0 {
			e.writeIndent()
		}
		e.encodeString(kv.key, e.colors.ObjectKey)
		e.writeByte(':', e.colors.Object)
		if e.indent != 0 {
			e.w.WriteByte(' ')
		}
		e.encode(kv.val)
	}
	e.depth -= e.indent
	if len(vs) > 0 && e.indent != 0 {
		e.writeIndent()
	}
	e.writeByte('}', e.colors.Object)
}

func (e *Encoder) writeIndent() {
	e.w.WriteByte('\n')
	if n := e.depth; n > 0 {
		if e.tab {
			const tabs = "\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t"
			for n > len(tabs) {
				e.w.Write([]byte(tabs))
				n -= len(tabs)
			}
			e.w.Write([]byte(tabs)[:n])
		} else {
			const spaces = "                                                                "
			for n > len(spaces) {
				e.w.Write([]byte(spaces))
				n -= len(spaces)
			}
			e.w.Write([]byte(spaces)[:n])
		}
	}
}

func (e *Encoder) writeByte(b byte, color []byte) {
	if e.wErr != nil {
		return
	}
	if color == nil {
		if err := e.w.WriteByte(b); err != nil {
			e.wErr = err
		}
	} else {
		if e.color {
			if _, err := e.w.Write(color); err != nil {
				e.wErr = err
			}
		}
		if err := e.w.WriteByte(b); err != nil {
			e.wErr = err
		}
		if e.color {
			if _, err := e.w.Write(e.colors.Reset); err != nil {
				e.wErr = err
			}
		}
	}
}

func (e *Encoder) write(bs []byte, color []byte) {
	if e.wErr != nil {
		return
	}
	if color == nil {
		if _, err := e.w.Write(bs); err != nil {
			e.wErr = err
		}
	} else {
		if e.color {
			if _, err := e.w.Write(color); err != nil {
				e.wErr = err
			}
		}
		if _, err := e.w.Write(bs); err != nil {
			e.wErr = err
		}
		if e.color {
			if _, err := e.w.Write(e.colors.Reset); err != nil {
				e.wErr = err
			}
		}
	}
}
