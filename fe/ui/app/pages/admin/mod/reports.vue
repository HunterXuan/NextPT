<template>
  <div class="min-h-[calc(100vh-4rem)] bg-slate-50 py-6 dark:bg-slate-950">
    <div class="w-full px-3 sm:px-4 lg:px-5">
      <div class="mb-4 grid gap-3 lg:grid-cols-[minmax(0,1fr)_auto]">
        <form class="grid gap-2 sm:grid-cols-[160px_180px_auto]" @submit.prevent="reloadFromFirstPage">
          <select v-model.number="query.status" class="h-10 rounded-md border border-slate-200 bg-white px-3 text-sm outline-none dark:border-slate-700 dark:bg-slate-900" @change="reloadFromFirstPage">
            <option v-for="option in statusOptions" :key="option.value" :value="option.value">{{ option.label }}</option>
          </select>
          <select v-model="query.targetType" class="h-10 rounded-md border border-slate-200 bg-white px-3 text-sm outline-none dark:border-slate-700 dark:bg-slate-900" @change="reloadFromFirstPage">
            <option v-for="option in targetTypeOptions" :key="option.value" :value="option.value">{{ option.label }}</option>
          </select>
          <UButton type="submit" color="primary" icon="i-lucide-search" :loading="pending">
            {{ $t('admin.mod.actions.filter') }}
          </UButton>
        </form>

        <div class="rounded-lg border border-slate-200 bg-white px-3 py-2 dark:border-slate-800 dark:bg-slate-900">
          <p class="text-xs text-slate-500 dark:text-slate-400">{{ $t('admin.mod.stats.total') }}</p>
          <p class="mt-1 text-lg font-semibold text-slate-950 dark:text-white">{{ numberFormatter.format(total) }}</p>
        </div>
      </div>

      <div class="grid gap-4 xl:grid-cols-[minmax(0,1fr)_380px]">
        <section class="overflow-hidden rounded-lg border border-slate-200 bg-white dark:border-slate-800 dark:bg-slate-900">
          <div class="flex items-center justify-between gap-3 border-b border-slate-200 px-4 py-3 dark:border-slate-800">
            <h2 class="text-sm font-semibold text-slate-950 dark:text-white">{{ $t('admin.mod.reports.list') }}</h2>
            <UButton color="neutral" variant="outline" icon="i-lucide-refresh-cw" :loading="pending" @click="loadReports">
              {{ $t('common.refresh') }}
            </UButton>
          </div>

          <div v-if="pending" class="space-y-2 p-4">
            <div v-for="item in 6" :key="item" class="h-20 animate-pulse rounded-md bg-slate-100 dark:bg-slate-800" />
          </div>

          <div v-else-if="errorMessage" class="flex flex-col items-center justify-center px-4 py-16 text-center">
            <UIcon name="i-lucide-circle-alert" class="size-9 text-red-500" />
            <p class="mt-3 text-sm font-medium text-slate-950 dark:text-white">{{ errorMessage }}</p>
          </div>

          <div v-else-if="reports.length === 0" class="flex flex-col items-center justify-center px-4 py-16 text-center">
            <UIcon name="i-lucide-inbox" class="size-9 text-slate-400" />
            <p class="mt-3 text-sm font-medium text-slate-950 dark:text-white">{{ $t('admin.mod.reports.empty') }}</p>
          </div>

          <div v-else class="overflow-x-auto">
            <table class="min-w-[920px] w-full table-fixed border-collapse text-left">
              <thead class="bg-slate-50 text-xs font-medium uppercase text-slate-500 dark:bg-slate-950/70 dark:text-slate-400">
                <tr>
                  <th class="w-[10%] border-b border-slate-200 px-4 py-3 dark:border-slate-800">ID</th>
                  <th class="w-[20%] border-b border-slate-200 px-4 py-3 dark:border-slate-800">{{ $t('admin.mod.table.target') }}</th>
                  <th class="w-[26%] border-b border-slate-200 px-4 py-3 dark:border-slate-800">{{ $t('admin.mod.table.reason') }}</th>
                  <th class="w-[14%] border-b border-slate-200 px-4 py-3 dark:border-slate-800">{{ $t('admin.mod.table.status') }}</th>
                  <th class="w-[18%] border-b border-slate-200 px-4 py-3 dark:border-slate-800">{{ $t('admin.mod.table.createdAt') }}</th>
                  <th class="w-[12%] border-b border-slate-200 px-4 py-3 text-right dark:border-slate-800">{{ $t('admin.mod.table.actions') }}</th>
                </tr>
              </thead>
              <tbody>
                <tr
                  v-for="report in reports"
                  :key="report.id"
                  class="border-b border-slate-200 last:border-b-0 dark:border-slate-800"
                  :class="selectedReport?.id === report.id ? 'bg-rose-50/70 dark:bg-rose-950/20' : 'hover:bg-slate-50 dark:hover:bg-slate-950/70'"
                >
                  <td class="px-4 py-3 text-sm font-medium text-slate-950 dark:text-white">#{{ report.id }}</td>
                  <td class="px-4 py-3 text-sm text-slate-600 dark:text-slate-300">
                    <p class="font-medium text-slate-950 dark:text-white">{{ targetTypeLabel(report.target_type) }}</p>
                    <p class="mt-1 text-xs text-slate-500">#{{ report.target_id }}</p>
                  </td>
                  <td class="px-4 py-3 text-sm text-slate-600 dark:text-slate-300">
                    <p class="line-clamp-2">{{ report.reason }}</p>
                    <p class="mt-1 text-xs text-slate-500">{{ $t('admin.mod.reports.reporter', { id: report.reporter_id }) }}</p>
                  </td>
                  <td class="px-4 py-3">
                    <UBadge :color="reportStatusColor(report.status)" variant="soft">{{ reportStatusLabel(report.status) }}</UBadge>
                  </td>
                  <td class="px-4 py-3 text-sm text-slate-600 dark:text-slate-300">{{ formatDateTime(report.created_at, locale) }}</td>
                  <td class="px-4 py-3 text-right">
                    <UButton color="neutral" variant="outline" size="sm" icon="i-lucide-check-check" :disabled="report.status !== 0" @click="selectReport(report)">
                      {{ $t('admin.mod.actions.resolve') }}
                    </UButton>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>

          <AppPager
            class="border-t border-slate-200 px-4 py-3 dark:border-slate-800"
            size="sm"
            :page="query.page"
            :total="total"
            :page-size="query.size"
            :disabled="pending"
            @page-change="changePage"
          />
        </section>

        <section class="rounded-lg border border-slate-200 bg-white dark:border-slate-800 dark:bg-slate-900">
          <div class="border-b border-slate-200 px-4 py-3 dark:border-slate-800">
            <h2 class="text-sm font-semibold text-slate-950 dark:text-white">
              {{ selectedReport ? $t('admin.mod.reports.resolveTitle', { id: selectedReport.id }) : $t('admin.mod.reports.select') }}
            </h2>
          </div>
          <form class="space-y-4 p-4" @submit.prevent="resolveReport">
            <label class="block">
              <span class="text-sm font-medium text-slate-700 dark:text-slate-200">{{ $t('admin.mod.form.result') }}</span>
              <select v-model.number="resolveForm.status" class="mt-1 h-10 w-full rounded-md border border-slate-200 bg-white px-3 text-sm outline-none dark:border-slate-700 dark:bg-slate-950" :disabled="!selectedReport || resolving">
                <option :value="1">{{ $t('admin.mod.status.accepted') }}</option>
                <option :value="2">{{ $t('admin.mod.status.rejected') }}</option>
              </select>
            </label>
            <label class="block">
              <span class="text-sm font-medium text-slate-700 dark:text-slate-200">{{ $t('admin.mod.form.comment') }}</span>
              <textarea v-model="resolveForm.comment" rows="5" class="mt-1 w-full resize-none rounded-md border border-slate-200 bg-white px-3 py-2 text-sm outline-none dark:border-slate-700 dark:bg-slate-950" :disabled="!selectedReport || resolving" />
            </label>
            <p v-if="resolveError" class="rounded-md border border-red-200 bg-red-50 px-3 py-2 text-sm text-red-700 dark:border-red-900 dark:bg-red-950/40 dark:text-red-200">{{ resolveError }}</p>
            <UButton type="submit" color="primary" icon="i-lucide-check" :loading="resolving" :disabled="!selectedReport">
              {{ $t('admin.mod.actions.submitResolve') }}
            </UButton>
          </form>
        </section>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ApiError } from '~/composables/useApi'
import type { AdminModReportItem } from '~/composables/useAdmin'
import { formatDateTime } from '~/utils/format'

definePageMeta({ layout: 'admin', middleware: 'admin' })

const { t, locale } = useI18n()
const toast = useToast()
const adminApi = useAdmin()

const reports = ref<AdminModReportItem[]>([])
const total = ref(0)
const pending = ref(false)
const resolving = ref(false)
const errorMessage = ref('')
const resolveError = ref('')
const selectedReport = ref<AdminModReportItem | null>(null)

const query = reactive({
  page: 1,
  size: 30,
  status: -1,
  targetType: ''
})

const resolveForm = reactive({
  status: 1,
  comment: ''
})

const numberFormatter = computed(() => new Intl.NumberFormat(locale.value))
const totalPages = computed(() => Math.max(1, Math.ceil(total.value / query.size)))
const statusOptions = computed(() => [
  { value: -1, label: t('admin.mod.status.all') },
  { value: 0, label: t('admin.mod.status.pending') },
  { value: 1, label: t('admin.mod.status.accepted') },
  { value: 2, label: t('admin.mod.status.rejected') }
])
const targetTypeOptions = computed(() => [
  { value: '', label: t('admin.mod.targetTypes.all') },
  { value: 'torrent', label: t('admin.mod.targetTypes.torrent') },
  { value: 'comment', label: t('admin.mod.targetTypes.comment') },
  { value: 'subtitle', label: t('admin.mod.targetTypes.subtitle') },
  { value: 'forum_topic', label: t('admin.mod.targetTypes.forumTopic') },
  { value: 'forum_reply', label: t('admin.mod.targetTypes.forumReply') }
])

useHead({ title: t('admin.mod.reports.title') })
onMounted(loadReports)

async function loadReports() {
  pending.value = true
  errorMessage.value = ''
  try {
    const data = await adminApi.listModReports(query)
    reports.value = data.list || []
    total.value = data.total || 0
  } catch (error) {
    errorMessage.value = error instanceof ApiError ? error.message : t('common.requestFailed')
  } finally {
    pending.value = false
  }
}

function reloadFromFirstPage() {
  query.page = 1
  loadReports()
}

function changePage(page: number) {
  query.page = Math.min(Math.max(1, page), totalPages.value)
  loadReports()
}

function selectReport(report: AdminModReportItem) {
  selectedReport.value = report
  resolveForm.status = 1
  resolveForm.comment = ''
  resolveError.value = ''
}

async function resolveReport() {
  if (!selectedReport.value) return
  resolving.value = true
  resolveError.value = ''
  try {
    await adminApi.resolveModReport(selectedReport.value.id, {
      status: resolveForm.status,
      comment: resolveForm.comment.trim()
    })
    toast.add({ title: t('admin.mod.reports.resolved'), color: 'success', icon: 'i-lucide-check-circle' })
    selectedReport.value = null
    await loadReports()
  } catch (error) {
    resolveError.value = error instanceof ApiError ? error.message : t('common.requestFailed')
  } finally {
    resolving.value = false
  }
}

function reportStatusLabel(status: number) {
  if (status === 0) return t('admin.mod.status.pending')
  if (status === 1) return t('admin.mod.status.accepted')
  if (status === 2) return t('admin.mod.status.rejected')
  return String(status)
}

function reportStatusColor(status: number) {
  if (status === 0) return 'warning'
  if (status === 1) return 'success'
  if (status === 2) return 'neutral'
  return 'neutral'
}

function targetTypeLabel(value: string) {
  return targetTypeOptions.value.find((item) => item.value === value)?.label || value || '-'
}
</script>
