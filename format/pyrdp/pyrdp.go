// Processes PyRDP replay files
// https://github.com/GoSecure/pyrdp
//
// Copyright (c) 2022-2023 GoSecure Inc.
// Copyright (c) 2024 Flare Systems
// Licensed under the MIT License
//
// Maintainer: Olivier Bilodeau <olivier.bilodeau@flare.io>
// Author: Lisandro Ubiedo
package pyrdp

import (
	"embed"
	"time"

	"github.com/wader/fq/format"
	"github.com/wader/fq/format/pyrdp/pdu"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
	"github.com/wader/fq/pkg/scalar"
)

//go:embed pyrdp.md
var pyrdpFS embed.FS

func init() {
	interp.RegisterFormat(
		format.PYRDP,
		&decode.Format{
			Description: "PyRDP Replay Files",
			DecodeFn:    decodePYRDP,
		})
	interp.RegisterFS(pyrdpFS)
}

func decodePYRDP(d *decode.D) any {
	d.Endian = decode.LittleEndian

	d.FieldArray("events", func(d *decode.D) {
		for !d.End() {
			d.FieldStruct("event", func(d *decode.D) {
				pos := d.Pos()

				size := d.FieldU64("size") // minus the length
				pduType := uint16(d.FieldU16("pdu_type", pdu.TypesMap))
				d.FieldU64("timestamp", scalar.UintActualUnixTimeDescription(time.Millisecond, time.RFC3339Nano))
				pduSize := int64(size - 18)

				pduParser, ok := pdu.ParsersMap[pduType]
				if !ok { // catch undeclared parsers
					if pduSize > 0 {
						d.FieldRawLen("data", pduSize*8)
					}
					return
				}
				parseFn, ok := pduParser.(func(d *decode.D, length int64))
				if !ok {
					return
				}
				parseFn(d, pduSize)

				curr := d.Pos() - pos
				d.FieldRawLen("extra", (int64(size)*8)-curr) // seek whatever is left
			})
		}
	})
	return nil
}
