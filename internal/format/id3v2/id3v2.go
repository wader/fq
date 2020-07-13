package id3v2

// https://id3.org/id3v2.3.0
// https://id3.org/id3v2.4.0-structure
// https://id3.org/id3v2.4.0-frames

import (
	"fmt"
	"fq/internal/decode"
)

var Tag = &decode.Format{
	Name: "id3v2",
	New:  func() decode.Decoder { return &TagDecoder{} },
}

// Decoder is ID3v2 tag decoder
type TagDecoder struct {
	decode.Common
}

func (d *TagDecoder) SyncSafeU32() uint64 {
	u := d.U32()
	// syncsafe integer is a number encoded
	// with 8th bit in each byte set to zero
	// 0aaaaaaa0bbbbbbb0ccccccc0ddddddd ->
	// 0000aaaaaaabbbbbbbcccccccddddddd
	return (((u & 0x7f000000) >> 3) |
		((u & 0x007f0000) >> 2) |
		((u & 0x00007f00) >> 1) |
		((u & 0x0000007f) >> 0))
}

func (d *TagDecoder) DecodeFrame(version int) uint64 {
	var id string
	switch version {
	case 2:
		id = d.UTF8(3)
		d.SeekRel(-3 * 8)
	case 3, 4:
		id = d.UTF8(4)
		d.SeekRel(-4 * 8)
	}

	d.FieldNoneFn(id, func() {

		switch version {
		case 2:
			// Frame ID   "XXX"
			// Frame size $xx xx xx
		case 3:
			// Frame ID   $xx xx xx xx  (four characters)
			// Size       $xx xx xx xx
			// Flags      $xx xx
		case 4:
			// Frame ID      $xx xx xx xx  (four characters)
			// Size      4 * %0xxxxxxx  (synchsafe integer)
			// Flags         $xx xx
		}
	})

	// TODO
	return 0
}

func (d *TagDecoder) DecodeFrames(version int, size uint64) {

	for size > 0 {
		// TODO: what was this about?
		var padding uint64
		for d.U8() == 0 {
			padding++
		}
		if padding > 0 {
			d.SeekRel(-int64(padding) * 8)
			d.FieldBytesLen("padding", padding)
		}

		size -= d.DecodeFrame(version)
	}

}

// Decode ID3v2
func (d *TagDecoder) Decode() {
	d.ValidateAtLeastBitsLeft(4 * 8)
	d.FieldValidateString("magic", "ID3")
	version := d.FieldU8("version")
	versionValid := version == 2 || version == 3 || version == 4
	if !versionValid {
		d.Invalid(fmt.Sprintf("unsupported version %d", version))
	}

	d.FieldU8("revision")
	var extendedHeader bool
	d.FieldNoneFn("flags", func() {
		d.FieldU1("unsynchronisation")
		extendedHeader = d.FieldBool("extended_header")
		d.FieldU1("experimental_indicator")
		d.FieldU5("unused")
	})
	size := d.FieldUFn("size,", func() (uint64, decode.NumberFormat, string) {
		return d.SyncSafeU32(), decode.NumberDecimal, ""
	})

	var extHeaderSize uint64
	if extendedHeader {
		switch version {
		case 3:
			extHeaderSize = d.FieldU32("size")
			d.FieldBytesLen("data", extHeaderSize)
		case 4:
			extHeaderSize = d.FieldUFn("size,", func() (uint64, decode.NumberFormat, string) {
				return d.SyncSafeU32(), decode.NumberDecimal, ""
			})
			// in v4 synchsafe integer includes itself
			d.FieldBytesLen("data", extHeaderSize-4)
		}
	}

	// TODO: unknownv version?

	// TODO: d.FieldBitBuf

	// d.FieldNoneFn("frames", func() {
	// 	sizeLeft := size

	// 	for sizeLeft > 0 {
	// 		// if d.PeekBits(8) == 0 {

	// 		// }

	// 	}

	// })

	d.FieldBytesLen("tags", size-extHeaderSize)
}
