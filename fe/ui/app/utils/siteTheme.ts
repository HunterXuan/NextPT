export const siteThemes = [
  { id: 'teal', primary: 'teal', neutral: 'gray', swatch: '#0d9488' },
  { id: 'green', primary: 'green', neutral: 'zinc', swatch: '#16a34a' },
  { id: 'sky', primary: 'sky', neutral: 'slate', swatch: '#0284c7' },
  { id: 'indigo', primary: 'indigo', neutral: 'slate', swatch: '#4f46e5' },
  { id: 'violet', primary: 'violet', neutral: 'zinc', swatch: '#7c3aed' },
  { id: 'rose', primary: 'rose', neutral: 'zinc', swatch: '#e11d48' },
  { id: 'amber', primary: 'amber', neutral: 'stone', swatch: '#b45309' },
  { id: 'graphite', primary: 'zinc', neutral: 'zinc', swatch: '#52525b' }
] as const

export type SiteTheme = typeof siteThemes[number]['id']
export const siteThemeStorageKey = 'nextpt.site-theme'

export function isSiteTheme(value: unknown): value is SiteTheme {
  return siteThemes.some(theme => theme.id === value)
}

export function resolveSiteTheme(value: unknown): SiteTheme {
  return isSiteTheme(value) ? value : 'sky'
}

const shades = [50, 100, 200, 300, 400, 500, 600, 700, 800, 900, 950]

// Root variables also reach teleported dialogs and cached workspace pages.
export const siteThemeStyles = siteThemes.map(theme => {
  const variables = shades.flatMap(shade => [
    `--ui-color-primary-${shade}:var(--color-${theme.primary}-${shade});`,
    `--ui-color-neutral-${shade}:var(--color-${theme.neutral}-${shade});`
  ]).join('')
  const lightShade = theme.id === 'graphite' ? 900 : ['teal', 'green', 'amber'].includes(theme.id) ? 700 : 600
  const darkShade = theme.id === 'graphite' ? 100 : 400
  return `html[data-site-theme="${theme.id}"]{${variables}--ui-primary:var(--ui-color-primary-${lightShade});}
html.dark[data-site-theme="${theme.id}"]{--ui-primary:var(--ui-color-primary-${darkShade});}`
}).join('\n')

export function siteThemeBootstrap(defaultTheme: SiteTheme): string {
  return `(function(){var theme=${JSON.stringify(defaultTheme)};try{var saved=localStorage.getItem(${JSON.stringify(siteThemeStorageKey)});if(${JSON.stringify(siteThemes.map(theme => theme.id))}.includes(saved))theme=saved;}catch(e){}document.documentElement.dataset.siteTheme=theme;})();`
}
