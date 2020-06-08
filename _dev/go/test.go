package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
	"strings"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")
var memprofile = flag.String("memprofile", "", "write memory profile to `file`")

type Common struct {
	Current *Field
	BitPos  int64
	Buf     []byte
}

type Type int

const (
	TypeNone = iota
	TypeInt
	TypeUint
	TypeString
	TypeData
)

type Value struct {
	Type Type
	Int  int64
	Uint uint64
	Data []byte
	Mime string
}

type Range struct {
	Start int64
	Stop  int64
}

type Field struct {
	Name     string
	Range    Range
	Display  string
	Value    Value
	Children []*Field
}

func readBits(buf []byte, bitPos int64, bits int) uint64 {
	var n uint64
	left := bits

	for left > 0 {
		bytePos := int(bitPos / 8)

		if bitPos%8 == 0 && left%8 == 0 {
			be := binary.BigEndian
			switch bits % 8 {
			case 1:
				n = n<<8 | uint64(buf[bytePos])
			case 2:
				n = n<<16 | uint64(be.Uint16(buf[bytePos:bytePos+1]))
			case 3:
				n = n<<24 |
					(uint64(be.Uint16(buf[bytePos:bytePos+1]))<<8 |
						uint64(buf[bytePos+2]))
			case 4:
				n = n<<32 |
					uint64(be.Uint32(buf[bytePos:bytePos+3]))
			case 5:
				n = n<<40 |
					(uint64(be.Uint32(buf[bytePos:bytePos+3]))<<8 |
						uint64(buf[bytePos+4]))
			case 6:
				n = n<<48 |
					(uint64(be.Uint32(buf[bytePos:bytePos+3]))<<16 |
						uint64(be.Uint16(buf[bytePos+4:bytePos+5])))
			case 7:
				n = n<<56 | (uint64(be.Uint32(buf[bytePos:bytePos+3]))<<40 |
					uint64(be.Uint16(buf[bytePos+4:bytePos+5]))<<8 |
					uint64(buf[bytePos+6]))
			case 8:
				// TODO: error if n != 0?
				n = binary.BigEndian.Uint64(buf[bytePos : bytePos+7])
			default:
				panic(fmt.Sprintf("unsupported byte length %d", bits))
			}
			bitPos += int64(left)
			left = 0
		} else {
			byteBitsLeft := int((8 - bitPos%8) % 8)
			if byteBitsLeft != 0 && left >= byteBitsLeft {
				n = n<<byteBitsLeft | (uint64(buf[bytePos]) & ((1 << byteBitsLeft) - 1))
				bitPos += int64(byteBitsLeft)
				left -= byteBitsLeft
			} else if left >= 8 {
				// TODO: more cases left >= 16 etc
				n = n<<8 | uint64(buf[bytePos])
				bitPos += int64(8)
				left -= 8
			} else {
				n = n<<left | (uint64(buf[bytePos]) >> (8 - left))
				bitPos += int64(left)
				left = 0
			}
		}
	}

	return n
}

func (c *Common) U(bits int, name string) Value {
	n := readBits(c.Buf, c.BitPos, bits)
	c.BitPos += int64(bits)
	return Value{Type: TypeUint, Uint: n}
}

func (c *Common) Field(name string, fn func() (Value, string)) (Value, string) {
	prev := c.Current

	f := &Field{Name: name}
	c.Current = f
	prev.Children = append(prev.Children, f)
	start := c.BitPos
	v, d := fn()
	stop := c.BitPos
	f.Range = Range{Start: start, Stop: stop - 1}
	f.Value = v
	f.Display = d

	c.Current = prev

	return v, d
}

func (c *Common) U8(name string) Value {
	f := &Field{Name: name}
	c.Current.Children = append(c.Current.Children, f)
	f.Range = Range{Start: c.BitPos, Stop: c.BitPos + 8}
	f.Value = Value{Type: TypeInt, Int: int64(c.Buf[c.BitPos/8])}
	c.BitPos += 8
	return f.Value
}

func (c *Common) EOF() bool {
	return c.BitPos/8 >= int64(len(c.Buf))
}

type FLAC struct {
	Common
}

func (f *FLAC) Unmarshl() {
	v, _ := f.Field("header", func() (Value, string) {
		f.U8("len")
		f.Field("flags", func() (Value, string) {
			f.U8("bla")
			f.Field("flags", func() (Value, string) {
				f.U8("bla")
				f.U8("bla1")
				f.U8("bla2")

				return Value{}, "tjo"
			})
			f.U8("bla4")

			return Value{}, "tjo"
		})
		f.U8("bla6")

		return Value{}, ""
	})

	f.U8("first frame")

	for !f.EOF() {
		f.U8("frame")
	}

	log.Printf("v: %#+v\n", v)
}

// --------------

func dump(f *Field, depth int) {
	indent := strings.Repeat("  ", depth)
	if (len(f.Children)) != 0 {
		fmt.Printf("%s%s: %#v {\n", indent, f.Name, f.Value)
		for _, c := range f.Children {
			dump(c, depth+1)
		}
		fmt.Printf("%s}\n", indent)
	} else {
		fmt.Printf("%s%s: %#v\n", indent, f.Name, f.Value)
	}
}

func main() {
	flag.Parse()

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		defer f.Close() // error handling omitted for example
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	// ... rest of the program ...

	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal("could not create memory profile: ", err)
		}
		defer f.Close() // error handling omitted for example
		runtime.GC()    // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatal("could not write memory profile: ", err)
		}
	}

	// s, err := ioutil.ReadFile(flag.Arg(0))
	// if err != nil {
	// 	panic(err)
	// }

	buf, err := ioutil.ReadFile(flag.Arg(0))
	if err != nil {
		panic(err)
	}

	f := &Field{Name: "root"}
	p := FLAC{Common: Common{Current: f, Buf: buf}}
	p.Unmarshl()

	// dump(f, 0)

}
