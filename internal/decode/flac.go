package decode

import "math/bits"

type FLAC struct {
	Common
}

const (
	metadataBlockTypeStreaminfo    = 0
	metadataBlockTypePadding       = 1
	metadataBlockTypeApplication   = 2
	metadataBlockTypeSeektable     = 3
	metadataBlockTypeVorbisComment = 4
	metadataBlockTypeCuesheet      = 5
	metadataBlockTypePicture       = 6
)

var metadataBlockNames = map[uint]string{
	metadataBlockTypeStreaminfo:    "Streaminfo",
	metadataBlockTypePadding:       "Padding",
	metadataBlockTypeApplication:   "Application",
	metadataBlockTypeSeektable:     "Seektable",
	metadataBlockTypeVorbisComment: "Vorbis comment",
	metadataBlockTypeCuesheet:      "Cuesheet",
	metadataBlockTypePicture:       "Picture",
}

const (
	blockingStrategyFixed    = 0
	blockingStrategyVariable = 1
)

var blockingStrategyNames = map[uint]string{
	blockingStrategyFixed:    "Fixed",
	blockingStrategyVariable: "Variable",
}

// TODO: generic enough?
func (f *FLAC) UTF8Uint() uint64 {
	n := f.U8()
	c := bits.LeadingZeros8(uint8(n))
	switch c {
	case 0:
		// nop
	case 1:
		// TODO: error
	default:
		n = n & ((1 << (8 - c - 1)) - 1)
		for i := 1; i < c; i++ {
			n = n<<6 | f.U8()&0x3f
		}
	}

	return n
}

func (f *FLAC) Unmarshl() {
	f.FieldUTF8(4, "Magic")

	// is used in frame
	var streamInfoSamepleRate uint64
	var streamInfoBitPerSample uint64

	lastBlock := false
	for !lastBlock {
		f.Field("metadatablock", func() (Value, string) {
			lastBlock = f.FieldU1("last_block") == 1
			typ, _ := f.Field("type", func() (Value, string) {
				t := f.U7()
				name := "Unknown"
				if s, ok := metadataBlockNames[uint(t)]; ok {
					name = s
				}
				return Value{Type: TypeUint, Uint: t}, name
			})
			length := f.FieldU24("length")

			switch typ.Uint {
			case metadataBlockTypeStreaminfo:
				f.FieldU16("minimum_block_size")
				f.FieldU16("maximum_block_size")
				f.FieldU24("minimum_frame_size")
				f.FieldU24("maximum_frame_size")
				streamInfoSamepleRate = f.FieldU(20, "sample_rate")
				f.Field("channels", func() (Value, string) {
					return Value{Type: TypeUint, Uint: f.U3() + 1}, ""
				})
				// TODO: uint64 with fn?
				streamInfoBitPerSampleValue, _ := f.Field("bits_per_sample", func() (Value, string) {
					return Value{Type: TypeUint, Uint: f.U5() + 1}, ""
				})
				streamInfoBitPerSample = streamInfoBitPerSampleValue.Uint
				f.FieldU(36, "total_samples_in_steam")
				f.FieldBytes(16, "MD5")
			default:
				f.FieldBytes(uint(length), "Data")
			}

			return Value{}, typ.Str
		})
	}

	for !f.EOF() {
		f.Field("frame", func() (Value, string) {
			// # <14> 11111111111110
			f.Field("sync", func() (Value, string) {
				n := f.U14()
				s := "correct"
				if n != 0b11111111111110 {
					s = "incorrect"
				}
				return Value{Type: TypeUint, Uint: n}, s
			})

			// # <1> Reserved
			// # 0 : mandatory value
			// # 1 : reserved for future use
			// TODO: name?
			f.Field("reserved0", func() (Value, string) {
				n := f.U1()
				s := "correct"
				if n != 0 {
					s = "incorrect"
				}
				return Value{Type: TypeUint, Uint: n}, s
			})

			// # <1> Blocking strategy:
			// # 0 : fixed-blocksize stream; frame header encodes the frame number
			// # 1 : variable-blocksize stream; frame header encodes the sample number
			blockingStrategy, _ := f.Field("blocking_strategy", func() (Value, string) {
				n := f.U1()
				return Value{Type: TypeUint, Uint: n}, blockingStrategyNames[uint(n)]
			})

			// # <4> Block size in inter-channel samples:
			// # 0000 : reserved
			// # 0001 : 192 samples
			// # 0010-0101 : 576 * (2^(n-2)) samples, i.e. 576/1152/2304/4608
			// # 0110 : get 8 bit (blocksize-1) from end of header
			// # 0111 : get 16 bit (blocksize-1) from end of header
			// # 1000-1111 : 256 * (2^(n-8)) samples, i.e. 256/512/1024/2048/4096/8192/16384/32768
			var blockSizeBits uint64
			blockSize, _ := f.Field("block_size", func() (Value, string) {
				blockSizeBits = f.U4()
				switch blockSizeBits {
				case 0:
					return Value{Type: TypeUint, Uint: 0}, "reserved"
				case 1:
					return Value{Type: TypeUint, Uint: 192}, ""
				case 2:
				case 3:
				case 4:
				case 5:
					return Value{Type: TypeUint, Uint: 576 * (1 << (blockSizeBits - 2))}, ""
				case 6:
					return Value{Type: TypeUint, Uint: 0}, "end of header 8 but"
				case 7:
					return Value{Type: TypeUint, Uint: 0}, "end of header 16 bit"
				default:
					return Value{Type: TypeUint, Uint: 256 * (1 << (blockSizeBits - 8))}, ""
				}
				panic("unreachable")
			})

			// set sample_rate_pos [bitreader::bytepos $br]
			// # <4> Sample rate:
			// # 0000 : get from STREAMINFO metadata block
			// # 0001 : 88.2kHz
			// # 0010 : 176.4kHz
			// # 0011 : 192kHz
			// # 0100 : 8kHz
			// # 0101 : 16kHz
			// # 0110 : 22.05kHz
			// # 0111 : 24kHz
			// # 1000 : 32kHz
			// # 1001 : 44.1kHz
			// # 1010 : 48kHz
			// # 1011 : 96kHz
			// # 1100 : get 8 bit sample rate (in kHz) from end of header
			// # 1101 : get 16 bit sample rate (in Hz) from end of header
			// # 1110 : get 16 bit sample rate (in tens of Hz) from end of header
			// # 1111 : invalid, to prevent sync-fooling string of 1s
			var sampleRateBits uint64
			sampleRate, _ := f.Field("sample_rate", func() (Value, string) {
				sampleRateBits = f.U4()
				switch sampleRateBits {
				case 0:
					return Value{Type: TypeUint, Uint: streamInfoSamepleRate}, "streaminfo"
				case 1:
					return Value{Type: TypeUint, Uint: 88200}, ""
				case 2:
					return Value{Type: TypeUint, Uint: 176000}, ""
				case 3:
					return Value{Type: TypeUint, Uint: 19200}, ""
				case 4:
					return Value{Type: TypeUint, Uint: 800}, ""
				case 5:
					return Value{Type: TypeUint, Uint: 1600}, ""
				case 6:
					return Value{Type: TypeUint, Uint: 22050}, ""
				case 7:
					return Value{Type: TypeUint, Uint: 44100}, ""
				case 8:
					return Value{Type: TypeUint, Uint: 32000}, ""
				case 9:
					return Value{Type: TypeUint, Uint: 44100}, ""
				case 10:
					return Value{Type: TypeUint, Uint: 48000}, ""
				case 11:
					return Value{Type: TypeUint, Uint: 96000}, ""
				case 12:
					return Value{}, "end of header (8 bit*1000)"
				case 13:
					return Value{}, "end of header (16 bit)"
				case 14:
					return Value{}, "end of header (16 bit*10)"
				default:
					return Value{}, "invalid"
				}
				panic("unreachable")
			})

			// # <4> Channel assignment
			// # 0000-0111 : (number of independent channels)-1. Where defined, the channel order follows SMPTE/ITU-R recommendations. The assignments are as follows:
			// # 1 channel: mono
			// # 2 channels: left, right
			// # 3 channels: left, right, center
			// # 4 channels: front left, front right, back left, back right
			// # 5 channels: front left, front right, front center, back/surround left, back/surround right
			// # 6 channels: front left, front right, front center, LFE, back/surround left, back/surround right
			// # 7 channels: front left, front right, front center, LFE, back center, side left, side right
			// # 8 channels: front left, front right, front center, LFE, back left, back right, side left, side right
			// # 1000 : left/side stereo: channel 0 is the left channel, channel 1 is the side(difference) channel
			// # 1001 : right/side stereo: channel 0 is the side(difference) channel, channel 1 is the right channel
			// # 1010 : mid/side stereo: channel 0 is the mid(average) channel, channel 1 is the side(difference) channel
			// # 1011-1111 : reserved
			var sideChannelIndex uint
			channels, _ := f.Field("channel_assignment", func() (Value, string) {
				switch f.U4() {
				case 0:
					return Value{Type: TypeUint, Uint: 1}, "mono"
				case 1:
					return Value{Type: TypeUint, Uint: 2}, "left, right"
				case 2:
					return Value{Type: TypeUint, Uint: 3}, "left, right, center"
				case 3:
					return Value{Type: TypeUint, Uint: 4}, "front left, front right, back left, back right"
				case 4:
					return Value{Type: TypeUint, Uint: 5}, "front left, front right, front center, back/surround left, back/surround right"
				case 5:
					return Value{Type: TypeUint, Uint: 6}, "front left, front right, front center, LFE, back/surround left, back/surround right"
				case 6:
					return Value{Type: TypeUint, Uint: 7}, "front left, front right, front center, LFE, back center, side left, side right"
				case 7:
					return Value{Type: TypeUint, Uint: 8}, "front left, front right, front center, LFE, back left, back right, side left, side right"
				case 8:
					sideChannelIndex = 1
					return Value{Type: TypeUint, Uint: 2}, "left/side"
				case 9:
					sideChannelIndex = 0
					return Value{Type: TypeUint, Uint: 2}, "side/right"
				case 10:
					sideChannelIndex = 1
					return Value{Type: TypeUint, Uint: 2}, "mid/side"
				default:
					return Value{}, "reserved"
				}
				panic("unreachable")
			})

			// # <3> Sample size in bits:
			// # 000 : get from STREAMINFO metadata block
			// # 001 : 8 bits per sample
			// # 010 : 12 bits per sample
			// # 011 : reserved
			// # 100 : 16 bits per sample
			// # 101 : 20 bits per sample
			// # 110 : 24 bits per sample
			// # 111 : reserved
			sampleSize, _ := f.Field("sample_size", func() (Value, string) {
				switch f.U4() {
				case 0:
					return Value{Type: TypeUint, Uint: streamInfoBitPerSample}, "streaminfo"
				case 1:
					return Value{Type: TypeUint, Uint: 8}, ""
				case 2:
					return Value{Type: TypeUint, Uint: 12}, ""
				case 3:
					return Value{}, "reserved"
				case 4:
					return Value{Type: TypeUint, Uint: 16}, ""
				case 5:
					return Value{Type: TypeUint, Uint: 20}, ""
				case 6:
					return Value{Type: TypeUint, Uint: 24}, ""
				case 7:
					return Value{}, "reserved"
				}
				panic("unreachable")
			})

			// # <1> Reserved:
			// # 0 : mandatory value
			// # 1 : reserved for future use
			f.Field("reserved1", func() (Value, string) {
				n := f.U1()
				s := "correct"
				if n != 0 {
					s = "incorrect"
				}
				return Value{Type: TypeUint, Uint: n}, s
			})

			f.Field("end_of_header", func() (Value, string) {
				// 	# if(variable blocksize)
				// 	#    <8-56>:"UTF-8" coded sample number (decoded number is 36 bits) [4]
				// 	# else
				// 	#    <8-48>:"UTF-8" coded frame number (decoded number is 31 bits) [4]
				switch blockingStrategy.Uint {
				case blockingStrategyVariable:
					f.Field("sample_number", func() (Value, string) {
						return Value{Type: TypeUint, Uint: f.UTF8Uint()}, ""
					})
				case blockingStrategyFixed:
					f.Field("frame_number", func() (Value, string) {
						return Value{Type: TypeUint, Uint: f.UTF8Uint()}, ""
					})

				}

				// 	# if(blocksize bits == 011x)
				// 	#    8/16 bit (blocksize-1)
				switch blockSizeBits {
				case 6:
					f.Field("block_size", func() (Value, string) {
						return Value{Type: TypeUint, Uint: f.U8() + 1}, ""
					})
				case 7:
					f.Field("block_size", func() (Value, string) {
						return Value{Type: TypeUint, Uint: f.U16() + 1}, ""
					})
				}

				// 	# if(sample rate bits == 11xx)
				// 	#    8/16 bit sample rate
				switch sampleRateBits {
				case 12:
					sampleRate, _ = f.Field("sample_rate", func() (Value, string) {
						return Value{Type: TypeUint, Uint: f.U8() * 1000}, ""
					})
				case 13:
					sampleRate, _ = f.Field("sample_rate", func() (Value, string) {
						return Value{Type: TypeUint, Uint: f.U16()}, ""
					})
				case 14:
					sampleRate, _ = f.Field("sample_rate", func() (Value, string) {
						return Value{Type: TypeUint, Uint: f.U16() * 10}, ""
					})
				}

				return Value{}, ""
			})

			f.FieldU8("crc")

			_ = blockSize
			_ = sampleSize
			_ = channels

			// # CRC-8 (polynomial = x^8 + x^2 + x^1 + x^0, initialized with 0) of everything before the crc, including the sync code
			// log::entry $l "CRC" {
			// 	set frame_end [expr [bitreader::bytepos $br]-1]
			// 	set ccrc [crc8 [bitreader::byterange $br $frame_start $frame_end]]
			// 	set crc [bitreader::uint $br 8]
			// 	set crc_correct "(correct)"
			// 	if {$crc != $ccrc} {
			// 		set crc_correct [format "(incorrect %.2x)" $ccrc]
			// 	}
			// 	list [format "%.2x %s" $crc $crc_correct]
			// }

			return Value{}, ""
		})

		break
	}
}
