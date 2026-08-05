<template>
  <UCard
    class="rounded-lg"
    :ui="{
      body: 'p-4 sm:p-4'
    }"
  >
    <div class="flex min-h-10 items-center justify-between gap-4">
      <div class="flex min-w-0 items-center gap-3">
        <span class="flex size-9 shrink-0 items-center justify-center rounded-md bg-sky-50 text-sky-700 dark:bg-sky-950 dark:text-sky-300">
          <UIcon name="i-lucide-shield-check" class="size-4" />
        </span>
        <div class="flex min-w-0 items-center gap-2.5">
          <h2 class="truncate text-sm font-semibold text-slate-950 dark:text-white">{{ t('user.twoStep.title') }}</h2>
          <UBadge :color="enabled ? 'success' : 'neutral'" variant="subtle" size="sm">
            {{ enabled ? t('user.twoStep.status.enabled') : t('user.twoStep.status.disabled') }}
          </UBadge>
        </div>
      </div>

      <UButton v-if="!enabled" color="primary" variant="soft" size="sm" icon="i-lucide-plus" @click="setupOpen = true">
        {{ t('user.twoStep.enable') }}
      </UButton>
      <div v-else class="flex shrink-0 items-center gap-1">
        <UTooltip :text="t('user.twoStep.regenerate')">
          <UButton color="neutral" variant="ghost" size="sm" icon="i-lucide-key-round" :aria-label="t('user.twoStep.regenerate')" @click="regenerateOpen = true" />
        </UTooltip>
        <UButton color="error" variant="soft" size="sm" icon="i-lucide-shield-off" @click="disableOpen = true">{{ t('user.twoStep.disable') }}</UButton>
      </div>
    </div>
  </UCard>

  <UModal :open="setupOpen" :title="setupData ? t('user.twoStep.setupTitle') : t('user.twoStep.enableTitle')" :content="setupModalContent" :ui="compactModalUi" @update:open="setSetupOpen">
    <template #body>
      <div v-if="!setupData" class="px-4 py-3">
        <UFormField size="sm" :label="t('user.twoStep.password')" required :error="setupError || undefined">
          <UInput v-model="setupPassword" class="w-full" icon="i-lucide-lock-keyhole" type="password" autocomplete="current-password" :disabled="setupPending" @keyup.enter="startSetup" />
        </UFormField>
      </div>
      <div v-else class="grid gap-4 px-4 py-4 sm:grid-cols-[112px_minmax(0,1fr)] sm:items-center">
        <div class="flex items-center justify-center rounded-md border border-slate-200 bg-slate-50 p-2 dark:border-slate-800 dark:bg-slate-950">
          <img :src="setupData.qrCodeDataUrl" :alt="t('user.twoStep.qrCode')" class="size-24 bg-white">
        </div>
        <div class="grid gap-3">
          <UFormField size="sm" :label="t('user.twoStep.manualKey')">
            <UInput :model-value="setupData.secret" readonly class="w-full font-mono text-xs" icon="i-lucide-key-round">
              <template #trailing>
                <UTooltip :text="t('common.copy')">
                  <UButton color="neutral" variant="ghost" size="xs" icon="i-lucide-copy" :aria-label="t('common.copy')" @click="copyText(setupData.secret)" />
                </UTooltip>
              </template>
            </UInput>
          </UFormField>
          <UFormField size="sm" :label="t('user.twoStep.authenticatorCode')" required :error="setupError || undefined">
            <UInput v-model="setupCode" class="w-full" icon="i-lucide-smartphone" inputmode="numeric" autocomplete="one-time-code" :disabled="confirmPending" @keyup.enter="confirmSetup" />
          </UFormField>
        </div>
      </div>
    </template>
    <template #footer>
      <div class="flex w-full justify-end">
        <UButton v-if="!setupData" color="primary" size="sm" trailing-icon="i-lucide-arrow-right" :loading="setupPending" :disabled="!setupPassword" @click="startSetup">{{ t('user.twoStep.continue') }}</UButton>
        <UButton v-else color="primary" size="sm" icon="i-lucide-shield-check" :loading="confirmPending" :disabled="!setupCode" @click="confirmSetup">{{ t('user.twoStep.confirm') }}</UButton>
      </div>
    </template>
  </UModal>

  <UModal :open="regenerateOpen" :title="t('user.twoStep.regenerate')" :content="compactModalContent" :ui="compactModalUi" @update:open="setRegenerateOpen">
    <template #body>
      <div class="px-4 py-3">
        <UFormField size="sm" :label="t('user.twoStep.code')" required :error="regenerateError || undefined">
          <UInput v-model="regenerateCode" class="w-full" icon="i-lucide-shield-keyhole" inputmode="numeric" autocomplete="one-time-code" :disabled="regeneratePending" @keyup.enter="regenerateRecoveryCodes" />
        </UFormField>
      </div>
    </template>
    <template #footer>
      <div class="flex w-full justify-end">
        <UButton color="primary" size="sm" icon="i-lucide-refresh-cw" :loading="regeneratePending" :disabled="!regenerateCode" @click="regenerateRecoveryCodes">{{ t('user.twoStep.regenerate') }}</UButton>
      </div>
    </template>
  </UModal>

  <UModal :open="disableOpen" :title="t('user.twoStep.disableTitle')" :content="compactModalContent" :ui="compactModalUi" @update:open="setDisableOpen">
    <template #body>
      <div class="grid gap-3 px-4 py-3">
        <UFormField size="sm" :label="t('user.twoStep.password')" required>
          <UInput v-model="disablePassword" class="w-full" icon="i-lucide-lock-keyhole" type="password" autocomplete="current-password" :disabled="disablePending" />
        </UFormField>
        <UFormField size="sm" :label="t('user.twoStep.code')" required :error="disableError || undefined">
          <UInput v-model="disableCode" class="w-full" icon="i-lucide-shield-keyhole" inputmode="numeric" autocomplete="one-time-code" :disabled="disablePending" @keyup.enter="disableTwoStepVerification" />
        </UFormField>
      </div>
    </template>
    <template #footer>
      <div class="flex w-full justify-end">
        <UButton color="error" size="sm" icon="i-lucide-shield-off" :loading="disablePending" :disabled="!disablePassword || !disableCode" @click="disableTwoStepVerification">{{ t('user.twoStep.disable') }}</UButton>
      </div>
    </template>
  </UModal>

  <UModal :open="recoveryOpen" :title="t('user.twoStep.recoveryTitle')" :content="recoveryModalContent" :ui="compactModalUi" @update:open="recoveryOpen = $event">
    <template #body>
      <div class="px-4 py-3">
        <div class="grid grid-cols-2 gap-px overflow-hidden rounded-md border border-slate-200 bg-slate-200 text-center font-mono text-sm dark:border-slate-700 dark:bg-slate-700">
          <span v-for="code in recoveryCodes" :key="code" class="bg-white px-3 py-2.5 dark:bg-slate-900">{{ code }}</span>
        </div>
      </div>
    </template>
    <template #footer>
      <div class="flex w-full items-center justify-between gap-3">
        <UTooltip :text="t('user.twoStep.copyCodes')">
          <UButton color="neutral" variant="ghost" size="sm" icon="i-lucide-copy" :aria-label="t('user.twoStep.copyCodes')" @click="copyText(recoveryCodes.join('\n'))" />
        </UTooltip>
        <UButton color="primary" size="sm" @click="finishRecoveryCodes">{{ t('user.twoStep.done') }}</UButton>
      </div>
    </template>
  </UModal>
</template>

<script setup lang="ts">
import type { AuthTwoStepSetupOut } from '~/composables/useAuth'
import { ApiError } from '~/composables/useApi'

const { t } = useI18n()
const toast = useToast()
const { user, setupTwoStep, confirmTwoStep, createTwoStepRecoveryCodes, disableTwoStep } = useAuth()

const compactModalContent = {
  style: {
    width: 'min(calc(100vw - 2rem), 22rem)',
    maxWidth: '22rem'
  }
}
const recoveryModalContent = {
  style: {
    width: 'min(calc(100vw - 2rem), 26rem)',
    maxWidth: '26rem'
  }
}
const compactModalUi = {
  header: 'min-h-0 p-0 sm:p-0 px-3 py-2.5 sm:px-3 sm:py-2.5',
  title: 'text-sm font-semibold',
  close: 'top-1.5 end-1.5',
  body: 'p-0 sm:p-0',
  footer: 'p-0 sm:p-0 px-3 py-2.5 sm:px-3 sm:py-2.5'
}

const enabled = computed(() => Boolean(user.value?.user.twoStepEnabled))
const setupOpen = ref(false)
const regenerateOpen = ref(false)
const disableOpen = ref(false)
const recoveryOpen = ref(false)
const setupPassword = ref('')
const setupCode = ref('')
const setupData = ref<AuthTwoStepSetupOut | null>(null)
const setupPending = ref(false)
const confirmPending = ref(false)
const regenerateCode = ref('')
const regeneratePending = ref(false)
const disablePassword = ref('')
const disableCode = ref('')
const disablePending = ref(false)
const setupError = ref('')
const regenerateError = ref('')
const disableError = ref('')
const recoveryCodes = ref<string[]>([])
const setupModalContent = computed(() => ({
  style: {
    width: setupData.value ? 'min(calc(100vw - 2rem), 30rem)' : 'min(calc(100vw - 2rem), 22rem)',
    maxWidth: setupData.value ? '30rem' : '22rem'
  }
}))

function setSetupOpen(value: boolean) {
  setupOpen.value = value
  if (!value) resetSetup()
}

function resetSetup() {
  setupPassword.value = ''
  setupCode.value = ''
  setupData.value = null
  setupError.value = ''
}

function setRegenerateOpen(value: boolean) {
  regenerateOpen.value = value
  if (!value) {
    regenerateCode.value = ''
    regenerateError.value = ''
  }
}

function setDisableOpen(value: boolean) {
  disableOpen.value = value
  if (!value) {
    disablePassword.value = ''
    disableCode.value = ''
    disableError.value = ''
  }
}

async function startSetup() {
  setupPending.value = true
  setupError.value = ''
  try {
    setupData.value = await setupTwoStep(setupPassword.value)
  } catch (error) {
    setupError.value = errorMessage(error)
  } finally {
    setupPending.value = false
  }
}

async function confirmSetup() {
  if (!setupData.value) return
  confirmPending.value = true
  setupError.value = ''
  try {
    const data = await confirmTwoStep(setupData.value.challenge, setupCode.value)
    recoveryCodes.value = data.recoveryCodes
    setupOpen.value = false
    resetSetup()
    recoveryOpen.value = true
  } catch (error) {
    setupError.value = errorMessage(error)
  } finally {
    confirmPending.value = false
  }
}

async function regenerateRecoveryCodes() {
  regeneratePending.value = true
  regenerateError.value = ''
  try {
    const data = await createTwoStepRecoveryCodes(regenerateCode.value)
    recoveryCodes.value = data.recoveryCodes
    setRegenerateOpen(false)
    recoveryOpen.value = true
  } catch (error) {
    regenerateError.value = errorMessage(error)
  } finally {
    regeneratePending.value = false
  }
}

async function disableTwoStepVerification() {
  disablePending.value = true
  disableError.value = ''
  try {
    await disableTwoStep(disablePassword.value, disableCode.value)
    setDisableOpen(false)
  } catch (error) {
    disableError.value = errorMessage(error)
  } finally {
    disablePending.value = false
  }
}

async function copyText(value: string) {
  if (typeof navigator === 'undefined' || !navigator.clipboard || !value) return
  await navigator.clipboard.writeText(value)
  toast.add({ title: t('common.copy'), color: 'success', icon: 'i-lucide-check' })
}

function finishRecoveryCodes() {
  recoveryOpen.value = false
  recoveryCodes.value = []
}

function errorMessage(error: unknown) {
  return error instanceof ApiError ? error.message : t('common.requestFailed')
}
</script>
