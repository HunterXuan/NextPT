<template>
  <div class="min-h-[calc(100vh-4rem)] bg-slate-50 py-6 dark:bg-slate-950">
    <div class="w-full px-3 sm:px-4 lg:px-5">
      <div class="mb-4 flex justify-end">
        <div class="flex items-center gap-2 rounded-lg border border-slate-200 bg-white px-3 py-2 dark:border-slate-800 dark:bg-slate-900">
          <UIcon name="i-lucide-shield-check" class="size-4 text-emerald-600 dark:text-emerald-400" />
          <span class="text-sm font-medium text-slate-700 dark:text-slate-200">{{ user?.role.name || $t('user.staff') }}</span>
        </div>
      </div>

      <div v-if="dashboardStats.length" class="mb-4 grid gap-3 sm:grid-cols-2 xl:grid-cols-4">
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

      <div v-if="adminEntryCards.length" class="grid gap-4 md:grid-cols-2 xl:grid-cols-4">
        <NuxtLink
          v-for="card in adminEntryCards"
          :key="card.key"
          :to="localePath(card.to)"
          :class="['group rounded-lg border border-slate-200 bg-white p-5 transition dark:border-slate-800 dark:bg-slate-900', card.hoverClass]"
        >
          <div class="flex items-center justify-between gap-3">
            <div class="flex items-center gap-3">
              <span :class="['flex size-10 items-center justify-center rounded-lg', card.iconClass]">
                <UIcon :name="card.icon" class="size-5" />
              </span>
              <div class="min-w-0">
                <h2 class="truncate text-base font-semibold text-slate-950 dark:text-white">{{ card.title }}</h2>
                <p class="mt-0.5 truncate text-sm text-slate-500 dark:text-slate-400">{{ card.description }}</p>
              </div>
            </div>
            <UIcon name="i-lucide-arrow-right" :class="['size-5 text-slate-400 transition', card.arrowClass]" />
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
const { user, hasPermission } = useAuth()
const adminApi = useAdmin()

const dashboardPending = ref(false)
const dashboard = reactive({
  users: 0,
  reports: 0,
  cheaters: 0,
  crons: 0
})

const numberFormatter = computed(() => new Intl.NumberFormat())
const dashboardStats = computed(() => compactDashboardItems([
  hasPermission(Permission.AdminIamUserManage) ? { key: 'users', label: t('admin.dashboard.stats.users'), icon: 'i-lucide-users', value: dashboard.users } : null,
  hasPermission(Permission.AdminModReportManage) ? { key: 'reports', label: t('admin.dashboard.stats.reports'), icon: 'i-lucide-flag', value: dashboard.reports } : null,
  hasPermission(Permission.AdminModCheaterManage) ? { key: 'cheaters', label: t('admin.dashboard.stats.cheaters'), icon: 'i-lucide-radar', value: dashboard.cheaters } : null,
  hasPermission(Permission.AdminSysCronManage) ? { key: 'crons', label: t('admin.dashboard.stats.crons'), icon: 'i-lucide-clock-3', value: dashboard.crons } : null
]))
const adminEntryCards = computed(() => compactDashboardItems([
  hasAnyPermission(Permission.AdminCatalogTorrentManage, Permission.AdminCatalogCategoryManage, Permission.AdminCatalogTagManage)
    ? {
        key: 'catalog',
        title: t('admin.catalog.title'),
        description: t('admin.catalog.entryDescription'),
        to: firstAllowedRoute([
          [Permission.AdminCatalogTorrentManage, '/admin/catalog/reviews'],
          [Permission.AdminCatalogCategoryManage, '/admin/catalog/categories'],
          [Permission.AdminCatalogTagManage, '/admin/catalog/tags']
        ]),
        icon: 'i-lucide-tags',
        iconClass: 'bg-sky-50 text-sky-700 dark:bg-sky-950 dark:text-sky-300',
        hoverClass: 'hover:border-sky-300 hover:bg-sky-50/70 dark:hover:border-sky-800 dark:hover:bg-sky-950/30',
        arrowClass: 'group-hover:text-sky-600 dark:group-hover:text-sky-300'
      }
    : null,
  hasAnyPermission(Permission.AdminForumCategoryManage, Permission.AdminForumNodeManage)
    ? {
        key: 'forum',
        title: t('admin.forum.title'),
        description: t('admin.forum.entryDescription'),
        to: firstAllowedRoute([
          [Permission.AdminForumCategoryManage, '/admin/forum/categories'],
          [Permission.AdminForumNodeManage, '/admin/forum/nodes']
        ]),
        icon: 'i-lucide-message-square-more',
        iconClass: 'bg-cyan-50 text-cyan-700 dark:bg-cyan-950 dark:text-cyan-300',
        hoverClass: 'hover:border-cyan-300 hover:bg-cyan-50/70 dark:hover:border-cyan-800 dark:hover:bg-cyan-950/30',
        arrowClass: 'group-hover:text-cyan-600 dark:group-hover:text-cyan-300'
      }
    : null,
  hasAnyPermission(Permission.AdminIamUserManage, Permission.AdminIamRoleManage, Permission.AdminIamInviteManage)
    ? {
        key: 'iam',
        title: t('admin.iam.title'),
        description: t('admin.iam.entryDescription'),
        to: firstAllowedRoute([
          [Permission.AdminIamUserManage, '/admin/iam/users'],
          [Permission.AdminIamRoleManage, '/admin/iam/roles'],
          [Permission.AdminIamInviteManage, '/admin/iam/invites']
        ]),
        icon: 'i-lucide-users',
        iconClass: 'bg-emerald-50 text-emerald-700 dark:bg-emerald-950 dark:text-emerald-300',
        hoverClass: 'hover:border-emerald-300 hover:bg-emerald-50/70 dark:hover:border-emerald-800 dark:hover:bg-emerald-950/30',
        arrowClass: 'group-hover:text-emerald-600 dark:group-hover:text-emerald-300'
      }
    : null,
  hasPermission(Permission.AdminSiteConfig)
    ? {
        key: 'site',
        title: t('admin.site.title'),
        description: t('admin.site.entryDescription'),
        to: '/admin/site/configs',
        icon: 'i-lucide-settings-2',
        iconClass: 'bg-indigo-50 text-indigo-700 dark:bg-indigo-950 dark:text-indigo-300',
        hoverClass: 'hover:border-indigo-300 hover:bg-indigo-50/70 dark:hover:border-indigo-800 dark:hover:bg-indigo-950/30',
        arrowClass: 'group-hover:text-indigo-600 dark:group-hover:text-indigo-300'
      }
    : null,
  hasAnyPermission(Permission.AdminModReportManage, Permission.AdminModCheaterManage)
    ? {
        key: 'moderation',
        title: t('admin.mod.title'),
        description: t('admin.mod.entryDescription'),
        to: firstAllowedRoute([
          [Permission.AdminModReportManage, '/admin/mod/reports'],
          [Permission.AdminModCheaterManage, '/admin/mod/cheaters']
        ]),
        icon: 'i-lucide-flag',
        iconClass: 'bg-rose-50 text-rose-700 dark:bg-rose-950 dark:text-rose-300',
        hoverClass: 'hover:border-rose-300 hover:bg-rose-50/70 dark:hover:border-rose-800 dark:hover:bg-rose-950/30',
        arrowClass: 'group-hover:text-rose-600 dark:group-hover:text-rose-300'
      }
    : null,
  hasPermission(Permission.AdminSiteAudit)
    ? {
        key: 'audit',
        title: t('admin.site.audits.title'),
        description: t('admin.site.audits.entryDescription'),
        to: '/admin/site/audits',
        icon: 'i-lucide-scroll-text',
        iconClass: 'bg-violet-50 text-violet-700 dark:bg-violet-950 dark:text-violet-300',
        hoverClass: 'hover:border-violet-300 hover:bg-violet-50/70 dark:hover:border-violet-800 dark:hover:bg-violet-950/30',
        arrowClass: 'group-hover:text-violet-600 dark:group-hover:text-violet-300'
      }
    : null,
  hasPermission(Permission.AdminSysCronManage)
    ? {
        key: 'system',
        title: t('admin.sys.crons.title'),
        description: t('admin.sys.entryDescription'),
        to: '/admin/sys/crons',
        icon: 'i-lucide-clock-3',
        iconClass: 'bg-amber-50 text-amber-700 dark:bg-amber-950 dark:text-amber-300',
        hoverClass: 'hover:border-amber-300 hover:bg-amber-50/70 dark:hover:border-amber-800 dark:hover:bg-amber-950/30',
        arrowClass: 'group-hover:text-amber-600 dark:group-hover:text-amber-300'
      }
    : null
]))

useHead({
  title: t('admin.title')
})

onMounted(loadDashboard)

async function loadDashboard() {
  dashboardPending.value = true
  try {
    const [users, reports, cheaters, crons] = await Promise.all([
      hasPermission(Permission.AdminIamUserManage) ? adminApi.listIamUsers({ page: 1, size: 1 }).catch(() => ({ total: 0 })) : Promise.resolve({ total: 0 }),
      hasPermission(Permission.AdminModReportManage) ? adminApi.listModReports({ page: 1, size: 1, status: 0 }).catch(() => ({ total: 0 })) : Promise.resolve({ total: 0 }),
      hasPermission(Permission.AdminModCheaterManage) ? adminApi.listModCheaters({ page: 1, size: 1, status: 0 }).catch(() => ({ total: 0 })) : Promise.resolve({ total: 0 }),
      hasPermission(Permission.AdminSysCronManage) ? adminApi.listSysCrons().catch(() => ({ list: [] })) : Promise.resolve({ list: [] })
    ])
    dashboard.users = users.total || 0
    dashboard.reports = reports.total || 0
    dashboard.cheaters = cheaters.total || 0
    dashboard.crons = crons.list?.length || 0
  } finally {
    dashboardPending.value = false
  }
}

function hasAnyPermission(...permissions: string[]) {
  return permissions.some((permission) => hasPermission(permission))
}

function firstAllowedRoute(candidates: Array<[string, string]>) {
  return candidates.find(([permission]) => hasPermission(permission))?.[1] || candidates[0]?.[1] || '/admin'
}

function compactDashboardItems<T>(items: Array<T | null>) {
  return items.filter((item): item is T => Boolean(item))
}
</script>
