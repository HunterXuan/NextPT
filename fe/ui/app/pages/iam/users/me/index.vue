<template>
  <div class="min-h-[calc(100vh-4rem)] bg-slate-50 py-6 dark:bg-slate-950">
    <div class="w-full px-3 sm:px-4 lg:px-5">
      <div v-if="loading" class="space-y-4">
        <div class="h-32 animate-pulse rounded-lg border border-slate-200 bg-white dark:border-slate-800 dark:bg-slate-900" />
        <div class="grid grid-cols-1 gap-3 md:grid-cols-2 xl:grid-cols-4">
          <div v-for="item in 4" :key="item" class="h-28 animate-pulse rounded-lg border border-slate-200 bg-white dark:border-slate-800 dark:bg-slate-900" />
        </div>
        <div class="grid gap-3 md:grid-cols-2 xl:grid-cols-3">
          <div v-for="item in 6" :key="item" class="h-28 animate-pulse rounded-lg border border-slate-200 bg-white dark:border-slate-800 dark:bg-slate-900" />
        </div>
      </div>

      <div v-else class="space-y-4">
        <IamUserIdentityCard :user="user" />
        <IamUserStatCards :user="user" :traffic="traffic" />
        <IamUserQuickLinks />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useAccounting, type TrafficSummary } from '~/composables/useAccounting'

definePageMeta({
  middleware: 'auth'
})

const { t } = useI18n()
const localePath = useLocalePath()
const { user, fetchUser } = useAuth()
const accounting = useAccounting()

const loading = ref(true)
const traffic = ref<TrafficSummary | null>(null)

onMounted(async () => {
  loading.value = true

  try {
    await fetchUser()
    await loadTraffic()
  } catch {
    await navigateTo(localePath('/login'))
  } finally {
    loading.value = false
  }
})

async function loadTraffic() {
  try {
    traffic.value = await accounting.getTraffic()
  } catch {
    traffic.value = null
  }
}

useSeoMeta({
  title: t('user.metaTitle'),
  robots: 'noindex, nofollow'
})
</script>
