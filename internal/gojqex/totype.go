// Some of these functions are based on gojq func.go functions
// TODO: maybe should be exported from gojq fq branch instead?
// The MIT License (MIT)
// Copyright (c) 2019-2021 itchyny

package gojqex

import (
	"math"
	"math/big"

	"github.com/wader/gojq"
)

func IsNull(x any) bool {
	switch x := x.(type) {
	case nil:
		return true
	case gojq.JQValue:
		return IsNull(x.JQValueToGoJQ())
	default:
		return false
	}
}

func ToGoJQValue(v any) (any, bool) {
	switch vv := v.(type) {
	case nil:
		return vv, true
	case bool:
		return vv, true
	case int:
		return vv, true
	case int64:
		if vv >= math.MinInt && vv <= math.MaxInt {
			return int(vv), true
		}
		return big.NewInt(vv), true
	case uint64:
		if vv <= math.MaxInt {
			return int(vv), true
		}
		return new(big.Int).SetUint64(vv), true
	case float32:
		return float64(vv), true
	case float64:
		return vv, true
	case *big.Int:
		if vv.IsInt64() {
			vv := vv.Int64()
			if vv >= math.MinInt && vv <= math.MaxInt {
				return int(vv), true
			}
			return vv, true
		} else if vv.IsUint64() {
			vv := vv.Uint64()
			if vv <= math.MaxInt {
				return int(vv), true
			}
			return vv, true
		}
		return vv, true
	case string:
		return vv, true
	case []byte:
		return string(vv), true
	case gojq.JQValue:
		return ToGoJQValue(vv.JQValueToGoJQ())
	case []any:
		vvs := make([]any, len(vv))
		for i, v := range vv {
			v, ok := ToGoJQValue(v)
			if !ok {
				return nil, false
			}
			vvs[i] = v
		}
		return vvs, true
	case map[string]any:
		vvs := make(map[string]any, len(vv))
		for k, v := range vv {
			v, ok := ToGoJQValue(v)
			if !ok {
				return nil, false
			}
			vvs[k] = v
		}
		return vvs, true
	default:
		return nil, false
	}
}
