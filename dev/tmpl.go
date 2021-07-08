package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"text/template"
)

func toInt(v interface{}) int {
	switch v := v.(type) {
	case float64:
		return int(v)
	case int:
		return v
	default:
		return 0
	}
}

func main() {
	funcMap := template.FuncMap{
		"xrange": func(args ...interface{}) (interface{}, error) {
			if len(args) < 2 {
				return nil, errors.New("need min and max argument")
			}

			min := toInt(args[0])
			max := toInt(args[1])
			var v []int
			for i := int(min); i <= int(max); i++ {
				v = append(v, i)
			}

			return v, nil
		},
	}

	data := map[string]interface{}{}
	if len(os.Args) > 1 {
		r, err := os.Open(os.Args[1])
		if err != nil {
			log.Fatalf("%s: %s", os.Args[1], err)
		}
		defer r.Close()
		json.NewDecoder(r).Decode(&data)
	}

	templateBytes, err := ioutil.ReadAll(os.Stdin)
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
