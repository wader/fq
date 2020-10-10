package tiff

// http://www.libpng.org/pub/png/spec/1.2/PNG-Contents.html
// https://ftp-osl.osuosl.org/pub/libpng/documents/pngext-1.5.0.html
// TODO: gps

import (
	"fq/pkg/decode"
	"fq/pkg/format/register"
)

var iccTag []*decode.Format

var File = register.Register(&decode.Format{
	Name:  "tiff",
	MIMEs: []string{"image/tiff"},
	New:   func() decode.Decoder { return &FileDecoder{} },
	Deps: []decode.Dep{
		{Names: []string{"icc"}, Formats: &iccTag},
	},
})

const littleEndian = 0x49492a00
const bigEndian = 0x4d4d002a

const (
	BYTE      = 1
	ASCII     = 2
	SHORT     = 3
	LONG      = 4
	RATIONAL  = 5
	UNDEFINED = 7
	SLONG     = 9
	SRATIONAL = 10
)

var typeNames = map[uint64]string{
	BYTE:      "BYTE",
	ASCII:     "ASCII",
	SHORT:     "SHORT",
	LONG:      "LONG",
	RATIONAL:  "RATIONAL",
	UNDEFINED: "UNDEFINED",
	SLONG:     "SLONG",
	SRATIONAL: "SRATIONAL",
}

// TODO: tiff 6.0 types
var typeByteSize = map[uint64]uint64{
	BYTE:      1,
	ASCII:     1,
	SHORT:     2,
	LONG:      4,
	RATIONAL:  4 + 4,
	UNDEFINED: 1,
	SLONG:     4,
	SRATIONAL: 4 + 4,
}

const (
	NewSubfileType               = 254
	SubfileType                  = 255
	ImageWidth                   = 256
	ImageLength                  = 257
	BitsPerSample                = 258
	Compression                  = 259
	PhotometricInterpretation    = 262
	Threshholding                = 263
	CellWidth                    = 264
	CellLength                   = 265
	FillOrder                    = 266
	DocumentName                 = 269
	ImageDescription             = 270
	Make                         = 271
	Model                        = 272
	StripOffsets                 = 273
	Orientation                  = 274
	SamplesPerPixel              = 277
	RowsPerStrip                 = 278
	StripByteCounts              = 279
	MinSampleValue               = 280
	MaxSampleValue               = 281
	XResolution                  = 282
	YResolution                  = 283
	PlanarConfiguration          = 284
	PageName                     = 285
	XPosition                    = 286
	YPosition                    = 287
	FreeOffsets                  = 288
	FreeByteCounts               = 289
	GrayResponseUnit             = 290
	GrayResponseCurve            = 291
	T4Options                    = 292
	T6Options                    = 293
	ResolutionUnit               = 296
	PageNumber                   = 297
	TransferFunction             = 301
	Software                     = 305
	DateTime                     = 306
	Artist                       = 315
	HostComputer                 = 316
	Predictor                    = 317
	WhitePoint                   = 318
	PrimaryChromaticities        = 319
	ColorMap                     = 320
	HalftoneHints                = 321
	TileWidth                    = 322
	TileLength                   = 323
	TileOffsets                  = 324
	TileByteCounts               = 325
	BadFaxLines                  = 326
	CleanFaxData                 = 327
	ConsecutiveBadFaxLines       = 328
	SubIFDs                      = 330
	InkSet                       = 332
	InkNames                     = 333
	NumberOfInks                 = 334
	DotRange                     = 336
	TargetPrinter                = 337
	ExtraSamples                 = 338
	SampleFormat                 = 339
	SMinSampleValue              = 340
	SMaxSampleValue              = 341
	TransferRange                = 342
	ClipPath                     = 343
	XClipPathUnits               = 344
	YClipPathUnits               = 345
	Indexed                      = 346
	JPEGTables                   = 347
	OPIProxy                     = 351
	GlobalParametersIFD          = 400
	ProfileType                  = 401
	FaxProfile                   = 402
	CodingMethods                = 403
	VersionYear                  = 404
	ModeNumber                   = 405
	Decode                       = 433
	DefaultImageColor            = 434
	JPEGProc                     = 512
	JPEGInterchangeFormat        = 513
	JPEGInterchangeFormatLength  = 514
	JPEGRestartInterval          = 515
	JPEGLosslessPredictors       = 517
	JPEGPointTransforms          = 518
	JPEGQTables                  = 519
	JPEGDCTables                 = 520
	JPEGACTables                 = 521
	YCbCrCoefficients            = 529
	YCbCrSubSampling             = 530
	YCbCrPositioning             = 531
	ReferenceBlackWhite          = 532
	StripRowCounts               = 559
	XMP                          = 700
	ImageRating                  = 18246
	ImageRatingPercent           = 18249
	ImageID                      = 32781
	WangAnnotation               = 32932
	CFARepeatPatternDim          = 33421
	CFAPattern                   = 33422
	BatteryLevel                 = 33423
	Copyright                    = 33432
	ExposureTime                 = 33434
	FNumber                      = 33437
	MDFileTag                    = 33445
	MDScalePixel                 = 33446
	MDColorTable                 = 33447
	MDLabName                    = 33448
	MDSampleInfo                 = 33449
	MDPrepDate                   = 33450
	MDPrepTime                   = 33451
	MDFileUnits                  = 33452
	ModelPixelScaleTag           = 33550
	IPTC_NAA                     = 33723
	INGRPacketDataTag            = 33918
	INGRFlagRegisters            = 33919
	IrasBTransformationMatrix    = 33920
	ModelTiepointTag             = 33922
	Site                         = 34016
	ColorSequence                = 34017
	IT8Header                    = 34018
	RasterPadding                = 34019
	BitsPerRunLength             = 34020
	BitsPerExtendedRunLength     = 34021
	ColorTable                   = 34022
	ImageColorIndicator          = 34023
	BackgroundColorIndicator     = 34024
	ImageColorValue              = 34025
	BackgroundColorValue         = 34026
	PixelIntensityRange          = 34027
	TransparencyIndicator        = 34028
	ColorCharacterization        = 34029
	HCUsage                      = 34030
	TrapIndicator                = 34031
	CMYKEquivalent               = 34032
	Reserved0                    = 34033
	Reserved1                    = 34034
	Reserved2                    = 34035
	ModelTransformationTag       = 34264
	Photoshop                    = 34377
	ExifIFD                      = 34665
	InterColorProfile            = 34675
	ImageLayer                   = 34732
	GeoKeyDirectoryTag           = 34735
	GeoDoubleParamsTag           = 34736
	GeoAsciiParamsTag            = 34737
	ExposureProgram              = 34850
	SpectralSensitivity          = 34852
	GPSInfo                      = 34853
	ISOSpeedRatings              = 34855
	OECF                         = 34856
	Interlace                    = 34857
	TimeZoneOffset               = 34858
	SelfTimeMode                 = 34859
	SensitivityType              = 34864
	StandardOutputSensitivity    = 34865
	RecommendedExposureIndex     = 34866
	ISOSpeed                     = 34867
	ISOSpeedLatitudeyyy          = 34868
	ISOSpeedLatitudezzz          = 34869
	HylaFAXFaxRecvParams         = 34908
	HylaFAXFaxSubAddress         = 34909
	HylaFAXFaxRecvTime           = 34910
	ExifVersion                  = 36864
	DateTimeOriginal             = 36867
	DateTimeDigitized            = 36868
	ComponentsConfiguration      = 37121
	CompressedBitsPerPixel       = 37122
	ShutterSpeedValue            = 37377
	ApertureValue                = 37378
	BrightnessValue              = 37379
	ExposureBiasValue            = 37380
	MaxApertureValue             = 37381
	SubjectDistance              = 37382
	MeteringMode                 = 37383
	LightSource                  = 37384
	Flash                        = 37385
	FocalLength                  = 37386
	FlashEnergy                  = 37387
	SpatialFrequencyResponse     = 37388
	Noise                        = 37389
	FocalPlaneXResolution        = 37390
	FocalPlaneYResolution        = 37391
	FocalPlaneResolutionUnit     = 37392
	ImageNumber                  = 37393
	SecurityClassification       = 37394
	ImageHistory                 = 37395
	SubjectLocation              = 37396
	ExposureIndex                = 37397
	TIFF_EPStandardID            = 37398
	SensingMethod                = 37399
	MakerNote                    = 37500
	UserComment                  = 37510
	SubsecTime                   = 37520
	SubsecTimeOriginal           = 37521
	SubsecTimeDigitized          = 37522
	ImageSourceData              = 37724
	XPTitle                      = 40091
	XPComment                    = 40092
	XPAuthor                     = 40093
	XPKeywords                   = 40094
	XPSubject                    = 40095
	FlashpixVersion              = 40960
	ColorSpace                   = 40961
	PixelXDimension              = 40962
	PixelYDimension              = 40963
	RelatedSoundFile             = 40964
	InteroperabilityIFD          = 40965
	FlashEnergy2                 = 41483
	SpatialFrequencyResponse2    = 41484
	FocalPlaneXResolution2       = 41486
	FocalPlaneYResolution2       = 41487
	FocalPlaneResolutionUnit2    = 41488
	SubjectLocation2             = 41492
	ExposureIndex2               = 41493
	SensingMethod2               = 41495
	FileSource                   = 41728
	SceneType                    = 41729
	CFAPattern2                  = 41730
	CustomRendered               = 41985
	ExposureMode                 = 41986
	WhiteBalance                 = 41987
	DigitalZoomRatio             = 41988
	FocalLengthIn35mmFilm        = 41989
	SceneCaptureType             = 41990
	GainControl                  = 41991
	Contrast                     = 41992
	Saturation                   = 41993
	Sharpness                    = 41994
	DeviceSettingDescription     = 41995
	SubjectDistanceRange         = 41996
	ImageUniqueID                = 42016
	CameraOwnerName              = 42032
	BodySerialNumber             = 42033
	LensSpecification            = 42034
	LensMake                     = 42035
	LensModel                    = 42036
	LensSerialNumber             = 42037
	GDAL_METADATA                = 42112
	GDAL_NODATA                  = 42113
	HDPixelFormat                = 48129
	HDTransformation             = 48130
	HDUncompressed               = 48131
	HDImageType                  = 48132
	HDImageWidth                 = 48256
	HDImageHeight                = 48257
	HDWidthResolution            = 48258
	HDHeightResolution           = 48259
	HDImageOffset                = 48320
	HDImageByteCount             = 48321
	HDAlphaOffset                = 48322
	HDAlphaByteCount             = 48323
	HDImageDataDiscard           = 48324
	HDAlphaDataDiscard           = 48325
	OceScanjobDescription        = 50215
	OceApplicationSelector       = 50216
	OceIdentificationNumber      = 50217
	OceImageLogicCharacteristics = 50218
	PrintImageMatching           = 50341
	DNGVersion                   = 50706
	DNGBackwardVersion           = 50707
	UniqueCameraModel            = 50708
	LocalizedCameraModel         = 50709
	CFAPlaneColor                = 50710
	CFALayout                    = 50711
	LinearizationTable           = 50712
	BlackLevelRepeatDim          = 50713
	BlackLevel                   = 50714
	BlackLevelDeltaH             = 50715
	BlackLevelDeltaV             = 50716
	WhiteLevel                   = 50717
	DefaultScale                 = 50718
	DefaultCropOrigin            = 50719
	DefaultCropSize              = 50720
	ColorMatrix1                 = 50721
	ColorMatrix2                 = 50722
	CameraCalibration1           = 50723
	CameraCalibration2           = 50724
	ReductionMatrix1             = 50725
	ReductionMatrix2             = 50726
	AnalogBalance                = 50727
	AsShotNeutral                = 50728
	AsShotWhiteXY                = 50729
	BaselineExposure             = 50730
	BaselineNoise                = 50731
	BaselineSharpness            = 50732
	BayerGreenSplit              = 50733
	LinearResponseLimit          = 50734
	CameraSerialNumber           = 50735
	LensInfo                     = 50736
	ChromaBlurRadius             = 50737
	AntiAliasStrength            = 50738
	ShadowScale                  = 50739
	DNGPrivateData               = 50740
	MakerNoteSafety              = 50741
	CalibrationIlluminant1       = 50778
	CalibrationIlluminant2       = 50779
	BestQualityScale             = 50780
	RawDataUniqueID              = 50781
	AliasLayerMetadata           = 50784
	OriginalRawFileName          = 50827
	OriginalRawFileData          = 50828
	ActiveArea                   = 50829
	MaskedAreas                  = 50830
	AsShotICCProfile             = 50831
	AsShotPreProfileMatrix       = 50832
	CurrentICCProfile            = 50833
	CurrentPreProfileMatrix      = 50834
	ColorimetricReference        = 50879
	CameraCalibrationSignature   = 50931
	ProfileCalibrationSignature  = 50932
	ExtraCameraProfiles          = 50933
	AsShotProfileName            = 50934
	NoiseReductionApplied        = 50935
	ProfileName                  = 50936
	ProfileHueSatMapDims         = 50937
	ProfileHueSatMapData1        = 50938
	ProfileHueSatMapData2        = 50939
	ProfileToneCurve             = 50940
	ProfileEmbedPolicy           = 50941
	ProfileCopyright             = 50942
	ForwardMatrix1               = 50964
	ForwardMatrix2               = 50965
	PreviewApplicationName       = 50966
	PreviewApplicationVersion    = 50967
	PreviewSettingsName          = 50968
	PreviewSettingsDigest        = 50969
	PreviewColorSpace            = 50970
	PreviewDateTime              = 50971
	RawImageDigest               = 50972
	OriginalRawFileDigest        = 50973
	SubTileBlockSize             = 50974
	RowInterleaveFactor          = 50975
	ProfileLookTableDims         = 50981
	ProfileLookTableData         = 50982
	OpcodeList1                  = 51008
	OpcodeList2                  = 51009
	OpcodeList3                  = 51022
	NoiseProfile                 = 51041
	OriginalDefaultFinalSize     = 51089
	OriginalBestQualityFinalSize = 51090
	OriginalDefaultCropSize      = 51091
	ProfileHueSatMapEncoding     = 51107
	ProfileLookTableEncoding     = 51108
	BaselineExposureOffset       = 51109
	DefaultBlackRender           = 51110
	NewRawImageDigest            = 51111
	RawToPreviewGain             = 51112
	DefaultUserCrop              = 51125
)

var tagNames = map[uint64]string{
	NewSubfileType:               "NewSubfileType",
	SubfileType:                  "SubfileType",
	ImageWidth:                   "ImageWidth",
	ImageLength:                  "ImageLength",
	BitsPerSample:                "BitsPerSample",
	Compression:                  "Compression",
	PhotometricInterpretation:    "PhotometricInterpretation",
	Threshholding:                "Threshholding",
	CellWidth:                    "CellWidth",
	CellLength:                   "CellLength",
	FillOrder:                    "FillOrder",
	DocumentName:                 "DocumentName",
	ImageDescription:             "ImageDescription",
	Make:                         "Make",
	Model:                        "Model",
	StripOffsets:                 "StripOffsets",
	Orientation:                  "Orientation",
	SamplesPerPixel:              "SamplesPerPixel",
	RowsPerStrip:                 "RowsPerStrip",
	StripByteCounts:              "StripByteCounts",
	MinSampleValue:               "MinSampleValue",
	MaxSampleValue:               "MaxSampleValue",
	XResolution:                  "XResolution",
	YResolution:                  "YResolution",
	PlanarConfiguration:          "PlanarConfiguration",
	PageName:                     "PageName",
	XPosition:                    "XPosition",
	YPosition:                    "YPosition",
	FreeOffsets:                  "FreeOffsets",
	FreeByteCounts:               "FreeByteCounts",
	GrayResponseUnit:             "GrayResponseUnit",
	GrayResponseCurve:            "GrayResponseCurve",
	T4Options:                    "T4Options",
	T6Options:                    "T6Options",
	ResolutionUnit:               "ResolutionUnit",
	PageNumber:                   "PageNumber",
	TransferFunction:             "TransferFunction",
	Software:                     "Software",
	DateTime:                     "DateTime",
	Artist:                       "Artist",
	HostComputer:                 "HostComputer",
	Predictor:                    "Predictor",
	WhitePoint:                   "WhitePoint",
	PrimaryChromaticities:        "PrimaryChromaticities",
	ColorMap:                     "ColorMap",
	HalftoneHints:                "HalftoneHints",
	TileWidth:                    "TileWidth",
	TileLength:                   "TileLength",
	TileOffsets:                  "TileOffsets",
	TileByteCounts:               "TileByteCounts",
	BadFaxLines:                  "BadFaxLines",
	CleanFaxData:                 "CleanFaxData",
	ConsecutiveBadFaxLines:       "ConsecutiveBadFaxLines",
	SubIFDs:                      "SubIFDs",
	InkSet:                       "InkSet",
	InkNames:                     "InkNames",
	NumberOfInks:                 "NumberOfInks",
	DotRange:                     "DotRange",
	TargetPrinter:                "TargetPrinter",
	ExtraSamples:                 "ExtraSamples",
	SampleFormat:                 "SampleFormat",
	SMinSampleValue:              "SMinSampleValue",
	SMaxSampleValue:              "SMaxSampleValue",
	TransferRange:                "TransferRange",
	ClipPath:                     "ClipPath",
	XClipPathUnits:               "XClipPathUnits",
	YClipPathUnits:               "YClipPathUnits",
	Indexed:                      "Indexed",
	JPEGTables:                   "JPEGTables",
	OPIProxy:                     "OPIProxy",
	GlobalParametersIFD:          "GlobalParametersIFD",
	ProfileType:                  "ProfileType",
	FaxProfile:                   "FaxProfile",
	CodingMethods:                "CodingMethods",
	VersionYear:                  "VersionYear",
	ModeNumber:                   "ModeNumber",
	Decode:                       "Decode",
	DefaultImageColor:            "DefaultImageColor",
	JPEGProc:                     "JPEGProc",
	JPEGInterchangeFormat:        "JPEGInterchangeFormat",
	JPEGInterchangeFormatLength:  "JPEGInterchangeFormatLength",
	JPEGRestartInterval:          "JPEGRestartInterval",
	JPEGLosslessPredictors:       "JPEGLosslessPredictors",
	JPEGPointTransforms:          "JPEGPointTransforms",
	JPEGQTables:                  "JPEGQTables",
	JPEGDCTables:                 "JPEGDCTables",
	JPEGACTables:                 "JPEGACTables",
	YCbCrCoefficients:            "YCbCrCoefficients",
	YCbCrSubSampling:             "YCbCrSubSampling",
	YCbCrPositioning:             "YCbCrPositioning",
	ReferenceBlackWhite:          "ReferenceBlackWhite",
	StripRowCounts:               "StripRowCounts",
	XMP:                          "XMP",
	ImageRating:                  "ImageRating",
	ImageRatingPercent:           "ImageRatingPercent",
	ImageID:                      "ImageID",
	WangAnnotation:               "WangAnnotation",
	CFARepeatPatternDim:          "CFARepeatPatternDim",
	CFAPattern:                   "CFAPattern",
	BatteryLevel:                 "BatteryLevel",
	Copyright:                    "Copyright",
	ExposureTime:                 "ExposureTime",
	FNumber:                      "FNumber",
	MDFileTag:                    "MDFileTag",
	MDScalePixel:                 "MDScalePixel",
	MDColorTable:                 "MDColorTable",
	MDLabName:                    "MDLabName",
	MDSampleInfo:                 "MDSampleInfo",
	MDPrepDate:                   "MDPrepDate",
	MDPrepTime:                   "MDPrepTime",
	MDFileUnits:                  "MDFileUnits",
	ModelPixelScaleTag:           "ModelPixelScaleTag",
	IPTC_NAA:                     "IPTC_NAA",
	INGRPacketDataTag:            "INGRPacketDataTag",
	INGRFlagRegisters:            "INGRFlagRegisters",
	IrasBTransformationMatrix:    "IrasBTransformationMatrix",
	ModelTiepointTag:             "ModelTiepointTag",
	Site:                         "Site",
	ColorSequence:                "ColorSequence",
	IT8Header:                    "IT8Header",
	RasterPadding:                "RasterPadding",
	BitsPerRunLength:             "BitsPerRunLength",
	BitsPerExtendedRunLength:     "BitsPerExtendedRunLength",
	ColorTable:                   "ColorTable",
	ImageColorIndicator:          "ImageColorIndicator",
	BackgroundColorIndicator:     "BackgroundColorIndicator",
	ImageColorValue:              "ImageColorValue",
	BackgroundColorValue:         "BackgroundColorValue",
	PixelIntensityRange:          "PixelIntensityRange",
	TransparencyIndicator:        "TransparencyIndicator",
	ColorCharacterization:        "ColorCharacterization",
	HCUsage:                      "HCUsage",
	TrapIndicator:                "TrapIndicator",
	CMYKEquivalent:               "CMYKEquivalent",
	Reserved0:                    "Reserved0",
	Reserved1:                    "Reserved1",
	Reserved2:                    "Reserved2",
	ModelTransformationTag:       "ModelTransformationTag",
	Photoshop:                    "Photoshop",
	ExifIFD:                      "ExifIFD",
	InterColorProfile:            "InterColorProfile",
	ImageLayer:                   "ImageLayer",
	GeoKeyDirectoryTag:           "GeoKeyDirectoryTag",
	GeoDoubleParamsTag:           "GeoDoubleParamsTag",
	GeoAsciiParamsTag:            "GeoAsciiParamsTag",
	ExposureProgram:              "ExposureProgram",
	SpectralSensitivity:          "SpectralSensitivity",
	GPSInfo:                      "GPSInfo",
	ISOSpeedRatings:              "ISOSpeedRatings",
	OECF:                         "OECF",
	Interlace:                    "Interlace",
	TimeZoneOffset:               "TimeZoneOffset",
	SelfTimeMode:                 "SelfTimeMode",
	SensitivityType:              "SensitivityType",
	StandardOutputSensitivity:    "StandardOutputSensitivity",
	RecommendedExposureIndex:     "RecommendedExposureIndex",
	ISOSpeed:                     "ISOSpeed",
	ISOSpeedLatitudeyyy:          "ISOSpeedLatitudeyyy",
	ISOSpeedLatitudezzz:          "ISOSpeedLatitudezzz",
	HylaFAXFaxRecvParams:         "HylaFAXFaxRecvParams",
	HylaFAXFaxSubAddress:         "HylaFAXFaxSubAddress",
	HylaFAXFaxRecvTime:           "HylaFAXFaxRecvTime",
	ExifVersion:                  "ExifVersion",
	DateTimeOriginal:             "DateTimeOriginal",
	DateTimeDigitized:            "DateTimeDigitized",
	ComponentsConfiguration:      "ComponentsConfiguration",
	CompressedBitsPerPixel:       "CompressedBitsPerPixel",
	ShutterSpeedValue:            "ShutterSpeedValue",
	ApertureValue:                "ApertureValue",
	BrightnessValue:              "BrightnessValue",
	ExposureBiasValue:            "ExposureBiasValue",
	MaxApertureValue:             "MaxApertureValue",
	SubjectDistance:              "SubjectDistance",
	MeteringMode:                 "MeteringMode",
	LightSource:                  "LightSource",
	Flash:                        "Flash",
	FocalLength:                  "FocalLength",
	FlashEnergy:                  "FlashEnergy",
	SpatialFrequencyResponse:     "SpatialFrequencyResponse",
	Noise:                        "Noise",
	FocalPlaneXResolution:        "FocalPlaneXResolution",
	FocalPlaneYResolution:        "FocalPlaneYResolution",
	FocalPlaneResolutionUnit:     "FocalPlaneResolutionUnit",
	ImageNumber:                  "ImageNumber",
	SecurityClassification:       "SecurityClassification",
	ImageHistory:                 "ImageHistory",
	SubjectLocation:              "SubjectLocation",
	ExposureIndex:                "ExposureIndex",
	TIFF_EPStandardID:            "TIFF_EPStandardID",
	SensingMethod:                "SensingMethod",
	MakerNote:                    "MakerNote",
	UserComment:                  "UserComment",
	SubsecTime:                   "SubsecTime",
	SubsecTimeOriginal:           "SubsecTimeOriginal",
	SubsecTimeDigitized:          "SubsecTimeDigitized",
	ImageSourceData:              "ImageSourceData",
	XPTitle:                      "XPTitle",
	XPComment:                    "XPComment",
	XPAuthor:                     "XPAuthor",
	XPKeywords:                   "XPKeywords",
	XPSubject:                    "XPSubject",
	FlashpixVersion:              "FlashpixVersion",
	ColorSpace:                   "ColorSpace",
	PixelXDimension:              "PixelXDimension",
	PixelYDimension:              "PixelYDimension",
	RelatedSoundFile:             "RelatedSoundFile",
	InteroperabilityIFD:          "InteroperabilityIFD",
	FlashEnergy2:                 "FlashEnergy2",
	SpatialFrequencyResponse2:    "SpatialFrequencyResponse2",
	FocalPlaneXResolution2:       "FocalPlaneXResolution2",
	FocalPlaneYResolution2:       "FocalPlaneYResolution2",
	FocalPlaneResolutionUnit2:    "FocalPlaneResolutionUnit2",
	SubjectLocation2:             "SubjectLocation2",
	ExposureIndex2:               "ExposureIndex2",
	SensingMethod2:               "SensingMethod2",
	FileSource:                   "FileSource",
	SceneType:                    "SceneType",
	CFAPattern2:                  "CFAPattern2",
	CustomRendered:               "CustomRendered",
	ExposureMode:                 "ExposureMode",
	WhiteBalance:                 "WhiteBalance",
	DigitalZoomRatio:             "DigitalZoomRatio",
	FocalLengthIn35mmFilm:        "FocalLengthIn35mmFilm",
	SceneCaptureType:             "SceneCaptureType",
	GainControl:                  "GainControl",
	Contrast:                     "Contrast",
	Saturation:                   "Saturation",
	Sharpness:                    "Sharpness",
	DeviceSettingDescription:     "DeviceSettingDescription",
	SubjectDistanceRange:         "SubjectDistanceRange",
	ImageUniqueID:                "ImageUniqueID",
	CameraOwnerName:              "CameraOwnerName",
	BodySerialNumber:             "BodySerialNumber",
	LensSpecification:            "LensSpecification",
	LensMake:                     "LensMake",
	LensModel:                    "LensModel",
	LensSerialNumber:             "LensSerialNumber",
	GDAL_METADATA:                "GDAL_METADATA",
	GDAL_NODATA:                  "GDAL_NODATA",
	HDPixelFormat:                "HDPixelFormat",
	HDTransformation:             "HDTransformation",
	HDUncompressed:               "HDUncompressed",
	HDImageType:                  "HDImageType",
	HDImageWidth:                 "HDImageWidth",
	HDImageHeight:                "HDImageHeight",
	HDWidthResolution:            "HDWidthResolution",
	HDHeightResolution:           "HDHeightResolution",
	HDImageOffset:                "HDImageOffset",
	HDImageByteCount:             "HDImageByteCount",
	HDAlphaOffset:                "HDAlphaOffset",
	HDAlphaByteCount:             "HDAlphaByteCount",
	HDImageDataDiscard:           "HDImageDataDiscard",
	HDAlphaDataDiscard:           "HDAlphaDataDiscard",
	OceScanjobDescription:        "OceScanjobDescription",
	OceApplicationSelector:       "OceApplicationSelector",
	OceIdentificationNumber:      "OceIdentificationNumber",
	OceImageLogicCharacteristics: "OceImageLogicCharacteristics",
	PrintImageMatching:           "PrintImageMatching",
	DNGVersion:                   "DNGVersion",
	DNGBackwardVersion:           "DNGBackwardVersion",
	UniqueCameraModel:            "UniqueCameraModel",
	LocalizedCameraModel:         "LocalizedCameraModel",
	CFAPlaneColor:                "CFAPlaneColor",
	CFALayout:                    "CFALayout",
	LinearizationTable:           "LinearizationTable",
	BlackLevelRepeatDim:          "BlackLevelRepeatDim",
	BlackLevel:                   "BlackLevel",
	BlackLevelDeltaH:             "BlackLevelDeltaH",
	BlackLevelDeltaV:             "BlackLevelDeltaV",
	WhiteLevel:                   "WhiteLevel",
	DefaultScale:                 "DefaultScale",
	DefaultCropOrigin:            "DefaultCropOrigin",
	DefaultCropSize:              "DefaultCropSize",
	ColorMatrix1:                 "ColorMatrix1",
	ColorMatrix2:                 "ColorMatrix2",
	CameraCalibration1:           "CameraCalibration1",
	CameraCalibration2:           "CameraCalibration2",
	ReductionMatrix1:             "ReductionMatrix1",
	ReductionMatrix2:             "ReductionMatrix2",
	AnalogBalance:                "AnalogBalance",
	AsShotNeutral:                "AsShotNeutral",
	AsShotWhiteXY:                "AsShotWhiteXY",
	BaselineExposure:             "BaselineExposure",
	BaselineNoise:                "BaselineNoise",
	BaselineSharpness:            "BaselineSharpness",
	BayerGreenSplit:              "BayerGreenSplit",
	LinearResponseLimit:          "LinearResponseLimit",
	CameraSerialNumber:           "CameraSerialNumber",
	LensInfo:                     "LensInfo",
	ChromaBlurRadius:             "ChromaBlurRadius",
	AntiAliasStrength:            "AntiAliasStrength",
	ShadowScale:                  "ShadowScale",
	DNGPrivateData:               "DNGPrivateData",
	MakerNoteSafety:              "MakerNoteSafety",
	CalibrationIlluminant1:       "CalibrationIlluminant1",
	CalibrationIlluminant2:       "CalibrationIlluminant2",
	BestQualityScale:             "BestQualityScale",
	RawDataUniqueID:              "RawDataUniqueID",
	AliasLayerMetadata:           "AliasLayerMetadata",
	OriginalRawFileName:          "OriginalRawFileName",
	OriginalRawFileData:          "OriginalRawFileData",
	ActiveArea:                   "ActiveArea",
	MaskedAreas:                  "MaskedAreas",
	AsShotICCProfile:             "AsShotICCProfile",
	AsShotPreProfileMatrix:       "AsShotPreProfileMatrix",
	CurrentICCProfile:            "CurrentICCProfile",
	CurrentPreProfileMatrix:      "CurrentPreProfileMatrix",
	ColorimetricReference:        "ColorimetricReference",
	CameraCalibrationSignature:   "CameraCalibrationSignature",
	ProfileCalibrationSignature:  "ProfileCalibrationSignature",
	ExtraCameraProfiles:          "ExtraCameraProfiles",
	AsShotProfileName:            "AsShotProfileName",
	NoiseReductionApplied:        "NoiseReductionApplied",
	ProfileName:                  "ProfileName",
	ProfileHueSatMapDims:         "ProfileHueSatMapDims",
	ProfileHueSatMapData1:        "ProfileHueSatMapData1",
	ProfileHueSatMapData2:        "ProfileHueSatMapData2",
	ProfileToneCurve:             "ProfileToneCurve",
	ProfileEmbedPolicy:           "ProfileEmbedPolicy",
	ProfileCopyright:             "ProfileCopyright",
	ForwardMatrix1:               "ForwardMatrix1",
	ForwardMatrix2:               "ForwardMatrix2",
	PreviewApplicationName:       "PreviewApplicationName",
	PreviewApplicationVersion:    "PreviewApplicationVersion",
	PreviewSettingsName:          "PreviewSettingsName",
	PreviewSettingsDigest:        "PreviewSettingsDigest",
	PreviewColorSpace:            "PreviewColorSpace",
	PreviewDateTime:              "PreviewDateTime",
	RawImageDigest:               "RawImageDigest",
	OriginalRawFileDigest:        "OriginalRawFileDigest",
	SubTileBlockSize:             "SubTileBlockSize",
	RowInterleaveFactor:          "RowInterleaveFactor",
	ProfileLookTableDims:         "ProfileLookTableDims",
	ProfileLookTableData:         "ProfileLookTableData",
	OpcodeList1:                  "OpcodeList1",
	OpcodeList2:                  "OpcodeList2",
	OpcodeList3:                  "OpcodeList3",
	NoiseProfile:                 "NoiseProfile",
	OriginalDefaultFinalSize:     "OriginalDefaultFinalSize",
	OriginalBestQualityFinalSize: "OriginalBestQualityFinalSize",
	OriginalDefaultCropSize:      "OriginalDefaultCropSize",
	ProfileHueSatMapEncoding:     "ProfileHueSatMapEncoding",
	ProfileLookTableEncoding:     "ProfileLookTableEncoding",
	BaselineExposureOffset:       "BaselineExposureOffset",
	DefaultBlackRender:           "DefaultBlackRender",
	NewRawImageDigest:            "NewRawImageDigest",
	RawToPreviewGain:             "RawToPreviewGain",
	DefaultUserCrop:              "DefaultUserCrop",
}

// FileDecoder is a TIFF decoder
type FileDecoder struct {
	decode.Common
}

// Decode TIFF file
func (d *FileDecoder) Decode() {
	switch d.PeekBits(32) {
	case littleEndian, bigEndian:
	default:
		d.Invalid("unknown endian")
	}
	var fu16 func(name string) uint64
	var fu32 func(name string) uint64
	var su32 func(name string) int64
	var u16 func() uint64

	endian := d.FieldUFn("endian", func() (uint64, decode.DisplayFormat, string) {
		endian := d.U32()
		d.SeekRel(-4 * 8)
		d.FieldUTF8("order", 2)
		// TODO: validate?
		d.FieldU16("integer_42")
		switch endian {
		case littleEndian:
			fu16 = d.FieldU16LE
			fu32 = d.FieldU32LE
			su32 = d.FieldS32LE
			u16 = d.U16LE
			fu32 = d.FieldU32LE
			return endian, decode.NumberHex, "little-endian"
		case bigEndian:
			fu16 = d.FieldU16BE
			fu32 = d.FieldU32BE
			su32 = d.FieldS32BE
			u16 = d.U16BE
			fu32 = d.FieldU32BE
			return endian, decode.NumberHex, "big-endian"
		}
		return endian, decode.NumberDecimal, "unknown"
	})

	ifdOffset := fu32("ifd_offset")

	d.MultiField("ifd", func() {
		// TODO: inf loop?
		for ifdOffset != 0 {
			d.SeekAbs(int64(ifdOffset) * 8)

			numberOfFields := fu16("number_of_field")
			for i := uint64(0); i < numberOfFields; i++ {
				d.Fields("ifd", func() {
					tag, _ := d.FieldStringMapFn("tag", tagNames, "unknown", u16)
					typ, typOk := d.FieldStringMapFn("type", typeNames, "unknown", u16)
					count := fu32("count")
					// TODO: short values stored in valueOffset directly?
					valueByteOffset := fu32("value_offset")

					if !typOk {
						return
					}

					valueByteSize := typeByteSize[typ] * count
					if valueByteSize <= 4 {
						// if value fits in offset itself use offset to value_offset
						valueByteOffset = uint64(d.Pos()/8) - 4
					}

					d.MultiField("values", func() {
						switch {
						case typ == UNDEFINED:
							switch tag {
							case InterColorProfile:
								d.FieldDecodeRange("icc", int64(valueByteOffset)*8, int64(valueByteSize)*8, iccTag...)
							default:
								// log.Printf("tag: %#+v\n", tag)
								// log.Printf("valueByteSize: %#+v\n", valueByteSize)
								d.FieldBitBufRange("value", int64(valueByteOffset)*8, int64(valueByteSize)*8)
							}
						case typ == ASCII:
							d.SubRangeFn(int64(valueByteOffset*8), int64(valueByteSize*8), func() {
								d.FieldUTF8("value", int64(valueByteSize))
							})
						case typ == BYTE:
							d.FieldBitBufRange("value", int64(valueByteOffset*8), int64(valueByteSize*8))
						default:
							// log.Printf("valueOffset: %d\n", valueByteOffset)
							// log.Printf("valueSize: %d\n", valueByteSize)
							d.SubRangeFn(int64(valueByteOffset*8), int64(valueByteSize*8), func() {
								for i := uint64(0); i < count; i++ {
									switch typ {
									// TODO: only some typ?
									// case BYTE:
									// 	d.FieldU8("value")
									case SHORT:
										fu16("value")
									case LONG:
										fu32("value")
									case RATIONAL:
										// TODO: endian? correct? unsigned 32:32 fixed point
										d.FieldUFP64("value")
									case SLONG:
										su32("value")
									case SRATIONAL:
										// TODO: endian? correct? signed 32:32 fixed point
										d.FieldFP64("value")
									default:
										panic("unknown type")
									}
								}
							})
						}
					})
				})
			}

			ifdOffset = fu32("next_ifd")
		}
	})

	_ = endian

}
