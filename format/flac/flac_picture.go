package flac

import (
	"github.com/wader/fq/format"
	"github.com/wader/fq/format/registry"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/scalar"
)

var images decode.Group

var pictureTypeNames = scalar.UToSymStr{
	0:  "Other",
	1:  "32x32_pixels",
	2:  "other_file_icon",
	3:  "cover_front)",
	4:  "cover_back",
	5:  "leaflet_page",
	6:  "media",
	7:  "lead_artist",
	8:  "artist",
	9:  "conductor",
	10: "band",
	11: "composer",
	12: "lyricist",
	13: "recording_location",
	14: "during_recording",
	15: "during_performance",
	16: "movie",
	17: "a_bright_colored_fish",
	18: "illustration",
	19: "artist_logotype",
	20: "publisher_logotype",
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

func pictureDecode(d *decode.D, in any) any {
	lenStr := func(name string) string {
		l := d.FieldU32(name + "_length")
		return d.FieldUTF8(name, int(l))
	}
	d.FieldU32("picture_type", pictureTypeNames)
	lenStr("mime")
	lenStr("description")
	d.FieldU32("width")
	d.FieldU32("height")
	d.FieldU32("color_depth")
	d.FieldU32("number_of_index_colors")
	pictureLen := d.FieldU32("picture_length")
	d.FieldFormatOrRawLen("picture_data", int64(pictureLen)*8, images, nil)

	return nil
}
