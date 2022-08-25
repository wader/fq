package tar

// https://www.gnu.org/software/tar/manual/html_node/Standard.html
// TODO: extensions?

import (
	"bytes"

	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
	"github.com/wader/fq/pkg/scalar"
)

var probeFormat decode.Group

func init() {
	interp.RegisterFormat(decode.Format{
		Name:        format.TAR,
		Description: "Tar archive",
		Groups:      []string{format.PROBE},
		DecodeFn:    tarDecode,
		Dependencies: []decode.Dependency{
			{Names: []string{format.PROBE}, Group: &probeFormat},
		},
	})
}

func tarDecode(d *decode.D, _ any) any {
	const blockBytes = 512
	const blockBits = blockBytes * 8

	mapTrimSpaceNull := scalar.ActualTrim(" \x00")
	blockPadding := func(d *decode.D) int64 {
		return (blockBits - (d.Pos() % blockBits)) % blockBits
	}

	// end marker is 512*2 zero bytes
	endMarker := [blockBytes * 2]byte{}
	var endMarkerStart int64
	var endMarkerEnd int64
	filesCount := 0

	d.FieldArray("files", func(d *decode.D) {
		for !d.End() {
			d.FieldStruct("file", func(d *decode.D) {
				d.FieldUTF8("name", 100, mapTrimSpaceNull)
				d.FieldUTF8NullFixedLen("mode", 8, scalar.TrySymUParseUint(8))
				d.FieldUTF8NullFixedLen("uid", 8, scalar.TrySymUParseUint(8))
				d.FieldUTF8NullFixedLen("gid", 8, scalar.TrySymUParseUint(8))
				sizeS := d.FieldScalarUTF8NullFixedLen("size", 12, scalar.TrySymUParseUint(8))
				if sizeS.Sym == nil {
					d.Fatalf("could not decode size")
				}
				size := int64(sizeS.SymU()) * 8
				d.FieldUTF8NullFixedLen("mtime", 12, scalar.TrySymUParseUint(8), scalar.DescriptionSymUUnixTime)
				d.FieldUTF8NullFixedLen("chksum", 8, scalar.TrySymUParseUint(8))
				d.FieldUTF8("typeflag", 1, mapTrimSpaceNull)
				d.FieldUTF8("linkname", 100, mapTrimSpaceNull)
				d.FieldUTF8("magic", 6, mapTrimSpaceNull, d.AssertStr("ustar"))
				d.FieldUTF8NullFixedLen("version", 2, scalar.TrySymUParseUint(8))
				d.FieldUTF8("uname", 32, mapTrimSpaceNull)
				d.FieldUTF8("gname", 32, mapTrimSpaceNull)
				d.FieldUTF8NullFixedLen("devmajor", 8, scalar.TrySymUParseUint(8))
				d.FieldUTF8NullFixedLen("devminor", 8, scalar.TrySymUParseUint(8))
				d.FieldUTF8("prefix", 155, mapTrimSpaceNull)
				d.FieldRawLen("header_block_padding", blockPadding(d), d.BitBufIsZero())

				d.FieldFormatOrRawLen("data", size, probeFormat, nil)

				d.FieldRawLen("data_block_padding", blockPadding(d), d.BitBufIsZero())
			})
			filesCount++

			if d.BitsLeft() >= int64(len(endMarker))*8 && bytes.Equal(d.PeekBytes(len(endMarker)), endMarker[:]) {
				endMarkerStart = d.Pos()
				// consensus seems to be to allow more than 2 zero blocks at end
				d.SeekRel(int64(len(endMarker)) * 8)
				zeroBlock := [blockBytes]byte{}
				for d.BitsLeft() >= blockBytes*8 && bytes.Equal(d.PeekBytes(blockBytes), zeroBlock[:]) {
					d.SeekRel(int64(len(zeroBlock)) * 8)
				}
				endMarkerEnd = d.Pos()
				break
			}
		}
	})
	endMarkerSize := endMarkerEnd - endMarkerStart
	if endMarkerSize > 0 {
		d.RangeFn(endMarkerStart, endMarkerSize, func(d *decode.D) {
			d.FieldRawLen("end_marker", d.BitsLeft())
		})
	}

	if filesCount == 0 {
		d.Errorf("no files found")
	}

	return nil
}
