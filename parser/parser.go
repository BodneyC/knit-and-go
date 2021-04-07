package parser

import (
	"fmt"

	"github.com/bodneyc/knit-and-go/ast"
	. "github.com/bodneyc/knit-and-go/lexer"
	"github.com/bodneyc/knit-and-go/util"

	log "github.com/sirupsen/logrus"
)

type Parser struct {
	lexer  Lexer
	locals map[string]ast.Node
	Root   ast.BlockStmt
}

func NewParser(lexer Lexer) *Parser {
	return &Parser{
		lexer:  lexer,
		locals: make(map[string]ast.Node),
		Root:   *ast.NewBlockStmt(),
	}
}

// ------------------ Expressions ------------------

func (p *Parser) parseSingleStitch(ident ast.Ident) (ast.Expr, error) {
	tp := p.peekIgnoreWs()
	args := ast.MakeBracketGroup()
	switch tp.Tok {
	case LEFT_PAREN_T:
		p.nextIgnoreWs() // Consume bracket
		var e error
		args, e = p.parseBracketGroup()
		if e != nil {
			return nil, fmt.Errorf("%s : %w", util.Fname(), e)
		}

	case IDENTIFIER_T, LEFT_BRACE_T, RIGHT_BRACE_T, COMMENT_T, NEW_LINE_T:
		break

	default:
		return nil, fmt.Errorf("%sInvalid token %v", util.Fname(), tp)
	}
	return ast.NewStitchExpr(ident, args), nil
}

// First '{' already consumed
func (p *Parser) parseRowExpr(firstToken TokenContainer, first bool) (*ast.RowExpr, error) {
	stitches := make([]ast.Expr, 0)

	wasFirst := first
	braced := false
	if firstToken.Tok == LEFT_BRACE_T && first {
		braced = true
	}

	for {
		var t TokenContainer
		if first {
			t = firstToken
			first = false
		} else {
			var e error
			t, e = p.nextIgnoreWs()
			if e != nil {
				return nil, fmt.Errorf("%s : %w", util.Fname(), e)
			}
		}
		switch t.Tok {
		case LEFT_BRACE_T:
			row, e := p.parseRowExpr(TokenContainer{}, false)
			if e != nil {
				return nil, fmt.Errorf("%s : %w", util.Fname(), e)
			}
			if wasFirst {
				return row, nil
			} else {
				stitches = append(stitches, row)
			}

		case IDENTIFIER_T:
			ident := ast.MakeIdent(t)
			stitch, e := p.parseSingleStitch(ident)
			if e != nil {
				return nil, fmt.Errorf("%s : %w", util.Fname(), e)
			}
			stitches = append(stitches, stitch)

		case RIGHT_BRACE_T:
			var args ast.BracketGroup
			var e error
			if p.peekIgnoreWs().Tok == LEFT_PAREN_T {
				p.nextIgnoreWs()
				args, e = p.parseBracketGroup()
				if e != nil {
					return nil, fmt.Errorf("%sError parsing bracket group for row", util.Fname())
				}
			}
			return ast.NewRowExpr(stitches, args), nil

		case NEW_LINE_T:
			if braced {
				return nil, fmt.Errorf("%sNewline before closing brace in row expression", util.Fname())
			}
			return ast.NewRowExpr(stitches, ast.MakeBracketGroup()), nil

		default:
			return nil, fmt.Errorf("%sInvalid token in row: %v", util.Fname(), t)
		}
	}
}

// First '{' already consumed
func (p *Parser) parseGroupExpr(lBrace TokenContainer) (*ast.GroupExpr, error) {
	lines := make([]ast.Stmt, 0)
	args := ast.MakeBracketGroup()
	var rBrace TokenContainer
	for {
		tp := p.peekIgnoreWsCr()
		if tp.Tok == RIGHT_BRACE_T {
			var e error
			rBrace, e = p.nextIgnoreWsCr()
			if e != nil {
				return nil, fmt.Errorf("%sError finding matching brace for groups: %w", util.Fname(), e)
			}
			if p.peekIgnoreWsCr().Tok == LEFT_PAREN_T {
				p.nextIgnoreWsCr() // Consume '('
				var e error
				args, e = p.parseBracketGroup()
				if e != nil {
					return nil, fmt.Errorf("%sError parsing bracket group for row", util.Fname())
				}
			}
			break
		}

		t, e := p.nextIgnoreWsCr()
		if e != nil {
			return nil, fmt.Errorf("%s : %w", util.Fname(), e)
		}

		line, e := p.parseLine(ast.CommentGroup{List: make([]ast.Comment, 0)}, t)
		if e != nil {
			return nil, fmt.Errorf("%sError parsing rows: %w", util.Fname(), e)
		}

		lines = append(lines, line)
	}
	return ast.NewGroupExpr(lBrace.Pos, rBrace.Pos, lines, args), nil
}

func (p *Parser) parseStitches() (ast.Expr, error) {
	t, e := p.nextIgnoreWs()
	if e != nil {
		return nil, fmt.Errorf("%s : %w", util.Fname(), e)
	}

	switch t.Tok {
	case LEFT_BRACE_T: // Row or group statement
		tp := p.peekIgnoreWs()

		switch tp.Tok {
		case NEW_LINE_T:
			p.nextIgnoreWs() // Consume '\n'
			s, e := p.parseGroupExpr(t)
			if e != nil {
				e = fmt.Errorf("%s : %w", util.Fname(), e)
			}
			return s, e

		case IDENTIFIER_T, LEFT_BRACE_T:
			s, e := p.parseRowExpr(TokenContainer{}, false)
			if e != nil {
				e = fmt.Errorf("%s : %w", util.Fname(), e)
			}
			return s, e

		default:
			return nil, fmt.Errorf("%sInvalid token following left brace: %v", util.Fname(), t)
		}

	case IDENTIFIER_T: // Row or identifier assignment
		s, e := p.parseRowExpr(t, true)
		if e != nil {
			e = fmt.Errorf("%s : %w", util.Fname(), e)
		}
		return s, e

	default:
		return nil, fmt.Errorf("%sInvalid stitch: %v", util.Fname(), t)
	}
}

// ------------------ Statements ------------------

func (p *Parser) parseRowStmt(desc ast.CommentGroup, firstToken TokenContainer, first bool) (ast.Stmt, error) {
	row, e := p.parseRowExpr(firstToken, first)
	if e != nil {
		return nil, fmt.Errorf("%s : %w", util.Fname(), e)
	}
	return ast.NewRowStmt(desc, *row), nil
}

func (p *Parser) parseGroupStmt(desc ast.CommentGroup, lBrace TokenContainer) (ast.Stmt, error) {
	group, e := p.parseGroupExpr(lBrace)
	if e != nil {
		return nil, fmt.Errorf("%s : %w", util.Fname(), e)
	}
	return ast.NewGroupStmt(desc, *group), nil
}

func (p *Parser) parseAssignment(desc ast.CommentGroup, ident ast.Ident) (ast.Stmt, error) {
	expr, e := p.parseStitches()
	if e != nil {
		return nil, fmt.Errorf("%s : %w", util.Fname(), e)
	}

	return ast.NewAssignStmt(desc, ident, expr), nil
}

func (p *Parser) parseAlias(desc ast.CommentGroup, lhs ast.Ident) (ast.Stmt, error) {
	t, e := p.nextIgnoreWs()
	if e != nil {
		return nil, e
	}
	switch t.Tok {
	case IDENTIFIER_T:
		rhs := ast.MakeIdent(t)
		return ast.NewAliasStmt(desc, lhs, rhs), nil
	default:
		return nil, fmt.Errorf("%sInvalid alias, %v", util.Fname(), t)
	}
}

// ------------------ Lines ------------------

func (p *Parser) parseIdentifierLine(desc ast.CommentGroup, firstToken TokenContainer) (ast.Stmt, error) {
	ident := ast.MakeIdent(firstToken)

	tp := p.peekIgnoreWs()

	switch tp.Tok {
	case ALIAS_T:
		p.nextIgnoreWs() // Consume ':='
		s, e := p.parseAlias(desc, ident)
		if e != nil {
			e = fmt.Errorf("%s : %w", util.Fname(), e)
		}
		return s, e

	case EQUALS_T:
		p.nextIgnoreWs() // Consume '='
		s, e := p.parseAssignment(desc, ident)
		if e != nil {
			e = fmt.Errorf("%s : %w", util.Fname(), e)
		}
		return s, e

	case LEFT_PAREN_T, LEFT_BRACE_T, IDENTIFIER_T, NEW_LINE_T:
		s, e := p.parseRowStmt(desc, firstToken, true)
		if e != nil {
			e = fmt.Errorf("%s : %w", util.Fname(), e)
		}
		return s, e

	default:
		return nil, fmt.Errorf("%sInvalid token following identifier %v", util.Fname(), tp)
	}
}

func (p *Parser) parseBraceLine(desc ast.CommentGroup, firstToken TokenContainer) (ast.Stmt, error) {
	tp := p.peekIgnoreWs()
	switch tp.Tok {
	case NEW_LINE_T:
		s, e := p.parseGroupStmt(desc, firstToken)
		if e != nil {
			e = fmt.Errorf("%s : %w", util.Fname(), e)
		}
		return s, e
	case IDENTIFIER_T, LEFT_BRACE_T:
		s, e := p.parseRowStmt(desc, firstToken, true)
		if e != nil {
			e = fmt.Errorf("%s : %w", util.Fname(), e)
		}
		return s, e
	default:
		return nil, fmt.Errorf("%sInvalid token following brace at SOL: %v", util.Fname(), tp)
	}
}

func (p *Parser) parseLine(desc ast.CommentGroup, firstToken TokenContainer) (ast.Stmt, error) {

	switch firstToken.Tok {
	case COMMENT_T: // Keep first case
		// parseComment will consume until no comments left
		desc = p.parseComment(firstToken)
		fallthrough

	case IDENTIFIER_T:
		s, e := p.parseIdentifierLine(desc, firstToken)
		if e != nil {
			e = fmt.Errorf("%s : %w", util.Fname(), e)
		}
		return s, e

	case LEFT_BRACE_T:
		s, e := p.parseBraceLine(desc, firstToken)
		if e != nil {
			e = fmt.Errorf("%s : %w", util.Fname(), e)
		}
		return s, e

	default:
		return nil, fmt.Errorf("%sIllegal start of line : %v", util.Fname(), firstToken)
	}
}

func (p *Parser) Parse() error {
	if t := p.peek(); t.Tok == COMMENT_T {
		p.Root.Desc = p.parseComment(p.next())
	}

	for {
		var desc ast.CommentGroup = ast.MakeCommentGroup()

		t := p.next()
		if t.Tok == EOF_T {
			return nil
		}

		if t.Tok == WHITE_SPACE_T|NEW_LINE_T {
			log.Debug("Line beginning with whitespace at ", t.Pos.Str())
			continue
		}

		if stmt, e := p.parseLine(desc, t); e != nil {
			return e
		} else {
			p.Root.Block = append(p.Root.Block, stmt)
		}
	}
}
