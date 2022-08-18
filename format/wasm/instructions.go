package wasm

import (
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/scalar"
)

func decodeBlockType(d *decode.D, name string) {
	b := d.PeekBytes(1)[0]
	switch b {
	case 0x40:
		d.FieldU8(name, scalar.Sym("ε"))
	case 0x6f, 0x70, 0x7b, 0x7c, 0x7d, 0x7e, 0x7f:
		d.FieldU8(name, valtypeToSymMapper)
	default:
		d.FieldSScalarFn(name, readSignedLEB128)
	}
}

func decodeMemArg(d *decode.D, name string) {
	d.FieldStruct(name, func(d *decode.D) {
		fieldU32(d, "a")
		fieldU32(d, "o")
	})
}

// expr ::= (in:instr)* 0x0B => in* end
func decodeExpr(d *decode.D, name string) {
	d.FieldArray(name, func(d *decode.D) {
		for {
			b := d.PeekBytes(1)[0]
			d.FieldStruct("in", decodeInstruction)
			if b == 0x0b {
				break
			}
		}
	})
}

func decodeInstruction(d *decode.D) {
	instr := d.PeekBytes(1)
	if len(instr) == 0 {
		return
	}

	i := instr[0]
	switch i {
	case 0x00:
		decodeInstructionWithoutOperands(d, "unreachable", i)
	case 0x01:
		decodeInstructionWithoutOperands(d, "nop", i)
	case 0x02:
		decodeBlock(d)
	case 0x03:
		decodeLoop(d)
	case 0x04:
		decodeIf(d)

	case 0x0b:
		decodeEnd(d)
	case 0x0c:
		decodeBr(d)
	case 0x0d:
		decodeBrIf(d)
	case 0x0e:
		decodeBrTable(d)
	case 0x0f:
		decodeInstructionWithoutOperands(d, "return", i)
	case 0x10:
		decodeCall(d)
	case 0x11:
		decodeCallIndirect(d)

	case 0x1a:
		decodeInstructionWithoutOperands(d, "drop", i)
	case 0x1b:
		decodeInstructionWithoutOperands(d, "select", i)
	case 0x1c:
		decodeSelectT(d)

	case 0x20:
		decodeLocalGet(d)
	case 0x21:
		decodeLocalSet(d)
	case 0x22:
		decodeLocalTee(d)
	case 0x23:
		decodeGlobalGet(d)
	case 0x24:
		decodeGlobalSet(d)

	case 0x25:
		decodeTableGet(d)
	case 0x26:
		decodeTableSet(d)

	case 0x28:
		decodeI32Load(d)
	case 0x29:
		decodeI64Load(d)
	case 0x2a:
		decodeF32Load(d)
	case 0x2b:
		decodeF64Load(d)
	case 0x2c:
		decodeI32Load8S(d)
	case 0x2d:
		decodeI32Load8U(d)
	case 0x2e:
		decodeI32Load16S(d)
	case 0x2f:
		decodeI32Load16U(d)
	case 0x30:
		decodeI64Load8S(d)
	case 0x31:
		decodeI64Load8U(d)
	case 0x32:
		decodeI64Load16S(d)
	case 0x33:
		decodeI64Load16U(d)
	case 0x34:
		decodeI64Load32S(d)
	case 0x35:
		decodeI64Load32U(d)
	case 0x36:
		decodeI32Store(d)
	case 0x37:
		decodeI64Store(d)
	case 0x38:
		decodeF32Store(d)
	case 0x39:
		decodeF64Store(d)
	case 0x3a:
		decodeI32Store8(d)
	case 0x3b:
		decodeI32Store16(d)
	case 0x3c:
		decodeI64Store8(d)
	case 0x3d:
		decodeI64Store16(d)
	case 0x3e:
		decodeI64Store32(d)

	case 0x3f:
		decodeMemorySize(d)
	case 0x40:
		decodeMemoryGrow(d)

	case 0x41:
		decodeI32Const(d)
	case 0x42:
		decodeI64Const(d)
	case 0x43:
		decodeF32Const(d)
	case 0x44:
		decodeF64Const(d)

	case 0x45:
		decodeInstructionWithoutOperands(d, "i32.eqz", i)
	case 0x46:
		decodeInstructionWithoutOperands(d, "i32.eq", i)
	case 0x47:
		decodeInstructionWithoutOperands(d, "i32.ne", i)
	case 0x48:
		decodeInstructionWithoutOperands(d, "i32.lt_s", i)
	case 0x49:
		decodeInstructionWithoutOperands(d, "i32.lt_u", i)
	case 0x4a:
		decodeInstructionWithoutOperands(d, "i32.gt_s", i)
	case 0x4b:
		decodeInstructionWithoutOperands(d, "i32.gt_u", i)
	case 0x4c:
		decodeInstructionWithoutOperands(d, "i32.le_s", i)
	case 0x4d:
		decodeInstructionWithoutOperands(d, "i32.le_u", i)
	case 0x4e:
		decodeInstructionWithoutOperands(d, "i32.ge_s", i)
	case 0x4f:
		decodeInstructionWithoutOperands(d, "i32.ge_u", i)

	case 0x50:
		decodeInstructionWithoutOperands(d, "i64.eqz", i)
	case 0x51:
		decodeInstructionWithoutOperands(d, "i64.eq", i)
	case 0x52:
		decodeInstructionWithoutOperands(d, "i64.ne", i)
	case 0x53:
		decodeInstructionWithoutOperands(d, "i64.lt_s", i)
	case 0x54:
		decodeInstructionWithoutOperands(d, "i64.lt_u", i)
	case 0x55:
		decodeInstructionWithoutOperands(d, "i64.gt_s", i)
	case 0x56:
		decodeInstructionWithoutOperands(d, "i64.gt_u", i)
	case 0x57:
		decodeInstructionWithoutOperands(d, "i64.le_s", i)
	case 0x58:
		decodeInstructionWithoutOperands(d, "i64.le_u", i)
	case 0x59:
		decodeInstructionWithoutOperands(d, "i64.ge_s", i)
	case 0x5a:
		decodeInstructionWithoutOperands(d, "i64.ge_u", i)

	case 0x5b:
		decodeInstructionWithoutOperands(d, "f32.eq", i)
	case 0x5c:
		decodeInstructionWithoutOperands(d, "f32.ne", i)
	case 0x5d:
		decodeInstructionWithoutOperands(d, "f32.lt", i)
	case 0x5e:
		decodeInstructionWithoutOperands(d, "f32.gt", i)
	case 0x5f:
		decodeInstructionWithoutOperands(d, "f32.le", i)
	case 0x60:
		decodeInstructionWithoutOperands(d, "f32.ge", i)

	case 0x61:
		decodeInstructionWithoutOperands(d, "f64.eq", i)
	case 0x62:
		decodeInstructionWithoutOperands(d, "f64.ne", i)
	case 0x63:
		decodeInstructionWithoutOperands(d, "f64.lt", i)
	case 0x64:
		decodeInstructionWithoutOperands(d, "f64.gt", i)
	case 0x65:
		decodeInstructionWithoutOperands(d, "f64.le", i)
	case 0x66:
		decodeInstructionWithoutOperands(d, "f64.ge", i)

	case 0x67:
		decodeInstructionWithoutOperands(d, "i32.clz", i)
	case 0x68:
		decodeInstructionWithoutOperands(d, "i32.ctz", i)
	case 0x69:
		decodeInstructionWithoutOperands(d, "i32.popcnt", i)
	case 0x6a:
		decodeInstructionWithoutOperands(d, "i32.add", i)
	case 0x6b:
		decodeInstructionWithoutOperands(d, "i32.sub", i)
	case 0x6c:
		decodeInstructionWithoutOperands(d, "i32.mul", i)
	case 0x6d:
		decodeInstructionWithoutOperands(d, "i32.div_s", i)
	case 0x6e:
		decodeInstructionWithoutOperands(d, "i32.div_u", i)
	case 0x6f:
		decodeInstructionWithoutOperands(d, "i32.rem_s", i)
	case 0x70:
		decodeInstructionWithoutOperands(d, "i32.rem_u", i)
	case 0x71:
		decodeInstructionWithoutOperands(d, "i32.and", i)
	case 0x72:
		decodeInstructionWithoutOperands(d, "i32.or", i)
	case 0x73:
		decodeInstructionWithoutOperands(d, "i32.xor", i)
	case 0x74:
		decodeInstructionWithoutOperands(d, "i32.shl", i)
	case 0x75:
		decodeInstructionWithoutOperands(d, "i32.shr_s", i)
	case 0x76:
		decodeInstructionWithoutOperands(d, "i32.shr_u", i)
	case 0x77:
		decodeInstructionWithoutOperands(d, "i32.rotl", i)
	case 0x78:
		decodeInstructionWithoutOperands(d, "i32.rotr", i)

	case 0x79:
		decodeInstructionWithoutOperands(d, "i64.clz", i)
	case 0x7a:
		decodeInstructionWithoutOperands(d, "i64.ctz", i)
	case 0x7b:
		decodeInstructionWithoutOperands(d, "i64.popcnt", i)
	case 0x7c:
		decodeInstructionWithoutOperands(d, "i64.add", i)
	case 0x7d:
		decodeInstructionWithoutOperands(d, "i64.sub", i)
	case 0x7e:
		decodeInstructionWithoutOperands(d, "i64.mul", i)
	case 0x7f:
		decodeInstructionWithoutOperands(d, "i64.div_s", i)
	case 0x80:
		decodeInstructionWithoutOperands(d, "i64.div_u", i)
	case 0x81:
		decodeInstructionWithoutOperands(d, "i64.rem_s", i)
	case 0x82:
		decodeInstructionWithoutOperands(d, "i64.rem_u", i)
	case 0x83:
		decodeInstructionWithoutOperands(d, "i64.and", i)
	case 0x84:
		decodeInstructionWithoutOperands(d, "i64.or", i)
	case 0x85:
		decodeInstructionWithoutOperands(d, "i64.xor", i)
	case 0x86:
		decodeInstructionWithoutOperands(d, "i64.shl", i)
	case 0x87:
		decodeInstructionWithoutOperands(d, "i64.shr_s", i)
	case 0x88:
		decodeInstructionWithoutOperands(d, "i64.shr_u", i)
	case 0x89:
		decodeInstructionWithoutOperands(d, "i64.rotl", i)
	case 0x8a:
		decodeInstructionWithoutOperands(d, "i64.rotr", i)

	case 0x8b:
		decodeInstructionWithoutOperands(d, "f32.abs", i)
	case 0x8c:
		decodeInstructionWithoutOperands(d, "f32.neg", i)
	case 0x8d:
		decodeInstructionWithoutOperands(d, "f32.ceil", i)
	case 0x8e:
		decodeInstructionWithoutOperands(d, "f32.floor", i)
	case 0x8f:
		decodeInstructionWithoutOperands(d, "f32.trunc", i)
	case 0x90:
		decodeInstructionWithoutOperands(d, "f32.nearest", i)
	case 0x91:
		decodeInstructionWithoutOperands(d, "f32.sqrt", i)
	case 0x92:
		decodeInstructionWithoutOperands(d, "f32.add", i)
	case 0x93:
		decodeInstructionWithoutOperands(d, "f32.sub", i)
	case 0x94:
		decodeInstructionWithoutOperands(d, "f32.mul", i)
	case 0x95:
		decodeInstructionWithoutOperands(d, "f32.div", i)
	case 0x96:
		decodeInstructionWithoutOperands(d, "f32.min", i)
	case 0x97:
		decodeInstructionWithoutOperands(d, "f32.max", i)
	case 0x98:
		decodeInstructionWithoutOperands(d, "f32.copysign", i)

	case 0x99:
		decodeInstructionWithoutOperands(d, "f64.abs", i)
	case 0x9a:
		decodeInstructionWithoutOperands(d, "f64.neg", i)
	case 0x9b:
		decodeInstructionWithoutOperands(d, "f64.ceil", i)
	case 0x9c:
		decodeInstructionWithoutOperands(d, "f64.floor", i)
	case 0x9d:
		decodeInstructionWithoutOperands(d, "f64.trunc", i)
	case 0x9e:
		decodeInstructionWithoutOperands(d, "f64.nearest", i)
	case 0x9f:
		decodeInstructionWithoutOperands(d, "f64.sqrt", i)
	case 0xa0:
		decodeInstructionWithoutOperands(d, "f64.add", i)
	case 0xa1:
		decodeInstructionWithoutOperands(d, "f64.sub", i)
	case 0xa2:
		decodeInstructionWithoutOperands(d, "f64.mul", i)
	case 0xa3:
		decodeInstructionWithoutOperands(d, "f64.div", i)
	case 0xa4:
		decodeInstructionWithoutOperands(d, "f64.min", i)
	case 0xa5:
		decodeInstructionWithoutOperands(d, "f64.max", i)
	case 0xa6:
		decodeInstructionWithoutOperands(d, "f64.copysign", i)

	case 0xa7:
		decodeInstructionWithoutOperands(d, "i32.wrap_i64", i)
	case 0xa8:
		decodeInstructionWithoutOperands(d, "i32.trunc_f32_s", i)
	case 0xa9:
		decodeInstructionWithoutOperands(d, "i32.trunc_f32_u", i)
	case 0xaa:
		decodeInstructionWithoutOperands(d, "i32.trunc_f64_s", i)
	case 0xab:
		decodeInstructionWithoutOperands(d, "i32.trunc_f64_u", i)
	case 0xac:
		decodeInstructionWithoutOperands(d, "i64.extend_i32_s", i)
	case 0xad:
		decodeInstructionWithoutOperands(d, "i64.extend_i32_u", i)
	case 0xae:
		decodeInstructionWithoutOperands(d, "i64.trunc_f32_s", i)
	case 0xaf:
		decodeInstructionWithoutOperands(d, "i64.trunc_f32_u", i)
	case 0xb0:
		decodeInstructionWithoutOperands(d, "i64.trunc_f64_s", i)
	case 0xb1:
		decodeInstructionWithoutOperands(d, "i64.trunc_f64_u", i)
	case 0xb2:
		decodeInstructionWithoutOperands(d, "f32.convert_i32_s", i)
	case 0xb3:
		decodeInstructionWithoutOperands(d, "f32.convert_i32_u", i)
	case 0xb4:
		decodeInstructionWithoutOperands(d, "f32.convert_i64_s", i)
	case 0xb5:
		decodeInstructionWithoutOperands(d, "f32.convert_i64_u", i)
	case 0xb6:
		decodeInstructionWithoutOperands(d, "f32.demote_f64", i)
	case 0xb7:
		decodeInstructionWithoutOperands(d, "f64.convert_i32_s", i)
	case 0xb8:
		decodeInstructionWithoutOperands(d, "f64.convert_i32_u", i)
	case 0xb9:
		decodeInstructionWithoutOperands(d, "f64.convert_i64_s", i)
	case 0xba:
		decodeInstructionWithoutOperands(d, "f64.convert_i64_u", i)
	case 0xbb:
		decodeInstructionWithoutOperands(d, "f64.promote_f32", i)
	case 0xbc:
		decodeInstructionWithoutOperands(d, "i32.reinterpret_f32", i)
	case 0xbd:
		decodeInstructionWithoutOperands(d, "i64.reinterpret_f64", i)
	case 0xbe:
		decodeInstructionWithoutOperands(d, "f32.reinterpret_i32", i)
	case 0xbf:
		decodeInstructionWithoutOperands(d, "f64.reinterpret_i64", i)

	case 0xc0:
		decodeInstructionWithoutOperands(d, "i32.extend8_s", i)
	case 0xc1:
		decodeInstructionWithoutOperands(d, "i32.extend16_s", i)
	case 0xc2:
		decodeInstructionWithoutOperands(d, "i64.extend8_s", i)
	case 0xc3:
		decodeInstructionWithoutOperands(d, "i64.extend16_s", i)
	case 0xc4:
		decodeInstructionWithoutOperands(d, "i64.extend32_s", i)

	case 0xd0:
		decodeRefNull(d)
	case 0xd1:
		decodeInstructionWithoutOperands(d, "ref.is_null", i)
	case 0xd2:
		decodeRefFunc(d)

	case 0xfc:
		decodePrefixedInstruction(d)

	case 0xfd:
		decodeVectorInstruction(d)
	default:
		d.Fatalf("unknown instruction: %#02x", instr[0])
	}
}

func decodeInstructionWithoutOperands(d *decode.D, name string, i byte) {
	d.FieldU8(name, d.AssertU(uint64(i)), scalar.ActualHex)
}

func decodeBlock(d *decode.D) {
	d.FieldU8("block", d.AssertU(0x02), scalar.ActualHex)
	decodeBlockType(d, "bt")
	d.FieldArray("instructions", func(d *decode.D) {
		for {
			b := d.PeekBytes(1)[0]
			if b == 0x0b {
				break
			}
			d.FieldStruct("instr", decodeInstruction)
		}
	})
	d.FieldU8("end", d.AssertU(0x0b))
}

func decodeLoop(d *decode.D) {
	d.FieldU8("loop", d.AssertU(0x03), scalar.ActualHex)
	decodeBlockType(d, "bt")
	d.FieldArray("instructions", func(d *decode.D) {
		for {
			b := d.PeekBytes(1)[0]
			if b == 0x0b {
				break
			}
			d.FieldStruct("instr", decodeInstruction)
		}
	})
	d.FieldU8("end", d.AssertU(0x0b))
}

func decodeIf(d *decode.D) {
	d.FieldU8("if", d.AssertU(0x04), scalar.ActualHex)
	decodeBlockType(d, "bt")
	elseClause := false
	d.FieldArray("in1", func(d *decode.D) {
		for {
			b := d.PeekBytes(1)[0]
			if b == 0x05 {
				elseClause = true
				break
			}
			if b == 0x0b {
				break
			}
			d.FieldStruct("instr", decodeInstruction)
		}
	})
	if elseClause {
		d.FieldU8("else", d.AssertU(0x05), scalar.ActualHex)
		d.FieldArray("in2", func(d *decode.D) {
			for {
				b := d.PeekBytes(1)[0]
				if b == 0x0b {
					break
				}
				d.FieldStruct("instr", decodeInstruction)
			}
		})
	}
	d.FieldU8("end", d.AssertU(0x0b))
}

func decodeEnd(d *decode.D) {
	d.FieldU8("end", d.AssertU(0x0b), scalar.ActualHex)
}

func decodeBr(d *decode.D) {
	d.FieldU8("br", d.AssertU(0x0c), scalar.ActualHex)
	decodeLabelIdx(d, "l")
}

func decodeBrIf(d *decode.D) {
	d.FieldU8("br_if", d.AssertU(0x0d), scalar.ActualHex)
	decodeLabelIdx(d, "l")
}

func decodeBrTable(d *decode.D) {
	d.FieldU8("br_table", d.AssertU(0x0e), scalar.ActualHex)
	decodeVec(d, "l", func(d *decode.D) {
		decodeLabelIdx(d, "l")
	})
	decodeLabelIdx(d, "lN")
}

func decodeCall(d *decode.D) {
	d.FieldU8("call", d.AssertU(0x10), scalar.ActualHex)
	decodeFuncIdx(d, "x")
}

func decodeCallIndirect(d *decode.D) {
	d.FieldU8("call", d.AssertU(0x11), scalar.ActualHex)
	decodeTypeIdx(d, "y")
	decodeTableIdx(d, "x")
}

func decodeSelectT(d *decode.D) {
	d.FieldU8("select", d.AssertU(0x1c), scalar.ActualHex)
	decodeVec(d, "t", func(d *decode.D) {
		decodeValType(d, "t")
	})
}

func decodeLocalGet(d *decode.D) {
	d.FieldU8("local.get", d.AssertU(0x20), scalar.ActualHex)
	decodeLocalIdx(d, "x")
}

func decodeLocalSet(d *decode.D) {
	d.FieldU8("local.set", d.AssertU(0x21), scalar.ActualHex)
	decodeLocalIdx(d, "x")
}

func decodeLocalTee(d *decode.D) {
	d.FieldU8("local.tee", d.AssertU(0x22), scalar.ActualHex)
	decodeLocalIdx(d, "x")
}

func decodeGlobalGet(d *decode.D) {
	d.FieldU8("global.get", d.AssertU(0x23), scalar.ActualHex)
	decodeGlobalIdx(d, "x")
}

func decodeGlobalSet(d *decode.D) {
	d.FieldU8("global.set", d.AssertU(0x24), scalar.ActualHex)
	decodeGlobalIdx(d, "x")
}

func decodeTableGet(d *decode.D) {
	d.FieldU8("table.get", d.AssertU(0x25), scalar.ActualHex)
	decodeTableIdx(d, "x")
}

func decodeTableSet(d *decode.D) {
	d.FieldU8("table.set", d.AssertU(0x26), scalar.ActualHex)
	decodeTableIdx(d, "x")
}

func decodeI32Load(d *decode.D) {
	d.FieldU8("i32.load", d.AssertU(0x28), scalar.ActualHex)
	decodeMemArg(d, "m")
}

func decodeI64Load(d *decode.D) {
	d.FieldU8("i64.load", d.AssertU(0x29), scalar.ActualHex)
	decodeMemArg(d, "m")
}

func decodeF32Load(d *decode.D) {
	d.FieldU8("f32.load", d.AssertU(0x2a), scalar.ActualHex)
	decodeMemArg(d, "m")
}

func decodeF64Load(d *decode.D) {
	d.FieldU8("f64.load", d.AssertU(0x2b), scalar.ActualHex)
	decodeMemArg(d, "m")
}

func decodeI32Load8S(d *decode.D) {
	d.FieldU8("i32.load8_s", d.AssertU(0x2c), scalar.ActualHex)
	decodeMemArg(d, "m")
}

func decodeI32Load8U(d *decode.D) {
	d.FieldU8("i32.load8_u", d.AssertU(0x2d), scalar.ActualHex)
	decodeMemArg(d, "m")
}

func decodeI32Load16S(d *decode.D) {
	d.FieldU8("i32.load16_s", d.AssertU(0x2e), scalar.ActualHex)
	decodeMemArg(d, "m")
}

func decodeI32Load16U(d *decode.D) {
	d.FieldU8("i32.load16_u", d.AssertU(0x2f), scalar.ActualHex)
	decodeMemArg(d, "m")
}

func decodeI64Load8S(d *decode.D) {
	d.FieldU8("i64.load8_s", d.AssertU(0x30), scalar.ActualHex)
	decodeMemArg(d, "m")
}

func decodeI64Load8U(d *decode.D) {
	d.FieldU8("i64.load8_u", d.AssertU(0x31), scalar.ActualHex)
	decodeMemArg(d, "m")
}

func decodeI64Load16S(d *decode.D) {
	d.FieldU8("i64.load16_s", d.AssertU(0x32), scalar.ActualHex)
	decodeMemArg(d, "m")
}

func decodeI64Load16U(d *decode.D) {
	d.FieldU8("i64.load16_u", d.AssertU(0x33), scalar.ActualHex)
	decodeMemArg(d, "m")
}

func decodeI64Load32S(d *decode.D) {
	d.FieldU8("i64.load32_s", d.AssertU(0x34), scalar.ActualHex)
	decodeMemArg(d, "m")
}

func decodeI64Load32U(d *decode.D) {
	d.FieldU8("i64.load32_u", d.AssertU(0x35), scalar.ActualHex)
	decodeMemArg(d, "m")
}

func decodeI32Store(d *decode.D) {
	d.FieldU8("i32.store", d.AssertU(0x36), scalar.ActualHex)
	decodeMemArg(d, "m")
}

func decodeI64Store(d *decode.D) {
	d.FieldU8("i64.store", d.AssertU(0x37), scalar.ActualHex)
	decodeMemArg(d, "m")
}

func decodeF32Store(d *decode.D) {
	d.FieldU8("f32.store", d.AssertU(0x38), scalar.ActualHex)
	decodeMemArg(d, "m")
}

func decodeF64Store(d *decode.D) {
	d.FieldU8("f64.store", d.AssertU(0x39), scalar.ActualHex)
	decodeMemArg(d, "m")
}

func decodeI32Store8(d *decode.D) {
	d.FieldU8("i32.store8", d.AssertU(0x3a), scalar.ActualHex)
	decodeMemArg(d, "m")
}

func decodeI32Store16(d *decode.D) {
	d.FieldU8("i32.store16", d.AssertU(0x3b), scalar.ActualHex)
	decodeMemArg(d, "m")
}

func decodeI64Store8(d *decode.D) {
	d.FieldU8("i64.store8", d.AssertU(0x3c), scalar.ActualHex)
	decodeMemArg(d, "m")
}

func decodeI64Store16(d *decode.D) {
	d.FieldU8("i64.store16", d.AssertU(0x3d), scalar.ActualHex)
	decodeMemArg(d, "m")
}

func decodeI64Store32(d *decode.D) {
	d.FieldU8("i64.store32", d.AssertU(0x3e), scalar.ActualHex)
	decodeMemArg(d, "m")
}

func decodeMemorySize(d *decode.D) {
	d.FieldU8("memory.size", d.AssertU(0x3f), scalar.ActualHex)
	d.FieldU8("reserved", d.AssertU(0x00), scalar.ActualHex)
}

func decodeMemoryGrow(d *decode.D) {
	d.FieldU8("memory.grow", d.AssertU(0x40), scalar.ActualHex)
	d.FieldU8("reserved", d.AssertU(0x00), scalar.ActualHex)
}

func decodeI32Const(d *decode.D) {
	d.FieldU8("i32.const", d.AssertU(0x41), scalar.ActualHex)
	fieldI32(d, "n")
}

func decodeI64Const(d *decode.D) {
	d.FieldU8("i64.const", d.AssertU(0x42), scalar.ActualHex)
	fieldI64(d, "n")
}

func decodeF32Const(d *decode.D) {
	d.FieldU8("f32.const", d.AssertU(0x43), scalar.ActualHex)
	d.FieldF32("z")
}

func decodeF64Const(d *decode.D) {
	d.FieldU8("f64.const", d.AssertU(0x44), scalar.ActualHex)
	d.FieldF64("z")
}

func decodeRefNull(d *decode.D) {
	d.FieldU8("ref.null", d.AssertU(0xd0), scalar.ActualHex)
	decodeRefType(d, "t")
}

func decodeRefFunc(d *decode.D) {
	d.FieldU8("ref.func", d.AssertU(0xd2), scalar.ActualHex)
	decodeFuncIdx(d, "x")
}

func decodePrefixedInstruction(d *decode.D) {
	d.FieldU8("prefix", d.AssertU(0xfc), scalar.ActualHex)
	s := peekUnsignedLEB128(d)
	v, ok := s.Actual.(uint64)
	if !ok {
		d.Fatalf("expected uint64 but got %t", s.Actual)
	}
	switch v {
	case 0:
		decodePrefixedInstructionWithoutOperands(d, "i32.trunc_sat_f32_s", v)
	case 1:
		decodePrefixedInstructionWithoutOperands(d, "i32.trunc_sat_f32_u", v)
	case 2:
		decodePrefixedInstructionWithoutOperands(d, "i32.trunc_sat_f64_s", v)
	case 3:
		decodePrefixedInstructionWithoutOperands(d, "i32.trunc_sat_f64_u", v)
	case 4:
		decodePrefixedInstructionWithoutOperands(d, "i64.trunc_sat_f32_s", v)
	case 5:
		decodePrefixedInstructionWithoutOperands(d, "i64.trunc_sat_f32_u", v)
	case 6:
		decodePrefixedInstructionWithoutOperands(d, "i64.trunc_sat_f64_s", v)
	case 7:
		decodePrefixedInstructionWithoutOperands(d, "i64.trunc_sat_f64_u", v)

	case 8:
		decodeMemoryInit(d)
	case 9:
		decodeDataDrop(d)
	case 10:
		decodeMemoryCopy(d)
	case 11:
		decodeMemoryFill(d)
	case 12:
		decodeTableInit(d)
	case 13:
		decodeElemDrop(d)
	case 14:
		decodeTableCopy(d)
	case 15:
		decodeTableGrow(d)
	case 16:
		decodeTableSize(d)
	case 17:
		decodeTableFill(d)

	default:
		d.Fatalf("unknown prefixed instruction: 0xfc %d", v)
	}
}

func decodePrefixedInstructionWithoutOperands(d *decode.D, name string, i uint64) {
	d.FieldUScalarFn(name, readUnsignedLEB128, d.AssertU(i))
}

func decodeMemoryInit(d *decode.D) {
	d.FieldUScalarFn("memory.init", readUnsignedLEB128, d.AssertU(8))
	decodeDataIdx(d, "x")
	d.FieldU8("reserved", scalar.ActualHex, d.AssertU(0))
}

func decodeDataDrop(d *decode.D) {
	d.FieldUScalarFn("data.drop", readUnsignedLEB128, d.AssertU(9))
	decodeDataIdx(d, "x")
}

func decodeMemoryCopy(d *decode.D) {
	d.FieldUScalarFn("memory.copy", readUnsignedLEB128, d.AssertU(10))
	d.FieldU8("reserved1", scalar.ActualHex, d.AssertU(0))
	d.FieldU8("reserved2", scalar.ActualHex, d.AssertU(0))
}

func decodeMemoryFill(d *decode.D) {
	d.FieldUScalarFn("memory.fill", readUnsignedLEB128, d.AssertU(11))
	d.FieldU8("reserved", scalar.ActualHex, d.AssertU(0))
}

func decodeTableInit(d *decode.D) {
	d.FieldUScalarFn("table.init", readUnsignedLEB128, d.AssertU(12))
	decodeElemIdx(d, "y")
	decodeTableIdx(d, "x")
}

func decodeElemDrop(d *decode.D) {
	d.FieldUScalarFn("elem.drop", readUnsignedLEB128, d.AssertU(13))
	decodeElemIdx(d, "x")
}

func decodeTableCopy(d *decode.D) {
	d.FieldUScalarFn("table.copy", readUnsignedLEB128, d.AssertU(14))
	decodeTableIdx(d, "x")
	decodeTableIdx(d, "y")
}

func decodeTableGrow(d *decode.D) {
	d.FieldUScalarFn("table.grow", readUnsignedLEB128, d.AssertU(15))
	decodeTableIdx(d, "x")
}

func decodeTableSize(d *decode.D) {
	d.FieldUScalarFn("table.size", readUnsignedLEB128, d.AssertU(16))
	decodeTableIdx(d, "x")
}

func decodeTableFill(d *decode.D) {
	d.FieldUScalarFn("table.fill", readUnsignedLEB128, d.AssertU(17))
	decodeTableIdx(d, "x")
}

func decodeVectorInstruction(d *decode.D) {
	d.FieldU8("prefix", d.AssertU(0xfd), scalar.ActualHex)
	s := peekUnsignedLEB128(d)
	v, ok := s.Actual.(uint64)
	if !ok {
		d.Fatalf("expected uint64 but got %t", s.Actual)
	}
	switch v {
	case 0:
		decodeV128Load(d)
	case 1:
		decodeV128Load8x8S(d)
	case 2:
		decodeV128Load8x8U(d)
	case 3:
		decodeV128Load16x4S(d)
	case 4:
		decodeV128Load16x4U(d)
	case 5:
		decodeV128Load32x2S(d)
	case 6:
		decodeV128Load32x2U(d)
	case 7:
		decodeV128Load8Splat(d)
	case 8:
		decodeV128Load16Splat(d)
	case 9:
		decodeV128Load32Splat(d)
	case 10:
		decodeV128Load64Splat(d)
	case 11:
		decodeV128Store(d)
	case 12:
		decodeV128Const(d)
	case 13:
		decodeI8x16Shuffle(d)
	case 14:
		decodeVectorInstructionWithoutOperands(d, "i8x16.swizzle", v)
	case 15:
		decodeVectorInstructionWithoutOperands(d, "i8x16.splat", v)
	case 16:
		decodeVectorInstructionWithoutOperands(d, "i16x8.splat", v)
	case 17:
		decodeVectorInstructionWithoutOperands(d, "i32x4.splat", v)
	case 18:
		decodeVectorInstructionWithoutOperands(d, "i64x2.splat", v)
	case 19:
		decodeVectorInstructionWithoutOperands(d, "f32x4.splat", v)
	case 20:
		decodeVectorInstructionWithoutOperands(d, "f64x2.splat", v)
	case 21:
		decodeI8x16ExtractLaneS(d)
	case 22:
		decodeI8x16ExtractLaneU(d)
	case 23:
		decodeI8x16ReplaceLane(d)
	case 24:
		decodeI16x8ExtractLaneS(d)
	case 25:
		decodeI16x8ExtractLaneU(d)
	case 26:
		decodeI16x8ReplaceLane(d)
	case 27:
		decodeI32x4ExtractLane(d)
	case 28:
		decodeI32x4ReplaceLane(d)
	case 29:
		decodeI64x2ExtractLane(d)
	case 30:
		decodeI64x2ReplaceLane(d)
	case 31:
		decodeF32x4ExtractLane(d)
	case 32:
		decodeF32x4ReplaceLane(d)
	case 33:
		decodeF64x2ExtractLane(d)
	case 34:
		decodeF64x2ReplaceLane(d)
	case 35:
		decodeVectorInstructionWithoutOperands(d, "i8x16.eq", v)
	case 36:
		decodeVectorInstructionWithoutOperands(d, "i8x16.ne", v)
	case 37:
		decodeVectorInstructionWithoutOperands(d, "i8x16.lt_s", v)
	case 38:
		decodeVectorInstructionWithoutOperands(d, "i8x16.lt_u", v)
	case 39:
		decodeVectorInstructionWithoutOperands(d, "i8x16.gt_s", v)
	case 40:
		decodeVectorInstructionWithoutOperands(d, "i8x16.gt_u", v)
	case 41:
		decodeVectorInstructionWithoutOperands(d, "i8x16.le_s", v)
	case 42:
		decodeVectorInstructionWithoutOperands(d, "i8x16.le_u", v)
	case 43:
		decodeVectorInstructionWithoutOperands(d, "i8x16.ge_s", v)
	case 44:
		decodeVectorInstructionWithoutOperands(d, "i8x16.ge_u", v)
	case 45:
		decodeVectorInstructionWithoutOperands(d, "i16x8.eq", v)
	case 46:
		decodeVectorInstructionWithoutOperands(d, "i16x8.ne", v)
	case 47:
		decodeVectorInstructionWithoutOperands(d, "i16x8.lt_s", v)
	case 48:
		decodeVectorInstructionWithoutOperands(d, "i16x8.lt_u", v)
	case 49:
		decodeVectorInstructionWithoutOperands(d, "i16x8.gt_s", v)
	case 50:
		decodeVectorInstructionWithoutOperands(d, "i16x8.gt_u", v)
	case 51:
		decodeVectorInstructionWithoutOperands(d, "i16x8.le_s", v)
	case 52:
		decodeVectorInstructionWithoutOperands(d, "i16x8.le_u", v)
	case 53:
		decodeVectorInstructionWithoutOperands(d, "i16x8.ge_s", v)
	case 54:
		decodeVectorInstructionWithoutOperands(d, "i16x8.ge_u", v)
	case 55:
		decodeVectorInstructionWithoutOperands(d, "i32x4.eq", v)
	case 56:
		decodeVectorInstructionWithoutOperands(d, "i32x4.ne", v)
	case 57:
		decodeVectorInstructionWithoutOperands(d, "i32x4.lt_s", v)
	case 58:
		decodeVectorInstructionWithoutOperands(d, "i32x4.lt_u", v)
	case 59:
		decodeVectorInstructionWithoutOperands(d, "i32x4.gt_s", v)
	case 60:
		decodeVectorInstructionWithoutOperands(d, "i32x4.gt_u", v)
	case 61:
		decodeVectorInstructionWithoutOperands(d, "i32x4.le_s", v)
	case 62:
		decodeVectorInstructionWithoutOperands(d, "i32x4.le_u", v)
	case 63:
		decodeVectorInstructionWithoutOperands(d, "i32x4.ge_s", v)
	case 64:
		decodeVectorInstructionWithoutOperands(d, "i32x4.ge_u", v)
	case 65:
		decodeVectorInstructionWithoutOperands(d, "f32x4.eq", v)
	case 66:
		decodeVectorInstructionWithoutOperands(d, "f32x4.ne", v)
	case 67:
		decodeVectorInstructionWithoutOperands(d, "f32x4.lt", v)
	case 68:
		decodeVectorInstructionWithoutOperands(d, "f32x4.gt", v)
	case 69:
		decodeVectorInstructionWithoutOperands(d, "f32x4.le", v)
	case 70:
		decodeVectorInstructionWithoutOperands(d, "f32x4.ge", v)
	case 71:
		decodeVectorInstructionWithoutOperands(d, "f64x2.eq", v)
	case 72:
		decodeVectorInstructionWithoutOperands(d, "f64x2.ne", v)
	case 73:
		decodeVectorInstructionWithoutOperands(d, "f64x2.lt", v)
	case 74:
		decodeVectorInstructionWithoutOperands(d, "f64x2.gt", v)
	case 75:
		decodeVectorInstructionWithoutOperands(d, "f64x2.le", v)
	case 76:
		decodeVectorInstructionWithoutOperands(d, "f64x2.ge", v)
	case 77:
		decodeVectorInstructionWithoutOperands(d, "v128.not", v)
	case 78:
		decodeVectorInstructionWithoutOperands(d, "v128.and", v)
	case 79:
		decodeVectorInstructionWithoutOperands(d, "v128.andnot", v)
	case 80:
		decodeVectorInstructionWithoutOperands(d, "v128.or", v)
	case 81:
		decodeVectorInstructionWithoutOperands(d, "v128.xor", v)
	case 82:
		decodeVectorInstructionWithoutOperands(d, "v128.bitselect", v)
	case 83:
		decodeVectorInstructionWithoutOperands(d, "v128.any_true", v)
	case 84:
		decodeV128Load8Lane(d)
	case 85:
		decodeV128Load16Lane(d)
	case 86:
		decodeV128Load32Lane(d)
	case 87:
		decodeV128Load64Lane(d)
	case 88:
		decodeV128Store8Lane(d)
	case 89:
		decodeV128Store16Lane(d)
	case 90:
		decodeV128Store32Lane(d)
	case 91:
		decodeV128Store64Lane(d)
	case 92:
		decodeV128Load32Zero(d)
	case 93:
		decodeV128Load64Zero(d)
	case 94:
		decodeVectorInstructionWithoutOperands(d, "f32x4.demote_f64x2_zero", v)
	case 95:
		decodeVectorInstructionWithoutOperands(d, "f64x2.promote_low_f32x4", v)
	case 96:
		decodeVectorInstructionWithoutOperands(d, "i8x16.abs", v)
	case 97:
		decodeVectorInstructionWithoutOperands(d, "i8x16.neg", v)
	case 98:
		decodeVectorInstructionWithoutOperands(d, "i8x16.popcnt", v)
	case 99:
		decodeVectorInstructionWithoutOperands(d, "i8x16.all_true", v)
	case 100:
		decodeVectorInstructionWithoutOperands(d, "i8x16.bitmask", v)
	case 101:
		decodeVectorInstructionWithoutOperands(d, "i8x16.narrow_i16x8_s", v)
	case 102:
		decodeVectorInstructionWithoutOperands(d, "i8x16.narrow_i16x8_u", v)
	case 103:
		decodeVectorInstructionWithoutOperands(d, "f32x4.ceil", v)
	case 104:
		decodeVectorInstructionWithoutOperands(d, "f32x4.floor", v)
	case 105:
		decodeVectorInstructionWithoutOperands(d, "f32x4.trunc", v)
	case 106:
		decodeVectorInstructionWithoutOperands(d, "f32x4.nearest", v)
	case 107:
		decodeVectorInstructionWithoutOperands(d, "i8x16.shl", v)
	case 108:
		decodeVectorInstructionWithoutOperands(d, "i8x16.shr_s", v)
	case 109:
		decodeVectorInstructionWithoutOperands(d, "i8x16.shr_u", v)
	case 110:
		decodeVectorInstructionWithoutOperands(d, "i8x16.add", v)
	case 111:
		decodeVectorInstructionWithoutOperands(d, "i8x16.add_sat_s", v)
	case 112:
		decodeVectorInstructionWithoutOperands(d, "i8x16.add_sat_u", v)
	case 113:
		decodeVectorInstructionWithoutOperands(d, "i8x16.sub", v)
	case 114:
		decodeVectorInstructionWithoutOperands(d, "i8x16.sub_sat_s", v)
	case 115:
		decodeVectorInstructionWithoutOperands(d, "i8x16.sub_sat_u", v)
	case 116:
		decodeVectorInstructionWithoutOperands(d, "f64x2.ceil", v)
	case 117:
		decodeVectorInstructionWithoutOperands(d, "f64x2.floor", v)
	case 118:
		decodeVectorInstructionWithoutOperands(d, "i8x16.min_s", v)
	case 119:
		decodeVectorInstructionWithoutOperands(d, "i8x16.min_u", v)
	case 120:
		decodeVectorInstructionWithoutOperands(d, "i8x16.max_s", v)
	case 121:
		decodeVectorInstructionWithoutOperands(d, "i8x16.max_u", v)
	case 122:
		decodeVectorInstructionWithoutOperands(d, "f64x2.trunc", v)
	case 123:
		decodeVectorInstructionWithoutOperands(d, "i8x16.avgr_u", v)
	case 124:
		decodeVectorInstructionWithoutOperands(d, "i16x8.extadd_pairwise_i8x16_s", v)
	case 125:
		decodeVectorInstructionWithoutOperands(d, "i16x8.extadd_pairwise_i8x16_u", v)
	case 126:
		decodeVectorInstructionWithoutOperands(d, "i32x4.extadd_pairwise_i16x8_s", v)
	case 127:
		decodeVectorInstructionWithoutOperands(d, "i32x4.extadd_pairwise_i16x8_u", v)
	case 128:
		decodeVectorInstructionWithoutOperands(d, "i16x8.abs", v)
	case 129:
		decodeVectorInstructionWithoutOperands(d, "i16x8.neg", v)
	case 130:
		decodeVectorInstructionWithoutOperands(d, "i16x8.q15mulr_sat_s", v)
	case 131:
		decodeVectorInstructionWithoutOperands(d, "i16x8.all_true", v)
	case 132:
		decodeVectorInstructionWithoutOperands(d, "i16x8.bitmask", v)
	case 133:
		decodeVectorInstructionWithoutOperands(d, "i16x8.narrow_i32x4_s", v)
	case 134:
		decodeVectorInstructionWithoutOperands(d, "i16x8.narrow_i32x4_u", v)
	case 135:
		decodeVectorInstructionWithoutOperands(d, "i16x8.extend_low_i8x16_s", v)
	case 136:
		decodeVectorInstructionWithoutOperands(d, "i16x8.extend_high_i8x16_s", v)
	case 137:
		decodeVectorInstructionWithoutOperands(d, "i16x8.extend_low_i8x16_u", v)
	case 138:
		decodeVectorInstructionWithoutOperands(d, "i16x8.extend_high_i8x16_u", v)
	case 139:
		decodeVectorInstructionWithoutOperands(d, "i16x8.shl", v)
	case 140:
		decodeVectorInstructionWithoutOperands(d, "i16x8.shr_s", v)
	case 141:
		decodeVectorInstructionWithoutOperands(d, "i16x8.shr_u", v)
	case 142:
		decodeVectorInstructionWithoutOperands(d, "i16x8.add", v)
	case 143:
		decodeVectorInstructionWithoutOperands(d, "i16x8.add_sat_s", v)
	case 144:
		decodeVectorInstructionWithoutOperands(d, "i16x8.add_sat_u", v)
	case 145:
		decodeVectorInstructionWithoutOperands(d, "i16x8.sub", v)
	case 146:
		decodeVectorInstructionWithoutOperands(d, "i16x8.sub_sat_s", v)
	case 147:
		decodeVectorInstructionWithoutOperands(d, "i16x8.sub_sat_u", v)
	case 148:
		decodeVectorInstructionWithoutOperands(d, "f64x2.nearest", v)
	case 149:
		decodeVectorInstructionWithoutOperands(d, "i16x8.mul", v)
	case 150:
		decodeVectorInstructionWithoutOperands(d, "i16x8.min_s", v)
	case 151:
		decodeVectorInstructionWithoutOperands(d, "i16x8.min_u", v)
	case 152:
		decodeVectorInstructionWithoutOperands(d, "i16x8.max_s", v)
	case 153:
		decodeVectorInstructionWithoutOperands(d, "i16x8.max_u", v)
	case 155:
		decodeVectorInstructionWithoutOperands(d, "i16x8.avgr_u", v)
	case 156:
		decodeVectorInstructionWithoutOperands(d, "i16x8.extmul_low_i8x16_s", v)
	case 157:
		decodeVectorInstructionWithoutOperands(d, "i16x8.extmul_high_i8x16_s", v)
	case 158:
		decodeVectorInstructionWithoutOperands(d, "i16x8.extmul_low_i8x16_u", v)
	case 159:
		decodeVectorInstructionWithoutOperands(d, "i16x8.extmul_high_i8x16_u", v)
	case 160:
		decodeVectorInstructionWithoutOperands(d, "i32x4.abs", v)
	case 161:
		decodeVectorInstructionWithoutOperands(d, "i32x4.neg", v)
	case 163:
		decodeVectorInstructionWithoutOperands(d, "i32x4.all_true", v)
	case 164:
		decodeVectorInstructionWithoutOperands(d, "i32x4.bitmask", v)
	case 167:
		decodeVectorInstructionWithoutOperands(d, "i32x4.extend_low_i16x8_s", v)
	case 168:
		decodeVectorInstructionWithoutOperands(d, "i32x4.extend_high_i16x8_s", v)
	case 169:
		decodeVectorInstructionWithoutOperands(d, "i32x4.extend_low_i16x8_u", v)
	case 170:
		decodeVectorInstructionWithoutOperands(d, "i32x4.extend_high_i16x8_u", v)
	case 171:
		decodeVectorInstructionWithoutOperands(d, "i32x4.shl", v)
	case 172:
		decodeVectorInstructionWithoutOperands(d, "i32x4.shr_s", v)
	case 173:
		decodeVectorInstructionWithoutOperands(d, "i32x4.shr_u", v)
	case 174:
		decodeVectorInstructionWithoutOperands(d, "i32x4.add", v)
	case 177:
		decodeVectorInstructionWithoutOperands(d, "i32x4.sub", v)
	case 181:
		decodeVectorInstructionWithoutOperands(d, "i32x4.mul", v)
	case 182:
		decodeVectorInstructionWithoutOperands(d, "i32x4.min_s", v)
	case 183:
		decodeVectorInstructionWithoutOperands(d, "i32x4.min_u", v)
	case 184:
		decodeVectorInstructionWithoutOperands(d, "i32x4.max_s", v)
	case 185:
		decodeVectorInstructionWithoutOperands(d, "i32x4.max_u", v)
	case 186:
		decodeVectorInstructionWithoutOperands(d, "i32x4.dot_i16x8_s", v)
	case 188:
		decodeVectorInstructionWithoutOperands(d, "i32x4.extmul_low_i16x8_s", v)
	case 189:
		decodeVectorInstructionWithoutOperands(d, "i32x4.extmul_high_i16x8_s", v)
	case 190:
		decodeVectorInstructionWithoutOperands(d, "i32x4.extmul_low_i16x8_u", v)
	case 191:
		decodeVectorInstructionWithoutOperands(d, "i32x4.extmul_high_i16x8_u", v)
	case 192:
		decodeVectorInstructionWithoutOperands(d, "i64x2.abs", v)
	case 193:
		decodeVectorInstructionWithoutOperands(d, "i64x2.neg", v)
	case 195:
		decodeVectorInstructionWithoutOperands(d, "i64x2.all_true", v)
	case 196:
		decodeVectorInstructionWithoutOperands(d, "i64x2.bitmask", v)
	case 199:
		decodeVectorInstructionWithoutOperands(d, "i64x2.extend_low_i32x4_s", v)
	case 200:
		decodeVectorInstructionWithoutOperands(d, "i64x2.extend_high_i32x4_s", v)
	case 201:
		decodeVectorInstructionWithoutOperands(d, "i64x2.extend_low_i32x4_u", v)
	case 202:
		decodeVectorInstructionWithoutOperands(d, "i64x2.extend_high_i32x4_u", v)
	case 203:
		decodeVectorInstructionWithoutOperands(d, "i64x2.shl", v)
	case 204:
		decodeVectorInstructionWithoutOperands(d, "i64x2.shr_s", v)
	case 205:
		decodeVectorInstructionWithoutOperands(d, "i64x2.shr_u", v)
	case 206:
		decodeVectorInstructionWithoutOperands(d, "i64x2.add", v)
	case 209:
		decodeVectorInstructionWithoutOperands(d, "i64x2.sub", v)
	case 213:
		decodeVectorInstructionWithoutOperands(d, "i64x2.mul", v)
	case 214:
		decodeVectorInstructionWithoutOperands(d, "i64x2.eq", v)
	case 215:
		decodeVectorInstructionWithoutOperands(d, "i64x2.ne", v)
	case 216:
		decodeVectorInstructionWithoutOperands(d, "i64x2.lt_s", v)
	case 217:
		decodeVectorInstructionWithoutOperands(d, "i64x2.gt_s", v)
	case 218:
		decodeVectorInstructionWithoutOperands(d, "i64x2.le_s", v)
	case 219:
		decodeVectorInstructionWithoutOperands(d, "i64x2.ge_s", v)
	case 220:
		decodeVectorInstructionWithoutOperands(d, "i64x2.extmul_low_i32x4_s", v)
	case 221:
		decodeVectorInstructionWithoutOperands(d, "i64x2.extmul_high_i32x4_s", v)
	case 222:
		decodeVectorInstructionWithoutOperands(d, "i64x2.extmul_low_i32x4_u", v)
	case 223:
		decodeVectorInstructionWithoutOperands(d, "i64x2.extmul_high_i32x4_u", v)
	case 224:
		decodeVectorInstructionWithoutOperands(d, "f32x4.abs", v)
	case 225:
		decodeVectorInstructionWithoutOperands(d, "f32x4.neg", v)
	case 227:
		decodeVectorInstructionWithoutOperands(d, "f32x4.sqrt", v)
	case 228:
		decodeVectorInstructionWithoutOperands(d, "f32x4.add", v)
	case 229:
		decodeVectorInstructionWithoutOperands(d, "f32x4.sub", v)
	case 230:
		decodeVectorInstructionWithoutOperands(d, "f32x4.mul", v)
	case 231:
		decodeVectorInstructionWithoutOperands(d, "f32x4.div", v)
	case 232:
		decodeVectorInstructionWithoutOperands(d, "f32x4.min", v)
	case 233:
		decodeVectorInstructionWithoutOperands(d, "f32x4.max", v)
	case 234:
		decodeVectorInstructionWithoutOperands(d, "f32x4.pmin", v)
	case 235:
		decodeVectorInstructionWithoutOperands(d, "f32x4.pmax", v)
	case 236:
		decodeVectorInstructionWithoutOperands(d, "f64x2.abs", v)
	case 237:
		decodeVectorInstructionWithoutOperands(d, "f64x2.neg", v)
	case 239:
		decodeVectorInstructionWithoutOperands(d, "f64x2.sqrt", v)
	case 240:
		decodeVectorInstructionWithoutOperands(d, "f64x2.add", v)
	case 241:
		decodeVectorInstructionWithoutOperands(d, "f64x2.sub", v)
	case 242:
		decodeVectorInstructionWithoutOperands(d, "f64x2.mul", v)
	case 243:
		decodeVectorInstructionWithoutOperands(d, "f64x2.div", v)
	case 244:
		decodeVectorInstructionWithoutOperands(d, "f64x2.min", v)
	case 245:
		decodeVectorInstructionWithoutOperands(d, "f64x2.max", v)
	case 246:
		decodeVectorInstructionWithoutOperands(d, "f64x2.pmin", v)
	case 247:
		decodeVectorInstructionWithoutOperands(d, "f64x2.pmax", v)
	case 248:
		decodeVectorInstructionWithoutOperands(d, "i32x4.trunc_sat_f32x4_s", v)
	case 249:
		decodeVectorInstructionWithoutOperands(d, "i32x4.trunc_sat_f32x4_u", v)
	case 250:
		decodeVectorInstructionWithoutOperands(d, "f32x4.convert_i32x4_s", v)
	case 251:
		decodeVectorInstructionWithoutOperands(d, "f32x4.convert_i32x4_u", v)
	case 252:
		decodeVectorInstructionWithoutOperands(d, "i32x4.trunc_sat_f64x2_s_zero", v)
	case 253:
		decodeVectorInstructionWithoutOperands(d, "i32x4.trunc_sat_f64x2_u_zero", v)
	case 254:
		decodeVectorInstructionWithoutOperands(d, "f64x2.convert_low_i32x4_s", v)
	case 255:
		decodeVectorInstructionWithoutOperands(d, "f64x2.convert_low_i32x4_u", v)
	default:
		d.Fatalf("unknown vector instruction: 0xfd %d", v)
	}
}

func decodeVectorInstructionWithoutOperands(d *decode.D, name string, i uint64) {
	d.FieldUScalarFn(name, readUnsignedLEB128, d.AssertU(i))
}

func decodeV128Load(d *decode.D) {
	d.FieldUScalarFn("v128.load", readUnsignedLEB128, d.AssertU(0))
	decodeMemArg(d, "m")
}

func decodeV128Load8x8S(d *decode.D) {
	d.FieldUScalarFn("v128.load8x8_s", readUnsignedLEB128, d.AssertU(1))
	decodeMemArg(d, "m")
}

func decodeV128Load8x8U(d *decode.D) {
	d.FieldUScalarFn("v128.load8x8_u", readUnsignedLEB128, d.AssertU(2))
	decodeMemArg(d, "m")
}

func decodeV128Load16x4S(d *decode.D) {
	d.FieldUScalarFn("v128.load16x4_s", readUnsignedLEB128, d.AssertU(3))
	decodeMemArg(d, "m")
}

func decodeV128Load16x4U(d *decode.D) {
	d.FieldUScalarFn("v128.load16x4_u", readUnsignedLEB128, d.AssertU(4))
	decodeMemArg(d, "m")
}

func decodeV128Load32x2S(d *decode.D) {
	d.FieldUScalarFn("v128.load32x2_s", readUnsignedLEB128, d.AssertU(5))
	decodeMemArg(d, "m")
}

func decodeV128Load32x2U(d *decode.D) {
	d.FieldUScalarFn("v128.load32x2_u", readUnsignedLEB128, d.AssertU(6))
	decodeMemArg(d, "m")
}

func decodeV128Load8Splat(d *decode.D) {
	d.FieldUScalarFn("v128.load8_splat", readUnsignedLEB128, d.AssertU(7))
	decodeMemArg(d, "m")
}

func decodeV128Load16Splat(d *decode.D) {
	d.FieldUScalarFn("v128.load16_splat", readUnsignedLEB128, d.AssertU(8))
	decodeMemArg(d, "m")
}

func decodeV128Load32Splat(d *decode.D) {
	d.FieldUScalarFn("v128.load32_splat", readUnsignedLEB128, d.AssertU(9))
	decodeMemArg(d, "m")
}

func decodeV128Load64Splat(d *decode.D) {
	d.FieldUScalarFn("v128.load64_splat", readUnsignedLEB128, d.AssertU(10))
	decodeMemArg(d, "m")
}

func decodeV128Store(d *decode.D) {
	d.FieldUScalarFn("v128.store", readUnsignedLEB128, d.AssertU(11))
	decodeMemArg(d, "m")
}

func decodeV128Const(d *decode.D) {
	d.FieldUScalarFn("v128.const", readUnsignedLEB128, d.AssertU(12))
	d.FieldRawLen("bytes", 16*8)
}

func decodeI8x16Shuffle(d *decode.D) {
	d.FieldUScalarFn("i8x16.shuffle", readUnsignedLEB128, d.AssertU(13))
	d.FieldArray("laneidx", func(d *decode.D) {
		for i := 0; i < 16; i++ {
			decodeLaneIdx(d, "l")
		}
	})
}

func decodeLaneIdx(d *decode.D, name string) {
	d.FieldU8(name)
}

func decodeI8x16ExtractLaneS(d *decode.D) {
	d.FieldUScalarFn("i8x16.extract_lane_s", readUnsignedLEB128, d.AssertU(21))
	decodeLaneIdx(d, "l")
}

func decodeI8x16ExtractLaneU(d *decode.D) {
	d.FieldUScalarFn("i8x16.extract_lane_u", readUnsignedLEB128, d.AssertU(22))
	decodeLaneIdx(d, "l")
}

func decodeI8x16ReplaceLane(d *decode.D) {
	d.FieldUScalarFn("i8x16.replace_lane", readUnsignedLEB128, d.AssertU(23))
	decodeLaneIdx(d, "l")
}

func decodeI16x8ExtractLaneS(d *decode.D) {
	d.FieldUScalarFn("i16x8.extract_lane_s", readUnsignedLEB128, d.AssertU(24))
	decodeLaneIdx(d, "l")
}

func decodeI16x8ExtractLaneU(d *decode.D) {
	d.FieldUScalarFn("i16x8.extract_lane_u", readUnsignedLEB128, d.AssertU(25))
	decodeLaneIdx(d, "l")
}

func decodeI16x8ReplaceLane(d *decode.D) {
	d.FieldUScalarFn("i16x8.replace_lane", readUnsignedLEB128, d.AssertU(26))
	decodeLaneIdx(d, "l")
}

func decodeI32x4ExtractLane(d *decode.D) {
	d.FieldUScalarFn("i32x4.extract_lane_u", readUnsignedLEB128, d.AssertU(27))
	decodeLaneIdx(d, "l")
}

func decodeI32x4ReplaceLane(d *decode.D) {
	d.FieldUScalarFn("i32x4.replace_lane", readUnsignedLEB128, d.AssertU(28))
	decodeLaneIdx(d, "l")
}

func decodeI64x2ExtractLane(d *decode.D) {
	d.FieldUScalarFn("i64x2.extract_lane_u", readUnsignedLEB128, d.AssertU(29))
	decodeLaneIdx(d, "l")
}

func decodeI64x2ReplaceLane(d *decode.D) {
	d.FieldUScalarFn("i64x2.replace_lane", readUnsignedLEB128, d.AssertU(30))
	decodeLaneIdx(d, "l")
}

func decodeF32x4ExtractLane(d *decode.D) {
	d.FieldUScalarFn("f32x4.extract_lane_u", readUnsignedLEB128, d.AssertU(31))
	decodeLaneIdx(d, "l")
}

func decodeF32x4ReplaceLane(d *decode.D) {
	d.FieldUScalarFn("f32x4.replace_lane", readUnsignedLEB128, d.AssertU(32))
	decodeLaneIdx(d, "l")
}

func decodeF64x2ExtractLane(d *decode.D) {
	d.FieldUScalarFn("f64x2.extract_lane_u", readUnsignedLEB128, d.AssertU(33))
	decodeLaneIdx(d, "l")
}

func decodeF64x2ReplaceLane(d *decode.D) {
	d.FieldUScalarFn("f64x2.replace_lane", readUnsignedLEB128, d.AssertU(34))
	decodeLaneIdx(d, "l")
}

func decodeV128Load8Lane(d *decode.D) {
	d.FieldUScalarFn("v128.load8_lane", readUnsignedLEB128, d.AssertU(84))
	decodeMemArg(d, "m")
	decodeLaneIdx(d, "l")
}

func decodeV128Load16Lane(d *decode.D) {
	d.FieldUScalarFn("v128.load16_lane", readUnsignedLEB128, d.AssertU(85))
	decodeMemArg(d, "m")
	decodeLaneIdx(d, "l")
}

func decodeV128Load32Lane(d *decode.D) {
	d.FieldUScalarFn("v128.load32_lane", readUnsignedLEB128, d.AssertU(86))
	decodeMemArg(d, "m")
	decodeLaneIdx(d, "l")
}

func decodeV128Load64Lane(d *decode.D) {
	d.FieldUScalarFn("v128.load64_lane", readUnsignedLEB128, d.AssertU(87))
	decodeMemArg(d, "m")
	decodeLaneIdx(d, "l")
}

func decodeV128Store8Lane(d *decode.D) {
	d.FieldUScalarFn("v128.store8_lane", readUnsignedLEB128, d.AssertU(88))
	decodeMemArg(d, "m")
	decodeLaneIdx(d, "l")
}

func decodeV128Store16Lane(d *decode.D) {
	d.FieldUScalarFn("v128.store16_lane", readUnsignedLEB128, d.AssertU(89))
	decodeMemArg(d, "m")
	decodeLaneIdx(d, "l")
}

func decodeV128Store32Lane(d *decode.D) {
	d.FieldUScalarFn("v128.store32_lane", readUnsignedLEB128, d.AssertU(90))
	decodeMemArg(d, "m")
	decodeLaneIdx(d, "l")
}

func decodeV128Store64Lane(d *decode.D) {
	d.FieldUScalarFn("v128.store64_lane", readUnsignedLEB128, d.AssertU(91))
	decodeMemArg(d, "m")
	decodeLaneIdx(d, "l")
}

func decodeV128Load32Zero(d *decode.D) {
	d.FieldUScalarFn("v128.load32_zero", readUnsignedLEB128, d.AssertU(92))
	decodeMemArg(d, "m")
}

func decodeV128Load64Zero(d *decode.D) {
	d.FieldUScalarFn("v128.load64_zero", readUnsignedLEB128, d.AssertU(93))
	decodeMemArg(d, "m")
}
