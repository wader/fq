// Copyright (c) 2022-2023 GoSecure Inc.
// Licensed under the MIT License
package pyrdp

import (
	"fmt"
	"os"
	"strings"

	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/scalar"
	"golang.org/x/text/encoding/unicode"
)

func toTextUTF16Fn(length int) func(d *decode.D) string {
	return func(d *decode.D) string {
		enc := unicode.UTF16(unicode.LittleEndian, unicode.UseBOM)
		decoder := enc.NewDecoder()

		decoded, _ := decoder.String(string(d.BytesLen(length)))
		return strings.TrimRight(decoded, "\x00")
	}
}

func printPos(d *decode.D) {
	fmt.Fprintf(os.Stderr, "Pos: %d\n", d.Pos())
}

var charMapper = scalar.UintFn(func(s scalar.Uint) (scalar.Uint, error) {
	char := s.Actual
	s.Sym = fmt.Sprintf("%c", int(char))
	return s, nil
})
