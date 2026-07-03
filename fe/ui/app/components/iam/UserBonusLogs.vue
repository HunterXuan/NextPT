<template>
  <UCard class="rounded-lg">
    <template #header>
      <div class="flex flex-col gap-3 sm:flex-row sm:items-start sm:justify-between">
        <div>
          <h2 class="text-base font-semibold text-slate-950 dark:text-white">{{ $t('user.bonusLogs.title') }}</h2>
          <p class="mt-1 text-sm text-slate-500 dark:text-slate-400">
            {{ $t('user.bonusLogs.summary', { total: numberFormatter.format(bonusLogTotal), hourly: formatBonus(hourlyBonus) }) }}
          </p>
        </div>
        <UButton color="neutral" variant="outline" size="sm" icon="i-lucide-refresh-cw" :loading="bonusLogsPending" @click="loadBonusLogs">
          {{ $t('common.refresh') }}
        </UButton>
      </div>
    </template>

    <div class="grid grid-cols-1 gap-3 sm:grid-cols-2">
      <div class="rounded-md border border-slate-200 px-3 py-3 dark:border-slate-800">
        <p class="text-xs text-slate-500 dark:text-slate-400">{{ $t('user.bonusLogs.hourly') }}</p>
        <p class="mt-1 text-lg font-semibold text-slate-950 dark:text-white">{{ formatBonus(hourlyBonus) }}</p>
      </div>
      <div class="rounded-md border border-slate-200 px-3 py-3 dark:border-slate-800">
        <p class="text-xs text-slate-500 dark:text-slate-400">{{ $t('user.bonusLogs.balance') }}</p>
        <p class="mt-1 text-lg font-semibold text-slate-950 dark:text-white">{{ formatBonus(user?.bonus || 0) }}</p>
      </div>
    </div>

    <div v-if="bonusLogsError" class="mt-4 flex flex-col items-center justify-center rounded-md border border-red-200 bg-red-50 px-4 py-8 text-center dark:border-red-900 dark:bg-red-950">
      <UIcon name="i-lucide-circle-alert" class="size-8 text-red-500" />
      <p class="mt-3 text-sm font-medium text-red-700 dark:text-red-200">{{ bonusLogsError }}</p>
      <UButton class="mt-5" color="neutral" variant="outline" size="sm" icon="i-lucide-refresh-cw" @click="loadBonusLogs">
        {{ $t('common.retry') }}
      </UButton>
    </div>
    <div v-else-if="bonusLogsPending && bonusLogs.length === 0" class="mt-4 space-y-3">
      <div v-for="item in 6" :key="item" class="h-14 animate-pulse rounded-md bg-slate-100 dark:bg-slate-800" />
    </div>
    <div v-else-if="bonusLogs.length === 0" class="mt-4 rounded-md border border-dashed border-slate-200 px-4 py-12 text-center dark:border-slate-800">
      <UIcon name="i-lucide-coins" class="mx-auto size-8 text-slate-400" />
      <p class="mt-3 text-sm text-slate-500 dark:text-slate-400">{{ $t('user.bonusLogs.empty') }}</p>
    </div>
    <div v-else class="mt-4 divide-y divide-slate-100 dark:divide-slate-800">
      <article v-for="log in bonusLogs" :key="log.id" class="grid gap-3 py-3 text-sm sm:grid-cols-[minmax(0,1fr)_160px] sm:items-center">
        <div class="min-w-0">
          <div class="flex min-w-0 items-center gap-2">
            <UBadge :color="log.amount >= 0 ? 'success' : 'error'" variant="soft">
              {{ log.amount >= 0 ? '+' : '' }}{{ formatBonus(log.amount) }}
            </UBadge>
            <p class="truncate font-medium text-slate-950 dark:text-white">{{ bonusActionLabel(log.action) }}</p>
          </div>
          <p class="mt-1 truncate text-xs text-slate-500 dark:text-slate-400">
            {{ log.remark || `${log.targetType || '-'} #${log.targetId || '-'}` }}
          </p>
        </div>
        <div class="text-xs text-slate-500 sm:text-right dark:text-slate-400">
          <p>{{ formatDateTime(log.createdAt, locale) }}</p>
          <p class="mt-1">{{ $t('user.bonusLogs.after', { balance: formatBonus(log.balanceAfter) }) }}</p>
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
  </UCard>
</template>

<script setup lang="ts">
import { ApiError } from '~/composables/useApi'
import { useEconomy, type BonusLogItem } from '~/composables/useEconomy'
import { formatDateTime } from '~/utils/format'

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

function bonusActionLabel(action?: string | null) {
  if (!action) return '-'
  const key = `user.bonusLogs.actions.${action}`
  const translated = t(key)
  return translated === key ? action : translated
}
</script>
