package tui

import (
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/dilsonrabelo/castor-prompt-builder/internal/parser"
)

func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.textArea.SetWidth(msg.Width - 6)
		return m, nil

	case tea.KeyMsg:
		switch m.screen {
		case screenSelectModel:
			return m.updateSelectModel(msg)
		case screenModelInfo:
			return m.updateModelInfo(msg)
		case screenSelectRole:
			return m.updateSelectRole(msg)
		case screenNarrative:
			return m.updateNarrative(msg)
		case screenGap:
			return m.updateGap(msg)
		case screenAskPhase:
			return m.updateAskPhase(msg)
		case screenDefinePhase:
			return m.updateDefinePhase(msg)
		case screenDone:
			return m.updateDone(msg)
		}
	}
	return m, nil
}

// --- Select Model ---

func (m AppModel) updateSelectModel(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit
	case "up", "k":
		if m.selectedModel > 0 {
			m.selectedModel--
		}
	case "down", "j":
		if m.selectedModel < len(m.models)-1 {
			m.selectedModel++
		}
	case "i":
		m.screen = screenModelInfo
	case "enter":
		m.screen = screenSelectRole
		m.selectedRole = 0
	}
	return m, nil
}

func (m AppModel) updateModelInfo(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit
	case "esc":
		m.screen = screenSelectModel
	case "enter":
		m.screen = screenSelectRole
		m.selectedRole = 0
	}
	return m, nil
}

// --- Select Role ---

func (m AppModel) updateSelectRole(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit
	case "esc":
		m.screen = screenSelectModel
	case "up", "k":
		if m.selectedRole > 0 {
			m.selectedRole--
		}
	case "down", "j":
		if m.selectedRole < len(m.roles)-1 {
			m.selectedRole++
		}
	case "enter":
		m.screen = screenNarrative
		m.textArea.Reset()
		m.textArea.Focus()
	}
	return m, nil
}

// --- Narrative ---

func (m AppModel) updateNarrative(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit
	case "esc":
		m.screen = screenSelectRole
		return m, nil
	case "ctrl+s":
		m.narrative = m.textArea.Value()
		if strings.TrimSpace(m.narrative) == "" {
			return m, nil
		}
		// prepara gaps do role selecionado
		role := m.roles[m.selectedRole]
		m.gaps = role.GapsComuns
		m.gapIndex = 0
		m.gapAnswers = make([]string, len(m.gaps))
		if len(m.gaps) > 0 {
			m.screen = screenGap
			m.textInput.Reset()
			m.textInput.Focus()
		} else {
			m.screen = screenAskPhase
		}
		return m, nil
	}
	var cmd tea.Cmd
	m.textArea, cmd = m.textArea.Update(msg)
	return m, cmd
}

// --- Gap ---

func (m AppModel) updateGap(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit
	case "esc":
		if m.gapIndex > 0 {
			m.gapIndex--
			m.textInput.SetValue(m.gapAnswers[m.gapIndex])
		} else {
			m.screen = screenNarrative
		}
		return m, nil
	case "enter", "tab":
		m.gapAnswers[m.gapIndex] = m.textInput.Value()
		m.gapIndex++
		if m.gapIndex >= len(m.gaps) {
			m.screen = screenAskPhase
		} else {
			m.textInput.Reset()
		}
		return m, nil
	}
	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

// --- Ask Phase ---

func (m AppModel) updateAskPhase(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit
	case "esc":
		if len(m.gaps) > 0 {
			m.screen = screenGap
			m.gapIndex = len(m.gaps) - 1
		} else {
			m.screen = screenNarrative
		}
	case "up", "k":
		if m.askPhaseChoice > 0 {
			m.askPhaseChoice--
		}
	case "down", "j":
		if m.askPhaseChoice < 1 {
			m.askPhaseChoice++
		}
	case "enter":
		if m.askPhaseChoice == 0 { // sim
			m.textInput.Reset()
			m.textInput.Placeholder = "Quantidade de fases (ex: 3)"
			m.textInput.Focus()
			m.screen = screenDefinePhase
			m.phaseIndex = -1 // -1 = coletando quantidade
		} else {
			m = m.buildAndSave()
		}
	}
	return m, nil
}

// --- Define Phase ---

func (m AppModel) updateDefinePhase(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit
	case "esc":
		if m.phaseIndex <= 0 {
			m.screen = screenAskPhase
			return m, nil
		}
		m.phaseIndex--
		return m, nil
	case "enter", "tab":
		if m.phaseIndex == -1 {
			// coletando quantidade
			n, err := strconv.Atoi(strings.TrimSpace(m.textInput.Value()))
			if err != nil || n < 1 {
				return m, nil
			}
			m.phaseCount = n
			m.phaseIndex = 0
			m.phaseEditField = 0
			// inicializa steps vazios
			steps := make([]parser.Step, n)
			m.values.Steps["fases"] = steps
			m.textInput.Reset()
			m.textInput.Placeholder = "Título da fase"
			return m, nil
		}
		if m.phaseEditField == 0 {
			// salvou título, vai pra descrição
			m.values.Steps["fases"][m.phaseIndex].Titulo = m.textInput.Value()
			m.phaseEditField = 1
			m.textArea.Reset()
			m.textArea.Focus()
			return m, nil
		}
		// salvou descrição
		m.values.Steps["fases"][m.phaseIndex].Descricao = m.textArea.Value()
		m.phaseIndex++
		if m.phaseIndex >= m.phaseCount {
			m = m.buildAndSave()
			return m, nil
		}
		m.phaseEditField = 0
		m.textInput.Reset()
		m.textInput.Placeholder = "Título da fase"
		return m, nil
	}

	if m.phaseIndex >= 0 && m.phaseEditField == 1 {
		var cmd tea.Cmd
		m.textArea, cmd = m.textArea.Update(msg)
		return m, cmd
	}
	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

// --- Done ---

func (m AppModel) updateDone(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q", "enter":
		return m, tea.Quit
	}
	return m, nil
}

