package schema

import (
	"encoding/json"
	"errors"
	"fmt"
)

const (
	ARRAY   = "array"
	BOOLEAN = "boolean"
	BYTES   = "bytes"
	DOUBLE  = "double"
	ENUM    = "enum"
	FIXED   = "fixed"
	FLOAT   = "float"
	INT     = "int"
	LONG    = "long"
	MAP     = "map"
	NULL    = "null"
	RECORD  = "record"
	STRING  = "string"
	UNION   = "union" // avro spec doesn't treat unions like this, but makes it easier for us
)

type SimplifiedSchema struct {
	Type        string            `json:"type"`
	Name        string            `json:"name"`
	LogicalType string            `json:"logicalType,omitempty"`
	Size        int               `json:"size,omitempty"`
	Scale       int               `json:"scale,omitempty"`
	Precision   int               `json:"precision,omitempty"`
	Items       *SimplifiedSchema `json:"items,omitempty"`
	Fields      []Field           `json:"fields,omitempty"`
	Symbols     []string          `json:"symbols,omitempty"`
	Values      *SimplifiedSchema `json:"values,omitempty"`
	UnionTypes  []SimplifiedSchema
	// Choosing not to handle Default as it adds a lot of complexity and this is used for showing the binary
	// representation of the data, not fully parsing it. See https://github.com/linkedin/goavro/blob/master/record.go
	// for how it could be handled.
}

type Field struct {
	Name string
	Type SimplifiedSchema
}

func FromSchemaString(schemaString string) (SimplifiedSchema, error) {
	var jsonSchema any
	if err := json.Unmarshal([]byte(schemaString), &jsonSchema); err != nil {
		return SimplifiedSchema{}, fmt.Errorf("failed to unmarshal header schema: %w", err)
	}

	return From(jsonSchema)
}

func From(schema any) (SimplifiedSchema, error) {
	if schema == nil {
		return SimplifiedSchema{}, errors.New("schema cannot be nil")
	}
	var s SimplifiedSchema
	switch v := schema.(type) {
	case []any:
		s.Type = UNION
		for _, i := range v {
			unionType, err := From(i)
			if err != nil {
				return s, fmt.Errorf("failed parsing union type: %w", err)
			}
			if unionType.Type == UNION {
				return s, errors.New("sub-unions are not supported")
			}
			s.UnionTypes = append(s.UnionTypes, unionType)
		}
	case string:
		s.Type = v
	case map[string]any:
		var err error
		if s.Type, err = getString(v, "type", true); err != nil {
			return s, err
		}
		if s.Name, err = getString(v, "name", false); err != nil {
			return s, err
		}
		if s.LogicalType, err = getString(v, "logicalType", false); err != nil {
			return s, err
		}
		if s.Scale, err = getInt(v, "scale", false); err != nil {
			return s, err
		}
		if s.Precision, err = getInt(v, "precision", false); err != nil {
			return s, err
		}
		if s.Size, err = getInt(v, "size", false); err != nil {
			return s, err
		}
		if s.Type == RECORD {
			if s.Fields, err = getFields(v); err != nil {
				return s, fmt.Errorf("failed parsing fields: %w", err)
			}
		} else if s.Type == ENUM {
			if s.Symbols, err = getSymbols(v); err != nil {
				return s, fmt.Errorf("failed parsing symbols: %w", err)
			}
		} else if s.Type == ARRAY {
			if s.Items, err = getSchema(v, "items"); err != nil {
				return s, fmt.Errorf("failed parsing items: %w", err)
			}
		} else if s.Type == MAP {
			if s.Values, err = getSchema(v, "values"); err != nil {
				return s, fmt.Errorf("failed parsing values: %w", err)
			}
		}
	default:
		return s, errors.New("unknown schema")
	}
	return s, nil
}

func getSchema(m map[string]any, key string) (*SimplifiedSchema, error) {
	vI, ok := m[key]
	if !ok {
		return nil, fmt.Errorf("%s not found", key)
	}
	v, err := From(vI)
	if err != nil {
		return nil, fmt.Errorf("failed parsing %s: %w", key, err)
	}
	return &v, nil
}

func getSymbols(m map[string]any) ([]string, error) {
	vI, ok := m["symbols"]
	if !ok {
		return nil, errors.New("symbols required for enum")
	}
	vA, ok := vI.([]any)
	if !ok {
		return nil, errors.New("symbols must be an array")
	}
	symbols := make([]string, len(vA))
	for i, entry := range vA {
		v, ok := entry.(string)
		if !ok {
			return nil, errors.New("symbols must be an array of strings")
		}
		symbols[i] = v
	}
	return symbols, nil
}

func getFields(m map[string]any) ([]Field, error) {
	var fields []Field
	var err error

	fieldsI, ok := m["fields"]
	if !ok {
		return fields, errors.New("no fields")
	}
	fieldsAI, ok := fieldsI.([]any)
	if !ok {
		return fields, errors.New("fields is not an array")
	}

	for _, fieldI := range fieldsAI {
		field, ok := fieldI.(map[string]any)
		if !ok {
			return fields, errors.New("field is not a json object")
		}
		var f Field
		f.Name, err = getString(field, "name", true)
		if err != nil {
			return fields, fmt.Errorf("failed parsing field name: %w", err)
		}
		t, ok := field["type"]
		if !ok {
			return fields, errors.New("field type must be a object")
		}

		if f.Type, err = From(t); err != nil {
			return fields, fmt.Errorf("failed parsing field %s type: %w", f.Name, err)
		}
		fields = append(fields, f)
	}
	return fields, nil
}

func getString(m map[string]any, key string, required bool) (string, error) {
	v, ok := m[key]
	if !ok {
		if required {
			return "", fmt.Errorf("%s is required", key)
		}
		return "", nil
	}
	s, ok := v.(string)
	if !ok {
		return "", fmt.Errorf("%s must be a string", key)
	}
	return s, nil
}

func getInt(m map[string]any, key string, required bool) (int, error) {
	v, ok := m[key]
	if !ok {
		if required {
			return 0, fmt.Errorf("%s is required", key)
		}
		return 0, nil
	}
	switch v := v.(type) {
	case int:
		return v, nil
	case int32:
		return int(v), nil
	case int64:
		return int(v), nil
	case float32:
		return int(v), nil
	case float64:
		return int(v), nil
	default:
		return 0, fmt.Errorf("%s must be a int", key)

	}
}
