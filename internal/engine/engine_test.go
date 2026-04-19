package engine

import (
	"strings"
	"testing"

	"github.com/dilsonrabelo/castor-prompt-builder/internal/parser"
)

func TestSimpleField(t *testing.T) {
	vals := NewValues()
	vals.Fields["role"] = "Analista de Dados"
	result := Render("Você é {{role}}.", vals)
	if result != "Você é Analista de Dados." {
		t.Errorf("got: %q", result)
	}
}

func TestIfPresent(t *testing.T) {
	vals := NewValues()
	vals.Fields["tom"] = "formal"
	result := Render("{{#if tom}}Tom: {{tom}}{{/if}}", vals)
	if !strings.Contains(result, "Tom: formal") {
		t.Errorf("got: %q", result)
	}
}

func TestIfAbsent(t *testing.T) {
	vals := NewValues()
	result := Render("before{{#if tom}}HIDDEN{{/if}}after", vals)
	if strings.Contains(result, "HIDDEN") {
		t.Errorf("expected hidden, got: %q", result)
	}
}

func TestEach(t *testing.T) {
	vals := NewValues()
	vals.Lists["canais"] = []string{"blog", "LinkedIn"}
	result := Render("{{#each canais}}{{.}}{{/each}}", vals)
	if !strings.Contains(result, "blog") || !strings.Contains(result, "LinkedIn") {
		t.Errorf("got: %q", result)
	}
}

func TestSteps(t *testing.T) {
	vals := NewValues()
	vals.Steps["fases"] = []parser.Step{
		{Titulo: "Diagnóstico", Descricao: "Analise as causas"},
		{Titulo: "Plano", Descricao: "Estruture o plano"},
	}
	result := Render("{{#steps fases}}## {{titulo}}\n{{descricao}}{{/steps}}", vals)
	if !strings.Contains(result, "## Diagnóstico") || !strings.Contains(result, "## Plano") {
		t.Errorf("got: %q", result)
	}
}
