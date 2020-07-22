package decode

import "io"

// TODO: move? better names
type FieldWriter interface {
	Write(w io.Writer) error
}

type FieldOutput struct {
	Name string
	New  func(f *Field) FieldWriter
}
