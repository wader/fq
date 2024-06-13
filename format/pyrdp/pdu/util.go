// Copyright (c) 2022-2023 GoSecure Inc.
// Copyright (c) 2024 Flare Systems
// Licensed under the MIT License
package pdu

import (
	"fmt"

	"github.com/wader/fq/pkg/scalar"
)

var charMapper = scalar.UintFn(func(s scalar.Uint) (scalar.Uint, error) {
	char := s.Actual
	s.Sym = fmt.Sprintf("%c", int(char))
	return s, nil
})
