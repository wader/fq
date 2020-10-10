// named rangex because of range is keyword
package rangex

import (
	"fmt"
	"fq/pkg/decode"
	"io"
)

var FieldOutput = &decode.FieldOutput{
	Name: "range",
	New:  func(f *decode.Field) decode.FieldWriter { return &FieldWriter{f: f} },
}

type FieldWriter struct {
	f *decode.Field
}

func (o *FieldWriter) Write(w io.Writer) error {
	o.f.WalkValues(func(v decode.Value) {
		start := v.Range.Start
		stop := v.Range.Stop
		w.Write([]byte(fmt.Sprintf("%d %d\n", start, stop)))

	})

	return nil
}
