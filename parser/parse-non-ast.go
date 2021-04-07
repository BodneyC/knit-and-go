package parser

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bodneyc/knit-and-go/ast"
	. "github.com/bodneyc/knit-and-go/lexer"
	"github.com/bodneyc/knit-and-go/util"
)

func (p *Parser) parseSizeMinus() (*ast.Size, error) {
	t, e := p.nextIgnoreWs()
	if e != nil {
		return nil, fmt.Errorf("%s : %w", util.Fname(), e)
	}
	if t.Tok != IDENTIFIER_T && t.Tok != NUMERIC_T {
		return nil, fmt.Errorf("%sValue following minus sign not valid", util.Fname())
	}
	size, e := p.parseSize(t)
	if e != nil {
		return nil, fmt.Errorf("%s : %w", util.Fname(), e)
	}
	size.Before = true
	return size, nil
}

func (p *Parser) parseSize(t TokenContainer) (*ast.Size, error) {
	var ni int64 = -1
	var nf float64 = -1.0
	if val, e := strconv.ParseInt(t.Str, 10, 16); e == nil {
		ni = int64(val)
	}
	if val, e := strconv.ParseFloat(t.Str, 64); e == nil {
		nf = float64(val)
	}

	var unit ast.MeasurementUnit
	tp := p.peekIgnoreWs()

	switch tp.Tok {
	case COMMA_T, RIGHT_PAREN_T:
		unit = ast.NOSIZE
	case FEET_T:
		p.nextIgnoreWs()
		unit = ast.FEET
	case INCHES_T:
		p.nextIgnoreWs()
		unit = ast.INCHES
	case IDENTIFIER_T: // 'mm' or 'cm'
		if strings.EqualFold(tp.Str, "mm") {
			unit = ast.MM
		} else if strings.EqualFold(tp.Str, "cm") {
			unit = ast.CM
		} else {
			return nil, fmt.Errorf("%sNo valid unit identifier found:\n  %#v", util.Fname(), tp)
		}
		p.nextIgnoreWs()
	default:
		return nil, fmt.Errorf("%sNo valid unit token found:\n  %#v", util.Fname(), tp)
	}
	return ast.NewSize(ni, nf, t, unit), nil
}

// LEFT_PAREN_T already consumed
func (p *Parser) parseBracketGroup() (ast.BracketGroup, error) {
	args := make([]ast.Expr, 0)
	for {
		t, e := p.nextIgnoreWs()
		if e != nil {
			return ast.BracketGroup{}, e
		}
		switch t.Tok {
		case IDENTIFIER_T:
			args = append(args, ast.NewIdent(t))

		case MINUS_T, NUMERIC_T:
			var size *ast.Size
			var e error
			if t.Tok == MINUS_T {
				size, e = p.parseSizeMinus()
			} else {
				size, e = p.parseSize(t)
			}
			if e != nil {
				return ast.BracketGroup{}, fmt.Errorf("%s : %w", util.Fname(), e)
			}
			args = append(args, size)

		case ASTERISK_T:
			args = append(args, ast.NewSizeAsterisk(t))

		case COMMA_T:
			continue

		case RIGHT_PAREN_T:
			return ast.BracketGroup{args}, nil

		default:
			return ast.BracketGroup{}, fmt.Errorf("%sInvalid token in bracket group: %v", util.Fname(), t)
		}
	}
}

func (p *Parser) parseComment(t TokenContainer) ast.CommentGroup {
	commentGroup := ast.CommentGroup{
		List: make([]ast.Comment, 0),
	}
	for {
		comment := ast.Comment{
			Semicolon: t.Pos,
			Text:      t.Str,
		}
		commentGroup.List = append(commentGroup.List, comment)
		if peeked := p.peek().Tok; peeked != COMMENT_T && peeked != NEW_LINE_T {
			break
		}
		t = p.next()
	}
	return commentGroup
}
