package yaml

// TODO: yaml type eval? walk eval?

import (
	"embed"
	"errors"
	"io"

	"github.com/wader/fq/format"
	"github.com/wader/fq/internal/gojqex"
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
		format.Yaml,
		&decode.Format{
			Description: "YAML Ain't Markup Language",
			ProbeOrder:  format.ProbeOrderTextFuzzy,
			Groups:      []*decode.Group{format.Probe},
			DecodeFn:    decodeYAML,
			Functions:   []string{"_todisplay"},
		})
	interp.RegisterFS(yamlFS)
	interp.RegisterFunc0("to_yaml", toYAML)
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
	s.Actual = gojqex.Normalize(r)

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

func toYAML(_ *interp.Interp, c any) any {
	b, err := yaml.Marshal(gojqex.Normalize(c))
	if err != nil {
		return err
	}
	return string(b)
}
