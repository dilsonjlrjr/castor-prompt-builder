---
id: race
nome: RACE
descricao: Ideal para tarefas com contexto rico e expectativa clara de entrega
campos:
  - id: role
    label: Papel
    tipo: text
    obrigatorio: true

  - id: action
    label: Ação
    tipo: textarea
    obrigatorio: true

  - id: context
    label: Contexto
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
    obrigatorio: false

  - id: canais
    label: Canais de distribuição
    tipo: multiselect
    opcoes:
      - blog
      - LinkedIn
      - email
      - Instagram
    obrigatorio: false

  - id: fases
    label: Fases de execução
    tipo: steps
    obrigatorio: false
    step_campos:
      - id: titulo
        label: Título da fase
        tipo: text
      - id: descricao
        label: O que deve ser feito nesta fase
        tipo: textarea

  - id: expectation
    label: Expectativa
    tipo: textarea
    obrigatorio: true
---

## Template de saída

Você é {{role}}.

{{action}}

O contexto é o seguinte: {{context}}

{{#if tom}}
Adote um tom {{tom}}.
{{/if}}

{{#steps fases}}
## {{titulo}}
{{descricao}}
{{/steps}}

{{#if canais}}
Considere os seguintes canais: {{#each canais}}{{.}}{{/each}}.
{{/if}}

Espera-se que {{expectation}}
