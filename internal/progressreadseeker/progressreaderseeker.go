package progressreadseeker

import (
	"io"
)

type ProgressFn func(approxReadBytes int64, totalSize int64)

type ProgressReaderSeeker struct {
	rs                  io.ReadSeeker
	pos                 int64
	totalSize           int64
	partitionSize       int64
	partitions          []bool
	partitionsReadCount int64
	progressFn          ProgressFn
}

func New(rs io.ReadSeeker, precision int64, totalSize int64, fn ProgressFn) *ProgressReaderSeeker {
	partitionSize := totalSize / precision
	if totalSize%precision != 0 {
		partitionSize++
	}
	return &ProgressReaderSeeker{
		rs:            rs,
		totalSize:     totalSize,
		partitionSize: partitionSize,
		partitions:    make([]bool, precision),
		progressFn:    fn,
	}
}

func (prs *ProgressReaderSeeker) Read(p []byte) (n int, err error) {
	n, err = prs.rs.Read(p)
	newPos := prs.pos + int64(n)
	lastPartitionsReadCount := prs.partitionsReadCount

	partStart := prs.pos / prs.partitionSize
	partEnd := newPos / prs.partitionSize

	for i := partStart; i < partEnd; i++ {
		if prs.partitions[i] {
			continue
		}
		prs.partitions[i] = true
		prs.partitionsReadCount++
	}

	if lastPartitionsReadCount != prs.partitionsReadCount {
		readBytes := prs.partitionSize * prs.partitionsReadCount
		if readBytes > prs.totalSize {
			readBytes = prs.totalSize
		}
		prs.progressFn(readBytes, prs.totalSize)
	}

	prs.pos = newPos

	return n, err
}

func (prs *ProgressReaderSeeker) Seek(offset int64, whence int) (int64, error) {
	pos, err := prs.rs.Seek(offset, whence)
	prs.pos = pos
	return pos, err
}
