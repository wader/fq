package opentimestamps

// https://opentimestamps.org/

import (
	"embed"
	"slices"

	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
	"github.com/wader/fq/pkg/scalar"
)

//go:embed opentimestamps.md
var otsFS embed.FS

func init() {
	interp.RegisterFormat(
		format.Opentimestamps,
		&decode.Format{
			Description: "OpenTimestamps file",
			DecodeFn:    decodeOTSFile,
			Groups:      []*decode.Group{format.Probe},
		})
	interp.RegisterFS(otsFS)
}

const (
	continuationByte = 0xff
	attestationTag   = 0x00
	appendTag        = 0xf0
	prependTag       = 0xf1
	reverseTag       = 0xf2
	hexlifyTag       = 0xf3
	sha1Tag          = 0x02
	ripemd160Tag     = 0x03
	sha256Tag        = 0x08
	keccak256Tag     = 0x67
)

var (
	headerMagic          = []byte{0x00, 0x4f, 0x70, 0x65, 0x6e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x73, 0x00, 0x00, 0x50, 0x72, 0x6f, 0x6f, 0x66, 0x00, 0xbf, 0x89, 0xe2, 0xe8, 0x84, 0xe8, 0x92, 0x94}
	calendarMagic uint64 = 0x83_df_e3_0d_2e_f9_0c_8e
	bitcoinMagic  uint64 = 0x05_88_96_0d_73_d7_19_01

	binaryTags = []byte{appendTag, prependTag}
)

var attestationMapper = scalar.UintMapSymStr{
	calendarMagic: "calendar",
	bitcoinMagic:  "bitcoin",
}

var tagMapper = scalar.UintMapSymStr{
	continuationByte: "continuation_byte",
	attestationTag:   "attestation",
	appendTag:        "append",
	prependTag:       "prepend",
	reverseTag:       "reverse",
	hexlifyTag:       "hexlify",
	sha1Tag:          "sha1",
	ripemd160Tag:     "ripemd160",
	sha256Tag:        "sha256",
	keccak256Tag:     "keccak256",
}

var digestSizes = map[byte]int64{
	sha1Tag:      20,
	ripemd160Tag: 20,
	sha256Tag:    32,
	keccak256Tag: 32,
}

func decodeOTSFile(d *decode.D) any {
	d.Endian = decode.BigEndian

	d.FieldRawLen("magic_bytes", int64(8*len(headerMagic)), d.AssertBitBuf(headerMagic))
	d.FieldUintFn("version", decodeVarInt)

	tag := d.FieldU8("digest_hash_algorithm", tagMapper)
	digestSize, ok := digestSizes[byte(tag)]
	if !ok {
		name := tagMapper[tag]
		d.Fatalf("hash algorithm not supported, got %x: '%s'", tag, name)
		return nil
	}
	d.FieldRawLen("digest", 8*digestSize, scalar.RawHex)

	d.FieldArray("operations", func(d *decode.D) {
		for d.NotEnd() {
			d.FieldStruct("operation", func(d *decode.D) {
				tag := d.FieldU8("tag", tagMapper)
				if tag == attestationTag {
					val := d.FieldU64BE("attestation_type", attestationMapper)
					n := d.FieldUintFn("attestation_varbytes_size", decodeVarInt)
					switch val {
					case bitcoinMagic:
						d.FieldUintFn("block", decodeVarInt)
					case calendarMagic:
						nurl := d.FieldUintFn("url_size", decodeVarInt)
						d.FieldUTF8("url", int(nurl))
					default:
						d.FieldRawLen("unknown_data", int64(n*8))
					}
				} else {
					if _, ok := tagMapper[tag]; ok {
						// read var bytes if argument
						if slices.Contains(binaryTags, byte(tag)) {
							n := d.FieldUintFn("argument_size", decodeVarInt)
							d.FieldRawLen("argument", int64(8*n), scalar.RawHex)
						}
					} else {
						d.Fatalf("unknown operation tag %x", tag)
					}
				}
			})
		}
	})

	return nil
}

func decodeVarInt(d *decode.D) uint64 {
	var value uint64 = 0
	var shift uint64 = 0

	for {
		b := d.U8()
		value |= (b & 0b01111111) << shift
		shift += 7
		if b&0b10000000 == 0 {
			break
		}
	}

	return value
}
