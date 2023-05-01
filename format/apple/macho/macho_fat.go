package macho

// https://github.com/aidansteele/osx-abi-macho-file-format-reference

import (
	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
	"github.com/wader/fq/pkg/scalar"
)

var machoFormat decode.Group

func init() {
	interp.RegisterFormat(
		format.MachO_Fat,
		&decode.Format{
			Description: "Fat Mach-O macOS executable (multi-architecture)",
			Groups:      []*decode.Group{format.Probe},
			DecodeFn:    machoFatDecode,
			Dependencies: []decode.Dependency{
				{Groups: []*decode.Group{format.MachO}, Out: &machoFormat},
			},
		})
}

const FAT_MAGIC = 0xcafe_babe

func machoFatDecode(d *decode.D) any {
	type ofile struct {
		offset int64
		size   int64
	}
	var ofiles []ofile

	d.FieldStruct("fat_header", func(d *decode.D) {
		d.FieldU32("magic", magicSymMapper, scalar.UintHex, d.UintAssert(FAT_MAGIC))

		narchs := d.FieldU32("narchs")
		d.FieldArray("archs", func(d *decode.D) {
			for i := 0; i < int(narchs); i++ {
				d.FieldStruct("arch", func(d *decode.D) {
					// beware cputype and cpusubtype changes from ofile header to fat header
					cpuType := d.FieldU32("cputype", cpuTypes, scalar.UintHex)
					d.FieldU32("cpusubtype", cpuSubTypes[cpuType], scalar.UintHex)
					offset := d.FieldU32("offset", scalar.UintHex)
					size := d.FieldU32("size")
					d.FieldU32("align")

					ofiles = append(ofiles, ofile{offset: int64(offset), size: int64(size)})
				})
			}
		})
	})

	d.FieldArray("files", func(d *decode.D) {
		for _, o := range ofiles {
			d.RangeFn(o.offset*8, o.size*8, func(d *decode.D) {
				d.FieldFormat("file", &machoFormat, nil)
			})
		}
	})

	return nil
}
