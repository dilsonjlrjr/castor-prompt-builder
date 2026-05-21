// Dicionário PT-BR — fonte de verdade.
// Mantenha a estrutura espelhada em en.ts e es.ts.
export const pt = {
  app: {
    name: 'CASTOR Builder',
    tagline: 'Prompt Builder',
    version: 'v1.2.0',
  },

  sidebar: {
    savedPrompts: 'Prompts salvos',
    tutorial:     'Ver tutorial',
    languageHint: 'Idioma',
  },

  steps: {
    model:     'Modelo',
    role:      'Papel',
    narrative: 'Tarefa',
    gap:       'Contexto',
    phase:     'Fases',
    result:    'Resultado',
  },

  topbar: {
    back: 'Voltar',
  },

  validation: {
    title:         'Validação de arquivos',
    statusRunning: 'Verificando arquivos...',
    statusOk:      'Tudo certo!',
    statusErrors:  '{n} problema{s} encontrado{s}',
    counter:       '{done} / {total} arquivos',
    badgeOk:       '✓ sem problemas',
    badgeErrors:   '✗ {n} com erro',
    continueOk:    'Continuar →',
    continueErr:   'Continuar mesmo assim →',
  },

  tutorial: {
    skip: 'Pular tutorial',
    prev: '← Anterior',
    next: 'Próximo →',
    cta:  '🚀 Começar a criar!',
    slides: {
      overview:  { title: 'Como funciona o CASTOR', subtitle: 'Visão geral · 6 etapas',
                   desc: 'O CASTOR guia você por um wizard simples. Cada etapa adiciona uma camada de contexto — no final, um prompt estruturado pronto para qualquer IA.',
                   tip:  'Você pode voltar a etapas anteriores a qualquer momento.' },
      model:     { title: 'Escolha o Modelo', subtitle: 'Passo 1 de 6',
                   desc: 'O modelo define a estrutura e a lógica do prompt. Cada um serve um tipo diferente de tarefa — escolha o que melhor se encaixa no seu objetivo.',
                   tip:  'Na dúvida, RACE é o mais versátil para a maioria dos casos.' },
      roles:     { title: 'Selecione os Papéis', subtitle: 'Passo 2 de 6',
                   desc: 'O papel define quem a IA deve "ser" para responder sua tarefa. Você pode combinar vários — o prompt incorpora as habilidades de todos.',
                   tip:  'Combinar papéis complementares (ex: Arquiteto + DevOps) enriquece muito o resultado.' },
      narrative: { title: 'Descreva sua Tarefa', subtitle: 'Passo 3 de 6',
                   desc: 'Escreva naturalmente, como explicaria a um colega. Sem formatação especial — o CASTOR distribui o contexto automaticamente nos campos certos.',
                   tip:  'Quanto mais contexto aqui, menos gaps serão perguntados depois.' },
      gaps:      { title: 'Preencha o Contexto', subtitle: 'Passo 4 de 6',
                   desc: 'O CASTOR identifica o que falta e faz perguntas direcionadas. Algumas vêm do modelo, outras dos papéis — e cada pergunta mostra de onde veio.',
                   tip:  'Campos opcionais podem ser pulados — aparecem no prompt como lacunas a considerar.' },
      phases:    { title: 'Defina as Fases', subtitle: 'Passo 5 de 6 · opcional',
                   desc: 'Para tarefas complexas, divida a execução em etapas sequenciais. Cada fase tem um título e uma descrição do que deve acontecer naquele momento.',
                   tip:  'Fases são ótimas para projetos de múltiplas entregas ou raciocínio encadeado.' },
      result:    { title: 'Prompt Pronto!', subtitle: 'Passo 6 de 6',
                   desc: 'O CASTOR monta o prompt com modelo, papéis, habilidades, contexto e fases — estruturado e pronto para usar em qualquer IA.',
                   tip:  'Copie e cole diretamente no ChatGPT, Claude, Gemini ou qualquer outra IA.' },
    },
  },

  resume: {
    message:  'Você tem um progresso salvo. Deseja continuar?',
    confirm:  'Sim, continuar',
    discard:  'Não, começar de novo',
  },

  model: {
    title:    'Escolha o framework',
    subtitle: 'Cada modelo estrutura o prompt de uma forma diferente. Escolha o que melhor se encaixa na sua tarefa.',
  },

  role: {
    title:        'Selecione o(s) papel(eis)',
    subtitle:     'Você pode combinar múltiplos papéis para um prompt mais rico.',
    searchPlaceholder: 'Buscar papel ou categoria...',
    selectedCount:'{n} selecionado{s}',
    continue:     'Continuar →',
  },

  narrative: {
    title:       'Descreva sua tarefa',
    subtitle:    'Escreva livremente. Quanto mais contexto, melhor o prompt gerado.',
    placeholder: 'Ex: Preciso criar um plano editorial para os próximos 3 meses. A empresa é uma startup B2B de SaaS que está com queda de engajamento no blog. O público-alvo são desenvolvedores e CTOs...',
    charCount:   '{n} caracteres',
    continue:    'Continuar →',
  },

  gap: {
    title:        'Contexto adicional',
    subtitle:     'Preencha as lacunas — ★ obrigatório, ✓ preenchido',
    expanded:     '▼ Expandido',
    collapsed:    '▶ Compacto',
    toggleEditor: '◧ Ambos',
    toggleEditorAlt: '◨ Editor',
    toggleEditorTitle:    'Mostrar lista de lacunas',
    toggleEditorAltTitle: 'Expandir área de resposta',
    prev:        '← Anterior',
    skipAll:     '⏭ Pular tudo',
    skipBlocked: 'Preencha os campos obrigatórios primeiro',
    skipTitle:   'Pular todas as lacunas',
    continue:    'Continuar →',
    next:        'Próxima →',
    sectionModel: 'Contexto do Modelo ({name})',
    noPending:    'Nenhuma lacuna pendente',
    filledOf:     '{done}/{total} preenchidos',
  },

  phase: {
    title:       'Fases de execução',
    subtitle:    'Divida o trabalho em etapas para um prompt mais estruturado.',
    enable:      'Definir fases',
    enableYes:   'Sim, quero definir fases',
    enableNo:    'Não, gerar direto',
    skip:        'Pular esta etapa',
    addPhase:    '+ Adicionar fase',
    phaseTitle:  'Título da fase',
    phaseDesc:   'O que deve acontecer nessa fase?',
    label:       'Fase {n}',
    removePhase: 'remover',
    generate:    'Gerar Prompt',
  },

  result: {
    success:    'Prompt gerado com sucesso!',
    errorTitle: 'Erro ao gerar prompt',
    save:       '💾 Salvar',
    saved:      '✓ Salvo!',
    copy:       '⎘ Copiar',
    copied:     '✓ Copiado!',
    restart:    '← Criar novo prompt',
    building:   'Gerando...',
  },

  prompts: {
    title:           'Prompts salvos',
    subtitle:        'Gerencie seus prompts gerados anteriormente.',
    counterOne:      '{n} prompt',
    counterMany:     '{n} prompts',
    searchPlaceholder: 'Buscar por título…',
    emptyTitle:      'Nenhum prompt salvo',
    emptyDesc:       'Salve prompts na tela de resultado para reaproveitá-los depois.',
    noSearchResults: 'Nenhum resultado para "{q}"',
    selectToView:    'Selecione um prompt para visualizar',
    loading:         'Carregando...',
    untitled:        'Sem título',
    deleteConfirm:   'Confirmar',
    deleteCancel:    'Cancelar',
    deleteTooltip:   'Excluir',
  },

  relativeDate: {
    now:    'agora',
    minute: 'há {n} min',
    hour:   'há {n} h',
    day:    'há {n} d',
  },

  // Rótulos amigáveis das categorias de papéis (chave = slug do diretório).
  categories: {
    arquitetura:  'Arquitetura',
    frontend:     'Frontend & Mobile',
    backend:      'Backend',
    devops:       'DevOps & Cloud',
    banco:        'Banco de Dados',
    dados:        'Dados & IA',
    gestao:       'Gestão',
    seguranca:    'Segurança',
    design:       'Design',
    marketing:    'Marketing',
    qa:           'QA & Testes',
    documentacao: 'Documentação',
    jornalismo:   'Jornalismo',
    direito:      'Direito',
    medicina:     'Medicina',
    veterinaria:  'Medicina Veterinária',
    educacao:     'Educação',
    financas:     'Finanças',
    rh:           'Recursos Humanos',
    engenharia:   'Engenharia',
    psicologia:   'Psicologia',
    saude:        'Saúde',
    vendas:       'Vendas',
  },
} as const

export type Dict = typeof pt
