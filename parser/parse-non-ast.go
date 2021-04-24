package parser

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bodneyc/knit-and-go/ast"
	. "github.com/bodneyc/knit-and-go/lexer"
	. "github.com/bodneyc/knit-and-go/util"
)

func (p *Parser) parseSizeExprMinus() (*ast.SizeExpr, error) {
	t, err := p.nextIgnoreWs()
	if err != nil {
		return nil, fmt.Errorf("%w%s", err, StackLine())
	}
	if t.Tok != IDENTIFIER_T && t.Tok != NUMERIC_T {
		return nil, fmt.Errorf("%sValue following minus sign not valid", StackLine())
	}
	size, err := p.parseSizeExpr(t)
	if err != nil {
		return nil, fmt.Errorf("%w%s", err, StackLine())
	}
	size.Before = true
	return size, nil
}

func (p *Parser) parseSizeExpr(t TokenContainer) (*ast.SizeExpr, error) {
	var ni int64 = -1
	var nf float64 = -1.0
	if val, err := strconv.ParseInt(t.Str, 10, 16); err == nil {
		ni = int64(val)
	}
	if val, err := strconv.ParseFloat(t.Str, 64); err == nil {
		nf = float64(val)
	}

	var unit ast.MeasurementUnit
	tp := p.peekIgnoreWs()

	switch tp.Tok {
	case COMMA_T, RIGHT_PAREN_T:
		unit = ast.NOUNIT
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
			return nil, fmt.Errorf("%sNo valid unit identifier found:\n  %#v", StackLine(), tp)
		}
		p.nextIgnoreWs()
	default:
		return nil, fmt.Errorf("%sNo valid unit token found:\n  %#v", StackLine(), tp)
	}
	return ast.NewSizeExpr(ni, nf, t, unit), nil
}

// LEFT_PAREN_T already consumed
func (p *Parser) parseBrackets() (ast.Brackets, error) {
	args := make([]ast.Expr, 0)
	for {
		t, err := p.nextIgnoreWs()
		if err != nil {
			return ast.Brackets{}, err
		}
		switch t.Tok {
		case IDENTIFIER_T:
			args = append(args, ast.NewIdentExpr(t))

		case MINUS_T, NUMERIC_T:
			var size *ast.SizeExpr
			var err error
			if t.Tok == MINUS_T {
				size, err = p.parseSizeExprMinus()
			} else {
				size, err = p.parseSizeExpr(t)
			}
			if err != nil {
				return ast.Brackets{}, fmt.Errorf("%s : %w", StackLine(), err)
			}
			args = append(args, size)

		case ASTERISK_T:
			args = append(args, ast.NewSizeExprAsterisk(t))

		case COMMA_T:
			continue

		case RIGHT_PAREN_T:
			return ast.Brackets{Args: args}, nil

		default:
			return ast.Brackets{}, fmt.Errorf("%sInvalid token in bracket group: %v", StackLine(), t)
		}
	}
}

func (p *Parser) parseCommentExpr(t TokenContainer) ast.CommentGroupExpr {
	commentGroup := ast.CommentGroupExpr{
		List: make([]ast.CommentExpr, 0),
	}
	for {
		comment := ast.CommentExpr{
			Semicolon: t.Pos,
			Str:       t.Str,
		}
		commentGroup.List = append(commentGroup.List, comment)
		peeked := p.peek()
		if peeked.Tok != COMMENT_T && peeked.Tok != NEW_LINE_T {
			break
		}
		t = p.next()
		if t.Tok == NEW_LINE_T {
			peeked = p.peek()
			if peeked.Tok == COMMENT_T {
				if t.Str == "\n\n" {
					break
				}
				t = p.next()
			}
		}
	}
	return commentGroup
}
