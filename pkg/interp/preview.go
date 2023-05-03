package interp

import (
	"bytes"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/wader/fq/internal/mathex"
	"github.com/wader/fq/internal/stringsex"
	"github.com/wader/fq/pkg/bitio"
	"github.com/wader/fq/pkg/scalar"
)

func previewValue(v any, df scalar.DisplayFormat) string {
	switch vv := v.(type) {
	case bool:
		if vv {
			return "true"
		}
		return "false"
	case int:
		// TODO: DisplayFormat is weird
		return mathex.PadFormatInt(int64(vv), df.FormatBase(), true, 0)
	case int64:
		// TODO: DisplayFormat is weird
		return mathex.PadFormatInt(vv, df.FormatBase(), true, 0)
	case uint64:
		return mathex.PadFormatUint(vv, df.FormatBase(), true, 0)
	case float64:
		// TODO: float32? better truncated to significant digits?
		return strconv.FormatFloat(vv, 'g', -1, 64)
	case string:
		s := strconv.Quote(stringsex.TrimN(vv, 50, "..."))
		// TODO: hack for https://github.com/golang/go/issues/52062
		// 0x7f used to be escaped as \u007f in 1.18 and lower, was changed to \x7f
		// remove once 1.18 is not supported
		if !bytes.Contains([]byte(s), []byte{0x7f}) {
			return s
		}
		return strings.ReplaceAll(s, `\u007f`, `\x7f`)
	case nil:
		return "null"
	case bitio.Reader,
		Binary:
		return "raw bits"
	case *big.Int:
		return mathex.PadFormatBigInt(vv, df.FormatBase(), true, 0)
	case map[string]any:
		return "{}"
	case []any:
		return "[]"

	default:
		panic(fmt.Sprintf("unreachable %v (%T)", v, v))
	}
}
