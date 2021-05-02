package parser

import (
	"fmt"

	"github.com/bodneyc/knit-and-go/ast"
	. "github.com/bodneyc/knit-and-go/lexer"
	. "github.com/bodneyc/knit-and-go/util"

	log "github.com/sirupsen/logrus"
)

type Parser struct {
	lexer Lexer
	Root  ast.BlockStmt
}

func (o *Parser) WalkForLocals(e *ast.EngineData) {
	o.Root.WalkForLocals(e)
}

func (o *Parser) WalkForLines(e *ast.EngineData) error {
	return o.Root.WalkForLines(e)
}

func NewParserFromBlockStmt(root ast.BlockStmt) *Parser {
	return &Parser{
		lexer: Lexer{},
		Root:  root,
	}
}

func NewParser(lexer Lexer) *Parser {
	return &Parser{
		lexer: lexer,
		Root:  *ast.NewBlockStmt(),
	}
}

// ------------------ Expressions ------------------

func (p *Parser) parseSingleStitch(ident ast.IdentExpr) (ast.Expr, error) {
	tp := p.peekIgnoreWs()
	args := ast.MakeBrackets()
	switch tp.Tok {
	case LEFT_PAREN_T:
		p.nextIgnoreWs() // Consume bracket
		var err error
		args, err = p.parseBrackets()
		if err != nil {
			return nil, fmt.Errorf("%w%s", err, StackLine())
		}

	case IDENTIFIER_T, LEFT_BRACE_T, RIGHT_BRACE_T, COMMENT_T, NEW_LINE_T:
		break

	default:
		return nil, fmt.Errorf("Invalid token %v%s", tp, StackLine())
	}
	return ast.NewStitchExpr(ident, args), nil
}

func (p *Parser) parseRowExpr(firstToken TokenContainer, first bool) (*ast.RowExpr, error) {
	stitches := make([]ast.Expr, 0)

	braced := firstToken.Tok == LEFT_BRACE_T

	for {
		var t TokenContainer
		if first {
			t = firstToken
			first = false
		} else {
			t, _ = p.nextIgnoreWs()
		}

		switch t.Tok {
		case LEFT_BRACE_T:
			row, err := p.parseRowExpr(t, false)
			if err != nil {
				return nil, fmt.Errorf("%w%s", err, StackLine())
			}
			stitches = append(stitches, row)

		case IDENTIFIER_T:
			ident := ast.MakeIdentExpr(t)
			stitch, err := p.parseSingleStitch(ident)
			if err != nil {
				return nil, fmt.Errorf("%w%s", err, StackLine())
			}
			stitches = append(stitches, stitch)
			// log.Debugf("stitches {%#v}\nstitch {%#v}\n", stitches, stitch)

		case RIGHT_BRACE_T:
			if !braced {
				return nil, fmt.Errorf("Found closing brace without opener%s", StackLine())
			}
			var args ast.Brackets
			var err error
			if p.peekIgnoreWs().Tok == LEFT_PAREN_T {
				p.nextIgnoreWs()
				args, err = p.parseBrackets()
				if err != nil {
					return nil, fmt.Errorf("Error parsing bracket group for row%s", StackLine())
				}
			}
			return ast.NewRowExpr(stitches, args), nil

		case NEW_LINE_T, EOF_T:
			if braced {
				return nil, fmt.Errorf("Newline before closing brace in row expression%s", StackLine())
			}
			return ast.NewRowExpr(stitches, ast.MakeBrackets()), nil

		default:
			return nil, fmt.Errorf("Invalid token in row: %v%s", t, StackLine())
		}
	}
}

// First '{' already consumed
func (p *Parser) parseGroupExpr(lBrace TokenContainer) (*ast.GroupExpr, error) {
	lines := make([]ast.Stmt, 0)
	args := ast.MakeBrackets()
	var rBrace TokenContainer
	for {
		tp := p.peekIgnoreWsCr()
		if tp.Tok == RIGHT_BRACE_T {
			var err error
			rBrace, err = p.nextIgnoreWsCr()
			if err != nil {
				return nil, fmt.Errorf("Error finding matching brace for groups: %w%s", err, StackLine())
			}
			if p.peekIgnoreWsCr().Tok == LEFT_PAREN_T {
				p.nextIgnoreWsCr() // Consume '('
				var err error
				args, err = p.parseBrackets()
				if err != nil {
					return nil, fmt.Errorf("Error parsing bracket group for row%s", StackLine())
				}
			}
			break
		}

		t, err := p.nextIgnoreWsCr()
		if err != nil {
			return nil, fmt.Errorf("%w%s", err, StackLine())
		}

		line, err := p.parseLine(ast.CommentGroupExpr{List: make([]ast.CommentExpr, 0)}, t)
		if err != nil {
			return nil, fmt.Errorf("Error parsing rows: %w%s", err, StackLine())
		}

		lines = append(lines, line)
	}
	return ast.NewGroupExpr(lBrace.Pos, rBrace.Pos, lines, args), nil
}

func (p *Parser) parseStitches() (ast.Expr, error) {
	t, err := p.nextIgnoreWs()
	if err != nil {
		return nil, fmt.Errorf("%w%s", err, StackLine())
	}

	switch t.Tok {
	case LEFT_BRACE_T: // Row or group statement
		tp := p.peekIgnoreWs()

		switch tp.Tok {
		case NEW_LINE_T:
			p.nextIgnoreWs() // Consume '\n'
			s, err := p.parseGroupExpr(t)
			if err != nil {
				err = fmt.Errorf("%w%s", err, StackLine())
			}
			return s, err

		case IDENTIFIER_T, LEFT_BRACE_T:
			s, err := p.parseRowExpr(t, false)
			if err != nil {
				err = fmt.Errorf("%w%s", err, StackLine())
			}
			return s, err

		default:
			return nil, fmt.Errorf("Invalid token following left brace: %v%s", t, StackLine())
		}

	case IDENTIFIER_T: // Row or identifier assignment
		s, err := p.parseRowExpr(t, true)
		if err != nil {
			err = fmt.Errorf("%w%s", err, StackLine())
		}
		return s, err

	default:
		return nil, fmt.Errorf("Invalid stitch: %v%s", t, StackLine())
	}
}

// ------------------ Statements ------------------

func (p *Parser) parseRowStmt(desc ast.CommentGroupExpr, firstToken TokenContainer, first bool) (ast.Stmt, error) {
	for {
		if firstToken.Tok == NEW_LINE_T {
			firstToken = p.next()
		}
		if firstToken.Tok == COMMENT_T {
			desc = p.parseCommentExpr(firstToken)
			var err error
			firstToken, err = p.nextIgnoreWs()
			if err != nil {
				return nil, fmt.Errorf("%w%s", err, StackLine())
			}
		} else {
			break
		}
	}

	// log.Debugf("%v%s", firstToken, StackLine())
	row, err := p.parseRowExpr(firstToken, first)
	if err != nil {
		return nil, fmt.Errorf("%w%s", err, StackLine())
	}
	return ast.NewRowStmt(desc, *row), nil
}

func (p *Parser) parseGroupStmt(desc ast.CommentGroupExpr, lBrace TokenContainer) (ast.Stmt, error) {
	group, err := p.parseGroupExpr(lBrace)
	if err != nil {
		return nil, fmt.Errorf("%w%s", err, StackLine())
	}
	return ast.NewGroupStmt(desc, *group), nil
}

func (p *Parser) parseAssignment(desc ast.CommentGroupExpr, ident ast.IdentExpr) (ast.Stmt, error) {
	expr, err := p.parseStitches()
	if err != nil {
		return nil, fmt.Errorf("%w%s", err, StackLine())
	}

	return ast.NewAssignStmt(desc, ident, expr), nil
}

func (p *Parser) parseAlias(desc ast.CommentGroupExpr, lhs ast.IdentExpr) (ast.Stmt, error) {
	t, err := p.nextIgnoreWs()
	if err != nil {
		return nil, fmt.Errorf("%w%s", err, StackLine())
	}
	switch t.Tok {
	case IDENTIFIER_T:
		rhs := ast.MakeIdentExpr(t)
		return ast.NewAliasStmt(desc, lhs, rhs), nil
	default:
		return nil, fmt.Errorf("Invalid alias, %v%s", t, StackLine())
	}
}

// ------------------ Lines ------------------

func (p *Parser) parseIdentExprLine(desc ast.CommentGroupExpr, firstToken TokenContainer) (ast.Stmt, error) {
	ident := ast.MakeIdentExpr(firstToken)

	tp := p.peekIgnoreWs()

	switch tp.Tok {
	case ALIAS_T:
		p.nextIgnoreWs() // Consume ':='
		s, err := p.parseAlias(desc, ident)
		if err != nil {
			err = fmt.Errorf("%w%s", err, StackLine())
		}
		return s, err

	case EQUALS_T:
		p.nextIgnoreWs() // Consume '='
		s, err := p.parseAssignment(desc, ident)
		if err != nil {
			err = fmt.Errorf("%w%s", err, StackLine())
		}
		return s, err

	case LEFT_PAREN_T, LEFT_BRACE_T, IDENTIFIER_T, NEW_LINE_T:
		s, err := p.parseRowStmt(desc, firstToken, true)
		if err != nil {
			err = fmt.Errorf("%w%s", err, StackLine())
		}
		return s, err

	default:
		return nil, fmt.Errorf("Invalid token following identifier %v%s", tp, StackLine())
	}
}

func (p *Parser) parseBraceLine(desc ast.CommentGroupExpr, firstToken TokenContainer, first bool) (ast.Stmt, error) {
	tp := p.peekIgnoreWs()
	switch tp.Tok {
	case NEW_LINE_T:
		s, err := p.parseGroupStmt(desc, firstToken)
		if err != nil {
			err = fmt.Errorf("%w%s", err, StackLine())
		}
		return s, err
	case IDENTIFIER_T, LEFT_BRACE_T:
		s, err := p.parseRowStmt(desc, firstToken, !(firstToken.Tok == LEFT_BRACE_T))
		if err != nil {
			err = fmt.Errorf("%w%s", err, StackLine())
		}
		return s, err
	default:
		return nil, fmt.Errorf("Invalid token following brace at SOL: %v%s", tp, StackLine())
	}
}

func (p *Parser) parseLine(desc ast.CommentGroupExpr, firstToken TokenContainer) (ast.Stmt, error) {

	for {
		if firstToken.Tok == NEW_LINE_T {
			firstToken = p.next()
		}
		if firstToken.Tok == COMMENT_T {
			desc = p.parseCommentExpr(firstToken)
			var err error
			firstToken, err = p.nextIgnoreWs()
			if err != nil {
				return nil, fmt.Errorf("%w%s", err, StackLine())
			}
		} else {
			break
		}
	}

	switch firstToken.Tok {
	case IDENTIFIER_T:
		s, err := p.parseIdentExprLine(desc, firstToken)
		if err != nil {
			err = fmt.Errorf("%w%s", err, StackLine())
		}
		return s, err

	case LEFT_BRACE_T:
		s, err := p.parseBraceLine(desc, firstToken, false)
		if err != nil {
			err = fmt.Errorf("%w%s", err, StackLine())
		}
		return s, err

	default:
		return nil, fmt.Errorf("Illegal start of line : %v%s", firstToken, StackLine())
	}
}

func (p *Parser) Parse() error {
	if t := p.peek(); t.Tok == COMMENT_T {
		p.Root.Desc = p.parseCommentExpr(p.next())
	}

	for {
		var desc ast.CommentGroupExpr = ast.MakeCommentGroupExpr()

		t := p.next()
		if t.Tok == EOF_T {
			return nil
		}

		if t.Tok == NEXT_SOURCE_T {
			log.Info("Moving to next source")
			if t := p.peek(); t.Tok == COMMENT_T {
				p.Root.Desc = p.parseCommentExpr(p.next())
			}
			continue
		}

		if t.Tok == WHITE_SPACE_T || t.Tok == NEW_LINE_T {
			log.Trace("Line beginning with whitespace at ", t.Pos.Str())
			continue
		}

		if stmt, err := p.parseLine(desc, t); err != nil {
			return fmt.Errorf("%w%s", err, StackLine())
		} else {
			p.Root.Block = append(p.Root.Block, stmt)
		}
	}
}
