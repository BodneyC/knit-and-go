package lexer

import (
	"bufio"
	"bytes"
	"io"

	// "strings"
	"unicode"

	log "github.com/sirupsen/logrus"
)

type Lexer struct {
	reader     *bufio.Reader
	pos        Position
	override   bool
	overridden TokenContainer
}

func NewLexer(reader io.Reader) *Lexer {
	return &Lexer{
		reader:     bufio.NewReader(reader),
		pos:        Position{Line: 1, Column: 0},
		override:   false,
		overridden: TokenContainer{},
	}
}

func (l *Lexer) Peek() TokenContainer {
	if !l.override {
		l.overridden = l.Next()
	}
	l.override = true
	return l.overridden
}

func (l *Lexer) Next() TokenContainer {
	if l.override {
		l.override = false
		return l.overridden
	}

	r := l.read()

	pos := l.pos

	if r == EOF_LITERAL {
		log.Trace("Lexer.Lex -> r: \":EOF:\"")
		return NewTokenContainer(pos, EOF_T, ":EOF:")
	}

	if r == COMMENT_LITERAL {
		log.Trace("Lexer.Lex -> r: ;")
		tok, str := l.lexComment()
		return NewTokenContainer(pos, tok, str)
	}

	if tok, _ := l.lexFor(r, isWhiteSpace, WHITE_SPACE_T); tok != ILLEGAL_T {
		return l.Next()
	}

	if tok, str := l.lexFor(r, isLineBreak, NEW_LINE_T); tok != ILLEGAL_T {
		return NewTokenContainer(pos, tok, str)
	}

	// isLetter, then isIdentifier
	if unicode.IsLetter(r) {
		log.Trace("Lexer.Lex -> r: [A-Za-z]")
		tok, str := l.lexIdentifier(r)
		return NewTokenContainer(pos, tok, str)
	}

	if tok, str := l.lexFor(r, isNumeric, NUMERIC_T); tok != ILLEGAL_T {
		return NewTokenContainer(pos, tok, str)
	}

	rp := l.read()
	if r == ':' && rp == '=' {
		return NewTokenContainer(pos, ALIAS_T, ":=")
	} else {
		l.unread(rp)
	}

	tok, str := l.lexGrammar(r)
	return NewTokenContainer(pos, tok, str)
}

func (l *Lexer) lexComment() (Token, string) {
	var buf bytes.Buffer

	r := l.read()

	for isNotEol(r) {
		buf.WriteRune(r)
		r = l.read()
	}

	l.unread(r)

	return COMMENT_T, buf.String()
}

type validator func(rune) bool

func (l *Lexer) lexFor(r rune, fn validator, t Token) (Token, string) {
	if !fn(r) {
		return ILLEGAL_T, ""
	}

	var buf bytes.Buffer

	for fn(r) {
		buf.WriteRune(r)
		r = l.read()
	}

	l.unread(r)

	return t, buf.String()
}

func (l *Lexer) lexGrammar(r rune) (Token, string) {
	log.Trace("Lexer.lexGrammar ->")

	var token Token

	switch r {
	case '-':
		token = MINUS_T
	case '*':
		token = ASTERISK_T
	case '\'':
		token = FEET_T
	case '"':
		token = INCHES_T
	case ',':
		token = COMMA_T
	case '(':
		token = LEFT_PAREN_T
	case ')':
		token = RIGHT_PAREN_T
	case '{':
		token = LEFT_BRACE_T
	case '}':
		token = RIGHT_BRACE_T
	case '=':
		token = EQUALS_T
	default:
		token = ILLEGAL_T
	}

	log.Trace("\t-> token: ", token)
	log.Trace("\t-> literal: ", string(r))
	return token, string(r)
}

func (l *Lexer) lexIdentifier(r rune) (Token, string) {
	log.Trace("Lexer.lexIdentifier ->")

	var buf bytes.Buffer

	for isIdentifier(r) {
		log.Trace("\t-> isIdentifier: ", string(r))
		if _, e := buf.WriteRune(r); e != nil {
			panic(e)
		}
		r = l.read()
	}

	l.unread(r)

	log.Trace("\t-> final r: ", r)

	return IDENTIFIER_T, buf.String()
}

func (l *Lexer) read() rune {
	if r, _, err := l.reader.ReadRune(); err != nil {
		return EOF_LITERAL
	} else {
		l.pos.inc(r)
		// log.Trace("Lexer.read -> ", string(r))
		return r
	}
}

func (l *Lexer) unread(r rune) {
	l.pos.dec(r)
	_ = l.reader.UnreadRune()
}

func isNotEol(r rune) bool {
	return !(r == '\r' || r == '\n' || r == EOF_LITERAL)
}

func isWhiteSpace(r rune) bool {
	switch r {
	case '\t', '\v', '\f', ' ', 0x85, 0xA0:
		return true
	}
	return false
}

func isNumeric(r rune) bool {
	return unicode.IsDigit(r) || r == '.'
}

func isLineBreak(r rune) bool {
	return r == '\n' || r == '\r'
}

func isIdentifier(r rune) bool {
	return unicode.IsLetter(r) || r == '-'
}
