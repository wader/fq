package mpeg

// TODO: rename? eia something?
// EIA-708 cc_data
// https://shop.cta.tech/products/digital-television-dtv-closed-captioning

import (
	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
)

func init() {
	interp.RegisterFormat(
		format.MPEG_CC_Data,
		&decode.Format{
			Description: "EIA-708 cc_data closed captioning data",
			DecodeFn:    mpegCcDataDecode,
		})
}

func mpegCcDataDecode(d *decode.D) any {
	d.FieldU1("reserved0")
	d.FieldBool("process_cc_data_flag")
	d.FieldU1("zero_bit")
	ccCount := d.FieldU5("cc_count")
	d.FieldU8("reserved1")
	d.FieldArray("cc", func(d *decode.D) {
		for i := 0; i < int(ccCount); i++ {
			d.FieldStruct("cc", func(d *decode.D) {
				d.FieldU1("one_bit")
				d.FieldU4("reserved0")
				d.FieldBool("cc_valid")
				d.FieldU2("cc_type")
				d.FieldU8("cc_data_1")
				d.FieldU8("cc_data_2")
			})
		}
	})

	return nil
}
