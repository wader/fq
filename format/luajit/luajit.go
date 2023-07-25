package luajit

// dump   = header proto+ 0U
// header = ESC 'L' 'J' versionB flagsU [namelenU nameB*]
// proto  = lengthU pdata
// pdata  = phead bcinsW* uvdataH* kgc* knum* [debugB*]
// phead  = flagsB numparamsB framesizeB numuvB numkgcU numknU numbcU
//          [debuglenU [firstlineU numlineU]]
// kgc    = kgctypeU { ktab | (loU hiU) | (rloU rhiU iloU ihiU) | strB* }
// knum   = intU0 | (loU1 hiU)
// ktab   = narrayU nhashU karray* khash*
// karray = ktabk
// khash  = ktabk ktabk
// ktabk  = ktabtypeU { intU | (loU hiU) | strB* }
//
// B = 8 bit, H = 16 bit, W = 32 bit, U = ULEB128 of W, U0/U1 = ULEB128 of W+1

// see:
//
//	* http://scm.zoomquiet.top/data/20131216145900/index.html
//	* https://github.com/LuaJIT/LuaJIT/blob/v2.1/src/lj_bcdump.h

import (
	"bytes"
	"embed"
	"encoding/binary"

	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
	"github.com/wader/fq/pkg/scalar"
)

// TODO: merge into scalar pkg
type fallbackUintMapSymStr struct {
	fallback string
	scalar.UintMapSymStr
}

func (m fallbackUintMapSymStr) MapUint(s scalar.Uint) (scalar.Uint, error) {
	s.Sym = m.fallback
	return m.UintMapSymStr.MapUint(s)
}

//go:embed luajit.md
var LuaJITFS embed.FS

func init() {
	interp.RegisterFormat(
		format.LuaJIT,
		&decode.Format{
			Description: "LuaJIT 2.0 bytecode",
			Groups:      []*decode.Group{format.Probe},
			DecodeFn:    LuaJITDecode,
		})
	interp.RegisterFS(LuaJITFS)
}

// reinterpret an int as a float
func u64tof64(u uint64) float64 {
	var buf [8]byte

	binary.BigEndian.PutUint64(buf[:], u)

	var f float64
	err := binary.Read(bytes.NewBuffer(buf[:]), binary.BigEndian, &f)
	if err != nil {
		panic(err)
	}

	return f
}

type DumpInfo struct {
	Strip     bool
	BigEndian bool
}

func LuaJITDecodeHeader(di *DumpInfo, d *decode.D) {
	d.FieldRawLen("magic", 3*8, d.AssertBitBuf([]byte{0x1b, 0x4c, 0x4a})) // ESC 'L' 'J'

	d.FieldU8("version")

	var flags uint64
	d.FieldStruct("flags", func(d *decode.D) {
		flags = d.FieldULEB128("raw")

		d.FieldValueBool("be", flags&0x01 > 0)
		d.FieldValueBool("strip", flags&0x02 > 0)
		d.FieldValueBool("ffi", flags&0x04 > 0)
		d.FieldValueBool("fr2", flags&0x08 > 0)
	})

	di.Strip = flags&0x2 > 0
	di.BigEndian = flags&0x1 > 0

	if !di.Strip {
		namelen := d.FieldU8("namelen")
		d.FieldUTF8("name", int(namelen))
	}
}

type jumpBias struct{}

func (j *jumpBias) MapUint(u scalar.Uint) (scalar.Uint, error) {
	u.Actual -= 0x8000
	return u, nil
}

func LuaJITDecodeBCIns(d *decode.D) {
	op := d.FieldU8("op", opcodes)

	d.FieldU8("a")

	if opcodes[int(op)].HasD() {
		if opcodes[int(op)].IsJump() {
			d.FieldU16("j", &jumpBias{})
		} else {
			d.FieldU16("d")
		}
	} else {
		d.FieldU8("c")
		d.FieldU8("b")
	}
}

func LuaJITDecodeNum(d *decode.D) {
	d.FieldAnyFn("value", func(d *decode.D) any {
		lo := d.ULEB128()
		hi := d.ULEB128()
		return u64tof64((hi << 32) + lo)
	})
}

func LuaJITDecodeKTabK(d *decode.D) {
	ktabtype := d.FieldULEB128("type", fallbackUintMapSymStr{
		fallback: "str",
		UintMapSymStr: scalar.UintMapSymStr{
			0: "nil",
			1: "false",
			2: "true",
			3: "int",
			4: "num",
		},
	})

	switch ktabtype {
	case 0:
		// nil
		d.FieldValueAny("value", nil)

	case 1:
		// false
		d.FieldValueBool("value", false)

	case 2:
		// true
		d.FieldValueBool("value", true)

	case 3:
		// int
		d.FieldULEB128("value")

	case 4:
		LuaJITDecodeNum(d)

	// ktabtype >= 5
	default:
		// str
		size := ktabtype - 5
		d.FieldUTF8("value", int(size))
	}
}

func LuaJITDecodeTab(d *decode.D) {
	narray := d.FieldULEB128("narray")
	nhash := d.FieldULEB128("nhash")

	d.FieldArray("array", func(d *decode.D) {
		for i := uint64(0); i < narray; i++ {
			d.FieldStruct("element", LuaJITDecodeKTabK)
		}
	})

	d.FieldArray("hash", func(d *decode.D) {
		for i := uint64(0); i < nhash; i++ {
			d.FieldStruct("pair", func(d *decode.D) {
				d.FieldStruct("key", LuaJITDecodeKTabK)
				d.FieldStruct("value", LuaJITDecodeKTabK)
			})
		}
	})
}

func LuaJITDecodeI64(d *decode.D) int64 {
	lo := d.ULEB128()
	hi := d.ULEB128()
	return int64((hi << 32) + lo)
}

func LuaJITDecodeU64(d *decode.D) uint64 {
	lo := d.ULEB128()
	hi := d.ULEB128()
	return (hi << 32) + lo
}

func LuaJITDecodeComplex(d *decode.D) {
	d.FieldAnyFn("real", func(d *decode.D) any {
		rlo := d.ULEB128()
		rhi := d.ULEB128()
		r := (rhi << 32) + rlo
		return u64tof64(r)
	})

	d.FieldAnyFn("imag", func(d *decode.D) any {
		ilo := d.ULEB128()
		ihi := d.ULEB128()
		i := (ihi << 32) + ilo
		return u64tof64(i)
	})
}

func LuaJITDecodeKGC(d *decode.D) {
	kgctype := d.FieldULEB128("type", fallbackUintMapSymStr{
		fallback: "str",
		UintMapSymStr: scalar.UintMapSymStr{
			0: "child",
			1: "tab",
			2: "i64",
			3: "u64",
			4: "complex",
		},
	})

	switch kgctype {
	case 0:
		// child

	case 1:
		LuaJITDecodeTab(d)

	case 2:
		LuaJITDecodeI64(d)

	case 3:
		LuaJITDecodeU64(d)

	case 4:
		// json does not support complex numbers,
		// so we use a struct{real: float64, imag: float64}
		d.FieldStruct("value", LuaJITDecodeComplex)

	// kgctype >= 5
	default:
		// str
		size := kgctype - 5
		d.FieldUTF8("value", int(size))
	}
}

func LuaJITDecodeKNum(d *decode.D) any {
	// knum = intU0 | (loU1 hiU)
	// ...
	// W = 32 bit, U = ULEB128 of W, U0/U1 = ULEB128 of W+1

	// intU0 encodes 33 bits : a (signed) int32, plus the LSB=0
	// loU1 encodes 33 bits : the lower half of a float64, plus the LSB=1
	// hiU encodes 32 bits : the higher half of the float64

	lo := d.ULEB128()
	if lo&1 == 0 {
		// we have an int32 (aka LuaJIT 'int')

		// drop the LSB
		lo_data := lo >> 1

		// downcast to 32bits. should not overflow
		lo_data_u32 := uint32(lo_data)

		// make it a signed integer
		lo_data_i32 := int32(lo_data_u32)

		// return a larger type to make fq happy
		return int64(lo_data_i32)
	} else {
		// we have float64 (aka LuaJIT 'number')

		hi := d.ULEB128()
		return u64tof64((hi << 32) + (lo >> 1))
	}
}

func LuaJITDecodeDebug(d *decode.D, debuglen uint64, numbc uint64) {
	d.FieldStruct("debug", func(d *decode.D) {
		d.FieldArray("lines", func(d *decode.D) {
			for i := uint64(0); i < numbc; i++ {
				d.FieldU8("value")
			}
		})

		// TODO: find out more about how to decode these strings
		d.FieldArray("annotations", func(d *decode.D) {
			i := numbc
			for i < debuglen {
				str := d.FieldUTF8Null("value")
				i += uint64(len(str) + 1)
			}
		})
	})
}

func LuaJITDecodeProto(di *DumpInfo, d *decode.D) {
	length := d.FieldULEB128("length")

	d.LimitedFn(8*int64(length), func(d *decode.D) {
		d.FieldStruct("pdata", func(d *decode.D) {
			var numuv uint64
			var numkgc uint64
			var numkn uint64
			var numbc uint64
			var debuglen uint64

			d.FieldStruct("phead", func(d *decode.D) {
				d.FieldU8("flags")
				d.FieldU8("numparams")
				d.FieldU8("framesize")
				numuv = d.FieldU8("numuv")
				numkgc = d.FieldULEB128("numkgc")
				numkn = d.FieldULEB128("numkn")
				numbc = d.FieldULEB128("numbc")

				debuglen = 0
				if !di.Strip {
					debuglen = d.FieldULEB128("debuglen")
					if debuglen > 0 {
						d.FieldULEB128("firstline")
						d.FieldULEB128("numline")
					}
				}
			})

			d.FieldArray("bcins", func(d *decode.D) {
				for i := uint64(0); i < numbc; i++ {
					d.FieldStruct("ins", func(d *decode.D) {
						LuaJITDecodeBCIns(d)
					})
				}
			})

			d.FieldArray("uvdata", func(d *decode.D) {
				for i := uint64(0); i < numuv; i++ {
					d.FieldU16("uv")
				}
			})

			d.FieldArray("kgc", func(d *decode.D) {
				for i := uint64(0); i < numkgc; i++ {
					d.FieldStruct("kgc", LuaJITDecodeKGC)
				}
			})

			d.FieldArray("knum", func(d *decode.D) {
				for i := uint64(0); i < numkn; i++ {
					d.FieldAnyFn("knum", LuaJITDecodeKNum)
				}
			})

			if !di.Strip {
				d.LimitedFn(8*int64(debuglen), func(d *decode.D) {
					LuaJITDecodeDebug(d, debuglen, numbc)
				})
			}
		})
	})
}

func LuaJITDecode(d *decode.D) any {
	di := DumpInfo{}

	d.FieldStruct("header", func(d *decode.D) {
		LuaJITDecodeHeader(&di, d)
	})

	if di.BigEndian {
		d.Endian = decode.BigEndian
	} else {
		d.Endian = decode.LittleEndian
	}

	d.FieldArray("proto", func(d *decode.D) {
		for {
			nextByte := d.PeekBytes(1)
			if bytes.Equal(nextByte, []byte{0}) {
				break
			}

			d.FieldStruct("proto", func(d *decode.D) {
				LuaJITDecodeProto(&di, d)
			})
		}
	})

	d.FieldRawLen("end", 8, d.AssertBitBuf([]byte{0}))

	return nil
}
