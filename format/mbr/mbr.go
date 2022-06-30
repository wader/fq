package mbr

import (
	"fmt"

	"github.com/wader/fq/format"
	"github.com/wader/fq/format/registry"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/scalar"
)

func init() {
	registry.MustRegister(decode.Format{
		Name:        format.MBR,
		Description: "Master Boot Record",
		DecodeFn:    mbrDecode,
	})
}

func decodePartitionTableEntry(d *decode.D) {
	d.FieldU8("boot_indicator", scalar.UToDescription{
		0x80: "active",
		0x00: "inactive",
	})
	d.FieldStrScalarFn("starting_chs_vals", decodeCHSBytes)
	d.FieldU8("partition_type", partitionTypes)
	d.FieldStrScalarFn("ending_chs_vals", decodeCHSBytes)
	d.FieldStrScalarFn("starting_sector", decodeCHSBytes)
	d.U8() // extra byte
	d.FieldScalarU32("partition_size")
}

// Because this is a fixed-sized table, I am opting to use a
// FieldStruct instead of a FieldArray
func decodePartitionTable(d *decode.D) {
	d.FieldStruct("entry_1", decodePartitionTableEntry)
	d.FieldStruct("entry_2", decodePartitionTableEntry)
	d.FieldStruct("entry_3", decodePartitionTableEntry)
	d.FieldStruct("entry_4", decodePartitionTableEntry)
}

// Source: https://thestarman.pcministry.com/asm/mbr/PartTables.htm#Decoding
func decodeCHSBytes(d *decode.D) scalar.S {
	head, _ := d.Bits(8)
	sectorHighBits, err := d.Bits(2)
	if err != nil {
		d.IOPanic(err, "chs")
	}
	sector, _ := d.Bits(6)
	cylinderLowerBits, err := d.Bits(8)
	if err != nil {
		d.IOPanic(err, "chs")
	}
	cylinder := (sectorHighBits << 2) | cylinderLowerBits
	return scalar.S{Actual: fmt.Sprintf("CHS(%x, %x, %x)", cylinder, head, sector)}
}

func mbrDecode(d *decode.D, in interface{}) interface{} {
	d.Endian = decode.LittleEndian

	d.FieldRawLen("code_area", 446*8)
	d.FieldStruct("partition_table", decodePartitionTable)
	d.FieldU16("boot_record_sig", scalar.ActualHex)
	//d.AssertU(0xaa55)
	return nil
}
