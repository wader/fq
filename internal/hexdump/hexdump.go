package hexdump

import (
	"io"
	"strings"

	"github.com/wader/fq/internal/columnwriter"
	"github.com/wader/fq/internal/mathx"
	"github.com/wader/fq/pkg/bitio"
)

type Dumper struct {
	addrLen          int
	addrBase         int
	lineBytes        int64
	columnW          *columnwriter.Writer
	separatorsW      io.Writer
	startOffset      int64
	offset           int64
	hexFn            func(b byte) string
	asciiFn          func(b byte) string
	dumpHeaderFn     func(s string) string
	dumpAddrFn       func(s string) string
	column           string
	hasWrittenHeader bool

	bitsBuf  []byte
	bitsBufN int64
}

// TODO: something more generic? bin, octal, arbitrary base?
// TODO: template for columns?
// TODO: merge with dump?
// TODO: replace addrLen with highest address and calc instead
// TODO: use dump options? config struct?
func New(w io.Writer, startOffset int64, addrLen int, addrBase int, lineBytes int,
	hexFn func(b byte) string,
	asciiFn func(b byte) string,
	dumpHeaderFn func(s string) string,
	dumpAddrFn func(s string) string,
	column string) *Dumper {
	cw := columnwriter.New(
		w,
		&columnwriter.MultiLineColumn{Width: addrLen},
		columnwriter.BarColumn(column),
		&columnwriter.MultiLineColumn{Width: lineBytes * 3},
		columnwriter.BarColumn(column),
		&columnwriter.MultiLineColumn{Width: lineBytes},
		columnwriter.BarColumn(column),
	)

	return &Dumper{
		addrLen:          addrLen,
		addrBase:         addrBase,
		lineBytes:        int64(lineBytes),
		columnW:          cw,
		separatorsW:      io.MultiWriter(cw.Columns[1], cw.Columns[3], cw.Columns[5]),
		startOffset:      startOffset,
		offset:           startOffset - startOffset%int64(lineBytes),
		hexFn:            hexFn,
		asciiFn:          asciiFn,
		dumpHeaderFn:     dumpHeaderFn,
		dumpAddrFn:       dumpAddrFn,
		column:           column,
		hasWrittenHeader: false,
		bitsBuf:          make([]byte, 1),
	}
}

func (d *Dumper) flush() error {
	if _, err := d.columnW.Columns[0].Write([]byte(
		d.dumpAddrFn(mathx.PadFormatInt(((d.offset-1)/d.lineBytes)*d.lineBytes, d.addrBase, true, d.addrLen)))); err != nil {
		return err
	}
	if _, err := d.separatorsW.Write([]byte(d.column)); err != nil {
		return err
	}
	if err := d.columnW.Flush(); err != nil {
		return err
	}
	return nil
}

func (d *Dumper) WriteBits(p []byte, nBits int64) (n int64, err error) {
	pos := int64(0)
	rBits := nBits
	if d.bitsBufN > 0 {
		r := min(8-d.bitsBufN, nBits)
		v := bitio.Read64(p, 0, r)
		bitio.Write64(v, r, d.bitsBuf, d.bitsBufN)

		d.bitsBufN += r

		if d.bitsBufN < 8 {
			return nBits, nil
		}
		if n, err := d.Write(d.bitsBuf); err != nil {
			return int64(n) * 8, err
		}
		pos = r
		rBits -= r
	}

	for rBits >= 8 {
		b := [1]byte{0}

		b[0] = byte(bitio.Read64(p, pos, 8))
		if n, err := d.Write(b[:]); err != nil {
			return int64(n) * 8, err
		}

		pos += 8
		rBits -= 8

	}

	if rBits > 0 {
		d.bitsBuf[0] = byte(bitio.Read64(p, pos, rBits)) << (8 - rBits)
		d.bitsBufN = rBits
	} else {
		d.bitsBufN = 0
	}

	return nBits, nil
}

func (d *Dumper) Write(p []byte) (n int, err error) {
	if !d.hasWrittenHeader {
		if _, err := d.separatorsW.Write([]byte(d.column)); err != nil {
			return 0, err
		}
		for i := int64(0); i < d.lineBytes; i++ {
			headerSB := &strings.Builder{}
			if _, err := headerSB.Write([]byte(mathx.PadFormatInt(i, d.addrBase, false, 2))); err != nil {
				return 0, err
			}
			if i < d.lineBytes-1 {
				if _, err := headerSB.Write([]byte(" ")); err != nil {
					return 0, err
				}
			}

			if _, err := d.columnW.Columns[2].Write([]byte(d.dumpHeaderFn(headerSB.String()))); err != nil {
				return 0, err
			}
		}
		if err := d.columnW.Flush(); err != nil {
			return 0, err
		}
		d.hasWrittenHeader = true
	}

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

			if _, err := d.columnW.Columns[2].Write([]byte(d.hexFn(b))); err != nil {
				return n, err
			}
			if d.offset%d.lineBytes != 0 {
				if _, err := d.columnW.Columns[2].Write([]byte(" ")); err != nil {
					return n, err
				}
			}
			if _, err := d.columnW.Columns[4].Write([]byte(d.asciiFn(b))); err != nil {
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
