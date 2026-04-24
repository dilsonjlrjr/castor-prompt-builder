package main

import (
	"fmt"
	"strings"

	"github.com/dilsonrabelo/castor-prompt-builder/pkg/parser"
)

// FileStatus resultado da validação de um arquivo.
type FileStatus struct {
	Arquivo  string `json:"arquivo"`
	Tipo     string `json:"tipo"`    // "model" | "role"
	OK       bool   `json:"ok"`
	Problema string `json:"problema,omitempty"`
}

// ValidateAll valida todos os models e roles carregados.
// Retorna um status por arquivo — ok ou com problema descrito.
func (a *App) ValidateAll() []FileStatus {
	var out []FileStatus

	for _, m := range a.models {
		issues := validateModel(m)
		fs := FileStatus{
			Arquivo: m.ID + ".md",
			Tipo:    "model",
			OK:      len(issues) == 0,
		}
		if len(issues) > 0 {
			fs.Problema = strings.Join(issues, "; ")
		}
		out = append(out, fs)
	}

	for _, r := range a.roles {
		issues := validateRole(r)
		fs := FileStatus{
			Arquivo: r.Categoria + "/" + r.ID + ".md",
			Tipo:    "role",
			OK:      len(issues) == 0,
		}
		if len(issues) > 0 {
			fs.Problema = strings.Join(issues, "; ")
		}
		out = append(out, fs)
	}

	return out
}

var validFieldTypes = map[string]bool{
	"text": true, "textarea": true, "select": true,
	"multiselect": true, "list": true, "steps": true,
}

func validateModel(m *parser.Model) []string {
	var issues []string
	if m.ID == "" {
		issues = append(issues, "id ausente")
	}
	if m.Nome == "" {
		issues = append(issues, "nome ausente")
	}
	if m.Template == "" {
		issues = append(issues, "template vazio")
	}
	if len(m.Campos) == 0 {
		issues = append(issues, "nenhum campo definido")
	}
	for _, c := range m.Campos {
		if c.ID == "" {
			issues = append(issues, "campo sem id")
			continue
		}
		if !validFieldTypes[string(c.Tipo)] {
			issues = append(issues, fmt.Sprintf("campo '%s': tipo inválido '%s'", c.ID, c.Tipo))
		}
		if (c.Tipo == parser.FieldSelect || c.Tipo == parser.FieldMultiselect) && len(c.Opcoes) == 0 {
			issues = append(issues, fmt.Sprintf("campo '%s': select sem opções", c.ID))
		}
	}
	if m.Template != "" && !strings.Contains(m.Template, "{{") {
		issues = append(issues, "template sem variáveis")
	}
	return issues
}

func validateRole(r *parser.Role) []string {
	var issues []string
	if r.ID == "" {
		issues = append(issues, "id ausente")
	}
	if r.Nome == "" {
		issues = append(issues, "nome ausente")
	}
	if r.Descricao == "" {
		issues = append(issues, "descrição ausente")
	}
	return issues
}
