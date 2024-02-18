package woff

// OpenType https://learn.microsoft.com/en-us/typography/opentype/

import (
	"bytes"
	"io"
	"time"

	"github.com/dsnet/compress/brotli"
	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/bitio"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
	"github.com/wader/fq/pkg/scalar"
)

func init() {
	interp.RegisterFormat(
		format.WOFF2,
		&decode.Format{
			Description: "Web Open Font Format version 2",
			Groups:      []*decode.Group{format.Probe},
			DecodeFn:    woff2Decode,
		})
}

var opentypeEpochDate = time.Date(1904, time.January, 1, 0, 0, 0, 0, time.UTC)

// WOFF2 1.3 UIntBase128 Data Type
func decodeUIntBase128(d *decode.D) uint64 {
	var accum uint32

	for i := 0; i < 5; i++ {
		dataByte := uint8(d.U8())

		if i == 0 && dataByte == 0x80 {
			d.Fatalf("no leading 0")
		}
		if accum&0xfe_00_00_00 != 0 {
			d.Fatalf("overflow")
		}

		accum = (accum << 7) | uint32(dataByte&0x7f)

		if dataByte&0x80 == 0 {
			return uint64(accum)
		}
	}

	d.Fatalf("exceeds 5 bytes")

	return 0
}

var knownTags = scalar.UintMapSymStr{
	0:  "cmap",
	1:  "head",
	2:  "hhea",
	3:  "hmtx",
	4:  "maxp",
	5:  "name",
	6:  "OS/2",
	7:  "post",
	8:  "cvt",
	9:  "fpgm",
	10: "glyf",
	11: "loca",
	12: "prep",
	13: "CFF",
	14: "VORG",
	15: "EBDT",
	16: "EBLC",
	17: "gasp",
	18: "hdmx",
	19: "kern",
	20: "LTSH",
	21: "PCLT",
	22: "VDMX",
	23: "vhea",
	24: "vmtx",
	25: "BASE",
	26: "GDEF",
	27: "GPOS",
	28: "GSUB",
	29: "EBSC",
	30: "JSTF",
	31: "MATH",
	32: "CBDT",
	33: "CBLC",
	34: "COLR",
	35: "CPAL",
	36: "SVG",
	37: "sbix",
	38: "acnt",
	39: "avar",
	40: "bdat",
	41: "bloc",
	42: "bsln",
	43: "cvar",
	44: "fdsc",
	45: "feat",
	46: "fmtx",
	47: "fvar",
	48: "gvar",
	49: "hsty",
	50: "just",
	51: "lcar",
	52: "mort",
	53: "morx",
	54: "opbd",
	55: "prop",
	56: "trak",
	57: "Zapf",
	58: "Silf",
	59: "Glat",
	60: "Gloc",
	61: "Feat",
	62: "Sill",
}

const tagGlyf = 10
const tagLoca = 11

const flavorTTCF = 0x74746366

func woff2Decode(d *decode.D) any {
	d.FieldUTF8("signature", 4, d.StrAssert("wOF2"))
	d.FieldU32("flavor", scalar.UintMapSymStr{
		flavorTTCF: "collection",
	})
	d.FieldU32("length")
	numTables := d.FieldU16("num_tables")
	d.FieldU16("reserved")
	d.FieldU32("total_sfnt_size")
	totalCompressSize := d.FieldU32("total_compressed_size")
	d.FieldU16("major_version")
	d.FieldU16("minor_version")
	d.FieldU32("meta_offset")
	d.FieldU32("meta_length")
	d.FieldU32("meta_orig_length")
	d.FieldU32("priv_offset")
	d.FieldU32("priv_length")

	type tableEntry struct {
		d                     *decode.D
		tag                   string
		transformationVersion uint64
		dataLen               int64
	}

	var tables []tableEntry

	d.FieldArray("tables", func(d *decode.D) {
		for i := uint64(0); i < numTables; i++ {
			d.FieldStruct("entry", func(d *decode.D) {
				transformationVersion := d.FieldU2("transformation_version")
				knownTag := d.FieldU6("known_tag", knownTags)
				var tag string
				if knownTag < 63 {
					tag = knownTags[knownTag]
				} else {
					tag = d.FieldUTF8("optional_tag", 4)
				}
				d.FieldValueStr("tag", tag)
				dataLen := d.FieldUintFn("orig_length", decodeUIntBase128)

				// For all tables in a font, except for 'glyf' and 'loca' tables, transformation version 0 indicates the null transform ...
				// For 'glyf' and 'loca' tables, transformation version 3 indicates the null transform ...
				glyfOrLoca := knownTag == tagGlyf || knownTag == tagLoca
				hasNullTransform :=
					(glyfOrLoca && transformationVersion == 0) ||
						(glyfOrLoca && transformationVersion == 3)

				if hasNullTransform {
					dataLen = d.FieldUintFn("transform_length", decodeUIntBase128)
				}

				tables = append(tables, tableEntry{
					d:                     d,
					tag:                   tag,
					transformationVersion: transformationVersion,
					dataLen:               int64(dataLen),
				})
			})
		}
	})

	// TODO: CollectionDirectory

	r := d.FieldRawLen("compressed", int64(totalCompressSize)*8)
	br, err := brotli.NewReader(bitio.NewIOReader(r), &brotli.ReaderConfig{})
	if err != nil {
		d.IOPanic(err, "brotli.NewReader")
	}
	brBuf := &bytes.Buffer{}
	_, err = io.Copy(brBuf, br)
	if err != nil {
		d.IOPanic(err, "brotli io.Copy")
	}

	left := brBuf.Bytes()
	for _, te := range tables {
		if len(left) < int(te.dataLen) {
			d.Fatalf("orig_len outside buffer")
		}

		data := bitio.NewBitReader(left[0:te.dataLen], -1)

		// TODO: move to own decoder?
		switch te.tag {
		case "name":
			// https://learn.microsoft.com/en-us/typography/opentype/spec/head
			te.d.FieldStructRootBitBufFn("data", data, func(d *decode.D) {
				version := d.FieldU16("version")
				count := d.FieldU16("count")
				storageOffset := d.FieldU16("storage_offset")

				d.FieldArray("records", func(d *decode.D) {
					for i := uint64(0); i < count; i++ {
						d.FieldStruct("record", func(d *decode.D) {
							d.FieldU16("platform_id")
							d.FieldU16("encoding_id")
							d.FieldU16("language_id")
							d.FieldU16("name_id")
							length := d.FieldU16("length")
							stringOffset := d.FieldU16("string_offset")
							d.RangeFn(int64(storageOffset+stringOffset)*8, int64(length)*8, func(d *decode.D) {
								d.FieldUTF16BE("value", int(length))
							})
						})
					}
				})

				// TODO: tags?
				_ = version
			})
		case "head":
			// https://learn.microsoft.com/en-us/typography/opentype/spec/head
			te.d.FieldStructRootBitBufFn("data", data, func(d *decode.D) {
				d.FieldU32("version")
				d.FieldU32("font_revision")
				d.FieldU32("checksum_adjustment", scalar.UintHex)
				d.FieldU32("magic_number", scalar.UintHex)
				d.FieldS16("flags")
				d.FieldS16("units_per_em")
				d.FieldU64("created", scalar.UintActualDateDescription(opentypeEpochDate, time.Second, time.RFC3339))
				d.FieldU64("modified", scalar.UintActualDateDescription(opentypeEpochDate, time.Second, time.RFC3339))
				d.FieldS16("x_min")
				d.FieldS16("y_min")
				d.FieldS16("x_max")
				d.FieldS16("y_max")
				d.FieldS16("mac_style")
				d.FieldS16("lowest_rec_ppem")
				d.FieldS16("font_direction_hint")
				d.FieldS16("index_to_loc_format")
				// d.FieldS16("glyph_data_format")
			})
		case "hhea":
			// https://learn.microsoft.com/en-us/typography/opentype/spec/hhea
			te.d.FieldStructRootBitBufFn("data", data, func(d *decode.D) {
				d.FieldU32("version")
				d.FieldS16("ascent")                  //	Distance from baseline of highest ascender
				d.FieldS16("descent")                 //	Distance from baseline of lowest descender
				d.FieldS16("line_gap")                //	typographic line gap
				d.FieldU16("advance_width_max")       //	must be consistent with horizontal metrics
				d.FieldS16("min_left_side_bearing")   //	must be consistent with horizontal metrics
				d.FieldS16("min_right_side_bearing")  //	must be consistent with horizontal metrics
				d.FieldS16("x_max_extent")            //	max(lsb + (xMax-xMin))
				d.FieldS16("caret_slope_rise")        //	used to calculate the slope of the caret (rise/run) set to 1 for vertical caret
				d.FieldS16("caret_slope_run")         //	0 for vertical
				d.FieldS16("caret_offset")            //	set value to 0 for non-slanted fonts
				d.FieldS16("reserved0")               //	set value to 0
				d.FieldS16("reserved1")               //	set value to 0
				d.FieldS16("reserved2")               //	set value to 0
				d.FieldS16("reserved3")               //	set value to 0
				d.FieldS16("metric_data_format")      //	0 for current format
				d.FieldU16("num_of_long_hor_metrics") // number of advance widths in metrics table
			})
		default:
			te.d.FieldRootBitBuf("data", data)
		}

		left = left[te.dataLen:]
	}

	return nil
}
