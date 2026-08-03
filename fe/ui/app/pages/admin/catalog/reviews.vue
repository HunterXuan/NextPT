<template>
  <div class="min-h-[calc(100vh-4rem)] bg-slate-50 py-6 dark:bg-slate-950">
    <div class="w-full px-3 sm:px-4 lg:px-5">
      <section class="overflow-hidden rounded-lg border border-slate-200 bg-white dark:border-slate-800 dark:bg-slate-900">
        <div class="border-b border-slate-200 px-4 py-3 dark:border-slate-800">
          <div class="flex flex-col gap-3 xl:flex-row xl:items-center xl:justify-between">
            <h2 class="text-sm font-semibold text-slate-950 dark:text-white">{{ $t('admin.catalog.reviews.list') }}</h2>
            <form class="grid min-w-0 gap-2 sm:grid-cols-[minmax(220px,1fr)_150px_180px]" @submit.prevent="reloadFromFirstPage">
              <div class="flex min-w-0">
                <UInput v-model="query.keyword" class="min-w-0 flex-1" size="lg" :ui="{ base: 'h-9 rounded-r-none' }" :placeholder="$t('admin.catalog.reviews.filters.keyword')" />
                <UButton type="submit" class="rounded-l-none" color="neutral" variant="soft" icon="i-lucide-search" :aria-label="$t('common.search')" />
              </div>
              <USelect v-model="query.status" size="lg" :ui="{ base: 'h-9 w-full' }" :items="statusOptions" value-key="value" @update:model-value="reloadFromFirstPage" />
              <USelect v-model="query.categoryId" size="lg" :ui="{ base: 'h-9 w-full' }" :items="categoryOptions" value-key="value" @update:model-value="reloadFromFirstPage" />
            </form>
          </div>
        </div>

        <div v-if="pending" class="space-y-2 p-4">
          <div v-for="item in 7" :key="item" class="h-[72px] animate-pulse rounded-md bg-slate-100 dark:bg-slate-800" />
        </div>
        <div v-else-if="errorMessage" class="flex flex-col items-center justify-center px-4 py-16 text-center">
          <UIcon name="i-lucide-circle-alert" class="size-9 text-red-500" />
          <p class="mt-3 text-sm font-medium text-slate-950 dark:text-white">{{ errorMessage }}</p>
        </div>
        <div v-else-if="reviews.length === 0" class="flex flex-col items-center justify-center px-4 py-16 text-center">
          <UIcon name="i-lucide-clipboard-check" class="size-9 text-slate-400" />
          <p class="mt-3 text-sm font-medium text-slate-950 dark:text-white">{{ $t('admin.catalog.reviews.empty') }}</p>
        </div>
        <div v-else class="overflow-x-auto">
          <table class="w-full min-w-[760px] table-fixed border-collapse text-left">
            <thead class="bg-slate-50 text-xs font-medium uppercase text-slate-500 dark:bg-slate-950/70 dark:text-slate-400">
              <tr>
                <th class="w-[42%] border-b border-slate-200 px-4 py-3 dark:border-slate-800">{{ $t('admin.catalog.reviews.table.torrent') }}</th>
                <th class="w-[15%] border-b border-slate-200 px-4 py-3 dark:border-slate-800">{{ $t('admin.catalog.reviews.table.status') }}</th>
                <th class="w-[20%] border-b border-slate-200 px-4 py-3 dark:border-slate-800">{{ $t('admin.catalog.reviews.table.uploader') }}</th>
                <th class="w-[23%] border-b border-slate-200 px-4 py-3 text-right dark:border-slate-800">{{ $t('admin.catalog.reviews.table.submittedAt') }}</th>
              </tr>
            </thead>
            <tbody>
              <tr
                v-for="review in reviews"
                :key="review.id"
                class="group cursor-pointer border-b border-slate-200 outline-none transition last:border-b-0 hover:bg-slate-50 focus-visible:bg-slate-50 dark:border-slate-800 dark:hover:bg-slate-950/70 dark:focus-visible:bg-slate-950/70"
                tabindex="0"
                role="link"
                @click="openReview(review)"
                @keydown.enter.prevent="openReview(review)"
              >
                <td class="px-4 py-3">
                  <p class="truncate text-sm font-semibold text-slate-950 transition-colors group-hover:text-sky-700 dark:text-white dark:group-hover:text-sky-300">{{ review.name }}</p>
                  <div class="mt-1 flex min-w-0 items-center gap-2 text-xs text-slate-500 dark:text-slate-400">
                    <span class="truncate">{{ categoryName(review.categoryId) }}</span>
                    <span class="text-slate-300 dark:text-slate-700">/</span>
                    <span class="shrink-0">{{ formatBytes(review.size) }}</span>
                  </div>
                </td>
                <td class="px-4 py-3"><UBadge :color="statusColor(review.status)" variant="soft">{{ statusLabel(review.status) }}</UBadge></td>
                <td class="px-4 py-3">
                  <IamUserPopover :user="review.owner" :fallback="ownerName(review)" show-avatar avatar-size="xs" class="truncate text-sm text-slate-700 dark:text-slate-200" />
                </td>
                <td class="px-4 py-3 text-right text-sm text-slate-600 dark:text-slate-300">
                  <UTooltip :text="formatDateTime(review.submittedAt || review.createdAt, locale)" :delay-duration="600">
                    <span>{{ relativeTime(review.submittedAt || review.createdAt) }}</span>
                  </UTooltip>
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <AppPager class="border-t border-slate-200 px-4 py-3 dark:border-slate-800" size="sm" :page="query.page" :total="total" :page-size="query.size" :disabled="pending" @page-change="changePage" />
      </section>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ApiError } from '~/composables/useApi'
import type { AdminCatalogTorrentReviewItem } from '~/composables/useAdmin'
import { TorrentStatus, type CatalogCategory } from '~/composables/useCatalogTorrents'
import { formatBytes, formatDateTime, localizeI18nName } from '~/utils/format'

definePageMeta({ layout: 'admin', middleware: 'admin' })

const { t, locale } = useI18n()
const localePath = useLocalePath()
const adminApi = useAdmin()
const catalogApi = useCatalogTorrents()

const reviews = ref<AdminCatalogTorrentReviewItem[]>([])
const categories = ref<CatalogCategory[]>([])
const total = ref(0)
const pending = ref(false)
const errorMessage = ref('')
const query = reactive({ page: 1, size: 30, keyword: '', status: TorrentStatus.Pending, categoryId: 0 })

const totalPages = computed(() => Math.max(1, Math.ceil(total.value / query.size)))
const statusOptions = computed(() => [
  { value: -1, label: t('admin.catalog.reviews.status.all') },
  { value: TorrentStatus.Pending, label: t('admin.catalog.reviews.status.pending') },
  { value: TorrentStatus.Published, label: t('admin.catalog.reviews.status.published') },
  { value: TorrentStatus.Rejected, label: t('admin.catalog.reviews.status.rejected') }
])
const categoryOptions = computed(() => [
  { value: 0, label: t('admin.catalog.reviews.filters.allCategories') },
  ...categories.value.map(category => ({ value: category.id, label: localizeI18nName(category.name, locale.value, category.slug) }))
])
const categoryMap = computed(() => new Map(categories.value.map(category => [category.id, category])))
const relativeTimeFormatter = computed(() => new Intl.RelativeTimeFormat(locale.value, { numeric: 'auto' }))

useHead({ title: t('admin.catalog.reviews.title') })
onMounted(async () => {
  await Promise.all([loadCategories(), loadReviews()])
})

async function loadCategories() {
  try {
    const data = await catalogApi.listCategories()
    categories.value = data.list || []
  } catch {
    categories.value = []
  }
}

async function loadReviews() {
  pending.value = true
  errorMessage.value = ''
  try {
    const data = await adminApi.listCatalogTorrentReviews({
      ...query,
      keyword: query.keyword.trim() || undefined,
      categoryId: query.categoryId || undefined
    })
    reviews.value = data.list || []
    total.value = data.total || 0
  } catch (error) {
    errorMessage.value = error instanceof ApiError ? error.message : t('common.requestFailed')
  } finally {
    pending.value = false
  }
}

function reloadFromFirstPage() {
  query.page = 1
  void loadReviews()
}

function changePage(page: number) {
  query.page = Math.min(Math.max(1, page), totalPages.value)
  void loadReviews()
}

function openReview(review: AdminCatalogTorrentReviewItem) {
  void navigateTo(localePath(`/catalog/torrents/${review.id}`))
}

function statusLabel(status: number) {
  if (status === TorrentStatus.Pending) return t('admin.catalog.reviews.status.pending')
  if (status === TorrentStatus.Published) return t('admin.catalog.reviews.status.published')
  if (status === TorrentStatus.Rejected) return t('admin.catalog.reviews.status.rejected')
  return String(status)
}

function statusColor(status: number) {
  if (status === TorrentStatus.Pending) return 'warning'
  if (status === TorrentStatus.Published) return 'success'
  return 'error'
}

function categoryName(id: number) {
  const category = categoryMap.value.get(id)
  return category ? localizeI18nName(category.name, locale.value, category.slug) : `#${id}`
}

function ownerName(review: AdminCatalogTorrentReviewItem) {
  return review.owner?.username || `#${review.owner?.id || 0}`
}

function relativeTime(value?: string | null) {
  if (!value) return '-'
  const diffSeconds = Math.round((new Date(value).getTime() - Date.now()) / 1000)
  const ranges: Array<[Intl.RelativeTimeFormatUnit, number]> = [['year', 31536000], ['month', 2592000], ['day', 86400], ['hour', 3600], ['minute', 60]]
  for (const [unit, seconds] of ranges) {
    if (Math.abs(diffSeconds) >= seconds) return relativeTimeFormatter.value.format(Math.round(diffSeconds / seconds), unit)
  }
  return relativeTimeFormatter.value.format(diffSeconds, 'second')
}
</script>
