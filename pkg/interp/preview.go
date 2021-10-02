package interp

import (
	"encoding/hex"
	"fmt"
	"strconv"

	"github.com/wader/fq/internal/num"
	"github.com/wader/fq/pkg/bitio"
	"github.com/wader/fq/pkg/decode"
)

func previewValue(v *decode.Value) string {
	switch vv := v.V.(type) {
	case decode.Array:
		return "[]"
	case decode.Struct:
		return v.Description
	case bool:
		if vv {
			return "true"
		}
		return "false"
	case int64:
		// TODO: DisplayFormat is weird
		return num.PadFormatInt(vv, decode.DisplayFormatToBase(v.DisplayFormat), true, 0)
	case uint64:
		return num.PadFormatUint(vv, decode.DisplayFormatToBase(v.DisplayFormat), true, 0)
	case float64:
		// TODO: float32? better truncated to significant digits?
		return strconv.FormatFloat(vv, 'g', -1, 64)
	case string:
		if len(vv) > 50 {
			return fmt.Sprintf("%q", vv[0:50]) + "..."
		}
		return fmt.Sprintf("%q", vv)
	case []byte:
		if len(vv) > 16 {
			return hex.EncodeToString(vv[0:16]) + "..."
		}
		return hex.EncodeToString(vv)
	case *bitio.Buffer:
		vvLen := vv.Len()
		if vvLen > 16*8 {
			bs, _ := vv.BytesRange(0, 16)
			return hex.EncodeToString(bs) + "..."
		}
		bs, _ := vv.BytesRange(0, int(bitio.BitsByteCount(vvLen)))
		return hex.EncodeToString(bs)
	case nil:
		return "none"
	default:
		panic("unreachable")
	}
}
