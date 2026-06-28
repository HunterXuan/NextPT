<template>
  <div class="min-h-[calc(100vh-4rem)] bg-slate-50 py-6 dark:bg-slate-950">
    <div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
      <div class="mb-6 flex flex-col gap-4 lg:flex-row lg:items-end lg:justify-between">
        <div>
          <UButton color="neutral" variant="ghost" icon="i-lucide-arrow-left" :to="localePath('/catalog/torrents')">
            {{ $t('catalog.torrents.detail.back') }}
          </UButton>
          <p class="mt-4 text-sm font-medium text-slate-500 dark:text-slate-400">{{ $t('catalog.subtitles.eyebrow') }}</p>
          <h1 class="mt-1 text-2xl font-semibold text-slate-950 dark:text-white">{{ $t('catalog.subtitles.title') }}</h1>
        </div>

        <div class="grid grid-cols-2 gap-3 sm:flex sm:items-center">
          <div class="rounded-lg border border-slate-200 bg-white px-4 py-3 dark:border-slate-800 dark:bg-slate-900">
            <p class="text-xs text-slate-500 dark:text-slate-400">{{ $t('catalog.subtitles.summary.total') }}</p>
            <p class="mt-1 text-lg font-semibold text-slate-950 dark:text-white">{{ numberFormatter.format(total) }}</p>
          </div>
          <UButton color="neutral" variant="outline" icon="i-lucide-refresh-cw" :loading="pending" @click="loadSubtitles">
            {{ $t('common.refresh') }}
          </UButton>
        </div>
      </div>

      <div class="overflow-hidden rounded-lg border border-slate-200 bg-white dark:border-slate-800 dark:bg-slate-900">
        <div class="hidden grid-cols-[minmax(0,1fr)_120px_150px_120px_150px_170px] gap-4 border-b border-slate-200 px-4 py-3 text-xs font-medium uppercase text-slate-500 lg:grid dark:border-slate-800 dark:text-slate-400">
          <span>{{ $t('catalog.subtitles.table.file') }}</span>
          <span>{{ $t('catalog.subtitles.table.torrent') }}</span>
          <span>{{ $t('catalog.subtitles.table.uploader') }}</span>
          <span>{{ $t('catalog.subtitles.table.size') }}</span>
          <span>{{ $t('catalog.subtitles.table.createdAt') }}</span>
          <span class="text-right">{{ $t('catalog.subtitles.table.actions') }}</span>
        </div>

        <div v-if="pending" class="divide-y divide-slate-200 dark:divide-slate-800">
          <div v-for="index in 6" :key="index" class="grid gap-4 px-4 py-4 lg:grid-cols-[minmax(0,1fr)_120px_150px_120px_150px_170px]">
            <div class="space-y-2">
              <div class="h-4 w-3/4 animate-pulse rounded bg-slate-200 dark:bg-slate-800" />
              <div class="h-3 w-1/2 animate-pulse rounded bg-slate-100 dark:bg-slate-800/70" />
            </div>
            <div class="h-4 w-20 animate-pulse rounded bg-slate-100 dark:bg-slate-800/70" />
            <div class="h-4 w-24 animate-pulse rounded bg-slate-100 dark:bg-slate-800/70" />
            <div class="h-4 w-16 animate-pulse rounded bg-slate-100 dark:bg-slate-800/70" />
            <div class="h-4 w-28 animate-pulse rounded bg-slate-100 dark:bg-slate-800/70" />
            <div class="ml-auto h-4 w-28 animate-pulse rounded bg-slate-100 dark:bg-slate-800/70" />
          </div>
        </div>

        <div v-else-if="errorMessage" class="flex flex-col items-center justify-center px-4 py-16 text-center">
          <UIcon name="i-lucide-circle-alert" class="size-9 text-red-500" />
          <p class="mt-3 text-sm font-medium text-slate-950 dark:text-white">{{ errorMessage }}</p>
          <UButton class="mt-5" color="neutral" variant="outline" icon="i-lucide-refresh-cw" @click="loadSubtitles">
            {{ $t('common.retry') }}
          </UButton>
        </div>

        <div v-else-if="subtitles.length === 0" class="flex flex-col items-center justify-center px-4 py-16 text-center">
          <UIcon name="i-lucide-captions-off" class="size-9 text-slate-400" />
          <p class="mt-3 text-sm font-medium text-slate-950 dark:text-white">{{ $t('catalog.subtitles.empty.title') }}</p>
          <p class="mt-1 max-w-md text-sm text-slate-500 dark:text-slate-400">{{ $t('catalog.subtitles.empty.description') }}</p>
        </div>

        <div v-else class="divide-y divide-slate-200 dark:divide-slate-800">
          <article
            v-for="subtitle in subtitles"
            :key="subtitle.id"
            class="grid gap-4 px-4 py-4 transition-colors hover:bg-slate-50 lg:grid-cols-[minmax(0,1fr)_120px_150px_120px_150px_170px] lg:items-center dark:hover:bg-slate-950/70"
          >
            <div class="min-w-0">
              <div class="flex min-w-0 flex-wrap items-center gap-2">
                <UIcon name="i-lucide-captions" class="size-4 shrink-0 text-sky-500" />
                <h2 class="min-w-0 truncate text-sm font-semibold text-slate-950 dark:text-white">{{ subtitle.fileName || `#${subtitle.id}` }}</h2>
                <UBadge color="neutral" variant="soft">{{ subtitle.language || '-' }}</UBadge>
              </div>
              <p class="mt-1 text-xs text-slate-500 dark:text-slate-400">#{{ subtitle.id }}</p>
            </div>

            <NuxtLink
              :to="localePath(`/catalog/torrents/${subtitle.torrentId}`)"
              class="inline-flex items-center gap-1 text-sm font-medium text-sky-700 hover:text-sky-800 dark:text-sky-300 dark:hover:text-sky-200"
            >
              <UIcon name="i-lucide-arrow-up-right" class="size-4" />
              #{{ subtitle.torrentId }}
            </NuxtLink>

            <p class="truncate text-sm text-slate-600 dark:text-slate-300">{{ subtitle.username || `#${subtitle.userId}` }}</p>
            <p class="text-sm font-medium text-slate-950 dark:text-white">{{ formatBytes(subtitle.size) }}</p>
            <p class="text-sm text-slate-500 dark:text-slate-400">{{ formatDateTime(subtitle.createdAt, locale) }}</p>

            <div class="flex flex-wrap items-center justify-start gap-2 lg:justify-end">
              <UButton color="neutral" variant="outline" size="xs" icon="i-lucide-download" :loading="downloadPendingId === subtitle.id" :disabled="downloadPendingId > 0" @click="handleDownload(subtitle)">
                {{ $t('catalog.subtitles.actions.download') }}
              </UButton>
              <UButton color="neutral" variant="ghost" size="xs" icon="i-lucide-flag" @click="startReport(subtitle.id)">
                {{ $t('catalog.subtitles.actions.report') }}
              </UButton>
            </div>

            <form v-if="activeReportId === subtitle.id" class="grid gap-2 rounded-md bg-slate-50 p-3 lg:col-span-6 dark:bg-slate-950" @submit.prevent="handleReport(subtitle.id)">
              <UTextarea v-model="reportReason" :rows="2" :placeholder="$t('catalog.subtitles.report.reason')" :disabled="reportPending" />
              <div class="flex justify-end gap-2">
                <UButton color="neutral" variant="ghost" size="xs" type="button" @click="activeReportId = 0">{{ $t('common.cancel') }}</UButton>
                <UButton color="error" variant="soft" size="xs" type="submit" :loading="reportPending" :disabled="reportReason.trim().length < 5">
                  {{ $t('catalog.subtitles.report.submit') }}
                </UButton>
              </div>
            </form>
          </article>
        </div>
      </div>

      <div class="mt-4 flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
        <p class="text-sm text-slate-500 dark:text-slate-400">
          {{ $t('catalog.subtitles.pagination.summary', { page: page, pages: totalPages }) }}
        </p>
        <div class="flex flex-wrap items-center gap-2">
          <label class="flex items-center gap-2 text-sm text-slate-500 dark:text-slate-400">
            <span>{{ $t('catalog.subtitles.pagination.pageSize') }}</span>
            <select
              v-model="selectedSize"
              class="h-9 rounded-md border border-slate-200 bg-white px-2 text-sm text-slate-950 outline-none transition focus:border-sky-400 focus:ring-2 focus:ring-sky-100 dark:border-slate-700 dark:bg-slate-900 dark:text-white dark:focus:border-sky-500 dark:focus:ring-sky-950"
              :disabled="pending"
              @change="handlePageSizeChange"
            >
              <option value="20">20</option>
              <option value="50">50</option>
              <option value="100">100</option>
            </select>
          </label>
          <UButton color="neutral" variant="outline" icon="i-lucide-chevron-left" :disabled="page <= 1 || pending" @click="goToPage(page - 1)">
            {{ $t('common.previous') }}
          </UButton>
          <UButton color="neutral" variant="outline" trailing-icon="i-lucide-chevron-right" :disabled="page >= totalPages || pending" @click="goToPage(page + 1)">
            {{ $t('common.next') }}
          </UButton>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ApiError } from '~/composables/useApi'
import type { SubtitleItem } from '~/composables/useCatalogTorrents'
import { formatBytes, formatDateTime } from '~/utils/format'

definePageMeta({
  middleware: 'auth'
})

const { t, locale } = useI18n()
const localePath = useLocalePath()
const route = useRoute()
const router = useRouter()
const toast = useToast()
const catalogTorrents = useCatalogTorrents()

const subtitles = ref<SubtitleItem[]>([])
const total = ref(0)
const page = ref(readPositiveIntQuery('page', 1))
const selectedSize = ref(String(readPageSizeQuery()))
const pending = ref(false)
const errorMessage = ref('')
const downloadPendingId = ref(0)
const activeReportId = ref(0)
const reportPending = ref(false)
const reportReason = ref('')

const numberFormatter = computed(() => new Intl.NumberFormat(locale.value))
const size = computed(() => Number(selectedSize.value))
const totalPages = computed(() => Math.max(1, Math.ceil(total.value / size.value)))

useHead(() => ({
  title: `${t('catalog.subtitles.metaTitle')} - NextPT`
}))

onMounted(loadSubtitles)

function readFirstQueryValue(key: string) {
  const value = route.query[key]
  return Array.isArray(value) ? value[0] : value
}

function readPositiveIntQuery(key: string, fallback: number) {
  const parsed = Number(readFirstQueryValue(key))
  return Number.isInteger(parsed) && parsed > 0 ? parsed : fallback
}

function readPageSizeQuery() {
  const parsed = readPositiveIntQuery('size', 20)
  return [20, 50, 100].includes(parsed) ? parsed : 20
}

async function loadSubtitles() {
  pending.value = true
  errorMessage.value = ''
  syncQuery()

  try {
    let data = await catalogTorrents.listAllSubtitles(page.value, size.value)
    const nextTotal = data.total || 0
    const nextTotalPages = Math.max(1, Math.ceil(nextTotal / size.value))
    if (page.value > nextTotalPages) {
      page.value = nextTotalPages
      data = await catalogTorrents.listAllSubtitles(page.value, size.value)
    }
    subtitles.value = data.list || []
    total.value = data.total || 0
  } catch (error) {
    subtitles.value = []
    total.value = 0
    errorMessage.value = error instanceof ApiError ? error.message : t('common.requestFailed')
  } finally {
    pending.value = false
  }
}

function goToPage(nextPage: number) {
  page.value = Math.min(Math.max(1, nextPage), totalPages.value)
  loadSubtitles()
}

function handlePageSizeChange() {
  page.value = 1
  loadSubtitles()
}

function syncQuery() {
  router.replace({
    query: {
      ...route.query,
      page: page.value > 1 ? String(page.value) : undefined,
      size: size.value !== 20 ? String(size.value) : undefined
    }
  })
}

async function handleDownload(subtitle: SubtitleItem) {
  if (downloadPendingId.value > 0) return

  downloadPendingId.value = subtitle.id
  try {
    const out = await catalogTorrents.downloadSubtitle(subtitle.id)
    downloadBlob(out.blob, out.filename || subtitle.fileName || `subtitle-${subtitle.id}`)
  } catch (error) {
    toast.add({
      title: error instanceof ApiError ? error.message : t('common.requestFailed'),
      color: 'error',
      icon: 'i-lucide-circle-alert'
    })
  } finally {
    downloadPendingId.value = 0
  }
}

function startReport(subtitleId: number) {
  activeReportId.value = activeReportId.value === subtitleId ? 0 : subtitleId
  reportReason.value = ''
}

async function handleReport(subtitleId: number) {
  const reason = reportReason.value.trim()
  if (reportPending.value || reason.length < 5) return

  reportPending.value = true
  try {
    await catalogTorrents.reportSubtitle(subtitleId, reason)
    activeReportId.value = 0
    reportReason.value = ''
    toast.add({
      title: t('catalog.subtitles.report.success'),
      color: 'success',
      icon: 'i-lucide-check-circle'
    })
  } catch (error) {
    toast.add({
      title: error instanceof ApiError ? error.message : t('common.requestFailed'),
      color: 'error',
      icon: 'i-lucide-circle-alert'
    })
  } finally {
    reportPending.value = false
  }
}

function downloadBlob(blob: Blob, filename: string) {
  const href = URL.createObjectURL(blob)
  const link = document.createElement('a')
  link.href = href
  link.download = filename
  document.body.appendChild(link)
  link.click()
  link.remove()
  URL.revokeObjectURL(href)
}
</script>
