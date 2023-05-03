package tiff

import "github.com/wader/fq/pkg/scalar"

const (
	NewSubfileType               = 0x00fe
	SubfileType                  = 0x00ff
	ImageWidth                   = 0x0100
	ImageLength                  = 0x0101
	BitsPerSample                = 0x0102
	Compression                  = 0x0103
	PhotometricInterpretation    = 0x0106
	Threshholding                = 0x0107
	CellWidth                    = 0x0108
	CellLength                   = 0x0109
	FillOrder                    = 0x010a
	DocumentName                 = 0x010d
	ImageDescription             = 0x010e
	Make                         = 0x010f
	Model                        = 0x0110
	StripOffsets                 = 0x0111
	Orientation                  = 0x0112
	SamplesPerPixel              = 0x0115
	RowsPerStrip                 = 0x0116
	StripByteCounts              = 0x0117
	MinSampleValue               = 0x0118
	MaxSampleValue               = 0x0119
	XResolution                  = 0x011a
	YResolution                  = 0x011b
	PlanarConfiguration          = 0x011c
	PageName                     = 0x011d
	XPosition                    = 0x011e
	YPosition                    = 0x011f
	FreeOffsets                  = 0x0120
	FreeByteCounts               = 0x0121
	GrayResponseUnit             = 0x0122
	GrayResponseCurve            = 0x0123
	T4Options                    = 0x0124
	T6Options                    = 0x0125
	ResolutionUnit               = 0x0128
	PageNumber                   = 0x0129
	TransferFunction             = 0x012d
	Software                     = 0x0131
	DateTime                     = 0x0132
	Artist                       = 0x013b
	HostComputer                 = 0x013c
	Predictor                    = 0x013d
	WhitePoint                   = 0x013e
	PrimaryChromaticities        = 0x013f
	ColorMap                     = 0x0140
	HalftoneHints                = 0x0141
	TileWidth                    = 0x0142
	TileLength                   = 0x0143
	TileOffsets                  = 0x0144
	TileByteCounts               = 0x0145
	BadFaxLines                  = 0x0146
	CleanFaxData                 = 0x0147
	ConsecutiveBadFaxLines       = 0x0148
	SubIFDs                      = 0x014a
	InkSet                       = 0x014c
	InkNames                     = 0x014d
	NumberOfInks                 = 0x014e
	DotRange                     = 0x0150
	TargetPrinter                = 0x0151
	ExtraSamples                 = 0x0152
	SampleFormat                 = 0x0153
	SMinSampleValue              = 0x0154
	SMaxSampleValue              = 0x0155
	TransferRange                = 0x0156
	ClipPath                     = 0x0157
	XClipPathUnits               = 0x0158
	YClipPathUnits               = 0x0159
	Indexed                      = 0x015a
	JPEGTables                   = 0x015b
	OPIProxy                     = 0x015f
	GlobalParametersIFD          = 0x0190
	ProfileType                  = 0x0191
	FaxProfile                   = 0x0192
	CodingMethods                = 0x0193
	VersionYear                  = 0x0194
	ModeNumber                   = 0x0195
	Decode                       = 0x01b1
	DefaultImageColor            = 0x01b2
	JPEGProc                     = 0x0200
	JPEGInterchangeFormat        = 0x0201
	JPEGInterchangeFormatLength  = 0x0202
	JPEGRestartInterval          = 0x0203
	JPEGLosslessPredictors       = 0x0205
	JPEGPointTransforms          = 0x0206
	JPEGQTables                  = 0x0207
	JPEGDCTables                 = 0x0208
	JPEGACTables                 = 0x0209
	YCbCrCoefficients            = 0x0211
	YCbCrSubSampling             = 0x0212
	YCbCrPositioning             = 0x0213
	ReferenceBlackWhite          = 0x0214
	StripRowCounts               = 0x022f
	XMP                          = 0x02bc
	ImageRating                  = 0x4746
	ImageRatingPercent           = 0x4749
	ImageID                      = 0x800d
	WangAnnotation               = 0x80a4
	CFARepeatPatternDim          = 0x828d
	CFAPattern                   = 0x828e
	BatteryLevel                 = 0x828f
	Copyright                    = 0x8298
	ExposureTime                 = 0x829a
	FNumber                      = 0x829d
	MDFileTag                    = 0x82a5
	MDScalePixel                 = 0x82a6
	MDColorTable                 = 0x82a7
	MDLabName                    = 0x82a8
	MDSampleInfo                 = 0x82a9
	MDPrepDate                   = 0x82aa
	MDPrepTime                   = 0x82ab
	MDFileUnits                  = 0x82ac
	ModelPixelScaleTag           = 0x830e
	IPTC_NAA                     = 0x83bb
	INGRPacketDataTag            = 0x847e
	INGRFlagRegisters            = 0x847f
	IrasBTransformationMatrix    = 0x8480
	ModelTiepointTag             = 0x8482
	Site                         = 0x84e0
	ColorSequence                = 0x84e1
	IT8Header                    = 0x84e2
	RasterPadding                = 0x84e3
	BitsPerRunLength             = 0x84e4
	BitsPerExtendedRunLength     = 0x84e5
	ColorTable                   = 0x84e6
	ImageColorIndicator          = 0x84e7
	BackgroundColorIndicator     = 0x84e8
	ImageColorValue              = 0x84e9
	BackgroundColorValue         = 0x84ea
	PixelIntensityRange          = 0x84eb
	TransparencyIndicator        = 0x84ec
	ColorCharacterization        = 0x84ed
	HCUsage                      = 0x84ee
	TrapIndicator                = 0x84ef
	CMYKEquivalent               = 0x84f0
	Reserved0                    = 0x84f1
	Reserved1                    = 0x84f2
	Reserved2                    = 0x84f3
	ModelTransformationTag       = 0x85d8
	Photoshop                    = 0x8649
	ExifIFD                      = 0x8769
	InterColorProfile            = 0x8773
	ImageLayer                   = 0x87ac
	GeoKeyDirectoryTag           = 0x87af
	GeoDoubleParamsTag           = 0x87b0
	GeoAsciiParamsTag            = 0x87b1
	ExposureProgram              = 0x8822
	SpectralSensitivity          = 0x8824
	GPSInfo                      = 0x8825
	ISOSpeedRatings              = 0x8827
	OECF                         = 0x8828
	Interlace                    = 0x8829
	TimeZoneOffset               = 0x882a
	SelfTimeMode                 = 0x882b
	SensitivityType              = 0x8830
	StandardOutputSensitivity    = 0x8831
	RecommendedExposureIndex     = 0x8832
	ISOSpeed                     = 0x8833
	ISOSpeedLatitudeyyy          = 0x8834
	ISOSpeedLatitudezzz          = 0x8835
	HylaFAXFaxRecvParams         = 0x885c
	HylaFAXFaxSubAddress         = 0x885d
	HylaFAXFaxRecvTime           = 0x885e
	ExifVersion                  = 0x9000
	DateTimeOriginal             = 0x9003
	DateTimeDigitized            = 0x9004
	ComponentsConfiguration      = 0x9101
	CompressedBitsPerPixel       = 0x9102
	ShutterSpeedValue            = 0x9201
	ApertureValue                = 0x9202
	BrightnessValue              = 0x9203
	ExposureBiasValue            = 0x9204
	MaxApertureValue             = 0x9205
	SubjectDistance              = 0x9206
	MeteringMode                 = 0x9207
	LightSource                  = 0x9208
	Flash                        = 0x9209
	FocalLength                  = 0x920a
	FlashEnergy                  = 0x920b
	SpatialFrequencyResponse     = 0x920c
	Noise                        = 0x920d
	FocalPlaneXResolution        = 0x920e
	FocalPlaneYResolution        = 0x920f
	FocalPlaneResolutionUnit     = 0x9210
	ImageNumber                  = 0x9211
	SecurityClassification       = 0x9212
	ImageHistory                 = 0x9213
	SubjectLocation              = 0x9214
	ExposureIndex                = 0x9215
	TIFF_EPStandardID            = 0x9216
	SensingMethod                = 0x9217
	MakerNote                    = 0x927c
	UserComment                  = 0x9286
	SubsecTime                   = 0x9290
	SubsecTimeOriginal           = 0x9291
	SubsecTimeDigitized          = 0x9292
	ImageSourceData              = 0x935c
	XPTitle                      = 0x9c9b
	XPComment                    = 0x9c9c
	XPAuthor                     = 0x9c9d
	XPKeywords                   = 0x9c9e
	XPSubject                    = 0x9c9f
	FlashpixVersion              = 0xa000
	ColorSpace                   = 0xa001
	PixelXDimension              = 0xa002
	PixelYDimension              = 0xa003
	RelatedSoundFile             = 0xa004
	InteroperabilityIFD          = 0xa005
	FlashEnergy2                 = 0xa20b
	SpatialFrequencyResponse2    = 0xa20c
	FocalPlaneXResolution2       = 0xa20e
	FocalPlaneYResolution2       = 0xa20f
	FocalPlaneResolutionUnit2    = 0xa210
	SubjectLocation2             = 0xa214
	ExposureIndex2               = 0xa215
	SensingMethod2               = 0xa217
	FileSource                   = 0xa300
	SceneType                    = 0xa301
	CFAPattern2                  = 0xa302
	CustomRendered               = 0xa401
	ExposureMode                 = 0xa402
	WhiteBalance                 = 0xa403
	DigitalZoomRatio             = 0xa404
	FocalLengthIn35mmFilm        = 0xa405
	SceneCaptureType             = 0xa406
	GainControl                  = 0xa407
	Contrast                     = 0xa408
	Saturation                   = 0xa409
	Sharpness                    = 0xa40a
	DeviceSettingDescription     = 0xa40b
	SubjectDistanceRange         = 0xa40c
	ImageUniqueID                = 0xa420
	CameraOwnerName              = 0xa430
	BodySerialNumber             = 0xa431
	LensSpecification            = 0xa432
	LensMake                     = 0xa433
	LensModel                    = 0xa434
	LensSerialNumber             = 0xa435
	GDAL_METADATA                = 0xa480
	GDAL_NODATA                  = 0xa481
	HDPixelFormat                = 0xbc01
	HDTransformation             = 0xbc02
	HDUncompressed               = 0xbc03
	HDImageType                  = 0xbc04
	HDImageWidth                 = 0xbc80
	HDImageHeight                = 0xbc81
	HDWidthResolution            = 0xbc82
	HDHeightResolution           = 0xbc83
	HDImageOffset                = 0xbcc0
	HDImageByteCount             = 0xbcc1
	HDAlphaOffset                = 0xbcc2
	HDAlphaByteCount             = 0xbcc3
	HDImageDataDiscard           = 0xbcc4
	HDAlphaDataDiscard           = 0xbcc5
	OceScanjobDescription        = 0xc427
	OceApplicationSelector       = 0xc428
	OceIdentificationNumber      = 0xc429
	OceImageLogicCharacteristics = 0xc42a
	PrintImageMatching           = 0xc4a5
	DNGVersion                   = 0xc612
	DNGBackwardVersion           = 0xc613
	UniqueCameraModel            = 0xc614
	LocalizedCameraModel         = 0xc615
	CFAPlaneColor                = 0xc616
	CFALayout                    = 0xc617
	LinearizationTable           = 0xc618
	BlackLevelRepeatDim          = 0xc619
	BlackLevel                   = 0xc61a
	BlackLevelDeltaH             = 0xc61b
	BlackLevelDeltaV             = 0xc61c
	WhiteLevel                   = 0xc61d
	DefaultScale                 = 0xc61e
	DefaultCropOrigin            = 0xc61f
	DefaultCropSize              = 0xc620
	ColorMatrix1                 = 0xc621
	ColorMatrix2                 = 0xc622
	CameraCalibration1           = 0xc623
	CameraCalibration2           = 0xc624
	ReductionMatrix1             = 0xc625
	ReductionMatrix2             = 0xc626
	AnalogBalance                = 0xc627
	AsShotNeutral                = 0xc628
	AsShotWhiteXY                = 0xc629
	BaselineExposure             = 0xc62a
	BaselineNoise                = 0xc62b
	BaselineSharpness            = 0xc62c
	BayerGreenSplit              = 0xc62d
	LinearResponseLimit          = 0xc62e
	CameraSerialNumber           = 0xc62f
	LensInfo                     = 0xc630
	ChromaBlurRadius             = 0xc631
	AntiAliasStrength            = 0xc632
	ShadowScale                  = 0xc633
	DNGPrivateData               = 0xc634
	MakerNoteSafety              = 0xc635
	CalibrationIlluminant1       = 0xc65a
	CalibrationIlluminant2       = 0xc65b
	BestQualityScale             = 0xc65c
	RawDataUniqueID              = 0xc65d
	AliasLayerMetadata           = 0xc660
	OriginalRawFileName          = 0xc68b
	OriginalRawFileData          = 0xc68c
	ActiveArea                   = 0xc68d
	MaskedAreas                  = 0xc68e
	AsShotICCProfile             = 0xc68f
	AsShotPreProfileMatrix       = 0xc690
	CurrentICCProfile            = 0xc691
	CurrentPreProfileMatrix      = 0xc692
	ColorimetricReference        = 0xc6bf
	CameraCalibrationSignature   = 0xc6f3
	ProfileCalibrationSignature  = 0xc6f4
	ExtraCameraProfiles          = 0xc6f5
	AsShotProfileName            = 0xc6f6
	NoiseReductionApplied        = 0xc6f7
	ProfileName                  = 0xc6f8
	ProfileHueSatMapDims         = 0xc6f9
	ProfileHueSatMapData1        = 0xc6fa
	ProfileHueSatMapData2        = 0xc6fb
	ProfileToneCurve             = 0xc6fc
	ProfileEmbedPolicy           = 0xc6fd
	ProfileCopyright             = 0xc6fe
	ForwardMatrix1               = 0xc714
	ForwardMatrix2               = 0xc715
	PreviewApplicationName       = 0xc716
	PreviewApplicationVersion    = 0xc717
	PreviewSettingsName          = 0xc718
	PreviewSettingsDigest        = 0xc719
	PreviewColorSpace            = 0xc71a
	PreviewDateTime              = 0xc71b
	RawImageDigest               = 0xc71c
	OriginalRawFileDigest        = 0xc71d
	SubTileBlockSize             = 0xc71e
	RowInterleaveFactor          = 0xc71f
	ProfileLookTableDims         = 0xc725
	ProfileLookTableData         = 0xc726
	OpcodeList1                  = 0xc740
	OpcodeList2                  = 0xc741
	OpcodeList3                  = 0xc74e
	NoiseProfile                 = 0xc761
	OriginalDefaultFinalSize     = 0xc791
	OriginalBestQualityFinalSize = 0xc792
	OriginalDefaultCropSize      = 0xc793
	ProfileHueSatMapEncoding     = 0xc7a3
	ProfileLookTableEncoding     = 0xc7a4
	BaselineExposureOffset       = 0xc7a5
	DefaultBlackRender           = 0xc7a6
	NewRawImageDigest            = 0xc7a7
	RawToPreviewGain             = 0xc7a8
	DefaultUserCrop              = 0xc7b5

	// TODO: where are these documented?
	FrameDelay       = 0x5100
	LoopCount        = 0x5101
	GlobalPalette    = 0x5102
	IndexBackground  = 0x5103
	IndexTransparent = 0x5104
	PixelUnit        = 0x5110
	PixelPerUnitX    = 0x5111
	PixelPerUnitY    = 0x5112
	PaletteHistogram = 0x5113
)

var tiffTagNames = scalar.UintMapSymStr{
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

	FrameDelay:       "FrameDelay",
	LoopCount:        "LoopCount",
	GlobalPalette:    "GlobalPalette",
	IndexBackground:  "IndexBackground",
	IndexTransparent: "IndexTransparent",
	PixelUnit:        "PixelUnit",
	PixelPerUnitX:    "PixelPerUnitX",
	PixelPerUnitY:    "PixelPerUnitY",
	PaletteHistogram: "PaletteHistogram",
}

const (
	GPSVersionID         = 0x0000
	GPSLatitudeRef       = 0x0001
	GPSLatitude          = 0x0002
	GPSLongitudeRef      = 0x0003
	GPSLongitude         = 0x0004
	GPSAltitudeRef       = 0x0005
	GPSAltitude          = 0x0006
	GPSTimeStamp         = 0x0007
	GPSSatellites        = 0x0008
	GPSStatus            = 0x0009
	GPSMeasureMode       = 0x000a
	GPSDOP               = 0x000b
	GPSSpeedRef          = 0x000c
	GPSSpeed             = 0x000d
	GPSTrackRef          = 0x000e
	GPSTrack             = 0x000f
	GPSImgDirectionRef   = 0x0010
	GPSImgDirection      = 0x0011
	GPSMapDatum          = 0x0012
	GPSDestLatitudeRef   = 0x0013
	GPSDestLatitude      = 0x0014
	GPSDestLongitudeRef  = 0x0015
	GPSDestLongitude     = 0x0016
	GPSDestBearingRef    = 0x0017
	GPSDestBearing       = 0x0018
	GPSDestDistanceRef   = 0x0019
	GPSDestDistance      = 0x001a
	GPSProcessingMethod  = 0x001b
	GPSAreaInformation   = 0x001c
	GPSDateStamp         = 0x001d
	GPSDifferential      = 0x001e
	GPSHPositioningError = 0x001f
)

var gpsInfoTagNames = scalar.UintMapSymStr{
	GPSVersionID:         "GPSVersionID",
	GPSLatitudeRef:       "GPSLatitudeRef",
	GPSLatitude:          "GPSLatitude",
	GPSLongitudeRef:      "GPSLongitudeRef",
	GPSLongitude:         "GPSLongitude",
	GPSAltitudeRef:       "GPSAltitudeRef",
	GPSAltitude:          "GPSAltitude",
	GPSTimeStamp:         "GPSTimeStamp",
	GPSSatellites:        "GPSSatellites",
	GPSStatus:            "GPSStatus",
	GPSMeasureMode:       "GPSMeasureMode",
	GPSDOP:               "GPSDOP",
	GPSSpeedRef:          "GPSSpeedRef",
	GPSSpeed:             "GPSSpeed",
	GPSTrackRef:          "GPSTrackRef",
	GPSTrack:             "GPSTrack",
	GPSImgDirectionRef:   "GPSImgDirectionRef",
	GPSImgDirection:      "GPSImgDirection",
	GPSMapDatum:          "GPSMapDatum",
	GPSDestLatitudeRef:   "GPSDestLatitudeRef",
	GPSDestLatitude:      "GPSDestLatitude",
	GPSDestLongitudeRef:  "GPSDestLongitudeRef",
	GPSDestLongitude:     "GPSDestLongitude",
	GPSDestBearingRef:    "GPSDestBearingRef",
	GPSDestBearing:       "GPSDestBearing",
	GPSDestDistanceRef:   "GPSDestDistanceRef",
	GPSDestDistance:      "GPSDestDistance",
	GPSProcessingMethod:  "GPSProcessingMethod",
	GPSAreaInformation:   "GPSAreaInformation",
	GPSDateStamp:         "GPSDateStamp",
	GPSDifferential:      "GPSDifferential",
	GPSHPositioningError: "GPSHPositioningError",
}
