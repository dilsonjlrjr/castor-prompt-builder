package tui

import (
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/dilsonrabelo/castor-prompt-builder/pkg/parser"
)

func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.textArea.SetWidth(msg.Width - 6)
		// textarea height: efectiveHeight − cabeçalho(4) − rodapé(2) − margens(4)
		taH := msg.Height - 10
		if taH < 4 {
			taH = 4
		}
		if taH > 12 {
			taH = 12
		}
		m.textArea.SetHeight(taH)
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
		m.modelInfoOffset = 0
	case "enter":
		m.screen = screenSelectRole
		m.roleCursor = 0
		m.selectedRoles = make(map[int]bool)
		m.roleSearch = ""
		m.textInput.Reset()
		m.textInput.Placeholder = "🔍 Buscar papel..."
		cmd := m.textInput.Focus()
		return m, cmd
	}
	return m, nil
}

func (m AppModel) updateModelInfo(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit
	case "esc":
		m.screen = screenSelectModel
	case "up", "k":
		if m.modelInfoOffset > 0 {
			m.modelInfoOffset--
		}
	case "down", "j":
		m.modelInfoOffset++
	case "enter":
		m.screen = screenSelectRole
		m.roleCursor = 0
		m.selectedRoles = make(map[int]bool)
		m.roleSearch = ""
		m.textInput.Reset()
		m.textInput.Placeholder = "🔍 Buscar papel..."
		cmd := m.textInput.Focus()
		return m, cmd
	}
	return m, nil
}

// --- Select Role (multi-select + busca) ---

func (m AppModel) filteredRoleIndices() []int {
	if m.roleSearch == "" {
		idxs := make([]int, len(m.roles))
		for i := range m.roles {
			idxs[i] = i
		}
		return idxs
	}
	q := strings.ToLower(m.roleSearch)
	var idxs []int
	for i, r := range m.roles {
		if strings.Contains(strings.ToLower(r.Nome), q) ||
			strings.Contains(strings.ToLower(r.Categoria), q) {
			idxs = append(idxs, i)
		}
	}
	return idxs
}

func (m AppModel) updateSelectRole(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	filtered := m.filteredRoleIndices()

	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit
	case "esc":
		if m.roleSearch != "" {
			m.roleSearch = ""
			m.textInput.Reset()
			m.roleCursor = 0
			return m, nil
		}
		m.screen = screenSelectModel
		return m, nil
	case "up", "k":
		if m.roleCursor > 0 {
			m.roleCursor--
		}
		return m, nil
	case "down", "j":
		if m.roleCursor < len(filtered)-1 {
			m.roleCursor++
		}
		return m, nil
	case " ":
		if len(filtered) > 0 && m.roleCursor < len(filtered) {
			globalIdx := filtered[m.roleCursor]
			m.selectedRoles[globalIdx] = !m.selectedRoles[globalIdx]
		}
		return m, nil
	case "enter":
		if len(m.selectedRoles) > 0 {
			m.screen = screenNarrative
			m.textArea.Reset()
			m.textArea.Focus()
		}
		return m, nil
	}

	// demais teclas → busca
	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)
	newSearch := m.textInput.Value()
	if newSearch != m.roleSearch {
		m.roleSearch = newSearch
		m.roleCursor = 0
	}
	return m, cmd
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
		// prepara gaps combinados dos papéis selecionados
		var allGaps []string
		for idx, sel := range m.selectedRoles {
			if sel {
				allGaps = append(allGaps, m.roles[idx].GapsComuns...)
			}
		}
		m.gaps = unique(allGaps)
		m.gapIndex = 0
		m.gapAnswers = make([]string, len(m.gaps))
		if len(m.gaps) > 0 {
			m.screen = screenGap
			m.textArea.Reset()
			m.textArea.SetHeight(4)
			m.textArea.Placeholder = "Digite sua resposta..."
			m.textArea.Focus()
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
			m.textArea.Reset()
			m.textArea.SetHeight(4)
			m.textArea.Placeholder = "Digite sua resposta..."
			m.textArea.SetValue(m.gapAnswers[m.gapIndex])
			m.textArea.Focus()
		} else {
			m.screen = screenNarrative
			m.textArea.SetHeight(8)
			m.textArea.Placeholder = "Descreva a tarefa livremente...\n(Ctrl+S para confirmar)"
			m.textArea.SetValue(m.narrative)
			m.textArea.Focus()
		}
		return m, nil
	case "ctrl+s":
		m.gapAnswers[m.gapIndex] = m.textArea.Value()
		m.gapIndex++
		if m.gapIndex >= len(m.gaps) {
			m.screen = screenAskPhase
		} else {
			m.textArea.Reset()
			m.textArea.SetHeight(4)
			m.textArea.Placeholder = "Digite sua resposta..."
			m.textArea.Focus()
		}
		return m, nil
	}
	var cmd tea.Cmd
	m.textArea, cmd = m.textArea.Update(msg)
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

