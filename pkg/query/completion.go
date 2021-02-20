package query

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/itchyny/gojq"
)

type CompletionType string

const (
	CompletionTypeIndex CompletionType = "index"
	CompletionTypeFunc  CompletionType = "func"
	CompletionTypeVar   CompletionType = "var"
	CompletionTypeNone  CompletionType = "none"
)

func BuildCompletionQuery(src string) (*gojq.Query, CompletionType, string) {
	if src == "" {
		return nil, CompletionTypeNone, ""
	}

	// HACK: if ending with "." or "$" append a test index that we remove later
	probePrefix := ""
	if len(src) > 0 && strings.HasSuffix(src, ".") || strings.HasSuffix(src, "$") {
		probePrefix = "x"
	}

	// log.Printf("src + probePrefix: %#+v\n", src+probePrefix)

	q, err := gojq.Parse(src + probePrefix)
	if err != nil {
		return nil, CompletionTypeNone, ""
	}

	cq, ct, prefix := transformToCompletionQuery(q)
	if probePrefix != "" {
		prefix = strings.TrimSuffix(prefix, probePrefix)
	}

	return cq, ct, prefix
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
	}

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
			q.Term.SuffixList = q.Term.SuffixList[0 : len(q.Term.SuffixList)-1]
			return q, CompletionTypeIndex, last.Index.Name
		}

		return nil, CompletionTypeNone, ""
	case gojq.TermTypeFunc:
		if len(q.Term.SuffixList) == 0 {
			if strings.HasPrefix(q.Term.Func.Name, "$") {
				return &gojq.Query{Term: &gojq.Term{Type: gojq.TermTypeIdentity}}, CompletionTypeVar, q.Term.Func.Name
			} else {
				return &gojq.Query{Term: &gojq.Term{Type: gojq.TermTypeIdentity}}, CompletionTypeFunc, q.Term.Func.Name
			}
		}

		return nil, CompletionTypeNone, ""
	default:
		return nil, CompletionTypeNone, ""

	}
}

func autoComplete(ctx context.Context, c interface{}, q *Query, line []rune, pos int) (newLine [][]rune, length int) {
	lineStr := string(line[0:pos])
	namesQuery, nameType, namePrefix := BuildCompletionQuery(lineStr)

	// log.Println("------")
	// log.Printf("namesQuery: %s\n", namesQuery)
	// log.Printf("namesType: %#+v\n", nameType)
	// log.Printf("namesPrefix: %#+v\n", namePrefix)

	if nameType == CompletionTypeNone {
		return [][]rune{}, pos
	}

	namesQueryStr := namesQuery.String()
	namePrefixReStr := jsonEscape("^" + regexp.QuoteMeta(namePrefix))

	src := ""
	switch nameType {
	case CompletionTypeIndex:
		src = fmt.Sprintf(`[[(%s) | keys?, _value_keys?] | add | unique | sort | .[] | strings | select(test(%s))]`,
			namesQueryStr, namePrefixReStr)
	case CompletionTypeFunc, CompletionTypeVar:
		src = fmt.Sprintf(`[%s | scope[] | select(test(%s))]`,
			namesQueryStr, namePrefixReStr)
	default:
		panic("unreachable")
	}

	log.Printf("src: %s\n", src)

	i, err := q.Eval(ctx, CompletionMode, c, src, DiscardOutput{})
	if err != nil {
		log.Printf("err: %#+v\n", err)
		return [][]rune{}, pos
	}

	var vss []interface{}
	for {
		vs, ok := i.Next()
		if !ok {
			break
		}
		log.Printf("vs: %#+v\n", vs)
		vss = append(vss, vs)
	}

	log.Printf("vss: %#+v\n", vss)

	shareLen := len(namePrefix)

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
