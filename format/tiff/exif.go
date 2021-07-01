package tiff

// https://exiftool.org/TagNames/EXIF.html

import (
	"fq/format"
	"fq/pkg/decode"
)

// currently just a alias for tiff

func init() {
	format.MustRegister(&decode.Format{
		Name:        format.EXIF,
		Description: "Exchangeable Image File Format",
		Groups:      []string{},
		DecodeFn:    tiffDecode,
	})
}
