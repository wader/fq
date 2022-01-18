package tar

// https://www.gnu.org/software/tar/manual/html_node/Standard.html
// TODO: extensions?

import (
	"bytes"
	"strconv"
	"strings"

	"github.com/wader/fq/format"
	"github.com/wader/fq/format/registry"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/scalar"
)

var probeFormat decode.Group

func init() {
	registry.MustRegister(decode.Format{
		Name:        format.TAR,
		Description: "Tar archive",
		Groups:      []string{format.PROBE},
		DecodeFn:    tarDecode,
		Dependencies: []decode.Dependency{
			{Names: []string{format.PROBE}, Group: &probeFormat},
		},
	})
}

func tarDecode(d *decode.D, in interface{}) interface{} {
	const blockBytes = 512
	const blockBits = blockBytes * 8

	mapTrimSpaceNull := scalar.Trim(" \x00")
	mapOctStrToSymU := scalar.Fn(func(s scalar.S) (scalar.S, error) {
		ts := strings.Trim(s.ActualStr(), " ")
		if ts != "" {
			n, err := strconv.ParseUint(ts, 8, 64)
			if err != nil {
				return s, err
			}
			s.Sym = n
		}
		return s, nil
	})
	blockPadding := func(d *decode.D) int64 {
		return (blockBits - (d.Pos() % blockBits)) % blockBits
	}

	// end marker is 512*2 zero bytes
	endMarker := [blockBytes * 2]byte{}
	endMarkerFound := false
	filesCount := 0

	d.FieldArray("files", func(d *decode.D) {
		for !d.End() {
			d.FieldStruct("file", func(d *decode.D) {
				d.FieldUTF8("name", 100, mapTrimSpaceNull)
				d.FieldUTF8NullFixedLen("mode", 8, mapOctStrToSymU)
				d.FieldUTF8NullFixedLen("uid", 8, mapOctStrToSymU)
				d.FieldUTF8NullFixedLen("gid", 8, mapOctStrToSymU)
				sizeS := d.FieldScalarUTF8NullFixedLen("size", 12, mapOctStrToSymU)
				if sizeS.Sym == nil {
					d.Fatalf("could not decode size")
				}
				size := int64(sizeS.SymU()) * 8
				d.FieldUTF8NullFixedLen("mtime", 12, mapOctStrToSymU)
				d.FieldUTF8NullFixedLen("chksum", 8, mapOctStrToSymU)
				d.FieldUTF8("typeflag", 1, mapTrimSpaceNull)
				d.FieldUTF8("linkname", 100, mapTrimSpaceNull)
				d.FieldUTF8("magic", 6, mapTrimSpaceNull, d.AssertStr("ustar"))
				d.FieldUTF8NullFixedLen("version", 2, mapOctStrToSymU)
				d.FieldUTF8("uname", 32, mapTrimSpaceNull)
				d.FieldUTF8("gname", 32, mapTrimSpaceNull)
				d.FieldUTF8NullFixedLen("devmajor", 8, mapOctStrToSymU)
				d.FieldUTF8NullFixedLen("devminor", 8, mapOctStrToSymU)
				d.FieldUTF8("prefix", 155, mapTrimSpaceNull)
				d.FieldRawLen("header_block_padding", blockPadding(d), d.BitBufIsZero())

				dv, _, _ := d.TryFieldFormatLen("data", size, probeFormat, nil)
				if dv == nil {
					d.FieldRawLen("data", size)
				}

				d.FieldRawLen("data_block_padding", blockPadding(d), d.BitBufIsZero())
			})
			filesCount++

			bs := d.PeekBytes(blockBytes * 2)
			if bytes.Equal(bs, endMarker[:]) {
				endMarkerFound = true
				break
			}
		}
	})
	if endMarkerFound {
		d.FieldRawLen("end_marker", int64(len(endMarker))*8)
	}

	if filesCount == 0 {
		d.Errorf("no files found")
	}

	return nil
}
