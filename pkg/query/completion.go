package query

import (
	"context"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/itchyny/gojq"
)

type CompletionType string

const (
	CompletionTypeIndex CompletionType = "index"
	CompletionTypeFunc  CompletionType = "func"
	CompletionTypeNone  CompletionType = "none"
)

func BuildCompletionQuery(src string) (*gojq.Query, CompletionType, string) {
	if src == "" {
		return nil, CompletionTypeNone, ""
	}

	// HACK: if ending with "." append a test index that we remove later
	probePrefix := ""
	if len(src) > 0 && strings.HasSuffix(src, ".") {
		probePrefix = "x"
	}

	q, err := gojq.Parse(src + probePrefix)
	if err != nil {
		return nil, CompletionTypeNone, ""
	}

	cq, ct, prefix := buildCompletionQuery(q)
	if prefix != "" && probePrefix != "" {
		prefix = strings.TrimPrefix(prefix, probePrefix)
	}

	return cq, ct, prefix
}

// find the right most term that is completeable
// return a query to find possible names and a prefix to filter by
func buildCompletionQuery(q *gojq.Query) (*gojq.Query, CompletionType, string) {
	switch q.Op {
	case gojq.OpPipe:
		r, ct, prefix := buildCompletionQuery(q.Right)
		if r == nil {
			return nil, ct, prefix
		}
		qc := *q
		qc.Right = r
		return &qc, ct, prefix
	default:
		switch q.Term.Type {
		case gojq.TermTypeIdentity:
			return q, CompletionTypeIndex, ""
		case gojq.TermTypeIndex:
			if len(q.Term.SuffixList) == 0 {
				if q.Term.Index.Start == nil {
					return &gojq.Query{Term: &gojq.Term{Type: gojq.TermTypeIdentity}}, CompletionTypeIndex, q.Term.Index.Name
				}
				return nil, CompletionTypeNone, ""
			}

			last := q.Term.SuffixList[len(q.Term.SuffixList)-1]
			if last.Index != nil && last.Index.Start == nil {
				qc := *q
				tc := *q.Term
				qc.Term = &tc
				qc.Term.SuffixList = qc.Term.SuffixList[0 : len(qc.Term.SuffixList)-1]
				return &qc, CompletionTypeIndex, last.Index.Name
			}

			return nil, CompletionTypeNone, ""
		case gojq.TermTypeFunc:
			if len(q.Term.SuffixList) == 0 {
				return nil, CompletionTypeFunc, q.Term.Func.Name
			}

			// TODO: refactor to share with index
			last := q.Term.SuffixList[len(q.Term.SuffixList)-1]
			if last.Index != nil && last.Index.Start == nil {
				qc := *q
				tc := *q.Term
				qc.Term = &tc
				qc.Term.SuffixList = qc.Term.SuffixList[0 : len(qc.Term.SuffixList)-1]
				return &qc, CompletionTypeIndex, last.Index.Name
			}
			return nil, CompletionTypeNone, ""

		default:
			return nil, CompletionTypeNone, ""
		}
	}
}

func autoComplete(ctx context.Context, q *Query, line []rune, pos int) (newLine [][]rune, length int) {
	lineStr := string(line[0:pos])
	namesQuery, namesType, namesPrefix := BuildCompletionQuery(lineStr)

	// log.Println("------")
	// log.Printf("namesQuery: %s\n", namesQuery)
	// log.Printf("namesType: %#+v\n", namesType)
	// log.Printf("namesPrefix: %#+v\n", namesPrefix)

	src := ""
	switch namesType {
	case CompletionTypeNone:
		return [][]rune{}, pos
	case CompletionTypeIndex:
		namesQueryStr := namesQuery.String()
		src = fmt.Sprintf(`[[(%s) | keys?, _value_keys?] | add | unique | sort | .[] | strings | select(test("^%s"))]`, namesQueryStr, namesPrefix)
	case CompletionTypeFunc:
		src = fmt.Sprintf(`[[builtins[] | split("/") | .[0]] | unique | sort | .[] | select(test("^%s"))]`, namesPrefix)
	default:
		panic("unreachable")
	}

	// log.Printf("src: %#+v\n", src)

	vss, err := q.Run(ctx, src, ioutil.Discard)
	if err != nil {
		// log.Printf("err: %#+v\n", err)
		return [][]rune{}, pos
	}

	shareLen := len(namesPrefix)

	vs := vss[0].([]interface{})
	var names []string
	for _, v := range vs {
		v, _ := v.(string)
		if v == "" {
			continue
		}
		names = append(names, v[shareLen:])
	}

	if len(names) <= 1 {
		shareLen = 0
	}

	// log.Printf("shareLen: %#+v\n", shareLen)
	// log.Printf("names: %#+v\n", names)

	var runeNames [][]rune
	for _, n := range names {
		runeNames = append(runeNames, []rune(n))
	}

	return runeNames, shareLen
}
