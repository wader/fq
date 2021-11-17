package flac

import (
	"github.com/wader/fq/format"
	"github.com/wader/fq/format/registry"
	"github.com/wader/fq/pkg/decode"
)

var images decode.Group

var pictureTypeNames = decode.UToStr{
	0:  "Other",
	1:  "32x32 pixels 'file icon' (PNG only)",
	2:  "Other file icon",
	3:  "Cover (front)",
	4:  "Cover (back)",
	5:  "Leaflet page",
	6:  "Media (e.g. label side of CD)",
	7:  "Lead artist/lead performer/soloist",
	8:  "Artist/performer",
	9:  "Conductor",
	10: "Band/Orchestra",
	11: "Composer",
	12: "Lyricist/text writer",
	13: "Recording Location",
	14: "During recording",
	15: "During performance",
	16: "Movie/video screen capture",
	17: "A bright colored fish",
	18: "Illustration",
	19: "Band/artist logotype",
	20: "Publisher/Studio logotype",
}

func init() {
	registry.MustRegister(decode.Format{
		Name:        format.FLAC_PICTURE,
		Description: "FLAC metadatablock picture",
		DecodeFn:    pictureDecode,
		Dependencies: []decode.Dependency{
			{Names: []string{format.IMAGE}, Group: &images},
		},
	})
}

func pictureDecode(d *decode.D, in interface{}) interface{} {
	lenStr := func(name string) string { //nolint:unparam
		l := d.FieldU32(name + "_length")
		return d.FieldUTF8(name, int(l))
	}
	d.FieldU32("picture_type", d.MapUToStrSym(pictureTypeNames))
	lenStr("mime")
	lenStr("description")
	d.FieldU32("width")
	d.FieldU32("height")
	d.FieldU32("color_depth")
	d.FieldU32("number_of_index_colors")
	pictureLen := d.FieldU32("picture_length")
	if dv, _, _ := d.FieldTryFormatLen("picture_data", int64(pictureLen)*8, images, nil); dv == nil {
		d.FieldRawLen("picture_data", int64(pictureLen)*8)
	}

	return nil
}
