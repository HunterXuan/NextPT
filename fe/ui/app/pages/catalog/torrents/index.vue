<template>
  <div class="min-h-[calc(100vh-4rem)] bg-slate-50 py-6 dark:bg-slate-950">
    <div class="w-full px-3 sm:px-4 lg:px-5">
      <section class="mb-3 rounded-lg border border-slate-200 bg-white p-3 dark:border-slate-800 dark:bg-slate-900">
        <div class="flex items-center gap-2">
          <form class="flex min-w-0 flex-1 items-center gap-2" @submit.prevent="handleSearchSubmit">
            <UInput
              v-model="keyword"
              class="min-w-0 flex-1"
              :ui="{ base: 'h-10' }"
              icon="i-lucide-search"
              :placeholder="$t('catalog.torrents.search.placeholder')"
              :disabled="pending"
            />
            <UButton
              type="submit"
              color="primary"
              icon="i-lucide-search"
              :loading="pending"
              class="h-10 w-10 shrink-0 justify-center p-0 sm:w-auto sm:px-3"
              :aria-label="$t('catalog.torrents.search.submit')"
            >
              <span class="hidden sm:inline">{{ $t('catalog.torrents.search.submit') }}</span>
            </UButton>
          </form>

          <div class="flex shrink-0 items-center gap-2">
            <UTooltip :text="$t('catalog.torrents.search.advanced')" :content="{ side: 'bottom', sideOffset: 8 }" :delay-duration="600">
              <span class="relative inline-flex">
                <UButton
                  type="button"
                  color="neutral"
                  :variant="advancedOpen || activeAdvancedFilterCount > 0 ? 'soft' : 'outline'"
                  icon="i-lucide-list-filter"
                  class="h-10 w-10 justify-center p-0"
                  :aria-label="$t('catalog.torrents.search.advanced')"
                  :aria-expanded="advancedOpen"
                  @click="advancedOpen = !advancedOpen"
                />
                <span
                  v-if="activeAdvancedFilterCount > 0"
                  class="pointer-events-none absolute -right-1 -top-1 flex size-4 items-center justify-center rounded-full bg-sky-600 text-[10px] font-semibold text-white ring-2 ring-white dark:ring-slate-900"
                >
                  {{ activeAdvancedFilterCount }}
                </span>
              </span>
            </UTooltip>
            <CatalogTorrentRssPopover
              :keyword="appliedKeyword"
              :category-ids="selectedCategoryIds"
              :tag-ids="appliedAdvancedFilters.tagIds"
              :promotion="appliedAdvancedFilters.promotion"
              :seed-status="appliedAdvancedFilters.seedStatus"
              :featured-only="appliedAdvancedFilters.featuredOnly"
              :min-size="appliedAdvancedFilters.minSize"
              :max-size="appliedAdvancedFilters.maxSize"
              :published-within="appliedAdvancedFilters.publishedWithin"
              :imdb-id="appliedAdvancedFilters.imdbId"
              :douban-id="appliedAdvancedFilters.doubanId"
              :bangumi-id="appliedAdvancedFilters.bangumiId"
              :tmdb-id="appliedAdvancedFilters.tmdbId"
              :tmdb-type="appliedAdvancedFilters.tmdbType"
            />
            <AppPermissionButton
              :permission="Permission.CatalogTorrentCreate"
              color="primary"
              variant="soft"
              icon="i-lucide-upload"
              :to="localePath('/catalog/torrents/upload')"
              class="h-10 w-10 shrink-0 justify-center p-0"
              :aria-label="$t('catalog.torrents.upload.action')"
              :tooltip="$t('catalog.torrents.upload.action')"
            />
          </div>
        </div>

        <div class="mt-3 -mx-1 overflow-x-auto px-1 pb-1">
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

        <div v-if="advancedOpen" class="mt-3 border-t border-slate-200 pt-3 dark:border-slate-800">
          <CatalogTorrentAdvancedFilters
            :value="appliedAdvancedFilters"
            :tag-groups="tagGroups"
            :category-ids="selectedCategoryIds"
            :pending="pending"
            @apply="applyAdvancedFilters"
          />
        </div>
      </section>

      <SiteAdvertisement placement="catalog_list" class="mb-3" />

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
                  <div v-if="torrent.tags?.length" class="mt-1 flex min-w-0 flex-wrap gap-1">
                    <span
                      v-for="tag in torrent.tags.slice(0, 3)"
                      :key="tag.id"
                      class="inline-flex h-5 max-w-28 items-center truncate rounded border border-slate-200 bg-slate-50 px-1.5 text-[10px] font-medium text-slate-500 dark:border-slate-700 dark:bg-slate-950 dark:text-slate-400"
                    >
                      {{ tagName(tag) }}
                    </span>
                    <span v-if="torrent.tags.length > 3" class="inline-flex h-5 items-center text-[10px] text-slate-400">+{{ torrent.tags.length - 3 }}</span>
                  </div>
                </div>
                <AppPermissionButton
                  :permission="Permission.CatalogTorrentDownload"
                  class="lg:hidden"
                  color="neutral"
                  variant="ghost"
                  size="xs"
                  icon="i-lucide-download"
                  :aria-label="$t('catalog.torrents.detail.actions.download')"
                  :tooltip="$t('catalog.torrents.detail.actions.download')"
                  :loading="downloadPendingId === torrent.id"
                  :disabled="downloadPendingId > 0"
                  @click="handleDownloadTorrent(torrent)"
                />
              </div>
              <div class="mt-1 flex flex-wrap items-center gap-x-3 gap-y-1 text-xs text-slate-500 dark:text-slate-400">
                <IamUserPopover
                  v-if="!torrent.anonymous && torrent.owner?.id"
                  :user="torrent.owner"
                  :fallback="torrentOwnerName(torrent)"
                  class="lg:hidden"
                />
                <span v-else class="lg:hidden">{{ torrentOwnerName(torrent) }}</span>
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
              <IamUserPopover
                v-if="!torrent.anonymous && torrent.owner?.id"
                :user="torrent.owner"
                :fallback="torrentOwnerPrimary(torrent)"
                class="truncate text-sm leading-5 text-slate-600 dark:text-slate-300"
              />
              <p v-else class="truncate text-sm leading-5 text-slate-600 dark:text-slate-300">{{ torrentOwnerPrimary(torrent) }}</p>
              <p v-if="torrentOwnerSecondary(torrent)" class="truncate text-xs leading-4 text-slate-400 dark:text-slate-500">
                {{ torrentOwnerSecondary(torrent) }}
              </p>
            </div>
            <div class="hidden justify-center lg:flex">
              <AppPermissionButton
                :permission="Permission.CatalogTorrentDownload"
                class="text-slate-400 hover:text-slate-700 dark:hover:text-slate-200"
                color="neutral"
                variant="ghost"
                size="xs"
                icon="i-lucide-download"
                :aria-label="$t('catalog.torrents.detail.actions.download')"
                :tooltip="$t('catalog.torrents.detail.actions.download')"
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
import type { CatalogCategory, CatalogTagGroup, CatalogTagItem, TorrentAdvancedFilters, TorrentListItem } from '~/composables/useCatalogTorrents'

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

const categories = ref<CatalogCategory[]>([])
const tagGroups = ref<CatalogTagGroup[]>([])
const torrents = ref<TorrentListItem[]>([])
const total = ref(0)
const pending = ref(false)
const errorMessage = ref('')
const downloadPendingId = ref(0)
const pageSizes = [20, 50, 100]

const page = ref(readPositiveIntQuery('page', 1))
const keyword = ref(readStringQuery('keyword'))
const appliedKeyword = ref(keyword.value)
const selectedCategoryIds = ref(readCategoryIdsQuery())
const selectedSize = ref(String(readPageSizeQuery()))
const appliedAdvancedFilters = reactive<TorrentAdvancedFilters>(readAdvancedFilters())
const advancedOpen = ref(advancedFilterCount(appliedAdvancedFilters) > 0)

const numberFormatter = computed(() => new Intl.NumberFormat(locale.value))
const relativeTimeFormatter = computed(() => new Intl.RelativeTimeFormat(locale.value, { numeric: 'auto' }))
const totalPages = computed(() => Math.max(1, Math.ceil(total.value / Number(selectedSize.value || 50))))
const canDownloadTorrent = computed(() => hasPermission(Permission.CatalogTorrentDownload))
const activeAdvancedFilterCount = computed(() => advancedFilterCount(appliedAdvancedFilters))

const categoryNameMap = computed(() => {
  const map = new Map<number, string>()
  for (const category of categories.value) {
    map.set(category.id, categoryDisplayName(category))
  }
  return map
})

onMounted(async () => {
  await Promise.all([loadCategories(), loadTagGroups(), loadTorrents()])
})

function readFirstQueryValue(key: string) {
  const value = route.query[key]
  return Array.isArray(value) ? value[0] : value
}

function readStringQuery(key: string) {
  return String(readFirstQueryValue(key) || '').trim()
}

function readBooleanQuery(key: string) {
  const value = readStringQuery(key).toLowerCase()
  return value === 'true' || value === '1'
}

function readAdvancedFilters(): TorrentAdvancedFilters {
  const promotion = readStringQuery('promotion')
  const seedStatus = readStringQuery('seedStatus')
  const publishedWithin = Number(readFirstQueryValue('publishedWithin'))
  const sort = readStringQuery('sort')
  return {
    promotion: ['promoted', 'normal', 'free', '2x', '2x_free', '50_percent', '2x_50_percent', '30_percent'].includes(promotion) ? promotion : 'all',
    seedStatus: ['seeded', 'unseeded'].includes(seedStatus) ? seedStatus : 'all',
    featuredOnly: readBooleanQuery('featuredOnly'),
    minSize: readNonnegativeIntegerQuery('minSize'),
    maxSize: readNonnegativeIntegerQuery('maxSize'),
    publishedWithin: [1, 7, 30, 90].includes(publishedWithin) ? publishedWithin : 0,
    sort: ['oldest', 'seeders', 'leechers', 'completed', 'size_asc', 'size_desc'].includes(sort) ? sort : 'newest',
    imdbId: readStringQuery('imdbId').toLowerCase(),
    doubanId: readStringQuery('doubanId'),
    bangumiId: readStringQuery('bangumiId'),
    tmdbId: readStringQuery('tmdbId'),
    tmdbType: ['movie', 'tv'].includes(readStringQuery('tmdbType')) ? readStringQuery('tmdbType') : 'all',
    tagIds: readIdListQuery('tagIds')
  }
}

function readNonnegativeIntegerQuery(key: string) {
  const value = Number(readFirstQueryValue(key))
  return Number.isSafeInteger(value) && value > 0 ? value : 0
}

function advancedFilterCount(value: TorrentAdvancedFilters) {
  return [
    value.promotion !== 'all',
    value.seedStatus !== 'all',
    value.featuredOnly,
    value.minSize > 0,
    value.maxSize > 0,
    value.publishedWithin > 0,
    value.sort !== 'newest',
    Boolean(value.imdbId),
    Boolean(value.doubanId),
    Boolean(value.bangumiId),
    Boolean(value.tmdbId),
    value.tagIds.length > 0
  ].filter(Boolean).length
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
  return readIdListQuery('categoryIds')
}

function readIdListQuery(key: string) {
  const rawValues = route.query[key] || route.query[`${key}[]`]
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

async function loadTagGroups() {
  try {
    const data = await catalogTorrents.listTagGroups()
    tagGroups.value = data.list || []
  } catch {
    tagGroups.value = []
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
      keyword: appliedKeyword.value,
      categoryIds: selectedCategoryIds.value,
      tagIds: appliedAdvancedFilters.tagIds,
      promotion: appliedAdvancedFilters.promotion,
      seedStatus: appliedAdvancedFilters.seedStatus,
      featuredOnly: appliedAdvancedFilters.featuredOnly,
      minSize: appliedAdvancedFilters.minSize,
      maxSize: appliedAdvancedFilters.maxSize,
      publishedWithin: appliedAdvancedFilters.publishedWithin,
      sort: appliedAdvancedFilters.sort,
      imdbId: appliedAdvancedFilters.imdbId,
      doubanId: appliedAdvancedFilters.doubanId,
      bangumiId: appliedAdvancedFilters.bangumiId,
      tmdbId: appliedAdvancedFilters.tmdbId,
      tmdbType: appliedAdvancedFilters.tmdbType
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
  appliedKeyword.value = keyword.value.trim()
  page.value = 1
  loadTorrents()
}

function applyAdvancedFilters(value: TorrentAdvancedFilters) {
  Object.assign(appliedAdvancedFilters, value)
  page.value = 1
  loadTorrents()
}

function toggleCategory(categoryId: number) {
  selectedCategoryIds.value = selectedCategoryIds.value.includes(categoryId)
    ? selectedCategoryIds.value.filter((id) => id !== categoryId)
    : [...selectedCategoryIds.value, categoryId]
  pruneUnavailableTagFilters()
  page.value = 1
  loadTorrents()
}

function clearCategories() {
  selectedCategoryIds.value = []
  page.value = 1
  loadTorrents()
}

function pruneUnavailableTagFilters() {
  if (!selectedCategoryIds.value.length || !tagGroups.value.length) return
  const available = new Set(tagGroups.value
    .filter(group => !group.categories?.length || group.categories.some(id => selectedCategoryIds.value.includes(id)))
    .flatMap(group => group.tags.map(tag => tag.id)))
  appliedAdvancedFilters.tagIds = appliedAdvancedFilters.tagIds.filter(id => available.has(id))
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
      keyword: appliedKeyword.value || undefined,
      categoryIds: selectedCategoryIds.value.length ? selectedCategoryIds.value.map(String) : undefined,
      promotion: appliedAdvancedFilters.promotion !== 'all' ? appliedAdvancedFilters.promotion : undefined,
      seedStatus: appliedAdvancedFilters.seedStatus !== 'all' ? appliedAdvancedFilters.seedStatus : undefined,
      featuredOnly: appliedAdvancedFilters.featuredOnly ? 'true' : undefined,
      minSize: appliedAdvancedFilters.minSize > 0 ? String(appliedAdvancedFilters.minSize) : undefined,
      maxSize: appliedAdvancedFilters.maxSize > 0 ? String(appliedAdvancedFilters.maxSize) : undefined,
      publishedWithin: appliedAdvancedFilters.publishedWithin > 0 ? String(appliedAdvancedFilters.publishedWithin) : undefined,
      sort: appliedAdvancedFilters.sort !== 'newest' ? appliedAdvancedFilters.sort : undefined,
      imdbId: appliedAdvancedFilters.imdbId || undefined,
      doubanId: appliedAdvancedFilters.doubanId || undefined,
      bangumiId: appliedAdvancedFilters.bangumiId || undefined,
      tmdbId: appliedAdvancedFilters.tmdbId || undefined,
      tmdbType: appliedAdvancedFilters.tmdbId && appliedAdvancedFilters.tmdbType !== 'all' ? appliedAdvancedFilters.tmdbType : undefined,
      tagIds: appliedAdvancedFilters.tagIds.length ? appliedAdvancedFilters.tagIds.map(String) : undefined
    }
  })
}

async function handleDownloadTorrent(torrent: TorrentListItem) {
  if (!canDownloadTorrent.value || downloadPendingId.value > 0) return

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

function tagName(tag: CatalogTagItem) {
  return localizeI18nName(tag.name, locale.value, tag.value)
}

function torrentOwnerName(torrent: TorrentListItem) {
  const ownerName = rawTorrentOwnerName(torrent)
  if (torrent.anonymous) {
    return torrent.owner?.id > 0 ? t('catalog.torrents.anonymousOwner', { name: ownerName }) : t('catalog.torrents.anonymous')
  }
  return torrent.owner?.id > 0 ? ownerName : '-'
}

function torrentOwnerPrimary(torrent: TorrentListItem) {
  if (torrent.anonymous) return t('catalog.torrents.anonymous')
  return torrent.owner?.id > 0 ? rawTorrentOwnerName(torrent) : '-'
}

function torrentOwnerSecondary(torrent: TorrentListItem) {
  if (!torrent.anonymous || !torrent.owner?.id) return ''
  return `(${rawTorrentOwnerName(torrent)})`
}

function rawTorrentOwnerName(torrent: TorrentListItem) {
  return torrent.owner?.username || (torrent.owner?.id > 0 ? `#${torrent.owner.id}` : '')
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
  return key >= 1 && key <= 6 ? t(`catalog.torrents.status.promotion.${key}`) : ''
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
  document.body.removeChild(link)
  URL.revokeObjectURL(href)
}
</script>
