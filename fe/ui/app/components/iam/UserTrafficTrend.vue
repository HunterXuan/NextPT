<template>
  <section class="overflow-hidden rounded-lg border border-slate-200 bg-white dark:border-slate-800 dark:bg-slate-900">
    <div class="border-b border-slate-200 px-5 py-4 dark:border-slate-800">
      <div class="flex flex-col gap-3 lg:flex-row lg:items-start lg:justify-between">
        <div class="min-w-0">
          <h2 class="text-base font-semibold text-slate-950 dark:text-white">{{ $t('user.trafficHistory.title') }}</h2>
          <p class="mt-1 text-sm text-slate-500 dark:text-slate-400">{{ $t('user.trafficHistory.subtitle') }}</p>
        </div>

        <div class="flex rounded-md border border-slate-200 bg-slate-50 p-1 dark:border-slate-800 dark:bg-slate-950">
          <button
            v-for="period in trafficHistoryPeriods"
            :key="period.value"
            type="button"
            class="inline-flex h-8 items-center gap-1.5 rounded px-3 text-sm font-medium transition disabled:cursor-not-allowed disabled:opacity-60"
            :class="trafficHistoryPeriod === period.value
              ? 'bg-white text-slate-950 shadow-sm dark:bg-slate-800 dark:text-white'
              : 'text-slate-600 hover:bg-white/70 dark:text-slate-300 dark:hover:bg-slate-900'"
            :disabled="trafficHistoryPending"
            @click="setTrafficHistoryPeriod(period.value)"
          >
            <UIcon :name="period.icon" class="size-4" />
            {{ period.label }}
          </button>
        </div>
      </div>
    </div>

    <div class="px-5 py-5">
      <p class="mb-3 text-sm font-medium text-slate-700 dark:text-slate-200">{{ trafficHistorySummaryTitle }}</p>
      <div class="grid grid-cols-2 gap-px overflow-hidden rounded-lg border border-slate-200 bg-slate-200 sm:grid-cols-4 dark:border-slate-800 dark:bg-slate-800">
        <div
          v-for="item in trafficHistorySummaryCards"
          :key="item.label"
          class="min-w-0 bg-slate-50/80 p-3 dark:bg-slate-950/70"
        >
          <div class="flex items-center gap-2">
            <span class="flex size-7 shrink-0 items-center justify-center rounded-md" :class="item.iconClass">
              <UIcon :name="item.icon" class="size-4" />
            </span>
            <p class="truncate text-xs text-slate-500 dark:text-slate-400">{{ item.label }}</p>
          </div>
          <p class="mt-2 truncate text-base font-semibold tabular-nums text-slate-950 dark:text-white">{{ item.value }}</p>
        </div>
      </div>

      <div class="mt-5 flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
        <p class="text-sm font-medium text-slate-700 dark:text-slate-200">{{ trafficHistoryChartTitle }}</p>
        <div class="flex w-fit rounded-md border border-slate-200 bg-slate-50 p-1 dark:border-slate-800 dark:bg-slate-950">
          <button
            v-for="mode in trafficHistoryModes"
            :key="mode.value"
            type="button"
            class="inline-flex h-8 items-center gap-1.5 rounded px-3 text-sm font-medium transition"
            :class="trafficChartMode === mode.value
              ? 'bg-white text-slate-950 shadow-sm dark:bg-slate-800 dark:text-white'
              : 'text-slate-600 hover:bg-white/70 dark:text-slate-300 dark:hover:bg-slate-900'"
            @click="trafficChartMode = mode.value"
          >
            <UIcon :name="mode.icon" class="size-4" />
            {{ mode.label }}
          </button>
        </div>
      </div>

      <div v-if="trafficHistoryError" class="mt-4 flex flex-col items-center justify-center rounded-md border border-red-200 bg-red-50 px-4 py-8 text-center dark:border-red-900 dark:bg-red-950">
        <UIcon name="i-lucide-circle-alert" class="size-8 text-red-500" />
        <p class="mt-3 text-sm font-medium text-red-700 dark:text-red-200">{{ trafficHistoryError }}</p>
        <UButton class="mt-5" color="neutral" variant="outline" size="sm" icon="i-lucide-refresh-cw" @click="loadTrafficHistory">
          {{ $t('common.retry') }}
        </UButton>
      </div>
      <div v-else-if="trafficHistoryPending && trafficHistory.length === 0" class="mt-4 h-80 animate-pulse rounded-md bg-slate-100 dark:bg-slate-800" />
      <div v-else-if="trafficChartItems.length === 0" class="mt-4 rounded-md border border-dashed border-slate-200 px-4 py-10 text-center dark:border-slate-800">
        <UIcon name="i-lucide-chart-no-axes-column" class="mx-auto size-8 text-slate-400" />
        <p class="mt-3 text-sm text-slate-500 dark:text-slate-400">{{ $t('user.trafficHistory.empty') }}</p>
      </div>
      <div v-else class="mt-4 h-80 rounded-md border border-slate-200 bg-white px-3 py-4 dark:border-slate-800 dark:bg-slate-950">
        <ClientOnly>
          <Line :key="trafficChartKey" :data="trafficChartData" :options="trafficChartOptions" />
          <template #fallback>
            <div class="h-full animate-pulse rounded-md bg-slate-100 dark:bg-slate-800" />
          </template>
        </ClientOnly>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
import {
  CategoryScale,
  Chart as ChartJS,
  Legend,
  LineElement,
  LinearScale,
  PointElement,
  Tooltip,
  type ChartData,
  type ChartOptions,
  type Plugin,
  type TooltipItem
} from 'chart.js'
import { Line } from 'vue-chartjs'

import { useAccounting, type TrafficHistoryItem, type TrafficHistoryPeriod } from '~/composables/useAccounting'
import { ApiError } from '~/composables/useApi'
import { formatBytes } from '~/utils/format'

type TrafficChartMode = 'traffic' | 'time' | 'bonus'
type TrafficChartSeriesKey = keyof Pick<TrafficHistoryItem, 'uploaded' | 'downloaded' | 'seedTime' | 'leechTime' | 'bonus'>

interface TrafficChartDatasetDefinition {
  key: TrafficChartSeriesKey
  label: string
  color: string
  backgroundColor: string
}

const trafficLegendBottomGap = 14

const trafficLegendGapPlugin: Plugin<'line'> = {
  id: 'trafficLegendGap',
  beforeInit(chart) {
    const legend = chart.legend
    if (!legend?.fit) return

    const originalFit = legend.fit.bind(legend)
    legend.fit = () => {
      originalFit()
      if (legend.options.display !== false && legend.position === 'top') {
        legend.height += trafficLegendBottomGap
      }
    }
  }
}

ChartJS.register(CategoryScale, LinearScale, PointElement, LineElement, Tooltip, Legend, trafficLegendGapPlugin)

const { t, locale } = useI18n()
const colorMode = useColorMode()
const accounting = useAccounting()

const trafficHistory = ref<TrafficHistoryItem[]>([])
const trafficHistoryPeriod = ref<TrafficHistoryPeriod>('daily')
const trafficChartMode = ref<TrafficChartMode>('traffic')
const trafficHistoryPending = ref(true)
const trafficHistoryError = ref('')
const numberFormatter = computed(() => new Intl.NumberFormat(locale.value))
const isDark = computed(() => colorMode.value === 'dark')

const chartTextColor = computed(() => isDark.value ? '#cbd5e1' : '#64748b')
const chartMutedColor = computed(() => isDark.value ? '#64748b' : '#94a3b8')
const chartGridColor = computed(() => isDark.value ? 'rgba(51, 65, 85, 0.75)' : 'rgba(226, 232, 240, 0.9)')
const chartSurfaceColor = computed(() => isDark.value ? '#020617' : '#ffffff')

const trafficHistoryPeriods = computed<Array<{ value: TrafficHistoryPeriod, label: string, icon: string }>>(() => [
  {
    value: 'daily',
    label: t('user.trafficHistory.periods.daily'),
    icon: 'i-lucide-calendar-days'
  },
  {
    value: 'monthly',
    label: t('user.trafficHistory.periods.monthly'),
    icon: 'i-lucide-calendar-range'
  }
])

const trafficHistoryModes = computed<Array<{ value: TrafficChartMode, label: string, icon: string }>>(() => [
  {
    value: 'traffic',
    label: t('user.trafficHistory.modes.traffic'),
    icon: 'i-lucide-arrow-up-down'
  },
  {
    value: 'time',
    label: t('user.trafficHistory.modes.time'),
    icon: 'i-lucide-timer'
  },
  {
    value: 'bonus',
    label: t('user.trafficHistory.modes.bonus'),
    icon: 'i-lucide-coins'
  }
])

const trafficChartDatasetDefinitions = computed<TrafficChartDatasetDefinition[]>(() => {
  if (trafficChartMode.value === 'time') {
    return [
      { key: 'seedTime', label: t('user.trafficHistory.seedTime'), color: '#10b981', backgroundColor: 'rgba(16, 185, 129, 0.12)' },
      { key: 'leechTime', label: t('user.trafficHistory.leechTime'), color: '#8b5cf6', backgroundColor: 'rgba(139, 92, 246, 0.12)' }
    ]
  }

  if (trafficChartMode.value === 'bonus') {
    return [
      { key: 'bonus', label: t('user.trafficHistory.bonus'), color: '#f59e0b', backgroundColor: 'rgba(245, 158, 11, 0.14)' }
    ]
  }

  return [
    { key: 'uploaded', label: t('user.trafficHistory.uploaded'), color: '#10b981', backgroundColor: 'rgba(16, 185, 129, 0.12)' },
    { key: 'downloaded', label: t('user.trafficHistory.downloaded'), color: '#0ea5e9', backgroundColor: 'rgba(14, 165, 233, 0.12)' }
  ]
})

const trafficChartItems = computed(() => {
  const limit = trafficHistoryPeriod.value === 'daily' ? 30 : 12
  return [...trafficHistory.value]
    .sort((a, b) => a.date.localeCompare(b.date))
    .slice(-limit)
})

const trafficHistoryTotals = computed(() => {
  return trafficChartItems.value.reduce(
    (total, item) => {
      total.uploaded += Number(item.uploaded || 0)
      total.downloaded += Number(item.downloaded || 0)
      total.activeTime += Number(item.seedTime || 0) + Number(item.leechTime || 0)
      total.bonus += Number(item.bonus || 0)
      return total
    },
    {
      uploaded: 0,
      downloaded: 0,
      activeTime: 0,
      bonus: 0
    }
  )
})

const trafficHistorySummaryCards = computed(() => [
  {
    label: t('user.trafficHistory.uploaded'),
    value: formatBytes(trafficHistoryTotals.value.uploaded),
    icon: 'i-lucide-upload',
    iconClass: 'bg-emerald-100 text-emerald-700 dark:bg-emerald-950 dark:text-emerald-300'
  },
  {
    label: t('user.trafficHistory.downloaded'),
    value: formatBytes(trafficHistoryTotals.value.downloaded),
    icon: 'i-lucide-download',
    iconClass: 'bg-sky-100 text-sky-700 dark:bg-sky-950 dark:text-sky-300'
  },
  {
    label: t('user.trafficHistory.activeTime'),
    value: formatDuration(trafficHistoryTotals.value.activeTime),
    icon: 'i-lucide-timer',
    iconClass: 'bg-violet-100 text-violet-700 dark:bg-violet-950 dark:text-violet-300'
  },
  {
    label: t('user.trafficHistory.bonus'),
    value: formatBonus(trafficHistoryTotals.value.bonus),
    icon: 'i-lucide-coins',
    iconClass: 'bg-amber-100 text-amber-700 dark:bg-amber-950 dark:text-amber-300'
  }
])

const trafficHistorySummaryTitle = computed(() => t('user.trafficHistory.summaryTitle', { range: trafficHistorySummaryRangeLabel.value }))

const trafficHistoryChartTitle = computed(() => t('user.trafficHistory.chartTitle', { range: trafficHistorySummaryRangeLabel.value }))

const trafficHistorySummaryRangeLabel = computed(() => {
  if (trafficHistoryPeriod.value === 'daily') {
    return t('user.trafficHistory.ranges.daily', { count: 30 })
  }

  return t('user.trafficHistory.ranges.monthly', { count: 12 })
})

const trafficChartSuggestedMax = computed(() => {
  const values = trafficChartItems.value.flatMap((item) => trafficChartDatasetDefinitions.value.map((series) => Number(item[series.key] || 0)))
  const maxValue = Math.max(0, ...values)
  if (trafficChartMode.value === 'traffic') return Math.max(1024, maxValue)
  if (trafficChartMode.value === 'time') return Math.max(60, maxValue)
  return Math.max(1, maxValue)
})

const trafficChartData = computed<ChartData<'line'>>(() => ({
  labels: trafficChartItems.value.map((item) => item.date),
  datasets: trafficChartDatasetDefinitions.value.map((series) => ({
    label: series.label,
    data: trafficChartItems.value.map((item) => Number(item[series.key] || 0)),
    borderColor: series.color,
    backgroundColor: series.backgroundColor,
    pointBackgroundColor: series.color,
    pointBorderColor: chartSurfaceColor.value,
    pointHoverBackgroundColor: series.color,
    pointHoverBorderColor: chartSurfaceColor.value,
    borderWidth: 2,
    pointBorderWidth: 2,
    pointHoverBorderWidth: 3,
    pointRadius: 3,
    pointHoverRadius: 5,
    tension: 0.35,
    fill: false
  }))
}))

const trafficChartOptions = computed<ChartOptions<'line'>>(() => ({
  responsive: true,
  maintainAspectRatio: false,
  animation: {
    duration: 180
  },
  interaction: {
    mode: 'index',
    intersect: false
  },
  plugins: {
    legend: {
      display: true,
      position: 'top',
      align: 'start',
      labels: {
        boxHeight: 8,
        boxWidth: 18,
        color: chartTextColor.value,
        font: {
          size: 12
        },
        padding: 16,
        usePointStyle: true
      }
    },
    tooltip: {
      mode: 'index',
      intersect: false,
      backgroundColor: isDark.value ? 'rgba(15, 23, 42, 0.96)' : 'rgba(255, 255, 255, 0.96)',
      bodyColor: isDark.value ? '#e2e8f0' : '#334155',
      borderColor: isDark.value ? '#334155' : '#e2e8f0',
      borderWidth: 1,
      padding: 10,
      titleColor: isDark.value ? '#f8fafc' : '#0f172a',
      callbacks: {
        label: (item: TooltipItem<'line'>) => `${item.dataset.label || ''}: ${formatTrafficChartValue(Number(item.parsed.y || 0))}`
      }
    }
  },
  scales: {
    x: {
      border: {
        display: false
      },
      grid: {
        display: false
      },
      ticks: {
        autoSkip: true,
        color: chartMutedColor.value,
        maxRotation: 0,
        maxTicksLimit: trafficHistoryPeriod.value === 'daily' ? 6 : 12
      }
    },
    y: {
      beginAtZero: true,
      border: {
        display: false
      },
      grid: {
        color: chartGridColor.value
      },
      suggestedMax: trafficChartSuggestedMax.value,
      ticks: {
        color: chartMutedColor.value,
        callback: (value) => formatTrafficChartValue(Number(value))
      }
    }
  }
}))

const trafficChartKey = computed(() => [
  trafficChartMode.value,
  trafficHistoryPeriod.value,
  isDark.value ? 'dark' : 'light',
  locale.value
].join(':'))

onMounted(loadTrafficHistory)

async function loadTrafficHistory() {
  if (trafficHistoryPending.value && trafficHistory.value.length > 0) return

  trafficHistoryPending.value = true
  trafficHistoryError.value = ''

  try {
    const data = await accounting.listTrafficHistory({
      period: trafficHistoryPeriod.value
    })
    trafficHistory.value = data.list || []
  } catch (error) {
    trafficHistory.value = []
    trafficHistoryError.value = error instanceof ApiError ? error.message : t('common.requestFailed')
  } finally {
    trafficHistoryPending.value = false
  }
}

function setTrafficHistoryPeriod(period: TrafficHistoryPeriod) {
  if (trafficHistoryPeriod.value === period) return

  trafficHistoryPeriod.value = period
  loadTrafficHistory()
}

function formatDuration(value?: number | null) {
  const seconds = Math.max(0, Number(value || 0))
  if (seconds < 60) {
    return t('user.duration.seconds', { count: numberFormatter.value.format(Math.floor(seconds)) })
  }

  const minutes = seconds / 60
  if (minutes < 60) {
    return t('user.duration.minutes', { count: numberFormatter.value.format(Math.floor(minutes)) })
  }

  const hours = minutes / 60
  if (hours < 24) {
    return t('user.duration.hours', { count: numberFormatter.value.format(Number(hours.toFixed(hours >= 10 ? 0 : 1))) })
  }

  const days = hours / 24
  return t('user.duration.days', { count: numberFormatter.value.format(Number(days.toFixed(days >= 10 ? 0 : 1))) })
}

function formatBonus(value?: number | null) {
  return numberFormatter.value.format(Number(value || 0))
}

function formatTrafficChartValue(value: number) {
  if (trafficChartMode.value === 'time') return formatDuration(value)
  if (trafficChartMode.value === 'bonus') return formatBonus(value)
  return formatBytes(value)
}
</script>
