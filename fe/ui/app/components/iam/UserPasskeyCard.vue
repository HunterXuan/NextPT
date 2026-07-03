<template>
  <UCard class="rounded-lg">
    <template #header>
      <div>
        <h2 class="text-base font-semibold text-slate-950 dark:text-white">{{ $t('user.passkey.title') }}</h2>
        <p class="mt-1 text-sm text-slate-500 dark:text-slate-400">{{ $t('user.passkey.subtitle') }}</p>
      </div>
    </template>

    <div class="rounded-md border border-slate-200 bg-slate-50 p-3 dark:border-slate-800 dark:bg-slate-950">
      <p class="break-all font-mono text-xs text-slate-700 dark:text-slate-300">{{ user?.passkey || '-' }}</p>
    </div>
    <UButton
      class="mt-3"
      color="neutral"
      variant="outline"
      icon="i-lucide-copy"
      block
      :disabled="!user?.passkey"
      @click="copyPasskey"
    >
      {{ $t('common.copy') }}
    </UButton>
  </UCard>
</template>

<script setup lang="ts">
const { t } = useI18n()
const toast = useToast()
const { user, fetchUser } = useAuth()

onMounted(async () => {
  if (!user.value) {
    await fetchUser()
  }
})

async function copyPasskey() {
  if (!user.value?.passkey || !navigator?.clipboard) return

  await navigator.clipboard.writeText(user.value.passkey)
  toast.add({
    title: t('user.passkey.copied'),
    color: 'success',
    icon: 'i-lucide-check-circle'
  })
}
</script>
