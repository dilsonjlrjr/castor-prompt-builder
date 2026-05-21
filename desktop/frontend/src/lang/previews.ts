import type { Lang } from '../lib/i18n'

// Dados ilustrativos exibidos nos previews do tutorial. Tudo é fake —
// serve só para mostrar visualmente como cada etapa parece.
export type PreviewData = {
  pipeline: {
    steps:    { icon: string; label: string; accent: string }[]
    features: { label: string; desc: string }[]
  }
  models:   { id: string; color: string; desc: string; tag: string }[]
  roles: {
    chips:    { nome: string; color: string; sel: boolean }[]
    catLabel: string
    grid:     { name: string; selected: boolean }[]
  }
  narrative: {
    paragraph: { type: 'text' | 'highlight'; content: string; color?: string }[]
    hints:     { label: string; val: string; c: string }[]
  }
  gaps: {
    question: string
    askedBy:  string
    answer:   string
    progress: string
  }
  phases: { n: number; title: string; desc: string; done: boolean }[]
  result: {
    title:    string
    roleHead: string
    roleBody: string
    ctxHead:  string
    ctxBody:  string
    skillsHead: string
    skillsBody: string
    copyButton: string
  }
}

const STEPS_BASE = [
  { icon: '📐', accent: '#f5a623' },
  { icon: '🎭', accent: '#a371f7' },
  { icon: '✍️', accent: '#58a6ff' },
  { icon: '💬', accent: '#3fb950' },
  { icon: '📋', accent: '#e06c75' },
  { icon: '🚀', accent: '#f5a623' },
] as const

const MODELS_BASE = [
  { id: 'RACE',   color: '#f5a623' },
  { id: 'RTF',    color: '#3fb950' },
  { id: 'RISEN',  color: '#a371f7' },
  { id: 'CREATE', color: '#58a6ff' },
] as const

const ROLE_CHIPS_BASE = [
  { color: '#a371f7', sel: true  },
  { color: '#a371f7', sel: true  },
  { color: '#a371f7', sel: false },
] as const

const HINT_COLORS = ['#58a6ff', '#3fb950', '#f5a623'] as const

const PREVIEWS: Record<Lang, PreviewData> = {

  // ──────────────────────────────────── PT ─────────────────────────────────
  pt: {
    pipeline: {
      steps: STEPS_BASE.map((s, i) => ({ ...s, label: ['Modelo','Papéis','Tarefa','Contexto','Fases','Prompt'][i] })),
      features: [
        { label: 'Estrutura clara',    desc: 'O modelo organiza as seções do prompt' },
        { label: 'Especialistas',      desc: 'Os papéis definem o estilo e contexto' },
        { label: 'Sem IA no processo', desc: 'Tudo é template engine e heurística' },
      ],
    },
    models: [
      { ...MODELS_BASE[0], desc: 'Contexto rico + entregável claro',    tag: 'Mais usado' },
      { ...MODELS_BASE[1], desc: 'Tarefas diretas e objetivas',          tag: ''           },
      { ...MODELS_BASE[2], desc: 'Steps detalhados com restrições',      tag: 'Complexo'   },
      { ...MODELS_BASE[3], desc: 'Conteúdo criativo com público e tom',  tag: 'Criativo'   },
    ],
    roles: {
      chips: [
        { ...ROLE_CHIPS_BASE[0], nome: 'Arquiteto Cloud' },
        { ...ROLE_CHIPS_BASE[1], nome: 'DevOps Engineer' },
        { ...ROLE_CHIPS_BASE[2], nome: 'QA Lead' },
      ],
      catLabel: '🏗️ Arquitetura',
      grid: [
        { name: 'Arquiteto Cloud',          selected: true  },
        { name: 'Arquiteto de Software',    selected: true  },
        { name: 'Arquiteto de Soluções',    selected: false },
        { name: 'Arquiteto de Microsserviços', selected: false },
      ],
    },
    narrative: {
      paragraph: [
        { type: 'highlight', content: 'Preciso criar uma API de pagamentos', color: '#58a6ff' },
        { type: 'text',      content: ' para o nosso e-commerce. Usamos Go e PostgreSQL, com pico de ' },
        { type: 'highlight', content: '10k req/s', color: '#3fb950' },
        { type: 'text',      content: ' em datas comemorativas. O sistema precisa integrar com o gateway da Stripe e suportar ' },
        { type: 'highlight', content: 'Pix, cartão e boleto', color: '#f5a623' },
        { type: 'text',      content: '. A equipe tem 4 devs sênior.' },
      ],
      hints: [
        { label: 'Stack detectada', val: 'Go + PostgreSQL', c: HINT_COLORS[0] },
        { label: 'Carga',           val: '10k req/s',       c: HINT_COLORS[1] },
        { label: 'Integrações',     val: 'Stripe, Pix',     c: HINT_COLORS[2] },
      ],
    },
    gaps: {
      question: 'Comunicação síncrona (REST/gRPC) ou assíncrona (Kafka/NATS)?',
      askedBy:  '🎭 Arquiteto Cloud',
      answer:   'REST síncrono para pagamentos, Kafka para eventos de confirmação...',
      progress: 'Pergunta 3 de 8',
    },
    phases: [
      { n: 1, title: 'Diagnóstico e Análise', desc: 'Avalie os requisitos de throughput, modelagem do banco e pontos de integração com o gateway.', done: true  },
      { n: 2, title: 'Implementação do Core', desc: 'Construa os endpoints principais de criação e consulta de transações com idempotência.',        done: false },
      { n: 3, title: 'Integrações e Testes',  desc: 'Conecte Stripe, Pix e boleto. Cobertura de testes de contrato e carga.',                          done: false },
    ],
    result: {
      title:      '# Prompt — Arquiteto Cloud',
      roleHead:   '## Papel',
      roleBody:   'Você é um Arquiteto Cloud sênior. Especialista em sistemas de alta disponibilidade, multi-region e otimização de custos em AWS/GCP/Azure...',
      ctxHead:    '## Contexto',
      ctxBody:    'API de pagamentos Go + PostgreSQL. Pico 10k req/s. Integração Stripe, Pix e boleto. Time de 4 devs sênior...',
      skillsHead: '## Habilidades relevantes',
      skillsBody: '- AWS / GCP / Azure\n- Infraestrutura como código\n- Kubernetes...',
      copyButton: '📋 Copiar prompt',
    },
  },

  // ──────────────────────────────────── EN ─────────────────────────────────
  en: {
    pipeline: {
      steps: STEPS_BASE.map((s, i) => ({ ...s, label: ['Model','Roles','Task','Context','Phases','Prompt'][i] })),
      features: [
        { label: 'Clear structure',  desc: 'The model organizes the prompt sections' },
        { label: 'Specialists',      desc: 'Roles define the style and context' },
        { label: 'No AI involved',   desc: 'Pure template engine and heuristics' },
      ],
    },
    models: [
      { ...MODELS_BASE[0], desc: 'Rich context + clear deliverable',    tag: 'Most used' },
      { ...MODELS_BASE[1], desc: 'Direct and objective tasks',          tag: ''          },
      { ...MODELS_BASE[2], desc: 'Detailed steps with constraints',     tag: 'Complex'   },
      { ...MODELS_BASE[3], desc: 'Creative content with audience/tone', tag: 'Creative'  },
    ],
    roles: {
      chips: [
        { ...ROLE_CHIPS_BASE[0], nome: 'Cloud Architect'  },
        { ...ROLE_CHIPS_BASE[1], nome: 'DevOps Engineer'  },
        { ...ROLE_CHIPS_BASE[2], nome: 'QA Lead'          },
      ],
      catLabel: '🏗️ Architecture',
      grid: [
        { name: 'Cloud Architect',         selected: true  },
        { name: 'Software Architect',      selected: true  },
        { name: 'Solutions Architect',     selected: false },
        { name: 'Microservices Architect', selected: false },
      ],
    },
    narrative: {
      paragraph: [
        { type: 'highlight', content: 'I need to build a payments API', color: '#58a6ff' },
        { type: 'text',      content: ' for our e-commerce. We use Go and PostgreSQL, with peaks of ' },
        { type: 'highlight', content: '10k req/s', color: '#3fb950' },
        { type: 'text',      content: ' on busy dates. The system must integrate with the Stripe gateway and support ' },
        { type: 'highlight', content: 'cards, Pix and bank slips', color: '#f5a623' },
        { type: 'text',      content: '. The team has 4 senior devs.' },
      ],
      hints: [
        { label: 'Detected stack', val: 'Go + PostgreSQL', c: HINT_COLORS[0] },
        { label: 'Load',           val: '10k req/s',       c: HINT_COLORS[1] },
        { label: 'Integrations',   val: 'Stripe, Pix',     c: HINT_COLORS[2] },
      ],
    },
    gaps: {
      question: 'Synchronous communication (REST/gRPC) or asynchronous (Kafka/NATS)?',
      askedBy:  '🎭 Cloud Architect',
      answer:   'Synchronous REST for payments, Kafka for confirmation events...',
      progress: 'Question 3 of 8',
    },
    phases: [
      { n: 1, title: 'Diagnosis and Analysis', desc: 'Assess throughput requirements, data modeling and integration points with the gateway.', done: true  },
      { n: 2, title: 'Core Implementation',    desc: 'Build the main endpoints for transaction creation and lookup with idempotency.',         done: false },
      { n: 3, title: 'Integrations and Tests', desc: 'Wire up Stripe, Pix and bank slips. Contract and load test coverage.',                    done: false },
    ],
    result: {
      title:      '# Prompt — Cloud Architect',
      roleHead:   '## Role',
      roleBody:   'You are a senior Cloud Architect. Specialist in high-availability systems, multi-region setups and cost optimization on AWS/GCP/Azure...',
      ctxHead:    '## Context',
      ctxBody:    'Payments API in Go + PostgreSQL. Peak 10k req/s. Stripe, Pix and bank slip integration. Team of 4 senior devs...',
      skillsHead: '## Relevant skills',
      skillsBody: '- AWS / GCP / Azure\n- Infrastructure as code\n- Kubernetes...',
      copyButton: '📋 Copy prompt',
    },
  },

  // ──────────────────────────────────── ES ─────────────────────────────────
  es: {
    pipeline: {
      steps: STEPS_BASE.map((s, i) => ({ ...s, label: ['Modelo','Roles','Tarea','Contexto','Fases','Prompt'][i] })),
      features: [
        { label: 'Estructura clara', desc: 'El modelo organiza las secciones del prompt' },
        { label: 'Especialistas',    desc: 'Los roles definen el estilo y el contexto' },
        { label: 'Sin IA en medio',  desc: 'Solo motor de plantillas y heurísticas' },
      ],
    },
    models: [
      { ...MODELS_BASE[0], desc: 'Contexto rico + entregable claro',     tag: 'Más usado' },
      { ...MODELS_BASE[1], desc: 'Tareas directas y objetivas',          tag: ''          },
      { ...MODELS_BASE[2], desc: 'Pasos detallados con restricciones',   tag: 'Complejo'  },
      { ...MODELS_BASE[3], desc: 'Contenido creativo con público y tono',tag: 'Creativo'  },
    ],
    roles: {
      chips: [
        { ...ROLE_CHIPS_BASE[0], nome: 'Arquitecto Cloud' },
        { ...ROLE_CHIPS_BASE[1], nome: 'DevOps Engineer'  },
        { ...ROLE_CHIPS_BASE[2], nome: 'QA Lead'          },
      ],
      catLabel: '🏗️ Arquitectura',
      grid: [
        { name: 'Arquitecto Cloud',           selected: true  },
        { name: 'Arquitecto de Software',     selected: true  },
        { name: 'Arquitecto de Soluciones',   selected: false },
        { name: 'Arquitecto de Microservicios', selected: false },
      ],
    },
    narrative: {
      paragraph: [
        { type: 'highlight', content: 'Necesito crear una API de pagos', color: '#58a6ff' },
        { type: 'text',      content: ' para nuestro e-commerce. Usamos Go y PostgreSQL, con picos de ' },
        { type: 'highlight', content: '10k req/s', color: '#3fb950' },
        { type: 'text',      content: ' en fechas señaladas. El sistema debe integrarse con la pasarela de Stripe y soportar ' },
        { type: 'highlight', content: 'tarjeta, Pix y boleta', color: '#f5a623' },
        { type: 'text',      content: '. El equipo tiene 4 devs sénior.' },
      ],
      hints: [
        { label: 'Stack detectado', val: 'Go + PostgreSQL', c: HINT_COLORS[0] },
        { label: 'Carga',           val: '10k req/s',       c: HINT_COLORS[1] },
        { label: 'Integraciones',   val: 'Stripe, Pix',     c: HINT_COLORS[2] },
      ],
    },
    gaps: {
      question: '¿Comunicación síncrona (REST/gRPC) o asíncrona (Kafka/NATS)?',
      askedBy:  '🎭 Arquitecto Cloud',
      answer:   'REST síncrono para pagos, Kafka para eventos de confirmación...',
      progress: 'Pregunta 3 de 8',
    },
    phases: [
      { n: 1, title: 'Diagnóstico y Análisis', desc: 'Evalúa los requisitos de throughput, modelado de datos y puntos de integración con la pasarela.', done: true  },
      { n: 2, title: 'Implementación del Core', desc: 'Construye los endpoints principales de creación y consulta de transacciones con idempotencia.',  done: false },
      { n: 3, title: 'Integraciones y Pruebas', desc: 'Conecta Stripe, Pix y boleta. Cobertura de pruebas de contrato y carga.',                          done: false },
    ],
    result: {
      title:      '# Prompt — Arquitecto Cloud',
      roleHead:   '## Rol',
      roleBody:   'Eres un Arquitecto Cloud sénior. Especialista en sistemas de alta disponibilidad, multi-región y optimización de costos en AWS/GCP/Azure...',
      ctxHead:    '## Contexto',
      ctxBody:    'API de pagos Go + PostgreSQL. Pico 10k req/s. Integración con Stripe, Pix y boleta. Equipo de 4 devs sénior...',
      skillsHead: '## Habilidades relevantes',
      skillsBody: '- AWS / GCP / Azure\n- Infraestructura como código\n- Kubernetes...',
      copyButton: '📋 Copiar prompt',
    },
  },
}

export function getPreviewData(l: Lang): PreviewData {
  return PREVIEWS[l] ?? PREVIEWS.pt
}
