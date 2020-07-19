package cli

import (
	"fq/pkg/decode"
	"io"
	"os"
)

type OS interface {
	Stdin() io.Reader
	Stdout() io.Writer
	Stderr() io.Writer
	Args() []string
	Open(name string) (*os.File, error)
}

type StandardOS struct{}

func (StandardOS) Stdin() io.Reader                   { return os.Stdin }
func (StandardOS) Stdout() io.Writer                  { return os.Stdout }
func (StandardOS) Stderr() io.Writer                  { return os.Stderr }
func (StandardOS) Args() []string                     { return os.Args }
func (StandardOS) Open(name string) (*os.File, error) { return os.Open(name) }

func StandardOSMain(formatsList ...[]*decode.Format) {
	if err := (Main{
		OS:          StandardOS{},
		FormatsList: formatsList,
	}).Run(); err != nil {
		os.Exit(1)
	}
}
