package cli

import (
	"context"
	"flag"
	"fmt"
	"fq"
	"fq/pkg/decode"
	"fq/pkg/osenv"
	"fq/pkg/query"
	"io"
	"os"
)

type StandardOS struct{}

func (StandardOS) Stdin() io.Reader                        { return os.Stdin }
func (StandardOS) Stdout() io.Writer                       { return os.Stdout }
func (StandardOS) Stderr() io.Writer                       { return os.Stderr }
func (StandardOS) Environ() []string                       { return os.Environ() }
func (StandardOS) Args() []string                          { return os.Args }
func (StandardOS) Open(name string) (io.ReadSeeker, error) { return os.Open(name) }

func StandardOSMain(r *decode.Registry) {
	if err := (Main{
		OS:       StandardOS{},
		Registry: r,
	}).Run(); err != nil {
		os.Exit(1)
	}
}

type Main struct {
	OS       osenv.OS
	Registry *decode.Registry
}

// Run cli main
func (m Main) Run() error {
	err := m.run()
	if err != nil && err != flag.ErrHelp {
		fmt.Fprintf(m.OS.Stderr(), "%s\n", err)
	}
	return err
}

func (m Main) run() error {
	filename := "asd"

	var args []interface{}
	for _, a := range m.OS.Args() {
		args = append(args, a)
	}

	q := query.NewQuery(query.QueryOptions{
		Variables: map[string]interface{}{
			"$FILENAME": filename,
			"$VERSION":  fq.Version,
			"$ARGS":     args,
		},
		Registry: m.Registry,

		Environ: m.OS.Environ, // TODO: func?
		Stdin:   m.OS.Stdin(),
		Open:    m.OS.Open,
	})

	runMode := query.ScriptMode

	i, err := q.Eval(context.Background(), runMode, nil, "main", query.WriterOutput{Ctx: context.Background(), W: m.OS.Stdout()}, nil)
	if err != nil {
		return err
	}
	for {
		v, ok := i.Next()
		if !ok {
			break
		} else if err, ok := v.(error); ok {
			fmt.Fprintln(m.OS.Stderr(), err)
			break
		} else if d, ok := v.([2]interface{}); ok {
			fmt.Fprintf(m.OS.Stdout(), "%s: %v\n", d[0], d[1])
		}
	}

	return nil
}
