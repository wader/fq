package bitio

func ReverseBytes(bs []byte) {
	l := len(bs)
	for i := 0; i < l/2; i++ {
		bs[i], bs[l-i-1] = bs[l-i-1], bs[i]
	}
}
