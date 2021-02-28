package interp_test

import (
	"fq/pkg/interp"
	"testing"
)

func TestBuildCompletionQuery(t *testing.T) {
	testCases := []struct {
		input          string
		expectedQuery  string
		expectedType   interp.CompletionType
		expectedPrefix string
	}{
		{"", "", interp.CompletionTypeNone, ""},
		{`.`, `.`, interp.CompletionTypeIndex, ``},
		{`.`, `.`, interp.CompletionTypeIndex, ``},
		{`.a`, `.`, interp.CompletionTypeIndex, `a`},
		{`.a.`, `.a`, interp.CompletionTypeIndex, ``},
		{`.a.b`, `.a`, interp.CompletionTypeIndex, `b`},
		{`.a.b.`, `.a.b`, interp.CompletionTypeIndex, ``},
		{` .a.b`, `.a`, interp.CompletionTypeIndex, `b`},
		{`.a | .b`, `.a | .`, interp.CompletionTypeIndex, `b`},
		{`.a | .b.c`, `.a | .b`, interp.CompletionTypeIndex, `c`},
		{`.a[]`, ``, interp.CompletionTypeNone, ``},
		{`.a[].b`, `.a[]`, interp.CompletionTypeIndex, `b`},
		{`.a[].b.c`, `.a[].b`, interp.CompletionTypeIndex, `c`},
		{`.a["b"]`, ``, interp.CompletionTypeNone, ``},
		{`.a["b"].c`, `.a["b"]`, interp.CompletionTypeIndex, `c`},
		{`.a[1:2]`, ``, interp.CompletionTypeNone, ``},
		{`.a[1:2].c`, `.a[1:2]`, interp.CompletionTypeIndex, `c`},
		{`a`, `.`, interp.CompletionTypeFunc, `a`},
		{`a | b`, `a | .`, interp.CompletionTypeFunc, `b`},
	}
	for _, tC := range testCases {
		t.Run(tC.input, func(t *testing.T) {
			actualQuery, actualType, actualPrefix := interp.BuildCompletionQuery(tC.input)
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
