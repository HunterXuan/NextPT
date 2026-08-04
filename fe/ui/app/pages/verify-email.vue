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
              <UIcon name="i-lucide-mail-check" class="size-5" />
            </span>
            <h1 class="mt-5 max-w-lg text-3xl font-semibold leading-tight">{{ $t('auth.verifyEmail.heroTitle') }}</h1>
            <p class="mt-4 max-w-xl text-sm leading-7 text-slate-200">{{ $t('auth.verifyEmail.heroDescription') }}</p>
          </div>
        </div>
      </section>

      <section class="flex flex-col justify-center p-6 sm:p-8 lg:p-10">
        <NuxtLink :to="localePath('/login')" class="mb-8 flex w-fit items-center gap-2 text-sm font-medium text-slate-500 hover:text-slate-950 dark:text-slate-400 dark:hover:text-white">
          <UIcon name="i-lucide-arrow-left" class="size-4" />
          {{ $t('auth.verifyEmail.backToLogin') }}
        </NuxtLink>

        <div v-if="state === 'verifying'" class="flex flex-col items-start">
          <span class="flex size-11 items-center justify-center rounded-lg bg-sky-50 text-sky-600 dark:bg-sky-950/40 dark:text-sky-300">
            <UIcon name="i-lucide-loader-circle" class="size-5 animate-spin" />
          </span>
          <h2 class="mt-5 text-2xl font-semibold text-slate-950 dark:text-white">{{ $t('auth.verifyEmail.verifyingTitle') }}</h2>
          <p class="mt-3 text-sm leading-6 text-slate-500 dark:text-slate-400">{{ $t('auth.verifyEmail.verifyingDescription') }}</p>
        </div>

        <div v-else-if="state === 'verified'" class="flex flex-col items-start">
          <span class="flex size-11 items-center justify-center rounded-lg bg-emerald-50 text-emerald-600 dark:bg-emerald-950/40 dark:text-emerald-300">
            <UIcon name="i-lucide-badge-check" class="size-5" />
          </span>
          <h2 class="mt-5 text-2xl font-semibold text-slate-950 dark:text-white">{{ $t('auth.verifyEmail.successTitle') }}</h2>
          <p class="mt-3 text-sm leading-6 text-slate-500 dark:text-slate-400">{{ $t('auth.verifyEmail.successDescription') }}</p>
          <UButton class="mt-7" color="primary" :to="localePath('/login')">
            {{ $t('auth.verifyEmail.login') }}
          </UButton>
        </div>

        <template v-else>
          <div>
            <p class="text-sm font-medium" :class="state === 'invalid' ? 'text-red-600 dark:text-red-400' : 'text-sky-600 dark:text-sky-400'">
              {{ $t(state === 'invalid' ? 'auth.verifyEmail.invalidEyebrow' : 'auth.verifyEmail.eyebrow') }}
            </p>
            <h2 class="mt-2 text-2xl font-semibold text-slate-950 dark:text-white">
              {{ $t(state === 'invalid' ? 'auth.verifyEmail.invalidTitle' : 'auth.verifyEmail.title') }}
            </h2>
            <p class="mt-2 text-sm leading-6 text-slate-500 dark:text-slate-400">
              {{ $t(state === 'invalid' ? 'auth.verifyEmail.invalidDescription' : 'auth.verifyEmail.subtitle') }}
            </p>
          </div>

          <form class="mt-8 space-y-4" @submit.prevent="handleResend">
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
              v-if="requestSent"
              color="success"
              variant="soft"
              icon="i-lucide-mail-check"
              :title="$t('auth.verifyEmail.sentTitle')"
              :description="$t('auth.verifyEmail.sentDescription')"
            />

            <UAlert
              v-else-if="errorMessage"
              color="error"
              variant="soft"
              icon="i-lucide-circle-alert"
              :title="errorMessage"
            />

            <UButton type="submit" color="primary" block :loading="pending" :disabled="!canResend">
              {{ $t('auth.verifyEmail.resend') }}
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

type VerificationState = 'waiting' | 'verifying' | 'verified' | 'invalid'

const { t } = useI18n()
const localePath = useLocalePath()
const route = useRoute()
const { requestEmailVerification, verifyEmail } = useAuth()
const verificationEmail = useState<string>('auth:verificationEmail', () => '')

const token = computed(() => typeof route.query.token === 'string' ? route.query.token.trim() : '')
const email = ref(verificationEmail.value)
const state = ref<VerificationState>(token.value ? 'verifying' : 'waiting')
const pending = ref(false)
const requestSent = ref(!token.value && Boolean(verificationEmail.value))
const errorMessage = ref('')
const canResend = computed(() => !pending.value && email.value.includes('@'))

async function handleResend() {
  if (!canResend.value) return
  pending.value = true
  requestSent.value = false
  errorMessage.value = ''
  try {
    await requestEmailVerification(email.value)
    verificationEmail.value = email.value
    requestSent.value = true
  } catch (error) {
    errorMessage.value = error instanceof ApiError ? error.message : t('common.requestFailed')
  } finally {
    pending.value = false
  }
}

onMounted(async () => {
  if (!token.value) return
  try {
    await verifyEmail({ token: token.value })
    verificationEmail.value = ''
    state.value = 'verified'
  } catch {
    state.value = 'invalid'
  }
})

useSeoMeta({
  title: t('auth.verifyEmail.metaTitle'),
  robots: 'noindex, nofollow'
})
</script>
