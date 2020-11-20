package cli

import (
	"bytes"
	"flag"
	"fmt"
	"fq/internal/bitio"
	"fq/pkg/decode"
	"fq/pkg/format"
	"fq/pkg/output"
	"io"
	"io/ioutil"
	"os"
	"sort"
	"strings"
)

type Main struct {
	OS       OS
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
	allFormats := m.Registry.MustAll()
	probeFormats := m.Registry.MustGroup(format.PROBE)

	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fs.SetOutput(m.OS.Stderr())
	dotFlag := fs.Bool("dot", false, "Output dot format graph (... | dot -Tsvg -o formats.svg)")
	forceFormatNameFlag := fs.String("f", "", "Force format")
	verboseFlag := fs.Bool("v", false, "Verbose output")
	outputFormatFlag := fs.String("o", "text", "Output format")
	fs.Usage = func() {
		maxNameLen := 0
		for _, f := range allFormats {
			if len(f.Name) > maxNameLen {
				maxNameLen = len(f.Name)
			}
		}

		formatsSorted := make([]*decode.Format, len(allFormats))
		copy(formatsSorted, allFormats)
		sort.Slice(formatsSorted, func(i, j int) bool {
			return formatsSorted[i].Name < formatsSorted[j].Name
		})

		pad := func(n int, s string) string { return strings.Repeat(" ", n-len(s)) }
		fmt.Fprintf(fs.Output(), "Usage: %s [FLAGS] FILE [EXP]\n", m.OS.Args()[0])
		fs.PrintDefaults()
		fmt.Fprintf(fs.Output(), "\n")
		fmt.Fprintf(fs.Output(), "Name:%s    MIME:\n", pad(maxNameLen, "Name:"))
		for _, f := range formatsSorted {
			fmt.Fprintf(fs.Output(), "%s%s    %s\n", f.Name, pad(maxNameLen, f.Name), strings.Join(f.MIMEs, ", "))
		}
	}
	if err := fs.Parse(m.OS.Args()[1:]); err != nil {
		return err
	}

	if *dotFlag {
		m.Registry.Dot(m.OS.Stdout())
		os.Exit(0)
	}

	var rs io.ReadSeeker
	if fs.Arg(0) != "" && fs.Arg(0) != "-" {
		f, err := m.OS.Open(fs.Arg(0))
		if err != nil {
			return err
		}
		defer f.Close()
		rs = f
	} else {
		buf, err := ioutil.ReadAll(m.OS.Stdin())
		if err != nil {
			return err
		}
		rs = bytes.NewReader(buf)
	}

	if *forceFormatNameFlag != "" {
		var err error
		probeFormats, err = m.Registry.Group(*forceFormatNameFlag)
		if err != nil {
			return fmt.Errorf("%s: %s", *forceFormatNameFlag, err)
		}
	}
	bb := bitio.NewBufferFromReadSeeker(rs)
	f, _, errs := decode.Probe(fs.Arg(0), bb, probeFormats)
	if *verboseFlag {
		for _, err := range errs {
			fmt.Fprintf(m.OS.Stderr(), "%s\n", err)
			if pe, ok := err.(*decode.DecodeError); ok {
				// if pe.PanicHandeled {
				fmt.Fprintf(m.OS.Stderr(), "%s", pe.PanicStack)
				// }
			}
		}
	}

	if f != nil {
		exp := fs.Arg(1)
		expValue, err := f.Eval(exp)
		if err != nil {
			return fmt.Errorf("%s: %s", exp, err)
		}

		var of *decode.FieldOutput
		for _, of = range output.All {
			if of.Name == *outputFormatFlag {
				break
			}
		}
		if of == nil {
			return fmt.Errorf("%s: unable to find output format", *outputFormatFlag)
		}

		if err := of.New(expValue).Write(m.OS.Stdout()); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("unable to probe format")
	}

	return nil
}
