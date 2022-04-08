package format

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
	BSD_LOOPBACK_FRAME  = "bsd_loopback_frame"
	BSON                = "bson"
	BZIP2               = "bzip2"
	CBOR                = "cbor"
	DNS                 = "dns"
	DNS_TCP             = "dns_tcp"
	ELF                 = "elf"
	ETHER8023_FRAME     = "ether8023_frame"
	EXIF                = "exif"
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
	MACHO               = "macho"
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
	UDP_DATAGRAM        = "udp_datagram"
	VORBIS_COMMENT      = "vorbis_comment"
	VORBIS_PACKET       = "vorbis_packet"
	VP8_FRAME           = "vp8_frame"
	VP9_CFM             = "vp9_cfm"
	VP9_FRAME           = "vp9_frame"
	VPX_CCR             = "vpx_ccr"
	WAV                 = "wav"
	WEBP                = "webp"
	XING                = "xing"
	ZIP                 = "zip"
)

// below are data types used to communicate between formats <FormatName>In/Out

type FlacStreamInfo struct {
	SampleRate           uint64
	BitPerSample         uint64
	TotalSamplesInStream uint64
	MD5                  []byte
}

type FlacStreaminfoOut struct {
	StreamInfo FlacStreamInfo
}

type FlacMetadatablockStreamInfo struct {
	SampleRate           uint64
	BitPerSample         uint64
	TotalSamplesInStream uint64
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
	SamplesBuf []byte
	StreamInfo FlacStreamInfo
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

type AvcIn struct {
	LengthSize uint64
}

type AvcDcrOut struct {
	LengthSize uint64
}

type HevcIn struct {
	LengthSize uint64
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
	ObjectType int
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

func (u UDPPayloadIn) MustIsPort(fn func(format string, a ...interface{}), ports ...int) {
	if !u.IsPort(ports...) {
		fn("incorrect udp port %t src:%d dst:%d", u.DestinationPort, u.SourcePort)
	}
}

type TCPStreamIn struct {
	IsClient        bool
	HasStart        bool
	HasEnd          bool
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

func (t TCPStreamIn) MustIsPort(fn func(format string, a ...interface{}), ports ...int) {
	if !t.IsPort(ports...) {
		fn("incorrect tcp port client %t src:%d dst:%d", t.IsClient, t.DestinationPort, t.SourcePort)
	}
}

type X86_64In struct {
	Base      int64
	SymLookup func(uint64) (string, uint64)
}

type ARM64In struct {
	Base      int64
	SymLookup func(uint64) (string, uint64)
}
