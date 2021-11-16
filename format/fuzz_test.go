//go:build fuzz

package format_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	_ "github.com/wader/fq/format/all"
	"github.com/wader/fq/format/registry"
	"github.com/wader/fq/pkg/interp"
)

type fuzzFS struct{}

func (fuzzFS) Open(name string) (fs.File, error) {
	return nil, fmt.Errorf("%s: file not found", name)
}

type fuzzTest struct {
	b []byte
}

type fuzzTestInput struct {
	interp.FileReader
	io.Writer
}

func (fuzzTestInput) IsTerminal() bool { return false }
func (fuzzTestInput) Size() (int, int) { return 120, 25 }

type fuzzTestOutput struct {
	io.Writer
}

func (o fuzzTestOutput) Size() (int, int) { return 120, 25 }
func (o fuzzTestOutput) IsTerminal() bool { return false }

func (ft *fuzzTest) Stdin() interp.Input {
	return fuzzTestInput{FileReader: interp.FileReader{R: bytes.NewBuffer(ft.b)}}
}
func (ft *fuzzTest) Stdout() interp.Output        { return fuzzTestOutput{os.Stdout} }
func (ft *fuzzTest) Stderr() interp.Output        { return fuzzTestOutput{os.Stderr} }
func (ft *fuzzTest) InterruptChan() chan struct{} { return nil }
func (ft *fuzzTest) Environ() []string            { return nil }
func (ft *fuzzTest) Args() []string {
	return []string{}
}
func (ft *fuzzTest) ConfigDir() (string, error) { return "/config", nil }
func (ft *fuzzTest) FS() fs.FS                  { return fuzzFS{} }
func (ft *fuzzTest) History() ([]string, error) { return nil, nil }

func (ft *fuzzTest) Readline(prompt string, complete func(line string, pos int) (newLine []string, shared int)) (string, error) {
	return "", io.EOF
}

func FuzzFormats(f *testing.F) {
	filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) == ".fqtest" {
			return nil
		}
		if filepath.Base(filepath.Dir(path)) != "testdata" {
			return nil
		}

		if st, err := os.Stat(path); err != nil || st.IsDir() {
			return nil
		}

		b, err := ioutil.ReadFile(path)
		if err != nil {
			panic(err)
		}

		f.Add(b)

		return nil
	})

	f.Fuzz(func(t *testing.T, b []byte) {
		fz := &fuzzTest{b: b}
		q, err := interp.New(fz, registry.Default)
		if err != nil {
			t.Fatal(err)
		}

		_ = q.Main(context.Background(), fz.Stdout(), "dev")
		// if err != nil {
		// 	// TODO: expect error
		// 	t.Fatal(err)
		// }
	})
}
