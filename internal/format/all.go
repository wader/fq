package format

import (
	"fq/internal/decode"
	"fq/internal/format/aac"
	"fq/internal/format/elf"
	"fq/internal/format/flac"
	"fq/internal/format/id3v1"
	"fq/internal/format/id3v11"
	"fq/internal/format/id3v2"
	"fq/internal/format/jpeg"
	"fq/internal/format/mp3"
	"fq/internal/format/mp4"
	"fq/internal/format/ogg"
	"fq/internal/format/png"
	"fq/internal/format/tar"
	"fq/internal/format/tiff"
	"fq/internal/format/vorbis"
)

var All = []*decode.Format{
	flac.File,
	flac.Picture,
	mp3.File,
	mp3.Frame,
	id3v11.Tag, // before id3v1 (TAG/TAG+ magic)
	id3v1.Tag,
	id3v2.Tag,
	elf.File,
	ogg.File,
	vorbis.Packet,
	jpeg.File,
	tar.File,
	mp4.File,
	aac.Frame,
	aac.ADTS,
	aac.Stream,
	png.File,
	tiff.File,
}
