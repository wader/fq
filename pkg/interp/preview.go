package interp

import (
	"fmt"
	"math/big"
	"strconv"

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
		return fmt.Sprintf("%q", stringsex.TrimN(vv, 50, "..."))
	case nil:
		return "null"
	case bitio.Reader:
		return "raw bits"
	case *big.Int:
		return mathex.PadFormatBigInt(vv, df.FormatBase(), true, 0)
	default:
		panic("unreachable")
	}
}
