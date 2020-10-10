package raw

import (
	"fq/pkg/decode"
	"io"
)

var FieldOutput = &decode.FieldOutput{
	Name: "raw",
	New:  func(f *decode.Field) decode.FieldWriter { return &FieldWriter{f: f} },
}

type FieldWriter struct {
	f *decode.Field
}

func (o *FieldWriter) Write(w io.Writer) error {
	// TODO: not byte aligned? pad with zeros
	// TODO: BytesRange version with padding?

	var err error
	o.f.WalkValues(func(v decode.Value) {
		if err != nil {
			// TODO: return false to stop walk? return err?
			return
		}
		_, err = io.Copy(w, v.BitBuf.Copy())
	})

	return err
}
