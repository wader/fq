package raw

import (
	"fq/pkg/decode"
	"io"
)

var FieldOutput = &decode.FieldOutput{
	Name: "raw",
	New:  func(v *decode.Value) decode.FieldWriter { return &FieldWriter{v: v} },
}

type FieldWriter struct {
	v *decode.Value
}

func (o *FieldWriter) Write(w io.Writer) error {
	// TODO: not byte aligned? pad with zeros
	// TODO: BytesRange version with padding?

	o.v.WalkPostOrder(func(v *decode.Value, depth int, rootDepth int) error {
		bb, _ := v.BitBuf.Copy()
		_, err := io.Copy(w, bb)
		return err
	})

	return nil

}
