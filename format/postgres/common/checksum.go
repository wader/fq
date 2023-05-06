package common

import "encoding/binary"

const (
	nSums        = 32
	fnvPrime     = 16777619
	nSumsSize    = 4 * nSums
	mainBlockLen = PageSize / nSumsSize
	RelSegSize   = 131072
)

var (
	checksumBaseOffsets = [nSums]uint32{
		0x5B1F36E9, 0xB8525960, 0x02AB50AA, 0x1DE66D2A,
		0x79FF467A, 0x9BB9F8A3, 0x217E7CD2, 0x83E13D2C,
		0xF8D4474F, 0xE39EB970, 0x42C6AE16, 0x993216FA,
		0x7B093B5D, 0x98DAFF3C, 0xF718902A, 0x0B1C9CDB,
		0xE58F764B, 0x187636BC, 0x5D7B3BB1, 0xE73DE7DE,
		0x92BEC979, 0xCCA6C0B2, 0x304A0979, 0x85AA43D4,
		0x783125BB, 0x6CA8EAA2, 0xE407EAC6, 0x4B5CFC3E,
		0x9FBF8C76, 0x15CA20BE, 0xF2CA9FD3, 0x959BD756,
	}
)

func checksumComp(checksum uint32, value uint32) uint32 {
	tmp := checksum ^ value
	checksum = tmp*fnvPrime ^ (tmp >> 17)
	return checksum
}

func pgChecksumBlock(page []byte) uint32 {
	sums := [nSums]uint32{0}
	for i := 0; i < nSums; i++ {
		sums[i] = checksumBaseOffsets[i]
	}
	result := uint32(0)

	for i := 0; i < mainBlockLen; i++ {
		for j := 0; j < nSums; j++ {
			v2 := binary.LittleEndian.Uint32(page[i*nSumsSize+j*4:])
			sums[j] = checksumComp(sums[j], v2)
		}
	}

	for i := 0; i < 2; i++ {
		for j := 0; j < nSums; j++ {
			sums[j] = checksumComp(sums[j], 0)
		}
	}

	for i := 0; i < nSums; i++ {
		result = result ^ sums[i]
	}

	return result
}

func CheckSumBlock(page []byte, blockNumber uint32) uint16 {
	// set pd_checksum to zero
	page[8] = 0
	page[9] = 0

	sum := pgChecksumBlock(page)
	sum1 := sum ^ blockNumber
	sum2 := uint16((sum1 % 65535) + 1)
	return sum2
}
