package caff

import (
	"bytes"
	"compress/flate"
	"io"

	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/bitio"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
	"github.com/wader/fq/pkg/scalar"
)

var probeGroup decode.Group

func init() {
	interp.RegisterFormat(
		format.CAFF,
		&decode.Format{
			Description: "Live2D Cubism archive",
			Groups:      []*decode.Group{format.Probe},
			DecodeFn:    decodeCAFF,
			DefaultInArg: format.CAFF_In{
				Uncompress: true,
			},
			Dependencies: []decode.Dependency{
				{Groups: []*decode.Group{format.Probe}, Out: &probeGroup},
			},
		})
}

const (
	imageFormatUnknown   = 0x00
	imageFormatPNG       = 0x01
	imageFormatNoPreview = 0x7F
)

const (
	colorTypeUnknown   = 0x00
	colorTypeARGB      = 0x01
	colorTypeRGB       = 0x02
	colorTypeNoPreview = 0x7F
)

const (
	compressOptionRaw   = 0x10
	compressOptionFast  = 0x21
	compressOptionSmall = 0x25
)

var imageFormatNames = scalar.UintMapSymStr{
	imageFormatUnknown:   "unknown",
	imageFormatPNG:       "png",
	imageFormatNoPreview: "no_preview",
}

var colorTypeNames = scalar.UintMapSymStr{
	colorTypeUnknown:   "unknown",
	colorTypeARGB:      "argb",
	colorTypeRGB:       "rgb",
	colorTypeNoPreview: "no_preview",
}

var compressOptionNames = scalar.UintMapSymStr{
	compressOptionRaw:   "raw",
	compressOptionFast:  "fast",
	compressOptionSmall: "small",
}

type fileInfoListEntry struct {
	filePath       string
	startPos       int64
	fileSize       int
	isObfuscated   bool
	compressOption uint8
}

func decodeVersion(d *decode.D) {
	d.FieldU8("major")
	d.FieldU8("minor")
	d.FieldU8("patch")
}

func decodeCAFF(d *decode.D) any {
	var ci format.CAFF_In
	d.ArgAs(&ci)

	var obfsKey uint64

	obfsU8 := func(d *decode.D) uint64 { return d.U8() ^ (obfsKey & 0xFF) }
	obfsU32 := func(d *decode.D) uint64 { return d.U32() ^ (obfsKey & 0xFFFFFFFF) }
	obfsU64 := func(d *decode.D) uint64 { return d.U64() ^ (obfsKey<<32 | obfsKey) }
	obfsBool := func(d *decode.D) bool { return obfsU8(d) != 0 }

	// "Big Endian Base 128" - LEB128's strange sibling
	obfsBEB128 := func(d *decode.D) (v uint64) {
		for {
			x := obfsU8(d)
			v <<= 7
			v |= (x & 0x7F)
			if (x >> 7) == 0 {
				return
			}
		}
	}

	obfsVarStr := func(d *decode.D) string {
		length := obfsBEB128(d)
		if length == 0 {
			return ""
		}

		raw := d.BytesLen(int(length))
		for i := uint64(0); i < length; i++ {
			raw[i] ^= byte(obfsKey)
		}
		return string(raw)
	}

	d.FieldUTF8("archive_id", 4, d.StrAssert("CAFF"))
	d.FieldStruct("archive_version", decodeVersion)
	d.FieldUTF8("format_id", 4)
	d.FieldStruct("format_version", decodeVersion)
	obfsKey = d.FieldU32("obfuscate_key")
	d.SeekRel(8 * 8)

	d.FieldStruct("preview_image", func(d *decode.D) {
		d.FieldU8("image_format", imageFormatNames)
		d.FieldU8("color_type", colorTypeNames)
		d.SeekRel(2 * 8)
		d.FieldU16("width")
		d.FieldU16("height")
		d.FieldU64("start_pos")
		d.FieldU32("file_size")
	})
	d.SeekRel(8 * 8)

	fileInfoListSize := d.FieldUintFn("file_info_map_size", obfsU32)
	fileInfoList := make([]fileInfoListEntry, int(fileInfoListSize))

	d.FieldArray("file_info_list", func(d *decode.D) {
		for i := uint64(0); i < fileInfoListSize; i++ {
			d.FieldStruct("file_info", func(d *decode.D) {
				var entry fileInfoListEntry

				entry.filePath = d.FieldStrFn("file_path", obfsVarStr)
				d.FieldStrFn("tag", obfsVarStr)
				entry.startPos = int64(d.FieldUintFn("start_pos", obfsU64))
				entry.fileSize = int(d.FieldUintFn("file_size", obfsU32))
				entry.isObfuscated = d.FieldBoolFn("is_obfuscated", obfsBool)
				entry.compressOption = uint8(d.FieldUintFn("compress_option", obfsU8, compressOptionNames))
				d.SeekRel(8 * 8)

				fileInfoList[int(i)] = entry
			})
		}
	})

	d.FieldArray("files", func(d *decode.D) {
		for _, entry := range fileInfoList {
			d.FieldStruct("file", func(d *decode.D) {
				d.SeekAbs(entry.startPos * 8)
				d.FieldValueStr("file_path", entry.filePath)
				d.FieldValueUint("file_size", uint64(entry.fileSize))
				d.FieldValueBool("is_obfuscated", entry.isObfuscated)
				d.FieldValueUint("compress_option", uint64(entry.compressOption), compressOptionNames)

				rawBytes := d.BytesLen(entry.fileSize)
				if entry.isObfuscated {
					for i, v := range rawBytes {
						rawBytes[i] = v ^ uint8(obfsKey)
					}
				}

				br := bitio.NewBitReader(rawBytes, -1)
				if !ci.Uncompress || entry.compressOption == compressOptionRaw {
					fieldName := "uncompressed"
					if entry.compressOption != compressOptionRaw {
						fieldName = "compressed"
					}

					value, _, err := d.TryFieldFormatBitBuf(fieldName, br, &probeGroup, format.Probe_In{})
					if value == nil && err != nil {
						d.FieldRootBitBuf(fieldName, br)
					}
				} else {
					d.FieldRootBitBuf("compressed", br)

					// Offset 0x26: skip ZIP entry header; there's nothing useful in it and it's always the same
					infBytes, err := io.ReadAll(flate.NewReader(bytes.NewReader(rawBytes[0x26:])))
					if err == nil {
						infBr := bitio.NewBitReader(infBytes, -1)
						_, _, _ = d.TryFieldFormatBitBuf("uncompressed", infBr, &probeGroup, format.Probe_In{})
					}
				}
			})
		}
	})

	d.SeekAbs(d.Len() - 2*8)
	d.FieldRawLen("guard_bytes", 2*8, d.AssertBitBuf([]byte{98, 99}))

	return nil
}
