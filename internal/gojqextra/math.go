// from gojq
// The MIT License (MIT)
// Copyright (c) 2019-2021 itchyny

package gojqextra

import "math/bits"

const (
	maxInt = 1<<(bits.UintSize-1) - 1 // math.MaxInt64 or math.MaxInt32
	minInt = -maxInt - 1              // math.MinInt64 or math.MinInt32
)
