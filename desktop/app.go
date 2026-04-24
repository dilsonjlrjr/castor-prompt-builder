package main

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/dilsonrabelo/castor-prompt-builder/pkg/engine"
	"github.com/dilsonrabelo/castor-prompt-builder/pkg/parser"
)

//go:embed bundled/*.md
var bundledFS embed.FS

// userDataDir retorna ~/.castorprompt
func userDataDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".castorprompt"), nil
}

// ensureUserDir cria ~/.castorprompt na primeira execução:
//   - models/ com os 4 modelos embutidos
//   - roles/ vazio (o usuário adiciona os seus)
//
// Retorna true se o diretório foi criado agora (primeira execução).
func ensureUserDir() bool {
	base, err := userDataDir()
	if err != nil {
		return false
	}
	// já existe → nada a fazer
	if _, err := os.Stat(base); err == nil {
		return false
	}
	modelsDir := filepath.Join(base, "models")
	rolesDir := filepath.Join(base, "roles")
	_ = os.MkdirAll(modelsDir, 0o755)
	_ = os.MkdirAll(rolesDir, 0o755)

	// escreve os modelos embutidos
	entries, _ := fs.ReadDir(bundledFS, "bundled")
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
			continue
		}
		data, err := bundledFS.ReadFile("bundled/" + e.Name())
		if err != nil {
			continue
		}
		_ = os.WriteFile(filepath.Join(modelsDir, e.Name()), data, 0o644)
	}
	return true
}

// execDir resolve o diretório base onde ficam models/ e roles/.
//
// Prioridade:
//  1. Raiz do projeto via cwd (dev: desktop/../ tem models/ E roles/)
//  2. ~/.castorprompt (produção, fallback)
func execDir() string {
	hasBoth := func(dir string) bool {
		abs, err := filepath.Abs(dir)
		if err != nil {
			return false
		}
		_, errM := os.Stat(filepath.Join(abs, "models"))
		_, errR := os.Stat(filepath.Join(abs, "roles"))
		return errM == nil && errR == nil
	}

	// 1. Candidatos baseados no cwd — forte sinal de dev
	if cwd, err := os.Getwd(); err == nil {
		for _, c := range []string{
			cwd,
			filepath.Join(cwd, ".."),
			filepath.Join(cwd, "../.."),
		} {
			if hasBoth(c) {
				if abs, err := filepath.Abs(c); err == nil {
					return abs
				}
			}
		}
	}

	// 2. ~/.castorprompt (produção)
	if base, err := userDataDir(); err == nil {
		if hasBoth(base) {
			return base
		}
	}

	return "."
}

// App é o struct principal exposto ao frontend via Wails.
type App struct {
	ctx      context.Context
	models   []*parser.Model
	roles    []*parser.Role
	firstRun bool
}

func NewApp() *App {
	return &App{}
}

// IsFirstRun retorna true se o app está sendo aberto pela primeira vez
// (o diretório ~/.castorprompt acabou de ser criado).
func (a *App) IsFirstRun() bool {
	return a.firstRun
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	a.firstRun = ensureUserDir()
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
	ID          string   `json:"id"`
	Nome        string   `json:"nome"`
	Categoria   string   `json:"categoria"`
	Tom         string   `json:"tom"`
	GapsComuns  []string `json:"gaps_comuns"`
	Habilidades []string `json:"habilidades"`
}

// GapAnswerDTO: FieldID preenchido para campos do modelo; vazio para gaps de papel.
type GapAnswerDTO struct {
	FieldID  string `json:"field_id,omitempty"`
	Pergunta string `json:"pergunta"`
	Resposta string `json:"resposta"`
	RoleNome string `json:"role_nome,omitempty"` // preenchido para gaps de papel
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
			ID:          r.ID,
			Nome:        r.Nome,
			Categoria:   r.Categoria,
			Tom:         r.Tom,
			GapsComuns:  r.GapsComuns,
			Habilidades: r.Habilidades,
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

	// gap answers: campos do modelo (FieldID preenchido)
	for _, ga := range req.GapAnswers {
		if strings.TrimSpace(ga.Resposta) == "" || ga.FieldID == "" {
			continue
		}
		var tipo parser.FieldType
		for _, c := range modelo.Campos {
			if c.ID == ga.FieldID {
				tipo = c.Tipo
				break
			}
		}
		if tipo == parser.FieldList || tipo == parser.FieldMultiselect {
			var items []string
			for _, line := range strings.Split(ga.Resposta, "\n") {
				if t := strings.TrimSpace(line); t != "" {
					items = append(items, t)
				}
			}
			if len(items) > 0 {
				vals.Lists[ga.FieldID] = items
			}
		} else {
			vals.Fields[ga.FieldID] = ga.Resposta
		}
	}

	// fases — injeta no vals.Steps["fases"] para {{#steps fases}}
	if len(req.Steps) > 0 {
		steps := make([]parser.Step, len(req.Steps))
		for i, s := range req.Steps {
			steps[i] = parser.Step{Titulo: s.Titulo, Descricao: s.Descricao}
		}
		vals.Steps["fases"] = steps
	}

	rendered := engine.Render(modelo.Template, vals)

	// ── seções extras (fora do template principal) ───────────────────────────

	var extras strings.Builder

	// fases: se o template não tem {{#steps}}, injeta como seção genérica
	if len(req.Steps) > 0 && !strings.Contains(modelo.Template, "{{#steps") {
		extras.WriteString("\n\n---\n## Fases de execução\n\n")
		for i, s := range req.Steps {
			extras.WriteString(fmt.Sprintf("### Fase %d — %s\n%s\n\n", i+1, s.Titulo, s.Descricao))
		}
	}

	// habilidades dos papéis selecionados (dedup)
	seenH := map[string]bool{}
	var habs []string
	for _, rid := range req.RoleIDs {
		for _, r := range a.roles {
			if r.ID == rid {
				for _, h := range r.Habilidades {
					if !seenH[h] {
						seenH[h] = true
						habs = append(habs, h)
					}
				}
				break
			}
		}
	}
	if len(habs) > 0 {
		extras.WriteString("\n\n---\n## Habilidades relevantes\n")
		for _, h := range habs {
			extras.WriteString("- " + h + "\n")
		}
	}

	// tom dos papéis
	var toms []string
	for _, rid := range req.RoleIDs {
		for _, r := range a.roles {
			if r.ID == rid && r.Tom != "" {
				toms = append(toms, r.Nome+": "+r.Tom)
				break
			}
		}
	}
	if len(toms) > 0 {
		extras.WriteString("\n\n---\n## Tom de comunicação\n")
		for _, t := range toms {
			extras.WriteString("- " + t + "\n")
		}
	}

	// gaps de papel (FieldID vazio) — contexto adicional fornecido pelo usuário
	var gapCtx []string
	for _, ga := range req.GapAnswers {
		if ga.FieldID == "" && strings.TrimSpace(ga.Resposta) != "" {
			label := ga.Pergunta
			if ga.RoleNome != "" {
				label = ga.RoleNome + " — " + ga.Pergunta
			}
			gapCtx = append(gapCtx, "**"+label+"**\n"+ga.Resposta)
		}
	}
	if len(gapCtx) > 0 {
		extras.WriteString("\n\n---\n## Contexto dos papéis\n\n")
		for _, g := range gapCtx {
			extras.WriteString(g + "\n\n")
		}
	}

	if extras.Len() > 0 {
		rendered += extras.String()
	}

	return BuildResultDTO{Conteudo: rendered}
}
