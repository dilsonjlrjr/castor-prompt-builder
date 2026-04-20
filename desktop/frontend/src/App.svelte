<script lang="ts">
  import { onMount } from 'svelte'
  import { GetModels, GetRoles, BuildPrompt } from '../wailsjs/go/main/App.js'

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

  let narrative = ''
  let gaps: string[] = []
  let gapIndex = 0
  let gapAnswers: string[] = []

  let usePhases = false
  let phases: Step[] = []

  let resultContent = ''
  let resultPath    = ''
  let resultError   = ''
  let building      = false

  // ---- ciclo de vida ----
  onMount(async () => {
    models = await GetModels()
    roles  = await GetRoles()
  })

  // ---- helpers ----
  $: categories = [...new Set(roles.map(r => r.categoria))].sort()

  $: filteredRoles = roleSearch.trim()
    ? roles.filter(r =>
        r.nome.toLowerCase().includes(roleSearch.toLowerCase()) ||
        r.categoria.toLowerCase().includes(roleSearch.toLowerCase()))
    : roles

  function rolesInCat(cat: string) {
    return filteredRoles.filter(r => r.categoria === cat)
  }

  function catLabel(cat: string): string {
    const map: Record<string, string> = {
      arquitetura: 'Arquitetura',
      frontend:    'Frontend & Mobile',
      backend:     'Backend',
      devops:      'DevOps & Cloud',
      banco:       'Banco de Dados',
      dados:       'Dados & IA',
      gestao:      'Gestão',
      seguranca:   'Segurança',
      design:      'Design',
      marketing:   'Marketing',
    }
    return map[cat] ?? cat
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
    const result = await BuildPrompt({
      model_id:    selectedModel!.id,
      role_ids:    [...selectedRoles],
      narrativa:   narrative,
      gap_answers: gaps.map((p, i) => ({ pergunta: p, resposta: gapAnswers[i] ?? '' })),
      steps:       usePhases ? phases : [],
    })
    resultContent = result.conteudo
    resultPath    = result.caminho
    resultError   = result.erro ?? ''
    building      = false
    screen        = 'result'
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

<!-- ====================================================
     LAYOUT PRINCIPAL
===================================================== -->
<div class="flex flex-col h-screen bg-[#0d0d14] text-[#c9d1d9] font-mono select-none">

  <!-- titlebar drag area -->
  <div class="h-8 flex-shrink-0" style="--wails-draggable:drag" />

  <!-- conteúdo centralizado -->
  <div class="flex-1 flex flex-col items-center justify-center px-8 pb-8 overflow-hidden">

    <!-- cabeçalho -->
    <div class="mb-8 text-center">
      <div class="text-xs tracking-[0.3em] text-[#6e7681] uppercase mb-1">Prompt Builder</div>
      <h1 class="text-2xl font-bold tracking-wide">
        <span class="text-[#f5a623]">CASTOR</span>
        <span class="text-[#e06b2e] ml-2">BUILDER</span>
      </h1>
    </div>

    <!-- ================================================
         TELA: SELECIONAR MODELO
    ================================================= -->
    {#if screen === 'model'}
      <p class="text-[#6e7681] text-sm mb-6 text-center max-w-md">
        Construa prompts estruturados para LLMs em segundos.<br>
        Escolha um framework, descreva sua tarefa e o CASTOR monta o prompt ideal.
      </p>

      <p class="text-[#f5a623] text-sm mb-4">Selecione o modelo de prompt:</p>

      <div class="flex flex-col gap-2 w-full max-w-lg">
        {#each models as m}
          <button
            on:click={() => confirmModel(m)}
            class="flex items-start gap-4 px-4 py-3 rounded-lg border border-[#1e1e30]
                   bg-[#13131f] hover:border-[#f5a623] hover:bg-[#1a1a28]
                   transition-all text-left group"
          >
            <span class="text-[#f5a623] font-bold w-16 flex-shrink-0 group-hover:text-[#f5a623]">
              {m.nome}
            </span>
            <span class="text-[#6e7681] text-sm">{m.descricao}</span>
          </button>
        {/each}
      </div>

    <!-- ================================================
         TELA: SELECIONAR PAPEL
    ================================================= -->
    {:else if screen === 'role'}
      <div class="flex items-center gap-3 mb-5 w-full max-w-lg">
        <button on:click={() => screen = 'model'} class="text-[#6e7681] hover:text-[#c9d1d9] text-sm">← voltar</button>
        <span class="text-[#f5a623] text-sm font-bold">{selectedModel?.nome}</span>
        {#if selectedRoles.size > 0}
          <span class="ml-auto text-xs bg-[#e06b2e] text-white px-2 py-0.5 rounded-full">
            {selectedRoles.size} selecionado(s)
          </span>
        {/if}
      </div>

      <p class="text-[#f5a623] text-sm mb-3 w-full max-w-lg">Selecione o(s) papel(eis):</p>

      <!-- busca -->
      <input
        bind:value={roleSearch}
        placeholder="🔍 Buscar papel..."
        class="w-full max-w-lg px-3 py-2 mb-4 rounded-lg border border-[#1e1e30]
               bg-[#13131f] text-[#c9d1d9] placeholder-[#6e7681] text-sm
               focus:outline-none focus:border-[#f5a623] transition-colors"
      />

      <!-- lista com scroll -->
      <div class="w-full max-w-lg flex-1 overflow-y-auto min-h-0">
        {#each categories as cat}
          {#if rolesInCat(cat).length > 0}
            <div class="text-[#6e7681] text-xs tracking-widest uppercase mt-4 mb-1 px-1">
              {catLabel(cat)}
            </div>
            {#each rolesInCat(cat) as role}
              <button
                on:click={() => toggleRole(role.id)}
                class="flex items-center gap-3 w-full px-3 py-2 rounded-md mb-1
                       hover:bg-[#1a1a28] transition-colors text-left text-sm
                       {selectedRoles.has(role.id) ? 'text-[#f5a623]' : 'text-[#c9d1d9]'}"
              >
                <span class="w-4 text-center flex-shrink-0">
                  {selectedRoles.has(role.id) ? '✓' : '○'}
                </span>
                {role.nome}
              </button>
            {/each}
          {/if}
        {/each}
      </div>

      <button
        on:click={confirmRoles}
        disabled={selectedRoles.size === 0}
        class="mt-5 px-6 py-2 rounded-lg bg-[#f5a623] text-black font-bold text-sm
               hover:bg-[#e09010] disabled:opacity-30 disabled:cursor-not-allowed transition-colors"
      >
        Confirmar →
      </button>

    <!-- ================================================
         TELA: NARRATIVA
    ================================================= -->
    {:else if screen === 'narrative'}
      <div class="flex items-center gap-3 mb-5 w-full max-w-lg">
        <button on:click={() => screen = 'role'} class="text-[#6e7681] hover:text-[#c9d1d9] text-sm">← voltar</button>
        <span class="text-[#f5a623] text-sm font-bold truncate">
          {[...selectedRoles].map(id => roles.find(r => r.id === id)?.nome).join(' + ')}
        </span>
      </div>

      <p class="text-[#f5a623] text-sm mb-3 w-full max-w-lg">Descreva a tarefa livremente:</p>

      <textarea
        bind:value={narrative}
        placeholder="O que você precisa fazer? Descreva com o máximo de contexto possível..."
        rows="8"
        class="w-full max-w-lg px-4 py-3 rounded-lg border border-[#1e1e30]
               bg-[#13131f] text-[#c9d1d9] placeholder-[#6e7681] text-sm resize-none
               focus:outline-none focus:border-[#f5a623] transition-colors leading-relaxed"
      />

      <button
        on:click={confirmNarrative}
        disabled={!narrative.trim()}
        class="mt-4 px-6 py-2 rounded-lg bg-[#f5a623] text-black font-bold text-sm
               hover:bg-[#e09010] disabled:opacity-30 disabled:cursor-not-allowed transition-colors"
      >
        Continuar →
      </button>

    <!-- ================================================
         TELA: GAPS
    ================================================= -->
    {:else if screen === 'gap'}
      <div class="text-xs text-[#6e7681] mb-2 w-full max-w-lg">
        lacuna {gapIndex + 1} de {gaps.length}
      </div>

      <p class="text-[#f5a623] text-sm mb-4 w-full max-w-lg">{gaps[gapIndex]}</p>

      <textarea
        bind:value={gapAnswers[gapIndex]}
        placeholder="Digite sua resposta... (deixe vazio para pular)"
        rows="4"
        class="w-full max-w-lg px-4 py-3 rounded-lg border border-[#1e1e30]
               bg-[#13131f] text-[#c9d1d9] placeholder-[#6e7681] text-sm resize-none
               focus:outline-none focus:border-[#f5a623] transition-colors leading-relaxed"
      />

      <div class="flex gap-3 mt-4">
        <button on:click={prevGap}
          class="px-4 py-2 rounded-lg border border-[#1e1e30] text-[#6e7681]
                 hover:border-[#f5a623] hover:text-[#c9d1d9] text-sm transition-colors">
          ← voltar
        </button>
        <button on:click={nextGap}
          class="px-6 py-2 rounded-lg bg-[#f5a623] text-black font-bold text-sm
                 hover:bg-[#e09010] transition-colors">
          {gapIndex < gaps.length - 1 ? 'Próxima →' : 'Continuar →'}
        </button>
      </div>

    <!-- ================================================
         TELA: FASES
    ================================================= -->
    {:else if screen === 'phase'}
      <p class="text-[#f5a623] text-sm mb-5 w-full max-w-lg">Deseja definir fases de execução?</p>

      <div class="flex gap-3 mb-6">
        <button
          on:click={() => usePhases = true}
          class="px-5 py-2 rounded-lg border text-sm transition-colors
                 {usePhases ? 'border-[#f5a623] text-[#f5a623] bg-[#1a1a28]' : 'border-[#1e1e30] text-[#6e7681] hover:border-[#f5a623]'}"
        >
          Sim, definir fases
        </button>
        <button
          on:click={() => { usePhases = false; phases = [] }}
          class="px-5 py-2 rounded-lg border text-sm transition-colors
                 {!usePhases ? 'border-[#f5a623] text-[#f5a623] bg-[#1a1a28]' : 'border-[#1e1e30] text-[#6e7681] hover:border-[#f5a623]'}"
        >
          Não, gerar direto
        </button>
      </div>

      {#if usePhases}
        <div class="w-full max-w-lg mb-4">
          {#each phases as phase, i}
            <div class="mb-3 p-3 rounded-lg border border-[#1e1e30] bg-[#13131f]">
              <div class="flex items-center gap-2 mb-2">
                <span class="text-[#6e7681] text-xs">Fase {i + 1}</span>
                <button on:click={() => phases = phases.filter((_, j) => j !== i)}
                  class="ml-auto text-[#6e7681] hover:text-[#f85149] text-xs">✕</button>
              </div>
              <input bind:value={phase.titulo} placeholder="Título da fase"
                class="w-full px-3 py-1.5 mb-2 rounded border border-[#1e1e30]
                       bg-[#0d0d14] text-sm text-[#c9d1d9] placeholder-[#6e7681]
                       focus:outline-none focus:border-[#f5a623]" />
              <textarea bind:value={phase.descricao} placeholder="Descrição..." rows="2"
                class="w-full px-3 py-1.5 rounded border border-[#1e1e30]
                       bg-[#0d0d14] text-sm text-[#c9d1d9] placeholder-[#6e7681] resize-none
                       focus:outline-none focus:border-[#f5a623]" />
            </div>
          {/each}
          <button on:click={() => phases = [...phases, { titulo: '', descricao: '' }]}
            class="w-full py-2 rounded-lg border border-dashed border-[#1e1e30]
                   text-[#6e7681] hover:border-[#f5a623] hover:text-[#f5a623] text-sm transition-colors">
            + Adicionar fase
          </button>
        </div>
      {/if}

      <button on:click={build} disabled={building}
        class="px-6 py-2 rounded-lg bg-[#f5a623] text-black font-bold text-sm
               hover:bg-[#e09010] disabled:opacity-50 transition-colors">
        {building ? 'Gerando...' : '⚡ Gerar Prompt'}
      </button>

    <!-- ================================================
         TELA: RESULTADO
    ================================================= -->
    {:else if screen === 'result'}
      {#if resultError}
        <div class="text-[#f85149] text-sm mb-4">✗ Erro: {resultError}</div>
      {:else}
        <div class="text-[#3fb950] text-sm mb-1 w-full max-w-2xl">✓ Prompt gerado com sucesso!</div>
        <div class="text-[#6e7681] text-xs mb-4 w-full max-w-2xl">{resultPath}</div>
      {/if}

      <pre class="w-full max-w-2xl flex-1 min-h-0 overflow-y-auto
                  bg-[#13131f] border border-[#1e1e30] rounded-lg
                  px-4 py-3 text-xs text-[#c9d1d9] leading-relaxed whitespace-pre-wrap"
      >{resultContent}</pre>

      <button on:click={restart}
        class="mt-5 px-6 py-2 rounded-lg border border-[#1e1e30] text-[#6e7681]
               hover:border-[#f5a623] hover:text-[#f5a623] text-sm transition-colors">
        ← Novo prompt
      </button>
    {/if}

  </div>
</div>
