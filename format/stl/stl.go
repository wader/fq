package stl

// see https://www.fabbers.com/tech/STL_Format
// see https://paulbourke.net/dataformats/stl/

import (
	"embed"

	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
)

//go:embed stl.md
var stlFS embed.FS

func init() {
	interp.RegisterFormat(
		format.STL,
		&decode.Format{
			Description: "Stereolithography",
			DecodeFn:    decodeSTL,
		})
	interp.RegisterFS(stlFS)
}

const (
	headerLength = 80
)

func decodeVector(d *decode.D) {
	d.FieldF32("x")
	d.FieldF32("y")
	d.FieldF32("z")
}

func decodeFacet(d *decode.D) {
	d.FieldStruct("normal", decodeVector)
	d.FieldStructNArray("vertices", "vertex", 3, decodeVector)
	// TODO support color of VisCAM and SolidView
	// TODO support color of Materialise Magics
	attributeByteCount := d.FieldU16("attribute_byte_count")
	if attributeByteCount > 0 {
		d.FieldRawLen("attribute", int64(attributeByteCount)*8)
	}
}

func decodeSTLModel(d *decode.D) {
	// TODO support color of Materialise Magics
	d.FieldUTF8NullFixedLen("header", headerLength)

	numFacets := d.FieldU32("num_facets")
	d.FieldStructNArray("facets", "facet", int64(numFacets), decodeFacet)
}

func decodeSTL(d *decode.D) any {
	d.Endian = decode.LittleEndian

	decodeSTLModel(d)

	return nil
}
