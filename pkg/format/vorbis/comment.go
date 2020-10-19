package vorbis

import (
	"encoding/base64"
	"fq/pkg/bitbuf"
	"fq/pkg/decode"
	"fq/pkg/format"
	"strings"
)

var flacPicture []*decode.Format

var Comment = format.MustRegister(&decode.Format{
	Name:      "vorbis_comment",
	New:       func() decode.Decoder { return &CommentDecoder{} },
	SkipProbe: true,
	Deps: []decode.Dep{
		{Names: []string{"flac_picture"}, Formats: &flacPicture},
	},
})

// CommentDecoder is a vorbis packet decoder
type CommentDecoder struct {
	decode.Common
}

// Decode vorbis comment
func (d *CommentDecoder) Decode() {
	lenStr := func(name string) string {
		len := d.FieldU32LE(name + "_length")
		return d.FieldUTF8(name, int64(len))
	}
	lenStr("vendor")
	userCommentListLength := d.FieldU32LE("user_comment_list_length")
	d.Array("user_comment", func() {
		for i := uint64(0); i < userCommentListLength; i++ {
			pair := lenStr("user_comment")
			pairParts := strings.SplitN(pair, "=", 2)
			if len(pairParts) == 2 {
				// METADATA_BLOCK_PICTURE=<base64>
				k, v := strings.ToUpper(pairParts[0]), pairParts[1]
				var metadataBlockPicture = "METADATA_BLOCK_PICTURE"
				if k == metadataBlockPicture {
					bs, err := base64.StdEncoding.DecodeString(v)
					if err == nil {
						bb, err := bitbuf.NewFromBytes(bs, 0)
						if err != nil {
							panic(err) // TODO: fixme
						}

						d.FieldDecodeBitBuf("picture",
							d.Pos()-int64(len(v))*8,
							int64(len(v)*8),
							bb,
							flacPicture,
						)
					} else {
						panic(err)
					}
				}
			}
		}
	})
}
