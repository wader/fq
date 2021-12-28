package codecs

const (
	NULL    = "null"
	BOOLEAN = "boolean"
	INT     = "int"
	LONG    = "long"
	FLOAT   = "float"
	DOUBLE  = "double"
	BYTES   = "bytes"
	STRING  = "string"
	RECORD  = "record"
)

type SimplifiedSchema struct {
	Type        string
	Name        *string
	Fields      []SimplifiedSchemaField
	Symbols     *[]string
	Items       *SimplifiedSchema
	LogicalType *string
	Scale       *int
	Precision   *int
	UnionTypes  []SimplifiedSchema
	//Choosing not to handle Default as it adds a lot of complexity and this is used for showing the binary
	//representation of the data, not fully parsing it. See https://github.com/linkedin/goavro/blob/master/record.go
	//for how it could be handled.
}

type SimplifiedSchemaField struct {
	Name string
	Type SimplifiedSchema
}
