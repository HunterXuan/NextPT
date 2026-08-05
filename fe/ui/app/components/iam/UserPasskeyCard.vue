<template>
  <UCard
    class="rounded-lg"
    :ui="{
      root: '',
      header: 'px-4 py-3 sm:px-5',
      body: 'p-4 sm:p-5'
    }"
  >
    <template #header>
      <div>
        <h2 class="text-sm font-semibold text-slate-950 dark:text-white">{{ $t('user.passkey.title') }}</h2>
        <p class="mt-0.5 text-xs text-slate-500 dark:text-slate-400">{{ $t('user.passkey.subtitle') }}</p>
      </div>
    </template>

    <div class="flex min-h-10 items-center rounded-md border border-slate-200 bg-slate-50 px-3 py-2 dark:border-slate-800 dark:bg-slate-950">
      <p class="break-all font-mono text-xs text-slate-700 dark:text-slate-300">{{ displayPasskey }}</p>
    </div>
    <div class="grid grid-cols-2 gap-2 pt-3 sm:grid-cols-3">
      <UButton
        type="button"
        color="neutral"
        variant="outline"
        :icon="showPasskey ? 'i-lucide-eye-off' : 'i-lucide-eye'"
        block
        :disabled="!user?.user.passkey"
        @click="showPasskey = !showPasskey"
      >
        {{ showPasskey ? $t('user.passkey.hide') : $t('user.passkey.show') }}
      </UButton>
      <UButton
        type="button"
        color="neutral"
        variant="outline"
        icon="i-lucide-copy"
        block
        :disabled="!user?.user.passkey"
        @click="copyPasskey"
      >
        {{ $t('common.copy') }}
      </UButton>
      <UPopover
        :open="resetConfirmOpen"
        :content="{ side: 'top', align: 'end', sideOffset: 8 }"
        :ui="{ content: 'w-64 p-3' }"
        @update:open="setResetConfirmOpen"
      >
        <UButton
          type="button"
          color="neutral"
          variant="outline"
          icon="i-lucide-refresh-cw"
          block
          :loading="resetPending"
          :disabled="resetPending"
        >
          {{ $t('user.passkey.reset') }}
        </UButton>

        <template #content="{ close }">
          <div class="space-y-3">
            <p class="text-sm font-medium text-slate-950 dark:text-white">
              {{ $t('user.passkey.resetConfirm') }}
            </p>
            <p class="text-xs leading-5 text-slate-500 dark:text-slate-400">
              {{ $t('user.passkey.resetDescription') }}
            </p>
            <div class="flex justify-end gap-2">
              <UButton color="neutral" variant="ghost" size="xs" type="button" @click="closeResetConfirm(close)">
                {{ $t('common.cancel') }}
              </UButton>
              <UButton color="warning" size="xs" type="button" :loading="resetPending" :disabled="resetPending" @click="confirmResetPasskey(close)">
                {{ $t('common.confirm') }}
              </UButton>
            </div>
          </div>
        </template>
      </UPopover>
    </div>
  </UCard>
</template>

<script setup lang="ts">
const { t } = useI18n()
const toast = useToast()
const { user, fetchUser, resetPasskey } = useAuth()
const showPasskey = ref(false)
const resetConfirmOpen = ref(false)
const resetPending = ref(false)

const displayPasskey = computed(() => {
  const passkey = user.value?.user.passkey || ''
  if (!passkey) return '-'
  if (showPasskey.value) return passkey
  return maskPasskey(passkey)
})

onMounted(async () => {
  if (!user.value) {
    await fetchUser()
  }
})

async function copyPasskey() {
  if (!user.value?.user.passkey || !navigator?.clipboard) return

  await navigator.clipboard.writeText(user.value.user.passkey)
  toast.add({
    title: t('user.passkey.copied'),
    color: 'success',
    icon: 'i-lucide-check-circle'
  })
}

function setResetConfirmOpen(open: boolean) {
  if (resetPending.value) return
  resetConfirmOpen.value = open
}

function closeResetConfirm(close?: () => void) {
  resetConfirmOpen.value = false
  close?.()
}

async function confirmResetPasskey(close?: () => void) {
  if (resetPending.value) return

  resetPending.value = true
  try {
    await resetPasskey()
    showPasskey.value = true
    toast.add({
      title: t('user.passkey.resetSuccess'),
      color: 'success',
      icon: 'i-lucide-check-circle'
    })
    closeResetConfirm(close)
  } catch (error) {
    toast.add({
      title: error instanceof Error && error.message ? error.message : t('common.requestFailed'),
      color: 'error',
      icon: 'i-lucide-circle-alert'
    })
  } finally {
    resetPending.value = false
  }
}

function maskPasskey(value: string) {
  if (value.length <= 8) return '*'.repeat(value.length)
  return `${value.slice(0, 4)} ${'*'.repeat(value.length - 8)} ${value.slice(-4)}`
}
</script>
