package tui

import (
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/dilsonrabelo/castor-prompt-builder/pkg/engine"
	"github.com/dilsonrabelo/castor-prompt-builder/pkg/parser"
)

func New(models []*parser.Model, roles []*parser.Role) AppModel {
	ti := textinput.New()
	ti.Placeholder = "Escreva aqui..."
	ti.Focus()

	ta := textarea.New()
	ta.Placeholder = "Descreva a tarefa livremente...\n(Ctrl+S para confirmar)"
	ta.SetWidth(50)
	ta.SetHeight(8)
	ta.Focus()

	return AppModel{
		screen:        screenSelectModel,
		models:        models,
		roles:         roles,
		values:        engine.NewValues(),
		textInput:     ti,
		textArea:      ta,
		selectedRoles: make(map[int]bool),
	}
}

func (m AppModel) Init() tea.Cmd {
	return textinput.Blink
}
