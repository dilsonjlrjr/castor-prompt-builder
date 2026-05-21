---
id: risen
nome: RISEN
descricao: Ideal when you need detailed steps with constraints
campos:
  - id: role
    label: Role
    tipo: text
    obrigatorio: true

  - id: input
    label: Input / Context
    tipo: textarea
    obrigatorio: true

  - id: fases
    label: Execution steps
    tipo: steps
    obrigatorio: false
    step_campos:
      - id: titulo
        label: Step title
        tipo: text
      - id: descricao
        label: What should happen
        tipo: textarea

  - id: expectation
    label: Delivery expectation
    tipo: textarea
    obrigatorio: true

  - id: narrowing
    label: Constraints and limits
    tipo: list
    obrigatorio: false

  - id: premissas
    label: Assumptions and conditions
    tipo: list
    obrigatorio: true
---
## Output template
You are {{role}}.
## Context
{{input}}
{{#steps fases}}
## {{titulo}}
{{descricao}}
{{/steps}}
## Expectation
{{expectation}}
{{#if narrowing}}
## Constraints
{{#each narrowing}}- {{.}}
{{/each}}
{{/if}}
{{#if premissas}}
## Assumptions
{{#each premissas}}- {{.}}
{{/each}}
{{/if}}
