package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dilsonrabelo/castor-prompt-builder/pkg/engine"
	"github.com/dilsonrabelo/castor-prompt-builder/pkg/parser"
)

// App é o struct principal exposto ao frontend via Wails.
type App struct {
	ctx    context.Context
	models []*parser.Model
	roles  []*parser.Role
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// Localiza as pastas roles/ e models/ relativas ao binário
	base := execDir()

	models, err := parser.LoadAllModels(filepath.Join(base, "models"))
	if err != nil {
		models = []*parser.Model{}
	}
	a.models = models

	roles, err := parser.LoadAllRoles(filepath.Join(base, "roles"))
	if err != nil {
		roles = []*parser.Role{}
	}
	a.roles = roles
}

// execDir localiza o diretório raiz do projeto buscando a pasta models/.
// Funciona tanto em dev (wails dev) quanto no .app distribuído.
func execDir() string {
	candidates := []string{}

	// 1) relativo ao cwd (wails dev roda de dentro de desktop/)
	if cwd, err := os.Getwd(); err == nil {
		candidates = append(candidates,
			cwd,
			filepath.Dir(cwd),           // desktop/ → project root
			filepath.Join(cwd, ".."),
			filepath.Join(cwd, "../.."),
		)
	}

	// 2) relativo ao executável
	if exe, err := os.Executable(); err == nil {
		dir := filepath.Dir(exe)
		candidates = append(candidates,
			dir,
			filepath.Join(dir, ".."),
			filepath.Join(dir, "../.."),
			filepath.Join(dir, "../../.."),
		)
	}

	for _, c := range candidates {
		abs, err := filepath.Abs(c)
		if err != nil {
			continue
		}
		if _, err := os.Stat(filepath.Join(abs, "models")); err == nil {
			return abs
		}
	}
	return "."
}

// ---- DTOs expostos ao frontend ----

type ModelDTO struct {
	ID       string `json:"id"`
	Nome     string `json:"nome"`
	Descricao string `json:"descricao"`
}

type RoleDTO struct {
	ID        string   `json:"id"`
	Nome      string   `json:"nome"`
	Categoria string   `json:"categoria"`
	Tom       string   `json:"tom"`
	GapsComuns []string `json:"gaps_comuns"`
}

type GapAnswerDTO struct {
	Pergunta string `json:"pergunta"`
	Resposta string `json:"resposta"`
}

type StepDTO struct {
	Titulo   string `json:"titulo"`
	Descricao string `json:"descricao"`
}

type BuildRequestDTO struct {
	ModelID    string         `json:"model_id"`
	RoleIDs    []string       `json:"role_ids"`
	Narrativa  string         `json:"narrativa"`
	GapAnswers []GapAnswerDTO `json:"gap_answers"`
	Steps      []StepDTO      `json:"steps"`
}

type BuildResultDTO struct {
	Conteudo  string `json:"conteudo"`
	Caminho   string `json:"caminho"`
	Erro      string `json:"erro,omitempty"`
}

// ---- Métodos expostos ao frontend ----

func (a *App) GetModels() []ModelDTO {
	out := make([]ModelDTO, len(a.models))
	for i, m := range a.models {
		out[i] = ModelDTO{ID: m.ID, Nome: m.Nome, Descricao: m.Descricao}
	}
	return out
}

func (a *App) GetRoles() []RoleDTO {
	out := make([]RoleDTO, len(a.roles))
	for i, r := range a.roles {
		out[i] = RoleDTO{
			ID:        r.ID,
			Nome:      r.Nome,
			Categoria: r.Categoria,
			Tom:       r.Tom,
			GapsComuns: r.GapsComuns,
		}
	}
	return out
}

func (a *App) BuildPrompt(req BuildRequestDTO) BuildResultDTO {
	// localiza modelo
	var modelo *parser.Model
	for _, m := range a.models {
		if m.ID == req.ModelID {
			modelo = m
			break
		}
	}
	if modelo == nil {
		return BuildResultDTO{Erro: "modelo não encontrado: " + req.ModelID}
	}

	// combina papéis selecionados
	var nomes, descs []string
	roleID := "papel"
	for _, rid := range req.RoleIDs {
		for _, r := range a.roles {
			if r.ID == rid {
				nomes = append(nomes, r.Nome)
				descs = append(descs, r.Descricao)
				if roleID == "papel" {
					roleID = r.ID
				}
				break
			}
		}
	}
	roleNome := strings.Join(nomes, " + ")
	if roleNome == "" {
		roleNome = "Especialista"
	}

	// monta values
	vals := engine.NewValues()
	vals.Fields["role"] = roleNome + ". " + strings.Join(descs, "\n\n")
	vals.Fields["action"] = req.Narrativa
	vals.Fields["context"] = req.Narrativa
	vals.Fields["task"] = req.Narrativa

	for i, ga := range req.GapAnswers {
		if strings.TrimSpace(ga.Resposta) != "" {
			vals.Fields[fmt.Sprintf("gap_%d", i)] = ga.Resposta
		}
	}

	if len(req.Steps) > 0 {
		steps := make([]parser.Step, len(req.Steps))
		for i, s := range req.Steps {
			steps[i] = parser.Step{Titulo: s.Titulo, Descricao: s.Descricao}
		}
		vals.Steps["fases"] = steps
	}

	rendered := engine.Render(modelo.Template, vals)

	return BuildResultDTO{Conteudo: rendered}
}
