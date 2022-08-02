package fit

// https://developer.garmin.com/fit/protocol/
import (
	"sort"

	"github.com/wader/fq/format"
	"github.com/wader/fq/format/registry"
	"github.com/wader/fq/pkg/decode"
)

type field_definition struct {
	number    uint8
	size      uint8
	base_type struct {
		endian_ability uint8
		number         uint8
	}
}

func init() {
	registry.MustRegister(decode.Format{
		Name:        format.FIT,
		Description: "Flexible and Interoperable Data Transfer (FIT) Protocol",
		DecodeFn:    decodeFIT,
	})
}

func decodeFieldDefinition(d *decode.D) field_definition {
	var field field_definition
	field.number = uint8(d.FieldU8("field_definition_number"))
	field.size = uint8(d.FieldU8("size"))
	field.base_type.endian_ability = uint8(d.FieldU1("base_type_endian_ability"))
	d.RawLen(2)
	field.base_type.number = uint8(d.FieldU5("base_type_number"))
	return field
}

func decodeDataMessage(d *decode.D, fields []field_definition) {
	sort.Slice(fields, func(p, q int) bool {
		return fields[p].number < fields[q].number
	})

	d.FieldArray("field", func(d *decode.D) {
		for _, e := range fields {
			switch e.base_type.number {
			case 0:
				if (e.size / 1) > 1 {
					d.FieldArray("values", func(d *decode.D) {
						array_size := e.size / 1
						for i := uint8(0); i < array_size; i++ {
							d.FieldU8("enum")
						}
					})
				} else {
					d.FieldU8("enum")
				}
			case 4:
				if (e.size / 2) > 1 {
					d.FieldArray("values", func(d *decode.D) {
						array_size := e.size / 2
						for i := uint8(0); i < array_size; i++ {
							d.FieldU16("uint16")
						}
					})
				} else {
					d.FieldU16("uint16")
				}
			case 6:
				if (e.size / 4) > 1 {
					d.FieldArray("values", func(d *decode.D) {
						array_size := e.size / 4
						for i := uint8(0); i < array_size; i++ {
							d.FieldU32("uint32")
						}
					})
				} else {
					d.FieldU16("uint32")
				}
			case 12:
				if (e.size / 4) > 1 {
					d.FieldArray("values", func(d *decode.D) {
						array_size := e.size / 4
						for i := uint8(0); i < array_size; i++ {
							d.FieldU32("uint32z")
						}
					})
				} else {
					d.FieldU16("uint32z")
				}
			}
		}
	})
}

func decodeDefinitionMessage(d *decode.D) []field_definition {
	d.FieldU8("reserved")
	d.FieldU8("architecture")
	d.FieldU16("global_message_number")
	num_of_fields := d.FieldU8("fields")
	var fields []field_definition
	d.FieldArray("field_definitions", func(d *decode.D) {
		for i := uint64(0); i < num_of_fields; i++ {
			d.FieldStruct("field_definition", func(d *decode.D) {
				fields = append(fields, decodeFieldDefinition(d))
			})
		}
	})
	return fields
}

func decodeFIT(d *decode.D, in interface{}) interface{} {
	d.Endian = decode.LittleEndian
	size := d.FieldU8("header_size")
	if size < 12 {
		d.Fatalf("Header size too small < 12")
	}
	d.FieldU8("protocol_version")
	d.FieldU16("profile_version")
	d.FieldU32("data_size")
	d.FieldUTF8("data_type", 4, d.AssertStr(".FIT"))
	d.FieldU16("crc")
	d.FieldArray("records", func(d *decode.D) {
		var fields []field_definition
		for i := uint64(0); i < 4; i++ {
			d.FieldStruct("record_header", func(d *decode.D) {
				d.FieldU1("normal_header")
				messge_type := d.FieldU1("message_type")
				d.FieldU1("message_type_specific")
				d.FieldU1("reserved")
				d.FieldU4("local_message_type")

				if messge_type == 1 {
					fields = nil
					d.FieldStruct("definition_message", func(d *decode.D) {
						fields = append(fields, decodeDefinitionMessage(d)...)
					})
				} else {
					d.FieldStruct("data_message", func(d *decode.D) {
						decodeDataMessage(d, fields)
					})
				}
			})
		}
	})
	return nil
}
