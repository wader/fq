package bitcoin

import (
	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
	"github.com/wader/fq/pkg/scalar"
)

type opcodeEntry struct {
	r [2]byte
	s scalar.Uint
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

func (ops opcodeEntries) MapUint(s scalar.Uint) (scalar.Uint, error) {
	u := s.Actual
	if fe, ok := ops.lookup(byte(u)); ok {
		s = fe.s
		s.Actual = u
	}
	return s, nil
}

func init() {
	interp.RegisterFormat(
		format.Bitcoin_Script,
		&decode.Format{
			Description: "Bitcoin script",
			DecodeFn:    decodeBitcoinScript,
			RootArray:   true,
			RootName:    "opcodes",
		})
}

func decodeBitcoinScript(d *decode.D) any {
	// based on https://en.bitcoin.it/wiki/Script
	opcodeEntries := opcodeEntries{
		{r: [2]byte{0x00, 0x00}, s: scalar.Uint{Sym: "false"}},
		// TODO: name op code?
		{r: [2]byte{0x01, 0x4b}, s: scalar.Uint{Sym: "pushself"}, d: func(d *decode.D, opcode byte) {
			d.FieldRawLen("arg", int64(opcode)*8)
		}},
		{r: [2]byte{0x04c, 0x4e}, s: scalar.Uint{Sym: "pushdata1"}, d: func(d *decode.D, opcode byte) {
			argLen := d.FieldU8("arg_length")
			d.FieldRawLen("arg", int64(argLen)*8)
		}},
		{r: [2]byte{0x04c, 0x4e}, s: scalar.Uint{Sym: "pushdata2"}, d: func(d *decode.D, opcode byte) {
			argLen := d.FieldU16("arg_length")
			d.FieldRawLen("arg", int64(argLen)*8)
		}},
		{r: [2]byte{0x04c, 0x4e}, s: scalar.Uint{Sym: "pushdata4"}, d: func(d *decode.D, opcode byte) {
			argLen := d.FieldU32("arg_length")
			d.FieldRawLen("arg", int64(argLen)*8)
		}},
		{r: [2]byte{0x4f, 0x4f}, s: scalar.Uint{Sym: "1negate"}},
		{r: [2]byte{0x51, 0x51}, s: scalar.Uint{Sym: "true"}},
		// TODO: name
		{r: [2]byte{0x52, 0x60}, s: scalar.Uint{Sym: "push"}, d: func(d *decode.D, opcode byte) {
			d.FieldValueUint("arg", uint64(opcode-0x50))
		}},
		{r: [2]byte{0x61, 0x61}, s: scalar.Uint{Sym: "nop"}},
		{r: [2]byte{0x62, 0x62}, s: scalar.Uint{Sym: "ver"}},
		{r: [2]byte{0x63, 0x63}, s: scalar.Uint{Sym: "if"}},
		{r: [2]byte{0x64, 0x64}, s: scalar.Uint{Sym: "notif"}},
		{r: [2]byte{0x65, 0x65}, s: scalar.Uint{Sym: "verif"}},
		{r: [2]byte{0x66, 0x66}, s: scalar.Uint{Sym: "vernotif"}},
		{r: [2]byte{0x67, 0x67}, s: scalar.Uint{Sym: "else"}},
		{r: [2]byte{0x68, 0x68}, s: scalar.Uint{Sym: "endif"}},
		{r: [2]byte{0x69, 0x69}, s: scalar.Uint{Sym: "verify"}},
		{r: [2]byte{0x6a, 0x6a}, s: scalar.Uint{Sym: "return"}},
		{r: [2]byte{0x6b, 0x6b}, s: scalar.Uint{Sym: "toaltstack"}},
		{r: [2]byte{0x6c, 0x6c}, s: scalar.Uint{Sym: "fromaltstack"}},
		{r: [2]byte{0x6d, 0x6d}, s: scalar.Uint{Sym: "2drop"}},
		{r: [2]byte{0x6e, 0x6e}, s: scalar.Uint{Sym: "2dup"}},
		{r: [2]byte{0x6f, 0x6f}, s: scalar.Uint{Sym: "3dup"}},
		{r: [2]byte{0x70, 0x70}, s: scalar.Uint{Sym: "2over"}},
		{r: [2]byte{0x71, 0x71}, s: scalar.Uint{Sym: "2rot"}},
		{r: [2]byte{0x72, 0x72}, s: scalar.Uint{Sym: "2swap"}},
		{r: [2]byte{0x73, 0x73}, s: scalar.Uint{Sym: "ifdup"}},
		{r: [2]byte{0x74, 0x74}, s: scalar.Uint{Sym: "depth"}},
		{r: [2]byte{0x75, 0x75}, s: scalar.Uint{Sym: "drop"}},
		{r: [2]byte{0x76, 0x76}, s: scalar.Uint{Sym: "dup"}},
		{r: [2]byte{0x77, 0x77}, s: scalar.Uint{Sym: "nip"}},
		{r: [2]byte{0x78, 0x78}, s: scalar.Uint{Sym: "over"}},
		{r: [2]byte{0x79, 0x79}, s: scalar.Uint{Sym: "pick"}},
		{r: [2]byte{0x7a, 0x7a}, s: scalar.Uint{Sym: "roll"}},
		{r: [2]byte{0x7b, 0x7b}, s: scalar.Uint{Sym: "rot"}},
		{r: [2]byte{0x7c, 0x7c}, s: scalar.Uint{Sym: "swap"}},
		{r: [2]byte{0x7d, 0x7d}, s: scalar.Uint{Sym: "tuck"}},
		{r: [2]byte{0x7e, 0x7e}, s: scalar.Uint{Sym: "cat"}},
		{r: [2]byte{0x7f, 0x7f}, s: scalar.Uint{Sym: "split"}},
		{r: [2]byte{0x80, 0x80}, s: scalar.Uint{Sym: "num2bin"}},
		{r: [2]byte{0x81, 0x81}, s: scalar.Uint{Sym: "bin2num"}},
		{r: [2]byte{0x82, 0x82}, s: scalar.Uint{Sym: "size"}},
		{r: [2]byte{0x83, 0x83}, s: scalar.Uint{Sym: "invert"}},
		{r: [2]byte{0x84, 0x84}, s: scalar.Uint{Sym: "and"}},
		{r: [2]byte{0x85, 0x85}, s: scalar.Uint{Sym: "or"}},
		{r: [2]byte{0x86, 0x86}, s: scalar.Uint{Sym: "xor"}},
		{r: [2]byte{0x87, 0x87}, s: scalar.Uint{Sym: "equal"}},
		{r: [2]byte{0x88, 0x88}, s: scalar.Uint{Sym: "equalverify"}},
		{r: [2]byte{0x89, 0x89}, s: scalar.Uint{Sym: "reserved1"}},
		{r: [2]byte{0x8a, 0x8a}, s: scalar.Uint{Sym: "reserved2"}},
		{r: [2]byte{0x8b, 0x8b}, s: scalar.Uint{Sym: "1add"}},
		{r: [2]byte{0x8c, 0x8c}, s: scalar.Uint{Sym: "1sub"}},
		{r: [2]byte{0x8d, 0x8d}, s: scalar.Uint{Sym: "2mul"}},
		{r: [2]byte{0x8e, 0x8e}, s: scalar.Uint{Sym: "2div"}},
		{r: [2]byte{0x8f, 0x8f}, s: scalar.Uint{Sym: "negate"}},
		{r: [2]byte{0x90, 0x90}, s: scalar.Uint{Sym: "abs"}},
		{r: [2]byte{0x91, 0x91}, s: scalar.Uint{Sym: "not"}},
		{r: [2]byte{0x92, 0x92}, s: scalar.Uint{Sym: "0notequal"}},
		{r: [2]byte{0x93, 0x93}, s: scalar.Uint{Sym: "add"}},
		{r: [2]byte{0x94, 0x94}, s: scalar.Uint{Sym: "sub"}},
		{r: [2]byte{0x95, 0x95}, s: scalar.Uint{Sym: "mul"}},
		{r: [2]byte{0x96, 0x96}, s: scalar.Uint{Sym: "div"}},
		{r: [2]byte{0x97, 0x97}, s: scalar.Uint{Sym: "mod"}},
		{r: [2]byte{0x98, 0x98}, s: scalar.Uint{Sym: "lshift"}},
		{r: [2]byte{0x99, 0x99}, s: scalar.Uint{Sym: "rshift"}},
		{r: [2]byte{0x9a, 0x9a}, s: scalar.Uint{Sym: "booland"}},
		{r: [2]byte{0x9b, 0x9b}, s: scalar.Uint{Sym: "boolor"}},
		{r: [2]byte{0x9c, 0x9c}, s: scalar.Uint{Sym: "numequal"}},
		{r: [2]byte{0x9d, 0x9d}, s: scalar.Uint{Sym: "numequalverify"}},
		{r: [2]byte{0x9e, 0x9e}, s: scalar.Uint{Sym: "numnotequal"}},
		{r: [2]byte{0x9f, 0x9f}, s: scalar.Uint{Sym: "lessthan"}},
		{r: [2]byte{0xa0, 0xa0}, s: scalar.Uint{Sym: "greaterthan"}},
		{r: [2]byte{0xa1, 0xa1}, s: scalar.Uint{Sym: "lessthanorequal"}},
		{r: [2]byte{0xa2, 0xa2}, s: scalar.Uint{Sym: "greaterthanorequal"}},
		{r: [2]byte{0xa3, 0xa3}, s: scalar.Uint{Sym: "min"}},
		{r: [2]byte{0xa4, 0xa4}, s: scalar.Uint{Sym: "max"}},
		{r: [2]byte{0xa5, 0xa5}, s: scalar.Uint{Sym: "within"}},
		{r: [2]byte{0xa6, 0xa6}, s: scalar.Uint{Sym: "ripemd160"}},
		{r: [2]byte{0xa7, 0xa7}, s: scalar.Uint{Sym: "sha1"}},
		{r: [2]byte{0xa8, 0xa8}, s: scalar.Uint{Sym: "sha256"}},
		{r: [2]byte{0xa9, 0xa9}, s: scalar.Uint{Sym: "hash160"}},
		{r: [2]byte{0xaa, 0xaa}, s: scalar.Uint{Sym: "hash256"}},
		{r: [2]byte{0xab, 0xab}, s: scalar.Uint{Sym: "codeseparator"}},
		{r: [2]byte{0xac, 0xac}, s: scalar.Uint{Sym: "checksig"}},
		{r: [2]byte{0xad, 0xad}, s: scalar.Uint{Sym: "checksigverify"}},
		{r: [2]byte{0xae, 0xae}, s: scalar.Uint{Sym: "checkmultisig"}},
		{r: [2]byte{0xaf, 0xaf}, s: scalar.Uint{Sym: "checkmultisigverify"}},
		{r: [2]byte{0xb0, 0xb0}, s: scalar.Uint{Sym: "nop1"}},
		{r: [2]byte{0xb1, 0xb1}, s: scalar.Uint{Sym: "nop2"}},
		{r: [2]byte{0xb1, 0xb1}, s: scalar.Uint{Sym: "checklocktimeverify"}},
		{r: [2]byte{0xb2, 0xb2}, s: scalar.Uint{Sym: "nop3"}},
		{r: [2]byte{0xb2, 0xb2}, s: scalar.Uint{Sym: "checksequenceverify"}},
		{r: [2]byte{0xb3, 0xb3}, s: scalar.Uint{Sym: "nop4"}},
		{r: [2]byte{0xb4, 0xb4}, s: scalar.Uint{Sym: "nop5"}},
		{r: [2]byte{0xb5, 0xb5}, s: scalar.Uint{Sym: "nop6"}},
		{r: [2]byte{0xb6, 0xb6}, s: scalar.Uint{Sym: "nop7"}},
		{r: [2]byte{0xb7, 0xb7}, s: scalar.Uint{Sym: "nop8"}},
		{r: [2]byte{0xb8, 0xb8}, s: scalar.Uint{Sym: "nop9"}},
		{r: [2]byte{0xb9, 0xb9}, s: scalar.Uint{Sym: "nop10"}},
		{r: [2]byte{0xba, 0xba}, s: scalar.Uint{Sym: "checkdatasig"}},
		{r: [2]byte{0xbb, 0xbb}, s: scalar.Uint{Sym: "checkdatasigverif"}},
		{r: [2]byte{0xfa, 0xfa}, s: scalar.Uint{Sym: "smallinteger"}},
		{r: [2]byte{0xfb, 0xfb}, s: scalar.Uint{Sym: "pubkeys"}},
		{r: [2]byte{0xfc, 0xfc}, s: scalar.Uint{Sym: "unknown252"}},
		{r: [2]byte{0xfd, 0xfd}, s: scalar.Uint{Sym: "pubkeyhash"}},
		{r: [2]byte{0xfe, 0xfe}, s: scalar.Uint{Sym: "pubkey"}},
		{r: [2]byte{0xff, 0xff}, s: scalar.Uint{Sym: "invalidopcode"}},
	}

	for !d.End() {
		opcode := byte(d.PeekUintBits(8))
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
