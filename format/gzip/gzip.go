package gz

// https://tools.ietf.org/html/rfc1952
// TODO: test name, comment etc
// TODO: verify isize?

import (
	"compress/flate"
	"hash/crc32"
	"io"
	"time"

	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
	"github.com/wader/fq/pkg/scalar"
)

var probeFormat decode.Group

func init() {
	interp.RegisterFormat(decode.Format{
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

var compressionMethodNames = scalar.UintMapSymStr{
	delfateMethod: "deflate",
}

var osNames = scalar.UintMapSymStr{
	0:  "fat",
	1:  "amiga",
	2:  "vms",
	3:  "unix",
	4:  "vm_cms",
	5:  "atari_tOS",
	6:  "hpfs",
	7:  "Mmcintosh",
	8:  "z_system",
	9:  "cpm",
	10: "tops_20",
	11: "ntfs",
	12: "qdos",
	13: "acorn_riscos",
}

var deflateExtraFlagsNames = scalar.UintMapSymStr{
	2: "slow",
	4: "fast",
}

func gzDecode(d *decode.D, _ any) any {
	d.Endian = decode.LittleEndian

	d.FieldRawLen("identification", 2*8, d.AssertBitBuf([]byte("\x1f\x8b")))
	compressionMethod := d.FieldU8("compression_method", compressionMethodNames)
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
	d.FieldU32("mtime", scalar.UintActualUnixTime(time.RFC3339))
	switch compressionMethod {
	case delfateMethod:
		d.FieldU8("extra_flags", deflateExtraFlagsNames)
	default:
		d.FieldU8("extra_flags")
	}
	d.FieldU8("os", osNames)
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
		d.FieldRawLen("header_crc", 16, scalar.RawHex)
	}

	var rFn func(r io.Reader) io.Reader
	switch compressionMethod {
	case delfateMethod:
		// bitio.NewIOReadSeeker implements io.ByteReader so that deflate don't do own
		// buffering and might read more than needed messing up knowing compressed size
		rFn = func(r io.Reader) io.Reader { return flate.NewReader(r) }
	}

	if rFn != nil {
		readCompressedSize, uncompressedBR, dv, _, _ := d.TryFieldReaderRangeFormat("uncompressed", d.Pos(), d.BitsLeft(), rFn, probeFormat, nil)
		if uncompressedBR != nil {
			if dv == nil {
				d.FieldRootBitBuf("uncompressed", uncompressedBR)
			}
			d.FieldRawLen("compressed", readCompressedSize)
			crc32W := crc32.NewIEEE()
			// TODO: cleanup clone
			d.CopyBits(crc32W, d.CloneReadSeeker(uncompressedBR))
			d.FieldU32("crc32", d.UintValidateBytes(crc32W.Sum(nil)), scalar.UintHex)
			d.FieldU32("isize")
		}
	}

	return nil
}
