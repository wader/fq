package shquote

import (
	"strings"
)

func Split(s string) []string {
	type splitState int
	const (
		Word splitState = iota
		SingleQuote
		DoubleQuote
		DoubleQuoteEscape
	)
	type word string
	type singleQuote string
	type doubleQuote string
	type delim struct{}
	var tokens []interface{}

	sb := &strings.Builder{}
	ss := Word
	for _, r := range s {
		switch ss {
		case DoubleQuoteEscape:
			if r != '"' {
				sb.WriteRune('\\')
			}
			sb.WriteRune(r)
			ss = DoubleQuote
			continue
		}

		switch r {
		case '\'':
			switch ss {
			case Word:
				tokens = append(tokens, word(sb.String()))
				sb.Reset()
				ss = SingleQuote
			case SingleQuote:
				tokens = append(tokens, singleQuote(sb.String()))
				sb.Reset()
				ss = Word
			default:
				sb.WriteRune(r)
			}
		case '"':
			switch ss {
			case Word:
				tokens = append(tokens, word(sb.String()))
				sb.Reset()
				ss = DoubleQuote
			case DoubleQuote:
				tokens = append(tokens, doubleQuote(sb.String()))
				sb.Reset()
				ss = Word
			default:
				sb.WriteRune(r)
			}
		case '\\':
			switch ss {
			case DoubleQuote:
				ss = DoubleQuoteEscape
			default:
				sb.WriteRune(r)
			}
		case ' ':
			switch ss {
			case Word:
				tokens = append(tokens, word(sb.String()))
				sb.Reset()
				tokens = append(tokens, delim{})
			default:
				sb.WriteRune(r)
			}
		default:
			sb.WriteRune(r)
		}
	}
	if sb.Len() > 0 {
		switch ss {
		case Word:
			tokens = append(tokens, word(sb.String()))
		case SingleQuote:
			tokens = append(tokens, singleQuote(sb.String()))
		case DoubleQuote:
			tokens = append(tokens, doubleQuote(sb.String()))
		}
	}
	tokens = append(tokens, delim{})

	var as []string
	sb.Reset()
	for _, t := range tokens {
		switch tt := t.(type) {
		case doubleQuote:
			sb.WriteString(string(tt))
		case singleQuote:
			sb.WriteString(string(tt))
		case word:
			sb.WriteString(string(tt))
		case delim:
			if sb.Len() > 0 {
				as = append(as, sb.String())
				sb.Reset()
			}
		}
	}

	return as
}
