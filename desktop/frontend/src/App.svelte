<script lang="ts">
  import { onMount } from 'svelte'
  import { fly } from 'svelte/transition'
  import { GetModels, GetRoles, BuildPrompt } from '../wailsjs/go/main/App.js'
  import { main } from '../wailsjs/go/models'

  // ---- tipos ----
  type Model  = { id: string; nome: string; descricao: string }
  type Role   = { id: string; nome: string; categoria: string; gaps_comuns: string[] }
  type Step   = { titulo: string; descricao: string }
  type Screen = 'model' | 'role' | 'narrative' | 'gap' | 'phase' | 'result'

  // ---- estado ----
  let models:  Model[] = []
  let roles:   Role[]  = []
  let screen:  Screen  = 'model'

  let selectedModel: Model | null = null
  let selectedRoles: Set<string>  = new Set()
  let roleSearch = ''

  let narrative   = ''
  let gaps: string[] = []
  let gapIndex    = 0
  let gapAnswers: string[] = []

  let usePhases = false
  let phases: Step[] = []

  let resultContent = ''
  let resultPath    = ''
  let resultError   = ''
  let building      = false
  let copied        = false

  // ---- ciclo de vida ----
  onMount(async () => {
    models = await GetModels()
    roles  = await GetRoles()
  })

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
  }
  const CAT_LABEL: Record<string, string> = {
    arquitetura: 'Arquitetura',   frontend:  'Frontend & Mobile',
    backend:     'Backend',       devops:    'DevOps & Cloud',
    banco:       'Banco de Dados',dados:     'Dados & IA',
    gestao:      'Gestão',        seguranca: 'Segurança',
    design:      'Design',        marketing: 'Marketing',
  }
  function catLabel(cat: string) { return CAT_LABEL[cat] ?? cat }
  function catIcon(cat: string)  { return CAT_ICON[cat]  ?? '📁' }

  // modelo: badge colorida por ID
  const MODEL_COLOR: Record<string, string> = {
    race:   '#f5a623', rtf:    '#3fb950',
    risen:  '#a371f7', create: '#58a6ff',
  }
  function modelColor(id: string) { return MODEL_COLOR[id] ?? '#6e7681' }

  // ---- navegação ----
  function goBack() {
    if (screen === 'role')      { screen = 'model';     return }
    if (screen === 'narrative') { screen = 'role';      return }
    if (screen === 'gap')       { prevGap();            return }
    if (screen === 'phase')     { screen = gaps.length > 0 ? 'gap' : 'narrative'; gapIndex = gaps.length - 1; return }
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
    const allGaps = [...selectedRoles]
      .flatMap(id => roles.find(r => r.id === id)?.gaps_comuns ?? [])
    gaps      = [...new Set(allGaps)]
    gapIndex  = 0
    gapAnswers = Array(gaps.length).fill('')
    screen = gaps.length > 0 ? 'gap' : 'phase'
  }

  function nextGap() {
    gapIndex++
    if (gapIndex >= gaps.length) screen = 'phase'
  }

  function prevGap() {
    if (gapIndex > 0) { gapIndex--; return }
    screen = 'narrative'
  }

  async function build() {
    building = true
    const result = await BuildPrompt(main.BuildRequestDTO.createFrom({
      model_id:    selectedModel!.id,
      role_ids:    [...selectedRoles],
      narrativa:   narrative,
      gap_answers: gaps.map((p, i) => ({ pergunta: p, resposta: gapAnswers[i] ?? '' })),
      steps:       usePhases ? phases : [],
    }))
    resultContent = result.conteudo
    resultPath    = result.caminho
    resultError   = result.erro ?? ''
    building      = false
    screen        = 'result'
  }

  async function copyResult() {
    await navigator.clipboard.writeText(resultContent)
    copied = true
    setTimeout(() => copied = false, 2000)
  }

  function restart() {
    selectedModel = null
    selectedRoles = new Set()
    narrative     = ''
    gaps          = []
    gapAnswers    = []
    phases        = []
    usePhases     = false
    resultContent = ''
    resultPath    = ''
    resultError   = ''
    screen        = 'model'
  }
</script>

<!-- ============================================================
     SHELL
============================================================= -->
<div class="flex flex-col h-screen bg-[#0a0a12] text-[#c9d1d9] font-mono select-none overflow-hidden">

  <!-- titlebar drag area -->
  <div class="h-8 flex-shrink-0" style="--wails-draggable:drag" />

  <!-- main content -->
  <div class="flex-1 flex overflow-hidden">

    <!-- ── SIDEBAR esquerda ── -->
    <aside class="w-52 flex-shrink-0 flex flex-col items-center pt-6 pb-8 px-4
                  border-r border-[#1a1a28] bg-[#0d0d18]">

      <!-- mascote -->
      <img src="/castor.png" alt="Castor"
           class="w-24 h-24 object-contain mb-3
                  drop-shadow-[0_0_18px_rgba(245,166,35,0.5)]" />

      <div class="text-[10px] tracking-[0.3em] text-[#4a5060] uppercase mb-0.5">Prompt Builder</div>
      <h1 class="text-base font-bold tracking-wider mb-8">
        <span class="text-[#f5a623]">CASTOR</span>
        <span class="text-[#e06b2e]"> BUILDER</span>
      </h1>

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

      <!-- versão -->
      <div class="mt-auto text-[10px] text-[#2a2a40]">v0.1.0</div>
    </aside>

    <!-- ── CONTEÚDO PRINCIPAL ── -->
    <main class="flex-1 flex flex-col overflow-hidden">

      <!-- topbar: voltar + breadcrumb -->
      {#if screen !== 'model' && screen !== 'result'}
        <div class="flex items-center gap-3 px-8 pt-6 pb-0 flex-shrink-0">
          <button on:click={goBack}
            class="flex items-center gap-1.5 text-xs text-[#4a5060]
                   hover:text-[#f5a623] transition-colors group">
            <span class="text-base leading-none group-hover:-translate-x-0.5 transition-transform">←</span>
            <span>voltar</span>
          </button>
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
        </div>
      {/if}

      <!-- scroll area -->
      <div class="flex-1 overflow-y-auto px-8 py-6">

        <!-- ============================================================
             TELA 1 — MODELO
        ============================================================= -->
        {#if screen === 'model'}
          <div in:fly={{ y: 16, duration: 200 }}>
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
             TELA 4 — GAPS
        ============================================================= -->
        {:else if screen === 'gap'}
          <div in:fly={{ y: 16, duration: 200 }}>
            <!-- progress bar -->
            <div class="flex items-center gap-3 mb-6">
              <div class="flex-1 h-1 rounded-full bg-[#1a1a28] overflow-hidden">
                <div class="h-full bg-[#f5a623] rounded-full transition-all duration-500"
                     style="width: {((gapIndex + 1) / gaps.length) * 100}%"></div>
              </div>
              <span class="text-xs text-[#6e7681] flex-shrink-0">
                {gapIndex + 1} / {gaps.length}
              </span>
            </div>

            <h2 class="text-lg font-bold mb-1">Contexto adicional</h2>
            <p class="text-sm text-[#6e7681] mb-5">
              Preencha o que souber — lacunas vazias virarão avisos no prompt final.
            </p>

            <!-- pergunta com destaque -->
            <div class="flex gap-3 p-4 rounded-xl border border-[#f5a623]/20 bg-[#f5a623]/5 mb-4">
              <span class="text-[#f5a623] text-lg flex-shrink-0">?</span>
              <p class="text-sm text-[#e0e6f0] leading-relaxed">{gaps[gapIndex]}</p>
            </div>

            <textarea bind:value={gapAnswers[gapIndex]}
              placeholder="Digite sua resposta... (deixe vazio para pular)"
              rows="5"
              class="w-full px-4 py-3 rounded-xl border border-[#1a1a28]
                     bg-[#0d0d18] text-[#c9d1d9] placeholder-[#3a3a50] text-sm
                     resize-none focus:outline-none focus:border-[#f5a623]/60
                     transition-colors leading-relaxed" />

            <div class="flex justify-end mt-4">
              <button on:click={nextGap}
                class="px-6 py-2.5 rounded-xl bg-[#f5a623] text-black font-bold text-sm
                       hover:bg-[#e09010] transition-all active:scale-[0.98]">
                {gapIndex < gaps.length - 1 ? 'Próxima →' : 'Continuar →'}
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
                <button on:click={copyResult}
                  class="flex items-center gap-1.5 text-xs px-3 py-1.5 rounded-lg border transition-all
                         {copied
                           ? 'border-[#3fb950]/40 bg-[#3fb950]/10 text-[#3fb950]'
                           : 'border-[#1a1a28] text-[#6e7681] hover:border-[#f5a623]/40 hover:text-[#f5a623]'}">
                  {copied ? '✓ Copiado!' : '⎘ Copiar'}
                </button>
              </div>
              {#if resultPath}
                <p class="text-[10px] text-[#3a3a50] mb-3 truncate">📄 {resultPath}</p>
              {/if}
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
        {/if}

      </div><!-- /scroll area -->
    </main>

  </div>
</div>
