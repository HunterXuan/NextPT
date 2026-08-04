<template>
  <div class="flex min-h-screen items-center bg-slate-50 px-4 py-6 pt-16 sm:py-8 sm:pt-20 lg:py-10 dark:bg-slate-950">
    <div class="mx-auto grid min-h-[680px] w-full max-w-6xl overflow-hidden rounded-lg border border-slate-200 bg-white shadow-sm lg:grid-cols-[minmax(0,1fr)_420px] dark:border-slate-800 dark:bg-slate-900">
      <section class="relative hidden min-h-full overflow-hidden lg:block">
        <img src="/images/home-hero.png" alt="" class="absolute inset-0 h-full w-full object-cover">
        <div class="absolute inset-0 bg-slate-950/68" />
        <div class="relative flex h-full flex-col justify-between p-10 text-white">
          <NuxtLink :to="localePath('/')" class="flex w-fit items-center gap-3">
            <span class="flex size-9 items-center justify-center rounded-lg bg-white text-slate-950">
              <UIcon name="i-lucide-radio-tower" class="size-5" />
            </span>
            <span class="text-base font-semibold">NextPT</span>
          </NuxtLink>

          <div>
            <p class="text-sm font-medium text-sky-200">{{ $t('auth.login.eyebrow') }}</p>
            <h1 class="mt-3 max-w-lg text-3xl font-semibold leading-tight">{{ $t('auth.login.heroTitle') }}</h1>
            <p class="mt-4 max-w-xl text-sm leading-7 text-slate-200">{{ $t('auth.login.heroDescription') }}</p>

            <div class="mt-8 grid gap-3">
              <div v-for="item in benefits" :key="item.title" class="flex items-start gap-3">
                <span class="mt-0.5 flex size-8 shrink-0 items-center justify-center rounded-md bg-white/10 text-sky-200">
                  <UIcon :name="item.icon" class="size-4" />
                </span>
                <div>
                  <p class="text-sm font-semibold">{{ item.title }}</p>
                  <p class="mt-1 text-sm leading-6 text-slate-300">{{ item.description }}</p>
                </div>
              </div>
            </div>
          </div>
        </div>
      </section>

      <section class="flex flex-col justify-center p-6 sm:p-8 lg:p-10">
        <NuxtLink :to="localePath('/')" class="mb-8 flex w-fit items-center gap-2 text-sm font-medium text-slate-500 hover:text-slate-950 lg:hidden dark:text-slate-400 dark:hover:text-white">
          <UIcon name="i-lucide-arrow-left" class="size-4" />
          NextPT
        </NuxtLink>

        <div>
          <p class="text-sm font-medium text-sky-600 dark:text-sky-400">{{ $t('auth.login.eyebrow') }}</p>
          <h2 class="mt-2 text-2xl font-semibold text-slate-950 dark:text-white">{{ $t('auth.login.title') }}</h2>
          <p class="mt-2 text-sm leading-6 text-slate-500 dark:text-slate-400">{{ $t('auth.login.subtitle') }}</p>
        </div>

        <form class="mt-8 space-y-4" @submit.prevent="handleSubmit">
          <UFormField :label="$t('auth.fields.username')" required :error="fieldError">
            <UInput
              v-model="form.username"
              class="w-full"
              icon="i-lucide-user-round"
              autocomplete="username"
              :disabled="pending"
            />
          </UFormField>

          <UFormField :label="$t('auth.fields.password')" required>
            <div class="grid grid-cols-[minmax(0,1fr)_auto] gap-2">
              <UInput
                v-model="form.password"
                class="w-full"
                icon="i-lucide-lock-keyhole"
                autocomplete="current-password"
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
            <div class="mt-2 flex items-center justify-between gap-3">
              <NuxtLink :to="localePath('/verify-email')" class="text-xs font-medium text-slate-500 hover:text-slate-950 dark:text-slate-400 dark:hover:text-white">
                {{ $t('auth.login.verifyEmail') }}
              </NuxtLink>
              <NuxtLink :to="localePath('/forgot-password')" class="text-xs font-medium text-sky-600 hover:text-sky-700 dark:text-sky-400 dark:hover:text-sky-300">
                {{ $t('auth.login.forgotPassword') }}
              </NuxtLink>
            </div>
          </UFormField>

          <UAlert
            v-if="errorMessage"
            color="error"
            variant="soft"
            icon="i-lucide-circle-alert"
            :title="errorMessage"
          />

          <UButton type="submit" color="primary" block :loading="pending" :disabled="!canSubmit">
            {{ $t('auth.login.submit') }}
          </UButton>
        </form>

        <div class="mt-6 flex items-center justify-between gap-3 border-t border-slate-200 pt-5 text-sm dark:border-slate-800">
          <span class="text-slate-500 dark:text-slate-400">{{ $t('auth.login.noAccount') }}</span>
          <NuxtLink :to="localePath('/register')" class="font-medium text-sky-600 hover:text-sky-700 dark:text-sky-400">
            {{ $t('auth.register.action') }}
          </NuxtLink>
        </div>
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

      if (isLoggedIn.value) {
        return navigateTo(localePath('/iam/users/me'), { replace: true })
      }
    }
  ]
})

const { t } = useI18n()
const localePath = useLocalePath()
const route = useRoute()
const toast = useToast()
const { login } = useAuth()

const form = reactive({
  username: '',
  password: ''
})
const pending = ref(false)
const showPassword = ref(false)
const errorMessage = ref('')

const benefits = computed(() => [
  {
    title: t('auth.login.benefits.passkey.title'),
    description: t('auth.login.benefits.passkey.description'),
    icon: 'i-lucide-key-round'
  },
  {
    title: t('auth.login.benefits.stats.title'),
    description: t('auth.login.benefits.stats.description'),
    icon: 'i-lucide-chart-no-axes-combined'
  },
  {
    title: t('auth.login.benefits.community.title'),
    description: t('auth.login.benefits.community.description'),
    icon: 'i-lucide-messages-square'
  }
])

const canSubmit = computed(() => form.username.trim().length > 0 && form.password.length > 0)
const fieldError = computed(() => errorMessage.value || undefined)

async function handleSubmit() {
  if (!canSubmit.value) return

  pending.value = true
  errorMessage.value = ''

  try {
    await login(form.username.trim(), form.password)
    toast.add({
      title: t('auth.login.success'),
      color: 'success',
      icon: 'i-lucide-check-circle'
    })

    const redirect = typeof route.query.redirect === 'string' ? route.query.redirect : localePath('/iam/users/me')
    await navigateTo(redirect)
  } catch (error) {
    errorMessage.value = error instanceof ApiError ? error.message : t('common.requestFailed')
  } finally {
    pending.value = false
  }
}

useSeoMeta({
  title: t('auth.login.metaTitle'),
  robots: 'noindex, nofollow'
})
</script>
