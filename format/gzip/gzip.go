package gz

// https://tools.ietf.org/html/rfc1952
// TODO: test name, comment etc
// TODO: verify isize?

import (
	"bytes"
	"compress/flate"
	"hash/crc32"
	"io"

	"github.com/wader/fq/format"
	"github.com/wader/fq/format/registry"
	"github.com/wader/fq/pkg/bitio"
	"github.com/wader/fq/pkg/decode"
)

var probeFormat decode.Group

func init() {
	registry.MustRegister(decode.Format{
		Name:        format.GZIP,
		Description: "gzip compression",
		Groups:      []string{format.PROBE},
		DecodeFn:    gzDecode,
		Dependencies: []decode.Dependency{
			{Names: []string{format.PROBE}, Group: &probeFormat},
		},
	})
}

const delfateMethod = 8

var compressionMethodNames = decode.UToStr{
	delfateMethod: "deflate",
}

var osNames = decode.UToStr{
	0:  "FAT filesystem (MS-DOS, OS/2, NT/Win32)",
	1:  "Amiga",
	2:  "VMS (or OpenVMS)",
	3:  "Unix",
	4:  "VM/CMS",
	5:  "Atari TOS",
	6:  "HPFS filesystem (OS/2, NT)",
	7:  "Macintosh",
	8:  "Z-System",
	9:  "CP/M",
	10: " TOPS-20",
	11: " NTFS filesystem (NT)",
	12: " QDOS",
	13: " Acorn RISCOS",
}

var deflateExtraFlagsNames = decode.UToStr{
	2: "slow",
	4: "fast",
}

func gzDecode(d *decode.D, in interface{}) interface{} {
	d.FieldRawLen("identification", 2*8, d.AssertBitBuf([]byte("\x1f\x8b")))
	compressionMethod := d.FieldU8("compression_method", d.MapUToStrSym(compressionMethodNames))
	hasHeaderCRC := false
	hasExtra := false
	hasName := false
	hasComment := false
	d.FieldStruct("flags", func(d *decode.D) {
		d.FieldBool("text")
		hasHeaderCRC = d.FieldBool("header_crc")
		hasExtra = d.FieldBool("extra")
		hasName = d.FieldBool("name")
		hasComment = d.FieldBool("comment")
		d.FieldU3("reserved")
	})
	d.FieldU32LE("mtime") // TODO: unix time
	switch compressionMethod {
	case delfateMethod:
		d.FieldU8("extra_flags", d.MapUToStrSym(deflateExtraFlagsNames))
	default:
		d.FieldU8("extra_flags")
	}
	d.FieldU8("os", d.MapUToStrSym(osNames))
	if hasExtra {
		// TODO:
		xLen := d.FieldU16("xlen")
		d.FieldRawLen("extra_fields", int64(xLen*8))
	}
	if hasName {
		d.FieldUTF8Null("name")
	}
	if hasComment {
		d.FieldUTF8Null("comment")
	}
	if hasHeaderCRC {
		// TODO: validate
		d.FieldRawLen("header_crc", 16, d.RawHex)
	}

	compressedLen := d.BitsLeft() - ((4 + 4) * 8) // len-(crc32+isize)
	compressedBB := d.FieldRawLen("compressed", compressedLen)
	crc32W := crc32.NewIEEE()

	switch compressionMethod {
	case delfateMethod:
		deflateR := flate.NewReader(compressedBB)
		uncompressed := &bytes.Buffer{}
		if _, err := d.Copy(io.MultiWriter(uncompressed, crc32W), deflateR); err != nil {
			d.Fatalf(err.Error())
		}
		uncompressedBB := bitio.NewBufferFromBytes(uncompressed.Bytes(), -1)
		dv, _, _ := d.FieldTryFormatBitBuf("uncompressed", uncompressedBB, probeFormat, nil)
		if dv == nil {
			d.FieldRootBitBuf("uncompressed", uncompressedBB)
		}
	default:
		d.FieldRawLen("compressed", compressedLen)
	}

	d.FieldRawLen("crc32", 32, d.ValidateBitBuf(bitio.ReverseBytes(crc32W.Sum(nil))), d.RawHex)
	d.FieldU32LE("isize")

	return nil
}
