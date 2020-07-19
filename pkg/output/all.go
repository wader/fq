package output

import (
	"fq/pkg/decode"
	"fq/pkg/output/json"
	"fq/pkg/output/text"
)

var All = map[string]*decode.FieldOutput{
	json.FieldOutput.Name: json.FieldOutput,
	text.FieldOutput.Name: text.FieldOutput,
}
