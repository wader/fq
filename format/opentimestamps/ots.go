package opentimestamps

import (
	"encoding/binary"
	"errors"
	"io"

	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
	"github.com/wader/fq/pkg/scalar"
)

func init() {
	interp.RegisterFormat(
		format.OpenTimestamps,
		&decode.Format{
			Description:  "OpenTimestamps file",
			Dependencies: nil,
			DecodeFn:     decodeOTSFile,
		})
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

var (
	headerMagic  = []byte{0x00, 0x4f, 0x70, 0x65, 0x6e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x73, 0x00, 0x00, 0x50, 0x72, 0x6f, 0x6f, 0x66, 0x00, 0xbf, 0x89, 0xe2, 0xe8, 0x84, 0xe8, 0x92, 0x94}
	pendingMagic = binary.BigEndian.Uint64([]byte{0x83, 0xdf, 0xe3, 0x0d, 0x2e, 0xf9, 0x0c, 0x8e})
	bitcoinMagic = binary.BigEndian.Uint64([]byte{0x05, 0x88, 0x96, 0x0d, 0x73, 0xd7, 0x19, 0x01})
)

var attestationMapper = scalar.UintMapSymStr{
	pendingMagic: "calendar",
	bitcoinMagic: "bitcoin",
}

var opMapper = scalar.UintMapSymStr{
	0xf0: "append",
	0xf1: "prepend",
	0xf2: "reverse",
	0xf3: "hexlify",
	0x02: "sha1",
	0x03: "ripemd160",
	0x08: "sha256",
	0x67: "keccak256",
	0x00: "attestation",
	0xff: "continuation_byte",
}

func decodeOTSFile(d *decode.D) any {
	d.Endian = decode.BigEndian

	d.FieldRawLen("magic_bytes", int64(8*len(headerMagic)), d.AssertBitBuf(headerMagic))
	d.FieldUintFn("version", decodeVarInt)
	tag := d.FieldU8("digest_hash_algorithm", opMapper,
		scalar.UintDescription("algorithm used to hash the source file"))
	if tag != 8 {
		name := opMapper[tag]
		d.Errorf("only sha256 supported, got %x: %s", tag, name)
		return nil
	}
	d.FieldRawLen("digest", 8*32, scalar.RawHex,
		scalar.BitBufDescription("hash of the source file"))

	d.FieldArray("instructions", func(d *decode.D) {
		for {
			if b, err := d.TryPeekBytes(1); errors.Is(err, io.EOF) {
				break
			} else if b[0] == 0x00 {
				d.FieldStruct("attestation", func(d *decode.D) {
					d.FieldU8("attestation_tag", scalar.UintMapSymBool{0x00: true})
					val := d.FieldU64BE("attestation_type", attestationMapper)
					d.FieldUintFn("attestation_varbytes_size", decodeVarInt)
					switch val {
					case bitcoinMagic:
						d.FieldUintFn("block", decodeVarInt,
							scalar.UintDescription("bitcoin block height to check for the merkle root"))
					case pendingMagic:
						nurl := d.FieldUintFn("url_size", decodeVarInt)
						d.FieldUTF8("url", int(nurl),
							scalar.StrDescription("url of the calendar server to get the final proof"))
					default:
						d.Errorf("unknown attestation tag %x", val)
					}
				})
			} else if b[0] == 0xff {
				d.FieldStruct("continuation_byte", func(d *decode.D) {
					d.FieldU8("continuation_byte", scalar.UintMapSymBool{0xff: true},
						scalar.UintDescription("tells we should continue reading after the next attestation block"))
				})
			} else {
				d.FieldStruct("instruction", func(d *decode.D) {
					tag := d.FieldU8("op", opMapper)
					if name, ok := opMapper[tag]; ok {
						// read var bytes if argument
						if name == "append" || name == "prepend" {
							n := d.FieldUintFn("argument_size", decodeVarInt)
							d.FieldRawLen("argument", int64(8*n), scalar.RawHex)
						}
					} else {
						d.Errorf("unknown operation tag %x", tag)
					}
				})
			}
		}
	})

	return nil
}
