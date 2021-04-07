package lexer

const EOF_LITERAL rune = rune(0)
const LF_LITERAL rune = rune(10)
const COMMENT_LITERAL rune = ';'

type TokenContainer struct {
	Pos Position
	Tok Token
	Str string
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
