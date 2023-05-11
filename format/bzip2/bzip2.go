package bzip2

// https://en.wikipedia.org/wiki/Bzip2
// https://github.com/dsnet/compress/blob/master/doc/bzip2-format.pdf
// TODO: multiple streams, possible to figure out length of compressed? use footer magic?
// TODO: empty file, no streams

import (
	"compress/bzip2"
	"encoding/binary"
	"hash/crc32"
	"io"
	"math/bits"

	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/bitio"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
	"github.com/wader/fq/pkg/scalar"
)

var probeGroup decode.Group

func init() {
	interp.RegisterFormat(
		format.Bzip2,
		&decode.Format{
			Description: "bzip2 compression",
			Groups:      []*decode.Group{format.Probe},
			DecodeFn:    bzip2Decode,
			Dependencies: []decode.Dependency{
				{Groups: []*decode.Group{format.Probe}, Out: &probeGroup},
			},
		})
}

const blockMagic = 0x31_41_59_26_53_59
const footerMagic = 0x17_72_45_38_50_90

type bitFlipReader struct {
	r io.Reader
}

func (bfr bitFlipReader) Read(p []byte) (n int, err error) {
	n, err = bfr.r.Read(p)
	for i := 0; i < n; i++ {
		p[i] = bits.Reverse8(p[i])
	}
	return n, err
}

func bzip2Decode(d *decode.D) any {
	// moreStreams := true

	// d.FieldArray("streams", func(d *decode.D) {
	// 	for moreStreams {
	// d.FieldStruct("stream", func(d *decode.D) {

	var blockCRCValue *decode.Value
	var streamCRCN uint32

	d.FieldUTF8("magic", 2, d.StrAssert("BZ"))
	d.FieldU8("version")
	d.FieldU8("hundred_k_blocksize")

	d.FieldStruct("block", func(d *decode.D) {
		// if d.PeekBits(48) != blockHeaderMagic {
		// 	moreStreams = false
		// 	return
		// }
		d.FieldU48("magic", d.UintAssert(blockMagic), scalar.UintHex)
		d.FieldU32("crc", scalar.UintHex)
		blockCRCValue = d.FieldGet("crc")
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
		d.FieldRawLen("syncmapl2", int64(ranges)*16)
		numTrees := d.FieldU3("num_trees")
		selectorsUsed := d.FieldU15("num_sels")
		selectorsI := uint64(0)
		d.FieldArrayLoop("selector_list", func() bool { return selectorsI < selectorsUsed }, func(d *decode.D) {
			d.FieldU1("selector")
			selectorsI++
		})
		treesI := uint64(0)
		d.FieldArrayLoop("trees", func() bool { return treesI < numTrees }, func(d *decode.D) {
			d.FieldUintFn("tree", func(d *decode.D) uint64 {
				l := d.U5()
				if !d.Bool() {
					return l
				}
				if d.Bool() {
					l--
				} else {
					l++
				}
				return l
			})
			treesI++
		})
	})

	compressedStart := d.Pos()

	readCompressedSize, uncompressedBR, dv, _, _ :=
		d.TryFieldReaderRangeFormat("uncompressed", 0, d.Len(), bzip2.NewReader, &probeGroup, format.Probe_In{})
	if uncompressedBR != nil {
		if dv == nil {
			d.FieldRootBitBuf("uncompressed", uncompressedBR)
		}

		blockCRC32W := crc32.NewIEEE()
		d.Copy(blockCRC32W, bitFlipReader{bitio.NewIOReader(uncompressedBR)})
		blockCRC32N := bits.Reverse32(binary.BigEndian.Uint32(blockCRC32W.Sum(nil)))
		_ = blockCRCValue.TryUintScalarFn(d.UintValidate(uint64(blockCRC32N)))
		streamCRCN = blockCRC32N ^ ((streamCRCN << 1) | (streamCRCN >> 31))

		// HACK: bzip2.NewReader will read from start of whole buffer and then we figure out compressedSize ourself
		// "It is important to note that none of the fields within a StreamBlock or StreamFooter are necessarily byte-aligned"
		const footerByteSize = 10
		compressedSize := (readCompressedSize - compressedStart) - footerByteSize*8
		for i := 0; i < 8; i++ {
			d.SeekAbs(compressedStart + compressedSize)
			if d.PeekUintBits(48) == footerMagic {
				break
			}
			compressedSize--
		}
		d.SeekAbs(compressedStart)

		d.FieldRawLen("compressed", compressedSize)

		d.FieldStruct("footer", func(d *decode.D) {
			d.FieldU48("magic", d.UintAssert(footerMagic), scalar.UintHex)
			// TODO: crc of block crcs
			d.FieldU32("crc", scalar.UintHex, d.UintValidate(uint64(streamCRCN)))
			d.FieldRawLen("padding", int64(d.ByteAlignBits()))
		})
	}

	// 		moreStreams = false
	// 	}
	// })

	return nil
}
