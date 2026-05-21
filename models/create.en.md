---
id: create
nome: CREATE
descricao: Creative content with audience and tone constraints
campos:
  - id: context
    label: Context
    tipo: textarea
    obrigatorio: true

  - id: role
    label: Role
    tipo: text
    obrigatorio: true

  - id: examples
    label: Reference examples
    tipo: list
    obrigatorio: false

  - id: audience
    label: Target audience
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
      - inspiring
    obrigatorio: false

  - id: expectation
    label: Delivery expectation
    tipo: textarea
    obrigatorio: true

  - id: premissas
    label: Assumptions and conditions
    tipo: list
    obrigatorio: true
---
## Output template
You are {{role}}.
## Context
{{context}}
## Target audience
{{audience}}
{{#if tom}}
Adopt a {{tom}} tone.
{{/if}}
{{#if examples}}
## References
{{#each examples}}- {{.}}
{{/each}}
{{/if}}
{{#if premissas}}
## Assumptions
{{#each premissas}}- {{.}}
{{/each}}
{{/if}}
## Expected deliverable
{{expectation}}
