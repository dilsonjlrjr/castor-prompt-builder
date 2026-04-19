# 🦫 CASTOR BUILDER

> CLI interativo para construir prompts estruturados e semanticamente ricos — sem IA no processo.

---

## O que é

CASTOR BUILDER é uma ferramenta de linha de comando com interface TUI (Terminal UI) que guia você na criação de prompts para LLMs. Você escolhe um **modelo de estrutura**, seleciona um **papel profissional**, descreve a tarefa em linguagem livre e o CASTOR monta um prompt pronto para uso.

Nenhuma chamada de API. Nenhuma dependência de IA. Tudo é parsing, heurística e template engine próprio.

---

## Instalação

```bash
git clone https://github.com/dilsonrabelo/castor-prompt-builder
cd castor-prompt-builder
go build -o castor .
./castor
```

**Requisitos:** Go 1.21+

---

## Como usar

Ao rodar `./castor`, um wizard interativo é iniciado:

```
[1] Selecionar modelo    →  RTF, RACE, RISEN, CREATE
[2] Selecionar papel     →  banco de especialistas em roles/
[3] Inserir narrativa    →  descreva a tarefa em linguagem livre
[4] Preencher lacunas    →  campos obrigatórios não detectados viram perguntas
[5] Definir fases?       →  opcional — divide a execução em steps
[6] Prompt gerado        →  salvo automaticamente em prompts/
```

### Navegação

| Tecla | Ação |
|---|---|
| `↑` / `↓` | Navegar entre opções |
| `Enter` | Confirmar / avançar |
| `Ctrl+S` | Confirmar textarea |
| `Esc` | Voltar tela anterior |
| `q` | Sair |

---

## Modelos disponíveis

| Modelo | Quando usar |
|---|---|
| **RTF** | Tarefas simples e diretas — Role, Task, Format |
| **RACE** | Contexto rico com entregável claro — Role, Action, Context, Expectation |
| **RISEN** | Steps detalhados com restrições — Role, Input, Steps, Expectation, Narrowing |
| **CREATE** | Conteúdo criativo com público e tom — Context, Role, Examples, Audience, Tone, Expectation |

---

## Prompt gerado — exemplo de saída

```markdown
# Prompt — Especialista em Marketing de Conteúdo
_Modelo: RACE | Gerado em: 2024-04-18_

---

Você é um Especialista em Marketing de Conteúdo. Profissional com sólida
experiência em estratégia editorial B2B e B2C...

Crie um plano editorial para os próximos 3 meses.

O contexto é o seguinte: Startup B2B de SaaS com queda de 15% no engajamento.

Adote um tom analítico e estratégico.

## Fase 1 — Diagnóstico
Analise as possíveis causas da queda...

## Expectativa
Documento com temas, frequência e canais por mês.

---
> ⚠️ Para refinar este prompt considere informar:
> - Orçamento disponível para produção de conteúdo?
```

Os prompts são salvos em `prompts/` com nomenclatura automática:

```
prompts/20240418_marketing_content_plano-editorial-3-meses.md
```

---

## Estrutura do projeto

```
castor-prompt-builder/
├── main.go
├── assets/
│   └── castor.png          ← mascote (opcional, substitui o fallback)
├── roles/                  ← banco de papéis profissionais
│   ├── marketing_content.md
│   ├── data_analyst.md
│   └── devops_engineer.md
├── models/                 ← modelos de estrutura de prompt
│   ├── rtf.md
│   ├── race.md
│   ├── risen.md
│   └── create.md
├── prompts/                ← gerados automaticamente
└── internal/
    ├── parser/             ← leitura dos .md (frontmatter YAML + corpo)
    ├── engine/             ← template engine ({{campo}}, {{#if}}, {{#steps}})
    └── tui/                ← wizard bubbletea + mascote chafa-go
```

---

## Adicionando papéis

Crie um arquivo em `roles/` seguindo o formato:

```markdown
---
id: nome_do_papel
nome: Nome exibido na lista
tom: estilo de escrita
habilidades:
  - habilidade 1
  - habilidade 2
gaps_comuns:
  - Qual o público-alvo?
  - Existe prazo definido?
---

Descrição rica do papel — vai direto no prompt gerado.
```

---

## Adicionando modelos

Crie um arquivo em `models/` seguindo o formato:

```markdown
---
id: meu_modelo
nome: Meu Modelo
descricao: Quando usar este modelo
campos:
  - id: campo1
    label: Nome do campo
    tipo: textarea        # text | textarea | select | multiselect | list | steps
    obrigatorio: true
---

## Template de saída

Você é {{role}}.

{{campo1}}

{{#if campo_opcional}}
Considere: {{campo_opcional}}
{{/if}}
```

### Tipos de campo

| Tipo | Comportamento |
|---|---|
| `text` | Linha única |
| `textarea` | Multi-linha (Ctrl+S confirma) |
| `select` | Lista numerada, escolha uma opção |
| `multiselect` | Lista numerada, escolha várias |
| `list` | Adiciona itens livres um a um |
| `steps` | Coleta título + descrição por fase |

### Diretivas de template

| Diretiva | Comportamento |
|---|---|
| `{{campo}}` | Substitui pelo valor do campo |
| `{{#if campo}}...{{/if}}` | Renderiza só se preenchido |
| `{{#steps campo}}...{{/steps}}` | Itera fases — expõe `{{titulo}}` e `{{descricao}}` |
| `{{#each campo}}...{{/each}}` | Itera lista — expõe `{{.}}` |

---

## Mascote

O CASTOR usa [chafa-go](https://github.com/ploMP4/chafa-go) para renderizar pixel art no terminal com true color. Para usar uma imagem personalizada, coloque um arquivo PNG em `assets/castor.png`. Se o arquivo não existir, o fallback em lipgloss é usado automaticamente.

---

## Stack

| Lib | Papel |
|---|---|
| [bubbletea](https://github.com/charmbracelet/bubbletea) | Engine TUI (arquitetura Elm) |
| [bubbles](https://github.com/charmbracelet/bubbles) | Componentes: list, textarea, textinput |
| [lipgloss](https://github.com/charmbracelet/lipgloss) | Estilos, cores, layout |
| [chafa-go](https://github.com/ploMP4/chafa-go) | Renderização de pixel art no terminal |
| [yaml.v3](https://gopkg.in/yaml.v3) | Parser de frontmatter |

---

## Licença

MIT
