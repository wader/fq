package raw

import (
	"fmt"
	"fq/pkg/decode"
	"io"
)

var FieldOutput = &decode.FieldOutput{
	Name: "raw",
	New:  func(v interface{}) decode.FieldWriter { return &FieldWriter{v: v} },
}

type FieldWriter struct {
	v interface{}
}

func (o *FieldWriter) Write(w io.Writer) error {
	// TODO: not byte aligned? pad with zeros
	// TODO: BytesRange version with padding?

	switch v := o.v.(type) {
	case *decode.Field:
		v.WalkValues(func(v decode.Value) error {
			_, err := io.Copy(w, v.BitBuf.Copy())
			return err
		})
	case []*decode.Field:
		for _, f := range v {
			f.WalkValues(func(v decode.Value) error {
				_, err := io.Copy(w, v.BitBuf.Copy())
				return err
			})
		}
	case []decode.Value:
		for _, ve := range v {
			_, err := io.Copy(w, ve.BitBuf.Copy())
			return err
		}
	case decode.Value:
		_, err := io.Copy(w, v.BitBuf.Copy())
		return err
	default:
		_, err := fmt.Fprintf(w, "%s", v)
		return err
	}

	return nil

}
