package lexer

import (
	"fmt"

	"github.com/bodneyc/knit-and-go/util"
	"github.com/sirupsen/logrus"
)

const EOF_LITERAL rune = rune(0)
const LF_LITERAL rune = rune(10)
const COMMENT_LITERAL rune = ';'

type TokenContainer struct {
	Pos Position
	Tok Token
	Str string
}

func (t *TokenContainer) Fields() logrus.Fields {
	s := t.Str
	if len(s) > 6 {
		s = fmt.Sprint(s[0:3], "...")
	}
	return logrus.Fields{
		"pos": t.Pos,
		"tok": t.Tok,
		"str": s,
	}
}

func (t TokenContainer) String() string {
	return fmt.Sprintf(
		"{%d:%d, %d, \"%s\"}",
		t.Pos.Line, t.Pos.Column, t.Tok, util.JsonEscape(t.Str))
}

func NewTokenContainer(pos Position, tok Token, str string) TokenContainer {
	return TokenContainer{
		Pos: pos,
		Tok: tok,
		Str: str,
	}
}

/// Enum: Lexing tokens
type Token int

const (
	ILLEGAL_T Token = iota
	EOF_T
	NEXT_SOURCE_T
	WHITE_SPACE_T
	NEW_LINE_T

	COMMENT_T
	IDENTIFIER_T
	NUMERIC_T

	MINUS_T
	ASTERISK_T
	FEET_T
	INCHES_T

	COMMA_T
	LEFT_PAREN_T
	RIGHT_PAREN_T
	LEFT_BRACE_T
	RIGHT_BRACE_T
	EQUALS_T
	ALIAS_T
)
