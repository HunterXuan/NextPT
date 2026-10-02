<template>
  <UApp :locale="uiLocale">
    <NuxtLoadingIndicator color="#0ea5e9" />
    <NuxtLayout>
      <NuxtPage
        :key="isLoggedIn ? 'workspace' : 'public'"
        :keepalive="isLoggedIn ? { max: 32 } : false"
        :page-key="pageKey"
        :transition="false"
      />
    </NuxtLayout>
  </UApp>
</template>

<script setup lang="ts">
import * as locales from '@nuxt/ui/locale'

const { locale } = useI18n()
const { isLoggedIn } = useAuth()
const { pageKey } = useWorkspacePageCache()

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
