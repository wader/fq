package gojqx

import (
	"encoding/json"
	"math"
	"math/big"

	"github.com/wader/gojq"
)

// from gojq
func NormalizeNumbers(v any) any {
	switch v := v.(type) {
	case json.Number:
		return gojq.ParseNumber(v)
	case *big.Int:
		if v.IsInt64() {
			if i := v.Int64(); math.MinInt <= i && i <= math.MaxInt {
				return int(i)
			}
		}
		return v
	case int64:
		if math.MinInt <= v && v <= math.MaxInt {
			return int(v)
		}
		return big.NewInt(v)
	case int32:
		return int(v)
	case int16:
		return int(v)
	case int8:
		return int(v)
	case uint:
		if v <= math.MaxInt {
			return int(v)
		}
		return new(big.Int).SetUint64(uint64(v))
	case uint64:
		if v <= math.MaxInt {
			return int(v)
		}
		return new(big.Int).SetUint64(v)
	case uint32:
		if uint64(v) <= math.MaxInt {
			return int(v)
		}
		return new(big.Int).SetUint64(uint64(v))
	case uint16:
		return int(v)
	case uint8:
		return int(v)
	case float32:
		return float64(v)
	case []any:
		for i, x := range v {
			v[i] = NormalizeNumbers(x)
		}
		return v
	case map[string]any:
		for k, x := range v {
			v[k] = NormalizeNumbers(x)
		}
		return v
	default:
		return v
	}
}
