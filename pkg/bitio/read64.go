package bitio

import (
	"encoding/binary"
	"fmt"
)

// Read64 read nBits bits large unsigned integer from buf starting from firstBit.
// Integer is read most significant bit first.
func Read64(buf []byte, firstBit int, nBits int) uint64 {
	if nBits > 64 {
		panic(fmt.Sprintf("unsupported bit length %d", nBits))
	}

	be := binary.BigEndian
	var n uint64
	bitPos := firstBit
	bitsLeft := nBits

	for bitsLeft > 0 {
		bytePos, byteBitPos := bitPos>>3, bitPos&0x7 // / % 8

		if byteBitPos == 0 && bitsLeft&0x7 == 0 {
			bytesLeft := bitsLeft >> 3
			// BCE: let compiler know the bounds
			nBuf := buf[bytePos : bytePos+bytesLeft : bytePos+bytesLeft]

			// bitPos and bitsLeft are byte aligned
			// BCE: for some reason -1 helps remove check for some cases
			switch bytesLeft - 1 {
			case 0:
				n = n<<8 | uint64(nBuf[0])
			case 1:
				n = n<<16 | uint64(be.Uint16(nBuf))
			case 2:
				n = n<<24 |
					(uint64(be.Uint16(nBuf))<<8 |
						uint64(nBuf[2]))
			case 3:
				n = n<<32 |
					uint64(be.Uint32(nBuf))
			case 4:
				n = n<<40 |
					(uint64(be.Uint32(nBuf))<<8 |
						uint64(nBuf[4]))
			case 5:
				n = n<<48 |
					(uint64(be.Uint32(nBuf))<<16 |
						uint64(be.Uint16(nBuf[4:6])))
			case 6:
				n = n<<56 | (uint64(be.Uint32(nBuf))<<24 |
					uint64(be.Uint16(nBuf[4:6]))<<8 |
					uint64(nBuf[6]))
			case 7:
				n = be.Uint64(nBuf)
			}
			// done
			return n
		} else {
			b := buf[bytePos]

			if byteBitPos == 0 {
				// bitPos is byte aligned but not bitsLeft
				if bitsLeft >= 8 {
					// TODO: more cases left >= 16 etc
					n = n<<8 | uint64(b)
					bitPos += 8
					bitsLeft -= 8
				} else {
					n = n<<bitsLeft | (uint64(b) >> (8 - bitsLeft))
					// done
					return n
				}
			} else {
				// neither bitPos or bitsLeft byte aligned
				byteBitsLeft := (8 - byteBitPos) & 0x7
				if bitsLeft >= byteBitsLeft {
					n = n<<byteBitsLeft | (uint64(b) & ((1 << byteBitsLeft) - 1))
					bitPos += byteBitsLeft
					bitsLeft -= byteBitsLeft
				} else {
					n = n<<bitsLeft | (uint64(b)&((1<<byteBitsLeft)-1))>>(byteBitsLeft-bitsLeft)
					// done
					return n
				}
			}
		}
	}

	return n
}

func Uint64ReverseBytes(nBits int, n uint64) uint64 {
	switch {
	case nBits <= 8:
		return n
	case nBits <= 16:
		return uint64(n&0xff00>>8 | n&0xff<<8)
	case nBits <= 24:
		return uint64(n&0xff<<16 | n&0xff00 | n&0xff0000>>16)
	case nBits <= 32:
		return uint64(n&0xff<<24 | n&0xff00<<8 | n&0xff0000>>8 | n&0xff000000>>24)
	case nBits <= 40:
		return uint64(n&0xff<<32 | n&0xff00<<16 | n&0xff0000 | n&0xff000000>>16 | n&0xff00000000>>32)
	case nBits <= 48:
		return uint64(n&0xff<<40 | n&0xff00<<24 | n&0xff0000<<8 | n&0xff000000>>8 | n&0xff00000000>>24 | n&0xff0000000000>>40)
	case nBits <= 56:
		return uint64(n&0xff<<48 | n&0xff00<<32 | n&0xff0000<<16 | n&0xff000000 | n&0xff00000000>>16 | n&0xff0000000000>>32 | n&0xff000000000000>>48)
	case nBits <= 64:
		return uint64(n&0xff<<56 | n&0xff00<<40 | n&0xff0000<<24 | n&0xff000000<<8 | n&0xff00000000>>8 | n&0xff0000000000>>24 | n&0xff000000000000>>40 | n&0xff00000000000000>>56)
	default:
		panic(fmt.Sprintf("unsupported bit length %d", nBits))
	}
}
