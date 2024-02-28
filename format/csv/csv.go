package csv

// TODO: error, throw error always? no decode value with gap etc? -d csv from_csv
// TODO: header row field count mismatch error, csv reader takes care of check atm. can use FieldsPerRecord -1
// TODO: row object keys mismatch writer
// TODO: lazy quotes?
// TODO: comment in writer? string elements?
// TODO: to_csv objects
// TODO: to_csv opts help
// TODO: go maps are random order, now sorts headers
// TODO: option aliases?
// TODO: snake_case option?

import (
	"bytes"
	"embed"
	"errors"
	"fmt"
	"io"
	"sort"

	"github.com/wader/fq/format"
	"github.com/wader/fq/internal/csvex"
	"github.com/wader/fq/internal/gojqex"
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
				Delimiter:        ",",
				Comment:          "",
				QuoteChar:        `"`,
				Header:           true,
				SkipInitialSpace: false,
			},
			Functions: []string{"_todisplay"},
		})
	interp.RegisterFS(csvFS)
	interp.RegisterFunc1("_to_csv", toCSV)
}

func decodeCSV(d *decode.D) any {
	var ci format.CSV_In
	d.ArgAs(&ci)

	br := d.RawLen(d.Len())
	r := csvex.NewReader(bitio.NewIOReader(br))
	r.LazyQuotes = true
	if ci.Delimiter != "" {
		r.Comma = rune(ci.Delimiter[0])
	} else if ci.Comma != "" {
		r.Comma = rune(ci.Comma[0])
	}
	if ci.Comment != "" {
		r.Comment = rune(ci.Comment[0])
	} else {
		r.Comment = 0
	}
	if ci.QuoteChar != "" {
		r.Quote = rune(ci.QuoteChar[0])
	} else {
		r.Quote = '"'
	}
	r.TrimLeadingSpace = ci.SkipInitialSpace

	row := 1
	var rvs []any

	var headers []string
	for {
		r, err := r.Read()
		if errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			return err
		}

		if ci.Header {
			if headers == nil {
				// TODO: duplicate headers?
				headers = append(headers, r...)
			} else {
				obj := map[string]any{}
				for i, s := range r {
					h := headers[i]
					obj[h] = s
				}
				rvs = append(rvs, obj)
			}
		} else {
			var vs []any
			for _, s := range r {
				vs = append(vs, s)
			}
			rvs = append(rvs, vs)
		}

		row++
	}

	d.Value.V = &scalar.Any{Actual: rvs}
	d.Value.Range.Len = d.Len()

	return nil
}

type ToCSVOpts struct {
	Comma     string // alias for Delimiter
	Delimiter string
	QuoteChar string
	Header    bool
}

func toCSV(_ *interp.Interp, c []any, opts ToCSVOpts) any {
	b := &bytes.Buffer{}
	w := csvex.NewWriter(b)
	if opts.Delimiter != "" {
		w.Comma = rune(opts.Delimiter[0])
	} else if opts.Comma != "" {
		w.Comma = rune(opts.Comma[0])
	}
	if opts.QuoteChar != "" {
		w.Quote = rune(opts.QuoteChar[0])
	} else {
		w.Quote = '"'
	}

	seenObject := 0
	seenArrays := 0
	var headers []string

	for _, row := range c {
		switch row.(type) {
		case []any:
			if seenObject > 0 {
				return fmt.Errorf("mixed row types, expected row to be an object, got %s", gojqex.TypeErrorPreview(row))
			}

			rs, ok := gojqex.Cast[[]any](row)
			if !ok {
				return fmt.Errorf("expected row to be an array, got %s", gojqex.TypeErrorPreview(row))
			}
			vs, ok := gojqex.NormalizeToStrings(rs).([]any)
			if !ok {
				panic("not array")
			}
			var ss []string
			for _, v := range vs {
				s, ok := v.(string)
				if !ok {
					return fmt.Errorf("expected row record to be scalars, got %s", gojqex.TypeErrorPreview(v))
				}
				ss = append(ss, s)
			}
			if err := w.Write(ss); err != nil {
				return err
			}

			seenArrays++
		case map[string]any:
			if seenArrays > 0 {
				return fmt.Errorf("mixed row types, expected row to be an array, got %s", gojqex.TypeErrorPreview(row))
			}

			rm, ok := gojqex.Cast[map[string]any](row)
			if !ok {
				return fmt.Errorf("expected row to be an object, got %s", gojqex.TypeErrorPreview(row))
			}
			vm, ok := gojqex.NormalizeToStrings(rm).(map[string]any)
			if !ok {
				panic("not object")
			}

			if headers == nil {
				// TODO: maps are random order in go
				for k := range vm {
					headers = append(headers, k)
				}
				sort.Strings(headers)

				if err := w.Write(headers); err != nil {
					return err
				}
			}

			var ss []string
			keysFound := 0
			for _, k := range headers {
				s, ok := vm[k].(string)
				if !ok {
					return fmt.Errorf("expected row object to have a %q key, %s", k, gojqex.TypeErrorPreview(row))
				}
				ss = append(ss, s)
				keysFound++
			}
			// TODO: what keys are extra/missing
			if keysFound < len(headers) {
				return fmt.Errorf("expected row object has missing keys %s", gojqex.TypeErrorPreview(row))
			} else if keysFound > len(headers) {
				return fmt.Errorf("expected row object has extra keys %s", gojqex.TypeErrorPreview(row))
			}

			if err := w.Write(ss); err != nil {
				return err
			}

			seenObject++
		}

	}
	w.Flush()

	return b.String()
}
