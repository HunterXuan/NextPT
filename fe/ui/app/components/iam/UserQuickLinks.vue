<template>
  <section class="grid gap-3 md:grid-cols-3">
    <NuxtLink
      v-for="item in items"
      :key="item.to"
      class="group min-h-32 rounded-lg border border-neutral-200 bg-white p-4 transition hover:border-primary-200 hover:bg-primary-50/50 dark:border-neutral-800 dark:bg-neutral-900 dark:hover:border-primary-900 dark:hover:bg-primary-950/30"
      :to="localePath(item.to)"
    >
      <div class="flex h-full flex-col justify-between gap-4">
        <div class="flex items-start gap-3">
          <span class="flex size-9 shrink-0 items-center justify-center rounded-md bg-neutral-100 text-neutral-600 transition group-hover:bg-primary-100 group-hover:text-primary-700 dark:bg-neutral-800 dark:text-neutral-300 dark:group-hover:bg-primary-950 dark:group-hover:text-primary-200">
            <UIcon :name="item.icon" class="size-4" />
          </span>
          <div class="min-w-0">
            <p class="font-medium text-neutral-950 dark:text-white">{{ item.label }}</p>
            <p class="mt-1 line-clamp-2 text-sm leading-6 text-neutral-500 dark:text-neutral-400">{{ item.description }}</p>
          </div>
        </div>
        <div class="flex items-center justify-between border-t border-neutral-100 pt-3 text-sm font-medium text-neutral-500 dark:border-neutral-800 dark:text-neutral-400">
          <span>{{ $t('user.overview.enter') }}</span>
          <span class="flex size-9 shrink-0 items-center justify-center rounded-md bg-neutral-100 text-neutral-600 transition group-hover:bg-primary-100 group-hover:text-primary-700 dark:bg-neutral-800 dark:text-neutral-300 dark:group-hover:bg-primary-950 dark:group-hover:text-primary-200">
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
