package interp

import (
	"fq/internal/asciiwriter"
	"fq/internal/hexdump"
	"fq/internal/hexpairwriter"
	"fq/internal/num"
	"fq/pkg/bitio"
	"io"
)

func hexdumpRange(bbr bufferRange, w io.Writer, opts Options) error {
	bitsByteAlign := bbr.r.Start % 8
	bb, err := bbr.bb.BitBufRange(bbr.r.Start-bitsByteAlign, bbr.r.Len+bitsByteAlign)
	if err != nil {
		return err
	}

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
	defer hw.Close()
	if _, err = io.Copy(hw, bb); err != nil {
		return err
	}

	return nil
}
