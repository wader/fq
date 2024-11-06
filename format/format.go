package format

import (
	"path/filepath"
	"strings"

	"github.com/wader/fq/pkg/decode"
)

// TODO: do before-format somehow and topology sort?
const (
	ProbeOrderBinUnique = 0   // binary with unlikely overlap
	ProbeOrderBinFuzzy  = 100 // binary with possible overlap
	ProbeOrderTextJSON  = 200 // text json has prio as yaml overlap
	ProbeOrderTextFuzzy = 300 // text with possible overlap
)

// TODO: move to group package somehow?

type Probe_In struct {
	IsProbe  bool
	Filename string
}

// Use HasExt("cer", "CeR", ...)
func (pi Probe_In) HasExt(ss ...string) bool {
	ext := filepath.Ext(pi.Filename)
	if ext == "" {
		return false
	}
	ext = ext[1:]
	for _, s := range ss {
		if strings.EqualFold(s, ext) {
			return true
		}
	}
	return false
}

type Probe_Args_In struct {
	IsProbeArgs bool
	Filename    string
	DecodeGroup string
}

var (
	Image          = &decode.Group{Name: "image"}
	INET_Packet    = &decode.Group{Name: "inet_packet", DefaultInArg: INET_Packet_In{}} // ex: ipv4
	IP_Packet      = &decode.Group{Name: "ip_packet", DefaultInArg: INET_Packet_In{}}   // ex: tcp
	Link_Frame     = &decode.Group{Name: "link_frame", DefaultInArg: Link_Frame_In{}}   // ex: ethernet
	MP3_Frame_Tags = &decode.Group{Name: "mp3_frame_tags"}
	Probe          = &decode.Group{Name: "probe", DefaultInArg: Probe_In{}}
	Probe_Args     = &decode.Group{Name: "probe_args", DefaultInArg: Probe_Args_In{}}
	TCP_Stream     = &decode.Group{Name: "tcp_stream", DefaultInArg: TCP_Stream_In{}}   // ex: http
	UDP_Payload    = &decode.Group{Name: "udp_payload", DefaultInArg: UDP_Payload_In{}} // ex: dns

	Bits  = &decode.Group{Name: "bits"}
	Bytes = &decode.Group{Name: "bytes"}

	AAC_Frame           = &decode.Group{Name: "aac_frame"}
	ADTS                = &decode.Group{Name: "adts"}
	ADTS_Frame          = &decode.Group{Name: "adts_frame"}
	AIFF                = &decode.Group{Name: "aiff"}
	AMF0                = &decode.Group{Name: "amf0"}
	Apev2               = &decode.Group{Name: "apev2"}
	Apple_Bookmark      = &decode.Group{Name: "apple_bookmark"}
	AR                  = &decode.Group{Name: "ar"}
	ASN1_BER            = &decode.Group{Name: "asn1_ber"}
	AV1_CCR             = &decode.Group{Name: "av1_ccr"}
	AV1_Frame           = &decode.Group{Name: "av1_frame"}
	AV1_OBU             = &decode.Group{Name: "av1_obu"}
	AVC_Annexb          = &decode.Group{Name: "avc_annexb"}
	AVC_AU              = &decode.Group{Name: "avc_au"}
	AVC_DCR             = &decode.Group{Name: "avc_dcr"}
	AVC_NALU            = &decode.Group{Name: "avc_nalu"}
	AVC_PPS             = &decode.Group{Name: "avc_pps"}
	AVC_SEI             = &decode.Group{Name: "avc_sei"}
	AVC_SPS             = &decode.Group{Name: "avc_sps"}
	AVI                 = &decode.Group{Name: "avi"}
	Avro_Ocf            = &decode.Group{Name: "avro_ocf"}
	Bencode             = &decode.Group{Name: "bencode"}
	Bitcoin_Blkdat      = &decode.Group{Name: "bitcoin_blkdat"}
	Bitcoin_Block       = &decode.Group{Name: "bitcoin_block"}
	Bitcoin_Script      = &decode.Group{Name: "bitcoin_script"}
	Bitcoin_Transaction = &decode.Group{Name: "bitcoin_transaction"}
	Bplist              = &decode.Group{Name: "bplist"}
	BSD_Loopback_Frame  = &decode.Group{Name: "bsd_loopback_frame"}
	BSON                = &decode.Group{Name: "bson"}
	Bzip2               = &decode.Group{Name: "bzip2"}
	CAFF                = &decode.Group{Name: "caff"}
	CBOR                = &decode.Group{Name: "cbor"}
	CSV                 = &decode.Group{Name: "csv"}
	DNS                 = &decode.Group{Name: "dns"}
	DNS_TCP             = &decode.Group{Name: "dns_tcp"}
	ELF                 = &decode.Group{Name: "elf"}
	Ether_8023_Frame    = &decode.Group{Name: "ether8023_frame"}
	Exif                = &decode.Group{Name: "exif"}
	Fairplay_SPC        = &decode.Group{Name: "fairplay_spc"}
	FIT                 = &decode.Group{Name: "fit"}
	FLAC                = &decode.Group{Name: "flac"}
	FLAC_Frame          = &decode.Group{Name: "flac_frame"}
	FLAC_Metadatablock  = &decode.Group{Name: "flac_metadatablock"}
	FLAC_Metadatablocks = &decode.Group{Name: "flac_metadatablocks"}
	FLAC_Picture        = &decode.Group{Name: "flac_picture"}
	FLAC_Streaminfo     = &decode.Group{Name: "flac_streaminfo"}
	FLV                 = &decode.Group{Name: "flv"}
	GIF                 = &decode.Group{Name: "gif"}
	Gzip                = &decode.Group{Name: "gzip"}
	HEVC_Annexb         = &decode.Group{Name: "hevc_annexb"}
	HEVC_AU             = &decode.Group{Name: "hevc_au"}
	HEVC_DCR            = &decode.Group{Name: "hevc_dcr"}
	HEVC_NALU           = &decode.Group{Name: "hevc_nalu"}
	HEVC_PPS            = &decode.Group{Name: "hevc_pps"}
	HEVC_SPS            = &decode.Group{Name: "hevc_sps"}
	HEVC_VPS            = &decode.Group{Name: "hevc_vps"}
	HTML                = &decode.Group{Name: "html"}
	ICC_Profile         = &decode.Group{Name: "icc_profile"}
	ICMP                = &decode.Group{Name: "icmp"}
	ICMPv6              = &decode.Group{Name: "icmpv6"}
	ID3v1               = &decode.Group{Name: "id3v1"}
	ID3v11              = &decode.Group{Name: "id3v11"}
	ID3v2               = &decode.Group{Name: "id3v2"}
	IPv4Packet          = &decode.Group{Name: "ipv4_packet"}
	IPv6Packet          = &decode.Group{Name: "ipv6_packet"}
	JP2C                = &decode.Group{Name: "jp2c"}
	JPEG                = &decode.Group{Name: "jpeg"}
	JSON                = &decode.Group{Name: "json"}
	JSONL               = &decode.Group{Name: "jsonl"}
	LevelDB_Descriptor  = &decode.Group{Name: "leveldb_descriptor"}
	LevelDB_LDB         = &decode.Group{Name: "leveldb_table"}
	LevelDB_LOG         = &decode.Group{Name: "leveldb_log"}
	LuaJIT              = &decode.Group{Name: "luajit"}
	MachO               = &decode.Group{Name: "macho"}
	MachO_Fat           = &decode.Group{Name: "macho_fat"}
	Markdown            = &decode.Group{Name: "markdown"}
	Matroska            = &decode.Group{Name: "matroska"}
	MIDI                = &decode.Group{Name: "midi"}
	MOC3                = &decode.Group{Name: "moc3"}
	MP3                 = &decode.Group{Name: "mp3"}
	MP3_Frame           = &decode.Group{Name: "mp3_frame"}
	MP3_Frame_VBRI      = &decode.Group{Name: "mp3_frame_vbri"}
	MP3_Frame_XING      = &decode.Group{Name: "mp3_frame_xing"}
	MP4                 = &decode.Group{Name: "mp4"}
	MPEG_ASC            = &decode.Group{Name: "mpeg_asc"}
	MPEG_ES             = &decode.Group{Name: "mpeg_es"}
	MPEG_PES_Packet     = &decode.Group{Name: "mpeg_pes_packet"}
	MPEG_SPU            = &decode.Group{Name: "mpeg_spu"}
	MPEG_TS             = &decode.Group{Name: "mpeg_ts"}
	MPES_PES            = &decode.Group{Name: "mpeg_pes"}
	MsgPack             = &decode.Group{Name: "msgpack"}
	Negentropy          = &decode.Group{Name: "negentropy"}
	NES                 = &decode.Group{Name: "nes"}
	Ogg                 = &decode.Group{Name: "ogg"}
	Ogg_Page            = &decode.Group{Name: "ogg_page"}
	Opentimestamps      = &decode.Group{Name: "opentimestamps"}
	Opus_Packet         = &decode.Group{Name: "opus_packet"}
	PCAP                = &decode.Group{Name: "pcap"}
	PCAPNG              = &decode.Group{Name: "pcapng"}
	Pg_BTree            = &decode.Group{Name: "pg_btree"}
	Pg_Control          = &decode.Group{Name: "pg_control"}
	Pg_Heap             = &decode.Group{Name: "pg_heap"}
	PNG                 = &decode.Group{Name: "png"}
	Prores_Frame        = &decode.Group{Name: "prores_frame"}
	Protobuf            = &decode.Group{Name: "protobuf"}
	ProtobufWidevine    = &decode.Group{Name: "protobuf_widevine"}
	PSSH_Playready      = &decode.Group{Name: "pssh_playready"}
	RTMP                = &decode.Group{Name: "rtmp"}
	SLL_Packet          = &decode.Group{Name: "sll_packet"}
	SLL2_Packet         = &decode.Group{Name: "sll2_packet"}
	TAP                 = &decode.Group{Name: "tap"}
	TAR                 = &decode.Group{Name: "tar"}
	TCP_Segment         = &decode.Group{Name: "tcp_segment"}
	TIFF                = &decode.Group{Name: "tiff"}
	TLS                 = &decode.Group{Name: "tls"}
	TOML                = &decode.Group{Name: "toml"}
	Tzif                = &decode.Group{Name: "tzif"}
	TZX                 = &decode.Group{Name: "tzx"}
	UDP_Datagram        = &decode.Group{Name: "udp_datagram"}
	Vorbis_Comment      = &decode.Group{Name: "vorbis_comment"}
	Vorbis_Packet       = &decode.Group{Name: "vorbis_packet"}
	VP8_Frame           = &decode.Group{Name: "vp8_frame"}
	VP9_CFM             = &decode.Group{Name: "vp9_cfm"}
	VP9_Frame           = &decode.Group{Name: "vp9_frame"}
	VPX_CCR             = &decode.Group{Name: "vpx_ccr"}
	WASM                = &decode.Group{Name: "wasm"}
	WAV                 = &decode.Group{Name: "wav"}
	WebP                = &decode.Group{Name: "webp"}
	XML                 = &decode.Group{Name: "xml"}
	YAML                = &decode.Group{Name: "yaml"}
	Zip                 = &decode.Group{Name: "zip"}
)

// below are data types used to communicate between formats <FormatName>In/Out

type AAC_Frame_In struct {
	ObjectType int `doc:"Audio object type"`
}
type AVC_AU_In struct {
	LengthSize uint64 `doc:"Length value size"`
}

type AVC_DCR_Out struct {
	LengthSize uint64
}

type CAFF_In struct {
	Uncompress bool `doc:"Uncompress and probe files"`
}

type FLAC_Frame_In struct {
	SamplesBuf    []byte
	BitsPerSample int `doc:"Bits per sample"`
}

type FLAC_Frame_Out struct {
	SamplesBuf    []byte
	Samples       uint64
	Channels      int
	BitsPerSample int
}
type FLAC_Stream_Info struct {
	SampleRate           uint64
	BitsPerSample        uint64
	TotalSamplesInStream uint64
	MD5                  []byte
}

type FLAC_Streaminfo_Out struct {
	StreamInfo FLAC_Stream_Info
}

type FLAC_Metadatablock_Out struct {
	IsLastBlock   bool
	HasStreamInfo bool
	StreamInfo    FLAC_Stream_Info
}

type FLAC_Metadatablocks_Out struct {
	HasStreamInfo bool
	StreamInfo    FLAC_Stream_Info
}

type HEVC_AU_In struct {
	LengthSize uint64 `doc:"Length value size"`
}

type HEVC_DCR_Out struct {
	LengthSize uint64
}

type Ogg_Page_Out struct {
	IsLastPage         bool
	IsFirstPage        bool
	IsContinuedPacket  bool
	StreamSerialNumber uint32
	SequenceNo         uint32
	Segments           [][]byte
}

type Protobuf_In struct {
	Message ProtoBufMessage
}

type Matroska_In struct {
	DecodeSamples bool `doc:"Decode samples"`
}

type MP3_In struct {
	MaxUniqueHeaderConfigs int `doc:"Max number of unique frame header configs allowed"`
	MaxUnknown             int `doc:"Max percent (0-100) unknown bits"`
	MaxSyncSeek            int `doc:"Max byte distance to next sync"`
}

type MP3_Frame_Out struct {
	MPEGVersion      int
	ProtectionAbsent bool
	BitRate          int
	SampleRate       int
	ChannelsIndex    int
	ChannelModeIndex int
}

type MPEG_Decoder_Config struct {
	ObjectType    int
	ASCObjectType int
}

type MPEG_ES_Out struct {
	DecoderConfigs []MPEG_Decoder_Config
}

type MPEG_ASC_Out struct {
	ObjectType int
}

type Link_Frame_In struct {
	Type           int
	IsLittleEndian bool // pcap endian etc
}

type INET_Packet_In struct {
	EtherType int
}

type IP_Packet_In struct {
	Protocol int
}

type UDP_Payload_In struct {
	SourcePort      int
	DestinationPort int
}

func (u UDP_Payload_In) IsPort(ports ...int) bool {
	for _, p := range ports {
		if u.DestinationPort == p || u.SourcePort == p {
			return true
		}
	}
	return false
}

func (u UDP_Payload_In) MustIsPort(fn func(format string, a ...any), ports ...int) {
	if !u.IsPort(ports...) {
		fn("incorrect udp port %t src:%d dst:%d", u.DestinationPort, u.SourcePort)
	}
}

type TCP_Stream_In struct {
	IsClient        bool
	HasStart        bool
	HasEnd          bool
	SkippedBytes    uint64
	SourcePort      int
	DestinationPort int
}

type TCP_Stream_Out struct {
	PostFn func(peerIn any)
	InArg  any
}

func (t TCP_Stream_In) IsPort(ports ...int) bool {
	for _, p := range ports {
		if (t.IsClient && t.DestinationPort == p) ||
			(!t.IsClient && t.SourcePort == p) {
			return true
		}
	}
	return false
}

func (t TCP_Stream_In) MustIsPort(fn func(format string, a ...any), ports ...int) {
	if !t.IsPort(ports...) {
		fn("incorrect tcp port client %t src:%d dst:%d", t.IsClient, t.DestinationPort, t.SourcePort)
	}
}

type MP4_In struct {
	DecodeSamples  bool `doc:"Decode samples"`
	AllowTruncated bool `doc:"Allow box to be truncated"`
}

type AVI_In struct {
	DecodeSamples        bool `doc:"Decode samples"`
	DecodeExtendedChunks bool `doc:"Decode extended chunks"`
}

type Zip_In struct {
	Uncompress bool `doc:"Uncompress and probe files"`
}

type XML_In struct {
	Seq             bool   `doc:"Use seq attribute to preserve element order"`
	Array           bool   `doc:"Decode as nested arrays"`
	AttributePrefix string `doc:"Prefix for attribute keys"`
}

type HTML_In struct {
	Seq             bool   `doc:"Use seq attribute to preserve element order"`
	Array           bool   `doc:"Decode as nested arrays"`
	AttributePrefix string `doc:"Prefix for attribute keys"`
}

type CSV_In struct {
	Comma   string `doc:"Separator character"`
	Comment string `doc:"Comment line character"`
}

type Bitcoin_Block_In struct {
	HasHeader bool `doc:"Has blkdat header"`
}

type TLS_In struct {
	Keylog string `doc:"NSS Key Log content"`
}

type Pg_Control_In struct {
	Flavour string `doc:"PostgreSQL flavour: postgres14, pgproee14.., postgres10"`
}

type Pg_Heap_In struct {
	Flavour string `doc:"PostgreSQL flavour: postgres14, pgproee14.., postgres10"`
	Page    int    `doc:"First page number in file, default is 0"`
	Segment int    `doc:"Segment file number (16790.1 is 1), default is 0"`
}

type Pg_BTree_In struct {
	Page int `doc:"First page number in file, default is 0"`
}
