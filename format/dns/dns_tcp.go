package dns

import (
	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
)

func init() {
	interp.RegisterFormat(
		format.DnsTcp,
		&decode.Format{
			Description: "DNS packet (TCP)",
			Groups:      []*decode.Group{format.TcpStream},
			DecodeFn:    dnsTCPDecode,
		})
}

func dnsTCPDecode(d *decode.D) any {
	var tsi format.TCPStreamIn
	if d.ArgAs(&tsi) {
		tsi.MustIsPort(d.Fatalf, format.TCPPortDomain)
	}
	return dnsDecode(d, true)
}
