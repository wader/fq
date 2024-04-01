// Some of these functions are based on gojq func.go functions
// TODO: maybe should be exported from gojq fq branch instead?
// The MIT License (MIT)
// Copyright (c) 2019-2021 itchyny

package gojqx

import (
	"fmt"
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

func ToGoJQValue(v any) (any, error) {
	return ToGoJQValueFn(v, func(v any) (any, error) {
		switch v := v.(type) {
		case gojq.JQValue:
			return v.JQValueToGoJQ(), nil
		default:
			return nil, fmt.Errorf("not a JQValue")
		}
	})
}

func ToGoJQValueFn(v any, valueFn func(v any) (any, error)) (any, error) {
	switch vv := v.(type) {
	case nil:
		return vv, nil
	case bool:
		return vv, nil
	case int:
		return vv, nil
	case int64:
		if vv >= math.MinInt && vv <= math.MaxInt {
			return int(vv), nil
		}
		return big.NewInt(vv), nil
	case uint64:
		if vv <= math.MaxInt {
			return int(vv), nil
		}
		return new(big.Int).SetUint64(vv), nil
	case float32:
		return float64(vv), nil
	case float64:
		return vv, nil
	case *big.Int:
		if vv.IsInt64() {
			vv := vv.Int64()
			if vv >= math.MinInt && vv <= math.MaxInt {
				return int(vv), nil
			}
		}
		return vv, nil
	case string:
		return vv, nil
	case []byte:
		return string(vv), nil
	case []any:
		vvs := make([]any, len(vv))
		for i, v := range vv {
			v, err := ToGoJQValueFn(v, valueFn)
			if err != nil {
				return nil, err
			}
			vvs[i] = v
		}
		return vvs, nil
	case map[string]any:
		vvs := make(map[string]any, len(vv))
		for k, v := range vv {
			v, err := ToGoJQValueFn(v, valueFn)
			if err != nil {
				return nil, err
			}
			vvs[k] = v
		}
		return vvs, nil
	case error:
		return nil, vv
	default:
		nv, err := valueFn(vv)
		if err != nil {
			return nil, err
		}

		return ToGoJQValueFn(nv, valueFn)
	}
}
