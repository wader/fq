package main

import (
	"encoding/binary"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"runtime/pprof"

	"github.com/glycerine/zygomys/zygo"
	"go.starlark.net/starlark"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")
var memprofile = flag.String("memprofile", "", "write memory profile to `file`")

type Parser struct {
	z           *zygo.Zlisp
	e           binary.ByteOrder
	size        int64
	pos         int64
	buf         []byte
	predeclared starlark.StringDict
}

func (p *Parser) u8(env *zygo.Zlisp, name string, args []zygo.Sexp) (zygo.Sexp, error) {
	if len(args) != 1 {
		return zygo.SexpNull, zygo.WrongNargs
	}
	var n string
	switch t := args[0].(type) {
	case *zygo.SexpStr:
		n = t.S
	default:
		return zygo.SexpNull, errors.New(fmt.Sprintf("1st argument of %v should be a string", name))
	}

	_ = n

	c := p.pos
	p.pos++
	return &zygo.SexpInt{Val: int64(p.buf[c])}, nil
}

func (p *Parser) eof(env *zygo.Zlisp, name string, args []zygo.Sexp) (zygo.Sexp, error) {
	return &zygo.SexpBool{Val: p.pos >= p.size}, nil
}

func (p *Parser) Run(buf []byte, s string) error {
	p.buf = buf

	if err := p.z.LoadString(s); err != nil {
		return err
	}

	p.pos = 0
	p.size = int64(len(buf))

	e, err := p.z.Run()
	if err != nil {
		return err
	}

	log.Printf("e: %#+v\n", e)

	//	json.NewEncoder(os.Stdout).Encode(globals)

	return nil
}

func NewParser() *Parser {
	p := &Parser{
		e: binary.BigEndian,
	}

	p.z = zygo.NewZlispSandbox()
	p.z.StandardSetup()
	p.z.AddFunction("u8", p.u8)
	p.z.AddFunction("eof", p.eof)

	return p
}

func main() {
	flag.Parse()

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		defer f.Close() // error handling omitted for example
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	// ... rest of the program ...

	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal("could not create memory profile: ", err)
		}
		defer f.Close() // error handling omitted for example
		runtime.GC()    // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatal("could not write memory profile: ", err)
		}
	}

	s, err := ioutil.ReadFile(flag.Arg(0))
	if err != nil {
		panic(err)
	}

	buf, err := ioutil.ReadFile(flag.Arg(1))
	if err != nil {
		panic(err)
	}

	p := NewParser()
	if err := p.Run(buf, string(s)); err != nil {
		panic(err)
	}
}
