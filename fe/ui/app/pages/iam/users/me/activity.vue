<template>
  <div class="min-h-[calc(100vh-4rem)] bg-slate-50 py-6 dark:bg-slate-950">
    <div class="w-full px-3 sm:px-4 lg:px-5">
      <div class="space-y-4">
        <IamUserTrafficTrend />
        <IamUserPeers />
        <IamUserSnatches />
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
  title: t('user.nav.activity'),
  robots: 'noindex, nofollow'
})
</script>
