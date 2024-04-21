//go:build exclude

// tool to convert go fuzz input files to bytes
// Usage: cat format/testdata/fuzz/FuzzFormats/144bde49b40c90fd05d302ec90b6ddb2b6d6aea553bad520a8b954797e40fe72 | go run dev/fuzzbytes.go | go run .
package main

import (
	"bytes"
	"io"
	"os"
	"strconv"
)

func main() {
	bs, err := io.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}

	// Input looks like this:
	// go test fuzz v1
	// []byte("...")
	prefix := []byte("[]byte(")
	start := bytes.Index(bs, prefix) + len(prefix)
	end := len(bs) - 2
	s, err := strconv.Unquote(string(bs[start:end]))
	if err != nil {
		panic(err)
	}

	if _, err := os.Stdout.Write([]byte(s)); err != nil {
		panic(err)
	}
}
