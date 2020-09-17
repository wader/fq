package bitio

import (
	"encoding/binary"
	"fmt"
)

// ReadBits read nBits bits large unsigned integer from buf starting from firstBit.
// Integer is read most significant bit first.
func ReadBits(buf []byte, firstBit int, nBits int) uint64 {
	if nBits > 64 {
		panic(fmt.Sprintf("unsupported bit length %d", nBits))
	}

	var n uint64
	bitPos := firstBit
	bitsLeft := nBits

	for bitsLeft > 0 {
		bytePos, byteBitPos := bitPos>>3, bitPos&0x7 // / % 8

		if byteBitPos == 0 && bitsLeft&0x7 == 0 {
			// bitPos and bitsLeft are byte aligned
			be := binary.BigEndian
			switch bitsLeft >> 3 {
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
				n = be.Uint64(buf[bytePos : bytePos+8])
			default:
				panic("unreachable")
			}
			// done
			break
		} else {
			if byteBitPos == 0 {
				// bitPos is byte aligned but not bitsLeft
				if bitsLeft >= 8 {
					// TODO: more cases left >= 16 etc
					n = n<<8 | uint64(buf[bytePos])
					bitPos += 8
					bitsLeft -= 8
				} else {
					n = n<<bitsLeft | (uint64(buf[bytePos]) >> (8 - bitsLeft))
					// done
					break
				}
			} else {
				// neither bitPos or bitsLeft byte aligned
				byteBitsLeft := (8 - byteBitPos) & 0x7
				if bitsLeft >= byteBitsLeft {
					n = n<<byteBitsLeft | (uint64(buf[bytePos]) & ((1 << byteBitsLeft) - 1))
					bitPos += byteBitsLeft
					bitsLeft -= byteBitsLeft
				} else {
					n = n<<bitsLeft | (uint64(buf[bytePos])&((1<<byteBitsLeft)-1))>>(byteBitsLeft-bitsLeft)
					// done
					break
				}
			}
		}
	}

	return n
}
