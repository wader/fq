package mathx

import (
	"cmp"
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

func DigitsInBase[T Integer](n T, basePrefix bool, base int) int {
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

func PadFormatInt[T Signed](i T, base int, basePrefix bool, width int) string {
	return padFormatNumber(strconv.FormatInt(int64(i), base), base, basePrefix, width)
}

func PadFormatUint[T Unsigned](i T, base int, basePrefix bool, width int) string {
	return padFormatNumber(strconv.FormatUint(uint64(i), base), base, basePrefix, width)
}

func PadFormatBigInt(i *big.Int, base int, basePrefix bool, width int) string {
	return padFormatNumber(i.Text(base), base, basePrefix, width)
}

func Clamp[T cmp.Ordered](a, b, v T) T {
	return max(a, min(b, v))
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
	return fmt.Sprintf("%s-%s", Bits(r.Start).StringByteBits(base), Bits(r.Start+r.Len).StringByteBits(base))
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
func ZigZag[U Unsigned, S Signed](n U) S {
	return S(n>>1 ^ -(n & 1))
}
