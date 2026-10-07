<template>
  <div class="min-h-screen bg-neutral-50 text-neutral-950 antialiased dark:bg-neutral-950 dark:text-white">
    <div class="fixed right-4 top-4 z-50 flex items-center gap-1 rounded-lg border border-neutral-200 bg-white/88 p-1 shadow-sm backdrop-blur dark:border-neutral-800 dark:bg-neutral-950/88">
      <AppLocaleSwitcher />
      <AppThemeToggle />
    </div>
    <slot />
  </div>
</template>
