---
id: create
nome: CREATE
descricao: Contenido creativo con restricciones de público y tono
campos:
  - id: context
    label: Contexto
    tipo: textarea
    obrigatorio: true

  - id: role
    label: Rol
    tipo: text
    obrigatorio: true

  - id: examples
    label: Ejemplos de referencia
    tipo: list
    obrigatorio: false

  - id: audience
    label: Público objetivo
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
      - inspirador
    obrigatorio: false

  - id: expectation
    label: Expectativa de entrega
    tipo: textarea
    obrigatorio: true

  - id: premissas
    label: Premisas y condiciones
    tipo: list
    obrigatorio: true
---
## Plantilla de salida
Eres {{role}}.
## Contexto
{{context}}
## Público objetivo
{{audience}}
{{#if tom}}
Adopta un tono {{tom}}.
{{/if}}
{{#if examples}}
## Referencias
{{#each examples}}- {{.}}
{{/each}}
{{/if}}
{{#if premissas}}
## Premisas
{{#each premissas}}- {{.}}
{{/each}}
{{/if}}
## Entregable esperado
{{expectation}}
