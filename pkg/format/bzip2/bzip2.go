package bzip2

// https://en.wikipedia.org/wiki/Bzip2
// https://github.com/dsnet/compress/blob/master/doc/bzip2-format.pdf
// TODO: multiple streams, possible to figure out length of compressed? use footer magic?
// TODO: empty file, no streams

import (
	"bytes"
	"compress/bzip2"
	"fq/pkg/bitio"
	"fq/pkg/decode"
	"fq/pkg/format"
	"hash/crc32"
	"io"
)

var probeFormat []*decode.Format

func init() {
	format.MustRegister(&decode.Format{
		Name:        format.BZIP2,
		Description: "bzip2 compression",
		Groups:      []string{format.PROBE},
		DecodeFn:    gzDecode,
		Dependencies: []decode.Dependency{
			{Names: []string{format.PROBE}, Formats: &probeFormat},
		},
	})
}

func gzDecode(d *decode.D, in interface{}) interface{} {
	// moreStreams := true

	// d.FieldArrayFn("streams", func(d *decode.D) {
	// 	for moreStreams {
	// d.FieldStructFn("stream", func(d *decode.D) {

	d.FieldValidateUTF8("magic", "BZ")
	d.FieldU8("version")
	d.FieldU8("hundred_k_blocksize")

	d.FieldStructFn("block", func(d *decode.D) {
		const blockHeaderMagic = 0x31_41_59_26_53_59
		// if d.PeekBits(48) != blockHeaderMagic {
		// 	moreStreams = false
		// 	return
		// }
		d.FieldValidateUFn("compressed_magic", blockHeaderMagic, d.U48)
		d.FieldU32("crc")
		d.FieldU1("randomised")
		d.FieldU24("origptr")
		d.FieldU16("syncmapl1")

		d.SeekRel(-16)
		ranges := 0
		for i := 0; i < 16; i++ {
			if d.Bool() {
				ranges++
			}
		}
		d.FieldBitBufLen("syncmapl2", int64(ranges)*16)
		numTrees := d.FieldU3("num_trees")
		selectorsUsed := d.FieldU15("num_sels")
		selectorsI := uint64(0)
		d.FieldArrayLoopFn("selector_list", func() bool { return selectorsI < selectorsUsed }, func(d *decode.D) {
			d.FieldU1("selector")
			selectorsI++
		})
		treesI := uint64(0)
		d.FieldArrayLoopFn("trees", func() bool { return treesI < numTrees }, func(d *decode.D) {
			d.FieldUFn("tree", func() (uint64, decode.DisplayFormat, string) {
				l := d.U5()
				if !d.Bool() {
					return l, decode.NumberDecimal, ""
				}
				if d.Bool() {
					l--
				} else {
					l++
				}
				return l, decode.NumberDecimal, ""
			})
			treesI++
		})
	})

	compressedBB := d.BitBufRange(0, d.Len())
	deflateR := bzip2.NewReader(compressedBB)
	uncompressed := &bytes.Buffer{}
	crc32W := crc32.NewIEEE()
	if _, err := io.Copy(io.MultiWriter(uncompressed, crc32W), deflateR); err != nil {
		d.Invalid(err.Error())
	}
	// calculatedCRC32 := crc32W.Sum(nil)
	uncompressedBB := bitio.NewBufferFromBytes(uncompressed.Bytes(), -1)
	dv, _, _ := d.FieldTryDecodeBitBuf("uncompressed", uncompressedBB, probeFormat)
	if dv == nil {
		d.FieldRootBitBuf("uncompressed", uncompressedBB)
	}

	// if calculatedCRC32 != nil {
	// 	d.FieldChecksumLen("crc32", 32, calculatedCRC32, decode.LittleEndian)
	// } else {
	// 	d.FieldU32LE("crc32")
	// }

	// d.FieldU48("footer_magic")
	// d.FieldU32("crc")
	// byte align padding
	// })

	// 		moreStreams = false

	// 	}
	// })

	return nil
}
