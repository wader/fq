package interp

import (
	"fmt"
	"strconv"

	"github.com/wader/fq/internal/num"
	"github.com/wader/fq/pkg/bitio"
	"github.com/wader/fq/pkg/scalar"
)

func previewValue(v interface{}, df scalar.DisplayFormat) string {
	switch vv := v.(type) {
	case bool:
		if vv {
			return "true"
		}
		return "false"
	case int:
		// TODO: DisplayFormat is weird
		return num.PadFormatInt(int64(vv), df.FormatBase(), true, 0)
	case int64:
		// TODO: DisplayFormat is weird
		return num.PadFormatInt(vv, df.FormatBase(), true, 0)
	case uint64:
		return num.PadFormatUint(vv, df.FormatBase(), true, 0)
	case float64:
		// TODO: float32? better truncated to significant digits?
		return strconv.FormatFloat(vv, 'g', -1, 64)
	case string:
		if len(vv) > 50 {
			return fmt.Sprintf("%q", vv[0:50]) + "..."
		}
		return fmt.Sprintf("%q", vv)
	case *bitio.Buffer:
		return "raw bits"
	default:
		panic("unreachable")
	}
}
