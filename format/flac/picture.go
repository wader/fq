package flac

import (
	"fq/format"
	"fq/format/registry"
	"fq/pkg/decode"
)

var images []*decode.Format

func init() {
	registry.MustRegister(&decode.Format{
		Name:        format.FLAC_PICTURE,
		Description: "FLAC metadatablock picture",
		DecodeFn:    pictureDecode,
		Dependencies: []decode.Dependency{
			{Names: []string{format.IMAGE}, Formats: &images},
		},
	})
}

func pictureDecode(d *decode.D, in interface{}) interface{} {
	lenStr := func(name string) string { //nolint:unparam
		len := d.FieldU32(name + "_length")
		return d.FieldUTF8(name, int(len))
	}
	d.FieldU32("picture_type")
	lenStr("mime")
	lenStr("description")
	d.FieldU32("width")
	d.FieldU32("height")
	d.FieldU32("color_depth")
	d.FieldU32("number_of_index_colors")
	pictureLen := d.FieldU32("picture_length")
	d.FieldDecodeLen("picture_data", int64(pictureLen)*8, images)

	return nil
}
