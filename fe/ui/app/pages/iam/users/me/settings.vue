<template>
  <div class="min-h-[calc(100vh-4rem)] bg-slate-50 py-6 dark:bg-slate-950">
    <div class="w-full px-3 sm:px-4 lg:px-5">
      <div class="grid gap-4 xl:grid-cols-[minmax(0,1fr)_380px] xl:items-start">
        <main class="min-w-0 space-y-4">
          <IamUserProfileForm />
          <IamUserInvitesPanel />
        </main>

        <aside class="space-y-4 xl:sticky xl:top-20">
          <IamUserPasskeyCard />
          <IamUserSecurityForm />
        </aside>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  middleware: 'auth'
})

const { t } = useI18n()

useSeoMeta({
  title: t('user.nav.settings'),
  robots: 'noindex, nofollow'
})
</script>
