<template>
  <section class="overflow-hidden rounded-lg border border-slate-200 bg-white dark:border-slate-800 dark:bg-slate-900">
    <div class="border-b border-slate-200 px-5 py-4 dark:border-slate-800">
      <div class="flex flex-col gap-3 lg:flex-row lg:items-center lg:justify-between">
        <div class="min-w-0">
          <h2 class="text-base font-semibold text-slate-950 dark:text-white">{{ $t('user.torrents.title') }}</h2>
          <p class="mt-1 text-sm text-slate-500 dark:text-slate-400">
            {{ $t('user.torrents.summary', { total: numberFormatter.format(total) }) }}
          </p>
        </div>

        <div class="overflow-x-auto [scrollbar-width:none] [&::-webkit-scrollbar]:hidden">
          <div class="inline-flex min-w-max rounded-md bg-slate-100 p-1 dark:bg-slate-800/80">
            <button
              v-for="filter in filters"
              :key="filter.value"
              type="button"
              class="inline-flex h-8 items-center rounded px-3 text-sm font-medium transition disabled:cursor-not-allowed disabled:opacity-60"
              :class="status === filter.value
                ? 'bg-white text-slate-950 shadow-sm dark:bg-slate-950 dark:text-white'
                : 'text-slate-600 hover:text-slate-950 dark:text-slate-300 dark:hover:text-white'"
              :disabled="pending"
              @click="setStatus(filter.value)"
            >
              {{ filter.label }}
            </button>
          </div>
        </div>
      </div>
    </div>

    <div v-if="errorMessage" class="m-5 flex flex-col items-center justify-center rounded-md border border-red-200 bg-red-50 px-4 py-8 text-center dark:border-red-900 dark:bg-red-950">
      <UIcon name="i-lucide-circle-alert" class="size-8 text-red-500" />
      <p class="mt-3 text-sm font-medium text-red-700 dark:text-red-200">{{ errorMessage }}</p>
    </div>
    <div v-else-if="pending && items.length === 0" class="space-y-3 p-5">
      <div v-for="item in 4" :key="item" class="h-20 animate-pulse rounded-md bg-slate-100 dark:bg-slate-800" />
    </div>
    <div v-else-if="items.length === 0" class="m-5 rounded-md border border-dashed border-slate-200 px-4 py-12 text-center dark:border-slate-800">
      <UIcon name="i-lucide-upload-cloud" class="mx-auto size-8 text-slate-400" />
      <p class="mt-3 text-sm text-slate-500 dark:text-slate-400">{{ $t('user.torrents.empty') }}</p>
    </div>
    <div v-else>
      <div class="hidden grid-cols-[minmax(0,1fr)_112px_120px_168px_88px] gap-4 border-b border-slate-200 bg-slate-50 px-5 py-2.5 text-xs font-medium text-slate-500 lg:grid dark:border-slate-800 dark:bg-slate-950/60 dark:text-slate-400">
        <span>{{ $t('user.torrents.columns.torrent') }}</span>
        <span class="text-center">{{ $t('user.torrents.columns.status') }}</span>
        <span class="text-right">{{ $t('user.torrents.columns.size') }}</span>
        <span class="text-right">{{ $t('user.torrents.columns.time') }}</span>
        <span class="text-right">{{ $t('user.torrents.columns.actions') }}</span>
      </div>

      <div class="divide-y divide-slate-100 dark:divide-slate-800">
        <article
          v-for="item in items"
          :key="item.id"
          class="grid gap-3 px-5 py-4 outline-none transition lg:grid-cols-[minmax(0,1fr)_112px_120px_168px_88px] lg:items-center lg:gap-4"
          :class="item.banned ? 'cursor-default opacity-75' : 'cursor-pointer hover:bg-slate-50/70 focus-visible:bg-slate-50/70 dark:hover:bg-slate-950/50 dark:focus-visible:bg-slate-950/50'"
          :tabindex="item.banned ? -1 : 0"
          :role="item.banned ? undefined : 'link'"
          @click="openTorrent(item)"
          @keydown.enter.prevent="openTorrent(item)"
        >
          <div class="min-w-0">
            <p class="truncate text-sm font-semibold text-slate-950 dark:text-white">{{ item.name }}</p>
            <p v-if="item.subTitle" class="mt-1 truncate text-xs text-slate-500 dark:text-slate-400">{{ item.subTitle }}</p>
            <p v-if="item.status === TorrentStatus.Rejected && item.reviewComment" class="mt-1.5 truncate text-xs text-red-600 dark:text-red-300">
              {{ item.reviewComment }}
            </p>
            <p class="mt-1.5 text-xs text-slate-500 lg:hidden dark:text-slate-400">
              {{ $t('user.torrents.files', { count: numberFormatter.format(item.fileCount) }) }}
            </p>
          </div>

          <div class="flex items-center justify-between gap-3 lg:justify-center">
            <span class="text-xs text-slate-500 lg:hidden dark:text-slate-400">{{ $t('user.torrents.columns.status') }}</span>
            <UBadge :color="statusColor(item)" variant="soft" size="sm">{{ statusLabel(item) }}</UBadge>
          </div>
          <div class="flex items-center justify-between gap-3 text-sm lg:block lg:text-right">
            <span class="text-xs text-slate-500 lg:hidden dark:text-slate-400">{{ $t('user.torrents.columns.size') }}</span>
            <div>
              <p class="font-medium tabular-nums text-slate-950 dark:text-white">{{ formatBytes(item.size) }}</p>
              <p class="mt-0.5 hidden text-xs text-slate-500 lg:block dark:text-slate-400">{{ $t('user.torrents.files', { count: numberFormatter.format(item.fileCount) }) }}</p>
            </div>
          </div>
          <div class="flex items-center justify-between gap-3 text-sm lg:block lg:text-right">
            <span class="text-xs text-slate-500 lg:hidden dark:text-slate-400">{{ $t('user.torrents.columns.time') }}</span>
            <div>
              <p class="font-medium tabular-nums text-slate-950 dark:text-white">{{ formatDateTime(itemTime(item), locale) }}</p>
              <p class="mt-0.5 text-xs text-slate-500 dark:text-slate-400">{{ timeLabel(item) }}</p>
            </div>
          </div>
          <div class="flex items-center justify-end gap-1" @click.stop @keydown.stop>
            <UTooltip v-if="!item.banned" :text="$t('user.torrents.edit')" :delay-duration="600">
              <UButton
                color="neutral"
                variant="ghost"
                size="xs"
                icon="i-lucide-pen-line"
                :aria-label="$t('user.torrents.edit')"
                :to="localePath(`/catalog/torrents/${item.id}/edit`)"
              />
            </UTooltip>
            <UTooltip v-if="!item.banned && item.status === TorrentStatus.Rejected" :text="$t('user.torrents.resubmit')" :delay-duration="600">
              <UButton
                color="primary"
                variant="ghost"
                size="xs"
                icon="i-lucide-send"
                :aria-label="$t('user.torrents.resubmit')"
                :loading="resubmittingId === item.id"
                :disabled="resubmittingId !== null"
                @click="resubmit(item)"
              />
            </UTooltip>
          </div>
        </article>
      </div>
    </div>

    <AppPager
      v-if="total > 0"
      class="border-t border-slate-200 px-5 py-4 dark:border-slate-800"
      size="sm"
      :page="page"
      :total="total"
      :page-size="pageSize"
      :disabled="pending"
      @page-change="goToPage"
    />
  </section>
</template>

<script setup lang="ts">
import { ApiError } from '~/composables/useApi'
import { TorrentStatus, type TorrentMineItem, type TorrentMineStatusFilter } from '~/composables/useCatalogTorrents'
import { formatBytes, formatDateTime } from '~/utils/format'

const { t, locale } = useI18n()
const localePath = useLocalePath()
const toast = useToast()
const catalog = useCatalogTorrents()

const items = ref<TorrentMineItem[]>([])
const status = ref<TorrentMineStatusFilter>(-1)
const page = ref(1)
const pageSize = 10
const total = ref(0)
const pending = ref(true)
const errorMessage = ref('')
const resubmittingId = ref<number | null>(null)

const numberFormatter = computed(() => new Intl.NumberFormat(locale.value))
const totalPages = computed(() => Math.max(1, Math.ceil(total.value / pageSize)))
const filters = computed<Array<{ value: TorrentMineStatusFilter, label: string }>>(() => [
  { value: -1, label: t('user.torrents.filters.all') },
  { value: TorrentStatus.Pending, label: t('user.torrents.filters.pending') },
  { value: TorrentStatus.Published, label: t('user.torrents.filters.published') },
  { value: TorrentStatus.Rejected, label: t('user.torrents.filters.rejected') }
])

onMounted(loadItems)

async function loadItems() {
  pending.value = true
  errorMessage.value = ''
  try {
    let data = await catalog.listMine(status.value, page.value, pageSize)
    const nextPages = Math.max(1, Math.ceil((data.total || 0) / pageSize))
    if (page.value > nextPages) {
      page.value = nextPages
      data = await catalog.listMine(status.value, page.value, pageSize)
    }
    items.value = data.list || []
    total.value = data.total || 0
  } catch (error) {
    items.value = []
    total.value = 0
    errorMessage.value = error instanceof ApiError ? error.message : t('common.requestFailed')
  } finally {
    pending.value = false
  }
}

function setStatus(value: TorrentMineStatusFilter) {
  if (status.value === value) return
  status.value = value
  page.value = 1
  void loadItems()
}

function goToPage(value: number) {
  page.value = Math.min(Math.max(1, value), totalPages.value)
  void loadItems()
}

function openTorrent(item: TorrentMineItem) {
  if (item.banned) return
  void navigateTo(localePath(`/catalog/torrents/${item.id}`))
}

async function resubmit(item: TorrentMineItem) {
  if (resubmittingId.value !== null) return
  resubmittingId.value = item.id
  try {
    await catalog.resubmitTorrent(item.id)
    toast.add({ title: t('user.torrents.resubmitSuccess'), color: 'success', icon: 'i-lucide-check-circle' })
    await loadItems()
  } catch (error) {
    toast.add({ title: error instanceof ApiError ? error.message : t('common.requestFailed'), color: 'error', icon: 'i-lucide-circle-alert' })
  } finally {
    resubmittingId.value = null
  }
}

function statusLabel(item: TorrentMineItem) {
  if (item.banned) return t('user.torrents.status.banned')
  if (item.status === TorrentStatus.Pending) return t('user.torrents.status.pending')
  if (item.status === TorrentStatus.Published) return t('user.torrents.status.published')
  if (item.status === TorrentStatus.Rejected) return t('user.torrents.status.rejected')
  return t('user.torrents.status.unknown')
}

function statusColor(item: TorrentMineItem) {
  if (item.banned || item.status === TorrentStatus.Rejected) return 'error'
  if (item.status === TorrentStatus.Pending) return 'warning'
  return 'success'
}

function itemTime(item: TorrentMineItem) {
  if (item.status === TorrentStatus.Published) return item.publishedAt || item.createdAt
  return item.submittedAt || item.updatedAt || item.createdAt
}

function timeLabel(item: TorrentMineItem) {
  return item.status === TorrentStatus.Published ? t('user.torrents.publishedAt') : t('user.torrents.submittedAt')
}
</script>
