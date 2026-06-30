<template>
  <div class="min-h-[calc(100vh-4rem)] bg-slate-50 py-6 dark:bg-slate-950">
    <div class="w-full px-3 sm:px-4 lg:px-5">
      <section class="overflow-hidden rounded-lg border border-slate-200 bg-white dark:border-slate-800 dark:bg-slate-900">
        <div class="flex flex-col gap-3 border-b border-slate-200 px-4 py-3 sm:flex-row sm:items-center sm:justify-between dark:border-slate-800">
          <div>
            <h2 class="text-sm font-semibold text-slate-950 dark:text-white">{{ $t('admin.site.audits.list') }}</h2>
            <p class="mt-1 text-xs text-slate-500 dark:text-slate-400">{{ $t('admin.site.audits.total', { total: numberFormatter.format(total) }) }}</p>
          </div>
          <UButton color="neutral" variant="outline" icon="i-lucide-refresh-cw" :loading="pending" @click="loadAudits">
            {{ $t('common.refresh') }}
          </UButton>
        </div>

        <div v-if="pending" class="space-y-2 p-4">
          <div v-for="item in 8" :key="item" class="h-16 animate-pulse rounded-md bg-slate-100 dark:bg-slate-800" />
        </div>

        <div v-else-if="errorMessage" class="flex flex-col items-center justify-center px-4 py-16 text-center">
          <UIcon name="i-lucide-circle-alert" class="size-9 text-red-500" />
          <p class="mt-3 text-sm font-medium text-slate-950 dark:text-white">{{ errorMessage }}</p>
        </div>

        <div v-else-if="audits.length === 0" class="flex flex-col items-center justify-center px-4 py-16 text-center">
          <UIcon name="i-lucide-inbox" class="size-9 text-slate-400" />
          <p class="mt-3 text-sm font-medium text-slate-950 dark:text-white">{{ $t('admin.site.audits.empty') }}</p>
        </div>

        <div v-else class="overflow-x-auto">
          <table class="min-w-[1040px] w-full table-fixed border-collapse text-left">
            <thead class="bg-slate-50 text-xs font-medium uppercase text-slate-500 dark:bg-slate-950/70 dark:text-slate-400">
              <tr>
                <th class="w-[8%] border-b border-slate-200 px-4 py-3 dark:border-slate-800">ID</th>
                <th class="w-[14%] border-b border-slate-200 px-4 py-3 dark:border-slate-800">{{ $t('admin.site.audits.table.user') }}</th>
                <th class="w-[16%] border-b border-slate-200 px-4 py-3 dark:border-slate-800">{{ $t('admin.site.audits.table.action') }}</th>
                <th class="w-[18%] border-b border-slate-200 px-4 py-3 dark:border-slate-800">{{ $t('admin.site.audits.table.target') }}</th>
                <th class="w-[22%] border-b border-slate-200 px-4 py-3 dark:border-slate-800">{{ $t('admin.site.audits.table.detail') }}</th>
                <th class="w-[10%] border-b border-slate-200 px-4 py-3 dark:border-slate-800">{{ $t('admin.site.audits.table.level') }}</th>
                <th class="w-[12%] border-b border-slate-200 px-4 py-3 dark:border-slate-800">{{ $t('admin.site.audits.table.time') }}</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="item in audits" :key="item.id" class="border-b border-slate-200 last:border-b-0 hover:bg-slate-50 dark:border-slate-800 dark:hover:bg-slate-950/70">
                <td class="px-4 py-3 text-sm font-medium text-slate-950 dark:text-white">#{{ item.id }}</td>
                <td class="px-4 py-3 text-sm text-slate-600 dark:text-slate-300">#{{ item.userId }}</td>
                <td class="px-4 py-3 text-sm font-medium text-slate-950 dark:text-white">{{ item.action || '-' }}</td>
                <td class="px-4 py-3 text-sm text-slate-600 dark:text-slate-300">{{ item.targetType || '-' }} #{{ item.targetId || 0 }}</td>
                <td class="px-4 py-3 text-sm text-slate-600 dark:text-slate-300">
                  <p class="line-clamp-2">{{ item.detail || '-' }}</p>
                  <p class="mt-1 text-xs text-slate-500">{{ item.ip || '-' }}</p>
                </td>
                <td class="px-4 py-3">
                  <UBadge :color="auditLevelColor(item.level)" variant="soft">{{ auditLevelLabel(item.level) }}</UBadge>
                </td>
                <td class="px-4 py-3 text-sm text-slate-600 dark:text-slate-300">{{ formatDateTime(item.createdAt, locale) }}</td>
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
          :page-size-options="pageSizes"
          :disabled="pending"
          @page-change="changePage"
          @page-size-change="changePageSize"
        />
      </section>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ApiError } from '~/composables/useApi'
import type { AdminSiteAuditItem } from '~/composables/useAdmin'
import { formatDateTime } from '~/utils/format'

definePageMeta({ layout: 'admin', middleware: 'admin' })

const { t, locale } = useI18n()
const adminApi = useAdmin()

const audits = ref<AdminSiteAuditItem[]>([])
const total = ref(0)
const pending = ref(false)
const errorMessage = ref('')
const query = reactive({ page: 1, size: 30 })
const pageSizes = [20, 30, 50, 100]
const numberFormatter = computed(() => new Intl.NumberFormat(locale.value))
const totalPages = computed(() => Math.max(1, Math.ceil(total.value / query.size)))

useHead({ title: `${t('admin.site.audits.title')} - NextPT` })
onMounted(loadAudits)

async function loadAudits() {
  pending.value = true
  errorMessage.value = ''
  try {
    const data = await adminApi.listSiteAudits(query)
    audits.value = data.list || []
    total.value = data.total || 0
  } catch (error) {
    errorMessage.value = error instanceof ApiError ? error.message : t('common.requestFailed')
  } finally {
    pending.value = false
  }
}

function reloadFromFirstPage() {
  query.page = 1
  loadAudits()
}

function changePageSize(size: number) {
  query.size = size
  reloadFromFirstPage()
}

function changePage(page: number) {
  query.page = Math.min(Math.max(1, page), totalPages.value)
  loadAudits()
}

function auditLevelLabel(level: number) {
  if (level === 2) return t('admin.site.audits.level.admin')
  if (level === 1) return t('admin.site.audits.level.mod')
  return t('admin.site.audits.level.normal')
}

function auditLevelColor(level: number) {
  if (level === 2) return 'error'
  if (level === 1) return 'warning'
  return 'neutral'
}
</script>
