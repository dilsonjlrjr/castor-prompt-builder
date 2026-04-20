package engine

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/dilsonrabelo/castor-prompt-builder/pkg/parser"
)

// Values contém todos os valores preenchidos pelo usuário
type Values struct {
	Fields map[string]string          // text, textarea, select
	Lists  map[string][]string        // multiselect, list
	Steps  map[string][]parser.Step   // steps
}

func NewValues() *Values {
	return &Values{
		Fields: make(map[string]string),
		Lists:  make(map[string][]string),
		Steps:  make(map[string][]parser.Step),
	}
}

var (
	reSimple = regexp.MustCompile(`\{\{(\w+)\}\}`)
	reIf     = regexp.MustCompile(`(?s)\{\{#if (\w+)\}\}(.*?)\{\{/if\}\}`)
	reSteps  = regexp.MustCompile(`(?s)\{\{#steps (\w+)\}\}(.*?)\{\{/steps\}\}`)
	reEach   = regexp.MustCompile(`(?s)\{\{#each (\w+)\}\}(.*?)\{\{/each\}\}`)
	reCurrent = regexp.MustCompile(`\{\{\.\}\}`)
)

// Render renderiza o template com os valores fornecidos
func Render(template string, vals *Values) string {
	// {{#steps campo}}...{{/steps}}
	result := reSteps.ReplaceAllStringFunc(template, func(match string) string {
		sub := reSteps.FindStringSubmatch(match)
		field, block := sub[1], sub[2]
		steps, ok := vals.Steps[field]
		if !ok || len(steps) == 0 {
			return ""
		}
		var sb strings.Builder
		for i, step := range steps {
			chunk := block
			chunk = strings.ReplaceAll(chunk, "{{titulo}}", step.Titulo)
			chunk = strings.ReplaceAll(chunk, "{{descricao}}", step.Descricao)
			chunk = strings.ReplaceAll(chunk, "{{index}}", fmt.Sprintf("%d", i+1))
			sb.WriteString(strings.TrimSpace(chunk))
			sb.WriteString("\n\n")
		}
		return sb.String()
	})

	// {{#each campo}}...{{/each}}
	result = reEach.ReplaceAllStringFunc(result, func(match string) string {
		sub := reEach.FindStringSubmatch(match)
		field, block := sub[1], sub[2]
		items, ok := vals.Lists[field]
		if !ok || len(items) == 0 {
			return ""
		}
		var sb strings.Builder
		for i, item := range items {
			chunk := reCurrent.ReplaceAllString(block, item)
			sb.WriteString(chunk)
			if i < len(items)-1 {
				sb.WriteString(", ")
			}
		}
		return sb.String()
	})

	// {{#if campo}}...{{/if}}
	result = reIf.ReplaceAllStringFunc(result, func(match string) string {
		sub := reIf.FindStringSubmatch(match)
		field, block := sub[1], sub[2]
		if v, ok := vals.Fields[field]; ok && v != "" {
			return strings.TrimSpace(block)
		}
		if items, ok := vals.Lists[field]; ok && len(items) > 0 {
			return strings.TrimSpace(block)
		}
		if steps, ok := vals.Steps[field]; ok && len(steps) > 0 {
			return strings.TrimSpace(block)
		}
		return ""
	})

	// {{campo}}
	result = reSimple.ReplaceAllStringFunc(result, func(match string) string {
		sub := reSimple.FindStringSubmatch(match)
		field := sub[1]
		if v, ok := vals.Fields[field]; ok {
			return v
		}
		// campo não preenchido: retorna vazio (nunca deixa {{placeholder}} no output)
		return ""
	})

	// limpa linhas duplas extras
	reBlank := regexp.MustCompile(`\n{3,}`)
	result = reBlank.ReplaceAllString(result, "\n\n")

	return strings.TrimSpace(result)
}
