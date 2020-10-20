package main

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"text/template"
)

func main() {

	// First we create a FuncMap with which to register the function.
	funcMap := template.FuncMap{
		// The name "title" is what the function will be called in the template text.
		"title": strings.Title,
		"xrange": func(args ...interface{}) (interface{}, error) {
			if len(args) < 2 {
				return nil, errors.New("need min and max argument")
			}

			min, minOk := args[0].(int)
			max, maxOk := args[1].(int)
			if !minOk || !maxOk {
				return nil, errors.New("min and max must be int")
			}

			var v []int
			for i := min; i <= max; i++ {
				v = append(v, i)
			}

			return v, nil
		},
		"map": func(args ...interface{}) (interface{}, error) {
			if len(args)%2 != 0 {
				return nil, errors.New("need even number of key value arguments")
			}

			v := map[interface{}]interface{}{}
			for i := 0; i < len(args)/2; i++ {
				v[args[i*2]] = args[i*2+1]
			}

			return v, nil
		},
		"array": func(args ...interface{}) (interface{}, error) {
			return args, nil
		},
	}

	templateText, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}

	// Create a template, add the function map, and parse the text.
	tmpl, err := template.New("titleTest").Funcs(funcMap).Parse(string(templateText))
	if err != nil {
		log.Fatalf("parsing: %s", err)
	}

	// Run the template to verify the output.
	err = tmpl.Execute(os.Stdout, "the go programming language")
	if err != nil {
		log.Fatalf("execution: %s", err)
	}
}
