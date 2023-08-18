package moc3

// https://github.com/OpenL2D/moc3ingbird/blob/master/src/moc3.hexpat

import (
	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
	"github.com/wader/fq/pkg/scalar"
)

func init() {
	interp.RegisterFormat(
		format.MOC3,
		&decode.Format{
			Description: "MOC3 file",
			Groups:      []*decode.Group{format.Probe},
			DecodeFn:    decodeMOC3,
		})
}

const (
	moc3Version3_00_00 = 1
	moc3Version3_03_00 = 2
	moc3Version4_00_00 = 3
	moc3Version4_02_00 = 4
)

var moc3VersionNames = scalar.UintMap{
	moc3Version3_00_00: {Sym: "3_00_00", Description: "3.0.00 - 3.2.07"},
	moc3Version3_03_00: {Sym: "3_03_00", Description: "3.3.00 - 3.3.03"},
	moc3Version4_00_00: {Sym: "4_00_00", Description: "4.0.00 - 4.1.05"},
	moc3Version4_02_00: {Sym: "4_02_00", Description: "4.2.00 - 4.2.02"},
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
	parts                       int64
	deformers                   int64
	warpDeformers               int64
	rotationDeformers           int64
	artMeshes                   int64
	parameters                  int64
	partKeyforms                int64
	warpDeformerKeyforms        int64
	rotationDeformerKeyforms    int64
	artMeshKeyforms             int64
	keyformPositions            int64
	parameterBindingIndices     int64
	keyformBindings             int64
	parameterBindings           int64
	keys                        int64
	uvs                         int64
	positionIndices             int64
	drawableMasks               int64
	drawOrderGroups             int64
	drawOrderGroupObjects       int64
	glue                        int64
	glueInfo                    int64
	glueKeyforms                int64
	keyformMultiplyColors       int64
	keyformScreenColors         int64
	blendShapeParameterBindings int64
	blendShapeKeyformBindings   int64
	blendShapesWarpDeformers    int64
	blendShapesArtMeshes        int64
	blendShapeConstraintIndices int64
	blendShapeConstraints       int64
	blendShapeConstraintValues  int64
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

	parameterOffsetsV4_2 struct {
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
		blendShapeConstraintIndexSourcesBeginIndices int64
		blendShapeConstraintIndexSourcesCounts       int64
		keyformSourcesBlendShapeIndices              int64
		keyformSourcesBlendShapeCounts               int64
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
}

func decodeMOC3(d *decode.D) any {
	d.FieldUTF8("magic", 4, d.StrAssert("MOC3"))
	version := d.FieldU8("version", moc3VersionNames)
	isBigEndian := d.FieldBoolFn("is_big_endian", func(d *decode.D) bool { return d.U8() != 0 })

	if !isBigEndian {
		d.Endian = decode.LittleEndian
	}
	d.SeekRel(58 * 8)

	var sectionOffsets sectionOffsetTable
	d.FieldStruct("section_offsets", func(d *decode.D) {
		sectionOffsets.countInfo = int64(d.FieldU32("count_info"))
		sectionOffsets.canvasInfo = int64(d.FieldU32("canvas_info"))

		d.FieldStruct("parts", func(d *decode.D) {
			sectionOffsets.parts.runtimeSpace0 = int64(d.FieldU32("runtime_space0"))
			sectionOffsets.parts.ids = int64(d.FieldU32("ids"))
			sectionOffsets.parts.keyformBindingSourcesIndices = int64(d.FieldU32("keyform_binding_sources_indices"))
			sectionOffsets.parts.keyformSourcesBeginIndices = int64(d.FieldU32("keyform_sources_begin_indices"))
			sectionOffsets.parts.keyformSourcesCounts = int64(d.FieldU32("keyform_sources_counts"))
			sectionOffsets.parts.isVisible = int64(d.FieldU32("is_visible"))
			sectionOffsets.parts.isEnabled = int64(d.FieldU32("is_enabled"))
			sectionOffsets.parts.parentPartIndices = int64(d.FieldU32("parent_part_indices"))
		})

		d.FieldStruct("deformers", func(d *decode.D) {
			sectionOffsets.deformers.runtimeSpace0 = int64(d.FieldU32("runtime_space0"))
			sectionOffsets.deformers.ids = int64(d.FieldU32("ids"))
			sectionOffsets.deformers.keyformBindingSourcesIndices = int64(d.FieldU32("keyform_binding_sources_indices"))
			sectionOffsets.deformers.isVisible = int64(d.FieldU32("is_visible"))
			sectionOffsets.deformers.isEnabled = int64(d.FieldU32("is_enabled"))
			sectionOffsets.deformers.parentPartIndices = int64(d.FieldU32("parent_part_indices"))
			sectionOffsets.deformers.parentDeformerIndices = int64(d.FieldU32("parent_deformer_indices"))
			sectionOffsets.deformers.types = int64(d.FieldU32("types"))
			sectionOffsets.deformers.specificSourcesIndices = int64(d.FieldU32("specific_sources_indices"))
		})

		d.FieldStruct("warp_deformers", func(d *decode.D) {
			sectionOffsets.warpDeformers.keyformBindingSourcesIndices = int64(d.FieldU32("keyform_binding_sources_indices"))
			sectionOffsets.warpDeformers.keyformSourcesBeginIndices = int64(d.FieldU32("keyform_sources_begin_indices"))
			sectionOffsets.warpDeformers.keyformSourcesCounts = int64(d.FieldU32("keyform_sources_counts"))
			sectionOffsets.warpDeformers.vertexCounts = int64(d.FieldU32("vertex_counts"))
			sectionOffsets.warpDeformers.rows = int64(d.FieldU32("rows"))
			sectionOffsets.warpDeformers.columns = int64(d.FieldU32("columns"))
		})

		d.FieldStruct("rotation_deformers", func(d *decode.D) {
			sectionOffsets.rotationDeformers.keyformBindingSourcesIndices = int64(d.FieldU32("keyform_binding_sources_indices"))
			sectionOffsets.rotationDeformers.keyformSourcesBeginIndices = int64(d.FieldU32("keyform_sources_begin_indices"))
			sectionOffsets.rotationDeformers.keyformSourcesCounts = int64(d.FieldU32("keyform_sources_counts"))
			sectionOffsets.rotationDeformers.baseAngles = int64(d.FieldU32("base_angles"))
		})

		d.FieldStruct("art_meshes", func(d *decode.D) {
			sectionOffsets.artMeshes.runtimeSpace0 = int64(d.FieldU32("runtime_space0"))
			sectionOffsets.artMeshes.runtimeSpace1 = int64(d.FieldU32("runtime_space1"))
			sectionOffsets.artMeshes.runtimeSpace2 = int64(d.FieldU32("runtime_space2"))
			sectionOffsets.artMeshes.runtimeSpace3 = int64(d.FieldU32("runtime_space3"))
			sectionOffsets.artMeshes.ids = int64(d.FieldU32("ids"))
			sectionOffsets.artMeshes.keyformBindingSourcesIndices = int64(d.FieldU32("keyform_binding_sources_indices"))
			sectionOffsets.artMeshes.keyformSourcesBeginIndices = int64(d.FieldU32("keyform_sources_begin_indices"))
			sectionOffsets.artMeshes.keyformSourcesCounts = int64(d.FieldU32("keyform_sources_counts"))
			sectionOffsets.artMeshes.isVisible = int64(d.FieldU32("is_visible"))
			sectionOffsets.artMeshes.isEnabled = int64(d.FieldU32("is_enabled"))
			sectionOffsets.artMeshes.parentPartIndices = int64(d.FieldU32("parent_part_indices"))
			sectionOffsets.artMeshes.parentDeformerIndices = int64(d.FieldU32("parent_deformer_indices"))
			sectionOffsets.artMeshes.textureNos = int64(d.FieldU32("texture_nos"))
			sectionOffsets.artMeshes.drawableFlags = int64(d.FieldU32("drawable_flags"))
			sectionOffsets.artMeshes.vertexCounts = int64(d.FieldU32("vertex_counts"))
			sectionOffsets.artMeshes.uvSourcesBeginIndices = int64(d.FieldU32("uv_sources_begin_indices"))
			sectionOffsets.artMeshes.positionIndexSourcesBeginIndices = int64(d.FieldU32("position_index_sources_begin_indices"))
			sectionOffsets.artMeshes.positionIndexSourcesCounts = int64(d.FieldU32("position_index_sources_counts"))
			sectionOffsets.artMeshes.drawableMaskSourcesBeginIndices = int64(d.FieldU32("drawable_mask_sources_begin_indices"))
			sectionOffsets.artMeshes.drawableMaskSourcesCounts = int64(d.FieldU32("drawable_mask_sources_counts"))
		})

		d.FieldStruct("parameters", func(d *decode.D) {
			sectionOffsets.parameters.runtimeSpace0 = int64(d.FieldU32("runtime_space0"))
			sectionOffsets.parameters.ids = int64(d.FieldU32("ids"))
			sectionOffsets.parameters.maxValues = int64(d.FieldU32("max_values"))
			sectionOffsets.parameters.minValues = int64(d.FieldU32("min_values"))
			sectionOffsets.parameters.defaultValues = int64(d.FieldU32("default_values"))
			sectionOffsets.parameters.isRepeat = int64(d.FieldU32("is_repeat"))
			sectionOffsets.parameters.decimalPlaces = int64(d.FieldU32("decimal_places"))
			sectionOffsets.parameters.parameterBindingSourcesBeginIndices = int64(d.FieldU32("parameter_binding_sources_begin_indices"))
			sectionOffsets.parameters.parameterBindingSourcesCounts = int64(d.FieldU32("parameter_binding_sources_counts"))
		})

		d.FieldStruct("part_keyforms", func(d *decode.D) {
			sectionOffsets.partKeyforms.drawOrders = int64(d.FieldU32("draw_orders"))
		})

		d.FieldStruct("warp_deformer_keyforms", func(d *decode.D) {
			sectionOffsets.warpDeformerKeyforms.opacities = int64(d.FieldU32("opacities"))
			sectionOffsets.warpDeformerKeyforms.keyformPositionSourcesBeginIndices = int64(d.FieldU32("keyform_position_sources_begin_indices"))
		})

		d.FieldStruct("rotation_deformer_keyforms", func(d *decode.D) {
			sectionOffsets.rotationDeformerKeyforms.opacities = int64(d.FieldU32("opacities"))
			sectionOffsets.rotationDeformerKeyforms.angles = int64(d.FieldU32("angles"))
			sectionOffsets.rotationDeformerKeyforms.originX = int64(d.FieldU32("origin_x"))
			sectionOffsets.rotationDeformerKeyforms.originY = int64(d.FieldU32("origin_y"))
			sectionOffsets.rotationDeformerKeyforms.scales = int64(d.FieldU32("scales"))
			sectionOffsets.rotationDeformerKeyforms.isReflectX = int64(d.FieldU32("is_reflect_x"))
			sectionOffsets.rotationDeformerKeyforms.isReflectY = int64(d.FieldU32("is_reflect_y"))
		})

		d.FieldStruct("art_mesh_keyforms", func(d *decode.D) {
			sectionOffsets.artMeshKeyforms.opacities = int64(d.FieldU32("opacities"))
			sectionOffsets.artMeshKeyforms.drawOrders = int64(d.FieldU32("draw_orders"))
			sectionOffsets.artMeshKeyforms.keyformPositionSourcesBeginIndices = int64(d.FieldU32("keyform_position_sources_begin_indices"))
		})

		d.FieldStruct("keyform_positions", func(d *decode.D) {
			sectionOffsets.keyformPositions.xys = int64(d.FieldU32("xys"))
		})

		d.FieldStruct("parameter_binding_indices", func(d *decode.D) {
			sectionOffsets.parameterBindingIndices.bindingSourcesIndices = int64(d.FieldU32("binding_sources_indices"))
		})

		d.FieldStruct("keyform_bindings", func(d *decode.D) {
			sectionOffsets.keyformBindings.parameterBindingIndexSourcesBeginIndices = int64(d.FieldU32("parameter_binding_index_sources_begin_indices"))
			sectionOffsets.keyformBindings.parameterBindingIndexSourcesCounts = int64(d.FieldU32("parameter_binding_index_sources_counts"))
		})

		d.FieldStruct("parameter_bindings", func(d *decode.D) {
			sectionOffsets.parameterBindings.keysSourcesBeginIndices = int64(d.FieldU32("keys_sources_begin_indices"))
			sectionOffsets.parameterBindings.keysSourcesCounts = int64(d.FieldU32("keys_sources_counts"))
		})

		d.FieldStruct("keys", func(d *decode.D) {
			sectionOffsets.keys.values = int64(d.FieldU32("values"))
		})

		d.FieldStruct("uvs", func(d *decode.D) {
			sectionOffsets.UVs.uvs = int64(d.FieldU32("uvs"))
		})

		d.FieldStruct("position_indices", func(d *decode.D) {
			sectionOffsets.positionIndices.indices = int64(d.FieldU32("indices"))
		})

		d.FieldStruct("drawable_masks", func(d *decode.D) {
			sectionOffsets.drawableMasks.artMeshSourcesIndices = int64(d.FieldU32("art_mesh_sources_indices"))
		})

		d.FieldStruct("draw_order_groups", func(d *decode.D) {
			sectionOffsets.drawOrderGroups.objectSourcesBeginIndices = int64(d.FieldU32("object_sources_begin_indices"))
			sectionOffsets.drawOrderGroups.objectSourcesCounts = int64(d.FieldU32("object_sources_counts"))
			sectionOffsets.drawOrderGroups.objectSourcesTotalCounts = int64(d.FieldU32("object_sources_total_counts"))
			sectionOffsets.drawOrderGroups.maximumDrawOrders = int64(d.FieldU32("maximum_draw_orders"))
			sectionOffsets.drawOrderGroups.minimumDrawOrders = int64(d.FieldU32("minimum_draw_orders"))
		})

		d.FieldStruct("draw_order_group_objects", func(d *decode.D) {
			sectionOffsets.drawOrderGroupObjects.types = int64(d.FieldU32("types"))
			sectionOffsets.drawOrderGroupObjects.indices = int64(d.FieldU32("indices"))
			sectionOffsets.drawOrderGroupObjects.selfIndices = int64(d.FieldU32("self_indices"))
		})

		d.FieldStruct("glue", func(d *decode.D) {
			sectionOffsets.glue.runtimeSpace0 = int64(d.FieldU32("runtime_space0"))
			sectionOffsets.glue.ids = int64(d.FieldU32("ids"))
			sectionOffsets.glue.keyformBindingSourcesIndices = int64(d.FieldU32("keyform_binding_sources_indices"))
			sectionOffsets.glue.keyformSourcesBeginIndices = int64(d.FieldU32("keyform_sources_begin_indices"))
			sectionOffsets.glue.keyformSourcesCounts = int64(d.FieldU32("keyform_sources_counts"))
			sectionOffsets.glue.artMeshIndicesA = int64(d.FieldU32("art_mesh_indices_a"))
			sectionOffsets.glue.artMeshIndicesB = int64(d.FieldU32("art_mesh_indices_b"))
			sectionOffsets.glue.glueInfoSourcesBeginIndices = int64(d.FieldU32("glue_info_sources_begin_indices"))
			sectionOffsets.glue.glueInfoSourcesCounts = int64(d.FieldU32("glue_info_sources_counts"))
		})

		d.FieldStruct("glue_info", func(d *decode.D) {
			sectionOffsets.glueInfo.weights = int64(d.FieldU32("weights"))
			sectionOffsets.glueInfo.positionIndices = int64(d.FieldU32("position_indices"))
		})

		if version < moc3Version3_03_00 {
			return
		}

		d.FieldStruct("glue_keyforms", func(d *decode.D) {
			sectionOffsets.glueKeyforms.intensities = int64(d.FieldU32("intensities"))
		})

		if version < moc3Version4_02_00 {
			return
		}

		d.FieldStruct("warp_deformers_v3_3", func(d *decode.D) {
			sectionOffsets.warpDeformersV3_3.isQuadSource = int64(d.FieldU32("is_quad_source"))
		})

		d.FieldStruct("parameter_extensions", func(d *decode.D) {
			sectionOffsets.parameterExtensions.runtimeSpace0 = int64(d.FieldU32("runtime_space0"))
			sectionOffsets.parameterExtensions.keysSourcesBeginIndices = int64(d.FieldU32("keys_sources_begin_indices"))
			sectionOffsets.parameterExtensions.keysSourcesCounts = int64(d.FieldU32("keys_sources_counts"))
		})

		d.FieldStruct("warp_deformers_v4_2", func(d *decode.D) {
			sectionOffsets.warpDeformersV4_2.keyformColorSourcesBeginIndices = int64(d.FieldU32("keyform_color_sources_begin_indices"))
		})

		d.FieldStruct("rotation_deformers_v4_2", func(d *decode.D) {
			sectionOffsets.rotationDeformersV4_2.keyformColorSourcesBeginIndices = int64(d.FieldU32("keyform_color_sources_begin_indices"))
		})

		d.FieldStruct("art_meshes_v4_2", func(d *decode.D) {
			sectionOffsets.artMeshesV4_2.keyformColorSourcesBeginIndices = int64(d.FieldU32("keyform_color_sources_begin_indices"))
		})

		d.FieldStruct("keyform_multiply_colors", func(d *decode.D) {
			sectionOffsets.keyformMultiplyColors.r = int64(d.FieldU32("r"))
			sectionOffsets.keyformMultiplyColors.g = int64(d.FieldU32("g"))
			sectionOffsets.keyformMultiplyColors.b = int64(d.FieldU32("b"))
		})

		d.FieldStruct("keyform_screen_colors", func(d *decode.D) {
			sectionOffsets.keyformScreenColors.r = int64(d.FieldU32("r"))
			sectionOffsets.keyformScreenColors.g = int64(d.FieldU32("g"))
			sectionOffsets.keyformScreenColors.b = int64(d.FieldU32("b"))
		})

		d.FieldStruct("parameter_offsets_v4_2", func(d *decode.D) {
			sectionOffsets.parameterOffsetsV4_2.parameterTypes = int64(d.FieldU32("parameter_types"))
			sectionOffsets.parameterOffsetsV4_2.blendShapeParameterBindingSourcesBeginIndices = int64(d.FieldU32("blend_shape_parameter_binding_sources_begin_indices"))
			sectionOffsets.parameterOffsetsV4_2.blendShapeParameterBindingSourcesCounts = int64(d.FieldU32("blend_shape_parameter_binding_sources_counts"))
		})

		d.FieldStruct("blend_shape_parameter_bindings", func(d *decode.D) {
			sectionOffsets.blendShapeParameterBindings.keysSourcesBeginIndices = int64(d.FieldU32("keys_sources_begin_indices"))
			sectionOffsets.blendShapeParameterBindings.keysSourcesCounts = int64(d.FieldU32("keys_sources_counts"))
			sectionOffsets.blendShapeParameterBindings.baseKeyIndices = int64(d.FieldU32("base_key_indices"))
		})

		d.FieldStruct("blend_shape_keyform_bindings", func(d *decode.D) {
			sectionOffsets.blendShapeKeyformBindings.parameterBindingSourcesIndices = int64(d.FieldU32("parameter_binding_sources_indices"))
			sectionOffsets.blendShapeKeyformBindings.blendShapeConstraintIndexSourcesBeginIndices = int64(d.FieldU32("blend_shape_constraint_index_sources_begin_indices"))
			sectionOffsets.blendShapeKeyformBindings.blendShapeConstraintIndexSourcesCounts = int64(d.FieldU32("blend_shape_constraint_index_sources_counts"))
			sectionOffsets.blendShapeKeyformBindings.keyformSourcesBlendShapeIndices = int64(d.FieldU32("keyform_sources_blend_shape_indices"))
			sectionOffsets.blendShapeKeyformBindings.keyformSourcesBlendShapeCounts = int64(d.FieldU32("keyform_sources_blend_shape_counts"))
		})

		d.FieldStruct("blend_shapes_warp_deformers", func(d *decode.D) {
			sectionOffsets.blendShapesWarpDeformers.targetIndices = int64(d.FieldU32("target_indices"))
			sectionOffsets.blendShapesWarpDeformers.blendShapeKeyformBindingSourcesBeginIndices = int64(d.FieldU32("blend_shape_keyform_binding_sources_begin_indices"))
			sectionOffsets.blendShapesWarpDeformers.blendShapeKeyformBindingSourcesCounts = int64(d.FieldU32("blend_shape_keyform_binding_sources_counts"))
		})

		d.FieldStruct("blend_shapes_art_meshes", func(d *decode.D) {
			sectionOffsets.blendShapesArtMeshes.targetIndices = int64(d.FieldU32("target_indices"))
			sectionOffsets.blendShapesArtMeshes.blendShapeKeyformBindingSourcesBeginIndices = int64(d.FieldU32("blend_shape_keyform_binding_sources_begin_indices"))
			sectionOffsets.blendShapesArtMeshes.blendShapeKeyformBindingSourcesCounts = int64(d.FieldU32("blend_shape_keyform_binding_sources_counts"))
		})

		d.FieldStruct("blend_shape_constraint_indices", func(d *decode.D) {
			sectionOffsets.blendShapeConstraintIndices.blendShapeConstraintSourcesIndices = int64(d.FieldU32("blend_shape_constraint_sources_indices"))
		})

		d.FieldStruct("blend_shape_constraints", func(d *decode.D) {
			sectionOffsets.blendShapeConstraints.parameterIndices = int64(d.FieldU32("parameter_indices"))
			sectionOffsets.blendShapeConstraints.blendShapeConstraintValueSourcesBeginIndices = int64(d.FieldU32("blend_shape_constraint_value_sources_begin_indices"))
			sectionOffsets.blendShapeConstraints.blendShapeConstraintValueSourcesCounts = int64(d.FieldU32("blend_shape_constraint_value_sources_counts"))
		})

		d.FieldStruct("blend_shape_constraint_values", func(d *decode.D) {
			sectionOffsets.blendShapeConstraintValues.keys = int64(d.FieldU32("keys"))
			sectionOffsets.blendShapeConstraintValues.weights = int64(d.FieldU32("weights"))
		})
	})

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

			if version < moc3Version4_02_00 {
				return
			}

			countInfo.glueKeyforms = int64(d.FieldU32("glue_keyforms"))
			countInfo.keyformMultiplyColors = int64(d.FieldU32("keyform_multiply_colors"))
			countInfo.keyformScreenColors = int64(d.FieldU32("keyform_screen_colors"))
			countInfo.blendShapeParameterBindings = int64(d.FieldU32("blend_shape_parameter_bindings"))
			countInfo.blendShapeKeyformBindings = int64(d.FieldU32("blend_shape_keyform_bindings"))
			countInfo.blendShapesWarpDeformers = int64(d.FieldU32("blend_shapes_warp_deformers"))
			countInfo.blendShapesArtMeshes = int64(d.FieldU32("blend_shapes_art_meshes"))
			countInfo.blendShapeConstraintIndices = int64(d.FieldU32("blend_shape_constraint_indices"))
			countInfo.blendShapeConstraints = int64(d.FieldU32("blend_shape_constraints"))
			countInfo.blendShapeConstraintValues = int64(d.FieldU32("blend_shape_constraint_values"))
		})

		d.SeekAbs(sectionOffsets.canvasInfo * 8)
		d.FieldStruct("canvas_info", func(d *decode.D) {
			d.FieldF32("pixels_per_unit")
			d.FieldF32("origin_x")
			d.FieldF32("origin_y")
			d.FieldF32("canvas_width")
			d.FieldF32("canvas_height")
			d.FieldU8("canvas_flags")
		})

		d.FieldStruct("parts", func(d *decode.D) {
			d.SeekAbs(sectionOffsets.parts.ids * 8)
			d.FieldArray("ids", func(d *decode.D) {
				for i := int64(0); i < countInfo.parts; i++ {
					d.FieldUTF8NullFixedLen("id", 64)
				}
			})

			d.SeekAbs(sectionOffsets.parts.keyformBindingSourcesIndices * 8)
			d.FieldArray("keyform_binding_sources_indices", func(d *decode.D) {
				for i := int64(0); i < countInfo.parts; i++ {
					d.FieldS32("element")
				}
			})

			d.SeekAbs(sectionOffsets.parts.keyformSourcesBeginIndices * 8)
			d.FieldArray("keyform_sources_begin_indices", func(d *decode.D) {
				for i := int64(0); i < countInfo.parts; i++ {
					d.FieldS32("element")
				}
			})

			d.SeekAbs(sectionOffsets.parts.keyformSourcesCounts * 8)
			d.FieldArray("keyform_sources_counts", func(d *decode.D) {
				for i := int64(0); i < countInfo.parts; i++ {
					d.FieldS32("element")
				}
			})

			d.SeekAbs(sectionOffsets.parts.isVisible * 8)
			d.FieldArray("is_visible", func(d *decode.D) {
				for i := int64(0); i < countInfo.parts; i++ {
					d.FieldBoolFn("element", func(d *decode.D) bool { return d.U32() != 0 })
				}
			})

			d.SeekAbs(sectionOffsets.parts.isEnabled * 8)
			d.FieldArray("is_enabled", func(d *decode.D) {
				for i := int64(0); i < countInfo.parts; i++ {
					d.FieldBoolFn("element", func(d *decode.D) bool { return d.U32() != 0 })
				}
			})

			d.SeekAbs(sectionOffsets.parts.parentPartIndices * 8)
			d.FieldArray("parent_part_indices", func(d *decode.D) {
				for i := int64(0); i < countInfo.parts; i++ {
					d.FieldS32("element")
				}
			})
		})

		d.FieldStruct("deformers", func(d *decode.D) {
			d.SeekAbs(sectionOffsets.deformers.ids * 8)
			d.FieldArray("ids", func(d *decode.D) {
				for i := int64(0); i < countInfo.deformers; i++ {
					d.FieldUTF8NullFixedLen("id", 64)
				}
			})

			d.SeekAbs(sectionOffsets.deformers.keyformBindingSourcesIndices * 8)
			d.FieldArray("keyform_binding_sources_indices", func(d *decode.D) {
				for i := int64(0); i < countInfo.deformers; i++ {
					d.FieldS32("element")
				}
			})

			d.SeekAbs(sectionOffsets.deformers.isVisible * 8)
			d.FieldArray("is_visible", func(d *decode.D) {
				for i := int64(0); i < countInfo.deformers; i++ {
					d.FieldBoolFn("element", func(d *decode.D) bool { return d.U32() != 0 })
				}
			})

			d.SeekAbs(sectionOffsets.deformers.isEnabled * 8)
			d.FieldArray("is_enabled", func(d *decode.D) {
				for i := int64(0); i < countInfo.deformers; i++ {
					d.FieldBoolFn("element", func(d *decode.D) bool { return d.U32() != 0 })
				}
			})

			d.SeekAbs(sectionOffsets.deformers.parentPartIndices * 8)
			d.FieldArray("parent_part_indices", func(d *decode.D) {
				for i := int64(0); i < countInfo.deformers; i++ {
					d.FieldS32("element")
				}
			})

			d.SeekAbs(sectionOffsets.deformers.parentDeformerIndices * 8)
			d.FieldArray("parent_deformer_indices", func(d *decode.D) {
				for i := int64(0); i < countInfo.deformers; i++ {
					d.FieldS32("element")
				}
			})

			d.SeekAbs(sectionOffsets.deformers.types * 8)
			d.FieldArray("types", func(d *decode.D) {
				for i := int64(0); i < countInfo.deformers; i++ {
					d.FieldU32("type", deformerTypeNames)
				}
			})

			d.SeekAbs(sectionOffsets.deformers.specificSourcesIndices * 8)
			d.FieldArray("specific_sources_indices", func(d *decode.D) {
				for i := int64(0); i < countInfo.deformers; i++ {
					d.FieldS32("element")
				}
			})
		})

		d.FieldStruct("warp_deformers", func(d *decode.D) {
			d.SeekAbs(sectionOffsets.warpDeformers.keyformBindingSourcesIndices * 8)
			d.FieldArray("keyform_binding_sources_indices", func(d *decode.D) {
				for i := int64(0); i < countInfo.warpDeformers; i++ {
					d.FieldS32("element")
				}
			})

			d.SeekAbs(sectionOffsets.warpDeformers.keyformSourcesBeginIndices * 8)
			d.FieldArray("keyform_sources_begin_indices", func(d *decode.D) {
				for i := int64(0); i < countInfo.warpDeformers; i++ {
					d.FieldS32("element")
				}
			})

			d.SeekAbs(sectionOffsets.warpDeformers.keyformSourcesCounts * 8)
			d.FieldArray("keyform_sources_counts", func(d *decode.D) {
				for i := int64(0); i < countInfo.warpDeformers; i++ {
					d.FieldS32("element")
				}
			})

			d.SeekAbs(sectionOffsets.warpDeformers.vertexCounts * 8)
			d.FieldArray("vertex_counts", func(d *decode.D) {
				for i := int64(0); i < countInfo.warpDeformers; i++ {
					d.FieldS32("element")
				}
			})

			d.SeekAbs(sectionOffsets.warpDeformers.rows * 8)
			d.FieldArray("rows", func(d *decode.D) {
				for i := int64(0); i < countInfo.warpDeformers; i++ {
					d.FieldU32("element")
				}
			})

			d.SeekAbs(sectionOffsets.warpDeformers.columns * 8)
			d.FieldArray("columns", func(d *decode.D) {
				for i := int64(0); i < countInfo.warpDeformers; i++ {
					d.FieldU32("element")
				}
			})
		})

		d.FieldStruct("rotation_deformers", func(d *decode.D) {
			d.SeekAbs(sectionOffsets.rotationDeformers.keyformBindingSourcesIndices * 8)
			d.FieldArray("keyform_binding_sources_indices", func(d *decode.D) {
				for i := int64(0); i < countInfo.rotationDeformers; i++ {
					d.FieldS32("element")
				}
			})

			d.SeekAbs(sectionOffsets.rotationDeformers.keyformSourcesBeginIndices * 8)
			d.FieldArray("keyform_sources_begin_indices", func(d *decode.D) {
				for i := int64(0); i < countInfo.rotationDeformers; i++ {
					d.FieldS32("element")
				}
			})

			d.SeekAbs(sectionOffsets.rotationDeformers.keyformSourcesCounts * 8)
			d.FieldArray("keyform_sources_counts", func(d *decode.D) {
				for i := int64(0); i < countInfo.rotationDeformers; i++ {
					d.FieldS32("element")
				}
			})

			d.SeekAbs(sectionOffsets.rotationDeformers.baseAngles * 8)
			d.FieldArray("base_angles", func(d *decode.D) {
				for i := int64(0); i < countInfo.rotationDeformers; i++ {
					d.FieldF32("element")
				}
			})
		})

		d.FieldStruct("art_meshes", func(d *decode.D) {
			d.SeekAbs(sectionOffsets.artMeshes.ids * 8)
			d.FieldArray("ids", func(d *decode.D) {
				for i := int64(0); i < countInfo.artMeshes; i++ {
					d.FieldUTF8NullFixedLen("id", 64)
				}
			})

			d.SeekAbs(sectionOffsets.artMeshes.keyformBindingSourcesIndices * 8)
			d.FieldArray("keyform_binding_sources_indices", func(d *decode.D) {
				for i := int64(0); i < countInfo.artMeshes; i++ {
					d.FieldS32("element")
				}
			})

			d.SeekAbs(sectionOffsets.artMeshes.keyformSourcesBeginIndices * 8)
			d.FieldArray("keyform_sources_begin_indices", func(d *decode.D) {
				for i := int64(0); i < countInfo.artMeshes; i++ {
					d.FieldS32("element")
				}
			})

			d.SeekAbs(sectionOffsets.artMeshes.keyformSourcesCounts * 8)
			d.FieldArray("keyform_sources_counts", func(d *decode.D) {
				for i := int64(0); i < countInfo.artMeshes; i++ {
					d.FieldS32("element")
				}
			})

			d.SeekAbs(sectionOffsets.artMeshes.isVisible * 8)
			d.FieldArray("is_visible", func(d *decode.D) {
				for i := int64(0); i < countInfo.artMeshes; i++ {
					d.FieldBoolFn("element", func(d *decode.D) bool { return d.U32() != 0 })
				}
			})

			d.SeekAbs(sectionOffsets.artMeshes.isEnabled * 8)
			d.FieldArray("is_enabled", func(d *decode.D) {
				for i := int64(0); i < countInfo.artMeshes; i++ {
					d.FieldBoolFn("element", func(d *decode.D) bool { return d.U32() != 0 })
				}
			})

			d.SeekAbs(sectionOffsets.artMeshes.parentPartIndices * 8)
			d.FieldArray("parent_part_indices", func(d *decode.D) {
				for i := int64(0); i < countInfo.artMeshes; i++ {
					d.FieldS32("element")
				}
			})

			d.SeekAbs(sectionOffsets.artMeshes.parentDeformerIndices * 8)
			d.FieldArray("parent_deformer_indices", func(d *decode.D) {
				for i := int64(0); i < countInfo.artMeshes; i++ {
					d.FieldS32("element")
				}
			})

			d.SeekAbs(sectionOffsets.artMeshes.textureNos * 8)
			d.FieldArray("texture_nos", func(d *decode.D) {
				for i := int64(0); i < countInfo.artMeshes; i++ {
					d.FieldU32("element")
				}
			})

			d.SeekAbs(sectionOffsets.artMeshes.drawableFlags * 8)
			d.FieldArray("drawable_flags", func(d *decode.D) {
				for i := int64(0); i < countInfo.artMeshes; i++ {
					d.FieldStruct("flags", func(d *decode.D) {
						d.FieldU4("reserved")
						d.FieldBool("is_inverted")
						d.FieldBool("is_double_sided")
						d.FieldU2("blend_mode", blendModeNames)
					})
				}
			})

			d.SeekAbs(sectionOffsets.artMeshes.vertexCounts * 8)
			d.FieldArray("vertex_counts", func(d *decode.D) {
				for i := int64(0); i < countInfo.artMeshes; i++ {
					d.FieldS32("element")
				}
			})

			d.SeekAbs(sectionOffsets.artMeshes.uvSourcesBeginIndices * 8)
			d.FieldArray("uv_sources_begin_indices", func(d *decode.D) {
				for i := int64(0); i < countInfo.artMeshes; i++ {
					d.FieldS32("element")
				}
			})

			d.SeekAbs(sectionOffsets.artMeshes.positionIndexSourcesBeginIndices * 8)
			d.FieldArray("position_index_sources_begin_indices", func(d *decode.D) {
				for i := int64(0); i < countInfo.artMeshes; i++ {
					d.FieldS32("element")
				}
			})

			d.SeekAbs(sectionOffsets.artMeshes.positionIndexSourcesCounts * 8)
			d.FieldArray("position_index_sources_counts", func(d *decode.D) {
				for i := int64(0); i < countInfo.artMeshes; i++ {
					d.FieldS32("element")
				}
			})

			d.SeekAbs(sectionOffsets.artMeshes.drawableMaskSourcesBeginIndices * 8)
			d.FieldArray("drawable_mask_sources_begin_indices", func(d *decode.D) {
				for i := int64(0); i < countInfo.artMeshes; i++ {
					d.FieldS32("element")
				}
			})

			d.SeekAbs(sectionOffsets.artMeshes.drawableMaskSourcesCounts * 8)
			d.FieldArray("drawable_mask_sources_counts", func(d *decode.D) {
				for i := int64(0); i < countInfo.artMeshes; i++ {
					d.FieldS32("element")
				}
			})
		})

		d.FieldStruct("parameters", func(d *decode.D) {
			d.SeekAbs(sectionOffsets.parameters.ids * 8)
			d.FieldArray("ids", func(d *decode.D) {
				for i := int64(0); i < countInfo.parameters; i++ {
					d.FieldUTF8NullFixedLen("id", 64)
				}
			})

			d.SeekAbs(sectionOffsets.parameters.maxValues * 8)
			d.FieldArray("max_values", func(d *decode.D) {
				for i := int64(0); i < countInfo.parameters; i++ {
					d.FieldF32("element")
				}
			})

			d.SeekAbs(sectionOffsets.parameters.minValues * 8)
			d.FieldArray("min_values", func(d *decode.D) {
				for i := int64(0); i < countInfo.parameters; i++ {
					d.FieldF32("element")
				}
			})

			d.SeekAbs(sectionOffsets.parameters.defaultValues * 8)
			d.FieldArray("default_values", func(d *decode.D) {
				for i := int64(0); i < countInfo.parameters; i++ {
					d.FieldF32("element")
				}
			})

			d.SeekAbs(sectionOffsets.parameters.isRepeat * 8)
			d.FieldArray("is_repeat", func(d *decode.D) {
				for i := int64(0); i < countInfo.parameters; i++ {
					d.FieldBoolFn("element", func(d *decode.D) bool { return d.U32() != 0 })
				}
			})

			d.SeekAbs(sectionOffsets.parameters.decimalPlaces * 8)
			d.FieldArray("decimal_places", func(d *decode.D) {
				for i := int64(0); i < countInfo.parameters; i++ {
					d.FieldU32("element")
				}
			})

			d.SeekAbs(sectionOffsets.parameters.parameterBindingSourcesBeginIndices * 8)
			d.FieldArray("parameter_binding_sources_begin_indices", func(d *decode.D) {
				for i := int64(0); i < countInfo.parameters; i++ {
					d.FieldS32("element")
				}
			})

			d.SeekAbs(sectionOffsets.parameters.parameterBindingSourcesCounts * 8)
			d.FieldArray("parameter_binding_sources_counts", func(d *decode.D) {
				for i := int64(0); i < countInfo.parameters; i++ {
					d.FieldS32("element")
				}
			})
		})

		d.FieldStruct("part_keyforms", func(d *decode.D) {
			d.SeekAbs(sectionOffsets.partKeyforms.drawOrders * 8)
			d.FieldArray("draw_orders", func(d *decode.D) {
				for i := int64(0); i < countInfo.partKeyforms; i++ {
					d.FieldF32("element")
				}
			})
		})

		d.FieldStruct("warp_deformer_keyforms", func(d *decode.D) {
			d.SeekAbs(sectionOffsets.warpDeformerKeyforms.opacities * 8)
			d.FieldArray("opacities", func(d *decode.D) {
				for i := int64(0); i < countInfo.warpDeformerKeyforms; i++ {
					d.FieldF32("element")
				}
			})

			d.SeekAbs(sectionOffsets.warpDeformerKeyforms.keyformPositionSourcesBeginIndices * 8)
			d.FieldArray("keyform_position_sources_begin_indices", func(d *decode.D) {
				for i := int64(0); i < countInfo.warpDeformerKeyforms; i++ {
					d.FieldS32("element")
				}
			})
		})

		d.FieldStruct("rotation_deformer_keyforms", func(d *decode.D) {
			d.SeekAbs(sectionOffsets.rotationDeformerKeyforms.opacities * 8)
			d.FieldArray("opacities", func(d *decode.D) {
				for i := int64(0); i < countInfo.rotationDeformerKeyforms; i++ {
					d.FieldF32("element")
				}
			})

			d.SeekAbs(sectionOffsets.rotationDeformerKeyforms.angles * 8)
			d.FieldArray("angles", func(d *decode.D) {
				for i := int64(0); i < countInfo.rotationDeformerKeyforms; i++ {
					d.FieldF32("element")
				}
			})

			d.SeekAbs(sectionOffsets.rotationDeformerKeyforms.originX * 8)
			d.FieldArray("origin_x", func(d *decode.D) {
				for i := int64(0); i < countInfo.rotationDeformerKeyforms; i++ {
					d.FieldF32("element")
				}
			})

			d.SeekAbs(sectionOffsets.rotationDeformerKeyforms.originY * 8)
			d.FieldArray("origin_y", func(d *decode.D) {
				for i := int64(0); i < countInfo.rotationDeformerKeyforms; i++ {
					d.FieldF32("element")
				}
			})

			d.SeekAbs(sectionOffsets.rotationDeformerKeyforms.scales * 8)
			d.FieldArray("scales", func(d *decode.D) {
				for i := int64(0); i < countInfo.rotationDeformerKeyforms; i++ {
					d.FieldF32("element")
				}
			})

			d.SeekAbs(sectionOffsets.rotationDeformerKeyforms.isReflectX * 8)
			d.FieldArray("is_reflect_x", func(d *decode.D) {
				for i := int64(0); i < countInfo.rotationDeformerKeyforms; i++ {
					d.FieldBoolFn("element", func(d *decode.D) bool { return d.U32() != 0 })
				}
			})

			d.SeekAbs(sectionOffsets.rotationDeformerKeyforms.isReflectY * 8)
			d.FieldArray("is_reflect_y", func(d *decode.D) {
				for i := int64(0); i < countInfo.rotationDeformerKeyforms; i++ {
					d.FieldBoolFn("element", func(d *decode.D) bool { return d.U32() != 0 })
				}
			})
		})

		d.FieldStruct("art_mesh_keyforms", func(d *decode.D) {
			d.SeekAbs(sectionOffsets.artMeshKeyforms.opacities * 8)
			d.FieldArray("opacities", func(d *decode.D) {
				for i := int64(0); i < countInfo.artMeshKeyforms; i++ {
					d.FieldF32("element")
				}
			})

			d.SeekAbs(sectionOffsets.artMeshKeyforms.drawOrders * 8)
			d.FieldArray("draw_orders", func(d *decode.D) {
				for i := int64(0); i < countInfo.artMeshKeyforms; i++ {
					d.FieldF32("element")
				}
			})

			d.SeekAbs(sectionOffsets.artMeshKeyforms.keyformPositionSourcesBeginIndices * 8)
			d.FieldArray("keyform_position_sources_begin_indices", func(d *decode.D) {
				for i := int64(0); i < countInfo.artMeshKeyforms; i++ {
					d.FieldS32("element")
				}
			})
		})

		d.FieldStruct("keyform_positions", func(d *decode.D) {
			d.SeekAbs(sectionOffsets.keyformPositions.xys * 8)
			d.FieldArray("xys", func(d *decode.D) {
				for i := int64(0); i < countInfo.keyformPositions; i++ {
					d.FieldF32("element")
				}
			})
		})

		d.FieldStruct("parameter_binding_indices", func(d *decode.D) {
			d.SeekAbs(sectionOffsets.parameterBindingIndices.bindingSourcesIndices * 8)
			d.FieldArray("binding_sources_indices", func(d *decode.D) {
				for i := int64(0); i < countInfo.parameterBindingIndices; i++ {
					d.FieldS32("element")
				}
			})
		})

		d.FieldStruct("keyform_bindings", func(d *decode.D) {
			d.SeekAbs(sectionOffsets.keyformBindings.parameterBindingIndexSourcesBeginIndices * 8)
			d.FieldArray("parameter_binding_index_sources_begin_indices", func(d *decode.D) {
				for i := int64(0); i < countInfo.keyformBindings; i++ {
					d.FieldS32("element")
				}
			})

			d.SeekAbs(sectionOffsets.keyformBindings.parameterBindingIndexSourcesCounts * 8)
			d.FieldArray("parameter_binding_index_sources_counts", func(d *decode.D) {
				for i := int64(0); i < countInfo.keyformBindings; i++ {
					d.FieldS32("element")
				}
			})
		})

		d.FieldStruct("parameter_bindings", func(d *decode.D) {
			d.SeekAbs(sectionOffsets.parameterBindings.keysSourcesBeginIndices * 8)
			d.FieldArray("keys_sources_begin_indices", func(d *decode.D) {
				for i := int64(0); i < countInfo.parameterBindings; i++ {
					d.FieldS32("element")
				}
			})

			d.SeekAbs(sectionOffsets.parameterBindings.keysSourcesCounts * 8)
			d.FieldArray("keys_sources_counts", func(d *decode.D) {
				for i := int64(0); i < countInfo.parameterBindings; i++ {
					d.FieldS32("element")
				}
			})
		})

		d.FieldStruct("keys", func(d *decode.D) {
			d.SeekAbs(sectionOffsets.keys.values * 8)
			d.FieldArray("values", func(d *decode.D) {
				for i := int64(0); i < countInfo.keys; i++ {
					d.FieldF32("element")
				}
			})
		})

		d.FieldStruct("uvs", func(d *decode.D) {
			d.SeekAbs(sectionOffsets.UVs.uvs * 8)
			d.FieldArray("uvs", func(d *decode.D) {
				for i := int64(0); i < countInfo.uvs; i++ {
					d.FieldF32("element")
				}
			})
		})

		d.FieldStruct("position_indices", func(d *decode.D) {
			d.SeekAbs(sectionOffsets.positionIndices.indices * 8)
			d.FieldArray("indices", func(d *decode.D) {
				for i := int64(0); i < countInfo.positionIndices; i++ {
					d.FieldS16("element")
				}
			})
		})

		d.FieldStruct("drawable_masks", func(d *decode.D) {
			d.SeekAbs(sectionOffsets.drawableMasks.artMeshSourcesIndices * 8)
			d.FieldArray("art_mesh_sources_indices", func(d *decode.D) {
				for i := int64(0); i < countInfo.drawableMasks; i++ {
					d.FieldS32("element")
				}
			})
		})

		d.FieldStruct("draw_order_groups", func(d *decode.D) {
			d.SeekAbs(sectionOffsets.drawOrderGroups.objectSourcesBeginIndices * 8)
			d.FieldArray("object_sources_begin_indices", func(d *decode.D) {
				for i := int64(0); i < countInfo.drawOrderGroups; i++ {
					d.FieldS32("element")
				}
			})

			d.SeekAbs(sectionOffsets.drawOrderGroups.objectSourcesCounts * 8)
			d.FieldArray("object_sources_counts", func(d *decode.D) {
				for i := int64(0); i < countInfo.drawOrderGroups; i++ {
					d.FieldS32("element")
				}
			})

			d.SeekAbs(sectionOffsets.drawOrderGroups.objectSourcesTotalCounts * 8)
			d.FieldArray("object_sources_total_counts", func(d *decode.D) {
				for i := int64(0); i < countInfo.drawOrderGroups; i++ {
					d.FieldS32("element")
				}
			})

			d.SeekAbs(sectionOffsets.drawOrderGroups.maximumDrawOrders * 8)
			d.FieldArray("maximum_draw_orders", func(d *decode.D) {
				for i := int64(0); i < countInfo.drawOrderGroups; i++ {
					d.FieldU32("element")
				}
			})

			d.SeekAbs(sectionOffsets.drawOrderGroups.minimumDrawOrders * 8)
			d.FieldArray("minimum_draw_orders", func(d *decode.D) {
				for i := int64(0); i < countInfo.drawOrderGroups; i++ {
					d.FieldU32("element")
				}
			})
		})

		d.FieldStruct("draw_order_group_objects", func(d *decode.D) {
			d.SeekAbs(sectionOffsets.drawOrderGroupObjects.types * 8)
			d.FieldArray("types", func(d *decode.D) {
				for i := int64(0); i < countInfo.drawOrderGroupObjects; i++ {
					d.FieldU32("type", drawOrderGroupObjectTypeNames)
				}
			})

			d.SeekAbs(sectionOffsets.drawOrderGroupObjects.indices * 8)
			d.FieldArray("indices", func(d *decode.D) {
				for i := int64(0); i < countInfo.drawOrderGroupObjects; i++ {
					d.FieldS32("element")
				}
			})

			d.SeekAbs(sectionOffsets.drawOrderGroupObjects.selfIndices * 8)
			d.FieldArray("self_indices", func(d *decode.D) {
				for i := int64(0); i < countInfo.drawOrderGroupObjects; i++ {
					d.FieldS32("element")
				}
			})
		})

		d.FieldStruct("glue", func(d *decode.D) {
			d.SeekAbs(sectionOffsets.glue.ids * 8)
			d.FieldArray("ids", func(d *decode.D) {
				for i := int64(0); i < countInfo.glue; i++ {
					d.FieldUTF8NullFixedLen("id", 64)
				}
			})

			d.SeekAbs(sectionOffsets.glue.keyformBindingSourcesIndices * 8)
			d.FieldArray("keyform_binding_sources_indices", func(d *decode.D) {
				for i := int64(0); i < countInfo.glue; i++ {
					d.FieldS32("element")
				}
			})

			d.SeekAbs(sectionOffsets.glue.keyformSourcesBeginIndices * 8)
			d.FieldArray("keyform_sources_begin_indices", func(d *decode.D) {
				for i := int64(0); i < countInfo.glue; i++ {
					d.FieldS32("element")
				}
			})

			d.SeekAbs(sectionOffsets.glue.keyformSourcesCounts * 8)
			d.FieldArray("keyform_sources_counts", func(d *decode.D) {
				for i := int64(0); i < countInfo.glue; i++ {
					d.FieldS32("element")
				}
			})

			d.SeekAbs(sectionOffsets.glue.artMeshIndicesA * 8)
			d.FieldArray("art_mesh_indices_a", func(d *decode.D) {
				for i := int64(0); i < countInfo.glue; i++ {
					d.FieldS32("element")
				}
			})

			d.SeekAbs(sectionOffsets.glue.artMeshIndicesB * 8)
			d.FieldArray("art_mesh_indices_b", func(d *decode.D) {
				for i := int64(0); i < countInfo.glue; i++ {
					d.FieldS32("element")
				}
			})

			d.SeekAbs(sectionOffsets.glue.glueInfoSourcesBeginIndices * 8)
			d.FieldArray("glue_info_sources_begin_indices", func(d *decode.D) {
				for i := int64(0); i < countInfo.glue; i++ {
					d.FieldS32("element")
				}
			})

			d.SeekAbs(sectionOffsets.glue.glueInfoSourcesCounts * 8)
			d.FieldArray("glue_info_sources_counts", func(d *decode.D) {
				for i := int64(0); i < countInfo.glue; i++ {
					d.FieldS32("element")
				}
			})
		})

		d.FieldStruct("glue_info", func(d *decode.D) {
			d.SeekAbs(sectionOffsets.glueInfo.weights * 8)
			d.FieldArray("weights", func(d *decode.D) {
				for i := int64(0); i < countInfo.glueInfo; i++ {
					d.FieldF32("element")
				}
			})

			d.SeekAbs(sectionOffsets.glueInfo.positionIndices * 8)
			d.FieldArray("position_indices", func(d *decode.D) {
				for i := int64(0); i < countInfo.glueInfo; i++ {
					d.FieldS16("element")
				}
			})
		})

		if version < moc3Version3_03_00 {
			return
		}

		d.FieldStruct("glue_keyforms", func(d *decode.D) {
			d.SeekAbs(sectionOffsets.glueKeyforms.intensities * 8)
			d.FieldArray("intensities", func(d *decode.D) {
				for i := int64(0); i < countInfo.glueKeyforms; i++ {
					d.FieldF32("element")
				}
			})
		})

		if version < moc3Version4_02_00 {
			return
		}

		d.FieldStruct("warp_deformers_v3_3", func(d *decode.D) {

			d.SeekAbs(sectionOffsets.warpDeformersV3_3.isQuadSource * 8)
			d.FieldArray("is_quad_source", func(d *decode.D) {
				for i := int64(0); i < countInfo.warpDeformers; i++ {
					d.FieldBoolFn("element", func(d *decode.D) bool { return d.U32() != 0 })
				}
			})
		})

		d.FieldStruct("parameter_extensions", func(d *decode.D) {
			d.SeekAbs(sectionOffsets.parameterExtensions.keysSourcesBeginIndices * 8)
			d.FieldArray("keys_sources_begin_indices", func(d *decode.D) {
				for i := int64(0); i < countInfo.parameters; i++ {
					d.FieldS32("element")
				}
			})

			d.SeekAbs(sectionOffsets.parameterExtensions.keysSourcesCounts * 8)
			d.FieldArray("keys_sources_counts", func(d *decode.D) {
				for i := int64(0); i < countInfo.parameters; i++ {
					d.FieldS32("element")
				}
			})
		})

		d.FieldStruct("warp_deformers_v4_2", func(d *decode.D) {
			d.SeekAbs(sectionOffsets.warpDeformersV4_2.keyformColorSourcesBeginIndices * 8)
			d.FieldArray("keyform_color_sources_begin_indices", func(d *decode.D) {
				for i := int64(0); i < countInfo.warpDeformers; i++ {
					d.FieldS32("element")
				}
			})
		})

		d.FieldStruct("rotation_deformers_v4_2", func(d *decode.D) {
			d.SeekAbs(sectionOffsets.rotationDeformersV4_2.keyformColorSourcesBeginIndices * 8)
			d.FieldArray("keyform_color_sources_begin_indices", func(d *decode.D) {
				for i := int64(0); i < countInfo.rotationDeformers; i++ {
					d.FieldS32("element")
				}
			})
		})

		d.FieldStruct("art_meshes_v4_2", func(d *decode.D) {
			d.SeekAbs(sectionOffsets.artMeshesV4_2.keyformColorSourcesBeginIndices * 8)
			d.FieldArray("keyform_color_sources_begin_indices", func(d *decode.D) {
				for i := int64(0); i < countInfo.artMeshes; i++ {
					d.FieldS32("element")
				}
			})
		})

		d.FieldStruct("keyform_multiply_colors", func(d *decode.D) {
			d.SeekAbs(sectionOffsets.keyformMultiplyColors.r * 8)
			d.FieldArray("r", func(d *decode.D) {
				for i := int64(0); i < countInfo.keyformMultiplyColors; i++ {
					d.FieldF32("element")
				}
			})

			d.SeekAbs(sectionOffsets.keyformMultiplyColors.g * 8)
			d.FieldArray("g", func(d *decode.D) {
				for i := int64(0); i < countInfo.keyformMultiplyColors; i++ {
					d.FieldF32("element")
				}
			})

			d.SeekAbs(sectionOffsets.keyformMultiplyColors.b * 8)
			d.FieldArray("b", func(d *decode.D) {
				for i := int64(0); i < countInfo.keyformMultiplyColors; i++ {
					d.FieldF32("element")
				}
			})
		})

		d.FieldStruct("keyform_screen_colors", func(d *decode.D) {
			d.SeekAbs(sectionOffsets.keyformScreenColors.r * 8)
			d.FieldArray("r", func(d *decode.D) {
				for i := int64(0); i < countInfo.keyformScreenColors; i++ {
					d.FieldF32("element")
				}
			})

			d.SeekAbs(sectionOffsets.keyformScreenColors.g * 8)
			d.FieldArray("g", func(d *decode.D) {
				for i := int64(0); i < countInfo.keyformScreenColors; i++ {
					d.FieldF32("element")
				}
			})

			d.SeekAbs(sectionOffsets.keyformScreenColors.b * 8)
			d.FieldArray("b", func(d *decode.D) {
				for i := int64(0); i < countInfo.keyformScreenColors; i++ {
					d.FieldF32("element")
				}
			})
		})

		d.FieldStruct("parameter_offsets_v4_2", func(d *decode.D) {
			d.SeekAbs(sectionOffsets.parameterOffsetsV4_2.parameterTypes * 8)
			d.FieldArray("parameter_types", func(d *decode.D) {
				for i := int64(0); i < countInfo.parameters; i++ {
					d.FieldU32("type", parameterTypeNames)
				}
			})

			d.SeekAbs(sectionOffsets.parameterOffsetsV4_2.blendShapeParameterBindingSourcesBeginIndices * 8)
			d.FieldArray("blend_shape_parameter_binding_sources_begin_indices", func(d *decode.D) {
				for i := int64(0); i < countInfo.parameters; i++ {
					d.FieldS32("element")
				}
			})

			d.SeekAbs(sectionOffsets.parameterOffsetsV4_2.blendShapeParameterBindingSourcesCounts * 8)
			d.FieldArray("blend_shape_parameter_binding_sources_counts", func(d *decode.D) {
				for i := int64(0); i < countInfo.parameters; i++ {
					d.FieldS32("element")
				}
			})
		})

		d.FieldStruct("blend_shape_parameter_bindings", func(d *decode.D) {
			d.SeekAbs(sectionOffsets.blendShapeParameterBindings.keysSourcesBeginIndices * 8)
			d.FieldArray("keys_sources_begin_indices", func(d *decode.D) {
				for i := int64(0); i < countInfo.blendShapeParameterBindings; i++ {
					d.FieldS32("element")
				}
			})

			d.SeekAbs(sectionOffsets.blendShapeParameterBindings.keysSourcesCounts * 8)
			d.FieldArray("keys_sources_counts", func(d *decode.D) {
				for i := int64(0); i < countInfo.blendShapeParameterBindings; i++ {
					d.FieldS32("element")
				}
			})

			d.SeekAbs(sectionOffsets.blendShapeParameterBindings.baseKeyIndices * 8)
			d.FieldArray("base_key_indices", func(d *decode.D) {
				for i := int64(0); i < countInfo.blendShapeParameterBindings; i++ {
					d.FieldS32("element")
				}
			})
		})

		d.FieldStruct("blend_shape_keyform_bindings", func(d *decode.D) {
			d.SeekAbs(sectionOffsets.blendShapeKeyformBindings.parameterBindingSourcesIndices * 8)
			d.FieldArray("parameter_binding_sources_indices", func(d *decode.D) {
				for i := int64(0); i < countInfo.blendShapeKeyformBindings; i++ {
					d.FieldS32("element")
				}
			})

			d.SeekAbs(sectionOffsets.blendShapeKeyformBindings.blendShapeConstraintIndexSourcesBeginIndices * 8)
			d.FieldArray("blend_shape_constraint_index_sources_begin_indices", func(d *decode.D) {
				for i := int64(0); i < countInfo.blendShapeKeyformBindings; i++ {
					d.FieldS32("element")
				}
			})

			d.SeekAbs(sectionOffsets.blendShapeKeyformBindings.blendShapeConstraintIndexSourcesCounts * 8)
			d.FieldArray("blend_shape_constraint_index_sources_counts", func(d *decode.D) {
				for i := int64(0); i < countInfo.blendShapeKeyformBindings; i++ {
					d.FieldS32("element")
				}
			})

			d.SeekAbs(sectionOffsets.blendShapeKeyformBindings.keyformSourcesBlendShapeIndices * 8)
			d.FieldArray("keyform_sources_blend_shape_indices", func(d *decode.D) {
				for i := int64(0); i < countInfo.blendShapeKeyformBindings; i++ {
					d.FieldS32("element")
				}
			})

			d.SeekAbs(sectionOffsets.blendShapeKeyformBindings.keyformSourcesBlendShapeCounts * 8)
			d.FieldArray("keyform_sources_blend_shape_counts", func(d *decode.D) {
				for i := int64(0); i < countInfo.blendShapeKeyformBindings; i++ {
					d.FieldS32("element")
				}
			})
		})

		d.FieldStruct("blend_shapes_warp_deformers", func(d *decode.D) {
			d.SeekAbs(sectionOffsets.blendShapesWarpDeformers.targetIndices * 8)
			d.FieldArray("target_indices", func(d *decode.D) {
				for i := int64(0); i < countInfo.blendShapesWarpDeformers; i++ {
					d.FieldS32("element")
				}
			})

			d.SeekAbs(sectionOffsets.blendShapesWarpDeformers.blendShapeKeyformBindingSourcesBeginIndices * 8)
			d.FieldArray("blend_shape_keyform_binding_sources_begin_indices", func(d *decode.D) {
				for i := int64(0); i < countInfo.blendShapesWarpDeformers; i++ {
					d.FieldS32("element")
				}
			})

			d.SeekAbs(sectionOffsets.blendShapesWarpDeformers.blendShapeKeyformBindingSourcesCounts * 8)
			d.FieldArray("blend_shape_keyform_binding_sources_counts", func(d *decode.D) {
				for i := int64(0); i < countInfo.blendShapesWarpDeformers; i++ {
					d.FieldS32("element")
				}
			})
		})

		d.FieldStruct("blend_shapes_art_meshes", func(d *decode.D) {
			d.SeekAbs(sectionOffsets.blendShapesArtMeshes.targetIndices * 8)
			d.FieldArray("target_indices", func(d *decode.D) {
				for i := int64(0); i < countInfo.blendShapesArtMeshes; i++ {
					d.FieldS32("element")
				}
			})

			d.SeekAbs(sectionOffsets.blendShapesArtMeshes.blendShapeKeyformBindingSourcesBeginIndices * 8)
			d.FieldArray("blend_shape_keyform_binding_sources_begin_indices", func(d *decode.D) {
				for i := int64(0); i < countInfo.blendShapesArtMeshes; i++ {
					d.FieldS32("element")
				}
			})

			d.SeekAbs(sectionOffsets.blendShapesArtMeshes.blendShapeKeyformBindingSourcesCounts * 8)
			d.FieldArray("blend_shape_keyform_binding_sources_counts", func(d *decode.D) {
				for i := int64(0); i < countInfo.blendShapesArtMeshes; i++ {
					d.FieldS32("element")
				}
			})
		})

		d.FieldStruct("blend_shape_constraint_indices", func(d *decode.D) {
			d.SeekAbs(sectionOffsets.blendShapeConstraintIndices.blendShapeConstraintSourcesIndices * 8)
			d.FieldArray("blend_shape_constraint_sources_indices", func(d *decode.D) {
				for i := int64(0); i < countInfo.blendShapeConstraintIndices; i++ {
					d.FieldS32("element")
				}
			})
		})

		d.FieldStruct("blend_shape_constraints", func(d *decode.D) {
			d.SeekAbs(sectionOffsets.blendShapeConstraints.parameterIndices * 8)
			d.FieldArray("parameter_indices", func(d *decode.D) {
				for i := int64(0); i < countInfo.blendShapeConstraints; i++ {
					d.FieldS32("element")
				}
			})

			d.SeekAbs(sectionOffsets.blendShapeConstraints.blendShapeConstraintValueSourcesBeginIndices * 8)
			d.FieldArray("blend_shape_constraint_value_sources_begin_indices", func(d *decode.D) {
				for i := int64(0); i < countInfo.blendShapeConstraints; i++ {
					d.FieldS32("element")
				}
			})

			d.SeekAbs(sectionOffsets.blendShapeConstraints.blendShapeConstraintValueSourcesCounts * 8)
			d.FieldArray("blend_shape_constraint_value_sources_counts", func(d *decode.D) {
				for i := int64(0); i < countInfo.blendShapeConstraints; i++ {
					d.FieldS32("element")
				}
			})
		})

		d.FieldStruct("blend_shape_constraint_values", func(d *decode.D) {
			d.SeekAbs(sectionOffsets.blendShapeConstraintValues.keys * 8)
			d.FieldArray("keys", func(d *decode.D) {
				for i := int64(0); i < countInfo.blendShapeConstraintValues; i++ {
					d.FieldF32("element")
				}
			})

			d.SeekAbs(sectionOffsets.blendShapeConstraintValues.weights * 8)
			d.FieldArray("weights", func(d *decode.D) {
				for i := int64(0); i < countInfo.blendShapeConstraintValues; i++ {
					d.FieldF32("element")
				}
			})
		})
	})

	return nil
}
