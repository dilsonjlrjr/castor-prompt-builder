package tui

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
	"unicode"

	"github.com/dilsonrabelo/castor-prompt-builder/pkg/engine"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

func slugify(s string) string {
	b := make([]byte, len(s)*4)
	n, _ := transform.NewReader(strings.NewReader(s), norm.NFD).Read(b)
	var buf bytes.Buffer
	for _, r := range string(b[:n]) {
		if !unicode.Is(unicode.Mn, r) {
			buf.WriteRune(r)
		}
	}
	result := strings.ToLower(buf.String())
	re := regexp.MustCompile(`[^a-z0-9]+`)
	result = re.ReplaceAllString(result, "-")
	result = strings.Trim(result, "-")
	words := strings.Split(result, "-")
	if len(words) > 6 {
		words = words[:6]
	}
	return strings.Join(words, "-")
}

func (m AppModel) buildAndSave() AppModel {
	model := m.models[m.selectedModel]

	// combina múltiplos papéis selecionados
	var nomes, descs []string
	roleID := "papel"
	for idx := range m.roles {
		if !m.selectedRoles[idx] {
			continue
		}
		r := m.roles[idx]
		nomes = append(nomes, r.Nome)
		descs = append(descs, r.Descricao)
		if roleID == "papel" {
			roleID = r.ID
		}
	}
	roleNome := strings.Join(nomes, " + ")
	roleDesc := strings.Join(descs, "\n\n")
	if roleNome == "" {
		roleNome = "Especialista"
		roleDesc = ""
	}

	// ── valores auto-mapeados ────────────────────────────────────────────────
	m.values.Fields["role"]    = roleNome + ". " + roleDesc
	m.values.Fields["action"]  = m.narrative
	m.values.Fields["context"] = m.narrative
	m.values.Fields["task"]    = m.narrative
	m.values.Fields["input"]   = m.narrative

	// ── gap answers: campos do modelo → campo real; gaps de papel → seção extra depois ──
	for i, ga := range m.gaps {
		if i >= len(m.gapAnswers) {
			break
		}
		ans := strings.TrimSpace(m.gapAnswers[i])
		if ans == "" || ga.FieldID == "" {
			continue // role gaps são injetados em seção extra
		}
		if ga.Tipo == "list" || ga.Tipo == "multiselect" {
			var items []string
			for _, line := range strings.Split(ans, "\n") {
				if t := strings.TrimSpace(line); t != "" {
					items = append(items, t)
				}
			}
			if len(items) > 0 {
				m.values.Lists[ga.FieldID] = items
			}
		} else {
			m.values.Fields[ga.FieldID] = ans
		}
	}

	rendered := engine.Render(model.Template, m.values)

	// ── seções extras (após o template) ─────────────────────────────────────
	var extras strings.Builder

	// fases: fallback para modelos sem {{#steps}}
	if len(m.values.Steps["fases"]) > 0 && !strings.Contains(model.Template, "{{#steps") {
		extras.WriteString("\n\n---\n## Fases de execução\n\n")
		for i, s := range m.values.Steps["fases"] {
			extras.WriteString(fmt.Sprintf("### Fase %d — %s\n%s\n\n", i+1, s.Titulo, s.Descricao))
		}
	}

	// habilidades dos papéis (dedup)
	seenH := map[string]bool{}
	var habs []string
	for idx := range m.roles {
		if !m.selectedRoles[idx] {
			continue
		}
		for _, h := range m.roles[idx].Habilidades {
			if !seenH[h] {
				seenH[h] = true
				habs = append(habs, h)
			}
		}
	}
	if len(habs) > 0 {
		extras.WriteString("\n\n---\n## Habilidades relevantes\n")
		for _, h := range habs {
			extras.WriteString("- " + h + "\n")
		}
	}

	// tom dos papéis
	var toms []string
	for idx := range m.roles {
		if !m.selectedRoles[idx] {
			continue
		}
		r := m.roles[idx]
		if r.Tom != "" {
			toms = append(toms, r.Nome+": "+r.Tom)
		}
	}
	if len(toms) > 0 {
		extras.WriteString("\n\n---\n## Tom de comunicação\n")
		for _, t := range toms {
			extras.WriteString("- " + t + "\n")
		}
	}

	// contexto dos papéis: gaps_comuns respondidos
	var gapCtx []string
	for i, ga := range m.gaps {
		if i >= len(m.gapAnswers) {
			break
		}
		ans := strings.TrimSpace(m.gapAnswers[i])
		if ga.FieldID == "" && ans != "" {
			label := ga.Pergunta
			if ga.RoleNome != "" {
				label = ga.RoleNome + " — " + ga.Pergunta
			}
			gapCtx = append(gapCtx, "**"+label+"**\n"+ans)
		}
	}
	if len(gapCtx) > 0 {
		extras.WriteString("\n\n---\n## Contexto dos papéis\n\n")
		for _, g := range gapCtx {
			extras.WriteString(g + "\n\n")
		}
	}

	if extras.Len() > 0 {
		rendered += extras.String()
	}

	// ── gaps obrigatórios não respondidos → aviso no rodapé ─────────────────
	var unanswered []string
	for i, ga := range m.gaps {
		if i >= len(m.gapAnswers) {
			break
		}
		if ga.FieldID != "" && ga.Obrigatorio && strings.TrimSpace(m.gapAnswers[i]) == "" {
			unanswered = append(unanswered, ga.Pergunta)
		}
	}

	// ── salvar arquivo ───────────────────────────────────────────────────────
	date := time.Now().Format("20060102")
	slug := slugify(m.narrative)
	filename := fmt.Sprintf("%s_%s_%s.md", date, roleID, slug)

	if err := os.MkdirAll("prompts", 0755); err != nil {
		m.err = err
		m.screen = screenDone
		return m
	}

	path := filepath.Join("prompts", filename)

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("# Prompt — %s\n", roleNome))
	sb.WriteString(fmt.Sprintf("_Modelo: %s | Gerado em: %s_\n\n", model.Nome, time.Now().Format("2006-01-02")))
	sb.WriteString("---\n\n")
	sb.WriteString(rendered)

	if len(unanswered) > 0 {
		sb.WriteString("\n\n---\n")
		sb.WriteString("> ⚠️ Para refinar este prompt considere informar:\n")
		for _, q := range unanswered {
			sb.WriteString(fmt.Sprintf("> - %s\n", q))
		}
	}

	if err := os.WriteFile(path, []byte(sb.String()), 0644); err != nil {
		m.err = err
	} else {
		m.savedPath = path
	}

	m.screen = screenDone
	return m
}

// unique remove duplicatas preservando ordem
func unique(ss []string) []string {
	seen := make(map[string]bool)
	var out []string
	for _, s := range ss {
		if !seen[s] {
			seen[s] = true
			out = append(out, s)
		}
	}
	return out
}
