package mathx

import "math/big"

var BigIntOne = big.NewInt(1)

func BigIntSetBytesSigned(n *big.Int, buf []byte) *big.Int {
	n.SetBytes(buf)
	if len(buf) > 0 && buf[0]&0x80 > 0 {
		n.Sub(n, new(big.Int).Lsh(BigIntOne, uint(len(buf))*8))
	}
	return n
}
