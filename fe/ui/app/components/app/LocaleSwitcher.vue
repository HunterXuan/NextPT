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
const { locale, locales, setLocale } = useI18n()

const languageItems = computed(() => [
  (locales.value as any[]).map((item) => ({
    label: item.name,
    onSelect: () => handleLanguageSwitch(item.code)
  }))
])

async function handleLanguageSwitch(code: string) {
  if (locale.value === code) return
  await setLocale(code)
}
</script>
