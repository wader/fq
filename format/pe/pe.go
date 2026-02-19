package pe

import (
	"fmt"

	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
)

var msDosStubGroup decode.Group
var coffGroup decode.Group

func init() {
	interp.RegisterFormat(
		format.PE,
		&decode.Format{
			Description: "Portable Executable",
			Groups:      []*decode.Group{format.Probe},
			Dependencies: []decode.Dependency{
				{Groups: []*decode.Group{format.MSDOS_Stub}, Out: &msDosStubGroup},
				{Groups: []*decode.Group{format.COFF}, Out: &coffGroup},
			},
			DecodeFn: peDecode,
		})
}

func peDecode(d *decode.D) any {
	_, v := d.FieldFormat("ms_dos_stub", &msDosStubGroup, nil)
	msDOSOut, ok := v.(format.MS_DOS_Out)
	if !ok {
		panic(fmt.Sprintf("expected MS_DOS_Out got %#+v", v))
	}
	d.FieldFormat("coff", &coffGroup, format.COFF_In{FilePointerOffset: msDOSOut.LFANew})

	return nil
}
