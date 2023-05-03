package wasm

// https://webassembly.github.io/spec/core/

import (
	"embed"
	"math"
	"sync"

	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
	"github.com/wader/fq/pkg/scalar"
)

//go:embed wasm.md
var wasmFS embed.FS

func init() {
	interp.RegisterFormat(
		format.WASM,
		&decode.Format{
			Description: "WebAssembly Binary Format",
			DecodeFn:    decodeWASM,
			Groups:      []*decode.Group{format.Probe},
		})
	interp.RegisterFS(wasmFS)
}

const (
	sectionIDCustom    = 0x00
	sectionIDType      = 0x01
	sectionIDImport    = 0x02
	sectionIDFunction  = 0x03
	sectionIDTable     = 0x04
	sectionIDMemory    = 0x05
	sectionIDGlobal    = 0x06
	sectionIDExport    = 0x07
	sectionIDStart     = 0x08
	sectionIDElement   = 0x09
	sectionIDCode      = 0x0a
	sectionIDData      = 0x0b
	sectionIDDataCount = 0x0c
)

// Decode vector.
//
//	vec(B) ::= n:u32 (x:B)^n => x^n
func decodeVec(d *decode.D, name string, fn func(d *decode.D)) {
	d.FieldStruct(name, func(d *decode.D) {
		n := fieldU32(d, "n")
		d.FieldArray("x", func(d *decode.D) {
			for i := uint64(0); i < n; i++ {
				fn(d)
			}
		})
	})
}

func decodeVecByte(d *decode.D, name string) {
	d.FieldStruct(name, func(d *decode.D) {
		n := fieldU32(d, "n")
		d.FieldRawLen("x", int64(n)*8)
	})
}

func fieldU32(d *decode.D, name string) uint64 {
	return d.FieldULEB128(name, d.UintAssertRange(0, math.MaxUint32))
}

func fieldI32(d *decode.D, name string) int64 {
	return d.FieldSLEB128(name, d.SintAssertRange(math.MinInt32, math.MaxInt32))
}

func fieldI64(d *decode.D, name string) int64 {
	return d.FieldSLEB128(name)
}

// Decode name.
//
//	name ::= b*:vec(byte) => name (if utf8(name) = b*)
func decodeName(d *decode.D, name string) {
	d.FieldStruct(name, func(d *decode.D) {
		n := fieldU32(d, "n")
		if n > math.MaxInt {
			d.Fatalf("invalid length of custom section name")
		}
		d.FieldUTF8("b", int(n))
	})
}

// Decode reftype.
//
//	reftype ::= 0x70 => funcref
//	         |  0x6F => externref
func decodeRefType(d *decode.D, name string) {
	d.FieldU8(name, reftypeTagToSym)
}

// Decode valtype.
//
//	valtype ::= t:numtype => t
//	         |  t:vectype => t
//	         |  t:reftype => t
//
//nolint:unparam
func decodeValType(d *decode.D, name string) {
	d.FieldU8(name, valtypeToSymMapper, scalar.UintHex)
}

// Decode resulttype.
//
//	resulttype ::= t*:vec(valtype) => [t*]
func decodeResultType(d *decode.D, name string) {
	decodeVec(d, name, func(d *decode.D) {
		decodeValType(d, "t")
	})
}

// Decode functype.
//
//	functype ::= 0x60 rt1:resulttype rt2:resulttype => rt1 -> rt2
func decodeFuncType(d *decode.D, name string) {
	d.FieldStruct(name, func(d *decode.D) {
		d.FieldU8("tag", d.UintAssert(0x60), scalar.UintHex)
		decodeResultType(d, "rt1")
		decodeResultType(d, "rt2")
	})
}

// Decode limits.
//
//	limits ::= 0x00 n:u32       => {min: n, max: ε}
//	        |  0x01 n:u32 m:u32 => {min: n, max: m}
func decodeLimits(d *decode.D, name string) {
	d.FieldStruct(name, func(d *decode.D) {
		tag := d.FieldU8("tag", scalar.UintHex)
		switch tag {
		case 0x00:
			fieldU32(d, "n")
		case 0x01:
			fieldU32(d, "n")
			fieldU32(d, "m")
		default:
			d.Fatalf("unknown limits type")
		}
	})
}

// Decode memtype.
//
//	memtype ::= lim:limits => lim
func decodeMemType(d *decode.D, name string) {
	d.FieldStruct(name, func(d *decode.D) {
		decodeLimits(d, "lim")
	})
}

// Decode tabletype.
//
//	tabletype ::= et:reftype lim:limits => lim et
func decodeTableType(d *decode.D, name string) {
	d.FieldStruct(name, func(d *decode.D) {
		decodeRefType(d, "et")
		decodeLimits(d, "lim")
	})
}

// Decode globaltype.
//
//	globaltype ::= t:valtype m:mut => m t
//	mut        ::= 0x00            => const
//	            |  0x01            => var
func decodeGlobalType(d *decode.D, name string) {
	d.FieldStruct(name, func(d *decode.D) {
		decodeValType(d, "t")
		d.FieldU8("m", mutToSym)
	})
}

// Decode typeidx.
//
//	typeidx ::= x:u32 => x
func decodeTypeIdx(d *decode.D, name string) {
	fieldU32(d, name)
}

// Decode funcidx.
//
//	funcidx ::= x:u32 => x
func decodeFuncIdx(d *decode.D, name string) {
	fieldU32(d, name)
}

// Decode tableidx.
//
//	tableidx ::= x:u32 => x
func decodeTableIdx(d *decode.D, name string) {
	fieldU32(d, name)
}

// Decode memidx.
//
//	memidx ::= x:u32 => x
func decodeMemIdx(d *decode.D, name string) {
	fieldU32(d, name)
}

// Decode globalidx.
//
//	globalidx ::= x:u32 => x
func decodeGlobalIdx(d *decode.D, name string) {
	fieldU32(d, name)
}

// Decode elemidx.
//
//	elemidx ::= x:u32 => x
func decodeElemIdx(d *decode.D, name string) {
	fieldU32(d, name)
}

// Decode dataidx.
//
//	dataidx ::= x:u32 => x
func decodeDataIdx(d *decode.D, name string) {
	fieldU32(d, name)
}

// Decode localidx.
//
//	localidx ::= x:u32 => x
func decodeLocalIdx(d *decode.D, name string) {
	fieldU32(d, name)
}

// Decode labelidx.
//
//	labelidx ::= l:u32 => l
func decodeLabelIdx(d *decode.D, name string) {
	fieldU32(d, name)
}

// Decode custom section.
//
//	customsec ::= section_0(custom)
//	custom    ::= name byte*
func decodeCustomSection(d *decode.D) {
	decodeName(d, "name")
	d.FieldRawLen("bytes", d.BitsLeft())
}

// Decode type section.
//
//	typesec ::= ft*:section_1(vec(functype)) => ft*
func decodeTypeSection(d *decode.D) {
	decodeVec(d, "ft", func(d *decode.D) {
		decodeFuncType(d, "ft")
	})
}

// Decode import section.
//
//	importsec  ::= im*:section_2(vec(import))    => im*
//	import     ::= mod:name nm:name d:importdesc => {module mod, name nm, desc d}
//	importdesc ::= 0x00 x:typeidx                => func x
//	            |  0x01 tt:tabletype             => table tt
//	            |  0x02 mt:memtype               => mem mt
//	            |  0x03 gt:globaltype            => global gt
func decodeImportSection(d *decode.D) {
	decodeVec(d, "im", func(d *decode.D) {
		decodeImportSegment(d, "im")
	})
}

func decodeImportSegment(d *decode.D, name string) {
	d.FieldStruct(name, func(d *decode.D) {
		decodeName(d, "mod")
		decodeName(d, "nm")
		decodeImportDesc(d, "d")
	})
}

func decodeImportDesc(d *decode.D, name string) {
	d.FieldStruct(name, func(d *decode.D) {
		tag := d.FieldU8("tag", importdescTagToSym, scalar.UintHex)
		switch tag {
		case 0x00:
			decodeTypeIdx(d, "x")
		case 0x01:
			decodeTableType(d, "tt")
		case 0x02:
			decodeMemType(d, "mt")
		case 0x03:
			decodeGlobalType(d, "gt")
		default:
			d.Fatalf("unknown import desc")
		}
	})
}

// Decode function section.
//
//	funcsec ::= x*:section_3(vec(typeidx)) => x*
func decodeFunctionSection(d *decode.D) {
	decodeVec(d, "x", func(d *decode.D) {
		decodeTypeIdx(d, "x")
	})
}

// Decode table section.
//
//	tablesec ::= tab*:section_4(vec(table)) => tab*
//	table    ::= tt:tabletype               => {type tt}
func decodeTableSection(d *decode.D) {
	decodeVec(d, "tab", func(d *decode.D) {
		decodeTableType(d, "tab")
	})
}

// Decode memory section.
//
//	memsec ::= mem*:section_5(vec(mem)) => mem*
//	mem    ::= mt:memtype               => {type mt}
func decodeMemorySection(d *decode.D) {
	decodeVec(d, "mem", func(d *decode.D) {
		decodeMemType(d, "mem")
	})
}

// Decode global section.
//
//	globalsec ::= glob*:section_6(vec(global)) => glob*
//	global    ::= gt:globaltype e:expr         => {type gt, init e}
func decodeGlobalSection(d *decode.D) {
	decodeVec(d, "glob", func(d *decode.D) {
		decodeGlobal(d, "glob")
	})
}

func decodeGlobal(d *decode.D, name string) {
	d.FieldStruct(name, func(d *decode.D) {
		decodeGlobalType(d, "gt")
		decodeExpr(d, "e")
	})
}

// Decode export section.
//
//	exportsec  ::= ex*:section_7(vec(export)) => ex*
//	export     ::= nm:name d:exportdesc       => {name nm, desc d}
//	exportdesc ::= 0x00 x:funcidx             => func x
//	            |  0x01 x:tableidx            => table x
//	            |  0x02 x:memidx              => mem x
//	            |  0x03 x:globalidx           => global x
func decodeExportSection(d *decode.D) {
	decodeVec(d, "ex", func(d *decode.D) {
		decodeExport(d, "ex")
	})
}

func decodeExport(d *decode.D, name string) {
	d.FieldStruct(name, func(d *decode.D) {
		decodeName(d, "nm")
		decodeExportDesc(d, "d")
	})
}

func decodeExportDesc(d *decode.D, name string) {
	d.FieldStruct(name, func(d *decode.D) {
		tag := d.FieldU8("tag", exportdescTagToSym, scalar.UintHex)
		switch tag {
		case 0x00:
			decodeFuncIdx(d, "x")
		case 0x01:
			decodeTableIdx(d, "x")
		case 0x02:
			decodeMemIdx(d, "x")
		case 0x03:
			decodeGlobalIdx(d, "x")
		default:
			d.Fatalf("unknown export desc")
		}
	})
}

// Decode start section.
//
//	startsec ::= st?:section_8(start) => st?
//	start    ::= x:funcidx            => {func x}
func decodeStartSection(d *decode.D) {
	decodeStart(d, "st")
}

func decodeStart(d *decode.D, name string) {
	d.FieldStruct(name, func(d *decode.D) {
		decodeFuncIdx(d, "x")
	})
}

// Decode element section.
//
//	elemsec ::= seg*:section_9(vec(elem))                           => seg*
//	elem    ::= 0:u32 e:expr y*:vec(funcidx)                        => {type funcref, init ((ref.func y) end)*, mode active {table 0, offset e}}
//	         |  1:u32 et:elemkind y*:vec(funcidx)                   => {type et, init ((ref.func y) end)*, mode passive}
//	         |  2:u32 x:tableidx e:expr et:elemkind y*:vec(funcidx) => {type et, init ((ref.func y) end)*, mode active {table x, offset e}}
//	         |  3:u32 et:elemkind y*:vec(funcidx)                   => {type et, init ((ref.func y) end)*, mode declarative}
//	         |  4:u32 e:expr el*:vec(expr)                          => {type funcref, init el*, mode active {table 0, offset e}}
//	         |  5:u32 et:reftype el*:vec(expr)                      => {type et, init el*, mode passive}
//	         |  6:u32 x:tableidx e:expr et:reftype el*:vec(expr)    => {type et, init el*, mode active {table x, offset e}}
//	         |  7:u32 et:reftype el*:vec(expr)                      => {type et, init el*, mode declarative}
func decodeElementSection(d *decode.D) {
	decodeVec(d, "seg", func(d *decode.D) {
		decodeElem(d, "seg")
	})
}

func decodeElem(d *decode.D, name string) {
	d.FieldStruct(name, func(d *decode.D) {
		tag := fieldU32(d, "tag")
		switch tag {
		case 0:
			decodeExpr(d, "e")
			decodeVec(d, "y", func(d *decode.D) {
				decodeFuncIdx(d, "y")
			})
		case 1, 3:
			decodeElemKind(d, "et")
			decodeVec(d, "y", func(d *decode.D) {
				decodeFuncIdx(d, "y")
			})
		case 2:
			decodeTableIdx(d, "x")
			decodeExpr(d, "e")
			decodeElemKind(d, "et")
			decodeVec(d, "y", func(d *decode.D) {
				decodeFuncIdx(d, "y")
			})
		case 4:
			decodeExpr(d, "e")
			decodeVec(d, "el", func(d *decode.D) {
				decodeExpr(d, "el")
			})
		case 5, 7:
			decodeRefType(d, "et")
			decodeVec(d, "el", func(d *decode.D) {
				decodeExpr(d, "el")
			})
		case 6:
			decodeTableIdx(d, "x")
			decodeExpr(d, "e")
			decodeRefType(d, "et")
			decodeVec(d, "el", func(d *decode.D) {
				decodeExpr(d, "el")
			})
		default:
			d.Fatalf("unknown elem type")
		}
	})
}

func decodeElemKind(d *decode.D, name string) {
	d.FieldU8(name, d.UintAssert(0x00), elemkindTagToSym)
}

// Decode code section.
//
//	codesec ::= code*:section_10(vec(code)) => code*
//	code    ::= size:u32 code:func          => code              (if size = ||func||)
//	func    ::= (t*)*:vec(locals) e:expr    => concat((t*)*),e   (if |concat((t*)*)| < 2^32)
//	locals  ::= n:u32 t:valtype             => t^n
func decodeCodeSection(d *decode.D) {
	decodeVec(d, "code", func(d *decode.D) {
		decodeCode(d, "code")
	})
}

func decodeCode(d *decode.D, name string) {
	d.FieldStruct(name, func(d *decode.D) {
		size := fieldU32(d, "size")
		d.FramedFn(int64(size)*8, func(d *decode.D) {
			decodeFunc(d, "code")
		})
	})
}

func decodeFunc(d *decode.D, name string) {
	d.FieldStruct(name, func(d *decode.D) {
		decodeVec(d, "t", func(d *decode.D) {
			decodeLocals(d, "t")
		})
		decodeExpr(d, "e")
	})
}

func decodeLocals(d *decode.D, name string) {
	d.FieldStruct(name, func(d *decode.D) {
		fieldU32(d, "n")
		decodeValType(d, "t")
	})
}

// Decode data section.
//
//	datasec ::= seg*:section_11(vec(data))         => seg*
//	data    ::= 0:u32 e:expr b*:vec(byte)          => {init b*, mode active {memory 0, offset e}}
//	         |  1:u32 b*:vec(byte)                 => {init b*, mode passive}
//	         |  2:u32 x:memidx e:expr b*:vec(byte) => {init b*, mode active {memory x, offset e}}
func decodeDataSection(d *decode.D) {
	decodeVec(d, "seg", func(d *decode.D) {
		decodeDataSegment(d, "seg")
	})
}

func decodeDataSegment(d *decode.D, name string) {
	d.FieldStruct(name, func(d *decode.D) {
		tag := fieldU32(d, "tag")
		switch tag {
		case 0:
			decodeExpr(d, "e")
			decodeVecByte(d, "b")
		case 1:
			decodeVecByte(d, "b")
		case 2:
			decodeMemIdx(d, "x")
			decodeExpr(d, "e")
			decodeVecByte(d, "b")
		default:
			d.Fatalf("unknown data segment type")
		}
	})
}

// Decode data count section.
//
//	datacountsec ::= n?:section_12(u32) => n?
func decodeDataCountSection(d *decode.D) {
	fieldU32(d, "n")
}

func decodeWASMModule(d *decode.D) {
	d.FieldRawLen("magic", 4*8, d.AssertBitBuf([]byte("\x00asm")))
	d.FieldU32("version")
	d.FieldArray("sections", func(d *decode.D) {
		for d.BitsLeft() > 0 {
			d.FieldStruct("section", func(d *decode.D) {
				sectionID := d.FieldU8("id", sectionIDToSym)
				size := fieldU32(d, "size")
				if size > math.MaxInt64/8 {
					d.Fatalf("invalid section size")
				}
				d.FramedFn(int64(size)*8, func(d *decode.D) {
					switch sectionID {
					case sectionIDCustom:
						d.FieldStruct("content", decodeCustomSection)
					case sectionIDType:
						d.FieldStruct("content", decodeTypeSection)
					case sectionIDImport:
						d.FieldStruct("content", decodeImportSection)
					case sectionIDFunction:
						d.FieldStruct("content", decodeFunctionSection)
					case sectionIDTable:
						d.FieldStruct("content", decodeTableSection)
					case sectionIDMemory:
						d.FieldStruct("content", decodeMemorySection)
					case sectionIDGlobal:
						d.FieldStruct("content", decodeGlobalSection)
					case sectionIDExport:
						d.FieldStruct("content", decodeExportSection)
					case sectionIDStart:
						d.FieldStruct("content", decodeStartSection)
					case sectionIDElement:
						d.FieldStruct("element", decodeElementSection)
					case sectionIDCode:
						d.FieldStruct("content", decodeCodeSection)
					case sectionIDData:
						d.FieldStruct("content", decodeDataSection)
					case sectionIDDataCount:
						d.FieldStruct("content", decodeDataCountSection)
					default:
						d.FieldRawLen("value", d.BitsLeft())
					}
				})
			})
		}
	})
}

var instrMapOnce sync.Once

const (
	opcodeBlock Opcode = 0x02
	opcodeLoop  Opcode = 0x03
	opcodeIf    Opcode = 0x04
)

func decodeWASM(d *decode.D) any {
	d.Endian = decode.LittleEndian

	// delayed initialization to break initialization reference cycle
	instrMapOnce.Do(func() {
		instrMap[opcodeBlock] = instructionInfo{mnemonic: "block", f: decodeBlock}
		instrMap[opcodeLoop] = instructionInfo{mnemonic: "loop", f: decodeLoop}
		instrMap[opcodeIf] = instructionInfo{mnemonic: "if", f: decodeIf}
	})

	decodeWASMModule(d)

	return nil
}
