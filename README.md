<p align="center">
  <img src="assets/castor.png" width="140" alt="Castor mascot" />
</p>

<h1 align="center">CASTOR Builder</h1>

<p align="center">
  <strong>Build structured prompts for LLMs — no AI required.</strong><br/>
  Choose a role. Describe your task. Get a production-ready prompt.
</p>

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat-square&logo=go&logoColor=white" />
  <img src="https://img.shields.io/badge/Wails-v2-a371f7?style=flat-square" />
  <img src="https://img.shields.io/badge/platform-macOS%20%7C%20Windows%20%7C%20Linux-f5a623?style=flat-square" />
  <img src="https://img.shields.io/badge/license-MIT-3fb950?style=flat-square" />
</p>

---

## What is this?

CASTOR Builder is a **wizard-driven prompt engineering tool** that turns a free-text task description into a structured, semantically rich LLM prompt — entirely offline, no AI API calls.

Pick a **role** (e.g. "Senior Backend Engineer in Go"), describe what you need, answer a few context questions, optionally define execution phases, and get a `.md` file ready to paste into any LLM.

**Two flavors:**

| Version | When to use |
|---|---|
| **Desktop** (Wails v2) | Visual UI, cross-platform native app |
| **TUI** (bubbletea) | Lives in your terminal, scriptable, zero friction |

---

## Features

- **47 roles** across 12 professional categories
- **4 proven prompt models** (RTF, RACE, RISEN, CREATE)
- **Gap detection** — unfilled fields become interactive questions
- **Multi-role** — combine roles; skills and tone merge automatically
- **Execution phases** — structure prompts into numbered steps
- **Zero AI** — pure Go template engine, fully offline
- **Cross-platform** — macOS universal, Windows amd64, Linux
- **Extensible** — add your own roles and models as plain `.md` files

---

## Quick Start

### Desktop (GUI)

Download the latest release for your platform from the [Releases](https://github.com/dilsonrabelo/castor-prompt-builder/releases) page.

| Platform | File |
|---|---|
| macOS (Universal) | `CASTOR-Builder-darwin-universal.dmg` |
| Windows | `CASTOR-Builder-windows-amd64.exe` |
| Linux | `CASTOR-Builder-linux-amd64` |

On first launch, bundled models and roles are extracted to `~/.castorprompt/`.

### TUI (Terminal)

```bash
git clone https://github.com/dilsonrabelo/castor-prompt-builder
cd castor-prompt-builder
go build -o castor-tui ./cmd/tui
./castor-tui
```

> Requires Go 1.21+

---

## The Wizard — 6 Steps

```
1. Select model    → RTF / RACE / RISEN / CREATE
2. Select role(s)  → search by name or category, multi-select
3. Write narrative → free-text description of your task
4. Fill gaps       → unanswered required fields become questions
5. Define phases?  → optionally structure into numbered steps
6. Done            → prompt saved to prompts/<date>_<role>_<slug>.md
```

### Keyboard shortcuts (TUI)

| Key | Action |
|---|---|
| `↑` / `↓` or `k` / `j` | Navigate |
| `Enter` | Confirm / advance |
| `Space` | Toggle role selection |
| `Ctrl+S` | Confirm textarea |
| `Esc` | Go back |
| `q` | Quit |

---

## Prompt Models

### RTF — Role, Task, Feature
Simple and direct. Best for focused, single-output tasks.

### RACE — Role, Action, Context, Expectation
Rich context with a clear deliverable. Best for complex briefs.

### RISEN — Role, Input, Steps, Expectation, Narrowing
Explicit step-by-step reasoning with constraints. Best for long tasks.

### CREATE — Context, Role, Examples, Audience, Tone, Expectation
Creative output with audience and tone control. Best for content creation.

---

## Role Categories

| Category | Roles |
|---|---|
| **Arquitetura** | Software Architect, Solutions Architect, Cloud Architect, Microservices Architect |
| **Backend** | Node.js, Go, Python, Java |
| **Frontend** | React, Vue |
| **Mobile** | React Native, Flutter |
| **DevOps / SRE** | SRE Engineer, Platform Engineer, Cloud AWS Engineer |
| **Banco de Dados** | SQL DBA, NoSQL DBA |
| **Dados** | Data Architect, Data Engineer, Data Scientist, ML Engineer, BI Analyst, Analytics Engineer |
| **Design** | UX Designer, Product Designer |
| **Gestão** | Scrum Master, Product Owner, Tech Lead, Product Manager, Business Analyst |
| **Segurança** | Security Engineer, Pentester, Compliance Analyst |
| **QA** | QA Engineer, QA Automation, QA Lead |
| **Marketing** | Content Specialist, Copywriter, SEO Specialist, Growth Analyst |
| **Documentação** | Technical Writer |
| **Especialistas** | Go, Java, .NET, Python |

---

## Adding Custom Roles

Create a `.md` file in `~/.castorprompt/roles/<category>/my_role.md`:

```markdown
---
id: my_role
nome: My Custom Role
categoria: custom
tom: analytical and direct
habilidades:
  - skill one
  - skill two
gaps_comuns:
  - What is the target audience?
  - What is the expected output format?
---

Free-text description of this role. Goes directly into the generated prompt.
```

Restart the app — your role appears in the list immediately.

---

## Adding Custom Models

Create a `.md` file in `~/.castorprompt/models/my_model.md`:

```markdown
---
id: my_model
nome: MY MODEL
descricao: Short description of when to use this model
campos:
  - id: role
    label: Role
    tipo: text
    obrigatorio: true
  - id: context
    label: Context
    tipo: textarea
    obrigatorio: true
---

You are {{role}}.

{{context}}

{{#if tom}}
Use a {{tom}} tone.
{{/if}}
```

### Template directives

| Directive | Behavior |
|---|---|
| `{{field}}` | Plain field substitution |
| `{{#if field}}...{{/if}}` | Render block only if field is filled |
| `{{#steps field}}...{{/steps}}` | Iterate execution phases — exposes `{{titulo}}` and `{{descricao}}` |
| `{{#each field}}...{{/each}}` | Iterate list/multiselect — exposes `{{.}}` for current item |

### Field types

| Type | UI behavior | Use for |
|---|---|---|
| `text` | Single line input | Title, name, tone |
| `textarea` | Multi-line editor | Narrative, context, description |
| `select` | Numbered list, pick one | Output format, communication tone |
| `multiselect` | Numbered list, pick many | Channels, audiences |
| `list` | Free items, one per line | Constraints, extra skills |
| `steps` | Title + description per phase | Execution phases |

---

## Building from Source

### Desktop (Wails v2)

```bash
# Install Wails CLI
go install github.com/wailsapp/wails/v2/cmd/wails@latest

# Production build
cd desktop && wails build

# Development mode (hot reload)
cd desktop && wails dev
```

### TUI

```bash
go build -o castor-tui ./cmd/tui
./castor-tui
```

---

## Generated Prompt — Example Output

```markdown
# Prompt — Senior Go Engineer
_Model: RACE | Generated: 2024-04-18_

---

You are a Senior Go Engineer. Expert in distributed systems, REST/gRPC APIs
and high-performance backends in Go...

Refactor the authentication service to use JWT with refresh tokens.

The context is: Startup SaaS B2B. Currently uses sessions stored in Redis.
Target: stateless auth, 10k concurrent users.

---
## Relevant skills
- Distributed systems
- JWT / OAuth 2.0
- Redis, PostgreSQL

---
## Communication tone
- Senior Go Engineer: technical and direct

---
## Execution phases

### Phase 1 — Audit
Map existing session logic and identify coupling points.

### Phase 2 — Implementation
Replace session middleware with JWT middleware, implement refresh flow.

### Phase 3 — Testing
Load test with k6, validate 10k concurrent connections.
```

Saved as:
```
prompts/20240418_especialista-go_refactor-auth-jwt.md
```

---

## Project Structure

```
castor-prompt-builder/
├── cmd/tui/            ← TUI entry point
├── desktop/            ← Wails v2 desktop app
│   ├── app.go          ← Go backend (BuildPrompt, LoadModels, LoadRoles)
│   ├── bundled/        ← models + roles embedded in binary
│   └── frontend/       ← Svelte + TypeScript + Tailwind UI
├── internal/
│   ├── tui/            ← bubbletea screens (Model / Update / View)
│   ├── engine/         ← template renderer
│   └── parser/         ← .md frontmatter parser (yaml.v3)
├── models/             ← prompt model definitions (4 built-in)
├── roles/              ← role definitions (47 roles, 12 categories)
└── prompts/            ← generated prompts output (gitignored)
```

---

## Tech Stack

| Component | Library |
|---|---|
| Desktop framework | [Wails v2](https://wails.io) |
| Frontend | Svelte + TypeScript + Tailwind CSS |
| TUI engine | [bubbletea](https://github.com/charmbracelet/bubbletea) |
| TUI components | [bubbles](https://github.com/charmbracelet/bubbles) |
| TUI styling | [lipgloss](https://github.com/charmbracelet/lipgloss) |
| YAML parsing | `gopkg.in/yaml.v3` |
| Unicode slugify | `golang.org/x/text` |

---

## License

MIT — see [LICENSE](LICENSE).

---

<p align="center">
  Made with 🦫 and Go &nbsp;·&nbsp;
  <a href="https://github.com/dilsonrabelo/castor-prompt-builder">GitHub</a>
</p>
