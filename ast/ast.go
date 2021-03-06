package ast

import (
	"strings"

	. "github.com/bodneyc/knit-and-go/lexer"
)

type Node interface {
	Pos() Position
	WalkForLocals(*EngineData)
}

type Expr interface {
	Node
	exprNode()
	Text(*EngineData) string
	WalkForLines(*EngineData, *LineContainer) error
}

type Stmt interface {
	Node
	stmtNode()
	WalkForLines(*EngineData) error
}

// ----- Bracket group, implements nothing -----

type Brackets struct {
	Args []Expr `json:"args"`
}

func (o *Brackets) WalkForLines(e *EngineData, lc *LineContainer) error {
	lc.Args = append(lc.Args, o.TextSlice(e)...)
	return nil
}

func (o *Brackets) GetSizeInt() int {
	i := 1
	for _, arg := range o.Args {
		switch arg.(type) {
		case *SizeExpr:
			argSize := arg.(*SizeExpr)
			i = argSize.GetSizeInt()
		}
	}
	return i
}

func (o *Brackets) GetSizeText(e *EngineData) string {
	if len(o.Args) == 1 {
		return o.Args[0].Text(e)
	}
	var s []string
	for _, arg := range o.Args {
		s = append(s, arg.Text(e))
	}
	return strings.Join(s, " ")
}

func (o *Brackets) TextSlice(e *EngineData) []string {
	s := make([]string, 0)
	for _, arg := range o.Args {
		s = append(s, arg.Text(e))
	}
	return s
}

func MakeBrackets() Brackets {
	return Brackets{
		Args: make([]Expr, 0),
	}
}

func (o *Brackets) WalkForLocals(e *EngineData) {
	for _, args := range o.Args {
		args.WalkForLocals(e)
	}
}
