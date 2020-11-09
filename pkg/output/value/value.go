package value

import (
	"fq/pkg/decode"
	"io"
)

var FieldOutput = &decode.FieldOutput{
	Name: "value",
	New:  func(v *decode.Value) decode.FieldWriter { return &FieldWriter{v: v} },
}

type FieldWriter struct {
	v *decode.Value
}

func (o *FieldWriter) write(w io.Writer, v interface{}) error {
	return o.v.WalkPreOrder(func(v *decode.Value, depth int, rootDepth int) error {
		if _, err := w.Write([]byte(v.RawString())); err != nil {
			return err
		}
		if _, err := w.Write([]byte("\n")); err != nil {
			return err
		}
		return nil
	})
}

func (o *FieldWriter) Write(w io.Writer) error {
	return o.write(w, o.v)
}
