package fit

import (
	"embed"
	"fmt"
	"sort"

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
			Description: "ANT FIT",
			Groups:      []*decode.Group{format.Probe},
			DecodeFn:    decodeFIT,
		})

	interp.RegisterFS(fitFS)
}

var fitCRCTable = [16]uint16{
	0x0000, 0xCC01, 0xD801, 0x1400, 0xF001, 0x3C00, 0x2800, 0xE401,
	0xA001, 0x6C00, 0x7800, 0xB401, 0x5000, 0x9C01, 0x8801, 0x4400,
}

func calcCRC(bytes []byte) uint16 {
	var crc uint16
	crc = 0
	for i := 0; i < len(bytes); i++ {
		// compute checksum of lower four bits of byte
		var checkByte = bytes[i]
		var tmp = fitCRCTable[crc&0xF]
		crc = (crc >> 4) & 0x0FFF
		crc = crc ^ tmp ^ fitCRCTable[checkByte&0xF]
		tmp = fitCRCTable[crc&0xF]
		crc = (crc >> 4) & 0x0FFF
		crc = crc ^ tmp ^ fitCRCTable[(checkByte>>4)&0xF]
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

type devFieldDefMap map[uint64]map[uint64]mappers.FieldDef
type localFieldDefMap map[uint64]map[uint64]mappers.FieldDef
type localMsgIsDevDef map[uint64]bool

func fitDecodeFileHeader(d *decode.D, fc *fitContext) {
	frameStart := d.Pos()

	headerSize := d.FieldU8("headerSize")
	d.FieldU8("protocolVersion")
	d.FieldU16("profileVersion")
	dataSize := d.FieldU32("dataSize")

	d.FieldRawLen("dataType", 4*8, d.AssertBitBuf([]byte(".FIT")))
	if headerSize == 14 {
		headerCRC := calcCRC(d.BytesRange(frameStart, int(headerSize)-2))
		d.FieldU16("crc", d.UintValidate(uint64(headerCRC)))
	}
	fc.headerSize = headerSize
	fc.dataSize = dataSize
}

func fitDecodeDataRecordHeader(d *decode.D, drc *dataRecordContext) {
	headerType := d.FieldU1("headerType", scalar.UintMapDescription{0: "Normal header", 1: "Compressed header"})
	drc.compressed = headerType == 1
	if drc.compressed {
		localMessageType := d.FieldU2("localMessageType")
		d.FieldU32("timeOffset")
		drc.localMessageType = localMessageType
		drc.data = true
	} else {
		mTypeIsDef := d.FieldU1("messageType", scalar.UintMapDescription{0: "Data message", 1: "Definition message"})
		hasDeveloperFields := d.FieldBool("hasDeveloperFields")
		d.FieldBool("reserved")
		localMessageType := d.FieldU4("localMessageType")

		drc.hasDeveloperFields = hasDeveloperFields
		drc.localMessageType = localMessageType
		drc.data = mTypeIsDef == 0
	}
}

func fitDecodeDefinitionMessage(d *decode.D, drc *dataRecordContext, lmfd localFieldDefMap, dmfd devFieldDefMap, isDevMap localMsgIsDevDef) {
	d.FieldU8("reserved")
	endian := d.FieldU8("architecture")
	switch endian {
	case 0:
		d.Endian = decode.LittleEndian
	case 1:
		d.Endian = decode.BigEndian
	default:
		d.Fatalf("Unknown endian %d", endian)
	}
	messageNo := d.FieldU16("globalMessageNumber", mappers.TypeDefMap["mesg_num"])
	if messageNo == 206 { // developer field_description
		isDevMap[drc.localMessageType] = true
	} else {
		isDevMap[drc.localMessageType] = false
	}
	numFields := d.FieldU8("fields")
	lmfd[drc.localMessageType] = make(map[uint64]mappers.FieldDef, numFields)
	d.FieldArray("fieldDefinitions", func(d *decode.D) {
		for i := uint64(0); i < numFields; i++ {
			d.FieldStruct("fieldDefinition", func(d *decode.D) {
				fieldDefNo := d.FieldU8("fieldDefNo", mappers.FieldDefMap[messageNo])
				size := d.FieldU8("size")
				baseType := d.FieldU8("baseType", mappers.TypeDefMap["fit_base_type"])

				var typ = mappers.TypeDefMap["fit_base_type"][baseType].Name
				fDefLookup, isSet := mappers.FieldDefMap[messageNo][fieldDefNo]
				if isSet {
					var foundName = fDefLookup.Name
					lmfd[drc.localMessageType][i] = mappers.FieldDef{Name: foundName, Type: typ, Size: size, Format: fDefLookup.Type, Unit: fDefLookup.Unit, Scale: fDefLookup.Scale, Offset: fDefLookup.Offset}
				} else {
					var foundName = fmt.Sprintf("UNKNOWN_%d", fieldDefNo)
					lmfd[drc.localMessageType][i] = mappers.FieldDef{Name: foundName, Type: typ, Size: size, Format: "unknown"}
				}
			})
		}
	})
	if drc.hasDeveloperFields {
		numDevFields := d.FieldU8("devFields")

		d.FieldArray("devFieldDefinitions", func(d *decode.D) {
			for i := numFields; i < (numDevFields + numFields); i++ {
				d.FieldStruct("devFieldDefinition", func(d *decode.D) {
					fieldDefNo := d.FieldU8("fieldDefNo")
					size := d.FieldU8("size")
					devDataIdx := d.FieldU8("devDataIdx")

					fDefLookup, isSet := dmfd[devDataIdx][fieldDefNo]

					if isSet {
						var foundName = fDefLookup.Name
						lmfd[drc.localMessageType][i] = mappers.FieldDef{Name: foundName, Type: fDefLookup.Type, Size: size, Unit: fDefLookup.Unit, Scale: fDefLookup.Scale, Offset: fDefLookup.Offset}
					} else {
						var foundName = fmt.Sprintf("UNKNOWN_%d", fieldDefNo)
						lmfd[drc.localMessageType][i] = mappers.FieldDef{Name: foundName, Type: "UNKNOWN", Size: size, Format: "unknown"}
					}
				})
			}
		})
	}

}

func ensureDevFieldMap(dmfd devFieldDefMap, devIdx uint64) {
	_, devIsSet := dmfd[devIdx]

	if !devIsSet {
		dmfd[devIdx] = make(map[uint64]mappers.FieldDef)
	}
}

func fieldUint(fieldFn func(string, ...scalar.UintMapper) uint64, expectedSize uint64, fDef mappers.FieldDef, uintFormatter scalar.UintFn, fdc *fileDescriptionContext) {
	var val uint64

	if fDef.Size != expectedSize {
		arrayCount := fDef.Size / expectedSize
		for i := uint64(0); i < arrayCount; i++ {
			fieldFn(fmt.Sprintf("%s_%d", fDef.Name, i), uintFormatter)
		}
	} else {
		if uintFormatter != nil {
			val = fieldFn(fDef.Name, uintFormatter)
		} else {
			val = fieldFn(fDef.Name)
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

func fieldSint(fieldFn func(string, ...scalar.SintMapper) int64, expectedSize uint64, fDef mappers.FieldDef, sintFormatter scalar.SintFn) {
	if fDef.Size != expectedSize {
		arrayCount := fDef.Size / expectedSize
		for i := uint64(0); i < arrayCount; i++ {
			fieldFn(fmt.Sprintf("%s_%d", fDef.Name, i), sintFormatter)
		}
	} else {
		if sintFormatter != nil {
			fieldFn(fDef.Name, sintFormatter)
		} else {
			fieldFn(fDef.Name)
		}
	}
}

func fieldFloat(fieldFn func(string, ...scalar.FltMapper) float64, expectedSize uint64, fDef mappers.FieldDef, floatFormatter scalar.FltFn) {
	if fDef.Size != expectedSize {
		arrayCount := fDef.Size / expectedSize
		for i := uint64(0); i < arrayCount; i++ {
			fieldFn(fmt.Sprintf("%s_%d", fDef.Name, i), floatFormatter)
		}
	} else {
		fieldFn(fDef.Name)
	}
}

func fieldString(d *decode.D, fDef mappers.FieldDef, fdc *fileDescriptionContext) {
	val := d.FieldUTF8NullFixedLen(fDef.Name, int(fDef.Size), scalar.StrMapSymStr{"": "[invalid]"})

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
	keys := make([]int, len(lmfd[drc.localMessageType]))
	i := 0
	for k := range lmfd[drc.localMessageType] {
		keys[i] = int(k)
		i++
	}
	sort.Ints(keys)

	isDevDep := isDevMap[drc.localMessageType]

	for _, k := range keys {
		fDef := lmfd[drc.localMessageType][uint64(k)]

		var uintFormatter = mappers.GetUintFormatter(fDef)
		var sintFormatter = mappers.GetSintFormatter(fDef)
		var floatFormatter = mappers.GetFloatFormatter(fDef)

		switch fDef.Type {
		case "enum", "uint8", "uint8z", "byte":
			fieldUint(d.FieldU8, 1, fDef, uintFormatter, &fdc)
		case "uint16", "uint16z":
			fieldUint(d.FieldU16, 2, fDef, uintFormatter, &fdc)
		case "uint32", "uint32z":
			fieldUint(d.FieldU32, 4, fDef, uintFormatter, &fdc)
		case "uint64", "uint64z":
			fieldUint(d.FieldU64, 8, fDef, uintFormatter, &fdc)
		case "sint8":
			fieldSint(d.FieldS8, 1, fDef, sintFormatter)
		case "sint16":
			fieldSint(d.FieldS16, 2, fDef, sintFormatter)
		case "sint32":
			fieldSint(d.FieldS32, 4, fDef, sintFormatter)
		case "sint64":
			fieldSint(d.FieldS64, 8, fDef, sintFormatter)
		case "float32":
			fieldFloat(d.FieldF32, 4, fDef, floatFormatter)
		case "float64":
			fieldFloat(d.FieldF64, 8, fDef, floatFormatter)
		case "string":
			fieldString(d, fDef, &fdc)
		default:
			d.Fatalf("Unknown type %s", fDef.Type)
		}
	}

	if isDevDep {
		ensureDevFieldMap(dmfd, fdc.devIdx)
		dmfd[fdc.devIdx][fdc.fDefNo] = mappers.FieldDef{Name: fdc.name, Type: fdc.typ, Unit: fdc.unit, Scale: 0, Offset: 0}
	}
}

func decodeFIT(d *decode.D) any {
	var fc fitContext
	d.Endian = decode.LittleEndian

	var lmfd localFieldDefMap = make(localFieldDefMap)
	var dmfd devFieldDefMap = make(devFieldDefMap)
	var isDevMap localMsgIsDevDef = make(localMsgIsDevDef)

	d.FieldStruct("header", func(d *decode.D) { fitDecodeFileHeader(d, &fc) })

	d.FieldArray("dataRecords", func(d *decode.D) {
		for d.Pos() < int64((fc.headerSize+fc.dataSize)*8) {
			d.FieldStruct("dataRecord", func(d *decode.D) {
				var drc dataRecordContext
				d.FieldStruct("dataRecordHeader", func(d *decode.D) { fitDecodeDataRecordHeader(d, &drc) })
				switch drc.data {
				case true:
					d.FieldStruct("dataMessage", func(d *decode.D) { fitDecodeDataMessage(d, &drc, lmfd, dmfd, isDevMap) })
				case false:
					d.FieldStruct("definitionMessage", func(d *decode.D) { fitDecodeDefinitionMessage(d, &drc, lmfd, dmfd, isDevMap) })
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
	d.FieldU16("crc", d.UintValidate(uint64(fileCRC)))

	return nil
}
