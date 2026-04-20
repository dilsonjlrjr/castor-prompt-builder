/** @type {import('tailwindcss').Config} */
export default {
  content: ['./src/**/*.{svelte,ts,js,html}'],
  darkMode: 'class',
  theme: {
    extend: {
      colors: {
        castor: {
          bg:      '#0d0d14',
          surface: '#13131f',
          border:  '#1e1e30',
          primary: '#f5a623',
          accent:  '#e06b2e',
          text:    '#c9d1d9',
          muted:   '#6e7681',
          success: '#3fb950',
          error:   '#f85149',
        },
      },
      fontFamily: {
        mono: ['JetBrains Mono', 'Fira Code', 'ui-monospace', 'monospace'],
      },
    },
  },
  plugins: [],
}
