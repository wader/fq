package inet

// https://en.wikipedia.org/wiki/Transmission_Control_Protocol

import (
	"github.com/wader/fq/format"
	"github.com/wader/fq/format/registry"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/scalar"
)

func init() {
	registry.MustRegister(decode.Format{
		Name:        format.TCP_SEGMENT,
		Description: "Transmission control protocol segment",
		DecodeFn:    decodeTCP,
	})
}

const (
	tcpOptionEnd = 0
	tcpOptionNop = 1
)

var tcpOptionsMap = scalar.UToScalar{
	tcpOptionEnd: {Sym: "end", Description: "End of options list"},
	tcpOptionNop: {Sym: "nop", Description: "No operation"},
	2:            {Sym: "maxseg", Description: "Maximum segment size"},
	3:            {Sym: "winscale", Description: "Window scale"},
	4:            {Sym: "sack_permitted", Description: "Selective Acknowledgement permitted"},
	5:            {Sym: "sack", Description: "Selective ACKnowledgement"},
	8:            {Sym: "timestamp", Description: "Timestamp and echo of previous timestamp"},
}

func decodeTCP(d *decode.D, in interface{}) interface{} {
	d.FieldU16("source_port", format.TCPPortMap)
	d.FieldU16("destination_port", format.TCPPortMap)
	d.FieldU32("sequence_number")
	d.FieldU32("acknowledgment_number")
	dataOffset := d.FieldU4("data_offset")
	d.FieldU3("reserved")
	d.FieldBool("ns")
	d.FieldBool("cwr")
	d.FieldBool("ece")
	d.FieldBool("urg")
	d.FieldBool("ack")
	d.FieldBool("psh")
	d.FieldBool("rst")
	d.FieldBool("syn")
	d.FieldBool("fin")
	d.FieldU16("window_size")
	// checksumStart := d.Pos()
	d.FieldU16("checksum", scalar.Hex)
	// checksumEnd := d.Pos()
	d.FieldU16("urgent_pointer")
	optionsLen := (int64(dataOffset) - 5) * 8 * 4
	if optionsLen > 0 {
		d.LenFn(optionsLen, func(d *decode.D) {
			d.FieldArray("options", func(d *decode.D) {
				for !d.End() {
					d.FieldStruct("option", func(d *decode.D) {
						kind := d.FieldU8("kind", tcpOptionsMap)
						switch kind {
						case tcpOptionEnd, tcpOptionNop:
						default:
							l := d.FieldU8("length")
							d.FieldRawLen("data", (int64(l-2))*8)
						}
					})
				}
			})
		})
	}

	// TODO: need to pass ipv4 pseudo header somehow
	// tcpChecksum := &checksum.IPv4{}
	// d.MustCopy(tcpChecksum, d.BitBufRange(0, checksumStart))
	// d.MustCopy(tcpChecksum, d.BitBufRange(checksumEnd, d.Len()-checksumEnd))
	// _ = d.FieldMustGet("checksum").TryScalarFn(d.ValidateUBytes(tcpChecksum.Sum(nil)), scalar.Hex)

	d.FieldRawLen("data", d.BitsLeft())

	return nil
}
