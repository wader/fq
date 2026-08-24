package mosaic

import (
	"embed"

	"github.com/wader/fq/format"
	"github.com/wader/fq/format/ebml"
	"github.com/wader/fq/format/mosaic/ebml_mosaic"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
	"github.com/wader/fq/pkg/scalar"
)

//go:embed mosaic.md
var mosaicFS embed.FS

func init() {
	interp.RegisterFormat(
		format.MOSAIC,
		&decode.Format{
			Description: "MOdular Storage of Archived and Indexed Contents",
			Groups:      []*decode.Group{format.Probe},
			DecodeFn:    mosaicDecode,
		})
	interp.RegisterFS(mosaicFS)
}

func decodeMaster(d *decode.D, bitsLimit int64, elm *ebml.Master) {
	tagEndBit := d.Pos() + bitsLimit

	for d.Pos() < tagEndBit && !d.End() {
		var childElm ebml.Element
		elm.PeekNextElement(d, &childElm)

		if childElm.IsSingleton() {
			d.FieldStruct(childElm.GetName(), func(d *decode.D) {
				decodeElement(d, childElm)
			})
		} else {
			arrayName := childElm.GetName()
			d.FieldArray(arrayName, func(d *decode.D) {
				// loop over siblings as long as they are the same element
				for d.Pos() < tagEndBit && !d.End() {
					elm.PeekNextElement(d, &childElm)
					if childElm.GetName() != arrayName {
						break
					}
					d.FieldStruct(arrayName, func(d *decode.D) {
						decodeElement(d, childElm)
					})
				}
			})
		}
	}
}

// For master elements whose children are not grouped by type
func decodeElements(d *decode.D, bitsLimit int64, elm *ebml.Master) {
	tagEndBit := d.Pos() + bitsLimit

	d.FieldArray("elements", func(d *decode.D) {
		var childElm ebml.Element
		for d.Pos() < tagEndBit && !d.End() {

			elm.PeekNextElement(d, &childElm)
			d.FieldStruct(childElm.GetName(), func(d *decode.D) {
				decodeElement(d, childElm)
			})
		}
	})
}

func decodeElement(d *decode.D, childElm ebml.Element) {
	tagID := d.FieldUintFn("id", ebml.DecodeRawVint, ebml.ElementIDMapper(childElm))
	d.FieldValueStr("type", childElm.GetType())

	const maxStringTagSize = 100 * 1024 * 1024
	tagSize := d.FieldUintFn("size", ebml.DecodeVint, scalar.UintMapDescription{
		0xffffffffffffff: "Unknown size",
	})

	// assert sane tag size
	// strings are limited because they are read into memory
	switch childElm.(type) {
	case *ebml.Integer,
		*ebml.Uinteger,
		*ebml.Float:
		if tagSize > 8 {
			d.Fatalf("invalid tagSize %d for number type", tagSize)
		}
	case *ebml.String,
		*ebml.UTF8:
		if tagSize > maxStringTagSize {
			d.Errorf("tagSize %d > maxStringTagSize %d", tagSize, maxStringTagSize)
		}
	case *ebml.Unknown,
		*ebml.Binary,
		*ebml.Date,
		*ebml.Master:
		// nop
	}

	switch childElm := childElm.(type) {
	case *ebml.Unknown:
		d.FieldRawLen("data", int64(tagSize)*8)
	case *ebml.Integer:
		var sm []scalar.SintMapper
		if childElm.Enums != nil {
			sm = append(sm, scalar.SintFn(func(s scalar.Sint) (scalar.Sint, error) {
				if e, ok := childElm.Enums[s.Actual]; ok {
					s.Sym = e.Name
					s.Description = e.Description
				}
				return s, nil
			}))
		}
		d.FieldS("value", int(tagSize)*8, sm...)
	case *ebml.Uinteger:
		var sm []scalar.UintMapper
		if childElm.Enums != nil {
			sm = append(sm, scalar.UintFn(func(s scalar.Uint) (scalar.Uint, error) {
				if e, ok := childElm.Enums[s.Actual]; ok {
					s.Sym = e.Name
					s.Description = e.Description
				}
				return s, nil
			}))
		}
		d.FieldU("value", int(tagSize)*8, sm...)
	case *ebml.Float:
		d.FieldF("value", int(tagSize)*8)
	case *ebml.String:
		var sm []scalar.StrMapper
		sm = append(sm, scalar.StrFn(func(s scalar.Str) (scalar.Str, error) {
			if e, ok := childElm.Enums[s.Actual]; ok {
				s.Sym = e.Name
				s.Description = e.Description
			}
			return s, nil
		}))
		v := d.FieldUTF8("value", int(tagSize), sm...)
		if tagID == ebml.DocTypeID && v != "mosaic" {
			d.Errorf("EBML doctype is not mosaic")
		}
	case *ebml.UTF8:
		d.FieldUTF8NullFixedLen("value", int(tagSize))
	case *ebml.Date: // that's not expected in MOSAIC
		d.FieldRawLen("value", int64(tagSize)*8)
	case *ebml.Binary:
		d.FieldRawLen("value", int64(tagSize)*8)

	case *ebml.Master:
		if tagID == ebml_mosaic.IdxUnrolledID {
			decodeElements(d, int64(tagSize)*8, childElm)
		} else {
			decodeMaster(d, int64(tagSize)*8, childElm)
		}
	}
}

func mosaicDecode(d *decode.D) any {
	ebmlHeaderID := uint64(0x1a45dfa3)
	if d.PeekUintBits(32) != ebmlHeaderID {
		d.Fatalf("no EBML header found")
	}

	decodeMaster(d, d.BitsLeft(), ebml_mosaic.RootElement)

	return nil
}
