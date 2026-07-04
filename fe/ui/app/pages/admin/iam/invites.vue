<template>
  <div class="min-h-[calc(100vh-4rem)] bg-slate-50 py-6 dark:bg-slate-950">
    <div class="w-full px-3 sm:px-4 lg:px-5">
      <div class="grid gap-4 xl:grid-cols-[minmax(0,560px)_320px] xl:items-start">
        <section class="rounded-lg border border-slate-200 bg-white dark:border-slate-800 dark:bg-slate-900">
          <div class="border-b border-slate-200 px-4 py-3 dark:border-slate-800">
            <h2 class="text-sm font-semibold text-slate-950 dark:text-white">{{ $t('admin.iam.invites.form.title') }}</h2>
          </div>

          <form class="grid gap-4 p-4" @submit.prevent="grantInvites">
            <label class="block">
              <span class="text-sm font-medium text-slate-700 dark:text-slate-200">{{ $t('admin.iam.invites.form.amount') }}</span>
              <input v-model.number="form.amount" type="number" min="1" step="1" class="mt-1 h-10 w-full rounded-md border border-slate-200 bg-white px-3 text-sm outline-none dark:border-slate-700 dark:bg-slate-950" :disabled="saving">
            </label>

            <label class="flex items-center justify-between gap-4 rounded-md border border-slate-200 px-3 py-2 dark:border-slate-800">
              <span class="text-sm font-medium text-slate-700 dark:text-slate-200">{{ $t('admin.iam.invites.form.isTemp') }}</span>
              <input v-model="form.isTemp" type="checkbox" class="size-4 rounded border-slate-300 text-emerald-600 focus:ring-emerald-500 dark:border-slate-600" :disabled="saving">
            </label>

            <label class="block">
              <span class="text-sm font-medium text-slate-700 dark:text-slate-200">{{ $t('admin.iam.invites.form.expireAt') }}</span>
              <input v-model="form.expireAt" type="datetime-local" class="mt-1 h-10 w-full rounded-md border border-slate-200 bg-white px-3 text-sm outline-none dark:border-slate-700 dark:bg-slate-950" :disabled="saving || !form.isTemp">
              <span class="mt-1 block text-xs text-slate-500 dark:text-slate-400">{{ expireHint }}</span>
            </label>

            <div class="rounded-md border border-slate-200 bg-slate-50 px-3 py-2 dark:border-slate-800 dark:bg-slate-950/60">
              <dl class="space-y-2 text-sm">
                <div class="flex items-center justify-between gap-3">
                  <dt class="text-slate-500 dark:text-slate-400">{{ $t('admin.iam.invites.preview.amount') }}</dt>
                  <dd class="font-medium text-slate-950 dark:text-white">{{ numberFormatter.format(normalizedAmount) }}</dd>
                </div>
                <div class="flex items-center justify-between gap-3">
                  <dt class="text-slate-500 dark:text-slate-400">{{ $t('admin.iam.invites.preview.type') }}</dt>
                  <dd class="font-medium text-slate-950 dark:text-white">{{ inviteTypeLabel }}</dd>
                </div>
                <div class="flex items-center justify-between gap-3">
                  <dt class="text-slate-500 dark:text-slate-400">{{ $t('admin.iam.invites.preview.expireAt') }}</dt>
                  <dd class="min-w-0 truncate font-medium text-slate-950 dark:text-white">{{ expirePreview }}</dd>
                </div>
              </dl>
            </div>

            <p v-if="errorMessage" class="rounded-md border border-red-200 bg-red-50 px-3 py-2 text-sm text-red-700 dark:border-red-900 dark:bg-red-950/40 dark:text-red-200">{{ errorMessage }}</p>

            <div class="flex justify-end">
              <UButton type="submit" color="primary" icon="i-lucide-ticket-plus" :loading="saving" :disabled="!canSubmit">
                {{ $t('admin.iam.invites.form.submit') }}
              </UButton>
            </div>
          </form>
        </section>

        <aside class="space-y-3">
          <section class="rounded-lg border border-slate-200 bg-white p-4 dark:border-slate-800 dark:bg-slate-900">
            <h2 class="text-sm font-semibold text-slate-950 dark:text-white">{{ $t('admin.iam.invites.guide.title') }}</h2>
            <ul class="mt-3 space-y-2 text-sm leading-6 text-slate-600 dark:text-slate-300">
              <li class="flex gap-2">
                <span class="mt-2 size-1.5 shrink-0 rounded-full bg-slate-300 dark:bg-slate-600" />
                <span>{{ $t('admin.iam.invites.guide.amount') }}</span>
              </li>
              <li class="flex gap-2">
                <span class="mt-2 size-1.5 shrink-0 rounded-full bg-slate-300 dark:bg-slate-600" />
                <span>{{ $t('admin.iam.invites.guide.temporary') }}</span>
              </li>
              <li class="flex gap-2">
                <span class="mt-2 size-1.5 shrink-0 rounded-full bg-slate-300 dark:bg-slate-600" />
                <span>{{ $t('admin.iam.invites.guide.delivery') }}</span>
              </li>
            </ul>
          </section>
        </aside>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ApiError } from '~/composables/useApi'

definePageMeta({ layout: 'admin', middleware: 'admin' })

const { t, locale } = useI18n()
const toast = useToast()
const adminApi = useAdmin()

const saving = ref(false)
const errorMessage = ref('')
const form = reactive({
  amount: 1,
  isTemp: false,
  expireAt: ''
})
const numberFormatter = computed(() => new Intl.NumberFormat(locale.value))
const normalizedAmount = computed(() => Math.max(0, Number(form.amount || 0)))
const canSubmit = computed(() => normalizedAmount.value > 0 && (!form.isTemp || Boolean(form.expireAt)) && !saving.value)
const inviteTypeLabel = computed(() => form.isTemp ? t('admin.iam.invites.preview.temporary') : t('admin.iam.invites.preview.permanent'))
const expireHint = computed(() => form.isTemp ? t('admin.iam.invites.form.expireAtRequired') : t('admin.iam.invites.form.expireAtDisabled'))
const expirePreview = computed(() => {
  if (!form.isTemp) return t('admin.iam.invites.preview.never')
  return form.expireAt ? new Date(form.expireAt).toLocaleString(locale.value) : t('admin.iam.invites.preview.notSet')
})

useHead({ title: t('admin.iam.invites.title') })

async function grantInvites() {
  if (!canSubmit.value) return
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
