package json

import (
	"bytes"
	"embed"

	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
)

//go:embed jsonl.jq
var jsonlFS embed.FS

// TODO: not strictly JSONL, allows any whitespace between

func init() {
	interp.RegisterFormat(decode.Format{
		Name:        format.JSONL,
		Description: "JavaScript Object Notation Lines",
		ProbeOrder:  format.ProbeOrderTextFuzzy,
		Groups:      []string{format.PROBE},
		DecodeFn:    decodeJSONL,
		Functions:   []string{"_todisplay"},
	})
	interp.RegisterFS(jsonlFS)
	interp.RegisterFunc0("tojsonl", toJSONL)
}

func decodeJSONL(d *decode.D, _ any) any {
	return decodeJSONEx(d, true)
}

func toJSONL(i *interp.Interp, c []any) any {
	cj := makeEncoder(ToJSONOpts{})
	bb := &bytes.Buffer{}

	for _, v := range c {
		if err := cj.Marshal(v, bb); err != nil {
			return err
		}
		bb.WriteByte('\n')
	}

	return bb.String()
}
