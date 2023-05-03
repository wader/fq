package mpeg

// http://www.mpucoder.com/DVD/spu.html
// http://sam.zoy.org/writings/dvd/subtitles/
// TODO: still some unknown data before and after pixel data

import (
	"fmt"
	"strings"

	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
	"github.com/wader/fq/pkg/scalar"
)

func init() {
	interp.RegisterFormat(
		format.MPEG_SPU,
		&decode.Format{
			Description: "Sub Picture Unit (DVD subtitle)",
			DecodeFn:    spuDecode,
		})
}

const (
	CMD_END    = 0xff
	FSTA_DSP   = 0x00
	STA_DSP    = 0x01
	STP_DSP    = 0x02
	SET_COLOR  = 0x03
	SET_CONTR  = 0x04
	SET_DAREA  = 0x05
	SET_DSPXA  = 0x06
	CHG_COLCON = 0x07
)

var commandNames = scalar.UintMapSymStr{
	CMD_END:    "CMD_END",
	FSTA_DSP:   "FSTA_DSP",
	STA_DSP:    "STA_DSP",
	STP_DSP:    "STP_DSP",
	SET_COLOR:  "SET_COLOR",
	SET_CONTR:  "SET_CONTR",
	SET_DAREA:  "SET_DAREA",
	SET_DSPXA:  "SET_DSPXA",
	CHG_COLCON: "CHG_COLCON",
}

func rleValue(d *decode.D) (uint64, uint64, int) {
	p := uint(d.PeekUintBits(8))

	// match zero prefix
	switch {
	case p&0b1111_1100 == 0:
		// 000000nnnnnnnncc
		d.U6()
		return d.U8(), d.U2(), 16
	case p&0b1111_0000 == 0:
		// 0000nnnnnncc
		d.U4()
		return d.U6(), d.U2(), 12
	case p&0b1100_0000 == 0:
		// 00nnnncc
		d.U2()
		return d.U4(), d.U2(), 8
	default:
		// nncc
		return d.U2(), d.U2(), 4
	}
}

//nolint:unparam
func decodeLines(d *decode.D, lines int, width int) []string {
	var ls []string

	for i := 0; i < lines; i++ {
		l := ""
		for x := 0; x < width; {
			n, c, b := rleValue(d)
			pixel := " "
			if c != 0 {
				pixel = fmt.Sprintf("%d", c)
			}

			if n == 0 && b == 16 {
				l += strings.Repeat(pixel, width-len(l))
				break
			}

			x += int(n)
		}
		if d.ByteAlignBits() > 0 {
			d.U(d.ByteAlignBits())
		}

		ls = append(ls, l)
	}

	return ls
}

func spuDecode(d *decode.D) any {
	d.FieldU16("size")
	dcsqtOffset := d.FieldU16("dcsqt_offset")

	// to catch infinite loops
	offsetSeen := map[uint64]struct{}{}

	d.SeekAbs(int64(dcsqtOffset) * 8)
	d.FieldArray("dcsqt", func(d *decode.D) {
		lastDCSQ := false

		for !lastDCSQ {
			d.FieldStruct("dcsq", func(d *decode.D) {
				dcsqStart := uint64(d.Pos() / 8)
				d.FieldU16("delay")
				offset := d.FieldU16("offset")
				if offset == dcsqStart {
					lastDCSQ = true
					return
				}

				if _, ok := offsetSeen[offset]; ok {
					d.Fatalf("dcsqt loop detected for %d", offset)
				}
				offsetSeen[offset] = struct{}{}

				var pxdTFOffset int64
				var pxdBFOffset int64
				var width uint64
				var height uint64

				d.FieldArray("commands", func(d *decode.D) {
					seenEnd := false
					for !seenEnd {
						d.FieldStruct("command", func(d *decode.D) {
							cmd := d.FieldU8("type", commandNames)
							switch cmd {
							case CMD_END:
								seenEnd = true
							case FSTA_DSP:
								// no args
							case STA_DSP:
								// no args
							case STP_DSP:
								// no args
							case SET_COLOR:
								d.FieldU4("a0")
								d.FieldU4("a1")
								d.FieldU4("a2")
								d.FieldU4("a3")
							case SET_CONTR:
								d.FieldU4("a0")
								d.FieldU4("a1")
								d.FieldU4("a2")
								d.FieldU4("a3")
							case SET_DAREA:
								startX := d.FieldU12("start_x")
								endX := d.FieldU12("end_x")
								startY := d.FieldU12("start_y")
								endY := d.FieldU12("end_y")
								width = endX - startX + 1
								height = endY - startY + 1
							case SET_DSPXA:
								pxdTFOffset = int64(d.FieldU16("offset_top_field")) * 8
								pxdBFOffset = int64(d.FieldU16("offset_bottom_field")) * 8
							case CHG_COLCON:
								size := d.FieldU16("size")
								// TODO
								d.FieldRawLen("data", int64(size)*8)
							default:
								d.Fatalf("unknown command %d", cmd)
							}
						})
					}
				})

				halfHeight := int(height) / 2
				// var tLines []string
				// var bLines []string

				if pxdTFOffset != 0 {
					d.SeekAbs(pxdTFOffset)
					/*tLines*/ _ = decodeLines(d, halfHeight, int(width))
					d.RangeFn(pxdTFOffset, d.Pos()-pxdTFOffset, func(d *decode.D) {
						d.FieldRawLen("top_pixels", d.BitsLeft())
					})
				}
				if pxdBFOffset != 0 {
					d.SeekAbs(pxdBFOffset)
					/*bLines*/ _ = decodeLines(d, halfHeight, int(width))
					d.RangeFn(pxdBFOffset, d.Pos()-pxdBFOffset, func(d *decode.D) {
						d.FieldRawLen("bottom_pixels", d.BitsLeft())
					})

				}

				// var lines []string
				// for i := 0; i < halfHeight; i++ {
				// 	lines = append(lines, tLines[i])
				// 	lines = append(lines, bLines[i])
				// }

				// for _, l := range lines {
				// 	log.Printf("l: '%s'\n", l)
				// }

				d.SeekAbs(int64(offset) * 8)
			})
		}
	})

	return nil
}
