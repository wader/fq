package moc3

// https://github.com/OpenL2D/moc3ingbird/blob/master/src/moc3.hexpat

import (
	"embed"

	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
	"github.com/wader/fq/pkg/scalar"
)

//go:embed moc3.md
var moc3FS embed.FS

func init() {
	interp.RegisterFormat(
		format.MOC3,
		&decode.Format{
			Description: "MOC3 file",
			Groups:      []*decode.Group{format.Probe},
			DecodeFn:    decodeMOC3,
		})
	interp.RegisterFS(moc3FS)
}

const (
	moc3AlignBytes = 64
	moc3AlignBits  = moc3AlignBytes * 8
)

const (
	moc3Version3_00_00 = 1
	moc3Version3_03_00 = 2
	moc3Version4_00_00 = 3
	moc3Version4_02_00 = 4
	moc3Version5_00_00 = 5
)

var boolNToSym = scalar.UintMapSymBool{
	0: false,
	1: true,
}

var moc3VersionNames = scalar.UintMap{
	moc3Version3_00_00: {Sym: "V3_00_00", Description: "3.0.00 - 3.2.07"},
	moc3Version3_03_00: {Sym: "V3_03_00", Description: "3.3.00 - 3.3.03"},
	moc3Version4_00_00: {Sym: "V4_00_00", Description: "4.0.00 - 4.1.05"},
	moc3Version4_02_00: {Sym: "V4_02_00", Description: "4.2.00 - 4.2.02"},
	moc3Version5_00_00: {Sym: "V5_00_00", Description: "5.0.00"},
}

var deformerTypeNames = scalar.UintMapSymStr{
	0: "warp",
	1: "rotation",
}

var blendModeNames = scalar.UintMapSymStr{
	0: "normal",
	1: "additive",
	2: "multiplicative",
}

var drawOrderGroupObjectTypeNames = scalar.UintMapSymStr{
	0: "art_mesh",
	1: "part",
}

var parameterTypeNames = scalar.UintMapSymStr{
	0: "normal",
	1: "blend_shape",
}

type countInfoTable struct {
	parts                        int64
	deformers                    int64
	warpDeformers                int64
	rotationDeformers            int64
	artMeshes                    int64
	parameters                   int64
	partKeyforms                 int64
	warpDeformerKeyforms         int64
	rotationDeformerKeyforms     int64
	artMeshKeyforms              int64
	keyformPositions             int64
	parameterBindingIndices      int64
	keyformBindings              int64
	parameterBindings            int64
	keys                         int64
	uvs                          int64
	positionIndices              int64
	drawableMasks                int64
	drawOrderGroups              int64
	drawOrderGroupObjects        int64
	glue                         int64
	glueInfo                     int64
	glueKeyforms                 int64
	keyformMultiplyColors        int64
	keyformScreenColors          int64
	blendShapeParameterBindings  int64
	blendShapeKeyformBindings    int64
	blendShapesWarpDeformers     int64
	blendShapesArtMeshes         int64
	blendShapeConstraintIndices  int64
	blendShapeConstraints        int64
	blendShapeConstraintValues   int64
	blendShapesParts             int64
	blendShapesRotationDeformers int64
	blendShapesGlue              int64
}

type sectionOffsetTable struct {
	countInfo  int64
	canvasInfo int64

	parts struct {
		runtimeSpace0                int64
		ids                          int64
		keyformBindingSourcesIndices int64
		keyformSourcesBeginIndices   int64
		keyformSourcesCounts         int64
		isVisible                    int64
		isEnabled                    int64
		parentPartIndices            int64
	}

	deformers struct {
		runtimeSpace0                int64
		ids                          int64
		keyformBindingSourcesIndices int64
		isVisible                    int64
		isEnabled                    int64
		parentPartIndices            int64
		parentDeformerIndices        int64
		types                        int64
		specificSourcesIndices       int64
	}

	warpDeformers struct {
		keyformBindingSourcesIndices int64
		keyformSourcesBeginIndices   int64
		keyformSourcesCounts         int64
		vertexCounts                 int64
		rows                         int64
		columns                      int64
	}

	rotationDeformers struct {
		keyformBindingSourcesIndices int64
		keyformSourcesBeginIndices   int64
		keyformSourcesCounts         int64
		baseAngles                   int64
	}

	artMeshes struct {
		runtimeSpace0                    int64
		runtimeSpace1                    int64
		runtimeSpace2                    int64
		runtimeSpace3                    int64
		ids                              int64
		keyformBindingSourcesIndices     int64
		keyformSourcesBeginIndices       int64
		keyformSourcesCounts             int64
		isVisible                        int64
		isEnabled                        int64
		parentPartIndices                int64
		parentDeformerIndices            int64
		textureNos                       int64
		drawableFlags                    int64
		vertexCounts                     int64
		uvSourcesBeginIndices            int64
		positionIndexSourcesBeginIndices int64
		positionIndexSourcesCounts       int64
		drawableMaskSourcesBeginIndices  int64
		drawableMaskSourcesCounts        int64
	}

	parameters struct {
		runtimeSpace0                       int64
		ids                                 int64
		maxValues                           int64
		minValues                           int64
		defaultValues                       int64
		isRepeat                            int64
		decimalPlaces                       int64
		parameterBindingSourcesBeginIndices int64
		parameterBindingSourcesCounts       int64
	}

	partKeyforms struct {
		drawOrders int64
	}

	warpDeformerKeyforms struct {
		opacities                          int64
		keyformPositionSourcesBeginIndices int64
	}

	rotationDeformerKeyforms struct {
		opacities  int64
		angles     int64
		originX    int64
		originY    int64
		scales     int64
		isReflectX int64
		isReflectY int64
	}

	artMeshKeyforms struct {
		opacities                          int64
		drawOrders                         int64
		keyformPositionSourcesBeginIndices int64
	}

	keyformPositions struct {
		xys int64
	}

	parameterBindingIndices struct {
		bindingSourcesIndices int64
	}

	keyformBindings struct {
		parameterBindingIndexSourcesBeginIndices int64
		parameterBindingIndexSourcesCounts       int64
	}

	parameterBindings struct {
		keysSourcesBeginIndices int64
		keysSourcesCounts       int64
	}

	keys struct {
		values int64
	}

	UVs struct {
		uvs int64
	}

	positionIndices struct {
		indices int64
	}

	drawableMasks struct {
		artMeshSourcesIndices int64
	}

	drawOrderGroups struct {
		objectSourcesBeginIndices int64
		objectSourcesCounts       int64
		objectSourcesTotalCounts  int64
		maximumDrawOrders         int64
		minimumDrawOrders         int64
	}

	drawOrderGroupObjects struct {
		types       int64
		indices     int64
		selfIndices int64
	}

	glue struct {
		runtimeSpace0                int64
		ids                          int64
		keyformBindingSourcesIndices int64
		keyformSourcesBeginIndices   int64
		keyformSourcesCounts         int64
		artMeshIndicesA              int64
		artMeshIndicesB              int64
		glueInfoSourcesBeginIndices  int64
		glueInfoSourcesCounts        int64
	}

	glueInfo struct {
		weights         int64
		positionIndices int64
	}

	glueKeyforms struct {
		intensities int64
	}

	warpDeformersV3_3 struct {
		isQuadSource int64
	}

	parameterExtensions struct {
		runtimeSpace0           int64
		keysSourcesBeginIndices int64
		keysSourcesCounts       int64
	}

	warpDeformersV4_2 struct {
		keyformColorSourcesBeginIndices int64
	}

	rotationDeformersV4_2 struct {
		keyformColorSourcesBeginIndices int64
	}

	artMeshesV4_2 struct {
		keyformColorSourcesBeginIndices int64
	}

	keyformMultiplyColors struct {
		r int64
		g int64
		b int64
	}

	keyformScreenColors struct {
		r int64
		g int64
		b int64
	}

	parametersV4_2 struct {
		parameterTypes                                int64
		blendShapeParameterBindingSourcesBeginIndices int64
		blendShapeParameterBindingSourcesCounts       int64
	}

	blendShapeParameterBindings struct {
		keysSourcesBeginIndices int64
		keysSourcesCounts       int64
		baseKeyIndices          int64
	}

	blendShapeKeyformBindings struct {
		parameterBindingSourcesIndices               int64
		keyformSourcesBlendShapeIndices              int64
		keyformSourcesBlendShapeCounts               int64
		blendShapeConstraintIndexSourcesBeginIndices int64
		blendShapeConstraintIndexSourcesCounts       int64
	}

	blendShapesWarpDeformers struct {
		targetIndices                               int64
		blendShapeKeyformBindingSourcesBeginIndices int64
		blendShapeKeyformBindingSourcesCounts       int64
	}

	blendShapesArtMeshes struct {
		targetIndices                               int64
		blendShapeKeyformBindingSourcesBeginIndices int64
		blendShapeKeyformBindingSourcesCounts       int64
	}

	blendShapeConstraintIndices struct {
		blendShapeConstraintSourcesIndices int64
	}

	blendShapeConstraints struct {
		parameterIndices                             int64
		blendShapeConstraintValueSourcesBeginIndices int64
		blendShapeConstraintValueSourcesCounts       int64
	}

	blendShapeConstraintValues struct {
		keys    int64
		weights int64
	}

	warpDeformerKeyformsV5_0 struct {
		keyformMultiplyColorSourcesBeginIndices int64
		keyformScreenColorSourcesBeginIndices   int64
	}

	rotationDeformerKeyformsV5_0 struct {
		keyformMultiplyColorSourcesBeginIndices int64
		keyformScreenColorSourcesBeginIndices   int64
	}

	artMeshKeyformsV5_0 struct {
		keyformMultiplyColorSourcesBeginIndices int64
		keyformScreenColorSourcesBeginIndices   int64
	}

	blendShapesParts struct {
		targetIndices                               int64
		blendShapeKeyformBindingSourcesBeginIndices int64
		blendShapeKeyformBindingSourcesCounts       int64
	}

	blendShapesRotationDeformers struct {
		targetIndices                               int64
		blendShapeKeyformBindingSourcesBeginIndices int64
		blendShapeKeyformBindingSourcesCounts       int64
	}

	blendShapesGlue struct {
		targetIndices                               int64
		blendShapeKeyformBindingSourcesBeginIndices int64
		blendShapeKeyformBindingSourcesCounts       int64
	}
}

func decodeMOC3(d *decode.D) any {
	fieldAlignedNArray := func(d *decode.D, name string, n int64, fn func(d *decode.D)) {
		d.FieldStruct(name, func(d *decode.D) {
			d.FieldArray("array", func(d *decode.D) {
				for i := int64(0); i < n; i++ {
					fn(d)
				}
			})

			var padding int64
			if n != 0 {
				padding = int64(d.AlignBits(moc3AlignBits))
			}
			d.FieldRawLen("padding", padding)
		})
	}

	fieldRuntimeSpace := func(d *decode.D, name string, n int64, pad bool) {
		paddedSize := n * 64
		if pad && paddedSize%moc3AlignBits != 0 {
			paddedSize = (paddedSize + moc3AlignBits) / moc3AlignBits * moc3AlignBits
		}

		d.FieldRawLen(name, paddedSize)
	}

	d.FieldUTF8("magic", 4, d.StrAssert("MOC3"))
	version := d.FieldU8("version", moc3VersionNames)
	isBigEndian := d.FieldU8("is_big_endian", boolNToSym) != 0

	if !isBigEndian {
		d.Endian = decode.LittleEndian
	}
	d.FieldRawLen("unused0", 58*8)

	var sectionOffsets sectionOffsetTable
	d.FramedFn(0x280*8, func(d *decode.D) {
		d.FieldStruct("section_offsets", func(d *decode.D) {
			sectionOffsets.countInfo = int64(d.FieldU32("count_info", scalar.UintHex))
			sectionOffsets.canvasInfo = int64(d.FieldU32("canvas_info", scalar.UintHex))

			d.FieldStruct("parts", func(d *decode.D) {
				sectionOffsets.parts.runtimeSpace0 = int64(d.FieldU32("runtime_space0", scalar.UintHex))
				sectionOffsets.parts.ids = int64(d.FieldU32("ids", scalar.UintHex))
				sectionOffsets.parts.keyformBindingSourcesIndices = int64(d.FieldU32("keyform_binding_sources_indices", scalar.UintHex))
				sectionOffsets.parts.keyformSourcesBeginIndices = int64(d.FieldU32("keyform_sources_begin_indices", scalar.UintHex))
				sectionOffsets.parts.keyformSourcesCounts = int64(d.FieldU32("keyform_sources_counts", scalar.UintHex))
				sectionOffsets.parts.isVisible = int64(d.FieldU32("is_visible", scalar.UintHex))
				sectionOffsets.parts.isEnabled = int64(d.FieldU32("is_enabled", scalar.UintHex))
				sectionOffsets.parts.parentPartIndices = int64(d.FieldU32("parent_part_indices", scalar.UintHex))
			})

			d.FieldStruct("deformers", func(d *decode.D) {
				sectionOffsets.deformers.runtimeSpace0 = int64(d.FieldU32("runtime_space0", scalar.UintHex))
				sectionOffsets.deformers.ids = int64(d.FieldU32("ids", scalar.UintHex))
				sectionOffsets.deformers.keyformBindingSourcesIndices = int64(d.FieldU32("keyform_binding_sources_indices", scalar.UintHex))
				sectionOffsets.deformers.isVisible = int64(d.FieldU32("is_visible", scalar.UintHex))
				sectionOffsets.deformers.isEnabled = int64(d.FieldU32("is_enabled", scalar.UintHex))
				sectionOffsets.deformers.parentPartIndices = int64(d.FieldU32("parent_part_indices", scalar.UintHex))
				sectionOffsets.deformers.parentDeformerIndices = int64(d.FieldU32("parent_deformer_indices", scalar.UintHex))
				sectionOffsets.deformers.types = int64(d.FieldU32("types", scalar.UintHex))
				sectionOffsets.deformers.specificSourcesIndices = int64(d.FieldU32("specific_sources_indices", scalar.UintHex))
			})

			d.FieldStruct("warp_deformers", func(d *decode.D) {
				sectionOffsets.warpDeformers.keyformBindingSourcesIndices = int64(d.FieldU32("keyform_binding_sources_indices", scalar.UintHex))
				sectionOffsets.warpDeformers.keyformSourcesBeginIndices = int64(d.FieldU32("keyform_sources_begin_indices", scalar.UintHex))
				sectionOffsets.warpDeformers.keyformSourcesCounts = int64(d.FieldU32("keyform_sources_counts", scalar.UintHex))
				sectionOffsets.warpDeformers.vertexCounts = int64(d.FieldU32("vertex_counts", scalar.UintHex))
				sectionOffsets.warpDeformers.rows = int64(d.FieldU32("rows", scalar.UintHex))
				sectionOffsets.warpDeformers.columns = int64(d.FieldU32("columns", scalar.UintHex))
			})

			d.FieldStruct("rotation_deformers", func(d *decode.D) {
				sectionOffsets.rotationDeformers.keyformBindingSourcesIndices = int64(d.FieldU32("keyform_binding_sources_indices", scalar.UintHex))
				sectionOffsets.rotationDeformers.keyformSourcesBeginIndices = int64(d.FieldU32("keyform_sources_begin_indices", scalar.UintHex))
				sectionOffsets.rotationDeformers.keyformSourcesCounts = int64(d.FieldU32("keyform_sources_counts", scalar.UintHex))
				sectionOffsets.rotationDeformers.baseAngles = int64(d.FieldU32("base_angles", scalar.UintHex))
			})

			d.FieldStruct("art_meshes", func(d *decode.D) {
				sectionOffsets.artMeshes.runtimeSpace0 = int64(d.FieldU32("runtime_space0", scalar.UintHex))
				sectionOffsets.artMeshes.runtimeSpace1 = int64(d.FieldU32("runtime_space1", scalar.UintHex))
				sectionOffsets.artMeshes.runtimeSpace2 = int64(d.FieldU32("runtime_space2", scalar.UintHex))
				sectionOffsets.artMeshes.runtimeSpace3 = int64(d.FieldU32("runtime_space3", scalar.UintHex))
				sectionOffsets.artMeshes.ids = int64(d.FieldU32("ids", scalar.UintHex))
				sectionOffsets.artMeshes.keyformBindingSourcesIndices = int64(d.FieldU32("keyform_binding_sources_indices", scalar.UintHex))
				sectionOffsets.artMeshes.keyformSourcesBeginIndices = int64(d.FieldU32("keyform_sources_begin_indices", scalar.UintHex))
				sectionOffsets.artMeshes.keyformSourcesCounts = int64(d.FieldU32("keyform_sources_counts", scalar.UintHex))
				sectionOffsets.artMeshes.isVisible = int64(d.FieldU32("is_visible", scalar.UintHex))
				sectionOffsets.artMeshes.isEnabled = int64(d.FieldU32("is_enabled", scalar.UintHex))
				sectionOffsets.artMeshes.parentPartIndices = int64(d.FieldU32("parent_part_indices", scalar.UintHex))
				sectionOffsets.artMeshes.parentDeformerIndices = int64(d.FieldU32("parent_deformer_indices", scalar.UintHex))
				sectionOffsets.artMeshes.textureNos = int64(d.FieldU32("texture_nos", scalar.UintHex))
				sectionOffsets.artMeshes.drawableFlags = int64(d.FieldU32("drawable_flags", scalar.UintHex))
				sectionOffsets.artMeshes.vertexCounts = int64(d.FieldU32("vertex_counts", scalar.UintHex))
				sectionOffsets.artMeshes.uvSourcesBeginIndices = int64(d.FieldU32("uv_sources_begin_indices", scalar.UintHex))
				sectionOffsets.artMeshes.positionIndexSourcesBeginIndices = int64(d.FieldU32("position_index_sources_begin_indices", scalar.UintHex))
				sectionOffsets.artMeshes.positionIndexSourcesCounts = int64(d.FieldU32("position_index_sources_counts", scalar.UintHex))
				sectionOffsets.artMeshes.drawableMaskSourcesBeginIndices = int64(d.FieldU32("drawable_mask_sources_begin_indices", scalar.UintHex))
				sectionOffsets.artMeshes.drawableMaskSourcesCounts = int64(d.FieldU32("drawable_mask_sources_counts", scalar.UintHex))
			})

			d.FieldStruct("parameters", func(d *decode.D) {
				sectionOffsets.parameters.runtimeSpace0 = int64(d.FieldU32("runtime_space0", scalar.UintHex))
				sectionOffsets.parameters.ids = int64(d.FieldU32("ids", scalar.UintHex))
				sectionOffsets.parameters.maxValues = int64(d.FieldU32("max_values", scalar.UintHex))
				sectionOffsets.parameters.minValues = int64(d.FieldU32("min_values", scalar.UintHex))
				sectionOffsets.parameters.defaultValues = int64(d.FieldU32("default_values", scalar.UintHex))
				sectionOffsets.parameters.isRepeat = int64(d.FieldU32("is_repeat", scalar.UintHex))
				sectionOffsets.parameters.decimalPlaces = int64(d.FieldU32("decimal_places", scalar.UintHex))
				sectionOffsets.parameters.parameterBindingSourcesBeginIndices = int64(d.FieldU32("parameter_binding_sources_begin_indices", scalar.UintHex))
				sectionOffsets.parameters.parameterBindingSourcesCounts = int64(d.FieldU32("parameter_binding_sources_counts", scalar.UintHex))
			})

			d.FieldStruct("part_keyforms", func(d *decode.D) {
				sectionOffsets.partKeyforms.drawOrders = int64(d.FieldU32("draw_orders", scalar.UintHex))
			})

			d.FieldStruct("warp_deformer_keyforms", func(d *decode.D) {
				sectionOffsets.warpDeformerKeyforms.opacities = int64(d.FieldU32("opacities", scalar.UintHex))
				sectionOffsets.warpDeformerKeyforms.keyformPositionSourcesBeginIndices = int64(d.FieldU32("keyform_position_sources_begin_indices", scalar.UintHex))
			})

			d.FieldStruct("rotation_deformer_keyforms", func(d *decode.D) {
				sectionOffsets.rotationDeformerKeyforms.opacities = int64(d.FieldU32("opacities", scalar.UintHex))
				sectionOffsets.rotationDeformerKeyforms.angles = int64(d.FieldU32("angles", scalar.UintHex))
				sectionOffsets.rotationDeformerKeyforms.originX = int64(d.FieldU32("origin_x", scalar.UintHex))
				sectionOffsets.rotationDeformerKeyforms.originY = int64(d.FieldU32("origin_y", scalar.UintHex))
				sectionOffsets.rotationDeformerKeyforms.scales = int64(d.FieldU32("scales", scalar.UintHex))
				sectionOffsets.rotationDeformerKeyforms.isReflectX = int64(d.FieldU32("is_reflect_x", scalar.UintHex))
				sectionOffsets.rotationDeformerKeyforms.isReflectY = int64(d.FieldU32("is_reflect_y", scalar.UintHex))
			})

			d.FieldStruct("art_mesh_keyforms", func(d *decode.D) {
				sectionOffsets.artMeshKeyforms.opacities = int64(d.FieldU32("opacities", scalar.UintHex))
				sectionOffsets.artMeshKeyforms.drawOrders = int64(d.FieldU32("draw_orders", scalar.UintHex))
				sectionOffsets.artMeshKeyforms.keyformPositionSourcesBeginIndices = int64(d.FieldU32("keyform_position_sources_begin_indices", scalar.UintHex))
			})

			d.FieldStruct("keyform_positions", func(d *decode.D) {
				sectionOffsets.keyformPositions.xys = int64(d.FieldU32("xys", scalar.UintHex))
			})

			d.FieldStruct("parameter_binding_indices", func(d *decode.D) {
				sectionOffsets.parameterBindingIndices.bindingSourcesIndices = int64(d.FieldU32("binding_sources_indices", scalar.UintHex))
			})

			d.FieldStruct("keyform_bindings", func(d *decode.D) {
				sectionOffsets.keyformBindings.parameterBindingIndexSourcesBeginIndices = int64(d.FieldU32("parameter_binding_index_sources_begin_indices", scalar.UintHex))
				sectionOffsets.keyformBindings.parameterBindingIndexSourcesCounts = int64(d.FieldU32("parameter_binding_index_sources_counts", scalar.UintHex))
			})

			d.FieldStruct("parameter_bindings", func(d *decode.D) {
				sectionOffsets.parameterBindings.keysSourcesBeginIndices = int64(d.FieldU32("keys_sources_begin_indices", scalar.UintHex))
				sectionOffsets.parameterBindings.keysSourcesCounts = int64(d.FieldU32("keys_sources_counts", scalar.UintHex))
			})

			d.FieldStruct("keys", func(d *decode.D) {
				sectionOffsets.keys.values = int64(d.FieldU32("values", scalar.UintHex))
			})

			d.FieldStruct("uvs", func(d *decode.D) {
				sectionOffsets.UVs.uvs = int64(d.FieldU32("uvs", scalar.UintHex))
			})

			d.FieldStruct("position_indices", func(d *decode.D) {
				sectionOffsets.positionIndices.indices = int64(d.FieldU32("indices", scalar.UintHex))
			})

			d.FieldStruct("drawable_masks", func(d *decode.D) {
				sectionOffsets.drawableMasks.artMeshSourcesIndices = int64(d.FieldU32("art_mesh_sources_indices", scalar.UintHex))
			})

			d.FieldStruct("draw_order_groups", func(d *decode.D) {
				sectionOffsets.drawOrderGroups.objectSourcesBeginIndices = int64(d.FieldU32("object_sources_begin_indices", scalar.UintHex))
				sectionOffsets.drawOrderGroups.objectSourcesCounts = int64(d.FieldU32("object_sources_counts", scalar.UintHex))
				sectionOffsets.drawOrderGroups.objectSourcesTotalCounts = int64(d.FieldU32("object_sources_total_counts", scalar.UintHex))
				sectionOffsets.drawOrderGroups.maximumDrawOrders = int64(d.FieldU32("maximum_draw_orders", scalar.UintHex))
				sectionOffsets.drawOrderGroups.minimumDrawOrders = int64(d.FieldU32("minimum_draw_orders", scalar.UintHex))
			})

			d.FieldStruct("draw_order_group_objects", func(d *decode.D) {
				sectionOffsets.drawOrderGroupObjects.types = int64(d.FieldU32("types", scalar.UintHex))
				sectionOffsets.drawOrderGroupObjects.indices = int64(d.FieldU32("indices", scalar.UintHex))
				sectionOffsets.drawOrderGroupObjects.selfIndices = int64(d.FieldU32("self_indices", scalar.UintHex))
			})

			d.FieldStruct("glue", func(d *decode.D) {
				sectionOffsets.glue.runtimeSpace0 = int64(d.FieldU32("runtime_space0", scalar.UintHex))
				sectionOffsets.glue.ids = int64(d.FieldU32("ids", scalar.UintHex))
				sectionOffsets.glue.keyformBindingSourcesIndices = int64(d.FieldU32("keyform_binding_sources_indices", scalar.UintHex))
				sectionOffsets.glue.keyformSourcesBeginIndices = int64(d.FieldU32("keyform_sources_begin_indices", scalar.UintHex))
				sectionOffsets.glue.keyformSourcesCounts = int64(d.FieldU32("keyform_sources_counts", scalar.UintHex))
				sectionOffsets.glue.artMeshIndicesA = int64(d.FieldU32("art_mesh_indices_a", scalar.UintHex))
				sectionOffsets.glue.artMeshIndicesB = int64(d.FieldU32("art_mesh_indices_b", scalar.UintHex))
				sectionOffsets.glue.glueInfoSourcesBeginIndices = int64(d.FieldU32("glue_info_sources_begin_indices", scalar.UintHex))
				sectionOffsets.glue.glueInfoSourcesCounts = int64(d.FieldU32("glue_info_sources_counts", scalar.UintHex))
			})

			d.FieldStruct("glue_info", func(d *decode.D) {
				sectionOffsets.glueInfo.weights = int64(d.FieldU32("weights", scalar.UintHex))
				sectionOffsets.glueInfo.positionIndices = int64(d.FieldU32("position_indices", scalar.UintHex))
			})

			d.FieldStruct("glue_keyforms", func(d *decode.D) {
				sectionOffsets.glueKeyforms.intensities = int64(d.FieldU32("intensities", scalar.UintHex))
			})

			if version >= moc3Version3_03_00 {
				d.FieldStruct("warp_deformers_v3_3", func(d *decode.D) {
					sectionOffsets.warpDeformersV3_3.isQuadSource = int64(d.FieldU32("is_quad_source", scalar.UintHex))
				})
			}

			if version >= moc3Version4_02_00 {
				d.FieldStruct("parameter_extensions", func(d *decode.D) {
					sectionOffsets.parameterExtensions.runtimeSpace0 = int64(d.FieldU32("runtime_space0", scalar.UintHex))
					sectionOffsets.parameterExtensions.keysSourcesBeginIndices = int64(d.FieldU32("keys_sources_begin_indices", scalar.UintHex))
					sectionOffsets.parameterExtensions.keysSourcesCounts = int64(d.FieldU32("keys_sources_counts", scalar.UintHex))
				})

				d.FieldStruct("warp_deformers_v4_2", func(d *decode.D) {
					sectionOffsets.warpDeformersV4_2.keyformColorSourcesBeginIndices = int64(d.FieldU32("keyform_color_sources_begin_indices", scalar.UintHex))
				})

				d.FieldStruct("rotation_deformers_v4_2", func(d *decode.D) {
					sectionOffsets.rotationDeformersV4_2.keyformColorSourcesBeginIndices = int64(d.FieldU32("keyform_color_sources_begin_indices", scalar.UintHex))
				})

				d.FieldStruct("art_meshes_v4_2", func(d *decode.D) {
					sectionOffsets.artMeshesV4_2.keyformColorSourcesBeginIndices = int64(d.FieldU32("keyform_color_sources_begin_indices", scalar.UintHex))
				})

				d.FieldStruct("keyform_multiply_colors", func(d *decode.D) {
					sectionOffsets.keyformMultiplyColors.r = int64(d.FieldU32("r", scalar.UintHex))
					sectionOffsets.keyformMultiplyColors.g = int64(d.FieldU32("g", scalar.UintHex))
					sectionOffsets.keyformMultiplyColors.b = int64(d.FieldU32("b", scalar.UintHex))
				})

				d.FieldStruct("keyform_screen_colors", func(d *decode.D) {
					sectionOffsets.keyformScreenColors.r = int64(d.FieldU32("r", scalar.UintHex))
					sectionOffsets.keyformScreenColors.g = int64(d.FieldU32("g", scalar.UintHex))
					sectionOffsets.keyformScreenColors.b = int64(d.FieldU32("b", scalar.UintHex))
				})

				d.FieldStruct("parameters_v4_2", func(d *decode.D) {
					sectionOffsets.parametersV4_2.parameterTypes = int64(d.FieldU32("parameter_types", scalar.UintHex))
					sectionOffsets.parametersV4_2.blendShapeParameterBindingSourcesBeginIndices = int64(d.FieldU32("blend_shape_parameter_binding_sources_begin_indices", scalar.UintHex))
					sectionOffsets.parametersV4_2.blendShapeParameterBindingSourcesCounts = int64(d.FieldU32("blend_shape_parameter_binding_sources_counts", scalar.UintHex))
				})

				d.FieldStruct("blend_shape_parameter_bindings", func(d *decode.D) {
					sectionOffsets.blendShapeParameterBindings.keysSourcesBeginIndices = int64(d.FieldU32("keys_sources_begin_indices", scalar.UintHex))
					sectionOffsets.blendShapeParameterBindings.keysSourcesCounts = int64(d.FieldU32("keys_sources_counts", scalar.UintHex))
					sectionOffsets.blendShapeParameterBindings.baseKeyIndices = int64(d.FieldU32("base_key_indices", scalar.UintHex))
				})

				d.FieldStruct("blend_shape_keyform_bindings", func(d *decode.D) {
					sectionOffsets.blendShapeKeyformBindings.parameterBindingSourcesIndices = int64(d.FieldU32("parameter_binding_sources_indices", scalar.UintHex))
					sectionOffsets.blendShapeKeyformBindings.keyformSourcesBlendShapeIndices = int64(d.FieldU32("keyform_sources_blend_shape_indices", scalar.UintHex))
					sectionOffsets.blendShapeKeyformBindings.keyformSourcesBlendShapeCounts = int64(d.FieldU32("keyform_sources_blend_shape_counts", scalar.UintHex))
					sectionOffsets.blendShapeKeyformBindings.blendShapeConstraintIndexSourcesBeginIndices = int64(d.FieldU32("blend_shape_constraint_index_sources_begin_indices", scalar.UintHex))
					sectionOffsets.blendShapeKeyformBindings.blendShapeConstraintIndexSourcesCounts = int64(d.FieldU32("blend_shape_constraint_index_sources_counts", scalar.UintHex))
				})

				d.FieldStruct("blend_shapes_warp_deformers", func(d *decode.D) {
					sectionOffsets.blendShapesWarpDeformers.targetIndices = int64(d.FieldU32("target_indices", scalar.UintHex))
					sectionOffsets.blendShapesWarpDeformers.blendShapeKeyformBindingSourcesBeginIndices = int64(d.FieldU32("blend_shape_keyform_binding_sources_begin_indices", scalar.UintHex))
					sectionOffsets.blendShapesWarpDeformers.blendShapeKeyformBindingSourcesCounts = int64(d.FieldU32("blend_shape_keyform_binding_sources_counts", scalar.UintHex))
				})

				d.FieldStruct("blend_shapes_art_meshes", func(d *decode.D) {
					sectionOffsets.blendShapesArtMeshes.targetIndices = int64(d.FieldU32("target_indices", scalar.UintHex))
					sectionOffsets.blendShapesArtMeshes.blendShapeKeyformBindingSourcesBeginIndices = int64(d.FieldU32("blend_shape_keyform_binding_sources_begin_indices", scalar.UintHex))
					sectionOffsets.blendShapesArtMeshes.blendShapeKeyformBindingSourcesCounts = int64(d.FieldU32("blend_shape_keyform_binding_sources_counts", scalar.UintHex))
				})

				d.FieldStruct("blend_shape_constraint_indices", func(d *decode.D) {
					sectionOffsets.blendShapeConstraintIndices.blendShapeConstraintSourcesIndices = int64(d.FieldU32("blend_shape_constraint_sources_indices", scalar.UintHex))
				})

				d.FieldStruct("blend_shape_constraints", func(d *decode.D) {
					sectionOffsets.blendShapeConstraints.parameterIndices = int64(d.FieldU32("parameter_indices", scalar.UintHex))
					sectionOffsets.blendShapeConstraints.blendShapeConstraintValueSourcesBeginIndices = int64(d.FieldU32("blend_shape_constraint_value_sources_begin_indices", scalar.UintHex))
					sectionOffsets.blendShapeConstraints.blendShapeConstraintValueSourcesCounts = int64(d.FieldU32("blend_shape_constraint_value_sources_counts", scalar.UintHex))
				})

				d.FieldStruct("blend_shape_constraint_values", func(d *decode.D) {
					sectionOffsets.blendShapeConstraintValues.keys = int64(d.FieldU32("keys", scalar.UintHex))
					sectionOffsets.blendShapeConstraintValues.weights = int64(d.FieldU32("weights", scalar.UintHex))
				})
			}

			if version >= moc3Version5_00_00 {
				d.FieldStruct("warp_deformer_keyforms_v5_0", func(d *decode.D) {
					sectionOffsets.warpDeformerKeyformsV5_0.keyformMultiplyColorSourcesBeginIndices = int64(d.FieldU32("keyform_multiply_color_sources_begin_indices", scalar.UintHex))
					sectionOffsets.warpDeformerKeyformsV5_0.keyformScreenColorSourcesBeginIndices = int64(d.FieldU32("keyform_screen_color_sources_begin_indices", scalar.UintHex))
				})

				d.FieldStruct("rotation_deformer_keyforms_v5_0", func(d *decode.D) {
					sectionOffsets.rotationDeformerKeyformsV5_0.keyformMultiplyColorSourcesBeginIndices = int64(d.FieldU32("keyform_multiply_color_sources_begin_indices", scalar.UintHex))
					sectionOffsets.rotationDeformerKeyformsV5_0.keyformScreenColorSourcesBeginIndices = int64(d.FieldU32("keyform_screen_color_sources_begin_indices", scalar.UintHex))
				})

				d.FieldStruct("art_mesh_keyforms_v5_0", func(d *decode.D) {
					sectionOffsets.artMeshKeyformsV5_0.keyformMultiplyColorSourcesBeginIndices = int64(d.FieldU32("keyform_multiply_color_sources_begin_indices", scalar.UintHex))
					sectionOffsets.artMeshKeyformsV5_0.keyformScreenColorSourcesBeginIndices = int64(d.FieldU32("keyform_screen_color_sources_begin_indices", scalar.UintHex))
				})

				d.FieldStruct("blend_shapes_parts", func(d *decode.D) {
					sectionOffsets.blendShapesParts.targetIndices = int64(d.FieldU32("target_indices", scalar.UintHex))
					sectionOffsets.blendShapesParts.blendShapeKeyformBindingSourcesBeginIndices = int64(d.FieldU32("blend_shape_keyform_binding_sources_begin_indices", scalar.UintHex))
					sectionOffsets.blendShapesParts.blendShapeKeyformBindingSourcesCounts = int64(d.FieldU32("blend_shape_keyform_binding_sources_counts", scalar.UintHex))
				})

				d.FieldStruct("blend_shapes_rotation_deformers", func(d *decode.D) {
					sectionOffsets.blendShapesRotationDeformers.targetIndices = int64(d.FieldU32("target_indices", scalar.UintHex))
					sectionOffsets.blendShapesRotationDeformers.blendShapeKeyformBindingSourcesBeginIndices = int64(d.FieldU32("blend_shape_keyform_binding_sources_begin_indices", scalar.UintHex))
					sectionOffsets.blendShapesRotationDeformers.blendShapeKeyformBindingSourcesCounts = int64(d.FieldU32("blend_shape_keyform_binding_sources_counts", scalar.UintHex))
				})

				d.FieldStruct("blend_shapes_glue", func(d *decode.D) {
					sectionOffsets.blendShapesGlue.targetIndices = int64(d.FieldU32("target_indices", scalar.UintHex))
					sectionOffsets.blendShapesGlue.blendShapeKeyformBindingSourcesBeginIndices = int64(d.FieldU32("blend_shape_keyform_binding_sources_begin_indices", scalar.UintHex))
					sectionOffsets.blendShapesGlue.blendShapeKeyformBindingSourcesCounts = int64(d.FieldU32("blend_shape_keyform_binding_sources_counts", scalar.UintHex))
				})
			}

			d.FieldRawLen("reserved", d.BitsLeft())
		})
	})

	d.FieldRawLen("runtime_address_map", 0x480*8)
	d.FieldRawLen("unused1", 0x7c0*8-d.Pos())

	var countInfo countInfoTable
	d.FieldStruct("sections", func(d *decode.D) {
		d.SeekAbs(sectionOffsets.countInfo * 8)
		d.FieldStruct("count_info", func(d *decode.D) {
			countInfo.parts = int64(d.FieldU32("parts"))
			countInfo.deformers = int64(d.FieldU32("deformers"))
			countInfo.warpDeformers = int64(d.FieldU32("warp_deformers"))
			countInfo.rotationDeformers = int64(d.FieldU32("rotation_deformers"))
			countInfo.artMeshes = int64(d.FieldU32("art_meshes"))
			countInfo.parameters = int64(d.FieldU32("parameters"))
			countInfo.partKeyforms = int64(d.FieldU32("part_keyforms"))
			countInfo.warpDeformerKeyforms = int64(d.FieldU32("warp_deformer_keyforms"))
			countInfo.rotationDeformerKeyforms = int64(d.FieldU32("rotation_deformer_keyforms"))
			countInfo.artMeshKeyforms = int64(d.FieldU32("art_mesh_keyforms"))
			countInfo.keyformPositions = int64(d.FieldU32("keyform_positions"))
			countInfo.parameterBindingIndices = int64(d.FieldU32("parameter_binding_indices"))
			countInfo.keyformBindings = int64(d.FieldU32("keyform_bindings"))
			countInfo.parameterBindings = int64(d.FieldU32("parameter_bindings"))
			countInfo.keys = int64(d.FieldU32("keys"))
			countInfo.uvs = int64(d.FieldU32("uvs"))
			countInfo.positionIndices = int64(d.FieldU32("position_indices"))
			countInfo.drawableMasks = int64(d.FieldU32("drawable_masks"))
			countInfo.drawOrderGroups = int64(d.FieldU32("draw_order_groups"))
			countInfo.drawOrderGroupObjects = int64(d.FieldU32("draw_order_group_objects"))
			countInfo.glue = int64(d.FieldU32("glue"))
			countInfo.glueInfo = int64(d.FieldU32("glue_info"))
			countInfo.glueKeyforms = int64(d.FieldU32("glue_keyforms"))

			if version >= moc3Version4_02_00 {
				countInfo.keyformMultiplyColors = int64(d.FieldU32("keyform_multiply_colors"))
				countInfo.keyformScreenColors = int64(d.FieldU32("keyform_screen_colors"))
				countInfo.blendShapeParameterBindings = int64(d.FieldU32("blend_shape_parameter_bindings"))
				countInfo.blendShapeKeyformBindings = int64(d.FieldU32("blend_shape_keyform_bindings"))
				countInfo.blendShapesWarpDeformers = int64(d.FieldU32("blend_shapes_warp_deformers"))
				countInfo.blendShapesArtMeshes = int64(d.FieldU32("blend_shapes_art_meshes"))
				countInfo.blendShapeConstraintIndices = int64(d.FieldU32("blend_shape_constraint_indices"))
				countInfo.blendShapeConstraints = int64(d.FieldU32("blend_shape_constraints"))
				countInfo.blendShapeConstraintValues = int64(d.FieldU32("blend_shape_constraint_values"))
			}

			if version >= moc3Version5_00_00 {
				countInfo.blendShapesParts = int64(d.FieldU32("blend_shapes_parts"))
				countInfo.blendShapesRotationDeformers = int64(d.FieldU32("blend_shapes_rotation_deformers"))
				countInfo.blendShapesGlue = int64(d.FieldU32("blend_shapes_glue"))
			}

			var reserved int64
			if version >= moc3Version5_00_00 {
				reserved = int64(d.AlignBits(256 * 8))
			} else {
				reserved = int64(d.AlignBits(128 * 8))
			}
			d.FieldRawLen("reserved", reserved)
		})

		d.SeekAbs(sectionOffsets.canvasInfo * 8)
		d.FieldStruct("canvas_info", func(d *decode.D) {
			d.FieldF32("pixels_per_unit")
			d.FieldF32("origin_x")
			d.FieldF32("origin_y")
			d.FieldF32("canvas_width")
			d.FieldF32("canvas_height")
			d.FieldStruct("canvas_flags", func(d *decode.D) {
				d.FieldU7("reserved")
				d.FieldBool("reverse_y_coordinate")
			})

			d.FieldRawLen("padding", int64(d.AlignBits(moc3AlignBits)))
		})

		d.FieldStruct("parts", func(d *decode.D) {
			count := countInfo.parts

			d.SeekAbs(sectionOffsets.parts.runtimeSpace0 * 8)
			fieldRuntimeSpace(d, "runtime_space0", count, false)

			d.SeekAbs(sectionOffsets.parts.ids * 8)
			fieldAlignedNArray(d, "ids", count, func(d *decode.D) { d.FieldUTF8NullFixedLen("id", 64) })

			d.SeekAbs(sectionOffsets.parts.keyformBindingSourcesIndices * 8)
			fieldAlignedNArray(d, "keyform_binding_sources_indices", count, func(d *decode.D) { d.FieldS32("index") })

			d.SeekAbs(sectionOffsets.parts.keyformSourcesBeginIndices * 8)
			fieldAlignedNArray(d, "keyform_sources_begin_indices", count, func(d *decode.D) { d.FieldS32("index") })

			d.SeekAbs(sectionOffsets.parts.keyformSourcesCounts * 8)
			fieldAlignedNArray(d, "keyform_sources_counts", count, func(d *decode.D) { d.FieldS32("index") })

			d.SeekAbs(sectionOffsets.parts.isVisible * 8)
			fieldAlignedNArray(d, "is_visible", count, func(d *decode.D) { d.FieldU32("visible", boolNToSym) })

			d.SeekAbs(sectionOffsets.parts.isEnabled * 8)
			fieldAlignedNArray(d, "is_enabled", count, func(d *decode.D) { d.FieldU32("enabled", boolNToSym) })

			d.SeekAbs(sectionOffsets.parts.parentPartIndices * 8)
			fieldAlignedNArray(d, "parent_part_indices", count, func(d *decode.D) { d.FieldS32("index") })
		})

		d.FieldStruct("deformers", func(d *decode.D) {
			count := countInfo.deformers

			d.SeekAbs(sectionOffsets.deformers.runtimeSpace0 * 8)
			fieldRuntimeSpace(d, "runtime_space0", count, false)

			d.SeekAbs(sectionOffsets.deformers.ids * 8)
			fieldAlignedNArray(d, "ids", count, func(d *decode.D) { d.FieldUTF8NullFixedLen("id", 64) })

			d.SeekAbs(sectionOffsets.deformers.keyformBindingSourcesIndices * 8)
			fieldAlignedNArray(d, "keyform_binding_sources_indices", count, func(d *decode.D) { d.FieldS32("index") })

			d.SeekAbs(sectionOffsets.deformers.isVisible * 8)
			fieldAlignedNArray(d, "is_visible", count, func(d *decode.D) { d.FieldU32("visible", boolNToSym) })

			d.SeekAbs(sectionOffsets.deformers.isEnabled * 8)
			fieldAlignedNArray(d, "is_enabled", count, func(d *decode.D) { d.FieldU32("enabled", boolNToSym) })

			d.SeekAbs(sectionOffsets.deformers.parentPartIndices * 8)
			fieldAlignedNArray(d, "parent_part_indices", count, func(d *decode.D) { d.FieldS32("index") })

			d.SeekAbs(sectionOffsets.deformers.parentDeformerIndices * 8)
			fieldAlignedNArray(d, "parent_deformer_indices", count, func(d *decode.D) { d.FieldS32("index") })

			d.SeekAbs(sectionOffsets.deformers.types * 8)
			fieldAlignedNArray(d, "types", count, func(d *decode.D) { d.FieldU32("type", deformerTypeNames) })

			d.SeekAbs(sectionOffsets.deformers.specificSourcesIndices * 8)
			fieldAlignedNArray(d, "specific_sources_indices", count, func(d *decode.D) { d.FieldS32("index") })
		})

		d.FieldStruct("warp_deformers", func(d *decode.D) {
			count := countInfo.warpDeformers

			d.SeekAbs(sectionOffsets.warpDeformers.keyformBindingSourcesIndices * 8)
			fieldAlignedNArray(d, "keyform_binding_sources_indices", count, func(d *decode.D) { d.FieldS32("index") })

			d.SeekAbs(sectionOffsets.warpDeformers.keyformSourcesBeginIndices * 8)
			fieldAlignedNArray(d, "keyform_sources_begin_indices", count, func(d *decode.D) { d.FieldS32("index") })

			d.SeekAbs(sectionOffsets.warpDeformers.keyformSourcesCounts * 8)
			fieldAlignedNArray(d, "keyform_sources_counts", count, func(d *decode.D) { d.FieldS32("index") })

			d.SeekAbs(sectionOffsets.warpDeformers.vertexCounts * 8)
			fieldAlignedNArray(d, "vertex_counts", count, func(d *decode.D) { d.FieldS32("count") })

			d.SeekAbs(sectionOffsets.warpDeformers.rows * 8)
			fieldAlignedNArray(d, "rows", count, func(d *decode.D) { d.FieldU32("row") })

			d.SeekAbs(sectionOffsets.warpDeformers.columns * 8)
			fieldAlignedNArray(d, "columns", count, func(d *decode.D) { d.FieldU32("column") })

			if version >= moc3Version3_03_00 {
				d.SeekAbs(sectionOffsets.warpDeformersV3_3.isQuadSource * 8)
				fieldAlignedNArray(d, "is_quad_source", count, func(d *decode.D) { d.FieldU32("quad_source", boolNToSym) })
			}

			if version >= moc3Version4_02_00 {
				d.SeekAbs(sectionOffsets.warpDeformersV4_2.keyformColorSourcesBeginIndices * 8)
				fieldAlignedNArray(d, "keyform_color_sources_begin_indices", count, func(d *decode.D) { d.FieldS32("index") })
			}
		})

		d.FieldStruct("rotation_deformers", func(d *decode.D) {
			count := countInfo.rotationDeformers

			d.SeekAbs(sectionOffsets.rotationDeformers.keyformBindingSourcesIndices * 8)
			fieldAlignedNArray(d, "keyform_binding_sources_indices", count, func(d *decode.D) { d.FieldS32("index") })

			d.SeekAbs(sectionOffsets.rotationDeformers.keyformSourcesBeginIndices * 8)
			fieldAlignedNArray(d, "keyform_sources_begin_indices", count, func(d *decode.D) { d.FieldS32("index") })

			d.SeekAbs(sectionOffsets.rotationDeformers.keyformSourcesCounts * 8)
			fieldAlignedNArray(d, "keyform_sources_counts", count, func(d *decode.D) { d.FieldS32("index") })

			d.SeekAbs(sectionOffsets.rotationDeformers.baseAngles * 8)
			fieldAlignedNArray(d, "base_angles", count, func(d *decode.D) { d.FieldF32("angle") })

			if version >= moc3Version4_02_00 {
				d.SeekAbs(sectionOffsets.rotationDeformersV4_2.keyformColorSourcesBeginIndices * 8)
				fieldAlignedNArray(d, "keyform_color_sources_begin_indices", count, func(d *decode.D) { d.FieldS32("index") })
			}
		})

		d.FieldStruct("art_meshes", func(d *decode.D) {
			count := countInfo.artMeshes

			d.SeekAbs(sectionOffsets.artMeshes.runtimeSpace0 * 8)
			fieldRuntimeSpace(d, "runtime_space0", count, true)

			d.SeekAbs(sectionOffsets.artMeshes.runtimeSpace1 * 8)
			fieldRuntimeSpace(d, "runtime_space1", count, true)

			d.SeekAbs(sectionOffsets.artMeshes.runtimeSpace2 * 8)
			fieldRuntimeSpace(d, "runtime_space2", count, true)

			d.SeekAbs(sectionOffsets.artMeshes.runtimeSpace3 * 8)
			fieldRuntimeSpace(d, "runtime_space3", count, false)

			d.SeekAbs(sectionOffsets.artMeshes.ids * 8)
			fieldAlignedNArray(d, "ids", count, func(d *decode.D) { d.FieldUTF8NullFixedLen("id", 64) })

			d.SeekAbs(sectionOffsets.artMeshes.keyformBindingSourcesIndices * 8)
			fieldAlignedNArray(d, "keyform_binding_sources_indices", count, func(d *decode.D) { d.FieldS32("index") })

			d.SeekAbs(sectionOffsets.artMeshes.keyformSourcesBeginIndices * 8)
			fieldAlignedNArray(d, "keyform_sources_begin_indices", count, func(d *decode.D) { d.FieldS32("index") })

			d.SeekAbs(sectionOffsets.artMeshes.keyformSourcesCounts * 8)
			fieldAlignedNArray(d, "keyform_sources_counts", count, func(d *decode.D) { d.FieldS32("index") })

			d.SeekAbs(sectionOffsets.artMeshes.isVisible * 8)
			fieldAlignedNArray(d, "is_visible", count, func(d *decode.D) { d.FieldU32("visible", boolNToSym) })

			d.SeekAbs(sectionOffsets.artMeshes.isEnabled * 8)
			fieldAlignedNArray(d, "is_enabled", count, func(d *decode.D) { d.FieldU32("enabled", boolNToSym) })

			d.SeekAbs(sectionOffsets.artMeshes.parentPartIndices * 8)
			fieldAlignedNArray(d, "parent_part_indices", count, func(d *decode.D) { d.FieldS32("index") })

			d.SeekAbs(sectionOffsets.artMeshes.parentDeformerIndices * 8)
			fieldAlignedNArray(d, "parent_deformer_indices", count, func(d *decode.D) { d.FieldS32("index") })

			d.SeekAbs(sectionOffsets.artMeshes.textureNos * 8)
			fieldAlignedNArray(d, "texture_nos", count, func(d *decode.D) { d.FieldU32("texture_no") })

			d.SeekAbs(sectionOffsets.artMeshes.drawableFlags * 8)
			fieldAlignedNArray(d, "drawable_flags", count, func(d *decode.D) {
				d.FieldStruct("flags", func(d *decode.D) {
					d.FieldU4("reserved")
					d.FieldBool("is_inverted")
					d.FieldBool("is_double_sided")
					d.FieldU2("blend_mode", blendModeNames)
				})
			})

			d.SeekAbs(sectionOffsets.artMeshes.vertexCounts * 8)
			fieldAlignedNArray(d, "vertex_counts", count, func(d *decode.D) { d.FieldS32("count") })

			d.SeekAbs(sectionOffsets.artMeshes.uvSourcesBeginIndices * 8)
			fieldAlignedNArray(d, "uv_sources_begin_indices", count, func(d *decode.D) { d.FieldS32("index") })

			d.SeekAbs(sectionOffsets.artMeshes.positionIndexSourcesBeginIndices * 8)
			fieldAlignedNArray(d, "position_index_sources_begin_indices", count, func(d *decode.D) { d.FieldS32("index") })

			d.SeekAbs(sectionOffsets.artMeshes.positionIndexSourcesCounts * 8)
			fieldAlignedNArray(d, "position_index_sources_counts", count, func(d *decode.D) { d.FieldS32("count") })

			d.SeekAbs(sectionOffsets.artMeshes.drawableMaskSourcesBeginIndices * 8)
			fieldAlignedNArray(d, "drawable_mask_sources_begin_indices", count, func(d *decode.D) { d.FieldS32("index") })

			d.SeekAbs(sectionOffsets.artMeshes.drawableMaskSourcesCounts * 8)
			fieldAlignedNArray(d, "drawable_mask_sources_counts", count, func(d *decode.D) { d.FieldS32("count") })

			if version >= moc3Version4_02_00 {
				d.SeekAbs(sectionOffsets.artMeshesV4_2.keyformColorSourcesBeginIndices * 8)
				fieldAlignedNArray(d, "keyform_color_sources_begin_indices", count, func(d *decode.D) { d.FieldS32("index") })
			}
		})

		d.FieldStruct("parameters", func(d *decode.D) {
			count := countInfo.parameters

			d.SeekAbs(sectionOffsets.parameters.runtimeSpace0 * 8)
			fieldRuntimeSpace(d, "runtime_space0", count, false)

			d.SeekAbs(sectionOffsets.parameters.ids * 8)
			fieldAlignedNArray(d, "ids", count, func(d *decode.D) { d.FieldUTF8NullFixedLen("id", 64) })

			d.SeekAbs(sectionOffsets.parameters.maxValues * 8)
			fieldAlignedNArray(d, "max_values", count, func(d *decode.D) { d.FieldF32("value") })

			d.SeekAbs(sectionOffsets.parameters.minValues * 8)
			fieldAlignedNArray(d, "min_values", count, func(d *decode.D) { d.FieldF32("value") })

			d.SeekAbs(sectionOffsets.parameters.defaultValues * 8)
			fieldAlignedNArray(d, "default_values", count, func(d *decode.D) { d.FieldF32("value") })

			d.SeekAbs(sectionOffsets.parameters.isRepeat * 8)
			fieldAlignedNArray(d, "is_repeat", count, func(d *decode.D) { d.FieldU32("repeat", boolNToSym) })

			d.SeekAbs(sectionOffsets.parameters.decimalPlaces * 8)
			fieldAlignedNArray(d, "decimal_places", count, func(d *decode.D) { d.FieldU32("value") })

			d.SeekAbs(sectionOffsets.parameters.parameterBindingSourcesBeginIndices * 8)
			fieldAlignedNArray(d, "parameter_binding_sources_begin_indices", count, func(d *decode.D) { d.FieldS32("index") })

			d.SeekAbs(sectionOffsets.parameters.parameterBindingSourcesCounts * 8)
			fieldAlignedNArray(d, "parameter_binding_sources_counts", count, func(d *decode.D) { d.FieldS32("count") })

			if version >= moc3Version4_02_00 {
				d.SeekAbs(sectionOffsets.parametersV4_2.parameterTypes * 8)
				fieldAlignedNArray(d, "parameter_types", count, func(d *decode.D) { d.FieldU32("type", parameterTypeNames) })

				d.SeekAbs(sectionOffsets.parametersV4_2.blendShapeParameterBindingSourcesBeginIndices * 8)
				fieldAlignedNArray(d, "blend_shape_parameter_binding_sources_begin_indices", count, func(d *decode.D) { d.FieldS32("index") })

				d.SeekAbs(sectionOffsets.parametersV4_2.blendShapeParameterBindingSourcesCounts * 8)
				fieldAlignedNArray(d, "blend_shape_parameter_binding_sources_counts", count, func(d *decode.D) { d.FieldS32("count") })
			}
		})

		d.FieldStruct("part_keyforms", func(d *decode.D) {
			d.SeekAbs(sectionOffsets.partKeyforms.drawOrders * 8)
			fieldAlignedNArray(d, "draw_orders", countInfo.partKeyforms, func(d *decode.D) { d.FieldF32("draw_order") })
		})

		d.FieldStruct("warp_deformer_keyforms", func(d *decode.D) {
			count := countInfo.warpDeformerKeyforms

			d.SeekAbs(sectionOffsets.warpDeformerKeyforms.opacities * 8)
			fieldAlignedNArray(d, "opacities", count, func(d *decode.D) { d.FieldF32("opacity") })

			d.SeekAbs(sectionOffsets.warpDeformerKeyforms.keyformPositionSourcesBeginIndices * 8)
			fieldAlignedNArray(d, "keyform_position_sources_begin_indices", count, func(d *decode.D) { d.FieldS32("index") })

			if version >= moc3Version5_00_00 {
				d.SeekAbs(sectionOffsets.warpDeformerKeyformsV5_0.keyformMultiplyColorSourcesBeginIndices * 8)
				fieldAlignedNArray(d, "keyform_multiply_color_sources_begin_indices", count, func(d *decode.D) { d.FieldS32("index") })

				d.SeekAbs(sectionOffsets.warpDeformerKeyformsV5_0.keyformScreenColorSourcesBeginIndices * 8)
				fieldAlignedNArray(d, "keyform_screen_color_sources_begin_indices", count, func(d *decode.D) { d.FieldS32("index") })
			}
		})

		d.FieldStruct("rotation_deformer_keyforms", func(d *decode.D) {
			count := countInfo.rotationDeformerKeyforms

			d.SeekAbs(sectionOffsets.rotationDeformerKeyforms.opacities * 8)
			fieldAlignedNArray(d, "opacities", count, func(d *decode.D) { d.FieldF32("value") })

			d.SeekAbs(sectionOffsets.rotationDeformerKeyforms.angles * 8)
			fieldAlignedNArray(d, "angles", count, func(d *decode.D) { d.FieldF32("angle") })

			d.SeekAbs(sectionOffsets.rotationDeformerKeyforms.originX * 8)
			fieldAlignedNArray(d, "origin_x", count, func(d *decode.D) { d.FieldF32("x") })

			d.SeekAbs(sectionOffsets.rotationDeformerKeyforms.originY * 8)
			fieldAlignedNArray(d, "origin_y", count, func(d *decode.D) { d.FieldF32("y") })

			d.SeekAbs(sectionOffsets.rotationDeformerKeyforms.scales * 8)
			fieldAlignedNArray(d, "scales", count, func(d *decode.D) { d.FieldF32("value") })

			d.SeekAbs(sectionOffsets.rotationDeformerKeyforms.isReflectX * 8)
			fieldAlignedNArray(d, "is_reflect_x", count, func(d *decode.D) { d.FieldU32("reflect_x", boolNToSym) })

			d.SeekAbs(sectionOffsets.rotationDeformerKeyforms.isReflectY * 8)
			fieldAlignedNArray(d, "is_reflect_y", count, func(d *decode.D) { d.FieldU32("reflect_y", boolNToSym) })

			if version >= moc3Version5_00_00 {
				d.SeekAbs(sectionOffsets.rotationDeformerKeyformsV5_0.keyformMultiplyColorSourcesBeginIndices * 8)
				fieldAlignedNArray(d, "keyform_multiply_color_sources_begin_indices", count, func(d *decode.D) { d.FieldS32("index") })

				d.SeekAbs(sectionOffsets.rotationDeformerKeyformsV5_0.keyformScreenColorSourcesBeginIndices * 8)
				fieldAlignedNArray(d, "keyform_screen_color_sources_begin_indices", count, func(d *decode.D) { d.FieldS32("index") })
			}
		})

		d.FieldStruct("art_mesh_keyforms", func(d *decode.D) {
			count := countInfo.artMeshKeyforms

			d.SeekAbs(sectionOffsets.artMeshKeyforms.opacities * 8)
			fieldAlignedNArray(d, "opacities", count, func(d *decode.D) { d.FieldF32("value") })

			d.SeekAbs(sectionOffsets.artMeshKeyforms.drawOrders * 8)
			fieldAlignedNArray(d, "draw_orders", count, func(d *decode.D) { d.FieldF32("draw_order") })

			d.SeekAbs(sectionOffsets.artMeshKeyforms.keyformPositionSourcesBeginIndices * 8)
			fieldAlignedNArray(d, "keyform_position_sources_begin_indices", count, func(d *decode.D) { d.FieldS32("index") })

			if version >= moc3Version5_00_00 {
				d.SeekAbs(sectionOffsets.artMeshKeyformsV5_0.keyformMultiplyColorSourcesBeginIndices * 8)
				fieldAlignedNArray(d, "keyform_multiply_color_sources_begin_indices", count, func(d *decode.D) { d.FieldS32("index") })

				d.SeekAbs(sectionOffsets.artMeshKeyformsV5_0.keyformScreenColorSourcesBeginIndices * 8)
				fieldAlignedNArray(d, "keyform_screen_color_sources_begin_indices", count, func(d *decode.D) { d.FieldS32("index") })
			}
		})

		d.FieldStruct("keyform_positions", func(d *decode.D) {
			d.SeekAbs(sectionOffsets.keyformPositions.xys * 8)
			fieldAlignedNArray(d, "xys", countInfo.keyformPositions/2, func(d *decode.D) { d.FieldF32("x"); d.FieldF32("y") })
		})

		d.FieldStruct("parameter_binding_indices", func(d *decode.D) {
			d.SeekAbs(sectionOffsets.parameterBindingIndices.bindingSourcesIndices * 8)
			fieldAlignedNArray(d, "binding_sources_indices", countInfo.parameterBindingIndices, func(d *decode.D) { d.FieldS32("index") })
		})

		d.FieldStruct("keyform_bindings", func(d *decode.D) {
			count := countInfo.keyformBindings

			d.SeekAbs(sectionOffsets.keyformBindings.parameterBindingIndexSourcesBeginIndices * 8)
			fieldAlignedNArray(d, "parameter_binding_index_sources_begin_indices", count, func(d *decode.D) { d.FieldS32("index") })

			d.SeekAbs(sectionOffsets.keyformBindings.parameterBindingIndexSourcesCounts * 8)
			fieldAlignedNArray(d, "parameter_binding_index_sources_counts", count, func(d *decode.D) { d.FieldS32("count") })
		})

		d.FieldStruct("parameter_bindings", func(d *decode.D) {
			count := countInfo.parameterBindings

			d.SeekAbs(sectionOffsets.parameterBindings.keysSourcesBeginIndices * 8)
			fieldAlignedNArray(d, "keys_sources_begin_indices", count, func(d *decode.D) { d.FieldS32("index") })

			d.SeekAbs(sectionOffsets.parameterBindings.keysSourcesCounts * 8)
			fieldAlignedNArray(d, "keys_sources_counts", count, func(d *decode.D) { d.FieldS32("count") })
		})

		d.FieldStruct("keys", func(d *decode.D) {
			d.SeekAbs(sectionOffsets.keys.values * 8)
			fieldAlignedNArray(d, "values", countInfo.keys, func(d *decode.D) { d.FieldF32("value") })
		})

		d.FieldStruct("uvs", func(d *decode.D) {
			d.SeekAbs(sectionOffsets.UVs.uvs * 8)
			fieldAlignedNArray(d, "uvs", countInfo.uvs/2, func(d *decode.D) { d.FieldF32("u"); d.FieldF32("v") })
		})

		d.FieldStruct("position_indices", func(d *decode.D) {
			d.SeekAbs(sectionOffsets.positionIndices.indices * 8)
			fieldAlignedNArray(d, "indices", countInfo.positionIndices, func(d *decode.D) { d.FieldS16("index") })
		})

		d.FieldStruct("drawable_masks", func(d *decode.D) {
			d.SeekAbs(sectionOffsets.drawableMasks.artMeshSourcesIndices * 8)
			fieldAlignedNArray(d, "art_mesh_sources_indices", countInfo.drawableMasks, func(d *decode.D) { d.FieldS32("index") })
		})

		d.FieldStruct("draw_order_groups", func(d *decode.D) {
			count := countInfo.drawOrderGroups

			d.SeekAbs(sectionOffsets.drawOrderGroups.objectSourcesBeginIndices * 8)
			fieldAlignedNArray(d, "object_sources_begin_indices", count, func(d *decode.D) { d.FieldS32("index") })

			d.SeekAbs(sectionOffsets.drawOrderGroups.objectSourcesCounts * 8)
			fieldAlignedNArray(d, "object_sources_counts", count, func(d *decode.D) { d.FieldS32("count") })

			d.SeekAbs(sectionOffsets.drawOrderGroups.objectSourcesTotalCounts * 8)
			fieldAlignedNArray(d, "object_sources_total_counts", count, func(d *decode.D) { d.FieldS32("count") })

			d.SeekAbs(sectionOffsets.drawOrderGroups.maximumDrawOrders * 8)
			fieldAlignedNArray(d, "maximum_draw_orders", count, func(d *decode.D) { d.FieldU32("draw_order") })

			d.SeekAbs(sectionOffsets.drawOrderGroups.minimumDrawOrders * 8)
			fieldAlignedNArray(d, "minimum_draw_orders", count, func(d *decode.D) { d.FieldU32("draw_order") })
		})

		d.FieldStruct("draw_order_group_objects", func(d *decode.D) {
			count := countInfo.drawOrderGroupObjects

			d.SeekAbs(sectionOffsets.drawOrderGroupObjects.types * 8)
			fieldAlignedNArray(d, "types", count, func(d *decode.D) { d.FieldU32("type", drawOrderGroupObjectTypeNames) })

			d.SeekAbs(sectionOffsets.drawOrderGroupObjects.indices * 8)
			fieldAlignedNArray(d, "indices", count, func(d *decode.D) { d.FieldS32("index") })

			d.SeekAbs(sectionOffsets.drawOrderGroupObjects.selfIndices * 8)
			fieldAlignedNArray(d, "self_indices", count, func(d *decode.D) { d.FieldS32("index") })
		})

		d.FieldStruct("glue", func(d *decode.D) {
			count := countInfo.glue

			d.SeekAbs(sectionOffsets.glue.runtimeSpace0 * 8)
			fieldRuntimeSpace(d, "runtime_space0", count, false)

			d.SeekAbs(sectionOffsets.glue.ids * 8)
			fieldAlignedNArray(d, "ids", count, func(d *decode.D) { d.FieldUTF8NullFixedLen("id", 64) })

			d.SeekAbs(sectionOffsets.glue.keyformBindingSourcesIndices * 8)
			fieldAlignedNArray(d, "keyform_binding_sources_indices", count, func(d *decode.D) { d.FieldS32("index") })

			d.SeekAbs(sectionOffsets.glue.keyformSourcesBeginIndices * 8)
			fieldAlignedNArray(d, "keyform_sources_begin_indices", count, func(d *decode.D) { d.FieldS32("index") })

			d.SeekAbs(sectionOffsets.glue.keyformSourcesCounts * 8)
			fieldAlignedNArray(d, "keyform_sources_counts", count, func(d *decode.D) { d.FieldS32("index") })

			d.SeekAbs(sectionOffsets.glue.artMeshIndicesA * 8)
			fieldAlignedNArray(d, "art_mesh_indices_a", count, func(d *decode.D) { d.FieldS32("index") })

			d.SeekAbs(sectionOffsets.glue.artMeshIndicesB * 8)
			fieldAlignedNArray(d, "art_mesh_indices_b", count, func(d *decode.D) { d.FieldS32("index") })

			d.SeekAbs(sectionOffsets.glue.glueInfoSourcesBeginIndices * 8)
			fieldAlignedNArray(d, "glue_info_sources_begin_indices", count, func(d *decode.D) { d.FieldS32("index") })

			d.SeekAbs(sectionOffsets.glue.glueInfoSourcesCounts * 8)
			fieldAlignedNArray(d, "glue_info_sources_counts", count, func(d *decode.D) { d.FieldS32("count") })
		})

		d.FieldStruct("glue_info", func(d *decode.D) {
			count := countInfo.glueInfo

			d.SeekAbs(sectionOffsets.glueInfo.weights * 8)
			fieldAlignedNArray(d, "weights", count, func(d *decode.D) { d.FieldF32("value") })

			d.SeekAbs(sectionOffsets.glueInfo.positionIndices * 8)
			fieldAlignedNArray(d, "position_indices", count, func(d *decode.D) { d.FieldS16("index") })
		})

		d.FieldStruct("glue_keyforms", func(d *decode.D) {
			d.SeekAbs(sectionOffsets.glueKeyforms.intensities * 8)
			fieldAlignedNArray(d, "intensities", countInfo.glueKeyforms, func(d *decode.D) { d.FieldF32("value") })
		})

		if version >= moc3Version4_02_00 {
			d.FieldStruct("parameter_extensions", func(d *decode.D) {
				count := countInfo.parameters

				d.SeekAbs(sectionOffsets.parameterExtensions.runtimeSpace0 * 8)
				fieldRuntimeSpace(d, "runtime_space0", count, false)

				d.SeekAbs(sectionOffsets.parameterExtensions.keysSourcesBeginIndices * 8)
				fieldAlignedNArray(d, "keys_sources_begin_indices", count, func(d *decode.D) { d.FieldS32("index") })

				d.SeekAbs(sectionOffsets.parameterExtensions.keysSourcesCounts * 8)
				fieldAlignedNArray(d, "keys_sources_counts", count, func(d *decode.D) { d.FieldS32("count") })
			})

			d.FieldStruct("keyform_multiply_colors", func(d *decode.D) {
				count := countInfo.keyformMultiplyColors

				d.SeekAbs(sectionOffsets.keyformMultiplyColors.r * 8)
				fieldAlignedNArray(d, "r", count, func(d *decode.D) { d.FieldF32("value") })

				d.SeekAbs(sectionOffsets.keyformMultiplyColors.g * 8)
				fieldAlignedNArray(d, "g", count, func(d *decode.D) { d.FieldF32("value") })

				d.SeekAbs(sectionOffsets.keyformMultiplyColors.b * 8)
				fieldAlignedNArray(d, "b", count, func(d *decode.D) { d.FieldF32("value") })
			})

			d.FieldStruct("keyform_screen_colors", func(d *decode.D) {
				count := countInfo.keyformScreenColors

				d.SeekAbs(sectionOffsets.keyformScreenColors.r * 8)
				fieldAlignedNArray(d, "r", count, func(d *decode.D) { d.FieldF32("value") })

				d.SeekAbs(sectionOffsets.keyformScreenColors.g * 8)
				fieldAlignedNArray(d, "g", count, func(d *decode.D) { d.FieldF32("value") })

				d.SeekAbs(sectionOffsets.keyformScreenColors.b * 8)
				fieldAlignedNArray(d, "b", count, func(d *decode.D) { d.FieldF32("value") })
			})

			d.FieldStruct("blend_shape_parameter_bindings", func(d *decode.D) {
				count := countInfo.blendShapeParameterBindings

				d.SeekAbs(sectionOffsets.blendShapeParameterBindings.keysSourcesBeginIndices * 8)
				fieldAlignedNArray(d, "keys_sources_begin_indices", count, func(d *decode.D) { d.FieldS32("index") })

				d.SeekAbs(sectionOffsets.blendShapeParameterBindings.keysSourcesCounts * 8)
				fieldAlignedNArray(d, "keys_sources_counts", count, func(d *decode.D) { d.FieldS32("count") })

				d.SeekAbs(sectionOffsets.blendShapeParameterBindings.baseKeyIndices * 8)
				fieldAlignedNArray(d, "base_key_indices", count, func(d *decode.D) { d.FieldS32("index") })
			})

			d.FieldStruct("blend_shape_keyform_bindings", func(d *decode.D) {
				count := countInfo.blendShapeKeyformBindings

				d.SeekAbs(sectionOffsets.blendShapeKeyformBindings.parameterBindingSourcesIndices * 8)
				fieldAlignedNArray(d, "parameter_binding_sources_indices", count, func(d *decode.D) { d.FieldS32("index") })

				d.SeekAbs(sectionOffsets.blendShapeKeyformBindings.keyformSourcesBlendShapeIndices * 8)
				fieldAlignedNArray(d, "keyform_sources_blend_shape_indices", count, func(d *decode.D) { d.FieldS32("index") })

				d.SeekAbs(sectionOffsets.blendShapeKeyformBindings.keyformSourcesBlendShapeCounts * 8)
				fieldAlignedNArray(d, "keyform_sources_blend_shape_counts", count, func(d *decode.D) { d.FieldS32("count") })

				d.SeekAbs(sectionOffsets.blendShapeKeyformBindings.blendShapeConstraintIndexSourcesBeginIndices * 8)
				fieldAlignedNArray(d, "blend_shape_constraint_index_sources_begin_indices", count, func(d *decode.D) { d.FieldS32("index") })

				d.SeekAbs(sectionOffsets.blendShapeKeyformBindings.blendShapeConstraintIndexSourcesCounts * 8)
				fieldAlignedNArray(d, "blend_shape_constraint_index_sources_counts", count, func(d *decode.D) { d.FieldS32("count") })
			})

			d.FieldStruct("blend_shapes_warp_deformers", func(d *decode.D) {
				count := countInfo.blendShapesWarpDeformers

				d.SeekAbs(sectionOffsets.blendShapesWarpDeformers.targetIndices * 8)
				fieldAlignedNArray(d, "target_indices", count, func(d *decode.D) { d.FieldS32("index") })

				d.SeekAbs(sectionOffsets.blendShapesWarpDeformers.blendShapeKeyformBindingSourcesBeginIndices * 8)
				fieldAlignedNArray(d, "blend_shape_keyform_binding_sources_begin_indices", count, func(d *decode.D) { d.FieldS32("index") })

				d.SeekAbs(sectionOffsets.blendShapesWarpDeformers.blendShapeKeyformBindingSourcesCounts * 8)
				fieldAlignedNArray(d, "blend_shape_keyform_binding_sources_counts", count, func(d *decode.D) { d.FieldS32("count") })
			})

			d.FieldStruct("blend_shapes_art_meshes", func(d *decode.D) {
				count := countInfo.blendShapesArtMeshes

				d.SeekAbs(sectionOffsets.blendShapesArtMeshes.targetIndices * 8)
				fieldAlignedNArray(d, "target_indices", count, func(d *decode.D) { d.FieldS32("index") })

				d.SeekAbs(sectionOffsets.blendShapesArtMeshes.blendShapeKeyformBindingSourcesBeginIndices * 8)
				fieldAlignedNArray(d, "blend_shape_keyform_binding_sources_begin_indices", count, func(d *decode.D) { d.FieldS32("index") })

				d.SeekAbs(sectionOffsets.blendShapesArtMeshes.blendShapeKeyformBindingSourcesCounts * 8)
				fieldAlignedNArray(d, "blend_shape_keyform_binding_sources_counts", count, func(d *decode.D) { d.FieldS32("count") })
			})

			d.FieldStruct("blend_shape_constraint_indices", func(d *decode.D) {
				d.SeekAbs(sectionOffsets.blendShapeConstraintIndices.blendShapeConstraintSourcesIndices * 8)
				fieldAlignedNArray(d, "blend_shape_constraint_sources_indices", countInfo.blendShapeConstraintIndices, func(d *decode.D) { d.FieldS32("index") })
			})

			d.FieldStruct("blend_shape_constraints", func(d *decode.D) {
				count := countInfo.blendShapeConstraints

				d.SeekAbs(sectionOffsets.blendShapeConstraints.parameterIndices * 8)
				fieldAlignedNArray(d, "parameter_indices", count, func(d *decode.D) { d.FieldS32("index") })

				d.SeekAbs(sectionOffsets.blendShapeConstraints.blendShapeConstraintValueSourcesBeginIndices * 8)
				fieldAlignedNArray(d, "blend_shape_constraint_value_sources_begin_indices", count, func(d *decode.D) { d.FieldS32("index") })

				d.SeekAbs(sectionOffsets.blendShapeConstraints.blendShapeConstraintValueSourcesCounts * 8)
				fieldAlignedNArray(d, "blend_shape_constraint_value_sources_counts", count, func(d *decode.D) { d.FieldS32("count") })
			})

			d.FieldStruct("blend_shape_constraint_values", func(d *decode.D) {
				count := countInfo.blendShapeConstraintValues

				d.SeekAbs(sectionOffsets.blendShapeConstraintValues.keys * 8)
				fieldAlignedNArray(d, "keys", count, func(d *decode.D) { d.FieldF32("value") })

				d.SeekAbs(sectionOffsets.blendShapeConstraintValues.weights * 8)
				fieldAlignedNArray(d, "weights", count, func(d *decode.D) { d.FieldF32("value") })
			})
		}

		if version >= moc3Version5_00_00 {
			d.FieldStruct("blend_shapes_parts", func(d *decode.D) {
				count := countInfo.blendShapesParts

				d.SeekAbs(sectionOffsets.blendShapesParts.targetIndices * 8)
				fieldAlignedNArray(d, "target_indices", count, func(d *decode.D) { d.FieldS32("index") })

				d.SeekAbs(sectionOffsets.blendShapesParts.blendShapeKeyformBindingSourcesBeginIndices * 8)
				fieldAlignedNArray(d, "blend_shape_keyform_binding_sources_begin_indices", count, func(d *decode.D) { d.FieldS32("index") })

				d.SeekAbs(sectionOffsets.blendShapesParts.blendShapeKeyformBindingSourcesCounts * 8)
				fieldAlignedNArray(d, "blend_shape_keyform_binding_sources_counts", count, func(d *decode.D) { d.FieldS32("count") })
			})

			d.FieldStruct("blend_shapes_rotation_deformers", func(d *decode.D) {
				count := countInfo.blendShapesRotationDeformers

				d.SeekAbs(sectionOffsets.blendShapesRotationDeformers.targetIndices * 8)
				fieldAlignedNArray(d, "target_indices", count, func(d *decode.D) { d.FieldS32("index") })

				d.SeekAbs(sectionOffsets.blendShapesRotationDeformers.blendShapeKeyformBindingSourcesBeginIndices * 8)
				fieldAlignedNArray(d, "blend_shape_keyform_binding_sources_begin_indices", count, func(d *decode.D) { d.FieldS32("index") })

				d.SeekAbs(sectionOffsets.blendShapesRotationDeformers.blendShapeKeyformBindingSourcesCounts * 8)
				fieldAlignedNArray(d, "blend_shape_keyform_binding_sources_counts", count, func(d *decode.D) { d.FieldS32("count") })
			})

			d.FieldStruct("blend_shapes_glue", func(d *decode.D) {
				count := countInfo.blendShapesGlue

				d.SeekAbs(sectionOffsets.blendShapesGlue.targetIndices * 8)
				fieldAlignedNArray(d, "target_indices", count, func(d *decode.D) { d.FieldS32("index") })

				d.SeekAbs(sectionOffsets.blendShapesGlue.blendShapeKeyformBindingSourcesBeginIndices * 8)
				fieldAlignedNArray(d, "blend_shape_keyform_binding_sources_begin_indices", count, func(d *decode.D) { d.FieldS32("index") })

				d.SeekAbs(sectionOffsets.blendShapesGlue.blendShapeKeyformBindingSourcesCounts * 8)
				fieldAlignedNArray(d, "blend_shape_keyform_binding_sources_counts", count, func(d *decode.D) { d.FieldS32("count") })
			})
		}
	})

	return nil
}
