<template>
  <div class="min-h-[calc(100vh-4rem)] bg-slate-50 px-4 py-12 dark:bg-slate-950">
    <div class="mx-auto grid max-w-5xl grid-cols-1 gap-8 lg:grid-cols-[1fr_400px]">
      <section class="hidden rounded-lg border border-slate-200 bg-white p-8 lg:block dark:border-slate-800 dark:bg-slate-900">
        <UIcon name="i-lucide-shield-check" class="size-8 text-emerald-500" />
        <h1 class="mt-6 text-3xl font-semibold text-slate-950 dark:text-white">{{ $t('auth.login.heroTitle') }}</h1>
        <p class="mt-4 max-w-xl text-sm leading-7 text-slate-600 dark:text-slate-300">{{ $t('auth.login.heroDescription') }}</p>
        <div class="mt-8 grid grid-cols-1 gap-3">
          <div
            v-for="item in highlights"
            :key="item.title"
            class="rounded-lg border border-slate-200 p-4 dark:border-slate-800"
          >
            <p class="text-sm font-semibold text-slate-950 dark:text-white">{{ item.title }}</p>
            <p class="mt-1 text-sm text-slate-500 dark:text-slate-400">{{ item.description }}</p>
          </div>
        </div>
      </section>

      <UCard class="rounded-lg">
        <template #header>
          <div>
            <h2 class="text-xl font-semibold text-slate-950 dark:text-white">{{ $t('auth.login.title') }}</h2>
            <p class="mt-1 text-sm text-slate-500 dark:text-slate-400">{{ $t('auth.login.subtitle') }}</p>
          </div>
        </template>

        <form class="space-y-4" @submit.prevent="handleSubmit">
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
            <UInput
              v-model="form.password"
              class="w-full"
              icon="i-lucide-lock-keyhole"
              autocomplete="current-password"
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
            {{ $t('auth.login.submit') }}
          </UButton>
        </form>

        <template #footer>
          <div class="flex items-center justify-between gap-3 text-sm">
            <span class="text-slate-500 dark:text-slate-400">{{ $t('auth.login.noAccount') }}</span>
            <NuxtLink :to="localePath('/register')" class="font-medium text-sky-600 hover:text-sky-700 dark:text-sky-400">
              {{ $t('auth.register.action') }}
            </NuxtLink>
          </div>
        </template>
      </UCard>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ApiError } from '~/composables/useApi'

definePageMeta({
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

const highlights = computed(() => [
  {
    title: t('auth.login.highlights.identity.title'),
    description: t('auth.login.highlights.identity.description')
  },
  {
    title: t('auth.login.highlights.tracker.title'),
    description: t('auth.login.highlights.tracker.description')
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
  title: t('auth.login.title'),
  robots: 'noindex, nofollow'
})
</script>
