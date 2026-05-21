import { writable, derived, get } from 'svelte/store'
import { pt } from '../lang/pt'
import { en } from '../lang/en'
import { es } from '../lang/es'

export type Lang = 'pt' | 'en' | 'es'

export const LANG_META: Record<Lang, { flag: string; label: string; nativeLabel: string }> = {
  pt: { flag: '🇧🇷', label: 'Portuguese', nativeLabel: 'Português' },
  en: { flag: '🇺🇸', label: 'English',    nativeLabel: 'English'    },
  es: { flag: '🇪🇸', label: 'Spanish',    nativeLabel: 'Español'    },
}

const DICTS = { pt, en, es } as const
const STORAGE_KEY = 'castor_lang'

// Detecta o idioma: 1) localStorage, 2) navigator.language, 3) PT fallback.
function detectLang(): Lang {
  try {
    const stored = localStorage.getItem(STORAGE_KEY) as Lang | null
    if (stored && stored in DICTS) return stored
  } catch {}
  const browser = (typeof navigator !== 'undefined' ? navigator.language : 'pt').toLowerCase()
  if (browser.startsWith('en')) return 'en'
  if (browser.startsWith('es')) return 'es'
  return 'pt'
}

export const lang = writable<Lang>(detectLang())

// Persiste a escolha do idioma sempre que muda.
lang.subscribe(l => {
  try { localStorage.setItem(STORAGE_KEY, l) } catch {}
})

// Resolve "tutorial.slide1.title" em dicionário aninhado.
function resolve(dict: any, key: string): string | undefined {
  let cur: any = dict
  for (const part of key.split('.')) {
    if (cur == null || typeof cur !== 'object') return undefined
    cur = cur[part]
  }
  return typeof cur === 'string' ? cur : undefined
}

// Interpola {variaveis} numa string.
function interpolate(s: string, vars?: Record<string, string | number>): string {
  if (!vars) return s
  return s.replace(/\{(\w+)\}/g, (_, k) => (k in vars ? String(vars[k]) : `{${k}}`))
}

// Store derivada com a função de tradução. Reativa: chamadas em template
// (`$t('chave')`) atualizam automaticamente ao trocar de idioma.
export const t = derived(lang, ($lang) => {
  const dict = DICTS[$lang] ?? DICTS.pt
  return (key: string, vars?: Record<string, string | number>): string => {
    const v = resolve(dict, key) ?? resolve(DICTS.pt, key) ?? key
    return interpolate(v, vars)
  }
})

export function setLang(l: Lang) { lang.set(l) }
export function getLang(): Lang { return get(lang) }
