package bitio

// IOReadSeeker is a io.ReadSeeker that reads from a bitio.ReadSeeker.
// Unaligned byte at EOF will be zero bit padded.
type IOReadSeeker struct {
	IOReader
	s    Seeker
	sPos int64
}

// NewIOReadSeeker return a new bitio.IOReadSeeker.
func NewIOReadSeeker(rs ReadSeeker) *IOReadSeeker {
	return &IOReadSeeker{
		IOReader: IOReader{r: rs},
		s:        rs,
	}
}

func (r *IOReadSeeker) Read(p []byte) (n int, err error) {
	n, err = r.IOReader.Read(p)
	r.sPos += int64(n)
	return n, err
}

func (r *IOReadSeeker) Seek(offset int64, whence int) (int64, error) {
	n, err := r.s.SeekBits(offset*8, whence)
	// TODO: reset last error on seek. some nicer way?
	r.IOReader.rErr = nil
	if n != r.sPos {
		r.b.Reset()
		r.sPos = n / 8
	}
	if err != nil {
		return n / 8, err
	}

	return n / 8, err
}

func (r *IOReadSeeker) Unwrap() any {
	return r.IOReader
}
