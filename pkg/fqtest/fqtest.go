package fqtest

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/wader/fq/format/registry"
	"github.com/wader/fq/internal/deepequal"
	"github.com/wader/fq/internal/script"
	"github.com/wader/fq/pkg/interp"
)

var writeActual = os.Getenv("WRITE_ACTUAL") != ""

func testDecodedTestCaseRun(t *testing.T, registry *registry.Registry, cr *script.CaseRun) {
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

	if writeActual {
		return
	}

	deepequal.Error(t, "exitcode", cr.ExpectedExitCode, cr.ActualExitCode)
	deepequal.Error(t, "stdout", cr.ToExpectedStdout(), cr.ActualStdoutBuf.String())
	deepequal.Error(t, "stderr", cr.ToExpectedStderr(), cr.ActualStderrBuf.String())
}

func TestPath(t *testing.T, registry *registry.Registry) {
	cs := []*script.Case{}

	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) != ".fqtest" {
			return nil
		}

		t.Run(path, func(t *testing.T) {
			b, err := ioutil.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			c := script.ParseCases(string(b))

			cs = append(cs, c)
			c.Path = path

			for _, p := range c.Parts {
				cr, ok := p.(*script.CaseRun)
				if !ok {
					continue
				}

				t.Run(strconv.Itoa(cr.LineNr)+":"+cr.Command, func(t *testing.T) {
					testDecodedTestCaseRun(t, registry, cr)
					c.WasRun = true
				})
			}
		})

		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	if writeActual {
		for _, c := range cs {
			if !c.WasRun {
				continue
			}
			if err := ioutil.WriteFile(c.Path, []byte(c.ToActual()), 0644); err != nil { //nolint:gosec
				t.Error(err)
			}
		}
	}
}
