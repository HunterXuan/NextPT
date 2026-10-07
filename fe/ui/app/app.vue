<template>
  <UApp :locale="uiLocale">
    <NuxtLoadingIndicator color="var(--ui-primary)" />
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
import { siteThemeStyles, siteThemeBootstrap } from '~/utils/siteTheme'

const { defaultTheme, restoreTheme } = useSiteTheme()
useHead({
  style: [{ key: 'site-themes', innerHTML: siteThemeStyles }],
  script: [{ key: 'site-theme-init', innerHTML: siteThemeBootstrap(defaultTheme.value), tagPosition: 'head' }]
})
onMounted(restoreTheme)

const { locale } = useI18n()
const { isLoggedIn } = useAuth()
const { pageKey } = useWorkspacePageCache()
const { siteName, siteFavicon } = useSiteBranding()

const uiLocale = computed(() => {
  if (locale.value === 'zh-CN') return locales.zh_cn
  if (locale.value === 'zh-TW') return locales.zh_tw
  return locales.en
})

const htmlLang = computed(() => locale.value)

useHead({
  titleTemplate: (title) => (title ? `${title} | ${siteName.value}` : siteName.value),
  htmlAttrs: {
    lang: htmlLang
  },
  link: [{ key: 'icon', rel: 'icon', href: siteFavicon }]
})

useSeoMeta({
  ogType: 'website',
  ogSiteName: () => siteName.value,
  twitterCard: 'summary'
})
</script>
