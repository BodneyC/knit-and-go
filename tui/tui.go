package tui

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bodneyc/knit-and-go/ast"
	"github.com/bodneyc/knit-and-go/util"

	ui "github.com/gizak/termui/v3"
	w "github.com/gizak/termui/v3/widgets"
	log "github.com/sirupsen/logrus"
)

type KnownTokens string

const (
	CAST_ON_KT  KnownTokens = "cast-on"
	CAST_OFF_KT             = "cast-off"
	USE_KT                  = "use"
)

type Screen struct {
	engine *ast.Engine
	keymapsPar,
	blockDescPar,
	groupDescPar,
	rowDescPar,
	rowCtrPar,
	stitchCtrPar,
	stateCtrPar,
	prevRow,
	nextRow,
	currentRowPar,
	argsPar *w.Paragraph
}

func NewScreen(engine *ast.Engine) *Screen {
	return &Screen{engine: engine}
}

func prettyRowWithHighlight(state *ast.CurrentState) string {
	var s []string
	for idx, fragment := range state.Lc.Row {
		if fragment == "{" || (len(state.Lc.Row) > idx+1 && state.Lc.Row[idx+1][0] == '}') {
			if idx == state.Ctr.StitchPhrase {
				fragment = fmt.Sprintf("[%s](fg:red)", fragment)
			}
			s = append(s, fragment)
		} else {
			if idx == state.Ctr.StitchPhrase {
				fragment = fmt.Sprintf("[%s](fg:red)", fragment)
			}
			if idx == len(state.Lc.Row)-1 {
				s = append(s, fmt.Sprintf("%s", fragment))
			} else {
				s = append(s, fmt.Sprintf("%s,", fragment))
			}
		}
	}
	return strings.Join(s, " ")
}

func (s *Screen) paragraphSetup() {
	s.keymapsPar = w.NewParagraph()
	s.keymapsPar.Title = "Keymaps"
	s.keymapsPar.TitleStyle.Modifier = ui.ModifierBold
	s.keymapsPar.Text = `q: quit
j: next
k: previous
l: right
h: left
s: stitch up
S: stitch down
r: row up
R: row down
x: ctr reset
^s: save`

	s.blockDescPar = w.NewParagraph()
	s.blockDescPar.Title = "Descriptions"
	s.blockDescPar.TitleStyle.Modifier = ui.ModifierBold
	s.blockDescPar.BorderBottom = false

	s.groupDescPar = w.NewParagraph()
	s.groupDescPar.Title = "Group description:"
	s.groupDescPar.TitleStyle.Modifier = ui.ModifierBold
	s.groupDescPar.BorderTop = false
	s.groupDescPar.BorderBottom = false

	s.rowDescPar = w.NewParagraph()
	s.rowDescPar.Title = "Row description:"
	s.rowDescPar.TitleStyle.Modifier = ui.ModifierBold
	s.rowDescPar.BorderTop = false

	s.rowCtrPar = w.NewParagraph()
	s.rowCtrPar.Title = "Row counter"
	s.rowCtrPar.TitleStyle.Modifier = ui.ModifierBold

	s.stitchCtrPar = w.NewParagraph()
	s.stitchCtrPar.Title = "Stitch counter"
	s.stitchCtrPar.TitleStyle.Modifier = ui.ModifierBold

	s.stateCtrPar = w.NewParagraph()
	s.stateCtrPar.Title = "Page counter"
	s.stateCtrPar.TitleStyle.Modifier = ui.ModifierBold

	s.currentRowPar = w.NewParagraph()
	s.currentRowPar.Title = "Current row"
	s.currentRowPar.TitleStyle.Modifier = ui.ModifierBold
	s.currentRowPar.BorderBottom = false

	s.argsPar = w.NewParagraph()
	s.argsPar.Title = "Args:"
	s.argsPar.TitleStyle.Modifier = ui.ModifierBold
	s.argsPar.BorderTop = false

	s.nextRow = w.NewParagraph()
	s.nextRow.Title = "Next row"
	s.nextRow.TitleStyle.Modifier = ui.ModifierBold

	s.prevRow = w.NewParagraph()
	s.prevRow.Title = "Previous row"
	s.prevRow.TitleStyle.Modifier = ui.ModifierBold
}

func (s *Screen) setParagraphs(state *ast.CurrentState) error {
	s.blockDescPar.Text = state.Desc.Block
	s.groupDescPar.Text = state.Desc.Group
	s.rowDescPar.Text = state.Desc.Row
	s.rowCtrPar.Text = strconv.Itoa(state.Ctr.Row)
	s.stitchCtrPar.Text = strconv.Itoa(state.Ctr.Stitch)
	s.stateCtrPar.Text = fmt.Sprintf("%d/%d", s.engine.StateIdx, len(s.engine.States)-1)
	s.currentRowPar.Text = prettyRowWithHighlight(state)
	s.argsPar.Text = strings.Join(state.Lc.Args, ", ")

	if s.engine.StateIdx-1 >= 0 {
		s.prevRow.Text = s.engine.States[s.engine.StateIdx-1].HistRow
	} else {
		s.prevRow.Text = "Start of pattern"
	}
	if s.engine.StateIdx+1 < len(s.engine.States) {
		s.nextRow.Text = s.engine.States[s.engine.StateIdx+1].HistRow
	} else {
		s.nextRow.Text = "End of pattern"
	}

	return nil
}

func (s *Screen) Run() (*util.LogrusCalls, error) {
	if err := ui.Init(); err != nil {
		log.Fatalf("Failed to initialize termui: %v", err)
	}
	defer ui.Close()

	s.paragraphSetup()

	state := &s.engine.States[s.engine.StateIdx]
	s.setParagraphs(state)

	width, height := ui.TerminalDimensions()
	grid := ui.NewGrid()
	grid.SetRect(0, 0, width, height)

	descGrid := ui.NewGrid()
	descGrid.Set()

	grid.Set(
		ui.NewRow(0.5,
			ui.NewCol(0.8,
				// Descriptions
				ui.NewRow(0.4, s.blockDescPar),
				ui.NewRow(0.3, s.groupDescPar),
				ui.NewRow(0.3, s.rowDescPar),
			),
			ui.NewCol(0.2, s.keymapsPar),
		),
		ui.NewRow(0.4,
			ui.NewCol(1.0,
				// Rows
				ui.NewRow(0.25, s.prevRow),
				ui.NewRow(0.3, s.currentRowPar),
				ui.NewRow(0.2, s.argsPar),
				ui.NewRow(0.25, s.nextRow),
			),
		),
		ui.NewRow(0.1,
			// Counters
			ui.NewCol(1.0/3, s.stitchCtrPar),
			ui.NewCol(1.0/3, s.rowCtrPar),
			ui.NewCol(1.0/3, s.stateCtrPar),
		),
	)

	ui.Render(grid)

	logCalls := util.NewLogrusCalls()

	events := ui.PollEvents()
	for {
		e := <-events
		switch e.ID {
		case "q", "<C-c>":
			return logCalls, nil

		case "n", "j", "<Down>":
			state = s.engine.NextState()
			logCalls.Trace = append(logCalls.Trace, util.MakeLogrusCall(
				log.WithField("state", s.engine.StateIdx),
				"Moved to next state",
			))

		case "p", "N", "k", "<Up>":
			state = s.engine.PrevState()
			logCalls.Trace = append(logCalls.Trace, util.MakeLogrusCall(
				log.WithField("state", s.engine.StateIdx),
				"Moved to prev state",
			))

		case "<C-s>":
			if err := s.engine.WriteEngine(); err != nil {
				return logCalls, err
			}
			logCalls.Info = append(logCalls.Info, util.MakeLogrusCall(
				log.WithField("statesfile", s.engine.StatesFile),
				"States saved",
			))

		case "l", "<Right>":
			if len(state.Lc.Row) > state.Ctr.StitchPhrase+1 {
				state.Ctr.StitchPhrase += 1
				logCalls.Trace = append(logCalls.Trace, util.MakeLogrusCall(
					log.WithField("stitch", s.engine.States[s.engine.StateIdx].Ctr.StitchPhrase),
					"Moved to right stitch",
				))
			} else {
				logCalls.Trace = append(logCalls.Trace, util.MakeLogrusCall(
					log.WithField("stitch", s.engine.States[s.engine.StateIdx].Ctr.StitchPhrase),
					"Already at rightmost stitch",
				))
			}

		case "h", "<Left>":
			if state.Ctr.StitchPhrase-1 >= 0 {
				state.Ctr.StitchPhrase -= 1
				logCalls.Trace = append(logCalls.Trace, util.MakeLogrusCall(
					log.WithField("stitch", s.engine.States[s.engine.StateIdx].Ctr.StitchPhrase),
					"Moved to left stitch",
				))
			} else {
				logCalls.Trace = append(logCalls.Trace, util.MakeLogrusCall(
					log.WithField("stitch", s.engine.States[s.engine.StateIdx].Ctr.StitchPhrase),
					"Already at leftmost stitch",
				))
			}

		case "s":
			state.Ctr.Stitch += 1
			logCalls.Trace = append(logCalls.Trace, util.MakeLogrusCall(
				log.WithField("stitch", s.engine.States[s.engine.StateIdx].Ctr.Stitch),
				"Increased stitch",
			))

		case "S":
			if state.Ctr.Stitch-1 >= 0 {
				state.Ctr.Stitch -= 1
				logCalls.Trace = append(logCalls.Trace, util.MakeLogrusCall(
					log.WithField("stitch", s.engine.States[s.engine.StateIdx].Ctr.Stitch),
					"Decreased stitch counter",
				))
			} else {
				logCalls.Trace = append(logCalls.Trace, util.MakeLogrusCall(
					log.WithField("stitch", s.engine.States[s.engine.StateIdx].Ctr.Stitch),
					"Cannot decrease stitch counter further",
				))
			}

		case "r":
			state.Ctr.Row += 1
			logCalls.Trace = append(logCalls.Trace, util.MakeLogrusCall(
				log.WithField("row", s.engine.States[s.engine.StateIdx].Ctr.Row),
				"Increased row counter",
			))

		case "R":
			if state.Ctr.Row-1 >= 0 {
				state.Ctr.Row -= 1
				logCalls.Trace = append(logCalls.Trace, util.MakeLogrusCall(
					log.WithField("row", s.engine.States[s.engine.StateIdx].Ctr.Row),
					"Decreased row counter",
				))
			} else {
				logCalls.Trace = append(logCalls.Trace, util.MakeLogrusCall(
					log.WithField("row", s.engine.States[s.engine.StateIdx].Ctr.Row),
					"Cannot decrease row counter further",
				))
			}

		case "x":
			state.Ctr.Stitch = 0
			state.Ctr.Row = 0
			logCalls.Trace = append(logCalls.Trace, util.MakeLogrusCall(
				log.WithFields(log.Fields{
					"row": s.engine.States[s.engine.StateIdx].Ctr.Row,
					"stitch": s.engine.States[s.engine.StateIdx].Ctr.Stitch,
				}),
				"Reset counters",
			))

		case "<Resize>":
			payload := e.Payload.(ui.Resize)
			grid.SetRect(0, 0, payload.Width, payload.Height)
			logCalls.Trace = append(logCalls.Trace, util.MakeLogrusCall(
				log.WithFields(log.Fields{
					"width": payload.Width,
					"height": payload.Height,
				}),
				"Screen resize",
			))
			ui.Clear()
			ui.Render(grid)
		}

		s.setParagraphs(state)
		ui.Render(grid)
	}
}
