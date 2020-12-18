package hexdump

import (
	"fq/internal/asciiwriter"
	"fq/internal/columnwriter"
	"fq/internal/hexpairwriter"
	"fq/internal/num"
	"io"
)

type Dumper struct {
	addrLen   int
	lineBytes int64
	cw        *columnwriter.Writer
	pw        io.Writer
	n         int64
}

func New(w io.Writer, addrLen int, lineBytes int) *Dumper {
	cw := columnwriter.New(w, []int{addrLen, 1, lineBytes*3 - 1, 1, lineBytes, 1})
	return &Dumper{
		addrLen:   addrLen,
		lineBytes: int64(lineBytes),
		cw:        cw,
		pw:        io.MultiWriter(cw.Columns[1], cw.Columns[3], cw.Columns[5]),
	}
}

func (d *Dumper) flush() error {
	if _, err := d.cw.Columns[0].Write([]byte(
		num.PadFormatInt(((d.n-1)/d.lineBytes)*d.lineBytes, 16, d.addrLen))); err != nil {
		return err
	}
	if _, err := d.pw.Write([]byte("|")); err != nil {
		return err
	}
	if err := d.cw.Flush(); err != nil {
		return err
	}
	return nil
}

func (d *Dumper) Write(p []byte) (n int, err error) {
	for len(p) > 0 {
		cl := d.n % d.lineBytes
		if cl == 0 {
			cl = d.lineBytes
		}
		if cl > int64(len(p)) {
			cl = int64(len(p))
		}

		ps := p[0:cl]
		for _, b := range ps {
			d.n++

			if _, err := d.cw.Columns[2].Write([]byte(hexpairwriter.Pair(b))); err != nil {
				return n, err
			}
			if d.n%d.lineBytes != 0 {
				if _, err := d.cw.Columns[2].Write([]byte(" ")); err != nil {
					return n, err
				}
			}
			if _, err := d.cw.Columns[4].Write([]byte(asciiwriter.SafeASCII(b))); err != nil {
				return n, err
			}
			n++
		}

		if d.n%d.lineBytes == 0 {
			if err := d.flush(); err != nil {
				return n, err
			}
		}

		p = p[cl:]
	}

	return n, nil
}

func (d *Dumper) Close() error {
	if d.n == 0 || d.n%d.lineBytes != 0 {
		return d.flush()
	}
	return nil
}
