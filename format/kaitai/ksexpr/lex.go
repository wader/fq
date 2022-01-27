package ksexpr

// TODO: end many times?
// TODO: cleanup
// TODO: redo error handling and messages

import (
	"fmt"
	"strconv"
	"strings"
)

type parseError struct {
	id    int
	token Token
}

func (p parseError) Error() string {
	switch p.id {
	case -1:
		return "unexpected end"
	case tokUnterminatedString:
		return fmt.Sprintf("no closing quote for string starting at position %d", p.token.Span.Start+1)
	default:
		return fmt.Sprintf("unexpected %s at position %d", p.token.Str, p.token.Span.Start+1)
	}
}

type Span struct {
	Start int
	Stop  int
}

type Token struct {
	Str  string
	Span Span
	V    any
}

type yyLex struct {
	s           []rune
	start       int
	pos         int
	lastTokenId int
	lastToken   Token

	result Node
	err    error
}

var lexTokenNames = map[int]string{
	tokNumber:    "tokNumber",
	tokIdent:     "tokIdent",
	tokString:    "tokString",
	tokLessEq:    "tokLessEq",
	tokGreaterEq: "tokGreaterEq",
	tokEqEq:      "tokEqEq",
	tokNotEq:     "tokNotEq",
	tokBSL:       "tokBSL",
	tokBSR:       "tokBSR",
	tokBAnd:      "tokBAnd",
	tokBOr:       "tokBOr",
	tokNot:       "tokNot",
	tokAnd:       "tokAnd",
	tokOr:        "tokOr",
	tokTrue:      "tokTrue",
	tokFalse:     "tokFalse",
}

func lexTokenToString(t int) string {
	if s, ok := lexTokenNames[t]; ok {
		return s
	} else if t == -1 {
		return "end"
	}
	return fmt.Sprintf("'%c'", t)
}

func (l *yyLex) tok() string { return string(l.s[l.start:l.pos]) }
func (l *yyLex) end() bool   { return l.pos >= len(l.s) }
func (l *yyLex) peek(n int) string {
	if l.pos+n >= len(l.s) {
		n = len(l.s) - l.pos
	}
	return string(l.s[l.pos : l.pos+n])
}

func (l *yyLex) consume(n int) string {
	s := l.peek(n)
	l.pos += n
	return s
}

func (l *yyLex) consumeWhile(f func(s string) bool) {
	for !l.end() {
		r := l.peek(1)
		if !f(r) {
			break
		}
		l.consume(1)
	}
}

func isOctal(s string) bool { return s[0] >= '0' && s[0] <= '7' }
func isDigit(s string) bool { return s[0] >= '0' && s[0] <= '9' }
func isHexDigit(s string) bool {
	return isDigit(s) ||
		(s[0] >= 'a' && s[0] <= 'f') ||
		(s[0] >= 'A' && s[0] <= 'F')
}
func isHexDigitRest(s string) bool {
	return isHexDigit(s) || s[0] == '_'
}

func isIdentStart(s string) bool { return s[0] == '_' || (s[0] >= 'a' && s[0] <= 'z') }
func isIdentRest(s string) bool  { return isIdentStart(s) || (s[0] >= '0' && s[0] <= '9') }

var escapeMap = map[string]string{
	`0`: "\x00",
	`a`: "\a",
	`b`: "\b",
	`t`: "\t",
	`n`: "\n",
	`v`: "\v",
	`f`: "\f",
	`r`: "\r",
	`e`: "\x1b",
	`"`: `"`,
	`'`: "'",
	`\`: `\`,
}

func (l *yyLex) consumeString(q string, escape bool) (string, int) {
	sb := &strings.Builder{}

	l.consume(1)
	for {
		switch l.peek(1) {
		case q:
			l.consume(1)
			return sb.String(), tokString
		case "":
			return "", tokUnterminatedString
		case `\`:
			l.consume(1)
			if escape {
				// TODO:
				// a\0b -> a 0 b
				// 1\02 -> 1 2
				p := l.peek(1)
				switch {
				case escapeMap[p] != "":
					sb.WriteString(escapeMap[p])
					l.consume(1)
				case isDigit(p):
					// \1, \12, \123 octal
					// TODO: eof
					ob := &strings.Builder{}
					for i := 0; i < 3 && isOctal(l.peek(1)); i++ {
						ob.WriteString(l.consume(1))
					}
					n, err := strconv.ParseInt(ob.String(), 8, 64)
					if err != nil {
						return "", tokError // TODO: unknown octal
					}
					sb.WriteRune(rune(n))
				case p == `u`:
					// \uffff utf16 unicode codepoint
					// TODO: surrogate pair
					// TODO: eof
					l.consume(1)
					n, err := strconv.ParseInt(l.consume(4), 16, 64)
					if err != nil {
						return "", tokError // TODO: unknown octal
					}
					sb.WriteRune(rune(n))
				default:
					return "", tokError // TODO: unknown escape
				}
			} else {
				sb.WriteString(`\`)
			}
		default:
			sb.WriteString(l.consume(1))
		}
	}
}

func (l *yyLex) Lex(lval *yySymType) (t int) {
	l.consumeWhile(func(s string) bool { return s == " " || s == "\t" || s == "\n" })

	l.start = l.pos
	defer func() {
		lval.token.Str = l.tok()
		lval.token.Span = Span{Start: l.start, Stop: l.pos}
		l.lastTokenId = t
		l.lastToken = lval.token
	}()

	if l.end() {
		return -1
	}

	p := l.peek(2)
	switch {
	case p == "<=":
		l.consume(2)
		return tokLessEq
	case p == "<<":
		l.consume(2)
		return tokBSL
	case p == ">=":
		l.consume(2)
		return tokGreaterEq
	case p == ">>":
		l.consume(2)
		return tokBSR
	case p == "==":
		l.consume(2)
		return tokEqEq
	case p == "!=":
		l.consume(2)
		return tokNotEq
	case p == "::":
		l.consume(2)
		return tokColonColon
	case isDigit(string(p[0])):
		switch p {
		case "0x",
			"0X",
			"0b",
			"0B",
			"0o",
			"0O":
			l.consume(2)
		}
		l.consumeWhile(isHexDigitRest) // 0-9,a-f,_

		// TODO: 0x123.123
		// TODO: 123E5
		// TODO: invalid check?
		if l.peek(1) == "." {
			l.consume(1)
			l.consumeWhile(isDigit)
			switch l.peek(1) {
			case "e", "E":
				l.consume(1)
				l.consumeWhile(isDigit)
			}
			f, err := strconv.ParseFloat(l.tok(), 64)
			if err != nil {
				return tokNumber
			}
			lval.token.V = f
			return tokNumber
		}

		n, _, err := strToInteger(l.tok(), 0, true)
		if err != nil {
			return tokNumber
		}
		lval.token.V = n

		return tokNumber
	case isIdentStart(p):
		l.consumeWhile(isIdentRest)
		switch l.tok() {
		case "true":
			lval.token.V = true
			return tokTrue
		case "false":
			lval.token.V = false
			return tokFalse
		case "not":
			return tokNot
		case "and":
			return tokAnd
		case "or":
			return tokOr
		default:
			return tokIdent
		}
	case p[0] == '"', p[0] == '\'':
		q := string(p[0])
		escape := q == `"`
		s, t := l.consumeString(q, escape)
		lval.token.V = s
		return t
	default:
		return int(l.consume(1)[0])
	}
}

func (l *yyLex) Error(s string) {
	l.err = parseError{
		id:    l.lastTokenId,
		token: l.lastToken,
	}
}

type LexToken struct {
	Token Token
	Err   error
	Name  string
}

func Lex(s string) []LexToken {
	var ts []LexToken

	l := &yyLex{s: []rune(s)}
	var lval yySymType
	for {
		t := l.Lex(&lval)
		ts = append(ts, LexToken{
			Token: lval.token,
			Err:   l.err,
			Name:  lexTokenToString(t),
		})
		if t == -1 {
			break
		}
	}

	return ts
}
