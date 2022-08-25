package dns

import (
	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
)

func init() {
	interp.RegisterFormat(decode.Format{
		Name:        format.DNS_TCP,
		Description: "DNS packet (TCP)",
		Groups:      []string{format.TCP_STREAM},
		DecodeFn:    dnsTCPDecode,
	})
}

func dnsTCPDecode(d *decode.D, in any) any {
	if tsi, ok := in.(format.TCPStreamIn); ok {
		tsi.MustIsPort(d.Fatalf, format.TCPPortDomain)
	}
	return dnsDecode(d, true)
}
