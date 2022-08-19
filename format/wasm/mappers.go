package wasm

import (
	"github.com/wader/fq/pkg/scalar"
)

var sectionIDToSym = scalar.UToSymStr{
	sectionIDCustom:    "custom section",
	sectionIDType:      "type section",
	sectionIDImport:    "import section",
	sectionIDFunction:  "function section",
	sectionIDTable:     "table section",
	sectionIDMemory:    "memory section",
	sectionIDGlobal:    "global section",
	sectionIDExport:    "export section",
	sectionIDStart:     "start section",
	sectionIDElement:   "element section",
	sectionIDCode:      "code section",
	sectionIDData:      "data section",
	sectionIDDataCount: "data count section",
}

// valtype ::= t:numtype => t
//          |  t:vectype => t
//          |  t:reftype => t
//
// numtype ::= 0x7F => i32
//          |  0x7E => i64
//          |  0x7D => f32
//          |  0x7C => f64
//
// vectype ::= 0x7B => v128
//
// reftype ::= 0x70 => funcref
//          |  0x6F => externref
var valtypeToSymMapper = scalar.UToSymStr{
	0x7f: "i32",
	0x7e: "i64",
	0x7d: "f32",
	0x7c: "f64",
	0x7b: "v128",
	0x70: "funcref",
	0x6f: "externref",
}

// importdesc ::= 0x00 x:typeidx     => func x
//             |  0x01 tt:tabletype  => table tt
//             |  0x02 mt:memtype    => mem mt
//             |  0x03 gt:globaltype => global gt
var importdescTagToSym = scalar.UToSymStr{
	0x00: "func",
	0x01: "table",
	0x02: "mem",
	0x03: "global",
}

// exportdesc ::= 0x00 x:funcidx   => func x
//             |  0x01 x:tableidx  => table x
//             |  0x02 x:memidx    => mem x
//             |  0x03 x:globalidx => global x
var exportdescTagToSym = scalar.UToSymStr{
	0x00: "funcidx",
	0x01: "tableidx",
	0x02: "memidx",
	0x03: "globalidx",
}

// reftype ::= 0x70 => funcref
//          |  0x6F => externref
var reftypeTagToSym = scalar.UToSymStr{
	0x70: "funcref",
	0x6f: "externref",
}

// mut ::= 0x00 => const
//      |  0x01 => var
var mutToSym = scalar.UToSymStr{
	0x00: "const",
	0x01: "var",
}

// elemkind ::= 0x00 => funcref
var elemkindTagToSym = scalar.UToSymStr{
	0x00: "funcref",
}
