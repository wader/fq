package av1

import (
	"fq/format"
	"fq/format/registry"
	"fq/pkg/decode"
)

func init() {
	registry.MustRegister(&decode.Format{
		Name:        format.AV1_OBU,
		Description: "AV1 Open Bitstream Unit",
		DecodeFn:    obuDecode,
	})
}

const (
	OBU_SEQUENCE_HEADER        = 1
	OBU_TEMPORAL_DELIMITER     = 2
	OBU_FRAME_HEADER           = 3
	OBU_TILE_GROUP             = 4
	OBU_METADATA               = 5
	OBU_FRAME                  = 6
	OBU_REDUNDANT_FRAME_HEADER = 7
	OBU_TILE_LIST              = 8
	OBU_PADDING                = 15
)

var obuTypeNames = map[uint64]string{
	OBU_SEQUENCE_HEADER:        "OBU_SEQUENCE_HEADER",
	OBU_TEMPORAL_DELIMITER:     "OBU_TEMPORAL_DELIMITER",
	OBU_FRAME_HEADER:           "OBU_FRAME_HEADER",
	OBU_TILE_GROUP:             "OBU_TILE_GROUP",
	OBU_METADATA:               "OBU_METADATA",
	OBU_FRAME:                  "OBU_FRAME",
	OBU_REDUNDANT_FRAME_HEADER: "OBU_REDUNDANT_FRAME_HEADER",
	OBU_TILE_LIST:              "OBU_TILE_LIST",
	OBU_PADDING:                "OBU_PADDING",
}

func leb128(d *decode.D) uint64 {
	var v uint64
	for i := 0; i < 8; i++ {
		b := d.U8()
		v = v | (b&0x7f)<<(i*7)
		if b&0x80 == 0 {
			break
		}
	}
	return v
}

func fieldLeb128(d *decode.D, name string) uint64 {
	return d.FieldUFn(name, func() (uint64, decode.DisplayFormat, string) {
		return leb128(d), decode.NumberDecimal, ""
	})
}

func obuDecode(d *decode.D, in interface{}) interface{} {
	var obuType uint64
	var obuSize int64
	hasExtension := false
	hasSizeField := false

	d.FieldStructFn("header", func(d *decode.D) {
		d.FieldU1("forbidden_bit")
		obuType, _ = d.FieldStringMapFn("type", obuTypeNames, "Reserved", d.U4, decode.NumberDecimal)
		hasExtension = d.FieldBool("extension_flag")
		hasSizeField = d.FieldBool("has_size_field")
		d.FieldU1("reserved_1bit")
		if hasExtension {
			d.FieldU3("temporal_id")
			d.FieldU2("spatial_id")
			d.FieldU3("extension_header_reserved_3bits")
		}
	})

	if hasSizeField {
		obuSize = int64(fieldLeb128(d, "size"))
	} else {
		obuSize = d.BitsLeft() / 8
		if hasExtension {
			obuSize--
		}
	}

	_ = obuType

	if d.BitsLeft() > 0 {
		d.FieldBitBufLen("data", obuSize*8)
	}

	return nil
}
