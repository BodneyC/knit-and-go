package tui

import "github.com/bodneyc/knit-and-go/ast"

type KnownTokens string

const (
	CAST_ON_KT  KnownTokens = "cast-on"
	CAST_OFF_KT             = "cast-off"
	USE_KT                  = "use"
)

type Screen struct {
	engine *ast.Engine
}

func NewScreen(engine *ast.Engine) *Screen {
	return &Screen{engine: engine}
}
