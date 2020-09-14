package format

import (
	"fq/pkg/decode"
	"fq/pkg/format/aac"
	"fq/pkg/format/ape"
	"fq/pkg/format/elf"
	"fq/pkg/format/flac"
	"fq/pkg/format/icc"
	"fq/pkg/format/id3v1"
	"fq/pkg/format/id3v11"
	"fq/pkg/format/id3v2"
	"fq/pkg/format/jpeg"
	"fq/pkg/format/mp3"
	"fq/pkg/format/mp4"
	"fq/pkg/format/ogg"
	"fq/pkg/format/png"
	"fq/pkg/format/tar"
	"fq/pkg/format/tiff"
	"fq/pkg/format/vorbis"
	"fq/pkg/format/wav"
)

// All formats
var All = []*decode.Format{
	flac.File,
	flac.Picture,
	mp3.File,
	mp3.Frame,
	mp3.XingHeader,
	id3v11.Tag,
	id3v1.Tag,
	id3v2.Tag,
	elf.File,
	ogg.File,
	ogg.Page,
	vorbis.Packet,
	jpeg.File,
	tar.File,
	mp4.File,
	aac.Frame,
	aac.ADTS,
	aac.Stream,
	png.File,
	tiff.File,
	ape.TagV2,
	wav.File,
	icc.Tag,
}
