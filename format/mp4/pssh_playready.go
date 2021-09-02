package mp4

import (
	"github.com/wader/fq/format"
	"github.com/wader/fq/format/registry"
	"github.com/wader/fq/pkg/decode"
)

func init() {
	registry.MustRegister(&decode.Format{
		Name:        format.PSSH_PLAYREADY,
		Description: "PlayReady PSSH",
		DecodeFn:    playreadyPsshDecode,
	})
}

const (
	recordTypeRightsManagementHeader = 1
	recordTypeLicenseStore           = 2
)

var recordTypeNames = map[uint64]string{
	recordTypeRightsManagementHeader: "Rights management header",
	recordTypeLicenseStore:           "License store",
}

func playreadyPsshDecode(d *decode.D, in interface{}) interface{} {
	d.Endian = decode.LittleEndian

	d.FieldU32("size")
	count := d.FieldU16("count")
	i := uint64(0)
	d.FieldStructArrayLoopFn("records", "record", func() bool { return i < count }, func(d *decode.D) {
		recordType, _ := d.FieldStringMapFn("type", recordTypeNames, "Unknown", d.U16, decode.NumberDecimal)
		recordLen := d.FieldU16("len")
		switch recordType {
		case recordTypeRightsManagementHeader, recordTypeLicenseStore:
			d.FieldUTF16LE("xml", int(recordLen))
		default:
			d.FieldBitBufLen("data", int64(recordLen)*8)
		}
		i++
	})

	return nil
}
