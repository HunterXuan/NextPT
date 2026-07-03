<template>
  <section class="overflow-hidden rounded-lg border border-slate-200 bg-white dark:border-slate-800 dark:bg-slate-900">
    <div class="border-b border-slate-200 px-5 py-4 dark:border-slate-800">
      <div class="flex flex-col gap-3 lg:flex-row lg:items-start lg:justify-between">
        <div class="min-w-0">
          <h2 class="text-base font-semibold text-slate-950 dark:text-white">{{ $t('user.snatches.title') }}</h2>
          <p class="mt-1 text-sm text-slate-500 dark:text-slate-400">
            {{ $t('user.snatches.summary', { total: numberFormatter.format(snatchTotal) }) }}
          </p>
        </div>

        <div class="flex rounded-md border border-slate-200 bg-slate-50 p-1 dark:border-slate-800 dark:bg-slate-950">
          <button
            v-for="filter in snatchFilters"
            :key="filter.value"
            type="button"
            class="inline-flex h-8 items-center gap-1.5 rounded px-3 text-sm font-medium transition disabled:cursor-not-allowed disabled:opacity-60"
            :class="snatchStatus === filter.value
              ? 'bg-white text-slate-950 shadow-sm dark:bg-slate-800 dark:text-white'
              : 'text-slate-600 hover:bg-white/70 dark:text-slate-300 dark:hover:bg-slate-900'"
            :disabled="snatchesPending"
            @click="setSnatchStatus(filter.value)"
          >
            <UIcon :name="filter.icon" class="size-4" />
            {{ filter.label }}
          </button>
        </div>
      </div>
    </div>

    <div v-if="snatchesError" class="m-5 flex flex-col items-center justify-center rounded-md border border-red-200 bg-red-50 px-4 py-8 text-center dark:border-red-900 dark:bg-red-950">
      <UIcon name="i-lucide-circle-alert" class="size-8 text-red-500" />
      <p class="mt-3 text-sm font-medium text-red-700 dark:text-red-200">{{ snatchesError }}</p>
      <UButton class="mt-5" color="neutral" variant="outline" size="sm" icon="i-lucide-refresh-cw" @click="loadSnatches">
        {{ $t('common.retry') }}
      </UButton>
    </div>
    <div v-else-if="snatchesPending && snatches.length === 0" class="space-y-3 p-5">
      <div v-for="item in 6" :key="item" class="h-20 animate-pulse rounded-md bg-slate-100 dark:bg-slate-800" />
    </div>
    <div v-else-if="snatches.length === 0" class="m-5 rounded-md border border-dashed border-slate-200 px-4 py-12 text-center dark:border-slate-800">
      <UIcon name="i-lucide-inbox" class="mx-auto size-8 text-slate-400" />
      <p class="mt-3 text-sm text-slate-500 dark:text-slate-400">{{ $t('user.snatches.empty') }}</p>
    </div>
    <div v-else>
      <div class="hidden grid-cols-[minmax(0,1fr)_88px_112px_112px_148px_140px] gap-4 border-b border-slate-200 bg-slate-50 px-5 py-2.5 text-xs font-medium text-slate-500 lg:grid dark:border-slate-800 dark:bg-slate-950/60 dark:text-slate-400">
        <span>{{ $t('user.snatches.torrent') }}</span>
        <span class="text-center">{{ $t('user.snatches.status') }}</span>
        <span class="text-right">{{ $t('user.snatches.uploaded') }}</span>
        <span class="text-right">{{ $t('user.snatches.downloaded') }}</span>
        <span class="text-right">{{ $t('user.snatches.duration') }}</span>
        <span class="text-right">{{ $t('user.snatches.lastAction') }}</span>
      </div>

      <div class="divide-y divide-slate-100 dark:divide-slate-800">
        <article
          v-for="snatch in snatches"
          :key="snatch.id"
          class="grid gap-3 px-5 py-4 transition hover:bg-slate-50/70 lg:grid-cols-[minmax(0,1fr)_88px_112px_112px_148px_140px] lg:items-center lg:gap-4 dark:hover:bg-slate-950/50"
        >
          <div class="min-w-0">
            <NuxtLink
              :to="localePath(`/catalog/torrents/${snatch.torrentId}`)"
              class="block truncate text-sm font-semibold text-slate-950 hover:text-sky-700 dark:text-white dark:hover:text-sky-300"
            >
              {{ snatch.torrentName || $t('catalog.torrents.detail.titleFallback', { id: snatch.torrentId }) }}
            </NuxtLink>
            <div class="mt-1.5 text-xs text-slate-500 dark:text-slate-400">
              <span class="tabular-nums">{{ formatBytes(snatch.torrentSize) }}</span>
            </div>
          </div>

          <div class="grid grid-cols-2 gap-x-4 gap-y-3 text-xs sm:grid-cols-3 lg:contents">
            <div class="lg:text-center">
              <p class="text-slate-500 lg:hidden dark:text-slate-400">{{ $t('user.snatches.status') }}</p>
              <UBadge :color="snatch.isFinished ? 'success' : 'primary'" variant="soft" size="sm">
                {{ snatch.isFinished ? $t('user.snatches.finished') : $t('user.snatches.unfinished') }}
              </UBadge>
            </div>
            <div class="lg:text-right">
              <p class="text-slate-500 lg:hidden dark:text-slate-400">{{ $t('user.snatches.uploaded') }}</p>
              <p class="mt-1 font-medium tabular-nums text-slate-950 lg:mt-0 dark:text-white">{{ formatBytes(snatch.uploaded) }}</p>
            </div>
            <div class="lg:text-right">
              <p class="text-slate-500 lg:hidden dark:text-slate-400">{{ $t('user.snatches.downloaded') }}</p>
              <p class="mt-1 font-medium tabular-nums text-slate-950 lg:mt-0 dark:text-white">{{ formatBytes(snatch.downloaded) }}</p>
            </div>
            <div class="lg:text-right">
              <p class="text-slate-500 lg:hidden dark:text-slate-400">{{ $t('user.snatches.duration') }}</p>
              <p class="mt-1 font-medium tabular-nums text-slate-950 lg:mt-0 dark:text-white">{{ $t('user.snatches.seed') }} {{ formatDuration(snatch.seedTime) }}</p>
              <p class="mt-0.5 tabular-nums text-slate-500 dark:text-slate-400">{{ $t('user.snatches.leech') }} {{ formatDuration(snatch.leechTime) }}</p>
            </div>
            <div class="lg:text-right">
              <p class="text-slate-500 lg:hidden dark:text-slate-400">{{ $t('user.snatches.lastAction') }}</p>
              <p class="mt-1 font-medium text-slate-950 lg:mt-0 dark:text-white">{{ formatDateTime(snatch.lastActionAt, locale) }}</p>
            </div>
          </div>
        </article>
      </div>
    </div>

    <AppPager
      v-if="snatchTotal > 0"
      class="border-t border-slate-200 px-5 py-4 dark:border-slate-800"
      size="sm"
      :page="snatchPage"
      :total="snatchTotal"
      :page-size="snatchSize"
      :disabled="snatchesPending"
      @page-change="goToSnatchPage"
    />
  </section>
</template>

<script setup lang="ts">
import { useAccounting, type SnatchItem } from '~/composables/useAccounting'
import { ApiError } from '~/composables/useApi'
import { formatBytes, formatDateTime } from '~/utils/format'

type SnatchStatus = 'all' | 'finished' | 'unfinished'

const { t, locale } = useI18n()
const localePath = useLocalePath()
const accounting = useAccounting()

const snatches = ref<SnatchItem[]>([])
const snatchTotal = ref(0)
const snatchPage = ref(1)
const snatchSize = 20
const snatchStatus = ref<SnatchStatus>('all')
const snatchesPending = ref(true)
const snatchesError = ref('')

const numberFormatter = computed(() => new Intl.NumberFormat(locale.value))
const snatchTotalPages = computed(() => Math.max(1, Math.ceil(snatchTotal.value / snatchSize)))

const snatchFilters = computed<Array<{ value: SnatchStatus, label: string, icon: string }>>(() => [
  {
    value: 'all',
    label: t('user.snatches.filters.all'),
    icon: 'i-lucide-list'
  },
  {
    value: 'finished',
    label: t('user.snatches.filters.finished'),
    icon: 'i-lucide-circle-check'
  },
  {
    value: 'unfinished',
    label: t('user.snatches.filters.unfinished'),
    icon: 'i-lucide-clock'
  }
])

onMounted(loadSnatches)

async function loadSnatches() {
  if (snatchesPending.value && snatches.value.length > 0) return

  snatchesPending.value = true
  snatchesError.value = ''

  try {
    let data = await accounting.listSnatches({
      page: snatchPage.value,
      size: snatchSize,
      isFinished: snatchStatus.value === 'all' ? undefined : snatchStatus.value === 'finished'
    })

    const nextTotal = data.total || 0
    const nextTotalPages = Math.max(1, Math.ceil(nextTotal / snatchSize))
    if (snatchPage.value > nextTotalPages) {
      snatchPage.value = nextTotalPages
      data = await accounting.listSnatches({
        page: snatchPage.value,
        size: snatchSize,
        isFinished: snatchStatus.value === 'all' ? undefined : snatchStatus.value === 'finished'
      })
    }

    snatches.value = data.list || []
    snatchTotal.value = data.total || 0
  } catch (error) {
    snatches.value = []
    snatchTotal.value = 0
    snatchesError.value = error instanceof ApiError ? error.message : t('common.requestFailed')
  } finally {
    snatchesPending.value = false
  }
}

function setSnatchStatus(status: SnatchStatus) {
  if (snatchStatus.value === status) return

  snatchStatus.value = status
  snatchPage.value = 1
  loadSnatches()
}

function goToSnatchPage(page: number) {
  snatchPage.value = Math.min(Math.max(1, page), snatchTotalPages.value)
  loadSnatches()
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
</script>
