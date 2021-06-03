package tiff

// https://exiftool.org/TagNames/EXIF.html

import (
	"fq/pkg/decode"
	"fq/pkg/format"
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
