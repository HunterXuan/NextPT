<template>
  <div class="min-h-[calc(100vh-4rem)] bg-slate-50 py-6 dark:bg-slate-950">
    <div class="w-full px-3 sm:px-4 lg:px-5">
      <div class="mb-4 flex justify-end">
        <div class="flex items-center gap-2 rounded-lg border border-slate-200 bg-white px-3 py-2 dark:border-slate-800 dark:bg-slate-900">
          <UIcon name="i-lucide-shield-check" class="size-4 text-emerald-600 dark:text-emerald-400" />
          <span class="text-sm font-medium text-slate-700 dark:text-slate-200">{{ user?.role.name || $t('user.staff') }}</span>
        </div>
      </div>

      <div class="mb-4 grid gap-3 sm:grid-cols-2 xl:grid-cols-4">
        <div v-for="stat in dashboardStats" :key="stat.key" class="rounded-lg border border-slate-200 bg-white px-4 py-3 dark:border-slate-800 dark:bg-slate-900">
          <div class="flex items-center justify-between gap-3">
            <p class="text-xs font-medium text-slate-500 dark:text-slate-400">{{ stat.label }}</p>
            <UIcon :name="stat.icon" class="size-4 text-slate-400" />
          </div>
          <p class="mt-2 text-2xl font-semibold text-slate-950 dark:text-white">
            <span v-if="dashboardPending" class="inline-block h-7 w-16 animate-pulse rounded bg-slate-200 dark:bg-slate-800" />
            <span v-else>{{ numberFormatter.format(stat.value) }}</span>
          </p>
        </div>
      </div>

      <div class="grid gap-4 md:grid-cols-2 xl:grid-cols-4">
        <NuxtLink
          :to="localePath('/admin/catalog/categories')"
          class="group rounded-lg border border-slate-200 bg-white p-5 transition hover:border-sky-300 hover:bg-sky-50/70 dark:border-slate-800 dark:bg-slate-900 dark:hover:border-sky-800 dark:hover:bg-sky-950/30"
        >
          <div class="flex items-center justify-between gap-3">
            <div class="flex items-center gap-3">
              <span class="flex size-10 items-center justify-center rounded-lg bg-sky-50 text-sky-700 dark:bg-sky-950 dark:text-sky-300">
                <UIcon name="i-lucide-tags" class="size-5" />
              </span>
              <div>
                <h2 class="text-base font-semibold text-slate-950 dark:text-white">{{ $t('admin.catalog.categories.title') }}</h2>
                <p class="mt-0.5 text-sm text-slate-500 dark:text-slate-400">{{ $t('admin.catalog.title') }}</p>
              </div>
            </div>
            <UIcon name="i-lucide-arrow-right" class="size-5 text-slate-400 transition group-hover:text-sky-600 dark:group-hover:text-sky-300" />
          </div>
        </NuxtLink>

        <NuxtLink
          :to="localePath('/admin/forum/categories')"
          class="group rounded-lg border border-slate-200 bg-white p-5 transition hover:border-cyan-300 hover:bg-cyan-50/70 dark:border-slate-800 dark:bg-slate-900 dark:hover:border-cyan-800 dark:hover:bg-cyan-950/30"
        >
          <div class="flex items-center gap-3">
            <span class="flex size-10 items-center justify-center rounded-lg bg-cyan-50 text-cyan-700 dark:bg-cyan-950 dark:text-cyan-300">
              <UIcon name="i-lucide-message-square-more" class="size-5" />
            </span>
            <div class="min-w-0 flex-1">
              <h2 class="text-base font-semibold text-slate-950 dark:text-white">{{ $t('admin.forum.title') }}</h2>
              <p class="mt-0.5 text-sm text-slate-500 dark:text-slate-400">{{ $t('admin.forum.entryDescription') }}</p>
            </div>
            <UIcon name="i-lucide-arrow-right" class="size-5 text-slate-400 transition group-hover:text-cyan-600 dark:group-hover:text-cyan-300" />
          </div>
        </NuxtLink>

        <NuxtLink
          :to="localePath('/admin/iam/users')"
          class="group rounded-lg border border-slate-200 bg-white p-5 transition hover:border-emerald-300 hover:bg-emerald-50/70 dark:border-slate-800 dark:bg-slate-900 dark:hover:border-emerald-800 dark:hover:bg-emerald-950/30"
        >
          <div class="flex items-center gap-3">
            <span class="flex size-10 items-center justify-center rounded-lg bg-emerald-50 text-emerald-700 dark:bg-emerald-950 dark:text-emerald-300">
              <UIcon name="i-lucide-users" class="size-5" />
            </span>
            <div class="min-w-0 flex-1">
              <h2 class="text-base font-semibold text-slate-950 dark:text-white">{{ $t('admin.iam.title') }}</h2>
              <p class="mt-0.5 text-sm text-slate-500 dark:text-slate-400">{{ $t('admin.iam.entryDescription') }}</p>
            </div>
            <UIcon name="i-lucide-arrow-right" class="size-5 text-slate-400 transition group-hover:text-emerald-600 dark:group-hover:text-emerald-300" />
          </div>
        </NuxtLink>

        <NuxtLink
          :to="localePath('/admin/site/configs')"
          class="group rounded-lg border border-slate-200 bg-white p-5 transition hover:border-indigo-300 hover:bg-indigo-50/70 dark:border-slate-800 dark:bg-slate-900 dark:hover:border-indigo-800 dark:hover:bg-indigo-950/30"
        >
          <div class="flex items-center gap-3">
            <span class="flex size-10 items-center justify-center rounded-lg bg-indigo-50 text-indigo-700 dark:bg-indigo-950 dark:text-indigo-300">
              <UIcon name="i-lucide-settings-2" class="size-5" />
            </span>
            <div class="min-w-0 flex-1">
              <h2 class="text-base font-semibold text-slate-950 dark:text-white">{{ $t('admin.site.title') }}</h2>
              <p class="mt-0.5 text-sm text-slate-500 dark:text-slate-400">{{ $t('admin.site.entryDescription') }}</p>
            </div>
            <UIcon name="i-lucide-arrow-right" class="size-5 text-slate-400 transition group-hover:text-indigo-600 dark:group-hover:text-indigo-300" />
          </div>
        </NuxtLink>

        <NuxtLink
          :to="localePath('/admin/mod/reports')"
          class="group rounded-lg border border-slate-200 bg-white p-5 transition hover:border-rose-300 hover:bg-rose-50/70 dark:border-slate-800 dark:bg-slate-900 dark:hover:border-rose-800 dark:hover:bg-rose-950/30"
        >
          <div class="flex items-center gap-3">
            <span class="flex size-10 items-center justify-center rounded-lg bg-rose-50 text-rose-700 dark:bg-rose-950 dark:text-rose-300">
              <UIcon name="i-lucide-flag" class="size-5" />
            </span>
            <div class="min-w-0 flex-1">
              <h2 class="text-base font-semibold text-slate-950 dark:text-white">{{ $t('admin.mod.title') }}</h2>
              <p class="mt-0.5 text-sm text-slate-500 dark:text-slate-400">{{ $t('admin.mod.entryDescription') }}</p>
            </div>
            <UIcon name="i-lucide-arrow-right" class="size-5 text-slate-400 transition group-hover:text-rose-600 dark:group-hover:text-rose-300" />
          </div>
        </NuxtLink>

        <NuxtLink
          :to="localePath('/admin/site/audits')"
          class="group rounded-lg border border-slate-200 bg-white p-5 transition hover:border-violet-300 hover:bg-violet-50/70 dark:border-slate-800 dark:bg-slate-900 dark:hover:border-violet-800 dark:hover:bg-violet-950/30"
        >
          <div class="flex items-center gap-3">
            <span class="flex size-10 items-center justify-center rounded-lg bg-violet-50 text-violet-700 dark:bg-violet-950 dark:text-violet-300">
              <UIcon name="i-lucide-scroll-text" class="size-5" />
            </span>
            <div class="min-w-0 flex-1">
              <h2 class="text-base font-semibold text-slate-950 dark:text-white">{{ $t('admin.site.audits.title') }}</h2>
              <p class="mt-0.5 text-sm text-slate-500 dark:text-slate-400">{{ $t('admin.site.audits.entryDescription') }}</p>
            </div>
            <UIcon name="i-lucide-arrow-right" class="size-5 text-slate-400 transition group-hover:text-violet-600 dark:group-hover:text-violet-300" />
          </div>
        </NuxtLink>

        <NuxtLink
          :to="localePath('/admin/sys/crons')"
          class="group rounded-lg border border-slate-200 bg-white p-5 transition hover:border-amber-300 hover:bg-amber-50/70 dark:border-slate-800 dark:bg-slate-900 dark:hover:border-amber-800 dark:hover:bg-amber-950/30"
        >
          <div class="flex items-center gap-3">
            <span class="flex size-10 items-center justify-center rounded-lg bg-amber-50 text-amber-700 dark:bg-amber-950 dark:text-amber-300">
              <UIcon name="i-lucide-clock-3" class="size-5" />
            </span>
            <div class="min-w-0 flex-1">
              <h2 class="text-base font-semibold text-slate-950 dark:text-white">{{ $t('admin.sys.crons.title') }}</h2>
              <p class="mt-0.5 text-sm text-slate-500 dark:text-slate-400">{{ $t('admin.sys.entryDescription') }}</p>
            </div>
            <UIcon name="i-lucide-arrow-right" class="size-5 text-slate-400 transition group-hover:text-amber-600 dark:group-hover:text-amber-300" />
          </div>
        </NuxtLink>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  layout: 'admin',
  middleware: 'admin'
})

const { t } = useI18n()
const localePath = useLocalePath()
const { user } = useAuth()
const adminApi = useAdmin()

const dashboardPending = ref(false)
const dashboard = reactive({
  users: 0,
  reports: 0,
  cheaters: 0,
  crons: 0
})

const numberFormatter = computed(() => new Intl.NumberFormat())
const dashboardStats = computed(() => [
  { key: 'users', label: t('admin.dashboard.stats.users'), icon: 'i-lucide-users', value: dashboard.users },
  { key: 'reports', label: t('admin.dashboard.stats.reports'), icon: 'i-lucide-flag', value: dashboard.reports },
  { key: 'cheaters', label: t('admin.dashboard.stats.cheaters'), icon: 'i-lucide-radar', value: dashboard.cheaters },
  { key: 'crons', label: t('admin.dashboard.stats.crons'), icon: 'i-lucide-clock-3', value: dashboard.crons }
])

useHead({
  title: t('admin.title')
})

onMounted(loadDashboard)

async function loadDashboard() {
  dashboardPending.value = true
  try {
    const [users, reports, cheaters, crons] = await Promise.all([
      adminApi.listIamUsers({ page: 1, size: 1 }).catch(() => ({ total: 0 })),
      adminApi.listModReports({ page: 1, size: 1, status: 0 }).catch(() => ({ total: 0 })),
      adminApi.listModCheaters({ page: 1, size: 1, status: 0 }).catch(() => ({ total: 0 })),
      adminApi.listSysCrons().catch(() => ({ list: [] }))
    ])
    dashboard.users = users.total || 0
    dashboard.reports = reports.total || 0
    dashboard.cheaters = cheaters.total || 0
    dashboard.crons = crons.list?.length || 0
  } finally {
    dashboardPending.value = false
  }
}
</script>
