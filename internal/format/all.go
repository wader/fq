package format

import (
	"fq/internal/decode"
	"fq/internal/format/flac"
	"fq/internal/format/id3v1"
	"fq/internal/format/id3v2"
	"fq/internal/format/mp3"
	"fq/internal/format/mp3frame"
)

var All = []*decode.Register{
	flac.Register,
	mp3.Register,
	mp3frame.Register,
	id3v1.Register,
	id3v2.Register,
}
