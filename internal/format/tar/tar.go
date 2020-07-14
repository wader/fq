package tar

// https://www.gnu.org/software/tar/manual/html_node/Standard.html
// TODO: extensions?

import (
	"bytes"
	"fmt"
	"fq/internal/decode"
	"strconv"
	"strings"
)

var File = &decode.Format{
	Name: "tar",
	MIME: "",
	New:  func() decode.Decoder { return &FileDecoder{} },
}

// Decoder is a tar decoder
type FileDecoder struct {
	decode.Common
}

// Decode tar
func (d *FileDecoder) Decode() {
	strFn := func(name string, nBytes uint64) string {
		return d.FieldStrFn(name, func() (string, string) {
			s := d.UTF8(nBytes)
			ts := strings.Trim(s, "\x00")
			return ts, ""
		})
	}
	numStrFn := func(name string, nBytes uint64) uint64 {
		return d.FieldUFn(name, func() (uint64, decode.NumberFormat, string) {
			s := d.UTF8(nBytes)
			ts := strings.Trim(s, "0 \x00")
			if ts == "" {
				return 0, decode.NumberDecimal, s
			}
			n, err := strconv.ParseUint(ts, 8, 64)
			if err != nil {
				d.Invalid(fmt.Sprintf("failed to parse %s number %s: %s", name, ts, err))
			}
			return n, decode.NumberDecimal, s
		})
	}
	blockPaddingFn := func() {
		const blockBits = 512 * 8
		blockPadding := (blockBits - (d.Pos() % blockBits)) % blockBits
		if blockPadding > 0 {
			d.FieldValidateZeroPadding("block_padding", blockPadding)
		}
	}

	// 512*2 zero bytes
	endMarker := [512 * 2]byte{}
	for !d.End() {
		name := d.UTF8(100)
		d.SeekRel(-100 * 8)
		d.FieldNoneFn(name, func() {
			strFn("name", 100)
			numStrFn("mode", 8)
			numStrFn("uid", 8)
			numStrFn("gid", 8)
			size := numStrFn("size", 12)
			numStrFn("mtime", 12)
			numStrFn("chksum", 8)
			strFn("typeflag", 1)
			strFn("linkname", 100)
			strFn("magic", 6)
			numStrFn("version", 2)
			strFn("uname", 32)
			strFn("gname", 32)
			numStrFn("devmajor", 8)
			numStrFn("devminor", 8)
			strFn("prefix", 155)
			blockPaddingFn()
			if size > 0 {
				d.FieldDecodeLen("data", size*8)
			}
			blockPaddingFn()
		})
		bs := d.PeekBytes(512 * 2)
		if bytes.Compare(bs, endMarker[:]) == 0 {
			d.FieldBytesLen("end_marker", 512*2)
			break
		}
	}
}
