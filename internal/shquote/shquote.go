package shquote

import (
	"strings"
)

type Token struct {
	Start     int
	End       int
	Str       string
	Separator bool
}

func Parse(s string) []Token {
	type splitState int
	const (
		whitespace splitState = iota
		escape
		word
		singleQuote
		doubleQuote
		doubleQuoteEscape
	)
	var tokens []Token

	sb := &strings.Builder{}
	ss := whitespace
	start := 0
	for i, r := range s {
		switch ss {
		case escape:
			sb.WriteRune(r)
			ss = word
			continue
		case doubleQuoteEscape:
			if r != '"' {
				sb.WriteRune('\\')
			}
			sb.WriteRune(r)
			ss = doubleQuote
			continue
		}

		switch r {
		case '\'':
			switch ss {
			case whitespace:
				ss = singleQuote
				start = i
			case word:
				tokens = append(tokens, Token{Start: start, End: i, Str: sb.String()})
				sb.Reset()
				ss = singleQuote
				start = i
			case singleQuote:
				tokens = append(tokens, Token{Start: start, End: i, Str: sb.String()})
				sb.Reset()
				ss = whitespace
				start = i
			default:
				sb.WriteRune(r)
			}
		case '"':
			switch ss {
			case whitespace:
				ss = doubleQuote
				start = i
			case word:
				tokens = append(tokens, Token{Start: start, End: i, Str: sb.String()})
				sb.Reset()
				ss = doubleQuote
				start = i
			case doubleQuote:
				tokens = append(tokens, Token{Start: start, End: i, Str: sb.String()})
				sb.Reset()
				ss = whitespace
				start = i
			default:
				sb.WriteRune(r)
			}
		case '\\':
			switch ss {
			case whitespace, word:
				ss = escape
				start = i
			case doubleQuote:
				ss = doubleQuoteEscape
				start = i
			default:
				sb.WriteRune(r)
			}
		case ' ':
			switch ss {
			case whitespace:
				tokens = append(tokens, Token{Separator: true})
			case word:
				tokens = append(tokens, Token{Start: start, End: i, Str: sb.String()})
				tokens = append(tokens, Token{Separator: true})
				sb.Reset()
				ss = whitespace
				start = i
			default:
				sb.WriteRune(r)
			}
		default:
			switch ss {
			case whitespace:
				ss = word
				start = i
			}
			sb.WriteRune(r)
		}
	}
	if sb.Len() > 0 {
		switch ss {
		case word, singleQuote, doubleQuote:
			tokens = append(tokens, Token{Start: start, End: len(s), Str: sb.String()})
		}
	}
	tokens = append(tokens, Token{Separator: true})

	var filtered []Token
	prevSep := true
	for _, t := range tokens {
		if prevSep && t.Separator {
			prevSep = true
			continue
		}
		prevSep = t.Separator
		filtered = append(filtered, t)
	}
	if len(filtered) > 0 && filtered[len(filtered)-1].Separator {
		filtered = filtered[0 : len(filtered)-1]
	}

	return filtered
}

func Split(s string) []string {
	var ss []string
	sb := &strings.Builder{}

	sb.Reset()
	for _, t := range Parse(s) {
		sb.WriteString(t.Str)
		if t.Separator && sb.Len() > 0 {
			ss = append(ss, sb.String())
			sb.Reset()
		}
	}

	if sb.Len() > 0 {
		ss = append(ss, sb.String())
	}

	return ss
}
