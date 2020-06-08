package decode

import (
	"encoding/binary"
	"fmt"
)

// ReadBits read a bits large unsigned interger from buf starting from bitPos.
// Integer is read most significant bit first.
func ReadBits(buf []byte, bitPos uint64, bits uint) uint64 {
	var n uint64
	left := bits

	if bits > 64 {
		panic(fmt.Sprintf("unsupported bit length %d", bits))
	}

	// log.Printf("bits: %#+v\n", bits)

	for left > 0 {
		bytePos := bitPos >> 3     // / 8
		byteBitPos := bitPos & 0x7 // % 8

		// log.Println("------")
		// log.Printf("n: %x\n", n)
		// log.Printf("left: %#+v\n", left)
		// log.Printf("bitPos: %d\n", bitPos)
		// log.Printf("bytePos: %#+v\n", bytePos)
		// log.Printf("byteBitPos: %#+v\n", byteBitPos)

		if byteBitPos == 0 && left%8 == 0 {
			be := binary.BigEndian
			switch left / 8 {
			case 1:
				n = n<<8 | uint64(buf[bytePos])
			case 2:
				n = n<<16 | uint64(be.Uint16(buf[bytePos:bytePos+2]))
			case 3:
				n = n<<24 |
					(uint64(be.Uint16(buf[bytePos:bytePos+2]))<<8 |
						uint64(buf[bytePos+2]))
			case 4:
				n = n<<32 |
					uint64(be.Uint32(buf[bytePos:bytePos+4]))
			case 5:
				n = n<<40 |
					(uint64(be.Uint32(buf[bytePos:bytePos+4]))<<8 |
						uint64(buf[bytePos+4]))
			case 6:
				n = n<<48 |
					(uint64(be.Uint32(buf[bytePos:bytePos+4]))<<16 |
						uint64(be.Uint16(buf[bytePos+4:bytePos+6])))
			case 7:
				n = n<<56 | (uint64(be.Uint32(buf[bytePos:bytePos+4]))<<24 |
					uint64(be.Uint16(buf[bytePos+4:bytePos+6]))<<8 |
					uint64(buf[bytePos+6]))
			case 8:
				// TODO: error if n != 0?
				n = binary.BigEndian.Uint64(buf[bytePos : bytePos+8])
			default:
				panic("unreachable")
			}
			// done
			break
		} else {
			byteBitsLeft := uint((8 - byteBitPos) & 0x7)
			// log.Printf("byteBitsLeft: %#+v\n", byteBitsLeft)
			// log.Printf("buf[bytePos]: %#+v\n", buf[bytePos])

			if byteBitsLeft == 0 {
				if left >= 8 {
					// TODO: more cases left >= 16 etc
					n = n<<8 | uint64(buf[bytePos])
					bitPos += uint64(8)
					left -= 8
				} else {
					n = n<<left | (uint64(buf[bytePos]) >> (8 - left))
					bitPos += uint64(left)
					left = 0
				}
			} else {
				if left >= byteBitsLeft {
					n = n<<byteBitsLeft | (uint64(buf[bytePos]) & ((1 << byteBitsLeft) - 1))
					bitPos += uint64(byteBitsLeft)
					left -= byteBitsLeft
				} else {
					n = n<<left | (uint64(buf[bytePos])&((1<<byteBitsLeft)-1))>>(byteBitsLeft-left)
					bitPos += uint64(left)
					// done
					break
				}
			}
		}
	}

	return n
}
