package decode

import "io"

func MustCopy(r io.Writer, w io.Reader) int64 {
	n, err := io.Copy(r, w)
	if err != nil {
		panic(err)
	}
	return n
}
