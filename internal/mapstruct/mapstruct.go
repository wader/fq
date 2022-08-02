// Package mapstruct maps struct <-> JSON using came case <-> snake case
// also set default values based on struct tags
package mapstruct

// TODO: implement own version as we don't need much?

import (
	"regexp"
	"strings"

	"github.com/creasty/defaults"
	"github.com/mitchellh/mapstructure"
)

var camelToSnakeRe = regexp.MustCompile(`[[:lower:]][[:upper:]]`)

// "AaaBbb" -> "aaa_bbb"
func CamelToSnake(s string) string {
	return strings.ToLower(camelToSnakeRe.ReplaceAllStringFunc(s, func(s string) string {
		return s[0:1] + "_" + s[1:2]
	}))
}

func ToStruct(m any, v any) error {
	_ = defaults.Set(v)
	ms, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		MatchName: func(mapKey, fieldName string) bool {
			return CamelToSnake(fieldName) == mapKey
		},
		TagName: "mapstruct",
		Result:  v,
	})
	if err != nil {
		return err
	}
	if err := ms.Decode(m); err != nil {
		return err
	}

	return nil
}

func CamelCase(v any) any {
	switch vv := v.(type) {
	case map[string]any:
		n := map[string]any{}
		for k, v := range vv {
			n[CamelToSnake(k)] = CamelCase(v)
		}
		return n
	case []any:
		n := make([]any, len(vv))
		for i, v := range vv {
			n[i] = CamelCase(v)
		}
		return n
	}
	return v
}

func ToMap(v any) (map[string]any, error) {
	m := map[string]any{}
	ms, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result: &m,
	})
	if err != nil {
		return nil, err
	}
	if err := ms.Decode(v); err != nil {
		return nil, err
	}
	m, ok := CamelCase(m).(map[string]any)
	if !ok {
		panic("not map")
	}

	return m, nil
}
