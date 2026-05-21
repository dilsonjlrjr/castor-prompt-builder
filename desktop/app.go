package main

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/dilsonrabelo/castor-prompt-builder/pkg/engine"
	"github.com/dilsonrabelo/castor-prompt-builder/pkg/parser"
)

const (
	promptFileExt         = ".md"
	promptTimestampFormat = "20060102_150405"
	promptDisplayFormat   = "2006-01-02 15:04"
	defaultPromptTitle    = "Sem título"
)

//go:embed bundled
var bundledFS embed.FS

// userDataDir retorna ~/.castorprompt
func userDataDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".castorprompt"), nil
}

// ensureUserDir cria ~/.castorprompt na primeira execução e extrai:
//   - models/ com os 4 modelos embutidos
//   - roles/ com todos os papéis embutidos (organizados por categoria)
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
	_ = os.MkdirAll(filepath.Join(base, "models"), 0o755)
	_ = os.MkdirAll(filepath.Join(base, "roles"), 0o755)

	// extrai modelos embutidos → ~/.castorprompt/models/
	entries, _ := fs.ReadDir(bundledFS, "bundled")
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
			continue
		}
		data, _ := bundledFS.ReadFile("bundled/" + e.Name())
		_ = os.WriteFile(filepath.Join(base, "models", e.Name()), data, 0o644)
	}

	// extrai roles embutidos → ~/.castorprompt/roles/<categoria>/arquivo.md
	_ = fs.WalkDir(bundledFS, "bundled/roles", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, ".md") {
			return nil
		}
		// path = "bundled/roles/arquitetura/arquivo.md"
		rel := strings.TrimPrefix(path, "bundled/roles/")
		dest := filepath.Join(base, "roles", rel)
		_ = os.MkdirAll(filepath.Dir(dest), 0o755)
		data, _ := bundledFS.ReadFile(path)
		_ = os.WriteFile(dest, data, 0o644)
		return nil
	})

	return true
}

// findSubdir procura por um subdiretório (ex: "models" ou "roles") em vários
// candidatos e retorna o caminho absoluto do primeiro que existir com conteúdo.
// Candidatos verificados, em ordem:
//  1. Diretório do executável          (produção: roles/ ao lado do .exe)
//  2. cwd e pais (../  ../../)         (dev: projeto na raiz)
//  3. ~/.castorprompt                  (produção: fallback padrão)
func findSubdir(name string) string {
	hasDir := func(parent string) (string, bool) {
		abs, err := filepath.Abs(parent)
		if err != nil {
			return "", false
		}
		target := filepath.Join(abs, name)
		info, err := os.Stat(target)
		if err != nil || !info.IsDir() {
			return "", false
		}
		// verifica se tem ao menos um arquivo .md dentro (não vazio)
		entries, _ := os.ReadDir(target)
		for _, e := range entries {
			if !e.IsDir() && strings.HasSuffix(e.Name(), ".md") {
				return target, true
			}
			if e.IsDir() {
				subs, _ := os.ReadDir(filepath.Join(target, e.Name()))
				for _, s := range subs {
					if strings.HasSuffix(s.Name(), ".md") {
						return target, true
					}
				}
			}
		}
		return "", false
	}

	// 1. Ao lado do executável
	if exePath, err := os.Executable(); err == nil {
		if p, ok := hasDir(filepath.Dir(exePath)); ok {
			return p
		}
	}

	// 2. cwd e pais (dev mode)
	if cwd, err := os.Getwd(); err == nil {
		for _, c := range []string{cwd, filepath.Join(cwd, ".."), filepath.Join(cwd, "../..")} {
			if p, ok := hasDir(c); ok {
				return p
			}
		}
	}

	// 3. ~/.castorprompt
	if base, err := userDataDir(); err == nil {
		if p, ok := hasDir(base); ok {
			return p
		}
	}

	return name // fallback relativo
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

	base, _ := userDataDir()

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

// ---- Lookups internos ----

func (a *App) findModel(id string) *parser.Model {
	for _, m := range a.models {
		if m.ID == id {
			return m
		}
	}
	return nil
}

func (a *App) findRole(id string) *parser.Role {
	for _, r := range a.roles {
		if r.ID == id {
			return r
		}
	}
	return nil
}

// resolveRoles devolve, na ordem dos IDs, os papéis encontrados (ignora IDs inválidos).
func (a *App) resolveRoles(ids []string) []*parser.Role {
	out := make([]*parser.Role, 0, len(ids))
	for _, id := range ids {
		if r := a.findRole(id); r != nil {
			out = append(out, r)
		}
	}
	return out
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
	modelo := a.findModel(req.ModelID)
	if modelo == nil {
		return BuildResultDTO{Erro: "modelo não encontrado: " + req.ModelID}
	}

	selected := a.resolveRoles(req.RoleIDs)

	nomes := make([]string, 0, len(selected))
	descs := make([]string, 0, len(selected))
	for _, r := range selected {
		nomes = append(nomes, r.Nome)
		descs = append(descs, r.Descricao)
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

	// habilidades dos papéis selecionados (dedup, preserva ordem de descoberta)
	seenH := map[string]bool{}
	var habs []string
	for _, r := range selected {
		for _, h := range r.Habilidades {
			if !seenH[h] {
				seenH[h] = true
				habs = append(habs, h)
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
	for _, r := range selected {
		if r.Tom != "" {
			toms = append(toms, r.Nome+": "+r.Tom)
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

// ---- Gerenciamento de prompts salvos ----

// PromptMetaDTO representa um prompt salvo (listagem).
type PromptMetaDTO struct {
	ID        string `json:"id"`        // nome do arquivo sem extensao
	Titulo    string `json:"titulo"`    // extraido do conteudo (# Prompt — ...)
	Data      string `json:"data"`      // data de criacao
}

// PromptDTO representa um prompt completo.
type PromptDTO struct {
	ID        string `json:"id"`
	Conteudo  string `json:"conteudo"`
}

func (a *App) promptsDir() string {
	base, _ := userDataDir()
	dir := filepath.Join(base, "prompts")
	_ = os.MkdirAll(dir, 0o755)
	return dir
}

func (a *App) promptPath(id string) string {
	return filepath.Join(a.promptsDir(), id+promptFileExt)
}

func writePrompt(path, conteudo string) {
	_ = os.WriteFile(path, []byte(conteudo), 0o644)
}

// SavePrompt grava o conteúdo em ~/.castorprompt/prompts/<timestamp>.md
func (a *App) SavePrompt(conteudo string) PromptMetaDTO {
	now := time.Now()
	id := now.Format(promptTimestampFormat)
	writePrompt(a.promptPath(id), conteudo)
	return PromptMetaDTO{
		ID:     id,
		Titulo: extrairTitulo(conteudo),
		Data:   now.Format(promptDisplayFormat),
	}
}

// ListPrompts lista prompts salvos ordenados por data de modificação (mais recentes primeiro).
func (a *App) ListPrompts() []PromptMetaDTO {
	dir := a.promptsDir()
	entries, err := os.ReadDir(dir)
	if err != nil {
		return []PromptMetaDTO{}
	}

	type item struct {
		meta    PromptMetaDTO
		modTime time.Time
	}
	items := make([]item, 0, len(entries))
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), promptFileExt) {
			continue
		}
		data, _ := os.ReadFile(filepath.Join(dir, e.Name()))
		info, _ := e.Info()
		mt := time.Time{}
		dataStr := ""
		if info != nil {
			mt = info.ModTime()
			dataStr = mt.Format(promptDisplayFormat)
		}
		items = append(items, item{
			meta: PromptMetaDTO{
				ID:     strings.TrimSuffix(e.Name(), promptFileExt),
				Titulo: extrairTitulo(string(data)),
				Data:   dataStr,
			},
			modTime: mt,
		})
	}
	sort.Slice(items, func(i, j int) bool { return items[i].modTime.After(items[j].modTime) })

	out := make([]PromptMetaDTO, len(items))
	for i, it := range items {
		out[i] = it.meta
	}
	return out
}

// GetPrompt retorna o conteúdo completo de um prompt.
func (a *App) GetPrompt(id string) PromptDTO {
	data, _ := os.ReadFile(a.promptPath(id))
	return PromptDTO{ID: id, Conteudo: string(data)}
}

// UpdatePrompt sobrescreve o conteúdo de um prompt existente.
func (a *App) UpdatePrompt(id string, conteudo string) {
	writePrompt(a.promptPath(id), conteudo)
}

// DeletePrompt remove um prompt salvo.
func (a *App) DeletePrompt(id string) {
	_ = os.Remove(a.promptPath(id))
}

// extrairTitulo retorna o texto do primeiro heading H1 do conteúdo,
// ou um título padrão quando ausente.
func extrairTitulo(conteudo string) string {
	for _, l := range strings.SplitN(conteudo, "\n", 3) {
		l = strings.TrimSpace(l)
		if strings.HasPrefix(l, "# ") {
			return strings.TrimPrefix(l, "# ")
		}
	}
	return defaultPromptTitle
}
