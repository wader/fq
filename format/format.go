package format

//nolint:revive
const (
	ALL = "all"

	PROBE       = "probe"
	IMAGE       = "image"
	TCP_STREAM  = "tcp_stream"
	UDP_PAYLOAD = "udp_payload"
	LINK_FRAME  = "link_frame"

	RAW  = "raw"
	JSON = "json"
	BSON = "bson"

	DNS                = "dns"
	DNS_TCP            = "dns_tcp"
	ETHER8023_FRAME    = "ether8023_frame"
	BSD_LOOPBACK_FRAME = "bsd_loopback_frame"
	SLL_PACKET         = "sll_packet"
	SLL2_PACKET        = "sll2_packet"
	IPV4_PACKET        = "ipv4_packet"
	UDP_DATAGRAM       = "udp_datagram"
	TCP_SEGMENT        = "tcp_segment"
	ICMP               = "icmp"

	ELF = "elf"
	AR  = "ar"

	AAC_FRAME           = "aac_frame"
	ADTS                = "adts"
	ADTS_FRAME          = "adts_frame"
	APEV2               = "apev2"
	AV1_CCR             = "av1_ccr"
	AV1_FRAME           = "av1_frame"
	AV1_OBU             = "av1_obu"
	BENCODE             = "bencode"
	BZIP2               = "bzip2"
	EXIF                = "exif"
	FLAC                = "flac"
	FLAC_FRAME          = "flac_frame"
	FLAC_METADATABLOCK  = "flac_metadatablock"
	FLAC_METADATABLOCKS = "flac_metadatablocks"
	FLAC_STREAMINFO     = "flac_streaminfo"
	FLAC_PICTURE        = "flac_picture"
	FLV                 = "flv" // TODO:
	GIF                 = "gif"
	GZIP                = "gzip"
	ICC_PROFILE         = "icc_profile"
	ID3V1               = "id3v1"
	ID3V11              = "id3v11"
	ID3V2               = "id3v2"
	JPEG                = "jpeg"
	MATROSKA            = "matroska"
	MP3                 = "mp3"
	MP3_FRAME           = "mp3_frame"
	XING                = "xing"
	MP4                 = "mp4"
	MPEG_ASC            = "mpeg_asc"
	AVC_ANNEXB          = "avc_annexb"
	AVC_DCR             = "avc_dcr"
	AVC_SPS             = "avc_sps"
	AVC_PPS             = "avc_pps"
	AVC_SEI             = "avc_sei"
	AVC_NALU            = "avc_nalu"
	AVC_AU              = "avc_au"
	HEVC_ANNEXB         = "hevc_annexb"
	HEVC_AU             = "hevc_au"
	HEVC_NALU           = "hevc_nalu"
	HEVC_DCR            = "hevc_dcr"
	MPEG_ES             = "mpeg_es"
	MPEG_PES            = "mpeg_pes"
	MPEG_PES_PACKET     = "mpeg_pes_packet"
	MPEG_SPU            = "mpeg_spu"
	MPEG_TS             = "mpeg_ts"
	OGG                 = "ogg"
	OGG_PAGE            = "ogg_page"
	OPUS_PACKET         = "opus_packet"
	PCAP                = "pcap"
	PCAPNG              = "pcapng"
	PNG                 = "png"
	PROTOBUF            = "protobuf"
	PROTOBUF_WIDEVINE   = "protobuf_widevine"
	PSSH_PLAYREADY      = "pssh_playready"
	TAR                 = "tar"
	TIFF                = "tiff"
	VORBIS_COMMENT      = "vorbis_comment"
	VORBIS_PACKET       = "vorbis_packet"
	VP8_FRAME           = "vp8_frame"
	VP9_FRAME           = "vp9_frame"
	VP9_CFM             = "vp9_cfm"
	VPX_CCR             = "vpx_ccr"
	WAV                 = "wav"
	WEBP                = "webp"
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

type In struct {
	SourcePort      int
	DestinationPort int
}

type LinkFrameIn struct {
	Type         int
	LittleEndian bool // pcap endian etc
}

type UDPPayloadIn struct {
	SourcePort      int
	DestinationPort int
}

type TCPStreamIn struct {
	SourcePort      int
	DestinationPort int
}

type X86_64In struct {
	Base      int64
	SymLookup func(uint64) (string, uint64)
}

type ARM64In struct {
	Base      int64
	SymLookup func(uint64) (string, uint64)
}
