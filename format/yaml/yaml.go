package yaml

// TODO: yaml type eval? walk eval?

import (
	"bytes"
	"embed"
	"errors"
	"io"

	"github.com/wader/fq/format"
	"github.com/wader/fq/internal/gojqx"
	"github.com/wader/fq/pkg/bitio"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
	"github.com/wader/fq/pkg/scalar"
	"gopkg.in/yaml.v3"
)

//go:embed yaml.jq
var yamlFS embed.FS

func init() {
	interp.RegisterFormat(
		format.YAML,
		&decode.Format{
			Description: "YAML Ain't Markup Language",
			ProbeOrder:  format.ProbeOrderTextFuzzy,
			Groups:      []*decode.Group{format.Probe},
			DecodeFn:    decodeYAML,
			Functions:   []string{"_todisplay"},
		})
	interp.RegisterFS(yamlFS)
	interp.RegisterFunc1("_to_yaml", toYAML)
}

func decodeYAML(d *decode.D) any {
	br := d.RawLen(d.Len())
	var r any

	yd := yaml.NewDecoder(bitio.NewIOReader(br))
	if err := yd.Decode(&r); err != nil {
		d.Fatalf("%s", err)
	}
	if err := yd.Decode(new(any)); !errors.Is(err, io.EOF) {
		d.Fatalf("trialing data after top-level value")
	}

	var s scalar.Any
	s.Actual = gojqx.Normalize(r)

	switch s.Actual.(type) {
	case map[string]any,
		[]any:
	default:
		d.Fatalf("root not object or array")
	}

	d.Value.V = &s
	d.Value.Range.Len = d.Len()

	return nil
}

type ToYAMLOpts struct {
	Indent int `default:"4"` // 4 is default for gopkg.in/yaml.v3
}

func toYAML(_ *interp.Interp, c any, opts ToYAMLOpts) any {
	b := &bytes.Buffer{}
	e := yaml.NewEncoder(b)
	// yaml.SetIndent panics if < 0
	if opts.Indent >= 0 {
		e.SetIndent(opts.Indent)
	}
	if err := e.Encode(gojqx.Normalize(c)); err != nil {
		return err
	}

	return b.String()
}
