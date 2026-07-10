<template>
  <section class="overflow-hidden rounded-lg border border-slate-200 bg-white dark:border-slate-800 dark:bg-slate-900">
    <div class="flex flex-col gap-3 border-b border-slate-200 px-4 py-3 sm:flex-row sm:items-center sm:justify-between dark:border-slate-800">
      <div class="min-w-0">
        <h2 class="truncate text-sm font-semibold text-slate-950 dark:text-white">{{ titleText }}</h2>
        <p class="mt-0.5 truncate text-xs text-slate-500 dark:text-slate-400">{{ $t('user.loginLogs.subtitle') }}</p>
      </div>
      <div class="inline-flex h-8 shrink-0 items-center rounded-md bg-slate-100 p-0.5 dark:bg-slate-800">
        <button
          v-for="option in resultOptions"
          :key="option.value"
          type="button"
          class="h-7 rounded px-2.5 text-xs font-medium transition"
          :class="resultFilter === option.value
            ? 'bg-white text-slate-950 shadow-sm dark:bg-slate-950 dark:text-white'
            : 'text-slate-500 hover:text-slate-950 dark:text-slate-400 dark:hover:text-white'"
          @click="changeResultFilter(option.value)"
        >
          {{ option.label }}
        </button>
      </div>
    </div>

    <div v-if="pending" class="space-y-2 p-4">
      <div v-for="item in 4" :key="item" class="h-[72px] animate-pulse rounded-md bg-slate-100 dark:bg-slate-800" />
    </div>

    <div v-else-if="errorMessage" class="flex flex-col items-center justify-center px-4 py-10 text-center">
      <UIcon name="i-lucide-circle-alert" class="size-8 text-red-500" />
      <p class="mt-3 text-sm font-medium text-slate-950 dark:text-white">{{ errorMessage }}</p>
    </div>

    <div v-else-if="logs.length === 0" class="flex flex-col items-center justify-center px-4 py-10 text-center">
      <UIcon name="i-lucide-inbox" class="size-8 text-slate-400" />
      <p class="mt-3 text-sm font-medium text-slate-950 dark:text-white">{{ $t('user.loginLogs.empty') }}</p>
    </div>

    <div v-else class="divide-y divide-slate-200 dark:divide-slate-800">
      <article
        v-for="item in logs"
        :key="item.id"
        class="grid gap-3 px-4 py-3 sm:grid-cols-[minmax(0,1fr)_auto] sm:items-start"
      >
        <div class="flex min-w-0 items-start gap-3">
          <span
            class="mt-0.5 flex size-9 shrink-0 items-center justify-center rounded-md"
            :class="resultIconClass(item.result)"
          >
            <UIcon :name="resultIcon(item.result)" class="size-4" />
          </span>
          <div class="min-w-0 flex-1">
            <div class="flex flex-wrap items-center gap-2">
              <span class="text-sm font-semibold text-slate-950 dark:text-white">{{ loginResultLabel(item.result) }}</span>
              <span v-if="admin" class="rounded bg-slate-100 px-1.5 py-0.5 text-xs font-medium text-slate-500 dark:bg-slate-800 dark:text-slate-400">#{{ item.userId }}</span>
              <span class="rounded bg-slate-100 px-1.5 py-0.5 font-mono text-xs text-slate-600 dark:bg-slate-800 dark:text-slate-300">{{ item.ip || '-' }}</span>
              <span v-if="item.result !== 1" class="text-xs font-medium text-red-600 dark:text-red-300">
                {{ failReasonLabel(item.failReason) }}
              </span>
            </div>
            <UTooltip :text="item.userAgent || '-'" :content="{ side: 'top', sideOffset: 8 }" :delay-duration="600">
              <p class="mt-1.5 truncate text-xs leading-5 text-slate-500 dark:text-slate-400">{{ item.userAgent || '-' }}</p>
            </UTooltip>
          </div>
        </div>

        <UTooltip :text="formatDateTime(item.createdAt, locale)" :content="{ side: 'top', sideOffset: 8 }" :delay-duration="600">
          <time class="text-left text-xs text-slate-500 sm:pt-1 sm:text-right dark:text-slate-400">
            {{ formatDateTime(item.createdAt, locale) }}
          </time>
        </UTooltip>
      </article>
    </div>

    <div class="border-t border-slate-200 px-4 py-3 dark:border-slate-800">
      <AppPager
        size="sm"
        :page="query.page"
        :total="total"
        :page-size="query.size"
        :page-size-options="pageSizes"
        :disabled="pending"
        @page-change="changePage"
        @page-size-change="changePageSize"
      />
    </div>
  </section>
</template>

<script setup lang="ts">
import { ApiError } from '~/composables/useApi'
import { formatDateTime } from '~/utils/format'

type LoginLogItem = {
  id: number
  userId: number
  ip: string
  userAgent: string
  result: number
  failReason: string
  createdAt?: string | null
}

type ResultFilter = 'all' | 'success' | 'fail'

const props = withDefaults(defineProps<{
  admin?: boolean
  userId?: number | null
  title?: string
}>(), {
  admin: false,
  userId: null,
  title: ''
})

const { t, locale } = useI18n()
const authApi = useAuth()
const adminApi = useAdmin()

const logs = ref<LoginLogItem[]>([])
const total = ref(0)
const pending = ref(false)
const errorMessage = ref('')
const resultFilter = ref<ResultFilter>('all')
const query = reactive({
  page: 1,
  size: 10
})

const pageSizes = [10, 20, 50]
const titleText = computed(() => props.title || t('user.loginLogs.title'))
const resultOptions = computed(() => [
  { value: 'all' as const, label: t('user.loginLogs.filters.all') },
  { value: 'success' as const, label: t('user.loginLogs.filters.success') },
  { value: 'fail' as const, label: t('user.loginLogs.filters.fail') }
])

onMounted(loadLogs)

watch(
  () => [props.admin, props.userId],
  () => {
    query.page = 1
    loadLogs()
  }
)

function resultQueryValue() {
  if (resultFilter.value === 'success') return 1
  if (resultFilter.value === 'fail') return 0
  return undefined
}

function changeResultFilter(value: ResultFilter) {
  if (resultFilter.value === value) return
  resultFilter.value = value
  query.page = 1
  loadLogs()
}

function changePage(page: number) {
  query.page = page
  loadLogs()
}

function changePageSize(size: number) {
  query.size = size
  query.page = 1
  loadLogs()
}

async function loadLogs() {
  pending.value = true
  errorMessage.value = ''
  try {
    const result = resultQueryValue()
    const data = props.admin
      ? await adminApi.listIamLoginLogs({
          userId: props.userId || undefined,
          result,
          page: query.page,
          size: query.size
        })
      : await authApi.listLoginLogs({
          result,
          page: query.page,
          size: query.size
        })
    logs.value = data.list || []
    total.value = data.total || 0
  } catch (error: unknown) {
    logs.value = []
    total.value = 0
    errorMessage.value = error instanceof ApiError ? error.message : t('common.requestFailed')
  } finally {
    pending.value = false
  }
}

function loginResultLabel(result: number) {
  return result === 1 ? t('user.loginLogs.results.success') : t('user.loginLogs.results.fail')
}

function resultIcon(result: number) {
  return result === 1 ? 'i-lucide-check' : 'i-lucide-x'
}

function resultIconClass(result: number) {
  return result === 1
    ? 'bg-emerald-50 text-emerald-600 dark:bg-emerald-950/40 dark:text-emerald-300'
    : 'bg-red-50 text-red-600 dark:bg-red-950/40 dark:text-red-300'
}

function failReasonLabel(reason: string) {
  const labels: Record<string, string> = {
    user_not_found: t('user.loginLogs.failReasons.userNotFound'),
    invalid_password: t('user.loginLogs.failReasons.invalidPassword'),
    account_unavailable: t('user.loginLogs.failReasons.accountUnavailable'),
    role_missing: t('user.loginLogs.failReasons.roleMissing'),
    token_create_failed: t('user.loginLogs.failReasons.tokenCreateFailed')
  }
  return labels[reason] || reason || '-'
}
</script>
