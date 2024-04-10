package jpeg

import (
	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
	"github.com/wader/fq/pkg/scalar"
)

func init() {
	interp.RegisterFormat(
		format.JP2C,
		&decode.Format{
			Description: "JPEG 2000 codestream",
			Groups:      []*decode.Group{format.Probe, format.Image},
			DecodeFn:    jp2cDecode,
			RootName:    "segments",
			RootArray:   true,
		})
}

const (
	JP2_SOC = 0xff_4f /* start of codestream */
	JP2_SOT = 0xff_90 /* start of tile */
	JP2_SOD = 0xff_93 /* start of data */
	JP2_EOC = 0xff_d9 /* end of codestream */
	/* fixed information segment */
	JP2_SIZ = 0xff_51 /* image and tile size */
	/* functional segments */
	JP2_COD = 0xff_52 /* coding style default */
	JP2_COC = 0xff_53 /* coding style component */
	JP2_RGN = 0xff_5e /* region of interest */
	JP2_QCD = 0xff_5c /* quantization default */
	JP2_QCC = 0xff_5d /* quantization component */
	JP2_POC = 0xff_5f /* progression order change */
	/* pointer segments */
	JP2_TLM = 0xff_55 /* tile-part lengths */
	JP2_PLM = 0xff_57 /* packet length (main header) */
	JP2_PLT = 0xff_58 /* packet length (tile-part header) */
	JP2_PPM = 0xff_60 /* packed packet headers (main header) */
	JP2_PPT = 0xff_61 /* packet packet headers (tile-part header) */
	/* bitstream internal markers and segments */
	JP2_SOP = 0xff_91 /* start of packet */
	JP2_EPH = 0xff_92 /* end of packet header */
	/* informational segments */
	JP2_CRG = 0xff_63 /* component registration */
	JP2_COM = 0xff_64 /* comment */
)

var jp2Markers = scalar.UintMap{
	JP2_SOC: {Sym: "soc", Description: "Start of codestream"},
	JP2_SOT: {Sym: "sot", Description: "Start of tile"},
	JP2_SOD: {Sym: "sod", Description: "Start of data"},
	JP2_EOC: {Sym: "eoc", Description: "End of codestream"},
	JP2_SIZ: {Sym: "siz", Description: "Image and tile size"},
	JP2_COD: {Sym: "cod", Description: "Coding style default"},
	JP2_COC: {Sym: "coc", Description: "Coding style component"},
	JP2_RGN: {Sym: "rgn", Description: "Region of interest"},
	JP2_QCD: {Sym: "qcd", Description: "Quantization default"},
	JP2_QCC: {Sym: "qcc", Description: "Quantization component"},
	JP2_POC: {Sym: "poc", Description: "Progression order change"},
	JP2_TLM: {Sym: "tlm", Description: "Tile-part lengths"},
	JP2_PLM: {Sym: "plm", Description: "Packet length (main header)"},
	JP2_PLT: {Sym: "plt", Description: "Packet length (tile-part header)"},
	JP2_PPM: {Sym: "ppm", Description: "Packed packet headers (main header)"},
	JP2_PPT: {Sym: "ppt", Description: "Packet packet headers (tile-part header)"},
	JP2_SOP: {Sym: "sop", Description: "Start of packet"},
	JP2_EPH: {Sym: "eph", Description: "End of packet header"},
	JP2_CRG: {Sym: "crg", Description: "Component registration"},
	JP2_COM: {Sym: "com", Description: "Comment"},
}

func jp2cDecode(d *decode.D) any {
	if d.PeekUintBits(16) != JP2_SOC {
		d.Fatalf("no SOC marker")
	}

	seenSOC := false
	seenSIZ := false
	seenEOC := false

	for !seenEOC && !d.End() {
		d.FieldStruct("segment", func(d *decode.D) {
			marker := d.FieldU16("marker", jp2Markers, scalar.UintHex)
			switch marker {
			case JP2_SOC:
				// zero length
				seenSOC = true
				return
			case JP2_SOD:
				l, _ := d.PeekFind(16, 8, d.BitsLeft(), func(v uint64) bool {
					return v == JP2_SOT || v == JP2_EOC
				})
				d.FieldRawLen("data", l)
			case JP2_SOT:
				d.FieldU16("l_sot")
				d.FieldU16("i_sot")
				d.FieldU32("p_sot")
				d.FieldU8("tp_sot")
				d.FieldU8("tn_sot")
			case JP2_SIZ:
				seenSIZ = true
				d.FieldU16("l_siz")
				d.FieldU16("r_siz")
				d.FieldU32("x_siz")
				d.FieldU32("y_siz")
				d.FieldU32("xo_siz")
				d.FieldU32("yo_siz")
				d.FieldU32("xt_siz")
				d.FieldU32("yt_siz")
				d.FieldU32("xto_siz")
				d.FieldU32("yto_siz")
				cSiz := d.FieldU16("c_siz")
				d.FieldArray("components", func(d *decode.D) {
					for i := 0; i < int(cSiz); i++ {
						d.FieldStruct("component", func(d *decode.D) {
							d.FieldU8("s_sizi")
							d.FieldU8("xr_sizi")
							d.FieldU8("yr_sizi")
						})
					}
				})

			case JP2_COM:
				length := d.FieldU16("length")
				d.FieldU16("r_cme")
				d.FieldUTF8("data", int(length-4))
			case JP2_EOC:
				// zero length
				seenEOC = true
				return
			default:
				length := d.FieldU16("length")
				d.FieldRawLen("data", int64(length-2)*8)
			}
		})
	}

	if !(seenSOC && seenSIZ && seenEOC) {
		d.Fatalf("SOC, SIZ or EOC marker not found")
	}

	return nil
}
