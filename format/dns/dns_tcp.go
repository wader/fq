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
	if tsi, ok := in.(format.TCPStreamIn); ok {
		tsi.MustIsPort(d.Fatalf, format.TCPPortDomain, format.TCPPortDomain)
	}
	return dnsDecode(d, true)
}
