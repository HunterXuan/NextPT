<template>
  <section class="grid gap-3 md:grid-cols-3">
    <NuxtLink
      v-for="item in items"
      :key="item.to"
      class="group min-h-32 rounded-lg border border-slate-200 bg-white p-4 transition hover:border-sky-200 hover:bg-sky-50/50 dark:border-slate-800 dark:bg-slate-900 dark:hover:border-sky-900 dark:hover:bg-sky-950/30"
      :to="localePath(item.to)"
    >
      <div class="flex h-full flex-col justify-between gap-4">
        <div class="flex items-start gap-3">
          <span class="flex size-9 shrink-0 items-center justify-center rounded-md bg-slate-100 text-slate-600 transition group-hover:bg-sky-100 group-hover:text-sky-700 dark:bg-slate-800 dark:text-slate-300 dark:group-hover:bg-sky-950 dark:group-hover:text-sky-200">
            <UIcon :name="item.icon" class="size-4" />
          </span>
          <div class="min-w-0">
            <p class="font-medium text-slate-950 dark:text-white">{{ item.label }}</p>
            <p class="mt-1 line-clamp-2 text-sm leading-6 text-slate-500 dark:text-slate-400">{{ item.description }}</p>
          </div>
        </div>
        <div class="flex items-center justify-between border-t border-slate-100 pt-3 text-sm font-medium text-slate-500 dark:border-slate-800 dark:text-slate-400">
          <span>{{ $t('user.overview.enter') }}</span>
          <span class="flex size-9 shrink-0 items-center justify-center rounded-md bg-slate-100 text-slate-600 transition group-hover:bg-sky-100 group-hover:text-sky-700 dark:bg-slate-800 dark:text-slate-300 dark:group-hover:bg-sky-950 dark:group-hover:text-sky-200">
            <UIcon name="i-lucide-arrow-right" class="size-4" />
          </span>
        </div>
      </div>
    </NuxtLink>
  </section>
</template>

<script setup lang="ts">
const { t } = useI18n()
const localePath = useLocalePath()

const items = computed(() => [
  {
    label: t('user.nav.activity'),
    description: t('user.overview.links.activity'),
    icon: 'i-lucide-chart-line',
    to: '/iam/users/me/activity'
  },
  {
    label: t('user.bonusLogs.title'),
    description: t('user.overview.links.bonus'),
    icon: 'i-lucide-coins',
    to: '/iam/users/me/bonus'
  },
  {
    label: t('user.nav.settings'),
    description: t('user.overview.links.settings'),
    icon: 'i-lucide-settings',
    to: '/iam/users/me/settings'
  }
])
</script>
