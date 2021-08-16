package vorbis

import (
	"encoding/base64"
	"fq/format"
	"fq/format/registry"
	"fq/pkg/bitio"
	"fq/pkg/decode"
	"strings"
)

var flacPicture []*decode.Format

func init() {
	registry.MustRegister(&decode.Format{
		Name:        format.VORBIS_COMMENT,
		Description: "Vorbis comment",
		DecodeFn:    commentDecode,
		Dependencies: []decode.Dependency{
			{Names: []string{format.FLAC_PICTURE}, Formats: &flacPicture},
		},
	})
}

func commentDecode(d *decode.D, in interface{}) interface{} {
	vendorLen := d.FieldU32LE("vendor_length")
	d.FieldUTF8("vendor", int(vendorLen))
	userCommentListLength := d.FieldU32LE("user_comment_list_length")
	i := uint64(0)
	d.FieldStructArrayLoopFn("user_comments", "user_comment", func() bool { return i < userCommentListLength }, func(d *decode.D) {
		userCommentLength := d.FieldU32LE("length")
		userComment := d.FieldUTF8("vendor", int(userCommentLength))
		pairParts := strings.SplitN(userComment, "=", 2)
		if len(pairParts) == 2 {
			k, v := strings.ToUpper(pairParts[0]), pairParts[1]
			var metadataBlockPicture = "METADATA_BLOCK_PICTURE"
			if k == metadataBlockPicture {
				// METADATA_BLOCK_PICTURE=<base64-flac-picture-metadatablock>
				bs, err := base64.StdEncoding.DecodeString(v)
				if err == nil {
					bb := bitio.NewBufferFromBytes(bs, -1)
					d.FieldFormatBitBuf("picture", bb, flacPicture)
				} else {
					panic(err)
				}
			}
		}
		i++
	})

	return nil
}
