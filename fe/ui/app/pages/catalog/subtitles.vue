<template>
  <div class="min-h-[calc(100vh-4rem)] bg-slate-50 py-6 dark:bg-slate-950">
    <div class="w-full px-3 sm:px-4 lg:px-5">
      <section class="overflow-hidden rounded-lg border border-slate-200 bg-white dark:border-slate-800 dark:bg-slate-900">
        <div class="hidden grid-cols-[minmax(0,1fr)_86px_92px_126px_96px_112px_76px] items-center gap-2 border-b border-slate-200 bg-slate-50 px-3 py-2 text-xs font-semibold text-slate-500 lg:grid dark:border-slate-800 dark:bg-slate-950/60 dark:text-slate-400">
          <span>{{ $t('catalog.subtitles.table.file') }}</span>
          <span class="text-center">{{ $t('catalog.subtitles.table.language') }}</span>
          <span class="text-center">{{ $t('catalog.subtitles.table.torrent') }}</span>
          <span>{{ $t('catalog.subtitles.table.uploader') }}</span>
          <span class="text-right">{{ $t('catalog.subtitles.table.size') }}</span>
          <span class="text-right">{{ $t('catalog.subtitles.table.createdAt') }}</span>
          <span class="text-center">{{ $t('catalog.subtitles.table.actions') }}</span>
        </div>

        <div v-if="pending" class="divide-y divide-slate-200 dark:divide-slate-800">
          <div v-for="index in 8" :key="index" class="grid gap-3 px-3 py-3 lg:grid-cols-[minmax(0,1fr)_86px_92px_126px_96px_112px_76px] lg:items-center lg:gap-2">
            <div class="space-y-2">
              <div class="h-4 w-3/4 animate-pulse rounded bg-slate-200 dark:bg-slate-800" />
              <div class="h-3 w-1/2 animate-pulse rounded bg-slate-100 dark:bg-slate-800/70" />
            </div>
            <div class="h-5 w-14 animate-pulse rounded bg-slate-100 lg:justify-self-center dark:bg-slate-800/70" />
            <div class="h-4 w-16 animate-pulse rounded bg-slate-100 lg:justify-self-center dark:bg-slate-800/70" />
            <div class="h-4 w-24 animate-pulse rounded bg-slate-100 dark:bg-slate-800/70" />
            <div class="h-4 w-16 animate-pulse rounded bg-slate-100 lg:justify-self-end dark:bg-slate-800/70" />
            <div class="h-4 w-20 animate-pulse rounded bg-slate-100 lg:justify-self-end dark:bg-slate-800/70" />
            <div class="h-8 w-16 animate-pulse rounded-md bg-slate-100 lg:justify-self-center dark:bg-slate-800/70" />
          </div>
        </div>

        <div v-else-if="errorMessage" class="flex flex-col items-center justify-center px-4 py-16 text-center">
          <UIcon name="i-lucide-circle-alert" class="size-9 text-red-500" />
          <p class="mt-3 text-sm font-medium text-slate-950 dark:text-white">{{ errorMessage }}</p>
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
            class="grid gap-3 px-3 py-3 transition-colors hover:bg-slate-50 lg:grid-cols-[minmax(0,1fr)_86px_92px_126px_96px_112px_76px] lg:items-center lg:gap-2 lg:py-2.5 dark:hover:bg-slate-950/70"
          >
            <div class="min-w-0">
              <h2 class="truncate text-sm font-medium leading-5 text-slate-950 dark:text-white" :title="subtitleFileName(subtitle)">
                {{ subtitleFileName(subtitle) }}
              </h2>
              <div class="mt-1 flex flex-wrap items-center gap-x-3 gap-y-1 text-xs text-slate-500 lg:hidden dark:text-slate-400">
                <span class="inline-flex items-center rounded border border-slate-200 px-1.5 py-0.5 font-medium text-slate-600 dark:border-slate-700 dark:text-slate-300">
                  {{ subtitleLanguage(subtitle) }}
                </span>
                <span>{{ formatBytes(subtitle.size) }}</span>
                <span :title="formatDateTime(subtitle.createdAt, locale)">{{ formatDateOnly(subtitle.createdAt, locale) }}</span>
                <span>{{ uploaderName(subtitle) }}</span>
              </div>
            </div>

            <div class="hidden justify-center lg:flex">
              <UBadge color="neutral" variant="soft">{{ subtitleLanguage(subtitle) }}</UBadge>
            </div>

            <NuxtLink
              :to="localePath(`/catalog/torrents/${subtitle.torrentId}`)"
              class="hidden min-w-0 justify-center gap-1 text-sm font-medium text-sky-700 hover:text-sky-800 lg:inline-flex dark:text-sky-300 dark:hover:text-sky-200"
              :title="$t('catalog.subtitles.actions.openTorrent')"
            >
              <UIcon name="i-lucide-arrow-up-right" class="size-4" />
              #{{ subtitle.torrentId }}
            </NuxtLink>

            <p class="hidden truncate text-sm text-slate-600 lg:block dark:text-slate-300" :title="uploaderName(subtitle)">
              {{ uploaderName(subtitle) }}
            </p>
            <p class="hidden text-right text-sm font-medium tabular-nums text-slate-950 lg:block dark:text-white">{{ formatBytes(subtitle.size) }}</p>
            <p class="hidden text-right text-sm text-slate-500 lg:block dark:text-slate-400" :title="formatDateTime(subtitle.createdAt, locale)">
              {{ formatDateOnly(subtitle.createdAt, locale) }}
            </p>

            <div class="hidden justify-center gap-1 lg:flex">
              <AppPermissionButton
                :permission="Permission.CatalogSubtitleDownload"
                color="neutral"
                variant="ghost"
                size="xs"
                icon="i-lucide-download"
                :tooltip="$t('catalog.subtitles.actions.download')"
                :loading="downloadPendingId === subtitle.id"
                :disabled="downloadPendingId > 0"
                :aria-label="$t('catalog.subtitles.actions.download')"
                @click="handleDownload(subtitle)"
              />
              <UTooltip
                :text="$t('catalog.subtitles.actions.report')"
                :content="{ side: 'top', sideOffset: 8 }"
                :delay-duration="120"
              >
                <UButton
                  color="neutral"
                  :variant="activeReportId === subtitle.id ? 'soft' : 'ghost'"
                  size="xs"
                  icon="i-lucide-flag"
                  :aria-label="$t('catalog.subtitles.actions.report')"
                  @click="startReport(subtitle.id)"
                />
              </UTooltip>
            </div>

            <div class="flex flex-wrap items-center gap-2 lg:hidden">
              <UButton color="neutral" variant="outline" size="xs" icon="i-lucide-arrow-up-right" :to="localePath(`/catalog/torrents/${subtitle.torrentId}`)">
                #{{ subtitle.torrentId }}
              </UButton>
              <AppPermissionButton
                :permission="Permission.CatalogSubtitleDownload"
                color="neutral"
                variant="outline"
                size="xs"
                icon="i-lucide-download"
                :tooltip="$t('catalog.subtitles.actions.download')"
                :loading="downloadPendingId === subtitle.id"
                :disabled="downloadPendingId > 0"
                @click="handleDownload(subtitle)"
              >
                {{ $t('catalog.subtitles.actions.download') }}
              </AppPermissionButton>
              <UButton color="neutral" :variant="activeReportId === subtitle.id ? 'soft' : 'ghost'" size="xs" icon="i-lucide-flag" @click="startReport(subtitle.id)">
                {{ $t('catalog.subtitles.actions.report') }}
              </UButton>
            </div>

            <form v-if="activeReportId === subtitle.id" class="grid gap-2 rounded-md border border-slate-200 bg-slate-50 p-3 lg:col-span-7 dark:border-slate-800 dark:bg-slate-950" @submit.prevent="handleReport(subtitle.id)">
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
      </section>

      <AppPager
        v-if="total > 0"
        class="mt-4"
        :page="page"
        :total="total"
        :page-size="size"
        :page-size-options="pageSizes"
        :disabled="pending"
        @page-change="goToPage"
        @page-size-change="handlePageSizeChange"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ApiError } from '~/composables/useApi'
import type { SubtitleItem } from '~/composables/useCatalogTorrents'
import { formatBytes, formatDateOnly, formatDateTime } from '~/utils/format'

definePageMeta({
  middleware: 'auth'
})

const { t, locale } = useI18n()
const localePath = useLocalePath()
const route = useRoute()
const router = useRouter()
const toast = useToast()
const catalogTorrents = useCatalogTorrents()
const { hasPermission } = useAuth()

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
const pageSizes = [20, 50, 100]
const subtitleLanguageLabels: Record<string, string> = {
  'zh-CN': '简体中文',
  'zh-TW': '繁體中文',
  'en-US': 'English',
  'ja-JP': '日本語',
  'ko-KR': '한국어',
  other: 'Other'
}

const size = computed(() => Number(selectedSize.value))
const totalPages = computed(() => Math.max(1, Math.ceil(total.value / size.value)))
const canDownloadSubtitle = computed(() => hasPermission(Permission.CatalogSubtitleDownload))

useHead(() => ({
  title: t('catalog.subtitles.metaTitle')
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
  const parsed = readPositiveIntQuery('size', 50)
  return [20, 50, 100].includes(parsed) ? parsed : 50
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

function handlePageSizeChange(nextSize: number) {
  selectedSize.value = String(nextSize)
  page.value = 1
  loadSubtitles()
}

function syncQuery() {
  router.replace({
    query: {
      ...route.query,
      page: page.value > 1 ? String(page.value) : undefined,
      size: size.value !== 50 ? String(size.value) : undefined
    }
  })
}

function subtitleFileName(subtitle: SubtitleItem) {
  return subtitle.fileName || `subtitle-${subtitle.id}`
}

function subtitleLanguage(subtitle: SubtitleItem) {
  if (!subtitle.language) return '-'
  return subtitleLanguageLabels[subtitle.language] || subtitle.language
}

function uploaderName(subtitle: SubtitleItem) {
  if (subtitle.anonymous) {
    const name = rawUploaderName(subtitle)
    return name ? t('catalog.torrents.anonymousOwner', { name }) : t('catalog.torrents.anonymous')
  }
  return subtitle.uploader?.id > 0 ? rawUploaderName(subtitle) : '-'
}

function rawUploaderName(subtitle: SubtitleItem) {
  return subtitle.uploader?.username || (subtitle.uploader?.id > 0 ? `#${subtitle.uploader.id}` : '')
}

async function handleDownload(subtitle: SubtitleItem) {
  if (!canDownloadSubtitle.value || downloadPendingId.value > 0) return

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
