package value

import (
	"fq/pkg/decode"
	"io"
)

var FieldOutput = &decode.FieldOutput{
	Name: "value",
	New:  func(f *decode.Field) decode.FieldWriter { return &FieldWriter{f: f} },
}

type FieldWriter struct {
	f *decode.Field
}

func (o *FieldWriter) Write(w io.Writer) error {
	w.Write([]byte(o.f.Value.RawString()))
	w.Write([]byte("\n"))
	return nil
}
