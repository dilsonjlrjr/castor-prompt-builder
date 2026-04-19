---
id: rtf
nome: RTF
descricao: Ideal para tarefas simples e diretas
campos:
  - id: role
    label: Papel
    tipo: text
    obrigatorio: true

  - id: task
    label: Tarefa
    tipo: textarea
    obrigatorio: true

  - id: format
    label: Formato de saída
    tipo: select
    opcoes:
      - lista com marcadores
      - texto corrido
      - tabela
      - código
    obrigatorio: false
---

## Template de saída

Você é {{role}}.

{{task}}

{{#if format}}
Apresente o resultado em formato de {{format}}.
{{/if}}
