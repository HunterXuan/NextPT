<template>
  <div class="min-h-[calc(100vh-4rem)] bg-slate-50 py-6 dark:bg-slate-950">
    <div class="w-full px-3 sm:px-4 lg:px-5">
      <section class="rounded-lg border border-slate-200 bg-white dark:border-slate-800 dark:bg-slate-900">
        <div class="flex flex-col gap-3 border-b border-slate-200 p-3 xl:flex-row xl:items-center dark:border-slate-800">
          <div class="flex min-w-0 gap-2 overflow-x-auto pb-1 xl:pb-0">
            <button v-for="item in viewOptions" :key="item.value" type="button" :class="filterButtonClass(view === item.value)" :disabled="pending" @click="setView(item.value)">
              {{ item.label }}
            </button>
          </div>

          <div class="flex min-w-0 flex-col gap-2 sm:flex-row sm:items-center xl:ml-auto">
            <div class="grid grid-cols-2 gap-2 sm:flex sm:items-center">
              <USelect
                :model-value="requestType"
                class="w-full sm:w-32"
                size="lg"
                :ui="{ base: 'h-9 w-full' }"
                :items="requestTypeOptions"
                value-key="value"
                :disabled="pending"
                @update:model-value="changeRequestType"
              />
              <USelect
                :model-value="statusValue"
                class="w-full sm:w-32"
                size="lg"
                :ui="{ base: 'h-9 w-full' }"
                :items="statusSelectOptions"
                value-key="value"
                :disabled="pending"
                @update:model-value="changeStatus"
              />
            </div>

            <div class="flex min-w-0 gap-2 sm:flex-1 xl:flex-none">
              <form class="grid min-w-0 flex-1 grid-cols-[minmax(0,1fr)_36px] items-center gap-2 xl:w-[420px] 2xl:w-[520px]" @submit.prevent="handleSearch">
                <UInput
                  v-model="keyword"
                  class="min-w-0"
                  :ui="{ base: 'h-9' }"
                  :placeholder="$t('catalog.requests.search.placeholder')"
                  :disabled="pending"
                />
                <UTooltip :text="$t('catalog.requests.search.submit')" :delay-duration="600">
                  <UButton type="submit" color="primary" icon="i-lucide-search" class="h-9 w-9 justify-center" :loading="pending" :aria-label="$t('catalog.requests.search.submit')" />
                </UTooltip>
              </form>

              <AppPermissionButton
                :permission="Permission.CatalogRequestCreate"
                color="primary"
                variant="soft"
                icon="i-lucide-plus"
                class="h-9 w-9 shrink-0 justify-center"
                :to="localePath('/catalog/requests/create')"
                :tooltip="$t('catalog.requests.create.action')"
                :aria-label="$t('catalog.requests.create.action')"
              />
            </div>
          </div>
        </div>

        <div v-if="errorMessage" class="border-b border-red-200 bg-red-50 px-4 py-3 text-sm text-red-700 dark:border-red-900 dark:bg-red-950/40 dark:text-red-200">
          {{ errorMessage }}
        </div>

        <div v-if="pending && requests.length === 0" class="divide-y divide-slate-100 dark:divide-slate-800">
          <div v-for="index in 6" :key="index" class="h-16 animate-pulse bg-slate-50/70 dark:bg-slate-950/40" />
        </div>

        <div v-else-if="requests.length === 0" class="px-4 py-16 text-center">
          <UIcon name="i-lucide-inbox" class="mx-auto size-9 text-slate-300 dark:text-slate-600" />
          <h2 class="mt-3 text-sm font-semibold text-slate-950 dark:text-white">{{ $t('catalog.requests.empty.title') }}</h2>
          <p class="mt-1 text-sm text-slate-500 dark:text-slate-400">{{ $t('catalog.requests.empty.description') }}</p>
        </div>

        <div v-else>
          <div class="hidden grid-cols-[minmax(280px,1fr)_120px_130px_140px_180px_170px] items-center gap-4 border-b border-slate-200 bg-slate-50/70 px-4 py-2 text-xs font-medium text-slate-500 xl:grid dark:border-slate-800 dark:bg-slate-950/40 dark:text-slate-400">
            <span>{{ $t('catalog.requests.fields.title') }}</span>
            <span class="text-center">{{ $t('catalog.requests.fields.status') }}</span>
            <span class="text-center">{{ $t('catalog.requests.fields.reward') }}</span>
            <span>{{ $t('catalog.requests.fields.publishedAt') }}</span>
            <span>{{ $t('catalog.requests.fields.publisher') }}</span>
            <span class="text-right">{{ $t('catalog.requests.fields.claimer') }}</span>
          </div>

          <div class="divide-y divide-slate-100 dark:divide-slate-800">
            <NuxtLink
              v-for="item in requests"
              :key="item.id"
              :to="localePath(`/catalog/requests/${item.id}`)"
              class="group grid gap-2 px-4 py-2.5 outline-none transition-colors hover:bg-slate-50/70 focus-visible:bg-sky-50/70 xl:grid-cols-[minmax(280px,1fr)_120px_130px_140px_180px_170px] xl:items-center xl:gap-4 dark:hover:bg-slate-950/50 dark:focus-visible:bg-sky-950/30"
            >
              <div class="min-w-0">
                <div class="flex min-w-0 items-center gap-2">
                  <UBadge :color="item.requestType === CatalogRequestType.Reseed ? 'warning' : 'primary'" variant="soft" size="sm">
                    {{ requestTypeLabel(item.requestType) }}
                  </UBadge>
                  <span class="min-w-0 flex-1 truncate text-sm font-semibold text-slate-950 transition-colors group-hover:text-sky-600 dark:text-white dark:group-hover:text-sky-400">
                    {{ item.title }}
                  </span>
                </div>
              </div>

              <div class="grid grid-cols-2 gap-x-4 gap-y-2 border-t border-slate-100 pt-2 xl:contents dark:border-slate-800">
                <div class="flex items-center justify-between xl:justify-center">
                  <span class="text-xs text-slate-500 xl:hidden dark:text-slate-400">{{ $t('catalog.requests.fields.status') }}</span>
                  <UBadge :color="requestStatusColor(item.status)" variant="soft" size="sm">{{ requestStatusLabel(item.status) }}</UBadge>
                </div>

                <div class="flex items-center justify-between xl:flex xl:justify-center">
                  <span class="text-xs text-slate-500 xl:hidden dark:text-slate-400">{{ $t('catalog.requests.fields.reward') }}</span>
                  <span class="inline-flex items-center gap-1 text-sm font-semibold text-amber-600 dark:text-amber-400">
                    <UIcon name="i-lucide-coins" class="size-4" />
                    {{ numberFormatter.format(item.rewardAmount) }}
                  </span>
                </div>

                <div class="flex items-center justify-between gap-2 xl:justify-start">
                  <span class="text-xs text-slate-500 xl:hidden dark:text-slate-400">{{ $t('catalog.requests.fields.publishedAt') }}</span>
                  <UTooltip :text="formatDateTime(item.createdAt, locale)" :delay-duration="600">
                    <span class="text-sm text-slate-600 dark:text-slate-300">{{ formatRelativeDateTime(item.createdAt, locale) }}</span>
                  </UTooltip>
                </div>

                <div class="flex min-w-0 items-center justify-between gap-2 xl:justify-start">
                  <span class="text-xs text-slate-500 xl:hidden dark:text-slate-400">{{ $t('catalog.requests.fields.publisher') }}</span>
                  <div class="flex min-w-0 items-center gap-2">
                    <IamUserAvatar :user="item.requester" size="xs" />
                    <span class="truncate text-sm text-slate-700 dark:text-slate-200">{{ item.requester.username || `#${item.requester.id}` }}</span>
                  </div>
                </div>

                <div class="col-span-2 flex min-w-0 items-center justify-between gap-2 xl:col-span-1 xl:justify-end">
                  <span class="text-xs text-slate-500 xl:hidden dark:text-slate-400">{{ $t('catalog.requests.fields.claimer') }}</span>
                  <div v-if="item.claimer" class="flex min-w-0 items-center gap-2">
                    <IamUserAvatar :user="item.claimer" size="xs" />
                    <span class="truncate text-sm text-slate-700 dark:text-slate-200">{{ item.claimer.username || `#${item.claimer.id}` }}</span>
                  </div>
                  <span v-else class="text-sm text-slate-400">{{ $t('catalog.requests.list.unclaimed') }}</span>
                </div>
              </div>
            </NuxtLink>
          </div>
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
import { formatDateTime, formatRelativeDateTime } from '~/utils/format'

definePageMeta({ middleware: 'auth' })

const { t, locale } = useI18n()
const localePath = useLocalePath()
const route = useRoute()
const router = useRouter()
const requestApi = useCatalogRequests()

const requests = ref<CatalogRequestListItem[]>([])
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
const requestTypeOptions = computed(() => [
  { value: 0, label: t('catalog.requests.filters.allTypes') },
  { value: CatalogRequestType.Torrent, label: t('catalog.requests.types.torrent') },
  { value: CatalogRequestType.Reseed, label: t('catalog.requests.types.reseed') }
])
const statusSelectOptions = computed(() => [
  { value: 'all', label: t('catalog.requests.filters.allStatuses') },
  ...statusOptions.value.map(item => ({ value: String(item.value), label: item.label }))
])

useHead(() => ({ title: t('catalog.requests.metaTitle') }))

onMounted(async () => {
  await loadRequests()
})

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

function changeRequestType(value: unknown) {
  const nextValue = Number(value)
  requestType.value = Number.isFinite(nextValue) ? nextValue : 0
  applyFilters()
}

function changeStatus(value: unknown) {
  statusValue.value = value == null ? 'all' : String(value)
  applyFilters()
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
    'h-8 shrink-0 rounded-md border px-3 text-sm font-medium transition-colors disabled:cursor-not-allowed disabled:opacity-60',
    active
      ? 'border-slate-950 bg-slate-950 text-white dark:border-white dark:bg-white dark:text-slate-950'
      : 'border-slate-200 text-slate-600 hover:bg-slate-50 dark:border-slate-700 dark:text-slate-300 dark:hover:bg-slate-950'
  ]
}
</script>
