package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"go.starlark.net/starlark"
)

type Parser struct {
	e           binary.ByteOrder
	r           io.ReadSeeker
	predeclared starlark.StringDict
}

func (p *Parser) uint8(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	name := ""
	if err := starlark.UnpackArgs(b.Name(), args, kwargs, "s?", &name); err != nil {
		return nil, err
	}
	var ba [1]byte
	_, _ = p.r.Read(ba[:])

	return starlark.MakeInt(int(ba[0])), nil
}

func NewParser() *Parser {
	p := &Parser{
		e: binary.BigEndian,
	}

	p.predeclared = starlark.StringDict{
		"uint8": starlark.NewBuiltin("uint8", p.uint8),
	}

	// vm.Set("uint8", func(call otto.FunctionCall) otto.Value {
	// 	var b [1]byte
	// 	_, _ = p.r.Read(b[:])
	// 	// TODO: err
	// 	result, _ := vm.ToValue(b[0])
	// 	return result
	// })
	// vm.Set("uint16", func(call otto.FunctionCall) otto.Value {
	// 	var b [2]byte
	// 	_, _ = p.r.Read(b[:])
	// 	n := p.e.Uint16(b[:])
	// 	// TODO: err
	// 	result, _ := vm.ToValue(n)
	// 	return result
	// })
	// vm.Set("uint24", func(call otto.FunctionCall) otto.Value {
	// 	var b [4]byte
	// 	// read into last 3 bytes then treat it as uint32
	// 	_, _ = p.r.Read(b[1:])
	// 	n := p.e.Uint32(b[:])
	// 	// TODO: err
	// 	result, _ := vm.ToValue(n)
	// 	return result
	// })
	// vm.Set("uint32", func(call otto.FunctionCall) otto.Value {
	// 	var b [4]byte
	// 	_, _ = p.r.Read(b[:])
	// 	n := p.e.Uint32(b[:])
	// 	// TODO: err
	// 	result, _ := vm.ToValue(n)
	// 	return result
	// })
	// vm.Set("ascii", func(call otto.FunctionCall) otto.Value {
	// 	n, _ := call.Argument(0).ToInteger()
	// 	b := make([]byte, n)
	// 	io.ReadFull(p.r, b)
	// 	// TODO: err
	// 	result, _ := vm.ToValue(string(b))
	// 	return result
	// })

	// vm.Set("bytes", func(call otto.FunctionCall) otto.Value {
	// 	n, _ := call.Argument(0).ToInteger()
	// 	b := make([]byte, n)
	// 	io.ReadFull(p.r, b)
	// 	// TODO: err
	// 	result, _ := vm.ToValue(string(b))
	// 	return result
	// })

	// vm.Set("range", func(call otto.FunctionCall) otto.Value {
	// 	n, _ := call.Argument(0).ToInteger()
	// 	var a []int64
	// 	for i := int64(0); i < n; i++ {
	// 		a = append(a, i)
	// 	}
	// 	result, _ := vm.ToValue(a)
	// 	return result
	// })

	// vm.Set("big_endian", func(call otto.FunctionCall) otto.Value {
	// 	p.e = binary.BigEndian
	// 	return otto.NullValue()
	// })

	// vm.Set("little_endian", func(call otto.FunctionCall) otto.Value {
	// 	p.e = binary.LittleEndian
	// 	return otto.NullValue()
	// })

	return p
}

func (p *Parser) Run(r io.ReadSeeker, s string) {
	p.r = r
	thread := &starlark.Thread{
		Name:  "example",
		Print: func(_ *starlark.Thread, msg string) { fmt.Println(msg) },
	}

	// Execute a program.
	globals, err := starlark.ExecFile(thread, os.Args[1], nil, p.predeclared)
	if err != nil {
		if evalErr, ok := err.(*starlark.EvalError); ok {
			log.Fatal(evalErr.Backtrace())
		}
		log.Fatal(err)
	}

	json.NewEncoder(os.Stdout).Encode(globals)
}

func main() {
	s, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}

	f, err := os.Open(os.Args[2])
	if err != nil {
		panic(err)
	}

	p := NewParser()
	p.Run(f, string(s))

}
