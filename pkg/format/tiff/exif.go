package tiff

// https://exiftool.org/TagNames/EXIF.html

import (
	"fq/pkg/decode"
	"fq/pkg/format"
)

// currently just a alias for tiff

var exifIccProfile []*decode.Format

func init() {
	format.MustRegister(&decode.Format{
		Name:        format.EXIF,
		Description: "Exchangeable Image File Format",
		Groups:      []string{},
		DecodeFn:    tiffDecode,
		Dependencies: []decode.Dependency{
			{Names: []string{format.ICC_PROFILE}, Formats: &exifIccProfile},
		},
	})
}
