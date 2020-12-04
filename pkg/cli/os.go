package cli

import (
	"fq/pkg/decode"
	"io"
	"os"
)

type ReadSeekCloser interface {
	io.ReadSeeker
	io.Closer
}

type OS interface {
	Stdin() io.Reader
	Stdout() io.Writer
	Stderr() io.Writer
	Args() []string
	// returned io.ReadSeeker can optionally implement io.Closer
	Open(name string) (io.ReadSeeker, error)
}

type StandardOS struct{}

func (StandardOS) Stdin() io.Reader                        { return os.Stdin }
func (StandardOS) Stdout() io.Writer                       { return os.Stdout }
func (StandardOS) Stderr() io.Writer                       { return os.Stderr }
func (StandardOS) Args() []string                          { return os.Args }
func (StandardOS) Open(name string) (io.ReadSeeker, error) { return os.Open(name) }

func StandardOSMain(r *decode.Registry) {
	if err := (Main{
		OS:       StandardOS{},
		Registry: r,
	}).Run(); err != nil {
		os.Exit(1)
	}
}
