package hexdump

import (
	"fq/internal/asciiwriter"
	"fq/internal/columnwriter"
	"fq/internal/hexpairwriter"
	"fq/internal/num"
	"io"
	"strings"
)

type Dumper struct {
	addrLen     int
	lineBytes   int64
	columnW     *columnwriter.Writer
	separatorsW io.Writer
	startOffset int64
	offset      int64
}

func New(w io.Writer, startOffset int64, addrLen int, lineBytes int) *Dumper {
	cw := columnwriter.New(w, []int{addrLen, 1, lineBytes*3 - 1, 1, lineBytes, 1})
	return &Dumper{
		addrLen:     addrLen,
		lineBytes:   int64(lineBytes),
		columnW:     cw,
		separatorsW: io.MultiWriter(cw.Columns[1], cw.Columns[3], cw.Columns[5]),
		startOffset: startOffset,
		offset:      startOffset - startOffset%int64(lineBytes),
	}
}

func (d *Dumper) flush() error {
	if _, err := d.columnW.Columns[0].Write([]byte(
		num.PadFormatInt(((d.offset-1)/d.lineBytes)*d.lineBytes, 16, d.addrLen))); err != nil {
		return err
	}
	if _, err := d.separatorsW.Write([]byte("|")); err != nil {
		return err
	}
	if err := d.columnW.Flush(); err != nil {
		return err
	}
	return nil
}

func (d *Dumper) Write(p []byte) (n int, err error) {
	if d.offset < d.startOffset {
		r := int(d.startOffset - d.offset)
		if _, err := d.columnW.Columns[2].Write([]byte(strings.Repeat("   ", r))); err != nil {
			return n, err
		}
		if _, err := d.columnW.Columns[4].Write([]byte(strings.Repeat(" ", r))); err != nil {
			return n, err
		}
		d.offset = d.startOffset
	}

	for len(p) > 0 {
		cl := d.lineBytes - d.offset%d.lineBytes
		if cl == 0 {
			cl = d.lineBytes
		}
		if cl > int64(len(p)) {
			cl = int64(len(p))
		}

		ps := p[0:cl]
		for _, b := range ps {
			d.offset++

			if _, err := d.columnW.Columns[2].Write([]byte(hexpairwriter.Pair(b))); err != nil {
				return n, err
			}
			if d.offset%d.lineBytes != 0 {
				if _, err := d.columnW.Columns[2].Write([]byte(" ")); err != nil {
					return n, err
				}
			}
			if _, err := d.columnW.Columns[4].Write([]byte(asciiwriter.SafeASCII(b))); err != nil {
				return n, err
			}
			n++
		}

		if d.offset%d.lineBytes == 0 {
			if err := d.flush(); err != nil {
				return n, err
			}
		}

		p = p[cl:]
	}

	return n, nil
}

func (d *Dumper) Close() error {
	if d.offset == 0 || d.offset%d.lineBytes != 0 {
		return d.flush()
	}
	return nil
}
