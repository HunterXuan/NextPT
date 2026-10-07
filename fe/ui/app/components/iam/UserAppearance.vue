<script setup lang="ts">
import { siteThemes } from '~/utils/siteTheme'

const { t } = useI18n()
const { preference, defaultTheme, setTheme } = useSiteTheme()
const colorMode = useColorMode()
const modes = [
  { value: 'light', icon: 'i-lucide-sun' },
  { value: 'dark', icon: 'i-lucide-moon' },
  { value: 'system', icon: 'i-lucide-monitor' }
]
</script>

<template>
  <section class="bg-default p-5 sm:p-6">
    <h2 class="text-base font-semibold text-highlighted">{{ t('appearance.title') }}</h2>
    <div class="mt-6 flex flex-wrap items-center justify-between gap-3">
      <h3 class="text-sm font-medium text-default">{{ t('appearance.palette') }}</h3>
      <UButton
        icon="i-lucide-rotate-ccw"
        color="neutral"
        variant="ghost"
        size="sm"
        :disabled="preference === 'default'"
        @click="setTheme('default')"
      >{{ t('appearance.siteDefault', { name: t(`appearance.themes.${defaultTheme}`) }) }}</UButton>
    </div>
    <div class="mt-3 grid grid-cols-2 gap-3 lg:grid-cols-4" role="group" :aria-label="t('appearance.palette')">
      <button
        v-for="theme in siteThemes"
        :key="theme.id"
        type="button"
        class="relative flex min-h-16 items-center gap-3 rounded-md border px-3 py-3 text-left transition-colors focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary"
        :class="(preference === 'default' ? defaultTheme : preference) === theme.id
          ? 'border-primary bg-primary/5'
          : 'border-default hover:bg-elevated'"
        :aria-pressed="(preference === 'default' ? defaultTheme : preference) === theme.id"
        @click="setTheme(theme.id)"
      >
        <span class="size-7 shrink-0 rounded-full ring-1 ring-black/10 dark:ring-white/10" :style="{ backgroundColor: theme.swatch }" />
        <span class="text-sm font-medium text-highlighted">{{ t(`appearance.themes.${theme.id}`) }}</span>
        <UIcon v-if="(preference === 'default' ? defaultTheme : preference) === theme.id" name="i-lucide-check" class="ml-auto size-4 shrink-0 text-primary" />
      </button>
    </div>
    <div class="mt-6 border-t border-default pt-5">
      <h3 class="mb-3 text-sm font-medium text-default">{{ t('appearance.mode') }}</h3>
      <ClientOnly>
        <div class="flex w-fit flex-wrap gap-1 rounded-md bg-elevated p-1" role="group" :aria-label="t('appearance.mode')">
          <button
            v-for="mode in modes"
            :key="mode.value"
            type="button"
            class="flex items-center gap-2 rounded px-3 py-2 text-sm transition-colors focus-visible:outline-2 focus-visible:outline-primary"
            :class="colorMode.preference === mode.value ? 'app-selected shadow-sm' : 'text-muted hover:text-highlighted'"
            :aria-pressed="colorMode.preference === mode.value"
            @click="colorMode.preference = mode.value"
          >
            <UIcon :name="mode.icon" class="size-4" />
            {{ t(`appearance.modes.${mode.value}`) }}
          </button>
        </div>
        <template #fallback><div class="h-11" /></template>
      </ClientOnly>
    </div>
  </section>
</template>
