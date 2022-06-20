// Package mapstruct maps struct <-> JSON using came case <-> snake case
// also set default values based on struct tags
package mapstruct

// TODO: implement own version as we don't need much?

import (
	"regexp"
	"strings"

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

func camelCaseMap(m map[string]any) map[string]any {
	nm := map[string]any{}
	for k, v := range m {
		if vm, ok := v.(map[string]any); ok {
			v = camelCaseMap(vm)
		}
		nm[CamelToSnake(k)] = v
	}
	return nm
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

	return camelCaseMap(m), nil
}
