<template>
  <UCard
    class="rounded-lg"
    :ui="{
      header: 'px-4 py-3 sm:px-5',
      body: 'p-4 sm:p-5'
    }"
  >
    <template #header>
      <div>
        <h2 class="text-sm font-semibold text-slate-950 dark:text-white">{{ $t('user.security.title') }}</h2>
        <p class="mt-0.5 text-xs text-slate-500 dark:text-slate-400">{{ $t('user.security.subtitle') }}</p>
      </div>
    </template>

    <form class="grid gap-4" @submit.prevent="handlePasswordChange">
      <UFormField :label="$t('user.security.oldPassword')" required>
        <UInput
          v-model="passwordForm.oldPassword"
          class="w-full"
          icon="i-lucide-lock-keyhole"
          :type="showOldPassword ? 'text' : 'password'"
          autocomplete="current-password"
          :disabled="passwordPending"
        >
          <template #trailing>
            <UButton
              type="button"
              color="neutral"
              variant="ghost"
              size="xs"
              :icon="showOldPassword ? 'i-lucide-eye-off' : 'i-lucide-eye'"
              :aria-label="showOldPassword ? $t('common.hidePassword') : $t('common.showPassword')"
              :disabled="passwordPending"
              @click="showOldPassword = !showOldPassword"
            />
          </template>
        </UInput>
      </UFormField>

      <div class="grid gap-4 md:grid-cols-2">
        <UFormField :label="$t('user.security.newPassword')" required>
          <UInput
            v-model="passwordForm.newPassword"
            class="w-full"
            icon="i-lucide-key-round"
            :type="showNewPassword ? 'text' : 'password'"
            autocomplete="new-password"
            :disabled="passwordPending"
          >
            <template #trailing>
              <UButton
                type="button"
                color="neutral"
                variant="ghost"
                size="xs"
                :icon="showNewPassword ? 'i-lucide-eye-off' : 'i-lucide-eye'"
                :aria-label="showNewPassword ? $t('common.hidePassword') : $t('common.showPassword')"
                :disabled="passwordPending"
                @click="showNewPassword = !showNewPassword"
              />
            </template>
          </UInput>
        </UFormField>

        <UFormField :label="$t('user.security.confirmPassword')" required>
          <UInput
            v-model="passwordForm.confirmPassword"
            class="w-full"
            icon="i-lucide-key-round"
            :type="showConfirmPassword ? 'text' : 'password'"
            autocomplete="new-password"
            :disabled="passwordPending"
          >
            <template #trailing>
              <UButton
                type="button"
                color="neutral"
                variant="ghost"
                size="xs"
                :icon="showConfirmPassword ? 'i-lucide-eye-off' : 'i-lucide-eye'"
                :aria-label="showConfirmPassword ? $t('common.hidePassword') : $t('common.showPassword')"
                :disabled="passwordPending"
                @click="showConfirmPassword = !showConfirmPassword"
              />
            </template>
          </UInput>
        </UFormField>
      </div>

      <UAlert
        v-if="passwordMismatch"
        color="warning"
        variant="soft"
        icon="i-lucide-circle-alert"
        :title="$t('user.security.passwordMismatch')"
      />

      <div class="flex justify-end">
        <UButton type="submit" color="primary" icon="i-lucide-key-round" :loading="passwordPending" :disabled="!canChangePassword">
          {{ $t('user.security.submit') }}
        </UButton>
      </div>
    </form>
  </UCard>
</template>

<script setup lang="ts">
import { ApiError } from '~/composables/useApi'

const { t } = useI18n()
const localePath = useLocalePath()
const toast = useToast()
const { changePassword } = useAuth()

const passwordPending = ref(false)
const showOldPassword = ref(false)
const showNewPassword = ref(false)
const showConfirmPassword = ref(false)
const passwordForm = reactive({
  oldPassword: '',
  newPassword: '',
  confirmPassword: ''
})

const passwordMismatch = computed(() => {
  return passwordForm.confirmPassword.length > 0 && passwordForm.newPassword !== passwordForm.confirmPassword
})
const canChangePassword = computed(() => {
  return passwordForm.oldPassword.length > 0
    && passwordForm.newPassword.length >= 6
    && passwordForm.newPassword === passwordForm.confirmPassword
})

async function handlePasswordChange() {
  if (!canChangePassword.value) return

  passwordPending.value = true

  try {
    await changePassword({
      oldPassword: passwordForm.oldPassword,
      newPassword: passwordForm.newPassword
    })
    toast.add({
      title: t('user.security.changed'),
      color: 'success',
      icon: 'i-lucide-check-circle'
    })
    await navigateTo(localePath('/login'))
  } catch (error) {
    toast.add({
      title: error instanceof ApiError ? error.message : t('common.requestFailed'),
      color: 'error',
      icon: 'i-lucide-circle-alert'
    })
  } finally {
    passwordPending.value = false
  }
}
</script>
