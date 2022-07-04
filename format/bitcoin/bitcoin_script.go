package bitcoin

import (
	"github.com/wader/fq/format"
	"github.com/wader/fq/format/registry"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/scalar"
)

type opcodeEntry struct {
	r [2]byte
	s scalar.S
	d func(d *decode.D, opcode byte)
}

type opcodeEntries []opcodeEntry

func (ops opcodeEntries) lookup(u byte) (opcodeEntry, bool) {
	for _, fe := range ops {
		if u >= fe.r[0] && u <= fe.r[1] {
			return fe, true
		}
	}
	return opcodeEntry{}, false
}

func (ops opcodeEntries) MapScalar(s scalar.S) (scalar.S, error) {
	u := s.ActualU()
	if fe, ok := ops.lookup(byte(u)); ok {
		s = fe.s
		s.Actual = u
	}
	return s, nil
}

func init() {
	registry.MustRegister(decode.Format{
		Name:        format.BITCOIN_SCRIPT,
		Description: "Bitcoin script",
		DecodeFn:    decodeBitcoinScript,
		RootArray:   true,
		RootName:    "opcodes",
	})
}

func decodeBitcoinScript(d *decode.D, in interface{}) interface{} {
	// based on https://en.bitcoin.it/wiki/Script
	opcodeEntries := opcodeEntries{
		{r: [2]byte{0x00, 0x00}, s: scalar.S{Sym: "false"}},
		// TODO: name op code?
		{r: [2]byte{0x01, 0x4b}, s: scalar.S{Sym: "pushself"}, d: func(d *decode.D, opcode byte) {
			d.FieldRawLen("arg", int64(opcode)*8)
		}},
		{r: [2]byte{0x04c, 0x4e}, s: scalar.S{Sym: "pushdata1"}, d: func(d *decode.D, opcode byte) {
			argLen := d.FieldU8("arg_length")
			d.FieldRawLen("arg", int64(argLen)*8)
		}},
		{r: [2]byte{0x04c, 0x4e}, s: scalar.S{Sym: "pushdata2"}, d: func(d *decode.D, opcode byte) {
			argLen := d.FieldU16("arg_length")
			d.FieldRawLen("arg", int64(argLen)*8)
		}},
		{r: [2]byte{0x04c, 0x4e}, s: scalar.S{Sym: "pushdata4"}, d: func(d *decode.D, opcode byte) {
			argLen := d.FieldU32("arg_length")
			d.FieldRawLen("arg", int64(argLen)*8)
		}},
		{r: [2]byte{0x4f, 0x4f}, s: scalar.S{Sym: "1negate"}},
		{r: [2]byte{0x51, 0x51}, s: scalar.S{Sym: "true"}},
		// TODO: name
		{r: [2]byte{0x52, 0x60}, s: scalar.S{Sym: "push"}, d: func(d *decode.D, opcode byte) {
			d.FieldValueU("arg", uint64(opcode-0x50))
		}},
		{r: [2]byte{0x61, 0x61}, s: scalar.S{Sym: "nop"}},
		{r: [2]byte{0x62, 0x62}, s: scalar.S{Sym: "ver"}},
		{r: [2]byte{0x63, 0x63}, s: scalar.S{Sym: "if"}},
		{r: [2]byte{0x64, 0x64}, s: scalar.S{Sym: "notif"}},
		{r: [2]byte{0x65, 0x65}, s: scalar.S{Sym: "verif"}},
		{r: [2]byte{0x66, 0x66}, s: scalar.S{Sym: "vernotif"}},
		{r: [2]byte{0x67, 0x67}, s: scalar.S{Sym: "else"}},
		{r: [2]byte{0x68, 0x68}, s: scalar.S{Sym: "endif"}},
		{r: [2]byte{0x69, 0x69}, s: scalar.S{Sym: "verify"}},
		{r: [2]byte{0x6a, 0x6a}, s: scalar.S{Sym: "return"}},
		{r: [2]byte{0x6b, 0x6b}, s: scalar.S{Sym: "toaltstack"}},
		{r: [2]byte{0x6c, 0x6c}, s: scalar.S{Sym: "fromaltstack"}},
		{r: [2]byte{0x6d, 0x6d}, s: scalar.S{Sym: "2drop"}},
		{r: [2]byte{0x6e, 0x6e}, s: scalar.S{Sym: "2dup"}},
		{r: [2]byte{0x6f, 0x6f}, s: scalar.S{Sym: "3dup"}},
		{r: [2]byte{0x70, 0x70}, s: scalar.S{Sym: "2over"}},
		{r: [2]byte{0x71, 0x71}, s: scalar.S{Sym: "2rot"}},
		{r: [2]byte{0x72, 0x72}, s: scalar.S{Sym: "2swap"}},
		{r: [2]byte{0x73, 0x73}, s: scalar.S{Sym: "ifdup"}},
		{r: [2]byte{0x74, 0x74}, s: scalar.S{Sym: "depth"}},
		{r: [2]byte{0x75, 0x75}, s: scalar.S{Sym: "drop"}},
		{r: [2]byte{0x76, 0x76}, s: scalar.S{Sym: "dup"}},
		{r: [2]byte{0x77, 0x77}, s: scalar.S{Sym: "nip"}},
		{r: [2]byte{0x78, 0x78}, s: scalar.S{Sym: "over"}},
		{r: [2]byte{0x79, 0x79}, s: scalar.S{Sym: "pick"}},
		{r: [2]byte{0x7a, 0x7a}, s: scalar.S{Sym: "roll"}},
		{r: [2]byte{0x7b, 0x7b}, s: scalar.S{Sym: "rot"}},
		{r: [2]byte{0x7c, 0x7c}, s: scalar.S{Sym: "swap"}},
		{r: [2]byte{0x7d, 0x7d}, s: scalar.S{Sym: "tuck"}},
		{r: [2]byte{0x7e, 0x7e}, s: scalar.S{Sym: "cat"}},
		{r: [2]byte{0x7f, 0x7f}, s: scalar.S{Sym: "split"}},
		{r: [2]byte{0x80, 0x80}, s: scalar.S{Sym: "num2bin"}},
		{r: [2]byte{0x81, 0x81}, s: scalar.S{Sym: "bin2num"}},
		{r: [2]byte{0x82, 0x82}, s: scalar.S{Sym: "size"}},
		{r: [2]byte{0x83, 0x83}, s: scalar.S{Sym: "invert"}},
		{r: [2]byte{0x84, 0x84}, s: scalar.S{Sym: "and"}},
		{r: [2]byte{0x85, 0x85}, s: scalar.S{Sym: "or"}},
		{r: [2]byte{0x86, 0x86}, s: scalar.S{Sym: "xor"}},
		{r: [2]byte{0x87, 0x87}, s: scalar.S{Sym: "equal"}},
		{r: [2]byte{0x88, 0x88}, s: scalar.S{Sym: "equalverify"}},
		{r: [2]byte{0x89, 0x89}, s: scalar.S{Sym: "reserved1"}},
		{r: [2]byte{0x8a, 0x8a}, s: scalar.S{Sym: "reserved2"}},
		{r: [2]byte{0x8b, 0x8b}, s: scalar.S{Sym: "1add"}},
		{r: [2]byte{0x8c, 0x8c}, s: scalar.S{Sym: "1sub"}},
		{r: [2]byte{0x8d, 0x8d}, s: scalar.S{Sym: "2mul"}},
		{r: [2]byte{0x8e, 0x8e}, s: scalar.S{Sym: "2div"}},
		{r: [2]byte{0x8f, 0x8f}, s: scalar.S{Sym: "negate"}},
		{r: [2]byte{0x90, 0x90}, s: scalar.S{Sym: "abs"}},
		{r: [2]byte{0x91, 0x91}, s: scalar.S{Sym: "not"}},
		{r: [2]byte{0x92, 0x92}, s: scalar.S{Sym: "0notequal"}},
		{r: [2]byte{0x93, 0x93}, s: scalar.S{Sym: "add"}},
		{r: [2]byte{0x94, 0x94}, s: scalar.S{Sym: "sub"}},
		{r: [2]byte{0x95, 0x95}, s: scalar.S{Sym: "mul"}},
		{r: [2]byte{0x96, 0x96}, s: scalar.S{Sym: "div"}},
		{r: [2]byte{0x97, 0x97}, s: scalar.S{Sym: "mod"}},
		{r: [2]byte{0x98, 0x98}, s: scalar.S{Sym: "lshift"}},
		{r: [2]byte{0x99, 0x99}, s: scalar.S{Sym: "rshift"}},
		{r: [2]byte{0x9a, 0x9a}, s: scalar.S{Sym: "booland"}},
		{r: [2]byte{0x9b, 0x9b}, s: scalar.S{Sym: "boolor"}},
		{r: [2]byte{0x9c, 0x9c}, s: scalar.S{Sym: "numequal"}},
		{r: [2]byte{0x9d, 0x9d}, s: scalar.S{Sym: "numequalverify"}},
		{r: [2]byte{0x9e, 0x9e}, s: scalar.S{Sym: "numnotequal"}},
		{r: [2]byte{0x9f, 0x9f}, s: scalar.S{Sym: "lessthan"}},
		{r: [2]byte{0xa0, 0xa0}, s: scalar.S{Sym: "greaterthan"}},
		{r: [2]byte{0xa1, 0xa1}, s: scalar.S{Sym: "lessthanorequal"}},
		{r: [2]byte{0xa2, 0xa2}, s: scalar.S{Sym: "greaterthanorequal"}},
		{r: [2]byte{0xa3, 0xa3}, s: scalar.S{Sym: "min"}},
		{r: [2]byte{0xa4, 0xa4}, s: scalar.S{Sym: "max"}},
		{r: [2]byte{0xa5, 0xa5}, s: scalar.S{Sym: "within"}},
		{r: [2]byte{0xa6, 0xa6}, s: scalar.S{Sym: "ripemd160"}},
		{r: [2]byte{0xa7, 0xa7}, s: scalar.S{Sym: "sha1"}},
		{r: [2]byte{0xa8, 0xa8}, s: scalar.S{Sym: "sha256"}},
		{r: [2]byte{0xa9, 0xa9}, s: scalar.S{Sym: "hash160"}},
		{r: [2]byte{0xaa, 0xaa}, s: scalar.S{Sym: "hash256"}},
		{r: [2]byte{0xab, 0xab}, s: scalar.S{Sym: "codeseparator"}},
		{r: [2]byte{0xac, 0xac}, s: scalar.S{Sym: "checksig"}},
		{r: [2]byte{0xad, 0xad}, s: scalar.S{Sym: "checksigverify"}},
		{r: [2]byte{0xae, 0xae}, s: scalar.S{Sym: "checkmultisig"}},
		{r: [2]byte{0xaf, 0xaf}, s: scalar.S{Sym: "checkmultisigverify"}},
		{r: [2]byte{0xb0, 0xb0}, s: scalar.S{Sym: "nop1"}},
		{r: [2]byte{0xb1, 0xb1}, s: scalar.S{Sym: "nop2"}},
		{r: [2]byte{0xb1, 0xb1}, s: scalar.S{Sym: "checklocktimeverify"}},
		{r: [2]byte{0xb2, 0xb2}, s: scalar.S{Sym: "nop3"}},
		{r: [2]byte{0xb2, 0xb2}, s: scalar.S{Sym: "checksequenceverify"}},
		{r: [2]byte{0xb3, 0xb3}, s: scalar.S{Sym: "nop4"}},
		{r: [2]byte{0xb4, 0xb4}, s: scalar.S{Sym: "nop5"}},
		{r: [2]byte{0xb5, 0xb5}, s: scalar.S{Sym: "nop6"}},
		{r: [2]byte{0xb6, 0xb6}, s: scalar.S{Sym: "nop7"}},
		{r: [2]byte{0xb7, 0xb7}, s: scalar.S{Sym: "nop8"}},
		{r: [2]byte{0xb8, 0xb8}, s: scalar.S{Sym: "nop9"}},
		{r: [2]byte{0xb9, 0xb9}, s: scalar.S{Sym: "nop10"}},
		{r: [2]byte{0xba, 0xba}, s: scalar.S{Sym: "checkdatasig"}},
		{r: [2]byte{0xbb, 0xbb}, s: scalar.S{Sym: "checkdatasigverif"}},
		{r: [2]byte{0xfa, 0xfa}, s: scalar.S{Sym: "smallinteger"}},
		{r: [2]byte{0xfb, 0xfb}, s: scalar.S{Sym: "pubkeys"}},
		{r: [2]byte{0xfc, 0xfc}, s: scalar.S{Sym: "unknown252"}},
		{r: [2]byte{0xfd, 0xfd}, s: scalar.S{Sym: "pubkeyhash"}},
		{r: [2]byte{0xfe, 0xfe}, s: scalar.S{Sym: "pubkey"}},
		{r: [2]byte{0xff, 0xff}, s: scalar.S{Sym: "invalidopcode"}},
	}

	for !d.End() {
		opcode := byte(d.PeekBits(8))
		ope, ok := opcodeEntries.lookup(opcode)
		if !ok {
			d.Fatalf("unknown opcode %x", opcode)
		}

		d.FieldStruct("opcode", func(d *decode.D) {
			d.FieldU8("op", opcodeEntries)
			if ope.d != nil {
				ope.d(d, opcode)
			}
		})
	}

	return nil
}
