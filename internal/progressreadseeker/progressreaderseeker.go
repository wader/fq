package progressreadseeker

// // TODO: move
// prs := &progressReaderSeeker{RS: rs, Length: bEnd, ProgressFn: func(pos, length int64) {
// 	fmt.Fprintf(os.Stderr, " %.1f%%\r", float64(pos*100)/float64(length))
// }}

// prs := newProgressReaderSeeker2(rs, bEnd, func(readBytes int64, length int64) {
// 	fmt.Fprintf(os.Stderr, " %.1f%%\r", (float64(readBytes)/float64(length))*100)
// })

import (
	"io"
)

const progressPrecision = 1024

type progressReaderSeeker struct {
	rs                  io.ReadSeeker
	pos                 int64
	length              int64
	partitionSize       int64
	partitions          []bool
	partitionsReadCount int64
	progressFn          func(readBytes int64, length int64)
}

func New(rs io.ReadSeeker, length int64, fn func(pos int64, length int64)) *progressReaderSeeker {
	partitionSize := length / progressPrecision
	if length%progressPrecision != 0 {
		partitionSize++
	}
	return &progressReaderSeeker{
		rs:            rs,
		length:        length,
		partitionSize: partitionSize,
		partitions:    make([]bool, progressPrecision),
		progressFn:    fn,
	}
}

func (prs *progressReaderSeeker) Read(p []byte) (n int, err error) {
	n, err = prs.rs.Read(p)
	newPos := prs.pos + int64(n)
	lastPartitionsReadCount := prs.partitionsReadCount

	partStart := prs.pos / prs.partitionSize
	partEnd := newPos / prs.partitionSize

	// log.Printf("prs.length: len=%d partitionSize=%d %d-%d pos %d->%d\n", prs.length, prs.partitionSize, partStart, partEnd, prs.pos, newPos)

	for i := partStart; i <= partEnd; i++ {
		if prs.partitions[i] {
			continue
		}
		prs.partitions[i] = true
		prs.partitionsReadCount++
	}

	if lastPartitionsReadCount != prs.partitionsReadCount {
		readBytes := prs.partitionSize * prs.partitionsReadCount
		if readBytes > prs.length {
			readBytes = prs.length
		}
		prs.progressFn(readBytes, prs.length)
	}

	prs.pos = newPos

	return n, err
}

func (prs *progressReaderSeeker) Seek(offset int64, whence int) (int64, error) {
	pos, err := prs.rs.Seek(offset, whence)
	prs.pos = pos
	return pos, err
}
