package tui

import (
	"fmt"
	"strings"
)

func (m AppModel) View() string {
	switch m.screen {
	case screenSelectModel:
		return m.viewSelectModel()
	case screenSelectRole:
		return m.viewSelectRole()
	case screenNarrative:
		return m.viewNarrative()
	case screenGap:
		return m.viewGap()
	case screenAskPhase:
		return m.viewAskPhase()
	case screenDefinePhase:
		return m.viewDefinePhase()
	case screenDone:
		return m.viewDone()
	}
	return ""
}

func header(badge string) string {
	title := styleHeader.Render(" CASTOR BUILDER ")
	if badge != "" {
		title += "  " + styleBadge.Render(badge)
	}
	return title
}

func (m AppModel) viewSelectModel() string {
	var sb strings.Builder
	sb.WriteString(styleTitle.Render(castor) + "\n\n")
	sb.WriteString(header("") + "\n\n")
	sb.WriteString(styleSubtitle.Render("Selecione o modelo de prompt:") + "\n\n")

	for i, mod := range m.models {
		cursor := "  "
		style := styleNormal
		if i == m.selectedModel {
			cursor = "> "
			style = styleSelected
		}
		line := fmt.Sprintf("%s%s", cursor, mod.Nome)
		if mod.Descricao != "" {
			line += styleMuted.Render(fmt.Sprintf("  — %s", mod.Descricao))
		}
		sb.WriteString(style.Render(line) + "\n")
	}

	sb.WriteString("\n" + styleHelp.Render("↑↓ navegar   Enter selecionar   q sair"))
	return sb.String()
}

func (m AppModel) viewSelectRole() string {
	model := m.models[m.selectedModel]
	var sb strings.Builder
	sb.WriteString(header("modelo: "+model.Nome) + "\n\n")
	sb.WriteString(styleSubtitle.Render("Selecione o papel:") + "\n\n")

	for i, role := range m.roles {
		cursor := "  "
		style := styleNormal
		if i == m.selectedRole {
			cursor = "> "
			style = styleSelected
		}
		sb.WriteString(style.Render(fmt.Sprintf("%s%s", cursor, role.Nome)) + "\n")
	}

	sb.WriteString("\n" + styleHelp.Render("↑↓ navegar   Enter selecionar   Esc voltar"))
	return sb.String()
}

func (m AppModel) viewNarrative() string {
	role := m.roles[m.selectedRole]
	var sb strings.Builder
	sb.WriteString(header("papel: "+role.Nome) + "\n\n")
	sb.WriteString(styleSubtitle.Render("Descreva a tarefa livremente:") + "\n\n")
	sb.WriteString(styleBorder.Render(m.textArea.View()) + "\n\n")
	sb.WriteString(styleHelp.Render("Ctrl+S confirmar   Esc voltar"))
	return sb.String()
}

func (m AppModel) viewGap() string {
	total := len(m.gaps)
	current := m.gapIndex + 1
	var sb strings.Builder
	sb.WriteString(header(fmt.Sprintf("gaps %d de %d", current, total)) + "\n\n")
	sb.WriteString(styleSubtitle.Render(m.gaps[m.gapIndex]) + "\n\n")
	sb.WriteString(styleBorder.Render(m.textInput.View()) + "\n\n")
	sb.WriteString(styleHelp.Render("Enter próximo   Esc voltar   (vazio = pular)"))
	return sb.String()
}

func (m AppModel) viewAskPhase() string {
	options := []string{"Sim, definir fases", "Não, gerar direto"}
	var sb strings.Builder
	sb.WriteString(header("") + "\n\n")
	sb.WriteString(styleSubtitle.Render("Deseja definir fases de execução?") + "\n\n")

	for i, opt := range options {
		cursor := "  "
		style := styleNormal
		if i == m.askPhaseChoice {
			cursor = "> "
			style = styleSelected
		}
		sb.WriteString(style.Render(fmt.Sprintf("%s%s", cursor, opt)) + "\n")
	}

	sb.WriteString("\n" + styleHelp.Render("↑↓ navegar   Enter confirmar   Esc voltar"))
	return sb.String()
}

func (m AppModel) viewDefinePhase() string {
	var sb strings.Builder

	if m.phaseIndex == -1 {
		sb.WriteString(header("fases") + "\n\n")
		sb.WriteString(styleSubtitle.Render("Quantas fases?") + "\n\n")
		sb.WriteString(styleBorder.Render(m.textInput.View()) + "\n\n")
		sb.WriteString(styleHelp.Render("Enter confirmar"))
		return sb.String()
	}

	phase := m.phaseIndex + 1
	sb.WriteString(header(fmt.Sprintf("fase %d de %d", phase, m.phaseCount)) + "\n\n")

	if m.phaseEditField == 0 {
		sb.WriteString(styleSubtitle.Render("Título da fase:") + "\n\n")
		sb.WriteString(styleBorder.Render(m.textInput.View()) + "\n\n")
		sb.WriteString(styleHelp.Render("Enter próximo campo   Esc voltar"))
	} else {
		sb.WriteString(styleSubtitle.Render("Descrição da fase:") + "\n\n")
		sb.WriteString(styleBorder.Render(m.textArea.View()) + "\n\n")
		sb.WriteString(styleHelp.Render("Enter avançar   Esc voltar"))
	}
	return sb.String()
}

func (m AppModel) viewDone() string {
	var sb strings.Builder
	sb.WriteString(header("") + "\n\n")

	if m.err != nil {
		sb.WriteString(styleError.Render("✗ Erro ao salvar: "+m.err.Error()) + "\n")
	} else {
		sb.WriteString(styleSuccess.Render("✓ Prompt gerado com sucesso!") + "\n\n")
		sb.WriteString(styleSubtitle.Render("Salvo em:") + "\n")
		sb.WriteString(styleSelected.Render("  "+m.savedPath) + "\n")
	}

	sb.WriteString("\n" + styleHelp.Render("Enter / q para sair"))
	return sb.String()
}
