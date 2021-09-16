package tar

// https://www.gnu.org/software/tar/manual/html_node/Standard.html
// TODO: extensions?

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/wader/fq/format"
	"github.com/wader/fq/format/registry"
	"github.com/wader/fq/pkg/decode"
)

var probeFormat []*decode.Format

func init() {
	registry.MustRegister(&decode.Format{
		Name:        format.TAR,
		Description: "Tar archive",
		Groups:      []string{format.PROBE},
		DecodeFn:    tarDecode,
		Dependencies: []decode.Dependency{
			{Names: []string{format.PROBE}, Formats: &probeFormat},
		},
	})
}

func tarDecode(d *decode.D, in interface{}) interface{} {
	str := func(nBytes int) string {
		s := d.UTF8(nBytes)
		ts := strings.Trim(s, " \x00")
		return ts
	}
	fieldStr := func(d *decode.D, name string, nBytes int) string {
		return d.FieldStrFn(name, func() (string, string) {
			return str(nBytes), ""
		})
	}
	fieldNumStr := func(d *decode.D, name string, nBytes int) uint64 {
		return d.FieldUFn(name, func() (uint64, decode.DisplayFormat, string) {
			ts := strings.TrimLeft(str(nBytes), "0 \x00")
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
	fieldBlockPadding := func(d *decode.D, name string) {
		const blockBits = 512 * 8
		blockPadding := (blockBits - (d.Pos() % blockBits)) % blockBits
		if blockPadding > 0 {
			d.FieldValidateZeroPadding(name, int(blockPadding))
		}
	}

	// 512*2 zero bytes
	endMarker := [512 * 2]byte{}
	foundEndMarker := false

	d.FieldArrayFn("files", func(d *decode.D) {
		for !d.End() {
			d.FieldStructFn("file", func(d *decode.D) {
				fieldStr(d, "name", 100)
				fieldNumStr(d, "mode", 8)
				fieldNumStr(d, "uid", 8)
				fieldNumStr(d, "gid", 8)
				size := fieldNumStr(d, "size", 12)
				fieldNumStr(d, "mtime", 12)
				fieldNumStr(d, "chksum", 8)
				fieldStr(d, "typeflag", 1)
				fieldStr(d, "linkname", 100)
				magic := fieldStr(d, "magic", 6)
				if magic != "ustar" {
					d.Invalid(fmt.Sprintf("invalid magic %s", magic))
				}
				fieldNumStr(d, "version", 2)
				fieldStr(d, "uname", 32)
				fieldStr(d, "gname", 32)
				fieldNumStr(d, "devmajor", 8)
				fieldNumStr(d, "devminor", 8)
				fieldStr(d, "prefix", 155)
				fieldBlockPadding(d, "header_block_padding")
				if size > 0 {
					dv, _, _ := d.FieldTryFormatLen("data", int64(size)*8, probeFormat, nil)
					if dv == nil {
						d.FieldBitBufLen("data", int64(size)*8)
					}
				}
				fieldBlockPadding(d, "data_block_padding")
			})

			bs := d.PeekBytes(512 * 2)
			if bytes.Equal(bs, endMarker[:]) {
				foundEndMarker = true
				break
			}
		}
	})
	d.FieldBitBufLen("end_marker", 512*2*8)

	if !foundEndMarker {
		d.Invalid("no files found")
	}

	return nil
}
