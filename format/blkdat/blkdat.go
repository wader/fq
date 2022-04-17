package bson

// https://learnmeabitcoin.com/technical/blkdat

import (
	"github.com/wader/fq/format"
	"github.com/wader/fq/format/registry"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/scalar"
)

func init() {
	registry.MustRegister(decode.Format{
		Name:        format.BLKDAT,
		Description: "Bitcoin blk.dat",
		Groups:      []string{format.PROBE},
		DecodeFn:    decodeBlkDat,
	})
}

// <= 0xfc	              12
// <= 0xffff	          fd1234              Prefix with fd, and the next 2 bytes is the VarInt (in little-endian).
// <= 0xffffffff	      fe12345678          Prefix with fe, and the next 4 bytes is the VarInt (in little-endian).
// <= 0xffffffffffffffff  ff1234567890abcdef  Prefix with ff, and the next 8 bytes is the VarInt (in little-endian).
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

func decodeBlkDat(d *decode.D, in interface{}) interface{} {
	d.FieldU32("magic", scalar.UToSymStr{
		0xf9beb4d9: "mainnet",
		0x0b110907: "testnet3",
		0xfabfb5da: "regtest",
	}, scalar.Hex, d.AssertU(0xf9beb4d9, 0x0b110907, 0xfabfb5da))

	d.Endian = decode.LittleEndian

	size := d.FieldU32("size")

	d.FramedFn(int64(size)*8, func(d *decode.D) {
		d.FieldStruct("block_header", func(d *decode.D) {
			d.FieldU32("version")
			d.FieldRawLen("previous_block_hash", 32*8)
			d.FieldRawLen("merkle_root", 32*8)
			d.FieldU32("time")
			d.FieldU32("bits")
			d.FieldU32("nonce")

		})

		txCount := d.FieldUFn("tx_count", decodeVarInt)

		d.FieldArray("transactions", func(d *decode.D) {
			for i := uint64(0); i < txCount; i++ {
				d.FieldStruct("transaction", func(d *decode.D) {
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
								d.FieldRawLen("txid", 32*8)
								d.FieldU32("vout")
								scriptSigSize := d.FieldUFn("scriptsig_size", decodeVarInt)
								d.FieldRawLen("scriptsig", int64(scriptSigSize)*8)
								d.FieldU32("sequence")
							})
						}
					})
					outputCount := d.FieldUFn("output_count", decodeVarInt)
					d.FieldArray("outputs", func(d *decode.D) {
						for i := uint64(0); i < outputCount; i++ {
							d.FieldStruct("output", func(d *decode.D) {
								d.FieldRawLen("value", 8*8)
								scriptSigSize := d.FieldUFn("scriptpub_size", decodeVarInt)
								d.FieldRawLen("scriptpub", int64(scriptSigSize)*8)
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
				})

			}
		})

	})

	// d.FieldSr("block_header", int64(size)*8)

	return nil
}
