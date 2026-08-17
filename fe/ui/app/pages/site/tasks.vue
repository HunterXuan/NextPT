<template>
  <div class="min-h-[calc(100vh-4rem)] bg-slate-50 py-6 dark:bg-slate-950">
    <div class="w-full px-3 sm:px-4 lg:px-5">
      <div v-if="pending" class="grid gap-3 md:grid-cols-2 2xl:grid-cols-3">
        <div v-for="index in 6" :key="index" class="h-52 animate-pulse rounded-lg border border-slate-200 bg-white dark:border-slate-800 dark:bg-slate-900" />
      </div>
      <div v-else-if="errorMessage" class="border-y border-red-200 bg-red-50 px-4 py-10 text-center text-sm text-red-700 dark:border-red-950 dark:bg-red-950/30 dark:text-red-200">{{ errorMessage }}</div>
      <div v-else-if="tasks.length === 0" class="flex items-center justify-center gap-2 border-y border-slate-200 px-4 py-10 text-sm text-slate-500 dark:border-slate-800 dark:text-slate-400">
        <UIcon name="i-lucide-clipboard-list" class="size-4 text-slate-400 dark:text-slate-500" />
        <p>{{ $t('site.tasks.empty') }}</p>
      </div>
      <div v-else :class="taskCategories.length > 2 ? 'items-start lg:grid lg:grid-cols-[220px_minmax(0,1fr)] lg:gap-4' : ''">
        <nav v-if="taskCategories.length > 2" class="mb-4 flex gap-2 overflow-x-auto rounded-lg border border-slate-200 bg-white p-1 dark:border-slate-800 dark:bg-slate-900 lg:mb-0 lg:block lg:space-y-1 lg:overflow-visible">
          <button
            v-for="category in taskCategories"
            :key="category"
            type="button"
            class="flex min-w-max items-center gap-2 rounded-md px-3 py-2 text-left text-sm transition lg:w-full lg:min-w-0"
            :class="activeCategory === category
              ? 'bg-slate-100 text-slate-950 dark:bg-slate-800 dark:text-white'
              : 'text-slate-600 hover:bg-slate-50 hover:text-slate-950 dark:text-slate-300 dark:hover:bg-slate-950 dark:hover:text-white'"
            :aria-pressed="activeCategory === category"
            @click="activeCategory = category"
          >
            <span
              class="flex size-8 shrink-0 items-center justify-center rounded-md"
              :class="activeCategory === category
                ? 'bg-white text-slate-950 dark:bg-slate-950 dark:text-white'
                : 'bg-slate-100 text-slate-500 dark:bg-slate-800 dark:text-slate-400'"
            >
              <UIcon :name="categoryIcon(category)" class="size-4" />
            </span>
            <span class="font-medium">{{ categoryLabel(category) }}</span>
          </button>
        </nav>

        <div class="grid items-start gap-3 md:grid-cols-2 2xl:grid-cols-3">
          <article v-for="task in filteredTasks" :key="task.key" class="rounded-lg border border-slate-200 bg-white p-4 transition-shadow hover:shadow-sm dark:border-slate-800 dark:bg-slate-900">
            <div class="flex items-start justify-between gap-3">
              <div class="min-w-0">
                <p class="mb-1 text-xs font-medium text-slate-500 dark:text-slate-400">{{ categoryLabel(taskCategory(task)) }}</p>
                <div class="flex flex-wrap items-center gap-2">
                  <h2 class="truncate text-sm font-semibold text-slate-950 dark:text-white">{{ taskName(task) }}</h2>
                  <UBadge :color="statusColor(task)" variant="soft" size="sm">{{ statusLabel(task) }}</UBadge>
                </div>
                <p v-if="taskDescription(task)" class="mt-1 line-clamp-2 text-xs leading-5 text-slate-500 dark:text-slate-400">{{ taskDescription(task) }}</p>
              </div>
              <div class="flex shrink-0 items-center gap-1 rounded-md bg-amber-50 px-2 py-1 text-xs font-semibold tabular-nums text-amber-700 dark:bg-amber-950/60 dark:text-amber-300">
                <UIcon name="i-lucide-coins" class="size-3.5" />
                <span>{{ rewardsLabel(task) }}</span>
              </div>
            </div>

            <div v-if="task.userTask" class="mt-4">
              <div class="flex flex-wrap items-center justify-between gap-x-3 gap-y-1 text-xs">
                <span class="font-medium tabular-nums text-slate-700 dark:text-slate-200">{{ progressLabel(task) }}</span>
                <span v-if="task.userTask.cycleEndedAt" class="text-slate-500 dark:text-slate-400">{{ $t('site.tasks.endsAt', { time: formatDateTime(task.userTask.cycleEndedAt, locale) }) }}</span>
              </div>
              <div class="mt-2 h-1.5 overflow-hidden rounded-full bg-slate-100 dark:bg-slate-800">
                <div class="h-full rounded-full bg-emerald-500 transition-[width]" :style="{ width: `${progressPercent(task)}%` }" />
              </div>
            </div>
            <div v-else class="mt-4 text-xs text-slate-500 dark:text-slate-400">{{ taskRuleLabel(task) }}</div>

            <div v-if="!task.userTask || task.userTask.status === 1" class="mt-4 flex justify-end border-t border-slate-100 pt-3 dark:border-slate-800">
              <UButton v-if="!task.userTask" color="primary" variant="soft" size="sm" icon="i-lucide-hand" :loading="workingKey === task.key" :disabled="!task.enabled || workingKey !== ''" @click="claim(task)">{{ $t('site.tasks.claim') }}</UButton>
              <UButton v-else color="success" variant="soft" size="sm" icon="i-lucide-gift" :loading="workingKey === task.key" :disabled="workingKey !== ''" @click="claimReward(task)">{{ $t('site.tasks.claimReward') }}</UButton>
            </div>
          </article>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { SiteTask } from '~/composables/useSite'
import type { AuthRoleItem } from '~/composables/useAuth'
import { ApiError } from '~/composables/useApi'
import { formatBytes, formatDateTime, localizeI18nName } from '~/utils/format'

definePageMeta({ middleware: 'auth' })
const { t, locale } = useI18n()
const toast = useToast()
const siteApi = useSite()
const { listRoles } = useAuth()
const tasks = ref<SiteTask[]>([])
const roles = ref<AuthRoleItem[]>([])
const pending = ref(false)
const errorMessage = ref('')
const workingKey = ref('')
const activeCategory = ref<TaskCategory>('all')

type TaskCategory = 'all' | 'publishing' | 'sharing' | 'growth'

const taskCategories = computed<TaskCategory[]>(() => {
  const categories: TaskCategory[] = ['all']
  for (const category of ['publishing', 'sharing', 'growth'] as const) {
    if (tasks.value.some(task => taskCategory(task) === category)) categories.push(category)
  }
  return categories
})
const filteredTasks = computed(() => activeCategory.value === 'all'
  ? tasks.value
  : tasks.value.filter(task => taskCategory(task) === activeCategory.value))

useSeoMeta({ title: t('site.tasks.title'), robots: 'noindex, nofollow' })
onMounted(load)

async function load() {
  pending.value = true
  errorMessage.value = ''
  try {
    const [taskData, roleData] = await Promise.all([
      siteApi.listTasks(),
      listRoles().catch(() => ({ roles: [] }))
    ])
    tasks.value = taskData.list || []
    roles.value = roleData.roles || []
    if (!taskCategories.value.includes(activeCategory.value)) activeCategory.value = 'all'
  } catch (error: unknown) {
    errorMessage.value = error instanceof ApiError ? error.message : t('common.requestFailed')
  } finally {
    pending.value = false
  }
}
async function claim(task: SiteTask) {
  workingKey.value = task.key
  try {
    await siteApi.claimTask(task.key)
    toast.add({ title: t('site.tasks.claimed') })
    await load()
  } finally { workingKey.value = '' }
}
async function claimReward(task: SiteTask) {
  if (!task.userTask) return
  workingKey.value = task.key
  try {
    await siteApi.claimTaskReward(task.userTask.id)
    toast.add({ title: t('site.tasks.rewarded') })
    await load()
  } finally { workingKey.value = '' }
}
function taskName(task: SiteTask) { return localizeI18nName(task.nameI18n, locale.value, task.key) }
function taskDescription(task: SiteTask) {
  const description = localizeI18nName(task.descriptionI18n, locale.value, '')
  return description === taskName(task) ? '' : description
}
function statusLabel(task: SiteTask) {
  if (!task.userTask) return task.enabled ? t('site.tasks.status.available') : t('site.tasks.status.disabled')
  return t(`site.tasks.status.${['active', 'completed', 'rewarded', 'expired'][task.userTask.status] || 'active'}`)
}
function statusColor(task: SiteTask) {
  const status = task.userTask?.status
  if (status === 1 || status === 2) return 'success'
  if (status === 3 || !task.enabled) return 'neutral'
  return 'primary'
}
function progressPercent(task: SiteTask) {
  const item = task.userTask
  if (!item || item.target <= 0) return 0
  return Math.min(100, Math.round(item.progress / item.target * 100))
}
function progressLabel(task: SiteTask) {
  const item = task.userTask
  if (!item) return ''
  if (task.rule.type === 'iam.role_level_reached') return taskRuleLabel(task)
  if (task.rule.type === 'tracker.uploaded') return `${formatBytes(item.progress)} / ${formatBytes(item.target)}`
  if (task.rule.type === 'tracker.seed_duration') return `${formatHours(item.progress)} / ${formatHours(item.target)}`
  return `${item.progress} / ${item.target}`
}
function rewardsLabel(task: SiteTask) { return task.rewards.map(reward => t(`site.tasks.rewards.${reward.type}`, { amount: reward.amount })).join(' · ') }
function taskRuleLabel(task: SiteTask) {
  if (task.rule.type === 'catalog.torrent_published') return t('site.tasks.rules.torrents', { value: task.rule.target })
  if (task.rule.type === 'tracker.seed_duration') return t('site.tasks.rules.seedHours', { value: task.rule.target })
  if (task.rule.type === 'tracker.uploaded') return t('site.tasks.rules.uploaded', { value: task.rule.target })
  if (task.rule.type === 'iam.role_level_reached') {
    const role = roles.value.find(item => item.level === task.rule.target)
    return role
      ? t('site.tasks.rules.roleLevel', { value: localizeI18nName(role.nameI18N, locale.value, role.name) })
      : t('site.tasks.rules.roleLevelFallback')
  }
  return ''
}
function taskCategory(task: SiteTask): Exclude<TaskCategory, 'all'> {
  if (task.rule.type === 'catalog.torrent_published') return 'publishing'
  if (task.rule.type === 'iam.role_level_reached') return 'growth'
  return 'sharing'
}
function categoryLabel(category: TaskCategory) { return t(`site.tasks.categories.${category}`) }
function categoryIcon(category: TaskCategory) {
  if (category === 'publishing') return 'i-lucide-upload'
  if (category === 'sharing') return 'i-lucide-share-2'
  if (category === 'growth') return 'i-lucide-trending-up'
  return 'i-lucide-layout-grid'
}
function formatHours(seconds: number) { return t('site.tasks.hours', { value: Math.floor(seconds / 3600) }) }
</script>
