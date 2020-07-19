package cli

import (
	"flag"
	"fmt"
	"fq/pkg/bitbuf"
	"fq/pkg/decode"
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
		for _, f := range fs {
			merged = append(merged, f)
		}
	}
	registry := decode.NewRegistryWithFormats(merged)

	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fs.SetOutput(m.OS.Stderr())
	forceFormatNameFlag := fs.String("f", "", "Force format")
	verboseFlag := fs.Bool("v", false, "Verbose output")
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

	r := m.OS.Stdin()
	if fs.Arg(0) != "" && fs.Arg(0) != "-" {
		f, err := m.OS.Open(fs.Arg(0))
		if err != nil {
			return err
		}
		defer f.Close()
		r = f
	}
	buf, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	var forceFormats []*decode.Format
	if *forceFormatNameFlag != "" {
		forceFormat := registry.FindFormat(*forceFormatNameFlag)
		if forceFormat == nil {
			return fmt.Errorf("%s: found not find format", *forceFormatNameFlag)
		}
		forceFormats = append(forceFormats, forceFormat)
	}
	bb := bitbuf.NewFromBytes(buf)
	d, errs := registry.Probe(nil, fs.Arg(0), decode.Range{Start: 0, Stop: bb.Len}, bitbuf.NewFromBytes(buf), forceFormats)
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

	if d != nil {
		f := d.Root()
		exp := fs.Arg(1)
		if _, err := f.Eval(m.OS.Stdout(), exp); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("unable to probe format")
	}

	return nil
}
