<template>
  <div>
    <section class="relative overflow-hidden border-b border-slate-200 bg-slate-50 dark:border-slate-800 dark:bg-slate-950">
      <div class="absolute inset-0 app-grid-bg opacity-70" />
      <div class="relative mx-auto grid min-h-[calc(100vh-4rem)] max-w-7xl grid-cols-1 gap-10 px-4 py-12 sm:px-6 lg:grid-cols-[1fr_440px] lg:px-8 lg:py-16">
        <div class="flex flex-col justify-center">
          <div class="mb-5 inline-flex w-fit items-center gap-2 rounded-md border border-slate-200 bg-white px-3 py-1.5 text-xs font-medium text-slate-600 dark:border-slate-800 dark:bg-slate-900 dark:text-slate-300">
            <UIcon name="i-lucide-lock-keyhole" class="size-4 text-emerald-500" />
            {{ $t('home.eyebrow') }}
          </div>

          <h1 class="max-w-3xl text-4xl font-semibold leading-tight text-slate-950 sm:text-5xl lg:text-6xl dark:text-white">
            {{ $t('home.title') }}
          </h1>
          <p class="mt-5 max-w-2xl text-base leading-7 text-slate-600 sm:text-lg dark:text-slate-300">
            {{ $t('home.description') }}
          </p>

          <div class="mt-8 flex flex-col gap-3 sm:flex-row">
            <UButton
              color="primary"
              size="lg"
              :icon="isLoggedIn ? 'i-lucide-library' : 'i-lucide-user-plus'"
              :to="isLoggedIn ? localePath('/catalog/torrents') : localePath('/register')"
            >
              {{ isLoggedIn ? $t('nav.catalog') : $t('auth.register.action') }}
            </UButton>
            <UButton
              color="neutral"
              variant="outline"
              size="lg"
              icon="i-lucide-log-in"
              :to="isLoggedIn ? localePath('/iam/users/me') : localePath('/login')"
            >
              {{ isLoggedIn ? $t('home.openCenter') : $t('auth.login.action') }}
            </UButton>
          </div>

          <dl class="mt-10 grid max-w-2xl grid-cols-3 gap-3">
            <div
              v-for="item in stats"
              :key="item.label"
              class="rounded-lg border border-slate-200 bg-white/82 p-4 dark:border-slate-800 dark:bg-slate-900/72"
            >
              <dt class="text-xs font-medium text-slate-500 dark:text-slate-400">{{ item.label }}</dt>
              <dd class="mt-2 text-xl font-semibold text-slate-950 dark:text-white">{{ item.value }}</dd>
            </div>
          </dl>
        </div>

        <div class="flex items-center">
          <div class="app-surface w-full rounded-lg p-4">
            <div class="mb-4 flex items-center justify-between border-b border-slate-200 pb-3 dark:border-slate-800">
              <div>
                <p class="text-sm font-semibold text-slate-950 dark:text-white">{{ $t('home.panel.title') }}</p>
                <p class="text-xs text-slate-500 dark:text-slate-400">{{ $t('home.panel.subtitle') }}</p>
              </div>
              <UBadge color="success" variant="soft">{{ $t('home.panel.badge') }}</UBadge>
            </div>

            <div class="space-y-3">
              <div
                v-for="item in modules"
                :key="item.title"
                class="rounded-lg border border-slate-200 bg-white p-4 dark:border-slate-800 dark:bg-slate-950"
              >
                <div class="flex items-start gap-3">
                  <span class="flex size-10 shrink-0 items-center justify-center rounded-md" :class="item.iconClass">
                    <UIcon :name="item.icon" class="size-5" />
                  </span>
                  <div class="min-w-0 flex-1">
                    <div class="flex items-center justify-between gap-3">
                      <h2 class="text-sm font-semibold text-slate-950 dark:text-white">{{ item.title }}</h2>
                      <span class="text-xs font-medium text-slate-500 dark:text-slate-400">{{ item.meta }}</span>
                    </div>
                    <p class="mt-1 text-sm leading-6 text-slate-600 dark:text-slate-300">{{ item.description }}</p>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>

    <section class="bg-white py-12 dark:bg-slate-950">
      <div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
        <div class="grid grid-cols-1 gap-4 md:grid-cols-3">
          <div
            v-for="item in workflows"
            :key="item.title"
            class="rounded-lg border border-slate-200 p-5 dark:border-slate-800"
          >
            <UIcon :name="item.icon" class="mb-4 size-5 text-sky-500" />
            <h2 class="text-base font-semibold text-slate-950 dark:text-white">{{ item.title }}</h2>
            <p class="mt-2 text-sm leading-6 text-slate-600 dark:text-slate-300">{{ item.description }}</p>
          </div>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
const { t } = useI18n()
const localePath = useLocalePath()
const { isLoggedIn } = useAuth()

const stats = computed(() => [
  { label: t('home.stats.tracker'), value: '24/7' },
  { label: t('home.stats.ratio'), value: 'Ratio' },
  { label: t('home.stats.community'), value: 'Forum' }
])

const modules = computed(() => [
  {
    title: t('home.modules.catalog.title'),
    description: t('home.modules.catalog.description'),
    meta: t('home.modules.catalog.meta'),
    icon: 'i-lucide-library',
    iconClass: 'bg-sky-100 text-sky-700 dark:bg-sky-950 dark:text-sky-300'
  },
  {
    title: t('home.modules.tracker.title'),
    description: t('home.modules.tracker.description'),
    meta: t('home.modules.tracker.meta'),
    icon: 'i-lucide-radio-tower',
    iconClass: 'bg-emerald-100 text-emerald-700 dark:bg-emerald-950 dark:text-emerald-300'
  },
  {
    title: t('home.modules.forum.title'),
    description: t('home.modules.forum.description'),
    meta: t('home.modules.forum.meta'),
    icon: 'i-lucide-message-square-text',
    iconClass: 'bg-amber-100 text-amber-700 dark:bg-amber-950 dark:text-amber-300'
  }
])

const workflows = computed(() => [
  {
    title: t('home.workflow.profile.title'),
    description: t('home.workflow.profile.description'),
    icon: 'i-lucide-user-cog'
  },
  {
    title: t('home.workflow.passkey.title'),
    description: t('home.workflow.passkey.description'),
    icon: 'i-lucide-key-round'
  },
  {
    title: t('home.workflow.next.title'),
    description: t('home.workflow.next.description'),
    icon: 'i-lucide-list-checks'
  }
])

useSeoMeta({
  title: t('home.metaTitle'),
  description: t('home.metaDescription')
})
</script>
