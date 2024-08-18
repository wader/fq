package midi

import (
	"github.com/wader/fq/pkg/scalar"
)

const (
	keyCMajor      = 0x0001
	keyGMajor      = 0x0101
	keyDMajor      = 0x0201
	keyAMajor      = 0x0301
	keyEMajor      = 0x0401
	keyBMajor      = 0x0501
	keyFSharpMajor = 0x0601
	keyCSharpMajor = 0x0701
	keyFMajor      = 0x0801
	keyBFlatMajor  = 0x0901
	keyEFlatMajor  = 0x0a01
	keyAFlatMajor  = 0x0b01
	keyDFlatMajor  = 0x0c01
	keyGFlatMajor  = 0x0d01
	keyCFlatMajor  = 0x0e01

	keyAMinor      = 0x0000
	keyEMinor      = 0x0100
	keyBMinor      = 0x0200
	keyASharpMinor = 0x0300
	keyDSharpMinor = 0x0400
	keyGSharpMinor = 0x0500
	keyCSharpMinor = 0x0600
	keyFSharpMinor = 0x0700
	keyDMinor      = 0x0800
	keyGMinor      = 0x0900
	keyCMinor      = 0x0a00
	keyFMinor      = 0x0b00
	keyBFlatMinor  = 0x0c00
	keyEFlatMinor  = 0x0d00
	keyAFlatMinor  = 0x0e00
)

var keys = scalar.UintMapSymStr{
	keyCMajor:      "C major",
	keyGMajor:      "G major",
	keyDMajor:      "D major",
	keyAMajor:      "A major",
	keyEMajor:      "E major",
	keyBMajor:      "B major",
	keyFSharpMajor: "F♯ major",
	keyCSharpMajor: "C♯ major",
	keyFMajor:      "F major",
	keyBFlatMajor:  "B♭ major",
	keyEFlatMajor:  "E♭ major",
	keyAFlatMajor:  "A♭ major",
	keyDFlatMajor:  "D♭ major",
	keyGFlatMajor:  "G♭ major",
	keyCFlatMajor:  "C♭ major",

	keyAMinor:      "A minor",
	keyEMinor:      "E minor",
	keyBMinor:      "B minor",
	keyASharpMinor: "A♯ minor",
	keyDSharpMinor: "D♯ minor",
	keyGSharpMinor: "G♯ minor",
	keyCSharpMinor: "C♯ minor",
	keyFSharpMinor: "F♯ minor",
	keyDMinor:      "D minor",
	keyGMinor:      "G minor",
	keyCMinor:      "C minor",
	keyFMinor:      "F minor",
	keyBFlatMinor:  "B♭ minor",
	keyEFlatMinor:  "E♭ minor",
	keyAFlatMinor:  "A♭ minor",
}
