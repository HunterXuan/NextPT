<template>
  <header class="sticky top-0 z-50 border-b border-slate-200 bg-white/88 backdrop-blur dark:border-slate-800 dark:bg-slate-950/88">
    <div class="mx-auto flex h-16 max-w-7xl items-center justify-between px-4 sm:px-6 lg:px-8">
      <NuxtLink :to="localePath('/')" class="flex items-center gap-3">
        <span class="flex size-9 items-center justify-center rounded-lg bg-slate-950 text-white dark:bg-white dark:text-slate-950">
          <UIcon name="i-lucide-radio-tower" class="size-5" />
        </span>
        <span class="text-base font-semibold tracking-normal text-slate-950 dark:text-white">NextPT</span>
      </NuxtLink>

      <div class="flex items-center gap-1.5">
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
const localePath = useLocalePath()
</script>
