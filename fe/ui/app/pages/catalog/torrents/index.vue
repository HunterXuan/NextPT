<template>
  <div class="min-h-[calc(100vh-4rem)] bg-slate-50 py-6 dark:bg-slate-950">
    <div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
      <div class="mb-6 flex flex-col gap-4 lg:flex-row lg:items-end lg:justify-between">
        <div>
          <p class="text-sm font-medium text-slate-500 dark:text-slate-400">{{ $t('catalog.torrents.eyebrow') }}</p>
          <h1 class="mt-1 text-2xl font-semibold text-slate-950 dark:text-white">{{ $t('catalog.torrents.title') }}</h1>
        </div>

        <div class="grid grid-cols-2 gap-3 sm:flex sm:items-center">
          <div class="rounded-lg border border-slate-200 bg-white px-4 py-3 dark:border-slate-800 dark:bg-slate-900">
            <p class="text-xs text-slate-500 dark:text-slate-400">{{ $t('catalog.torrents.summary.total') }}</p>
            <p class="mt-1 text-lg font-semibold text-slate-950 dark:text-white">{{ numberFormatter.format(total) }}</p>
          </div>
          <div class="rounded-lg border border-slate-200 bg-white px-4 py-3 dark:border-slate-800 dark:bg-slate-900">
            <p class="text-xs text-slate-500 dark:text-slate-400">{{ $t('catalog.torrents.summary.categories') }}</p>
            <p class="mt-1 text-lg font-semibold text-slate-950 dark:text-white">{{ numberFormatter.format(categories.length) }}</p>
          </div>
          <UButton class="col-span-2 sm:col-span-1" color="primary" icon="i-lucide-upload" :to="localePath('/catalog/torrents/upload')">
            {{ $t('catalog.torrents.upload.action') }}
          </UButton>
          <UButton class="col-span-2 sm:col-span-1" color="neutral" variant="outline" icon="i-lucide-captions" :to="localePath('/catalog/subtitles')">
            {{ $t('catalog.subtitles.title') }}
          </UButton>
        </div>
      </div>

      <div class="mb-4 rounded-lg border border-slate-200 bg-white p-3 dark:border-slate-800 dark:bg-slate-900">
        <form class="grid grid-cols-1 gap-3 md:grid-cols-[minmax(0,1fr)_auto_auto_auto]" @submit.prevent="handleSearchSubmit">
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
          <UButton
            type="button"
            color="neutral"
            variant="outline"
            icon="i-lucide-sliders-horizontal"
            :trailing-icon="advancedOpen ? 'i-lucide-chevron-up' : 'i-lucide-chevron-down'"
            @click="advancedOpen = !advancedOpen"
          >
            {{ $t('catalog.torrents.search.advanced') }}
          </UButton>
          <UButton
            type="button"
            color="neutral"
            variant="outline"
            icon="i-lucide-refresh-cw"
            :loading="pending"
            @click="loadTorrents"
          >
            {{ $t('common.refresh') }}
          </UButton>
        </form>

        <div v-if="advancedOpen" class="mt-4 border-t border-slate-200 pt-4 dark:border-slate-800">
          <div class="flex flex-col gap-2 sm:flex-row sm:items-center sm:justify-between">
            <p class="text-xs font-medium text-slate-500 dark:text-slate-400">{{ $t('catalog.torrents.filters.category') }}</p>
            <UButton
              v-if="selectedCategoryIds.length > 0"
              color="neutral"
              variant="ghost"
              size="xs"
              icon="i-lucide-x"
              @click="clearCategories"
            >
              {{ $t('catalog.torrents.filters.clearCategories') }}
            </UButton>
          </div>

          <div class="mt-3 grid grid-cols-2 gap-2 sm:grid-cols-3 lg:grid-cols-5">
            <label
              v-for="category in categories"
              :key="category.id"
              class="flex h-10 cursor-pointer items-center gap-2 rounded-md border px-3 text-sm transition"
              :class="selectedCategoryIds.includes(category.id)
                ? 'border-sky-300 bg-sky-50 text-sky-800 dark:border-sky-700 dark:bg-sky-950 dark:text-sky-200'
                : 'border-slate-200 bg-white text-slate-700 hover:border-slate-300 dark:border-slate-700 dark:bg-slate-950 dark:text-slate-200 dark:hover:border-slate-600'"
            >
              <input
                type="checkbox"
                class="size-4 rounded border-slate-300 text-sky-600 focus:ring-sky-500 dark:border-slate-600"
                :checked="selectedCategoryIds.includes(category.id)"
                :disabled="pending"
                @change="toggleCategory(category.id)"
              >
              <span class="truncate">{{ categoryDisplayName(category) }}</span>
            </label>
          </div>

          <p v-if="categories.length === 0" class="mt-3 text-sm text-slate-500 dark:text-slate-400">
            {{ $t('catalog.torrents.filters.noCategories') }}
          </p>
        </div>
      </div>

      <div class="overflow-hidden rounded-lg border border-slate-200 bg-white dark:border-slate-800 dark:bg-slate-900">
        <div class="hidden grid-cols-[minmax(0,1fr)_128px_150px_112px_150px] gap-4 border-b border-slate-200 px-4 py-3 text-xs font-medium uppercase text-slate-500 md:grid dark:border-slate-800 dark:text-slate-400">
          <span>{{ $t('catalog.torrents.table.torrent') }}</span>
          <span>{{ $t('catalog.torrents.table.size') }}</span>
          <span>{{ $t('catalog.torrents.table.activity') }}</span>
          <span>{{ $t('catalog.torrents.table.completed') }}</span>
          <span>{{ $t('catalog.torrents.table.createdAt') }}</span>
        </div>

        <div v-if="pending" class="divide-y divide-slate-200 dark:divide-slate-800">
          <div v-for="index in 6" :key="index" class="grid gap-4 px-4 py-4 md:grid-cols-[minmax(0,1fr)_128px_150px_112px_150px]">
            <div class="space-y-2">
              <div class="h-4 w-3/4 animate-pulse rounded bg-slate-200 dark:bg-slate-800" />
              <div class="h-3 w-1/2 animate-pulse rounded bg-slate-100 dark:bg-slate-800/70" />
            </div>
            <div class="h-4 w-20 animate-pulse rounded bg-slate-100 dark:bg-slate-800/70" />
            <div class="h-4 w-24 animate-pulse rounded bg-slate-100 dark:bg-slate-800/70" />
            <div class="h-4 w-12 animate-pulse rounded bg-slate-100 dark:bg-slate-800/70" />
            <div class="h-4 w-28 animate-pulse rounded bg-slate-100 dark:bg-slate-800/70" />
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
            class="grid gap-4 px-4 py-4 transition-colors hover:bg-slate-50 md:grid-cols-[minmax(0,1fr)_128px_150px_112px_150px] md:items-center dark:hover:bg-slate-950/70"
          >
            <div class="min-w-0">
              <div class="flex flex-wrap items-center gap-2">
                <UBadge color="neutral" variant="soft">{{ categoryName(torrent.categoryId) }}</UBadge>
                <UBadge :color="torrent.type === 1 ? 'primary' : 'neutral'" variant="subtle">
                  {{ torrent.type === 1 ? $t('catalog.torrents.types.multi') : $t('catalog.torrents.types.single') }}
                </UBadge>
              </div>
              <h2 class="mt-2 truncate text-sm font-semibold">
                <NuxtLink
                  :to="localePath(`/catalog/torrents/${torrent.id}`)"
                  class="text-slate-950 hover:text-sky-700 dark:text-white dark:hover:text-sky-300"
                >
                  {{ torrent.name || `#${torrent.id}` }}
                </NuxtLink>
              </h2>
              <p v-if="torrent.subTitle" class="mt-1 truncate text-sm text-slate-500 dark:text-slate-400">
                {{ torrent.subTitle }}
              </p>
              <p class="mt-2 text-xs text-slate-500 dark:text-slate-400">
                {{ torrentOwnerName(torrent) }}
              </p>
            </div>

            <div>
              <p class="text-sm font-medium text-slate-950 dark:text-white">{{ formatBytes(torrent.size) }}</p>
              <p class="mt-1 text-xs text-slate-500 dark:text-slate-400">
                {{ $t('catalog.torrents.fileCount', { count: numberFormatter.format(torrent.fileCount) }) }}
              </p>
            </div>

            <div class="flex items-center gap-3 text-sm">
              <span class="inline-flex items-center gap-1 text-emerald-600 dark:text-emerald-400">
                <UIcon name="i-lucide-arrow-up" class="size-4" />
                {{ numberFormatter.format(torrent.seeders) }}
              </span>
              <span class="inline-flex items-center gap-1 text-sky-600 dark:text-sky-400">
                <UIcon name="i-lucide-arrow-down" class="size-4" />
                {{ numberFormatter.format(torrent.leechers) }}
              </span>
            </div>

            <p class="text-sm font-medium text-slate-950 dark:text-white">
              {{ numberFormatter.format(torrent.snatched) }}
            </p>

            <p class="text-sm text-slate-500 dark:text-slate-400">
              {{ formatDateTime(torrent.createdAt, locale) }}
            </p>
          </article>
        </div>
      </div>

      <div class="mt-4 flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
        <p class="text-sm text-slate-500 dark:text-slate-400">
          {{ $t('catalog.torrents.pagination.summary', { page: page, pages: totalPages }) }}
        </p>
        <div class="flex flex-wrap items-center gap-2">
          <label class="flex items-center gap-2 text-sm text-slate-500 dark:text-slate-400">
            <span>{{ $t('catalog.torrents.pagination.pageSize') }}</span>
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
          <UButton
            color="neutral"
            variant="outline"
            icon="i-lucide-chevron-left"
            :disabled="page <= 1 || pending"
            @click="goToPage(page - 1)"
          >
            {{ $t('common.previous') }}
          </UButton>
          <UButton
            color="neutral"
            variant="outline"
            trailing-icon="i-lucide-chevron-right"
            :disabled="page >= totalPages || pending"
            @click="goToPage(page + 1)"
          >
            {{ $t('common.next') }}
          </UButton>
        </div>
      </div>
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
const catalogTorrents = useCatalogTorrents()

const categories = ref<CatalogCategory[]>([])
const torrents = ref<TorrentListItem[]>([])
const total = ref(0)
const pending = ref(false)
const errorMessage = ref('')

const page = ref(readPositiveIntQuery('page', 1))
const keyword = ref(readStringQuery('keyword'))
const selectedCategoryIds = ref(readCategoryIdsQuery())
const selectedSize = ref(String(readPageSizeQuery()))
const advancedOpen = ref(selectedCategoryIds.value.length > 0)

const numberFormatter = computed(() => new Intl.NumberFormat(locale.value))
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
  const rawValues = route.query.categoryIds
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

function handlePageSizeChange() {
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
      ...route.query,
      page: page.value > 1 ? String(page.value) : undefined,
      size: selectedSize.value !== '50' ? selectedSize.value : undefined,
      keyword: keyword.value.trim() || undefined,
      categoryIds: selectedCategoryIds.value.length ? selectedCategoryIds.value.join(',') : undefined
    }
  })
}

function categoryDisplayName(category: CatalogCategory) {
  return localizeI18nName(category.name, locale.value, category.slug || `#${category.id}`)
}

function categoryName(categoryId: number) {
  return categoryNameMap.value.get(categoryId) || t('catalog.torrents.unknownCategory')
}

function torrentOwnerName(torrent: TorrentListItem) {
  const ownerName = torrent.ownerName || `#${torrent.ownerId}`
  if (torrent.anonymous) {
    return torrent.ownerId > 0 ? t('catalog.torrents.anonymousOwner', { name: ownerName }) : t('catalog.torrents.anonymous')
  }
  return torrent.ownerId > 0 ? ownerName : '-'
}

useSeoMeta({
  title: t('catalog.torrents.metaTitle'),
  robots: 'noindex, nofollow'
})
</script>
