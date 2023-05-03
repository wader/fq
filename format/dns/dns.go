package dns

// https://datatracker.ietf.org/doc/html/rfc1035
// https://github.com/Forescout/namewreck/blob/main/rfc/draft-dashevskyi-dnsrr-antipatterns-00.txt

import (
	"net"
	"strings"

	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
	"github.com/wader/fq/pkg/scalar"
)

func init() {
	interp.RegisterFormat(
		format.DNS,
		&decode.Format{
			Description: "DNS packet",
			Groups:      []*decode.Group{format.UDP_Payload},
			DecodeFn:    dnsUDPDecode,
		})
}

const (
	classIN = 1
)

var classNames = scalar.UintRangeToScalar{
	{Range: [2]uint64{0x0000, 0x0000}, S: scalar.Uint{Sym: "reserved", Description: "Reserved"}},
	{Range: [2]uint64{classIN, classIN}, S: scalar.Uint{Sym: "in", Description: "Internet"}},
	{Range: [2]uint64{0x0002, 0x0002}, S: scalar.Uint{Sym: "unassigned", Description: "Unassigned"}},
	{Range: [2]uint64{0x0003, 0x0003}, S: scalar.Uint{Sym: "chaos", Description: "Chaos"}},
	{Range: [2]uint64{0x0004, 0x0004}, S: scalar.Uint{Sym: "hesiod", Description: "Hesiod"}},
	{Range: [2]uint64{0x0005, 0x00fd}, S: scalar.Uint{Sym: "unassigned", Description: "Unassigned"}},
	{Range: [2]uint64{0x00fe, 0x00fe}, S: scalar.Uint{Sym: "qclass_none", Description: "QCLASS NONE"}},
	{Range: [2]uint64{0x00ff, 0x00ff}, S: scalar.Uint{Sym: "qclass_any", Description: "QCLASS ANY"}},
	{Range: [2]uint64{0x0100, 0xfeff}, S: scalar.Uint{Sym: "unassigned", Description: "Unassigned"}},
	{Range: [2]uint64{0xff00, 0xfffe}, S: scalar.Uint{Sym: "private", Description: "Reserved for Private Use"}},
	{Range: [2]uint64{0xffff, 0xffff}, S: scalar.Uint{Sym: "reserved", Description: "Reserved"}},
}

const (
	typeA     = 1
	typeNS    = 2
	typeCNAME = 5
	typeSOA   = 6
	typePTR   = 12
	typeTXT   = 16
	typeAAAA  = 28
)

var typeNames = scalar.UintMapSymStr{
	typeA:     "a",
	typeAAAA:  "aaaa",
	18:        "afsdb",
	42:        "apl",
	257:       "caa",
	60:        "cdnskey",
	59:        "cds",
	37:        "cert",
	typeCNAME: "cname",
	62:        "csync",
	49:        "dhcid",
	32769:     "dlv",
	39:        "dname",
	48:        "dnskey",
	43:        "ds",
	108:       "eui48",
	109:       "eui64",
	13:        "hinfo",
	55:        "hip",
	45:        "ipseckey",
	25:        "key",
	36:        "kx",
	29:        "loc",
	15:        "mx",
	35:        "naptr",
	typeNS:    "ns",
	47:        "nsec",
	50:        "nsec3",
	51:        "nsec3_param",
	61:        "openpgp_key",
	typePTR:   "ptr",
	46:        "rrsig",
	17:        "rp",
	24:        "sig",
	53:        "smimea",
	typeSOA:   "soa",
	33:        "srv",
	44:        "sshfp",
	32768:     "ta",
	249:       "tkey",
	52:        "tlsa",
	250:       "tsig",
	typeTXT:   "txt",
	256:       "uri",
	63:        "zonemd",
	64:        "svcb",
	65:        "https",
}

var rcodeNames = scalar.UintMap{
	0:  {Sym: "no_error", Description: "No error"},
	1:  {Sym: "form_err", Description: "Format error"},
	2:  {Sym: "serv_fail", Description: "Server failure"},
	3:  {Sym: "nx_domain", Description: "Non-Existent Domain"},
	4:  {Sym: "no_tiimpl", Description: "Not implemented"},
	5:  {Sym: "refused", Description: "Refused"},
	6:  {Sym: "yx_domain", Description: "DescriptionName Exists when it should not"}, // RFC 2136
	7:  {Sym: "yxrr_set", Description: "RR Set Exists when it should not"},           // RFC 2136
	8:  {Sym: "nxrr_set", Description: "RR Set that should exist does not"},          // RFC 2136
	9:  {Sym: "not_auth", Description: "Server Not Authoritative for zone"},          // RFC 2136
	10: {Sym: "not_zone", Description: "Name not contained in zone"},                 // RFC 2136
	// collision in RFCs
	// 16: {Sym: "badvers", Description: "Bad OPT Version"},           // RFC 2671
	16: {Sym: "bad_sig", Description: "TSIG Signature Failure"},        // RFC 2845
	17: {Sym: "bad_key", Description: "Key not recognized"},            // RFC 2845
	18: {Sym: "bad_time", Description: "Signature out of time window"}, // RFC 2845
	19: {Sym: "bad_mode", Description: "Bad TKEY Mode"},                // RFC 2930
	20: {Sym: "bad_name", Description: "Duplicate key name"},           // RFC 2930
	21: {Sym: "bad_alg", Description: "Algorithm not supported"},       // RFC 2930
}

func decodeAStr(d *decode.D) string {
	return net.IP(d.BytesLen(4)).String()
}

func decodeAAAAStr(d *decode.D) string {
	return net.IP(d.BytesLen(16)).String()
}

func fieldDecodeLabel(d *decode.D, pointerOffset int64, name string) {
	var endPos int64
	const maxJumps = 100
	jumpCount := 0

	d.FieldStruct(name, func(d *decode.D) {
		var ls []string
		d.FieldArray("labels", func(d *decode.D) {
			seenTermintor := false
			for !seenTermintor {
				d.FieldStruct("label", func(d *decode.D) {
					if d.PeekUintBits(2) == 0b11 {
						d.FieldU2("is_pointer")
						pointer := d.FieldU14("pointer")
						if endPos == 0 {
							endPos = d.Pos()
						}
						jumpCount++
						if jumpCount > maxJumps {
							d.Fatalf("label has more than %d jumps", maxJumps)
						}
						d.SeekAbs(int64(pointer)*8 + pointerOffset)
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

func dnsDecodeRR(d *decode.D, pointerOffset int64, resp bool, count uint64, name string, structName string) {
	d.FieldArray(name, func(d *decode.D) {
		for i := uint64(0); i < count; i++ {
			d.FieldStruct(structName, func(d *decode.D) {
				fieldDecodeLabel(d, pointerOffset, "name")
				typ := d.FieldU16("type", typeNames)
				class := d.FieldU16("class", classNames)
				if resp {
					d.FieldU32("ttl")
					rdLength := d.FieldU16("rdlength")
					d.FramedFn(int64(rdLength)*8, func(d *decode.D) {
						// TODO: all only for classIN?
						switch {
						case class == classIN && typ == typeA:
							d.FieldStrFn("address", decodeAStr)
						case typ == typeNS:
							fieldDecodeLabel(d, pointerOffset, "ns")
						case typ == typeCNAME:
							fieldDecodeLabel(d, pointerOffset, "cname")
						case typ == typeSOA:
							fieldDecodeLabel(d, pointerOffset, "mname")
							fieldDecodeLabel(d, pointerOffset, "rname")
							d.FieldU32("serial")
							d.FieldU32("refresh")
							d.FieldU32("retry")
							d.FieldU32("expire")
							d.FieldU32("minimum")
						case typ == typePTR:
							fieldDecodeLabel(d, pointerOffset, "ptr")
						case typ == typeTXT:
							var ss []string
							d.FieldStruct("txt", func(d *decode.D) {
								d.FieldArray("strings", func(d *decode.D) {
									for !d.End() {
										ss = append(ss, d.FieldUTF8ShortString("string"))
									}
								})
								d.FieldValueStr("value", strings.Join(ss, ""))
							})
						case class == classIN && typ == typeAAAA:
							d.FieldStrFn("address", decodeAAAAStr)
						default:
							d.FieldUTF8("rdata", int(rdLength))
						}
					})
				}
			})
		}
	})
}

func dnsDecode(d *decode.D, hasLengthHeader bool) any {
	pointerOffset := int64(0)
	d.FieldStruct("header", func(d *decode.D) {
		if hasLengthHeader {
			pointerOffset = 16
			d.FieldU16("length")
		}
		d.FieldU16("id")
		d.FieldU1("qr", scalar.UintMapSymStr{
			0: "query",
			1: "response",
		})
		d.FieldU4("opcode", scalar.UintMapSymStr{
			0: "query",
			1: "iquery",
			2: "status",
			4: "notify", // RFC 1996
			5: "update", // RFC 2136
		})
		d.FieldBool("authoritative_answer")
		d.FieldBool("truncation")
		d.FieldBool("recursion_desired")
		d.FieldBool("recursion_available")
		d.FieldU3("z")
		d.FieldU4("rcode", rcodeNames)
	})

	qdCount := d.FieldU16("qd_count")
	anCount := d.FieldU16("an_count")
	nsCount := d.FieldU16("ns_count")
	arCount := d.FieldU16("ar_count")
	dnsDecodeRR(d, pointerOffset, false, qdCount, "questions", "question")
	dnsDecodeRR(d, pointerOffset, true, anCount, "answers", "answer")
	dnsDecodeRR(d, pointerOffset, true, nsCount, "nameservers", "nameserver")
	dnsDecodeRR(d, pointerOffset, true, arCount, "additionals", "additional")

	return nil
}

func dnsUDPDecode(d *decode.D) any {
	var upi format.UDP_Payload_In
	if d.ArgAs(&upi) {
		upi.MustIsPort(d.Fatalf, format.UDPPortDomain, format.UDPPortMDNS)
	}

	return dnsDecode(d, false)
}
