package postgres

import (
	"github.com/wader/fq/format"
	"github.com/wader/fq/format/registry"
	"github.com/wader/fq/pkg/decode"
	_ "github.com/wader/fq/pkg/scalar"
)

func init() {
	registry.MustRegister(decode.Format{
		Name:        format.PGWALPAGE,
		Description: "PostgreSQL write-ahead page",
		DecodeFn:    walpageDecode,
	})
}

//const XLOG_BLCKSZ = 8192

func walpageDecode(d *decode.D, in interface{}) interface{} {

	d.Endian = decode.LittleEndian

	pageHeaders := d.FieldArrayValue("XLogPageHeaders")
	_ = pageHeaders.FieldStruct("XLogPageHeaderData", decodeXLogPageHeaderData)

	return nil
}
