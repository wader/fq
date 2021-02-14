package osenv

import (
	"io"
)

type OS interface {
	Stdin() io.Reader
	Stdout() io.Writer
	Stderr() io.Writer
	Args() []string
	Environ() []string
	// returned io.ReadSeeker can optionally implement io.Closer
	Open(name string) (io.ReadSeeker, error)
}
