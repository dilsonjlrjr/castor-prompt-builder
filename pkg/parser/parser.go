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
	// extrai apenas o bloco após o cabeçalho do template (em qualquer idioma)
	model.Template = extractTemplate(body)
	return &model, nil
}

// Cabeçalhos aceitos para a seção de template em cada idioma suportado.
var templateMarkers = []string{
	"## Template de saída",   // pt
	"## Output template",     // en
	"## Plantilla de salida", // es
}

func extractTemplate(body string) string {
	for _, marker := range templateMarkers {
		if idx := strings.Index(body, marker); idx >= 0 {
			return strings.TrimSpace(body[idx+len(marker):])
		}
	}
	return body
}

// LoadAllRoles carrega roles em PT (compatibilidade — usa LoadAllRolesLang).
func LoadAllRoles(rolesDir string) ([]*Role, error) {
	return LoadAllRolesLang(rolesDir, "pt")
}

// LoadAllModels carrega modelos em PT (compatibilidade — usa LoadAllModelsLang).
func LoadAllModels(modelsDir string) ([]*Model, error) {
	return LoadAllModelsLang(modelsDir, "pt")
}

// stripLangSuffix devolve (base, lang) para "arquiteto_cloud.en.md" → ("arquiteto_cloud", "en"),
// e ("arquiteto_cloud", "") quando não há sufixo. Sufixos aceitos: pt, en, es.
func stripLangSuffix(name string) (base, lang string) {
	if !strings.HasSuffix(name, ".md") {
		return name, ""
	}
	noExt := strings.TrimSuffix(name, ".md")
	if dot := strings.LastIndex(noExt, "."); dot > 0 {
		switch cand := noExt[dot+1:]; cand {
		case "pt", "en", "es":
			return noExt[:dot], cand
		}
	}
	return noExt, ""
}

// resolveLangFile retorna o caminho preferido para uma base name no idioma pedido:
// 1) <base>.<lang>.md, 2) <base>.md (fallback PT). String vazia se nada existir.
func resolveLangFile(dir, base, lang string) string {
	if lang != "" {
		p := filepath.Join(dir, base+"."+lang+".md")
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	p := filepath.Join(dir, base+".md")
	if _, err := os.Stat(p); err == nil {
		return p
	}
	return ""
}

// LoadAllRolesLang carrega roles recursivamente, escolhendo a variante de idioma
// quando existir (<base>.<lang>.md) e caindo no PT (<base>.md) caso contrário.
func LoadAllRolesLang(rolesDir, lang string) ([]*Role, error) {
	if lang == "" {
		lang = "pt"
	}
	// dir → conjunto de base names únicos
	bases := map[string]map[string]bool{}
	err := filepath.WalkDir(rolesDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(path, ".md") {
			return nil
		}
		dir := filepath.Dir(path)
		base, _ := stripLangSuffix(d.Name())
		if bases[dir] == nil {
			bases[dir] = map[string]bool{}
		}
		bases[dir][base] = true
		return nil
	})
	if err != nil {
		return nil, err
	}
	var roles []*Role
	for dir, set := range bases {
		for base := range set {
			path := resolveLangFile(dir, base, lang)
			if path == "" {
				continue
			}
			r, err := LoadRole(path)
			if err != nil {
				return nil, err
			}
			roles = append(roles, r)
		}
	}
	sort.Slice(roles, func(i, j int) bool {
		if roles[i].Categoria != roles[j].Categoria {
			return roles[i].Categoria < roles[j].Categoria
		}
		return roles[i].Nome < roles[j].Nome
	})
	return roles, nil
}

// LoadAllModelsLang carrega modelos no idioma pedido (mesmo esquema dos roles).
func LoadAllModelsLang(modelsDir, lang string) ([]*Model, error) {
	if lang == "" {
		lang = "pt"
	}
	entries, err := os.ReadDir(modelsDir)
	if err != nil {
		return nil, err
	}
	bases := map[string]bool{}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
			continue
		}
		base, _ := stripLangSuffix(e.Name())
		bases[base] = true
	}
	var models []*Model
	for base := range bases {
		path := resolveLangFile(modelsDir, base, lang)
		if path == "" {
			continue
		}
		m, err := LoadModel(path)
		if err != nil {
			return nil, err
		}
		models = append(models, m)
	}
	sort.Slice(models, func(i, j int) bool { return models[i].ID < models[j].ID })
	return models, nil
}
