package fqtest

import (
	"context"
	"io/ioutil"
	"os"
	"strconv"
	"testing"

	"github.com/wader/fq/internal/difftest"
	"github.com/wader/fq/internal/script"
	"github.com/wader/fq/pkg/interp"
	"github.com/wader/fq/pkg/registry"
)

func TestPath(t *testing.T, registry *registry.Registry) {
	difftest.TestWithOptions(t, difftest.Options{
		Path:        ".",
		Pattern:     "*.fqtest",
		ColorDiff:   os.Getenv("DIFF_COLOR") != "",
		WriteOutput: os.Getenv("WRITE_ACTUAL") != "",
		Fn: func(t *testing.T, path, input string) (string, string, error) {
			b, err := ioutil.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			c := script.ParseCases(string(b))
			c.Path = path

			for _, p := range c.Parts {
				cr, ok := p.(*script.CaseRun)
				if !ok {
					continue
				}

				t.Run(strconv.Itoa(cr.LineNr)+"/"+cr.Command, func(t *testing.T) {
					cr.WasRun = true

					i, err := interp.New(cr, registry)
					if err != nil {
						t.Fatal(err)
					}

					err = i.Main(context.Background(), cr.Stdout(), "dev")
					if err != nil {
						if ex, ok := err.(interp.Exiter); ok { //nolint:errorlint
							cr.ActualExitCode = ex.ExitCode()
						}
					}
				})
			}

			return path, c.ToActual(), nil
		},
	})
}
