package parser

import (
	"fmt"

	. "github.com/bodneyc/knit-and-go/lexer"
	. "github.com/bodneyc/knit-and-go/util"

	log "github.com/sirupsen/logrus"
)

func (p *Parser) next() TokenContainer {
	t := p.lexer.Next()
	log.WithFields(t.Fields()).Trace("[Parser.next]")
	return t
}

func (p *Parser) nextIgnoreWsCr() (TokenContainer, error) {
	for {
		t := p.lexer.Next()
		log.WithFields(t.Fields()).Trace("[Parser.nextIgnoreWsCr]")
		if t.Tok == EOF_T {
			return t, fmt.Errorf(":EOF: before next token%s", StackLine())
		}
		if t.Tok != WHITE_SPACE_T && t.Tok != NEW_LINE_T {
			return t, nil
		}
	}
}

func (p *Parser) nextIgnoreWs() (TokenContainer, error) {
	for {
		t := p.lexer.Next()
		log.WithFields(t.Fields()).Trace("[Parser.nextIgnoreWs]")
		if t.Tok == EOF_T {
			return t, fmt.Errorf(":EOF: before next token%s", StackLine())
		}
		if t.Tok != WHITE_SPACE_T {
			return t, nil
		}
	}
}

func (p *Parser) peek() TokenContainer {
	t := p.lexer.Peek()
	log.WithFields(t.Fields()).Trace("[Parser.peek]")
	return t
}

func (p *Parser) peekIgnoreWsCr() TokenContainer {
	for {
		tp := p.lexer.Peek()
		if tp.Tok != WHITE_SPACE_T && tp.Tok != NEW_LINE_T {
			log.WithFields(tp.Fields()).Trace("[Parser.peekIgnoreWsCr]")
			return tp
		}
		p.lexer.Next()
	}
}

func (p *Parser) peekIgnoreWs() TokenContainer {
	for {
		if tp := p.lexer.Peek(); tp.Tok != WHITE_SPACE_T {
			log.WithFields(tp.Fields()).Trace("[Parser.peekIgnoreWs]")
			return tp
		}
		p.lexer.Next()
	}
}
