package parser

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

func splitFrontmatter(content string) (string, string, error) {
	content = strings.TrimSpace(content)
	if !strings.HasPrefix(content, "---") {
		return "", content, nil
	}
	parts := strings.SplitN(content, "---", 3)
	if len(parts) < 3 {
		return "", content, fmt.Errorf("frontmatter malformado")
	}
	return strings.TrimSpace(parts[1]), strings.TrimSpace(parts[2]), nil
}

// LoadRole lê um arquivo .md de role e retorna a struct Role
func LoadRole(path string) (*Role, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	fm, body, err := splitFrontmatter(string(data))
	if err != nil {
		return nil, err
	}
	var role Role
	if err := yaml.NewDecoder(bytes.NewBufferString(fm)).Decode(&role); err != nil {
		return nil, fmt.Errorf("yaml role %s: %w", path, err)
	}
	role.Descricao = body
	return &role, nil
}

// LoadModel lê um arquivo .md de modelo e retorna a struct Model
func LoadModel(path string) (*Model, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	fm, body, err := splitFrontmatter(string(data))
	if err != nil {
		return nil, err
	}
	var model Model
	if err := yaml.NewDecoder(bytes.NewBufferString(fm)).Decode(&model); err != nil {
		return nil, fmt.Errorf("yaml model %s: %w", path, err)
	}
	// extrai apenas o bloco após "## Template de saída"
	if idx := strings.Index(body, "## Template de saída"); idx >= 0 {
		model.Template = strings.TrimSpace(body[idx+len("## Template de saída"):])
	} else {
		model.Template = body
	}
	return &model, nil
}

// LoadAllRoles carrega todos os .md do diretório rolesDir
func LoadAllRoles(rolesDir string) ([]*Role, error) {
	entries, err := filepath.Glob(filepath.Join(rolesDir, "*.md"))
	if err != nil {
		return nil, err
	}
	var roles []*Role
	for _, path := range entries {
		r, err := LoadRole(path)
		if err != nil {
			return nil, err
		}
		roles = append(roles, r)
	}
	return roles, nil
}

// LoadAllModels carrega todos os .md do diretório modelsDir
func LoadAllModels(modelsDir string) ([]*Model, error) {
	entries, err := filepath.Glob(filepath.Join(modelsDir, "*.md"))
	if err != nil {
		return nil, err
	}
	var models []*Model
	for _, path := range entries {
		m, err := LoadModel(path)
		if err != nil {
			return nil, err
		}
		models = append(models, m)
	}
	return models, nil
}
