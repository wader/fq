package tar

// https://www.gnu.org/software/tar/manual/html_node/Standard.html
// TODO: extensions?

import (
	"bytes"
	"fmt"
	"fq/pkg/decode"
	"fq/pkg/format"
	"strconv"
	"strings"
)

var all []*decode.Format

var File = format.MustRegister(&decode.Format{
	Name:  "tar",
	MIMEs: []string{"application/x-tar"},
	New:   func() decode.Decoder { return &FileDecoder{} },
	Deps: []decode.Dep{
		{Names: []string{"probeable"}, Formats: &all},
	},
})

// Decoder is a tar decoder
type FileDecoder struct {
	decode.Common
}

// Decode tar
func (d *FileDecoder) Decode() {
	str := func(nBytes int64) string {
		s := d.UTF8(nBytes)
		ts := strings.Trim(s, "\x00")
		return ts
	}
	fieldStr := func(name string, nBytes int64) string {
		return d.FieldStrFn(name, func() (string, string) {
			return str(nBytes), ""
		})
	}
	fieldNumStr := func(name string, nBytes int64) uint64 {
		return d.FieldUFn(name, func() (uint64, decode.DisplayFormat, string) {
			ts := strings.Trim(str(nBytes), "0 \x00")
			if ts == "" {
				return 0, decode.NumberDecimal, ts
			}
			n, err := strconv.ParseUint(ts, 8, 64)
			if err != nil {
				d.Invalid(fmt.Sprintf("failed to parse %s number %s: %s", name, ts, err))
			}
			return n, decode.NumberDecimal, ts
		})
	}
	fieldBlockPadding := func(name string) {
		const blockBits = 512 * 8
		blockPadding := (blockBits - (d.Pos() % blockBits)) % blockBits
		if blockPadding > 0 {
			d.FieldValidateZeroPadding(name, blockPadding)
		}
	}

	// 512*2 zero bytes
	endMarker := [512 * 2]byte{}
	validFiles := 0

	d.Array("file", func() {
		for !d.End() {
			name := str(100)
			d.SeekRel(-100 * 8)
			d.FieldNoneFn(name, func() {
				fieldStr("name", 100)
				fieldNumStr("mode", 8)
				fieldNumStr("uid", 8)
				fieldNumStr("gid", 8)
				size := fieldNumStr("size", 12)
				fieldNumStr("mtime", 12)
				fieldNumStr("chksum", 8)
				fieldStr("typeflag", 1)
				fieldStr("linkname", 100)
				fieldStr("magic", 6)
				fieldNumStr("version", 2)
				fieldStr("uname", 32)
				fieldStr("gname", 32)
				fieldNumStr("devmajor", 8)
				fieldNumStr("devminor", 8)
				fieldStr("prefix", 155)
				fieldBlockPadding("header_block_padding")
				if size > 0 {
					d.FieldDecodeLen("data", int64(size)*8, all)
				}
				fieldBlockPadding("data_block_padding")
			})
			validFiles++

			bs := d.PeekBytes(512 * 2)
			if bytes.Equal(bs, endMarker[:]) {
				d.FieldBitBufLen("end_marker", 512*2*8)
				break
			}
		}
	})

	if validFiles == 0 {
		d.Invalid("no files found")
	}
}
