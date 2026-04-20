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
	var nomes, descs, gaps []string
	roleID := "papel"
	for idx := range m.roles {
		if !m.selectedRoles[idx] {
			continue
		}
		r := m.roles[idx]
		nomes = append(nomes, r.Nome)
		descs = append(descs, r.Descricao)
		gaps = append(gaps, r.GapsComuns...)
		if roleID == "papel" {
			roleID = r.ID
		}
	}
	roleNome := strings.Join(nomes, " + ")
	roleDesc := strings.Join(descs, "\n\n")
	if roleNome == "" {
		// fallback seguro
		roleNome = "Especialista"
		roleDesc = ""
	}
	m.gaps = unique(gaps)

	// preenche values
	m.values.Fields["role"] = roleNome + ". " + roleDesc
	m.values.Fields["action"] = m.narrative
	m.values.Fields["context"] = m.narrative
	m.values.Fields["task"] = m.narrative

	// gap answers
	for i, ans := range m.gapAnswers {
		if i < len(m.gaps) && ans != "" {
			key := fmt.Sprintf("gap_%d", i)
			m.values.Fields[key] = ans
		}
	}

	rendered := engine.Render(model.Template, m.values)

	// monta gaps não respondidos
	var unanswered []string
	for i, ans := range m.gapAnswers {
		if strings.TrimSpace(ans) == "" && i < len(m.gaps) {
			unanswered = append(unanswered, m.gaps[i])
		}
	}

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
