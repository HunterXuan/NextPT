<template>
  <UCard class="rounded-lg">
    <template #header>
      <div>
        <h2 class="text-base font-semibold text-slate-950 dark:text-white">{{ $t('user.profile.title') }}</h2>
        <p class="mt-1 text-sm text-slate-500 dark:text-slate-400">{{ $t('user.profile.subtitle') }}</p>
      </div>
    </template>

    <form class="grid grid-cols-1 gap-4" @submit.prevent="handleProfileSave">
      <UFormField :label="$t('user.profile.avatar')">
        <UInput
          v-model="profileForm.avatar"
          class="w-full"
          icon="i-lucide-image"
          :disabled="profilePending || loading"
        />
      </UFormField>

      <UFormField :label="$t('user.profile.signature')">
        <UInput
          v-model="profileForm.signature"
          class="w-full"
          icon="i-lucide-pen-line"
          :disabled="profilePending || loading"
        />
      </UFormField>

      <UFormField :label="$t('user.profile.info')">
        <UTextarea
          v-model="profileForm.info"
          class="w-full"
          :rows="6"
          :disabled="profilePending || loading"
        />
      </UFormField>

      <div class="flex justify-end">
        <UButton type="submit" color="primary" :loading="profilePending" :disabled="loading">
          {{ $t('common.save') }}
        </UButton>
      </div>
    </form>
  </UCard>
</template>

<script setup lang="ts">
import { ApiError } from '~/composables/useApi'

const { t } = useI18n()
const toast = useToast()
const { user, fetchUser, updateProfile } = useAuth()

const loading = ref(true)
const profilePending = ref(false)

const profileForm = reactive({
  avatar: '',
  signature: '',
  info: ''
})

watch(
  user,
  (value) => {
    profileForm.avatar = value?.avatar || ''
    profileForm.signature = value?.signature || ''
    profileForm.info = value?.info || ''
  },
  { immediate: true }
)

onMounted(async () => {
  loading.value = true
  try {
    await fetchUser()
  } finally {
    loading.value = false
  }
})

async function handleProfileSave() {
  profilePending.value = true

  try {
    await updateProfile({
      avatar: profileForm.avatar.trim(),
      signature: profileForm.signature.trim(),
      info: profileForm.info.trim()
    })
    toast.add({
      title: t('user.profile.saved'),
      color: 'success',
      icon: 'i-lucide-check-circle'
    })
  } catch (error) {
    toast.add({
      title: error instanceof ApiError ? error.message : t('common.requestFailed'),
      color: 'error',
      icon: 'i-lucide-circle-alert'
    })
  } finally {
    profilePending.value = false
  }
}
</script>
