<template>
  <UCard class="rounded-lg">
    <template #header>
      <div>
        <h2 class="text-base font-semibold text-slate-950 dark:text-white">{{ $t('user.security.title') }}</h2>
        <p class="mt-1 text-sm text-slate-500 dark:text-slate-400">{{ $t('user.security.subtitle') }}</p>
      </div>
    </template>

    <form class="grid grid-cols-1 gap-4 md:grid-cols-2" @submit.prevent="handlePasswordChange">
      <UFormField :label="$t('user.security.oldPassword')" required>
        <UInput
          v-model="passwordForm.oldPassword"
          class="w-full"
          type="password"
          autocomplete="current-password"
          :disabled="passwordPending"
        />
      </UFormField>

      <UFormField :label="$t('user.security.newPassword')" required>
        <UInput
          v-model="passwordForm.newPassword"
          class="w-full"
          type="password"
          autocomplete="new-password"
          :disabled="passwordPending"
        />
      </UFormField>

      <div class="flex justify-end md:col-span-2">
        <UButton type="submit" color="neutral" variant="outline" :loading="passwordPending" :disabled="!canChangePassword">
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
const passwordForm = reactive({
  oldPassword: '',
  newPassword: ''
})

const canChangePassword = computed(() => passwordForm.oldPassword.length > 0 && passwordForm.newPassword.length >= 6)

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
