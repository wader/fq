package columnwriter

import (
	"bytes"
	"io"
	"unicode/utf8"
)

type Column interface {
	io.Writer
	Lines() int
	PreFlush()
	FlushLine(w io.Writer, lines int, lastColumn bool) error
	Reset()
}

var _ Column = (*MultiLineColumn)(nil)

type MultiLineColumn struct {
	Width   int
	Wrap    bool
	LenFn   func(s string) int
	SliceFn func(s string, start, stop int) string

	lines []string
	buf   bytes.Buffer
}

func (c *MultiLineColumn) divideString(s string, l int) []string {
	var ss []string
	parts := c.lenFn(s) / l
	for i := 0; i < parts; i++ {
		ss = append(ss, c.sliceFn(s, i*l, (i+1)*l))
	}
	if len(s)%l != 0 {
		ss = append(ss, c.sliceFn(s, parts*l, -1))
	}

	return ss
}

// TODO: fn assume fixed width runes
func (c *MultiLineColumn) lenFn(s string) int {
	if c.LenFn != nil {
		return c.LenFn(s)
	}
	return utf8.RuneCountInString(s)
}

func (c *MultiLineColumn) sliceFn(s string, start, stop int) string {
	if c.LenFn != nil {
		return c.SliceFn(s, start, stop)
	}
	if stop == -1 {
		return string(([]rune(s))[start:])
	}
	return string(([]rune(s))[start:stop])
}

func (c *MultiLineColumn) Write(p []byte) (int, error) {
	bb := &c.buf

	bb.Write(p)

	b := bb.Bytes()
	pos := 0

	for {
		i := indexByteSet(b[pos:], []byte{'\n'})
		if i < 0 {
			break
		}

		line := string([]rune(string(b[pos : pos+i])))
		if c.Wrap && c.Width != -1 && c.lenFn(line) > c.Width {
			c.lines = append(c.lines, c.divideString(line, c.Width)...)
		} else {
			c.lines = append(c.lines, line)
		}

		pos += i + 1
	}
	bb.Reset()
	bb.Write(b[pos:])

	return len(p), nil
}

func (c *MultiLineColumn) Lines() int { return len(c.lines) }

func (c *MultiLineColumn) PreFlush() {
	if c.buf.Len() > 0 {
		_, _ = c.Write([]byte{'\n'})
	}
}

func (c *MultiLineColumn) FlushLine(w io.Writer, lineNr int, lastColumn bool) error {
	var s string
	if lineNr < len(c.lines) {
		s = c.lines[lineNr]
		if c.Width != -1 && c.lenFn(s) > c.Width {
			s = c.sliceFn(s, 0, c.Width)
		}
	}

	if _, err := w.Write([]byte(s)); err != nil {
		return err
	}

	if !lastColumn && c.Width != -1 {
		l := c.lenFn(s)
		if l < c.Width {
			n := c.Width - l
			for n > 0 {
				const whitespace = "                                                                                "
				r := n
				if r > len(whitespace) {
					r = len(whitespace)
				}
				n -= r

				if _, err := w.Write([]byte(whitespace[0:r])); err != nil {
					return err
				}
			}

		}
	}

	return nil
}

func (c *MultiLineColumn) Reset() {
	c.lines = nil
	c.buf.Reset()
}

type BarColumn string

var _ Column = (*BarColumn)(nil)

func (c BarColumn) Write(p []byte) (int, error) { return len(p), nil } // TODO: can be removed?
func (c BarColumn) Lines() int                  { return 1 }
func (c BarColumn) PreFlush()                   {}
func (c BarColumn) FlushLine(w io.Writer, lineNr int, lastColumn bool) error {
	if _, err := w.Write([]byte(c)); err != nil {
		return err
	}
	return nil
}
func (c BarColumn) Reset() {}

// Writer maintins multiple column io.Writer:s. On Flush() row align them.
type Writer struct {
	Columns []Column

	w io.Writer
}

func indexByteSet(s []byte, cs []byte) int {
	ri := -1

	for _, c := range cs {
		i := bytes.IndexByte(s, c)
		if i != -1 && (ri == -1 || i < ri) {
			ri = i
		}
	}

	return ri
}

func New(w io.Writer, columns ...Column) *Writer {
	return &Writer{
		Columns: columns,
		w:       w,
	}
}

func (w *Writer) Flush() error {
	maxLines := 0
	for _, c := range w.Columns {
		l := c.Lines()
		if l > maxLines {
			maxLines = l
		}
	}

	for _, c := range w.Columns {
		c.PreFlush()
	}

	for line := 0; line < maxLines; line++ {
		for ci, c := range w.Columns {
			if err := c.FlushLine(w.w, line, ci == len(w.Columns)-1); err != nil {
				return err
			}
		}
		if _, err := w.w.Write([]byte{'\n'}); err != nil {
			return err
		}
	}

	for _, c := range w.Columns {
		c.Reset()
	}

	return nil
}
