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
	Readline(prompt string, complete func(line string, pos int) (newLine []string, shared int)) (string, error)
}
