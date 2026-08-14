// Package mapstruct maps struct <-> JSON using came case <-> snake case
// also set default values based on struct tags
package mapstruct

// TODO: implement own version as we don't need much?

import (
	"reflect"
	"strings"
	"sync"

	"github.com/creasty/defaults"
	"github.com/mitchellh/mapstructure"
)

// "AaaBbb" -> "aaa_bbb"
func CamelToSnake(s string) string {
	if v, ok := camelToSnakeCache.Load(s); ok {
		s, _ := v.(string)
		return s
	}
	r := camelToSnakeSlow(s)
	camelToSnakeCache.Store(s, r)
	return r
}

var camelToSnakeCache sync.Map

func camelToSnakeSlow(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	for i := range len(s) {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			if i > 0 && s[i-1] >= 'a' && s[i-1] <= 'z' {
				b.WriteByte('_')
			}
			b.WriteByte(c + ('a' - 'A'))
		} else {
			b.WriteByte(c)
		}
	}
	return b.String()
}

func ToStruct(m any, v any) error {
	_ = defaults.Set(v)
	ms, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		MatchName: func(mapKey, fieldName string) bool {
			return CamelToSnake(fieldName) == mapKey
		},
		TagName: "mapstruct",
		Squash:  true,
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
		Squash: true,
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

func StructFieldsSquashed(v any) []reflect.StructField {
	var fs []reflect.StructField

	var r func(v any)
	r = func(v any) {
		st := reflect.TypeOf(v)
		sv := reflect.ValueOf(v)
		for i := range st.NumField() {
			f := st.Field(i)
			if f.Anonymous {
				r(sv.Field(i).Interface())
			} else {
				fs = append(fs, f)
			}
		}
	}
	r(v)

	return fs
}
