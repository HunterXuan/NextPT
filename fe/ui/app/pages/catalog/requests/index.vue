<template>
  <div class="min-h-[calc(100vh-4rem)] bg-slate-50 py-6 dark:bg-slate-950">
    <div class="w-full px-3 sm:px-4 lg:px-5">
      <section class="rounded-lg border border-slate-200 bg-white dark:border-slate-800 dark:bg-slate-900">
        <div class="flex flex-col gap-3 border-b border-slate-200 p-3 sm:flex-row sm:items-center dark:border-slate-800">
          <form class="grid min-w-0 flex-1 grid-cols-[minmax(0,1fr)_40px] items-center gap-2 sm:grid-cols-[minmax(0,1fr)_auto]" @submit.prevent="handleSearch">
            <UInput
              v-model="keyword"
              class="min-w-0 flex-1"
              :ui="{ base: 'h-10' }"
              icon="i-lucide-search"
              :placeholder="$t('catalog.requests.search.placeholder')"
              :disabled="pending"
            />
            <UTooltip :text="$t('catalog.requests.search.submit')" :delay-duration="600">
              <UButton type="submit" color="primary" icon="i-lucide-search" class="h-10 w-10 justify-center sm:w-auto" :loading="pending" :aria-label="$t('catalog.requests.search.submit')">
                <span class="hidden sm:inline">{{ $t('catalog.requests.search.submit') }}</span>
              </UButton>
            </UTooltip>
          </form>
          <AppPermissionButton
            :permission="Permission.CatalogRequestCreate"
            color="primary"
            variant="soft"
            icon="i-lucide-plus"
            class="h-10 shrink-0"
            :to="localePath('/catalog/requests/create')"
            :tooltip="$t('catalog.requests.create.action')"
          >
            {{ $t('catalog.requests.create.action') }}
          </AppPermissionButton>
        </div>

        <div class="flex flex-col gap-3 border-b border-slate-200 p-3 lg:flex-row lg:items-center lg:justify-between dark:border-slate-800">
          <div class="flex min-w-0 gap-2 overflow-x-auto pb-1 lg:pb-0">
            <button v-for="item in viewOptions" :key="item.value" type="button" :class="filterButtonClass(view === item.value)" @click="setView(item.value)">
              {{ item.label }}
            </button>
          </div>
          <div class="grid grid-cols-2 gap-2 sm:flex sm:items-center">
            <select v-model.number="requestType" class="h-9 min-w-0 w-full rounded-md border border-slate-200 bg-white px-2.5 text-sm text-slate-700 outline-none dark:border-slate-700 dark:bg-slate-950 dark:text-slate-200" @change="applyFilters">
              <option :value="0">{{ $t('catalog.requests.filters.allTypes') }}</option>
              <option :value="CatalogRequestType.Torrent">{{ $t('catalog.requests.types.torrent') }}</option>
              <option :value="CatalogRequestType.Reseed">{{ $t('catalog.requests.types.reseed') }}</option>
            </select>
            <select v-model="statusValue" class="h-9 min-w-0 w-full rounded-md border border-slate-200 bg-white px-2.5 text-sm text-slate-700 outline-none dark:border-slate-700 dark:bg-slate-950 dark:text-slate-200" @change="applyFilters">
              <option value="all">{{ $t('catalog.requests.filters.allStatuses') }}</option>
              <option v-for="item in statusOptions" :key="item.value" :value="String(item.value)">{{ item.label }}</option>
            </select>
          </div>
        </div>

        <div v-if="errorMessage" class="border-b border-red-200 bg-red-50 px-4 py-3 text-sm text-red-700 dark:border-red-900 dark:bg-red-950/40 dark:text-red-200">
          {{ errorMessage }}
        </div>

        <div v-if="pending && requests.length === 0" class="divide-y divide-slate-100 dark:divide-slate-800">
          <div v-for="index in 6" :key="index" class="h-20 animate-pulse bg-slate-50/70 dark:bg-slate-950/40" />
        </div>

        <div v-else-if="requests.length === 0" class="px-4 py-16 text-center">
          <UIcon name="i-lucide-inbox" class="mx-auto size-9 text-slate-300 dark:text-slate-600" />
          <h2 class="mt-3 text-sm font-semibold text-slate-950 dark:text-white">{{ $t('catalog.requests.empty.title') }}</h2>
          <p class="mt-1 text-sm text-slate-500 dark:text-slate-400">{{ $t('catalog.requests.empty.description') }}</p>
        </div>

        <div v-else class="divide-y divide-slate-100 dark:divide-slate-800">
          <article v-for="item in requests" :key="item.id" class="grid gap-3 px-4 py-3 transition-colors hover:bg-slate-50/70 lg:grid-cols-[minmax(0,1fr)_110px_130px_170px] lg:items-center dark:hover:bg-slate-950/50">
            <div class="min-w-0">
              <div class="flex min-w-0 items-center gap-2">
                <UBadge :color="item.requestType === CatalogRequestType.Reseed ? 'warning' : 'primary'" variant="soft" size="sm">
                  {{ requestTypeLabel(item.requestType) }}
                </UBadge>
                <span v-if="categoryName(item.categoryId)" class="shrink-0 text-xs text-slate-500 dark:text-slate-400">{{ categoryName(item.categoryId) }}</span>
              </div>
              <NuxtLink :to="localePath(`/catalog/requests/${item.id}`)" class="mt-1.5 block truncate text-sm font-semibold text-slate-950 outline-none hover:text-sky-600 focus-visible:text-sky-600 dark:text-white dark:hover:text-sky-400">
                {{ item.title }}
              </NuxtLink>
              <div class="mt-1 flex flex-wrap items-center gap-x-2 text-xs text-slate-500 dark:text-slate-400">
                <span>{{ item.requester.username || `#${item.requester.id}` }}</span>
                <span>/</span>
                <UTooltip :text="formatDateTime(item.updatedAt, locale)" :delay-duration="600">
                  <span>{{ formatRelativeDateTime(item.updatedAt, locale) }}</span>
                </UTooltip>
              </div>
            </div>

            <div class="flex items-center justify-between lg:block lg:text-right">
              <span class="text-xs text-slate-500 lg:hidden dark:text-slate-400">{{ $t('catalog.requests.fields.reward') }}</span>
              <span class="inline-flex items-center gap-1 text-sm font-semibold text-amber-600 dark:text-amber-400">
                <UIcon name="i-lucide-coins" class="size-4" />
                {{ numberFormatter.format(item.rewardAmount) }}
              </span>
            </div>

            <div class="flex items-center justify-between lg:justify-end">
              <span class="text-xs text-slate-500 lg:hidden dark:text-slate-400">{{ $t('catalog.requests.fields.status') }}</span>
              <UBadge :color="requestStatusColor(item.status)" variant="soft" size="sm">{{ requestStatusLabel(item.status) }}</UBadge>
            </div>

            <div class="flex min-w-0 items-center justify-between gap-2 lg:justify-end">
              <span class="text-xs text-slate-500 lg:hidden dark:text-slate-400">{{ $t('catalog.requests.fields.claimer') }}</span>
              <div v-if="item.claimer" class="flex min-w-0 items-center gap-2">
                <IamUserAvatar :user="item.claimer" size="xs" />
                <span class="truncate text-sm text-slate-700 dark:text-slate-200">{{ item.claimer.username || `#${item.claimer.id}` }}</span>
              </div>
              <span v-else class="text-sm text-slate-400">-</span>
            </div>
          </article>
        </div>
      </section>

      <AppPager
        class="mt-4"
        :page="page"
        :total="total"
        :page-size="pageSize"
        :page-size-options="pageSizeOptions"
        :disabled="pending"
        @page-change="goToPage"
        @page-size-change="changePageSize"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ApiError } from '~/composables/useApi'
import { CatalogRequestStatus, CatalogRequestType, type CatalogRequestListItem, type CatalogRequestView } from '~/composables/useCatalogRequests'
import type { CatalogCategory } from '~/composables/useCatalogTorrents'
import { formatDateTime, formatRelativeDateTime, localizeI18nName } from '~/utils/format'

definePageMeta({ middleware: 'auth' })

const { t, locale } = useI18n()
const localePath = useLocalePath()
const route = useRoute()
const router = useRouter()
const requestApi = useCatalogRequests()
const catalogApi = useCatalogTorrents()

const requests = ref<CatalogRequestListItem[]>([])
const categories = ref<CatalogCategory[]>([])
const total = ref(0)
const pending = ref(false)
const errorMessage = ref('')
const keyword = ref(readStringQuery('keyword'))
const appliedKeyword = ref(keyword.value)
const requestType = ref(readNumberQuery('requestType'))
const statusValue = ref(readStringQuery('status') || 'all')
const view = ref<CatalogRequestView>(readViewQuery())
const page = ref(Math.max(1, readNumberQuery('page') || 1))
const pageSize = ref([20, 50, 100].includes(readNumberQuery('size')) ? readNumberQuery('size') : 20)
const pageSizeOptions = [20, 50, 100]
const numberFormatter = computed(() => new Intl.NumberFormat(locale.value, { maximumFractionDigits: 1 }))

const viewOptions = computed(() => [
  { value: 'all' as const, label: t('catalog.requests.views.all') },
  { value: 'created' as const, label: t('catalog.requests.views.created') },
  { value: 'claimed' as const, label: t('catalog.requests.views.claimed') }
])
const statusOptions = computed(() => Object.values(CatalogRequestStatus).map(value => ({ value, label: requestStatusLabel(value) })))

useHead(() => ({ title: t('catalog.requests.metaTitle') }))

onMounted(async () => {
  await Promise.all([loadCategories(), loadRequests()])
})

async function loadCategories() {
  try {
    categories.value = (await catalogApi.listCategories()).list || []
  } catch {
    categories.value = []
  }
}

async function loadRequests() {
  pending.value = true
  errorMessage.value = ''
  appliedKeyword.value = keyword.value.trim()
  syncQuery()
  try {
    const data = await requestApi.listRequests({
      page: page.value,
      size: pageSize.value,
      keyword: appliedKeyword.value,
      requestType: requestType.value || undefined,
      status: statusValue.value === 'all' ? undefined : Number(statusValue.value),
      view: view.value
    })
    requests.value = data.list || []
    total.value = data.total || 0
  } catch (error) {
    requests.value = []
    total.value = 0
    errorMessage.value = error instanceof ApiError ? error.message : t('common.requestFailed')
  } finally {
    pending.value = false
  }
}

function handleSearch() {
  page.value = 1
  loadRequests()
}

function applyFilters() {
  page.value = 1
  loadRequests()
}

function setView(value: CatalogRequestView) {
  if (view.value === value) return
  view.value = value
  applyFilters()
}

function goToPage(value: number) {
  page.value = value
  loadRequests()
}

function changePageSize(value: number) {
  pageSize.value = value
  page.value = 1
  loadRequests()
}

function syncQuery() {
  const query: Record<string, string> = {}
  if (page.value > 1) query.page = String(page.value)
  if (pageSize.value !== 20) query.size = String(pageSize.value)
  if (appliedKeyword.value) query.keyword = appliedKeyword.value
  if (requestType.value) query.requestType = String(requestType.value)
  if (statusValue.value !== 'all') query.status = statusValue.value
  if (view.value !== 'all') query.view = view.value
  router.replace({ query })
}

function readStringQuery(key: string) {
  const value = route.query[key]
  return String(Array.isArray(value) ? value[0] || '' : value || '')
}

function readNumberQuery(key: string) {
  const value = Number(readStringQuery(key))
  return Number.isFinite(value) && value >= 0 ? value : 0
}

function readViewQuery(): CatalogRequestView {
  const value = readStringQuery('view')
  return value === 'created' || value === 'claimed' ? value : 'all'
}

function categoryName(id: number) {
  const category = categories.value.find(item => item.id === id)
  return category ? localizeI18nName(category.name, locale.value) : ''
}

function requestTypeLabel(value: number) {
  return value === CatalogRequestType.Reseed ? t('catalog.requests.types.reseed') : t('catalog.requests.types.torrent')
}

function requestStatusLabel(value: number) {
  return t(`catalog.requests.status.${value}`)
}

function requestStatusColor(value: number): 'neutral' | 'primary' | 'warning' | 'success' | 'error' {
  if (value === CatalogRequestStatus.Open) return 'primary'
  if (value === CatalogRequestStatus.Claimed || value === CatalogRequestStatus.Submitted) return 'warning'
  if (value === CatalogRequestStatus.Completed) return 'success'
  if (value === CatalogRequestStatus.Cancelled) return 'error'
  return 'neutral'
}

function filterButtonClass(active: boolean) {
  return [
    'h-8 shrink-0 rounded-md border px-3 text-sm font-medium transition-colors',
    active
      ? 'border-slate-950 bg-slate-950 text-white dark:border-white dark:bg-white dark:text-slate-950'
      : 'border-slate-200 text-slate-600 hover:bg-slate-50 dark:border-slate-700 dark:text-slate-300 dark:hover:bg-slate-950'
  ]
}
</script>
