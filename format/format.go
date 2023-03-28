package format

import "github.com/wader/fq/pkg/decode"

// TODO: do before-format somehow and topology sort?
const (
	ProbeOrderBinUnique = 0   // binary with unlikely overlap
	ProbeOrderBinFuzzy  = 100 // binary with possible overlap
	ProbeOrderTextJSON  = 200 // text json has prio as yaml overlap
	ProbeOrderTextFuzzy = 300 // text with possible overlap
)

// TODO: move to group package somehow?

var (
	All = &decode.Group{Name: "all"}

	Image        = &decode.Group{Name: "image"}
	Probe        = &decode.Group{Name: "probe"}
	LinkFrame    = &decode.Group{Name: "link_frame", DefaultInArg: LinkFrameIn{}}   // ex: ethernet
	InetPacket   = &decode.Group{Name: "inet_packet", DefaultInArg: InetPacketIn{}} // ex: ipv4
	IpPacket     = &decode.Group{Name: "ip_packet", DefaultInArg: InetPacketIn{}}   // ex: tcp
	TcpStream    = &decode.Group{Name: "tcp_stream", DefaultInArg: TCPStreamIn{}}   // ex: http
	UdpPayload   = &decode.Group{Name: "udp_payload", DefaultInArg: UDPPayloadIn{}} // ex: dns
	Mp3FrameTags = &decode.Group{Name: "mp3_frame_tags"}

	Bytes = &decode.Group{Name: "bytes"}
	Bits  = &decode.Group{Name: "bits"}

	AacFrame           = &decode.Group{Name: "aac_frame"}
	Adts               = &decode.Group{Name: "adts"}
	AdtsFrame          = &decode.Group{Name: "adts_frame"}
	Aiff               = &decode.Group{Name: "aiff"}
	Amf0               = &decode.Group{Name: "amf0"}
	Apev2              = &decode.Group{Name: "apev2"}
	AppleBookmark      = &decode.Group{Name: "apple_bookmark"}
	Ar                 = &decode.Group{Name: "ar"}
	Asn1Ber            = &decode.Group{Name: "asn1_ber"}
	Av1Ccr             = &decode.Group{Name: "av1_ccr"}
	Av1Frame           = &decode.Group{Name: "av1_frame"}
	Av1Obu             = &decode.Group{Name: "av1_obu"}
	AvcAnnexb          = &decode.Group{Name: "avc_annexb"}
	AvcAu              = &decode.Group{Name: "avc_au"}
	AvcDcr             = &decode.Group{Name: "avc_dcr"}
	AvcNalu            = &decode.Group{Name: "avc_nalu"}
	AvcPps             = &decode.Group{Name: "avc_pps"}
	AvcSei             = &decode.Group{Name: "avc_sei"}
	AvcSps             = &decode.Group{Name: "avc_sps"}
	Avi                = &decode.Group{Name: "avi"}
	AvroOcf            = &decode.Group{Name: "avro_ocf"}
	Bencode            = &decode.Group{Name: "bencode"}
	BitcoinBlkdat      = &decode.Group{Name: "bitcoin_blkdat"}
	BitcoinBlock       = &decode.Group{Name: "bitcoin_block"}
	BitcoinScript      = &decode.Group{Name: "bitcoin_script"}
	BitcoinTransaction = &decode.Group{Name: "bitcoin_transaction"}
	Bplist             = &decode.Group{Name: "bplist"}
	BsdLoopbackFrame   = &decode.Group{Name: "bsd_loopback_frame"}
	Bson               = &decode.Group{Name: "bson"}
	Bzip2              = &decode.Group{Name: "bzip2"}
	Cbor               = &decode.Group{Name: "cbor"}
	Csv                = &decode.Group{Name: "csv"}
	Dns                = &decode.Group{Name: "dns"}
	DnsTcp             = &decode.Group{Name: "dns_tcp"}
	Elf                = &decode.Group{Name: "elf"}
	Ether8023Frame     = &decode.Group{Name: "ether8023_frame"}
	Exif               = &decode.Group{Name: "exif"}
	FairplaySpc        = &decode.Group{Name: "fairplay_spc"}
	Flac               = &decode.Group{Name: "flac"}
	FlacFrame          = &decode.Group{Name: "flac_frame"}
	FlacMetadatablock  = &decode.Group{Name: "flac_metadatablock"}
	FlacMetadatablocks = &decode.Group{Name: "flac_metadatablocks"}
	FlacPicture        = &decode.Group{Name: "flac_picture"}
	FlacStreaminfo     = &decode.Group{Name: "flac_streaminfo"}
	Flv                = &decode.Group{Name: "flv"}
	Gif                = &decode.Group{Name: "gif"}
	Gzip               = &decode.Group{Name: "gzip"}
	HevcAnnexb         = &decode.Group{Name: "hevc_annexb"}
	HevcAu             = &decode.Group{Name: "hevc_au"}
	HevcDcr            = &decode.Group{Name: "hevc_dcr"}
	HevcNalu           = &decode.Group{Name: "hevc_nalu"}
	HevcPps            = &decode.Group{Name: "hevc_pps"}
	HevcSps            = &decode.Group{Name: "hevc_sps"}
	HevcVps            = &decode.Group{Name: "hevc_vps"}
	Html               = &decode.Group{Name: "html"}
	IccProfile         = &decode.Group{Name: "icc_profile"}
	Icmp               = &decode.Group{Name: "icmp"}
	Icmpv6             = &decode.Group{Name: "icmpv6"}
	Id3v1              = &decode.Group{Name: "id3v1"}
	Id3v11             = &decode.Group{Name: "id3v11"}
	Id3v2              = &decode.Group{Name: "id3v2"}
	Ipv4Packet         = &decode.Group{Name: "ipv4_packet"}
	Ipv6Packet         = &decode.Group{Name: "ipv6_packet"}
	Jpeg               = &decode.Group{Name: "jpeg"}
	Json               = &decode.Group{Name: "json"}
	Jsonl              = &decode.Group{Name: "jsonl"}
	Macho              = &decode.Group{Name: "macho"}
	MachoFat           = &decode.Group{Name: "macho_fat"}
	Markdown           = &decode.Group{Name: "markdown"}
	Matroska           = &decode.Group{Name: "matroska"}
	Mp3                = &decode.Group{Name: "mp3"}
	Mp3Frame           = &decode.Group{Name: "mp3_frame"}
	Mp3FrameVbri       = &decode.Group{Name: "mp3_frame_vbri"}
	Mp3FrameXing       = &decode.Group{Name: "mp3_frame_xing"}
	Mp4                = &decode.Group{Name: "mp4"}
	MpegAsc            = &decode.Group{Name: "mpeg_asc"}
	MpegEs             = &decode.Group{Name: "mpeg_es"}
	MpegPes            = &decode.Group{Name: "mpeg_pes"}
	MpegPesPacket      = &decode.Group{Name: "mpeg_pes_packet"}
	MpegSpu            = &decode.Group{Name: "mpeg_spu"}
	MpegTs             = &decode.Group{Name: "mpeg_ts"}
	Msgpack            = &decode.Group{Name: "msgpack"}
	Ogg                = &decode.Group{Name: "ogg"}
	OggPage            = &decode.Group{Name: "ogg_page"}
	OpusPacket         = &decode.Group{Name: "opus_packet"}
	Pcap               = &decode.Group{Name: "pcap"}
	Pcapng             = &decode.Group{Name: "pcapng"}
	Png                = &decode.Group{Name: "png"}
	ProresFrame        = &decode.Group{Name: "prores_frame"}
	Protobuf           = &decode.Group{Name: "protobuf"}
	ProtobufWidevine   = &decode.Group{Name: "protobuf_widevine"}
	PsshPlayready      = &decode.Group{Name: "pssh_playready"}
	Rtmp               = &decode.Group{Name: "rtmp"}
	SllPacket          = &decode.Group{Name: "sll_packet"}
	Sll2Packet         = &decode.Group{Name: "sll2_packet"}
	Tar                = &decode.Group{Name: "tar"}
	TcpSegment         = &decode.Group{Name: "tcp_segment"}
	Tiff               = &decode.Group{Name: "tiff"}
	Tls                = &decode.Group{Name: "tls"}
	Toml               = &decode.Group{Name: "toml"}
	Tzif               = &decode.Group{Name: "tzif"}
	UdpDatagram        = &decode.Group{Name: "udp_datagram"}
	VorbisComment      = &decode.Group{Name: "vorbis_comment"}
	VorbisPacket       = &decode.Group{Name: "vorbis_packet"}
	Vp8Frame           = &decode.Group{Name: "vp8_frame"}
	Vp9Cfm             = &decode.Group{Name: "vp9_cfm"}
	Vp9Frame           = &decode.Group{Name: "vp9_frame"}
	VpxCcr             = &decode.Group{Name: "vpx_ccr"}
	Wasm               = &decode.Group{Name: "wasm"}
	Wav                = &decode.Group{Name: "wav"}
	Webp               = &decode.Group{Name: "webp"}
	Xml                = &decode.Group{Name: "xml"}
	Yaml               = &decode.Group{Name: "yaml"}
	Zip                = &decode.Group{Name: "zip"}
)

// below are data types used to communicate between formats <FormatName>In/Out

type AACFrameIn struct {
	ObjectType int `doc:"Audio object type"`
}
type AvcAuIn struct {
	LengthSize uint64 `doc:"Length value size"`
}

type AvcDcrOut struct {
	LengthSize uint64
}
type FlacFrameIn struct {
	SamplesBuf    []byte
	BitsPerSample int `doc:"Bits per sample"`
}

type FlacFrameOut struct {
	SamplesBuf    []byte
	Samples       uint64
	Channels      int
	BitsPerSample int
}
type FlacStreamInfo struct {
	SampleRate           uint64
	BitsPerSample        uint64
	TotalSamplesInStream uint64
	MD5                  []byte
}

type FlacStreaminfoOut struct {
	StreamInfo FlacStreamInfo
}

type FlacMetadatablockOut struct {
	IsLastBlock   bool
	HasStreamInfo bool
	StreamInfo    FlacStreamInfo
}

type FlacMetadatablocksOut struct {
	HasStreamInfo bool
	StreamInfo    FlacStreamInfo
}

type HevcAuIn struct {
	LengthSize uint64 `doc:"Length value size"`
}

type HevcDcrOut struct {
	LengthSize uint64
}

type OggPageOut struct {
	IsLastPage         bool
	IsFirstPage        bool
	IsContinuedPacket  bool
	StreamSerialNumber uint32
	SequenceNo         uint32
	Segments           [][]byte
}

type ProtoBufIn struct {
	Message ProtoBufMessage
}

type MatroskaIn struct {
	DecodeSamples bool `doc:"Decode samples"`
}

type Mp3In struct {
	MaxUniqueHeaderConfigs int `doc:"Max number of unique frame header configs allowed"`
	MaxUnknown             int `doc:"Max percent (0-100) unknown bits"`
	MaxSyncSeek            int `doc:"Max byte distance to next sync"`
}

type MP3FrameOut struct {
	MPEGVersion      int
	ProtectionAbsent bool
	BitRate          int
	SampleRate       int
	ChannelsIndex    int
	ChannelModeIndex int
}

type MpegDecoderConfig struct {
	ObjectType    int
	ASCObjectType int
}

type MpegEsOut struct {
	DecoderConfigs []MpegDecoderConfig
}

type MPEGASCOut struct {
	ObjectType int
}

type LinkFrameIn struct {
	Type           int
	IsLittleEndian bool // pcap endian etc
}

type InetPacketIn struct {
	EtherType int
}

type IPPacketIn struct {
	Protocol int
}

type UDPPayloadIn struct {
	SourcePort      int
	DestinationPort int
}

func (u UDPPayloadIn) IsPort(ports ...int) bool {
	for _, p := range ports {
		if u.DestinationPort == p || u.SourcePort == p {
			return true
		}
	}
	return false
}

func (u UDPPayloadIn) MustIsPort(fn func(format string, a ...any), ports ...int) {
	if !u.IsPort(ports...) {
		fn("incorrect udp port %t src:%d dst:%d", u.DestinationPort, u.SourcePort)
	}
}

type TCPStreamIn struct {
	IsClient        bool
	HasStart        bool
	HasEnd          bool
	SkippedBytes    uint64
	SourcePort      int
	DestinationPort int
}

type TCPStreamOut struct {
	PostFn func(peerIn any)
	InArg  any
}

func (t TCPStreamIn) IsPort(ports ...int) bool {
	for _, p := range ports {
		if (t.IsClient && t.DestinationPort == p) ||
			(!t.IsClient && t.SourcePort == p) {
			return true
		}
	}
	return false
}

func (t TCPStreamIn) MustIsPort(fn func(format string, a ...any), ports ...int) {
	if !t.IsPort(ports...) {
		fn("incorrect tcp port client %t src:%d dst:%d", t.IsClient, t.DestinationPort, t.SourcePort)
	}
}

type Mp4In struct {
	DecodeSamples  bool `doc:"Decode samples"`
	AllowTruncated bool `doc:"Allow box to be truncated"`
}

type AviIn struct {
	DecodeSamples bool `doc:"Decode samples"`
}

type ZipIn struct {
	Uncompress bool `doc:"Uncompress and probe files"`
}

type XMLIn struct {
	Seq             bool   `doc:"Use seq attribute to preserve element order"`
	Array           bool   `doc:"Decode as nested arrays"`
	AttributePrefix string `doc:"Prefix for attribute keys"`
}

type HTMLIn struct {
	Seq             bool   `doc:"Use seq attribute to preserve element order"`
	Array           bool   `doc:"Decode as nested arrays"`
	AttributePrefix string `doc:"Prefix for attribute keys"`
}

type CSVLIn struct {
	Comma   string `doc:"Separator character"`
	Comment string `doc:"Comment line character"`
}

type BitCoinBlockIn struct {
	HasHeader bool `doc:"Has blkdat header"`
}

type TLSIn struct {
	Keylog string `doc:"NSS Key Log content"`
}
