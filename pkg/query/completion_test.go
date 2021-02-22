package query_test

import (
	"fmt"
	"fq/pkg/query"
	"testing"
)

func TestBuildCompletionQuery(t *testing.T) {
	testCases := []struct {
		input          string
		expectedQuery  string
		expectedType   query.CompletionType
		expectedPrefix string
	}{
		{"", "", query.CompletionTypeNone, ""},
		{`.`, `.`, query.CompletionTypeIndex, ``},
		{`.`, `.`, query.CompletionTypeIndex, ``},
		{`.a`, `.`, query.CompletionTypeIndex, `a`},
		{`.a.`, `.a`, query.CompletionTypeIndex, ``},
		{`.a.b`, `.a`, query.CompletionTypeIndex, `b`},
		{`.a.b.`, `.a.b`, query.CompletionTypeIndex, ``},
		{` .a.b`, `.a`, query.CompletionTypeIndex, `b`},
		{`.a | .b`, `.a | .`, query.CompletionTypeIndex, `b`},
		{`.a | .b.c`, `.a | .b`, query.CompletionTypeIndex, `c`},
		{`.a[]`, ``, query.CompletionTypeNone, ``},
		{`.a[].b`, `.a[]`, query.CompletionTypeIndex, `b`},
		{`.a[].b.c`, `.a[].b`, query.CompletionTypeIndex, `c`},
		{`.a["b"]`, ``, query.CompletionTypeNone, ``},
		{`.a["b"].c`, `.a["b"]`, query.CompletionTypeIndex, `c`},
		{`.a[1:2]`, ``, query.CompletionTypeNone, ``},
		{`.a[1:2].c`, `.a[1:2]`, query.CompletionTypeIndex, `c`},
		{`a`, `.`, query.CompletionTypeFunc, `a`},
		{`a | b`, `a | .`, query.CompletionTypeFunc, `b`},
	}
	for _, tC := range testCases {
		t.Run(tC.input, func(t *testing.T) {
			actualQuery, actualType, actualPrefix := query.BuildCompletionQuery(tC.input)
			actualQueryStr := ""
			if actualQuery != nil {
				actualQueryStr = actualQuery.String()
			}

			if tC.expectedQuery != actualQueryStr {
				t.Errorf("expected query %q, got query %q", tC.expectedQuery, actualQueryStr)
			}
			if tC.expectedType != actualType {
				t.Errorf("expected type %s, got type %s", tC.expectedType, actualType)
			}
			if tC.expectedPrefix != actualPrefix {
				t.Errorf("expected prefix %s, got prefix %q", tC.expectedPrefix, actualPrefix)
			}
		})
	}
}

func TestSharedPrefix(t *testing.T) {
	testCases := []struct {
		vs     []string
		shared string
	}{
		{vs: []string{}, shared: ""},
		{vs: []string{""}, shared: ""},
		{vs: []string{"a"}, shared: "a"},
		{vs: []string{"a", "a"}, shared: "a"},
		{vs: []string{"aa", "ab"}, shared: "a"},
		{vs: []string{"aa", "aa"}, shared: "aa"},
		{vs: []string{"abc", "abc", "ab"}, shared: "ab"},
		{vs: []string{"ab", "abc", "ab"}, shared: "ab"},
		{vs: []string{"abc", "abc", "abd"}, shared: "ab"},
	}
	for _, tC := range testCases {
		t.Run(fmt.Sprintf("%v", tC.vs), func(t *testing.T) {
			actual := query.SharedPrefix(tC.vs)
			if tC.shared != actual {
				t.Errorf("expected %v, got %v", tC.shared, actual)
			}
		})
	}
}
