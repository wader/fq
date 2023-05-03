package wasm

import (
	"github.com/wader/fq/pkg/scalar"
)

var sectionIDToSym = scalar.UintMapSymStr{
	sectionIDCustom:    "custom_section",
	sectionIDType:      "type_section",
	sectionIDImport:    "import_section",
	sectionIDFunction:  "function_section",
	sectionIDTable:     "table_section",
	sectionIDMemory:    "memory_section",
	sectionIDGlobal:    "global_section",
	sectionIDExport:    "export_section",
	sectionIDStart:     "start_section",
	sectionIDElement:   "element_section",
	sectionIDCode:      "code_section",
	sectionIDData:      "data_section",
	sectionIDDataCount: "data_count_section",
}

// A map to convert valtypes to symbols.
//
//	valtype ::= t:numtype => t
//	         |  t:vectype => t
//	         |  t:reftype => t
//
//	numtype ::= 0x7F => i32
//	         |  0x7E => i64
//	         |  0x7D => f32
//	         |  0x7C => f64
//
//	vectype ::= 0x7B => v128
//
//	reftype ::= 0x70 => funcref
//	         |  0x6F => externref
var valtypeToSymMapper = scalar.UintMapSymStr{
	0x7f: "i32",
	0x7e: "i64",
	0x7d: "f32",
	0x7c: "f64",
	0x7b: "v128",
	0x70: "funcref",
	0x6f: "externref",
}

// A map to convert tags of importdesc to symbols.
//
//	importdesc ::= 0x00 x:typeidx     => func x
//	            |  0x01 tt:tabletype  => table tt
//	            |  0x02 mt:memtype    => mem mt
//	            |  0x03 gt:globaltype => global gt
var importdescTagToSym = scalar.UintMapSymStr{
	0x00: "func",
	0x01: "table",
	0x02: "mem",
	0x03: "global",
}

// A map to convert tags of exportdesc to symbols.
//
//	exportdesc ::= 0x00 x:funcidx   => func x
//	            |  0x01 x:tableidx  => table x
//	            |  0x02 x:memidx    => mem x
//	            |  0x03 x:globalidx => global x
var exportdescTagToSym = scalar.UintMapSymStr{
	0x00: "funcidx",
	0x01: "tableidx",
	0x02: "memidx",
	0x03: "globalidx",
}

// A map to convert reftypes to symbols.
//
//	reftype ::= 0x70 => funcref
//	         |  0x6F => externref
var reftypeTagToSym = scalar.UintMapSymStr{
	0x70: "funcref",
	0x6f: "externref",
}

// A map to convert mut to symbols.
//
//	mut ::= 0x00 => const
//	     |  0x01 => var
var mutToSym = scalar.UintMapSymStr{
	0x00: "const",
	0x01: "var",
}

// A map to convert elemkind to symbols.
//
//	elemkind ::= 0x00 => funcref
var elemkindTagToSym = scalar.UintMapSymStr{
	0x00: "funcref",
}
