<template>
  <footer class="border-t border-slate-200 bg-slate-50/70 dark:border-slate-800 dark:bg-slate-950">
    <div class="mx-auto flex max-w-7xl flex-col gap-4 px-4 py-8 text-sm text-slate-500 sm:flex-row sm:items-center sm:justify-between sm:px-6 lg:px-8 dark:text-slate-400">
      <div class="flex items-center gap-3">
        <span class="flex size-8 items-center justify-center rounded-md bg-slate-900 text-white dark:bg-white dark:text-slate-950">
          <UIcon name="i-lucide-radio-tower" class="size-4" />
        </span>
        <div>
          <p class="font-semibold text-slate-900 dark:text-white">NextPT</p>
          <p class="text-xs">{{ $t('footer.tagline') }}</p>
        </div>
      </div>

      <p class="text-xs">Copyright © {{ year }} NextPT</p>
    </div>
  </footer>
</template>

<script setup lang="ts">
const year = new Date().getFullYear()
</script>
