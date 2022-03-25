package http

import (
	"github.com/wader/fq/format"
	"github.com/wader/fq/internal/lazyre"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
	"github.com/wader/fq/pkg/scalar"
)

func init() {
	interp.RegisterFormat(
		format.TextProto,
		&decode.Format{
			Description: "Generic text-based protocol (HTTP,SMTP-like)",
			RootArray:   true,
			DecodeFn:    decodeTextProto,
			DefaultInArg: format.TextProto_In{
				Name: "pair",
			},
		})
}

// TODO: line folding correct?
// TODO: move to decode? also make d.FieldArray/Struct return T?

var textprotoLineRE = &lazyre.RE{S: `` +
	(`(?P<name>[\w-]+:)`) +
	(`(?P<value>` +
		`\s*` + // eagerly skip leading whitespace
		`(?:` +
		`.*?(?:\r?\n[\t ].*?)*` +
		`)` +
		`\r?\n` +
		`)` +
		``)}

func decodeTextProto(d *decode.D) any {
	var tpi format.TextProto_In
	d.ArgAs(&tpi)

	m := map[string][]string{}

	for !d.End() {
		c := d.PeekBytes(1)[0]
		if c == '\n' || c == '\r' {
			break
		}

		d.FieldStruct(tpi.Name, func(d *decode.D) {
			cm := map[string]string{}
			// TODO: don't strip :?
			d.FieldRE(textprotoLineRE.Must(), &cm, scalar.StrActualTrim(" :\r\n"))
			name := cm["name"]
			value := cm["value"]
			m[name] = append(m[name], value)
		})
	}

	return format.TextProto_Out{Pairs: m}
}
