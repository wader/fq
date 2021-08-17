package interp

import (
	"io"

	"github.com/wader/fq/internal/asciiwriter"
	"github.com/wader/fq/internal/hexdump"
	"github.com/wader/fq/internal/hexpairwriter"
	"github.com/wader/fq/internal/num"
	"github.com/wader/fq/pkg/bitio"
)

func hexdumpRange(bbr bufferRange, w io.Writer, opts Options) error {
	bitsByteAlign := bbr.r.Start % 8
	bb, err := bbr.bb.BitBufRange(bbr.r.Start-bitsByteAlign, bbr.r.Len+bitsByteAlign)
	if err != nil {
		return err
	}

	br := bitio.NewSectionBitReader(bb, 0, bb.Len())

	d := opts.Decorator
	hw := hexdump.New(
		w,
		(bbr.r.Start-bitsByteAlign)/8,
		num.DigitsInBase(bitio.BitsByteCount(bbr.r.Stop()+bitsByteAlign), true, opts.AddrBase),
		opts.AddrBase,
		opts.LineBytes,
		func(b byte) string { return d.ByteColor(b).Wrap(hexpairwriter.Pair(b)) },
		func(b byte) string { return d.ByteColor(b).Wrap(asciiwriter.SafeASCII(b)) },
		func(s string) string { return d.DumpHeader.Wrap(s) },
		func(s string) string { return d.DumpAddr.Wrap(s) },
		d.Column,
	)
	aw := &bitio.AlignBitWriter{W: hw, N: 8}
	// TODO: ugly, AlignBitWriter take Closer? some other way?
	defer hw.Close()
	defer aw.Close()
	if _, err = bitio.Copy(aw, br); err != nil {
		return err
	}

	return nil
}
