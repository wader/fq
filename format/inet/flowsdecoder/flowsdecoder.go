package flowsdecoder

import (
	"bytes"
	"encoding/binary"
	"net"

	"github.com/google/gopacket"
	"github.com/google/gopacket/ip4defrag"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/reassembly"
)

type IPEndpoint struct {
	IP   net.IP
	Port int
}

type TCPConnection struct {
	ClientEndpoint IPEndpoint
	ServerEnpoint  IPEndpoint
	ClientStream   *bytes.Buffer
	ServerStream   *bytes.Buffer

	tcpstate       *reassembly.TCPSimpleFSM
	optchecker     reassembly.TCPOptionCheck
	net, transport gopacket.Flow
}

func (t *TCPConnection) Accept(tcp *layers.TCP, ci gopacket.CaptureInfo, dir reassembly.TCPFlowDirection, nextSeq reassembly.Sequence, start *bool, ac reassembly.AssemblerContext) bool {
	// has ok state?
	if !t.tcpstate.CheckState(tcp, dir) {
		// TODO: handle err?
		return false
	}
	// has ok options?
	if err := t.optchecker.Accept(tcp, ci, dir, nextSeq, start); err != nil {
		// TODO: handle err?
		return false
	}
	// TODO: checksum?

	// accept
	return true
}

func (t *TCPConnection) ReassembledSG(sg reassembly.ScatterGather, ac reassembly.AssemblerContext) {
	dir, _, _, _ := sg.Info()
	length, _ := sg.Lengths()

	data := sg.Fetch(length)

	switch dir {
	case reassembly.TCPDirClientToServer:
		t.ClientStream.Write(data)
	case reassembly.TCPDirServerToClient:
		t.ServerStream.Write(data)
	}
}

func (t *TCPConnection) ReassemblyComplete(ac reassembly.AssemblerContext) bool {
	// do not remove the connection to allow last ACK
	return false
}

type IPV4Reassembled struct {
	SourceIP      net.IP
	DestinationIP net.IP
	Datagram      []byte
}

func (fd *Decoder) New(net, transport gopacket.Flow, tcp *layers.TCP, ac reassembly.AssemblerContext) reassembly.Stream {
	fsmOptions := reassembly.TCPSimpleFSMOptions{
		SupportMissingEstablishment: true,
	}
	// TODO: get ip layer somehow?
	stream := &TCPConnection{
		ClientEndpoint: IPEndpoint{
			IP:   append([]byte(nil), net.Src().Raw()...),
			Port: int(binary.BigEndian.Uint16(transport.Src().Raw())),
		},
		ServerEnpoint: IPEndpoint{
			IP:   append([]byte(nil), net.Dst().Raw()...),
			Port: int(binary.BigEndian.Uint16(transport.Dst().Raw())),
		},
		ClientStream: &bytes.Buffer{},
		ServerStream: &bytes.Buffer{},

		net:        net,
		transport:  transport,
		tcpstate:   reassembly.NewTCPSimpleFSM(fsmOptions),
		optchecker: reassembly.NewTCPOptionCheck(),
	}

	fd.TCPConnections = append(fd.TCPConnections, stream)

	return stream
}

type Decoder struct {
	TCPConnections []*TCPConnection
	IPV4Reassbled  []IPV4Reassembled

	ipv4Defrag   *ip4defrag.IPv4Defragmenter
	tcpAssembler *reassembly.Assembler
}

func New() *Decoder {
	flowDecoder := &Decoder{}
	streamPool := reassembly.NewStreamPool(flowDecoder)
	tcpAssembler := reassembly.NewAssembler(streamPool)
	flowDecoder.tcpAssembler = tcpAssembler
	flowDecoder.ipv4Defrag = ip4defrag.NewIPv4Defragmenter()

	return flowDecoder
}

func (fd *Decoder) SLLPacket(bs []byte) {
	fd.packet(gopacket.NewPacket(bs, layers.LayerTypeLinuxSLL, gopacket.Lazy))
}

func (fd *Decoder) EthernetFrame(bs []byte) {
	fd.packet(gopacket.NewPacket(bs, layers.LayerTypeEthernet, gopacket.Lazy))
}

func (fd *Decoder) packet(p gopacket.Packet) {
	// TODO: linkType
	ip4Layer := p.Layer(layers.LayerTypeIPv4)
	if ip4Layer != nil {
		ip4, _ := ip4Layer.(*layers.IPv4)
		l := ip4.Length
		newip4, err := fd.ipv4Defrag.DefragIPv4(ip4)
		if err != nil {
			panic(err)
		} else if newip4 != nil {
			// TODO: correct way to detect finished reassemble?
			if newip4.Length != l {
				// TODO: better way to reconstruct package?
				sb := gopacket.NewSerializeBuffer()
				b, _ := sb.PrependBytes(len(newip4.Payload))
				copy(b, newip4.Payload)
				_ = newip4.SerializeTo(sb, gopacket.SerializeOptions{
					FixLengths:       true,
					ComputeChecksums: true,
				})

				fd.IPV4Reassbled = append(fd.IPV4Reassbled, IPV4Reassembled{
					SourceIP:      ip4.SrcIP,
					DestinationIP: ip4.DstIP,
					Datagram:      sb.Bytes(),
				})

				pb, ok := p.(gopacket.PacketBuilder)
				if !ok {
					panic("not a PacketBuilder")
				}
				nextDecoder := newip4.NextLayerType()
				_ = nextDecoder.Decode(newip4.Payload, pb)
			}
		}
	}

	tcp := p.Layer(layers.LayerTypeTCP)
	if tcp != nil {
		tcp, _ := tcp.(*layers.TCP)
		fd.tcpAssembler.Assemble(p.NetworkLayer().NetworkFlow(), tcp)
	}
}

func (fd *Decoder) Flush() {
	fd.tcpAssembler.FlushAll()
}
