package ioextra

import (
	"errors"
	"io"
)

func SeekerEnd(s io.Seeker) (int64, error) {
	cPos, err := s.Seek(0, io.SeekCurrent)
	if err != nil {
		return 0, err
	}
	epos, err := s.Seek(0, io.SeekEnd)
	if err != nil {
		return 0, err
	}
	if _, err := s.Seek(cPos, io.SeekStart); err != nil {
		return 0, err
	}

	return epos, nil
}

type ReadErrSeeker struct{ io.Reader }

func (r *ReadErrSeeker) Seek(offset int64, whence int) (int64, error) {
	return 0, errors.New("seek")
}
