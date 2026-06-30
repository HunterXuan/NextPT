<template>
  <UPopover
    :open="open"
    :content="{ side: 'top', align: 'start', sideOffset: 8 }"
    :ui="{ content: 'w-64 p-3' }"
    @update:open="setOpen"
  >
    <UButton
      type="button"
      color="neutral"
      variant="ghost"
      size="xs"
      icon="i-lucide-coins"
      :disabled="disabled || pending"
      :aria-label="ariaLabel || t('common.reward.confirm')"
      :title="ariaLabel || t('common.reward.confirm')"
    >
      {{ numberFormatter.format(count) }}
    </UButton>

    <template #content="{ close }">
      <div class="space-y-3">
        <div class="grid grid-cols-3 gap-2">
          <button
            v-for="amount in normalizedPresets"
            :key="amount"
            type="button"
            :class="amountButtonClass(amount)"
            :disabled="pending"
            @click="selectedAmount = amount"
          >
            {{ formatRewardAmount(amount) }}
          </button>
        </div>

        <div v-if="selectedAmount > 0" class="rounded-md bg-slate-50 p-2 dark:bg-slate-950">
          <p class="text-sm font-medium text-slate-950 dark:text-white">
            {{ t(confirmTitleKey, { amount: formatRewardAmount(selectedAmount) }) }}
          </p>
          <div class="mt-3 flex justify-end gap-2">
            <UButton color="neutral" variant="ghost" size="xs" type="button" @click="closePopover(close)">
              {{ t('common.cancel') }}
            </UButton>
            <UButton
              color="primary"
              size="xs"
              type="button"
              icon="i-lucide-coins"
              :loading="pending"
              :disabled="pending"
              @click="confirmReward(close)"
            >
              {{ t(confirmKey) }}
            </UButton>
          </div>
        </div>
      </div>
    </template>
  </UPopover>
</template>

<script setup lang="ts">
type RewardSubmitter = (amount: number) => Promise<void>

const props = withDefaults(defineProps<{
  count?: number
  presets?: number[]
  disabled?: boolean
  ariaLabel?: string
  confirmTitleKey?: string
  confirmKey?: string
  successKey?: string
  submitReward: RewardSubmitter
}>(), {
  count: 0,
  presets: () => [20, 50, 100, 200, 500, 1000],
  disabled: false,
  ariaLabel: '',
  confirmTitleKey: 'common.reward.confirmTitle',
  confirmKey: 'common.reward.confirm',
  successKey: 'common.reward.success'
})

const emit = defineEmits<{
  success: [amount: number]
}>()

const { t, locale } = useI18n()
const toast = useToast()

const open = ref(false)
const pending = ref(false)
const selectedAmount = ref(0)

const numberFormatter = computed(() => new Intl.NumberFormat(locale.value))
const normalizedPresets = computed(() => {
  return Array.from(new Set(props.presets
    .map((amount) => Number(amount))
    .filter((amount) => Number.isFinite(amount) && amount > 0)))
})

function setOpen(value: boolean) {
  if (pending.value) return
  open.value = value
  if (!value) {
    selectedAmount.value = 0
  }
}

function closePopover(close?: () => void) {
  selectedAmount.value = 0
  open.value = false
  close?.()
}

async function confirmReward(close?: () => void) {
  const amount = selectedAmount.value
  if (!Number.isFinite(amount) || amount <= 0 || pending.value) return

  pending.value = true
  try {
    await props.submitReward(amount)
    emit('success', amount)
    toast.add({
      title: t(props.successKey),
      color: 'success',
      icon: 'i-lucide-check-circle'
    })
    closePopover(close)
  } catch (error) {
    toast.add({
      title: errorMessage(error),
      color: 'error',
      icon: 'i-lucide-circle-alert'
    })
  } finally {
    pending.value = false
  }
}

function amountButtonClass(amount: number) {
  const base = 'h-8 rounded-md border px-2 text-sm font-medium transition disabled:cursor-not-allowed disabled:opacity-60'
  if (selectedAmount.value === amount) {
    return `${base} border-primary bg-primary/10 text-primary`
  }
  return `${base} border-slate-200 bg-white text-slate-700 hover:border-sky-200 hover:bg-sky-50 hover:text-sky-700 dark:border-slate-800 dark:bg-slate-900 dark:text-slate-200 dark:hover:border-sky-900 dark:hover:bg-sky-950 dark:hover:text-sky-300`
}

function formatRewardAmount(value: number) {
  return numberFormatter.value.format(Number(value || 0))
}

function errorMessage(error: unknown) {
  if (error instanceof Error && error.message) return error.message
  return t('common.requestFailed')
}
</script>
