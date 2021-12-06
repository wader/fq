package pcap

import (
	"fmt"

	"github.com/wader/fq/format"
	"github.com/wader/fq/format/inet/flowsdecoder"
	"github.com/wader/fq/pkg/bitio"
	"github.com/wader/fq/pkg/decode"
)

// TODO: is shared between pcap and pcapng
var linkToFormat = map[int]*decode.Group{
	format.LinkTypeETHERNET:   &pcapngEther8023Format,
	format.LinkTypeLINUX_SLL:  &pcapngSLLPacketFormat,
	format.LinkTypeLINUX_SLL2: &pcapngSLL2PacketFormat,
}

var linkToDecodeFn = map[int]func(fd *flowsdecoder.Decoder, bs []byte) error{
	format.LinkTypeETHERNET:  (*flowsdecoder.Decoder).EthernetFrame,
	format.LinkTypeLINUX_SLL: (*flowsdecoder.Decoder).SLLPacket,
	format.LinkTypeLINUX_SLL2: func(fd *flowsdecoder.Decoder, bs []byte) error {
		if len(bs) < 20 {
			// TODO: too short sll packet, error somehow?
			return fmt.Errorf("packet too short %d", len(bs))
		}

		// TODO: gopacket does not support SLL2 atm so convert SLL to SSL2
		nbs := []byte{
			0, bs[10], // packet type
			bs[8], bs[9], // arphdr
			0, bs[11], // link layer address length
			bs[12], bs[13], bs[14], bs[15], bs[16], bs[17], bs[18], bs[19], //  link layer address
			bs[0], bs[1], // protocol type
		}
		nbs = append(nbs, bs[20:]...)

		return fd.SLLPacket(nbs)
	},
}

func fieldFlows(d *decode.D, fd *flowsdecoder.Decoder, tcpStreamFormat decode.Group, ipv4PacketFormat decode.Group) {
	d.FieldArray("ipv4_reassembled", func(d *decode.D) {
		for _, p := range fd.IPV4Reassembled {
			bb := bitio.NewBufferFromBytes(p.Datagram, -1)
			if dv, _, _ := d.TryFieldFormatBitBuf(
				"ipv4_packet",
				bb,
				ipv4PacketFormat,
				nil,
			); dv == nil {
				d.FieldRootBitBuf("ipv4_packet", bb)
			}
		}
	})

	d.FieldArray("tcp_connections", func(d *decode.D) {
		for _, s := range fd.TCPConnections {
			d.FieldStruct("flow", func(d *decode.D) {
				d.FieldValueStr("source_ip", s.ClientEndpoint.IP.String())
				d.FieldValueU("source_port", uint64(s.ClientEndpoint.Port), format.TCPPortMap)
				d.FieldValueStr("destination_ip", s.ServerEndpoint.IP.String())
				d.FieldValueU("destination_port", uint64(s.ServerEndpoint.Port), format.TCPPortMap)
				csBB := bitio.NewBufferFromBytes(s.ClientToServer.Bytes(), -1)
				if dv, _, _ := d.TryFieldFormatBitBuf(
					"client_stream",
					csBB,
					tcpStreamFormat,
					format.TCPStreamIn{
						SourcePort:      s.ClientEndpoint.Port,
						DestinationPort: s.ServerEndpoint.Port,
					},
				); dv == nil {
					d.FieldRootBitBuf("client_stream", csBB)
				}

				scBB := bitio.NewBufferFromBytes(s.ServerToClient.Bytes(), -1)
				if dv, _, _ := d.TryFieldFormatBitBuf(
					"server_stream",
					scBB,
					tcpStreamFormat,
					format.TCPStreamIn{
						SourcePort:      s.ClientEndpoint.Port,
						DestinationPort: s.ServerEndpoint.Port,
					},
				); dv == nil {
					d.FieldRootBitBuf("server_stream", scBB)
				}
			})
		}
	})
}
