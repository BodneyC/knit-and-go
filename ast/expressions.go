package ast

import (
	"fmt"
	"strings"

	. "github.com/bodneyc/knit-and-go/lexer"
	"github.com/bodneyc/knit-and-go/util"
)

// ------------------ CommentExpr ------------------

type CommentExpr struct {
	Semicolon Position `json:"semicolon"`
	Str       string   `json:"text"`
}

func (o *CommentExpr) exprNode()     {}
func (o *CommentExpr) Pos() Position { return o.Semicolon }

func (o *CommentExpr) WalkForLines(e *EngineData, lc *LineContainer) error { return nil }

func (o *CommentExpr) WalkForLocals(e *EngineData) {}
func (o *CommentExpr) Text(e *EngineData) string   { return o.Str }

// ---------------- CommentGroupExpr ---------------

type CommentGroupExpr struct {
	List []CommentExpr `json:"list"`
}

func MakeCommentGroupExpr() CommentGroupExpr {
	return CommentGroupExpr{
		List: make([]CommentExpr, 0),
	}
}

func (o *CommentGroupExpr) exprNode()     {}
func (o *CommentGroupExpr) Pos() Position { return o.List[0].Pos() }

func (o *CommentGroupExpr) WalkForLocals(e *EngineData) {}

func (o *CommentGroupExpr) WalkForLines(e *EngineData, lc *LineContainer) error { return nil }

func (o *CommentGroupExpr) Text(e *EngineData) string {
	return strings.Join(o.TextSlice(e), "\n")
}
func (o *CommentGroupExpr) TextSlice(e *EngineData) []string {
	if o == nil {
		return []string{}
	}
	comments := make([]string, len(o.List))
	for i, c := range o.List {
		comments[i] = c.Text(e)
	}
	lines := make([]string, 0, 10)
	for _, c := range comments {
		if len(c) > 0 && c[0] == ';' {
			c = c[1:]
		}
		lines = append(lines, strings.Trim(c, "\n\t "))
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
	return lines
}

// ----------------- IdentExpr ----------------

type IdentExpr struct {
	At   Position `json:"at"`
	Name string   `json:"name"`
}

func MakeIdentExpr(t TokenContainer) IdentExpr {
	return IdentExpr{
		At:   t.Pos,
		Name: t.Str,
	}
}

func NewIdentExpr(t TokenContainer) *IdentExpr {
	return &IdentExpr{
		At:   t.Pos,
		Name: t.Str,
	}
}

func (o *IdentExpr) exprNode()     {}
func (o *IdentExpr) Pos() Position { return o.At }

// Only for assignment values
func (o *IdentExpr) WalkForLines(e *EngineData, lc *LineContainer) error {
	if assign := e.checkAssigns(o); assign != nil {
		if err := (*assign).WalkForLines(e, lc); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("No assignment for %#v", o)
	}
	return nil
}

func (o *IdentExpr) AliasForLines(e *EngineData, lc *LineContainer, size string) error {
	alias := e.checkAliases(*o)
	aliasText := alias.Text(e)
	fragment := aliasText
	if size != "" {
		fragment = fmt.Sprintf("%s %s", fragment, size)
	}
	lc.Row = append(lc.Row, fragment)
	return nil
}

func (o *IdentExpr) WalkForLocals(e *EngineData) {}

// No assignment should exist if calling this function
func (o *IdentExpr) Text(e *EngineData) string {
	return e.checkAliases(*o).Name
}

// --------------- Single stitch ---------------

type StitchExpr struct {
	At   Position  `json:"at"`
	Id   IdentExpr `json:"stitch"`
	Args Brackets  `json:"args"`
}

func NewStitchExpr(ident IdentExpr, args Brackets) *StitchExpr {
	return &StitchExpr{
		At:   ident.Pos(),
		Id:   ident,
		Args: args,
	}
}

func (o *StitchExpr) exprNode()     {}
func (o *StitchExpr) Pos() Position { return o.At }

func (o *StitchExpr) WalkForLocals(e *EngineData) {}
func (o *StitchExpr) Text(e *EngineData) string   { return "" }

func (o *StitchExpr) WalkForLines(e *EngineData, lc *LineContainer) error {
	if assign := e.checkAssigns(&o.Id); assign != nil {
		o.Args.WalkForLines(e, lc)
		(*assign).WalkForLines(e, lc)
	} else {
		size := o.Args.GetSizeText(e)
		o.Id.AliasForLines(e, lc, size)
	}
	return nil
}

// ----------------- SizeExpr ----------------

type MeasurementUnit int

const (
	NOUNIT MeasurementUnit = iota
	MM
	CM
	INCHES
	FEET
	ASTERISK
)

type SizeExpr struct {
	At     Position        `json:"at"`
	Ni     int64           `json:"ni,string"`
	Nf     float64         `json:"nf,string"`
	Id     IdentExpr       `json:"nid"`
	Before bool            `json:"before"`
	Unit   MeasurementUnit `json:"unit,string"`
}

func NewSizeExpr(ni int64, nf float64, nid TokenContainer, unit MeasurementUnit) *SizeExpr {
	return &SizeExpr{
		At:     nid.Pos,
		Ni:     ni,
		Nf:     nf,
		Id:     MakeIdentExpr(nid),
		Before: false,
		Unit:   unit,
	}
}

func NewSizeExprAsterisk(asterisk TokenContainer) *SizeExpr {
	return &SizeExpr{
		At:     asterisk.Pos,
		Ni:     -1,
		Nf:     -1,
		Before: false,
		Id:     IdentExpr{},
		Unit:   ASTERISK,
	}
}

func (o *SizeExpr) exprNode()     {}
func (o *SizeExpr) Pos() Position { return o.At }

func (o *SizeExpr) WalkForLocals(e *EngineData) {}

func (o *SizeExpr) WalkForLines(e *EngineData, lc *LineContainer) error { return nil }

func (o *SizeExpr) GetSizeInt() int {
	if o.Unit == NOUNIT && !o.Before {
		return int(o.Ni)
	}
	return -1
}

func (o *SizeExpr) Text(e *EngineData) string {
	var s string

	if o.Unit == ASTERISK {
		if o.Before {
			return "before end of row"
		} else {
			return "to end of row"
		}
	}

	// Something should be in `s` unless `ASTERISK`
	if o.Ni != -1 {
		s = fmt.Sprintf("%d", o.Ni)
	} else if o.Nf != -1 {
		s = fmt.Sprintf("%.2f", o.Nf)
	} else if o.Id.Name != "" {
		s = e.checkAliases(o.Id).Name
	}

	// Append unit specifier
	switch o.Unit {
	case NOUNIT:
		s = fmt.Sprintf("%s", s)
	case MM:
		s = fmt.Sprintf("%smm", s)
	case CM:
		s = fmt.Sprintf("%scm", s)
	case INCHES:
		s = fmt.Sprintf("%s\"", s)
	case FEET:
		s = fmt.Sprintf("%s'", s)
	}

	if o.Before {
		s = fmt.Sprintf("until %s", s)
	}

	return s
}

// -------------------- Row --------------------

type RowExpr struct {
	Stitches []Expr   `json:"stitches"`
	Args     Brackets `json:"args"`
}

func NewRowExpr(exprs []Expr, args Brackets) *RowExpr {
	return &RowExpr{
		Stitches: exprs,
		Args:     args,
	}
}

func (o *RowExpr) exprNode()     {}
func (o *RowExpr) Pos() Position { return o.Stitches[0].Pos() }

func (o *RowExpr) Text(e *EngineData) string { return "" }

func (o *RowExpr) WalkForLines(e *EngineData, lc *LineContainer) error {
	if !e.nestedRow {
		o.Args.WalkForLines(e, lc)
	}
	for _, stitch := range o.Stitches {
		switch stitch.(type) {
		case *RowExpr:
			rowExpr := stitch.(*RowExpr)
			lc.Row = append(lc.Row, "{")
			e.nestedRow = true
			e.nestedLevel += 1
			if err := stitch.WalkForLines(e, lc); err != nil {
				return fmt.Errorf("%w%s", err, util.StackLine())
			}
			e.nestedLevel -= 1
			e.nestedRow = false
			endWrap := "}"
			if size := rowExpr.Args.GetSizeText(e); size != "" {
				endWrap = fmt.Sprintf("%s %s", endWrap, size)
			}
			lc.Row = append(lc.Row, endWrap)
		default:
			if err := stitch.WalkForLines(e, lc); err != nil {
				return fmt.Errorf("%w%s", err, util.StackLine())
			}
		}
	}
	return nil
}

func (o *RowExpr) WalkForLocals(e *EngineData) {
	for _, stitch := range o.Stitches {
		stitch.WalkForLocals(e)
	}
}

// ------------------- Group -------------------

type GroupExpr struct {
	LBrace Position `json:"lbrace"`
	RBrace Position `json:"rbrace"`
	Args   Brackets `json:"args"`
	Lines  []Stmt   `json:"lines"`
}

func NewGroupExpr(l Position, r Position, lines []Stmt, args Brackets) *GroupExpr {
	return &GroupExpr{
		LBrace: l,
		RBrace: r,
		Lines:  lines,
		Args:   args,
	}
}

func (o *GroupExpr) exprNode()     {}
func (o *GroupExpr) Pos() Position { return o.LBrace }

func (o *GroupExpr) Text(e *EngineData) string { return "" }

func (o *GroupExpr) WalkForLines(e *EngineData, lc *LineContainer) error {
	for _, line := range o.Lines {
		if err := line.WalkForLines(e); err != nil {
			return fmt.Errorf("%w%s", err, util.StackLine())
		}
	}
	return nil
}

func (o *GroupExpr) WalkForLocals(e *EngineData) {
	for _, line := range o.Lines {
		line.WalkForLocals(e)
	}
}
