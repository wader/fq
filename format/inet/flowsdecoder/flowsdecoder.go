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
	ServerEndpoint IPEndpoint
	ClientToServer *bytes.Buffer
	ServerToClient *bytes.Buffer

	tcpState   *reassembly.TCPSimpleFSM
	optChecker reassembly.TCPOptionCheck
	net        gopacket.Flow
	transport  gopacket.Flow
}

func (t *TCPConnection) Accept(tcp *layers.TCP, ci gopacket.CaptureInfo, dir reassembly.TCPFlowDirection, nextSeq reassembly.Sequence, start *bool, ac reassembly.AssemblerContext) bool {
	// has ok state?
	if !t.tcpState.CheckState(tcp, dir) {
		// TODO: handle err?
		return false
	}
	// has ok options?
	if err := t.optChecker.Accept(tcp, ci, dir, nextSeq, start); err != nil {
		// TODO: handle err?
		return false
	}
	// TODO: checksum?

	// accept
	return true
}

func (t *TCPConnection) ReassembledSG(sg reassembly.ScatterGather, ac reassembly.AssemblerContext) {
	dir, _, _, skip := sg.Info()
	length, _ := sg.Lengths()

	if skip != 0 {
		// stream has missing bytes
		return
	}

	data := sg.Fetch(length)

	switch dir {
	case reassembly.TCPDirClientToServer:
		t.ClientToServer.Write(data)
	case reassembly.TCPDirServerToClient:
		t.ServerToClient.Write(data)
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
		ServerEndpoint: IPEndpoint{
			IP:   append([]byte(nil), net.Dst().Raw()...),
			Port: int(binary.BigEndian.Uint16(transport.Dst().Raw())),
		},
		ClientToServer: &bytes.Buffer{},
		ServerToClient: &bytes.Buffer{},

		net:        net,
		transport:  transport,
		tcpState:   reassembly.NewTCPSimpleFSM(fsmOptions),
		optChecker: reassembly.NewTCPOptionCheck(),
	}

	fd.TCPConnections = append(fd.TCPConnections, stream)

	return stream
}

type Decoder struct {
	TCPConnections  []*TCPConnection
	IPV4Reassembled []IPV4Reassembled

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

func (fd *Decoder) SLLPacket(bs []byte) error {
	return fd.packet(gopacket.NewPacket(bs, layers.LayerTypeLinuxSLL, gopacket.Lazy))
}

func (fd *Decoder) EthernetFrame(bs []byte) error {
	return fd.packet(gopacket.NewPacket(bs, layers.LayerTypeEthernet, gopacket.Lazy))
}

func (fd *Decoder) packet(p gopacket.Packet) error {
	// TODO: linkType
	ip4Layer := p.Layer(layers.LayerTypeIPv4)
	if ip4Layer != nil {
		ip4, _ := ip4Layer.(*layers.IPv4)
		l := ip4.Length
		newIPv4, err := fd.ipv4Defrag.DefragIPv4(ip4)
		if err != nil {
			return err
		} else if newIPv4 != nil {
			// TODO: correct way to detect finished reassemble?
			if newIPv4.Length != l {
				// TODO: better way to reconstruct package?
				sb := gopacket.NewSerializeBuffer()
				b, _ := sb.PrependBytes(len(newIPv4.Payload))
				copy(b, newIPv4.Payload)
				if err := newIPv4.SerializeTo(sb, gopacket.SerializeOptions{
					FixLengths:       true,
					ComputeChecksums: true,
				}); err != nil {
					return err
				}

				fd.IPV4Reassembled = append(fd.IPV4Reassembled, IPV4Reassembled{
					SourceIP:      ip4.SrcIP,
					DestinationIP: ip4.DstIP,
					Datagram:      sb.Bytes(),
				})

				pb, ok := p.(gopacket.PacketBuilder)
				if !ok {
					panic("not a PacketBuilder")
				}
				nextDecoder := newIPv4.NextLayerType()
				if err := nextDecoder.Decode(newIPv4.Payload, pb); err != nil {
					return err
				}
			}
		}
	}

	tcp := p.Layer(layers.LayerTypeTCP)
	if tcp != nil {
		tcp, _ := tcp.(*layers.TCP)
		fd.tcpAssembler.Assemble(p.NetworkLayer().NetworkFlow(), tcp)
	}

	return nil
}

func (fd *Decoder) Flush() {
	fd.tcpAssembler.FlushAll()
}
