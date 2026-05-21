---
id: rtf
nome: RTF
descricao: Ideal para tareas simples y directas
campos:
  - id: role
    label: Rol
    tipo: text
    obrigatorio: true

  - id: task
    label: Tarea
    tipo: textarea
    obrigatorio: true

  - id: format
    label: Formato de salida
    tipo: select
    opcoes:
      - lista con viñetas
      - texto corrido
      - tabla
      - código
    obrigatorio: false

  - id: premissas
    label: Premisas y condiciones
    tipo: list
    obrigatorio: true
---
## Plantilla de salida
Eres {{role}}.
{{task}}
{{#if format}}
Presenta el resultado en formato de {{format}}.
{{/if}}
{{#if premissas}}
Considera las siguientes premisas:
{{#each premissas}}- {{.}}
{{/each}}
{{/if}}
