package group

import (
	"fq/pkg/decode"
	"fq/pkg/format/jpeg"
	"fq/pkg/format/png"
	"fq/pkg/format/tiff"
)

var Images = []*decode.Format{
	jpeg.File,
	png.File,
	tiff.File,
}
