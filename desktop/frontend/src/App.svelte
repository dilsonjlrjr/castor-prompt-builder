<script lang="ts">
  import { onMount } from 'svelte'
  import { fly, fade } from 'svelte/transition'
  import { cubicOut } from 'svelte/easing'
  import { GetModels, GetRoles, BuildPrompt, IsFirstRun, ValidateAll, SavePrompt, ListPrompts, DeletePrompt, GetPrompt } from '../wailsjs/go/main/App.js'
  import { main } from '../wailsjs/go/models'

  // ---- tipos ----
  type Campo  = { id: string; label: string; tipo: string; obrigatorio: boolean; opcoes?: string[] }
  type Model  = { id: string; nome: string; descricao: string; campos: Campo[] }
  type Role   = { id: string; nome: string; categoria: string; gaps_comuns: string[] }
  type Step   = { titulo: string; descricao: string }
  type Screen = 'model' | 'role' | 'narrative' | 'gap' | 'phase' | 'result' | 'prompts'

  // Gap unificado: campo do modelo (tem fieldId) ou gap de papel (fieldId = '')
  type Gap = { fieldId: string; pergunta: string; tipo: string; opcoes: string[]; obrigatorio: boolean; roleNome?: string; answer: string }
  // ContextSection agrupa gaps — do modelo ou de um papel
  type ContextSection = { title: string; kind: string; gaps: Gap[] }
  type FileStatus = { arquivo: string; tipo: string; ok: boolean; problema?: string }
  type PromptMeta = { id: string; titulo: string; data: string }

  // ---- estado ----
  let models:  Model[] = []
  let roles:   Role[]  = []
  let screen:  Screen  = 'model'

  let selectedModel: Model | null = null
  let selectedRoles: Set<string>  = new Set()
  let roleSearch = ''

  let narrative = ''
  let contextSections: ContextSection[] = []

  let usePhases = false
  let phases: Step[] = []

  let resultContent = ''
  let resultError   = ''
  let building      = false
  let copied        = false
  let saved         = false

  // ---- prompts salvos ----
  let promptList: PromptMeta[] = []
  let promptViewId: string | null = null
  let promptViewContent = ''
  let promptSearch = ''
  let promptToDelete: string | null = null
  let promptViewCopied = false
  let previousScreen: Screen = 'model'

  // ---- validação de arquivos ----
  let showValidation   = false
  let validationItems: FileStatus[] = []
  let validationVisible = 0   // quantos já aparecem na lista (animado)
  let validationDone   = false
  $: validationErrors = validationItems.filter(f => !f.ok)

  // ---- tutorial de uso ----
  let showTutorial  = false
  let tutorialSlide = 0

  const TUTORIAL_SLIDES = [
    {
      step:     0,
      icon:     '🗺️',
      accent:   '#f5a623',
      title:    'Como funciona o CASTOR',
      subtitle: 'Visão geral · 6 etapas',
      desc:     'O CASTOR guia você por um wizard simples. Cada etapa adiciona uma camada de contexto — no final, um prompt estruturado pronto para qualquer IA.',
      tip:      'Você pode voltar a etapas anteriores a qualquer momento.',
      preview:  'pipeline',
    },
    {
      step:     1,
      icon:     '📐',
      accent:   '#f5a623',
      title:    'Escolha o Modelo',
      subtitle: 'Passo 1 de 6',
      desc:     'O modelo define a estrutura e a lógica do prompt. Cada um serve um tipo diferente de tarefa — escolha o que melhor se encaixa no seu objetivo.',
      tip:      'Na dúvida, RACE é o mais versátil para a maioria dos casos.',
      preview:  'models',
    },
    {
      step:     2,
      icon:     '🎭',
      accent:   '#a371f7',
      title:    'Selecione os Papéis',
      subtitle: 'Passo 2 de 6',
      desc:     'O papel define quem a IA deve "ser" para responder sua tarefa. Você pode combinar vários — o prompt incorpora as habilidades de todos.',
      tip:      'Combinar papéis complementares (ex: Arquiteto + DevOps) enriquece muito o resultado.',
      preview:  'roles',
    },
    {
      step:     3,
      icon:     '✍️',
      accent:   '#58a6ff',
      title:    'Descreva sua Tarefa',
      subtitle: 'Passo 3 de 6',
      desc:     'Escreva naturalmente, como explicaria a um colega. Sem formatação especial — o CASTOR distribui o contexto automaticamente nos campos certos.',
      tip:      'Quanto mais contexto aqui, menos gaps serão perguntados depois.',
      preview:  'narrative',
    },
    {
      step:     4,
      icon:     '💬',
      accent:   '#3fb950',
      title:    'Preencha o Contexto',
      subtitle: 'Passo 4 de 6',
      desc:     'O CASTOR identifica o que falta e faz perguntas direcionadas. Algumas vêm do modelo, outras dos papéis — e cada pergunta mostra de onde veio.',
      tip:      'Campos opcionais podem ser pulados — aparecem no prompt como lacunas a considerar.',
      preview:  'gaps',
    },
    {
      step:     5,
      icon:     '📋',
      accent:   '#e06c75',
      title:    'Defina as Fases',
      subtitle: 'Passo 5 de 6 · opcional',
      desc:     'Para tarefas complexas, divida a execução em etapas sequenciais. Cada fase tem um título e uma descrição do que deve acontecer naquele momento.',
      tip:      'Fases são ótimas para projetos de múltiplas entregas ou raciocínio encadeado.',
      preview:  'phases',
    },
    {
      step:     6,
      icon:     '🚀',
      accent:   '#f5a623',
      title:    'Prompt Pronto!',
      subtitle: 'Passo 6 de 6',
      desc:     'O CASTOR monta o prompt com modelo, papéis, habilidades, contexto e fases — estruturado e pronto para usar em qualquer IA.',
      tip:      'Copie e cole diretamente no ChatGPT, Claude, Gemini ou qualquer outra IA.',
      preview:  'result',
    },
  ]

  function nextTutorial() {
    if (tutorialSlide < TUTORIAL_SLIDES.length - 1) tutorialSlide++
    else showTutorial = false
  }
  function prevTutorial() {
    if (tutorialSlide > 0) tutorialSlide--
  }

  async function runValidation() {
    showValidation   = true
    validationDone   = false
    validationVisible = 0
    validationItems  = []
    const results = await ValidateAll()
    validationItems = results
    // anima revelação item a item
    for (let i = 0; i < results.length; i++) {
      await new Promise(r => setTimeout(r, 60))
      validationVisible = i + 1
    }
    validationDone = true
    // auto-fecha se sem erros
    if (validationErrors.length === 0) {
      await new Promise(r => setTimeout(r, 800))
      showValidation = false
    }
  }

  // ---- ciclo de vida ----
  onMount(async () => {
    models = await GetModels()
    roles  = await GetRoles()
    const first = await IsFirstRun()
    if (first) { tutorialSlide = 0; showTutorial = true }
    await runValidation()
    hasSavedState = (() => {
      const raw = localStorage.getItem(SAVE_KEY)
      if (!raw) return false
      try {
        const s = JSON.parse(raw)
        return !!s.selectedModelId
      } catch { return false }
    })()
  })

  // ---- persistência de estado ----
  const SAVE_KEY = 'castor_wizard_state'

  function saveState() {
    const state = {
      screen,
      selectedModelId: selectedModel?.id ?? null,
      selectedRoleIds: [...selectedRoles],
      narrative,
      contextSections: contextSections.map(s => ({
        ...s,
        gaps: s.gaps.map(g => ({ fieldId: g.fieldId, pergunta: g.pergunta, tipo: g.tipo, opcoes: g.opcoes, obrigatorio: g.obrigatorio, roleNome: g.roleNome, answer: g.answer }))
      })),
      usePhases,
      phases,
    }
    localStorage.setItem(SAVE_KEY, JSON.stringify(state))
  }

  function clearState() {
    localStorage.removeItem(SAVE_KEY)
  }

  $: if (screen && screen !== 'result' && screen !== 'prompts') saveState()

  function resumeState() {
    const raw = localStorage.getItem(SAVE_KEY)
    if (!raw) return false
    try {
      const s = JSON.parse(raw)
      if (s.selectedModelId) selectedModel = models.find(m => m.id === s.selectedModelId) ?? null
      if (s.selectedRoleIds) selectedRoles = new Set(s.selectedRoleIds)
      if (s.narrative) narrative = s.narrative
      if (s.contextSections) contextSections = s.contextSections.map((sec: any) => ({ ...sec, gaps: sec.gaps.map((g: any) => ({ ...g, answer: g.answer ?? '' })) }))
      if (s.usePhases !== undefined) usePhases = s.usePhases
      if (s.phases) phases = s.phases
      if (s.screen && s.screen !== 'model' && s.screen !== 'prompts') {
        screen = s.screen
        allExpanded = false
        // expand all sections on resume
        collapsedSections = {}
      }
      return true
    } catch { return false }
  }

  function dismissResume() {
    clearState()
    hasSavedState = false
  }

  let hasSavedState = false

  // ---- stepper ----
  const STEPS: { id: Screen; label: string }[] = [
    { id: 'model',     label: 'Modelo'    },
    { id: 'role',      label: 'Papel'     },
    { id: 'narrative', label: 'Tarefa'    },
    { id: 'gap',       label: 'Contexto'  },
    { id: 'phase',     label: 'Fases'     },
    { id: 'result',    label: 'Resultado' },
  ]
  $: stepIndex = STEPS.findIndex(s => s.id === screen)

  // ---- helpers ----
  $: categories = [...new Set(roles.map(r => r.categoria))].sort()

  $: filteredRoles = roleSearch.trim()
    ? roles.filter(r =>
        r.nome.toLowerCase().includes(roleSearch.toLowerCase()) ||
        r.categoria.toLowerCase().includes(roleSearch.toLowerCase()))
    : roles

  // rolesByCategory é $: para que o Svelte rastreie filteredRoles como dependência
  // e re-renderize o template quando a busca mudar
  $: rolesByCategory = categories.reduce((acc, cat) => {
    acc[cat] = filteredRoles.filter(r => r.categoria === cat)
    return acc
  }, {} as Record<string, Role[]>)

  const CAT_ICON: Record<string, string> = {
    arquitetura: '🏗️', frontend: '🎨', backend: '⚙️',
    devops: '☁️', banco: '🗄️', dados: '📊',
    gestao: '📋', seguranca: '🔐', design: '✏️', marketing: '📣',
    qa: '🧪', documentacao: '📝',
  }
  const CAT_LABEL: Record<string, string> = {
    arquitetura: 'Arquitetura',   frontend:  'Frontend & Mobile',
    backend:     'Backend',       devops:    'DevOps & Cloud',
    banco:       'Banco de Dados',dados:     'Dados & IA',
    gestao:      'Gestão',        seguranca: 'Segurança',
    design:      'Design',        marketing: 'Marketing',
    qa:          'QA & Testes',   documentacao: 'Documentação',
  }
  function catLabel(cat: string) { return CAT_LABEL[cat] ?? cat }
  function catIcon(cat: string)  { return CAT_ICON[cat]  ?? '📁' }

  // modelo: badge colorida por ID
  const MODEL_COLOR: Record<string, string> = {
    race:   '#f5a623', rtf:    '#3fb950',
    risen:  '#a371f7', create: '#58a6ff',
  }
  function modelColor(id: string) { return MODEL_COLOR[id] ?? '#6e7681' }

  // Ordenação de gaps: obrigatórias vazias → preenchidas → opcionais vazias
  const gapPriority = (g: Gap) => g.answer ? 1 : (g.obrigatorio ? 0 : 2)
  const byGapPriority = (a: Gap, b: Gap) => gapPriority(a) - gapPriority(b)

  // ---- navegação ----
  const BACK_MAP: Partial<Record<Screen, () => Screen>> = {
    role:      () => 'model',
    narrative: () => 'role',
    gap:       () => 'narrative',
    phase:     () => (contextSections.length > 0 ? 'gap' : 'narrative'),
  }

  function goBack() {
    if (screen === 'prompts') {
      // Navegação em 2 níveis: primeiro fecha a visualização do prompt,
      // depois sai da tela voltando à origem (nunca permanece em 'prompts').
      if (promptViewId) { promptViewId = null; return }
      const target = previousScreen && previousScreen !== 'prompts' ? previousScreen : 'model'
      screen = target
      return
    }
    const next = BACK_MAP[screen]
    if (next) screen = next()
  }

  function toggleRole(id: string) {
    const s = new Set(selectedRoles)
    s.has(id) ? s.delete(id) : s.add(id)
    selectedRoles = s
  }

  function confirmModel(m: Model) {
    selectedModel = m
    selectedRoles = new Set()
    roleSearch    = ''
    screen = 'role'
  }

  function confirmRoles() {
    if (selectedRoles.size === 0) return
    narrative = ''
    screen    = 'narrative'
  }

  function confirmNarrative() {
    if (!narrative.trim()) return

    const sections: ContextSection[] = []

    // 1. Todos os campos do modelo (obrigatorios + opcionais) em uma secao
    const modelGaps: Gap[] = (selectedModel?.campos ?? []).map(c => ({
      fieldId:  c.id,
      pergunta:   c.label,
      tipo:       c.tipo,
      opcoes:     c.opcoes ?? [],
      obrigatorio: c.obrigatorio,
      answer:      '',
    }))
    modelGaps.sort(byGapPriority)
    if (modelGaps.length > 0) {
      sections.push({ title: 'Contexto do Modelo (' + (selectedModel?.nome ?? '') + ')', kind: 'model', gaps: modelGaps })
    }

    // 2. Gaps de papel — uma seção por papel
    for (const id of selectedRoles) {
      const role = roles.find(r => r.id === id)
      if (!role || !role.gaps_comuns?.length) continue
      const roleGaps: Gap[] = role.gaps_comuns.map(q => ({
        fieldId:  '',
        pergunta:   q,
        tipo:       'textarea',
        opcoes:     [],
        obrigatorio: false,
        roleNome:   role.nome,
        answer:     '',
      }))
      sections.push({ title: role.nome, kind: 'role', gaps: roleGaps })
    }

    contextSections = sections
    selSecIdx        = 0
    selGapIdx        = 0
    collapsedSections = {}
    expandEditor      = false
    allExpanded       = true
    screen = sections.length > 0 ? 'gap' : 'phase'
  }

  let gapError = false
  let visibleSections: ContextSection[] = []
  let collapsedSections: Record<string, boolean> = {}
  let expandEditor = false
  let selSecIdx = 0
  let selGapIdx = 0

  $: visibleSections = contextSections

  $: sortedVisibleSections = visibleSections.map(sec => ({
    ...sec,
    gaps: [...sec.gaps].sort(byGapPriority),
  }))

  $: selectedGap = (() => {
    const sec = sortedVisibleSections[selSecIdx]
    if (!sec) return null
    return { gap: sec.gaps[selGapIdx], sectionTitle: sec.title, sectionKind: sec.kind }
  })()

  $: totalCount = (() => {
    let n = 0
    for (const s of visibleSections) {
      if (!collapsedSections[s.title]) n += s.gaps.length
    }
    return n
  })()

  $: answeredCount = (() => {
    let n = 0
    for (const s of visibleSections) {
      if (!collapsedSections[s.title]) n += s.gaps.filter(g => g.answer).length
    }
    return n
  })()

  $: rawTotal = (() => {
    let n = 0
    for (const s of visibleSections) n += s.gaps.length
    return n
  })()

  // helper functions
  function cycleSelectOption(gap: Gap) {
    const idx = gap.opcoes.indexOf(gap.answer)
    gap.answer = gap.opcoes[(idx + 1) % gap.opcoes.length]
  }

  function toggleMultiOption(gap: Gap, opt: string) {
    const parts = gap.answer ? gap.answer.split(', ').filter(p => p) : []
    const found = parts.indexOf(opt)
    if (found >= 0) parts.splice(found, 1)
    else parts.push(opt)
    gap.answer = parts.join(', ')
  }

  function hasOpt(gap: Gap, opt: string): boolean {
    return (gap.answer ?? '').split(', ').includes(opt)
  }

  function formatAnswer(g: Gap): string {
    if (!g.answer) return ''
    if (g.tipo === 'select' || g.tipo === 'multiselect') return g.answer
    return g.answer.length > 30 ? g.answer.slice(0, 30) + '...' : g.answer
  }

  function gapIcon(g: Gap): string {
    if (g.answer) return '✓'
    if (g.obrigatorio) return '★'
    return '·'
  }

  function gapIconColor(g: Gap): string {
    if (g.answer) return 'text-[#3fb950]'
    if (g.obrigatorio) return 'text-[#f85149]'
    return 'text-[#3a3a50]'
  }

  function toggleSection(title: string) {
    collapsedSections[title] = !collapsedSections[title]
    collapsedSections = collapsedSections
  }

  let allExpanded = true

  function toggleAllSections() {
    if (allExpanded) {
      for (const s of visibleSections) collapsedSections[s.title] = true
    } else {
      collapsedSections = {}
    }
    allExpanded = !allExpanded
    collapsedSections = collapsedSections
  }

  function selectGap(si: number, gi: number) {
    selSecIdx = si
    selGapIdx = gi
  }

  function nextGap() {
    const vis = sortedVisibleSections
    if (vis.length === 0) { screen = 'phase'; return }
    // try next in current section
    const curSec = vis[selSecIdx]
    if (curSec && !collapsedSections[curSec.title] && selGapIdx < curSec.gaps.length - 1) {
      selGapIdx++
      // skip to next required empty if possible
      for (let gi = selGapIdx; gi < curSec.gaps.length; gi++) {
        if (curSec.gaps[gi].obrigatorio && !curSec.gaps[gi].answer) { selGapIdx = gi; return }
      }
      return
    }
    // move to next section
    for (let si = selSecIdx + 1; si < vis.length; si++) {
      if (collapsedSections[vis[si].title]) continue
      if (vis[si].gaps.length > 0) {
        selSecIdx = si
        selGapIdx = 0
        for (let gi = 0; gi < vis[si].gaps.length; gi++) {
          if (vis[si].gaps[gi].obrigatorio && !vis[si].gaps[gi].answer) { selGapIdx = gi; return }
        }
        return
      }
    }
    screen = 'phase'
  }

  function prevGap() {
    if (selGapIdx > 0) { selGapIdx--; return }
    const vis = sortedVisibleSections
    for (let si = selSecIdx - 1; si >= 0; si--) {
      if (collapsedSections[vis[si].title]) continue
      if (vis[si].gaps.length > 0) {
        selSecIdx = si
        selGapIdx = vis[si].gaps.length - 1
        return
      }
    }
  }

  $: hasUnansweredRequired = (() => {
    for (const s of visibleSections) {
      if (collapsedSections[s.title]) continue
      for (const g of s.gaps) {
        if (g.obrigatorio && !g.answer.trim()) return true
      }
    }
    return false
  })()

  function skipAll() {
    if (hasUnansweredRequired) return
    screen = 'phase'
  }

  async function build() {
    building = true
    const allGapAnswers = contextSections.flatMap(s => s.gaps.map(g => ({
      field_id: g.fieldId,
      pergunta: g.pergunta,
      resposta: g.answer,
      role_nome: g.roleNome ?? '',
    })))
    const result = await BuildPrompt(main.BuildRequestDTO.createFrom({
      model_id:    selectedModel!.id,
      role_ids:    [...selectedRoles],
      narrativa:   narrative,
      gap_answers: allGapAnswers,
      steps:       usePhases ? phases : [],
    }))
    resultContent = result.conteudo
    resultError   = result.erro ?? ''
    building      = false
    screen        = 'result'
    clearState()
  }

  async function copyResult() {
    await navigator.clipboard.writeText(resultContent)
    copied = true
    setTimeout(() => copied = false, 2000)
  }

  async function saveToLibrary() {
    await SavePrompt(resultContent)
    saved = true
    setTimeout(() => saved = false, 2000)
  }

  async function loadPrompts() {
    promptList = await ListPrompts()
  }

  async function viewPrompt(id: string) {
    promptViewId = id
    promptViewContent = ''
    promptToDelete = null
    promptViewCopied = false
    const p = await GetPrompt(id)
    promptViewContent = p.conteudo
  }

  function openPrompts() {
    // Toggle: se já está em prompts, volta à origem.
    if (screen === 'prompts') { goBack(); return }
    previousScreen = screen
    promptViewId = null
    promptSearch = ''
    promptToDelete = null
    loadPrompts()
    screen = 'prompts'
  }

  async function deletePrompt(id: string) {
    await DeletePrompt(id)
    if (promptViewId === id) promptViewId = null
    promptToDelete = null
    await loadPrompts()
  }

  function requestDelete(id: string) {
    promptToDelete = promptToDelete === id ? null : id
  }

  async function copyPromptView() {
    if (!promptViewContent) return
    await navigator.clipboard.writeText(promptViewContent)
    promptViewCopied = true
    setTimeout(() => promptViewCopied = false, 2000)
  }

  // Filtra prompts por título (case-insensitive). Busca vazia devolve a lista.
  $: filteredPrompts = promptSearch.trim()
    ? promptList.filter(p => (p.titulo || '').toLowerCase().includes(promptSearch.toLowerCase()))
    : promptList

  // Converte "2026-05-21 11:53" em "há 5 min", "há 2 h", "há 3 d" ou data absoluta.
  function relativeDate(d: string): string {
    if (!d) return ''
    const ts = Date.parse(d.replace(' ', 'T'))
    if (isNaN(ts)) return d
    const diff = Date.now() - ts
    const min = Math.floor(diff / 60_000)
    if (min < 1) return 'agora'
    if (min < 60) return `há ${min} min`
    const h = Math.floor(min / 60)
    if (h < 24) return `há ${h} h`
    const days = Math.floor(h / 24)
    if (days < 7) return `há ${days} d`
    return d
  }

  function restart() {
    selectedModel    = null
    selectedRoles    = new Set()
    narrative        = ''
    contextSections  = []
    selSecIdx        = 0
    selGapIdx        = 0
    collapsedSections = {}
    expandEditor      = false
    allExpanded       = true
    phases           = []
    usePhases        = false
    resultContent    = ''
    resultError      = ''
    screen           = 'model'
    clearState()
  }
</script>

<!-- ============================================================
     SHELL
============================================================= -->
<div class="flex flex-col h-screen bg-[#0a0a12] text-[#c9d1d9] font-mono select-none overflow-hidden">

  <!-- titlebar drag area -->
  <div class="h-8 flex-shrink-0" style="--wails-draggable:drag" />

  <!-- ============================================================
       VALIDAÇÃO OVERLAY
  ============================================================= -->
  {#if showValidation}
    <div class="fixed inset-0 z-50 flex items-center justify-center"
         in:fade={{ duration: 250 }} out:fade={{ duration: 200 }}>
      <div class="absolute inset-0 bg-[#04040a]/85 backdrop-blur-md"></div>

      <div class="relative z-10 w-[500px] rounded-2xl border border-[#1e1e30]
                  bg-[#0d0d18] shadow-2xl overflow-hidden"
           in:fly={{ y: 20, duration: 300, easing: cubicOut }}>

        <!-- header -->
        <div class="px-8 pt-8 pb-4">
          <p class="text-[10px] tracking-[0.25em] uppercase text-[#f5a623] font-semibold mb-1">
            {validationDone ? (validationErrors.length === 0 ? 'Tudo certo!' : `${validationErrors.length} problema${validationErrors.length > 1 ? 's' : ''} encontrado${validationErrors.length > 1 ? 's' : ''}`) : 'Verificando arquivos...'}
          </p>
          <h2 class="text-lg font-bold text-[#e8eaf0]">Validação de arquivos</h2>
        </div>

        <!-- gauge -->
        <div class="px-8 mb-4">
          <div class="h-1.5 w-full rounded-full bg-[#1a1a28] overflow-hidden">
            <div class="h-full rounded-full transition-all duration-300"
                 style="width:{validationItems.length ? (validationVisible/validationItems.length)*100 : 0}%;
                        background:{validationDone && validationErrors.length > 0 ? '#f85149' : '#f5a623'}"></div>
          </div>
          <div class="flex justify-between mt-1">
            <span class="text-[10px] text-[#4a5060]">{validationVisible} / {validationItems.length} arquivos</span>
            {#if validationDone}
              <span class="text-[10px] font-semibold"
                    class:text-[#3fb950]={validationErrors.length === 0}
                    class:text-[#f85149]={validationErrors.length > 0}>
                {validationErrors.length === 0 ? '✓ sem problemas' : `✗ ${validationErrors.length} com erro`}
              </span>
            {/if}
          </div>
        </div>

        <!-- lista de arquivos -->
        <div class="px-8 pb-4 max-h-64 overflow-y-auto flex flex-col gap-0.5">
          {#each validationItems.slice(0, validationVisible) as item}
            <div class="flex items-start gap-2.5 py-1.5 text-xs
                        {!item.ok ? 'text-[#f85149]' : 'text-[#6e7681]'}">
              <span class="flex-shrink-0 mt-0.5 font-bold">
                {item.ok ? '✓' : '✗'}
              </span>
              <div class="min-w-0">
                <div class="flex items-center gap-2">
                  <span class="font-mono truncate">{item.arquivo}</span>
                  <span class="text-[9px] px-1.5 py-0.5 rounded flex-shrink-0
                               {item.tipo === 'model'
                                 ? 'bg-[#f5a623]/15 text-[#f5a623]'
                                 : 'bg-[#a371f7]/15 text-[#a371f7]'}">
                    {item.tipo}
                  </span>
                </div>
                {#if !item.ok && item.problema}
                  <p class="text-[#f85149]/80 mt-0.5 leading-snug">{item.problema}</p>
                {/if}
              </div>
            </div>
          {/each}
        </div>

        <!-- footer -->
        {#if validationDone}
          <div class="flex justify-end px-8 py-5 border-t border-[#1a1a28]">
            <button on:click={() => showValidation = false}
              class="px-6 py-2 rounded-lg text-sm font-bold text-black transition-all
                     hover:brightness-110 active:scale-[0.97]
                     {validationErrors.length > 0 ? 'bg-[#f85149]' : 'bg-[#3fb950]'}">
              {validationErrors.length > 0 ? `Continuar mesmo assim →` : 'Continuar →'}
            </button>
          </div>
        {/if}
      </div>
    </div>
  {/if}

  <!-- ============================================================
       TUTORIAL OVERLAY
  ============================================================= -->
  {#if showTutorial}
    {@const ts = TUTORIAL_SLIDES[tutorialSlide]}
    <div class="fixed inset-0 z-50 flex items-center justify-center"
         in:fade={{ duration: 300 }} out:fade={{ duration: 200 }}>

      <div class="absolute inset-0 bg-[#04040a]/85 backdrop-blur-md"></div>

      <div class="relative z-10 w-[640px] rounded-2xl border border-[#1e1e30]
                  bg-[#0d0d18] shadow-2xl overflow-hidden"
           in:fly={{ y: 24, duration: 350, easing: cubicOut }}>

        <!-- barra de progresso -->
        <div class="h-0.5 w-full bg-[#1a1a28]">
          <div class="h-full transition-all duration-500 rounded-full"
               style="width:{((tutorialSlide + 1) / TUTORIAL_SLIDES.length) * 100}%;
                      background:{ts.accent}"></div>
        </div>

        <!-- header -->
        <div class="px-8 pt-7 pb-4 flex items-center gap-4">
          <div class="relative flex-shrink-0">
            <div class="absolute inset-0 rounded-full blur-xl opacity-40"
                 style="background:{ts.accent}; transform:scale(2)"></div>
            <div class="relative text-4xl leading-none select-none"
                 style="filter:drop-shadow(0 0 16px {ts.accent}99)">
              {ts.icon}
            </div>
          </div>
          <div>
            <p class="text-[10px] tracking-[0.2em] uppercase font-semibold mb-0.5"
               style="color:{ts.accent}">{ts.subtitle}</p>
            <h2 class="text-xl font-bold text-[#e8eaf0]">{ts.title}</h2>
          </div>
          <button on:click={() => showTutorial = false}
            class="ml-auto text-[#3a3a50] hover:text-[#6e7681] transition-colors text-lg leading-none">✕</button>
        </div>

        <!-- desc -->
        <div class="px-8 pb-4">
          <p class="text-sm text-[#8a94a8] leading-relaxed">{ts.desc}</p>
        </div>

        <!-- preview área -->
        <div class="mx-8 mb-4 rounded-xl border border-[#1e1e2e] bg-[#08080f] overflow-hidden">

          <!-- PIPELINE overview -->
          {#if ts.preview === 'pipeline'}
            <div class="p-5">
              <div class="flex items-stretch gap-0">
                {#each [
                  { icon:'📐', label:'Modelo',   accent:'#f5a623' },
                  { icon:'🎭', label:'Papéis',   accent:'#a371f7' },
                  { icon:'✍️', label:'Tarefa',   accent:'#58a6ff' },
                  { icon:'💬', label:'Contexto', accent:'#3fb950' },
                  { icon:'📋', label:'Fases',    accent:'#e06c75' },
                  { icon:'🚀', label:'Prompt',   accent:'#f5a623' },
                ] as s, si}
                  <div class="flex-1 flex flex-col items-center gap-1.5 relative">
                    <div class="w-10 h-10 rounded-xl flex items-center justify-center text-xl
                                border border-[#1e1e2e]"
                         style="background:{s.accent}12; border-color:{s.accent}25">
                      {s.icon}
                    </div>
                    <span class="text-[10px] font-semibold" style="color:{s.accent}">{s.label}</span>
                    {#if si < 5}
                      <div class="absolute top-5 left-[calc(50%+20px)] right-0 h-px"
                           style="background:linear-gradient(90deg,{s.accent}40,transparent)"></div>
                    {/if}
                  </div>
                {/each}
              </div>
              <div class="mt-4 pt-3 border-t border-[#1a1a28] grid grid-cols-3 gap-2">
                {#each [
                  { label:'Estrutura clara', desc:'O modelo organiza as seções do prompt' },
                  { label:'Especialistas', desc:'Os papéis definem o estilo e contexto' },
                  { label:'Sem IA no processo', desc:'Tudo é template engine e heurística' },
                ] as feat}
                  <div class="rounded-lg p-2.5 bg-[#0d0d18] border border-[#1e1e2e]">
                    <p class="text-[10px] font-bold text-[#c9d1d9] mb-0.5">{feat.label}</p>
                    <p class="text-[9px] text-[#4a5060] leading-snug">{feat.desc}</p>
                  </div>
                {/each}
              </div>
            </div>

          <!-- MODELS preview -->
          {:else if ts.preview === 'models'}
            <div class="p-4 flex flex-col gap-2">
              {#each [
                { id:'RACE',   color:'#f5a623', desc:'Contexto rico + entregável claro',    tag:'Mais usado'   },
                { id:'RTF',    color:'#3fb950', desc:'Tarefas diretas e objetivas',          tag:''             },
                { id:'RISEN',  color:'#a371f7', desc:'Steps detalhados com restrições',      tag:'Complexo'     },
                { id:'CREATE', color:'#58a6ff', desc:'Conteúdo criativo com público e tom',  tag:'Criativo'     },
              ] as m, mi}
                <div class="flex items-center gap-3 px-3 py-2.5 rounded-lg border transition-all
                            {mi === 0
                              ? 'border-[#f5a623]/40 bg-[#f5a623]/5'
                              : 'border-[#1e1e2e] bg-transparent opacity-60'}">
                  <div class="w-1.5 h-1.5 rounded-full flex-shrink-0"
                       style="background:{m.color}; {mi !== 0 ? 'opacity:0.4' : ''}"></div>
                  <span class="text-xs font-bold w-14 flex-shrink-0" style="color:{m.color}">{m.id}</span>
                  <span class="text-xs text-[#6e7681] flex-1">{m.desc}</span>
                  {#if m.tag}
                    <span class="text-[9px] px-1.5 py-0.5 rounded font-semibold flex-shrink-0"
                          style="background:{m.color}18; color:{m.color}">{m.tag}</span>
                  {/if}
                  {#if mi === 0}
                    <span class="text-[#f5a623] text-xs">✓</span>
                  {/if}
                </div>
              {/each}
            </div>

          <!-- ROLES preview -->
          {:else if ts.preview === 'roles'}
            <div class="p-4">
              <div class="flex gap-2 mb-3 flex-wrap">
                {#each [
                  { nome:'Arquiteto Cloud', color:'#a371f7', sel: true  },
                  { nome:'DevOps Engineer', color:'#a371f7', sel: true  },
                  { nome:'QA Lead',         color:'#a371f7', sel: false },
                ] as r}
                  <div class="flex items-center gap-1.5 px-2.5 py-1 rounded-lg border text-xs font-semibold
                              {r.sel
                                ? 'border-[#a371f7]/40 bg-[#a371f7]/10 text-[#a371f7]'
                                : 'border-[#2a2a42] text-[#4a5060]'}">
                    {#if r.sel}<span>✓</span>{:else}<span>+</span>{/if}
                    {r.nome}
                  </div>
                {/each}
              </div>
              <div class="border-t border-[#1a1a28] pt-3">
                <p class="text-[9px] uppercase tracking-widest text-[#3a3a50] mb-2">🏗️ Arquitetura</p>
                <div class="grid grid-cols-2 gap-1.5">
                  {#each ['Arquiteto Cloud ✓','Arquiteto de Software ✓','Arquiteto de Soluções','Arquiteto de Microsserviços'] as r}
                    <div class="text-[10px] px-2 py-1 rounded border
                                {r.includes('✓')
                                  ? 'border-[#a371f7]/30 bg-[#a371f7]/8 text-[#a371f7]'
                                  : 'border-[#1e1e2e] text-[#4a5060]'}">
                      {r}
                    </div>
                  {/each}
                </div>
              </div>
            </div>

          <!-- NARRATIVE preview -->
          {:else if ts.preview === 'narrative'}
            <div class="p-4">
              <div class="rounded-lg border border-[#2a2a42] bg-[#0d0d18] p-3 font-mono text-xs
                          text-[#8a94a8] leading-relaxed mb-3">
                <span class="text-[#58a6ff]">Preciso criar uma API de pagamentos</span> para o nosso
                e-commerce. Usamos Go e PostgreSQL, com pico de <span class="text-[#3fb950]">10k req/s</span>
                em datas comemorativas. O sistema precisa integrar com o gateway da Stripe e suportar
                <span class="text-[#f5a623]">Pix, cartão e boleto</span>. A equipe tem 4 devs sênior.
                <span class="text-[#4a5060] animate-pulse">█</span>
              </div>
              <div class="flex gap-2 flex-wrap">
                {#each [
                  { label:'Stack detectada', val:'Go + PostgreSQL', c:'#58a6ff' },
                  { label:'Carga', val:'10k req/s', c:'#3fb950' },
                  { label:'Integrações', val:'Stripe, Pix', c:'#f5a623' },
                ] as hint}
                  <div class="flex items-center gap-1.5 px-2 py-1 rounded text-[10px] border"
                       style="border-color:{hint.c}25; background:{hint.c}10; color:{hint.c}">
                    <span class="opacity-60">{hint.label}:</span>
                    <span class="font-bold">{hint.val}</span>
                  </div>
                {/each}
              </div>
            </div>

          <!-- GAPS preview -->
          {:else if ts.preview === 'gaps'}
            <div class="p-4 flex flex-col gap-3">
              <div class="rounded-lg border border-[#3fb950]/25 bg-[#3fb950]/5 p-3">
                <div class="flex items-start justify-between gap-2 mb-2">
                  <p class="text-xs text-[#c9d1d9] font-medium leading-snug">
                    Comunicação síncrona (REST/gRPC) ou assíncrona (Kafka/NATS)?
                  </p>
                  <span class="text-[#3fb950] text-base flex-shrink-0">?</span>
                </div>
                <div class="flex items-center gap-1.5">
                  <span class="text-[9px] px-2 py-0.5 rounded-full border font-semibold
                               bg-[#a371f7]/12 text-[#a371f7] border-[#a371f7]/20">
                    🎭 Arquiteto Cloud
                  </span>
                </div>
              </div>
              <div class="rounded-lg border border-[#1e1e2e] bg-[#0d0d18] px-3 py-2 text-xs
                          text-[#4a5060] font-mono">
                REST síncrono para pagamentos, Kafka para eventos de confirmação...
                <span class="text-[#3fb950] animate-pulse">█</span>
              </div>
              <div class="flex justify-between items-center text-[10px]">
                <span class="text-[#3a3a50]">Pergunta 3 de 8</span>
                <span class="text-[#3fb950]">▓▓▓░░░░░ 37%</span>
              </div>
            </div>

          <!-- PHASES preview -->
          {:else if ts.preview === 'phases'}
            <div class="p-4 flex flex-col gap-2">
              {#each [
                { n:1, title:'Diagnóstico e Análise', desc:'Avalie os requisitos de throughput, modelagem do banco e pontos de integração com o gateway.', done: true  },
                { n:2, title:'Implementação do Core', desc:'Construa os endpoints principais de criação e consulta de transações com idempotência.', done: false },
                { n:3, title:'Integrações e Testes',  desc:'Conecte Stripe, Pix e boleto. Cobertura de testes de contrato e carga.', done: false },
              ] as ph}
                <div class="flex gap-3 rounded-lg border p-3
                            {ph.done
                              ? 'border-[#e06c75]/30 bg-[#e06c75]/5'
                              : 'border-[#1e1e2e] opacity-50'}">
                  <div class="w-5 h-5 rounded-full flex items-center justify-center text-[10px]
                               font-bold flex-shrink-0 mt-0.5"
                       style="background:{ph.done ? '#e06c75' : '#1e1e2e'}; color:{ph.done ? '#fff' : '#3a3a50'}">
                    {ph.n}
                  </div>
                  <div>
                    <p class="text-xs font-bold text-[#c9d1d9] mb-0.5">{ph.title}</p>
                     <p class="text-[10px] text-[#4a5060] leading-snug">{ph.desc}</p>
                  </div>
                </div>
              {/each}
            </div>

          <!-- RESULT preview -->
          {:else if ts.preview === 'result'}
            <div class="p-4">
              <div class="rounded-lg border border-[#1e1e2e] bg-[#0d0d18] p-3 font-mono text-[10px]
                          text-[#6e7681] leading-relaxed max-h-40 overflow-hidden relative">
                <div class="text-[#f5a623] font-bold mb-1"># Prompt — Arquiteto Cloud</div>
                <div class="text-[#3fb950] mb-1">## Papel</div>
                <div class="mb-2">Você é um Arquiteto Cloud sênior. Especialista em sistemas de
                alta disponibilidade, multi-region e otimização de custos em AWS/GCP/Azure...</div>
                <div class="text-[#3fb950] mb-1">## Contexto</div>
                <div class="mb-2">API de pagamentos Go + PostgreSQL. Pico 10k req/s. Integração
                Stripe, Pix e boleto. Time de 4 devs sênior...</div>
                <div class="text-[#a371f7] mb-1">## Habilidades relevantes</div>
                <div>- AWS / GCP / Azure&#10;- Infraestrutura como código&#10;- Kubernetes...</div>
                <div class="absolute bottom-0 inset-x-0 h-10 bg-gradient-to-t from-[#0d0d18]"></div>
              </div>
              <div class="flex justify-end mt-3">
                <div class="flex items-center gap-2 px-4 py-2 rounded-lg text-xs font-bold
                            bg-[#f5a623] text-black cursor-default">
                  📋 Copiar prompt
                </div>
              </div>
            </div>
          {/if}

        </div>

        <!-- dica -->
        <div class="mx-8 mb-5 flex items-start gap-2.5 px-3 py-2.5 rounded-lg
                    bg-[#f5a623]/6 border border-[#f5a623]/15">
          <span class="text-sm flex-shrink-0">💡</span>
          <p class="text-[11px] text-[#8a94a8] leading-relaxed">{ts.tip}</p>
        </div>

        <!-- footer -->
        <div class="flex items-center justify-between px-8 py-5 border-t border-[#1a1a28]">

          <!-- dots -->
          <div class="flex gap-1.5">
            {#each TUTORIAL_SLIDES as _, i}
              <button on:click={() => tutorialSlide = i}
                class="rounded-full transition-all duration-300"
                style="width:{i === tutorialSlide ? '20px' : '6px'};
                       height:6px;
                       background:{i === tutorialSlide ? ts.accent : '#2a2a42'}">
              </button>
            {/each}
          </div>

          <!-- nav -->
          <div class="flex items-center gap-3">
            {#if tutorialSlide > 0}
              <button on:click={prevTutorial}
                class="px-4 py-2 rounded-lg text-xs text-[#4a5060]
                       hover:text-[#c9d1d9] transition-colors border border-transparent
                       hover:border-[#2a2a42]">
                ← Anterior
              </button>
            {:else}
              <button on:click={() => showTutorial = false}
                class="px-4 py-2 rounded-lg text-xs text-[#3a3a50]
                       hover:text-[#6e7681] transition-colors">
                Pular tutorial
              </button>
            {/if}

            <button on:click={nextTutorial}
              class="px-5 py-2 rounded-lg text-sm font-bold text-black transition-all
                     hover:brightness-110 active:scale-[0.97]"
              style="background:{ts.accent}">
              {tutorialSlide < TUTORIAL_SLIDES.length - 1 ? 'Próximo →' : '🚀 Começar a criar!'}
            </button>
          </div>
        </div>

      </div>
    </div>
  {/if}

  <!-- main content -->
  <div class="flex-1 flex overflow-hidden">

    <!-- ── SIDEBAR esquerda ── -->
    <aside class="w-56 flex-shrink-0 flex flex-col items-center pt-7 pb-7 px-5
                  border-r border-[#1a1a28] bg-[#0d0d18]">

      <!-- mascote -->
      <img src="/castor.png" alt="Castor"
           class="w-24 h-24 object-contain mb-3
                  drop-shadow-[0_0_18px_rgba(245,166,35,0.5)]" />

      <div class="text-[10px] tracking-[0.3em] text-[#4a5060] uppercase mb-1">Prompt Builder</div>
      <h1 class="text-base font-bold tracking-wider mb-7">
        <span class="text-[#f5a623]">CASTOR</span>
        <span class="text-[#e06b2e]"> BUILDER</span>
      </h1>

      <!-- biblioteca de prompts -->
      <button on:click={openPrompts}
        class="w-full flex items-center gap-2.5 px-3.5 py-2.5 mb-7 rounded-lg
               border border-[#1a1a28] text-[#6e7681] text-xs
               hover:border-[#a371f7]/40 hover:text-[#a371f7] transition-all
               {screen === 'prompts' ? 'border-[#a371f7]/40 text-[#a371f7] bg-[#a371f7]/8' : ''}">
        <span class="text-sm">📚</span>
        Prompts salvos
      </button>

      <!-- stepper vertical -->
      <nav class="w-full flex flex-col gap-1">
        {#each STEPS as step, i}
          {@const done    = i < stepIndex}
          {@const active  = i === stepIndex}
          {@const locked  = i > stepIndex}
          <div class="flex items-center gap-2.5 px-2 py-1.5 rounded-lg transition-all
                      {active ? 'bg-[#1a1a30]' : ''}">
            <!-- círculo -->
            <div class="w-5 h-5 rounded-full flex-shrink-0 flex items-center justify-center text-[10px] font-bold
                        transition-all
                        {done   ? 'bg-[#f5a623] text-black' : ''}
                        {active ? 'bg-[#f5a623]/20 border border-[#f5a623] text-[#f5a623]' : ''}
                        {locked ? 'bg-[#1a1a28] border border-[#2a2a40] text-[#3a3a50]' : ''}">
              {#if done}✓{:else}{i + 1}{/if}
            </div>
            <!-- label -->
            <span class="text-xs transition-colors
                         {active ? 'text-[#f5a623] font-semibold' : ''}
                         {done   ? 'text-[#6e7681]' : ''}
                         {locked ? 'text-[#2a2a40]' : ''}">
              {step.label}
            </span>
          </div>
          <!-- conector -->
          {#if i < STEPS.length - 1}
            <div class="ml-4 w-px h-3 {i < stepIndex ? 'bg-[#f5a623]/40' : 'bg-[#1a1a28]'}"></div>
          {/if}
        {/each}
      </nav>

      <!-- tutorial + versão -->
      <div class="mt-auto flex flex-col items-center gap-4 w-full pt-6">
        <button on:click={() => { tutorialSlide = 0; showTutorial = true }}
          class="w-full flex items-center justify-center gap-2.5 px-4 py-3 rounded-xl
                 bg-[#f5a623] text-black text-[12px] font-bold tracking-wide
                 hover:bg-[#e09010] active:scale-[0.97] transition-all
                 shadow-[0_0_16px_rgba(245,166,35,0.25)] hover:shadow-[0_0_24px_rgba(245,166,35,0.4)]
                 group">
          <span class="text-base group-hover:scale-110 transition-transform">🗺️</span>
          Ver tutorial
        </button>
        <div class="text-[10px] tracking-wider text-[#2a2a40]">v0.1.0</div>
      </div>
    </aside>

    <!-- ── CONTEÚDO PRINCIPAL ── -->
    <main class="flex-1 flex flex-col overflow-hidden">

      <!-- topbar: voltar + breadcrumb -->
      {#if screen !== 'model'}
        <div class="flex items-center gap-3 px-8 pt-6 pb-0 flex-shrink-0">
          <button on:click={screen === 'result' ? restart : goBack}
            class="flex items-center gap-1.5 text-xs font-medium px-3 py-1.5 rounded-lg
                   border border-[#2a2a42] text-[#8a94a8]
                   hover:border-[#f5a623]/40 hover:text-[#f5a623] hover:bg-[#f5a623]/5
                   transition-all group">
            <span class="text-sm leading-none group-hover:-translate-x-0.5 transition-transform">←</span>
            <span>{screen === 'result' ? 'Novo prompt' : 'Voltar'}</span>
          </button>
          {#if screen !== 'result' && screen !== 'prompts'}
            <span class="text-[#2a2a40]">/</span>
            <span class="text-xs text-[#6e7681]">
              {#if selectedModel}<span class="text-[#f5a623]">{selectedModel.nome}</span>{/if}
              {#if selectedRoles.size > 0}
                <span class="text-[#2a2a40] mx-1">›</span>
                <span class="text-[#c9d1d9]">
                  {[...selectedRoles].map(id => roles.find(r => r.id === id)?.nome).filter(Boolean).join(' + ')}
                </span>
              {/if}
            </span>
          {/if}
        </div>
      {/if}

      <!-- scroll area -->
      <div class="flex-1 overflow-y-auto px-8 py-6 flex flex-col min-h-0">

        <!-- ============================================================
             TELA 1 — MODELO
        ============================================================= -->
        {#if screen === 'model'}
          <div in:fly={{ y: 16, duration: 200 }}>
            {#if hasSavedState}
              <div class="flex items-center justify-between p-3 mb-5 rounded-xl border border-[#f5a623]/30 bg-[#f5a623]/8">
                <span class="text-xs text-[#f5a623]">Você tem um progresso salvo. Deseja continuar?</span>
                <div class="flex gap-2">
                  <button on:click={() => { resumeState(); hasSavedState = false }}
                    class="px-4 py-1.5 rounded-lg bg-[#f5a623] text-black text-xs font-bold
                           hover:bg-[#e09010] transition-all">
                    Sim, continuar
                  </button>
                  <button on:click={dismissResume}
                    class="px-4 py-1.5 rounded-lg border border-[#2a2a42] text-[#6e7681] text-xs
                           hover:border-[#f5a623]/40 transition-all">
                    Não, começar novo
                  </button>
                </div>
              </div>
            {/if}
            <h2 class="text-lg font-bold mb-1">Escolha o framework</h2>
            <p class="text-sm text-[#6e7681] mb-6">
              Cada modelo estrutura o prompt de uma forma diferente.<br>
              Escolha o que melhor se encaixa na sua tarefa.
            </p>

            <div class="grid grid-cols-2 gap-3">
              {#each models as m}
                <button on:click={() => confirmModel(m)}
                  class="relative flex flex-col gap-2 p-4 rounded-xl border border-[#1a1a28]
                         bg-[#0d0d18] hover:border-[#f5a623]/60 hover:bg-[#111120]
                         transition-all text-left group overflow-hidden">
                  <!-- glow de fundo sutil -->
                  <div class="absolute inset-0 opacity-0 group-hover:opacity-100 transition-opacity rounded-xl"
                       style="background: radial-gradient(ellipse at top left, {modelColor(m.id)}10, transparent 70%)"></div>
                  <!-- badge ID -->
                  <span class="text-xs font-bold px-2 py-0.5 rounded-md self-start"
                        style="background:{modelColor(m.id)}22; color:{modelColor(m.id)}">
                    {m.id.toUpperCase()}
                  </span>
                  <span class="text-sm font-semibold text-[#e0e6f0] group-hover:text-white transition-colors">
                    {m.nome}
                  </span>
                  <span class="text-xs text-[#6e7681] leading-relaxed">{m.descricao}</span>
                  <span class="absolute top-3 right-3 text-[#2a2a40] group-hover:text-[#f5a623]/60
                               text-lg transition-colors">→</span>
                </button>
              {/each}
            </div>
          </div>

        <!-- ============================================================
             TELA 2 — PAPEL
        ============================================================= -->
        {:else if screen === 'role'}
          <div in:fly={{ y: 16, duration: 200 }}>
            <div class="flex items-center justify-between mb-1">
              <h2 class="text-lg font-bold">Selecione o(s) papel(eis)</h2>
              {#if selectedRoles.size > 0}
                <span class="text-xs px-2.5 py-1 rounded-full bg-[#f5a623]/15 text-[#f5a623] border border-[#f5a623]/30">
                  {selectedRoles.size} selecionado{selectedRoles.size > 1 ? 's' : ''}
                </span>
              {/if}
            </div>
            <p class="text-sm text-[#6e7681] mb-4">
              Você pode combinar múltiplos papéis para um prompt mais rico.
            </p>

            <!-- busca -->
            <div class="relative mb-4">
              <span class="absolute left-3 top-1/2 -translate-y-1/2 text-[#4a5060] text-sm">⌕</span>
              <input bind:value={roleSearch}
                placeholder="Buscar papel ou categoria..."
                class="w-full pl-8 pr-4 py-2 rounded-lg border border-[#1a1a28]
                       bg-[#0d0d18] text-sm text-[#c9d1d9] placeholder-[#4a5060]
                       focus:outline-none focus:border-[#f5a623]/60 transition-colors" />
            </div>

            <!-- lista categorizada -->
            <div class="flex flex-col gap-0.5 mb-5">
              {#each categories as cat (cat)}
                {#if (rolesByCategory[cat] ?? []).length > 0}
                  <div class="flex items-center gap-2 mt-4 mb-1.5">
                    <span class="text-base">{catIcon(cat)}</span>
                    <span class="text-[10px] font-semibold tracking-widest uppercase text-[#4a5060]">
                      {catLabel(cat)}
                    </span>
                  </div>
                  {#each rolesByCategory[cat] as role (role.id)}
                    <button on:click={() => toggleRole(role.id)}
                      class="flex items-center gap-3 w-full px-3 py-2 rounded-lg
                             border transition-all text-left text-sm
                             {selectedRoles.has(role.id)
                               ? 'border-[#f5a623]/40 bg-[#f5a623]/8 text-[#f5a623]'
                               : 'border-transparent hover:border-[#1a1a28] hover:bg-[#0d0d18] text-[#c9d1d9]'}">
                      <span class="w-4 h-4 rounded flex-shrink-0 flex items-center justify-center text-[10px] transition-all
                                   {selectedRoles.has(role.id) ? 'bg-[#f5a623] text-black font-bold' : 'border border-[#2a2a40] text-[#4a5060]'}">
                        {selectedRoles.has(role.id) ? '✓' : ''}
                      </span>
                      <span class="truncate">{role.nome}</span>
                    </button>
                  {/each}
                {/if}
              {/each}
            </div>

            <div class="sticky bottom-0 pt-3 pb-1 bg-[#0a0a12]/80 backdrop-blur-sm">
              <button on:click={confirmRoles} disabled={selectedRoles.size === 0}
                class="w-full py-2.5 rounded-xl bg-[#f5a623] text-black font-bold text-sm
                       hover:bg-[#e09010] disabled:opacity-30 disabled:cursor-not-allowed
                       transition-all active:scale-[0.98]">
                Continuar →
              </button>
            </div>
          </div>

        <!-- ============================================================
             TELA 3 — NARRATIVA
        ============================================================= -->
        {:else if screen === 'narrative'}
          <div in:fly={{ y: 16, duration: 200 }}>
            <h2 class="text-lg font-bold mb-1">Descreva sua tarefa</h2>
            <p class="text-sm text-[#6e7681] mb-5">
              Escreva livremente. Quanto mais contexto, melhor o prompt gerado.
            </p>

            <textarea bind:value={narrative}
              placeholder="Ex: Preciso criar um plano editorial para os próximos 3 meses. A empresa é uma startup B2B de SaaS que está com queda de engajamento no blog. O público-alvo são desenvolvedores e CTOs..."
              rows="10"
              class="w-full px-4 py-3 rounded-xl border border-[#1a1a28]
                     bg-[#0d0d18] text-[#c9d1d9] placeholder-[#3a3a50] text-sm
                     resize-none focus:outline-none focus:border-[#f5a623]/60
                     transition-colors leading-relaxed" />

            <div class="flex items-center justify-between mt-4">
              <span class="text-xs text-[#3a3a50]">{narrative.length} caracteres</span>
              <button on:click={confirmNarrative} disabled={!narrative.trim()}
                class="px-6 py-2.5 rounded-xl bg-[#f5a623] text-black font-bold text-sm
                       hover:bg-[#e09010] disabled:opacity-30 disabled:cursor-not-allowed
                       transition-all active:scale-[0.98]">
                Continuar →
              </button>
            </div>
          </div>

        <!-- ============================================================
             TELA 4 — GAPS (two-column editor)
        ============================================================= -->
        {:else if screen === 'gap'}
          <div class="flex flex-col h-full" in:fly={{ y: 16, duration: 200 }}>
            <!-- header -->
            <div class="flex-shrink-0 mb-3">
              <h2 class="text-lg font-bold mb-1">Contexto adicional</h2>
              <p class="text-sm text-[#6e7681]">
                Preencha as lacunas — ★ obrigatório, ✓ preenchido
              </p>
            </div>

            <!-- top bar: collapse all + progress -->
            <div class="flex-shrink-0 flex items-center justify-between mb-3">
              <div class="flex items-center gap-2">
                <button on:click={toggleAllSections}
                  class="text-[10px] px-3 py-1 rounded-full border font-semibold transition-all
                         {allExpanded ? 'border-[#f5a623]/40 text-[#f5a623] bg-[#f5a623]/10' : 'border-[#1a1a28] text-[#4a5060] hover:text-[#f5a623] hover:border-[#f5a623]/40'}">
                  {allExpanded ? '▼ Expandido' : '▶ Compacto'}
                </button>
              </div>
              <!-- progress + layout toggle -->
              <div class="flex items-center gap-2">
                <button on:click={() => expandEditor = !expandEditor}
                  class="text-[10px] px-2 py-1 rounded border transition-all
                         {expandEditor ? 'border-[#f5a623]/40 text-[#f5a623] bg-[#f5a623]/10' : 'border-[#1a1a28] text-[#4a5060] hover:text-[#f5a623] hover:border-[#f5a623]/40'}"
                  title={expandEditor ? 'Mostrar lista de lacunas' : 'Expandir área de resposta'}>
                  {expandEditor ? '◧ Ambos' : '◨ Editor'}
                </button>
                <div class="w-24 h-1.5 rounded-full bg-[#1a1a28] overflow-hidden">
                  <div class="h-full bg-[#f5a623] rounded-full transition-all duration-300"
                       style="width: {totalCount > 0 ? (answeredCount / totalCount) * 100 : 0}%"></div>
                </div>
                <span class="text-[10px] text-[#4a5060]">{answeredCount}/{totalCount}</span>
              </div>
            </div>

            <!-- two-column body -->
            <div class="flex-1 flex gap-4 min-h-0">
              <!-- LEFT: gap list -->
              {#if !expandEditor}
              <div class="w-[38%] flex-shrink-0 overflow-y-auto rounded-xl border border-[#1a1a28] bg-[#0d0d18]">
                {#if rawTotal === 0}
                  <div class="p-4 text-xs text-[#4a5060] text-center">Nenhuma lacuna pendente</div>
                {:else}
                  <div class="flex flex-col">
                    {#each sortedVisibleSections as sec, si}
                      {@const collapsed = collapsedSections[sec.title] === true}
                      <!-- section header -->
                      <button
                        on:click={() => toggleSection(sec.title)}
                        class="flex items-center gap-2 px-3 py-1.5 text-[10px] font-semibold tracking-wider uppercase border-b transition-colors w-full text-left
                               {sec.kind === 'model'
                                 ? 'text-[#f5a623]/80 border-[#f5a623]/20 bg-[#f5a623]/5'
                                 : 'text-[#a371f7]/80 border-[#1a1a28]'}">
                        <span class="text-[#6e7681]">{collapsed ? '▶' : '▼'}</span>
                        <span>{sec.kind === 'model' ? '📐 ' : '🎭 '}{sec.title}</span>
                        <span class="ml-auto text-[#3a3a50]">{sec.gaps.length}</span>
                      </button>
                      {#if collapsed}
                        <div class="px-3 py-1 text-[10px] text-[#3a3a50] border-b border-[#1a1a28]">
                          {sec.gaps.filter(g => g.answer).length}/{sec.gaps.length} preenchidos
                        </div>
                      {:else}
                        {#each sec.gaps as g, gi}
                          <button
                            on:click={() => selectGap(si, gi)}
                            class="flex items-center gap-2 px-3 py-2 text-left text-xs transition-all
                                   border-b border-[#0e0e1a] last:border-b-0
                                   {si === selSecIdx && gi === selGapIdx
                                     ? 'bg-[#f5a623]/10 border-l-2 border-l-[#f5a623]'
                                     : 'border-l-2 border-l-transparent hover:bg-[#1a1a28]'}">
                            <span class="flex-shrink-0 w-4 text-center text-xs {gapIconColor(g)}">
                              {gapIcon(g)}
                            </span>
                            <span class="flex-1 truncate {si === selSecIdx && gi === selGapIdx ? 'text-[#f5a623] font-medium' : 'text-[#c9d1d9]'}">
                              {g.pergunta}
                            </span>
                            {#if g.answer}
                              <span class="flex-shrink-0 text-[10px] text-[#3fb950]/70 truncate max-w-[80px]">
                                {formatAnswer(g)}
                              </span>
                            {/if}
                          </button>
                        {/each}
                      {/if}
                    {/each}
                  </div>
                {/if}
              </div>
              {/if}

              <!-- RIGHT: active gap editor -->
              <div class="flex-1 flex flex-col min-w-0 overflow-y-auto rounded-xl border border-[#1a1a28] bg-[#0d0d18] p-4">
                {#if selectedGap}
                  <div class="mb-3">
                    <div class="flex items-start justify-between gap-2">
                      <p class="text-sm text-[#e0e6f0] font-medium leading-relaxed">
                        {selectedGap.gap.obrigatorio ? '★ ' : ''}{selectedGap.gap.pergunta}
                      </p>
                      {#if selectedGap.gap.obrigatorio}
                        <span class="flex-shrink-0 text-[10px] px-2 py-0.5 rounded-full bg-[#f85149]/15 text-[#f85149] border border-[#f85149]/20">
                          obrigatório
                        </span>
                      {/if}
                    </div>
                    {#if selectedGap.sectionKind === 'role'}
                      <span class="inline-block mt-1.5 text-[10px] px-2 py-0.5 rounded-full bg-[#a371f7]/15 text-[#a371f7] border border-[#a371f7]/20">
                        {selectedGap.sectionTitle}
                      </span>
                    {/if}
                  </div>

                  {#if selectedGap.gap.tipo === 'select' && selectedGap.gap.opcoes.length > 0}
                    <div class="flex flex-wrap gap-2">
                      {#each selectedGap.gap.opcoes as opt}
                        <button
                          on:click={() => selectedGap.gap.answer = opt}
                          class="px-4 py-2 rounded-lg border text-sm transition-all
                                 {selectedGap.gap.answer === opt
                                   ? 'border-[#f5a623] bg-[#f5a623]/15 text-[#f5a623] font-medium'
                                   : 'border-[#1a1a28] text-[#6e7681] hover:border-[#f5a623]/40'}">
                          {opt}
                        </button>
                      {/each}
                    </div>
                  {:else if selectedGap.gap.tipo === 'multiselect' && selectedGap.gap.opcoes.length > 0}
                    <div class="flex flex-wrap gap-2">
                      {#each selectedGap.gap.opcoes as opt}
                        <button
                          on:click={() => toggleMultiOption(selectedGap.gap, opt)}
                          class="px-4 py-2 rounded-lg border text-sm transition-all
                                 {hasOpt(selectedGap.gap, opt)
                                   ? 'border-[#f5a623] bg-[#f5a623]/15 text-[#f5a623] font-medium'
                                   : 'border-[#1a1a28] text-[#6e7681] hover:border-[#f5a623]/40'}">
                          {hasOpt(selectedGap.gap, opt) ? '✓ ' : ''}{opt}
                        </button>
                      {/each}
                    </div>
                  {:else}
                    <textarea
                      bind:value={selectedGap.gap.answer}
                      placeholder={selectedGap.gap.obrigatorio ? 'Campo obrigatório — preencha para continuar' : 'Digite sua resposta...'}
                      rows="6"
                      class="w-full flex-1 px-4 py-3 rounded-lg border text-sm
                             bg-[#0a0a12] text-[#c9d1d9] placeholder-[#3a3a50]
                             resize-none focus:outline-none transition-colors leading-relaxed
                             border-[#1a1a28] focus:border-[#f5a623]/60" />
                  {/if}
                {:else}
                  <div class="flex-1 flex items-center justify-center text-xs text-[#4a5060]">
                    Selecione uma lacuna na lista ao lado
                  </div>
                {/if}
              </div>
            </div>

            <!-- bottom bar -->
            <div class="flex-shrink-0 flex justify-between items-center mt-3 pt-3 border-t border-[#1a1a28]">
              <div class="flex gap-2">
                <button on:click={prevGap} disabled={selSecIdx === 0 && selGapIdx === 0}
                  class="px-4 py-2 rounded-lg border text-xs transition-all
                         disabled:opacity-20 disabled:cursor-not-allowed
                         {selSecIdx > 0 || selGapIdx > 0
                           ? 'border-[#f5a623]/40 text-[#f5a623] hover:bg-[#f5a623]/10'
                           : 'border-[#1a1a28] text-[#3a3a50]'}">
                  ← Anterior
                </button>
                <button on:click={skipAll}
                  class="px-4 py-2 rounded-lg border text-xs transition-all
                         {hasUnansweredRequired
                           ? 'border-[#f85149]/30 text-[#f85149]/50 cursor-not-allowed'
                           : 'border-[#1a1a28] text-[#4a5060] hover:border-[#f5a623]/40 hover:text-[#f5a623]'}"
                  disabled={hasUnansweredRequired}
                  title={hasUnansweredRequired ? 'Preencha os campos obrigatórios primeiro' : 'Pular todas as lacunas'}>
                  ⏭ Pular tudo
                </button>
              </div>
              <span class="text-[10px] text-[#4a5060]">
                {totalCount}
              </span>
              <button on:click={nextGap}
                class="px-6 py-2 rounded-xl bg-[#f5a623] text-black font-bold text-sm
                       hover:bg-[#e09010] transition-all active:scale-[0.98]">
                {answeredCount >= totalCount ? 'Continuar →' : 'Próxima →'}
              </button>
            </div>
          </div>

        <!-- ============================================================
             TELA 5 — FASES
        ============================================================= -->
        {:else if screen === 'phase'}
          <div in:fly={{ y: 16, duration: 200 }}>
            <h2 class="text-lg font-bold mb-1">Fases de execução</h2>
            <p class="text-sm text-[#6e7681] mb-5">
              Divida o trabalho em etapas para um prompt mais estruturado.
            </p>

            <!-- toggle -->
            <div class="flex gap-2 mb-6">
              <button on:click={() => usePhases = true}
                class="flex-1 py-2.5 rounded-xl border text-sm font-medium transition-all
                       {usePhases
                         ? 'border-[#f5a623]/50 bg-[#f5a623]/10 text-[#f5a623]'
                         : 'border-[#1a1a28] text-[#6e7681] hover:border-[#f5a623]/30'}">
                ✦ Sim, quero definir fases
              </button>
              <button on:click={() => { usePhases = false; phases = [] }}
                class="flex-1 py-2.5 rounded-xl border text-sm font-medium transition-all
                       {!usePhases
                         ? 'border-[#f5a623]/50 bg-[#f5a623]/10 text-[#f5a623]'
                         : 'border-[#1a1a28] text-[#6e7681] hover:border-[#f5a623]/30'}">
                ⚡ Não, gerar direto
              </button>
            </div>

            {#if usePhases}
              <div class="flex flex-col gap-3 mb-5">
                {#each phases as phase, i}
                  <div class="p-4 rounded-xl border border-[#1a1a28] bg-[#0d0d18]">
                    <div class="flex items-center justify-between mb-3">
                      <span class="text-xs font-semibold text-[#f5a623]/70 tracking-wider uppercase">
                        Fase {i + 1}
                      </span>
                      <button on:click={() => phases = phases.filter((_, j) => j !== i)}
                        class="text-[#3a3a50] hover:text-[#f85149] text-xs transition-colors">
                        ✕ remover
                      </button>
                    </div>
                    <input bind:value={phase.titulo} placeholder="Título da fase"
                      class="w-full px-3 py-2 mb-2 rounded-lg border border-[#1a1a28]
                             bg-[#0a0a12] text-sm text-[#c9d1d9] placeholder-[#3a3a50]
                             focus:outline-none focus:border-[#f5a623]/50 transition-colors" />
                    <textarea bind:value={phase.descricao} placeholder="O que deve ser feito nesta fase..." rows="2"
                      class="w-full px-3 py-2 rounded-lg border border-[#1a1a28]
                             bg-[#0a0a12] text-sm text-[#c9d1d9] placeholder-[#3a3a50] resize-none
                             focus:outline-none focus:border-[#f5a623]/50 transition-colors" />
                  </div>
                {/each}
                <button on:click={() => phases = [...phases, { titulo: '', descricao: '' }]}
                  class="w-full py-3 rounded-xl border border-dashed border-[#2a2a40]
                         text-[#4a5060] hover:border-[#f5a623]/40 hover:text-[#f5a623]/70
                         text-sm transition-all">
                  + Adicionar fase
                </button>
              </div>
            {/if}

            <button on:click={build} disabled={building}
              class="w-full py-3 rounded-xl font-bold text-sm transition-all active:scale-[0.99]
                     {building
                       ? 'bg-[#f5a623]/40 text-black/50 cursor-wait'
                       : 'bg-[#f5a623] text-black hover:bg-[#e09010]'}">
              {building ? '⏳ Gerando...' : '⚡ Gerar Prompt'}
            </button>
          </div>

        <!-- ============================================================
             TELA 6 — RESULTADO
        ============================================================= -->
        {:else if screen === 'result'}
          <div class="flex flex-col h-full" in:fly={{ y: 16, duration: 200 }}>
            {#if resultError}
              <div class="flex items-start gap-3 p-4 rounded-xl border border-[#f85149]/30 bg-[#f85149]/8 mb-4">
                <span class="text-[#f85149] text-lg">✗</span>
                <div>
                  <p class="text-sm font-semibold text-[#f85149]">Erro ao gerar prompt</p>
                  <p class="text-xs text-[#c9d1d9] mt-1">{resultError}</p>
                </div>
              </div>
            {:else}
              <div class="flex items-center justify-between mb-2">
                <div class="flex items-center gap-2">
                  <span class="text-[#3fb950]">✓</span>
                  <span class="text-sm font-semibold text-[#3fb950]">Prompt gerado com sucesso!</span>
                </div>
                <div class="flex gap-2">
                  <button on:click={saveToLibrary}
                    class="flex items-center gap-1.5 text-xs px-3 py-1.5 rounded-lg border transition-all
                           {saved
                             ? 'border-[#3fb950]/40 bg-[#3fb950]/10 text-[#3fb950]'
                             : 'border-[#1a1a28] text-[#6e7681] hover:border-[#a371f7]/40 hover:text-[#a371f7]'}">
                    {saved ? '✓ Salvo!' : '💾 Salvar'}
                  </button>
                  <button on:click={copyResult}
                    class="flex items-center gap-1.5 text-xs px-3 py-1.5 rounded-lg border transition-all
                           {copied
                             ? 'border-[#3fb950]/40 bg-[#3fb950]/10 text-[#3fb950]'
                             : 'border-[#1a1a28] text-[#6e7681] hover:border-[#f5a623]/40 hover:text-[#f5a623]'}">
                    {copied ? '✓ Copiado!' : '⎘ Copiar'}
                  </button>
                </div>
              </div>
            {/if}

            <pre class="flex-1 min-h-0 overflow-y-auto rounded-xl border border-[#1a1a28]
                        bg-[#0d0d18] px-5 py-4 text-xs text-[#c9d1d9] leading-relaxed
                        whitespace-pre-wrap font-mono">{resultContent}</pre>

            <button on:click={restart}
              class="mt-4 w-full py-2.5 rounded-xl border border-[#1a1a28] text-[#6e7681] text-sm
                     hover:border-[#f5a623]/40 hover:text-[#f5a623] transition-all">
              ← Criar novo prompt
            </button>
          </div>
        {:else if screen === 'prompts'}
          <div class="flex flex-col h-full" in:fly={{ y: 16, duration: 200 }}>
            <!-- header com contador -->
            <div class="flex items-end justify-between mb-1">
              <h2 class="text-lg font-bold">Prompts salvos</h2>
              {#if promptList.length > 0}
                <span class="text-[11px] text-[#4a5060]">
                  {promptList.length} {promptList.length === 1 ? 'prompt' : 'prompts'}
                </span>
              {/if}
            </div>
            <p class="text-sm text-[#6e7681] mb-4">
              Gerencie seus prompts gerados anteriormente.
            </p>

            <div class="flex-1 flex gap-4 min-h-0">
              <!-- coluna esquerda: busca + lista -->
              <div class="w-[40%] flex-shrink-0 flex flex-col gap-3 min-h-0">

                {#if promptList.length > 0}
                  <div class="relative">
                    <input
                      type="text"
                      bind:value={promptSearch}
                      placeholder="Buscar por título…"
                      class="w-full pl-8 pr-3 py-2 rounded-lg text-xs
                             bg-[#0d0d18] border border-[#1a1a28]
                             text-[#c9d1d9] placeholder:text-[#3a3a50]
                             focus:outline-none focus:border-[#a371f7]/50 transition-colors" />
                    <span class="absolute left-2.5 top-1/2 -translate-y-1/2 text-[#3a3a50] text-xs">🔍</span>
                    {#if promptSearch}
                      <button on:click={() => promptSearch = ''}
                        class="absolute right-2 top-1/2 -translate-y-1/2 text-[#3a3a50] hover:text-[#c9d1d9] text-xs">
                        ✕
                      </button>
                    {/if}
                  </div>
                {/if}

                <div class="flex-1 overflow-y-auto rounded-xl border border-[#1a1a28] bg-[#0d0d18] min-h-0">
                  {#if promptList.length === 0}
                    <div class="h-full flex flex-col items-center justify-center gap-2 p-6 text-center">
                      <span class="text-2xl opacity-50">📭</span>
                      <p class="text-xs text-[#6e7681]">Nenhum prompt salvo</p>
                      <p class="text-[10px] text-[#4a5060] leading-relaxed">
                        Salve prompts na tela de resultado para reaproveitá-los depois.
                      </p>
                    </div>
                  {:else if filteredPrompts.length === 0}
                    <div class="h-full flex items-center justify-center p-4 text-xs text-[#4a5060]">
                      Nenhum resultado para "{promptSearch}"
                    </div>
                  {:else}
                    <div class="flex flex-col">
                      {#each filteredPrompts as p}
                        <div class="flex items-center justify-between border-b border-[#0e0e1a] last:border-b-0
                                    transition-all
                                    {promptViewId === p.id ? 'bg-[#a371f7]/10 border-l-2 border-l-[#a371f7]' : 'border-l-2 border-l-transparent hover:bg-[#1a1a28]'}">
                          <button
                            on:click={() => viewPrompt(p.id)}
                            class="flex-1 min-w-0 px-3 py-2.5 text-left text-xs">
                            <p class="truncate {promptViewId === p.id ? 'text-[#a371f7] font-medium' : 'text-[#c9d1d9]'}">
                              {p.titulo || 'Sem título'}
                            </p>
                            <p class="text-[10px] text-[#4a5060]" title={p.data}>{relativeDate(p.data)}</p>
                          </button>

                          {#if promptToDelete === p.id}
                            <div class="flex items-center gap-1 pr-2">
                              <button on:click={() => deletePrompt(p.id)}
                                class="text-[10px] px-2 py-1 rounded bg-[#f85149]/15 border border-[#f85149]/40
                                       text-[#f85149] hover:bg-[#f85149]/25 transition-colors">
                                Confirmar
                              </button>
                              <button on:click={() => promptToDelete = null}
                                class="text-[10px] px-2 py-1 rounded text-[#6e7681] hover:text-[#c9d1d9] transition-colors">
                                Cancelar
                              </button>
                            </div>
                          {:else}
                            <button on:click={() => requestDelete(p.id)}
                              title="Excluir"
                              class="flex-shrink-0 mr-2 p-1 text-[#3a3a50] hover:text-[#f85149] text-xs transition-colors">
                              ✕
                            </button>
                          {/if}
                        </div>
                      {/each}
                    </div>
                  {/if}
                </div>
              </div>

              <!-- coluna direita: visualizador -->
              <div class="flex-1 flex flex-col gap-3 min-h-0">
                {#if promptViewId}
                  <div class="flex justify-end">
                    <button on:click={copyPromptView}
                      class="flex items-center gap-1.5 text-[11px] px-3 py-1.5 rounded-lg border transition-all
                             {promptViewCopied
                               ? 'border-[#3fb950]/40 bg-[#3fb950]/10 text-[#3fb950]'
                               : 'border-[#1a1a28] text-[#6e7681] hover:border-[#f5a623]/40 hover:text-[#f5a623]'}">
                      {promptViewCopied ? '✓ Copiado!' : '⎘ Copiar'}
                    </button>
                  </div>
                  <div class="flex-1 overflow-y-auto rounded-xl border border-[#1a1a28] bg-[#0d0d18] min-h-0">
                    <pre class="p-5 text-xs text-[#c9d1d9] leading-relaxed
                                whitespace-pre-wrap font-mono">{promptViewContent || 'Carregando...'}</pre>
                  </div>
                {:else}
                  <div class="flex-1 flex items-center justify-center p-4 text-xs text-[#4a5060]
                              rounded-xl border border-[#1a1a28] bg-[#0d0d18] min-h-0">
                    Selecione um prompt para visualizar
                  </div>
                {/if}
              </div>
            </div>
          </div>
        {/if}

      </div><!-- /scroll area -->
    </main>

  </div>
</div>
