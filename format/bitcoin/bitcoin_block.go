package bitcoin

// https://learnmeabitcoin.com/technical/blkdat

import (
	"fmt"
	"time"

	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
	"github.com/wader/fq/pkg/scalar"
)

var bitcoinTranscationGroup decode.Group

func init() {
	interp.RegisterFormat(
		format.Bitcoin_Block,
		&decode.Format{
			Description: "Bitcoin block",
			Dependencies: []decode.Dependency{
				{Groups: []*decode.Group{format.Bitcoin_Transaction}, Out: &bitcoinTranscationGroup},
			},
			DecodeFn: decodeBitcoinBlock,
			DefaultInArg: format.Bitcoin_Block_In{
				HasHeader: false,
			},
		})
}

var rawHexReverse = scalar.BitBufFn(func(s scalar.BitBuf) (scalar.BitBuf, error) {
	return scalar.RawSym(s, -1, func(b []byte) string {
		decode.ReverseBytes(b)
		return fmt.Sprintf("%x", b)
	})
})

func decodeBitcoinBlock(d *decode.D) any {
	var bbi format.Bitcoin_Block_In
	d.ArgAs(&bbi)

	size := d.BitsLeft()

	if bbi.HasHeader {
		magic := d.PeekUintBits(32)
		switch magic {
		case 0xf9beb4d9,
			0x0b110907,
			0xfabfb5da:
			d.FieldU32("magic", scalar.UintMapSymStr{
				0xf9beb4d9: "mainnet",
				0x0b110907: "testnet3",
				0xfabfb5da: "regtest",
			}, scalar.UintHex)
			size = int64(d.FieldU32LE("size")) * 8
		default:
			d.Fatalf("unknown magic %x", magic)
		}
	}

	d.Endian = decode.LittleEndian

	d.FramedFn(size, func(d *decode.D) {
		d.FieldStruct("header", func(d *decode.D) {
			d.FieldU32("version", scalar.UintHex)
			d.FieldRawLen("previous_block_hash", 32*8, rawHexReverse)
			d.FieldRawLen("merkle_root", 32*8, rawHexReverse)
			d.FieldU32("time", scalar.UintActualUnixTimeDescription(time.Second, time.RFC3339))
			d.FieldU32("bits", scalar.UintHex)
			d.FieldU32("nonce", scalar.UintHex)
		})

		// TODO: remove? support header only decode this way?
		if d.BitsLeft() == 0 {
			return
		}

		txCount := d.FieldUintFn("tx_count", decodeVarInt)
		d.FieldArray("transactions", func(d *decode.D) {
			for i := uint64(0); i < txCount; i++ {
				d.FieldFormat("transaction", &bitcoinTranscationGroup, nil)
			}
		})
	})

	return nil
}
