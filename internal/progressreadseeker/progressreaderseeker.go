package progressreadseeker

// // TODO: move
// prs := &progressReaderSeeker{RS: rs, Length: bEnd, ProgressFn: func(pos, length int64) {
// 	fmt.Fprintf(os.Stderr, " %.1f%%\r", float64(pos*100)/float64(length))
// }}

// prs := newProgressReaderSeeker2(rs, bEnd, func(readBytes int64, length int64) {
// 	fmt.Fprintf(os.Stderr, " %.1f%%\r", (float64(readBytes)/float64(length))*100)
// })

import "io"

const progressPrecision = 1024

type progressReaderSeeker struct {
	rs             io.ReadSeeker
	pos            int64
	length         int64
	partitions     []bool
	partitionsRead int64
	progressFn     func(readBytes int64, length int64)
}

func newProgressReaderSeeker(rs io.ReadSeeker, length int64, fn func(pos int64, length int64)) *progressReaderSeeker {
	return &progressReaderSeeker{
		rs:         rs,
		length:     length,
		partitions: make([]bool, progressPrecision+1),
		progressFn: fn,
	}
}

func (prs *progressReaderSeeker) Read(p []byte) (n int, err error) {
	n, err = prs.rs.Read(p)
	newPos := prs.pos + int64(n)
	lastPartitionsRead := prs.partitionsRead
	partitionSize := prs.length / (progressPrecision - 1)

	partStart := prs.pos / partitionSize
	partEnd := newPos / partitionSize
	for i := partStart; i <= partEnd; i++ {
		if prs.partitions[i] {
			continue
		}
		prs.partitions[i] = true
		prs.partitionsRead++
	}

	if lastPartitionsRead != prs.partitionsRead {
		readBytes := partitionSize * prs.partitionsRead
		if readBytes > prs.length {
			readBytes = prs.length
		}
		prs.progressFn(readBytes, prs.length)
	}

	prs.pos = newPos

	return n, err
}
