package columnwriter

import (
	"bytes"
	"io"
	"strings"
)

type column struct {
	width int
	lines []string
	buf   bytes.Buffer
}

type Writer struct {
	columns []column
	current int
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
	var columns []column
	for _, w := range widths {
		columns = append(columns, column{width: w})
	}

	return &Writer{
		columns: columns,
		current: 0,
		w:       w,
	}
}

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

func (w *Writer) Row() {
	maxLines := 0
	for _, c := range w.columns {
		if len(c.lines) > maxLines {
			maxLines = len(c.lines)
		}
	}

	for i := 0; i < maxLines; i++ {
		for _, c := range w.columns {
			var s string
			if i < len(c.lines) {
				s = c.lines[i]
			}

			if c.width != -1 {
				if len(s) > c.width {
					s = s[0:c.width]
				} else if len(s) < c.width {
					s += strings.Repeat(" ", c.width-len(s))
				}
			}

			w.w.Write([]byte(s))
		}
		w.w.Write([]byte{'\n'})
	}

	for _, c := range w.columns {
		c.lines = nil
		c.buf.Reset()
	}
}
