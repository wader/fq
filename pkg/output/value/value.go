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
	o.f.WalkValues(func(v decode.Value) {
		w.Write([]byte(v.RawString()))
		w.Write([]byte("\n"))
	})
	return nil
}
