<template>
  <div class="min-h-screen bg-slate-50 text-slate-950 antialiased dark:bg-slate-950 dark:text-white">
    <AppHeader />

    <div class="mx-auto grid max-w-[1600px] lg:grid-cols-[260px_minmax(0,1fr)]">
      <aside class="hidden border-r border-slate-200 bg-white lg:block dark:border-slate-800 dark:bg-slate-900">
        <div class="sticky top-16 h-[calc(100vh-4rem)] overflow-y-auto px-4 py-5">
          <div class="mb-5 rounded-md border border-slate-200 px-3 py-3 dark:border-slate-800">
            <p class="text-xs font-medium uppercase text-slate-500 dark:text-slate-400">{{ $t('admin.eyebrow') }}</p>
            <p class="mt-1 truncate text-sm font-semibold text-slate-950 dark:text-white">{{ user?.roleName || $t('user.staff') }}</p>
          </div>

          <nav class="space-y-5">
            <section v-for="section in navSections" :key="section.key">
              <h2 class="px-2 text-xs font-semibold uppercase text-slate-400 dark:text-slate-500">{{ section.label }}</h2>
              <div class="mt-2 space-y-1">
                <NuxtLink
                  v-for="item in section.items"
                  :key="item.to"
                  :to="localePath(item.to)"
                  class="flex min-h-10 items-center gap-3 rounded-md px-3 text-sm font-medium transition-colors"
                  :class="item.active ? 'bg-slate-950 text-white dark:bg-white dark:text-slate-950' : 'text-slate-600 hover:bg-slate-100 hover:text-slate-950 dark:text-slate-300 dark:hover:bg-slate-800 dark:hover:text-white'"
                >
                  <UIcon :name="item.icon" class="size-4 shrink-0" />
                  <span class="truncate">{{ item.label }}</span>
                </NuxtLink>
              </div>
            </section>
          </nav>
        </div>
      </aside>

      <main class="min-w-0">
        <div class="border-b border-slate-200 bg-white px-4 py-3 lg:hidden dark:border-slate-800 dark:bg-slate-900">
          <div class="flex gap-2 overflow-x-auto pb-1">
            <NuxtLink
              v-for="item in flatNavItems"
              :key="item.to"
              :to="localePath(item.to)"
              class="flex h-9 shrink-0 items-center gap-2 rounded-md border px-3 text-sm font-medium"
              :class="item.active ? 'border-slate-950 bg-slate-950 text-white dark:border-white dark:bg-white dark:text-slate-950' : 'border-slate-200 bg-white text-slate-600 dark:border-slate-800 dark:bg-slate-900 dark:text-slate-300'"
            >
              <UIcon :name="item.icon" class="size-4" />
              {{ item.label }}
            </NuxtLink>
          </div>
        </div>

        <slot />
      </main>
    </div>
  </div>
</template>

<script setup lang="ts">
const { t } = useI18n()
const localePath = useLocalePath()
const route = useRoute()
const { user } = useAuth()

interface AdminNavItem {
  label: string
  to: string
  icon: string
  active: boolean
}

function isActive(path: string) {
  const localized = localePath(path)
  return route.path === localized || route.path.startsWith(`${localized}/`)
}

const navSections = computed(() => [
  {
    key: 'overview',
    label: t('admin.nav.overview'),
    items: [
      { label: t('admin.dashboard.title'), to: '/admin', icon: 'i-lucide-layout-dashboard', active: route.path === localePath('/admin') }
    ]
  },
  {
    key: 'catalog',
    label: t('admin.nav.catalog'),
    items: [
      { label: t('admin.catalog.categories.title'), to: '/admin/catalog/categories', icon: 'i-lucide-tags', active: isActive('/admin/catalog/categories') }
    ]
  },
  {
    key: 'forum',
    label: t('admin.nav.forum'),
    items: [
      { label: t('admin.forum.categories.title'), to: '/admin/forum/categories', icon: 'i-lucide-folder-tree', active: isActive('/admin/forum/categories') },
      { label: t('admin.forum.nodes.title'), to: '/admin/forum/nodes', icon: 'i-lucide-panels-top-left', active: isActive('/admin/forum/nodes') }
    ]
  },
  {
    key: 'iam',
    label: t('admin.nav.iam'),
    items: [
      { label: t('admin.iam.users.title'), to: '/admin/iam/users', icon: 'i-lucide-users', active: isActive('/admin/iam/users') },
      { label: t('admin.iam.roles.title'), to: '/admin/iam/roles', icon: 'i-lucide-shield-check', active: isActive('/admin/iam/roles') },
      { label: t('admin.iam.invites.title'), to: '/admin/iam/invites', icon: 'i-lucide-ticket-plus', active: isActive('/admin/iam/invites') }
    ]
  },
  {
    key: 'moderation',
    label: t('admin.nav.moderation'),
    items: [
      { label: t('admin.mod.reports.title'), to: '/admin/mod/reports', icon: 'i-lucide-flag', active: isActive('/admin/mod/reports') },
      { label: t('admin.mod.cheaters.title'), to: '/admin/mod/cheaters', icon: 'i-lucide-radar', active: isActive('/admin/mod/cheaters') }
    ]
  },
  {
    key: 'site',
    label: t('admin.nav.site'),
    items: [
      { label: t('admin.site.configs.title'), to: '/admin/site/configs', icon: 'i-lucide-settings-2', active: isActive('/admin/site/configs') },
      { label: t('admin.site.audits.title'), to: '/admin/site/audits', icon: 'i-lucide-scroll-text', active: isActive('/admin/site/audits') }
    ]
  },
  {
    key: 'system',
    label: t('admin.nav.system'),
    items: [
      { label: t('admin.sys.crons.title'), to: '/admin/sys/crons', icon: 'i-lucide-clock-3', active: isActive('/admin/sys/crons') }
    ]
  }
])

const flatNavItems = computed<AdminNavItem[]>(() => navSections.value.flatMap((section) => section.items))
</script>
