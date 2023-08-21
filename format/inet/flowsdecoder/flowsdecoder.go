package flowsdecoder

// TODO: option to not allow missing syn/ack?

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"

	"github.com/gopacket/gopacket"
	"github.com/gopacket/gopacket/ip4defrag"
	"github.com/gopacket/gopacket/layers"
	"github.com/gopacket/gopacket/reassembly"
)

type TCPEndpoint struct {
	IP   net.IP
	Port int
}

type TCPDirection struct {
	Endpoint     TCPEndpoint
	HasStart     bool
	HasEnd       bool
	Buffer       *bytes.Buffer
	SkippedBytes uint64
}

type TCPConnection struct {
	Client     *TCPDirection
	Server     *TCPDirection
	tcpState   *reassembly.TCPSimpleFSM
	optChecker *reassembly.TCPOptionCheck
	net        gopacket.Flow
	transport  gopacket.Flow
}

func (t *TCPConnection) Accept(tcp *layers.TCP, ci gopacket.CaptureInfo, dir reassembly.TCPFlowDirection, nextSeq reassembly.Sequence, start *bool, ac reassembly.AssemblerContext) bool {
	// has ok state?
	if !t.tcpState.CheckState(tcp, dir) {
		// TODO: handle err?
		return false
	}
	if t.optChecker != nil {
		// has ok options?
		if err := t.optChecker.Accept(tcp, ci, dir, nextSeq, start); err != nil {
			// TODO: handle err?
			return false
		}
	}
	// TODO: checksum?

	// accept
	return true
}

func (t *TCPConnection) ReassembledSG(sg reassembly.ScatterGather, ac reassembly.AssemblerContext) {
	dir, start, end, skip := sg.Info()
	length, _ := sg.Lengths()

	var d *TCPDirection
	switch dir {
	case reassembly.TCPDirClientToServer:
		d = t.Client
	case reassembly.TCPDirServerToClient:
		d = t.Server
	default:
		panic("unreachable")
	}

	if skip == -1 {
		// can't find where skip == -1 is documented but this is what gopacket reassemblydump does
		// to allow missing syn/ack
	} else if skip != 0 {
		// stream has missing bytes
		d.SkippedBytes += uint64(skip)
		return
	}

	d.HasStart = d.HasStart || start
	d.HasEnd = d.HasEnd || end

	data := sg.Fetch(length)

	d.Buffer.Write(data)
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
	// TODO: understand how gopacket handles broken/too short packets, seems like
	// we can get here when lots of things are missing, assume zero port for now
	var clientPort int
	if len(transport.Src().Raw()) == 2 {
		clientPort = int(binary.BigEndian.Uint16(transport.Src().Raw()))
	}
	var serverPort int
	if len(transport.Dst().Raw()) == 2 {
		serverPort = int(binary.BigEndian.Uint16(transport.Dst().Raw()))
	}

	stream := &TCPConnection{
		Client: &TCPDirection{
			Endpoint: TCPEndpoint{
				IP:   append([]byte(nil), net.Src().Raw()...),
				Port: clientPort,
			},
			Buffer: &bytes.Buffer{},
		},
		Server: &TCPDirection{
			Endpoint: TCPEndpoint{
				IP:   append([]byte(nil), net.Dst().Raw()...),
				Port: serverPort,
			},
			Buffer: &bytes.Buffer{},
		},

		net:       net,
		transport: transport,
		tcpState:  reassembly.NewTCPSimpleFSM(fsmOptions),
	}

	if fd.Options.CheckTCPOptions {
		c := reassembly.NewTCPOptionCheck()
		stream.optChecker = &c
	}

	fd.TCPConnections = append(fd.TCPConnections, stream)

	return stream
}

type Decoder struct {
	Options DecoderOptions

	TCPConnections  []*TCPConnection
	IPV4Reassembled []IPV4Reassembled

	ipv4Defrag   *ip4defrag.IPv4Defragmenter
	tcpAssembler *reassembly.Assembler
}

type DecoderOptions struct {
	CheckTCPOptions bool
}

func New(options DecoderOptions) *Decoder {
	flowDecoder := &Decoder{
		Options: options,
	}
	streamPool := reassembly.NewStreamPool(flowDecoder)
	tcpAssembler := reassembly.NewAssembler(streamPool)
	flowDecoder.tcpAssembler = tcpAssembler
	flowDecoder.ipv4Defrag = ip4defrag.NewIPv4Defragmenter()

	return flowDecoder
}

func (fd *Decoder) EthernetFrame(bs []byte) error {
	return fd.packet(gopacket.NewPacket(bs, layers.LayerTypeEthernet, gopacket.Lazy))
}

func (fd *Decoder) IPv4Packet(bs []byte) error {
	return fd.packet(gopacket.NewPacket(bs, layers.LayerTypeIPv4, gopacket.Lazy))
}

func (fd *Decoder) IPv6Packet(bs []byte) error {
	return fd.packet(gopacket.NewPacket(bs, layers.LayerTypeIPv6, gopacket.Lazy))
}

func (fd *Decoder) SLLPacket(bs []byte) error {
	return fd.packet(gopacket.NewPacket(bs, layers.LayerTypeLinuxSLL, gopacket.Lazy))
}

func (fd *Decoder) SLL2Packet(bs []byte) error {
	return fd.packet(gopacket.NewPacket(bs, layers.LayerTypeLinuxSLL2, gopacket.Lazy))
}

func (fd *Decoder) LoopbackFrame(bs []byte) error {
	return fd.packet(gopacket.NewPacket(bs, layers.LayerTypeLoopback, gopacket.Lazy))
}

// LinkTypeRAW IPv4 or Ipv6
func (fd *Decoder) RAWIPFrame(bs []byte) error {
	version := bs[0] >> 4
	switch version {
	case 4:
		return fd.IPv4Packet(bs)
	case 6:
		return fd.IPv6Packet(bs)
	}
	return fmt.Errorf("invalid ip version %v", version)
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

				// i think this replaces p with the newly defragmented ip packet and is
				// used below when reassembling tcp streams
				// see gopacket reassemblydump example
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
