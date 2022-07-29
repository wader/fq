package common

func TypeAlign(alignVal uint64, alignLen uint64) uint64 {
	return (alignLen + alignVal - 1) & ^(alignVal - 1)
}

func TypeAlign8(alignLen uint64) uint64 {
	return TypeAlign(8, alignLen)
}

func RoundDown(alignVal uint64, alignLen uint64) uint64 {
	return (alignLen / alignVal) * alignVal
}
