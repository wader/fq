package id3v2

import (
	"fq/internal/decode"
)

// Decoder is id3v2 decoder
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

// Decode id3v2
func (d *Decoder) Decode(opts decode.Options) bool {
	d.FieldUTF8(3, "magic")
	version := d.FieldU8("version")
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
			d.FieldBytes(uint(extHeaderSize), "data")
		case 4:
			extHeaderSize = d.FieldUFn("size,", func() (uint64, decode.Format, string) {
				return d.SyncSafeU32(), decode.FormatDecimal, ""
			})
			// in v24 synchsafe integer includes itself
			d.FieldBytes(uint(extHeaderSize)-4, "data")
		}
	}

	d.FieldBytes(uint(size-extHeaderSize), "tags")

	return true
}
