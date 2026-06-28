<template>
  <AppShell v-if="isLoggedIn">
    <slot />
  </AppShell>

  <div v-else class="min-h-screen bg-white text-slate-950 antialiased dark:bg-slate-950 dark:text-white">
    <AppHeader />
    <main>
      <slot />
    </main>
    <AppFooter />
  </div>
</template>

<script setup lang="ts">
const { isLoggedIn } = useAuth()
</script>
