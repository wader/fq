package inet

import (
	"github.com/wader/fq/format"
	"github.com/wader/fq/format/registry"
	"github.com/wader/fq/pkg/decode"
)

func init() {
	registry.MustRegister(decode.Format{
		Name:        format.TCP,
		Description: "Transmission Control Protocol",
		DecodeFn:    decodeTCP,
	})
}

func decodeTCP(d *decode.D, in interface{}) interface{} {
	d.FieldU16("source_port", d.MapUToScalar(tcpPortMap))
	d.FieldU16("destination_port", d.MapUToScalar(tcpPortMap))
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
	d.FieldU16("checksum", d.Hex)
	d.FieldU16("urgent_pointer")
	if dataOffset > 5 {
		d.FieldRawLen("options", (int64(dataOffset)-5)*8*4)
	}
	d.FieldRawLen("data", d.BitsLeft())

	return nil
}
