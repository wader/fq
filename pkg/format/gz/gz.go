package gz

// https://tools.ietf.org/html/rfc1952
// TODO: test name, comment etc

import (
	"bytes"
	"compress/flate"
	"fq/pkg/bitio"
	"fq/pkg/decode"
	"fq/pkg/format"
	"hash/crc32"
	"io"
)

var probeFormat []*decode.Format

func init() {
	format.MustRegister(&decode.Format{
		Name:        format.GZIP,
		Description: "GZIP",
		Groups:      []string{format.PROBE},
		MIMEs:       []string{"application/gzip"},
		DecodeFn:    gzDecode,
		Dependencies: []decode.Dependency{
			{Names: []string{format.PROBE}, Formats: &probeFormat},
		},
	})
}

const delfateMethod = 8

var osNames = map[uint64]string{
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

func fieldStrZeroTerminated(d *decode.D, name string) string {
	return d.FieldStrFn(name, func() (string, string) {
		return strZeroTerminated(d), ""
	})
}

func strZeroTerminated(d *decode.D) string {
	c := d.PeekFindByte(0, -1)
	s := d.UTF8(int(c))
	return s[:len(s)-1]
}

func gzDecode(d *decode.D, in interface{}) interface{} {
	d.FieldValidateUTF8("identification", "\x1f\x8b")
	compressionMethod := d.FieldUFn("compression_method", func() (uint64, decode.DisplayFormat, string) {
		n := d.U8()
		if n == delfateMethod {
			return n, decode.NumberDecimal, "deflate"
		}
		return n, decode.NumberDecimal, "unknown"
	})
	hasHeaderCRC := false
	hasExtra := false
	hasName := false
	hasComment := false
	d.FieldStructFn("flags", func(d *decode.D) {
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
		d.FieldUFn("extra_flags", func() (uint64, decode.DisplayFormat, string) {
			n := d.U8()
			switch n {
			case 2:
				return n, decode.NumberDecimal, "slow"
			case 4:
				return n, decode.NumberDecimal, "fast"
			default:
				return n, decode.NumberDecimal, "unknown"
			}
		})
	default:
		d.FieldU8("extra_flags")
	}
	d.FieldStringMapFn("os", osNames, "unknown", d.U8)
	if hasExtra {
		// TODO:
		xLen := d.FieldU16("xlen")
		d.FieldBitBufLen("extra_fields", int64(xLen*8))
	}
	if hasName {
		fieldStrZeroTerminated(d, "name")
	}
	if hasComment {
		fieldStrZeroTerminated(d, "comment")
	}
	if hasHeaderCRC {
		d.FieldU16LE("header_crc")
	}

	compressedLen := d.BitsLeft() - ((4 + 4) * 8) // len-(crc32+isize)
	compressedBB := d.FieldBitBufLen("compressed", compressedLen)
	var calculatedCRC32 []byte

	switch compressionMethod {
	case delfateMethod:
		deflateR := flate.NewReader(compressedBB)
		uncompressed := &bytes.Buffer{}
		crc32W := crc32.NewIEEE()
		if _, err := io.Copy(io.MultiWriter(uncompressed, crc32W), deflateR); err != nil {
			d.Invalid(err.Error())
		}
		calculatedCRC32 = crc32W.Sum(nil)
		uncompressedBB := bitio.NewBufferFromBytes(uncompressed.Bytes(), -1)
		v, _, _ := d.FieldTryDecodeBitBuf("uncompressed", uncompressedBB, probeFormat)
		if v == nil {
			d.FieldRootBitBuf("uncompressed", uncompressedBB)
		}
	default:
		d.FieldBitBufLen("compressed", compressedLen)
	}

	if calculatedCRC32 != nil {
		d.FieldChecksumLen("crc32", 32, calculatedCRC32, decode.LittleEndian)
	} else {
		d.FieldU32LE("crc32")
	}
	d.FieldU32LE("isize")

	return nil
}
