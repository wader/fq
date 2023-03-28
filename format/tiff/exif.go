package tiff

// https://exiftool.org/TagNames/EXIF.html
// TODO: JPEGInterchangeFormat/JPEGInterchangeFormatLength, seem to just after the exif tag usually?

import (
	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
)

// currently just a alias for tiff

func init() {
	interp.RegisterFormat(
		format.Exif,
		&decode.Format{
			Description: "Exchangeable Image File Format",
			Groups:      []*decode.Group{},
			DecodeFn:    tiffDecode,
		})
}
