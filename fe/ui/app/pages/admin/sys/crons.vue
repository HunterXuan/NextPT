<template>
  <div class="min-h-[calc(100vh-4rem)] bg-slate-50 py-6 dark:bg-slate-950">
    <div class="w-full px-3 sm:px-4 lg:px-5">
      <div class="grid gap-4 xl:grid-cols-[340px_minmax(0,1fr)]">
        <section class="overflow-hidden rounded-lg border border-slate-200 bg-white dark:border-slate-800 dark:bg-slate-900">
          <div class="flex items-center justify-between gap-3 border-b border-slate-200 px-4 py-3 dark:border-slate-800">
            <h2 class="text-sm font-semibold text-slate-950 dark:text-white">{{ $t('admin.sys.crons.list') }}</h2>
            <span class="text-xs text-slate-500 dark:text-slate-400">
              {{ $t('admin.sys.crons.tasksCount', { count: numberFormatter.format(crons.length) }) }}
            </span>
          </div>

          <div v-if="cronsPending" class="space-y-2 p-4">
            <div v-for="item in 5" :key="item" class="h-14 animate-pulse rounded-md bg-slate-100 dark:bg-slate-800" />
          </div>
          <div v-else-if="cronsError" class="px-4 py-10 text-center text-sm text-red-600 dark:text-red-300">{{ cronsError }}</div>
          <div v-else-if="crons.length === 0" class="px-4 py-10 text-center text-sm text-slate-500 dark:text-slate-400">{{ $t('admin.sys.crons.empty') }}</div>
          <div v-else class="space-y-1 p-2">
            <button
              v-for="cron in crons"
              :key="cron.name"
              type="button"
              class="group flex w-full items-center gap-3 rounded-md px-3 py-2 text-left transition-colors"
              :class="selectedName === cron.name ? 'bg-sky-50 ring-1 ring-sky-200 dark:bg-sky-950/30 dark:ring-sky-900' : 'hover:bg-slate-50 dark:hover:bg-slate-950/70'"
              @click="selectCron(cron.name)"
            >
              <span class="size-2.5 shrink-0 rounded-full" :class="cronStatusDotClass(cron.status)" />
              <UTooltip class="min-w-0 flex-1" :text="cron.name" :content="{ side: 'right', sideOffset: 8 }" :delay-duration="600">
                <p class="truncate font-mono text-[13px] font-semibold leading-6 text-slate-950 dark:text-white">{{ cron.name }}</p>
              </UTooltip>
              <span class="sr-only">{{ cronStatusLabel(cron.status) }}</span>
            </button>
          </div>
        </section>

        <section class="overflow-hidden rounded-lg border border-slate-200 bg-white dark:border-slate-800 dark:bg-slate-900">
          <div class="flex flex-col gap-3 border-b border-slate-200 px-4 py-3 sm:flex-row sm:items-center sm:justify-between dark:border-slate-800">
            <div class="flex min-w-0 items-center gap-3">
              <span class="flex size-9 shrink-0 items-center justify-center rounded-md bg-slate-100 text-slate-500 dark:bg-slate-800 dark:text-slate-300">
                <UIcon name="i-lucide-scroll-text" class="size-4" />
              </span>
              <div class="min-w-0">
                <div class="flex min-w-0 items-center gap-2">
                  <h2 class="truncate text-sm font-semibold text-slate-950 dark:text-white">{{ selectedName || $t('admin.sys.crons.logsTitle') }}</h2>
                  <UBadge v-if="selectedCron" :color="cronStatusColor(selectedCron.status)" variant="soft" class="shrink-0 whitespace-nowrap">{{ cronStatusLabel(selectedCron.status) }}</UBadge>
                </div>
                <p v-if="selectedCron" class="mt-1 truncate text-xs text-slate-500 dark:text-slate-400">
                  {{ $t('admin.sys.crons.registeredAt', { time: formatDateTime(selectedCron.registerTime, locale) }) }}
                </p>
              </div>
            </div>
            <div v-if="selectedCron" class="flex shrink-0 flex-wrap items-center gap-2">
              <div class="inline-flex h-9 items-center gap-0.5 rounded-md border border-slate-200 bg-slate-50 p-0.5 dark:border-slate-700 dark:bg-slate-950" :aria-label="$t('admin.sys.crons.filters.status')" role="group">
                <button
                  v-for="option in logStatusOptions"
                  :key="option.value"
                  type="button"
                  class="h-7 rounded px-2.5 text-xs font-medium transition disabled:cursor-not-allowed disabled:opacity-60"
                  :class="logQuery.status === option.value ? 'bg-slate-950 text-white shadow-sm dark:bg-white dark:text-slate-950' : 'text-slate-500 hover:bg-white hover:text-slate-950 dark:text-slate-400 dark:hover:bg-slate-900 dark:hover:text-white'"
                  :disabled="logsPending"
                  @click="setLogStatus(option.value)"
                >
                  {{ option.label }}
                </button>
              </div>
            </div>
          </div>

          <div v-if="!selectedName" class="flex flex-col items-center justify-center px-4 py-16 text-center">
            <UIcon name="i-lucide-clock-3" class="size-9 text-slate-400" />
            <p class="mt-3 text-sm font-medium text-slate-950 dark:text-white">{{ $t('admin.sys.crons.select') }}</p>
          </div>
          <div v-else-if="logsPending" class="space-y-2 p-4">
            <div v-for="item in 8" :key="item" class="h-16 animate-pulse rounded-md bg-slate-100 dark:bg-slate-800" />
          </div>
          <div v-else-if="logsError" class="px-4 py-10 text-center text-sm text-red-600 dark:text-red-300">{{ logsError }}</div>
          <div v-else-if="logs.length === 0" class="px-4 py-10 text-center text-sm text-slate-500 dark:text-slate-400">{{ $t('admin.sys.crons.logsEmpty') }}</div>
          <div v-else class="overflow-x-auto">
            <table class="min-w-[820px] w-full table-fixed border-collapse text-left">
              <thead class="bg-slate-50 text-xs font-medium uppercase text-slate-500 dark:bg-slate-950/70 dark:text-slate-400">
                <tr>
                  <th class="w-[10%] border-b border-slate-200 px-4 py-3 dark:border-slate-800">{{ $t('admin.sys.crons.table.status') }}</th>
                  <th class="w-[10%] border-b border-slate-200 px-4 py-3 text-right dark:border-slate-800">{{ $t('admin.sys.crons.table.duration') }}</th>
                  <th class="w-[26%] border-b border-slate-200 px-4 py-3 dark:border-slate-800">{{ $t('admin.sys.crons.table.node') }}</th>
                  <th class="w-[32%] border-b border-slate-200 px-4 py-3 dark:border-slate-800">{{ $t('admin.sys.crons.table.error') }}</th>
                  <th class="w-[22%] border-b border-slate-200 px-4 py-3 text-right dark:border-slate-800">{{ $t('admin.sys.crons.table.time') }}</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="log in logs" :key="log.id" class="border-b border-slate-200 last:border-b-0 hover:bg-slate-50 dark:border-slate-800 dark:hover:bg-slate-950/70">
                  <td class="px-4 py-2.5">
                    <div class="flex min-w-0 items-center gap-2">
                      <UBadge :color="cronLogStatusColor(log.status)" variant="soft" class="shrink-0 whitespace-nowrap">{{ cronLogStatusLabel(log.status) }}</UBadge>
                    </div>
                  </td>
                  <td class="whitespace-nowrap px-4 py-2.5 text-right font-mono text-sm text-slate-600 dark:text-slate-300">{{ formatDuration(log.durationMs) }}</td>
                  <td class="px-4 py-2.5 font-mono text-xs text-slate-600 dark:text-slate-300">
                    <UTooltip :text="log.nodeIp || '-'" :content="{ side: 'top', sideOffset: 8 }" :delay-duration="600">
                      <p class="truncate">{{ log.nodeIp || '-' }}</p>
                    </UTooltip>
                  </td>
                  <td class="px-4 py-2.5 text-sm text-slate-600 dark:text-slate-300">
                    <button
                      v-if="log.errorMessage"
                      type="button"
                      class="group flex min-w-0 items-start gap-2 text-left text-red-600 transition hover:text-red-700 dark:text-red-300 dark:hover:text-red-200"
                      @click="openError(log)"
                    >
                      <span class="line-clamp-2 min-w-0 flex-1 break-all">{{ log.errorMessage }}</span>
                      <UIcon name="i-lucide-maximize-2" class="mt-0.5 size-3.5 shrink-0 opacity-50 group-hover:opacity-100" />
                    </button>
                    <span v-else class="text-slate-300 dark:text-slate-600">-</span>
                  </td>
                  <td class="whitespace-nowrap px-4 py-2.5 text-right text-sm text-slate-600 dark:text-slate-300">
                    <UTooltip :text="formatDateTime(log.createdAt, locale)" :content="{ side: 'top', sideOffset: 8 }" :delay-duration="600">
                      <span>{{ formatDateTime(log.createdAt, locale) }}</span>
                    </UTooltip>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>

          <AppPager
            class="border-t border-slate-200 px-4 py-3 dark:border-slate-800"
            size="sm"
            :page="logQuery.page"
            :total="logTotal"
            :page-size="logQuery.size"
            :page-size-options="pageSizes"
            :disabled="logsPending || !selectedName"
            @page-change="changeLogPage"
            @page-size-change="changeLogPageSize"
          />
        </section>
      </div>
    </div>

    <UModal
      :open="errorModalOpen"
      :title="$t('admin.sys.crons.errorModal.title')"
      :description="selectedErrorLogDescription"
      :ui="{
        content: 'sm:max-w-2xl overflow-hidden',
        header: 'min-h-0 px-5 py-4 sm:px-5',
        body: 'p-0 sm:p-0',
        title: 'text-base font-semibold text-slate-950 dark:text-white',
        description: 'mt-1 text-sm text-slate-500 dark:text-slate-400',
        close: 'top-4 end-4'
      }"
      @update:open="setErrorModalOpen"
    >
      <template #body>
        <div>
          <div class="flex items-center justify-end border-b border-slate-200 bg-slate-50 px-5 py-2.5 dark:border-slate-800 dark:bg-slate-950/70">
            <UButton color="neutral" variant="soft" size="xs" icon="i-lucide-copy" :disabled="!selectedErrorLog?.errorMessage" @click="copyError">
              {{ $t('common.copy') }}
            </UButton>
          </div>
          <div class="max-h-[58vh] overflow-auto bg-slate-950">
            <pre class="min-h-40 whitespace-pre-wrap break-words px-5 py-4 font-mono text-xs leading-5 text-slate-100 selection:bg-sky-500/30">{{ selectedErrorLog?.errorMessage || '' }}</pre>
          </div>
        </div>
      </template>
    </UModal>
  </div>
</template>

<script setup lang="ts">
import { ApiError } from '~/composables/useApi'
import type { AdminSysCronItem, AdminSysCronLogItem } from '~/composables/useAdmin'
import { formatDateTime } from '~/utils/format'

definePageMeta({ layout: 'admin', middleware: 'admin' })

const { t, locale } = useI18n()
const adminApi = useAdmin()

const crons = ref<AdminSysCronItem[]>([])
const logs = ref<AdminSysCronLogItem[]>([])
const selectedName = ref('')
const cronsPending = ref(false)
const logsPending = ref(false)
const cronsError = ref('')
const logsError = ref('')
const logTotal = ref(0)
const logQuery = reactive({ page: 1, size: 30, status: -1 })
const pageSizes = [20, 30, 50, 100]
const errorModalOpen = ref(false)
const selectedErrorLog = ref<AdminSysCronLogItem | null>(null)
const numberFormatter = computed(() => new Intl.NumberFormat(locale.value))
const logTotalPages = computed(() => Math.max(1, Math.ceil(logTotal.value / logQuery.size)))
const selectedCron = computed(() => crons.value.find((item) => item.name === selectedName.value) || null)
const logStatusOptions = computed(() => [
  { value: -1, label: t('admin.sys.crons.filters.allStatuses') },
  { value: 0, label: cronLogStatusLabel(0) },
  { value: 1, label: cronLogStatusLabel(1) },
  { value: 2, label: cronLogStatusLabel(2) }
])
const selectedErrorLogDescription = computed(() => {
  if (!selectedErrorLog.value) return ''
  return `${selectedErrorLog.value.jobName} / ${t('admin.sys.crons.logId', { id: selectedErrorLog.value.id })}`
})

useHead({ title: t('admin.sys.crons.title') })
onMounted(loadCrons)

async function loadCrons() {
  cronsPending.value = true
  cronsError.value = ''
  try {
    const data = await adminApi.listSysCrons()
    crons.value = data.list || []
    if (!selectedName.value && crons.value[0]) {
      selectCron(crons.value[0].name)
    } else if (selectedName.value && !crons.value.some((item) => item.name === selectedName.value)) {
      selectedName.value = ''
      logs.value = []
      logTotal.value = 0
    }
  } catch (error) {
    cronsError.value = error instanceof ApiError ? error.message : t('common.requestFailed')
  } finally {
    cronsPending.value = false
  }
}

function selectCron(name: string) {
  if (selectedName.value === name && !logsError.value) return
  selectedName.value = name
  logQuery.page = 1
  loadLogs()
}

async function loadLogs() {
  if (!selectedName.value) return
  logsPending.value = true
  logsError.value = ''
  try {
    const data = await adminApi.listSysCronLogs(selectedName.value, logQuery)
    logs.value = data.list || []
    logTotal.value = data.total || 0
  } catch (error) {
    logsError.value = error instanceof ApiError ? error.message : t('common.requestFailed')
  } finally {
    logsPending.value = false
  }
}

function changeLogPage(page: number) {
  logQuery.page = Math.min(Math.max(1, page), logTotalPages.value)
  loadLogs()
}

function changeLogPageSize(size: number) {
  logQuery.size = size
  logQuery.page = 1
  loadLogs()
}

function setLogStatus(status: number) {
  if (logQuery.status === status) return
  logQuery.status = status
  logQuery.page = 1
  loadLogs()
}

function openError(log: AdminSysCronLogItem) {
  selectedErrorLog.value = log
  errorModalOpen.value = true
}

function setErrorModalOpen(open: boolean) {
  errorModalOpen.value = open
  if (!open) {
    selectedErrorLog.value = null
  }
}

async function copyError() {
  if (!selectedErrorLog.value?.errorMessage || typeof navigator === 'undefined' || !navigator.clipboard) return
  await navigator.clipboard.writeText(selectedErrorLog.value.errorMessage)
}

function formatDuration(durationMs: number) {
  const value = Math.max(0, Number(durationMs) || 0)
  if (value < 1000) {
    return `${numberFormatter.value.format(value)} ms`
  }
  const seconds = value / 1000
  return `${numberFormatter.value.format(Number(seconds.toFixed(seconds >= 10 ? 1 : 2)))} s`
}

function cronStatusLabel(status: number) {
  if (status === 1) return t('admin.sys.crons.status.running')
  if (status === 2) return t('admin.sys.crons.status.stopped')
  if (status === -1) return t('admin.sys.crons.status.closed')
  return t('admin.sys.crons.status.runnable')
}

function cronStatusColor(status: number) {
  if (status === 1) return 'primary'
  if (status === 2) return 'warning'
  if (status === -1) return 'error'
  return 'success'
}

function cronStatusDotClass(status: number) {
  if (status === 1) return 'bg-sky-500'
  if (status === 2) return 'bg-amber-500'
  if (status === -1) return 'bg-red-500'
  return 'bg-emerald-500'
}

function cronLogStatusLabel(status: number) {
  if (status === 1) return t('admin.sys.crons.logStatus.success')
  if (status === 2) return t('admin.sys.crons.logStatus.failed')
  return t('admin.sys.crons.logStatus.running')
}

function cronLogStatusColor(status: number) {
  if (status === 1) return 'success'
  if (status === 2) return 'error'
  return 'primary'
}
</script>
