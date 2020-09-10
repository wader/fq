package raw

import (
	"fq/pkg/decode"
	"io"
	"log"
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

	log.Printf("o.f: %#+v\n", o.f)

	log.Printf("o.f.BitBuf().Pos: %#+v\n", o.f.BitBuf().Pos)

	_, err := io.Copy(w, o.f.BitBuf().Copy())
	return err
}
