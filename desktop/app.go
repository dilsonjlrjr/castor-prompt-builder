package main

import (
	"context"
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

func execDir() string {
	candidates := []string{}
	if cwd, err := os.Getwd(); err == nil {
		candidates = append(candidates,
			cwd,
			filepath.Dir(cwd),
			filepath.Join(cwd, ".."),
			filepath.Join(cwd, "../.."),
		)
	}
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

// CampoDTO representa um campo de modelo que precisa de entrada do usuário.
// Campos já mapeados automaticamente (role, action, context, task, input, fases)
// são excluídos — o frontend os usa para montar a tela de gaps.
type CampoDTO struct {
	ID          string   `json:"id"`
	Label       string   `json:"label"`
	Tipo        string   `json:"tipo"`
	Obrigatorio bool     `json:"obrigatorio"`
	Opcoes      []string `json:"opcoes,omitempty"`
}

type ModelDTO struct {
	ID        string     `json:"id"`
	Nome      string     `json:"nome"`
	Descricao string     `json:"descricao"`
	Campos    []CampoDTO `json:"campos"`
}

type RoleDTO struct {
	ID         string   `json:"id"`
	Nome       string   `json:"nome"`
	Categoria  string   `json:"categoria"`
	Tom        string   `json:"tom"`
	GapsComuns []string `json:"gaps_comuns"`
}

// GapAnswerDTO: FieldID preenchido para campos do modelo; vazio para gaps de papel.
type GapAnswerDTO struct {
	FieldID  string `json:"field_id,omitempty"`
	Pergunta string `json:"pergunta"`
	Resposta string `json:"resposta"`
}

type StepDTO struct {
	Titulo    string `json:"titulo"`
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
	Conteudo string `json:"conteudo"`
	Caminho  string `json:"caminho"`
	Erro     string `json:"erro,omitempty"`
}

// campos do modelo mapeados automaticamente a partir da narrativa
var autoMapped = map[string]bool{
	"role": true, "action": true, "context": true,
	"task": true, "input": true,
}

// ---- Métodos expostos ao frontend ----

func (a *App) GetModels() []ModelDTO {
	out := make([]ModelDTO, len(a.models))
	for i, m := range a.models {
		var campos []CampoDTO
		for _, c := range m.Campos {
			// Exclui campos mapeados automaticamente e campos do tipo steps
			// (fases são coletadas em tela separada)
			if autoMapped[c.ID] || c.Tipo == parser.FieldSteps {
				continue
			}
			campos = append(campos, CampoDTO{
				ID:          c.ID,
				Label:       c.Label,
				Tipo:        string(c.Tipo),
				Obrigatorio: c.Obrigatorio,
				Opcoes:      c.Opcoes,
			})
		}
		out[i] = ModelDTO{ID: m.ID, Nome: m.Nome, Descricao: m.Descricao, Campos: campos}
	}
	return out
}

func (a *App) GetRoles() []RoleDTO {
	out := make([]RoleDTO, len(a.roles))
	for i, r := range a.roles {
		out[i] = RoleDTO{
			ID:         r.ID,
			Nome:       r.Nome,
			Categoria:  r.Categoria,
			Tom:        r.Tom,
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
	for _, rid := range req.RoleIDs {
		for _, r := range a.roles {
			if r.ID == rid {
				nomes = append(nomes, r.Nome)
				descs = append(descs, r.Descricao)
				break
			}
		}
	}
	roleNome := strings.Join(nomes, " + ")
	if roleNome == "" {
		roleNome = "Especialista"
	}

	// monta values — campos mapeados automaticamente
	vals := engine.NewValues()
	vals.Fields["role"]    = roleNome + ". " + strings.Join(descs, "\n\n")
	vals.Fields["action"]  = req.Narrativa
	vals.Fields["context"] = req.Narrativa
	vals.Fields["task"]    = req.Narrativa
	vals.Fields["input"]   = req.Narrativa

	// gap answers: usa FieldID quando disponível (campos do modelo),
	// caso contrário armazena como contexto extra
	for _, ga := range req.GapAnswers {
		if strings.TrimSpace(ga.Resposta) == "" {
			continue
		}
		if ga.FieldID != "" {
			vals.Fields[ga.FieldID] = ga.Resposta
		}
		// gaps de papel sem fieldId são contexto extra — não entram no template diretamente
	}

	// fases
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
