// Some of these functions are based on gojq func.go functions
// TODO: maybe should be exported from gojq fq branch instead?
// The MIT License (MIT)
// Copyright (c) 2019-2021 itchyny

package gojqextra

import (
	"math"
	"math/big"
	"strconv"

	"github.com/wader/gojq"
)

func ToString(x interface{}) (string, bool) {
	switch x := x.(type) {
	case string:
		return x, true
	case gojq.JQValue:
		return ToString(x.JQValueToGoJQ())
	default:
		return "", false
	}
}

func ToObject(x interface{}) (map[string]interface{}, bool) {
	switch x := x.(type) {
	case map[string]interface{}:
		return x, true
	case gojq.JQValue:
		return ToObject(x.JQValueToGoJQ())
	default:
		return nil, false
	}
}

func ToArray(x interface{}) ([]interface{}, bool) {
	switch x := x.(type) {
	case []interface{}:
		return x, true
	case gojq.JQValue:
		return ToArray(x.JQValueToGoJQ())
	default:
		return nil, false
	}
}

func ToBoolean(x interface{}) (bool, bool) {
	switch x := x.(type) {
	case bool:
		return x, true
	case gojq.JQValue:
		return ToBoolean(x.JQValueToGoJQ())
	default:
		return false, false
	}
}

func IsNull(x interface{}) bool {
	switch x := x.(type) {
	case nil:
		return true
	case gojq.JQValue:
		return IsNull(x.JQValueToGoJQ())
	default:
		return false
	}
}

func ToInt(x interface{}) (int, bool) {
	switch x := x.(type) {
	case int:
		return x, true
	case float64:
		return floatToInt(x), true
	case *big.Int:
		if x.IsInt64() {
			if i := x.Int64(); minInt <= i && i <= maxInt {
				return int(i), true
			}
		}
		if x.Sign() > 0 {
			return maxInt, true
		}
		return minInt, true
	case gojq.JQValue:
		return ToInt(x.JQValueToGoJQ())
	default:
		return 0, false
	}
}

func floatToInt(x float64) int {
	if minInt <= x && x <= maxInt {
		return int(x)
	}
	if x > 0 {
		return maxInt
	}
	return minInt
}

func ToFloat(x interface{}) (float64, bool) {
	switch x := x.(type) {
	case int:
		return float64(x), true
	case float64:
		return x, true
	case *big.Int:
		return bigToFloat(x), true
	case gojq.JQValue:
		return ToFloat(x.JQValueToGoJQ())
	default:
		return 0.0, false
	}
}

func bigToFloat(x *big.Int) float64 {
	if x.IsInt64() {
		return float64(x.Int64())
	}
	if f, err := strconv.ParseFloat(x.String(), 64); err == nil {
		return f
	}
	return math.Inf(x.Sign())
}
