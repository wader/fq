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

	// 512*2 zero bytes
	endMarker := [blockBytes * 2]byte{}
	foundEndMarker := false

	d.FieldArray("files", func(d *decode.D) {
		for !d.End() {
			d.FieldStruct("file", func(d *decode.D) {
				d.FieldUTF8("name", 100, mapTrimSpaceNull)
				d.FieldUTF8NullFixedLen("mode", 8, mapOctStrToSymU)
				d.FieldUTF8NullFixedLen("uid", 8, mapOctStrToSymU)
				d.FieldUTF8NullFixedLen("gid", 8, mapOctStrToSymU)
				size := d.FieldScalarUTF8NullFixedLen("size", 12, mapOctStrToSymU).SymU()
				d.FieldUTF8NullFixedLen("mtime", 12, mapOctStrToSymU)
				d.FieldUTF8NullFixedLen("chksum", 8, mapOctStrToSymU)
				d.FieldUTF8("typeflag", 1, mapTrimSpaceNull)
				d.FieldUTF8("linkname", 100, mapTrimSpaceNull)
				magic := d.FieldUTF8("magic", 6, mapTrimSpaceNull)
				if magic != "ustar" {
					d.Errorf("invalid magic %s", magic)
				}
				d.FieldUTF8NullFixedLen("version", 2, mapOctStrToSymU)
				d.FieldUTF8("uname", 32, mapTrimSpaceNull)
				d.FieldUTF8("gname", 32, mapTrimSpaceNull)
				d.FieldUTF8NullFixedLen("devmajor", 8, mapOctStrToSymU)
				d.FieldUTF8NullFixedLen("devminor", 8, mapOctStrToSymU)
				d.FieldUTF8("prefix", 155, mapTrimSpaceNull)
				d.FieldRawLen("header_block_padding", blockPadding(d), d.BitBufIsZero())

				dv, _, _ := d.TryFieldFormatLen("data", int64(size)*8, probeFormat, nil)
				if dv == nil {
					d.FieldRawLen("data", int64(size)*8)
				}

				d.FieldRawLen("data_block_padding", blockPadding(d), d.BitBufIsZero())
			})

			bs := d.PeekBytes(512 * 2)
			if bytes.Equal(bs, endMarker[:]) {
				foundEndMarker = true
				break
			}
		}
	})
	d.FieldRawLen("end_marker", 512*2*8)

	if !foundEndMarker {
		d.Errorf("no files found")
	}

	return nil
}
