package tiff

// https://exiftool.org/TagNames/EXIF.html

import (
	"fq/format"
	"fq/format/registry"
	"fq/pkg/decode"
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
