package lexer

import (
	"bufio"
	"bytes"
	"fmt"
	"os"

	// "strings"
	"unicode"

	log "github.com/sirupsen/logrus"
)

type Lexer struct {
	infiles    []string
	inputIdx   int
	file       *os.File
	reader     *bufio.Reader
	pos        Position
	override   bool
	overridden TokenContainer
}

func NewLexer(infiles []string) (*Lexer, error) {
	l := &Lexer{
		infiles:    infiles,
		inputIdx:   0,
		file:       nil,
		reader:     nil,
		pos:        Position{Line: 1, Column: 0},
		override:   false,
		overridden: TokenContainer{},
	}
	suc, err := l.readNextInput()
	if err != nil {
		return nil, err
	}
	if !suc {
		return nil, fmt.Errorf("No further input files")
	}
	return l, nil
}

func (l *Lexer) readNextInput() (bool, error) {
	if len(l.infiles) == l.inputIdx {
		return false, nil
	}

	if l.file != nil {
		if err := l.file.Close(); err != nil {
			return false, err
		}
	}

	log.WithField("infile", l.infiles[l.inputIdx]).Info("Attempting to open input")

	file, err := os.Open(l.infiles[l.inputIdx])
	if err != nil {
		return false, err
	}

	l.file = file
	l.reader = bufio.NewReader(file)

	l.inputIdx++

	return true, nil
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
		log.WithField("literal", "\":EOF:\"").Trace("[Lexer.Next]")
		b, err := l.readNextInput()
		if err != nil {
			log.Warn("Failed to open next input file", err)
		}
		if b {
			return NewTokenContainer(pos, NEXT_SOURCE_T, ":NEXT_SOURCE:")
		} else {
			return NewTokenContainer(pos, EOF_T, ":EOF:")
		}
	}

	if r == COMMENT_LITERAL {
		log.WithField("literal", ";").Trace("[Lexer.Next]")
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
		log.WithField("literal", "[A-Za-z]").Trace("[Lexer.Next]")
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

	log.WithFields(log.Fields{"literal": string(r), "token": token}).Trace("[Lexer.lexGrammar]")
	return token, string(r)
}

func (l *Lexer) lexIdentifier(r rune) (Token, string) {
	var buf bytes.Buffer

	for isIdentifier(r) {
		if _, err := buf.WriteRune(r); err != nil {
			panic(err)
		}
		r = l.read()
	}

	l.unread(r)

	log.WithField("literal", r).Trace("[Lexer.lexIdentifier]")

	return IDENTIFIER_T, buf.String()
}

func (l *Lexer) read() rune {
	if r, _, err := l.reader.ReadRune(); err != nil {
		return EOF_LITERAL
	} else {
		l.pos.inc(r)
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
