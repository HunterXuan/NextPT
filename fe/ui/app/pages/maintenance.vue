<template>
  <main class="flex min-h-screen items-center bg-slate-50 px-5 py-12 text-center dark:bg-slate-950">
    <section class="mx-auto w-full max-w-xl">
      <span class="mx-auto flex size-12 items-center justify-center rounded-lg bg-amber-50 text-amber-700 dark:bg-amber-950/50 dark:text-amber-300">
        <UIcon name="i-lucide-wrench" class="size-6" />
      </span>
      <h1 class="mt-5 text-xl font-semibold text-slate-950 dark:text-white">{{ $t('site.maintenance.title') }}</h1>
      <p class="mx-auto mt-3 max-w-md whitespace-pre-line text-sm leading-6 text-slate-600 dark:text-slate-300">{{ message }}</p>
      <div class="mt-6 flex flex-wrap justify-center gap-2">
        <UButton icon="i-lucide-refresh-cw" @click="retry">
          {{ $t('site.maintenance.retry') }}
        </UButton>
        <UButton color="neutral" variant="outline" icon="i-lucide-house" @click="goHome">
          {{ $t('site.maintenance.backHome') }}
        </UButton>
      </div>
    </section>
  </main>
</template>

<script setup lang="ts">
definePageMeta({ layout: false })

const route = useRoute()
const { t } = useI18n()

const message = computed(() => {
  const value = route.query.message
  if (typeof value === 'string' && value.trim()) return value.trim()
  return t('site.maintenance.defaultMessage')
})

const returnTo = computed(() => {
  const value = route.query.returnTo
  if (typeof value !== 'string' || !value.startsWith('/') || value.startsWith('//') || value.startsWith('/maintenance')) {
    return '/'
  }
  return value
})

function retry() {
  window.location.assign(returnTo.value)
}

function goHome() {
  window.location.assign('/')
}

useHead(() => ({ title: t('site.maintenance.title') }))
</script>
