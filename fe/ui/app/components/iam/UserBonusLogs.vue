<template>
  <UCard class="rounded-lg">
    <template #header>
      <div>
        <h2 class="text-base font-semibold text-slate-950 dark:text-white">{{ $t('user.bonusLogs.title') }}</h2>
        <p class="mt-1 text-sm text-slate-500 dark:text-slate-400">
          {{ $t('user.bonusLogs.description') }}
        </p>
      </div>
    </template>

    <div class="grid grid-cols-1 gap-3 sm:grid-cols-2">
      <div class="rounded-md border border-slate-200 bg-slate-50 px-4 py-3 dark:border-slate-800 dark:bg-slate-950/50">
        <div class="flex items-center justify-between gap-3">
          <p class="text-xs text-slate-500 dark:text-slate-400">{{ $t('user.bonusLogs.balance') }}</p>
          <UIcon name="i-lucide-wallet-cards" class="size-4 text-slate-400" />
        </div>
        <p class="mt-1.5 text-xl font-semibold tabular-nums text-slate-950 dark:text-white">{{ formatBonus(user?.stat.bonus || 0) }}</p>
      </div>
      <div class="rounded-md border border-slate-200 bg-slate-50 px-4 py-3 dark:border-slate-800 dark:bg-slate-950/50">
        <div class="flex items-center justify-between gap-3">
          <div class="flex items-center gap-1.5">
            <p class="text-xs text-slate-500 dark:text-slate-400">{{ $t('user.bonusLogs.hourly') }}</p>
            <UPopover :content="{ side: 'top', align: 'start', sideOffset: 8 }" :ui="{ content: 'w-80 p-4' }">
              <UButton
                color="neutral"
                variant="ghost"
                size="xs"
                icon="i-lucide-info"
                class="size-6 rounded-full text-slate-400 hover:text-slate-700 dark:hover:text-slate-200"
                :aria-label="$t('user.bonusLogs.formula.title')"
              />
              <template #content>
                <div class="space-y-3 text-sm">
                  <div>
                    <p class="font-semibold text-slate-950 dark:text-white">{{ $t('user.bonusLogs.formula.title') }}</p>
                    <p class="mt-1 text-xs leading-5 text-slate-500 dark:text-slate-400">{{ $t('user.bonusLogs.formula.description') }}</p>
                  </div>
                  <div class="space-y-2 rounded-md bg-slate-50 p-3 font-mono text-xs leading-5 text-slate-700 dark:bg-slate-950 dark:text-slate-200">
                    <p>{{ $t('user.bonusLogs.formula.score') }}</p>
                    <p>{{ $t('user.bonusLogs.formula.volume') }}</p>
                    <p>{{ $t('user.bonusLogs.formula.total') }}</p>
                  </div>
                  <ul class="space-y-1 text-xs leading-5 text-slate-500 dark:text-slate-400">
                    <li>{{ $t('user.bonusLogs.formula.variables.size') }}</li>
                    <li>{{ $t('user.bonusLogs.formula.variables.age') }}</li>
                    <li>{{ $t('user.bonusLogs.formula.variables.seeders') }}</li>
                    <li>{{ $t('user.bonusLogs.formula.variables.params') }}</li>
                  </ul>
                </div>
              </template>
            </UPopover>
          </div>
          <UIcon name="i-lucide-clock-3" class="size-4 text-slate-400" />
        </div>
        <p class="mt-1.5 text-xl font-semibold tabular-nums text-slate-950 dark:text-white">{{ formatBonus(hourlyBonus) }}</p>
      </div>
    </div>

    <div class="mt-5 border-t border-slate-200 pt-4 dark:border-slate-800">
      <div class="flex items-center justify-between gap-3">
        <div>
          <h3 class="text-sm font-semibold text-slate-950 dark:text-white">{{ $t('user.bonusLogs.detailTitle') }}</h3>
        </div>
      </div>

      <div v-if="bonusLogsError" class="mt-4 flex flex-col items-center justify-center rounded-md border border-red-200 bg-red-50 px-4 py-8 text-center dark:border-red-900 dark:bg-red-950">
        <UIcon name="i-lucide-circle-alert" class="size-8 text-red-500" />
        <p class="mt-3 text-sm font-medium text-red-700 dark:text-red-200">{{ bonusLogsError }}</p>
      </div>
      <div v-else-if="bonusLogsPending && bonusLogs.length === 0" class="mt-4 space-y-3">
        <div v-for="item in 6" :key="item" class="h-14 animate-pulse rounded-md bg-slate-100 dark:bg-slate-800" />
      </div>
      <div v-else-if="bonusLogs.length === 0" class="mt-4 rounded-md border border-dashed border-slate-200 px-4 py-12 text-center dark:border-slate-800">
        <UIcon name="i-lucide-coins" class="mx-auto size-8 text-slate-400" />
        <p class="mt-3 text-sm text-slate-500 dark:text-slate-400">{{ $t('user.bonusLogs.empty') }}</p>
      </div>
      <div v-else class="mt-2 divide-y divide-slate-100 dark:divide-slate-800">
        <article v-for="log in bonusLogs" :key="log.id" class="grid gap-3 py-3 text-sm sm:grid-cols-[minmax(0,1fr)_180px] sm:items-center">
          <div class="flex min-w-0 items-center gap-3">
            <div class="flex size-8 shrink-0 items-center justify-center rounded-md" :class="bonusActionIconClass(log)">
              <UIcon :name="bonusActionIcon(log.action)" class="size-4" />
            </div>
            <div class="min-w-0">
              <p class="truncate font-medium text-slate-950 dark:text-white">{{ bonusActionLabel(log) }}</p>
              <p v-if="bonusLogDescription(log)" class="mt-1 truncate text-xs text-slate-500 dark:text-slate-400">
                {{ bonusLogDescription(log) }}
              </p>
            </div>
          </div>
          <div class="sm:text-right">
            <p
              class="font-semibold tabular-nums"
              :class="log.amount >= 0 ? 'text-emerald-600 dark:text-emerald-300' : 'text-red-600 dark:text-red-300'"
            >
              {{ log.amount >= 0 ? '+' : '' }}{{ formatBonus(log.amount) }}
            </p>
            <p class="mt-1 text-xs text-slate-500 dark:text-slate-400">{{ $t('user.bonusLogs.after', { balance: formatBonus(log.balanceAfter) }) }}</p>
            <UTooltip
              :text="formatDateTime(log.createdAt, locale)"
              :content="{ side: 'top', sideOffset: 8 }"
              :delay-duration="120"
            >
              <p class="mt-1 text-xs text-slate-500 dark:text-slate-400">{{ formatDateOnly(log.createdAt, locale) }}</p>
            </UTooltip>
          </div>
        </article>
      </div>

      <AppPager
        class="mt-4 border-t border-slate-200 pt-4 dark:border-slate-800"
        size="sm"
        :page="bonusLogPage"
        :total="bonusLogTotal"
        :page-size="bonusLogSize"
        :disabled="bonusLogsPending"
        @page-change="goToBonusLogPage"
      />
    </div>
  </UCard>
</template>

<script setup lang="ts">
import { ApiError } from '~/composables/useApi'
import { useEconomy, type BonusLogItem } from '~/composables/useEconomy'
import { formatDateOnly, formatDateTime } from '~/utils/format'

const { t, locale } = useI18n()
const { user, fetchUser } = useAuth()
const economy = useEconomy()

const bonusLogs = ref<BonusLogItem[]>([])
const bonusLogTotal = ref(0)
const bonusLogPage = ref(1)
const bonusLogSize = 10
const hourlyBonus = ref(0)
const bonusLogsPending = ref(true)
const bonusLogsError = ref('')

const numberFormatter = computed(() => new Intl.NumberFormat(locale.value))
const bonusLogTotalPages = computed(() => Math.max(1, Math.ceil(bonusLogTotal.value / bonusLogSize)))

onMounted(async () => {
  await Promise.all([
    fetchUser(),
    loadHourlyBonus(),
    loadBonusLogs()
  ])
})

async function loadHourlyBonus() {
  try {
    const data = await economy.getHourlyBonus()
    hourlyBonus.value = data.hourlyBonus || 0
  } catch {
    hourlyBonus.value = 0
  }
}

async function loadBonusLogs() {
  if (bonusLogsPending.value && bonusLogs.value.length > 0) return

  bonusLogsPending.value = true
  bonusLogsError.value = ''

  try {
    let data = await economy.listBonusLogs({
      page: bonusLogPage.value,
      size: bonusLogSize
    })

    const nextTotal = data.total || 0
    const nextTotalPages = Math.max(1, Math.ceil(nextTotal / bonusLogSize))
    if (bonusLogPage.value > nextTotalPages) {
      bonusLogPage.value = nextTotalPages
      data = await economy.listBonusLogs({
        page: bonusLogPage.value,
        size: bonusLogSize
      })
    }

    bonusLogs.value = data.list || []
    bonusLogTotal.value = data.total || 0
  } catch (error) {
    bonusLogs.value = []
    bonusLogTotal.value = 0
    bonusLogsError.value = error instanceof ApiError ? error.message : t('common.requestFailed')
  } finally {
    bonusLogsPending.value = false
  }
}

function goToBonusLogPage(page: number) {
  bonusLogPage.value = Math.min(Math.max(1, page), bonusLogTotalPages.value)
  loadBonusLogs()
}

function formatBonus(value?: number | string | null) {
  return numberFormatter.value.format(Number(value || 0))
}

function bonusActionLabel(log: BonusLogItem) {
  if (isRewardTarget(log.targetType)) {
    if (log.action === 'transfer_sent') return t('user.bonusLogs.actions.reward_sent')
    if (log.action === 'transfer_received') return t('user.bonusLogs.actions.reward_received')
  }

  const action = log.action
  if (!action) return '-'
  const key = `user.bonusLogs.actions.${action}`
  const translated = t(key)
  return translated === key ? action : translated
}

function bonusActionIcon(action?: string | null) {
  switch (action) {
    case 'transfer_sent':
      return 'i-lucide-arrow-up-right'
    case 'transfer_received':
      return 'i-lucide-arrow-down-left'
    case 'seed_bonus':
      return 'i-lucide-sprout'
    case 'admin_adjustment':
      return 'i-lucide-shield'
    case 'request_escrow':
      return 'i-lucide-lock-keyhole'
    case 'request_reward':
      return 'i-lucide-circle-check'
    case 'request_refund':
      return 'i-lucide-undo-2'
    case 'shop_purchase':
      return 'i-lucide-shopping-bag'
    case 'forum_topic':
      return 'i-lucide-message-square'
    case 'forum_reply':
      return 'i-lucide-message-circle'
    default:
      return 'i-lucide-coins'
  }
}

function bonusActionIconClass(log: BonusLogItem) {
  if (log.action === 'seed_bonus') return 'bg-emerald-50 text-emerald-600 dark:bg-emerald-950/50 dark:text-emerald-300'
  if (log.action === 'request_reward' || log.action === 'request_refund') return 'bg-emerald-50 text-emerald-600 dark:bg-emerald-950/50 dark:text-emerald-300'
  if (log.amount < 0) return 'bg-red-50 text-red-600 dark:bg-red-950/50 dark:text-red-300'
  return 'bg-amber-50 text-amber-600 dark:bg-amber-950/50 dark:text-amber-300'
}

function bonusLogDescription(log: BonusLogItem) {
  if (log.action === 'seed_bonus' && log.period?.trim()) return t('user.bonusLogs.period', { period: log.period.trim() })
  if (log.targetType && log.targetId) return bonusTargetLabel(log.targetType, log.targetId)
  if (log.remark?.trim()) return log.remark.trim()
  return ''
}

function isRewardTarget(targetType?: string | null) {
  return ['catalog_torrent', 'catalog_comment', 'forum_topic', 'forum_reply'].includes(targetType || '')
}

function bonusTargetLabel(targetType: string, targetId: number) {
  const key = `user.bonusLogs.targetTypes.${targetType}`
  const translated = t(key)
  const label = translated === key ? targetType : translated
  return `${label} #${targetId}`
}
</script>
