<template>
  <UDropdownMenu :items="languageItems" :content="{ align: 'end' }">
    <UButton
      color="neutral"
      variant="ghost"
      icon="i-lucide-languages"
      :aria-label="$t('common.switchLanguage')"
    />
  </UDropdownMenu>
</template>

<script setup lang="ts">
import { normalizeLocaleCode } from '~/utils/locale'

const { locale, locales, setLocale } = useI18n()

const languageItems = computed(() => [
  (locales.value as any[]).map((item) => ({
    label: item.name,
    onSelect: () => handleLanguageSwitch(item.code)
  }))
])

async function handleLanguageSwitch(code: string) {
  const supportedCodes = (locales.value as any[]).map((item) => item.code).filter(Boolean)
  const normalized = normalizeLocaleCode(code, supportedCodes, supportedCodes[0] || locale.value)
  if (locale.value === normalized) return
  await setLocale(normalized)
}
</script>
