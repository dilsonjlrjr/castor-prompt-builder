---
id: risen
nome: RISEN
descricao: Ideal cuando necesitas pasos detallados con restricciones
campos:
  - id: role
    label: Rol
    tipo: text
    obrigatorio: true

  - id: input
    label: Input / Contexto
    tipo: textarea
    obrigatorio: true

  - id: fases
    label: Pasos de ejecución
    tipo: steps
    obrigatorio: false
    step_campos:
      - id: titulo
        label: Título del paso
        tipo: text
      - id: descricao
        label: Qué debe ocurrir
        tipo: textarea

  - id: expectation
    label: Expectativa de entrega
    tipo: textarea
    obrigatorio: true

  - id: narrowing
    label: Restricciones y límites
    tipo: list
    obrigatorio: false

  - id: premissas
    label: Premisas y condiciones
    tipo: list
    obrigatorio: true
---
## Plantilla de salida
Eres {{role}}.
## Contexto
{{input}}
{{#steps fases}}
## {{titulo}}
{{descricao}}
{{/steps}}
## Expectativa
{{expectation}}
{{#if narrowing}}
## Restricciones
{{#each narrowing}}- {{.}}
{{/each}}
{{/if}}
{{#if premissas}}
## Premisas
{{#each premissas}}- {{.}}
{{/each}}
{{/if}}
