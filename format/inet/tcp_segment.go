package inet

// https://en.wikipedia.org/wiki/Transmission_Control_Protocol

import (
	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
	"github.com/wader/fq/pkg/scalar"
)

func init() {
	interp.RegisterFormat(
		format.TCP_Segment,
		&decode.Format{
			Description: "Transmission control protocol segment",
			Groups:      []*decode.Group{format.IP_Packet},
			DecodeFn:    decodeTCP,
		})
}

const (
	tcpOptionEnd           = 0
	tcpOptionNop           = 1
	tcpOptionMSS           = 2
	tcpOptionWinscale      = 3
	tcpOptionSackPermitted = 4
	tcpOptionSack          = 5
	tcpOptionTimestamp     = 8
)

var tcpOptionsMap = scalar.UintMap{
	tcpOptionEnd:           {Sym: "end", Description: "End of options list"},
	tcpOptionNop:           {Sym: "nop", Description: "No operation"},
	tcpOptionMSS:           {Sym: "mss", Description: "Maximum segment size"},
	tcpOptionWinscale:      {Sym: "winscale", Description: "Window scale"},
	tcpOptionSackPermitted: {Sym: "sack_permitted", Description: "Selective Acknowledgement permitted"},
	tcpOptionSack:          {Sym: "sack", Description: "Selective Acknowledgement"},
	tcpOptionTimestamp:     {Sym: "timestamp", Description: "Timestamp and echo of previous timestamp"},
}

func decodeTCP(d *decode.D) any {
	var ipi format.IP_Packet_In
	if d.ArgAs(&ipi) && ipi.Protocol != format.IPv4ProtocolTCP {
		d.Fatalf("incorrect protocol %d", ipi.Protocol)
	}

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
	d.FieldU16("checksum", scalar.UintHex)
	// checksumEnd := d.Pos()
	d.FieldU16("urgent_pointer")
	optionsLen := (int64(dataOffset) - 5) * 8 * 4
	if optionsLen > 0 {
		d.FramedFn(optionsLen, func(d *decode.D) {
			d.FieldArray("options", func(d *decode.D) {
				for !d.End() {
					d.FieldStruct("option", func(d *decode.D) {
						kind := d.FieldU8("kind", tcpOptionsMap)
						switch kind {
						case tcpOptionEnd, tcpOptionNop:
							// has no length or data
						default:
							l := d.FieldU8("length")
							switch kind {
							case tcpOptionMSS:
								d.FieldU16("size")
							case tcpOptionWinscale:
								d.FieldU8("shift")
							case tcpOptionSackPermitted:
								// none
							case tcpOptionSack:
								d.FramedFn((int64(l-2))*8, func(d *decode.D) {
									d.FieldArray("blocks", func(d *decode.D) {
										for !d.End() {
											d.FieldStruct("block", func(d *decode.D) {
												d.FieldU32("left_edge")
												d.FieldU32("right_edge")
											})
										}
									})
								})
							case tcpOptionTimestamp:
								d.FieldU32("value")
								d.FieldU32("echo_reply")
							default:
								d.FieldRawLen("data", (int64(l-2))*8)
							}
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
	// _ = d.FieldMustGet("checksum").TryScalarFn(d.UintValidateBytes(tcpChecksum.Sum(nil)), scalar.Hex)

	d.FieldRawLen("payload", d.BitsLeft())

	return nil
}
