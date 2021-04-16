package dns

// TODO: https://github.com/Forescout/namewreck/blob/main/rfc/draft-dashevskyi-dnsrr-antipatterns-00.txt

import (
	"fmt"
	"fq/pkg/decode"
	"fq/pkg/format"
	"strings"
)

func init() {
	format.MustRegister(&decode.Format{
		Name:        format.DNS,
		Description: "DNS packet",
		DecodeFn:    dnsDecode,
	})
}

// TODO: type consts
// TODO: aaaa,a rddata

var classNames = map[[2]uint64]string{
	{0x0000, 0x0000}: "Reserved",
	{0x0001, 0x0001}: "IN",
	{0x0002, 0x0002}: "Unassigned",
	{0x0003, 0x0003}: "Chaos",
	{0x0004, 0x0004}: "Hesiod",
	{0x0005, 0x00fd}: "Unassigned",
	{0x00fe, 0x00fe}: "QCLASS NONE",
	{0x00ff, 0x00ff}: "QCLASS ANY",
	{0x0100, 0xfeff}: "Unassigned",
	{0xff00, 0xfffe}: "Reserved for Private Use",
	{0xffff, 0xffff}: "Reserved",
}

const (
	typeCNAME = 5
)

var typeNames = map[uint64]string{
	1:         "A",
	28:        "AAAA",
	18:        "AFSDB",
	42:        "APL",
	257:       "CAA",
	60:        "CDNSKEY",
	59:        "CDS",
	37:        "CERT",
	typeCNAME: "CNAME",
	62:        "CSYNC",
	49:        "DHCID",
	32769:     "DLV",
	39:        "DNAME",
	48:        "DNSKEY",
	43:        "DS",
	108:       "EUI48",
	109:       "EUI64",
	13:        "HINFO",
	55:        "HIP",
	45:        "IPSECKEY",
	25:        "KEY",
	36:        "KX",
	29:        "LOC",
	15:        "MX",
	35:        "NAPTR",
	2:         "NS",
	47:        "NSEC",
	50:        "NSEC3",
	51:        "NSEC3PARAM",
	61:        "OPENPGPKEY",
	12:        "PTR",
	46:        "RRSIG",
	17:        "RP",
	24:        "SIG",
	53:        "SMIMEA",
	6:         "SOA",
	33:        "SRV",
	44:        "SSHFP",
	32768:     "TA",
	249:       "TKEY",
	52:        "TLSA",
	250:       "TSIG",
	16:        "TXT",
	256:       "URI",
	63:        "ZONEMD",
	64:        "SVCB",
	65:        "HTTPS",
}

var rcodeNames = map[uint64]string{
	0: "No error",
	1: "Format error",
	2: "Server failure",
	3: "Name error",
	4: "Not implemented",
	5: "Refused",
}

func fieldDecodeLabel(d *decode.D, name string) {
	var endPos int64
	const maxJumps = 1000
	jumpCount := 0

	d.FieldStructFn(name, func(d *decode.D) {
		var ls []string
		d.FieldArrayFn("labels", func(d *decode.D) {
			seenTermintor := false
			for !seenTermintor {
				d.FieldStructFn("label", func(d *decode.D) {
					if d.PeekBits(2) == 0b11 {
						d.FieldU2("is_pointer")
						pointer := d.FieldU14("pointer")
						if endPos == 0 {
							endPos = d.Pos()
						}
						jumpCount++
						if jumpCount > maxJumps {
							d.Invalid(fmt.Sprintf("label has more than %d jumps", maxJumps))
						}
						d.SeekAbs(int64(pointer * 8))
					}

					l := d.FieldU8("length")
					if l == 0 {
						seenTermintor = true
						return
					}
					ls = append(ls, d.FieldUTF8("value", int(l)))
				})
			}
		})
		d.FieldValueStr("value", strings.Join(ls, "."), "")
	})

	if endPos != 0 {
		d.SeekAbs(endPos)
	}
}

func fieldDecodeRR(d *decode.D, count uint64, name string, structName string) {
	d.FieldArrayFn(name, func(d *decode.D) {
		for i := uint64(0); i < count; i++ {
			d.FieldStructFn(structName, func(d *decode.D) {
				fieldDecodeLabel(d, "name")
				typ, _ := d.FieldStringMapFn("type", typeNames, "Unknown", d.U16, decode.NumberDecimal)
				d.FieldStringRangeMapFn("class", classNames, "Unknown", d.U16, decode.NumberDecimal)
				d.FieldU32("ttl")
				// TODO: pointer?
				rdLength := d.FieldU16("rd_length")

				switch typ {
				case typeCNAME:
					fieldDecodeLabel(d, "cname")
				default:
					d.FieldUTF8("rddata", int(rdLength))
				}

			})
		}
	})
}

func dnsDecode(d *decode.D, in interface{}) interface{} {
	d.FieldStructFn("header", func(d *decode.D) {
		d.FieldU16("id")
		d.FieldBool("query")
		d.FieldU4("opcode")
		d.FieldBool("authoritative_answer")
		d.FieldBool("truncation")
		d.FieldBool("recursion_desired")
		d.FieldBool("recursion_available")
		d.FieldU3("z")
		d.FieldStringMapFn("rcode", rcodeNames, "Unknown", d.U4, decode.NumberDecimal)
	})

	qdCount := d.FieldU16("qd_count")
	anCount := d.FieldU16("an_count")
	nsCount := d.FieldU16("ns_count")
	arCount := d.FieldU16("ar_count")

	d.FieldArrayFn("questions", func(d *decode.D) {
		for i := uint64(0); i < qdCount; i++ {
			d.FieldStructFn("question", func(d *decode.D) {
				fieldDecodeLabel(d, "name")
				d.FieldStringMapFn("type", typeNames, "Unknown", d.U16, decode.NumberDecimal)
				d.FieldStringRangeMapFn("class", classNames, "Unknown", d.U16, decode.NumberDecimal)
			})
		}
	})

	fieldDecodeRR(d, anCount, "answers", "answer")
	fieldDecodeRR(d, nsCount, "nameservers", "nameserver")
	fieldDecodeRR(d, arCount, "additionals", "additional")

	return nil
}
