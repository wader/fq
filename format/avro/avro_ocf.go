package avro

import (
	"bytes"
	"compress/flate"
	"embed"
	"hash/crc32"

	"github.com/golang/snappy"
	"github.com/wader/fq/format"
	"github.com/wader/fq/format/avro/decoders"
	"github.com/wader/fq/format/avro/schema"
	"github.com/wader/fq/pkg/bitio"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
	"github.com/wader/fq/pkg/scalar"
)

//go:embed avro_ocf.md
var avroOcfFS embed.FS

func init() {
	interp.RegisterFormat(
		format.Avro_Ocf,
		&decode.Format{
			Description: "Avro object container file",
			Groups:      []*decode.Group{format.Probe},
			DecodeFn:    decodeAvroOCF,
		})
	interp.RegisterFS(avroOcfFS)
}

type HeaderData struct {
	Schema schema.SimplifiedSchema
	Codec  string
	Sync   []byte
}

const headerSchemaSpec = `
{
  "type": "record",
  "name": "org.apache.avro.file.Header",
  "fields": [
   {"name": "meta", "type": {"type": "map", "values": "string"}},
   {"name": "sync", "type": {"type": "fixed", "name": "Sync", "size": 16}}
  ]
}`

func decodeHeader(d *decode.D) HeaderData {
	d.FieldRawLen("magic", 4*8, d.AssertBitBuf([]byte{'O', 'b', 'j', 1}))

	var headerData HeaderData

	headerSchema, err := schema.FromSchemaString(headerSchemaSpec)
	if err != nil {
		d.Fatalf("Failed to parse header schema: %v", err)
	}
	decodeHeaderFn, err := decoders.DecodeFnForSchema(headerSchema)
	if err != nil {
		d.Fatalf("failed to parse header: %v", err)
	}

	header := decodeHeaderFn("header", d)
	headerRecord, ok := header.(map[string]any)
	if !ok {
		d.Fatalf("header is not a map")
	}
	meta, ok := headerRecord["meta"].(map[string]any)
	if !ok {
		d.Fatalf("header.meta is not a map")
	}

	metaSchema, ok := meta["avro.schema"].(string)
	if !ok {
		d.Fatalf("missing meta avro.schema")
	}

	headerData.Schema, err = schema.FromSchemaString(metaSchema)
	if err != nil {
		d.Fatalf("failed to parse schema: %v", err)
	}
	if codec, ok := meta["avro.codec"]; ok {
		headerData.Codec, ok = codec.(string)
		if !ok {
			d.Fatalf("avro.codec is not a string")
		}
	} else {
		headerData.Codec = "null"
	}

	headerData.Sync, ok = headerRecord["sync"].([]byte)
	if !ok {
		d.Fatalf("header.sync is not a byte array")
	}
	return headerData
}

func decodeBlockCodec(d *decode.D, dataSize int64, codec string) *bytes.Buffer {
	bb := &bytes.Buffer{}
	if codec == "deflate" {
		br := d.FieldRawLen("compressed", dataSize*8)
		d.Copy(bb, flate.NewReader(bitio.NewIOReader(br)))
	} else if codec == "snappy" {
		// Everything but last 4 bytes which are the checksum
		n := dataSize - 4
		br := d.FieldRawLen("compressed", n*8)

		// This could be simplified to be similar to deflate, however snappy's reader only works for streaming frames,
		// not block data. See https://github.com/google/snappy/blob/main/framing_format.txt for details.
		compressed := make([]byte, n)
		if _, err := bitio.ReadFull(br, compressed, n*8); err != nil {
			d.Fatalf("failed reading compressed data %v", err)
		}
		decompressed, err := snappy.Decode(nil, compressed)
		if err != nil {
			d.Fatalf("failed decompressing data: %v", err)
		}
		d.Copy(bb, bytes.NewReader(decompressed))

		// Check the checksum
		crc32W := crc32.NewIEEE()
		d.Copy(crc32W, bytes.NewReader(bb.Bytes()))
		d.FieldU32("crc", d.UintValidateBytes(crc32W.Sum(nil)), scalar.UintHex)
	} else {
		// Unknown codec, just dump the compressed data.
		d.FieldRawLen("compressed", dataSize*8, scalar.BitBufDescription(codec+" encoded"))
		return nil
	}
	return bb
}

func decodeAvroOCF(d *decode.D) any {
	header := decodeHeader(d)

	decodeFn, err := decoders.DecodeFnForSchema(header.Schema)
	if err != nil {
		d.Fatalf("unable to create codec: %v", err)
	}

	d.FieldStructArrayLoop("blocks", "block", func() bool { return d.NotEnd() }, func(d *decode.D) {
		count := d.FieldSintFn("count", decoders.VarZigZag)
		if count <= 0 {
			return
		}
		size := d.FieldSintFn("size", decoders.VarZigZag)
		i := int64(0)

		if header.Codec != "null" {
			if bb := decodeBlockCodec(d, size, header.Codec); bb != nil {
				d.FieldArrayRootBitBufFn("data", bitio.NewBitReader(bb.Bytes(), -1), func(d *decode.D) {
					for ; i < count; i++ {
						decodeFn("data", d)
					}
				})
			}
		} else {
			d.FieldArrayLoop("data", func() bool { return i < count }, func(d *decode.D) {
				decodeFn("datum", d)
				i++
			})
		}
		d.FieldRawLen("sync", 16*8, d.AssertBitBuf(header.Sync))
	})

	return nil
}
