<template>
  <div class="min-h-[calc(100vh-4rem)] bg-slate-50 py-6 dark:bg-slate-950">
    <div class="w-full px-3 sm:px-4 lg:px-5">
      <section class="rounded-lg border border-slate-200 bg-white dark:border-slate-800 dark:bg-slate-900">
        <div class="border-b border-slate-200 px-4 py-3 dark:border-slate-800">
          <h2 class="text-sm font-semibold text-slate-950 dark:text-white">{{ $t('admin.iam.invites.form.title') }}</h2>
        </div>

        <form class="grid gap-4 p-4" @submit.prevent="grantInvites">
          <label class="block">
            <span class="text-sm font-medium text-slate-700 dark:text-slate-200">{{ $t('admin.iam.invites.form.amount') }}</span>
            <input v-model.number="form.amount" type="number" min="1" step="1" class="mt-1 h-10 w-full rounded-md border border-slate-200 bg-white px-3 text-sm outline-none dark:border-slate-700 dark:bg-slate-950" :disabled="saving">
          </label>

          <label class="flex items-center gap-2 text-sm font-medium text-slate-700 dark:text-slate-200">
            <input v-model="form.isTemp" type="checkbox" class="size-4 rounded border-slate-300 text-emerald-600 focus:ring-emerald-500 dark:border-slate-600" :disabled="saving">
            <span>{{ $t('admin.iam.invites.form.isTemp') }}</span>
          </label>

          <label class="block">
            <span class="text-sm font-medium text-slate-700 dark:text-slate-200">{{ $t('admin.iam.invites.form.expireAt') }}</span>
            <input v-model="form.expireAt" type="datetime-local" class="mt-1 h-10 w-full rounded-md border border-slate-200 bg-white px-3 text-sm outline-none dark:border-slate-700 dark:bg-slate-950" :disabled="saving || !form.isTemp">
          </label>

          <p v-if="errorMessage" class="rounded-md border border-red-200 bg-red-50 px-3 py-2 text-sm text-red-700 dark:border-red-900 dark:bg-red-950/40 dark:text-red-200">{{ errorMessage }}</p>

          <div class="flex justify-end">
            <UButton type="submit" color="primary" icon="i-lucide-ticket-plus" :loading="saving" :disabled="form.amount < 1">
              {{ $t('admin.iam.invites.form.submit') }}
            </UButton>
          </div>
        </form>
      </section>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ApiError } from '~/composables/useApi'

definePageMeta({ layout: 'admin', middleware: 'admin' })

const { t } = useI18n()
const toast = useToast()
const adminApi = useAdmin()

const saving = ref(false)
const errorMessage = ref('')
const form = reactive({
  amount: 1,
  isTemp: false,
  expireAt: ''
})

useHead({ title: t('admin.iam.invites.title') })

async function grantInvites() {
  saving.value = true
  errorMessage.value = ''
  try {
    await adminApi.grantIamInvites({
      amount: Number(form.amount || 0),
      isTemp: form.isTemp,
      expireAt: form.isTemp && form.expireAt ? new Date(form.expireAt).toISOString() : null
    })
    toast.add({ title: t('admin.iam.invites.form.success'), color: 'success', icon: 'i-lucide-check-circle' })
    form.amount = 1
    form.isTemp = false
    form.expireAt = ''
  } catch (error) {
    errorMessage.value = error instanceof ApiError ? error.message : t('common.requestFailed')
  } finally {
    saving.value = false
  }
}
</script>
