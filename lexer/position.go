package lexer

import "fmt"

type Position struct {
	Line   int
	Column int
}

func (p *Position) inc(r rune) {
	if isLineBreak(r) {
		p.Line++
		p.Column = 0
	} else {
		p.Column++
	}
}

func (p *Position) dec(r rune) {
	if isLineBreak(r) {
		p.Line--
		p.Column = -1
	} else {
		p.Column--
	}
}

func (p *Position) Str() string {
	return fmt.Sprint("(", p.Line, ":", p.Column, ")")
}

var NO_POS = Position{
	Line:   -1,
	Column: -1,
}
