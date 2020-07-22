package output

import (
	"fq/pkg/decode"
	"fq/pkg/output/json"
	"fq/pkg/output/rangex"
	"fq/pkg/output/raw"
	"fq/pkg/output/text"
	"fq/pkg/output/value"
)

var All = map[string]*decode.FieldOutput{
	value.FieldOutput.Name:  value.FieldOutput,
	text.FieldOutput.Name:   text.FieldOutput,
	json.FieldOutput.Name:   json.FieldOutput,
	raw.FieldOutput.Name:    raw.FieldOutput,
	rangex.FieldOutput.Name: rangex.FieldOutput,
}
