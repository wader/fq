package value

import (
	"fmt"
	"fq/pkg/decode"
	"io"
)

var FieldOutput = &decode.FieldOutput{
	Name: "value",
	New:  func(v interface{}) decode.FieldWriter { return &FieldWriter{v: v} },
}

type FieldWriter struct {
	v interface{}
}

func (o *FieldWriter) write(w io.Writer, v interface{}) error {
	switch v := o.v.(type) {
	case *decode.Field:
		return v.WalkValues(func(v decode.Value) error {
			if _, err := w.Write([]byte(v.RawString())); err != nil {
				return err
			}
			if _, err := w.Write([]byte("\n")); err != nil {
				return err
			}
			return nil
		})
	case []*decode.Field:
		for _, f := range v {
			if err := o.write(w, f); err != nil {
				return err
			}
		}
	case []decode.Value:
		for _, ve := range v {
			if err := o.write(w, ve); err != nil {
				return err
			}
		}
	case decode.Value:
		_, err := fmt.Fprintf(w, "%s", v.V)
		return err
	default:
		panic("unreachable")
	}

	return nil
}

func (o *FieldWriter) Write(w io.Writer) error {
	return o.write(w, o.v)
}
