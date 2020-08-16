package columnwriter

import (
	"bytes"
	"io"
	"strings"
)

type Column struct {
	Width int
	Lines []string
	Buf   bytes.Buffer
	Wrap  bool
}

func divideString(s string, l int) []string {
	var ss []string
	parts := len(s) / l
	for i := 0; i < parts; i++ {
		ss = append(ss, s[i*l:(i+1)*l])
	}
	if len(s)%l != 0 {
		ss = append(ss, s[parts*l:])
	}

	return ss
}

func (c *Column) Write(p []byte) (int, error) {
	bb := &c.Buf

	bb.Write(p)

	b := bb.Bytes()
	pos := 0

	for {
		i := indexByteSet(b[pos:], []byte{'\n'})
		if i < 0 {
			break
		}

		line := string([]rune(string(b[pos : pos+i])))
		if c.Wrap && len(line) > c.Width {
			c.Lines = append(c.Lines, divideString(line, c.Width)...)
		} else {
			c.Lines = append(c.Lines, line)
		}

		pos += i + 1
	}
	bb.Reset()
	bb.Write(b[pos:])

	return len(p), nil
}

func (c *Column) Flush() {
	if c.Buf.Len() > 0 {
		_, _ = c.Write([]byte{'\n'})
	}
}

type Writer struct {
	Columns []*Column
	w       io.Writer
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

func New(w io.Writer, widths []int) *Writer {
	var columns []*Column
	for _, w := range widths {
		columns = append(columns, &Column{Width: w})
	}

	return &Writer{
		Columns: columns,
		w:       w,
	}
}

/*
func (w *Writer) Write(p []byte) (int, error) {
	c := &w.columns[w.current]
	bb := &c.buf

	bb.Write(p)

	b := bb.Bytes()
	pos := 0

	for {
		i := indexByteSet(b[pos:], []byte{'\n'})
		if i < 0 {
			break
		}

		c.lines = append(c.lines, string([]rune(string(b[pos:pos+i]))))
		pos += i + 1
	}
	bb.Reset()
	bb.Write(b[pos:])

	return len(p), nil
}
*/

/*
func (w *Writer) Next() {
	c := &w.columns[w.current]
	if c.buf.Len() > 0 {
		w.Write([]byte{'\n'})
	}

	w.current++
	if w.current == len(w.columns) {
		// panic(fmt.Sprintf("column index %d > %d", w.current, len(w.columns)))
		w.Row()
		w.current = 0
	}
}
*/

func (w *Writer) Flush() error {
	for _, c := range w.Columns {
		c.Flush()
	}

	maxLines := 0
	for _, c := range w.Columns {
		lenLines := len(c.Lines)
		// if c.Wrap {
		// 	lenLines = 0
		// 	for _, l := range c.Lines {
		// 		lenLine := len(l)
		// 		wrappedLines := lenLines / c.Width
		// 		if lenLine%c.Width != 0 {
		// 			wrappedLines++
		// 		}
		// 		lenLines += wrappedLines
		// 	}
		// }

		if lenLines > maxLines {
			maxLines = len(c.Lines)
		}
	}

	for i := 0; i < maxLines; i++ {
		for _, c := range w.Columns {
			var s string
			if i < len(c.Lines) {
				s = c.Lines[i]
			}

			if c.Width != -1 {
				if len(s) > c.Width {
					s = s[0:c.Width]
				} else if len(s) < c.Width {
					s += strings.Repeat(" ", c.Width-len(s))
				}
			}

			if _, err := w.w.Write([]byte(s)); err != nil {
				return err
			}
		}
		if _, err := w.w.Write([]byte{'\n'}); err != nil {
			return err
		}
	}

	for _, c := range w.Columns {
		c.Lines = nil
		c.Buf.Reset()
	}

	return nil
}
