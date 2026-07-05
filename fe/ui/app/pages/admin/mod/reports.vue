<template>
  <div class="min-h-[calc(100vh-4rem)] bg-slate-50 py-6 dark:bg-slate-950">
    <div class="w-full px-3 sm:px-4 lg:px-5">
      <div class="grid gap-4 xl:grid-cols-[minmax(0,1fr)_480px] xl:items-start">
        <section class="overflow-hidden rounded-lg border border-slate-200 bg-white dark:border-slate-800 dark:bg-slate-900">
          <div class="border-b border-slate-200 px-4 py-3 dark:border-slate-800">
            <div class="flex flex-col gap-3 lg:flex-row lg:items-center lg:justify-between">
              <div>
                <h2 class="text-sm font-semibold text-slate-950 dark:text-white">{{ $t('admin.mod.reports.list') }}</h2>
              </div>
              <form class="grid gap-2 sm:grid-cols-[150px_180px]" @submit.prevent="reloadFromFirstPage">
                <select v-model.number="query.status" class="h-9 rounded-md border border-slate-200 bg-white px-3 text-sm outline-none transition focus:border-sky-300 dark:border-slate-700 dark:bg-slate-950 dark:focus:border-sky-700" @change="reloadFromFirstPage">
                  <option v-for="option in statusOptions" :key="option.value" :value="option.value">{{ option.label }}</option>
                </select>
                <select v-model="query.targetType" class="h-9 rounded-md border border-slate-200 bg-white px-3 text-sm outline-none transition focus:border-sky-300 dark:border-slate-700 dark:bg-slate-950 dark:focus:border-sky-700" @change="reloadFromFirstPage">
                  <option v-for="option in targetTypeOptions" :key="option.value" :value="option.value">{{ option.label }}</option>
                </select>
              </form>
            </div>
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
            <table class="min-w-[760px] w-full table-fixed border-collapse text-left">
              <thead class="bg-slate-50 text-xs font-medium uppercase text-slate-500 dark:bg-slate-950/70 dark:text-slate-400">
                <tr>
                  <th class="w-[30%] border-b border-slate-200 px-4 py-3 dark:border-slate-800">{{ $t('admin.mod.table.target') }}</th>
                  <th class="w-[36%] border-b border-slate-200 px-4 py-3 dark:border-slate-800">{{ $t('admin.mod.reports.reason') }}</th>
                  <th class="w-[12%] border-b border-slate-200 px-4 py-3 dark:border-slate-800">{{ $t('admin.mod.table.status') }}</th>
                  <th class="w-[22%] border-b border-slate-200 px-4 py-3 text-right dark:border-slate-800">{{ $t('admin.mod.table.createdAt') }}</th>
                </tr>
              </thead>
              <tbody>
                <tr
                  v-for="report in reports"
                  :key="report.id"
                  class="cursor-pointer border-b border-slate-200 outline-none transition last:border-b-0 focus-visible:bg-slate-50 dark:border-slate-800 dark:focus-visible:bg-slate-950/70"
                  :class="selectedReport?.id === report.id ? 'bg-sky-50/80 dark:bg-sky-950/25' : 'hover:bg-slate-50 dark:hover:bg-slate-950/70'"
                  tabindex="0"
                  @click="selectReport(report)"
                  @keydown.enter.prevent="selectReport(report)"
                >
                  <td class="px-4 py-3 text-sm text-slate-600 dark:text-slate-300">
                    <div class="flex min-w-0 items-center gap-3">
                      <span class="flex size-9 shrink-0 items-center justify-center rounded-md bg-slate-100 text-slate-500 dark:bg-slate-800 dark:text-slate-300">
                        <UIcon :name="targetTypeIcon(reportTargetType(report))" class="size-4" />
                      </span>
                      <div class="min-w-0">
                        <div class="flex min-w-0 items-center gap-2">
                          <p class="truncate font-semibold text-slate-950 dark:text-white">{{ reportTargetTitle(report) }}</p>
                          <UBadge v-if="isReportTargetDeleted(report)" color="neutral" variant="soft" size="sm" class="shrink-0 whitespace-nowrap">{{ $t('admin.mod.reports.targetDeleted') }}</UBadge>
                          <UBadge v-else-if="isReportTargetUnavailable(report)" color="warning" variant="soft" size="sm" class="shrink-0 whitespace-nowrap">{{ $t('admin.mod.reports.targetUnavailable') }}</UBadge>
                        </div>
                        <div class="mt-1 flex min-w-0 items-center gap-2 text-xs text-slate-500 dark:text-slate-400">
                          <span class="truncate">{{ targetTypeLabel(reportTargetType(report)) }} #{{ reportTargetId(report) }}</span>
                          <span class="text-slate-300 dark:text-slate-700">/</span>
                          <span class="truncate">{{ $t('admin.mod.reports.reportId', { id: report.id }) }}</span>
                        </div>
                      </div>
                    </div>
                  </td>
                  <td class="px-4 py-3 text-sm text-slate-600 dark:text-slate-300">
                    <p class="line-clamp-2 font-medium text-slate-800 dark:text-slate-100">{{ report.reason || '-' }}</p>
                    <div class="mt-2 flex min-w-0 items-center gap-2 text-xs text-slate-500 dark:text-slate-400">
                      <span>{{ $t('admin.mod.reports.reporterLabel') }}</span>
                      <IamUserAvatar :user="reportReporter(report)" size="xs" />
                      <span class="truncate">{{ userDisplayName(reportReporter(report), report.reporter_id) }}</span>
                    </div>
                  </td>
                  <td class="px-4 py-3 align-middle">
                    <UBadge :color="reportStatusColor(report.status)" variant="soft">{{ reportStatusLabel(report.status) }}</UBadge>
                  </td>
                  <td class="px-4 py-3 text-right align-middle text-sm text-slate-600 dark:text-slate-300">
                    <UTooltip :text="formatDateTime(report.created_at, locale)" :content="{ side: 'top', sideOffset: 8 }" :delay-duration="600">
                      <span>{{ relativeTime(report.created_at) }}</span>
                    </UTooltip>
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

        <section class="overflow-hidden rounded-lg border border-slate-200 bg-white dark:border-slate-800 dark:bg-slate-900 xl:sticky xl:top-20">
          <div v-if="!selectedReport" class="flex min-h-[360px] flex-col items-center justify-center px-6 py-12 text-center">
            <span class="flex size-12 items-center justify-center rounded-lg bg-slate-100 text-slate-400 dark:bg-slate-800 dark:text-slate-500">
              <UIcon name="i-lucide-mouse-pointer-2" class="size-5" />
            </span>
            <h2 class="mt-4 text-sm font-semibold text-slate-950 dark:text-white">{{ $t('admin.mod.reports.select') }}</h2>
            <p class="mt-2 max-w-xs text-sm leading-6 text-slate-500 dark:text-slate-400">{{ $t('admin.mod.reports.selectHint') }}</p>
          </div>

          <template v-else>
            <div class="flex items-center justify-between gap-3 border-b border-slate-200 px-4 py-3 dark:border-slate-800">
              <div class="min-w-0">
                <p class="text-xs text-slate-500 dark:text-slate-400">{{ $t('admin.mod.reports.current') }}</p>
                <h2 class="mt-1 text-sm font-semibold text-slate-950 dark:text-white">{{ $t('admin.mod.reports.reportId', { id: selectedReport.id }) }}</h2>
              </div>
              <UBadge :color="reportStatusColor(selectedReport.status)" variant="soft" class="whitespace-nowrap">{{ reportStatusLabel(selectedReport.status) }}</UBadge>
            </div>

            <div class="divide-y divide-slate-200 dark:divide-slate-800">
              <div class="space-y-5 p-4">
                <section class="space-y-2">
                  <div class="flex items-center justify-between gap-3">
                    <p class="text-xs font-medium text-slate-500 dark:text-slate-400">{{ $t('admin.mod.table.target') }}</p>
                    <UTooltip v-if="reportTargetPath(selectedReport)" :text="$t('admin.mod.reports.openTarget')" :content="{ side: 'top', sideOffset: 8 }" :delay-duration="600">
                      <UButton color="neutral" variant="ghost" size="xs" icon="i-lucide-external-link" :to="localePath(reportTargetPath(selectedReport))" :aria-label="$t('admin.mod.reports.openTarget')" />
                    </UTooltip>
                  </div>
                  <div class="flex items-start gap-3 border-l-2 border-slate-200 pl-3 dark:border-slate-700">
                    <span class="mt-0.5 flex size-9 shrink-0 items-center justify-center rounded-md bg-slate-100 text-slate-500 dark:bg-slate-800 dark:text-slate-300">
                      <UIcon :name="targetTypeIcon(reportTargetType(selectedReport))" class="size-4" />
                    </span>
                    <div class="min-w-0 flex-1">
                      <div class="flex min-w-0 items-center gap-2">
                        <p class="line-clamp-2 text-sm font-semibold leading-5 text-slate-950 dark:text-white">{{ reportTargetTitle(selectedReport) }}</p>
                        <UBadge v-if="isReportTargetDeleted(selectedReport)" color="neutral" variant="soft" size="sm" class="shrink-0 whitespace-nowrap">{{ $t('admin.mod.reports.targetDeleted') }}</UBadge>
                        <UBadge v-else-if="isReportTargetUnavailable(selectedReport)" color="warning" variant="soft" size="sm" class="shrink-0 whitespace-nowrap">{{ $t('admin.mod.reports.targetUnavailable') }}</UBadge>
                      </div>
                      <p class="mt-1 text-xs text-slate-500 dark:text-slate-400">
                        {{ targetTypeLabel(reportTargetType(selectedReport)) }} #{{ reportTargetId(selectedReport) }}
                      </p>
                      <div v-if="reportTarget(selectedReport).author?.id" class="mt-2 flex min-w-0 items-center gap-2">
                        <IamUserAvatar :user="reportTarget(selectedReport).author" size="xs" />
                        <span class="truncate text-xs text-slate-500 dark:text-slate-400">{{ $t('admin.mod.reports.authorLabel') }} {{ userDisplayName(reportTarget(selectedReport).author, reportTarget(selectedReport).author.id) }}</span>
                      </div>
                    </div>
                  </div>
                </section>

                <section class="space-y-2">
                  <p class="text-xs font-medium text-slate-500 dark:text-slate-400">{{ $t('admin.mod.reports.reason') }}</p>
                  <p class="rounded-md border border-amber-200 bg-amber-50/60 px-3 py-2 text-sm font-medium leading-6 text-slate-800 dark:border-amber-900/60 dark:bg-amber-950/20 dark:text-slate-100">{{ selectedReport.reason || '-' }}</p>
                </section>

                <dl class="grid grid-cols-2 gap-3 text-sm">
                  <div>
                    <dt class="text-xs text-slate-500 dark:text-slate-400">{{ $t('admin.mod.reports.reporterLabel') }}</dt>
                    <dd class="mt-1 flex min-w-0 items-center gap-2 font-medium text-slate-950 dark:text-white">
                      <IamUserAvatar :user="reportReporter(selectedReport)" size="xs" />
                      <span class="truncate">{{ userDisplayName(reportReporter(selectedReport), selectedReport.reporter_id) }}</span>
                    </dd>
                  </div>
                  <div>
                    <dt class="text-xs text-slate-500 dark:text-slate-400">{{ $t('admin.mod.table.createdAt') }}</dt>
                    <dd class="mt-1 font-medium text-slate-950 dark:text-white">{{ formatDateTime(selectedReport.created_at, locale) }}</dd>
                  </div>
                </dl>
              </div>

              <form v-if="selectedReport.status === reportStatus.pending" class="space-y-4 p-4" @submit.prevent="resolveReport">
                <div>
                  <span class="text-sm font-medium text-slate-700 dark:text-slate-200">{{ $t('admin.mod.form.result') }}</span>
                  <div class="mt-2 grid grid-cols-2 gap-2">
                    <button
                      type="button"
                      class="h-10 rounded-md border px-3 text-sm font-semibold transition"
                      :class="resolveForm.status === reportStatus.accepted ? 'border-emerald-300 bg-emerald-50 text-emerald-800 dark:border-emerald-800 dark:bg-emerald-950/30 dark:text-emerald-200' : 'border-slate-200 text-slate-600 hover:border-slate-300 dark:border-slate-800 dark:text-slate-300 dark:hover:border-slate-700'"
                      :disabled="resolving"
                      @click="resolveForm.status = reportStatus.accepted"
                    >
                      {{ $t('admin.mod.status.accepted') }}
                    </button>
                    <button
                      type="button"
                      class="h-10 rounded-md border px-3 text-sm font-semibold transition"
                      :class="resolveForm.status === reportStatus.rejected ? 'border-slate-300 bg-slate-100 text-slate-900 dark:border-slate-700 dark:bg-slate-800 dark:text-white' : 'border-slate-200 text-slate-600 hover:border-slate-300 dark:border-slate-800 dark:text-slate-300 dark:hover:border-slate-700'"
                      :disabled="resolving"
                      @click="resolveForm.status = reportStatus.rejected"
                    >
                      {{ $t('admin.mod.status.rejected') }}
                    </button>
                  </div>
                </div>

                <label class="block">
                  <span class="text-sm font-medium text-slate-700 dark:text-slate-200">{{ $t('admin.mod.form.comment') }}</span>
                  <textarea v-model="resolveForm.comment" rows="5" class="mt-2 w-full resize-none rounded-md border border-slate-200 bg-white px-3 py-2 text-sm outline-none transition focus:border-sky-300 dark:border-slate-700 dark:bg-slate-950 dark:focus:border-sky-700" :placeholder="$t('admin.mod.reports.commentPlaceholder')" :disabled="resolving" />
                </label>

                <p v-if="resolveError" class="rounded-md border border-red-200 bg-red-50 px-3 py-2 text-sm text-red-700 dark:border-red-900 dark:bg-red-950/40 dark:text-red-200">{{ resolveError }}</p>
                <div class="flex justify-end">
                  <UButton type="submit" color="primary" icon="i-lucide-check" :loading="resolving">
                    {{ $t('admin.mod.actions.submitResolve') }}
                  </UButton>
                </div>
              </form>

              <div v-else class="space-y-3 p-4">
                <p class="text-sm font-semibold text-slate-950 dark:text-white">{{ $t('admin.mod.reports.resolvedDetail') }}</p>
                <dl class="grid grid-cols-2 gap-3 text-sm">
                  <div>
                    <dt class="text-xs text-slate-500 dark:text-slate-400">{{ $t('admin.mod.form.result') }}</dt>
                    <dd class="mt-1">
                      <UBadge :color="reportStatusColor(selectedReport.status)" variant="soft" class="whitespace-nowrap">{{ reportStatusLabel(selectedReport.status) }}</UBadge>
                    </dd>
                  </div>
                  <div>
                    <dt class="text-xs text-slate-500 dark:text-slate-400">{{ $t('admin.mod.reports.handledBy') }}</dt>
                    <dd class="mt-1 flex min-w-0 items-center gap-2 font-medium text-slate-950 dark:text-white">
                      <IamUserAvatar v-if="selectedReport.dealt_by" :user="reportDealtUser(selectedReport)" size="xs" />
                      <span class="truncate">{{ selectedReport.dealt_by ? userDisplayName(reportDealtUser(selectedReport), selectedReport.dealt_by) : '-' }}</span>
                    </dd>
                  </div>
                  <div class="col-span-2">
                    <dt class="text-xs text-slate-500 dark:text-slate-400">{{ $t('admin.mod.reports.handledAt') }}</dt>
                    <dd class="mt-1 font-medium text-slate-950 dark:text-white">{{ selectedReport.dealt_at ? formatDateTime(selectedReport.dealt_at, locale) : '-' }}</dd>
                  </div>
                </dl>
                <div v-if="selectedReport.dealt_comment" class="space-y-2">
                  <p class="text-xs text-slate-500 dark:text-slate-400">{{ $t('admin.mod.form.comment') }}</p>
                  <p class="border-l-2 border-slate-200 pl-3 text-sm leading-6 text-slate-700 dark:border-slate-700 dark:text-slate-200">{{ selectedReport.dealt_comment }}</p>
                </div>
              </div>
            </div>
          </template>
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
const localePath = useLocalePath()
const toast = useToast()
const adminApi = useAdmin()
const { getTargetPath } = useAdminTargetLink()

const reports = ref<AdminModReportItem[]>([])
const total = ref(0)
const pending = ref(false)
const resolving = ref(false)
const errorMessage = ref('')
const resolveError = ref('')
const selectedReport = ref<AdminModReportItem | null>(null)
const reportStatus = {
  pending: 0,
  accepted: 1,
  rejected: 2
} as const

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

const relativeTimeFormatter = computed(() => new Intl.RelativeTimeFormat(locale.value, { numeric: 'auto' }))
const totalPages = computed(() => Math.max(1, Math.ceil(total.value / query.size)))
const statusOptions = computed(() => [
  { value: -1, label: t('admin.mod.status.all') },
  { value: reportStatus.pending, label: t('admin.mod.status.pending') },
  { value: reportStatus.accepted, label: t('admin.mod.status.accepted') },
  { value: reportStatus.rejected, label: t('admin.mod.status.rejected') }
])
const targetTypeOptions = computed(() => [
  { value: '', label: t('admin.mod.targetTypes.all') },
  { value: 'catalog_torrent', label: t('admin.mod.targetTypes.catalogTorrent') },
  { value: 'catalog_comment', label: t('admin.mod.targetTypes.catalogComment') },
  { value: 'catalog_subtitle', label: t('admin.mod.targetTypes.catalogSubtitle') },
  { value: 'forum_topic', label: t('admin.mod.targetTypes.forumTopic') },
  { value: 'forum_reply', label: t('admin.mod.targetTypes.forumReply') }
])
const targetTypeLabels = computed<Record<string, string>>(() => ({
  catalog_torrent: t('admin.mod.targetTypes.catalogTorrent'),
  catalog_comment: t('admin.mod.targetTypes.catalogComment'),
  catalog_subtitle: t('admin.mod.targetTypes.catalogSubtitle'),
  forum_topic: t('admin.mod.targetTypes.forumTopic'),
  forum_reply: t('admin.mod.targetTypes.forumReply')
}))

useHead({ title: t('admin.mod.reports.title') })
onMounted(loadReports)

async function loadReports() {
  pending.value = true
  errorMessage.value = ''
  const currentSelected = selectedReport.value
  try {
    const data = await adminApi.listModReports(query)
    reports.value = data.list || []
    total.value = data.total || 0
    if (currentSelected) {
      selectedReport.value = reports.value.find((report) => report.id === currentSelected.id) || currentSelected
    }
  } catch (error) {
    errorMessage.value = error instanceof ApiError ? error.message : t('common.requestFailed')
  } finally {
    pending.value = false
  }
}

function reloadFromFirstPage() {
  query.page = 1
  resetSelection()
  loadReports()
}

function changePage(page: number) {
  query.page = Math.min(Math.max(1, page), totalPages.value)
  resetSelection()
  loadReports()
}

function selectReport(report: AdminModReportItem) {
  selectedReport.value = report
  resolveForm.status = report.status === reportStatus.rejected ? reportStatus.rejected : reportStatus.accepted
  resolveForm.comment = ''
  resolveError.value = ''
}

function resetSelection() {
  selectedReport.value = null
  resolveForm.status = reportStatus.accepted
  resolveForm.comment = ''
  resolveError.value = ''
}

async function resolveReport() {
  if (!selectedReport.value) return
  const currentReport = selectedReport.value
  const nextStatus = resolveForm.status
  const nextComment = resolveForm.comment.trim()
  resolving.value = true
  resolveError.value = ''
  try {
    await adminApi.resolveModReport(currentReport.id, {
      status: nextStatus,
      comment: nextComment
    })
    toast.add({ title: t('admin.mod.reports.resolved'), color: 'success', icon: 'i-lucide-check-circle' })
    if (query.status >= 0 && query.status !== nextStatus) {
      resetSelection()
      await loadReports()
      return
    }
    selectedReport.value = {
      ...currentReport,
      status: nextStatus,
      dealt_comment: nextComment,
      dealt_at: new Date().toISOString()
    }
    await loadReports()
  } catch (error) {
    resolveError.value = error instanceof ApiError ? error.message : t('common.requestFailed')
  } finally {
    resolving.value = false
  }
}

function reportStatusLabel(status: number) {
  if (status === reportStatus.pending) return t('admin.mod.status.pending')
  if (status === reportStatus.accepted) return t('admin.mod.status.accepted')
  if (status === reportStatus.rejected) return t('admin.mod.status.rejected')
  return String(status)
}

function reportStatusColor(status: number) {
  if (status === reportStatus.pending) return 'warning'
  if (status === reportStatus.accepted) return 'success'
  if (status === reportStatus.rejected) return 'neutral'
  return 'neutral'
}

function targetTypeLabel(value: string) {
  return targetTypeLabels.value[value] || value || '-'
}

function targetTypeIcon(value: string) {
  if (value === 'catalog_torrent') return 'i-lucide-database'
  if (value === 'catalog_comment') return 'i-lucide-message-square'
  if (value === 'catalog_subtitle') return 'i-lucide-file-text'
  if (value === 'forum_topic') return 'i-lucide-message-circle'
  if (value === 'forum_reply') return 'i-lucide-reply'
  return 'i-lucide-flag'
}

function reportTarget(report: AdminModReportItem) {
  return report.target || {
    type: report.target_type,
    id: report.target_id,
    title: '',
    parent_type: '',
    parent_id: 0,
    status: 'normal',
    author: { id: 0, username: '', avatar: '' }
  }
}

function reportTargetType(report: AdminModReportItem) {
  return reportTarget(report).type || report.target_type
}

function reportTargetId(report: AdminModReportItem) {
  return reportTarget(report).id || report.target_id
}

function reportTargetTitle(report: AdminModReportItem) {
  const target = reportTarget(report)
  if (reportTargetStatus(report) !== 'normal') return `${targetTypeLabel(reportTargetType(report))} #${reportTargetId(report)}`
  return target.title || `${targetTypeLabel(reportTargetType(report))} #${reportTargetId(report)}`
}

function reportTargetPath(report: AdminModReportItem) {
  return getTargetPath(reportTarget(report))
}

function reportTargetStatus(report: AdminModReportItem) {
  return reportTarget(report).status || 'normal'
}

function isReportTargetDeleted(report: AdminModReportItem) {
  return reportTargetStatus(report) === 'deleted'
}

function isReportTargetUnavailable(report: AdminModReportItem) {
  return reportTargetStatus(report) === 'unavailable'
}

function reportReporter(report: AdminModReportItem) {
  return report.reporter || { id: report.reporter_id, username: '', avatar: '' }
}

function reportDealtUser(report: AdminModReportItem) {
  return report.dealt_user || { id: report.dealt_by, username: '', avatar: '' }
}

function userDisplayName(user: { id?: number, username?: string } | null | undefined, fallbackId?: number) {
  if (user?.username) return user.username
  const id = user?.id || fallbackId
  return id ? `#${id}` : '-'
}

function relativeTime(value?: string | null) {
  if (!value) return '-'
  const date = new Date(value)
  const timestamp = date.getTime()
  if (Number.isNaN(timestamp)) return '-'

  const diffSeconds = Math.round((timestamp - Date.now()) / 1000)
  const absSeconds = Math.abs(diffSeconds)
  if (absSeconds < 45) return relativeTimeFormatter.value.format(0, 'second')
  if (absSeconds < 45 * 60) return relativeTimeFormatter.value.format(Math.round(diffSeconds / 60), 'minute')
  if (absSeconds < 22 * 60 * 60) return relativeTimeFormatter.value.format(Math.round(diffSeconds / 60 / 60), 'hour')
  if (absSeconds < 30 * 24 * 60 * 60) return relativeTimeFormatter.value.format(Math.round(diffSeconds / 60 / 60 / 24), 'day')
  if (absSeconds < 12 * 30 * 24 * 60 * 60) return relativeTimeFormatter.value.format(Math.round(diffSeconds / 60 / 60 / 24 / 30), 'month')
  return relativeTimeFormatter.value.format(Math.round(diffSeconds / 60 / 60 / 24 / 365), 'year')
}
</script>
