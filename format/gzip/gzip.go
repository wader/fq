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
	"github.com/wader/fq/pkg/bitio"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
	"github.com/wader/fq/pkg/scalar"
)

var probeGroup decode.Group

func init() {
	interp.RegisterFormat(
		format.Gzip,
		&decode.Format{
			Description: "gzip compression",
			Groups:      []*decode.Group{format.Probe},
			DecodeFn:    gzipDecode,
			Dependencies: []decode.Dependency{
				{Groups: []*decode.Group{format.Probe}, Out: &probeGroup},
			},
		})
}

const deflateMethod = 8

var compressionMethodNames = scalar.UintMapSymStr{
	deflateMethod: "deflate",
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

func gzipDecodeMember(d *decode.D) bitio.ReaderAtSeeker {
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
	d.FieldU32("mtime", scalar.UintActualUnixTimeDescription(time.Second, time.RFC3339))
	switch compressionMethod {
	case deflateMethod:
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
	case deflateMethod:
		// bitio.NewIOReadSeeker implements io.ByteReader so that deflate don't do own
		// buffering and might read more than needed messing up knowing compressed size
		rFn = func(r io.Reader) io.Reader { return flate.NewReader(r) }
	}

	var uncompressedBR bitio.ReaderAtSeeker
	if rFn != nil {
		var readCompressedSize int64
		var err error
		readCompressedSize, uncompressedBR, err =
			d.FieldReaderRange("uncompressed", d.Pos(), d.BitsLeft(), rFn)
		if err != nil {
			d.IOPanic(err, "uncompressed", "FieldReaderRange")
		}
		d.FieldRawLen("compressed", readCompressedSize)
		crc32W := crc32.NewIEEE()
		// TODO: cleanup clone
		d.CopyBits(crc32W, d.CloneReadSeeker(uncompressedBR))
		d.FieldU32("crc32", d.UintValidateBytes(crc32W.Sum(nil)), scalar.UintHex)
		d.FieldU32("isize")
	} else {
		d.Fatalf("unknown compression method %d", compressionMethod)
	}

	return uncompressedBR
}

func gzipDecode(d *decode.D) any {
	d.Endian = decode.LittleEndian

	var brs []bitio.ReadAtSeeker
	d.FieldArray("members", func(d *decode.D) {
		for !d.End() {
			var br bitio.ReadAtSeeker
			d.FieldStruct("member", func(d *decode.D) {
				br = gzipDecodeMember(d)
			})
			brs = append(brs, br)
		}
	})

	if len(brs) == 0 {
		d.Fatalf("no members found")
	}

	cbr, err := bitio.NewMultiReader(brs...)
	if err != nil {
		d.IOPanic(err, "members", "NewMultiReader")
	}
	dv, _, _ := d.TryFieldFormatBitBuf("uncompressed", cbr, &probeGroup, format.Probe_In{})
	if dv == nil {
		d.FieldRootBitBuf("uncompressed", cbr)
	}

	return nil
}
