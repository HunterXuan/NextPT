<template>
  <div class="min-h-[calc(100vh-4rem)] bg-slate-50 py-6 dark:bg-slate-950">
    <div class="w-full px-3 sm:px-4 lg:px-5">
      <section class="overflow-hidden rounded-lg border border-slate-200 bg-white dark:border-slate-800 dark:bg-slate-900">
        <div class="border-b border-slate-200 px-4 py-3 dark:border-slate-800">
          <h2 class="text-sm font-semibold text-slate-950 dark:text-white">{{ $t('admin.site.audits.list') }}</h2>
          <form class="mt-3 grid gap-2 md:grid-cols-2 xl:grid-cols-[130px_150px_170px_120px_minmax(280px,320px)_auto] xl:items-center" @submit.prevent="applyFilters">
            <USelect v-model="filters.level" class="w-full" size="lg" :ui="{ base: 'h-9 w-full' }" :items="levelOptions" value-key="value" :aria-label="$t('admin.site.audits.filters.level')" />
            <USelect v-model="filters.action" class="w-full" size="lg" :ui="{ base: 'h-9 w-full' }" :items="actionOptions" value-key="value" :aria-label="$t('admin.site.audits.filters.action')" />
            <USelect v-model="filters.targetType" class="w-full" size="lg" :ui="{ base: 'h-9 w-full' }" :items="targetTypeOptions" value-key="value" :aria-label="$t('admin.site.audits.filters.targetType')" />
            <input v-model.trim="filters.userId" inputmode="numeric" class="h-9 w-full rounded-md border border-slate-200 bg-white px-3 text-sm outline-none transition focus:border-sky-300 dark:border-slate-700 dark:bg-slate-950 dark:focus:border-sky-700" :placeholder="$t('admin.site.audits.filters.userIdPlaceholder')" :aria-label="$t('admin.site.audits.filters.userId')">
            <div class="flex h-9 min-w-0 items-center overflow-hidden rounded-md border border-slate-200 bg-white transition focus-within:border-sky-300 dark:border-slate-700 dark:bg-slate-950 dark:focus-within:border-sky-700" :aria-label="`${$t('admin.site.audits.filters.startAt')} / ${$t('admin.site.audits.filters.endAt')}`" role="group">
              <input v-model="filters.startAt" type="date" class="h-full min-w-0 flex-1 bg-transparent px-3 text-sm outline-none [color-scheme:light] dark:[color-scheme:dark]" :aria-label="$t('admin.site.audits.filters.startAt')">
              <span class="h-4 w-px shrink-0 bg-slate-200 dark:bg-slate-700" />
              <input v-model="filters.endAt" type="date" class="h-full min-w-0 flex-1 bg-transparent px-3 text-sm outline-none [color-scheme:light] dark:[color-scheme:dark]" :aria-label="$t('admin.site.audits.filters.endAt')">
            </div>
            <div class="flex items-center gap-2">
              <button type="submit" class="inline-flex h-9 items-center justify-center rounded-md bg-slate-950 px-3 text-sm font-medium text-white transition hover:bg-slate-800 disabled:cursor-not-allowed disabled:opacity-60 dark:bg-white dark:text-slate-950 dark:hover:bg-slate-200" :disabled="pending">
                {{ $t('admin.site.audits.filters.apply') }}
              </button>
              <button type="button" class="inline-flex h-9 items-center justify-center rounded-md border border-slate-200 px-3 text-sm font-medium text-slate-600 transition hover:bg-slate-50 disabled:cursor-not-allowed disabled:opacity-60 dark:border-slate-700 dark:text-slate-300 dark:hover:bg-slate-800" :disabled="pending || !hasFilters" @click="clearFilters">
                {{ $t('common.clear') }}
              </button>
            </div>
          </form>
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
          <table class="min-w-[940px] w-full table-fixed border-collapse text-left">
            <thead class="bg-slate-50 text-xs font-medium uppercase text-slate-500 dark:bg-slate-950/70 dark:text-slate-400">
              <tr>
                <th class="w-[32%] border-b border-slate-200 px-4 py-3 dark:border-slate-800">{{ $t('admin.site.audits.table.event') }}</th>
                <th class="w-[17%] border-b border-slate-200 px-4 py-3 dark:border-slate-800">{{ $t('admin.site.audits.table.user') }}</th>
                <th class="w-[29%] border-b border-slate-200 px-4 py-3 dark:border-slate-800">{{ $t('admin.site.audits.table.detail') }}</th>
                <th class="w-[10%] border-b border-slate-200 px-4 py-3 dark:border-slate-800">{{ $t('admin.site.audits.table.level') }}</th>
                <th class="w-[12%] border-b border-slate-200 px-4 py-3 text-right dark:border-slate-800">{{ $t('admin.site.audits.table.time') }}</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="item in audits" :key="item.id" class="border-b border-slate-200 last:border-b-0 hover:bg-slate-50 dark:border-slate-800 dark:hover:bg-slate-950/70">
                <td class="px-4 py-3">
                  <div class="flex min-w-0 items-center gap-3">
                    <span class="flex size-9 shrink-0 items-center justify-center rounded-md bg-slate-100 text-slate-500 dark:bg-slate-800 dark:text-slate-300">
                      <UIcon :name="auditActionIcon(item.action)" class="size-4" />
                    </span>
                    <div class="min-w-0">
                      <div class="flex min-w-0 items-center gap-2">
                        <p class="truncate text-sm font-semibold text-slate-950 dark:text-white">{{ auditActionLabel(item.action) }}</p>
                      </div>
                      <p class="mt-1 truncate text-xs text-slate-500 dark:text-slate-400">
                        {{ auditTargetLabel(item.targetType) }}<span v-if="item.targetId"> #{{ item.targetId }}</span>
                        <span class="mx-2 text-slate-300 dark:text-slate-700">/</span>
                        {{ $t('admin.site.audits.recordId', { id: item.id }) }}
                      </p>
                    </div>
                  </div>
                </td>
                <td class="px-4 py-3 text-sm text-slate-600 dark:text-slate-300">
                  <div v-if="item.userId" class="flex min-w-0 items-center gap-2.5">
                    <IamUserAvatar :user="auditActor(item)" size="sm" />
                    <div class="min-w-0">
                      <p class="truncate font-medium text-slate-950 dark:text-white">{{ actorDisplayName(item) }}</p>
                      <p class="mt-1 truncate text-xs text-slate-500 dark:text-slate-400">{{ item.ip || '-' }}</p>
                    </div>
                  </div>
                  <div v-else class="flex min-w-0 items-center gap-2.5">
                    <span class="flex size-8 shrink-0 items-center justify-center rounded-md border border-slate-200 bg-slate-100 text-slate-500 dark:border-slate-700 dark:bg-slate-800 dark:text-slate-300">
                      <UIcon name="i-lucide-cpu" class="size-4" />
                    </span>
                    <div class="min-w-0">
                      <p class="truncate font-medium text-slate-950 dark:text-white">{{ $t('admin.site.audits.systemUser') }}</p>
                      <p class="mt-1 truncate text-xs text-slate-500 dark:text-slate-400">{{ item.ip || '-' }}</p>
                    </div>
                  </div>
                </td>
                <td class="px-4 py-3 text-sm text-slate-600 dark:text-slate-300">
                  <div v-if="item.detail" class="flex min-w-0 items-start gap-2">
                    <button type="button" class="line-clamp-2 min-w-0 flex-1 break-all text-left font-mono text-xs leading-5 text-slate-500 underline-offset-2 transition hover:text-slate-950 hover:underline focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-sky-300 dark:text-slate-400 dark:hover:text-white dark:focus-visible:ring-sky-700" @click="openDetail(item)">
                      {{ item.detail }}
                    </button>
                    <UTooltip :text="$t('admin.site.audits.viewDetail')" :content="{ side: 'top', sideOffset: 8 }" :delay-duration="300">
                      <UButton color="neutral" variant="ghost" size="xs" icon="i-lucide-maximize-2" :aria-label="$t('admin.site.audits.viewDetail')" @click="openDetail(item)" />
                    </UTooltip>
                  </div>
                  <p v-else class="text-slate-400 dark:text-slate-500">{{ $t('admin.site.audits.emptyDetail') }}</p>
                </td>
                <td class="px-4 py-3">
                  <UBadge :color="auditLevelColor(item.level)" variant="soft">{{ auditLevelLabel(item.level) }}</UBadge>
                </td>
                <td class="px-4 py-3 text-right text-sm text-slate-600 dark:text-slate-300">
                  <UTooltip :text="formatDateTime(item.createdAt, locale)" :content="{ side: 'top', sideOffset: 8 }" :delay-duration="600">
                    <span>{{ formatDateTime(item.createdAt, locale) }}</span>
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
          :page-size-options="pageSizes"
          :disabled="pending"
          @page-change="changePage"
          @page-size-change="changePageSize"
        />
      </section>
    </div>

    <UModal
      :open="detailModalOpen"
      :title="$t('admin.site.audits.detailModal.title')"
      :description="selectedAudit ? auditModalDescription(selectedAudit) : ''"
      :ui="{
        content: 'sm:max-w-2xl overflow-hidden',
        header: 'min-h-0 px-5 py-4 sm:px-5',
        body: 'p-0 sm:p-0',
        title: 'text-base font-semibold text-slate-950 dark:text-white',
        description: 'mt-1 text-sm text-slate-500 dark:text-slate-400',
        close: 'top-4 end-4'
      }"
      @update:open="setDetailModalOpen"
    >
      <template #body>
        <div>
          <div class="flex items-center justify-between gap-3 border-b border-slate-200 bg-slate-50 px-5 py-2.5 dark:border-slate-800 dark:bg-slate-950/70">
            <span class="inline-flex h-6 items-center rounded border border-slate-200 bg-white px-2 font-mono text-[11px] font-semibold tracking-wide text-slate-500 dark:border-slate-700 dark:bg-slate-900 dark:text-slate-300">
              {{ selectedAuditDetailType }}
            </span>
            <UButton color="neutral" variant="soft" size="xs" icon="i-lucide-copy" :disabled="!selectedAudit?.detail" @click="copyDetail">
              {{ $t('common.copy') }}
            </UButton>
          </div>
          <div class="max-h-[62vh] overflow-auto bg-slate-950">
            <pre class="min-h-56 whitespace-pre-wrap break-words px-5 py-4 font-mono text-xs leading-5 text-slate-100 selection:bg-sky-500/30">{{ selectedAuditDetailText }}</pre>
          </div>
        </div>
      </template>
    </UModal>
  </div>
</template>

<script setup lang="ts">
import { ApiError } from '~/composables/useApi'
import type { AdminSiteAuditItem, AdminSiteAuditListParams } from '~/composables/useAdmin'
import { formatDateTime } from '~/utils/format'

definePageMeta({ layout: 'admin', middleware: 'admin' })

const { t, te, locale } = useI18n()
const adminApi = useAdmin()

const audits = ref<AdminSiteAuditItem[]>([])
const total = ref(0)
const pending = ref(false)
const errorMessage = ref('')
const detailModalOpen = ref(false)
const selectedAudit = ref<AdminSiteAuditItem | null>(null)
const query = reactive({ page: 1, size: 30 })
const filters = reactive({
  level: '',
  action: '',
  targetType: '',
  userId: '',
  startAt: '',
  endAt: ''
})
const pageSizes = [20, 30, 50, 100]
const totalPages = computed(() => Math.max(1, Math.ceil(total.value / query.size)))
const hasFilters = computed(() => Boolean(
  filters.level
  || filters.action
  || filters.targetType
  || filters.userId
  || filters.startAt
  || filters.endAt
))
const levelOptions = computed(() => [
  { value: '', label: t('admin.site.audits.filters.allLevels') },
  { value: '0', label: auditLevelLabel(0) },
  { value: '1', label: auditLevelLabel(1) },
  { value: '2', label: auditLevelLabel(2) }
])
const actionOptions = computed(() => [
  { value: '', label: t('admin.site.audits.filters.allActions') },
  ...auditActions.map(action => ({ value: action, label: auditActionLabel(action) }))
])
const targetTypeOptions = computed(() => [
  { value: '', label: t('admin.site.audits.filters.allTargets') },
  ...auditTargetTypes.map(targetType => ({ value: targetType, label: auditTargetLabel(targetType) }))
])
const selectedAuditDetailText = computed(() => formatAuditDetail(selectedAudit.value?.detail || ''))
const selectedAuditDetailType = computed(() => isJsonAuditDetail(selectedAudit.value?.detail || '') ? 'JSON' : 'TEXT')

const auditActions = [
  'create',
  'update',
  'delete'
]
const auditTargetTypes = [
  'site_config',
  'catalog_category',
  'catalog_request',
  'catalog_torrent',
  'forum_category',
  'forum_node',
  'forum_topic',
  'iam_role',
  'iam_user',
  'iam_invite',
  'iam_session',
  'mod_report',
  'mod_cheater_log',
  'mod_user_log'
]

useHead({ title: t('admin.site.audits.title') })
onMounted(loadAudits)

async function loadAudits() {
  pending.value = true
  errorMessage.value = ''
  try {
    const data = await adminApi.listSiteAudits(buildAuditQuery())
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

function applyFilters() {
  reloadFromFirstPage()
}

function clearFilters() {
  filters.level = ''
  filters.action = ''
  filters.targetType = ''
  filters.userId = ''
  filters.startAt = ''
  filters.endAt = ''
  reloadFromFirstPage()
}

function buildAuditQuery(): AdminSiteAuditListParams {
  const params: AdminSiteAuditListParams = {
    page: query.page,
    size: query.size
  }
  if (filters.level !== '') {
    params.level = Number(filters.level)
  }
  const userId = Number(filters.userId)
  if (Number.isInteger(userId) && userId > 0) {
    params.userId = userId
  }
  if (filters.action) {
    params.action = filters.action
  }
  if (filters.targetType) {
    params.targetType = filters.targetType
  }
  if (filters.startAt) {
    params.startAt = `${filters.startAt} 00:00:00`
  }
  if (filters.endAt) {
    params.endAt = `${filters.endAt} 23:59:59`
  }
  return params
}

function changePageSize(size: number) {
  query.size = size
  reloadFromFirstPage()
}

function changePage(page: number) {
  query.page = Math.min(Math.max(1, page), totalPages.value)
  loadAudits()
}

function openDetail(item: AdminSiteAuditItem) {
  selectedAudit.value = item
  detailModalOpen.value = true
}

function setDetailModalOpen(open: boolean) {
  detailModalOpen.value = open
  if (!open) {
    selectedAudit.value = null
  }
}

function auditModalDescription(item: AdminSiteAuditItem) {
  return `${auditTargetLabel(item.targetType)}${item.targetId ? ` #${item.targetId}` : ''} / ${t('admin.site.audits.recordId', { id: item.id })}`
}

async function copyDetail() {
  if (!selectedAudit.value?.detail || typeof navigator === 'undefined' || !navigator.clipboard) return
  await navigator.clipboard.writeText(selectedAudit.value.detail)
}

function formatAuditDetail(detail: string) {
  if (!detail) return ''
  try {
    return JSON.stringify(JSON.parse(detail), null, 2)
  } catch {
    return detail
  }
}

function isJsonAuditDetail(detail: string) {
  if (!detail) return false
  try {
    JSON.parse(detail)
    return true
  } catch {
    return false
  }
}

function auditLevelLabel(level: number) {
  if (level === 2) return t('admin.site.audits.level.critical')
  if (level === 1) return t('admin.site.audits.level.important')
  return t('admin.site.audits.level.normal')
}

function auditLevelColor(level: number) {
  if (level === 2) return 'error'
  if (level === 1) return 'warning'
  return 'neutral'
}

function auditActionLabel(action: string) {
  return translateAuditValue('actions', action)
}

function auditTargetLabel(targetType: string) {
  return translateAuditValue('targets', targetType)
}

function translateAuditValue(group: string, value: string) {
  if (!value) return '-'
  const key = `admin.site.audits.${group}.${value}`
  return te(key) ? t(key) : humanizeKey(value)
}

function auditActionIcon(action: string) {
  const icons: Record<string, string> = {
    create: 'i-lucide-plus',
    update: 'i-lucide-pencil',
    delete: 'i-lucide-trash-2'
  }
  return icons[action] || 'i-lucide-scroll-text'
}

function auditActor(item: AdminSiteAuditItem) {
  return item.actor?.id ? item.actor : { id: item.userId, username: '', avatar: '' }
}

function actorDisplayName(item: AdminSiteAuditItem) {
  return item.actor?.username || (item.userId ? `#${item.userId}` : t('admin.site.audits.systemUser'))
}

function humanizeKey(value: string) {
  return value
    .replace(/_/g, ' ')
    .replace(/\b\w/g, letter => letter.toUpperCase())
}
</script>
