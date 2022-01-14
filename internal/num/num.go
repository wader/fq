package num

import (
	"fmt"
	"math"
	"math/big"
	"strconv"
	"strings"

	"github.com/wader/fq/pkg/ranges"
)

var BasePrefixMap = map[int]string{
	2:  "0b",
	8:  "0o",
	16: "0x",
}

func DigitsInBase(n int64, basePrefix bool, base int) int {
	prefixLen := 0
	if basePrefix {
		prefixLen = len(BasePrefixMap[base])
	}
	if n == 0 {
		return prefixLen + 1
	}
	return prefixLen + int(1+math.Floor(math.Log(float64(n))/math.Log(float64(base))))
}

func padFormatNumber(s string, base int, basePrefix bool, width int) string {
	prefixStr := ""
	if basePrefix {
		prefixStr = BasePrefixMap[base]
	}
	padStr := ""
	padN := width - len(s) - len(prefixStr)
	if padN > 0 {
		// TODO: something faster?
		padStr = strings.Repeat("0", padN)
	}
	return prefixStr + padStr + s
}

func PadFormatInt(i int64, base int, basePrefix bool, width int) string {
	return padFormatNumber(strconv.FormatInt(i, base), base, basePrefix, width)
}

func PadFormatUint(i uint64, base int, basePrefix bool, width int) string {
	return padFormatNumber(strconv.FormatUint(i, base), base, basePrefix, width)
}

func PadFormatBigInt(i *big.Int, base int, basePrefix bool, width int) string {
	return padFormatNumber(i.Text(base), base, basePrefix, width)
}

func MaxUInt64(a, b uint64) uint64 {
	if a < b {
		return b
	}
	return a
}

func MinUInt64(a, b uint64) uint64 {
	if a > b {
		return b
	}
	return a
}

func MaxInt64(a, b int64) int64 {
	if a < b {
		return b
	}
	return a
}

func MinInt64(a, b int64) int64 {
	if a > b {
		return b
	}
	return a
}

func MaxInt(a, b int) int {
	if a < b {
		return b
	}
	return a
}

func MinInt(a, b int) int {
	if a > b {
		return b
	}
	return a
}

func ClampInt(min, max, v int) int {
	return MaxInt(min, MinInt(max, v))
}

func ClampInt64(min, max, v int64) int64 {
	return MaxInt64(min, MinInt64(max, v))
}

type Bits uint64

func (b Bits) StringByteBits(base int) string {
	if b&0x7 != 0 {
		return BasePrefixMap[base] + strconv.FormatUint(uint64(b)>>3, base) + "." + strconv.FormatUint(uint64(b)&0x7, base)
	}
	return BasePrefixMap[base] + strconv.FormatUint(uint64(b>>3), base)
}

type BitRange ranges.Range

func (r BitRange) StringByteBits(base int) string {
	if r.Len == 0 {
		return fmt.Sprintf("%s-NA", Bits(r.Start).StringByteBits(base))
	}
	return fmt.Sprintf("%s-%s", Bits(r.Start).StringByteBits(base), Bits(r.Start+r.Len-1).StringByteBits(base))
}

func TwosComplement(nBits int, n uint64) int64 {
	if n&(1<<(nBits-1)) > 0 {
		// two's complement
		return -int64((^n & ((1 << nBits) - 1)) + 1)
	}
	return int64(n)
}

// decode zigzag encoded integer
// https://developers.google.com/protocol-buffers/docs/encoding
func ZigZag(n uint64) int64 {
	return int64(n>>1 ^ -(n & 1))
}
