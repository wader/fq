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

// NopCloseReadSeeker convinced type if using something not requiring close (http source etc)
type NopCloseReadSeeker struct{ io.ReadSeeker }

func (NopCloseReadSeeker) Close() error { return nil }

type OS interface {
	Stdin() io.Reader
	Stdout() io.Writer
	Stderr() io.Writer
	Args() []string
	Open(name string) (ReadSeekCloser, error)
}

type StandardOS struct{}

func (StandardOS) Stdin() io.Reader                         { return os.Stdin }
func (StandardOS) Stdout() io.Writer                        { return os.Stdout }
func (StandardOS) Stderr() io.Writer                        { return os.Stderr }
func (StandardOS) Args() []string                           { return os.Args }
func (StandardOS) Open(name string) (ReadSeekCloser, error) { return os.Open(name) }

func StandardOSMain(formatsList ...[]*decode.Format) {
	if err := (Main{
		OS:          StandardOS{},
		FormatsList: formatsList,
	}).Run(); err != nil {
		os.Exit(1)
	}
}
