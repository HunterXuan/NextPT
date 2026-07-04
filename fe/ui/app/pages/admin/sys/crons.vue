<template>
  <div class="min-h-[calc(100vh-4rem)] bg-slate-50 py-6 dark:bg-slate-950">
    <div class="w-full px-3 sm:px-4 lg:px-5">
      <div class="mb-4 flex justify-end">
        <UButton color="neutral" variant="outline" icon="i-lucide-refresh-cw" :loading="cronsPending || logsPending" @click="reloadAll">
          {{ $t('common.refresh') }}
        </UButton>
      </div>

      <div class="grid gap-4 xl:grid-cols-[360px_minmax(0,1fr)]">
        <section class="overflow-hidden rounded-lg border border-slate-200 bg-white dark:border-slate-800 dark:bg-slate-900">
          <div class="border-b border-slate-200 px-4 py-3 dark:border-slate-800">
            <h2 class="text-sm font-semibold text-slate-950 dark:text-white">{{ $t('admin.sys.crons.list') }}</h2>
          </div>

          <div v-if="cronsPending" class="space-y-2 p-4">
            <div v-for="item in 5" :key="item" class="h-14 animate-pulse rounded-md bg-slate-100 dark:bg-slate-800" />
          </div>
          <div v-else-if="cronsError" class="px-4 py-10 text-center text-sm text-red-600 dark:text-red-300">{{ cronsError }}</div>
          <div v-else-if="crons.length === 0" class="px-4 py-10 text-center text-sm text-slate-500 dark:text-slate-400">{{ $t('admin.sys.crons.empty') }}</div>
          <div v-else class="divide-y divide-slate-100 dark:divide-slate-800">
            <button
              v-for="cron in crons"
              :key="cron.name"
              type="button"
              class="grid w-full gap-2 px-4 py-3 text-left transition-colors"
              :class="selectedName === cron.name ? 'bg-indigo-50 dark:bg-indigo-950/30' : 'hover:bg-slate-50 dark:hover:bg-slate-950/70'"
              @click="selectCron(cron.name)"
            >
              <div class="flex min-w-0 items-center justify-between gap-3">
                <p class="truncate text-sm font-semibold text-slate-950 dark:text-white">{{ cron.name }}</p>
                <UBadge :color="cronStatusColor(cron.status)" variant="soft">{{ cronStatusLabel(cron.status) }}</UBadge>
              </div>
              <p class="text-xs text-slate-500 dark:text-slate-400">{{ formatDateTime(cron.registerTime, locale) }}</p>
            </button>
          </div>
        </section>

        <section class="overflow-hidden rounded-lg border border-slate-200 bg-white dark:border-slate-800 dark:bg-slate-900">
          <div class="flex flex-col gap-3 border-b border-slate-200 px-4 py-3 sm:flex-row sm:items-center sm:justify-between dark:border-slate-800">
            <div>
              <h2 class="text-sm font-semibold text-slate-950 dark:text-white">{{ selectedName || $t('admin.sys.crons.logsTitle') }}</h2>
              <p class="mt-1 text-xs text-slate-500 dark:text-slate-400">{{ $t('admin.sys.crons.logsTotal', { total: numberFormatter.format(logTotal) }) }}</p>
            </div>
            <UButton color="neutral" variant="outline" icon="i-lucide-refresh-cw" :loading="logsPending" :disabled="!selectedName" @click="loadLogs">
              {{ $t('common.refresh') }}
            </UButton>
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
            <table class="min-w-[860px] w-full table-fixed border-collapse text-left">
              <thead class="bg-slate-50 text-xs font-medium uppercase text-slate-500 dark:bg-slate-950/70 dark:text-slate-400">
                <tr>
                  <th class="w-[9%] border-b border-slate-200 px-4 py-3 dark:border-slate-800">ID</th>
                  <th class="w-[15%] border-b border-slate-200 px-4 py-3 dark:border-slate-800">{{ $t('admin.sys.crons.table.status') }}</th>
                  <th class="w-[14%] border-b border-slate-200 px-4 py-3 dark:border-slate-800">{{ $t('admin.sys.crons.table.duration') }}</th>
                  <th class="w-[18%] border-b border-slate-200 px-4 py-3 dark:border-slate-800">{{ $t('admin.sys.crons.table.node') }}</th>
                  <th class="w-[24%] border-b border-slate-200 px-4 py-3 dark:border-slate-800">{{ $t('admin.sys.crons.table.error') }}</th>
                  <th class="w-[20%] border-b border-slate-200 px-4 py-3 dark:border-slate-800">{{ $t('admin.sys.crons.table.time') }}</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="log in logs" :key="log.id" class="border-b border-slate-200 last:border-b-0 hover:bg-slate-50 dark:border-slate-800 dark:hover:bg-slate-950/70">
                  <td class="px-4 py-3 text-sm font-medium text-slate-950 dark:text-white">#{{ log.id }}</td>
                  <td class="px-4 py-3"><UBadge :color="cronLogStatusColor(log.status)" variant="soft">{{ cronLogStatusLabel(log.status) }}</UBadge></td>
                  <td class="px-4 py-3 text-sm text-slate-600 dark:text-slate-300">{{ numberFormatter.format(log.durationMs) }} ms</td>
                  <td class="px-4 py-3 text-sm text-slate-600 dark:text-slate-300">{{ log.nodeIp || '-' }}</td>
                  <td class="px-4 py-3 text-sm text-slate-600 dark:text-slate-300"><p class="line-clamp-2">{{ log.errorMessage || '-' }}</p></td>
                  <td class="px-4 py-3 text-sm text-slate-600 dark:text-slate-300">{{ formatDateTime(log.createdAt, locale) }}</td>
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
            :disabled="logsPending || !selectedName"
            @page-change="changeLogPage"
          />
        </section>
      </div>
    </div>
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
const logQuery = reactive({ page: 1, size: 30 })
const numberFormatter = computed(() => new Intl.NumberFormat(locale.value))
const logTotalPages = computed(() => Math.max(1, Math.ceil(logTotal.value / logQuery.size)))

useHead({ title: t('admin.sys.crons.title') })
onMounted(loadCrons)

async function reloadAll() {
  await loadCrons()
  if (selectedName.value) await loadLogs()
}

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
