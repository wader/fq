package mpeg

// ISO/IEC 23003-3
// MediaInfoLib/Source/MediaInfo/Audio/File_Usac.cpp
// TODO: still lacking lots of things

import (
	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
	"github.com/wader/fq/pkg/scalar"
)

func init() {
	interp.RegisterFormat(decode.Format{
		Name:        format.MPEG_USACC,
		Description: "MPEG Unified Speech and Audio Coding config",
		DecodeFn:    usaccDecoder,
	})
}

type usacCoreSBRFrameMapEntry struct {
	ratioIndex                  int
	outputFrameLengthDivided256 int
}

var usacCoreSBRFramehMap = map[uint64]usacCoreSBRFrameMapEntry{
	0: {ratioIndex: 0, outputFrameLengthDivided256: 3},
	1: {ratioIndex: 0, outputFrameLengthDivided256: 4},
	2: {ratioIndex: 2, outputFrameLengthDivided256: 8},
	3: {ratioIndex: 3, outputFrameLengthDivided256: 8},
	4: {ratioIndex: 1, outputFrameLengthDivided256: 16},
}

type usacCoreSBRFramehRatioIndexSymMapper map[uint64]usacCoreSBRFrameMapEntry

func (sm usacCoreSBRFramehRatioIndexSymMapper) MapScalar(s scalar.Uint) (scalar.Uint, error) {
	if e, ok := sm[uint64(s.Actual)]; ok {
		s.Sym = e.ratioIndex
	}
	return s, nil
}

type usacCoreSBRFramehOutputFrameLengthSymMapper map[uint64]usacCoreSBRFrameMapEntry

func (sm usacCoreSBRFramehOutputFrameLengthSymMapper) MapUint(s scalar.Uint) (scalar.Uint, error) {
	if e, ok := sm[s.Actual]; ok {
		s.Sym = e.outputFrameLengthDivided256
	}
	return s, nil
}

const (
	usacExtElementConfigFill         = 0
	usacExtElementConfigMPEGS        = 1
	usacExtElementConfigSAOC         = 2
	usacExtElementConfigAudioPreRoll = 3
	usacExtElementConfigUniDRC       = 4
)

var usacExtElementConfigMap = scalar.UintMap{
	usacExtElementConfigFill:         {Sym: "fill", Description: "Fill"},
	usacExtElementConfigMPEGS:        {Sym: "mpegs", Description: "MPEGS"},
	usacExtElementConfigSAOC:         {Sym: "saoc", Description: "SAOC"},
	usacExtElementConfigAudioPreRoll: {Sym: "audio_pre_roll", Description: "AudioPreRoll"},
	usacExtElementConfigUniDRC:       {Sym: "unidrc", Description: "UniDRC"},
}

func usacDecodeUniDRC(d *decode.D) {
	sampleRatePresent := d.FieldBool("sample_rate_present")
	if sampleRatePresent {
		d.FieldU18("sample_rate")
	}
	_ = d.FieldU7("downmix_instruction_count")
	descriptionBasicPresent := d.FieldBool("description_basic_present")
	if descriptionBasicPresent {
		d.FieldU3("coefficients_basic_count")
		d.FieldU4("instructions_basic_count")
	}
}

func usacDecoderExtElementConfig(d *decode.D) {
	extType := d.FieldUintFn("ext_type", decodeEscapeValueAddFn(4, 8, 16), usacExtElementConfigMap)
	d.FieldUintFn("length", decodeEscapeValueAddFn(4, 8, 16))
	var defaultLength uint64
	if d.FieldBool("default_length_present") {
		defaultLength = d.FieldUintFn("default_length", decodeEscapeValueAddFn(8, 16, 0))
	}
	// TODO: type
	d.FieldBool("payload_flag")
	_ = defaultLength

	switch extType {
	case usacExtElementConfigUniDRC:
		usacDecodeUniDRC(d)
	}
}

func usacDecoderCoreConfig(d *decode.D) {
	d.FieldBool("tw_mdct")
	d.FieldBool("noise_filling")
}

func usacDecoderSingleChannelElement(d *decode.D, coreSbrFrameLengthIndex uint64) {
	usacDecoderCoreConfig(d)
	if usacCoreSBRFramehMap[coreSbrFrameLengthIndex].ratioIndex > 0 {
		// TODO:
	}
}

const (
	usacElementTypeSCE = 0
	usacElementTypeCPE = 1
	usacElementTypeLFE = 2
	usacElementTypeEXT = 3
)

var usacElementTypeMap = scalar.UintMap{
	usacElementTypeSCE: {Sym: "sce", Description: "Single channel element"},
	usacElementTypeCPE: {Sym: "cpe", Description: "Channel pair element"},
	usacElementTypeLFE: {Sym: "lfe", Description: "Low frequency effect"},
	usacElementTypeEXT: {Sym: "ext", Description: "Extension"},
}

func usacDecoderConfig(d *decode.D, coreSbrFrameLengthIndex uint64) {
	numElements := d.FieldUintFn("num_elements", decodeEscapeValueAddFn(4, 8, 16), scalar.UintActualAdd(1))
	d.FieldArray("elements", func(d *decode.D) {
		for i := uint64(0); i < numElements; i++ {
			d.FieldStruct("element", func(d *decode.D) {
				switch d.FieldU2("type", usacElementTypeMap) {
				case usacElementTypeSCE:
					// TODO:
				case usacElementTypeCPE:
					usacDecoderSingleChannelElement(d, coreSbrFrameLengthIndex)
				case usacElementTypeLFE:
					// TODO:
				case usacElementTypeEXT:
					usacDecoderExtElementConfig(d)
				}
			})
		}
	})
}

var methodDefinitionSizes = map[uint64]int64{
	0: 8,
	1: 8,
	2: 8,
	3: 8,
	4: 8,
	5: 8,
	6: 8,
	7: 5,
	8: 2,
	9: 8,
}

// TODO: v1
func usacDecoderLoudnessInfo(d *decode.D, isAlbum bool) {
	d.FieldU6("drc_set_id")
	d.FieldU7("downmix_id")
	if d.FieldBool("sample_peak_level_present") {
		d.FieldU12("sample_peak_level")
	}
	if d.FieldBool("true_peak_level_present") {
		d.FieldU12("true_peak_level")
		d.FieldU4("measure_system")
		d.FieldU2("reliability")
	}
	measureCount := d.FieldU4("measure_count")
	d.FieldArray("measures", func(d *decode.D) {
		for i := uint64(0); i < measureCount; i++ {
			d.FieldStruct("measure", func(d *decode.D) {
				methodDefinition := d.FieldU4("method_definition")
				if methodDefinition >= uint64(len(methodDefinitionSizes)) {
					return
				}
				d.FieldU("method_value", int(methodDefinitionSizes[methodDefinition]))
				d.FieldU4("measure_system")
				d.FieldU2("reliability")
			})
		}
	})
}

func usacDecoderLoudnessInfoSet(d *decode.D) {
	albumCount := d.FieldU6("album_count")
	count := d.FieldU6("count")
	d.FieldArray("albums", func(d *decode.D) {
		for i := uint64(0); i < albumCount; i++ {
			usacDecoderLoudnessInfo(d, true)
		}
	})
	d.FieldArray("items", func(d *decode.D) {
		for i := uint64(0); i < count; i++ {
			d.FieldStruct("item", func(d *decode.D) { usacDecoderLoudnessInfo(d, false) })
		}
	})
	if d.BitsLeft() > 0 {
		d.FieldRawLen("padding", d.BitsLeft())
	}
}

const (
	usacConfigExtensionFill         = 0
	usacConfigExtensionLoudnessInfo = 2
	usacConfigExtensionStreamID     = 7
)

var usacConfigExtensionMap = scalar.UintMap{
	usacConfigExtensionFill:         {Sym: "fill", Description: "Fill"},
	usacConfigExtensionLoudnessInfo: {Sym: "loudness_info", Description: "Loudness info"},
	usacConfigExtensionStreamID:     {Sym: "stream_id", Description: "Stream ID"},
}

func usacDecoderConfigExtension(d *decode.D) {
	numConfigExtensions := d.FieldUintFn("num_config_extensions", decodeEscapeValueAddFn(2, 4, 8), scalar.UintActualAdd(1))
	d.FieldArray("extensions", func(d *decode.D) {
		for i := uint64(0); i < numConfigExtensions; i++ {
			d.FieldStruct("extension", func(d *decode.D) {
				typ := d.FieldUintFn("type", decodeEscapeValueAddFn(4, 8, 16), usacConfigExtensionMap)
				length := d.FieldUintFn("length", decodeEscapeValueAddFn(4, 8, 16))

				d.LimitedFn(int64(length)*8, func(d *decode.D) {
					switch typ {
					case usacConfigExtensionFill:
						if length > 0 {
							d.FieldRawLen("filling", int64(length)*8)
						}
					case usacConfigExtensionLoudnessInfo:
						d.FieldStruct("loudness", usacDecoderLoudnessInfoSet)
					case usacConfigExtensionStreamID:
						d.FieldU16("identifier")
					default:
						if length > 0 {
							d.FieldRawLen("data", int64(length)*8)
						}
					}
				})
			})
		}
	})
}

func usaccDecoder(d *decode.D) any {
	d.FieldUintFn("sampling_frequency", decodeEscapeValueAddFn(5, 24, 0), frequencyIndexHzMap)
	coreSbrFrameLengthIndex := d.FieldU3("core_sbr_frame_length", usacCoreSBRFramehOutputFrameLengthSymMapper(usacCoreSBRFramehMap))
	channelConfiguration := d.FieldU5("channel_configuration", channelConfigurationNames)
	if channelConfiguration == 0 {
		numOutChannels := d.FieldUintFn("num_out_channels", decodeEscapeValueAddFn(5, 8, 16))
		d.FieldArray("out_channels", func(d *decode.D) {
			for i := uint64(0); i < numOutChannels; i++ {
				d.FieldU5("out_channel_pos")
			}
		})
	}
	if coreSbrFrameLengthIndex >= uint64(len(usacCoreSBRFramehMap)) {
		// TODO:
		return nil
	}
	d.FieldStruct("decoder_config", func(d *decode.D) { usacDecoderConfig(d, coreSbrFrameLengthIndex) })
	if d.FieldBool("config_extension_present") {
		d.FieldStruct("config_extension", usacDecoderConfigExtension)
	}

	return nil
}
