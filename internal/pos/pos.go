package pos

import (
	"fmt"
	"strings"
)

func offsetToLineColumn(s string, offset int) (int, int) {
	co := 0
	line := 1
	for {
		no := strings.Index(s[co:], "\n")
		if no == -1 || co+no >= offset {
			return line, offset - co
		}
		co += no + 1
		line++
	}
}

type Pos struct {
	S      string
	Line   int
	Column int
}

func (p Pos) String() string {
	return fmt.Sprintf("%d:%d", p.Line, p.Column)
}

func NewFromOffset(s string, offset int) Pos {
	l, c := offsetToLineColumn(s, offset)
	return Pos{S: s, Line: l, Column: c}
}
