// +build gofuzzbeta

package format_test

import (
	"bytes"
	"context"
	"fmt"
	"fq/pkg/format"
	"fq/pkg/interp"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

type fuzzTest struct {
	b []byte
}

type fuzzTestOutput struct {
	io.Writer
}

func (o fuzzTestOutput) Size() (int, int) { return 120, 25 }
func (o fuzzTestOutput) IsTerminal() bool { return false }

func (ft *fuzzTest) Stdin() io.Reader         { return bytes.NewBuffer(ft.b) } // TODO: special file?
func (ft *fuzzTest) Stdout() interp.Output    { return fuzzTestOutput{os.Stdout} }
func (ft *fuzzTest) Stderr() io.Writer        { return os.Stderr }
func (ft *fuzzTest) Interrupt() chan struct{} { return nil }
func (ft *fuzzTest) Environ() []string        { return nil }
func (ft *fuzzTest) Args() []string {
	return []string{}
}
func (ft *fuzzTest) ConfigDir() (string, error) { return "/config", nil }
func (ft *fuzzTest) Open(name string) (io.ReadSeeker, error) {
	return nil, fmt.Errorf("%s: file not found", name)
}
func (ft *fuzzTest) Readline(prompt string, complete func(line string, pos int) (newLine []string, shared int)) (string, error) {
	return "", io.EOF
}

func FuzzFQTests(f *testing.F) {

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
		q, err := interp.New(fz, format.DefaultRegistry)
		if err != nil {
			t.Fatal(err)
		}

		err = q.Main(context.Background(), fz.Stdout(), "dev")
		if err != nil {
			// TODO: expect error
			t.Fatal(err)
		}
	})
}
