package format

// TODO: do before-format somehow and topology sort?
const (
	ProbeOrderBinUnique = 0   // binary with unlikely overlap
	ProbeOrderBinFuzzy  = 100 // binary with possible overlap
	ProbeOrderTextJSON  = 200 // text json has prio as yaml overlap
	ProbeOrderTextFuzzy = 300 // text with possible overlap
)

// TODO: change to CamelCase?
//
//nolint:revive
const (
	ALL = "all"

	IMAGE       = "image"
	PROBE       = "probe"
	LINK_FRAME  = "link_frame"  // ex: ethernet
	INET_PACKET = "inet_packet" // ex: ipv4
	IP_PACKET   = "ip_packet"   // ex: tcp
	TCP_STREAM  = "tcp_stream"  // ex: http
	UDP_PAYLOAD = "udp_payload" // ex: dns

	AAC_FRAME           = "aac_frame"
	ADTS                = "adts"
	ADTS_FRAME          = "adts_frame"
	AMF0                = "amf0"
	APEV2               = "apev2"
	AR                  = "ar"
	ASN1_BER            = "asn1_ber"
	AV1_CCR             = "av1_ccr"
	AV1_FRAME           = "av1_frame"
	AV1_OBU             = "av1_obu"
	AVC_ANNEXB          = "avc_annexb"
	AVC_AU              = "avc_au"
	AVC_DCR             = "avc_dcr"
	AVC_NALU            = "avc_nalu"
	AVC_PPS             = "avc_pps"
	AVC_SEI             = "avc_sei"
	AVC_SPS             = "avc_sps"
	AVRO_OCF            = "avro_ocf"
	BENCODE             = "bencode"
	BITCOIN_BLKDAT      = "bitcoin_blkdat"
	BITCOIN_BLOCK       = "bitcoin_block"
	BITCOIN_SCRIPT      = "bitcoin_script"
	BITCOIN_TRANSACTION = "bitcoin_transaction"
	BSD_LOOPBACK_FRAME  = "bsd_loopback_frame"
	BSON                = "bson"
	BZIP2               = "bzip2"
	CBOR                = "cbor"
	CSV                 = "csv"
	DNS                 = "dns"
	DNS_TCP             = "dns_tcp"
	ELF                 = "elf"
	ETHER8023_FRAME     = "ether8023_frame"
	EXIF                = "exif"
	FAIRPLAY_SPC        = "fairplay_spc"
	FLAC                = "flac"
	FLAC_FRAME          = "flac_frame"
	FLAC_METADATABLOCK  = "flac_metadatablock"
	FLAC_METADATABLOCKS = "flac_metadatablocks"
	FLAC_PICTURE        = "flac_picture"
	FLAC_STREAMINFO     = "flac_streaminfo"
	FLV                 = "flv" // TODO:
	GIF                 = "gif"
	GZIP                = "gzip"
	HEVC_ANNEXB         = "hevc_annexb"
	HEVC_AU             = "hevc_au"
	HEVC_DCR            = "hevc_dcr"
	HEVC_NALU           = "hevc_nalu"
	HEVC_PPS            = "hevc_pps"
	HEVC_SPS            = "hevc_sps"
	HEVC_VPS            = "hevc_vps"
	HTML                = "html"
	ICC_PROFILE         = "icc_profile"
	ICMP                = "icmp"
	ICMPV6              = "icmpv6"
	ID3V1               = "id3v1"
	ID3V11              = "id3v11"
	ID3V2               = "id3v2"
	IPV4_PACKET         = "ipv4_packet"
	IPV6_PACKET         = "ipv6_packet"
	JPEG                = "jpeg"
	JSON                = "json"
	JSONL               = "jsonl"
	MACHO               = "macho"
	MACHO_FAT           = "macho_fat"
	MATROSKA            = "matroska"
	MP3                 = "mp3"
	MP3_FRAME           = "mp3_frame"
	MP4                 = "mp4"
	MPEG_ASC            = "mpeg_asc"
	MPEG_ES             = "mpeg_es"
	MPEG_PES            = "mpeg_pes"
	MPEG_PES_PACKET     = "mpeg_pes_packet"
	MPEG_SPU            = "mpeg_spu"
	MPEG_TS             = "mpeg_ts"
	MSGPACK             = "msgpack"
	OGG                 = "ogg"
	OGG_PAGE            = "ogg_page"
	OPUS_PACKET         = "opus_packet"
	PCAP                = "pcap"
	PCAPNG              = "pcapng"
	PNG                 = "png"
	PROTOBUF            = "protobuf"
	PROTOBUF_WIDEVINE   = "protobuf_widevine"
	PSSH_PLAYREADY      = "pssh_playready"
	RAW                 = "raw"
	RTMP                = "rtmp"
	SLL_PACKET          = "sll_packet"
	SLL2_PACKET         = "sll2_packet"
	TAR                 = "tar"
	TCP_SEGMENT         = "tcp_segment"
	TIFF                = "tiff"
	TOML                = "toml"
	UDP_DATAGRAM        = "udp_datagram"
	VORBIS_COMMENT      = "vorbis_comment"
	VORBIS_PACKET       = "vorbis_packet"
	VP8_FRAME           = "vp8_frame"
	VP9_CFM             = "vp9_cfm"
	VP9_FRAME           = "vp9_frame"
	VPX_CCR             = "vpx_ccr"
	WASM                = "wasm"
	WAV                 = "wav"
	WEBP                = "webp"
	XING                = "xing"
	XML                 = "xml"
	YAML                = "yaml"
	ZIP                 = "zip"
)

// below are data types used to communicate between formats <FormatName>In/Out

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

type OggPageOut struct {
	IsLastPage         bool
	IsFirstPage        bool
	IsContinuedPacket  bool
	StreamSerialNumber uint32
	SequenceNo         uint32
	Segments           [][]byte
}

type AvcAuIn struct {
	LengthSize uint64 `doc:"Length value size"`
}

type AvcDcrOut struct {
	LengthSize uint64
}

type HevcAuIn struct {
	LengthSize uint64 `doc:"Length value size"`
}

type HevcDcrOut struct {
	LengthSize uint64
}

type ProtoBufIn struct {
	Message ProtoBufMessage
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

type AACFrameIn struct {
	ObjectType int `doc:"Audio object type"`
}

type Mp3In struct {
	MaxUniqueHeaderConfigs int `doc:"Max number of unique frame header configs allowed"`
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
	DecodeSamples  bool `doc:"Decode supported media samples"`
	AllowTruncated bool `doc:"Allow box to be truncated"`
}

type ZipIn struct {
	Uncompress bool `doc:"Uncompress and probe files"`
}

type XMLIn struct {
	Seq   bool `doc:"Use seq attribute to preserve element order"`
	Array bool `doc:"Decode as nested arrays"`
}

type HTMLIn struct {
	Seq   bool `doc:"Use seq attribute to preserve element order"`
	Array bool `doc:"Decode as nested arrays"`
}

type CSVLIn struct {
	Comma   string `doc:"Separator character"`
	Comment string `doc:"Comment line character"`
}
