package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/wader/fq/format/kaitai/ksexpr"
)

// used to make json numbers into int/float64
func normalizeNumbers(a any) any {
	switch a := a.(type) {
	case map[string]any:
		for k, v := range a {
			a[k] = normalizeNumbers(v)
		}
		return a
	case []any:
		for k, v := range a {
			a[k] = normalizeNumbers(v)
		}
		return a
	case json.Number:
		if strings.Contains(a.String(), ".") {
			f, _ := a.Float64()
			return f
		}
		// TODO: truncates
		i, _ := a.Int64()
		return int(i)
	default:
		return a
	}
}

type JSONVar struct {
	V any
}

func (v *JSONVar) String() string {
	if v.V != nil {
		b, _ := json.Marshal(v.V)
		return string(b)
	}
	return ""
}

func (v *JSONVar) Set(s string) error {
	jd := json.NewDecoder(bytes.NewBufferString(s))
	jd.UseNumber()
	if err := jd.Decode(&v.V); err != nil {
		return err
	}
	v.V = normalizeNumbers(v.V)
	return nil
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Evaluate kaitai expression\n")
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [FLAGS] [EXPR]\n", os.Args[0])
		flag.PrintDefaults()
	}
	lexFlag := flag.Bool("lex", false, "Lex expression")
	parseFlag := flag.Bool("parse", false, "Parse expression")
	var inputValue JSONVar
	flag.Var(&inputValue, "input", "Input JSON")

	flag.Parse()
	exprStr := flag.Arg(0)

	if exprStr == "" {
		flag.Usage()
		os.Exit(1)
	}

	if *lexFlag {
		for _, t := range ksexpr.Lex(exprStr) {
			fmt.Fprintf(
				os.Stderr,
				"%s %s (%d-%d) %v\n",
				t.Name, t.Token.Str, t.Token.Span.Start, t.Token.Span.Stop, t.Err,
			)
		}
		return
	} else {
		r, err := ksexpr.Parse(exprStr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "parse: %s\n", err)
			os.Exit(1)
		}

		if *parseFlag {
			je := json.NewEncoder(os.Stderr)
			je.SetIndent("", "  ")
			if err := je.Encode(&r); err != nil {
				fmt.Fprintf(os.Stderr, "%s", err)
				os.Exit(1)
			}
			return
		} else {
			v, err := r.Eval(ksexpr.ToValue(inputValue.V))
			if err != nil {
				fmt.Fprintf(os.Stderr, "eval: %s\n", err)
				os.Exit(1)
			}
			fmt.Fprintf(os.Stdout, "%#v\n", v)
		}
	}
}
