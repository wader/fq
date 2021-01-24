package query

import (
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
