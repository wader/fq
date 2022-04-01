package isa

import (
	"strings"

	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
	"github.com/wader/fq/pkg/scalar"
	"golang.org/x/arch/arm64/arm64asm"
)

func init() {
	interp.RegisterFormat(
		format.Arm64,
		&decode.Format{
			Description: "ARM64 instructions",
			DecodeFn:    decodeARM64,
			RootArray:   true,
			RootName:    "instructions",
		})
}

func decodeARM64(d *decode.D) any {
	var symLookup func(uint64) (string, uint64)
	var base int64
	var ai format.ARM64In

	if d.ArgAs(&ai) {
		symLookup = ai.SymLookup
		base = ai.Base
	}

	bb := d.BytesRange(0, int(d.BitsLeft()/8))
	// TODO: uint64?
	pc := base

	for !d.End() {
		d.FieldStruct("instruction", func(d *decode.D) {
			i, err := arm64asm.Decode(bb)
			if err != nil {
				d.Fatalf("failed to decode arm64 instruction: %s", err)
			}

			// TODO: other syntax
			d.FieldRawLen("opcode", int64(4)*8, scalar.BitBufSym(arm64asm.GoSyntax(i, uint64(pc), symLookup, nil)), scalar.RawHex)

			// TODO: Enc?
			d.FieldValueUint("op", uint64(i.Enc), scalar.UintSym(strings.ToLower(i.Op.String())), scalar.UintHex)

			bb = bb[4:]
			pc += int64(4)
		})

	}

	return nil
}
