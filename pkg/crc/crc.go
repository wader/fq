package crc

// TODO: lazy make table?

type Table [256]uint

func MakeTable(poly int, bits int) Table {
	table := [256]uint{}
	mask := uint((1 << bits) - 1)

	for i := 0; i < 256; i++ {
		// note sure about -8 for > 16 bit crc
		var crc uint = uint(i << (bits - 8))
		for j := 0; j < 8; j++ {
			if crc&(1<<(bits-1)) > 0 {
				crc = ((crc << 1) ^ uint(poly)) & mask
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
			c.Current = (c.Current<<8 ^ c.Table[(c.Current>>8)^uint(b)]) & 0xffff
		}
	default:
		panic("unsupported")
	}

	return len(p), nil
}

func (c *CRC) Sum(b []byte) []byte {
	switch c.Bits {
	case 8:
		return append(b, byte(c.Current))
	case 16:
		return append(b, byte(c.Current>>8), byte(c.Current))
	default:
		panic("unsupported")
	}

}
func (c *CRC) Reset()         { c.Current = 0 }
func (c *CRC) Size() int      { return c.Bits / 8 }
func (c *CRC) BlockSize() int { return c.Bits / 8 }
