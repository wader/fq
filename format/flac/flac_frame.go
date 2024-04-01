package flac

import (
	"encoding/binary"
	"math/bits"

	"github.com/wader/fq/format"
	"github.com/wader/fq/internal/mathx"
	"github.com/wader/fq/pkg/checksum"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
	"github.com/wader/fq/pkg/scalar"
)

func init() {
	interp.RegisterFormat(
		format.FLAC_Frame,
		&decode.Format{
			Description: "FLAC frame",
			DecodeFn:    frameDecode,
			DefaultInArg: format.FLAC_Frame_In{
				BitsPerSample: 16, // TODO: maybe should not have a default value?
			},
		})
}

const (
	SampleRateStreaminfo = 0b0000
)

const (
	SampleSizeStreaminfo = 0b000
)

const (
	BlockingStrategyFixed    = 0
	BlockingStrategyVariable = 1
)

var BlockingStrategyNames = scalar.UintMapSymStr{
	BlockingStrategyFixed:    "fixed",
	BlockingStrategyVariable: "variable",
}

const (
	BlockSizeEndOfHeader8  = 0b0110
	BlockSizeEndOfHeader16 = 0b0111
)

const (
	SampeleRateEndOfHeader8   = 0b1100
	SampeleRateEndOfHeader16  = 0b1101
	SampeleRateEndOfHeader160 = 0b1110
)

const (
	SubframeConstant = "constant"
	SubframeVerbatim = "verbatim"
	SubframeFixed    = "fixed"
	SubframeLPC      = "lpc"
	SubframeReserved = "reserved"
)

const (
	ChannelLeftSide  = 0b1000
	ChannelSideRight = 0b1001
	ChannelMidSide   = 0b1010
)

const (
	ResidualCodingMethodRice  = 0b00
	ResidualCodingMethodRice2 = 0b01
)

var ResidualCodingMethodMap = scalar.UintMap{
	ResidualCodingMethodRice:  scalar.Uint{Sym: uint64(4), Description: "rice"},
	ResidualCodingMethodRice2: scalar.Uint{Sym: uint64(5), Description: "rice2"},
}

// TODO: generic enough?
func utf8Uint(d *decode.D) uint64 {
	n := d.U8()
	// leading ones, bit negate and count zeroes
	c := bits.LeadingZeros8(^uint8(n))
	// 0b0xxxxxxx 1 byte
	// 0b110xxxxx 2 byte
	// 0b1110xxxx 3 byte
	// 0b11110xxx 4 byte
	switch c {
	case 0:
		// nop
	case 2, 3, 4:
		n = n & ((1 << (8 - c - 1)) - 1)
		for i := 1; i < c; i++ {
			n = n<<6 | d.U8()&0x3f
		}
	default:
		d.Errorf("invalid UTF8Uint")
	}
	return n
}

// in argument is an optional FlacFrameIn struct with stream info
func frameDecode(d *decode.D) any {
	frameStart := d.Pos()
	blockSize := 0
	channelAssignment := uint64(0)
	channels := 0
	sampleSize := 0
	sideChannelIndex := -1

	var ffi format.FLAC_Frame_In
	if d.ArgAs(&ffi) {
		sampleSize = ffi.BitsPerSample
	}

	d.FieldStruct("header", func(d *decode.D) {
		// <14> 11111111111110
		d.FieldU14("sync", d.UintAssert(0b11111111111110), scalar.UintBin)

		// <1> Reserved
		// 0 : mandatory value
		// 1 : reserved for future use
		d.FieldU1("reserved0", d.UintAssert(0))

		// <1> Blocking strategy:
		// 0 : fixed-blocksize stream; frame header encodes the frame number
		// 1 : variable-blocksize stream; frame header encodes the sample number
		blockingStrategy := d.FieldU1("blocking_strategy", BlockingStrategyNames)

		// <4> Block size in inter-channel samples:
		// 0000 : reserved
		// 0001 : 192 samples
		// 0010-0101 : 576 * (2^(n-2)) samples, i.e. 576/1152/2304/4608
		// 0110 : get 8 bit (blocksize-1) from end of header
		// 0111 : get 16 bit (blocksize-1) from end of header
		// 1000-1111 : 256 * (2^(n-8)) samples, i.e. 256/512/1024/2048/4096/8192/16384/32768
		var blockSizeMap = scalar.UintMap{
			0b0000: {Description: "reserved"},
			0b0001: {Sym: uint64(192)},
			0b0010: {Sym: uint64(576)},
			0b0011: {Sym: uint64(1152)},
			0b0100: {Sym: uint64(2304)},
			0b0101: {Sym: uint64(4608)},
			0b0110: {Description: "end of header (8 bit)"},
			0b0111: {Description: "end of header (16 bit)"},
			0b1000: {Sym: uint64(256)},
			0b1001: {Sym: uint64(512)},
			0b1010: {Sym: uint64(1024)},
			0b1011: {Sym: uint64(2048)},
			0b1100: {Sym: uint64(4096)},
			0b1101: {Sym: uint64(8192)},
			0b1110: {Sym: uint64(16384)},
			0b1111: {Sym: uint64(32768)},
		}
		blockSizeS := d.FieldScalarU4("block_size", blockSizeMap, scalar.UintBin)
		if blockSizeS.Sym != nil {
			blockSize = int(blockSizeS.SymUint())
		}

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
		var sampleRateMap = scalar.UintMap{
			0b0000: {Description: "from streaminfo"},
			0b0001: {Sym: uint64(88200)},
			0b0010: {Sym: uint64(176400)},
			0b0011: {Sym: uint64(192000)},
			0b0100: {Sym: uint64(8000)},
			0b0101: {Sym: uint64(16000)},
			0b0110: {Sym: uint64(22050)},
			0b0111: {Sym: uint64(24000)},
			0b1000: {Sym: uint64(32000)},
			0b1001: {Sym: uint64(44100)},
			0b1010: {Sym: uint64(48000)},
			0b1011: {Sym: uint64(96000)},
			0b1100: {Description: "end of header (8 bit*1000)"},
			0b1101: {Description: "end of header (16 bit)"},
			0b1110: {Description: "end of header (16 bit*10)"},
			0b1111: {Description: "invalid"},
		}
		sampleRateS := d.FieldScalarU4("sample_rate", sampleRateMap, scalar.UintBin)

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
		// TODO: extract to tables and cleanup
		var channelAssignmentMap = scalar.UintMap{
			0:      {Sym: uint64(1), Description: "mono"},
			1:      {Sym: uint64(2), Description: "lr"},
			2:      {Sym: uint64(3), Description: "lrc"},
			3:      {Sym: uint64(4), Description: "fl,fr,bl,br"},
			4:      {Sym: uint64(5), Description: "fl,fr,fc,back/surround left,back/surround right"},
			5:      {Sym: uint64(6), Description: "fl,fr,fc,lfe,back/surround left,back/surround right"},
			6:      {Sym: uint64(7), Description: "fl,fr,fc,lfe,back center,sl,sr"},
			7:      {Sym: uint64(8), Description: "fl,fr,fc,lfe,back left,br,sl,sr"},
			0b1000: {Sym: uint64(2), Description: "left/side stereo"},
			0b1001: {Sym: uint64(2), Description: "right/side stereo"},
			0b1010: {Sym: uint64(2), Description: "mid/side stereo"},
			0b1011: {Sym: nil, Description: "reserved"},
			0b1100: {Sym: nil, Description: "reserved"},
			0b1101: {Sym: nil, Description: "reserved"},
			0b1111: {Sym: nil, Description: "reserved"},
		}
		channelAssignmentUint := d.FieldScalarU4("channel_assignment", channelAssignmentMap)
		if channelAssignmentUint.Sym == nil {
			d.Fatalf("unknown number of channels")
		}
		channelAssignment = channelAssignmentUint.Actual
		channels = int(channelAssignmentUint.SymUint())
		switch channelAssignmentUint.Actual {
		case ChannelLeftSide:
			sideChannelIndex = 1
		case ChannelSideRight:
			sideChannelIndex = 0
		case ChannelMidSide:
			sideChannelIndex = 1
		}
		if sideChannelIndex != -1 {
			d.FieldValueUint("side_channel_index", uint64(sideChannelIndex))
		}

		// <3> Sample size in bits:
		// 000 : get from STREAMINFO metadata block
		// 001 : 8 bits per sample
		// 010 : 12 bits per sample
		// 011 : reserved
		// 100 : 16 bits per sample
		// 101 : 20 bits per sample
		// 110 : 24 bits per sample
		// 111 : reserved
		var sampleSizeMap = scalar.UintMap{
			0b000: {Description: "from streaminfo"},
			0b001: {Sym: uint64(8)},
			0b010: {Sym: uint64(12)},
			0b011: {Description: "reserved"},
			0b100: {Sym: uint64(16)},
			0b101: {Sym: uint64(20)},
			0b110: {Sym: uint64(24)},
			0b111: {Sym: uint64(32)},
		}
		sampleSizeS := d.FieldScalarU3("sample_size", sampleSizeMap, scalar.UintBin)
		switch sampleSizeS.Actual {
		case SampleSizeStreaminfo:
			sampleSize = ffi.BitsPerSample
		default:
			if sampleSizeS.Sym != nil {
				sampleSize = int(sampleSizeS.SymUint())
			}
		}

		// <1> Reserved:
		// 0 : mandatory value
		// 1 : reserved for future use
		d.FieldU1("reserved1", d.UintAssert(0))

		d.FieldStruct("end_of_header", func(d *decode.D) {
			// if(variable blocksize)
			//   <8-56>:"UTF-8" coded sample number (decoded number is 36 bits) [4]
			// else
			//   <8-48>:"UTF-8" coded frame number (decoded number is 31 bits) [4]
			// 0 : fixed-blocksize stream; frame header encodes the frame number
			// 1 : variable-blocksize stream; frame header encodes the sample number
			switch blockingStrategy {
			case BlockingStrategyVariable:
				d.FieldUintFn("sample_number", utf8Uint)
			case BlockingStrategyFixed:
				d.FieldUintFn("frame_number", utf8Uint)
			}

			// if(blocksize bits == 011x)
			//   8/16 bit (blocksize-1)
			// 0110 : get 8 bit (blocksize-1) from end of header
			// 0111 : get 16 bit (blocksize-1) from end of header
			switch blockSizeS.Actual {
			case BlockSizeEndOfHeader8:
				blockSize = int(d.FieldU8("block_size", scalar.UintActualAdd(1)))
			case BlockSizeEndOfHeader16:
				blockSize = int(d.FieldU16("block_size", scalar.UintActualAdd(1)))
			}

			// if(sample rate bits == 11xx)
			//   8/16 bit sample rate
			// 1100 : get 8 bit sample rate (in kHz) from end of header
			// 1101 : get 16 bit sample rate (in Hz) from end of header
			// 1110 : get 16 bit sample rate (in tens of Hz) from end of header
			switch sampleRateS.Actual {
			case SampeleRateEndOfHeader8:
				d.FieldUintFn("sample_rate", func(d *decode.D) uint64 { return d.U8() * 1000 })
			case SampeleRateEndOfHeader16:
				d.FieldU16("sample_rate")
			case SampeleRateEndOfHeader160:
				d.FieldUintFn("sample_rate", func(d *decode.D) uint64 { return d.U16() * 10 })
			}
		})

		headerCRC := &checksum.CRC{Bits: 8, Table: checksum.ATM8Table}
		d.CopyBits(headerCRC, d.BitBufRange(frameStart, d.Pos()-frameStart))
		d.FieldU8("crc", d.UintValidateBytes(headerCRC.Sum(nil)), scalar.UintHex)
	})

	var channelSamples [][]int64
	d.FieldArray("subframes", func(d *decode.D) {
		for channelIndex := 0; channelIndex < channels; channelIndex++ {
			d.FieldStruct("subframe", func(d *decode.D) {
				// <1> Zero bit padding, to prevent sync-fooling string of 1s
				d.FieldU1("zero_bit", d.UintAssert(0))

				// <6> Subframe type:
				// 000000 : SUBFRAME_CONSTANT
				// 000001 : SUBFRAME_VERBATIM
				// 00001x : reserved
				// 0001xx : reserved
				// 001xxx : if(xxx <= 4) SUBFRAME_FIXED, xxx=order ; else reserved
				// 01xxxx : reserved
				// 1xxxxx : SUBFRAME_LPC, xxxxx=order-1
				lpcOrder := -1
				var subframeTypeRangeMap = scalar.UintRangeToScalar{
					{Range: [2]uint64{0b000000, 0b000000}, S: scalar.Uint{Sym: SubframeConstant}},
					{Range: [2]uint64{0b000001, 0b000001}, S: scalar.Uint{Sym: SubframeVerbatim}},
					{Range: [2]uint64{0b000010, 0b000011}, S: scalar.Uint{Sym: SubframeReserved}},
					{Range: [2]uint64{0b000100, 0b000111}, S: scalar.Uint{Sym: SubframeReserved}},
					{Range: [2]uint64{0b001000, 0b001100}, S: scalar.Uint{Sym: SubframeFixed}},
					{Range: [2]uint64{0b001101, 0b001111}, S: scalar.Uint{Sym: SubframeReserved}},
					{Range: [2]uint64{0b010000, 0b011111}, S: scalar.Uint{Sym: SubframeReserved}},
					{Range: [2]uint64{0b100000, 0b111111}, S: scalar.Uint{Sym: SubframeLPC}},
				}
				subframeTypeUint := d.FieldScalarU6("subframe_type", subframeTypeRangeMap, scalar.UintBin)
				switch subframeTypeUint.SymStr() {
				case SubframeFixed:
					lpcOrder = int(subframeTypeUint.Actual & 0b111)
				case SubframeLPC:
					lpcOrder = int((subframeTypeUint.Actual & 0b11111) + 1)
				}
				if lpcOrder != -1 {
					d.FieldValueUint("lpc_order", uint64(lpcOrder))
				}

				// 'Wasted bits-per-sample' flag:
				// 0 : no wasted bits-per-sample in source subblock, k=0
				// 1 : k wasted bits-per-sample in source subblock, k-1 follows, unary coded; e.g. k=3 => 001 follows, k=7 => 0000001 follows.
				wastedBitsFlag := d.FieldU1("wasted_bits_flag")
				var wastedBitsK int
				if wastedBitsFlag != 0 {
					wastedBitsK = int(d.FieldUnary("wasted_bits_k", 0, scalar.UintActualAdd(1)))
				}

				subframeSampleSize := sampleSize - wastedBitsK
				if subframeSampleSize < 1 {
					d.Fatalf("subframeSampleSize %d < 1", subframeSampleSize)
				}
				// if channel is side, add en extra sample bit
				// https://github.com/xiph/flac/blob/37e675b777d4e0de53ac9ff69e2aea10d92e729c/src/libFLAC/stream_decoder.c#L2040
				if channelIndex == sideChannelIndex {
					subframeSampleSize++
				}
				d.FieldValueUint("subframe_sample_size", uint64(subframeSampleSize))

				decodeWarmupSamples := func(samples []int64, n int, sampleSize int) {
					if len(samples) < n {
						d.Fatalf("decodeWarmupSamples outside block size")
					}

					d.FieldArray("warmup_samples", func(d *decode.D) {
						for i := 0; i < n; i++ {
							samples[i] = d.FieldS("value", sampleSize)
						}
					})
				}

				decodeResiduals := func(samples []int64) {
					samplesLen := len(samples)
					n := 0

					// <2> Residual coding method:
					// 00 : partitioned Rice coding with 4-bit Rice parameter; RESIDUAL_CODING_METHOD_PARTITIONED_RICE follows
					// 01 : partitioned Rice coding with 5-bit Rice parameter; RESIDUAL_CODING_METHOD_PARTITIONED_RICE2 follows
					// 10-11 : reserved
					var riceEscape int
					var riceBits int
					residualCodingMethod := d.FieldU2("residual_coding_method", scalar.UintMap{
						0b00: scalar.Uint{Sym: uint64(4), Description: "rice"},
						0b01: scalar.Uint{Sym: uint64(5), Description: "rice2"},
					})
					switch residualCodingMethod {
					case ResidualCodingMethodRice:
						riceEscape = 0b1111
						riceBits = 4
					case ResidualCodingMethodRice2:
						riceEscape = 0b11111
						riceBits = 5
					}

					// <4> Partition order.
					partitionOrder := int(d.FieldU4("partition_order"))
					// There will be 2^order partitions.
					ricePartitions := 1 << partitionOrder
					d.FieldValueUint("rice_partitions", uint64(ricePartitions))

					d.FieldArray("partitions", func(d *decode.D) {
						for i := 0; i < ricePartitions; i++ {
							d.FieldStruct("partition", func(d *decode.D) {
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
								var count int
								if partitionOrder == 0 {
									count = blockSize - lpcOrder
								} else if i != 0 {
									count = blockSize / ricePartitions
								} else {
									count = (blockSize / ricePartitions) - lpcOrder
								}

								d.FieldValueUint("count", uint64(count))

								riceParameter := int(d.FieldU("rice_parameter", riceBits))

								if samplesLen < n+count {
									d.Fatalf("decodeResiduals outside block size")
								}
								if count < 0 {
									d.Fatalf("negative sample count %d", count)
								}

								if riceParameter == riceEscape {
									escapeSampleSize := int(d.FieldU5("escape_sample_size"))
									if escapeSampleSize == 0 {
										// Zero sample size, we can just skip ahead count samples as they are already zero. From spec:
										// Note that it is possible that the number of bits is 0, which means all residual samples in that partition have
										// a value of 0, and no bits code for the partition itself.
										n += count
									} else {
										d.RangeFn(d.Pos(), int64(count*escapeSampleSize), func(d *decode.D) {
											d.FieldRawLen("samples", int64(count*escapeSampleSize))
										})
										for j := 0; j < count; j++ {
											samples[n] = d.S(escapeSampleSize)
											n++
										}
									}
								} else {
									samplesStart := d.Pos()
									for j := 0; j < count; j++ {
										high := d.Unary(0)
										low := d.U(riceParameter)
										samples[n] = mathx.ZigZag[uint64, int64](high<<riceParameter | low)
										n++
									}
									samplesStop := d.Pos()
									d.RangeFn(samplesStart, samplesStop-samplesStart, func(d *decode.D) {
										d.FieldRawLen("samples", d.BitsLeft())
									})
								}
							})
						}
					})
				}

				// modifies input samples slice and returns it
				decodeLPC := func(lpcOrder int, samples []int64, coeffs []int64, shift int64) {
					for i := lpcOrder; i < len(samples); i++ {
						r := int64(0)
						for j := 0; j < len(coeffs); j++ {
							c := coeffs[j]
							s := samples[i-j-1]
							r += c * s
						}
						samples[i] = samples[i] + (r >> shift)
					}
				}

				samples := make([]int64, blockSize)
				switch subframeTypeUint.SymStr() {
				case SubframeConstant:
					// <n> Unencoded constant value of the subblock, n = frame's bits-per-sample.
					v := d.FieldS("value", subframeSampleSize)
					for i := 0; i < blockSize; i++ {
						samples[i] = v
					}
				case SubframeVerbatim:
					// <n*i> Unencoded subblock; n = frame's bits-per-sample, i = frame's blocksize.
					// TODO: refactor into some kind of FieldBitBufLenFn?
					d.RangeFn(d.Pos(), int64(blockSize*subframeSampleSize), func(d *decode.D) {
						d.FieldRawLen("samples", d.BitsLeft())
					})

					for i := 0; i < blockSize; i++ {
						samples[i] = d.S(subframeSampleSize)
					}
				case SubframeFixed:
					// <n> Unencoded warm-up samples (n = frame's bits-per-sample * predictor order).
					decodeWarmupSamples(samples, lpcOrder, subframeSampleSize)
					// Encoded residual
					decodeResiduals(samples[lpcOrder:])
					// http://www.hpl.hp.com/techreports/1999/HPL-1999-144.pdf
					fixedCoeffs := [][]int64{
						{},
						{1},
						{2, -1},
						{3, -3, 1},
						{4, -6, 4, -1},
					}
					coeffs := fixedCoeffs[lpcOrder]
					decodeLPC(lpcOrder, samples, coeffs, 0)
				case SubframeLPC:
					// <n> Unencoded warm-up samples (n = frame's bits-per-sample * lpc order).
					decodeWarmupSamples(samples, lpcOrder, subframeSampleSize)
					// <4> (Quantized linear predictor coefficients' precision in bits)-1 (1111 = invalid).
					precision := int(d.FieldU4("precision", scalar.UintActualAdd(1)))
					// <5> Quantized linear predictor coefficient shift needed in bits (NOTE: this number is signed two's-complement).
					shift := d.FieldS5("shift")
					if shift < 0 {
						d.Fatalf("negative LPC shift %d", shift)
					}
					// <n> Unencoded predictor coefficients (n = qlp coeff precision * lpc order) (NOTE: the coefficients are signed two's-complement).
					var coeffs []int64
					d.FieldArray("coefficients", func(d *decode.D) {
						for i := 0; i < lpcOrder; i++ {
							coeffs = append(coeffs, d.FieldS("value", precision))
						}
					})
					// Encoded residual
					decodeResiduals(samples[lpcOrder:])
					decodeLPC(lpcOrder, samples, coeffs, shift)
				}

				if wastedBitsK != 0 {
					for i := 0; i < len(samples); i++ {
						samples[i] <<= wastedBitsK
					}
				}

				channelSamples = append(channelSamples, samples)
			})
		}
	})

	// <?> Zero-padding to byte alignment.
	d.FieldU("byte_align", d.ByteAlignBits(), d.UintAssert(0))
	// <16> CRC-16 (polynomial = x^16 + x^15 + x^2 + x^0, initialized with 0) of everything before the crc, back to and including the frame header sync code
	footerCRC := &checksum.CRC{Bits: 16, Table: checksum.ANSI16Table}
	d.CopyBits(footerCRC, d.BitBufRange(frameStart, d.Pos()-frameStart))
	d.FieldRawLen("footer_crc", 16, d.ValidateBitBuf(footerCRC.Sum(nil)), scalar.RawHex)

	streamSamples := len(channelSamples[0])
	for j := 0; j < len(channelSamples); j++ {
		if streamSamples > len(channelSamples[j]) {
			d.Fatalf("different amount of samples in channels %d != %d", streamSamples, len(channelSamples[j]))
		}
	}

	// Transform mid/side channels into left, right
	// mid = (left + right)/2
	// side = left - right
	switch channelAssignment {
	case ChannelLeftSide:
		for i := 0; i < len(channelSamples[0]); i++ {
			channelSamples[1][i] = channelSamples[0][i] - channelSamples[1][i]
		}
	case ChannelSideRight:
		for i := 0; i < len(channelSamples[0]); i++ {
			channelSamples[0][i] = channelSamples[1][i] + channelSamples[0][i]
		}
	case ChannelMidSide:
		for i := 0; i < len(channelSamples[0]); i++ {
			m := channelSamples[0][i]
			s := channelSamples[1][i]
			m = m<<1 | s&1
			channelSamples[0][i] = (m + s) >> 1
			channelSamples[1][i] = (m - s) >> 1
		}
	default:
		// not stereo or no side channel
	}

	outSampleSize := sampleSize + (sampleSize % 8)

	bytesPerSample := outSampleSize / 8
	p := 0
	le := binary.LittleEndian

	interleavedSamplesBuf := ffi.SamplesBuf
	interleavedSamplesBufLen := len(channelSamples) * streamSamples * bytesPerSample
	// TODO: decode read buffer?
	// TODO: reuse buffer if possible
	if interleavedSamplesBuf == nil || len(interleavedSamplesBuf) < interleavedSamplesBufLen {
		interleavedSamplesBuf = make([]byte, interleavedSamplesBufLen)
	}

	// TODO: speedup by using more cache friendly memory layout for samples
	for i := 0; i < streamSamples; i++ {
		for j := 0; j < len(channelSamples); j++ {

			s := channelSamples[j][i]
			switch outSampleSize {
			case 8:
				interleavedSamplesBuf[p] = byte(s)
			case 16:
				le.PutUint16(interleavedSamplesBuf[p:], uint16(s))
			case 24:
				interleavedSamplesBuf[p] = byte(s)
				le.PutUint16(interleavedSamplesBuf[p+1:], uint16(s>>8))
			case 32:
				le.PutUint32(interleavedSamplesBuf[p:], uint32(s))
			}
			p += bytesPerSample
		}
	}

	return format.FLAC_Frame_Out{
		SamplesBuf:    interleavedSamplesBuf,
		Samples:       uint64(streamSamples),
		Channels:      channels,
		BitsPerSample: outSampleSize,
	}
}
