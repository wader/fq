//go:generate goyacc -o parse.go -v "" parse.go.y
package ksexpr

// TODO:
// redo infix/prefix with generics somehow?
// string and eval interface?
// string runes?
// context timeout? not needed? no loops?
// (string).to_i and prefixes, transpiles to Number.parseInt in js etc so "0b" prefix not supported etc
// == and arrays, byte arrays?
// _ for current?
// _parent for parent
// id._sizeof
// id._bitsizeof
// sizeof<u4>

import (
	"fmt"
	"strings"
)

type Caller interface {
	KSExprCall(ns []string, name string, args []any) (any, error)
}

type Indexer interface {
	KSExprIndex(index int) (any, error)
}

type Node interface {
	Eval(input any) (any, error)
	fmt.Stringer
}

type ConstNode Token

func (c ConstNode) String() string { return c.Str }
func (c ConstNode) Eval(input any) (any, error) {

	// TODO: redo, BigInt

	return ToValue(c.V), nil
}

type IdentNode struct {
	NS   []Token
	Name Token
}

func (i IdentNode) ns() []string {
	var vs []string
	for _, t := range i.NS {
		vs = append(vs, t.Str)
	}
	return vs
}

func (i IdentNode) String() string {
	if i.NS != nil {
		return fmt.Sprintf("%s::%s", strings.Join(i.ns(), "::"), i.Name.Str)
	}
	return i.Name.Str
}

func (i IdentNode) Eval(input any) (any, error) {
	if c, ok := input.(Caller); ok {
		v, err := c.KSExprCall(i.ns(), i.Name.Str, nil)
		if err != nil {
			return nil, err
		}
		return ToValue(v), nil
	}
	return nil, noSuchKeyError{input, fmt.Sprintf("%s::%s", strings.Join(i.ns(), "::"), i.Name.Str)}
}

type ArrayNode []Node

func (a ArrayNode) String() string {
	var es []string
	for _, e := range a {
		es = append(es, e.String())
	}
	return strings.Join(es, ", ")
}

func (a ArrayNode) Eval(input any) (any, error) {
	var vs []any
	for _, e := range a {
		v, err := e.Eval(input)
		if err != nil {
			return nil, err
		}
		vs = append(vs, v)
	}
	return Array(vs), nil
}

type TermNode struct {
	T        Node
	Trailers []Node
}

func (t TermNode) String() string {
	return t.T.String()
}

func (t TermNode) Eval(input any) (any, error) {
	tv, err := t.T.Eval(input)
	if err != nil {
		return nil, err
	}

	for _, tr := range t.Trailers {
		tv, err = tr.Eval(tv)
		if err != nil {
			return nil, err
		}
	}

	return tv, nil
}

type TrailerCallNode struct {
	Name Token
	Args []Node
}

func (t TrailerCallNode) String() string {
	var as []string
	for _, a := range t.Args {
		as = append(as, a.String())
	}
	return fmt.Sprintf("%s(%s)", t.Name.Str, strings.Join(as, ", "))
}

func (t TrailerCallNode) Eval(input any) (any, error) {
	var av []any
	for _, a := range t.Args {
		v, err := a.Eval(input)
		if err != nil {
			return nil, err
		}
		av = append(av, v)
	}

	if tc, ok := input.(Caller); ok {
		v, err := tc.KSExprCall(nil, t.Name.Str, av)
		if err != nil {
			return nil, err
		}
		return ToValue(v), nil
	}

	return nil, noSuchMethodError{input, t.Name.Str}
}

type TrailerIndexNode struct {
	Expr Node
}

func (t TrailerIndexNode) String() string {
	return fmt.Sprintf("[%s]", t.Expr)
}

func (t TrailerIndexNode) Eval(input any) (any, error) {
	ev, err := t.Expr.Eval(input)
	if err != nil {
		return nil, err
	}
	v, ok := ToInt(ev)
	if !ok {
		return nil, invalidIndexError{input, ev}
	}

	if tc, ok := input.(Indexer); ok {
		v, err := tc.KSExprIndex(v)
		if err != nil {
			return nil, err
		}
		return ToValue(v), nil
	}

	return nil, notIndexableError{input, v}
}

type PrefixOpNode struct {
	Op   PrefixOp
	Expr Node
}

func (o PrefixOpNode) String() string {
	return fmt.Sprintf("%s %s", prefixOpNames[o.Op], o.Expr)
}

func (o PrefixOpNode) Eval(input any) (any, error) {
	v, err := o.Expr.Eval(input)
	if err != nil {
		return nil, err
	}

	fn, ok := prefixOpFn[o.Op]
	if ok {
		v := fn(v)
		if err, ok := v.(error); ok {
			return nil, err
		}
		return v, nil
	}

	panic(fmt.Sprintf("unknown prefix op %d", o.Op))
}

type InfixOpNode struct {
	LHS Node
	Op  InfixOp
	RHS Node
}

func (o InfixOpNode) String() string {
	return fmt.Sprintf("%s %s %s", o.LHS, infixOpNames[o.Op], o.RHS)
}

func (o InfixOpNode) Eval(input any) (any, error) {
	lhs, err := o.LHS.Eval(input)
	if err != nil {
		return nil, err
	}
	rhs, err := o.RHS.Eval(input)
	if err != nil {
		return nil, err
	}

	fn, ok := infixOpFn[o.Op]
	if ok {
		v := fn(lhs, rhs)
		if err, ok := v.(error); ok {
			return nil, err
		}
		return v, nil
	}

	panic(fmt.Sprintf("unknown infix op %d", o.Op))
}

type TernaryNode struct {
	Expr      Node
	TrueExpr  Node
	FalseExpr Node
}

func (t TernaryNode) String() string {
	return fmt.Sprintf("%s ? %s : %s", t.Expr, t.TrueExpr, t.FalseExpr)
}

func (t TernaryNode) Eval(input any) (any, error) {
	ev, err := t.Expr.Eval(input)
	if err != nil {
		return nil, err
	}
	b, ok := ev.(Boolean)
	if !ok {
		return nil, fmt.Errorf("non-boolean ternary condition %s", ev)
	}
	ce := t.FalseExpr
	if b {
		ce = t.TrueExpr
	}
	v, err := ce.Eval(input)
	if err != nil {
		return nil, err
	}

	return v, nil
}

func Parse(s string) (Node, error) {
	l := &yyLex{s: []rune(s)}
	n := yyParse(l)
	if n > 0 {
		return nil, l.err
	}
	return l.result, nil

}
