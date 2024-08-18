package midi

import (
	"github.com/wader/fq/pkg/scalar"
)

const (
	keyCMajor      = 0x0000
	keyGMajor      = 0x0100
	keyDMajor      = 0x0200
	keyAMajor      = 0x0300
	keyEMajor      = 0x0400
	keyBMajor      = 0x0500
	keyFSharpMajor = 0x0600
	keyCSharpMajor = 0x0700
	keyFMajor      = 0xff00
	keyBFlatMajor  = 0xfe00
	keyEFlatMajor  = 0xfd00
	keyAFlatMajor  = 0xfc00
	keyDFlatMajor  = 0xfb00
	keyGFlatMajor  = 0xfa00
	keyCFlatMajor  = 0xf900

	keyAMinor      = 0x0001
	keyEMinor      = 0x0101
	keyBMinor      = 0x0201
	keyFSharpMinor = 0x0301
	keyCSharpMinor = 0x0401
	keyGSharpMinor = 0x0501
	keyDSharpMinor = 0x0601
	keyASharpMinor = 0x0701
	keyDMinor      = 0xff01
	keyGMinor      = 0xfe01
	keyCMinor      = 0xfd01
	keyFMinor      = 0xfc01
	keyBFlatMinor  = 0xfb01
	keyEFlatMinor  = 0xfa01
	keyAFlatMinor  = 0xf901
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
	keyFSharpMinor: "F♯ minor",
	keyCSharpMinor: "C♯ minor",
	keyGSharpMinor: "G♯ minor",
	keyDSharpMinor: "D♯ minor",
	keyASharpMinor: "A♯ minor",
	keyDMinor:      "D minor",
	keyGMinor:      "G minor",
	keyCMinor:      "C minor",
	keyFMinor:      "F minor",
	keyBFlatMinor:  "B♭ minor",
	keyEFlatMinor:  "E♭ minor",
	keyAFlatMinor:  "A♭ minor",
}
