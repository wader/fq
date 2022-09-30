package riff

import (
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/scalar"
)

type pathEntry struct {
	id   string
	data any
}

type path []pathEntry

func (p path) topData() any {
	if len(p) < 1 {
		return nil
	}
	return p[len(p)-1].data
}

func riffDecode(d *decode.D, path path, headFn func(d *decode.D, path path) (string, int64), chunkFn func(d *decode.D, id string, path path) (bool, any)) {
	id, size := headFn(d, path)

	d.FramedFn(size*8, func(d *decode.D) {
		hasChildren, data := chunkFn(d, id, path)
		if hasChildren {
			np := append(path, pathEntry{id: id, data: data})
			d.FieldArray("chunks", func(d *decode.D) {
				for !d.End() {
					d.FieldStruct("chunk", func(d *decode.D) {
						riffDecode(d, np, headFn, chunkFn)
					})
				}
			})
		}
	})

	wordAlgin := d.AlignBits(16)
	if wordAlgin != 0 {
		d.FieldRawLen("align", int64(wordAlgin))
	}
}

// TODO: sym name?
var chunkIDDescriptions = scalar.StrMapDescription{
	"LIST": "Chunk list",
	"JUNK": "Alignment",

	"idx1": "Index",
	"indx": "Base index",

	"avih": "AVI main header",
	"strh": "Stream header",
	"strf": "Stream format",
	"strn": "Stream name",
	"vprp": "Video properties",

	"dmlh": "Extended AVI header",

	"ISMP": "SMPTE timecode",
	"IDIT": "Time and date digitizing commenced",
	"IARL": "Archival Location. Indicates where the subject of the file is archived.",
	"IART": "Artist. Lists the artist of the original subject of the file",
	"ICMS": "Commissioned. Lists the name of the person or organization that commissioned the subject of the file",
	"ICMT": "Comments. Provides general comments about the file or the subject of the file",
	"ICOP": "Copyright. Records the copyright information for the file",
	"ICRD": "Creation date. Specifies the date the subject of the file was created.",
	"ICRP": "Cropped. Describes whether an image has been cropped and, if so, how it was cropped",
	"IDIM": "Dimensions. Specifies the size of the original subject of the file",
	"IDPI": "Dots Per Inch. Stores dots per inch setting of the digitizer used to produce the file",
	"IENG": "Engineer. Stores the name of the engineer who worked on the file. If there are multiple engineers, separate the names by a semicolon and a blank",
	"IGNR": "Genre. Describes the original work",
	"IKEY": "Keywords. Provides a list of keywords that refer to the file or subject of the file",
	"ILGT": "Lightness. Describes the changes in lightness settings on the digitizer required to produce the file.",
	"IMED": "Medium. Describes the original subject of the file",
	"INAM": "Name. Stores the title of the subject of the file",
	"IPLT": "Palette Setting. Specifies the number of colors requested when digitizing an image",
	"IPRD": "Product. Specifies the name of the title the file was originally intended for",
	"ISBJ": "Subject. Describes the contents of the file",
	"ISFT": "Software. Identifies the name of the software package used to create the file",
	"ISHP": "Sharpness. Identifies the changes in sharpness for the digitizer required to produce the file",
	"ISRC": "Source. Identifies the name of the person or organization who supplied the original subject of the file",
	"ISRF": "Source Form. Identifies the original form of the material that was digitized",
	"ITCH": "Technician. Identifies the technician who digitized the subject file",
}

func riffIsStringChunkID(id string) bool {
	switch id {
	case "strn",
		"ISMP",
		"IDIT",
		"IARL",
		"IART",
		"ICMS",
		"ICMT",
		"ICOP",
		"ICRD",
		"ICRP",
		"IDIM",
		"IDPI",
		"IENG",
		"IGNR",
		"IKEY",
		"ILGT",
		"IMED",
		"INAM",
		"IPLT",
		"IPRD",
		"ISBJ",
		"ISFT",
		"ISHP",
		"ISRC",
		"ISRF",
		"ITCH":
		return true
	default:
		return false
	}
}
