---
id: race
nome: RACE
descricao: Ideal for tasks with rich context and a clear delivery expectation
campos:
  - id: role
    label: Role
    tipo: text
    obrigatorio: true

  - id: action
    label: Action
    tipo: textarea
    obrigatorio: true

  - id: context
    label: Context
    tipo: textarea
    obrigatorio: true

  - id: tom
    label: Communication tone
    tipo: select
    opcoes:
      - formal
      - informal
      - technical
      - persuasive
    obrigatorio: false

  - id: canais
    label: Distribution channels
    tipo: multiselect
    opcoes:
      - blog
      - LinkedIn
      - email
      - Instagram
    obrigatorio: false

  - id: fases
    label: Execution phases
    tipo: steps
    obrigatorio: false
    step_campos:
      - id: titulo
        label: Phase title
        tipo: text
      - id: descricao
        label: What should happen in this phase
        tipo: textarea

  - id: expectation
    label: Expectation
    tipo: textarea
    obrigatorio: true

  - id: premissas
    label: Assumptions and conditions
    tipo: list
    obrigatorio: true
---
## Output template
You are {{role}}.
{{action}}
The context is the following: {{context}}
{{#if tom}}
Adopt a {{tom}} tone.
{{/if}}
{{#steps fases}}
## {{titulo}}
{{descricao}}
{{/steps}}
{{#if canais}}
Consider the following channels: {{#each canais}}{{.}}{{/each}}.
{{/if}}
{{#if premissas}}
## Assumptions
{{#each premissas}}- {{.}}
{{/each}}
{{/if}}
The expected outcome is that {{expectation}}
