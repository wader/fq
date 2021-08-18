package interp_test

import (
	"testing"

	"github.com/wader/fq/pkg/interp"
)

func TestBuildCompletionQuery(t *testing.T) {
	testCases := []struct {
		input          string
		expectedQuery  string
		expectedType   interp.CompletionType
		expectedPrefix string
	}{
		{"", "", interp.CompletionTypeError, ""},
		{`.`, `[.[] | . | keys?] | add`, interp.CompletionTypeIndex, ``},
		{`.`, `[.[] | . | keys?] | add`, interp.CompletionTypeIndex, ``},
		{`.a`, `[.[] | . | keys?] | add`, interp.CompletionTypeIndex, `a`},
		{`.a.`, `[.[] | .a | keys?] | add`, interp.CompletionTypeIndex, ``},
		{`.a.b`, `[.[] | .a | keys?] | add`, interp.CompletionTypeIndex, `b`},
		{`.a.b.`, `[.[] | .a.b | keys?] | add`, interp.CompletionTypeIndex, ``},
		{`.a.b`, `[.[] | .a | keys?] | add`, interp.CompletionTypeIndex, `b`},
		{`.a | .b`, `[.[] | .a | . | keys?] | add`, interp.CompletionTypeIndex, `b`},
		{`.a | .b.c`, `[.[] | .a | .b | keys?] | add`, interp.CompletionTypeIndex, `c`},
		{`.a[]`, `.a[]`, interp.CompletionTypeNone, ``},
		{`.a[].b`, `[.[] | .a[] | keys?] | add`, interp.CompletionTypeIndex, `b`},
		{`.a[].b.c`, `[.[] | .a[].b | keys?] | add`, interp.CompletionTypeIndex, `c`},
		{`.a["b"]`, `.a["b"]`, interp.CompletionTypeNone, ``},
		{`.a["b"].c`, `[.[] | .a["b"] | keys?] | add`, interp.CompletionTypeIndex, `c`},
		{`.a[1:2]`, `.a[1:2]`, interp.CompletionTypeNone, ``},
		{`.a[1:2].c`, `[.[] | .a[1:2] | keys?] | add`, interp.CompletionTypeIndex, `c`},
		{`a`, `[.[] | scope?] | add`, interp.CompletionTypeFunc, `a`},
		{`a | b`, `[.[] | a | scope?] | add`, interp.CompletionTypeFunc, `b`},
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
