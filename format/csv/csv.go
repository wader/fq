package csv

import (
	"bytes"
	"embed"
	"encoding/csv"
	"errors"
	"fmt"
	"io"

	"github.com/wader/fq/format"
	"github.com/wader/fq/internal/gojqx"
	"github.com/wader/fq/pkg/bitio"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
	"github.com/wader/fq/pkg/scalar"
)

//go:embed csv.jq
//go:embed csv.md
var csvFS embed.FS

func init() {
	interp.RegisterFormat(
		format.CSV,
		&decode.Format{
			Description: "Comma separated values",
			ProbeOrder:  format.ProbeOrderTextFuzzy,
			DecodeFn:    decodeCSV,
			DefaultInArg: format.CSV_In{
				Comma:   ",",
				Comment: "#",
			},
			Functions: []string{"_todisplay"},
		})
	interp.RegisterFS(csvFS)
	interp.RegisterFunc1("_to_csv", toCSV)
}

func decodeCSV(d *decode.D) any {
	var ci format.CSV_In
	d.ArgAs(&ci)

	var rvs []any
	br := d.RawLen(d.Len())
	r := csv.NewReader(bitio.NewIOReader(br))
	r.TrimLeadingSpace = true
	r.LazyQuotes = true
	if ci.Comma != "" {
		r.Comma = rune(ci.Comma[0])
	}
	if ci.Comment != "" {
		r.Comment = rune(ci.Comment[0])
	}
	for {
		r, err := r.Read()
		if errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			return err
		}
		var vs []any
		for _, s := range r {
			vs = append(vs, s)
		}
		rvs = append(rvs, vs)
	}

	d.Value.V = &scalar.Any{Actual: rvs}
	d.Value.Range.Len = d.Len()

	return nil
}

type ToCSVOpts struct {
	Comma string
}

func toCSV(_ *interp.Interp, c []any, opts ToCSVOpts) any {
	b := &bytes.Buffer{}
	w := csv.NewWriter(b)
	if opts.Comma != "" {
		w.Comma = rune(opts.Comma[0])
	}
	for _, row := range c {
		rs, ok := gojqx.Cast[[]any](row)
		if !ok {
			return fmt.Errorf("expected row to be an array, got %s", gojqx.TypeErrorPreview(row))
		}
		vs, ok := gojqx.NormalizeToStrings(rs).([]any)
		if !ok {
			panic("not array")
		}
		var ss []string
		for _, v := range vs {
			s, ok := v.(string)
			if !ok {
				return fmt.Errorf("expected row record to be scalars, got %s", gojqx.TypeErrorPreview(v))
			}
			ss = append(ss, s)
		}
		if err := w.Write(ss); err != nil {
			return err
		}
	}
	w.Flush()

	return b.String()
}
