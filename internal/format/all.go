package format

import (
	"fq/internal/decode"
	"fq/internal/format/aac"
	"fq/internal/format/elf"
	"fq/internal/format/flac"
	"fq/internal/format/flacpicture"
	"fq/internal/format/id3v1"
	"fq/internal/format/id3v11"
	"fq/internal/format/id3v2"
	"fq/internal/format/jpeg"
	"fq/internal/format/mp3"
	"fq/internal/format/mp3frame"
	"fq/internal/format/mp4"
	"fq/internal/format/ogg"
	"fq/internal/format/tar"
	"fq/internal/format/vorbis"
)

var All = []*decode.Register{
	flac.Register,
	mp3.Register,
	mp3frame.Register,
	id3v11.Register, // before id3v1 (TAG/TAG+ magic)
	id3v1.Register,
	id3v2.Register,
	elf.Register,
	ogg.Register,
	vorbis.Register,
	flacpicture.Register,
	jpeg.Register,
	tar.Register,
	mp4.Register,
	aac.Frame,
	aac.ADTS,
	aac.Stream,
}
