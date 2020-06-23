package flacpicture

// https://xiph.org/ogg/doc/framing.html

import (
	"fq/internal/bitbuf"
	"fq/internal/decode"
)

var Register = &decode.Register{
	Name: "flacpicture",
	MIME: "",
	New: func(common decode.Common) decode.Decoder {
		return &Decoder{
			Common: common,
		}
	},
	SkipProbe: true,
}

// Decoder is a flacpicture decoder
type Decoder struct {
	decode.Common
}

// Decode flacpicture
func (d *Decoder) Decode(opts decode.Options) {
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
	pictureBs := d.FieldBytesLen("picture_data", uint64(pictureLen))
	d.FieldDecodeBitBuf("picture", bitbuf.NewFromBytes(pictureBs), []string{"jpeg"})

	// uint32 "The picture type"
	// set mime_length [uint32 "MIME length"]
	// ascii_maybe_empty $mime_length "MIME type"
	// set desc_len [uint32 "Description length"]
	// ascii_maybe_empty $desc_len "Description"
	// uint32 "Width"
	// uint32 "Height"
	// uint32 "Color depth"
	// uint32 "Number of indexed colors"
	// set picture_len [uint32 "Picture length"]
	// bytes $picture_len "Picture data"
}
