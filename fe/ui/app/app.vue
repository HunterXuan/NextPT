<template>
  <UApp :locale="uiLocale">
    <NuxtLoadingIndicator color="#0ea5e9" />
    <NuxtLayout>
      <NuxtPage />
    </NuxtLayout>
  </UApp>
</template>

<script setup lang="ts">
import * as locales from '@nuxt/ui/locale'

const { locale } = useI18n()

const uiLocale = computed(() => {
  if (locale.value === 'zh-CN') return locales.zh_cn
  if (locale.value === 'zh-TW') return locales.zh_tw
  return locales.en
})

const htmlLang = computed(() => locale.value)

useHead({
  titleTemplate: (title) => (title ? `${title} | NextPT` : 'NextPT'),
  htmlAttrs: {
    lang: htmlLang
  }
})

useSeoMeta({
  ogType: 'website',
  ogSiteName: 'NextPT',
  twitterCard: 'summary'
})
</script>
