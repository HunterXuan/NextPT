<template>
  <div class="flex min-h-screen items-center bg-slate-50 px-4 py-6 pt-16 sm:py-8 sm:pt-20 lg:py-10 dark:bg-slate-950">
    <div class="mx-auto grid min-h-[620px] w-full max-w-5xl overflow-hidden rounded-lg border border-slate-200 bg-white shadow-sm lg:grid-cols-[minmax(0,1fr)_420px] dark:border-slate-800 dark:bg-slate-900">
      <section class="relative hidden min-h-full overflow-hidden lg:block">
        <img src="/images/home-hero.png" alt="" class="absolute inset-0 h-full w-full object-cover">
        <div class="absolute inset-0 bg-slate-950/72" />
        <div class="relative flex h-full flex-col justify-between p-10 text-white">
          <NuxtLink :to="localePath('/')" class="flex w-fit items-center gap-3">
            <span class="flex size-9 items-center justify-center rounded-lg bg-white text-slate-950">
              <UIcon name="i-lucide-radio-tower" class="size-5" />
            </span>
            <span class="text-base font-semibold">NextPT</span>
          </NuxtLink>

          <div>
            <span class="flex size-10 items-center justify-center rounded-lg bg-white/10 text-sky-200">
              <UIcon name="i-lucide-shield-keyhole" class="size-5" />
            </span>
            <h1 class="mt-5 max-w-lg text-3xl font-semibold leading-tight">{{ $t('auth.resetPassword.heroTitle') }}</h1>
            <p class="mt-4 max-w-xl text-sm leading-7 text-slate-200">{{ $t('auth.resetPassword.heroDescription') }}</p>
          </div>
        </div>
      </section>

      <section class="flex flex-col justify-center p-6 sm:p-8 lg:p-10">
        <NuxtLink :to="localePath('/login')" class="mb-8 flex w-fit items-center gap-2 text-sm font-medium text-slate-500 hover:text-slate-950 dark:text-slate-400 dark:hover:text-white">
          <UIcon name="i-lucide-arrow-left" class="size-4" />
          {{ $t('auth.resetPassword.backToLogin') }}
        </NuxtLink>

        <div v-if="completed" class="flex flex-col items-start">
          <span class="flex size-11 items-center justify-center rounded-lg bg-emerald-50 text-emerald-600 dark:bg-emerald-950/40 dark:text-emerald-300">
            <UIcon name="i-lucide-shield-check" class="size-5" />
          </span>
          <h2 class="mt-5 text-2xl font-semibold text-slate-950 dark:text-white">{{ $t('auth.resetPassword.successTitle') }}</h2>
          <p class="mt-3 text-sm leading-6 text-slate-500 dark:text-slate-400">{{ $t('auth.resetPassword.successDescription') }}</p>
          <UButton class="mt-7" color="primary" :to="localePath('/login')">
            {{ $t('auth.resetPassword.login') }}
          </UButton>
        </div>

        <div v-else-if="!token" class="flex flex-col items-start">
          <span class="flex size-11 items-center justify-center rounded-lg bg-red-50 text-red-600 dark:bg-red-950/40 dark:text-red-300">
            <UIcon name="i-lucide-link-2-off" class="size-5" />
          </span>
          <h2 class="mt-5 text-2xl font-semibold text-slate-950 dark:text-white">{{ $t('auth.resetPassword.invalidTitle') }}</h2>
          <p class="mt-3 text-sm leading-6 text-slate-500 dark:text-slate-400">{{ $t('auth.resetPassword.invalidDescription') }}</p>
          <UButton class="mt-7" color="neutral" variant="outline" :to="localePath('/forgot-password')">
            {{ $t('auth.resetPassword.requestAgain') }}
          </UButton>
        </div>

        <template v-else>
          <div>
            <p class="text-sm font-medium text-sky-600 dark:text-sky-400">{{ $t('auth.resetPassword.eyebrow') }}</p>
            <h2 class="mt-2 text-2xl font-semibold text-slate-950 dark:text-white">{{ $t('auth.resetPassword.title') }}</h2>
            <p class="mt-2 text-sm leading-6 text-slate-500 dark:text-slate-400">{{ $t('auth.resetPassword.subtitle') }}</p>
          </div>

          <form class="mt-8 space-y-4" @submit.prevent="handleSubmit">
            <UFormField :label="$t('auth.fields.newPassword')" required>
              <div class="grid grid-cols-[minmax(0,1fr)_auto] gap-2">
                <UInput
                  v-model="form.newPassword"
                  class="w-full"
                  icon="i-lucide-lock-keyhole"
                  autocomplete="new-password"
                  :type="showPassword ? 'text' : 'password'"
                  :disabled="pending"
                />
                <UButton
                  type="button"
                  color="neutral"
                  variant="outline"
                  :icon="showPassword ? 'i-lucide-eye-off' : 'i-lucide-eye'"
                  :aria-label="showPassword ? $t('common.hidePassword') : $t('common.showPassword')"
                  :disabled="pending"
                  @click="showPassword = !showPassword"
                />
              </div>
            </UFormField>

            <UFormField :label="$t('auth.fields.confirmPassword')" required :error="confirmError">
              <UInput
                v-model="form.confirmPassword"
                class="w-full"
                icon="i-lucide-lock-keyhole"
                autocomplete="new-password"
                :type="showPassword ? 'text' : 'password'"
                :disabled="pending"
              />
            </UFormField>

            <UAlert
              v-if="errorMessage"
              color="error"
              variant="soft"
              icon="i-lucide-circle-alert"
              :title="errorMessage"
            />

            <UButton type="submit" color="primary" block :loading="pending" :disabled="!canSubmit">
              {{ $t('auth.resetPassword.submit') }}
            </UButton>
          </form>
        </template>
      </section>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ApiError } from '~/composables/useApi'

definePageMeta({
  layout: 'auth',
  middleware: [
    () => {
      const localePath = useLocalePath()
      const { isLoggedIn } = useAuth()
      if (isLoggedIn.value) return navigateTo(localePath('/iam/users/me'), { replace: true })
    }
  ]
})

const { t } = useI18n()
const localePath = useLocalePath()
const route = useRoute()
const { resetPassword } = useAuth()

const token = computed(() => typeof route.query.token === 'string' ? route.query.token.trim() : '')
const form = reactive({ newPassword: '', confirmPassword: '' })
const pending = ref(false)
const completed = ref(false)
const showPassword = ref(false)
const errorMessage = ref('')
const passwordsMatch = computed(() => form.confirmPassword === '' || form.newPassword === form.confirmPassword)
const confirmError = computed(() => passwordsMatch.value ? undefined : t('auth.resetPassword.passwordMismatch'))
const canSubmit = computed(() => form.newPassword.length >= 6 && form.newPassword === form.confirmPassword)

async function handleSubmit() {
  if (!canSubmit.value || !token.value) return
  pending.value = true
  errorMessage.value = ''
  try {
    await resetPassword({ token: token.value, newPassword: form.newPassword })
    completed.value = true
  } catch (error) {
    errorMessage.value = error instanceof ApiError ? error.message : t('common.requestFailed')
  } finally {
    pending.value = false
  }
}

useSeoMeta({
  title: t('auth.resetPassword.metaTitle'),
  robots: 'noindex, nofollow'
})
</script>
