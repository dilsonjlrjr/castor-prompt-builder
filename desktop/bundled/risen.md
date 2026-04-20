---
id: risen
nome: RISEN
descricao: Ideal quando precisa de steps detalhados com restrições
campos:
  - id: role
    label: Papel
    tipo: text
    obrigatorio: true

  - id: input
    label: Input / Contexto
    tipo: textarea
    obrigatorio: true

  - id: fases
    label: Steps de execução
    tipo: steps
    obrigatorio: false
    step_campos:
      - id: titulo
        label: Título do step
        tipo: text
      - id: descricao
        label: O que deve ser feito
        tipo: textarea

  - id: expectation
    label: Expectativa de entrega
    tipo: textarea
    obrigatorio: true

  - id: narrowing
    label: Restrições e limitações
    tipo: list
    obrigatorio: false
---

## Template de saída

Você é {{role}}.

## Contexto
{{input}}

{{#steps fases}}
## {{titulo}}
{{descricao}}
{{/steps}}

## Expectativa
{{expectation}}

{{#if narrowing}}
## Restrições
{{#each narrowing}}- {{.}}
{{/each}}
{{/if}}
