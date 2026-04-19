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

	"github.com/dilsonrabelo/castor-prompt-builder/internal/engine"
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
	role := m.roles[m.selectedRole]
	model := m.models[m.selectedModel]

	// preenche values com role e narrative
	m.values.Fields["role"] = role.Nome
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
	filename := fmt.Sprintf("%s_%s_%s.md", date, role.ID, slug)

	if err := os.MkdirAll("prompts", 0755); err != nil {
		m.err = err
		m.screen = screenDone
		return m
	}

	path := filepath.Join("prompts", filename)

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("# Prompt — %s\n", role.Nome))
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
