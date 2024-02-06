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
	"golang.org/x/text/encoding"
)

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
		var byte = bytes[i]
		var tmp = fitCRCTable[crc&0xF]
		crc = (crc >> 4) & 0x0FFF
		crc = crc ^ tmp ^ fitCRCTable[byte&0xF]
		tmp = fitCRCTable[crc&0xF]
		crc = (crc >> 4) & 0x0FFF
		crc = crc ^ tmp ^ fitCRCTable[(byte>>4)&0xF]
	}

	return crc
}

type fitContext struct {
	dataSize   int
	headerSize int
}

type dataRecordContext struct {
	compressed         bool
	data               bool
	localMessageType   int
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

type fieldDef struct {
	name   string
	typ    string
	format string
	unit   string
	scale  float64
	offset int64
	size   int
}

type devFieldDefMap map[uint64]map[uint64]fieldDef
type localFieldDefMap map[uint64]map[uint64]fieldDef
type localMsgIsDevDef map[uint64]bool

func fitDecodeFileHeader(d *decode.D, fc *fitContext) {
	frameStart := d.Pos()

	// d.FieldStruct("ident", func(d *decode.D) {
	headerSize := d.FieldU8("headerSize")
	d.FieldU8("protocolVersion")
	d.FieldU16("profileVersion")
	dataSize := d.FieldU32("dataSize")

	d.FieldRawLen("dataType", 4*8, d.AssertBitBuf([]byte(".FIT")))
	if headerSize == 14 {
		headerCRC := calcCRC(d.BytesRange(frameStart, int(headerSize)-2))
		d.FieldU16("crc", d.UintValidate(uint64(headerCRC)))
	}
	fc.headerSize = int(headerSize)
	fc.dataSize = int(dataSize)
}

func fitDecodeDataRecordHeader(d *decode.D, drc *dataRecordContext) {
	drc.compressed = d.FieldBool("normalHeader", scalar.BoolMapDescription{false: "Normal header",
		true: "Compressed header"})
	if drc.compressed {
		localMessageType := d.FieldU2("localMessageType")
		d.FieldU32("timeOffset")
		drc.localMessageType = int(localMessageType)
		drc.data = true
	} else {
		mTypeIsDef := d.FieldBool("messageType", scalar.BoolMap{true: {Sym: 1, Description: "Definition message"},
			false: {Sym: 0, Description: "Data message"}})
		hasDeveloperFields := d.FieldBool("hasDeveloperFields")
		d.FieldBool("reserved")
		localMessageType := d.FieldU4("localMessageType")

		drc.hasDeveloperFields = hasDeveloperFields
		drc.localMessageType = int(localMessageType)
		drc.data = !mTypeIsDef
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
		isDevMap[uint64(drc.localMessageType)] = true
	} else {
		isDevMap[uint64(drc.localMessageType)] = false
	}
	numFields := d.FieldU8("fields")
	lmfd[uint64(drc.localMessageType)] = make(map[uint64]fieldDef, numFields)
	d.FieldArray("fieldDefinitions", func(d *decode.D) {
		for i := 0; i < int(numFields); i++ {
			d.FieldStruct("fieldDefinition", func(d *decode.D) {
				fieldDefNo := d.FieldU8("fieldDefNo", mappers.FieldDefMap[messageNo])
				size := d.FieldU8("size")
				baseType := d.FieldU8("baseType", mappers.TypeDefMap["fit_base_type"])

				var typ = mappers.TypeDefMap["fit_base_type"][baseType].Name
				fDefLookup, isSet := mappers.FieldDefMap[messageNo][fieldDefNo]
				if isSet {
					var foundName = fDefLookup.Name
					lmfd[uint64(drc.localMessageType)][uint64(i)] = fieldDef{name: foundName, typ: typ, size: int(size), format: fDefLookup.Type, unit: fDefLookup.Unit, scale: fDefLookup.Scale, offset: fDefLookup.Offset}
				} else {
					var foundName = fmt.Sprintf("UNKOWN_%d", fieldDefNo)
					lmfd[uint64(drc.localMessageType)][uint64(i)] = fieldDef{name: foundName, typ: typ, size: int(size), format: "unknown"}
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

					//baseType := d.FieldU8("baseType", mappers.TypeDefMap["fit_base_type"])

					//var typ = mappers.TypeDefMap["fit_base_type"][baseType].Name
					typ := dmfd[devDataIdx][fieldDefNo].typ
					//fDefLookup, isSet := mappers.FieldDefMap[messageNo][fieldDefNo]
					fDefLookup, isSet := dmfd[devDataIdx][fieldDefNo]
					if isSet {
						var foundName = fDefLookup.name
						lmfd[uint64(drc.localMessageType)][uint64(i)] = fieldDef{name: foundName, typ: typ, size: int(size), unit: fDefLookup.unit, scale: fDefLookup.scale, offset: fDefLookup.offset}
					} else {
						var foundName = fmt.Sprintf("UNKOWN_%d", fieldDefNo)
						lmfd[uint64(drc.localMessageType)][uint64(i)] = fieldDef{name: foundName, typ: typ, size: int(size), format: "unknown"}
					}
				})
			}
		})
	}

}

func ensureDevFieldMap(dmfd devFieldDefMap, devIdx uint64, fieldDefNo uint64) {
	_, devIsSet := dmfd[devIdx]
	if !devIsSet {
		dmfd[devIdx] = make(map[uint64]fieldDef)
	}
}

func fieldUint(fieldFn func(string, ...scalar.UintMapper) uint64, d *decode.D, fdc *fileDescriptionContext, expectedSize int, fDef fieldDef, uintFormatter scalar.UintFn) {
	var val uint64
	if fDef.size != expectedSize {
		d.FieldStr(fDef.name, fDef.size, encoding.Nop)
	} else {
		if uintFormatter != nil {
			val = fieldFn(fDef.name, uintFormatter)
		} else {
			val = fieldFn(fDef.name)
		}

		switch fDef.name {
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

func fieldSint(fieldFn func(string, ...scalar.SintMapper) int64, d *decode.D, fdc *fileDescriptionContext, expectedSize int, fDef fieldDef, sintFormatter scalar.SintFn) {
	//var val uint64
	if fDef.size != expectedSize {
		d.FieldStr(fDef.name, fDef.size, encoding.Nop)
	} else {
		if sintFormatter != nil {
			fieldFn(fDef.name, sintFormatter)
		} else {
			fieldFn(fDef.name)
		}
		//setDevFieldUint(isDevDep, dmfd, fDef.name, val);
	}
}

func fieldFloat(fieldFn func(string, ...scalar.FltMapper) float64, d *decode.D, fdc *fileDescriptionContext, expectedSize int, fDef fieldDef) {
	//var val uint64
	if fDef.size != expectedSize {
		d.FieldStr(fDef.name, fDef.size, encoding.Nop)
	} else {
		fieldFn(fDef.name)

		//setDevFieldUint(isDevDep, dmfd, fDef.name, val);
	}
}

func fieldString(d *decode.D, fdc *fileDescriptionContext, fDef fieldDef) {
	val := d.FieldUTF8NullFixedLen(fDef.name, fDef.size)

	switch fDef.name {
	case "field_name":
		fdc.name = val
	case "units":
		fdc.unit = val
	}
}

func fitDecodeDataMessage(d *decode.D, drc *dataRecordContext, lmfd localFieldDefMap, dmfd devFieldDefMap, isDevMap localMsgIsDevDef) {
	var fdc fileDescriptionContext
	keys := make([]int, len(lmfd[uint64(drc.localMessageType)]))
	i := 0
	for k := range lmfd[uint64(drc.localMessageType)] {
		keys[i] = int(k)
		i++
	}
	sort.Ints(keys)

	isDevDep := isDevMap[uint64(drc.localMessageType)]

	for _, k := range keys {
		fDef := lmfd[uint64(drc.localMessageType)][uint64(k)]

		var uintFormatter = mappers.GetUintFormatter(fDef.format, fDef.unit, fDef.scale, fDef.offset)
		var sintFormatter = mappers.GetSintFormatter(fDef.format, fDef.unit, fDef.scale, fDef.offset)

		switch fDef.typ {
		// case "byte":
		// d.FieldStr(fDef.name, fDef.size, encoding.Nop)
		case "enum", "uint8", "uint8z", "byte":
			fieldUint(d.FieldU8, d, &fdc, 1, fDef, uintFormatter)
		case "sint8":
			fieldSint(d.FieldS8, d, &fdc, 1, fDef, sintFormatter)
		case "sint16":
			fieldSint(d.FieldS16, d, &fdc, 2, fDef, sintFormatter)
		case "uint16", "uint16z":
			fieldUint(d.FieldU16, d, &fdc, 2, fDef, uintFormatter)
		case "sint32":
			fieldSint(d.FieldS32, d, &fdc, 4, fDef, sintFormatter)
		case "uint32", "uint32z":
			fieldUint(d.FieldU32, d, &fdc, 4, fDef, uintFormatter)
		case "float32":
			fieldFloat(d.FieldF32, d, &fdc, 4, fDef)
		case "float64":
			fieldFloat(d.FieldF64, d, &fdc, 8, fDef)
		case "sint64":
			fieldSint(d.FieldS64, d, &fdc, 4, fDef, sintFormatter)
		case "uint64", "uint64z":
			fieldUint(d.FieldU64, d, &fdc, 8, fDef, uintFormatter)
		case "string":
			fieldString(d, &fdc, fDef)
		default:
			d.Fatalf("Unknown type %s", fDef.typ)
		}
	}

	if isDevDep {
		ensureDevFieldMap(dmfd, fdc.devIdx, fdc.fDefNo)
		dmfd[fdc.devIdx][fdc.fDefNo] = fieldDef{name: fdc.name, typ: fdc.typ, unit: fdc.unit, scale: 0, offset: 0}
	}
}

func decodeFIT(d *decode.D) any {
	var fc fitContext
	d.Endian = decode.LittleEndian

	var lmfd localFieldDefMap = make(localFieldDefMap)
	var dmfd devFieldDefMap = make(devFieldDefMap)
	var isDevMap localMsgIsDevDef = make(localMsgIsDevDef)

	//decodeBSONDocument(d)
	d.FieldStruct("header", func(d *decode.D) { fitDecodeFileHeader(d, &fc) })

	d.FieldArray("dataRecords", func(d *decode.D) {
		// 	// headerPos := d.Pos()

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
		fileCRC = calcCRC(d.BytesRange(0, fc.dataSize+fc.headerSize)) // 12 byte header - CRC whole file except the CRC itself
	} else {
		fileCRC = calcCRC(d.BytesRange(14*8, fc.dataSize)) // 14 byte header - CRC everything below header except the CRC itself
	}
	d.FieldU16("crc", d.UintValidate(uint64(fileCRC)))

	return nil
}
