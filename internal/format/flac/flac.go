package flac

// TODO: more metablocks
// TODO: crc, md5
// TODO: 24 picture truncate

import (
	"fq/internal/decode"
	"math/bits"
)

var Register = &decode.Register{
	Name: "flac",
	MIME: "",
	New:  func() decode.Decoder { return &Decoder{} },
}

// Decoder is a FLAC decoder
type Decoder struct {
	decode.Common
}

const (
	MetadataBlockStreaminfo    = 0
	MetadataBlockPadding       = 1
	MetadataBlockApplication   = 2
	MetadataBlockSeektable     = 3
	MetadataBlockVorbisComment = 4
	MetadataBlockCuesheet      = 5
	MetadataBlockPicture       = 6
)

var metadataBlockNames = map[uint]string{
	MetadataBlockStreaminfo:    "Streaminfo",
	MetadataBlockPadding:       "Padding",
	MetadataBlockApplication:   "Application",
	MetadataBlockSeektable:     "Seektable",
	MetadataBlockVorbisComment: "Vorbis comment",
	MetadataBlockCuesheet:      "Cuesheet",
	MetadataBlockPicture:       "Picture",
}

const (
	BlockingStrategyFixed = iota
	BlockingStrategyVariable
)

var BlockingStrategyNames = map[uint]string{
	BlockingStrategyFixed:    "Fixed",
	BlockingStrategyVariable: "Variable",
}

const (
	SubframeConstant = iota
	SubframeVerbatim
	SubframeFixed
	SubframeLPC
)

var SubframeTypeNames = map[uint]string{
	SubframeConstant: "Constant",
	SubframeVerbatim: "Verbatim",
	SubframeFixed:    "Fixed",
	SubframeLPC:      "LPC",
}

// TODO: generic enough?
func (d *Decoder) UTF8Uint() uint64 {
	n := d.U8()
	// leading ones, bit negate and count zeroes
	c := bits.LeadingZeros8(^uint8(n))
	switch c {
	case 0:
		// nop
	case 1:
		// TODO: error
		panic("invalid UTF8Uint")
	default:
		n = n & ((1 << (8 - c - 1)) - 1)
		for i := 1; i < c; i++ {
			n = n<<6 | d.U8()&0x3f
		}
	}
	return n
}

// Decode FLAC
func (d *Decoder) Decode(opts decode.Options) bool {
	magic := d.FieldUTF8("magic", 4)
	if opts.Probe && magic == "fLaC" {
		return true
	}

	// is used in frame
	var streamInfoSamepleRate uint64
	var streamInfoBitPerSample uint64

	lastBlock := false
	for !lastBlock {
		d.FieldNoneFn("metadatablock", func() {
			lastBlock = d.FieldU1("last_block") == 1
			typ := d.FieldUFn("type", func() (uint64, decode.Format, string) {
				t := d.U7()
				name := "Unknown"
				if s, ok := metadataBlockNames[uint(t)]; ok {
					name = s
				}
				return t, decode.FormatDecimal, name
			})
			length := d.FieldU24("length")

			switch typ {
			case MetadataBlockStreaminfo:
				d.FieldU16("minimum_block_size")
				d.FieldU16("maximum_block_size")
				d.FieldU24("minimum_frame_size")
				d.FieldU24("maximum_frame_size")
				streamInfoSamepleRate = d.FieldU("sample_rate", 20)
				// <3> (number of channels)-1. FLAC supports from 1 to 8 channels
				d.FieldUFn("channels", func() (uint64, decode.Format, string) { return d.U3() + 1, decode.FormatDecimal, "" })
				// <5> (bits per sample)-1. FLAC supports from 4 to 32 bits per sample. Currently the reference encoder and decoders only support up to 24 bits per sample.
				streamInfoBitPerSample = d.FieldUFn("bits_per_sample", func() (uint64, decode.Format, string) {
					return d.U5() + 1, decode.FormatDecimal, ""
				})
				d.FieldU("total_samples_in_steam", 36)
				d.FieldBytes("md5", 16)
			default:
				d.FieldBytes("data", length)
			}
		})
	}

	for !d.End() {
		d.FieldNoneFn("frame", func() {
			// <14> 11111111111110
			d.FieldVerifyUFn("sync", 0b11111111111110, d.U14)

			// <1> Reserved
			// 0 : mandatory value
			// 1 : reserved for future use
			d.FieldVerifyUFn("reserved0", 0, d.U1)

			// <1> Blocking strategy:
			// 0 : fixed-blocksize stream; frame header encodes the frame number
			// 1 : variable-blocksize stream; frame header encodes the sample number
			blockingStrategy := d.FieldUFn("blocking_strategy", func() (uint64, decode.Format, string) {
				switch d.U1() {
				case 0:
					return BlockingStrategyFixed, decode.FormatDecimal, BlockingStrategyNames[BlockingStrategyFixed]
				default:
					return BlockingStrategyVariable, decode.FormatDecimal, BlockingStrategyNames[BlockingStrategyVariable]
				}
			})

			// <4> Block size in inter-channel samples:
			// 0000 : reserved
			// 0001 : 192 samples
			// 0010-0101 : 576 * (2^(n-2)) samples, i.e. 576/1152/2304/4608
			// 0110 : get 8 bit (blocksize-1) from end of header
			// 0111 : get 16 bit (blocksize-1) from end of header
			// 1000-1111 : 256 * (2^(n-8)) samples, i.e. 256/512/1024/2048/4096/8192/16384/32768
			var blockSizeBits uint64
			blockSize := d.FieldUFn("block_size", func() (uint64, decode.Format, string) {
				blockSizeBits = d.U4()
				switch blockSizeBits {
				case 0:
					return 0, decode.FormatDecimal, "reserved"
				case 1:
					return 192, decode.FormatDecimal, ""
				case 2, 3, 4, 5:
					return 576 * (1 << (blockSizeBits - 2)), decode.FormatDecimal, ""
				case 6:
					return 0, decode.FormatDecimal, "end of header (8 bit)"
				case 7:
					return 0, decode.FormatDecimal, "end of header (16 bit)"
				default:
					return 256 * (1 << (blockSizeBits - 8)), decode.FormatDecimal, ""
				}
			})

			// <4> Sample rate:
			// 0000 : get from STREAMINFO metadata block
			// 0001 : 88.2kHz
			// 0010 : 176.4kHz
			// 0011 : 192kHz
			// 0100 : 8kHz
			// 0101 : 16kHz
			// 0110 : 22.05kHz
			// 0111 : 24kHz
			// 1000 : 32kHz
			// 1001 : 44.1kHz
			// 1010 : 48kHz
			// 1011 : 96kHz
			// 1100 : get 8 bit sample rate (in kHz) from end of header
			// 1101 : get 16 bit sample rate (in Hz) from end of header
			// 1110 : get 16 bit sample rate (in tens of Hz) from end of header
			// 1111 : invalid, to prevent sync-fooling string of 1s
			var sampleRateBits uint64
			d.FieldUFn("sample_rate", func() (uint64, decode.Format, string) {
				sampleRateBits = d.U4()
				switch sampleRateBits {
				case 0:
					return streamInfoSamepleRate, decode.FormatDecimal, "streaminfo"
				case 1:
					return 88200, decode.FormatDecimal, ""
				case 2:
					return 176000, decode.FormatDecimal, ""
				case 3:
					return 19200, decode.FormatDecimal, ""
				case 4:
					return 800, decode.FormatDecimal, ""
				case 5:
					return 1600, decode.FormatDecimal, ""
				case 6:
					return 22050, decode.FormatDecimal, ""
				case 7:
					return 44100, decode.FormatDecimal, ""
				case 8:
					return 32000, decode.FormatDecimal, ""
				case 9:
					return 44100, decode.FormatDecimal, ""
				case 10:
					return 48000, decode.FormatDecimal, ""
				case 11:
					return 96000, decode.FormatDecimal, ""
				case 12:
					return 0, decode.FormatDecimal, "end of header (8 bit*1000)"
				case 13:
					return 0, decode.FormatDecimal, "end of header (16 bit)"
				case 14:
					return 0, decode.FormatDecimal, "end of header (16 bit*10)"
				default:
					return 0, decode.FormatDecimal, "invalid"
				}
			})

			// <4> Channel assignment
			// 0000-0111 : (number of independent channels)-1. Where defined, the channel order follows SMPTE/ITU-R recommendations. The assignments are as follows:
			// 1 channel: mono
			// 2 channels: left, right
			// 3 channels: left, right, center
			// 4 channels: front left, front right, back left, back right
			// 5 channels: front left, front right, front center, back/surround left, back/surround right
			// 6 channels: front left, front right, front center, LFE, back/surround left, back/surround right
			// 7 channels: front left, front right, front center, LFE, back center, side left, side right
			// 8 channels: front left, front right, front center, LFE, back left, back right, side left, side right
			// 1000 : left/side stereo: channel 0 is the left channel, channel 1 is the side(difference) channel
			// 1001 : right/side stereo: channel 0 is the side(difference) channel, channel 1 is the right channel
			// 1010 : mid/side stereo: channel 0 is the mid(average) channel, channel 1 is the side(difference) channel
			// 1011-1111 : reserved
			sideChannelIndex := -1
			channels := d.FieldUFn("channel_assignment", func() (uint64, decode.Format, string) {
				si, u, fmt, disp := func() (int, uint64, decode.Format, string) {
					switch d.U4() {
					case 0:
						return -1, 1, decode.FormatDecimal, "mono"
					case 1:
						return -1, 2, decode.FormatDecimal, "left, right"
					case 2:
						return -1, 3, decode.FormatDecimal, "left, right, center"
					case 3:
						return -1, 4, decode.FormatDecimal, "front left, front right, back left, back right"
					case 4:
						return -1, 5, decode.FormatDecimal, "front left, front right, front center, back/surround left, back/surround right"
					case 5:
						return -1, 6, decode.FormatDecimal, "front left, front right, front center, LFE, back/surround left, back/surround right"
					case 6:
						return -1, 7, decode.FormatDecimal, "front left, front right, front center, LFE, back center, side left, side right"
					case 7:
						return -1, 8, decode.FormatDecimal, "front left, front right, front center, LFE, back left, back right, side left, side right"
					case 8:
						sideChannelIndex = 1
						return -1, 2, decode.FormatDecimal, "left/side"
					case 9:
						sideChannelIndex = 0
						return -1, 2, decode.FormatDecimal, "side/right"
					case 10:
						sideChannelIndex = 1
						return -1, 2, decode.FormatDecimal, "mid/side"
					default:
						return -1, 0, decode.FormatDecimal, "reserved"
					}
				}()
				if si != -1 {
					sideChannelIndex = si
					d.FieldUFn("side_channel_index", func() (uint64, decode.Format, string) {
						return uint64(sideChannelIndex), decode.FormatDecimal, ""
					})
				}
				return u, fmt, disp
			})

			// <3> Sample size in bits:
			// 000 : get from STREAMINFO metadata block
			// 001 : 8 bits per sample
			// 010 : 12 bits per sample
			// 011 : reserved
			// 100 : 16 bits per sample
			// 101 : 20 bits per sample
			// 110 : 24 bits per sample
			// 111 : reserved
			sampleSize := d.FieldUFn("sample_size", func() (uint64, decode.Format, string) {
				switch d.U3() {
				case 0:
					return streamInfoBitPerSample, decode.FormatDecimal, "streaminfo"
				case 1:
					return 8, decode.FormatDecimal, ""
				case 2:
					return 12, decode.FormatDecimal, ""
				case 3:
					return 0, decode.FormatDecimal, "reserved"
				case 4:
					return 16, decode.FormatDecimal, ""
				case 5:
					return 20, decode.FormatDecimal, ""
				case 6:
					return 24, decode.FormatDecimal, ""
				case 7:
					return 0, decode.FormatDecimal, "reserved"
				}
				panic("unreachable")
			})

			// <1> Reserved:
			// 0 : mandatory value
			// 1 : reserved for future use
			d.FieldVerifyUFn("reserved1", 0, d.U1)

			d.FieldNoneFn("end_of_header", func() {
				// if(variable blocksize)
				//   <8-56>:"UTF-8" coded sample number (decoded number is 36 bits) [4]
				// else
				//   <8-48>:"UTF-8" coded frame number (decoded number is 31 bits) [4]
				switch blockingStrategy {
				case BlockingStrategyVariable:
					d.FieldUFn("sample_number", func() (uint64, decode.Format, string) {
						return d.UTF8Uint(), decode.FormatDecimal, ""
					})
				case BlockingStrategyFixed:
					d.FieldUFn("frame_number", func() (uint64, decode.Format, string) {
						return d.UTF8Uint(), decode.FormatDecimal, ""
					})
				}

				// if(blocksize bits == 011x)
				//   8/16 bit (blocksize-1)
				switch blockSizeBits {
				case 6:
					blockSize = d.FieldUFn("block_size", func() (uint64, decode.Format, string) {
						return d.U8() + 1, decode.FormatDecimal, ""
					})
				case 7:
					blockSize = d.FieldUFn("block_size", func() (uint64, decode.Format, string) {
						return d.U16() + 1, decode.FormatDecimal, ""
					})
				}

				// if(sample rate bits == 11xx)
				//   8/16 bit sample rate
				switch sampleRateBits {
				case 12:
					d.FieldUFn("sample_rate", func() (uint64, decode.Format, string) {
						return d.U8() * 1000, decode.FormatDecimal, ""
					})
				case 13:
					d.FieldUFn("sample_rate", func() (uint64, decode.Format, string) {
						return d.U16(), decode.FormatDecimal, ""
					})
				case 14:
					d.FieldUFn("sample_rate", func() (uint64, decode.Format, string) {
						return d.U16() * 10, decode.FormatDecimal, ""
					})
				}
			})

			// CRC-8 (polynomial = x^8 + x^2 + x^1 + x^0, initialized with 0) of everything before the crc, including the sync code
			d.FieldU8("crc")

			for channelIndex := 0; channelIndex < int(channels); channelIndex++ {
				d.FieldNoneFn("subframe", func() {
					// <1> Zero bit padding, to prevent sync-fooling string of 1s
					d.FieldVerifyUFn("zero_bit", 0, d.U1)

					// <6> Subframe type:
					// 000000 : SUBFRAME_CONSTANT
					// 000001 : SUBFRAME_VERBATIM
					// 00001x : reserved
					// 0001xx : reserved
					// 001xxx : if(xxx <= 4) SUBFRAME_FIXED, xxx=order ; else reserved
					// 01xxxx : reserved
					// 1xxxxx : SUBFRAME_LPC, xxxxx=order-1
					var lpcOrder uint64
					subframeType := d.FieldUFn("subframe_type", func() (uint64, decode.Format, string) {
						u, fmt, disp := func() (uint64, decode.Format, string) {
							bits := d.U6()
							switch bits {
							case 0:
								return SubframeConstant, decode.FormatDecimal, SubframeTypeNames[SubframeConstant]
							case 1:
								return SubframeVerbatim, decode.FormatDecimal, SubframeTypeNames[SubframeVerbatim]
							case 8, 9, 10, 11, 12:
								lpcOrder = bits & 0x7
								return SubframeFixed, decode.FormatDecimal, SubframeTypeNames[SubframeFixed]
							default:
								if bits&0x20 > 0 {
									lpcOrder = (bits & 0x1f) + 1
								} else {
									return 0, decode.FormatDecimal, "reserved"
								}
								return SubframeLPC, decode.FormatDecimal, SubframeTypeNames[SubframeLPC]
							}
						}()
						d.FieldUFn("lpc_order", func() (uint64, decode.Format, string) {
							return uint64(lpcOrder), decode.FormatDecimal, ""
						})
						return u, fmt, disp
					})

					// 'Wasted bits-per-sample' flag:
					// 0 : no wasted bits-per-sample in source subblock, k=0
					// 1 : k wasted bits-per-sample in source subblock, k-1 follows, unary coded; e.g. k=3 => 001 follows, k=7 => 0000001 follows.
					wastedBitsFlag := d.FieldU1("wasted_bits_flag")
					var wastedBitsK uint64
					if wastedBitsFlag != 0 {
						wastedBitsK = d.FieldUFn("wasted_bits_k", func() (uint64, decode.Format, string) {
							return uint64(d.Unary(0)) + 1, decode.FormatDecimal, ""
						})
					}

					subframeSampleSize := sampleSize - wastedBitsK
					// if channel is side, add en extra sample bit
					// https://github.com/xiph/flac/blob/37e675b777d4e0de53ac9ff69e2aea10d92e729c/src/libFLAC/stream_decoder.c#L2040
					if channelIndex == sideChannelIndex {
						subframeSampleSize++
					}
					d.FieldUFn("subframe_sample_size", func() (uint64, decode.Format, string) {
						return subframeSampleSize, decode.FormatDecimal, ""
					})

					decodeWarmupSamples := func(n uint64, sampleSize uint64) {
						d.FieldNoneFn("warmup_samples", func() {
							for i := uint64(0); i < n; i++ {
								d.FieldS("value", sampleSize)
							}
						})
					}

					decodeResiduals := func() {
						// <2> Residual coding method:
						// 00 : partitioned Rice coding with 4-bit Rice parameter; RESIDUAL_CODING_METHOD_PARTITIONED_RICE follows
						// 01 : partitioned Rice coding with 5-bit Rice parameter; RESIDUAL_CODING_METHOD_PARTITIONED_RICE2 follows
						// 10-11 : reserved
						var riceEscape uint64
						riceBits := d.FieldUFn("residual_coding_method", func() (uint64, decode.Format, string) {
							switch d.U2() {
							case 0:
								riceEscape = 15
								return 4, decode.FormatDecimal, "rice"
							case 1:
								riceEscape = 31
								return 5, decode.FormatDecimal, "rice2"
							default:
								return 0, decode.FormatDecimal, "reserved"
							}
						})

						// <4> Partition order.
						partitionOrder := d.FieldU4("partition_order")
						// There will be 2^order partitions.
						ricePartitions := uint64(1 << partitionOrder)
						d.FieldUFn("rice_partitions", func() (uint64, decode.Format, string) {
							return ricePartitions, decode.FormatDecimal, ""
						})

						for i := uint64(0); i < ricePartitions; i++ {
							d.FieldNoneFn("partition", func() {
								// Encoding parameter:
								// <4(+5)> Encoding parameter:
								// 0000-1110 : Rice parameter.
								// 1111 : Escape code, meaning the partition is in unencoded binary form using n bits per sample; n follows as a 5-bit number.
								// Or:
								// <5(+5)> Encoding parameter:
								// 00000-11110 : Rice parameter.
								// 11111 : Escape code, meaning the partition is in unencoded binary form using n bits per sample; n follows as a 5-bit number.
								// Encoded residual. The number of samples (n) in the partition is determined as follows:
								// if the partition order is zero, n = frame's blocksize - predictor order
								// else if this is not the first partition of the subframe, n = (frame's blocksize / (2^partition order))
								// else n = (frame's blocksize / (2^partition order)) - predictor order
								var count uint64
								if partitionOrder == 0 {
									count = blockSize - lpcOrder
								} else if i != 0 {
									count = blockSize / ricePartitions
								} else {
									count = (blockSize / ricePartitions) - lpcOrder
								}

								riceParameter := d.FieldU("rice_parameter", riceBits)
								if riceParameter == riceEscape {
									escapeSampleSize := d.FieldU5("escape_sample_size")
									d.FieldBytes("samples", count*escapeSampleSize)
								} else {
									d.FieldNoneFn("samples", func() {
										for j := uint64(0); j < count; j++ {
											high := d.Unary(0)
											_ = high
											low := d.U(riceParameter)
											_ = low
											// r = zigzag(high<<riceParameter | $low)
										}
									})
								}
							})
						}
					}

					switch subframeType {
					case SubframeConstant:
						// <n> Unencoded constant value of the subblock, n = frame's bits-per-sample.
						d.FieldS("value", subframeSampleSize)
					case SubframeVerbatim:
						// <n> Unencoded warm-up samples (n = frame's bits-per-sample * predictor order).
						d.FieldBytes("samples", blockSize*subframeSampleSize)
					case SubframeFixed:
						// <n> Unencoded warm-up samples (n = frame's bits-per-sample * predictor order).
						decodeWarmupSamples(lpcOrder, subframeSampleSize)
						// Encoded residual
						decodeResiduals()
					case SubframeLPC:
						// <n> Unencoded warm-up samples (n = frame's bits-per-sample * lpc order).
						decodeWarmupSamples(lpcOrder, subframeSampleSize)
						// <4> (Quantized linear predictor coefficients' precision in bits)-1 (1111 = invalid).
						precision := d.FieldUFn("precision", func() (uint64, decode.Format, string) {
							return d.U4() + 1, decode.FormatDecimal, ""
						})
						// <5> Quantized linear predictor coefficient shift needed in bits (NOTE: this number is signed two's-complement).
						d.FieldS5("shift")
						// <n> Unencoded predictor coefficients (n = qlp coeff precision * lpc order) (NOTE: the coefficients are signed two's-complement).
						d.FieldNoneFn("coefficients", func() {
							for i := uint64(0); i < lpcOrder; i++ {
								d.FieldS("value", precision)
							}
						})
						// Encoded residual
						decodeResiduals()
					}
				})
			}

			// <?> Zero-padding to byte alignment.
			d.FieldVerifyUFn("byte_align", 0, func() uint64 { return d.U(d.ByteAlignBits()) })
			// <16> CRC-16 (polynomial = x^16 + x^15 + x^2 + x^0, initialized with 0) of everything before the crc, back to and including the frame header sync code
			d.FieldU16("footer_crc")
		})
	}

	return true
}
