//go:build exclude

package main

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"text/template"
)

func toInt(v any) int {
	switch v := v.(type) {
	case float64:
		return int(v)
	case int:
		return v
	default:
		return 0
	}
}

func toString(v any) string {
	switch v := v.(type) {
	case string:
		return v
	case int:
		return strconv.Itoa(v)
	default:
		return ""
	}
}

func main() {
	funcMap := template.FuncMap{
		"xrange": func(args ...any) (any, error) {
			if len(args) < 2 {
				return nil, errors.New("need min and max argument")
			}

			min := toInt(args[0])
			max := toInt(args[1])
			var v []int
			for i := min; i < max; i++ {
				v = append(v, i)
			}

			return v, nil
		},
		"replace": func(args ...any) (any, error) {
			if len(args) < 3 {
				return nil, errors.New("need tmpl, old and new argument")
			}

			s := toString(args[0])
			o := toString(args[1])
			n := toString(args[2])

			return strings.Replace(s, o, n, -1), nil
		},
		"slice": func(args ...any) []any {
			return args
		},
		"map": func(args ...any) map[any]any {
			m := map[any]any{}
			for i := 0; i < len(args)/2; i++ {
				m[args[i*2]] = args[i*2+1]
			}
			return m
		},
	}

	data := map[string]any{}
	if len(os.Args) > 1 {
		r, err := os.Open(os.Args[1])
		if err != nil {
			log.Fatalf("%s: %s", os.Args[1], err)
		}
		defer r.Close()
		_ = json.NewDecoder(r).Decode(&data)
	}

	templateBytes, err := io.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}
	templateStr := string(templateBytes)

	tmpl, err := template.New("").Funcs(funcMap).Parse(templateStr)
	if err != nil {
		log.Fatalf("template.New: %s", err)
	}

	err = tmpl.Execute(os.Stdout, &data)
	if err != nil {
		log.Fatalf("tmpl.Execute: %s", err)
	}
}
