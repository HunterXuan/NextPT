<template>
  <div class="min-h-[calc(100vh-4rem)] bg-slate-50 py-6 dark:bg-slate-950">
    <div class="w-full px-3 sm:px-4 lg:px-5">
      <div class="grid gap-4 xl:grid-cols-[minmax(0,1fr)_320px] xl:items-start">
        <IamLoginLogs
          :key="loginLogKey"
          admin
          :user-id="appliedUserId"
          :title="$t('admin.iam.loginLogs.title')"
        />

        <aside class="rounded-lg border border-slate-200 bg-white dark:border-slate-800 dark:bg-slate-900">
          <div class="border-b border-slate-200 px-4 py-3 dark:border-slate-800">
            <h2 class="text-sm font-semibold text-slate-950 dark:text-white">{{ $t('admin.iam.loginLogs.filters.title') }}</h2>
          </div>

          <form class="space-y-3 p-4" @submit.prevent="applyFilters">
            <label class="block">
              <span class="sr-only">{{ $t('admin.iam.loginLogs.filters.userId') }}</span>
              <input
                v-model.trim="userIdInput"
                inputmode="numeric"
                class="h-10 w-full rounded-md border border-slate-200 bg-white px-3 text-sm text-slate-950 outline-none transition focus:border-sky-400 focus:ring-2 focus:ring-sky-100 dark:border-slate-700 dark:bg-slate-950 dark:text-white dark:focus:border-sky-500 dark:focus:ring-sky-950"
                :placeholder="$t('admin.iam.loginLogs.filters.userIdPlaceholder')"
              >
            </label>
            <div class="grid grid-cols-2 gap-2">
              <UButton type="submit" color="primary" icon="i-lucide-search" class="justify-center">
                {{ $t('admin.iam.loginLogs.filters.apply') }}
              </UButton>
              <UButton type="button" color="neutral" variant="outline" icon="i-lucide-rotate-ccw" class="justify-center" :disabled="!hasUserFilter" @click="clearFilters">
                {{ $t('admin.iam.loginLogs.filters.clear') }}
              </UButton>
            </div>
          </form>
        </aside>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({ layout: 'admin', middleware: 'admin' })

const { t } = useI18n()

const userIdInput = ref('')
const appliedUserId = ref<number | null>(null)
const hasUserFilter = computed(() => appliedUserId.value !== null || userIdInput.value.trim() !== '')
const loginLogKey = computed(() => appliedUserId.value ? `user-${appliedUserId.value}` : 'all')

useHead({ title: t('admin.iam.loginLogs.title') })

function applyFilters() {
  appliedUserId.value = normalizeUserId(userIdInput.value)
  userIdInput.value = appliedUserId.value ? String(appliedUserId.value) : ''
}

function clearFilters() {
  userIdInput.value = ''
  appliedUserId.value = null
}

function normalizeUserId(value: string) {
  const id = Number.parseInt(value, 10)
  return Number.isFinite(id) && id > 0 ? id : null
}
</script>
