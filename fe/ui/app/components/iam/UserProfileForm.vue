<template>
  <UCard class="rounded-lg">
    <template #header>
      <div>
        <h2 class="text-base font-semibold text-slate-950 dark:text-white">{{ $t('user.profile.title') }}</h2>
        <p class="mt-1 text-sm text-slate-500 dark:text-slate-400">{{ $t('user.profile.subtitle') }}</p>
      </div>
    </template>

    <form class="grid grid-cols-1 gap-4" @submit.prevent="handleProfileSave">
      <IamUserAvatar
        :id="user?.id"
        :username="user?.username"
        :avatar="profileForm.avatar"
        :alt="user?.username || $t('user.profile.avatar')"
        size="xl"
        class="ring-1 ring-slate-200 dark:ring-slate-800"
      />

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
        <UButton type="submit" color="primary" :loading="profilePending" :disabled="loading || !profileChanged">
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
const originalProfile = reactive({
  avatar: '',
  signature: '',
  info: ''
})

const normalizedProfile = computed(() => ({
  avatar: profileForm.avatar.trim(),
  signature: profileForm.signature.trim(),
  info: profileForm.info.trim()
}))

const profileChanged = computed(() => {
  return normalizedProfile.value.avatar !== originalProfile.avatar
    || normalizedProfile.value.signature !== originalProfile.signature
    || normalizedProfile.value.info !== originalProfile.info
})

watch(
  user,
  (value) => {
    const nextProfile = {
      avatar: value?.avatar || '',
      signature: value?.signature || '',
      info: value?.info || ''
    }

    profileForm.avatar = nextProfile.avatar
    profileForm.signature = nextProfile.signature
    profileForm.info = nextProfile.info
    originalProfile.avatar = nextProfile.avatar.trim()
    originalProfile.signature = nextProfile.signature.trim()
    originalProfile.info = nextProfile.info.trim()
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
  if (!profileChanged.value) return

  profilePending.value = true

  try {
    await updateProfile(normalizedProfile.value)
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
