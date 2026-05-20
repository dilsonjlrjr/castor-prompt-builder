package tui

import (
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/dilsonrabelo/castor-prompt-builder/pkg/engine"
	"github.com/dilsonrabelo/castor-prompt-builder/pkg/parser"
)

type screen int

const (
	screenSelectModel screen = iota
	screenModelInfo
	screenSelectRole
	screenNarrative
	screenContext
	screenAskPhase
	screenDefinePhase
	screenDone
)

type AppModel struct {
	screen screen
	width  int
	height int

	models []*parser.Model
	roles  []*parser.Role

	selectedModel int
	roleCursor    int
	selectedRoles map[int]bool
	roleSearch    string

	textInput textinput.Model
	textArea  textarea.Model

	narrative string
	values    *engine.Values

	// contexto (antigo screenGap)
	contextSections  []ContextSection
	contextSecIdx    int
	contextCursor    int
	contextEditing   bool
	contextMultiIdx  int    // indice do ciclo para multiselect
	contextRoleFilter int   // -1=todos, 0+=indice da secao de papel

	// scroll de telas
	modelInfoOffset int

	// fases
	askPhaseChoice int
	phaseCount     int
	phaseIndex     int
	phaseTitle     string
	phaseEditField int

	// resultado
	savedPath string
	err       error
}

// ContextGap representa uma pergunta de contexto com sua resposta.
type ContextGap struct {
	FieldID     string
	Pergunta    string
	RoleNome    string
	Obrigatorio bool
	Tipo        parser.FieldType
	Opcoes      []string
	Answer      string
}

// ContextSection agrupa lacunas de contexto — do modelo ou de um papel.
type ContextSection struct {
	Title string
	Kind  string // "model" ou "role"
	Gaps  []*ContextGap
}
