package vorbis

import (
	"encoding/base64"
	"fq/pkg/bitbuf"
	"fq/pkg/decode"
	"fq/pkg/format"
	"strings"
)

var flacPicture []*decode.Format

func init() {
	format.MustRegister(&decode.Format{
		Name:     format.VORBIS_COMMENT,
		DecodeFn: commentDecode,
		Deps: []decode.Dep{
			{Names: []string{format.FLAC_PICTURE}, Formats: &flacPicture},
		},
	})
}

func commentDecode(d *decode.D) interface{} {
	fieldLenStr := func(d *decode.D, name string) string {
		len := d.FieldU32LE(name + "_length")
		return d.FieldUTF8(name, int64(len))
	}
	fieldLenStr(d, "vendor")
	userCommentListLength := d.FieldU32LE("user_comment_list_length")
	d.FieldArrayFn("user_comment", func(d *decode.D) {
		for i := uint64(0); i < userCommentListLength; i++ {
			pair := fieldLenStr(d, "user_comment")
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

	return nil
}
