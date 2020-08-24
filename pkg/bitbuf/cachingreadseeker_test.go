package bitbuf_test

import (
	"bytes"
	"fq/pkg/bitbuf"
	"io"
	"log"
	"testing"
)

type readRecord struct {
	pLen int
	off  int64
	retP []byte
	err  error
}

type recordingReadSeeker struct {
	rs      io.ReadSeeker
	off     int64
	records []readRecord
}

func (r *recordingReadSeeker) Read(p []byte) (n int, err error) {
	n, err = r.rs.Read(p)
	retP := make([]byte, n)
	copy(retP, p)
	r.records = append(r.records, readRecord{
		pLen: len(p),
		off:  r.off,
		retP: retP,
		err:  err,
	})
	return n, err

}

func (r *recordingReadSeeker) Seek(offset int64, whence int) (int64, error) {
	off, err := r.rs.Seek(offset, whence)
	r.off = off
	return off, err
}

func TestNewReadAtCache(t *testing.T) {
	rrs := &recordingReadSeeker{rs: bytes.NewReader([]byte("abc"))}
	r := bitbuf.NewCachingReadSeeker(rrs, 2)

	b := make([]byte, 1)
	r.Read(b)
	r.Read(b)
	r.Read(b)

	log.Printf("b: %s\n", b)

	log.Printf("rrs.records: %#+v\n", rrs.records)
}
