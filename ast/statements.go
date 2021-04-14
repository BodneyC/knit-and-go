package ast

import (
	"fmt"

	. "github.com/bodneyc/knit-and-go/lexer"
	"github.com/bodneyc/knit-and-go/util"
	log "github.com/sirupsen/logrus"
)

// ------------------ AliasStmt ----------------

type AliasStmt struct {
	Lhs  IdentExpr        `json:"lhs"`
	Rhs  IdentExpr        `json:"rhs"`
	Desc CommentGroupExpr `json:"desc"`
}

func (o *AliasStmt) stmtNode()     {}
func (s *AliasStmt) Pos() Position { return s.Lhs.Pos() }

func (o *AliasStmt) WalkForLines(e *Engine) error { return nil }
func (o *AliasStmt) WalkForLocals(e *Engine) {
	e.aliases[o.Lhs.Name] = o.Rhs
}

func NewAliasStmt(desc CommentGroupExpr, lhs IdentExpr, rhs IdentExpr) *AliasStmt {
	return &AliasStmt{
		Lhs:  lhs,
		Rhs:  rhs,
		Desc: desc,
	}
}

// ----------------- AssignStmt ----------------

type AssignStmt struct {
	Lhs  IdentExpr        `json:"lhs"`
	Rhs  Expr             `json:"rhs"`
	Desc CommentGroupExpr `json:"desc"`
}

func (o *AssignStmt) stmtNode()     {}
func (s *AssignStmt) Pos() Position { return s.Lhs.Pos() }

func (o *AssignStmt) WalkForLines(e *Engine) error { return nil }
func (o *AssignStmt) WalkForLocals(e *Engine) {
	e.assigns[o.Lhs.Name] = &o.Rhs
	o.Rhs.WalkForLocals(e)
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

func (o *RowStmt) stmtNode()     {}
func (s *RowStmt) Pos() Position { return s.Row.Stitches[0].Pos() }

func (o *RowStmt) WalkForLines(e *Engine) error {
	log.Info("In RowStmt#WalkForLines")
	lc := MakeLineContainer()
	if err := o.Row.WalkForLines(e, &lc); err != nil {
		return fmt.Errorf("%w%s", err, util.StackLine())
	}
	if !lc.IsEmpty() {
		startLc, endLc := START_OF_ROW_LC, END_OF_ROW_LC
		startLc.Desc = o.Desc.TextSlice(e)
		e.Lines = append(e.Lines, startLc)
		e.Lines = append(e.Lines, lc)
		e.Lines = append(e.Lines, endLc)
	}
	return nil
}

func (o *RowStmt) WalkForLocals(e *Engine) {
	o.Row.WalkForLocals(e)
}

func NewRowStmt(desc CommentGroupExpr, row RowExpr) *RowStmt {
	return &RowStmt{
		Row:  row,
		Desc: desc,
	}
}

// ----------------- GroupStmt -----------------

type GroupStmt struct {
	Group GroupExpr        `json:"group"`
	Desc  CommentGroupExpr `json:"desc"`
}

func (o *GroupStmt) stmtNode()     {}
func (s *GroupStmt) Pos() Position { return s.Group.LBrace }

func (o *GroupStmt) WalkForLines(e *Engine) error {
	lc := START_OF_GROUP_LC
	lc.Desc = o.Desc.TextSlice(e)
	lc.Args = o.Group.Args.TextSlice(e)
	e.Lines = append(e.Lines, lc)
	o.Group.WalkForLines(e, &lc)
	e.Lines = append(e.Lines, END_OF_GROUP_LC)
	return nil
}

func (o *GroupStmt) WalkForLocals(e *Engine) {
	o.Group.WalkForLocals(e)
}

func NewGroupStmt(desc CommentGroupExpr, group GroupExpr) *GroupStmt {
	return &GroupStmt{
		Group: group,
		Desc:  desc,
	}
}

// ----------------- BlockStmt -----------------

type BlockStmt struct {
	Block []Stmt           `json:"block"`
	Start Position         `json:"start"`
	End   Position         `json:"end"`
	Desc  CommentGroupExpr `json:"desc"`
}

func (o *BlockStmt) stmtNode()     {}
func (o *BlockStmt) Pos() Position { return o.Start }

func (o *BlockStmt) WalkForLines(e *Engine) error {
	lc := START_OF_BLOCK_LC
	lc.Desc = o.Desc.TextSlice(e)
	e.Lines = append(e.Lines, lc)
	log.Tracef("%s", util.StackLine())
	for _, subblock := range o.Block {
		log.Tracef("%#v\n", subblock)
		subblock.WalkForLines(e)
	}
	e.Lines = append(e.Lines, END_OF_BLOCK_LC)
	return nil
}

func (o *BlockStmt) WalkForLocals(e *Engine) {
	for _, subblock := range o.Block {
		subblock.WalkForLocals(e)
	}
}

func NewBlockStmt() *BlockStmt {
	return &BlockStmt{
		Block: make([]Stmt, 0),
		Start: Position{Line: 0, Column: 1},
		End:   NO_POS,
		Desc:  CommentGroupExpr{},
	}
}

// Fuck... to do, to do big time

type ImportStmt struct {
	Doc         *CommentGroupExpr `json:"doc"`
	Name        *IdentExpr        `json:"name"`
	CommentExpr *CommentGroupExpr `json:"comment"`
}

func (o *ImportStmt) stmtNode()     {}
func (o *ImportStmt) Pos() Position { return o.Name.Pos() }

func (o *ImportStmt) WalkForLocals(e *Engine) {
	// TODO:
}
