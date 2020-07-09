package group

import (
	"fq/internal/decode"
	"fq/internal/format/jpeg"
	"fq/internal/format/png"
	"fq/internal/format/tiff"
)

var Images = []*decode.Format{
	jpeg.File,
	png.File,
	tiff.File,
}
