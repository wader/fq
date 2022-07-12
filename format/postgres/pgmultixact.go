package postgres

import (
	"github.com/wader/fq/format"
	"github.com/wader/fq/format/registry"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/scalar"
)

const BLCKSZ = 8192

func init() {
	registry.MustRegister(decode.Format{
		Name:        format.PGMULTIXACTOFF,
		Description: "PostgreSQL multixact offset file",
		DecodeFn:    mxOffsetDecode,
	})
	registry.MustRegister(decode.Format{
		Name:        format.PGMULTIXACTMEM,
		Description: "PostgreSQL multixact members file",
		DecodeFn:    mxMembersDecode,
	})
}

func mxOffsetDecode(d *decode.D, in interface{}) interface{} {
	d.Endian = decode.LittleEndian

	d.FieldArray("offsets", func(d *decode.D) {
		for {
			if d.End() {
				break
			}
			d.FieldU32("offset", scalar.Hex)

		}
	})
	return nil
}

var flags = scalar.UToScalar{
	0: {Sym: "ForKeyShare", Description: "For Key Share"},
	1: {Sym: "ForShare", Description: "For Share"},
	2: {Sym: "ForNoKeyUpdate", Description: "For No Key Update"},
	3: {Sym: "ForUpdate", Description: "For Update"},
	4: {Sym: "NoKeyUpdate", Description: "No Key Update"},
	5: {Sym: "Update", Description: "Update"},
}

func mxMembersDecode(d *decode.D, in interface{}) interface{} {
	var xidLen uint = 4
	var groupLen uint = 4 * (1 + xidLen)
	d.Endian = decode.LittleEndian

	m := d.FieldArrayValue("members")
	p := d.FieldArrayValue("paddings")

	for {
		var xacts []*decode.D = make([]*decode.D, 4)

		for i := 0; i < 4; i++ {
			xacts[i] = m.FieldStructValue("xact")
			xacts[i].FieldU8("status", flags)
		}

		for i := 0; i < 4; i++ {
			xacts[i].FieldU32("xid")
		}

		// Check if rest of bytes are padding before EOF
		if d.BitsLeft() < int64(groupLen*8) && d.BitsLeft() > 0 {
			p.FieldRawLen("padding", d.BitsLeft())
			break
		}

		// Check on EOF
		if d.End() {
			break
		}

		// Not EOF, let's check on block boundary
		if blkLeft := BLCKSZ - (uint(d.Pos())>>3)%BLCKSZ; blkLeft < groupLen {
			p.FieldRawLen("padding", int64(blkLeft*8))
		}
	}
	return nil
}
