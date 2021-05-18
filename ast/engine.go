package ast

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"
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
	Lc       LineContainer
	Desc     Descs
	Ctr      Counters
	HistRow  string
	GroupCtr int
	RowCtr   int
	GroupMax int
	RowMax   int
}

func MakeCurrentState() CurrentState {
	return CurrentState{
		Lc:       LineContainer{},
		Desc:     Descs{Row: "", Group: "", Block: ""},
		Ctr:      Counters{Stitch: 0, Row: 0, StitchPhrase: 0},
		HistRow:  "",
		GroupCtr: 1,
		RowCtr:   1,
		GroupMax: 0,
		RowMax:   0,
	}
}

func (o CurrentState) String() string {
	return fmt.Sprintf(`----------------------
Block desc.Title = ""
%s
Group desc (%d/%d):
%s
Row desc:
%s
Row:
%s
Args:
%s`,
		o.Desc.Block,
		o.GroupCtr,
		o.GroupMax,
		o.Desc.Group,
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
		defer tmpFile.Close()
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

type IdxAndArgs struct {
	idx  int
	args []string
}

func (e *Engine) FormStates() {
	nestedGroupCtr, nestedRowCtr := 0, 0
	groupStartArr := make([]IdxAndArgs, 0)
	rowStartArr := make([]IdxAndArgs, 0)
	state := MakeCurrentState()
	for i := 0; i < len(e.engineData.Lines); i++ {
		lc := e.engineData.Lines[i]
		if lc.rowIsEqual(START_OF_BLOCK_LC) {
			state.Desc.Block = strings.Join(lc.Desc, "\n")

		} else if lc.rowIsEqual(END_OF_BLOCK_LC) {
			break

		} else if lc.rowIsEqual(START_OF_GROUP_LC) {
			log.WithFields(log.Fields{
				"groupCtr": nestedGroupCtr,
				"desc":     shorten(lc.Desc),
				"args":     shorten(lc.Args),
			}).Debug("[Engine.FormStates] ", strings.Repeat("  ", nestedGroupCtr), "Start of group")
			if len(lc.Desc) != 0 {
				state.Desc.Group = strings.Join(lc.Desc, "\n")
			}
			groupStartArr = append(groupStartArr, IdxAndArgs{
				idx:  len(e.States),
				args: lc.Args,
			})
			if len(lc.Args) == 1 {
				if val, err := strconv.Atoi(lc.Args[0]); err == nil {
					state.GroupMax = val
				}
			}
			nestedGroupCtr += 1

		} else if lc.rowIsEqual(END_OF_GROUP_LC) {
			lastIdx := len(groupStartArr) - 1
			if lastIdx+1 != nestedGroupCtr {
				panic("End of group reached, no groupStartIdxArr")
			}
			var gidxAndArgs IdxAndArgs
			gidxAndArgs, groupStartArr = groupStartArr[lastIdx], groupStartArr[:lastIdx]
			gmax := e.States[gidxAndArgs.idx].GroupMax
			if gmax != 0 {
				slice := append(make([]CurrentState, 0), e.States[gidxAndArgs.idx:len(e.States)]...)
				for i := 0; i < gmax-1; i++ {
					for idx := range slice {
						slice[idx].GroupCtr += 1
					}
					e.States = append(e.States, slice...)
				}
			}
			nestedGroupCtr -= 1
			state.GroupMax = 1
			log.WithField("groupCtr", nestedGroupCtr).Debug("[Engine.FormStates] ", strings.Repeat("  ", nestedGroupCtr), "End of group")

		} else if lc.rowIsEqual(START_OF_ROW_LC) {
			log.WithFields(log.Fields{
				"rowCtr": nestedRowCtr,
				"desc":   shorten(lc.Desc),
			}).Debug("[Engine.FormStates] ", strings.Repeat("  ", nestedRowCtr), "Start of row")
			if len(lc.Desc) > 0 {
				state.Desc.Row = strings.Join(lc.Desc, "\n")
			}
			rowStartArr = append(rowStartArr, IdxAndArgs{
				idx:  len(e.States),
				args: lc.Args,
			})
			if len(lc.Args) == 1 {
				if val, err := strconv.Atoi(lc.Args[0]); err == nil {
					state.RowMax = val
				}
			}
			nestedRowCtr += 1

		} else if lc.rowIsEqual(END_OF_ROW_LC) {
			lastIdx := len(rowStartArr) - 1
			if lastIdx+1 != nestedRowCtr {
				panic("End of row reached, no rowStartArr")
			}
			var idxAndArgs IdxAndArgs
			idxAndArgs, rowStartArr = rowStartArr[lastIdx], rowStartArr[:lastIdx]
			rmax := e.States[idxAndArgs.idx].RowMax
			if rmax != 0 {
				slice := append(make([]CurrentState, 0), e.States[idxAndArgs.idx:len(e.States)]...)
				for i := 0; i < rmax-1; i++ {
					for idx := range slice {
						slice[idx].RowCtr += 1
					}
					e.States = append(e.States, slice...)
				}
			}
			nestedRowCtr -= 1
			state.RowMax = 1
			log.WithField("rowCtr", nestedRowCtr).Debug("[Engine.FormStates] ", strings.Repeat("  ", nestedRowCtr), "End of row")

		} else {
			if len(lc.Row) > 0 {
				state.Lc = lc
				state.HistRow = lc.prettyRow()
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
