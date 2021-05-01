package ast

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
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

// ------------------ LineContainer ------------------

type LineContainer struct {
	Desc []string
	Args []string
	Row  []string
}

func MakeLineContainer() LineContainer {
	return LineContainer{
		Desc: make([]string, 0),
		Args: make([]string, 0),
		Row:  make([]string, 0),
	}
}

func (o *LineContainer) rowIsEqual(lc LineContainer) bool {
	if len(o.Row) != len(lc.Row) {
		return false
	}
	for i, s := range o.Row {
		if s != lc.Row[i] {
			return false
		}
	}
	return true
}

func (o *LineContainer) prettyRow() string {
	var s []string
	for idx, fragment := range o.Row {
		if fragment == "{" || (len(o.Row) > idx+1 && o.Row[idx+1][0] == '}') {
			s = append(s, fragment)
		} else {
			if idx == len(o.Row)-1 {
				s = append(s, fmt.Sprintf("%s", fragment))
			} else {
				s = append(s, fmt.Sprintf("%s,", fragment))
			}
		}
	}
	return strings.Join(s, " ")
}

func (o *LineContainer) isEmpty() bool {
	return len(o.Desc) == 0 && len(o.Args) == 0 && len(o.Row) == 0
}

// ------------------ EngineData ------------------

type EngineData struct {
	Lines       []LineContainer
	aliases     map[string]IdentExpr
	assigns     map[string]*Expr
	nestedRow   bool
	nestedLevel int
}

func NewEngineData() *EngineData {
	return &EngineData{
		aliases:   make(map[string]IdentExpr),
		assigns:   make(map[string]*Expr),
		Lines:     make([]LineContainer, 0),
		nestedRow: false,
		// Nested level may be a little redundant
		nestedLevel: 1,
	}
}

func (e *EngineData) PrintLines() {
	for _, line := range e.Lines {
		if len(line.Desc) > 0 {
			fmt.Println(strings.Join(line.Desc, "\n"))
		}
		if len(line.Row) > 0 {
			fmt.Printf("row: %s\n", line.prettyRow())
		}
		if len(line.Args) > 0 {
			fmt.Printf(" args: %s\n", strings.Join(line.Args, ", "))
		}
		fmt.Println("---")
	}
}

func (e *EngineData) checkAliases(o IdentExpr) IdentExpr {
	if alias, ok := e.aliases[o.Name]; ok {
		return alias
	}
	return o
}

func (e *EngineData) checkAssigns(o *IdentExpr) *Expr {
	if assign, ok := e.assigns[o.Name]; ok {
		return assign
	}
	return nil
}

// ------------------ CurrentState ------------------

type Counters struct {
	Stitch, Row, StitchPhrase int
}

func (c *Counters) reset() {
	c.Stitch = 0
	c.Row = 0
	c.StitchPhrase = 0
}

type Descs struct {
	Row   string
	Group string
	Block string
}

type CurrentState struct {
	Lc             LineContainer
	Desc           Descs
	Ctr            Counters
	HistRow        string
	NestedRowCtr   int
	NestedGroupCtr int
}

func MakeCurrentState() CurrentState {
	return CurrentState{
		Lc:             LineContainer{},
		Desc:           Descs{Row: "", Group: "", Block: ""},
		Ctr:            Counters{Stitch: 0, Row: 0, StitchPhrase: 0},
		HistRow:        "",
		NestedRowCtr:   0,
		NestedGroupCtr: 0,
	}
}

func (o CurrentState) String() string {
	return fmt.Sprintf(`----------------------
Block desc.Title = ""
%s
Group desc [%d]:
%s
Row desc [%d]:
%s
Row:
%s
Args:
%s`,
		o.Desc.Block,
		o.NestedGroupCtr,
		o.Desc.Group,
		o.NestedRowCtr,
		o.Desc.Row,
		o.Lc.prettyRow(),
		strings.Join(o.Lc.Args, ", "),
	)
}

// ------------------ Engine ------------------

type Engine struct {
	States     []CurrentState
	StateIdx   int
	engineData *EngineData
	StatesFile string
}

func MakeEngine(e *EngineData, s string) Engine {
	return Engine{
		States:     make([]CurrentState, 0),
		StateIdx:   0,
		engineData: e,
		StatesFile: s,
	}
}

func MakeEngineFromStatesFile(StatesFile string) (Engine, error) {
	var engine Engine
	statesJson, err := ioutil.ReadFile(StatesFile)
	if err != nil {
		return engine, err
	}
	err = json.Unmarshal([]byte(statesJson), &engine)
	if err != nil {
		return engine, err
	}
	engine.StatesFile = StatesFile
	return engine, err
}

func (e *Engine) WriteEngine() error {
	if e.StatesFile == "" {
		tmpFile, err := ioutil.TempFile(".", "states.*.json")
		if err != nil {
			return err
		}
		e.StatesFile = tmpFile.Name()
	}
	engineJson, err := json.MarshalIndent(e, "", "  ")
	if err == nil {
		err = ioutil.WriteFile(e.StatesFile, engineJson, 0644)
	}
	return err
}

func (e *Engine) PrevState() *CurrentState {
	if e.StateIdx == 0 {
		return &e.States[e.StateIdx]
	}
	e.StateIdx -= 1
	return &e.States[e.StateIdx]
}

func (e *Engine) NextState() *CurrentState {
	if len(e.States)-1 == e.StateIdx {
		return &e.States[e.StateIdx]
	}
	e.StateIdx += 1
	return &e.States[e.StateIdx]
}

func (e *Engine) GotoState(idx int) (*CurrentState, error) {
	if idx >= 0 && len(e.States) > idx {
		e.StateIdx = idx
		return &e.States[e.StateIdx], nil
	}
	return nil, errors.New(fmt.Sprint("Invalid goto value: ", idx))
}

func shorten(desc []string) string {
	if len(desc) == 0 {
		return ""
	}
	if len(desc[0]) > 15 {
		return fmt.Sprint(desc[0][0:13], "...")
	}
	return desc[0]
}

func (e *Engine) FormStates() {
	state := MakeCurrentState()
	for i := 0; i < len(e.engineData.Lines); i++ {
		lc := e.engineData.Lines[i]
		if lc.rowIsEqual(START_OF_BLOCK_LC) {
			state.Desc.Block = strings.Join(lc.Desc, "\n")
		} else if lc.rowIsEqual(END_OF_BLOCK_LC) {
			break
		} else if lc.rowIsEqual(START_OF_GROUP_LC) {
			log.WithFields(log.Fields{
				"groupCtr": state.NestedGroupCtr,
				"desc":     shorten(lc.Desc),
			}).Trace("[Engine.FormStates] ", strings.Repeat("  ", state.NestedGroupCtr), "Start of group")
			if len(lc.Desc) > 0 {
				state.Desc.Group = strings.Join(lc.Desc, "\n")
			}
			state.NestedGroupCtr += 1
		} else if lc.rowIsEqual(END_OF_GROUP_LC) {
			state.NestedGroupCtr -= 1
			log.WithField("groupCtr", state.NestedGroupCtr).Trace("[Engine.FormStates] ", strings.Repeat("  ", state.NestedGroupCtr), "End of group")
		} else if lc.rowIsEqual(START_OF_ROW_LC) {
			log.WithFields(log.Fields{
				"rowCtr": state.NestedRowCtr,
				"desc":   shorten(lc.Desc),
			}).Trace("[Engine.FormStates] ", strings.Repeat("  ", state.NestedRowCtr), "Start of row")
			if len(lc.Desc) > 0 {
				state.Desc.Row = strings.Join(lc.Desc, "\n")
			}
			state.NestedRowCtr += 1
		} else if lc.rowIsEqual(END_OF_ROW_LC) {
			state.NestedRowCtr -= 1
			log.WithField("rowCtr", state.NestedRowCtr).Trace("[Engine.FormStates] ", strings.Repeat("  ", state.NestedRowCtr), "End of row")
		} else {
			if len(lc.Row) > 0 {
				state.Lc = lc
				state.HistRow = fmt.Sprintf("[%s](color:gray)", lc.prettyRow())
				e.States = append(e.States, state)
			}
		}
	}
}

func (e *Engine) PrintEngine() {
	for _, state := range e.States {
		fmt.Println(state)
	}
}
