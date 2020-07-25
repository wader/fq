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
	b, err := o.f.BitBuf().BytesRange(0, o.f.Range.Length()/8)
	if err != nil {
		return err
	}

	w.Write(b)

	return nil
}
