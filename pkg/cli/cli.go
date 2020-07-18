package cli

import (
	"flag"
	"fmt"
	"fq/pkg/bitbuf"
	"fq/pkg/decode"
	"io"
	"io/ioutil"
	"os"
	"sort"
	"strings"
)

type Main struct {
	Stdout      io.Writer
	Args        []string
	FormatsList [][]*decode.Format
}

func (m Main) Run() error {
	err := m.run()
	if err != nil {
		fmt.Fprintf(os.Stdout, "%s\n", err)
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

	fs := flag.NewFlagSet("", flag.ExitOnError)
	fs.SetOutput(m.Stdout)
	forceFormatNameFlag := fs.String("f", "", "Force format")
	verboseFlag := fs.Bool("v", false, "")
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

		fmt.Fprintf(m.Stdout, "Usage: %s FILE [EXP]\n", m.Args[0])
		fs.PrintDefaults()
		fmt.Fprintf(m.Stdout, "\n")
		fmt.Fprintf(m.Stdout, "Name:%s    MIME:\n", pad(maxNameLen, "Name:"))
		for _, f := range formatsSorted {
			fmt.Fprintf(m.Stdout, "%s%s    %s\n", f.Name, pad(maxNameLen, f.Name), strings.Join(f.MIMEs, ", "))
		}
	}
	if err := fs.Parse(m.Args); err != nil {
		return err
	}

	if fs.Arg(1) == "" {
		fs.Usage()
		os.Exit(1)
	}

	buf, err := ioutil.ReadFile(fs.Arg(1))
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
	d, errs := registry.Probe(nil, fs.Arg(1), decode.Range{Start: 0, Stop: bb.Len}, bitbuf.NewFromBytes(buf), forceFormats)
	if d == nil || *verboseFlag {
		for _, err := range errs {
			fmt.Fprintf(m.Stdout, "%s\n", err)
			if pe := err.(*decode.ProbeError); pe != nil {
				// if pe.PanicHandeled {
				fmt.Fprintf(m.Stdout, "%s", pe.PanicStack)
				// }
			}
		}
	}

	if d != nil {
		f := d.Root()
		exp := fs.Arg(2)
		if _, err := f.Eval(os.Stdout, exp); err != nil {
			return err
		}
	}

	return nil
}
