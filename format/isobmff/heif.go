package isobmff

// HEIF ISO/IEC 23008-12 (MPEG-H Part 12)

import (
	"cmp"
	"slices"

	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
)

// TODO: avc, jpeg and jpeg2000 in HEIF

var heifAV1CCRGroup decode.Group
var heifAV1FrameGroup decode.Group
var heifHEVCAUGroup decode.Group
var heifHEVCDCRGroup decode.Group

// TODO: not really used as icc is shared
var heifICCProfileGroup decode.Group

func init() {
	interp.RegisterFormat(
		format.HEIF,
		&decode.Format{
			Description: "High Efficiency Image Format",
			Groups: []*decode.Group{
				format.Probe,
				format.Image,
			},
			DecodeFn: heifDecode,
			DefaultInArg: format.HEIF_In{
				AllowTruncated: false,
			},
			Dependencies: []decode.Dependency{
				{Groups: []*decode.Group{format.AV1_CCR}, Out: &heifAV1CCRGroup},
				{Groups: []*decode.Group{format.AV1_Frame}, Out: &heifAV1FrameGroup},
				{Groups: []*decode.Group{format.HEVC_AU}, Out: &heifHEVCAUGroup},
				{Groups: []*decode.Group{format.HEVC_DCR}, Out: &heifHEVCDCRGroup},
				{Groups: []*decode.Group{format.ICC_Profile}, Out: &heifICCProfileGroup},
			},
		})
}

type ilocBoxExtent struct {
	offset uint64
	length uint64
}

type ilocBoxItem struct {
	id uint64
	// from ISO 14496-12 section 8.11.3 (iloc box)
	// file_offset: by the usual absolute file offsets into the file at data_reference_index; (construction_method == 0)
	// idat_offset: by box offsets into the idat box in the same meta box; neither the data_reference_index nor extent_index fields are used; (construction_method == 1)
	// item_offset: by item offset into the items indicated by the extent_index field, which is only used (currently) by this construction method. (construction_method == 2).
	constructionMethod uint64
	baseOffset         uint64
	extents            []ilocBoxExtent
}

type ilocBox struct {
	offsetSize     uint64
	lengthSize     uint64
	baseOffsetSize uint64
	items          []ilocBoxItem
}

type ipmaBoxAssociation struct {
	essential     bool
	propertyIndex int
}

type ipmaBoxEntry struct {
	itemID       uint64
	associations []ipmaBoxAssociation
}

type ipmaBox struct {
	entries []ipmaBoxEntry
}

type infeBox struct {
	itemID      uint64
	itemType    string
	contentType string
	itemName    string
}

func heifItems(d *decode.D, ctx *decodeContext) {
	meta := ctx.root.find("meta")
	if meta == nil {
		return
	}

	itemsCollected := slices.Collect(findAllData[*infeBox](meta, "iinf/infe"))
	slices.SortStableFunc(itemsCollected, func(a, b *infeBox) int {
		return cmp.Compare(a.itemID, b.itemID)
	})

	ipco := meta.find("iprp/ipco")
	iloc := findData[*ilocBox](meta, "iloc")

	d.FieldArray("items", func(d *decode.D) {
		for _, inf := range itemsCollected {
			d.FieldStruct("item", func(d *decode.D) {
				if inf == nil {
					return
				}

				d.FieldValueUint("id", inf.itemID)
				d.FieldValueStr("type", inf.itemType)
				if inf.contentType != "" {
					d.FieldValueStr("content_type", inf.contentType)
				}
				if inf.itemName != "" {
					d.FieldValueStr("name", inf.itemName)
				}

				ipma := findData[*ipmaBox](meta, "iprp/ipma")
				if ipma != nil {
					d.FieldArray("properties", func(d *decode.D) {
						for _, entry := range ipma.entries {
							if entry.itemID != inf.itemID {
								continue
							}
							for _, assoc := range entry.associations {
								d.FieldStruct("property", func(d *decode.D) {
									idx := assoc.propertyIndex - 1
									d.FieldValueUint("index", uint64(assoc.propertyIndex))
									d.FieldValueBool("essential", assoc.essential)
									if ipco != nil && idx >= 0 && idx < len(ipco.children) {
										prop := ipco.children[idx]
										d.FieldValueStr("type", prop.typ)
									}
								})
							}
						}
					})
				}

				if iloc != nil {
					var ilocItem *ilocBoxItem
					for i := range iloc.items {
						if iloc.items[i].id == inf.itemID {
							ilocItem = &iloc.items[i]
							break
						}
					}

					// TODO: only file offset for now and single extent
					if ilocItem != nil && ilocItem.constructionMethod == 0 && len(ilocItem.extents) == 1 {
						extent := ilocItem.extents[0]
						d.SeekAbs(int64(ilocItem.baseOffset+extent.offset)*8, func(d *decode.D) {
							switch inf.itemType {
							case "hvc1":
								d.FieldFormatLen("data", int64(extent.length)*8, &heifHEVCAUGroup, nil)
							case "av01":
								d.FieldFormatLen("data", int64(extent.length)*8, &heifAV1FrameGroup, nil)
							default:
								d.FieldRawLen("data", int64(extent.length)*8)
							}
						})
					}
				}
			})
		}
	})
}

func heifDecode(d *decode.D) any {
	var hi format.HEIF_In
	d.ArgAs(&hi)

	ctx := isobmffDecode(d, hi.AllowTruncated, func(firstType string, ftyp ftypBox) {
		switch firstType {
		case "ftyp":
		default:
			d.Errorf("type not ftyp")
		}
		switch ftyp.majorBrand {
		case "mif1", "mif3", "msf1", "avif":
			return
		}
		for _, b := range ftyp.minorBrands {
			switch b {
			case "mif1", "mif3", "msf1", "avif":
				return
			}
		}
		d.Errorf("not a HEIF brand (mif1/mif3/msf1/avif)")
	})

	heifItems(d, ctx)
	traks := slices.Collect(ctx.root.findAll("moov/trak"))
	if len(traks) > 0 {
		// TODO: replace MP4_In with some shared struct
		mp4Tracks(d, format.MP4_In{DecodeSamples: true}, traks, nil)
	}

	return nil
}
