<template>
  <div class="min-h-[calc(100vh-4rem)] bg-slate-50 py-6 dark:bg-slate-950">
    <div class="w-full px-3 sm:px-4 lg:px-5">
      <section class="mb-3 rounded-lg border border-slate-200 bg-white p-2.5 dark:border-slate-800 dark:bg-slate-900">
        <div class="grid gap-3 xl:grid-cols-[minmax(0,1fr)_auto] xl:items-start">
          <form class="grid gap-2 sm:grid-cols-[minmax(0,1fr)_auto]" @submit.prevent="handleSearchSubmit">
            <UInput
              v-model="keyword"
              class="w-full"
              icon="i-lucide-search"
              :placeholder="$t('catalog.torrents.search.placeholder')"
              :disabled="pending"
            />
            <UButton type="submit" color="primary" icon="i-lucide-search" :loading="pending">
              {{ $t('catalog.torrents.search.submit') }}
            </UButton>
          </form>

          <div class="grid grid-cols-2 gap-2 sm:flex sm:items-center sm:justify-end">
            <UButton
              color="neutral"
              variant="outline"
              icon="i-lucide-refresh-cw"
              :loading="pending"
              @click="loadTorrents"
            >
              {{ $t('common.refresh') }}
            </UButton>
            <UButton color="primary" icon="i-lucide-upload" :to="localePath('/catalog/torrents/upload')">
              {{ $t('catalog.torrents.upload.action') }}
            </UButton>
          </div>
        </div>

        <div class="mt-2 -mx-1 overflow-x-auto px-1 pb-1">
          <div class="flex min-w-max items-center gap-1.5">
            <button
              type="button"
              class="inline-flex h-7 shrink-0 items-center rounded-md border px-2.5 text-xs font-medium transition-colors"
              :class="selectedCategoryIds.length === 0
                ? 'border-slate-950 bg-slate-950 text-white dark:border-white dark:bg-white dark:text-slate-950'
                : 'border-slate-200 text-slate-600 hover:border-slate-300 hover:bg-slate-50 dark:border-slate-700 dark:text-slate-300 dark:hover:border-slate-600 dark:hover:bg-slate-950'"
              :disabled="pending"
              @click="clearCategories"
            >
              {{ $t('catalog.torrents.filters.allCategories') }}
            </button>
            <button
              v-for="category in categories"
              :key="category.id"
              type="button"
              class="inline-flex h-7 max-w-36 shrink-0 items-center rounded-md border px-2.5 text-xs font-medium transition-colors"
              :class="selectedCategoryIds.includes(category.id)
                ? 'border-sky-300 bg-sky-50 text-sky-800 dark:border-sky-700 dark:bg-sky-950 dark:text-sky-200'
                : 'border-slate-200 text-slate-600 hover:border-slate-300 hover:bg-slate-50 dark:border-slate-700 dark:text-slate-300 dark:hover:border-slate-600 dark:hover:bg-slate-950'"
              :disabled="pending"
              @click="toggleCategory(category.id)"
            >
              <span class="truncate">{{ categoryDisplayName(category) }}</span>
            </button>
            <UButton
              v-if="selectedCategoryIds.length > 0"
              class="shrink-0"
              color="neutral"
              variant="ghost"
              size="xs"
              icon="i-lucide-x"
              @click="clearCategories"
            >
              {{ $t('catalog.torrents.filters.clearCategories') }}
            </UButton>
          </div>
        </div>

        <p v-if="categories.length === 0" class="mt-3 text-sm text-slate-500 dark:text-slate-400">
          {{ $t('catalog.torrents.filters.noCategories') }}
        </p>
      </section>

      <section class="overflow-hidden rounded-lg border border-slate-200 bg-white dark:border-slate-800 dark:bg-slate-900">
        <div class="hidden grid-cols-[86px_minmax(0,1fr)_88px_48px_48px_56px_96px_110px_34px] items-center gap-2 border-b border-slate-200 bg-slate-50 px-2.5 py-2 text-xs font-semibold text-slate-500 lg:grid dark:border-slate-800 dark:bg-slate-950/60 dark:text-slate-400">
          <span>{{ $t('catalog.torrents.table.category') }}</span>
          <span>{{ $t('catalog.torrents.table.torrent') }}</span>
          <span class="text-right">{{ $t('catalog.torrents.table.size') }}</span>
          <span class="text-right tabular-nums">{{ $t('catalog.torrents.table.seeders') }}</span>
          <span class="text-right tabular-nums">{{ $t('catalog.torrents.table.leechers') }}</span>
          <span class="text-right tabular-nums">{{ $t('catalog.torrents.table.completed') }}</span>
          <span class="text-right">{{ $t('catalog.torrents.table.createdAt') }}</span>
          <span class="text-center">{{ $t('catalog.torrents.table.uploader') }}</span>
          <span class="text-center">{{ $t('catalog.torrents.table.actions') }}</span>
        </div>

        <div v-if="pending" class="divide-y divide-slate-200 dark:divide-slate-800">
          <div v-for="index in 8" :key="index" class="grid gap-2 px-2.5 py-2.5 lg:grid-cols-[86px_minmax(0,1fr)_88px_48px_48px_56px_96px_110px_34px] lg:items-center">
            <div class="h-5 w-16 animate-pulse rounded bg-slate-100 dark:bg-slate-800/70" />
            <div class="space-y-2">
              <div class="h-4 w-4/5 animate-pulse rounded bg-slate-200 dark:bg-slate-800" />
              <div class="h-3 w-2/5 animate-pulse rounded bg-slate-100 dark:bg-slate-800/70" />
            </div>
            <div class="h-4 w-16 animate-pulse rounded bg-slate-100 lg:justify-self-end dark:bg-slate-800/70" />
            <div class="h-4 w-8 animate-pulse rounded bg-slate-100 lg:justify-self-end dark:bg-slate-800/70" />
            <div class="h-4 w-8 animate-pulse rounded bg-slate-100 lg:justify-self-end dark:bg-slate-800/70" />
            <div class="h-4 w-8 animate-pulse rounded bg-slate-100 lg:justify-self-end dark:bg-slate-800/70" />
            <div class="h-4 w-24 animate-pulse rounded bg-slate-100 lg:justify-self-end dark:bg-slate-800/70" />
            <div class="h-4 w-20 animate-pulse rounded bg-slate-100 lg:justify-self-center dark:bg-slate-800/70" />
            <div class="h-8 w-8 animate-pulse rounded-md bg-slate-100 lg:justify-self-center dark:bg-slate-800/70" />
          </div>
        </div>

        <div v-else-if="errorMessage" class="flex flex-col items-center justify-center px-4 py-16 text-center">
          <UIcon name="i-lucide-circle-alert" class="size-9 text-red-500" />
          <p class="mt-3 text-sm font-medium text-slate-950 dark:text-white">{{ errorMessage }}</p>
          <UButton class="mt-5" color="neutral" variant="outline" icon="i-lucide-refresh-cw" @click="loadTorrents">
            {{ $t('common.retry') }}
          </UButton>
        </div>

        <div v-else-if="torrents.length === 0" class="flex flex-col items-center justify-center px-4 py-16 text-center">
          <UIcon name="i-lucide-inbox" class="size-9 text-slate-400" />
          <p class="mt-3 text-sm font-medium text-slate-950 dark:text-white">{{ $t('catalog.torrents.empty.title') }}</p>
          <p class="mt-1 max-w-md text-sm text-slate-500 dark:text-slate-400">{{ $t('catalog.torrents.empty.description') }}</p>
        </div>

        <div v-else class="divide-y divide-slate-200 dark:divide-slate-800">
          <article
            v-for="torrent in torrents"
            :key="torrent.id"
            class="grid gap-3 px-3 py-3 transition-colors hover:bg-slate-50 lg:grid-cols-[86px_minmax(0,1fr)_88px_48px_48px_56px_96px_110px_34px] lg:items-center lg:gap-2 lg:px-2.5 lg:py-2 dark:hover:bg-slate-950/70"
          >
            <div class="flex flex-wrap items-center gap-2 lg:block">
              <span class="inline-flex h-6 max-w-full items-center rounded border border-slate-200 bg-slate-50 px-2 text-[11px] font-medium text-slate-600 dark:border-slate-700 dark:bg-slate-950 dark:text-slate-300">
                <span class="truncate">{{ categoryName(torrent.categoryId) }}</span>
              </span>
            </div>

            <div class="min-w-0">
              <div class="flex min-w-0 items-start justify-between gap-3">
                <div class="min-w-0">
                  <h2 class="flex min-w-0 items-center gap-1.5 text-sm font-medium leading-5">
                    <NuxtLink
                      :to="localePath(`/catalog/torrents/${torrent.id}`)"
                      class="min-w-0 truncate text-slate-950 hover:text-sky-700 dark:text-white dark:hover:text-sky-300"
                    >
                      {{ torrent.name || `#${torrent.id}` }}
                    </NuxtLink>
                    <span
                      v-for="badge in torrentStatusBadges(torrent)"
                      :key="badge.key"
                      class="inline-flex h-5 shrink-0 items-center rounded px-1.5 text-[10px] font-semibold leading-none"
                      :class="badge.class"
                      :title="badge.title"
                    >
                      {{ badge.label }}
                    </span>
                  </h2>
                  <p v-if="torrent.subTitle" class="mt-0.5 truncate text-xs text-slate-500 dark:text-slate-400">
                    {{ torrent.subTitle }}
                  </p>
                </div>
                <UButton
                  class="lg:hidden"
                  color="neutral"
                  variant="ghost"
                  size="xs"
                  icon="i-lucide-download"
                  :aria-label="$t('catalog.torrents.detail.actions.download')"
                  :title="$t('catalog.torrents.detail.actions.download')"
                  :loading="downloadPendingId === torrent.id"
                  :disabled="downloadPendingId > 0"
                  @click="handleDownloadTorrent(torrent)"
                />
              </div>
              <div class="mt-1 flex flex-wrap items-center gap-x-3 gap-y-1 text-xs text-slate-500 dark:text-slate-400">
                <span class="lg:hidden">{{ torrentOwnerName(torrent) }}</span>
                <span class="lg:hidden" :title="formatDateTime(torrent.createdAt, locale)">{{ relativeDateTime(torrent.createdAt) }}</span>
              </div>
            </div>

            <div class="hidden text-right tabular-nums lg:block">
              <p class="text-sm font-medium text-slate-950 dark:text-white">{{ formatBytes(torrent.size) }}</p>
            </div>

            <div class="hidden text-right text-sm font-semibold tabular-nums text-emerald-600 lg:block dark:text-emerald-400">
              {{ numberFormatter.format(torrent.seeders) }}
            </div>
            <div class="hidden text-right text-sm font-semibold tabular-nums text-sky-600 lg:block dark:text-sky-400">
              {{ numberFormatter.format(torrent.leechers) }}
            </div>
            <div class="hidden text-right text-sm font-semibold tabular-nums text-slate-700 lg:block dark:text-slate-200">
              {{ numberFormatter.format(torrent.snatched) }}
            </div>
            <p class="hidden text-right text-sm text-slate-500 lg:block dark:text-slate-400" :title="formatDateTime(torrent.createdAt, locale)">
              {{ relativeDateTime(torrent.createdAt) }}
            </p>
            <div class="hidden min-w-0 text-center lg:block" :title="torrentOwnerName(torrent)">
              <p class="truncate text-sm leading-5 text-slate-600 dark:text-slate-300">
                {{ torrentOwnerPrimary(torrent) }}
              </p>
              <p v-if="torrentOwnerSecondary(torrent)" class="truncate text-xs leading-4 text-slate-400 dark:text-slate-500">
                {{ torrentOwnerSecondary(torrent) }}
              </p>
            </div>
            <div class="hidden justify-center lg:flex">
              <UButton
                class="text-slate-400 hover:text-slate-700 dark:hover:text-slate-200"
                color="neutral"
                variant="ghost"
                size="xs"
                icon="i-lucide-download"
                :aria-label="$t('catalog.torrents.detail.actions.download')"
                :title="$t('catalog.torrents.detail.actions.download')"
                :loading="downloadPendingId === torrent.id"
                :disabled="downloadPendingId > 0"
                @click="handleDownloadTorrent(torrent)"
              />
            </div>

            <div class="grid grid-cols-4 gap-2 rounded-md bg-slate-50 px-3 py-2 text-xs lg:hidden dark:bg-slate-950">
              <div>
                <p class="text-slate-500 dark:text-slate-400">{{ $t('catalog.torrents.table.size') }}</p>
                <p class="mt-1 font-medium text-slate-950 dark:text-white">{{ formatBytes(torrent.size) }}</p>
              </div>
              <div>
                <p class="text-slate-500 dark:text-slate-400">{{ $t('catalog.torrents.table.seeders') }}</p>
                <p class="mt-1 font-semibold text-emerald-600 dark:text-emerald-400">{{ numberFormatter.format(torrent.seeders) }}</p>
              </div>
              <div>
                <p class="text-slate-500 dark:text-slate-400">{{ $t('catalog.torrents.table.leechers') }}</p>
                <p class="mt-1 font-semibold text-sky-600 dark:text-sky-400">{{ numberFormatter.format(torrent.leechers) }}</p>
              </div>
              <div>
                <p class="text-slate-500 dark:text-slate-400">{{ $t('catalog.torrents.table.completed') }}</p>
                <p class="mt-1 font-semibold text-slate-700 dark:text-slate-200">{{ numberFormatter.format(torrent.snatched) }}</p>
              </div>
            </div>
          </article>
        </div>
      </section>

      <AppPager
        class="mt-4"
        :page="page"
        :total="total"
        :page-size="Number(selectedSize)"
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
import type { CatalogCategory, TorrentListItem } from '~/composables/useCatalogTorrents'

definePageMeta({
  middleware: 'auth'
})

const { t, locale } = useI18n()
const localePath = useLocalePath()
const route = useRoute()
const router = useRouter()
const toast = useToast()
const catalogTorrents = useCatalogTorrents()

const categories = ref<CatalogCategory[]>([])
const torrents = ref<TorrentListItem[]>([])
const total = ref(0)
const pending = ref(false)
const errorMessage = ref('')
const downloadPendingId = ref(0)
const pageSizes = [20, 50, 100]

const page = ref(readPositiveIntQuery('page', 1))
const keyword = ref(readStringQuery('keyword'))
const selectedCategoryIds = ref(readCategoryIdsQuery())
const selectedSize = ref(String(readPageSizeQuery()))

const numberFormatter = computed(() => new Intl.NumberFormat(locale.value))
const relativeTimeFormatter = computed(() => new Intl.RelativeTimeFormat(locale.value, { numeric: 'auto' }))
const totalPages = computed(() => Math.max(1, Math.ceil(total.value / Number(selectedSize.value || 50))))

const categoryNameMap = computed(() => {
  const map = new Map<number, string>()
  for (const category of categories.value) {
    map.set(category.id, categoryDisplayName(category))
  }
  return map
})

onMounted(async () => {
  await Promise.all([loadCategories(), loadTorrents()])
})

function readFirstQueryValue(key: string) {
  const value = route.query[key]
  return Array.isArray(value) ? value[0] : value
}

function readStringQuery(key: string) {
  return String(readFirstQueryValue(key) || '').trim()
}

function readPositiveIntQuery(key: string, fallback: number) {
  const parsed = Number(readFirstQueryValue(key))
  return Number.isInteger(parsed) && parsed > 0 ? parsed : fallback
}

function readPageSizeQuery() {
  const parsed = readPositiveIntQuery('size', 50)
  return [20, 50, 100].includes(parsed) ? parsed : 50
}

function readCategoryIdsQuery() {
  const rawValues = route.query.categoryIds || route.query['categoryIds[]']
  const values = Array.isArray(rawValues) ? rawValues : [rawValues]
  const ids = values
    .flatMap((value) => String(value || '').split(','))
    .map((value) => Number(value))
    .filter((value) => Number.isInteger(value) && value > 0)

  return [...new Set(ids)]
}

async function loadCategories() {
  try {
    const data = await catalogTorrents.listCategories()
    categories.value = data.list || []
  } catch {
    categories.value = []
  }
}

async function loadTorrents() {
  pending.value = true
  errorMessage.value = ''
  syncQuery()

  try {
    const data = await catalogTorrents.listTorrents({
      page: page.value,
      size: Number(selectedSize.value),
      keyword: keyword.value,
      categoryIds: selectedCategoryIds.value
    })
    torrents.value = data.list || []
    total.value = data.total || 0

    if (page.value > totalPages.value) {
      page.value = totalPages.value
      await loadTorrents()
    }
  } catch (error) {
    torrents.value = []
    total.value = 0
    errorMessage.value = error instanceof ApiError ? error.message : t('common.requestFailed')
  } finally {
    pending.value = false
  }
}

function handleSearchSubmit() {
  page.value = 1
  loadTorrents()
}

function toggleCategory(categoryId: number) {
  selectedCategoryIds.value = selectedCategoryIds.value.includes(categoryId)
    ? selectedCategoryIds.value.filter((id) => id !== categoryId)
    : [...selectedCategoryIds.value, categoryId]
  page.value = 1
  loadTorrents()
}

function clearCategories() {
  selectedCategoryIds.value = []
  page.value = 1
  loadTorrents()
}

function handlePageSizeChange(nextSize: number) {
  selectedSize.value = String(nextSize)
  page.value = 1
  loadTorrents()
}

function goToPage(nextPage: number) {
  page.value = Math.min(Math.max(1, nextPage), totalPages.value)
  loadTorrents()
}

function syncQuery() {
  router.replace({
    query: {
      page: page.value > 1 ? String(page.value) : undefined,
      size: selectedSize.value !== '50' ? selectedSize.value : undefined,
      keyword: keyword.value.trim() || undefined,
      categoryIds: selectedCategoryIds.value.length ? selectedCategoryIds.value.map(String) : undefined
    }
  })
}

async function handleDownloadTorrent(torrent: TorrentListItem) {
  if (downloadPendingId.value > 0) return

  downloadPendingId.value = torrent.id
  try {
    const out = await catalogTorrents.downloadTorrent(torrent.id)
    downloadBlob(out.blob, out.filename || fallbackTorrentFilename(torrent))
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

function categoryDisplayName(category: CatalogCategory) {
  return localizeI18nName(category.name, locale.value, category.slug || `#${category.id}`)
}

function categoryName(categoryId: number) {
  return categoryNameMap.value.get(categoryId) || t('catalog.torrents.unknownCategory')
}

function torrentOwnerName(torrent: TorrentListItem) {
  const ownerName = rawTorrentOwnerName(torrent)
  if (torrent.anonymous) {
    return torrent.ownerId > 0 ? t('catalog.torrents.anonymousOwner', { name: ownerName }) : t('catalog.torrents.anonymous')
  }
  return torrent.ownerId > 0 ? ownerName : '-'
}

function torrentOwnerPrimary(torrent: TorrentListItem) {
  if (torrent.anonymous) return t('catalog.torrents.anonymous')
  return torrent.ownerId > 0 ? rawTorrentOwnerName(torrent) : '-'
}

function torrentOwnerSecondary(torrent: TorrentListItem) {
  if (!torrent.anonymous || torrent.ownerId === 0) return ''
  return `(${rawTorrentOwnerName(torrent)})`
}

function rawTorrentOwnerName(torrent: TorrentListItem) {
  return torrent.ownerName || `#${torrent.ownerId}`
}

function torrentStatusBadges(torrent: TorrentListItem) {
  const badges = []
  if (torrent.isPinned) {
    badges.push({
      key: 'pinned',
      label: t('catalog.torrents.status.pinned'),
      class: 'bg-slate-900 text-white dark:bg-white dark:text-slate-950',
      title: t('catalog.torrents.status.pinned')
    })
  }
  if (torrent.isFeatured) {
    badges.push({
      key: 'featured',
      label: t('catalog.torrents.status.featured'),
      class: 'bg-amber-100 text-amber-800 dark:bg-amber-950 dark:text-amber-200',
      title: t('catalog.torrents.status.featured')
    })
  }

  const promotionLabel = torrentPromotionLabel(torrent.spState)
  if (promotionLabel && isPromotionActive(torrent)) {
    badges.push({
      key: `sp-${torrent.spState}`,
      label: promotionLabel,
      class: 'bg-emerald-100 text-emerald-800 dark:bg-emerald-950 dark:text-emerald-200',
      title: torrent.spExpireAt
        ? t('catalog.torrents.status.expiresAt', { time: formatDateTime(torrent.spExpireAt, locale.value) })
        : promotionLabel
    })
  }
  return badges
}

function torrentPromotionLabel(spState?: number | null) {
  const key = Number(spState || 0)
  return key >= 1 && key <= 7 ? t(`catalog.torrents.status.promotion.${key}`) : ''
}

function isPromotionActive(torrent: TorrentListItem) {
  if (!torrent.spState) return false
  if (!torrent.spExpireAt) return true

  const expireAt = parseDateTime(torrent.spExpireAt)
  return !expireAt || expireAt.getTime() > Date.now()
}

function relativeDateTime(value?: string | null) {
  const date = parseDateTime(value)
  if (!date) return '-'

  const diffSeconds = Math.round((date.getTime() - Date.now()) / 1000)
  const absSeconds = Math.abs(diffSeconds)

  if (absSeconds < 45) return relativeTimeFormatter.value.format(0, 'second')
  if (absSeconds < 45 * 60) return relativeTimeFormatter.value.format(Math.round(diffSeconds / 60), 'minute')
  if (absSeconds < 22 * 60 * 60) return relativeTimeFormatter.value.format(Math.round(diffSeconds / 60 / 60), 'hour')
  if (absSeconds < 26 * 60 * 60) return relativeTimeFormatter.value.format(Math.round(diffSeconds / 60 / 60 / 24), 'day')
  if (absSeconds < 30 * 24 * 60 * 60) return relativeTimeFormatter.value.format(Math.round(diffSeconds / 60 / 60 / 24), 'day')
  if (absSeconds < 12 * 30 * 24 * 60 * 60) return relativeTimeFormatter.value.format(Math.round(diffSeconds / 60 / 60 / 24 / 30), 'month')
  return relativeTimeFormatter.value.format(Math.round(diffSeconds / 60 / 60 / 24 / 365), 'year')
}

function parseDateTime(value?: string | null) {
  if (!value) return null

  const date = new Date(value)
  if (!Number.isNaN(date.getTime())) return date

  const normalizedDate = new Date(value.replace(' ', 'T'))
  return Number.isNaN(normalizedDate.getTime()) ? null : normalizedDate
}

function fallbackTorrentFilename(torrent: TorrentListItem) {
  const name = (torrent.name || `torrent-${torrent.id}`).replace(/[\\/:*?"<>|]+/g, '_').trim()
  return `${name || `torrent-${torrent.id}`}.torrent`
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

useSeoMeta({
  title: t('catalog.torrents.metaTitle'),
  robots: 'noindex, nofollow'
})
</script>
