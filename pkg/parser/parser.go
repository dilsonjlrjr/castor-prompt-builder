package parser

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
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

// LoadRole lê um arquivo .md de role e retorna a struct Role.
// categoria é derivada do nome do diretório pai (ex: "frontend", "devops").
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
	role.Categoria = filepath.Base(filepath.Dir(path))
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

// LoadAllRoles carrega recursivamente todos os .md de rolesDir e subdiretórios.
// Roles são ordenados por categoria e depois por nome.
func LoadAllRoles(rolesDir string) ([]*Role, error) {
	var roles []*Role
	err := filepath.WalkDir(rolesDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(path, ".md") {
			return nil
		}
		r, err := LoadRole(path)
		if err != nil {
			return err
		}
		roles = append(roles, r)
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Slice(roles, func(i, j int) bool {
		if roles[i].Categoria != roles[j].Categoria {
			return roles[i].Categoria < roles[j].Categoria
		}
		return roles[i].Nome < roles[j].Nome
	})
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
