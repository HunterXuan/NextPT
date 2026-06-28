<template>
  <div>
    <template v-if="!isLoggedIn">
      <section class="relative isolate overflow-hidden bg-slate-950 text-white">
        <img
          src="/images/home-hero.png"
          alt=""
          class="absolute inset-0 -z-10 h-full w-full object-cover"
        >
        <div class="absolute inset-0 -z-10 bg-slate-950/62" />

        <div class="mx-auto flex min-h-[calc(100svh-9rem)] max-w-7xl items-center px-4 py-16 sm:px-6 lg:min-h-[calc(100svh-4rem)] lg:px-8 lg:py-20">
          <div class="max-w-3xl">
            <p class="inline-flex items-center gap-2 rounded-md border border-white/15 bg-white/10 px-3 py-1.5 text-xs font-medium text-sky-100 backdrop-blur">
              <UIcon name="i-lucide-sparkles" class="size-4 text-amber-300" />
              {{ $t('home.guest.eyebrow') }}
            </p>
            <h1 class="mt-5 text-4xl font-semibold leading-tight sm:text-5xl lg:text-6xl">
              {{ $t('home.guest.title') }}
            </h1>
            <p class="mt-5 max-w-2xl text-base leading-7 text-slate-200 sm:text-lg">
              {{ $t('home.guest.description') }}
            </p>

            <div class="mt-8 flex flex-col gap-3 sm:flex-row">
              <UButton color="primary" size="lg" icon="i-lucide-user-plus" :to="localePath('/register')">
                {{ $t('home.guest.primaryAction') }}
              </UButton>
              <UButton color="neutral" variant="outline" size="lg" icon="i-lucide-log-in" :to="localePath('/login')">
                {{ $t('home.guest.secondaryAction') }}
              </UButton>
            </div>

            <dl class="mt-10 grid max-w-2xl grid-cols-1 gap-3 sm:grid-cols-3">
              <div v-for="item in guestStats" :key="item.label" class="rounded-lg border border-white/12 bg-white/10 px-4 py-3 backdrop-blur">
                <dt class="text-xs font-medium text-slate-300">{{ item.label }}</dt>
                <dd class="mt-1 text-lg font-semibold text-white">{{ item.value }}</dd>
              </div>
            </dl>
          </div>
        </div>
      </section>

      <section class="bg-white py-14 dark:bg-slate-950">
        <div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
          <div class="max-w-2xl">
            <p class="text-sm font-medium text-sky-600 dark:text-sky-400">{{ $t('home.guest.features.eyebrow') }}</p>
            <h2 class="mt-2 text-2xl font-semibold text-slate-950 dark:text-white">{{ $t('home.guest.features.title') }}</h2>
            <p class="mt-3 text-sm leading-6 text-slate-600 dark:text-slate-300">{{ $t('home.guest.features.description') }}</p>
          </div>

          <div class="mt-8 grid gap-4 md:grid-cols-2 xl:grid-cols-3">
            <article v-for="feature in guestFeatures" :key="feature.title" class="rounded-lg border border-slate-200 bg-white p-5 dark:border-slate-800 dark:bg-slate-900">
              <span class="flex size-10 items-center justify-center rounded-md" :class="feature.iconClass">
                <UIcon :name="feature.icon" class="size-5" />
              </span>
              <h3 class="mt-4 text-base font-semibold text-slate-950 dark:text-white">{{ feature.title }}</h3>
              <p class="mt-2 text-sm leading-6 text-slate-600 dark:text-slate-300">{{ feature.description }}</p>
            </article>
          </div>
        </div>
      </section>

      <section class="border-y border-slate-200 bg-slate-50 py-14 dark:border-slate-800 dark:bg-slate-900/40">
        <div class="mx-auto grid max-w-7xl gap-10 px-4 sm:px-6 lg:grid-cols-[360px_minmax(0,1fr)] lg:px-8">
          <div>
            <p class="text-sm font-medium text-sky-600 dark:text-sky-400">{{ $t('home.guest.flow.eyebrow') }}</p>
            <h2 class="mt-2 text-2xl font-semibold text-slate-950 dark:text-white">{{ $t('home.guest.flow.title') }}</h2>
            <p class="mt-3 text-sm leading-6 text-slate-600 dark:text-slate-300">{{ $t('home.guest.flow.description') }}</p>
          </div>

          <div class="grid gap-3 md:grid-cols-2">
            <article v-for="(step, index) in guestFlow" :key="step.title" class="rounded-lg border border-slate-200 bg-white p-5 dark:border-slate-800 dark:bg-slate-950">
              <p class="text-xs font-semibold text-sky-600 dark:text-sky-400">{{ numberFormatter.format(index + 1) }}</p>
              <h3 class="mt-2 text-base font-semibold text-slate-950 dark:text-white">{{ step.title }}</h3>
              <p class="mt-2 text-sm leading-6 text-slate-600 dark:text-slate-300">{{ step.description }}</p>
            </article>
          </div>
        </div>
      </section>

      <section class="bg-white py-14 dark:bg-slate-950">
        <div class="mx-auto grid max-w-7xl gap-10 px-4 sm:px-6 lg:grid-cols-[minmax(0,1fr)_420px] lg:px-8">
          <div>
            <p class="text-sm font-medium text-sky-600 dark:text-sky-400">{{ $t('home.guest.faq.eyebrow') }}</p>
            <h2 class="mt-2 text-2xl font-semibold text-slate-950 dark:text-white">{{ $t('home.guest.faq.title') }}</h2>
            <div class="mt-6 divide-y divide-slate-200 rounded-lg border border-slate-200 dark:divide-slate-800 dark:border-slate-800">
              <details v-for="item in guestFaqs" :key="item.question" class="group px-4 py-4">
                <summary class="flex cursor-pointer list-none items-center justify-between gap-4 text-sm font-semibold text-slate-950 dark:text-white">
                  {{ item.question }}
                  <UIcon name="i-lucide-chevron-down" class="size-4 shrink-0 text-slate-400 transition group-open:rotate-180" />
                </summary>
                <p class="mt-3 text-sm leading-6 text-slate-600 dark:text-slate-300">{{ item.answer }}</p>
              </details>
            </div>
          </div>

          <aside class="rounded-lg border border-slate-200 bg-slate-50 p-6 dark:border-slate-800 dark:bg-slate-900">
            <UIcon name="i-lucide-ticket-plus" class="size-8 text-sky-500" />
            <h2 class="mt-5 text-xl font-semibold text-slate-950 dark:text-white">{{ $t('home.guest.cta.title') }}</h2>
            <p class="mt-3 text-sm leading-6 text-slate-600 dark:text-slate-300">{{ $t('home.guest.cta.description') }}</p>
            <div class="mt-6 grid gap-3">
              <UButton color="primary" icon="i-lucide-user-plus" block :to="localePath('/register')">
                {{ $t('home.guest.primaryAction') }}
              </UButton>
              <UButton color="neutral" variant="outline" icon="i-lucide-log-in" block :to="localePath('/login')">
                {{ $t('home.guest.secondaryAction') }}
              </UButton>
            </div>
          </aside>
        </div>
      </section>
    </template>

    <template v-else>
      <section class="min-h-[calc(100vh-4rem)] bg-slate-50 py-6 dark:bg-slate-950">
        <div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
          <div class="mb-6 flex flex-col gap-4 lg:flex-row lg:items-end lg:justify-between">
            <div>
              <p class="text-sm font-medium text-slate-500 dark:text-slate-400">{{ $t('home.member.eyebrow') }}</p>
              <h1 class="mt-1 text-2xl font-semibold text-slate-950 dark:text-white">
                {{ $t('home.member.title', { name: user?.username || $t('nav.user') }) }}
              </h1>
              <p class="mt-2 max-w-2xl text-sm leading-6 text-slate-600 dark:text-slate-300">{{ $t('home.member.description') }}</p>
            </div>

            <div class="flex flex-wrap items-center gap-2">
              <UButton color="primary" icon="i-lucide-library" :to="localePath('/catalog/torrents')">
                {{ $t('nav.catalog') }}
              </UButton>
              <UButton color="neutral" variant="outline" icon="i-lucide-square-pen" :to="localePath('/forum/topics/create')">
                {{ $t('forum.actions.createTopic') }}
              </UButton>
            </div>
          </div>

          <div class="grid gap-4 md:grid-cols-2 xl:grid-cols-4">
            <NuxtLink
              v-for="action in memberActions"
              :key="action.to"
              :to="localePath(action.to)"
              class="group rounded-lg border border-slate-200 bg-white p-5 transition hover:border-sky-300 hover:bg-sky-50/70 dark:border-slate-800 dark:bg-slate-900 dark:hover:border-sky-800 dark:hover:bg-sky-950/30"
            >
              <div class="flex items-center justify-between gap-4">
                <span class="flex size-10 items-center justify-center rounded-md" :class="action.iconClass">
                  <UIcon :name="action.icon" class="size-5" />
                </span>
                <UIcon name="i-lucide-arrow-right" class="size-5 text-slate-400 transition group-hover:text-sky-600 dark:group-hover:text-sky-300" />
              </div>
              <h2 class="mt-4 text-base font-semibold text-slate-950 dark:text-white">{{ action.title }}</h2>
              <p class="mt-2 text-sm leading-6 text-slate-600 dark:text-slate-300">{{ action.description }}</p>
            </NuxtLink>
          </div>

          <div class="mt-6 grid gap-6 lg:grid-cols-[minmax(0,1fr)_360px]">
            <section class="rounded-lg border border-slate-200 bg-white p-5 dark:border-slate-800 dark:bg-slate-900">
              <h2 class="text-base font-semibold text-slate-950 dark:text-white">{{ $t('home.member.next.title') }}</h2>
              <div class="mt-4 grid gap-3 md:grid-cols-3">
                <div v-for="item in memberNext" :key="item.title" class="rounded-md border border-slate-200 p-4 dark:border-slate-800">
                  <UIcon :name="item.icon" class="size-5 text-sky-500" />
                  <h3 class="mt-3 text-sm font-semibold text-slate-950 dark:text-white">{{ item.title }}</h3>
                  <p class="mt-2 text-sm leading-6 text-slate-600 dark:text-slate-300">{{ item.description }}</p>
                </div>
              </div>
            </section>

            <aside class="rounded-lg border border-slate-200 bg-white p-5 dark:border-slate-800 dark:bg-slate-900">
              <h2 class="text-base font-semibold text-slate-950 dark:text-white">{{ $t('home.member.account.title') }}</h2>
              <dl class="mt-4 space-y-3 text-sm">
                <div class="flex items-center justify-between gap-3">
                  <dt class="text-slate-500 dark:text-slate-400">{{ $t('user.fields.role') }}</dt>
                  <dd class="font-medium text-slate-950 dark:text-white">{{ user?.roleName || '-' }}</dd>
                </div>
                <div class="flex items-center justify-between gap-3">
                  <dt class="text-slate-500 dark:text-slate-400">{{ $t('user.fields.id') }}</dt>
                  <dd class="font-medium text-slate-950 dark:text-white">#{{ user?.id || '-' }}</dd>
                </div>
              </dl>
              <UButton class="mt-5" color="neutral" variant="outline" icon="i-lucide-user-round" block :to="localePath('/iam/users/me')">
                {{ $t('home.member.account.action') }}
              </UButton>
            </aside>
          </div>
        </div>
      </section>
    </template>
  </div>
</template>

<script setup lang="ts">
const { t, locale } = useI18n()
const localePath = useLocalePath()
const { isLoggedIn, user } = useAuth()

const numberFormatter = computed(() => new Intl.NumberFormat(locale.value))

const guestStats = computed(() => [
  { label: t('home.guest.stats.privacy'), value: t('home.guest.stats.privacyValue') },
  { label: t('home.guest.stats.ratio'), value: t('home.guest.stats.ratioValue') },
  { label: t('home.guest.stats.community'), value: t('home.guest.stats.communityValue') }
])

const guestFeatures = computed(() => [
  {
    title: t('home.guest.featureItems.catalog.title'),
    description: t('home.guest.featureItems.catalog.description'),
    icon: 'i-lucide-library',
    iconClass: 'bg-sky-100 text-sky-700 dark:bg-sky-950 dark:text-sky-300'
  },
  {
    title: t('home.guest.featureItems.tracker.title'),
    description: t('home.guest.featureItems.tracker.description'),
    icon: 'i-lucide-radio-tower',
    iconClass: 'bg-emerald-100 text-emerald-700 dark:bg-emerald-950 dark:text-emerald-300'
  },
  {
    title: t('home.guest.featureItems.ratio.title'),
    description: t('home.guest.featureItems.ratio.description'),
    icon: 'i-lucide-scale',
    iconClass: 'bg-amber-100 text-amber-700 dark:bg-amber-950 dark:text-amber-300'
  },
  {
    title: t('home.guest.featureItems.forum.title'),
    description: t('home.guest.featureItems.forum.description'),
    icon: 'i-lucide-messages-square',
    iconClass: 'bg-cyan-100 text-cyan-700 dark:bg-cyan-950 dark:text-cyan-300'
  },
  {
    title: t('home.guest.featureItems.subtitle.title'),
    description: t('home.guest.featureItems.subtitle.description'),
    icon: 'i-lucide-captions',
    iconClass: 'bg-indigo-100 text-indigo-700 dark:bg-indigo-950 dark:text-indigo-300'
  },
  {
    title: t('home.guest.featureItems.account.title'),
    description: t('home.guest.featureItems.account.description'),
    icon: 'i-lucide-user-check',
    iconClass: 'bg-rose-100 text-rose-700 dark:bg-rose-950 dark:text-rose-300'
  }
])

const guestFlow = computed(() => [
  {
    title: t('home.guest.flow.items.invite.title'),
    description: t('home.guest.flow.items.invite.description')
  },
  {
    title: t('home.guest.flow.items.discover.title'),
    description: t('home.guest.flow.items.discover.description')
  },
  {
    title: t('home.guest.flow.items.seed.title'),
    description: t('home.guest.flow.items.seed.description')
  },
  {
    title: t('home.guest.flow.items.discuss.title'),
    description: t('home.guest.flow.items.discuss.description')
  }
])

const guestFaqs = computed(() => [
  {
    question: t('home.guest.faq.items.invite.question'),
    answer: t('home.guest.faq.items.invite.answer')
  },
  {
    question: t('home.guest.faq.items.ratio.question'),
    answer: t('home.guest.faq.items.ratio.answer')
  },
  {
    question: t('home.guest.faq.items.passkey.question'),
    answer: t('home.guest.faq.items.passkey.answer')
  },
  {
    question: t('home.guest.faq.items.forum.question'),
    answer: t('home.guest.faq.items.forum.answer')
  }
])

const memberActions = computed(() => [
  {
    title: t('home.member.actions.catalog.title'),
    description: t('home.member.actions.catalog.description'),
    to: '/catalog/torrents',
    icon: 'i-lucide-library',
    iconClass: 'bg-sky-100 text-sky-700 dark:bg-sky-950 dark:text-sky-300'
  },
  {
    title: t('home.member.actions.upload.title'),
    description: t('home.member.actions.upload.description'),
    to: '/catalog/torrents/upload',
    icon: 'i-lucide-upload',
    iconClass: 'bg-emerald-100 text-emerald-700 dark:bg-emerald-950 dark:text-emerald-300'
  },
  {
    title: t('home.member.actions.forum.title'),
    description: t('home.member.actions.forum.description'),
    to: '/forum',
    icon: 'i-lucide-messages-square',
    iconClass: 'bg-cyan-100 text-cyan-700 dark:bg-cyan-950 dark:text-cyan-300'
  },
  {
    title: t('home.member.actions.profile.title'),
    description: t('home.member.actions.profile.description'),
    to: '/iam/users/me',
    icon: 'i-lucide-user-round',
    iconClass: 'bg-amber-100 text-amber-700 dark:bg-amber-950 dark:text-amber-300'
  }
])

const memberNext = computed(() => [
  {
    title: t('home.member.next.items.search.title'),
    description: t('home.member.next.items.search.description'),
    icon: 'i-lucide-search'
  },
  {
    title: t('home.member.next.items.passkey.title'),
    description: t('home.member.next.items.passkey.description'),
    icon: 'i-lucide-key-round'
  },
  {
    title: t('home.member.next.items.community.title'),
    description: t('home.member.next.items.community.description'),
    icon: 'i-lucide-message-square-heart'
  }
])

useSeoMeta({
  title: t('home.metaTitle'),
  description: t('home.metaDescription')
})
</script>
