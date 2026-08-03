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
              <UIcon name="i-lucide-key-round" class="size-5" />
            </span>
            <h1 class="mt-5 max-w-lg text-3xl font-semibold leading-tight">{{ $t('auth.forgotPassword.heroTitle') }}</h1>
            <p class="mt-4 max-w-xl text-sm leading-7 text-slate-200">{{ $t('auth.forgotPassword.heroDescription') }}</p>
          </div>
        </div>
      </section>

      <section class="flex flex-col justify-center p-6 sm:p-8 lg:p-10">
        <NuxtLink :to="localePath('/login')" class="mb-8 flex w-fit items-center gap-2 text-sm font-medium text-slate-500 hover:text-slate-950 dark:text-slate-400 dark:hover:text-white">
          <UIcon name="i-lucide-arrow-left" class="size-4" />
          {{ $t('auth.forgotPassword.backToLogin') }}
        </NuxtLink>

        <div v-if="submitted" class="flex flex-col items-start">
          <span class="flex size-11 items-center justify-center rounded-lg bg-emerald-50 text-emerald-600 dark:bg-emerald-950/40 dark:text-emerald-300">
            <UIcon name="i-lucide-mail-check" class="size-5" />
          </span>
          <h2 class="mt-5 text-2xl font-semibold text-slate-950 dark:text-white">{{ $t('auth.forgotPassword.successTitle') }}</h2>
          <p class="mt-3 text-sm leading-6 text-slate-500 dark:text-slate-400">{{ $t('auth.forgotPassword.successDescription') }}</p>
          <UButton class="mt-7" color="primary" :to="localePath('/login')">
            {{ $t('auth.forgotPassword.backToLogin') }}
          </UButton>
        </div>

        <template v-else>
          <div>
            <p class="text-sm font-medium text-sky-600 dark:text-sky-400">{{ $t('auth.forgotPassword.eyebrow') }}</p>
            <h2 class="mt-2 text-2xl font-semibold text-slate-950 dark:text-white">{{ $t('auth.forgotPassword.title') }}</h2>
            <p class="mt-2 text-sm leading-6 text-slate-500 dark:text-slate-400">{{ $t('auth.forgotPassword.subtitle') }}</p>
          </div>

          <form class="mt-8 space-y-4" @submit.prevent="handleSubmit">
            <UFormField :label="$t('auth.fields.email')" required>
              <UInput
                v-model.trim="email"
                class="w-full"
                type="email"
                icon="i-lucide-mail"
                autocomplete="email"
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
              {{ $t('auth.forgotPassword.submit') }}
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
const { requestPasswordReset } = useAuth()

const email = ref('')
const pending = ref(false)
const submitted = ref(false)
const errorMessage = ref('')
const canSubmit = computed(() => email.value.includes('@'))

async function handleSubmit() {
  if (!canSubmit.value) return
  pending.value = true
  errorMessage.value = ''
  try {
    await requestPasswordReset(email.value)
    submitted.value = true
  } catch (error) {
    errorMessage.value = error instanceof ApiError ? error.message : t('common.requestFailed')
  } finally {
    pending.value = false
  }
}

useSeoMeta({
  title: t('auth.forgotPassword.metaTitle'),
  robots: 'noindex, nofollow'
})
</script>
