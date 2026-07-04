import type { ComputedRef, Ref } from 'vue'
import { normalizeLocaleCode } from '~/utils/locale'

interface NuxtI18nComposer {
  locale: Ref<string>
  localeCodes: ComputedRef<string[]>
  defaultLocale: string
  setLocale: (locale: string) => Promise<void>
}

export default defineNuxtPlugin({
  name: 'nextpt:locale-normalize',
  dependsOn: ['i18n:plugin'],
  async setup(nuxtApp) {
    const i18n = nuxtApp.$i18n as NuxtI18nComposer
    const supportedCodes = [...i18n.localeCodes.value]
    const normalized = normalizeLocaleCode(i18n.locale.value, supportedCodes, i18n.defaultLocale)

    if (i18n.locale.value === normalized) return

    await i18n.setLocale(normalized)
  }
})
