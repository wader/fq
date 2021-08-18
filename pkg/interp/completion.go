package interp

import (
	"context"
	"fmt"
	"strings"

	"github.com/wader/gojq"
)

type CompletionType string

const (
	CompletionTypeIndex CompletionType = "index"
	CompletionTypeFunc  CompletionType = "function"
	CompletionTypeVar   CompletionType = "variable"
	CompletionTypeNone  CompletionType = "none"
	CompletionTypeError CompletionType = "error"
)

func BuildCompletionQuery(src string) (*gojq.Query, CompletionType, string) {
	if src == "" {
		return nil, CompletionTypeError, ""
	}

	// HACK: if ending with "." or "$" append a test index that we remove later
	probePrefix := ""
	if len(src) > 0 && strings.HasSuffix(src, ".") || strings.HasSuffix(src, "$") {
		probePrefix = "x"
	}

	// log.Printf("src + probePrefix: %#+v\n", src+probePrefix)

	q, err := gojq.Parse(src + probePrefix)
	if err != nil {
		return nil, CompletionTypeError, ""
	}

	cq, ct, prefix := transformToCompletionQuery(q)
	if probePrefix != "" {
		prefix = strings.TrimSuffix(prefix, probePrefix)
	}

	if ct == CompletionTypeNone {
		return cq, ct, ""
	}

	// [.[] | cq | add]
	return &gojq.Query{
		Left: &gojq.Query{
			Term: &gojq.Term{
				Type: gojq.TermTypeArray,
				Array: &gojq.Array{
					Query: &gojq.Query{
						Left: &gojq.Query{
							Term: &gojq.Term{
								Type:       gojq.TermTypeIdentity,
								SuffixList: []*gojq.Suffix{{Iter: true}},
							},
						},
						Op:    gojq.OpPipe,
						Right: cq,
					},
				},
			},
		},
		Op: gojq.OpPipe,
		Right: &gojq.Query{
			Term: &gojq.Term{
				Type: gojq.TermTypeFunc,
				Func: &gojq.Func{Name: "add"},
			},
		},
	}, ct, prefix
}

// find the right most term that is completeable
// return a query to find possible names and a prefix to filter by
func transformToCompletionQuery(q *gojq.Query) (*gojq.Query, CompletionType, string) {
	// pipe, eq etc
	if q.Right != nil {
		r, ct, prefix := transformToCompletionQuery(q.Right)
		if r == nil {
			return nil, ct, prefix
		}
		q.Right = r
		return q, ct, prefix
	}

	keysFuncName := func(name string) string {
		if strings.HasPrefix(name, "_") {
			return "_extkeys"
		}
		return "keys"
	}

	optFunc := func(name string) *gojq.Query {
		return &gojq.Query{
			Term: &gojq.Term{
				Type:       gojq.TermTypeFunc,
				Func:       &gojq.Func{Name: name},
				SuffixList: []*gojq.Suffix{{Optional: true}},
			},
		}
	}

	// ... as ...
	if q.Term.SuffixList != nil {
		last := q.Term.SuffixList[len(q.Term.SuffixList)-1]
		if last.Bind != nil {
			r, ct, prefix := transformToCompletionQuery(last.Bind.Body)
			if r == nil {
				return nil, ct, prefix
			}
			last.Bind.Body = r
			return q, ct, prefix
		}
		if last.Index != nil && last.Index.Name != "" {
			prefix := last.Index.Name
			last.Index = nil
			return &gojq.Query{
				Left:  q,
				Op:    gojq.OpPipe,
				Right: optFunc(keysFuncName(prefix)),
			}, CompletionTypeIndex, prefix
		}
	}

	switch q.Term.Type { //nolint:exhaustive
	case gojq.TermTypeIdentity:
		return &gojq.Query{
			Left:  q,
			Op:    gojq.OpPipe,
			Right: optFunc(keysFuncName("")),
		}, CompletionTypeIndex, ""
	case gojq.TermTypeIndex:
		if len(q.Term.SuffixList) == 0 {
			if q.Term.Index.Start == nil {
				return &gojq.Query{
					Left:  &gojq.Query{Term: &gojq.Term{Type: gojq.TermTypeIdentity}},
					Op:    gojq.OpPipe,
					Right: optFunc(keysFuncName(q.Term.Index.Name)),
				}, CompletionTypeIndex, q.Term.Index.Name
			}
			return q, CompletionTypeNone, ""
		}

		last := q.Term.SuffixList[len(q.Term.SuffixList)-1]
		if last.Index != nil && last.Index.Start == nil {
			q.Term.SuffixList = q.Term.SuffixList[0 : len(q.Term.SuffixList)-1]
			return &gojq.Query{
				Left:  q,
				Op:    gojq.OpPipe,
				Right: optFunc(keysFuncName(last.Index.Name)),
			}, CompletionTypeIndex, last.Index.Name
		}

		return q, CompletionTypeNone, ""
	case gojq.TermTypeFunc:
		if len(q.Term.SuffixList) == 0 {
			if strings.HasPrefix(q.Term.Func.Name, "$") {
				return optFunc("scope"), CompletionTypeVar, q.Term.Func.Name
			} else {
				return optFunc("scope"), CompletionTypeFunc, q.Term.Func.Name
			}
		}

		return q, CompletionTypeNone, ""
	default:
		return q, CompletionTypeNone, ""

	}
}

func completeTrampoline(ctx context.Context, completeFn string, c interface{}, i *Interp, line string, pos int) (newLine []string, shared int, err error) {
	vs, err := i.EvalFuncValues(ctx, CompletionMode, c, completeFn, []interface{}{line, pos}, DiscardOutput{Ctx: ctx})
	if err != nil {
		return nil, pos, err
	}
	if len(vs) < 1 {
		return nil, pos, fmt.Errorf("no values")
	}
	v := vs[0]
	if vErr, ok := v.(error); ok {
		return nil, pos, vErr
	}

	// {abc: 123, abd: 123} | complete(".ab"; 3) will return {prefix: "ab", names: ["abc", "abd"]}

	var names []string
	var prefix string
	cm, ok := v.(map[string]interface{})
	if !ok {
		return nil, pos, fmt.Errorf("%v: complete function return value not an object", cm)
	}
	if namesV, ok := cm["names"].([]interface{}); ok {
		for _, name := range namesV {
			names = append(names, name.(string))
		}
	} else {
		return nil, pos, fmt.Errorf("%v: names missing in complete return object", cm)
	}
	if prefixV, ok := cm["prefix"]; ok {
		prefix, _ = prefixV.(string)
	} else {
		return nil, pos, fmt.Errorf("%v: prefix missing in complete return object", cm)
	}

	if len(names) == 0 {
		return nil, pos, nil
	}

	sharedLen := len(prefix)

	return names, sharedLen, nil
}
