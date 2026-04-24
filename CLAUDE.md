# CASTOR Builder — Context & Design

> Documento de referência do projeto. Atualizar conforme decisões evoluem.

---

## Visão Geral

App desktop (Wails v2 + Svelte + Go) que guia o usuário em um wizard para construir prompts estruturados e semanticamente ricos. O usuário escolhe um **modelo de prompt**, seleciona um ou mais **papéis** (especialistas), descreve a tarefa em linguagem livre, responde perguntas de contexto (gaps) e opcionalmente define fases de execução.

Sem IA no processo. Tudo é parsing, heurística e template engine próprio.

---

## Stack

| Camada | Tecnologia |
|---|---|
| Desktop framework | [Wails v2](https://wails.io) |
| Backend | Go |
| Frontend | Svelte + TypeScript + Tailwind (via Vite) |
| Parser frontmatter | `gopkg.in/yaml.v3` |
| Template engine | Próprio (`{{campo}}`, `{{#if}}`, `{{#steps}}`, `{{#each}}`) |

---

## Estrutura de Diretórios

```
castor-prompt-builder/
├── CLAUDE.md
├── desktop/                        ← projeto Wails
│   ├── main.go
│   ├── app.go                      ← lógica principal (GetModels, GetRoles, BuildPrompt)
│   ├── validate.go                 ← ValidateAll — valida models e roles ao iniciar
│   ├── bundled/                    ← models e roles embutidos via embed.FS
│   │   ├── models/                 ← rtf.md, race.md, risen.md, create.md
│   │   └── roles/                  ← papéis organizados por categoria/
│   ├── build/                      ← assets de build (ícones, manifests)
│   ├── frontend/
│   │   ├── src/App.svelte          ← toda a UI (wizard de 6 telas)
│   │   ├── wailsjs/go/main/        ← bindings gerados (App.js, App.d.ts)
│   │   └── wailsjs/go/models.ts   ← DTOs gerados
│   └── pkg/
│       ├── parser/                 ← leitura dos .md de roles e models
│       └── engine/                 ← template engine
├── dist/                           ← artefatos de distribuição
│   ├── castor-builder-macos.app
│   └── castor-builder-windows-amd64.exe
└── Makefile                        ← targets: macos, windows, icons
```

---

## Fluxo do App (Wizard)

```
[1] Validação        → tela animada ao iniciar, valida cada model/role
[2] Onboarding       → carrossel na primeira execução (IsFirstRun)
[3] Selecionar modelo → RTF / RACE / RISEN / CREATE
[4] Selecionar papéis → multi-seleção com busca por categoria
[5] Inserir narrativa → textarea livre
[6] Preencher gaps    → campos do modelo + gaps_comuns dos papéis (uma pergunta por tela)
[7] Definir fases     → opcional, título + descrição por fase
[8] Resultado         → prompt gerado, copiável
```

---

## Comportamento do BuildPrompt (app.go)

1. Localiza o modelo pelo `model_id`
2. Monta `vals` com os `gap_answers` que têm `field_id` não vazio (campos do modelo)
3. Chama `engine.Render(template, vals)` → prompt base
4. **Extras injetados após o render:**
   - **Fases**: se o usuário definiu steps e o template não tem `{{#steps`, injeta `## Fases de execução` com cada fase numerada
   - **Habilidades**: agrega `habilidades` de todos os papéis (dedup) → `## Habilidades relevantes`
   - **Tom**: agrega `tom` de cada papel → `## Tom de comunicação`
   - **Contexto dos papéis**: gaps de `gaps_comuns` respondidos pelo usuário (field_id vazio) → `## Contexto dos papéis` com atribuição (`RoleNome — Pergunta`)

---

## Especificação: Arquivo de Papel (`bundled/roles/**/*.md`)

```markdown
---
id: arquiteto_cloud
nome: Arquiteto Cloud
categoria: arquitetura
tom: técnico e pragmático
habilidades:
  - AWS / GCP / Azure
  - infraestrutura como código
gaps_comuns:
  - Qual o cloud provider principal?
  - Existe multi-cloud ou hybrid cloud?
---

Descrição rica do papel (vai no prompt gerado).
```

### Campos do frontmatter

| Campo | Tipo | Uso |
|---|---|---|
| `id` | string | Referência interna |
| `nome` | string | Exibido na UI |
| `categoria` | string | Agrupa na tela de seleção |
| `tom` | string | Injetado como `## Tom de comunicação` |
| `habilidades` | list | Injetadas como `## Habilidades relevantes` |
| `gaps_comuns` | list | Perguntas de contexto específicas do papel |

---

## Especificação: Arquivo de Modelo (`bundled/models/*.md`)

```markdown
---
id: race
nome: RACE
descricao: Ideal para tarefas com contexto rico e expectativa clara de entrega
campos:
  - id: role
    label: Papel
    tipo: text
    obrigatorio: true
  - id: action
    label: Ação
    tipo: textarea
    obrigatorio: true
  - id: fases
    label: Fases de execução
    tipo: steps
    obrigatorio: false
---

## Template de saída

Você é {{role}}.

{{action}}

{{#steps fases}}
## {{titulo}}
{{descricao}}
{{/steps}}
```

### Tipos de Campos

| Tipo | Comportamento |
|---|---|
| `text` | Input linha única |
| `textarea` | Input multi-linha |
| `select` | Escolhe uma opção |
| `multiselect` | Escolhe várias opções |
| `list` | Adiciona itens livres |
| `steps` | Coleta título + descrição por fase |

### Diretivas de Template

| Diretiva | Comportamento |
|---|---|
| `{{campo}}` | Substitui pelo valor do campo |
| `{{#if campo}}...{{/if}}` | Renderiza se campo preenchido |
| `{{#steps campo}}...{{/steps}}` | Itera lista de steps |
| `{{#each campo}}...{{/each}}` | Itera list/multiselect |

---

## Modelos Disponíveis

| ID | Nome | Quando usar |
|---|---|---|
| `rtf` | RTF — Role, Task, Feature | Tarefas simples e diretas |
| `race` | RACE — Role, Action, Context, Expectation | Contexto rico + entregável claro |
| `risen` | RISEN — Role, Input, Steps, Expectation, Narrowing | Steps detalhados com restrições |
| `create` | CREATE — Context, Role, Examples, Audience, Tone, Expectation | Conteúdo criativo |

---

## Diretório do Usuário (`~/.castorprompt`)

Criado na primeira execução (`IsFirstRun`). Modelos bundled são copiados para lá.  
O usuário pode adicionar papéis customizados em `~/.castorprompt/roles/`.

**Prioridade de carregamento (execDir):**
1. `./models` e `./roles` junto ao binário (dev mode)
2. `~/.castorprompt/models` e `~/.castorprompt/roles` (produção)

---

## Validação ao Inicializar

`ValidateAll()` roda ao montar o app, valida cada model e role carregado:
- Model: id, nome, template não vazio, campos definidos, tipos válidos, select com opções
- Role: id, nome, descrição

A tela de validação exibe cada arquivo com animação item a item. Erros bloqueiam e listam os problemas. Se tudo OK, fecha automaticamente em 800ms.

---

## DTOs (models.ts / app.go)

| DTO | Campos relevantes |
|---|---|
| `BuildRequestDTO` | `model_id`, `role_ids[]`, `narrativa`, `gap_answers[]`, `steps[]` |
| `GapAnswerDTO` | `field_id`, `pergunta`, `resposta`, `role_nome` |
| `BuildResultDTO` | `conteudo`, `caminho`, `erro` |
| `ModelDTO` | `id`, `nome`, `descricao`, `campos[]` |
| `RoleDTO` | `id`, `nome`, `categoria`, `tom`, `gaps_comuns[]`, `habilidades[]` |
| `FileStatus` | `arquivo`, `tipo`, `ok`, `problema` |

---

## Decisões Tomadas

- [x] Interface: app desktop com Wails v2 (abandonada TUI bubbletea)
- [x] Parser do frontmatter: `gopkg.in/yaml.v3`
- [x] Navegação: wizard linear com voltar
- [x] Gaps de papel: coletados, deduplicados entre roles, injetados com atribuição
- [x] Habilidades e tom dos papéis: injetados como seções extras no prompt
- [x] Fases: fallback genérico para modelos sem `{{#steps}}`
- [x] Validação animada ao inicializar
- [x] Onboarding carrossel na primeira execução
- [x] Diretório `~/.castorprompt` criado na primeira run com bundled models
- [x] Build: macOS universal + Windows amd64 via Makefile

---

## Diretrizes

- NUNCA ao fazer commit colocar o claude como autor
- Sempre usar o context-mode em todos os processos
- Sempre usar caveman
- Se perceber que as features são grandes, crie fases
- Ao fim de cada fase ou conclusão de atividade, faça commit usando git semântico
