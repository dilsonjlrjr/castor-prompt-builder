---
id: race
nome: RACE
descricao: Ideal para tareas con contexto rico y expectativa clara de entrega
campos:
  - id: role
    label: Rol
    tipo: text
    obrigatorio: true

  - id: action
    label: Acción
    tipo: textarea
    obrigatorio: true

  - id: context
    label: Contexto
    tipo: textarea
    obrigatorio: true

  - id: tom
    label: Tono de comunicación
    tipo: select
    opcoes:
      - formal
      - informal
      - técnico
      - persuasivo
    obrigatorio: false

  - id: canais
    label: Canales de distribución
    tipo: multiselect
    opcoes:
      - blog
      - LinkedIn
      - email
      - Instagram
    obrigatorio: false

  - id: fases
    label: Fases de ejecución
    tipo: steps
    obrigatorio: false
    step_campos:
      - id: titulo
        label: Título de la fase
        tipo: text
      - id: descricao
        label: Qué debe ocurrir en esta fase
        tipo: textarea

  - id: expectation
    label: Expectativa
    tipo: textarea
    obrigatorio: true

  - id: premissas
    label: Premisas y condiciones
    tipo: list
    obrigatorio: true
---
## Plantilla de salida
Eres {{role}}.
{{action}}
El contexto es el siguiente: {{context}}
{{#if tom}}
Adopta un tono {{tom}}.
{{/if}}
{{#steps fases}}
## {{titulo}}
{{descricao}}
{{/steps}}
{{#if canais}}
Considera los siguientes canales: {{#each canais}}{{.}}{{/each}}.
{{/if}}
{{#if premissas}}
## Premisas
{{#each premissas}}- {{.}}
{{/each}}
{{/if}}
Se espera que {{expectation}}
