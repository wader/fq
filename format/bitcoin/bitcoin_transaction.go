package bitcoin

// TODO: coinbase transaction

// https://learnmeabitcoin.com/technical/blkdat

import (
	"github.com/wader/fq/format"
	"github.com/wader/fq/format/registry"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/scalar"
)

var bitcoinScriptFormat decode.Group

func init() {
	registry.MustRegister(decode.Format{
		Name:        format.BITCOIN_TRANSACTION,
		Description: "Bitcoin transaction",
		Dependencies: []decode.Dependency{
			{Names: []string{format.BITCOIN_SCRIPT}, Group: &bitcoinScriptFormat},
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

func decodeBitcoinTranscation(d *decode.D, in interface{}) interface{} {
	d.Endian = decode.LittleEndian

	d.FieldU32("version")
	witness := false
	if d.PeekBits(8) == 0 {
		witness = true
		d.FieldU8("marker")
		d.FieldU8("flag")
	}
	inputCount := d.FieldUFn("input_count", decodeVarInt)
	d.FieldArray("inputs", func(d *decode.D) {
		for i := uint64(0); i < inputCount; i++ {
			d.FieldStruct("input", func(d *decode.D) {
				d.FieldRawLen("txid", 32*8, scalar.BytesToScalar{
					{Bytes: txIDCoinbaseBytes[:], Scalar: scalar.S{Description: "coinbase"}},
				}, rawHexReverse)
				d.FieldU32("vout")
				scriptSigSize := d.FieldUFn("scriptsig_size", decodeVarInt)
				d.FieldFormatOrRawLen("scriptsig", int64(scriptSigSize)*8, bitcoinScriptFormat, nil)
				// TODO: better way to know if there should be a valid script
				d.FieldU32("sequence", scalar.ActualHex)
			})
		}
	})
	outputCount := d.FieldUFn("output_count", decodeVarInt)
	d.FieldArray("outputs", func(d *decode.D) {
		for i := uint64(0); i < outputCount; i++ {
			d.FieldStruct("output", func(d *decode.D) {
				d.FieldU64("value")
				scriptSigSize := d.FieldUFn("scriptpub_size", decodeVarInt)
				// TODO: better way to know if there should be a valid script
				d.FieldFormatOrRawLen("scriptpub", int64(scriptSigSize)*8, bitcoinScriptFormat, nil)
			})
		}
	})

	if witness {
		d.FieldArray("witnesses", func(d *decode.D) {
			for i := uint64(0); i < inputCount; i++ {
				d.FieldStruct("witness", func(d *decode.D) {
					witnessSize := d.FieldUFn("witness_size", decodeVarInt)
					d.FieldArray("items", func(d *decode.D) {
						for j := uint64(0); j < witnessSize; j++ {
							itemSize := d.FieldUFn("item_size", decodeVarInt)
							d.FieldRawLen("item", int64(itemSize)*8)
						}
					})
				})
			}
		})
	}
	d.FieldU32("locktime")

	return nil
}
