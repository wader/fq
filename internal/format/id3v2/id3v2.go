package id3v2

import (
	"fmt"
	"fq/internal/decode"
)

var Register = &decode.Register{
	Name: "id3v2",
	MIME: "",
	New:  func(common decode.Common) decode.Decoder { return &Decoder{Common: common} },
}

// Decoder is ID3v2 decoder
type Decoder struct {
	decode.Common
}

func (d *Decoder) SyncSafeU32() uint64 {
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

// Decode ID3v2
func (d *Decoder) Decode(opts decode.Options) {
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
	size := d.FieldUFn("size,", func() (uint64, decode.Format, string) {
		return d.SyncSafeU32(), decode.FormatDecimal, ""
	})

	var extHeaderSize uint64
	if extendedHeader {
		switch version {
		case 3:
			extHeaderSize = d.FieldU32("size")
			d.FieldBytes("data", extHeaderSize)
		case 4:
			extHeaderSize = d.FieldUFn("size,", func() (uint64, decode.Format, string) {
				return d.SyncSafeU32(), decode.FormatDecimal, ""
			})
			// in v4 synchsafe integer includes itself
			d.FieldBytes("data", extHeaderSize-4)
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

	d.FieldBytes("tags", size-extHeaderSize)
}
