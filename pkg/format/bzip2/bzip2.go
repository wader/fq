package bzip2

// https://en.wikipedia.org/wiki/Bzip2
// TODO: test name, comment etc
// TODO: proable

import (
	"fq/pkg/decode"
	"fq/pkg/format"
)

var probeFormat []*decode.Format

func init() {
	format.MustRegister(&decode.Format{
		Name:        format.BZIP2,
		Description: "bzip2 compression",
		//Groups:      []string{format.PROBE},
		MIMEs:    []string{"application/gzip"},
		DecodeFn: gzDecode,
		Dependencies: []decode.Dependency{
			{Names: []string{format.PROBE}, Formats: &probeFormat},
		},
	})
}

func gzDecode(d *decode.D, in interface{}) interface{} {
	moreStreams := true

	d.FieldArrayFn("streams", func(d *decode.D) {
		for moreStreams {
			d.FieldStructFn("stream", func(d *decode.D) {
				d.FieldValidateUTF8("magic", "BZ")
				d.FieldU8("version")
				d.FieldU8("hundred_k_blocksize")

				d.FieldStructFn("block", func(d *decode.D) {
					const blockHeaderMagic = 0x31_41_59_26_53_59
					if d.PeekBits(48) != blockHeaderMagic {
						moreStreams = false
						return
					}
					d.FieldU48("compressed_magic")
					d.FieldU32("crc")
					d.FieldU1("randomised")
					d.FieldU24("origptr")
					d.FieldU16("huffman_used_map")

					d.SeekRel(-16)
					ranges := 0
					for i := 0; i < 16; i++ {
						if d.Bool() {
							ranges++
						}
					}
					d.FieldBitBufLen("huffman_used_bitmap", int64(ranges)*16)
					d.FieldU3("huffman_groups")
					selectorsUsed := d.FieldU15("selectors_used")
					selectorsI := uint64(0)
					d.FieldArrayLoopFn("selector_list", func() bool { return selectorsI < selectorsUsed }, func(d *decode.D) {
						d.FieldU1("selector")
						selectorsI++
					})
				})

				d.FieldU48("footer_magic")
				d.FieldU32("crc")
			})

			moreStreams = false

		}
	})

	// .huffman_used_map:16            = bitmap, of ranges of 16 bytes, present/not present
	// .huffman_used_bitmaps:0..256    = bitmap, of symbols used, present/not present (multiples of 16)
	// .huffman_groups:3               = 2..6 number of different Huffman tables in use
	// .selectors_used:15              = number of times that the Huffman tables are swapped (each 50 symbols)
	// *.selector_list:1..6            = zero-terminated bit runs (0..62) of MTF'ed Huffman table (*selectors_used)
	// .start_huffman_length:5         = 0..20 starting bit length for Huffman deltas
	// *.delta_bit_length:1..40        = 0=>next symbol; 1=>alter length
	// 												{ 1=>decrement length;  0=>increment length } (*(symbols+2)*groups)
	// .contents:2..∞                  = Huffman encoded data stream until end of block (max. 7372800 bit)

	// .eos_magic:48                   = 0x177245385090 (BCD sqrt(pi))
	// .crc:32                         = checksum for whole stream
	// .padding:0..7                   = align to whole byte

	return nil
}
