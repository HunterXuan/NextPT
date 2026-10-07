<template>
  <footer class="border-t border-neutral-200 bg-neutral-50/70 dark:border-neutral-800 dark:bg-neutral-950">
    <div class="mx-auto flex max-w-7xl flex-col gap-4 px-4 py-8 text-sm text-neutral-500 sm:flex-row sm:items-center sm:justify-between sm:px-6 lg:px-8 dark:text-neutral-400">
      <div class="flex items-center gap-3">
        <AppLogo small />
        <div>
          <p class="font-semibold text-neutral-900 dark:text-white">{{ siteName }}</p>
          <p class="text-xs">{{ $t('footer.tagline') }}</p>
        </div>
      </div>

      <p class="text-xs">Copyright © {{ year }} {{ siteName }}</p>
    </div>
  </footer>
</template>

<script setup lang="ts">
const { siteName } = useSiteBranding()
const year = new Date().getFullYear()
</script>
