package vorbis

import (
	"encoding/base64"
	"errors"
	"io"
	"strings"

	"github.com/wader/fq/format"
	"github.com/wader/fq/format/registry"
	"github.com/wader/fq/pkg/decode"
)

var flacPicture decode.Group

func init() {
	registry.MustRegister(decode.Format{
		Name:        format.VORBIS_COMMENT,
		Description: "Vorbis comment",
		DecodeFn:    commentDecode,
		Dependencies: []decode.Dependency{
			{Names: []string{format.FLAC_PICTURE}, Group: &flacPicture},
		},
	})
}

func commentDecode(d *decode.D, in interface{}) interface{} {
	vendorLen := d.FieldU32LE("vendor_length")
	d.FieldUTF8("vendor", int(vendorLen))
	userCommentListLength := d.FieldU32LE("user_comment_list_length")
	i := uint64(0)
	d.FieldStructArrayLoop("user_comments", "user_comment", func() bool { return i < userCommentListLength }, func(d *decode.D) {
		userCommentLength := d.FieldU32LE("length")
		userCommentStart := d.Pos()
		userComment := d.FieldUTF8("comment", int(userCommentLength))
		var metadataBlockPicturePreix = "METADATA_BLOCK_PICTURE="
		var metadataBlockPicturePrefixLower = "metadata_block_picture="

		if strings.HasPrefix(userComment, metadataBlockPicturePreix) ||
			strings.HasPrefix(userComment, metadataBlockPicturePrefixLower) {

			base64Offset := int64(len(metadataBlockPicturePreix)) * 8
			base64Len := int64(len(userComment))*8 - base64Offset

			rFn := func(r io.Reader) io.Reader { return base64.NewDecoder(base64.StdEncoding, r) }

			_, uncompressedBB, dv, _, err := d.TryFieldReaderRangeFormat("picture", userCommentStart+base64Offset, base64Len, rFn, flacPicture, nil)
			if dv == nil && errors.As(err, &decode.FormatsError{}) {
				d.FieldRootBitBuf("picture", uncompressedBB)
			}
		}
		i++
	})

	return nil
}
