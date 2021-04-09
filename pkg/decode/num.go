package decode

import "fq/internal/num"

// TODO: move to other package?
func ZigZag(n uint64) int64                    { return num.ZigZag(n) }
func TwosComplement(nBits int, n uint64) int64 { return num.TwosComplement(nBits, n) }
