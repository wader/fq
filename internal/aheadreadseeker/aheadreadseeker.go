package aheadreadseeker

import (
	"io"
)

// TODO: smarter cache? cache behind too somehow?

func min64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

type Reader struct {
	rs      io.ReadSeeker
	minRead int

	off       int64
	cache     []byte
	cacheOff  int64
	cacheUsed int
}

func New(rs io.ReadSeeker, minRead int) *Reader {
	return &Reader{
		rs:      rs,
		minRead: minRead,
	}
}

func (r *Reader) Read(p []byte) (n int, err error) {
	for {
		if r.off >= r.cacheOff && r.off < r.cacheOff+int64(r.cacheUsed) {
			d := r.off - r.cacheOff
			copyLen := min64(int64(r.cacheUsed)-d, int64(len(p)))
			copy(p[0:copyLen], r.cache[d:d+copyLen])
			r.off += copyLen

			return int(copyLen), nil
		}

		readBytes := len(p)
		if readBytes < r.minRead {
			readBytes = r.minRead
		}
		if readBytes > len(r.cache) {
			r.cache = make([]byte, readBytes)
		}

		n, err := io.ReadFull(r.rs, r.cache[0:readBytes])
		r.cacheOff = r.off
		r.cacheUsed = n
		if n == 0 || (err != io.ErrUnexpectedEOF && err == io.EOF) {
			return 0, err
		}
	}
}

func (r *Reader) Seek(offset int64, whence int) (int64, error) {
	var absOff int64
	var err error

	switch whence {
	case io.SeekStart:
		absOff = offset
	case io.SeekCurrent:
		absOff = r.off + offset
	case io.SeekEnd:
		absOff, err = r.rs.Seek(offset, whence)
		if err != nil {
			return 0, err
		}
	}

	if absOff >= r.cacheOff && absOff < r.cacheOff+int64(r.cacheUsed) {
		r.off = absOff
		return absOff, nil
	}

	_, err = r.rs.Seek(absOff, io.SeekStart)
	if err != nil {
		return 0, err
	}
	r.off = absOff
	r.cacheOff = 0
	r.cacheUsed = 0

	return absOff, nil
}
