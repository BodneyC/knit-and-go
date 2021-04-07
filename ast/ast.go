package ast

import (
	. "github.com/bodneyc/knit-and-go/lexer"

	"strings"
	"unicode"
)

type Node interface {
	Pos() Position
}

type Expr interface {
	Node
	exprNode()
}

type Stmt interface {
	Node
	stmtNode()
}

type Decl interface {
	Node
	declNode()
}

// ---------------------------------------------
// Expressions
// ---------------------------------------------

// ---------------- CommentGroup ---------------

type Comment struct {
	Semicolon Position `json:"semicolon"`
	Text      string   `json:"text"`
}

func (c *Comment) Pos() Position { return c.Semicolon }

type CommentGroup struct {
	List []Comment `json:"list"`
}

func MakeCommentGroup() CommentGroup {
	return CommentGroup{
		List: make([]Comment, 0),
	}
}

func (g *CommentGroup) Pos() Position { return g.List[0].Pos() }
func (g *CommentGroup) exprNode()     {}

func stripLeadingWhitespace(s string) string {
	i := len(s)
	for i > 0 && unicode.IsSpace(rune(s[i-1])) {
		i--
	}
	return s[0:i]
}

func (g *CommentGroup) Text() string {
	if g == nil {
		return ""
	}
	comments := make([]string, len(g.List))
	for i, c := range g.List {
		comments[i] = c.Text
	}
	lines := make([]string, 0, 10)
	for _, c := range comments {
		switch c[1] {
		case ';':
			c = c[1:]
			if len(c) == 0 {
				break
			}
			if c[0] == ' ' {
				c = c[1:]
				break
			}
		}
		lines = append(lines, stripLeadingWhitespace(c))
	}
	n := 0
	for _, line := range lines {
		if line != "" || n > 0 && lines[n-1] != "" {
			lines[n] = line
			n++
		}
	}
	lines = lines[0:n]
	if n > 0 && lines[n-1] != "" {
		lines = append(lines, "")
	}
	return strings.Join(lines, "\n")
}

// ----------------- Identifier ----------------

type Ident struct {
	At   Position `json:"at"`
	Name string   `json:"name"`
}

func (x *Ident) Pos() Position { return x.At }
func (g *Ident) exprNode()     {}

func MakeIdent(t TokenContainer) Ident {
	return Ident{
		At:   t.Pos,
		Name: t.Str,
	}
}

func NewIdent(t TokenContainer) *Ident {
	return &Ident{
		At:   t.Pos,
		Name: t.Str,
	}
}

// ----------------- Size ----------------

type MeasurementUnit int

const (
	NOSIZE MeasurementUnit = iota
	MM
	CM
	INCHES
	FEET
	ASTERISK
)

type Size struct {
	At     Position        `json:"at"`
	Ni     int64           `json:"ni,string"`
	Nf     float64         `json:"nf,string"`
	NId    Ident           `json:"nid"`
	Before bool            `json:"before"`
	Unit   MeasurementUnit `json:"unit,string"`
}

func (x *Size) Pos() Position { return x.At }
func (g *Size) exprNode()     {}

func NewSizeAsterisk(asterisk TokenContainer) *Size {
	return &Size{
		At:     asterisk.Pos,
		Ni:     -1,
		Nf:     -1,
		Before: false,
		NId:    Ident{},
		Unit:   ASTERISK,
	}
}

func NewSize(ni int64, nf float64, nid TokenContainer, unit MeasurementUnit) *Size {
	return &Size{
		At:     nid.Pos,
		Ni:     ni,
		Nf:     nf,
		NId:    MakeIdent(nid),
		Before: false,
		Unit:   unit,
	}
}

// ----- Bracket group, implements nothing -----

type BracketGroup struct {
	Args []Expr `json:"args"`
}

func MakeBracketGroup() BracketGroup {
	return BracketGroup{
		Args: make([]Expr, 0),
	}
}

// --------------- Single stitch ---------------

type StitchExpr struct {
	At     Position     `json:"at"`
	Stitch Ident        `json:"stitch"`
	Args   BracketGroup `json:"args"`
}

func (x *StitchExpr) Pos() Position { return x.At }
func (g *StitchExpr) exprNode()     {}

func NewStitchExpr(ident Ident, args BracketGroup) *StitchExpr {
	return &StitchExpr{
		At:     ident.Pos(),
		Stitch: ident,
		Args:   args,
	}
}

// -------------------- Row --------------------

type RowExpr struct {
	Stitches []Expr       `json:"stitches"`
	Args     BracketGroup `json:"args"`
}

func (x *RowExpr) Pos() Position { return x.Stitches[0].Pos() }
func (g *RowExpr) exprNode()     {}

func NewRowExpr(exprs []Expr, args BracketGroup) *RowExpr {
	return &RowExpr{
		Stitches: exprs,
		Args:     args,
	}
}

// ------------------- Group -------------------

type GroupExpr struct {
	LBrace Position     `json:"lbrace"`
	RBrace Position     `json:"rbrace"`
	Args   BracketGroup `json:"args"`
	Lines  []Stmt       `json:"lines"`
}

func (s *GroupExpr) Pos() Position { return s.LBrace }
func (g *GroupExpr) exprNode()     {}

func NewGroupExpr(l Position, r Position, lines []Stmt, args BracketGroup) *GroupExpr {
	return &GroupExpr{
		LBrace: l,
		RBrace: r,
		Lines:  lines,
		Args:   args,
	}
}

// ---------------------------------------------
// Statements
// ---------------------------------------------

// ------------------ AliasStmt ----------------

type AliasStmt struct {
	Lhs  Ident        `json:"lhs"`
	Rhs  Ident        `json:"rhs"`
	Desc CommentGroup `json:"desc"`
}

func (s *AliasStmt) Pos() Position { return s.Lhs.Pos() }
func (*AliasStmt) stmtNode()       {}

func NewAliasStmt(desc CommentGroup, lhs Ident, rhs Ident) *AliasStmt {
	return &AliasStmt{
		Lhs:  lhs,
		Rhs:  rhs,
		Desc: desc,
	}
}

// ----------------- AssignStmt ----------------

type AssignStmt struct {
	Lhs  Ident        `json:"lhs"`
	Rhs  Expr         `json:"rhs"`
	Desc CommentGroup `json:"desc"`
}

func (s *AssignStmt) Pos() Position { return s.Lhs.Pos() }
func (*AssignStmt) stmtNode()       {}

func NewAssignStmt(desc CommentGroup, ident Ident, expr Expr) *AssignStmt {
	return &AssignStmt{
		Lhs:  ident,
		Rhs:  expr,
		Desc: desc,
	}
}

// ------------------ RowStmt ------------------

type RowStmt struct {
	Row  RowExpr      `json:"row"`
	Desc CommentGroup `json:"desc"`
}

func (*RowStmt) stmtNode() {}

func (s *RowStmt) Pos() Position { return s.Row.Stitches[0].Pos() }

func NewRowStmt(desc CommentGroup, row RowExpr) *RowStmt {
	return &RowStmt{
		Row:  row,
		Desc: desc,
	}
}

// ----------------- GroupStmt -----------------

type GroupStmt struct {
	Group GroupExpr    `json:"group"`
	Desc  CommentGroup `json:"desc"`
}

func (*GroupStmt) stmtNode()       {}
func (s *GroupStmt) Pos() Position { return s.Group.LBrace }

func NewGroupStmt(desc CommentGroup, group GroupExpr) *GroupStmt {
	return &GroupStmt{
		Group: group,
		Desc:  desc,
	}
}

// ----------------- BlockStmt -----------------

type BlockStmt struct {
	Block  []Stmt       `json:"block"`
	Start  Position     `json:"start"`
	End    Position     `json:"end"`
	Length int64        `json:"length,string"`
	Desc   CommentGroup `json:"desc"`
}

func (*BlockStmt) stmtNode() {}

func (s *BlockStmt) Pos() Position { return s.Start }

func NewBlockStmt() *BlockStmt {
	return &BlockStmt{
		Block:  make([]Stmt, 0),
		Start:  Position{Line: 0, Column: 1},
		End:    NO_POS,
		Length: -1,
		Desc:   CommentGroup{},
	}
}

// Fuck... to do, to do big time

type (
	ImportSpec struct {
		Doc  *CommentGroup `json:"doc"`
		Name *Ident        `json:"name"`
		// Path    *ValueE `json:"path"`
		Comment *CommentGroup `json:"comment"`
	}
)

func (s *ImportSpec) Pos() Position {
	return s.Name.Pos()
}

func (*ImportSpec) specNode() {}
