<template>
  <div class="min-h-screen bg-slate-50 text-slate-950 antialiased dark:bg-slate-950 dark:text-white">
    <div class="fixed right-4 top-4 z-50 flex items-center gap-1 rounded-lg border border-slate-200 bg-white/88 p-1 shadow-sm backdrop-blur dark:border-slate-800 dark:bg-slate-950/88">
      <AppLocaleSwitcher />
      <AppThemeToggle />
    </div>
    <slot />
  </div>
</template>
