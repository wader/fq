package flac

import (
	"fq/internal/bitbuf"
	"fq/internal/decode"
)

var Picture = &decode.Register{
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
