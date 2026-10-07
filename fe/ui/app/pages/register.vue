<template>
  <div class="flex min-h-screen items-center bg-neutral-50 px-4 py-6 pt-16 sm:py-8 sm:pt-20 lg:py-10 dark:bg-neutral-950">
    <div class="mx-auto grid min-h-[720px] w-full max-w-6xl overflow-hidden rounded-lg border border-neutral-200 bg-white shadow-sm lg:grid-cols-[440px_minmax(0,1fr)] dark:border-neutral-800 dark:bg-neutral-900">
      <section class="flex flex-col justify-center p-6 sm:p-8 lg:p-10">
        <NuxtLink :to="localePath('/')" class="mb-8 flex w-fit items-center gap-2 text-sm font-medium text-neutral-500 hover:text-neutral-950 dark:text-neutral-400 dark:hover:text-white">
          <UIcon name="i-lucide-arrow-left" class="size-4" />
          {{ siteName }}
        </NuxtLink>

        <div>
          <p class="text-sm font-medium text-primary-600 dark:text-primary-400">{{ $t('auth.register.eyebrow') }}</p>
          <h1 class="mt-2 text-2xl font-semibold text-neutral-950 dark:text-white">{{ $t('auth.register.title', { siteName }) }}</h1>
          <p class="mt-2 text-sm leading-6 text-neutral-500 dark:text-neutral-400">{{ $t('auth.register.subtitle') }}</p>
        </div>

        <form class="mt-8 space-y-4" @submit.prevent="handleSubmit">
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
            <div class="grid grid-cols-[minmax(0,1fr)_auto] gap-2">
              <UInput
                v-model="form.password"
                class="w-full"
                icon="i-lucide-lock-keyhole"
                :type="showPassword ? 'text' : 'password'"
                autocomplete="new-password"
                :disabled="pending"
              />
              <UButton
                type="button"
                color="neutral"
                variant="outline"
                :icon="showPassword ? 'i-lucide-eye-off' : 'i-lucide-eye'"
                :aria-label="showPassword ? $t('common.hidePassword') : $t('common.showPassword')"
                :disabled="pending"
                @click="() => { showPassword = !showPassword }"
              />
            </div>
          </UFormField>

          <UFormField :label="$t('auth.fields.inviteHash')">
            <UInput
              v-model="form.inviteHash"
              class="w-full"
              icon="i-lucide-ticket"
              :disabled="pending"
            />
            <p class="mt-2 text-xs leading-5 text-neutral-500 dark:text-neutral-400">{{ $t('auth.register.inviteHint') }}</p>
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

        <div class="mt-6 flex items-center justify-between gap-3 border-t border-neutral-200 pt-5 text-sm dark:border-neutral-800">
          <span class="text-neutral-500 dark:text-neutral-400">{{ $t('auth.register.hasAccount') }}</span>
          <NuxtLink :to="localePath('/login')" class="font-medium text-primary-600 hover:text-primary-700 dark:text-primary-400">
            {{ $t('auth.login.action') }}
          </NuxtLink>
        </div>
      </section>

      <section class="relative hidden min-h-full overflow-hidden lg:block">
        <img src="/images/home-hero.png" alt="" class="absolute inset-0 h-full w-full object-cover">
        <div class="absolute inset-0 bg-neutral-950/70" />
        <div class="relative flex h-full flex-col justify-between p-10 text-white">
          <div class="flex items-center gap-3">
            <AppLogo dark />
            <span class="text-base font-semibold">{{ siteName }}</span>
          </div>

          <div>
            <p class="text-sm font-medium text-primary-200">{{ $t('auth.register.eyebrow') }}</p>
            <h2 class="mt-3 max-w-lg text-3xl font-semibold leading-tight">{{ $t('auth.register.heroTitle') }}</h2>
            <p class="mt-4 max-w-xl text-sm leading-7 text-neutral-200">{{ $t('auth.register.heroDescription') }}</p>

            <div class="mt-8 grid gap-3 sm:grid-cols-2">
              <div v-for="item in policies" :key="item.title" class="rounded-lg border border-white/12 bg-white/10 p-4 backdrop-blur">
                <UIcon :name="item.icon" class="size-5 text-primary-200" />
                <p class="mt-3 text-sm font-semibold">{{ item.title }}</p>
                <p class="mt-1 text-sm leading-6 text-neutral-300">{{ item.description }}</p>
              </div>
            </div>
          </div>
        </div>
      </section>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ApiError } from '~/composables/useApi'

const { siteName } = useSiteBranding()

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
const { register } = useAuth()
const verificationEmail = useState<string>('auth:verificationEmail', () => '')

const form = reactive({
  username: '',
  email: '',
  password: '',
  inviteHash: typeof route.query.invite === 'string' ? route.query.invite : ''
})
const pending = ref(false)
const showPassword = ref(false)
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
    verificationEmail.value = form.email.trim()
    await navigateTo(localePath('/verify-email'))
  } catch (error) {
    errorMessage.value = error instanceof ApiError ? error.message : t('common.requestFailed')
  } finally {
    pending.value = false
  }
}

useSeoMeta({
  title: t('auth.register.metaTitle'),
  robots: 'noindex, nofollow'
})
</script>
