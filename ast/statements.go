package ast

import (
	"fmt"

	. "github.com/bodneyc/knit-and-go/lexer"
	"github.com/bodneyc/knit-and-go/util"
)

// ------------------ AliasStmt ----------------

type AliasStmt struct {
	Lhs  IdentExpr        `json:"lhs"`
	Rhs  IdentExpr        `json:"rhs"`
	Desc CommentGroupExpr `json:"desc"`
}

func NewAliasStmt(desc CommentGroupExpr, lhs IdentExpr, rhs IdentExpr) *AliasStmt {
	return &AliasStmt{
		Lhs:  lhs,
		Rhs:  rhs,
		Desc: desc,
	}
}

func (s *AliasStmt) stmtNode()     {}
func (s *AliasStmt) Pos() Position { return s.Lhs.Pos() }

func (s *AliasStmt) WalkForLines(e *EngineData) error { return nil }
func (s *AliasStmt) WalkForLocals(e *EngineData) {
	e.aliases[s.Lhs.Name] = s.Rhs
}

// ----------------- AssignStmt ----------------

type AssignStmt struct {
	Lhs  IdentExpr        `json:"lhs"`
	Rhs  Expr             `json:"rhs"`
	Desc CommentGroupExpr `json:"desc"`
}

func (s *AssignStmt) stmtNode()     {}
func (s *AssignStmt) Pos() Position { return s.Lhs.Pos() }

func (s *AssignStmt) WalkForLines(e *EngineData) error { return nil }
func (s *AssignStmt) WalkForLocals(e *EngineData) {
	e.assigns[s.Lhs.Name] = &s.Rhs
	s.Rhs.WalkForLocals(e)
}

func NewAssignStmt(desc CommentGroupExpr, ident IdentExpr, expr Expr) *AssignStmt {
	return &AssignStmt{
		Lhs:  ident,
		Rhs:  expr,
		Desc: desc,
	}
}

// ------------------ RowStmt ------------------

type RowStmt struct {
	Row  RowExpr          `json:"row"`
	Desc CommentGroupExpr `json:"desc"`
}

func NewRowStmt(desc CommentGroupExpr, row RowExpr) *RowStmt {
	return &RowStmt{
		Row:  row,
		Desc: desc,
	}
}

func (s *RowStmt) stmtNode()     {}
func (s *RowStmt) Pos() Position { return s.Row.Stitches[0].Pos() }

func (s *RowStmt) WalkForLines(e *EngineData) error {
	startLc, endLc := START_OF_ROW_LC, END_OF_ROW_LC
	startLc.Desc = s.Desc.TextSlice(e)
	e.Lines = append(e.Lines, startLc)
	lc := MakeLineContainer()
	if err := s.Row.WalkForLines(e, &lc); err != nil {
		return fmt.Errorf("%w%s", err, util.StackLine())
	}
	if !lc.isEmpty() {
		e.Lines = append(e.Lines, lc)
	}
	e.Lines = append(e.Lines, endLc)
	return nil
}

func (s *RowStmt) WalkForLocals(e *EngineData) {
	s.Row.WalkForLocals(e)
}

// ----------------- GroupStmt -----------------

type GroupStmt struct {
	Group GroupExpr        `json:"group"`
	Desc  CommentGroupExpr `json:"desc"`
}

func NewGroupStmt(desc CommentGroupExpr, group GroupExpr) *GroupStmt {
	return &GroupStmt{
		Group: group,
		Desc:  desc,
	}
}

func (s *GroupStmt) stmtNode()     {}
func (s *GroupStmt) Pos() Position { return s.Group.LBrace }

func (s *GroupStmt) WalkForLines(e *EngineData) error {
	lc := START_OF_GROUP_LC
	lc.Desc = s.Desc.TextSlice(e)
	lc.Args = s.Group.Args.TextSlice(e)
	e.Lines = append(e.Lines, lc)
	s.Group.WalkForLines(e, &lc)
	e.Lines = append(e.Lines, END_OF_GROUP_LC)
	return nil
}

func (s *GroupStmt) WalkForLocals(e *EngineData) {
	s.Group.WalkForLocals(e)
}

// ----------------- BlockStmt -----------------

type BlockStmt struct {
	Block []Stmt           `json:"block"`
	Start Position         `json:"start"`
	End   Position         `json:"end"`
	Desc  CommentGroupExpr `json:"desc"`
}

func NewBlockStmt() *BlockStmt {
	return &BlockStmt{
		Block: make([]Stmt, 0),
		Start: Position{Line: 0, Column: 1},
		End:   NO_POS,
		Desc:  CommentGroupExpr{},
	}
}

func (s *BlockStmt) stmtNode()     {}
func (s *BlockStmt) Pos() Position { return s.Start }

func (s *BlockStmt) WalkForLines(e *EngineData) error {
	lc := START_OF_BLOCK_LC
	lc.Desc = s.Desc.TextSlice(e)
	e.Lines = append(e.Lines, lc)
	for _, subblock := range s.Block {
		subblock.WalkForLines(e)
	}
	e.Lines = append(e.Lines, END_OF_BLOCK_LC)
	return nil
}

func (s *BlockStmt) WalkForLocals(e *EngineData) {
	for _, subblock := range s.Block {
		subblock.WalkForLocals(e)
	}
}

// Fuck... to do, to do big time

type ImportStmt struct {
	Doc         *CommentGroupExpr `json:"doc"`
	Name        *IdentExpr        `json:"name"`
	CommentExpr *CommentGroupExpr `json:"comment"`
}

func (s *ImportStmt) stmtNode()     {}
func (s *ImportStmt) Pos() Position { return s.Name.Pos() }

func (s *ImportStmt) WalkForLocals(e *EngineData) {
	// TODO:
}
