package checksum

import (
	"fmt"
)

// TODO: lazy make table?

type Table [256]uint

func MakeTable(poly uint, bits int) Table {
	table := [256]uint{}
	mask := uint((1 << bits) - 1)

	for i := 0; i < 256; i++ {
		// note sure about -8 for > 16 bit crc
		crc := uint(i << (bits - 8))
		for j := 0; j < 8; j++ {
			if crc&(1<<(bits-1)) != 0 {
				crc = ((crc << 1) ^ poly) & mask
			} else {
				crc = (crc << 1) & mask
			}
		}
		table[i] = crc
	}

	return Table(table)
}

var ATM8Table = MakeTable(0x7, 8)
var ANSI16Table = MakeTable(0x8005, 16)
var Poly04c11db7Table = MakeTable(0x04c11db7, 32) // TODO: is this IEEE?
var IEEELETable = MakeTable(0xedb88320, 32)       // TODO: is this IEEE?

// CRC implements hash.Hash
type CRC struct {
	Bits    int
	Current uint
	Table   Table
}

func (c *CRC) Write(p []byte) (n int, err error) {
	switch c.Bits {
	case 8:
		for _, b := range p {
			c.Current = c.Table[c.Current^uint(b)]
		}
	case 16:
		for _, b := range p {
			c.Current = (c.Current<<8 ^ c.Table[(c.Current>>8)^uint(b)]) & 0xff_ff
		}
	case 32:
		for _, b := range p {
			c.Current = (c.Current<<8 ^ c.Table[(c.Current>>24)^uint(b)]) & 0xff_ff_ff_ff
		}
	default:
		panic(fmt.Sprintf("unsupported crc bit length %d", c.Bits))
	}

	return len(p), nil
}

func (c *CRC) Sum(b []byte) []byte {
	s := c.Current
	switch c.Bits {
	case 8:
		return append(b, byte(s))
	case 16:
		return append(b, byte(s>>8), byte(s))
	case 32:
		return append(b, byte(s>>24), byte(s>>16), byte(s>>8), byte(s))
	default:
		panic(fmt.Sprintf("unsupported crc bit length %d", c.Bits))
	}

}
func (c *CRC) Reset()         { c.Current = 0 }
func (c *CRC) Size() int      { return c.Bits / 8 }
func (c *CRC) BlockSize() int { return c.Bits / 8 }
