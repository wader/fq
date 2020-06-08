package columnwriter

import (
	"bytes"
	"fq/internal/ansi"
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

// Writer maintins multiple column io.Writer:s. On Flush() row align them.
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

func (w *Writer) Flush() error {
	for _, c := range w.Columns {
		c.Flush()
	}

	maxLines := 0
	for _, c := range w.Columns {
		lenLines := len(c.Lines)
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
				l := ansi.Len(s)
				if l < c.Width {
					s += strings.Repeat(" ", c.Width-l)
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
