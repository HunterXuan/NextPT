<template>
  <header class="sticky top-0 z-50 border-b border-slate-200 bg-white/88 backdrop-blur dark:border-slate-800 dark:bg-slate-950/88">
    <div class="mx-auto flex h-16 max-w-7xl items-center justify-between px-4 sm:px-6 lg:px-8">
      <div class="flex items-center gap-7">
        <NuxtLink :to="localePath('/')" class="flex items-center gap-3">
          <span class="flex size-9 items-center justify-center rounded-lg bg-slate-950 text-white dark:bg-white dark:text-slate-950">
            <UIcon name="i-lucide-radio-tower" class="size-5" />
          </span>
          <span class="text-base font-semibold tracking-normal text-slate-950 dark:text-white">NextPT</span>
        </NuxtLink>

        <nav class="hidden items-center gap-1 md:flex">
          <NuxtLink
            v-for="item in navItems"
            :key="item.to"
            :to="localePath(item.to)"
            class="rounded-md px-3 py-2 text-sm font-medium transition-colors"
            :class="item.active ? 'bg-slate-100 text-slate-950 dark:bg-slate-800 dark:text-white' : 'text-slate-600 hover:bg-slate-100 hover:text-slate-950 dark:text-slate-300 dark:hover:bg-slate-800 dark:hover:text-white'"
          >
            {{ item.label }}
          </NuxtLink>
        </nav>
      </div>

      <div class="flex items-center gap-1.5">
        <UDropdownMenu :items="languageItems" :content="{ align: 'end' }">
          <UButton
            color="neutral"
            variant="ghost"
            icon="i-lucide-languages"
            :aria-label="$t('common.switchLanguage')"
          />
        </UDropdownMenu>

        <AppThemeToggle />

        <UButton
          v-if="!isLoggedIn"
          class="hidden sm:inline-flex"
          color="neutral"
          variant="ghost"
          :to="localePath('/login')"
        >
          {{ $t('auth.login.action') }}
        </UButton>
        <UButton
          v-if="!isLoggedIn"
          color="primary"
          variant="solid"
          :to="localePath('/register')"
        >
          {{ $t('auth.register.action') }}
        </UButton>

        <UDropdownMenu v-if="isLoggedIn" :items="userMenuItems" :content="{ align: 'end' }">
          <UButton
            color="neutral"
            variant="ghost"
            icon="i-lucide-user-round"
            trailing-icon="i-lucide-chevron-down"
            :label="user?.username || $t('nav.user')"
            :loading="loggingOut"
          />
        </UDropdownMenu>
      </div>
    </div>
  </header>
</template>

<script setup lang="ts">
const { t, locale, locales, setLocale } = useI18n()
const localePath = useLocalePath()
const route = useRoute()
const { user, isLoggedIn, isStaff, logout } = useAuth()
const loggingOut = ref(false)

const navItems = computed(() => {
  const items = [
    {
      label: t('nav.home'),
      to: '/',
      active: route.path === localePath('/')
    }
  ]

  if (isLoggedIn.value) {
    items.push({
      label: t('nav.catalog'),
      to: '/catalog/torrents',
      active: route.path.startsWith(localePath('/catalog'))
    })

    items.push({
      label: t('nav.forum'),
      to: '/forum',
      active: route.path.startsWith(localePath('/forum'))
    })

    items.push({
      label: t('nav.user'),
      to: '/iam/users/me',
      active: route.path.startsWith(localePath('/iam/users/me'))
    })
  }

  if (isStaff.value) {
    items.push({
      label: t('nav.admin'),
      to: '/admin',
      active: route.path.startsWith(localePath('/admin'))
    })
  }

  return items
})

const userMenuItems = computed(() => [
  [
    {
      label: t('nav.user'),
      icon: 'i-lucide-user-round',
      onSelect: () => navigateTo(localePath('/iam/users/me'))
    },
    {
      label: t('nav.catalog'),
      icon: 'i-lucide-library',
      onSelect: () => navigateTo(localePath('/catalog/torrents'))
    },
    {
      label: t('nav.forum'),
      icon: 'i-lucide-messages-square',
      onSelect: () => navigateTo(localePath('/forum'))
    },
    {
      label: t('nav.uploadTorrent'),
      icon: 'i-lucide-upload',
      onSelect: () => navigateTo(localePath('/catalog/torrents/upload'))
    },
    {
      label: t('nav.subtitles'),
      icon: 'i-lucide-captions',
      onSelect: () => navigateTo(localePath('/catalog/subtitles'))
    },
    {
      label: t('nav.bookmarks'),
      icon: 'i-lucide-bookmark',
      onSelect: () => navigateTo(localePath('/catalog/bookmarks'))
    }
  ],
  [
    {
      label: t('auth.logout'),
      icon: 'i-lucide-log-out',
      onSelect: handleLogout
    }
  ]
])

const languageItems = computed(() => [
  (locales.value as any[]).map((item) => ({
    label: item.name,
    onSelect: () => handleLanguageSwitch(item.code)
  }))
])

async function handleLanguageSwitch(code: string) {
  if (locale.value === code) return
  await setLocale(code)
}

async function handleLogout() {
  loggingOut.value = true
  try {
    await logout()
    await navigateTo(localePath('/'))
  } finally {
    loggingOut.value = false
  }
}
</script>
