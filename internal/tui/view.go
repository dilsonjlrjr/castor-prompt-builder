package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m AppModel) View() string {
	switch m.screen {
	case screenSelectModel:
		return m.viewSelectModel()
	case screenModelInfo:
		return m.viewModelInfo()
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


func renderCastor() string {
	return lipgloss.JoinHorizontal(
		lipgloss.Center,
		renderMascote(),
		"    ",
		bigTitle(),
	)
}

func badge(txt string) string {
	return styleBadge.Render(" " + txt + " ")
}

// siglas mapeia o ID do modelo para o significado do acrônimo
var siglas = map[string]string{
	"rtf":    "Role · Task · Format",
	"race":   "Role · Action · Context · Expectation",
	"risen":  "Role · Input · Steps · Expectation · Narrowing",
	"create": "Context · Role · Examples · Audience · Tone · Expectation",
}

const propositoCastor = "Construa prompts estruturados para LLMs em segundos.\nEscolha um framework, descreva sua tarefa e o CASTOR monta o prompt ideal."

// --- Selecionar Modelo ---

func (m AppModel) viewSelectModel() string {
	var sb strings.Builder
	sb.WriteString(renderCastor() + "\n\n")

	// propósito
	sb.WriteString(lipgloss.NewStyle().
		Foreground(colorText).
		PaddingLeft(1).
		Render(propositoCastor) + "\n\n")

	sb.WriteString(styleSubtitle.Render("Selecione o modelo de prompt:") + "\n\n")

	for i, mod := range m.models {
		cursor := "  "
		style := styleNormal
		if i == m.selectedModel {
			cursor = styleSelected.Render("> ")
			style = styleSelected
		}
		// nome do modelo
		line := cursor + style.Render(mod.Nome)
		// significado do acrônimo
		if sig, ok := siglas[mod.ID]; ok {
			line += styleMuted.Render("  (" + sig + ")")
		}
		sb.WriteString(line + "\n")
		// descrição em linha separada, indentada
		if mod.Descricao != "" {
			sb.WriteString(styleMuted.Render("     "+mod.Descricao) + "\n")
		}
	}

	sb.WriteString("\n" + styleHelp.Render("↑↓ navegar   Enter selecionar   i mais informações   q sair"))
	return sb.String()
}

// --- Info do Modelo ---

var modelInfos = map[string]string{
	"rtf": `
█▀█ ▀█▀ █▀▀   Role · Task · Format
█▀▄  █  █▀    Tarefas simples e diretas

QUANDO USAR
  Ideal quando a tarefa é clara, o contexto é mínimo
  e você sabe exatamente o que quer.

CAMPOS
  • Papel       — quem vai executar a tarefa
  • Tarefa      — o que precisa ser feito (textarea)
  • Formato     — como a resposta deve ser apresentada
                  lista, tabela, código, texto corrido

PONTOS FORTES
  ✓ Rápido de preencher
  ✓ Sem ambiguidade
  ✓ Ótimo para tarefas repetitivas

EXEMPLOS DE USO
  → Resumir um artigo técnico em 5 pontos
  → Traduzir um texto para inglês formal
  → Listar os prós e contras de uma tecnologia
  → Responder uma dúvida específica de código

DICA
  Use RTF quando você já sabe o que quer.
  Se precisar dar contexto ou ter entregável complexo,
  prefira RACE ou RISEN.
`,

	"race": `
█▀█ █▀█ █▀▀ █▀▀   Role · Action · Context · Expectation
█▀▄ █▀█ █▄▄ ██▄   Contexto rico com entregável claro

QUANDO USAR
  Ideal quando a tarefa exige contexto de negócio,
  restrições ou expectativa de entregável específico.

CAMPOS
  • Papel        — especialista que vai executar
  • Ação         — o que deve ser feito
  • Contexto     — cenário, dados, histórico relevante
  • Tom          — formal, informal, técnico, persuasivo
  • Canais       — blog, LinkedIn, email, Instagram
  • Fases        — divisão da execução em etapas
  • Expectativa  — o que se espera como entrega final

PONTOS FORTES
  ✓ Captura o "por quê" da tarefa
  ✓ Permite dar muito contexto de forma estruturada
  ✓ Suporte a fases para entregas progressivas

EXEMPLOS DE USO
  → Plano editorial de 3 meses para startup B2B
  → Estratégia de go-to-market para novo produto
  → Análise de concorrência com foco em pricing
  → Campanha de email para reengajamento de leads

DICA
  Quanto mais contexto você der, melhor o resultado.
  Use o campo Contexto para dados reais:
  números, métricas, público, histórico da empresa.
`,

	"risen": `
█▀█ █ █▀ █▀▀ █▄  █   Role · Input · Steps · Expectation · Narrowing
█▀▄ █ ▄█ ██▄ █ ▀▄█   Steps detalhados com restrições

QUANDO USAR
  Ideal quando você quer controlar o raciocínio
  passo a passo e impor limites à resposta.

CAMPOS
  • Papel        — quem executa
  • Input        — dados/contexto de entrada
  • Steps        — etapas explícitas de execução
  • Expectativa  — entregável final esperado
  • Restrições   — o que NÃO deve ser feito ou limitações

PONTOS FORTES
  ✓ Máximo controle sobre o processo de raciocínio
  ✓ Restrições evitam respostas genéricas
  ✓ Steps garantem estrutura na resposta

EXEMPLOS DE USO
  → Auditar código em 4 etapas sem sugerir libs externas
  → Analisar dados financeiros com restrições de privacidade
  → Criar arquitetura de sistema com steps de decisão
  → Revisar proposta comercial considerando só o mercado BR

DICA
  As Restrições são o diferencial do RISEN.
  Seja específico: "não use jargão técnico",
  "limite a 500 palavras", "considere apenas dados de 2024".
`,

	"create": `
█▀▀ █▀█ █▀▀ █▀█ ▀█▀ █▀▀   Context · Role · Examples · Audience · Tone · Expectation
█▄▄ █▀▄ ██▄ █▀█  █  ██▄   Conteúdo criativo com público e tom definidos

QUANDO USAR
  Ideal para marketing, copywriting e conteúdo criativo
  onde público-alvo e tom são fatores críticos.

CAMPOS
  • Contexto     — cenário e motivação
  • Papel        — especialista criativo
  • Exemplos     — referências de conteúdo que você admira
  • Público      — quem vai consumir o conteúdo
  • Tom          — formal, informal, técnico, persuasivo, inspirador
  • Expectativa  — formato e entregável esperado

PONTOS FORTES
  ✓ Exemplos guiam o estilo sem ser restritivo
  ✓ Público bem definido = resposta mais relevante
  ✓ Ideal para times de marketing e conteúdo

EXEMPLOS DE USO
  → 5 headlines para campanha direcionada a CTOs
  → Roteiro de vídeo educativo para desenvolvedores júnior
  → Posts de LinkedIn para CEO de startup de healthtech
  → Landing page de produto SaaS tom persuasivo e técnico

DICA
  O campo Exemplos é o mais poderoso do CREATE.
  Cole 2-3 trechos de conteúdo que você considera
  referência — o modelo vai absorver o estilo.
`,
}

func (m AppModel) viewModelInfo() string {
	model := m.models[m.selectedModel]
	info, ok := modelInfos[model.ID]
	if !ok {
		info = "\nNenhuma informação adicional disponível para este modelo."
	}

	var sb strings.Builder
	sb.WriteString(styleHeader.Render(" "+model.Nome+" ") + "  " + badge(model.ID) + "\n")
	sb.WriteString(styleMuted.Render(model.Descricao) + "\n")
	sb.WriteString(strings.Repeat("─", 60) + "\n")
	sb.WriteString(styleNormal.Render(info) + "\n")
	sb.WriteString(styleHelp.Render("Esc voltar   Enter selecionar este modelo"))
	return sb.String()
}

// --- Selecionar Papel ---

func (m AppModel) viewSelectRole() string {
	model := m.models[m.selectedModel]
	var sb strings.Builder
	sb.WriteString(styleHeader.Render(" CASTOR BUILDER ") + "  " + badge("modelo: "+model.Nome) + "\n\n")
	sb.WriteString(styleSubtitle.Render("Selecione o papel:") + "\n\n")

	for i, role := range m.roles {
		cursor := "  "
		style := styleNormal
		if i == m.selectedRole {
			cursor = styleSelected.Render("> ")
			style = styleSelected
		}
		sb.WriteString(cursor + style.Render(role.Nome) + "\n")
	}

	sb.WriteString("\n" + styleHelp.Render("↑↓ navegar   Enter selecionar   Esc voltar"))
	return sb.String()
}

// --- Narrativa ---

func (m AppModel) viewNarrative() string {
	role := m.roles[m.selectedRole]
	var sb strings.Builder
	sb.WriteString(styleHeader.Render(" CASTOR BUILDER ") + "  " + badge("papel: "+role.Nome) + "\n\n")
	sb.WriteString(styleSubtitle.Render("Descreva a tarefa livremente:") + "\n\n")
	sb.WriteString(styleBorder.Render(m.textArea.View()) + "\n\n")
	sb.WriteString(styleHelp.Render("Ctrl+S confirmar   Esc voltar"))
	return sb.String()
}

// --- Gap ---

func (m AppModel) viewGap() string {
	total := len(m.gaps)
	current := m.gapIndex + 1
	var sb strings.Builder
	sb.WriteString(styleHeader.Render(" CASTOR BUILDER ") + "  " + badge(fmt.Sprintf("lacuna %d de %d", current, total)) + "\n\n")
	sb.WriteString(styleSubtitle.Render(m.gaps[m.gapIndex]) + "\n\n")
	sb.WriteString(styleBorder.Render(m.textInput.View()) + "\n\n")
	sb.WriteString(styleHelp.Render("Enter próximo   Esc voltar   (deixe vazio para pular)"))
	return sb.String()
}

// --- Fases ---

func (m AppModel) viewAskPhase() string {
	options := []string{"Sim, definir fases", "Não, gerar direto"}
	var sb strings.Builder
	sb.WriteString(styleHeader.Render(" CASTOR BUILDER ") + "\n\n")
	sb.WriteString(styleSubtitle.Render("Deseja definir fases de execução?") + "\n\n")

	for i, opt := range options {
		cursor := "  "
		style := styleNormal
		if i == m.askPhaseChoice {
			cursor = styleSelected.Render("> ")
			style = styleSelected
		}
		sb.WriteString(cursor + style.Render(opt) + "\n")
	}

	sb.WriteString("\n" + styleHelp.Render("↑↓ navegar   Enter confirmar   Esc voltar"))
	return sb.String()
}

func (m AppModel) viewDefinePhase() string {
	var sb strings.Builder

	if m.phaseIndex == -1 {
		sb.WriteString(styleHeader.Render(" CASTOR BUILDER ") + "  " + badge("fases") + "\n\n")
		sb.WriteString(styleSubtitle.Render("Quantas fases?") + "\n\n")
		sb.WriteString(styleBorder.Render(m.textInput.View()) + "\n\n")
		sb.WriteString(styleHelp.Render("Enter confirmar"))
		return sb.String()
	}

	phase := m.phaseIndex + 1
	sb.WriteString(styleHeader.Render(" CASTOR BUILDER ") + "  " + badge(fmt.Sprintf("fase %d de %d", phase, m.phaseCount)) + "\n\n")

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

// --- Concluído ---

func (m AppModel) viewDone() string {
	var sb strings.Builder
	sb.WriteString(styleHeader.Render(" CASTOR BUILDER ") + "\n\n")

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

// lipgloss usado via mascot.go — evita import não utilizado
var _ = lipgloss.Color("")
