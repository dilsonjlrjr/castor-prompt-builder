# Prompt Builder — Context & Design

> Documento de referência do projeto. Atualizar conforme decisões evoluem.

---

## Visão Geral

CLI em Go que recebe um **papel** (escolhido de um banco) e uma **narrativa livre**, extrai os campos necessários, preenche os gaps, opcionalmente faseia o processo e gera um prompt estruturado e semanticamente rico salvo em `.md`.

Sem IA no processo. Tudo é parsing, heurística e template engine próprio.

---

## Fluxo do Programa

```
1. Selecionar modelo de prompt   → lista os .md em /models
2. Selecionar papel              → lista os .md em /roles
3. Inserir narrativa             → textarea livre
4. Preencher gaps                → campos obrigatórios não extraídos viram perguntas
5. Fasear? (s/n)                 → se sim, coleta steps (título + descrição por fase)
6. Montar prompt                 → renderiza template do modelo com os valores
7. Salvar                        → /prompts/<data>_<papel>_<slug>.md
```

---

## Estrutura de Diretórios

```
prompt-builder/
├── main.go
├── context.md               ← este arquivo
├── roles/                   ← banco de papéis
│   ├── marketing_content.md
│   ├── data_analyst.md
│   └── devops_engineer.md
├── models/                  ← modelos de prompt
│   ├── rtf.md
│   ├── race.md
│   ├── risen.md
│   └── create.md
└── prompts/                 ← gerados pelo programa
    └── 20240418_marketing_content_plano_editorial.md
```

---

## Especificação: Arquivo de Papel (`/roles/*.md`)

```markdown
---
id: marketing_content
nome: Especialista em Marketing de Conteúdo
tom: analítico e estratégico
habilidades:
  - planejamento editorial
  - análise de métricas
  - SEO e copywriting
  - mapeamento de personas
gaps_comuns:
  - Qual o público-alvo?
  - Qual o tom de comunicação da marca?
  - Existe orçamento definido?
  - Quais canais já são utilizados?
---

Profissional com sólida experiência em estratégia editorial B2B e B2C.
Domina a criação de planos de conteúdo orientados a métricas, com foco
em engajamento, geração de leads e posicionamento de marca. Capaz de
alinhar narrativa de conteúdo ao funil de vendas e aos objetivos do negócio.
```

### Campos do frontmatter

| Campo | Tipo | Descrição |
|---|---|---|
| `id` | string | Referência interna, deve bater com o nome do arquivo |
| `nome` | string | Exibido na lista de seleção do CLI |
| `tom` | string | Instrução de estilo inserida no prompt |
| `habilidades` | list | Enriquece a descrição do Role no prompt |
| `gaps_comuns` | list | Perguntas específicas deste papel quando campo não preenchido |

O **corpo** do `.md` (abaixo do frontmatter) é a descrição rica do papel — vai direto no prompt gerado.

---

## Especificação: Arquivo de Modelo (`/models/*.md`)

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

  - id: context
    label: Contexto
    tipo: textarea
    obrigatorio: true

  - id: tom
    label: Tom de comunicação
    tipo: select
    opcoes:
      - formal
      - informal
      - técnico
      - persuasivo
    obrigatorio: false

  - id: canais
    label: Canais de distribuição
    tipo: multiselect
    opcoes:
      - blog
      - LinkedIn
      - email
      - Instagram
    obrigatorio: false

  - id: fases
    label: Fases de execução
    tipo: steps
    obrigatorio: false
    step_campos:
      - id: titulo
        label: Título da fase
        tipo: text
      - id: descricao
        label: O que deve ser feito nesta fase
        tipo: textarea

  - id: expectation
    label: Expectativa
    tipo: textarea
    obrigatorio: true
---

## Template de saída

Você é {{role}}.

{{action}}

O contexto é o seguinte: {{context}}

{{#if tom}}
Adote um tom {{tom}}.
{{/if}}

{{#steps fases}}
## {{titulo}}
{{descricao}}
{{/steps}}

{{#if canais}}
Considere os seguintes canais: {{#each canais}}{{.}}{{/each}}.
{{/if}}

Espera-se que {{expectation}}
```

---

## Tipos de Campos

| Tipo | Comportamento no CLI | Uso típico |
|---|---|---|
| `text` | Input livre, linha única | Título, tom livre |
| `textarea` | Input livre, multi-linha (Enter duplo encerra) | Narrativa, contexto, descrição |
| `select` | Lista numerada, escolhe uma | Tom, formato de saída |
| `multiselect` | Lista numerada, escolhe várias (ex: `1,3`) | Canais, públicos |
| `list` | Adiciona itens livres um a um (linha vazia encerra) | Restrições, habilidades extras |
| `steps` | Pergunta quantidade, depois coleta título + descrição por fase | Fases de execução |

---

## Diretivas de Template

| Diretiva | Comportamento |
|---|---|
| `{{campo}}` | Substitui pelo valor simples do campo |
| `{{#if campo}}...{{/if}}` | Renderiza bloco apenas se campo foi preenchido |
| `{{#steps campo}}...{{/steps}}` | Itera lista de steps, expõe `{{titulo}}` e `{{descricao}}` |
| `{{#each campo}}...{{/each}}` | Itera `list` ou `multiselect`, expõe `{{.}}` para o item atual |

---

## Modelos Planejados

| ID | Nome | Quando usar |
|---|---|---|
| `rtf` | RTF — Role, Task, Feature | Tarefas simples e diretas |
| `race` | RACE — Role, Action, Context, Expectation | Contexto rico + entregável claro |
| `risen` | RISEN — Role, Input, Steps, Expectation, Narrowing | Quando precisa de steps detalhados com restrições |
| `create` | CREATE — Context, Role, Examples, Audience, Tone, Expectation | Conteúdo criativo com restrições de público e tom |
| `cot` | Chain of Thought | Raciocínio passo a passo explícito |

---

## Prompt Gerado — Exemplo de Saída

```markdown
# Prompt — Especialista em Marketing de Conteúdo
_Modelo: RACE | Gerado em: 2024-04-18_

## Papel
Você é um Especialista em Marketing de Conteúdo. Profissional com sólida
experiência em estratégia editorial B2B e B2C. Domina a criação de planos
de conteúdo orientados a métricas, com foco em engajamento, geração de leads
e posicionamento de marca.

## Ação
Crie um plano editorial para os próximos 3 meses.

## Contexto
Startup B2B de SaaS enfrentando queda de 15% no engajamento do blog.
Público-alvo: desenvolvedores e CTOs de empresas mid-market.

## Fase 1 — Diagnóstico
Analise as possíveis causas da queda de engajamento considerando o perfil
da empresa e o público informado. Liste os 3 principais fatores com justificativa.

## Fase 2 — Plano Editorial
Com base no diagnóstico, estruture o plano editorial para 3 meses. Inclua
temas por mês, frequência de publicação e canais de distribuição.

## Fase 3 — Entregável Final
Consolide em um documento estruturado pronto para apresentação ao time de marketing.

## Expectativa
Documento com temas, frequência e canais de distribuição por mês.

---
> ⚠️ Para refinar este prompt considere informar:
> - Orçamento disponível para produção de conteúdo?
> - Existe time interno ou será freelancer?
```

---

## Convenção de Nomenclatura dos Prompts Salvos

```
/prompts/<YYYYMMDD>_<role_id>_<slug-da-narrativa>.md

Exemplo:
/prompts/20240418_marketing_content_plano-editorial-3-meses.md
```

O slug é gerado a partir das primeiras palavras da narrativa, em lowercase, separadas por hífen, sem acentos.

---

## Interface — TUI (Terminal UI)

### Stack escolhido

| Lib | Papel |
|---|---|
| [`bubbletea`](https://github.com/charmbracelet/bubbletea) | Engine principal, arquitetura Elm (Model / Update / View) |
| [`bubbles`](https://github.com/charmbracelet/bubbles) | Componentes prontos: list, textarea, textinput, spinner |
| [`lipgloss`](https://github.com/charmbracelet/lipgloss) | Estilo, cores, bordas, layout |

### Navegação

Wizard linear com possibilidade de **voltar** para tela anterior.

### Fluxo de telas

```
[1] Selecionar modelo
[2] Selecionar papel
[3] Inserir narrativa
[4] Preencher gaps        (uma pergunta por tela)
[5] Fasear? s/n
[6] Definir fases         (uma fase por tela, se aplicável)
[7] Confirmação + salvar
```

### Esboço das telas

```
┌─────────────────────────────────────┐
│  Prompt Builder                     │
├─────────────────────────────────────┤
│  Selecione o modelo                 │
│                                     │
│  > RACE                             │
│    RTF                              │
│    RISEN                            │
│                                     │
│  ↑↓ navegar   Enter selecionar      │
└─────────────────────────────────────┘

┌─────────────────────────────────────┐
│  Prompt Builder        modelo: RACE │
├─────────────────────────────────────┤
│  Selecione o papel                  │
│                                     │
│  > Especialista em Marketing        │
│    Analista de Dados                │
│    Engenheiro DevOps                │
│                                     │
│  ↑↓ navegar   Enter selecionar      │
└─────────────────────────────────────┘

┌─────────────────────────────────────┐
│  Prompt Builder          papel: MKT │
├─────────────────────────────────────┤
│  Narrativa                          │
│  ┌───────────────────────────────┐  │
│  │ A empresa é uma startup B2B   │  │
│  │ de SaaS com queda de 15%...   │  │
│  └───────────────────────────────┘  │
│                                     │
│  Ctrl+S confirmar   Esc voltar      │
└─────────────────────────────────────┘

┌─────────────────────────────────────┐
│  Gaps detectados        2 de 3      │
├─────────────────────────────────────┤
│  Qual o público-alvo?               │
│  ┌───────────────────────────────┐  │
│  │ Desenvolvedores e CTOs        │  │
│  └───────────────────────────────┘  │
│                                     │
│  Tab próximo   Esc voltar           │
└─────────────────────────────────────┘

┌─────────────────────────────────────┐
│  Fase 1 de 3                        │
├─────────────────────────────────────┤
│  Título                             │
│  ┌───────────────────────────────┐  │
│  │ Diagnóstico                   │  │
│  └───────────────────────────────┘  │
│  Descrição                          │
│  ┌───────────────────────────────┐  │
│  │ Analise as causas da queda... │  │
│  └───────────────────────────────┘  │
│                                     │
│  Tab próximo campo   Ctrl+S avançar │
└─────────────────────────────────────┘

┌─────────────────────────────────────┐
│  ✓ Prompt gerado!                   │
├─────────────────────────────────────┤
│  Salvo em:                          │
│  prompts/20240418_mkt_plano.md      │
│                                     │
│  > Ver prompt                       │
│    Novo prompt                      │
│    Sair                             │
└─────────────────────────────────────┘
```

---

## Estrutura do Projeto

```
prompt-builder/
├── main.go
├── go.mod
├── context.md
├── roles/
│   ├── marketing_content.md
│   ├── data_analyst.md
│   └── devops_engineer.md
├── models/
│   ├── rtf.md
│   ├── race.md
│   ├── risen.md
│   └── create.md
├── prompts/                          ← gerados pelo programa
│   └── 20240418_marketing_content_plano-editorial.md
└── internal/
    ├── parser/                       ← leitura dos .md de roles e models
    ├── engine/                       ← template engine ({{campo}}, {{#if}}, etc)
    └── tui/                          ← telas bubbletea
```

---

## Decisões Tomadas

- [x] Parser do frontmatter: `gopkg.in/yaml.v3`
- [x] Interface: TUI com `bubbletea` + `bubbles` + `lipgloss`
- [x] Navegação: wizard linear com voltar (Esc)
- [x] Gaps não respondidos: entram no prompt como seção `⚠️` (visível pra IA também)
- [x] Nome do arquivo salvo: gerado automaticamente ou usuário digita título?

## Diretrizes

- Inicialize git.
- NUNCA ao fazer commit colocar o claude como autor;
- Sempre usar o context-mode em todos os processos;
- Sempre usar caveman;
- Crie um mascote com bubbletea, um CASTOR. O nome do projeto á CASTOR BUILDER;
- Se perceber que as features são grandes crie fases;
- Ao fim de cada fase ou conclusão de atividade faça commit usando git semântico


