---
id: rtf
nome: RTF
descricao: Ideal for simple and direct tasks
campos:
  - id: role
    label: Role
    tipo: text
    obrigatorio: true

  - id: task
    label: Task
    tipo: textarea
    obrigatorio: true

  - id: format
    label: Output format
    tipo: select
    opcoes:
      - bullet list
      - prose
      - table
      - code
    obrigatorio: false

  - id: premissas
    label: Assumptions and conditions
    tipo: list
    obrigatorio: true
---
## Output template
You are {{role}}.
{{task}}
{{#if format}}
Present the result as {{format}}.
{{/if}}
{{#if premissas}}
Consider the following assumptions:
{{#each premissas}}- {{.}}
{{/each}}
{{/if}}
