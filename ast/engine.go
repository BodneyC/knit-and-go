package ast

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
)

var (
	START_OF_GROUP_LC = LineContainer{Row: []string{"START: GROUP"}}
	END_OF_GROUP_LC   = LineContainer{Row: []string{"END: GROUP"}}

	START_OF_ROW_LC = LineContainer{Row: []string{"START: ROW"}}
	END_OF_ROW_LC   = LineContainer{Row: []string{"END: ROW"}}

	START_OF_BLOCK_LC = LineContainer{Row: []string{"START: BLOCK"}}
	END_OF_BLOCK_LC   = LineContainer{Row: []string{"END: BLOCK"}}
)

type LineContainer struct {
	Desc   []string
	Args   []string
	Row    []string
	RowIdx int
}

func MakeLineContainer() LineContainer {
	return LineContainer{
		Desc:   make([]string, 0),
		Args:   make([]string, 0),
		Row:    make([]string, 0),
		RowIdx: 0,
	}
}

func (o *LineContainer) PrettyRow() string {
	var s []string
	for idx, fragment := range o.Row {
		if fragment == "[" || (len(o.Row) > idx+1 && o.Row[idx+1][0] == ']') {
			s = append(s, fragment)
		} else {
			s = append(s, fmt.Sprintf("%s,", fragment))
		}
	}
	return strings.Join(s, " ")
}

func (o *LineContainer) IsEmpty() bool {
	return len(o.Desc) == 0 && len(o.Args) == 0 && len(o.Row) == 0
}

func (o *LineContainer) String() string {
	var s strings.Builder
	s.WriteString(" -- Desc\n")
	for _, desc := range o.Desc {
		s.WriteString(fmt.Sprintf("  %s\n", desc))
	}
	s.WriteString(" -- Args\n")
	for _, arg := range o.Args {
		s.WriteString(fmt.Sprintf("  %s\n", arg))
	}
	s.WriteString(" -- Row\n")
	for _, row := range o.Row {
		s.WriteString(fmt.Sprintf("  %s\n", row))
	}
	return s.String()
}

type Engine struct {
	Lines       []LineContainer
	aliases     map[string]IdentExpr
	assigns     map[string]*Expr
	nestedRow   bool
	nestedLevel int
}

func NewEngine() *Engine {
	return &Engine{
		aliases:   make(map[string]IdentExpr),
		assigns:   make(map[string]*Expr),
		Lines:     make([]LineContainer, 0),
		nestedRow: false,
		// Nested level may be a little redundant
		nestedLevel: 1,
	}
}

func (e *Engine) PrintLines() {
	for _, line := range e.Lines {
		if len(line.Desc) > 0 {
			fmt.Println(strings.Join(line.Desc, "\n"))
		}
		if len(line.Row) > 0 {
			fmt.Printf("row: %s", line.PrettyRow())
		}
		if len(line.Args) > 0 {
			fmt.Printf(" args: %s\n", strings.Join(line.Args, ", "))
		}
		fmt.Println("---")
	}
}

func (e *Engine) PrintEngine() {
	fmt.Println("-- Aliases")
	for k, v := range e.aliases {
		fmt.Printf("%s : %v\n", k, v)
	}
	fmt.Println("-- Assignments")
	for k, v := range e.assigns {
		fmt.Printf("%s : %v\n", k, v)
	}
	fmt.Println("-- Lines")
	for _, v := range e.Lines {
		fmt.Println(v.String())
	}
}

func (e *Engine) checkAliases(o IdentExpr) IdentExpr {
	if alias, ok := e.aliases[o.Name]; ok {
		return alias
	}
	return o
}

func (e *Engine) checkAssigns(o *IdentExpr) *Expr {
	log.Debugf("Checking assignment for %s\n", o.Name)
	if assign, ok := e.assigns[o.Name]; ok {
		return assign
	}
	return nil
}
