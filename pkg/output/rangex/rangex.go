// named rangex because of range is keyword
package rangex

import (
	"fq/pkg/decode"
	"io"
)

var FieldOutput = &decode.FieldOutput{
	Name: "range",
	New:  func(v interface{}) decode.FieldWriter { return &FieldWriter{v: v} },
}

type FieldWriter struct {
	v interface{}
}

func (o *FieldWriter) Write(w io.Writer) error {
	// o.f.WalkValues(func(v decode.Value) {
	// 	start := v.Range.Start
	// 	stop := v.Range.Stop
	// 	w.Write([]byte(fmt.Sprintf("%d %d\n", start, stop)))

	// })

	return nil
}
