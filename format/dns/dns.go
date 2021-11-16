package dns

// TODO: https://github.com/Forescout/namewreck/blob/main/rfc/draft-dashevskyi-dnsrr-antipatterns-00.txt

import (
	"fmt"
	"net"
	"strings"

	"github.com/wader/fq/format"
	"github.com/wader/fq/format/registry"
	"github.com/wader/fq/pkg/decode"
)

func init() {
	registry.MustRegister(&decode.Format{
		Name:        format.DNS,
		Description: "DNS packet",
		DecodeFn:    dnsDecode,
	})
}

const (
	classIN = 1
)

var classNames = map[[2]uint64]decode.Scalar{
	{0x0000, 0x0000}:   {Sym: "Reserved", Description: "Reserved"},
	{classIN, classIN}: {Sym: "IN", Description: "Internet"},
	{0x0002, 0x0002}:   {Sym: "Unassigned", Description: "Unassigned"},
	{0x0003, 0x0003}:   {Sym: "Chaos", Description: "Chaos"},
	{0x0004, 0x0004}:   {Sym: "Hesiod", Description: "Hesiod"},
	{0x0005, 0x00fd}:   {Sym: "Unassigned", Description: "Unassigned"},
	{0x00fe, 0x00fe}:   {Sym: "QCLASS_NONE", Description: "QCLASS NONE"},
	{0x00ff, 0x00ff}:   {Sym: "QCLASS_ANY", Description: "QCLASS ANY"},
	{0x0100, 0xfeff}:   {Sym: "Unassigned", Description: "Unassigned"},
	{0xff00, 0xfffe}:   {Sym: "Private", Description: "Reserved for Private Use"},
	{0xffff, 0xffff}:   {Sym: "Reserved", Description: "Reserved"},
}

const (
	typeA     = 1
	typeAAAA  = 28
	typeCNAME = 5
)

var typeNames = decode.UToStr{
	typeA:     "A",
	typeAAAA:  "AAAA",
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

var rcodeNames = decode.UToScalar{
	0:  {Sym: "NoError", Description: "No error"},
	1:  {Sym: "FormErr", Description: "Format error"},
	2:  {Sym: "ServFail", Description: "Server failure"},
	3:  {Sym: "NXDomain", Description: "Non-Existent Domain"},
	4:  {Sym: "NotiImpl", Description: "Not implemented"},
	5:  {Sym: "Refused", Description: "Refused"},
	6:  {Sym: "YXDomain", Description: "DescriptionName Exists when it should not"}, // RFC 2136
	7:  {Sym: "YXRRSet", Description: "RR Set Exists when it should not"},           // RFC 2136
	8:  {Sym: "NXRRSet", Description: "RR Set that should exist does not"},          // RFC 2136
	9:  {Sym: "NotAuth", Description: "Server Not Authoritative for zone"},          // RFC 2136
	10: {Sym: "NotZone", Description: "Name not contained in zone"},                 // RFC 2136
	// collision in RFCs
	// 16: {Sym: "BADVERS", Description: "Bad OPT Version"},                            // RFC 2671
	16: {Sym: "BADSIG", Description: "TSIG Signature Failure"},        // RFC 2845
	17: {Sym: "BADKEY", Description: "Key not recognized"},            // RFC 2845
	18: {Sym: "BADTIME", Description: "Signature out of time window"}, // RFC 2845
	19: {Sym: "BADMODE", Description: "Bad TKEY Mode"},                // RFC 2930
	20: {Sym: "BADNAME", Description: "Duplicate key name"},           // RFC 2930
	21: {Sym: "BADALG", Description: "Algorithm not supported"},       // RFC 2930
}

func decodeINAStr(d *decode.D) string {
	return net.IP(d.BytesLen(4)).String()
}

func decodeINAAAAStr(d *decode.D) string {
	return net.IP(d.BytesLen(16)).String()
}

func fieldFormatLabel(d *decode.D, name string) {
	var endPos int64
	const maxJumps = 1000
	jumpCount := 0

	d.FieldStruct(name, func(d *decode.D) {
		var ls []string
		d.FieldArray("labels", func(d *decode.D) {
			seenTermintor := false
			for !seenTermintor {
				d.FieldStruct("label", func(d *decode.D) {
					if d.PeekBits(2) == 0b11 {
						d.FieldU2("is_pointer")
						pointer := d.FieldU14("pointer")
						if endPos == 0 {
							endPos = d.Pos()
						}
						jumpCount++
						if jumpCount > maxJumps {
							d.Fatal(fmt.Sprintf("label has more than %d jumps", maxJumps))
						}
						d.SeekAbs(int64(pointer) * 8)
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
		d.FieldValueStr("value", strings.Join(ls, "."))
	})

	if endPos != 0 {
		d.SeekAbs(endPos)
	}
}

func fieldFormatRR(d *decode.D, count uint64, name string, structName string) {
	d.FieldArray(name, func(d *decode.D) {
		for i := uint64(0); i < count; i++ {
			d.FieldStruct(structName, func(d *decode.D) {
				fieldFormatLabel(d, "name")
				typ := d.FieldU16("type", d.MapUToStr(typeNames))
				class := d.FieldU16("class", d.MapURangeToScalar(classNames))
				d.FieldU32("ttl")
				// TODO: pointer?
				rdLength := d.FieldU16("rdlength")

				switch {
				case typ == typeCNAME:
					fieldFormatLabel(d, "cname")
				case typ == typeA && class == classIN:
					d.FieldStrFn("address", decodeINAStr)
				case typ == typeAAAA && class == classIN:
					d.FieldStrFn("address", decodeINAAAAStr)
				default:
					d.FieldUTF8("rdata", int(rdLength))
				}
			})
		}
	})
}

func dnsDecode(d *decode.D, in interface{}) interface{} {
	d.FieldStruct("header", func(d *decode.D) {
		d.FieldU16("id")
		d.FieldBool("query", d.MapBoolToStr(decode.BoolToStr{
			true:  "Query",
			false: "Response",
		}))
		d.FieldU4("opcode", d.MapUToStr(decode.UToStr{
			0: "Query",
			1: "IQuery",
			2: "Status",
			4: "Notify", // RFC 1996
			5: "Update", // RFC 2136
		}))
		d.FieldBool("authoritative_answer")
		d.FieldBool("truncation")
		d.FieldBool("recursion_desired")
		d.FieldBool("recursion_available")
		d.FieldU3("z")
		d.FieldU4("rcode", d.MapUToScalar(rcodeNames))
	})

	qdCount := d.FieldU16("qd_count")
	anCount := d.FieldU16("an_count")
	nsCount := d.FieldU16("ns_count")
	arCount := d.FieldU16("ar_count")

	d.FieldArray("questions", func(d *decode.D) {
		for i := uint64(0); i < qdCount; i++ {
			d.FieldStruct("question", func(d *decode.D) {
				fieldFormatLabel(d, "name")
				d.FieldU16("type", d.MapUToStr(typeNames))
				d.FieldU16("class", d.MapURangeToScalar(classNames))
			})
		}
	})

	fieldFormatRR(d, anCount, "answers", "answer")
	fieldFormatRR(d, nsCount, "nameservers", "nameserver")
	fieldFormatRR(d, arCount, "additionals", "additional")

	return nil
}
