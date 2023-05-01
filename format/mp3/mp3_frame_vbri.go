package mp3

// https://www.codeproject.com/Articles/8295/MPEG-Audio-Frame-Header

import (
	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
	"github.com/wader/fq/pkg/scalar"
)

func init() {
	interp.RegisterFormat(
		format.MP3_Frame_VBRI,
		&decode.Format{
			Description: "MP3 frame Fraunhofer encoder variable bitrate tag",
			Groups:      []*decode.Group{format.MP3_Frame_Tags},
			DecodeFn:    mp3FrameTagVBRIDecode,
		})
}

func mp3FrameTagVBRIDecode(d *decode.D) any {
	d.FieldUTF8("header", 4, d.StrAssert("VBRI"))
	d.FieldU16("version_id")
	d.FieldU16("delay")
	d.FieldU16("quality")
	d.FieldU32("length", scalar.UintDescription("Number of bytes"))
	d.FieldU32("frames", scalar.UintDescription("Number of frames"))
	tocEntries := d.FieldU16("toc_entries", scalar.UintDescription("Number of entries within TOC table"))
	d.FieldU16("scale_factor", scalar.UintDescription("Scale factor of TOC table entries"))
	tocEntrySize := d.FieldU16("toc_entry_size", d.UintAssert(1, 2, 3, 4), scalar.UintDescription("Size per table entry"))
	d.FieldU16("frame_per_entry", scalar.UintDescription("Frames per table entry"))
	d.FieldArray("toc", func(d *decode.D) {
		for i := 0; i < int(tocEntries); i++ {
			d.FieldU("entry", int(tocEntrySize)*8)
		}
	})

	return nil
}
