package ctxreadseeker_test

import (
	"context"
	"fmt"
	"fq/internal/ctxreadseeker"
	"io"
	"reflect"
	"testing"
	"time"
)

type rwcRecorder struct {
	log []string
}

func (r *rwcRecorder) Read(p []byte) (n int, err error) {
	r.log = append(r.log, fmt.Sprintf("read %d", len(p)))
	return 0, err
}

func (r *rwcRecorder) Seek(offset int64, whence int) (n int64, err error) {
	r.log = append(r.log, fmt.Sprintf("seek %d %d", offset, whence))
	return 0, err
}

func (r *rwcRecorder) Close() (err error) {
	r.log = append(r.log, "close")
	return nil
}

func TestNormal(t *testing.T) {
	r := &rwcRecorder{}
	crs := ctxreadseeker.New(context.Background(), r)

	_, _ = crs.Seek(1, io.SeekStart)
	_, _ = crs.Read(make([]byte, 3))
	_ = crs.Close()

	expected := []string{"seek 1 0", "read 3", "close"}
	if !reflect.DeepEqual(r.log, expected) {
		t.Errorf("expected %v, got %v", expected, r.log)
	}
}

func TestCancel(t *testing.T) {
	r := &rwcRecorder{}
	c, cancelFn := context.WithCancel(context.Background())
	crs := ctxreadseeker.New(c, r)
	_, _ = crs.Seek(1, io.SeekStart)
	_, _ = crs.Read(make([]byte, 3))
	cancelFn()
	time.Sleep(10 * time.Millisecond)

	cancelActualReadN, cancelActualReadErr := crs.Read(make([]byte, 3))

	expected := []string{"seek 1 0", "read 3", "close"}
	if !reflect.DeepEqual(r.log, expected) {
		t.Errorf("expected %v, got %v", expected, r.log)
	}

	if cancelActualReadN != 0 {
		t.Errorf("expected 0, got %v", cancelActualReadN)
	}

	if cancelActualReadErr != context.Canceled {
		t.Errorf("expected cancel err, got %v", cancelActualReadErr)
	}

}
