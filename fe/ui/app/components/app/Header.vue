<template>
  <header class="sticky top-0 z-50 border-b border-neutral-200 bg-white/88 backdrop-blur dark:border-neutral-800 dark:bg-neutral-950/88">
    <div class="mx-auto flex h-16 max-w-7xl items-center justify-between px-4 sm:px-6 lg:px-8">
      <NuxtLink :to="localePath('/')" class="mr-3 flex min-w-0 items-center gap-3" :aria-label="siteName">
        <AppLogo />
        <span class="truncate text-base font-semibold tracking-normal text-neutral-950 dark:text-white">{{ siteName }}</span>
      </NuxtLink>

      <div class="flex shrink-0 items-center gap-1.5">
        <AppLocaleSwitcher />

        <AppThemeToggle />

        <UButton
          class="hidden sm:inline-flex"
          color="neutral"
          variant="ghost"
          :to="localePath('/login')"
        >
          {{ $t('auth.login.action') }}
        </UButton>
        <UButton
          color="primary"
          variant="solid"
          :to="localePath('/register')"
        >
          {{ $t('auth.register.action') }}
        </UButton>
      </div>
    </div>
  </header>
</template>

<script setup lang="ts">
const { siteName } = useSiteBranding()
const localePath = useLocalePath()
</script>
