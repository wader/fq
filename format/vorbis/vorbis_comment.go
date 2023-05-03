package vorbis

import (
	"encoding/base64"
	"io"
	"strings"

	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
)

var flacPicture decode.Group

func init() {
	interp.RegisterFormat(
		format.Vorbis_Comment,
		&decode.Format{
			Description: "Vorbis comment",
			DecodeFn:    commentDecode,
			Dependencies: []decode.Dependency{
				{Groups: []*decode.Group{format.FLAC_Picture}, Out: &flacPicture},
			},
		})
}

func commentDecode(d *decode.D) any {
	d.Endian = decode.LittleEndian

	vendorLen := d.FieldU32("vendor_length")
	d.FieldUTF8("vendor", int(vendorLen))
	userCommentListLength := d.FieldU32("user_comment_list_length")
	i := uint64(0)
	d.FieldStructArrayLoop("user_comments", "user_comment", func() bool { return i < userCommentListLength }, func(d *decode.D) {
		userCommentLength := d.FieldU32("length")
		userCommentStart := d.Pos()
		userComment := d.FieldUTF8("comment", int(userCommentLength))
		var metadataBlockPicturePrefix = "METADATA_BLOCK_PICTURE="
		var metadataBlockPicturePrefixLower = "metadata_block_picture="

		if strings.HasPrefix(userComment, metadataBlockPicturePrefix) ||
			strings.HasPrefix(userComment, metadataBlockPicturePrefixLower) {

			base64Offset := int64(len(metadataBlockPicturePrefix)) * 8
			base64Len := int64(len(userComment))*8 - base64Offset
			_, base64Br, dv, _, _ := d.TryFieldReaderRangeFormat(
				"picture",
				userCommentStart+base64Offset, base64Len,
				func(r io.Reader) io.Reader { return base64.NewDecoder(base64.StdEncoding, r) },
				&flacPicture, nil,
			)
			if dv == nil && base64Br != nil {
				d.FieldRootBitBuf("picture", base64Br)
			}
		}
		i++
	})

	return nil
}
