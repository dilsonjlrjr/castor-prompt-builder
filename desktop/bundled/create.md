---
id: create
nome: CREATE
descricao: Conteúdo criativo com restrições de público e tom
campos:
  - id: context
    label: Contexto
    tipo: textarea
    obrigatorio: true

  - id: role
    label: Papel
    tipo: text
    obrigatorio: true

  - id: examples
    label: Exemplos de referência
    tipo: list
    obrigatorio: false

  - id: audience
    label: Público-alvo
    tipo: textarea
    obrigatorio: true

  - id: tom
    label: Tom de comunicação
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
---

## Template de saída

Você é {{role}}.

## Contexto
{{context}}

## Público-alvo
{{audience}}

{{#if tom}}
Adote um tom {{tom}}.
{{/if}}

{{#if examples}}
## Referências
{{#each examples}}- {{.}}
{{/each}}
{{/if}}

## Entregável esperado
{{expectation}}
