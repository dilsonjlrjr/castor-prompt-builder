package tui

import (
	"sort"
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
		case screenContext:
			return m.updateContext(msg)
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

	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)
	newSearch := m.textInput.Value()
	if newSearch != m.roleSearch {
		m.roleSearch = newSearch
		m.roleCursor = 0
	}
	return m, cmd
}

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

		autoMapped := map[string]bool{
			"role": true, "action": true, "context": true,
			"task": true, "input": true,
		}

		var sections []ContextSection
		model := m.models[m.selectedModel]

		// 1. Todos os campos do modelo (obrigatorios + opcionais) em uma secao
		var modelGaps []*ContextGap
		for _, c := range model.Campos {
			if autoMapped[c.ID] || c.Tipo == parser.FieldSteps {
				continue
			}
			modelGaps = append(modelGaps, &ContextGap{
				FieldID:     c.ID,
				Pergunta:    c.Label,
				Obrigatorio: c.Obrigatorio,
				Tipo:        c.Tipo,
				Opcoes:      c.Opcoes,
			})
		}
		// obrigatorios primeiro, depois respondidos, depois opcionais vazios
		sort.SliceStable(modelGaps, func(i, j int) bool {
			ai, aj := modelGaps[i], modelGaps[j]
			// required empty = 0, answered = 1, optional empty = 2
			scoreI := 2
			if ai.Obrigatorio && ai.Answer == "" { scoreI = 0 }
			if ai.Answer != "" { scoreI = 1 }
			scoreJ := 2
			if aj.Obrigatorio && aj.Answer == "" { scoreJ = 0 }
			if aj.Answer != "" { scoreJ = 1 }
			return scoreI < scoreJ
		})
		if len(modelGaps) > 0 {
			sections = append(sections, ContextSection{
				Title: "Contexto do Modelo (" + model.Nome + ")",
				Kind:  "model",
				Gaps:  modelGaps,
			})
		}

		// 2. gaps_comuns dos papeis selecionados — agrupados por papel
		for idx, sel := range m.selectedRoles {
			if !sel {
				continue
			}
			r := m.roles[idx]
			if len(r.GapsComuns) == 0 {
				continue
			}
			var roleGaps []*ContextGap
			for _, q := range r.GapsComuns {
				roleGaps = append(roleGaps, &ContextGap{
					FieldID:  "",
					Pergunta: q,
					RoleNome: r.Nome,
					Tipo:     parser.FieldTextarea,
				})
			}
			sections = append(sections, ContextSection{
				Title: r.Nome,
				Kind:  "role",
				Gaps:  roleGaps,
			})
		}

		m.contextSections = sections
		m.contextSecIdx = 0
		m.contextCursor = 0
		m.contextEditing = false
		m.contextMultiIdx = 0
		m.contextRoleFilter = -1
		if len(m.contextSections) > 0 {
			m.screen = screenContext
			m.textArea.Reset()
			m.textArea.SetHeight(4)
			m.textArea.Placeholder = "Digite sua resposta..."
		} else {
			m.screen = screenAskPhase
		}
		return m, nil
	}
	var cmd tea.Cmd
	m.textArea, cmd = m.textArea.Update(msg)
	return m, cmd
}

func (m AppModel) updateContext(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.contextEditing {
		return m.updateContextEditing(msg)
	}

	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit
	case "esc":
		m.screen = screenNarrative
		m.textArea.SetHeight(8)
		m.textArea.Placeholder = "Descreva a tarefa livremente...\n(Ctrl+S para confirmar)"
		m.textArea.SetValue(m.narrative)
		m.textArea.Focus()
		return m, nil

	case "ctrl+s":
		m.screen = screenAskPhase
		return m, nil

	case "up", "k":
		if m.contextCursor > 0 {
			m.contextCursor--
		}
	case "down", "j":
		sec := m.activeSection()
		if sec != nil && m.contextCursor < len(sec.Gaps)-1 {
			m.contextCursor++
		}
	case "left", "h":
		m.moveSection(-1)
	case "right", "l":
		m.moveSection(1)
	case "tab":
		m.cycleRoleFilter()

	case "enter":
		gap := m.currentContextGap()
		if gap == nil {
			return m, nil
		}
		if gap.Tipo == parser.FieldSelect && len(gap.Opcoes) > 0 {
			idx := -1
			for i, o := range gap.Opcoes {
				if o == gap.Answer {
					idx = i
					break
				}
			}
			idx = (idx + 1) % len(gap.Opcoes)
			gap.Answer = gap.Opcoes[idx]
			return m, nil
		}
		if gap.Tipo == parser.FieldMultiselect && len(gap.Opcoes) > 0 {
			op := gap.Opcoes[m.contextMultiIdx%len(gap.Opcoes)]
			m.contextMultiIdx++
			parts := strings.Split(gap.Answer, ", ")
			var newParts []string
			found := false
			for _, p := range parts {
				p = strings.TrimSpace(p)
				if p == op {
					found = true
					continue
				}
				if p != "" {
					newParts = append(newParts, p)
				}
			}
			if !found {
				newParts = append(newParts, op)
			}
			gap.Answer = strings.Join(newParts, ", ")
			return m, nil
		}
		m.contextEditing = true
		m.textArea.Reset()
		m.textArea.SetHeight(4)
		m.textArea.Placeholder = "Digite sua resposta..."
		m.textArea.SetValue(gap.Answer)
		m.textArea.Focus()
		return m, nil
	}
	return m, nil
}

func (m AppModel) moveSection(delta int) {
	roleSecs := m.roleSectionIndices()
	if m.contextRoleFilter >= 0 {
		if len(roleSecs) == 0 {
			return
		}
		cur := m.contextRoleFilter
		cur += delta
		if cur < 0 {
			cur = len(roleSecs) - 1
		}
		if cur >= len(roleSecs) {
			cur = 0
		}
		m.contextRoleFilter = cur
		m.contextCursor = 0
		return
	}
	// filtro "todos": navega entre todas as secoes
	all := m.visibleSections()
	if len(all) == 0 {
		return
	}
	cur := m.contextSecIdx + delta
	if cur < 0 {
		cur = len(all) - 1
	}
	if cur >= len(all) {
		cur = 0
	}
	newSec := all[cur]
	// acha o indice real em contextSections
	for i, sec := range m.contextSections {
		if sec.Title == newSec.Title && sec.Kind == newSec.Kind {
			m.contextSecIdx = i
			break
		}
	}
	if len(newSec.Gaps) > 0 && m.contextCursor >= len(newSec.Gaps) {
		m.contextCursor = len(newSec.Gaps) - 1
	}
	if m.contextCursor < 0 {
		m.contextCursor = 0
	}
}

func (m AppModel) cycleRoleFilter() {
	roleSecs := m.roleSectionIndices()
	if len(roleSecs) == 0 {
		return
	}
	m.contextRoleFilter++
	if m.contextRoleFilter >= len(roleSecs) {
		m.contextRoleFilter = -1
	}
	m.contextSecIdx = 0
	m.contextCursor = 0
}

// roleSectionIndices retorna indices em contextSections das secoes kind="role"
func (m AppModel) roleSectionIndices() []int {
	var idxs []int
	for i, sec := range m.contextSections {
		if sec.Kind == "role" {
			idxs = append(idxs, i)
		}
	}
	return idxs
}

// visibleSections retorna secoes visiveis considerando o filtro de papel
func (m AppModel) visibleSections() []ContextSection {
	if m.contextRoleFilter < 0 {
		return m.contextSections
	}
	roleSecs := m.roleSectionIndices()
	if m.contextRoleFilter >= len(roleSecs) {
		return m.contextSections
	}
	visible := make([]ContextSection, 0)
	for i, sec := range m.contextSections {
		if sec.Kind == "model" {
			visible = append(visible, sec)
		} else if i == roleSecs[m.contextRoleFilter] {
			visible = append(visible, sec)
		}
	}
	return visible
}

// activeSection retorna a secao ativa (considerando filtro)
func (m AppModel) activeSection() *ContextSection {
	visible := m.visibleSections()
	if m.contextSecIdx >= len(visible) {
		return nil
	}
	return &visible[m.contextSecIdx]
}

func (m AppModel) updateContextEditing(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit
	case "esc":
		m.contextEditing = false
		return m, nil
	case "ctrl+s":
		gap := m.currentContextGap()
		if gap != nil {
			gap.Answer = m.textArea.Value()
		}
		m.contextEditing = false
		return m, nil
	}
	var cmd tea.Cmd
	m.textArea, cmd = m.textArea.Update(msg)
	return m, cmd
}

func (m AppModel) currentContextGap() *ContextGap {
	sec := m.activeSection()
	if sec == nil {
		return nil
	}
	if m.contextCursor < len(sec.Gaps) {
		return sec.Gaps[m.contextCursor]
	}
	return nil
}

func (m AppModel) updateAskPhase(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit
	case "esc":
		if len(m.contextSections) > 0 {
			m.screen = screenContext
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
		if m.askPhaseChoice == 0 {
			m.textInput.Reset()
			m.textInput.Placeholder = "Quantidade de fases (ex: 3)"
			m.textInput.Focus()
			m.screen = screenDefinePhase
			m.phaseIndex = -1
		} else {
			m = m.buildAndSave()
		}
	}
	return m, nil
}

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
			n, err := strconv.Atoi(strings.TrimSpace(m.textInput.Value()))
			if err != nil || n < 1 {
				return m, nil
			}
			m.phaseCount = n
			m.phaseIndex = 0
			m.phaseEditField = 0
			steps := make([]parser.Step, n)
			m.values.Steps["fases"] = steps
			m.textInput.Reset()
			m.textInput.Placeholder = "Título da fase"
			return m, nil
		}
		if m.phaseEditField == 0 {
			m.values.Steps["fases"][m.phaseIndex].Titulo = m.textInput.Value()
			m.phaseEditField = 1
			m.textArea.Reset()
			m.textArea.Focus()
			return m, nil
		}
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

func (m AppModel) updateDone(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q", "enter":
		return m, tea.Quit
	}
	return m, nil
}
