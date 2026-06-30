<template>
  <div class="min-h-[calc(100vh-4rem)] bg-slate-50 py-6 dark:bg-slate-950">
    <div class="w-full px-3 sm:px-4 lg:px-5">
      <div class="mb-4 flex justify-end">
        <UButton color="neutral" variant="outline" icon="i-lucide-refresh-cw" :loading="pending" @click="loadBookmarks">
          {{ $t('common.refresh') }}
        </UButton>
      </div>

      <div class="overflow-hidden rounded-lg border border-slate-200 bg-white dark:border-slate-800 dark:bg-slate-900">
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
          <UButton class="mt-5" color="neutral" variant="outline" icon="i-lucide-refresh-cw" @click="loadBookmarks">
            {{ $t('common.retry') }}
          </UButton>
        </div>

        <div v-else-if="torrents.length === 0" class="flex flex-col items-center justify-center px-4 py-16 text-center">
          <UIcon name="i-lucide-bookmark-x" class="size-9 text-slate-400" />
          <p class="mt-3 text-sm font-medium text-slate-950 dark:text-white">{{ $t('catalog.bookmarks.empty.title') }}</p>
          <p class="mt-1 max-w-md text-sm text-slate-500 dark:text-slate-400">{{ $t('catalog.bookmarks.empty.description') }}</p>
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
                <UBadge v-if="torrent.anonymous" color="neutral" variant="outline">
                  {{ torrentOwnerName(torrent) }}
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

            <p class="text-sm font-medium text-slate-950 dark:text-white">{{ numberFormatter.format(torrent.snatched) }}</p>
            <p class="text-sm text-slate-500 dark:text-slate-400">{{ formatDateTime(torrent.createdAt, locale) }}</p>
          </article>
        </div>
      </div>

      <AppPager
        class="mt-4"
        :page="page"
        :total="total"
        :page-size="size"
        :disabled="pending"
        @page-change="goToPage"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ApiError } from '~/composables/useApi'
import type { CatalogCategory, TorrentListItem } from '~/composables/useCatalogTorrents'
import { formatBytes, formatDateTime, localizeI18nName } from '~/utils/format'

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
const page = ref(readPositiveIntQuery('page', 1))
const size = 20
const pending = ref(false)
const errorMessage = ref('')
const numberFormatter = computed(() => new Intl.NumberFormat(locale.value))
const totalPages = computed(() => Math.max(1, Math.ceil(total.value / size)))

const categoryNameMap = computed(() => {
  const map = new Map<number, string>()
  for (const category of categories.value) {
    map.set(category.id, localizeI18nName(category.name, locale.value, category.slug || `#${category.id}`))
  }
  return map
})

useHead(() => ({
  title: `${t('catalog.bookmarks.metaTitle')} - NextPT`
}))

onMounted(async () => {
  await Promise.all([loadCategories(), loadBookmarks()])
})

function readFirstQueryValue(key: string) {
  const value = route.query[key]
  return Array.isArray(value) ? value[0] : value
}

function readPositiveIntQuery(key: string, fallback: number) {
  const parsed = Number(readFirstQueryValue(key))
  return Number.isInteger(parsed) && parsed > 0 ? parsed : fallback
}

async function loadCategories() {
  try {
    const data = await catalogTorrents.listCategories()
    categories.value = data.list || []
  } catch {
    categories.value = []
  }
}

async function loadBookmarks() {
  pending.value = true
  errorMessage.value = ''
  syncQuery()

  try {
    let data = await catalogTorrents.listBookmarks(page.value, size)
    const nextTotal = data.total || 0
    const nextTotalPages = Math.max(1, Math.ceil(nextTotal / size))
    if (page.value > nextTotalPages) {
      page.value = nextTotalPages
      data = await catalogTorrents.listBookmarks(page.value, size)
    }
    torrents.value = data.list || []
    total.value = data.total || 0
  } catch (error) {
    torrents.value = []
    total.value = 0
    errorMessage.value = error instanceof ApiError ? error.message : t('common.requestFailed')
  } finally {
    pending.value = false
  }
}

function goToPage(nextPage: number) {
  page.value = Math.min(Math.max(1, nextPage), totalPages.value)
  loadBookmarks()
}

function syncQuery() {
  router.replace({
    query: {
      ...route.query,
      page: page.value > 1 ? String(page.value) : undefined
    }
  })
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
</script>
