package fit

// https://developer.garmin.com/fit/protocol/
// https://developer.garmin.com/fit/cookbook/decoding-activity-files/

import (
	"embed"
	"fmt"
	"slices"

	"github.com/wader/fq/format"
	"github.com/wader/fq/format/fit/mappers"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
	"github.com/wader/fq/pkg/scalar"
)

//go:embed fit.md
var fitFS embed.FS

func init() {
	interp.RegisterFormat(
		format.FIT,
		&decode.Format{
			Description: "Garmin Flexible and Interoperable Data Transfer",
			Groups:      []*decode.Group{format.Probe},
			DecodeFn:    decodeFIT,
		})

	interp.RegisterFS(fitFS)
}

var fitCRCTable = [16]uint16{
	0x0000, 0xcc01, 0xd801, 0x1400, 0xf001, 0x3c00, 0x2800, 0xe401,
	0xa001, 0x6c00, 0x7800, 0xb401, 0x5000, 0x9c01, 0x8801, 0x4400,
}

func calcCRC(bytes []byte) uint16 {
	var crc uint16
	crc = 0
	for i := 0; i < len(bytes); i++ {
		// compute checksum of lower four bits of byte
		var checkByte = bytes[i]
		var tmp = fitCRCTable[crc&0xf]
		crc = (crc >> 4) & 0x0fff
		crc = crc ^ tmp ^ fitCRCTable[checkByte&0xf]
		tmp = fitCRCTable[crc&0xf]
		crc = (crc >> 4) & 0x0fff
		crc = crc ^ tmp ^ fitCRCTable[(checkByte>>4)&0xf]
	}

	return crc
}

type fitContext struct {
	dataSize   uint64
	headerSize uint64
}

type dataRecordContext struct {
	compressed         bool
	data               bool
	localMessageType   uint64
	hasDeveloperFields bool
}

type fileDescriptionContext struct {
	devIdx        uint64
	fDefNo        uint64
	typ           string
	name          string
	unit          string
	nativeFieldNo uint64
	nativeMsgNo   uint64
}

type valueType struct {
	value uint64
	typ   string
}

type devFieldDefMap map[uint64]map[uint64]mappers.FieldDef
type localFieldDefMap map[uint64]map[uint64]mappers.LocalFieldDef
type localMsgIsDevDef map[uint64]bool
type valueMap map[string]valueType

// expected size in bytes
var expectedSizeMap = map[string]uint64{
	"byte":    1,
	"enum":    1,
	"float32": 4,
	"float64": 8,
	"sint8":   1,
	"sint16":  2,
	"sint32":  4,
	"sint64":  8,
	"uint8":   1,
	"uint8z":  1,
	"uint16":  2,
	"uint16z": 2,
	"uint32":  4,
	"uint32z": 4,
	"uint64":  8,
	"uint64z": 8,
}

// "Magic" numbers
const (
	developerFieldDescMesgNo = 206 // Special data message used as dynamic field definition message
)

func fitDecodeFileHeader(d *decode.D, fc *fitContext) {
	frameStart := d.Pos()

	headerSize := d.FieldU8("header_size")
	d.FieldU8("protocol_version")
	d.FieldU16("profile_version")
	dataSize := d.FieldU32("data_size")

	d.FieldRawLen("data_type", 4*8, d.AssertBitBuf([]byte(".FIT")))
	if headerSize == 14 {
		headerCRC := calcCRC(d.BytesRange(frameStart, int(headerSize)-2))
		d.FieldU16("crc", d.UintValidate(uint64(headerCRC)), scalar.UintHex)
	}
	fc.headerSize = headerSize
	fc.dataSize = dataSize
}

func fitDecodeDataRecordHeader(d *decode.D, drc *dataRecordContext) {
	headerType := d.FieldU1("header_type", scalar.UintMapSymStr{0: "normal", 1: "compressed"})
	drc.compressed = headerType == 1
	if drc.compressed {
		localMessageType := d.FieldU2("local_message_type")
		d.FieldU32("time_offset")
		drc.localMessageType = localMessageType
		drc.data = true
	} else {
		mTypeIsDef := d.FieldU1("message_type", scalar.UintMapSymStr{0: "data", 1: "definition"})
		hasDeveloperFields := d.FieldBool("has_developer_fields")
		d.FieldBool("reserved")
		localMessageType := d.FieldU4("local_message_type")

		drc.hasDeveloperFields = hasDeveloperFields
		drc.localMessageType = localMessageType
		drc.data = mTypeIsDef == 0
	}
}

func fitDecodeDefinitionMessage(d *decode.D, drc *dataRecordContext, lmfd localFieldDefMap, dmfd devFieldDefMap, isDevMap localMsgIsDevDef) {
	d.FieldU8("reserved")
	endian := d.FieldU8("architecture", scalar.UintMapSymStr{0: "little_endian", 1: "big_endian"})
	switch endian {
	case 0:
		d.Endian = decode.LittleEndian
	case 1:
		d.Endian = decode.BigEndian
	default:
		d.Fatalf("Unknown endian %d", endian)
	}

	messageNo := d.FieldU16("global_message_number", mappers.TypeDefMap["mesg_num"])
	isDevMap[drc.localMessageType] = messageNo == developerFieldDescMesgNo

	numFields := d.FieldU8("fields")
	lmfd[drc.localMessageType] = make(map[uint64]mappers.LocalFieldDef, numFields)

	d.FieldArray("field_definitions", func(d *decode.D) {
		for i := uint64(0); i < numFields; i++ {
			d.FieldStruct("field_definition", func(d *decode.D) {
				fieldDefNo := d.FieldU8("field_definition_number", mappers.FieldDefMap[messageNo])
				size := d.FieldU8("size")
				baseType := d.FieldU8("base_type", mappers.TypeDefMap["fit_base_type"])

				var typ = mappers.TypeDefMap["fit_base_type"][baseType].Name
				fDefLookup, isSet := mappers.FieldDefMap[messageNo][fieldDefNo]
				if isSet {
					var foundName = fDefLookup.Name
					lmfd[drc.localMessageType][i] = mappers.LocalFieldDef{
						Name:             foundName,
						Type:             typ,
						Size:             size,
						Format:           fDefLookup.Type,
						Unit:             fDefLookup.Unit,
						Scale:            fDefLookup.Scale,
						Offset:           fDefLookup.Offset,
						GlobalFieldDef:   fDefLookup,
						GlobalMessageNo:  messageNo,
						GlobalFieldDefNo: fieldDefNo,
					}
				} else {
					var foundName = fmt.Sprintf("UNKNOWN_%d", fieldDefNo)
					lmfd[drc.localMessageType][i] = mappers.LocalFieldDef{
						Name:   foundName,
						Type:   typ,
						Size:   size,
						Format: "unknown",
					}
				}
			})
		}
	})
	if drc.hasDeveloperFields {
		numDevFields := d.FieldU8("developer_fields")

		d.FieldArray("developer_field_definitions", func(d *decode.D) {
			for i := numFields; i < (numDevFields + numFields); i++ {
				d.FieldStruct("developer_field_definition", func(d *decode.D) {
					fieldDefNo := d.FieldU8("field_definition_number")
					size := d.FieldU8("size")
					devDataIdx := d.FieldU8("developer_data_index")

					fDefLookup, isSet := dmfd[devDataIdx][fieldDefNo]

					if isSet {
						var foundName = fDefLookup.Name
						lmfd[drc.localMessageType][i] = mappers.LocalFieldDef{
							Name:   foundName,
							Type:   fDefLookup.Type,
							Size:   size,
							Unit:   fDefLookup.Unit,
							Scale:  fDefLookup.Scale,
							Offset: fDefLookup.Offset,
						}
					} else {
						var foundName = fmt.Sprintf("UNKNOWN_%d", fieldDefNo)
						lmfd[drc.localMessageType][i] = mappers.LocalFieldDef{
							Name:   foundName,
							Type:   "UNKNOWN",
							Size:   size,
							Format: "UNKNOWN",
						}
					}
				})
			}
		})
	}
}

func readUint(d *decode.D, fDef mappers.LocalFieldDef, valMap valueMap) {
	expectedSize := expectedSizeMap[fDef.Type]
	if fDef.Size != expectedSize {
		d.SeekRel(int64(fDef.Size) * 8) // skip over array types since they cannot be referenced by subfields
	} else {
		val := d.U(int(expectedSize) * 8)
		valMap[fDef.Name] = valueType{value: val, typ: fDef.Format}
	}
}

func fieldUint(d *decode.D, fDef mappers.LocalFieldDef, uintFormatter scalar.UintFn, fdc *fileDescriptionContext, valMap valueMap) {
	var val uint64
	expectedSize := expectedSizeMap[fDef.Type]

	if fDef.Size != expectedSize {
		arrayCount := fDef.Size / expectedSize
		for i := uint64(0); i < arrayCount; i++ {
			d.FieldU(fmt.Sprintf("%s_%d", fDef.Name, i), int(expectedSize)*8, uintFormatter)
		}
	} else {
		if fDef.GlobalFieldDef.HasSubField {
			var found = false
			if subFieldValueMap, ok := mappers.SubFieldDefMap[fDef.GlobalMessageNo][fDef.GlobalFieldDefNo]; ok {
				for k := range subFieldValueMap {
					if subFieldDef, ok := subFieldValueMap[k][mappers.TypeDefMap[valMap[k].typ][valMap[k].value].Name]; ok {
						subUintFormatter := mappers.GetUintFormatter(mappers.LocalFieldDef{
							Name:   subFieldDef.Name,
							Type:   fDef.Type,
							Size:   fDef.Size,
							Format: subFieldDef.Type,
							Unit:   subFieldDef.Unit,
							Scale:  subFieldDef.Scale,
							Offset: subFieldDef.Offset,
						})
						val = d.FieldU(subFieldDef.Name, int(expectedSize)*8, subUintFormatter)
						found = true
						continue
					}
				}
			}
			if !found { // SubField conditions could not be resolved
				val = d.FieldU(fDef.Name, int(expectedSize)*8, uintFormatter)
			}
		} else {
			val = d.FieldU(fDef.Name, int(expectedSize)*8, uintFormatter)
		}

		// Save developer field definitions
		switch fDef.Name {
		case "developer_data_index":
			fdc.devIdx = val
		case "field_definition_number":
			fdc.fDefNo = val
		case "fit_base_type_id":
			fdc.typ = mappers.TypeDefMap["fit_base_type"][val].Name
		case "native_field_num":
			fdc.nativeFieldNo = val
		case "native_mesg_num":
			fdc.nativeMsgNo = val
		}
	}
}

func fieldSint(d *decode.D, fDef mappers.LocalFieldDef, sintFormatter scalar.SintFn) {
	expectedSize := expectedSizeMap[fDef.Type]
	if fDef.Size != expectedSize {
		arrayCount := fDef.Size / expectedSize
		for i := uint64(0); i < arrayCount; i++ {
			d.FieldS(fmt.Sprintf("%s_%d", fDef.Name, i), int(expectedSize)*8, sintFormatter)
		}
	} else {
		d.FieldS(fDef.Name, int(expectedSize)*8, sintFormatter)
	}
}

func fieldFloat(d *decode.D, fDef mappers.LocalFieldDef, floatFormatter scalar.FltFn) {
	expectedSize := expectedSizeMap[fDef.Type]
	if fDef.Size != expectedSize {
		arrayCount := fDef.Size / expectedSize
		for i := uint64(0); i < arrayCount; i++ {
			d.FieldF(fmt.Sprintf("%s_%d", fDef.Name, i), int(expectedSize)*8, floatFormatter)
		}
	} else {
		d.FieldF(fDef.Name, int(expectedSize)*8)
	}
}

func fieldString(d *decode.D, fDef mappers.LocalFieldDef, fdc *fileDescriptionContext) {
	val := d.FieldUTF8NullFixedLen(fDef.Name, int(fDef.Size), scalar.StrMapDescription{"": "Invalid"})

	// Save developer field definitions
	switch fDef.Name {
	case "field_name":
		fdc.name = val
	case "units":
		fdc.unit = val
	}
}

func fitDecodeDataMessage(d *decode.D, drc *dataRecordContext, lmfd localFieldDefMap, dmfd devFieldDefMap, isDevMap localMsgIsDevDef) {
	var fdc fileDescriptionContext
	valMap := make(valueMap, len(lmfd[drc.localMessageType]))
	keys := make([]int, len(lmfd[drc.localMessageType]))
	i := 0
	for k := range lmfd[drc.localMessageType] {
		keys[i] = int(k)
		i++
	}
	slices.Sort(keys)

	isDevDep := isDevMap[drc.localMessageType]

	// pre read all integer fields and store them in the value map
	// so that they can be referenced by eventual subfields
	curPos := d.Pos()
	for _, k := range keys {
		fDef := lmfd[drc.localMessageType][uint64(k)]
		switch fDef.Type {
		case "enum", "byte", "uint8", "uint8z", "uint16", "uint16z", "uint32", "uint32z", "uint64", "uint64z":
			readUint(d, fDef, valMap)
		default:
			d.SeekRel(int64(fDef.Size) * 8)
		}
	}
	d.SeekRel(curPos - d.Pos())

	for _, k := range keys {
		fDef := lmfd[drc.localMessageType][uint64(k)]

		var uintFormatter = mappers.GetUintFormatter(fDef)
		var sintFormatter = mappers.GetSintFormatter(fDef)
		var floatFormatter = mappers.GetFloatFormatter(fDef)

		switch fDef.Type {
		case "enum", "byte", "uint8", "uint8z", "uint16", "uint16z", "uint32", "uint32z", "uint64", "uint64z":
			fieldUint(d, fDef, uintFormatter, &fdc, valMap)
		case "sint8", "sint16", "sint32", "sint64":
			fieldSint(d, fDef, sintFormatter)
		case "float32", "float64":
			fieldFloat(d, fDef, floatFormatter)
		case "string":
			fieldString(d, fDef, &fdc)
		default:
			d.Fatalf("Unknown type %s", fDef.Type)
		}
	}

	if isDevDep {
		if _, ok := dmfd[fdc.devIdx]; !ok {
			dmfd[fdc.devIdx] = make(map[uint64]mappers.FieldDef)
		}
		dmfd[fdc.devIdx][fdc.fDefNo] = mappers.FieldDef{
			Name:   fdc.name,
			Type:   fdc.typ,
			Unit:   fdc.unit,
			Scale:  0,
			Offset: 0,
		}
	}
}

func decodeFIT(d *decode.D) any {
	var fc fitContext
	d.Endian = decode.LittleEndian

	var lmfd localFieldDefMap = make(localFieldDefMap)
	var dmfd devFieldDefMap = make(devFieldDefMap)
	var isDevMap localMsgIsDevDef = make(localMsgIsDevDef)

	d.FieldStruct("header", func(d *decode.D) { fitDecodeFileHeader(d, &fc) })

	d.FieldArray("data_records", func(d *decode.D) {
		for d.Pos() < int64((fc.headerSize+fc.dataSize)*8) {
			d.FieldStruct("data_record", func(d *decode.D) {
				var drc dataRecordContext
				d.FieldStruct("record_header", func(d *decode.D) { fitDecodeDataRecordHeader(d, &drc) })
				switch drc.data {
				case true:
					d.FieldStruct("data_message", func(d *decode.D) { fitDecodeDataMessage(d, &drc, lmfd, dmfd, isDevMap) })
				case false:
					d.FieldStruct("definition_message", func(d *decode.D) { fitDecodeDefinitionMessage(d, &drc, lmfd, dmfd, isDevMap) })
				}
			})
		}
	})

	var fileCRC uint16
	if fc.headerSize == 12 {
		fileCRC = calcCRC(d.BytesRange(0, int(fc.dataSize+fc.headerSize))) // 12 byte header - CRC whole file except the CRC itself
	} else {
		fileCRC = calcCRC(d.BytesRange(14*8, int(fc.dataSize))) // 14 byte header - CRC everything below header except the CRC itself
	}
	d.FieldU16("crc", d.UintValidate(uint64(fileCRC)), scalar.UintHex)

	return nil
}
