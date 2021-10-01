package tiff

// https://exiftool.org/TagNames/EXIF.html
// TODO: JPEGInterchangeFormat/JPEGInterchangeFormatLength, seem to just after the exif tag usually?

import (
	"github.com/wader/fq/format"
	"github.com/wader/fq/format/registry"
	"github.com/wader/fq/pkg/decode"
)

// currently just a alias for tiff

func init() {
	registry.MustRegister(&decode.Format{
		Name:        format.EXIF,
		Description: "Exchangeable Image File Format",
		Groups:      []string{},
		DecodeFn:    tiffDecode,
	})
}
