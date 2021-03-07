package ioextra

import (
	"io"
)

func MustCopy(r io.Writer, w io.Reader) int64 {
	n, err := io.Copy(r, w)
	if err != nil {
		panic(err)
	}
	return n
}

func SeekerEnd(s io.Seeker) (int64, error) {
	cpos, err := s.Seek(0, io.SeekCurrent)
	if err != nil {
		return 0, err
	}
	epos, err := s.Seek(0, io.SeekEnd)
	if err != nil {
		return 0, err
	}
	if _, err := s.Seek(cpos, io.SeekStart); err != nil {
		return 0, err
	}

	return epos, nil
}
