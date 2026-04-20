package tui

import (
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/dilsonrabelo/castor-prompt-builder/internal/engine"
	"github.com/dilsonrabelo/castor-prompt-builder/internal/parser"
)

type screen int

const (
	screenSelectModel screen = iota
	screenModelInfo
	screenSelectRole
	screenNarrative
	screenGap
	screenAskPhase
	screenDefinePhase
	screenDone
)

// AppModel é o modelo principal do wizard bubbletea
type AppModel struct {
	screen screen
	width  int
	height int

	// dados carregados
	models []*parser.Model
	roles  []*parser.Role

	// seleções
	selectedModel int
	// multi-select de papéis
	roleCursor    int
	selectedRoles map[int]bool
	roleSearch    string

	// inputs
	textInput textinput.Model
	textArea  textarea.Model

	// valores coletados
	narrative    string
	values       *engine.Values
	gaps         []string   // perguntas pendentes
	gapIndex     int        // gap atual
	gapAnswers   []string   // respostas

	// fases
	askPhaseChoice int  // 0=não definido, 1=sim, 2=não
	phaseCount     int
	phaseIndex     int
	phaseTitle     string
	phaseEditField int // 0=titulo, 1=descricao

	// resultado
	savedPath string
	err       error
}
