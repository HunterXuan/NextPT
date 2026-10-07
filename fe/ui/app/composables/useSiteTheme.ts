import { isSiteTheme, resolveSiteTheme, siteThemeStorageKey, type SiteTheme } from '~/utils/siteTheme'

export function useSiteTheme() {
  const config = useRuntimeConfig()
  const defaultTheme = computed(() => resolveSiteTheme(config.public.siteTheme))
  const preference = useState<SiteTheme | 'default'>('site-theme-preference', () => 'default')
  const currentTheme = computed(() => preference.value === 'default' ? defaultTheme.value : preference.value)

  function applyTheme(value: unknown) {
    preference.value = isSiteTheme(value) ? value : 'default'
    if (import.meta.client) document.documentElement.dataset.siteTheme = currentTheme.value
  }

  function setTheme(value: SiteTheme | 'default') {
    applyTheme(value)
    if (!import.meta.client) return
    try {
      if (preference.value === 'default') localStorage.removeItem(siteThemeStorageKey)
      else localStorage.setItem(siteThemeStorageKey, preference.value)
    } catch {
      // Private browsers may disable storage; the current page can still change theme.
    }
  }

  function restoreTheme() {
    if (!import.meta.client) return
    try {
      applyTheme(localStorage.getItem(siteThemeStorageKey))
    } catch {
      applyTheme('default')
    }
  }

  return { preference: readonly(preference), currentTheme, defaultTheme, setTheme, restoreTheme }
}
