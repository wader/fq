package dns

import (
	"github.com/wader/fq/format"
	"github.com/wader/fq/format/registry"
	"github.com/wader/fq/pkg/decode"
)

func init() {
	registry.MustRegister(decode.Format{
		Name:        format.DNS_TCP,
		Description: "DNS packet (TCP)",
		DecodeFn:    dnsTCPDecode,
	})
}

func dnsTCPDecode(d *decode.D, in interface{}) interface{} {
	return dnsDecode(d, true)
}
