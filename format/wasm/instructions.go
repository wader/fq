package wasm

import (
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/scalar"
)

func decodeBlockType(d *decode.D, name string) {
	b := d.PeekBytes(1)[0]
	switch b {
	case 0x40:
		d.FieldU8(name, scalar.UintSym("ε"))
	case 0x6f, 0x70, 0x7b, 0x7c, 0x7d, 0x7e, 0x7f:
		d.FieldU8(name, valtypeToSymMapper)
	default:
		d.FieldSLEB128(name)
	}
}

func decodeMemArg(d *decode.D, name string) {
	d.FieldStruct(name, func(d *decode.D) {
		fieldU32(d, "a")
		fieldU32(d, "o")
	})
}

// Decode expr.
//
//	expr ::= (in:instr)* 0x0B => in* end
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

type Opcode uint64

type instructionInfo struct {
	mnemonic string
	f        func(d *decode.D) // function to decode operands
}

type instructionMap map[Opcode]instructionInfo

func (m instructionMap) MapUint(s scalar.Uint) (scalar.Uint, error) {
	opcode := s.Actual
	instr, found := m[Opcode(opcode)]
	if !found {
		return s, nil
	}

	s.Sym = instr.mnemonic
	return s, nil
}

var instrMap = instructionMap{
	0x00: {mnemonic: "unreachable"},
	0x01: {mnemonic: "nop"},

	// for the following 3 entries, `f` will be set in decodeWASM() to break initialization reference cycle.
	0x02: {mnemonic: "block"},
	0x03: {mnemonic: "loop"},
	0x04: {mnemonic: "if"},

	0x0b: {mnemonic: "end"},
	0x0c: {mnemonic: "br", f: decodeLabelIdxWithName("l")},
	0x0d: {mnemonic: "br_if", f: decodeLabelIdxWithName("l")},
	0x0e: {mnemonic: "br_table", f: decodeBrTable},
	0x0f: {mnemonic: "return"},
	0x10: {mnemonic: "call", f: decodeCall},
	0x11: {mnemonic: "call_indirect", f: decodeCallIndirect},

	0x1a: {mnemonic: "drop"},
	0x1b: {mnemonic: "select"},
	0x1c: {mnemonic: "select", f: decodeSelectT},

	0x20: {mnemonic: "local.get", f: decodeLocalIdxWithName("x")},
	0x21: {mnemonic: "local.set", f: decodeLocalIdxWithName("x")},
	0x22: {mnemonic: "local.tee", f: decodeLocalIdxWithName("x")},
	0x23: {mnemonic: "global.get", f: decodeGlobalIdxWithName("x")},
	0x24: {mnemonic: "global.set", f: decodeGlobalIdxWithName("x")},

	0x25: {mnemonic: "table.get", f: decodeTableIdxWithName("x")},
	0x26: {mnemonic: "table.set", f: decodeTableIdxWithName("x")},

	0x28: {mnemonic: "i32.load", f: decodeMemArgWithName("m")},
	0x29: {mnemonic: "i64.load", f: decodeMemArgWithName("m")},
	0x2a: {mnemonic: "f32.load", f: decodeMemArgWithName("m")},
	0x2b: {mnemonic: "f64.load", f: decodeMemArgWithName("m")},
	0x2c: {mnemonic: "i32.load8_s", f: decodeMemArgWithName("m")},
	0x2d: {mnemonic: "i32.load8_u", f: decodeMemArgWithName("m")},
	0x2e: {mnemonic: "i32.load16_s", f: decodeMemArgWithName("m")},
	0x2f: {mnemonic: "i32.load16_u", f: decodeMemArgWithName("m")},
	0x30: {mnemonic: "i64.load8_s", f: decodeMemArgWithName("m")},
	0x31: {mnemonic: "i64.load8_u", f: decodeMemArgWithName("m")},
	0x32: {mnemonic: "i64.load16_s", f: decodeMemArgWithName("m")},
	0x33: {mnemonic: "i64.load16_u", f: decodeMemArgWithName("m")},
	0x34: {mnemonic: "i64.load32_s", f: decodeMemArgWithName("m")},
	0x35: {mnemonic: "i64.load32_u", f: decodeMemArgWithName("m")},
	0x36: {mnemonic: "i32.store", f: decodeMemArgWithName("m")},
	0x37: {mnemonic: "i64.store", f: decodeMemArgWithName("m")},
	0x38: {mnemonic: "f32.store", f: decodeMemArgWithName("m")},
	0x39: {mnemonic: "f64.store", f: decodeMemArgWithName("m")},
	0x3a: {mnemonic: "i32.store8", f: decodeMemArgWithName("m")},
	0x3b: {mnemonic: "i32.store16", f: decodeMemArgWithName("m")},
	0x3c: {mnemonic: "i64.store8", f: decodeMemArgWithName("m")},
	0x3d: {mnemonic: "i64.store16", f: decodeMemArgWithName("m")},
	0x3e: {mnemonic: "i64.store32", f: decodeMemArgWithName("m")},

	0x3f: {mnemonic: "memory.size", f: decodeMemorySize},
	0x40: {mnemonic: "memory.grow", f: decodeMemoryGrow},

	0x41: {mnemonic: "i32.const", f: decodeI32Const},
	0x42: {mnemonic: "i64.const", f: decodeI64Const},
	0x43: {mnemonic: "f32.const", f: decodeF32Const},
	0x44: {mnemonic: "f64.const", f: decodeF64Const},

	0x45: {mnemonic: "32.eqz"},
	0x46: {mnemonic: "i32.eq"},
	0x47: {mnemonic: "i32.ne"},
	0x48: {mnemonic: "i32.lt_s"},
	0x49: {mnemonic: "i32.lt_u"},
	0x4a: {mnemonic: "i32.gt_s"},
	0x4b: {mnemonic: "i32.gt_u"},
	0x4c: {mnemonic: "i32.le_s"},
	0x4d: {mnemonic: "i32.le_u"},
	0x4e: {mnemonic: "i32.ge_s"},
	0x4f: {mnemonic: "i32.ge_u"},

	0x50: {mnemonic: "i64.eqz"},
	0x51: {mnemonic: "i64.eq"},
	0x52: {mnemonic: "i64.ne"},
	0x53: {mnemonic: "i64.lt_s"},
	0x54: {mnemonic: "i64.lt_u"},
	0x55: {mnemonic: "i64.gt_s"},
	0x56: {mnemonic: "i64.gt_u"},
	0x57: {mnemonic: "i64.le_s"},
	0x58: {mnemonic: "i64.le_u"},
	0x59: {mnemonic: "i64.ge_s"},
	0x5a: {mnemonic: "i64.ge_u"},

	0x5b: {mnemonic: "f32.eq"},
	0x5c: {mnemonic: "f32.ne"},
	0x5d: {mnemonic: "f32.lt"},
	0x5e: {mnemonic: "f32.gt"},
	0x5f: {mnemonic: "f32.le"},
	0x60: {mnemonic: "f32.ge"},

	0x61: {mnemonic: "f64.eq"},
	0x62: {mnemonic: "f64.ne"},
	0x63: {mnemonic: "f64.lt"},
	0x64: {mnemonic: "f64.gt"},
	0x65: {mnemonic: "f64.le"},
	0x66: {mnemonic: "f64.ge"},

	0x67: {mnemonic: "i32.clz"},
	0x68: {mnemonic: "i32.ctz"},
	0x69: {mnemonic: "i32.popcnt"},
	0x6a: {mnemonic: "i32.add"},
	0x6b: {mnemonic: "i32.sub"},
	0x6c: {mnemonic: "i32.mul"},
	0x6d: {mnemonic: "i32.div_s"},
	0x6e: {mnemonic: "i32.div_u"},
	0x6f: {mnemonic: "i32.rem_s"},
	0x70: {mnemonic: "i32.rem_u"},
	0x71: {mnemonic: "i32.and"},
	0x72: {mnemonic: "i32.or"},
	0x73: {mnemonic: "i32.xor"},
	0x74: {mnemonic: "i32.shl"},
	0x75: {mnemonic: "i32.shr_s"},
	0x76: {mnemonic: "i32.shr_u"},
	0x77: {mnemonic: "i32.rotl"},
	0x78: {mnemonic: "i32.rotr"},

	0x79: {mnemonic: "i64.clz"},
	0x7a: {mnemonic: "i64.ctz"},
	0x7b: {mnemonic: "i64.popcnt"},
	0x7c: {mnemonic: "i64.add"},
	0x7d: {mnemonic: "i64.sub"},
	0x7e: {mnemonic: "i64.mul"},
	0x7f: {mnemonic: "i64.div_s"},
	0x80: {mnemonic: "i64.div_u"},
	0x81: {mnemonic: "i64.rem_s"},
	0x82: {mnemonic: "i64.rem_u"},
	0x83: {mnemonic: "i64.and"},
	0x84: {mnemonic: "i64.or"},
	0x85: {mnemonic: "i64.xor"},
	0x86: {mnemonic: "i64.shl"},
	0x87: {mnemonic: "i64.shr_s"},
	0x88: {mnemonic: "i64.shr_u"},
	0x89: {mnemonic: "i64.rotl"},
	0x8a: {mnemonic: "i64.rotr"},

	0x8b: {mnemonic: "f32.abs"},
	0x8c: {mnemonic: "f32.neg"},
	0x8d: {mnemonic: "f32.ceil"},
	0x8e: {mnemonic: "f32.floor"},
	0x8f: {mnemonic: "f32.trunc"},
	0x90: {mnemonic: "f32.nearest"},
	0x91: {mnemonic: "f32.sqrt"},
	0x92: {mnemonic: "f32.add"},
	0x93: {mnemonic: "f32.sub"},
	0x94: {mnemonic: "f32.mul"},
	0x95: {mnemonic: "f32.div"},
	0x96: {mnemonic: "f32.min"},
	0x97: {mnemonic: "f32.max"},
	0x98: {mnemonic: "f32.copysign"},

	0x99: {mnemonic: "f64.abs"},
	0x9a: {mnemonic: "f64.neg"},
	0x9b: {mnemonic: "f64.ceil"},
	0x9c: {mnemonic: "f64.floor"},
	0x9d: {mnemonic: "f64.trunc"},
	0x9e: {mnemonic: "f64.nearest"},
	0x9f: {mnemonic: "f64.sqrt"},
	0xa0: {mnemonic: "f64.add"},
	0xa1: {mnemonic: "f64.sub"},
	0xa2: {mnemonic: "f64.mul"},
	0xa3: {mnemonic: "f64.div"},
	0xa4: {mnemonic: "f64.min"},
	0xa5: {mnemonic: "f64.max"},
	0xa6: {mnemonic: "f64.copysign"},

	0xa7: {mnemonic: "i32.wrap_i64"},
	0xa8: {mnemonic: "i32.trunc_f32_s"},
	0xa9: {mnemonic: "i32.trunc_f32_u"},
	0xaa: {mnemonic: "i32.trunc_f64_s"},
	0xab: {mnemonic: "i32.trunc_f64_u"},
	0xac: {mnemonic: "i64.extend_i32_s"},
	0xad: {mnemonic: "i64.extend_i32_u"},
	0xae: {mnemonic: "i64.trunc_f32_s"},
	0xaf: {mnemonic: "i64.trunc_f32_u"},
	0xb0: {mnemonic: "i64.trunc_f64_s"},
	0xb1: {mnemonic: "i64.trunc_f64_u"},
	0xb2: {mnemonic: "f32.convert_i32_s"},
	0xb3: {mnemonic: "f32.convert_i32_u"},
	0xb4: {mnemonic: "f32.convert_i64_s"},
	0xb5: {mnemonic: "f32.convert_i64_u"},
	0xb6: {mnemonic: "f32.demote_f64"},
	0xb7: {mnemonic: "f64.convert_i32_s"},
	0xb8: {mnemonic: "f64.convert_i32_u"},
	0xb9: {mnemonic: "f64.convert_i64_s"},
	0xba: {mnemonic: "f64.convert_i64_u"},
	0xbb: {mnemonic: "f64.promote_f32"},
	0xbc: {mnemonic: "i32.reinterpret_f32"},
	0xbd: {mnemonic: "i64.reinterpret_f64"},
	0xbe: {mnemonic: "f32.reinterpret_i32"},
	0xbf: {mnemonic: "f64.reinterpret_i64"},

	0xc0: {mnemonic: "i32.extend8_s"},
	0xc1: {mnemonic: "i32.extend16_s"},
	0xc2: {mnemonic: "i64.extend8_s"},
	0xc3: {mnemonic: "i64.extend16_s"},
	0xc4: {mnemonic: "i64.extend32_s"},

	0xd0: {mnemonic: "ref.null", f: decodeRefNull},
	0xd1: {mnemonic: "ref.is_null"},
	0xd2: {mnemonic: "ref.func", f: decodeRefFunc},

	0xfc: {mnemonic: "prefix", f: decodePrefixedInstruction},

	0xfd: {mnemonic: "vector", f: decodeVectorInstruction},
}

func decodeInstruction(d *decode.D) {
	opcode := decodeOpcode(d)
	instr := instrMap[opcode]
	if instr.f != nil {
		instr.f(d)
	}
}

func decodeOpcode(d *decode.D) Opcode {
	return Opcode(d.FieldU8("opcode", instrMap, scalar.UintHex))
}

func decodeElse(d *decode.D) {
	d.FieldU8("else", d.UintAssert(uint64(0x05)), scalar.UintHex)
}

func decodeEnd(d *decode.D) {
	d.FieldU8("end", d.UintAssert(uint64(0x0b)), scalar.UintHex)
}

func decodeBlock(d *decode.D) {
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
	decodeEnd(d)
}

func decodeLoop(d *decode.D) {
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
	decodeEnd(d)
}

func decodeIf(d *decode.D) {
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
		decodeElse(d)
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
	decodeEnd(d)
}

func decodeLabelIdxWithName(name string) func(d *decode.D) {
	return func(d *decode.D) {
		decodeLabelIdx(d, name)
	}
}

func decodeBrTable(d *decode.D) {
	decodeVec(d, "l", func(d *decode.D) {
		decodeLabelIdx(d, "l")
	})
	decodeLabelIdx(d, "lN")
}

func decodeCall(d *decode.D) {
	decodeFuncIdx(d, "x")
}

func decodeCallIndirect(d *decode.D) {
	decodeTypeIdx(d, "y")
	decodeTableIdx(d, "x")
}

func decodeSelectT(d *decode.D) {
	decodeVec(d, "t", func(d *decode.D) {
		decodeValType(d, "t")
	})
}

func decodeLocalIdxWithName(name string) func(d *decode.D) {
	return func(d *decode.D) {
		decodeLocalIdx(d, name)
	}
}

func decodeGlobalIdxWithName(name string) func(d *decode.D) {
	return func(d *decode.D) {
		decodeGlobalIdx(d, name)
	}
}

//nolint:unparam
func decodeTableIdxWithName(name string) func(d *decode.D) {
	return func(d *decode.D) {
		decodeTableIdx(d, name)
	}
}

//nolint:unparam
func decodeMemArgWithName(name string) func(d *decode.D) {
	return func(d *decode.D) {
		decodeMemArg(d, name)
	}
}

func decodeMemorySize(d *decode.D) {
	d.FieldU8("reserved", d.UintAssert(0x00), scalar.UintHex)
}

func decodeMemoryGrow(d *decode.D) {
	d.FieldU8("reserved", d.UintAssert(0x00), scalar.UintHex)
}

func decodeI32Const(d *decode.D) {
	fieldI32(d, "n")
}

func decodeI64Const(d *decode.D) {
	fieldI64(d, "n")
}

func decodeF32Const(d *decode.D) {
	d.FieldF32("z")
}

func decodeF64Const(d *decode.D) {
	d.FieldF64("z")
}

func decodeRefNull(d *decode.D) {
	decodeRefType(d, "t")
}

func decodeRefFunc(d *decode.D) {
	decodeFuncIdx(d, "x")
}

var prefixedInstrMap = instructionMap{
	0: {mnemonic: "i32.trunc_sat_f32_s"},
	1: {mnemonic: "i32.trunc_sat_f32_u"},
	2: {mnemonic: "i32.trunc_sat_f64_s"},
	3: {mnemonic: "i32.trunc_sat_f64_u"},
	4: {mnemonic: "i64.trunc_sat_f32_s"},
	5: {mnemonic: "i64.trunc_sat_f32_u"},
	6: {mnemonic: "i64.trunc_sat_f64_s"},
	7: {mnemonic: "i64.trunc_sat_f64_u"},

	8:  {mnemonic: "memory.init", f: decodeMemoryInit},
	9:  {mnemonic: "data.drop", f: decodeDataDrop},
	10: {mnemonic: "memory.copy", f: decodeMemoryCopy},
	11: {mnemonic: "memory.fill", f: decodeMemoryFill},
	12: {mnemonic: "table.init", f: decodeTableInit},
	13: {mnemonic: "elem.drop", f: decodeElemDrop},
	14: {mnemonic: "table.copy", f: decodeTableCopy},
	15: {mnemonic: "table.grow", f: decodeTableIdxWithName("x")},
	16: {mnemonic: "table.size", f: decodeTableIdxWithName("x")},
	17: {mnemonic: "table.fill", f: decodeTableIdxWithName("x")},
}

func decodePrefixedInstruction(d *decode.D) {
	opcode := decodePrefixedOpcode(d)
	instr := prefixedInstrMap[opcode]
	if instr.f != nil {
		instr.f(d)
	}
}

func decodePrefixedOpcode(d *decode.D) Opcode {
	return Opcode(d.FieldULEB128("p_opcode", prefixedInstrMap))
}

func decodeMemoryInit(d *decode.D) {
	decodeDataIdx(d, "x")
	d.FieldU8("reserved", scalar.UintHex, d.UintAssert(0))
}

func decodeDataDrop(d *decode.D) {
	decodeDataIdx(d, "x")
}

func decodeMemoryCopy(d *decode.D) {
	d.FieldU8("reserved1", scalar.UintHex, d.UintAssert(0))
	d.FieldU8("reserved2", scalar.UintHex, d.UintAssert(0))
}

func decodeMemoryFill(d *decode.D) {
	d.FieldU8("reserved", scalar.UintHex, d.UintAssert(0))
}

func decodeTableInit(d *decode.D) {
	decodeElemIdx(d, "y")
	decodeTableIdx(d, "x")
}

func decodeElemDrop(d *decode.D) {
	decodeElemIdx(d, "x")
}

func decodeTableCopy(d *decode.D) {
	decodeTableIdx(d, "x")
	decodeTableIdx(d, "y")
}

var vectorInstrMap = instructionMap{
	0:   {mnemonic: "v128.load", f: decodeMemArgWithName("m")},
	1:   {mnemonic: "v128.load8x8_s", f: decodeMemArgWithName("m")},
	2:   {mnemonic: "v128.load8x8_u", f: decodeMemArgWithName("m")},
	3:   {mnemonic: "v128.load16x4_s", f: decodeMemArgWithName("m")},
	4:   {mnemonic: "v128.load16x4_u", f: decodeMemArgWithName("m")},
	5:   {mnemonic: "v128.load32x2_s", f: decodeMemArgWithName("m")},
	6:   {mnemonic: "v128.load32x2_u", f: decodeMemArgWithName("m")},
	7:   {mnemonic: "v128.load8_splat", f: decodeMemArgWithName("m")},
	8:   {mnemonic: "v128.load16_splat", f: decodeMemArgWithName("m")},
	9:   {mnemonic: "v128.load32_splat", f: decodeMemArgWithName("m")},
	10:  {mnemonic: "v128.load64_splat", f: decodeMemArgWithName("m")},
	11:  {mnemonic: "v128.store", f: decodeMemArgWithName("m")},
	12:  {mnemonic: "v128.const", f: decodeV128ConstOperand},
	13:  {mnemonic: "i8x16.shuffle", f: decodeI8x16ShuffleOperand},
	14:  {mnemonic: "i8x16.swizzle"},
	15:  {mnemonic: "i8x16.splat"},
	16:  {mnemonic: "i16x8.splat"},
	17:  {mnemonic: "i32x4.splat"},
	18:  {mnemonic: "i64x2.splat"},
	19:  {mnemonic: "f32x4.splat"},
	20:  {mnemonic: "f64x2.splat"},
	21:  {mnemonic: "i8x16.extract_lane_s", f: decodeLaneIdxWithName("l")},
	22:  {mnemonic: "i8x16.extract_lane_u", f: decodeLaneIdxWithName("l")},
	23:  {mnemonic: "i8x16.replace_lane", f: decodeLaneIdxWithName("l")},
	24:  {mnemonic: "i16x8.extract_lane_s", f: decodeLaneIdxWithName("l")},
	25:  {mnemonic: "i16x8.extract_lane_u", f: decodeLaneIdxWithName("l")},
	26:  {mnemonic: "i16x8.replace_lane", f: decodeLaneIdxWithName("l")},
	27:  {mnemonic: "i32x4.extract_lane", f: decodeLaneIdxWithName("l")},
	28:  {mnemonic: "i32x4.replace_lane", f: decodeLaneIdxWithName("l")},
	29:  {mnemonic: "i64x2.extract_lane", f: decodeLaneIdxWithName("l")},
	30:  {mnemonic: "i64x2.replace_lane", f: decodeLaneIdxWithName("l")},
	31:  {mnemonic: "f32x4.extract_lane", f: decodeLaneIdxWithName("l")},
	32:  {mnemonic: "f32x4.replace_lane", f: decodeLaneIdxWithName("l")},
	33:  {mnemonic: "f64x2.extract_lane", f: decodeLaneIdxWithName("l")},
	34:  {mnemonic: "f64x2.replace_lane", f: decodeLaneIdxWithName("l")},
	35:  {mnemonic: "i8x16.eq"},
	36:  {mnemonic: "i8x16.ne"},
	37:  {mnemonic: "i8x16.lt_s"},
	38:  {mnemonic: "i8x16.lt_u"},
	39:  {mnemonic: "i8x16.gt_s"},
	40:  {mnemonic: "i8x16.gt_u"},
	41:  {mnemonic: "i8x16.le_s"},
	42:  {mnemonic: "i8x16.le_u"},
	43:  {mnemonic: "i8x16.ge_s"},
	44:  {mnemonic: "i8x16.ge_u"},
	45:  {mnemonic: "i16x8.eq"},
	46:  {mnemonic: "i16x8.ne"},
	47:  {mnemonic: "i16x8.lt_s"},
	48:  {mnemonic: "i16x8.lt_u"},
	49:  {mnemonic: "i16x8.gt_s"},
	50:  {mnemonic: "i16x8.gt_u"},
	51:  {mnemonic: "i16x8.le_s"},
	52:  {mnemonic: "i16x8.le_u"},
	53:  {mnemonic: "i16x8.ge_s"},
	54:  {mnemonic: "i16x8.ge_u"},
	55:  {mnemonic: "i32x4.eq"},
	56:  {mnemonic: "i32x4.ne"},
	57:  {mnemonic: "i32x4.lt_s"},
	58:  {mnemonic: "i32x4.lt_u"},
	59:  {mnemonic: "i32x4.gt_s"},
	60:  {mnemonic: "i32x4.gt_u"},
	61:  {mnemonic: "i32x4.le_s"},
	62:  {mnemonic: "i32x4.le_u"},
	63:  {mnemonic: "i32x4.ge_s"},
	64:  {mnemonic: "i32x4.ge_u"},
	65:  {mnemonic: "f32x4.eq"},
	66:  {mnemonic: "f32x4.ne"},
	67:  {mnemonic: "f32x4.lt"},
	68:  {mnemonic: "f32x4.gt"},
	69:  {mnemonic: "f32x4.le"},
	70:  {mnemonic: "f32x4.ge"},
	71:  {mnemonic: "f64x2.eq"},
	72:  {mnemonic: "f64x2.ne"},
	73:  {mnemonic: "f64x2.lt"},
	74:  {mnemonic: "f64x2.gt"},
	75:  {mnemonic: "f64x2.le"},
	76:  {mnemonic: "f64x2.ge"},
	77:  {mnemonic: "v128.not"},
	78:  {mnemonic: "v128.and"},
	79:  {mnemonic: "v128.andnot"},
	80:  {mnemonic: "v128.or"},
	81:  {mnemonic: "v128.xor"},
	82:  {mnemonic: "v128.bitselect"},
	83:  {mnemonic: "v128.any_true"},
	84:  {mnemonic: "v128.load8_lane", f: decodeMemArgAndLaneIdxWithNames("m", "l")},
	85:  {mnemonic: "v128.load16_lane", f: decodeMemArgAndLaneIdxWithNames("m", "l")},
	86:  {mnemonic: "v128.load32_lane", f: decodeMemArgAndLaneIdxWithNames("m", "l")},
	87:  {mnemonic: "v128.load64_lane", f: decodeMemArgAndLaneIdxWithNames("m", "l")},
	88:  {mnemonic: "v128.store8_lane", f: decodeMemArgAndLaneIdxWithNames("m", "l")},
	89:  {mnemonic: "v128.store16_lane", f: decodeMemArgAndLaneIdxWithNames("m", "l")},
	90:  {mnemonic: "v128.store32_lane", f: decodeMemArgAndLaneIdxWithNames("m", "l")},
	91:  {mnemonic: "v128.store64_lane", f: decodeMemArgAndLaneIdxWithNames("m", "l")},
	92:  {mnemonic: "v128.load32_zero", f: decodeMemArgWithName("m")},
	93:  {mnemonic: "v128.load64_zero", f: decodeMemArgWithName("m")},
	94:  {mnemonic: "f32x4.demote_f64x2_zero"},
	95:  {mnemonic: "f64x2.promote_low_f32x4"},
	96:  {mnemonic: "i8x16.abs"},
	97:  {mnemonic: "i8x16.neg"},
	98:  {mnemonic: "i8x16.popcnt"},
	99:  {mnemonic: "i8x16.all_true"},
	100: {mnemonic: "i8x16.bitmask"},
	101: {mnemonic: "i8x16.narrow_i16x8_s"},
	102: {mnemonic: "i8x16.narrow_i16x8_u"},
	103: {mnemonic: "f32x4.ceil"},
	104: {mnemonic: "f32x4.floor"},
	105: {mnemonic: "f32x4.trunc"},
	106: {mnemonic: "f32x4.nearest"},
	107: {mnemonic: "i8x16.shl"},
	108: {mnemonic: "i8x16.shr_s"},
	109: {mnemonic: "i8x16.shr_u"},
	110: {mnemonic: "i8x16.add"},
	111: {mnemonic: "i8x16.add_sat_s"},
	112: {mnemonic: "i8x16.add_sat_u"},
	113: {mnemonic: "i8x16.sub"},
	114: {mnemonic: "i8x16.sub_sat_s"},
	115: {mnemonic: "i8x16.sub_sat_u"},
	116: {mnemonic: "f64x2.ceil"},
	117: {mnemonic: "f64x2.floor"},
	118: {mnemonic: "i8x16.min_s"},
	119: {mnemonic: "i8x16.min_u"},
	120: {mnemonic: "i8x16.max_s"},
	121: {mnemonic: "i8x16.max_u"},
	122: {mnemonic: "f64x2.trunc"},
	123: {mnemonic: "i8x16.avgr_u"},
	124: {mnemonic: "i16x8.extadd_pairwise_i8x16_s"},
	125: {mnemonic: "i16x8.extadd_pairwise_i8x16_u"},
	126: {mnemonic: "i32x4.extadd_pairwise_i16x8_s"},
	127: {mnemonic: "i32x4.extadd_pairwise_i16x8_u"},
	128: {mnemonic: "i16x8.abs"},
	129: {mnemonic: "i16x8.neg"},
	130: {mnemonic: "i16x8.q15mulr_sat_s"},
	131: {mnemonic: "i16x8.all_true"},
	132: {mnemonic: "i16x8.bitmask"},
	133: {mnemonic: "i16x8.narrow_i32x4_s"},
	134: {mnemonic: "i16x8.narrow_i32x4_u"},
	135: {mnemonic: "i16x8.extend_low_i8x16_s"},
	136: {mnemonic: "i16x8.extend_high_i8x16_s"},
	137: {mnemonic: "i16x8.extend_low_i8x16_u"},
	138: {mnemonic: "i16x8.extend_high_i8x16_u"},
	139: {mnemonic: "i16x8.shl"},
	140: {mnemonic: "i16x8.shr_s"},
	141: {mnemonic: "i16x8.shr_u"},
	142: {mnemonic: "i16x8.add"},
	143: {mnemonic: "i16x8.add_sat_s"},
	144: {mnemonic: "i16x8.add_sat_u"},
	145: {mnemonic: "i16x8.sub"},
	146: {mnemonic: "i16x8.sub_sat_s"},
	147: {mnemonic: "i16x8.sub_sat_u"},
	148: {mnemonic: "f64x2.nearest"},
	149: {mnemonic: "i16x8.mul"},
	150: {mnemonic: "i16x8.min_s"},
	151: {mnemonic: "i16x8.min_u"},
	152: {mnemonic: "i16x8.max_s"},
	153: {mnemonic: "i16x8.max_u"},
	155: {mnemonic: "i16x8.avgr_u"},
	156: {mnemonic: "i16x8.extmul_low_i8x16_s"},
	157: {mnemonic: "i16x8.extmul_high_i8x16_s"},
	158: {mnemonic: "i16x8.extmul_low_i8x16_u"},
	159: {mnemonic: "i16x8.extmul_high_i8x16_u"},
	160: {mnemonic: "i32x4.abs"},
	161: {mnemonic: "i32x4.neg"},
	163: {mnemonic: "i32x4.all_true"},
	164: {mnemonic: "i32x4.bitmask"},
	167: {mnemonic: "i32x4.extend_low_i16x8_s"},
	168: {mnemonic: "i32x4.extend_high_i16x8_s"},
	169: {mnemonic: "i32x4.extend_low_i16x8_u"},
	170: {mnemonic: "i32x4.extend_high_i16x8_u"},
	171: {mnemonic: "i32x4.shl"},
	172: {mnemonic: "i32x4.shr_s"},
	173: {mnemonic: "i32x4.shr_u"},
	174: {mnemonic: "i32x4.add"},
	177: {mnemonic: "i32x4.sub"},
	181: {mnemonic: "i32x4.mul"},
	182: {mnemonic: "i32x4.min_s"},
	183: {mnemonic: "i32x4.min_u"},
	184: {mnemonic: "i32x4.max_s"},
	185: {mnemonic: "i32x4.max_u"},
	186: {mnemonic: "i32x4.dot_i16x8_s"},
	188: {mnemonic: "i32x4.extmul_low_i16x8_s"},
	189: {mnemonic: "i32x4.extmul_high_i16x8_s"},
	190: {mnemonic: "i32x4.extmul_low_i16x8_u"},
	191: {mnemonic: "i32x4.extmul_high_i16x8_u"},
	192: {mnemonic: "i64x2.abs"},
	193: {mnemonic: "i64x2.neg"},
	195: {mnemonic: "i64x2.all_true"},
	196: {mnemonic: "i64x2.bitmask"},
	199: {mnemonic: "i64x2.extend_low_i32x4_s"},
	200: {mnemonic: "i64x2.extend_high_i32x4_s"},
	201: {mnemonic: "i64x2.extend_low_i32x4_u"},
	202: {mnemonic: "i64x2.extend_high_i32x4_u"},
	203: {mnemonic: "i64x2.shl"},
	204: {mnemonic: "i64x2.shr_s"},
	205: {mnemonic: "i64x2.shr_u"},
	206: {mnemonic: "i64x2.add"},
	209: {mnemonic: "i64x2.sub"},
	213: {mnemonic: "i64x2.mul"},
	214: {mnemonic: "i64x2.eq"},
	215: {mnemonic: "i64x2.ne"},
	216: {mnemonic: "i64x2.lt_s"},
	217: {mnemonic: "i64x2.gt_s"},
	218: {mnemonic: "i64x2.le_s"},
	219: {mnemonic: "i64x2.ge_s"},
	220: {mnemonic: "i64x2.extmul_low_i32x4_s"},
	221: {mnemonic: "i64x2.extmul_high_i32x4_s"},
	222: {mnemonic: "i64x2.extmul_low_i32x4_u"},
	223: {mnemonic: "i64x2.extmul_high_i32x4_u"},
	224: {mnemonic: "f32x4.abs"},
	225: {mnemonic: "f32x4.neg"},
	227: {mnemonic: "f32x4.sqrt"},
	228: {mnemonic: "f32x4.add"},
	229: {mnemonic: "f32x4.sub"},
	230: {mnemonic: "f32x4.mul"},
	231: {mnemonic: "f32x4.div"},
	232: {mnemonic: "f32x4.min"},
	233: {mnemonic: "f32x4.max"},
	234: {mnemonic: "f32x4.pmin"},
	235: {mnemonic: "f32x4.pmax"},
	236: {mnemonic: "f64x2.abs"},
	237: {mnemonic: "f64x2.neg"},
	239: {mnemonic: "f64x2.sqrt"},
	240: {mnemonic: "f64x2.add"},
	241: {mnemonic: "f64x2.sub"},
	242: {mnemonic: "f64x2.mul"},
	243: {mnemonic: "f64x2.div"},
	244: {mnemonic: "f64x2.min"},
	245: {mnemonic: "f64x2.max"},
	246: {mnemonic: "f64x2.pmin"},
	247: {mnemonic: "f64x2.pmax"},
	248: {mnemonic: "i32x4.trunc_sat_f32x4_s"},
	249: {mnemonic: "i32x4.trunc_sat_f32x4_u"},
	250: {mnemonic: "f32x4.convert_i32x4_s"},
	251: {mnemonic: "f32x4.convert_i32x4_u"},
	252: {mnemonic: "i32x4.trunc_sat_f64x2_s_zero"},
	253: {mnemonic: "i32x4.trunc_sat_f64x2_u_zero"},
	254: {mnemonic: "f64x2.convert_low_i32x4_s"},
	255: {mnemonic: "f64x2.convert_low_i32x4_u"},
}

func decodeVectorInstruction(d *decode.D) {
	opcode := decodeVectorOpcode(d)
	instr := vectorInstrMap[opcode]
	df := instr.f
	if df != nil {
		df(d)
	}
}

func decodeVectorOpcode(d *decode.D) Opcode {
	return Opcode(d.FieldULEB128("v_opcode", vectorInstrMap))
}

func decodeV128ConstOperand(d *decode.D) {
	d.FieldRawLen("bytes", 16*8)
}

func decodeI8x16ShuffleOperand(d *decode.D) {
	d.FieldArray("laneidx", func(d *decode.D) {
		for i := 0; i < 16; i++ {
			decodeLaneIdx(d, "l")
		}
	})
}

func decodeLaneIdx(d *decode.D, name string) {
	d.FieldU8(name)
}

//nolint:unparam
func decodeLaneIdxWithName(name string) func(d *decode.D) {
	return func(d *decode.D) {
		decodeLaneIdx(d, name)
	}
}

//nolint:unparam
func decodeMemArgAndLaneIdxWithNames(memArgName, laneIdxName string) func(d *decode.D) {
	return func(d *decode.D) {
		decodeMemArg(d, memArgName)
		decodeLaneIdx(d, laneIdxName)
	}
}
