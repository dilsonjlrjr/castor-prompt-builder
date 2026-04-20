package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// maxHeight é o teto em linhas (~550px a ~16px/linha).
const maxHeight = 35

// effectiveHeight retorna a altura real a usar: mín(terminalHeight, maxHeight).
func (m AppModel) effectiveHeight() int {
	if m.height > 0 && m.height < maxHeight {
		return m.height
	}
	return maxHeight
}

// capHeight garante que o conteúdo tenha exatamente h linhas.
func capHeight(s string, h int) string {
	lines := strings.Split(s, "\n")
	if len(lines) > h {
		lines = lines[:h]
	}
	for len(lines) < h {
		lines = append(lines, "")
	}
	return strings.Join(lines, "\n")
}

func (m AppModel) View() string {
	h := m.effectiveHeight()
	switch m.screen {
	case screenSelectModel:
		return capHeight(m.viewSelectModel(), h)
	case screenModelInfo:
		return capHeight(m.viewModelInfo(), h)
	case screenSelectRole:
		return capHeight(m.viewSelectRole(), h)
	case screenNarrative:
		return capHeight(m.viewNarrative(), h)
	case screenGap:
		return capHeight(m.viewGap(), h)
	case screenAskPhase:
		return capHeight(m.viewAskPhase(), h)
	case screenDefinePhase:
		return capHeight(m.viewDefinePhase(), h)
	case screenDone:
		return capHeight(m.viewDone(), h)
	}
	return ""
}

// renderCastor adapta o cabeçalho ao espaço disponível.
// - largura >= mascote+título+4 : lado a lado (layout completo)
// - largura >= título            : só o título em pixel font
// - largura pequena              : texto simples centralizado
func renderCastor(availableWidth int) string {
	mascot := renderMascote()
	title := bigTitle()

	mascotW := 0
	for _, l := range strings.Split(mascot, "\n") {
		if w := lipgloss.Width(l); w > mascotW {
			mascotW = w
		}
	}
	titleW := 0
	for _, l := range strings.Split(title, "\n") {
		if w := lipgloss.Width(l); w > titleW {
			titleW = w
		}
	}

	needed := mascotW + 4 + titleW
	if availableWidth == 0 || availableWidth >= needed {
		return lipgloss.JoinHorizontal(lipgloss.Center, mascot, "    ", title)
	}
	if availableWidth >= titleW+2 {
		return title
	}
	// terminal muito estreito
	return lipgloss.NewStyle().
		Foreground(colorPrimary).
		Bold(true).
		Render("◆ CASTOR BUILDER ◆")
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
	sb.WriteString(renderCastor(m.width) + "\n\n")

	// propósito
	sb.WriteString(lipgloss.NewStyle().
		Foreground(colorText).
		PaddingLeft(1).
		Render(propositoCastor) + "\n\n")

	sb.WriteString(
		styleSubtitle.Render("Selecione o modelo de prompt:") +
			"  " + styleMuted.Render("(pressione i para mais informações)") + "\n\n")

	// largura máxima dos nomes para alinhamento tabular
	maxNome := 0
	for _, mod := range m.models {
		if len(mod.Nome) > maxNome {
			maxNome = len(mod.Nome)
		}
	}

	for i, mod := range m.models {
		cursor := "  "
		style := styleNormal
		if i == m.selectedModel {
			cursor = styleSelected.Render("> ")
			style = styleSelected
		}
		padding := strings.Repeat(" ", maxNome-len(mod.Nome)+2)
		line := cursor + style.Render(mod.Nome) + padding + styleMuted.Render(mod.Descricao)
		sb.WriteString(line + "\n")
	}

	sb.WriteString("\n" + styleHelp.Render("↑↓ navegar   Enter selecionar   q sair"))
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

	// área disponível: altura efetiva − cabeçalho(3) − rodapé(2)
	infoAreaH := m.effectiveHeight() - 5
	if infoAreaH < 5 {
		infoAreaH = 5
	}

	infoLines := strings.Split(styleNormal.Render(info), "\n")
	offset := m.modelInfoOffset
	if offset >= len(infoLines) {
		offset = len(infoLines) - 1
	}
	if offset < 0 {
		offset = 0
	}
	end := offset + infoAreaH
	if end > len(infoLines) {
		end = len(infoLines)
	}
	visibleInfo := strings.Join(infoLines[offset:end], "\n")

	var sb strings.Builder
	sb.WriteString(styleHeader.Render(" "+model.Nome+" ") + "  " + badge(model.ID) + "\n")
	sb.WriteString(styleMuted.Render(model.Descricao) + "\n")
	sb.WriteString(strings.Repeat("─", 60) + "\n")
	sb.WriteString(visibleInfo + "\n")
	if end < len(infoLines) {
		sb.WriteString(styleMuted.Render("  ↓ mais abaixo (j / ↓)") + "\n")
	}
	sb.WriteString(styleHelp.Render("↑↓ rolar   Esc voltar   Enter selecionar este modelo"))
	return sb.String()
}

// nomeCategoria mapeia o id do diretório para nome de exibição
var nomeCategoria = map[string]string{
	"arquitetura": "Arquitetura",
	"frontend":    "Frontend & Mobile",
	"backend":     "Backend",
	"devops":      "DevOps & Cloud",
	"banco":       "Banco de Dados",
	"dados":       "Dados & IA",
	"gestao":      "Gestão",
	"seguranca":   "Segurança",
	"design":      "Design",
	"marketing":   "Marketing",
}

// --- Selecionar Papel ---

func (m AppModel) viewSelectRole() string {
	model := m.models[m.selectedModel]
	filtered := m.filteredRoleIndices()

	selectedCount := 0
	for _, sel := range m.selectedRoles {
		if sel {
			selectedCount++
		}
	}
	titulo := "Selecione o(s) papel(eis):"
	if selectedCount > 0 {
		titulo += fmt.Sprintf("  %s", badge(fmt.Sprintf("%d selecionado(s)", selectedCount)))
	}

	var sb strings.Builder
	sb.WriteString(styleHeader.Render(" CASTOR BUILDER ") + "  " + badge("modelo: "+model.Nome) + "\n\n")
	sb.WriteString(styleSubtitle.Render(titulo) + "\n\n")
	sb.WriteString(styleBorder.Render(m.textInput.View()) + "\n\n")

	if len(filtered) == 0 {
		sb.WriteString(styleMuted.Render("  Nenhum papel encontrado.") + "\n")
		sb.WriteString("\n" + styleHelp.Render("↑↓ navegar   Espaço marcar/desmarcar   Enter confirmar   Esc limpar busca / voltar   (digite para buscar)"))
		return sb.String()
	}

	// área disponível para a lista: altura − cabeçalho(8) − rodapé(2) − indicadores(2)
	listH := m.effectiveHeight() - 12
	if listH < 4 {
		listH = 4
	}

	// monta todas as entradas (cabeçalhos de categoria + papéis)
	type entrada struct {
		texto    string
		isCursor bool
	}
	var entradas []entrada
	lastCat := ""
	cursorLinha := 0

	for listIdx, globalIdx := range filtered {
		role := m.roles[globalIdx]
		cat := role.Categoria
		if catNome, ok := nomeCategoria[cat]; ok {
			cat = catNome
		}
		if cat != lastCat {
			if len(entradas) > 0 {
				entradas = append(entradas, entrada{texto: ""})
			}
			entradas = append(entradas, entrada{
				texto: styleSelected.Render("  ── " + strings.ToUpper(cat) + " ──"),
			})
			lastCat = cat
		}

		cursor := "  "
		check := "[ ]"
		style := styleNormal
		if listIdx == m.roleCursor {
			cursor = "> "
			style = styleSelected
			cursorLinha = len(entradas)
		}
		if m.selectedRoles[globalIdx] {
			check = "[✓]"
		}
		entradas = append(entradas, entrada{
			texto:    cursor + style.Render(check+" "+role.Nome),
			isCursor: listIdx == m.roleCursor,
		})
	}

	// janela centrada no cursor
	half := listH / 2
	start := cursorLinha - half
	if start < 0 {
		start = 0
	}
	end := start + listH
	if end > len(entradas) {
		end = len(entradas)
		if start = end - listH; start < 0 {
			start = 0
		}
	}

	if start > 0 {
		sb.WriteString(styleMuted.Render("  ↑ mais acima") + "\n")
	}
	for _, e := range entradas[start:end] {
		sb.WriteString(e.texto + "\n")
	}
	if end < len(entradas) {
		sb.WriteString(styleMuted.Render("  ↓ mais abaixo") + "\n")
	}

	sb.WriteString("\n" + styleHelp.Render("↑↓ navegar   Espaço marcar/desmarcar   Enter confirmar   Esc limpar busca / voltar   (digite para buscar)"))
	return sb.String()
}

// --- Narrativa ---

func (m AppModel) selectedRoleNames() string {
	var nomes []string
	for idx, sel := range m.selectedRoles {
		if sel {
			nomes = append(nomes, m.roles[idx].Nome)
		}
	}
	if len(nomes) == 0 {
		return "papel"
	}
	return strings.Join(nomes, " + ")
}

func (m AppModel) viewNarrative() string {
	var sb strings.Builder
	sb.WriteString(styleHeader.Render(" CASTOR BUILDER ") + "  " + badge(m.selectedRoleNames()) + "\n\n")
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
