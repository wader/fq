package prores

// https://wiki.multimedia.cx/index.php/Apple_ProRes

import (
	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
	"github.com/wader/fq/pkg/scalar"
)

func init() {
	interp.RegisterFormat(
		format.Prores_Frame,
		&decode.Format{
			Description: "Apple ProRes frame",
			DecodeFn:    decodeProResFrame,
		})
}

func decodeProResFrame(d *decode.D) any {
	var size int64
	d.FieldStruct("container", func(d *decode.D) {
		size = int64(d.FieldU32("size"))
		d.FieldUTF8("type", 4, d.StrAssert("icpf"))
	})
	d.FramedFn((size-8)*8, func(d *decode.D) {
		d.FieldStruct("header", func(d *decode.D) {
			d.FieldU16("hdr_size")
			d.FieldU16("version")
			d.FieldUTF8("creator_id", 4)
			d.FieldU16("width")
			d.FieldU16("height")
			d.FieldStruct("frame_flags", func(d *decode.D) {
				d.FieldU2("chrominance_factor", scalar.UintMapSymStr{
					2: "422",
					3: "444",
				})
				d.FieldU2("unused0")
				d.FieldU2("frame_type", scalar.UintMapSymStr{
					0: "progressive",
					1: "interlaced_top_first",
					2: "interlaced_bottom_first",
				})
				d.FieldU2("unused1")
			})
			// TODO: more mappings
			d.FieldU8("reserved1")
			d.FieldU8("primaries")
			d.FieldU8("transf_func")
			d.FieldU8("color_matrix")
			d.FieldU4("src_pix_fmt")
			d.FieldU4("alpha_info")
			d.FieldU8("reserved2")
			d.FieldU8("q_mat_flags")
			d.FieldRawLen("q_mat_luma", 64*8)
			d.FieldRawLen("q_mat_chroma", 64*8)
		})
		d.FieldRawLen("picture_data", d.BitsLeft())
	})
	return nil
}
