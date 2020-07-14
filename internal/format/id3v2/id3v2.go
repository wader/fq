package id3v2

// https://id3.org/id3v2.3.0
// https://id3.org/id3v2.4.0-structure
// https://id3.org/id3v2.4.0-frames

import (
	"bytes"
	"fmt"
	"fq/internal/decode"
	"fq/internal/format/group"
	"log"
)

var Tag = &decode.Format{
	Name: "id3v2",
	New:  func() decode.Decoder { return &TagDecoder{} },
}

type encoding int

const (
	encodingISO8859_1 = 0
	encodingUTF16     = 1
	encodingUTF16BE   = 2
	encodingUTF8      = 3
)

// $00 ISO-8859-1 [ISO-8859-1]. Terminated with $00.
// $01 UTF-16 [UTF-16] encoded Unicode [UNICODE] with BOM. All
//     strings in the same frame SHALL have the same byteorder.
//     Terminated with $00 00.
// $02 UTF-16BE [UTF-16] encoded Unicode [UNICODE] without BOM.
//     Terminated with $00 00.
// $03 UTF-8 [UTF-8] encoded Unicode [UNICODE]. Terminated with $00.
var encodingNames = map[uint64]string{
	encodingISO8859_1: "ISO-8859-1",
	encodingUTF16:     "UTF-16",
	encodingUTF16BE:   "UTF-16BE",
	encodingUTF8:      "UTF-8",
}

var encodingLen = map[uint64]int{
	encodingISO8859_1: 1,
	encodingUTF16:     2,
	encodingUTF16BE:   2,
	encodingUTF8:      1,
}

var encodingToUTF8 = map[int]func(b []byte) string{
	encodingISO8859_1: func(b []byte) string {
		rs := make([]rune, len(b))
		for i, r := range b {
			rs[i] = rune(r)
		}
		return string(rs)
	},
	encodingUTF16: func(b []byte) string {
		beBOM := []byte("\xfe\xff")
		leBOM := []byte("\xff\xfe")
		var rs []rune
		switch {
		case bytes.HasPrefix(b, leBOM):
			b = b[2:]
			rs = make([]rune, len(b)/2)
			for i := 0; i < len(b)/2; i++ {
				rs[i] = rune(b[i*2] | b[i*2+1]<<8)
			}
		case bytes.HasPrefix(b, beBOM):
			b = b[2:]
			fallthrough
		default:
			rs = make([]rune, len(b)/2)
			for i := 0; i < len(b)/2; i++ {
				rs[i] = rune(b[i*2]<<8 | b[i*2+1])
			}
		}
		return string(rs)
	},
	encodingUTF16BE: func(b []byte) string {
		rs := make([]rune, len(b)/2)
		for i := 0; i < len(b)/2; i++ {
			rs[i] = rune(b[i*2]<<8 + b[i*2+1])
		}
		return string(rs)
	},
	encodingUTF8: func(b []byte) string {
		return string(b)
	},
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

func (d *TagDecoder) FieldSyncSafeU32(name string) uint64 {
	return d.FieldUFn(name, func() (uint64, decode.NumberFormat, string) {
		return d.SyncSafeU32(), decode.NumberDecimal, ""
	})
}

func (d *TagDecoder) TextNull(encoding int, name string) string {
	return d.FieldStrFn(name, func() (string, string) {
		nullLen := encodingLen[encodingUTF8]
		if n, ok := encodingLen[uint64(encoding)]; ok {
			nullLen = n
		}
		encodingFn := encodingToUTF8[encodingUTF8]
		if fn, ok := encodingToUTF8[encoding]; ok {
			encodingFn = fn
		}

		textLen := d.PeekFind(uint64(nullLen*8), 0, -1)/8 - uint64(nullLen)
		log.Printf("textLen: %#+v\n", textLen)
		textBs := d.BytesLen(textLen)
		d.SeekRel(int64(nullLen) * 8)

		return encodingFn(textBs), ""
	})
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

	var size uint64
	var dataSize uint64

	d.FieldNoneFn(id, func() {
		switch version {
		case 2:
			// Frame ID   "XXX"
			// Frame size $xx xx xx
			d.FieldUTF8("id", 3)
			dataSize = d.FieldU24("size")
			size = dataSize + 6
		case 3:
			// Frame ID   $xx xx xx xx  (four characters)
			// Size       $xx xx xx xx
			// Flags      $xx xx
			d.FieldUTF8("id", 4)
			dataSize = d.FieldU32("size")
			d.FieldU16("flags")
			size = dataSize + 10
		case 4:
			// Frame ID      $xx xx xx xx  (four characters)
			// Size      4 * %0xxxxxxx  (synchsafe integer)
			// Flags         $xx xx
			d.FieldUTF8("id", 4)
			dataSize = d.FieldSyncSafeU32("size")
			var headerLen uint64 = 10

			const flagUnsync = 0b10
			const flagDataLen = 0b1

			dataLenFlag := false
			d.FieldNoneFn("flags", func() {
				d.FieldU14("unused")
				d.FieldBool("unsync")
				dataLenFlag = d.FieldBool("data_length_indicator")
			})

			if dataLenFlag {
				d.FieldSyncSafeU32("data_length_indicator")
				dataSize -= 4
				headerLen = 4
			}

			size = dataSize + headerLen
		}

		frames := map[string]func(size uint64){
			// <Header for 'Attached picture', ID: "APIC">
			// Text encoding      $xx
			// MIME type          <text string> $00
			// Picture type       $xx
			// Description        <text string according to encoding> $00 (00)
			// Picture data       <binary data>
			"APIC": func(size uint64) {
				d.SubLen(size*8, func() {
					encoding := int(d.FieldStringMapFn("text_encoding", encodingNames, "unknown", d.U8))
					d.TextNull(encodingUTF8, "mime_type")
					d.FieldU8("picture_type") // TODO: table
					d.TextNull(encoding, "description")
					_, errs := d.FieldDecodeLen("picture", d.BitsLeft(), group.Images...)
					for _, err := range errs {
						log.Printf("err: %#+v\n", err)
					}
				})
			},
		}

		if fn, ok := frames[id]; ok {
			fn(dataSize)
		} else {
			d.FieldBytesLen("data", dataSize)
		}
	})

	// TODO
	return size
}

func (d *TagDecoder) DecodeFrames(version int, size uint64) {
	for size > 0 {
		for d.PeekBits(8) == 0 {
			d.FieldValidateZeroPadding("padding", size*8)
			return
		}

		size -= d.DecodeFrame(version)
	}

	// TODO: padding?
}

// Decode ID3v2
func (d *TagDecoder) Decode() {
	d.ValidateAtLeastBitsLeft(4 * 8)
	d.FieldValidateString("magic", "ID3")
	version := int(d.FieldU8("version"))
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
	size := d.FieldUFn("size", func() (uint64, decode.NumberFormat, string) {
		return d.SyncSafeU32(), decode.NumberDecimal, ""
	})

	var extHeaderSize uint64
	if extendedHeader {
		d.FieldNoneFn("extended_header", func() {
			switch version {
			case 3:
				extHeaderSize = d.FieldU32("size")
				d.FieldBytesLen("data", extHeaderSize)
			case 4:
				extHeaderSize = d.FieldUFn("size", func() (uint64, decode.NumberFormat, string) {
					return d.SyncSafeU32(), decode.NumberDecimal, ""
				})
				// in v4 synchsafe integer includes itself
				d.FieldBytesLen("data", extHeaderSize-4)
			}
		})
	}

	d.DecodeFrames(version, size)
}
