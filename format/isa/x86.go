package isa

import (
	"strings"

	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
	"github.com/wader/fq/pkg/scalar"
	"golang.org/x/arch/x86/x86asm"
)

func init() {
	// amd64?
	interp.RegisterFormat(
		format.X86_64,
		&decode.Format{
			Description: "x86-64 instructions",
			DecodeFn:    func(d *decode.D) any { return decodeX86(d, 64) },
			RootArray:   true,
			RootName:    "instructions",
		})
	interp.RegisterFormat(
		format.X86_32,
		&decode.Format{
			Description: "x86-32 instructions",
			DecodeFn:    func(d *decode.D) any { return decodeX86(d, 32) },
			RootArray:   true,
			RootName:    "instructions",
		})
	interp.RegisterFormat(
		format.X86_16,
		&decode.Format{
			Description: "x86-16 instructions",
			DecodeFn:    func(d *decode.D) any { return decodeX86(d, 16) },
			RootArray:   true,
			RootName:    "instructions",
		})
}

func decodeX86(d *decode.D, mode int) any {
	var symLookup func(uint64) (string, uint64)
	var base int64
	var xi format.X86_64In

	if d.ArgAs(&xi) {
		symLookup = xi.SymLookup
		base = xi.Base
	}

	bb := d.BytesRange(0, int(d.BitsLeft()/8))
	// TODO: uint64?
	pc := base

	for !d.End() {
		d.FieldStruct("instruction", func(d *decode.D) {
			i, err := x86asm.Decode(bb, mode)
			if err != nil {
				d.Fatalf("failed to decode x86 instruction: %s", err)
			}

			d.FieldRawLen("opcode", int64(i.Len)*8, scalar.BitBufSym(x86asm.IntelSyntax(i, uint64(pc), symLookup)), scalar.RawHex)

			// log.Printf("i.Len: %#+v\n", i.Len)
			// log.Printf("i.Opcode: %x\n", i.Opcode)
			// log.Printf("i: %#+v\n", i)

			// TODO: rebuild op lower?
			d.FieldValueUint("op", uint64(i.Opcode), scalar.UintSym(strings.ToLower(i.Op.String())), scalar.UintHex)

			bb = bb[i.Len:]
			pc += int64(i.Len)
		})

	}

	return nil
}
