package bitio_test

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"testing"

	"github.com/wader/fq/pkg/bitio"
)

func createMultiReader(numSegments int, segmentSize int) (*bitio.MultiReader, []byte, error) {
	data := make([]byte, numSegments*segmentSize)
	for i := range data {
		data[i] = byte(i % 251)
	}

	readers := make([]bitio.ReadAtSeeker, numSegments)
	for i := 0; i < numSegments; i++ {
		start := i * segmentSize
		end := start + segmentSize
		readers[i] = bitio.NewBitReader(data[start:end], -1)
	}

	mr, err := bitio.NewMultiReader(readers...)
	return mr, data, err
}

func TestMultiReaderEmpty(t *testing.T) {
	mr, err := bitio.NewMultiReader()
	if err != nil {
		t.Fatalf("unexpected error creating empty MultiReader: %v", err)
	}

	buf := make([]byte, 1)
	n, err := mr.ReadBitsAt(buf, 8, 0)
	if !errors.Is(err, io.EOF) || n != 0 {
		t.Fatalf("expected EOF and 0 bits, got n=%d, err=%v", n, err)
	}

	n, err = mr.ReadBitsAt(buf, 8, -1)
	if !errors.Is(err, io.EOF) || n != 0 {
		t.Fatalf("expected EOF and 0 bits for negative offset, got n=%d, err=%v", n, err)
	}

	n, err = mr.ReadBits(buf, 8)
	if !errors.Is(err, io.EOF) || n != 0 {
		t.Fatalf("expected EOF and 0 bits, got n=%d, err=%v", n, err)
	}

	pos, err := mr.SeekBits(0, io.SeekStart)
	if err != nil || pos != 0 {
		t.Fatalf("expected pos=0, err=nil, got pos=%d, err=%v", pos, err)
	}

	_, err = mr.SeekBits(1, io.SeekStart)
	if !errors.Is(err, bitio.ErrOffset) {
		t.Fatalf("expected ErrOffset, got %v", err)
	}
}

func TestMultiReaderEmptySegments(t *testing.T) {
	rEmpty1 := bitio.NewBitReader([]byte{}, 0)
	rEmpty2 := bitio.NewBitReader([]byte{}, 0)
	rData1 := bitio.NewBitReader([]byte{0xaa}, 8)
	rEmpty3 := bitio.NewBitReader([]byte{}, 0)
	rData2 := bitio.NewBitReader([]byte{0x55}, 8)
	rEmpty4 := bitio.NewBitReader([]byte{}, 0)

	mr, err := bitio.NewMultiReader(rEmpty1, rEmpty2, rData1, rEmpty3, rData2, rEmpty4)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	buf := make([]byte, 1)
	n, err := mr.ReadBitsAt(buf, 8, 0)
	if err != nil || n != 8 || buf[0] != 0xaa {
		t.Fatalf("read first byte failed: n=%d err=%v buf=%x", n, err, buf[0])
	}

	n, err = mr.ReadBitsAt(buf, 8, 8)
	if err != nil || n != 8 || buf[0] != 0x55 {
		t.Fatalf("read second byte failed: n=%d err=%v buf=%x", n, err, buf[0])
	}

	n, err = mr.ReadBitsAt(buf, 8, 16)
	if !errors.Is(err, io.EOF) || n != 0 {
		t.Fatalf("expected EOF at offset 16, got n=%d err=%v", n, err)
	}
}

func TestMultiReaderRandomAccess(t *testing.T) {
	segmentCounts := []int{1, 2, 3, 5, 16, 64, 256}
	segmentSize := 32

	for _, count := range segmentCounts {
		t.Run(fmt.Sprintf("%d_segments", count), func(t *testing.T) {
			mr, data, err := createMultiReader(count, segmentSize)
			if err != nil {
				t.Fatalf("failed to create MultiReader: %v", err)
			}

			buf := make([]byte, 4)

			rng := rand.New(rand.NewSource(42)) //nolint:gosec
			for i := 0; i < 200; i++ {
				// Pick a bit offset and nBits within a single segment so ReadBitsAt returns nBits
				segIdx := rng.Intn(count)
				segStart := int64(segIdx * segmentSize * 8)
				segBits := int64(segmentSize * 8)

				nBits := int64(rng.Intn(32) + 1)
				offsetInSeg := rng.Int63n(segBits - nBits + 1)
				bitOff := segStart + offsetInSeg

				n, err := mr.ReadBitsAt(buf, nBits, bitOff)
				if err != nil && !errors.Is(err, io.EOF) {
					t.Fatalf("ReadBitsAt failed at bitOff=%d, nBits=%d: %v", bitOff, nBits, err)
				}
				if n != nBits {
					t.Fatalf("expected %d bits, got %d", nBits, n)
				}

				// Verify bits read against SectionReader on whole data
				expectedReader := bitio.NewBitReader(data, -1)
				expectedBuf := make([]byte, 4)
				en, eErr := expectedReader.ReadBitsAt(expectedBuf, nBits, bitOff)
				if eErr != nil && !errors.Is(eErr, io.EOF) {
					t.Fatalf("expectedReader ReadBitsAt failed: %v", eErr)
				}
				if en != nBits {
					t.Fatalf("expectedReader read %d bits, want %d", en, nBits)
				}
				if !bytes.Equal(buf[:(nBits+7)/8], expectedBuf[:(nBits+7)/8]) {
					t.Fatalf("data mismatch at bitOff=%d nBits=%d: got %x, want %x", bitOff, nBits, buf, expectedBuf)
				}
			}

			// Also test full sequential read using ReadBits across all segment boundaries
			_, err = mr.SeekBits(0, io.SeekStart)
			if err != nil {
				t.Fatalf("seek to 0 failed: %v", err)
			}
			seqBuf := make([]byte, 1)
			for b := int64(0); b < int64(len(data)); b++ {
				n, err := mr.ReadBits(seqBuf, 8)
				if err != nil {
					t.Fatalf("sequential ReadBits failed at byte %d: %v", b, err)
				}
				if n != 8 {
					t.Fatalf("sequential ReadBits expected 8 bits, got %d", n)
				}
				if seqBuf[0] != data[b] {
					t.Fatalf("sequential ReadBits byte mismatch at %d: got %02x, want %02x", b, seqBuf[0], data[b])
				}
			}
			// Next read should return EOF
			n, err := mr.ReadBits(seqBuf, 8)
			if !errors.Is(err, io.EOF) || n != 0 {
				t.Fatalf("expected EOF at end of stream, got n=%d err=%v", n, err)
			}
		})
	}
}

func TestMultiReaderSeek(t *testing.T) {
	mr, _, err := createMultiReader(4, 16)
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}

	totalBits := int64(4 * 16 * 8) // 512 bits

	pos, err := mr.SeekBits(100, io.SeekStart)
	if err != nil || pos != 100 {
		t.Fatalf("SeekStart failed: pos=%d, err=%v", pos, err)
	}

	pos, err = mr.SeekBits(50, io.SeekCurrent)
	if err != nil || pos != 150 {
		t.Fatalf("SeekCurrent failed: pos=%d, err=%v", pos, err)
	}

	pos, err = mr.SeekBits(-50, io.SeekEnd)
	if err != nil || pos != totalBits-50 {
		t.Fatalf("SeekEnd failed: pos=%d, err=%v", pos, err)
	}

	// Invalid seeks
	_, err = mr.SeekBits(-1, io.SeekStart)
	if !errors.Is(err, bitio.ErrOffset) {
		t.Fatalf("expected ErrOffset for negative start seek, got %v", err)
	}

	_, err = mr.SeekBits(totalBits+1, io.SeekStart)
	if !errors.Is(err, bitio.ErrOffset) {
		t.Fatalf("expected ErrOffset for past-end start seek, got %v", err)
	}
}

func TestMultiReaderClone(t *testing.T) {
	mr, _, err := createMultiReader(3, 16)
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}

	_, err = mr.SeekBits(64, io.SeekStart)
	if err != nil {
		t.Fatalf("seek failed: %v", err)
	}

	cloned, err := mr.CloneReaderAtSeeker()
	if err != nil {
		t.Fatalf("clone failed: %v", err)
	}

	// Cloned reader should have pos reset to 0
	buf1 := make([]byte, 1)
	buf2 := make([]byte, 1)

	n1, err1 := mr.ReadBits(buf1, 8)
	n2, err2 := cloned.ReadBits(buf2, 8)
	if err1 != nil || n1 != 8 || err2 != nil || n2 != 8 {
		t.Fatalf("reads failed: %v %v", err1, err2)
	}

	// Original was at bit 64 (byte 8), cloned was at bit 0 (byte 0)
	if bytes.Equal(buf1, buf2) {
		t.Fatalf("expected different bytes at bit 64 vs 0, got both %x", buf1)
	}
}

func BenchmarkMultiReaderReadBitsAt(b *testing.B) {
	segmentCounts := []int{1, 2, 4, 8, 16, 64, 1024, 65536}
	segmentBytes := 16

	for _, count := range segmentCounts {
		b.Run(fmt.Sprintf("%d_segments", count), func(b *testing.B) {
			mr, data, err := createMultiReader(count, segmentBytes)
			if err != nil {
				b.Fatalf("failed to create MultiReader: %v", err)
			}

			totalBits := int64(len(data) * 8)
			offsets := make([]int64, 1024)
			rng := rand.New(rand.NewSource(42)) //nolint:gosec
			maxOffset := totalBits - 64
			for i := range offsets {
				if maxOffset > 0 {
					offsets[i] = rng.Int63n(maxOffset)
				} else {
					offsets[i] = 0
				}
			}

			buf := make([]byte, 8)
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				bitOff := offsets[i%len(offsets)]
				_, _ = mr.ReadBitsAt(buf, 64, bitOff)
			}
		})
	}
}
