package cli

import (
	"bytes"
	"flag"
	"fmt"
	"fq/pkg/bitbuf"
	"fq/pkg/decode"
	"fq/pkg/output"
	"io"
	"io/ioutil"
	"sort"
	"strings"
)

type Main struct {
	OS          OS
	FormatsList [][]*decode.Format
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
	var merged []*decode.Format
	for _, fs := range m.FormatsList {
		merged = append(merged, fs...)
	}
	registry := decode.NewRegistryWithFormats(merged)

	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fs.SetOutput(m.OS.Stderr())
	forceFormatNameFlag := fs.String("f", "", "Force format")
	verboseFlag := fs.Bool("v", false, "Verbose output")
	outputFormatFlag := fs.String("o", "text", "Output format")
	fs.Usage = func() {
		maxNameLen := 0
		for _, f := range registry.Formats {
			if len(f.Name) > maxNameLen {
				maxNameLen = len(f.Name)
			}
		}

		formatsSorted := make([]*decode.Format, len(registry.Formats))
		copy(formatsSorted, registry.Formats)
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

	var forceFormats []*decode.Format
	if *forceFormatNameFlag != "" {
		forceFormat := registry.FindFormat(*forceFormatNameFlag)
		if forceFormat == nil {
			return fmt.Errorf("%s: could not find format", *forceFormatNameFlag)
		}
		forceFormats = append(forceFormats, forceFormat)
	}
	bb, err := bitbuf.NewFromReadSeeker(rs, 0)
	if err != nil {
		panic(err)
	}
	f, _, _, errs := registry.Probe(nil, fs.Arg(0), decode.Range{Start: 0, Stop: bb.Len}, bb, forceFormats)
	if *verboseFlag {
		for _, err := range errs {
			fmt.Fprintf(m.OS.Stderr(), "%s\n", err)
			if pe := err.(*decode.ProbeError); pe != nil {
				// if pe.PanicHandeled {
				fmt.Fprintf(m.OS.Stderr(), "%s", pe.PanicStack)
				// }
			}
		}
	}

	if f != nil {
		exp := fs.Arg(1)
		expField, err := f.Eval(exp)
		if err != nil {
			return fmt.Errorf("%s: %s", exp, err)
		}

		var ow decode.FieldWriter
		for _, of := range output.All {
			if of.Name != *outputFormatFlag {
				continue
			}
			// TODO: multi?
			ow = of.New(expField.(*decode.Field))
		}
		if ow == nil {
			return fmt.Errorf("%s: unable to find output format", *outputFormatFlag)
		}

		if err := ow.Write(m.OS.Stdout()); err != nil {
			return err
		}

		// switch expType {
		// case decode.FieldExpTree:
		// 	ow.Write(m.OS.Stdout())
		// case decode.FieldExpValue:
		// 	fmt.Fprintf(m.OS.Stdout(), "%s", expField.Value.RawString())
		// case decode.FieldExpRange:
		// 	fmt.Fprintf(m.OS.Stdout(), "%s\n", expField.Range)
		// }

	} else {
		return fmt.Errorf("unable to probe format")
	}

	return nil
}
