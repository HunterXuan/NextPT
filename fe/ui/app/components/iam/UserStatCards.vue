<template>
  <section class="overflow-hidden rounded-lg border border-slate-200 bg-white dark:border-slate-800 dark:bg-slate-900">
    <div class="grid grid-cols-1 divide-y divide-slate-100 sm:grid-cols-2 sm:divide-x sm:divide-y-0 xl:grid-cols-4 dark:divide-slate-800">
      <div v-for="item in statCards" :key="item.label" class="flex items-center justify-between gap-3 p-4">
        <div class="min-w-0 space-y-1">
          <div class="flex items-center gap-2">
            <span class="size-2 rounded-full" :class="item.dotClass" />
            <p class="text-sm text-slate-500 dark:text-slate-400">{{ item.label }}</p>
          </div>
          <p class="truncate text-2xl font-semibold tabular-nums text-slate-950 dark:text-white">{{ item.value }}</p>
        </div>
        <span class="flex size-10 shrink-0 items-center justify-center rounded-md" :class="item.iconClass">
          <UIcon :name="item.icon" class="size-5" />
        </span>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
import type { TrafficSummary } from '~/composables/useAccounting'
import type { AuthUser } from '~/composables/useAuth'
import { formatBytes } from '~/utils/format'

const props = defineProps<{
  user?: AuthUser | null
  traffic?: TrafficSummary | null
}>()

const { t, locale } = useI18n()
const numberFormatter = computed(() => new Intl.NumberFormat(locale.value))

const trafficUploaded = computed(() => props.traffic?.uploaded ?? props.user?.uploaded ?? 0)
const trafficDownloaded = computed(() => props.traffic?.downloaded ?? props.user?.downloaded ?? 0)
const ratio = computed(() => props.traffic?.shareRatio ?? props.user?.shareRatio ?? 0)

const statCards = computed(() => [
  {
    label: t('user.stats.uploaded'),
    value: formatBytes(trafficUploaded.value),
    icon: 'i-lucide-upload',
    iconClass: 'bg-emerald-100 text-emerald-700 dark:bg-emerald-950 dark:text-emerald-300',
    dotClass: 'bg-emerald-500'
  },
  {
    label: t('user.stats.downloaded'),
    value: formatBytes(trafficDownloaded.value),
    icon: 'i-lucide-download',
    iconClass: 'bg-sky-100 text-sky-700 dark:bg-sky-950 dark:text-sky-300',
    dotClass: 'bg-sky-500'
  },
  {
    label: t('user.stats.ratio'),
    value: Number.isFinite(ratio.value) ? ratio.value.toFixed(2) : '-',
    icon: 'i-lucide-scale',
    iconClass: 'bg-amber-100 text-amber-700 dark:bg-amber-950 dark:text-amber-300',
    dotClass: 'bg-amber-500'
  },
  {
    label: t('user.stats.bonus'),
    value: numberFormatter.value.format(Number(props.user?.bonus || 0)),
    icon: 'i-lucide-coins',
    iconClass: 'bg-violet-100 text-violet-700 dark:bg-violet-950 dark:text-violet-300',
    dotClass: 'bg-violet-500'
  }
])
</script>
