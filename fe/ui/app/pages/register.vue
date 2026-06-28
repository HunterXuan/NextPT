<template>
  <div class="min-h-[calc(100vh-4rem)] bg-slate-50 px-4 py-12 dark:bg-slate-950">
    <div class="mx-auto grid max-w-5xl grid-cols-1 gap-8 lg:grid-cols-[400px_1fr]">
      <UCard class="rounded-lg">
        <template #header>
          <div>
            <h1 class="text-xl font-semibold text-slate-950 dark:text-white">{{ $t('auth.register.title') }}</h1>
            <p class="mt-1 text-sm text-slate-500 dark:text-slate-400">{{ $t('auth.register.subtitle') }}</p>
          </div>
        </template>

        <form class="space-y-4" @submit.prevent="handleSubmit">
          <UFormField :label="$t('auth.fields.username')" required>
            <UInput
              v-model="form.username"
              class="w-full"
              icon="i-lucide-user-round"
              autocomplete="username"
              :disabled="pending"
            />
          </UFormField>

          <UFormField :label="$t('auth.fields.email')" required>
            <UInput
              v-model="form.email"
              class="w-full"
              icon="i-lucide-mail"
              type="email"
              autocomplete="email"
              :disabled="pending"
            />
          </UFormField>

          <UFormField :label="$t('auth.fields.password')" required>
            <UInput
              v-model="form.password"
              class="w-full"
              icon="i-lucide-lock-keyhole"
              type="password"
              autocomplete="new-password"
              :disabled="pending"
            />
          </UFormField>

          <UFormField :label="$t('auth.fields.inviteHash')">
            <UInput
              v-model="form.inviteHash"
              class="w-full"
              icon="i-lucide-ticket"
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
            {{ $t('auth.register.submit') }}
          </UButton>
        </form>

        <template #footer>
          <div class="flex items-center justify-between gap-3 text-sm">
            <span class="text-slate-500 dark:text-slate-400">{{ $t('auth.register.hasAccount') }}</span>
            <NuxtLink :to="localePath('/login')" class="font-medium text-sky-600 hover:text-sky-700 dark:text-sky-400">
              {{ $t('auth.login.action') }}
            </NuxtLink>
          </div>
        </template>
      </UCard>

      <section class="rounded-lg border border-slate-200 bg-white p-8 dark:border-slate-800 dark:bg-slate-900">
        <UIcon name="i-lucide-sparkles" class="size-8 text-amber-500" />
        <h2 class="mt-6 text-3xl font-semibold text-slate-950 dark:text-white">{{ $t('auth.register.heroTitle') }}</h2>
        <p class="mt-4 max-w-xl text-sm leading-7 text-slate-600 dark:text-slate-300">{{ $t('auth.register.heroDescription') }}</p>

        <div class="mt-8 grid gap-3 sm:grid-cols-2">
          <div
            v-for="item in policies"
            :key="item.title"
            class="rounded-lg border border-slate-200 p-4 dark:border-slate-800"
          >
            <UIcon :name="item.icon" class="mb-3 size-5 text-sky-500" />
            <p class="text-sm font-semibold text-slate-950 dark:text-white">{{ item.title }}</p>
            <p class="mt-1 text-sm leading-6 text-slate-500 dark:text-slate-400">{{ item.description }}</p>
          </div>
        </div>
      </section>
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
const { register, login } = useAuth()

const form = reactive({
  username: '',
  email: '',
  password: '',
  inviteHash: typeof route.query.invite === 'string' ? route.query.invite : ''
})
const pending = ref(false)
const errorMessage = ref('')

const canSubmit = computed(() => (
  form.username.trim().length >= 3 &&
  form.email.includes('@') &&
  form.password.length >= 6
))

const policies = computed(() => [
  {
    title: t('auth.register.policies.identity.title'),
    description: t('auth.register.policies.identity.description'),
    icon: 'i-lucide-badge-check'
  },
  {
    title: t('auth.register.policies.ratio.title'),
    description: t('auth.register.policies.ratio.description'),
    icon: 'i-lucide-scale'
  },
  {
    title: t('auth.register.policies.community.title'),
    description: t('auth.register.policies.community.description'),
    icon: 'i-lucide-message-square-heart'
  },
  {
    title: t('auth.register.policies.security.title'),
    description: t('auth.register.policies.security.description'),
    icon: 'i-lucide-key-round'
  }
])

async function handleSubmit() {
  if (!canSubmit.value) return

  pending.value = true
  errorMessage.value = ''

  try {
    await register({
      username: form.username.trim(),
      email: form.email.trim(),
      password: form.password,
      inviteHash: form.inviteHash.trim()
    })
    await login(form.username.trim(), form.password)
    toast.add({
      title: t('auth.register.success'),
      color: 'success',
      icon: 'i-lucide-check-circle'
    })
    await navigateTo(localePath('/iam/users/me'))
  } catch (error) {
    errorMessage.value = error instanceof ApiError ? error.message : t('common.requestFailed')
  } finally {
    pending.value = false
  }
}

useSeoMeta({
  title: t('auth.register.title'),
  robots: 'noindex, nofollow'
})
</script>
