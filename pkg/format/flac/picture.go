package flac

import (
	"fq/pkg/decode"
	"fq/pkg/format/group"
)

var Picture = &decode.Format{
	Name:      "flac_picture",
	New:       func() decode.Decoder { return &PictureDecoder{} },
	SkipProbe: true,
}

// PictureDecoder is a FLAC picture decoder
type PictureDecoder struct {
	decode.Common
}

// PictureDecoder decodes a FLAC picture
func (d *PictureDecoder) Decode() {
	lenStr := func(name string) string {
		len := d.FieldU32(name + "_length")
		return d.FieldUTF8(name, len)
	}
	d.FieldU32("picture_type")
	lenStr("mime")
	lenStr("description")
	d.FieldU32("width")
	d.FieldU32("height")
	d.FieldU32("color_depth")
	d.FieldU32("number_of_index_colors")
	pictureLen := d.FieldU32("picture_length")
	d.FieldDecodeLen("picture_data", pictureLen*8, group.Images...)
}
