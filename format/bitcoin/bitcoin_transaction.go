package bitcoin

// TODO: coinbase transaction

// https://learnmeabitcoin.com/technical/blkdat

import (
	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
	"github.com/wader/fq/pkg/scalar"
)

var bitcoinScriptGroup decode.Group

func init() {
	interp.RegisterFormat(
		format.Bitcoin_Transaction,
		&decode.Format{
			Description: "Bitcoin transaction",
			Dependencies: []decode.Dependency{
				{Groups: []*decode.Group{format.Bitcoin_Script}, Out: &bitcoinScriptGroup},
			},
			DecodeFn: decodeBitcoinTranscation,
		})
}

// Prefix with fd, and the next 2 bytes is the VarInt (in little-endian).
// Prefix with fe, and the next 4 bytes is the VarInt (in little-endian).
// Prefix with ff, and the next 8 bytes is the VarInt (in little-endian).
func decodeVarInt(d *decode.D) uint64 {
	n := d.U8()
	switch n {
	case 0xfd:
		return d.U16()
	case 0xfe:
		return d.U32()
	case 0xff:
		return d.U64()
	default:
		return n
	}
}

// all zero
var txIDCoinbaseBytes = [32]byte{}

func decodeBitcoinTranscation(d *decode.D) any {
	d.Endian = decode.LittleEndian

	d.FieldU32("version")
	witness := false
	if d.PeekUintBits(8) == 0 {
		witness = true
		d.FieldU8("marker")
		d.FieldU8("flag")
	}
	inputCount := d.FieldUintFn("input_count", decodeVarInt)
	d.FieldArray("inputs", func(d *decode.D) {
		for i := uint64(0); i < inputCount; i++ {
			d.FieldStruct("input", func(d *decode.D) {
				d.FieldRawLen("txid", 32*8, scalar.RawBytesMap{
					{Bytes: txIDCoinbaseBytes[:], Scalar: scalar.BitBuf{Description: "coinbase"}},
				}, rawHexReverse)
				d.FieldU32("vout")
				scriptSigSize := d.FieldUintFn("scriptsig_size", decodeVarInt)
				d.FieldFormatOrRawLen("scriptsig", int64(scriptSigSize)*8, &bitcoinScriptGroup, nil)
				// TODO: better way to know if there should be a valid script
				d.FieldU32("sequence", scalar.UintHex)
			})
		}
	})
	outputCount := d.FieldUintFn("output_count", decodeVarInt)
	d.FieldArray("outputs", func(d *decode.D) {
		for i := uint64(0); i < outputCount; i++ {
			d.FieldStruct("output", func(d *decode.D) {
				d.FieldU64("value")
				scriptSigSize := d.FieldUintFn("scriptpub_size", decodeVarInt)
				// TODO: better way to know if there should be a valid script
				d.FieldFormatOrRawLen("scriptpub", int64(scriptSigSize)*8, &bitcoinScriptGroup, nil)
			})
		}
	})

	if witness {
		d.FieldArray("witnesses", func(d *decode.D) {
			for i := uint64(0); i < inputCount; i++ {
				d.FieldStruct("witness", func(d *decode.D) {
					witnessSize := d.FieldUintFn("witness_size", decodeVarInt)
					d.FieldArray("items", func(d *decode.D) {
						for j := uint64(0); j < witnessSize; j++ {
							d.FieldStruct("item", func(d *decode.D) {
								itemSize := d.FieldUintFn("item_size", decodeVarInt)
								d.FieldRawLen("item", int64(itemSize)*8)
							})
						}
					})
				})
			}
		})
	}
	d.FieldU32("locktime")

	return nil
}
