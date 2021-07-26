package bitio

import (
	"encoding/binary"
	"fmt"
)

// Read64 read nBits bits large unsigned integer from buf starting from firstBit.
// Integer is read most significant bit first.
func Read64(buf []byte, firstBit int, nBits int) uint64 {
	if nBits > 64 {
		panic(fmt.Sprintf("only supports =< 64 bits (%d)", nBits))
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

			// byteBitPos and bitsLeft are byte aligned
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
				// neither byteBitPos or bitsLeft byte aligned
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

func Write64(v uint64, nBits int, buf []byte, firstBit int) {
	if nBits > 64 {
		panic(fmt.Sprintf("only supports =< 64 bits (%d)", nBits))
	}

	be := binary.BigEndian
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
				nBuf[0] = byte(v)
			case 1:
				be.PutUint16(nBuf, uint16(v))
			case 2:
				be.PutUint16(nBuf, uint16(v>>8))
				nBuf[2] = byte(v)
			case 3:
				be.PutUint32(nBuf, uint32(v))
			case 4:
				be.PutUint32(nBuf, uint32(v>>8))
				nBuf[4] = byte(v)
			case 5:
				be.PutUint32(nBuf, uint32(v>>16))
				be.PutUint16(nBuf[4:6], uint16(v))
			case 6:
				be.PutUint32(nBuf, uint32(v>>24))
				be.PutUint16(nBuf[4:6], uint16(v>>8))
				nBuf[6] = byte(v)
			case 7:
				be.PutUint64(nBuf, v)
			}
			// done
			return
		} else {
			b := buf[bytePos]

			if byteBitPos == 0 {
				// byteBitPos is byte aligned but not bitsLeft
				if bitsLeft >= 8 {
					// TODO: more cases left >= 16 etc
					buf[bytePos] = byte(v >> (bitsLeft - 8))
					bitPos += 8
					bitsLeft -= 8
				} else {
					extraBits := 8 - bitsLeft
					buf[bytePos] = byte(v)<<extraBits | b&((1<<extraBits)-1)
					// done
					return
				}
			} else {
				// neither byteBitPos or bitsLeft byte aligned
				byteBitsLeft := (8 - byteBitPos) & 0x7
				if bitsLeft >= byteBitsLeft {
					bMask := byte((1<<byteBitPos)-1) << (8 - byteBitPos)
					buf[bytePos] = b&bMask | byte(v>>(bitsLeft-byteBitsLeft))
					bitPos += byteBitsLeft
					bitsLeft -= byteBitsLeft
				} else {
					extraBits := byteBitsLeft - bitsLeft
					bMask := byte(((1<<byteBitPos)-1)<<(8-byteBitPos) | ((1 << extraBits) - 1))
					buf[bytePos] = b&bMask | byte(v)<<extraBits
					// done
					return
				}
			}
		}
	}
}

func Uint64ReverseBytes(nBits int, n uint64) uint64 {
	switch {
	case nBits <= 8:
		return n
	case nBits <= 16:
		return n&0xff00>>8 | n&0xff<<8
	case nBits <= 24:
		return n&0xff<<16 | n&0xff00 | n&0xff0000>>16
	case nBits <= 32:
		return n&0xff<<24 | n&0xff00<<8 | n&0xff0000>>8 | n&0xff000000>>24
	case nBits <= 40:
		return n&0xff<<32 | n&0xff00<<16 | n&0xff0000 | n&0xff000000>>16 | n&0xff00000000>>32
	case nBits <= 48:
		return n&0xff<<40 | n&0xff00<<24 | n&0xff0000<<8 | n&0xff000000>>8 | n&0xff00000000>>24 | n&0xff0000000000>>40
	case nBits <= 56:
		return n&0xff<<48 | n&0xff00<<32 | n&0xff0000<<16 | n&0xff000000 | n&0xff00000000>>16 | n&0xff0000000000>>32 | n&0xff000000000000>>48
	case nBits <= 64:
		return n&0xff<<56 | n&0xff00<<40 | n&0xff0000<<24 | n&0xff000000<<8 | n&0xff00000000>>8 | n&0xff0000000000>>24 | n&0xff000000000000>>40 | n&0xff00000000000000>>56
	default:
		panic(fmt.Sprintf("unsupported bit length %d", nBits))
	}
}
