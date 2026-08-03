<template>
  <section class="rounded-lg border border-slate-200 bg-white p-4 dark:border-slate-800 dark:bg-slate-900">
    <div class="flex items-start justify-between gap-3">
      <div>
        <h2 class="text-sm font-semibold text-slate-950 dark:text-white">{{ title }}</h2>
        <p v-if="summaryText" class="mt-1 text-xs text-slate-500 dark:text-slate-400">
          {{ summaryText }}
        </p>
      </div>
      <UIcon :name="icon" class="mt-0.5 size-4 shrink-0" :class="iconClass" />
    </div>

    <div class="mt-3 grid grid-cols-3 gap-2">
      <UPopover
        v-for="amount in normalizedPresets"
        :key="amount"
        :open="activeAmount === amount"
        :content="{ side: 'top', align: 'center', sideOffset: 8 }"
        :ui="{ content: 'w-56 p-3' }"
        @update:open="setRewardPopover(amount, $event)"
      >
        <UButton
          type="button"
          size="xs"
          class="w-full justify-center"
          :color="activeAmount === amount ? 'primary' : 'neutral'"
          :variant="activeAmount === amount ? 'soft' : 'outline'"
          :disabled="rewardPending"
        >
          {{ formatRewardAmount(amount) }}
        </UButton>

        <template #content="{ close }">
          <div class="space-y-3">
            <p class="text-sm font-medium text-slate-950 dark:text-white">
              {{ t(confirmTitleKey, { amount: formatRewardAmount(amount) }) }}
            </p>
            <div class="flex justify-end gap-2">
              <UButton color="neutral" variant="ghost" size="xs" type="button" @click="closeRewardPopover(close)">
                {{ t('common.cancel') }}
              </UButton>
              <UButton
                color="primary"
                size="xs"
                type="button"
                icon="i-lucide-coins"
                :loading="rewardPending && activeAmount === amount"
                :disabled="rewardPending"
                @click="confirmReward(amount, close)"
              >
                {{ t(confirmKey) }}
              </UButton>
            </div>
          </div>
        </template>
      </UPopover>
    </div>

    <template v-if="hasRewardHistory">
      <div v-if="previewError" class="mt-3 rounded-md border border-red-200 bg-red-50 px-3 py-2 text-sm text-red-700 dark:border-red-900 dark:bg-red-950 dark:text-red-200">
        {{ previewError }}
      </div>
      <div v-else-if="previewPending && previewRewards.length === 0" class="mt-4 space-y-2">
        <div v-for="item in 3" :key="item" class="h-9 animate-pulse rounded-md bg-slate-100 dark:bg-slate-800" />
      </div>
      <div v-else-if="previewRewards.length === 0" class="mt-3 rounded-md border border-dashed border-slate-200 px-3 py-2 text-xs text-slate-500 dark:border-slate-800 dark:text-slate-400">
        {{ t(emptyKey) }}
      </div>
      <div v-else class="mt-4 max-h-72 divide-y divide-slate-100 overflow-auto dark:divide-slate-800">
        <div v-for="(reward, index) in previewRewards" :key="rewardKey(reward, index)" class="grid grid-cols-[minmax(0,1fr)_auto] gap-3 py-3 text-sm">
          <div class="flex min-w-0 items-center gap-2">
            <IamUserAvatar :user="reward.user" :alt="rewardDisplayName(reward)" size="xs" />
            <div class="min-w-0">
              <IamUserPopover :user="reward.user" :fallback="rewardDisplayName(reward)" class="truncate font-medium text-slate-950 dark:text-white" />
              <p class="mt-0.5 truncate text-xs text-slate-500 dark:text-slate-400">{{ rewardMetaText(reward) }}</p>
            </div>
          </div>
          <span class="self-center whitespace-nowrap text-sm font-semibold text-amber-600 dark:text-amber-300">{{ formatRewardAmount(reward.amount) }}</span>
        </div>
      </div>

      <div v-if="rewardTotal > previewRewards.length && previewRewards.length > 0" class="mt-3 flex items-center justify-between gap-3 border-t border-slate-100 pt-3 dark:border-slate-800">
        <p class="min-w-0 text-xs text-slate-500 dark:text-slate-400">
          {{ t(previewKey, { count: numberFormatter.format(previewRewards.length), total: numberFormatter.format(rewardTotal) }) }}
        </p>
        <UButton color="neutral" variant="ghost" size="xs" trailing-icon="i-lucide-arrow-right" @click="setModalOpen(true)">
          {{ t(viewAllKey) }}
        </UButton>
      </div>
    </template>
  </section>

  <UModal
    v-if="hasRewardHistory"
    :open="modalOpen"
    :title="allTitleText"
    :description="allDescriptionText"
    :ui="{ content: 'sm:max-w-xl' }"
    @update:open="setModalOpen"
  >
    <template #body>
      <div class="px-4 pb-4 sm:px-6 sm:pb-6">
        <div v-if="modalError" class="rounded-md border border-red-200 bg-red-50 px-3 py-2 text-sm text-red-700 dark:border-red-900 dark:bg-red-950 dark:text-red-200">
          {{ modalError }}
        </div>
        <div v-else-if="modalPending && modalRewards.length === 0" class="space-y-2">
          <div v-for="item in 5" :key="item" class="h-12 animate-pulse rounded-md bg-slate-100 dark:bg-slate-800" />
        </div>
        <div v-else-if="modalRewards.length === 0" class="rounded-md border border-dashed border-slate-200 px-4 py-10 text-center text-sm text-slate-500 dark:border-slate-800 dark:text-slate-400">
          {{ t(emptyKey) }}
        </div>
        <div v-else class="divide-y divide-slate-100 dark:divide-slate-800">
          <div v-for="(reward, index) in modalRewards" :key="rewardKey(reward, index)" class="grid grid-cols-[minmax(0,1fr)_auto] gap-3 py-3 text-sm">
            <div class="flex min-w-0 items-center gap-2">
              <IamUserAvatar :user="reward.user" :alt="rewardDisplayName(reward)" size="xs" />
              <div class="min-w-0">
                <IamUserPopover :user="reward.user" :fallback="rewardDisplayName(reward)" class="truncate font-medium text-slate-950 dark:text-white" />
                <p class="mt-0.5 truncate text-xs text-slate-500 dark:text-slate-400">{{ rewardMetaText(reward) }}</p>
              </div>
            </div>
            <span class="self-center whitespace-nowrap text-sm font-semibold text-amber-600 dark:text-amber-300">{{ formatRewardAmount(reward.amount) }}</span>
          </div>
        </div>

        <AppPager
          v-if="rewardTotal > modalSize"
          class="mt-4 border-t border-slate-200 pt-4 dark:border-slate-800"
          size="xs"
          :page="modalPage"
          :page-size="modalSize"
          :total="rewardTotal"
          :disabled="modalPending"
          @page-change="goToModalPage"
        />
      </div>
    </template>
  </UModal>
</template>

<script setup lang="ts">
import { formatDateOnly } from '~/utils/format'

interface RewardPanelItem {
  user?: {
    id?: number
    username?: string
  }
  amount: number
  rewardCount?: number
  lastRewardAt?: string
}

interface RewardListResult {
  list?: RewardPanelItem[]
  total?: number
}

type RewardLoader = (page: number, size: number) => Promise<RewardListResult>
type RewardSubmitter = (amount: number) => Promise<void>

const props = withDefaults(defineProps<{
  title: string
  sourceKey?: string | number
  icon?: string
  iconClass?: string
  presets?: number[]
  previewSize?: number
  modalSize?: number
  showHistory?: boolean
  summaryKey?: string
  confirmTitleKey?: string
  confirmKey?: string
  successKey?: string
  emptyKey?: string
  previewKey?: string
  viewAllKey?: string
  allTitleKey?: string
  allDescriptionKey?: string
  loadRewards?: RewardLoader
  submitReward: RewardSubmitter
}>(), {
  icon: 'i-lucide-coins',
  iconClass: 'text-amber-500',
  presets: () => [20, 50, 100, 200, 500, 1000],
  previewSize: 5,
  modalSize: 10,
  showHistory: true,
  summaryKey: 'common.reward.summary',
  confirmTitleKey: 'common.reward.confirmTitle',
  confirmKey: 'common.reward.confirm',
  successKey: 'common.reward.success',
  emptyKey: 'common.reward.empty',
  previewKey: 'common.reward.preview',
  viewAllKey: 'common.reward.viewAll',
  allTitleKey: 'common.reward.allTitle',
  allDescriptionKey: 'common.reward.allDescription'
})

const { t, locale } = useI18n()
const toast = useToast()

const previewRewards = ref<RewardPanelItem[]>([])
const modalRewards = ref<RewardPanelItem[]>([])
const rewardTotal = ref(0)
const activeAmount = ref(0)
const rewardPending = ref(false)
const previewPending = ref(false)
const modalPending = ref(false)
const previewError = ref('')
const modalError = ref('')
const modalOpen = ref(false)
const modalPage = ref(1)
let previewLoadToken = 0
let modalLoadToken = 0

const numberFormatter = computed(() => new Intl.NumberFormat(locale.value))
const hasRewardHistory = computed(() => props.showHistory && typeof props.loadRewards === 'function')
const normalizedPresets = computed(() => {
  return Array.from(new Set(props.presets
    .map((amount) => Number(amount))
    .filter((amount) => Number.isFinite(amount) && amount > 0)))
})
const summaryText = computed(() => {
  if (!props.summaryKey || !hasRewardHistory.value) return ''
  return t(props.summaryKey, { count: numberFormatter.value.format(rewardTotal.value) })
})
const allTitleText = computed(() => t(props.allTitleKey))
const allDescriptionText = computed(() => t(props.allDescriptionKey, { count: numberFormatter.value.format(rewardTotal.value) }))

onMounted(() => {
  void loadPreviewRewards()
})

watch(() => props.sourceKey, () => {
  resetRewardHistory()
  void loadPreviewRewards()
})

watch(hasRewardHistory, (enabled) => {
  resetRewardHistory()
  if (enabled) {
    void loadPreviewRewards()
  }
})

function resetRewardHistory() {
  previewLoadToken += 1
  modalLoadToken += 1
  previewRewards.value = []
  modalRewards.value = []
  rewardTotal.value = 0
  previewPending.value = false
  modalPending.value = false
  previewError.value = ''
  modalError.value = ''
  modalPage.value = 1
  modalOpen.value = false
}

async function loadPreviewRewards(force = false) {
  const loader = props.loadRewards
  if (!hasRewardHistory.value || !loader || (previewPending.value && !force)) return

  const token = previewLoadToken + 1
  previewLoadToken = token
  previewPending.value = true
  previewError.value = ''
  try {
    const data = await loader(1, props.previewSize)
    if (token !== previewLoadToken) return
    previewRewards.value = data.list || []
    rewardTotal.value = data.total || 0
  } catch (error) {
    if (token !== previewLoadToken) return
    previewRewards.value = []
    rewardTotal.value = 0
    previewError.value = errorMessage(error)
  } finally {
    if (token === previewLoadToken) {
      previewPending.value = false
    }
  }
}

function setModalOpen(open: boolean) {
  modalOpen.value = open
  if (!open) return

  modalPage.value = 1
  void loadModalPage(1)
}

function goToModalPage(page: number) {
  void loadModalPage(page)
}

async function loadModalPage(page = modalPage.value, force = false) {
  const loader = props.loadRewards
  if (!hasRewardHistory.value || !loader || (modalPending.value && !force)) return

  const token = modalLoadToken + 1
  modalLoadToken = token
  modalPage.value = Math.max(1, page)
  modalPending.value = true
  modalError.value = ''
  try {
    const data = await loader(modalPage.value, props.modalSize)
    if (token !== modalLoadToken) return
    modalRewards.value = data.list || []
    rewardTotal.value = data.total || 0

    const totalPages = Math.max(1, Math.ceil(rewardTotal.value / props.modalSize))
    if (modalPage.value > totalPages) {
      modalPage.value = totalPages
      const nextData = await loader(modalPage.value, props.modalSize)
      if (token !== modalLoadToken) return
      modalRewards.value = nextData.list || []
      rewardTotal.value = nextData.total || 0
    }
  } catch (error) {
    if (token !== modalLoadToken) return
    modalRewards.value = []
    modalError.value = errorMessage(error)
  } finally {
    if (token === modalLoadToken) {
      modalPending.value = false
    }
  }
}

function setRewardPopover(amount: number, open: boolean) {
  if (rewardPending.value) return
  if (open) {
    activeAmount.value = amount
  } else if (activeAmount.value === amount) {
    activeAmount.value = 0
  }
}

function closeRewardPopover(close?: () => void) {
  activeAmount.value = 0
  close?.()
}

async function confirmReward(amount: number, close?: () => void) {
  if (!Number.isFinite(amount) || amount <= 0 || rewardPending.value) return

  activeAmount.value = amount
  rewardPending.value = true
  try {
    await props.submitReward(amount)
    toast.add({
      title: t(props.successKey),
      color: 'success',
      icon: 'i-lucide-check-circle'
    })
    await loadPreviewRewards(true)
    if (modalOpen.value) {
      await loadModalPage(modalPage.value, true)
    }
    closeRewardPopover(close)
  } catch (error) {
    toast.add({
      title: errorMessage(error),
      color: 'error',
      icon: 'i-lucide-circle-alert'
    })
  } finally {
    rewardPending.value = false
  }
}

function rewardDisplayName(reward: RewardPanelItem) {
  return reward.user?.username || `#${reward.user?.id || 0}`
}

function rewardMetaText(reward: RewardPanelItem) {
  const count = Math.max(1, Number(reward.rewardCount || 0))
  const countText = count === 1
    ? t('common.reward.once')
    : t('common.reward.times', { count: numberFormatter.value.format(count) })
  if (!reward.lastRewardAt) return countText
  return `${countText} · ${formatDateOnly(reward.lastRewardAt, locale.value)}`
}

function rewardKey(reward: RewardPanelItem, index: number) {
  return `${reward.user?.id || 0}-${reward.lastRewardAt || 'none'}-${index}`
}

function formatRewardAmount(value: number) {
  return numberFormatter.value.format(Number(value || 0))
}

function errorMessage(error: unknown) {
  if (error instanceof Error && error.message) return error.message
  return t('common.requestFailed')
}
</script>
