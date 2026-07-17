<template>
  <div class="min-h-screen bg-slate-50 text-slate-950 antialiased dark:bg-slate-950 dark:text-white" style="--app-sticky-top: 7.5rem;">
    <div
      v-if="mobileSidebarOpen"
      class="fixed inset-0 z-40 bg-slate-950/40 backdrop-blur-sm lg:hidden"
      @click="mobileSidebarOpen = false"
    />

    <aside
      class="fixed inset-y-0 left-0 z-50 hidden border-r border-slate-200 transition-[width] duration-200 lg:block dark:border-slate-800"
      :class="sidebarCollapsed ? 'w-16' : 'w-64'"
    >
      <AppShellSidebar
        :sections="navSections"
        :user="shellUser"
        :user-menu-items="userMenuItems"
        :return-action="adminReturnAction"
        :logging-out="loggingOut"
        :collapsed="sidebarCollapsed"
        @navigate="mobileSidebarOpen = false"
      />
    </aside>

    <aside
      class="fixed inset-y-0 left-0 z-50 w-72 max-w-[86vw] border-r border-slate-200 shadow-2xl transition-transform duration-200 lg:hidden dark:border-slate-800"
      :class="mobileSidebarOpen ? 'translate-x-0' : '-translate-x-full'"
    >
      <AppShellSidebar
        :sections="navSections"
        :user="shellUser"
        :user-menu-items="userMenuItems"
        :return-action="adminReturnAction"
        :logging-out="loggingOut"
        show-close
        @close="mobileSidebarOpen = false"
        @navigate="mobileSidebarOpen = false"
      />
    </aside>

    <div class="min-w-0 transition-[padding-left] duration-200" :class="sidebarCollapsed ? 'lg:pl-16' : 'lg:pl-64'">
      <header class="sticky top-0 z-30 border-b border-slate-200 bg-white/88 backdrop-blur dark:border-slate-800 dark:bg-slate-950/88">
        <div class="flex h-16 items-center justify-between gap-3 px-3 sm:px-4 lg:px-5">
          <div class="flex min-w-0 items-center gap-3">
            <UButton
              class="lg:hidden"
              color="neutral"
              variant="ghost"
              icon="i-lucide-menu"
              :aria-label="activeItemLabel"
              @click="mobileSidebarOpen = true"
            />
            <UButton
              class="hidden lg:inline-flex"
              color="neutral"
              variant="ghost"
              :icon="sidebarCollapsed ? 'i-lucide-panel-left-open' : 'i-lucide-panel-left-close'"
              :aria-label="sidebarCollapsed ? $t('common.expandSidebar') : $t('common.collapseSidebar')"
              :title="sidebarCollapsed ? $t('common.expandSidebar') : $t('common.collapseSidebar')"
              @click="toggleSidebar"
            />
            <div class="min-w-0">
              <p class="truncate text-sm font-semibold text-slate-950 dark:text-white">{{ activeItemLabel }}</p>
            </div>
          </div>

          <div class="flex shrink-0 items-center gap-1.5">
            <UTooltip :text="$t('site.messages.title')" :content="{ side: 'bottom', sideOffset: 8 }" :delay-duration="600">
              <UButton
                v-if="props.mode === 'app' && hasPermission(Permission.SiteMessageRead)"
                class="relative"
                color="neutral"
                variant="ghost"
                icon="i-lucide-bell"
                :aria-label="$t('site.messages.title')"
                :to="localePath('/site/messages')"
              >
                <span
                  v-if="unreadMessageCount > 0"
                  class="absolute -right-0.5 -top-0.5 flex min-w-4 items-center justify-center rounded-full bg-red-500 px-1 text-[10px] font-semibold leading-4 text-white"
                >
                  {{ unreadMessageCount > 99 ? '99+' : unreadMessageCount }}
                </span>
              </UButton>
            </UTooltip>

            <UDropdownMenu :items="languageItems" :content="{ align: 'end' }">
              <UButton
                color="neutral"
                variant="ghost"
                icon="i-lucide-languages"
                :aria-label="$t('common.switchLanguage')"
              />
            </UDropdownMenu>

            <AppThemeToggle />
          </div>
        </div>
        <AppShellTabs :mode="props.mode" :current-title="activeItemLabel" @refresh="refreshCurrentPage" />
      </header>

      <main :key="pageRefreshKey" class="min-w-0">
        <slot />
      </main>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { DropdownMenuItem } from '@nuxt/ui'
import { normalizeLocaleCode } from '~/utils/locale'

type ShellMode = 'app' | 'admin'

interface ShellNavItem {
  label: string
  to: string
  icon: string
  active: boolean
}

interface ShellNavSection {
  key: string
  label: string
  items: ShellNavItem[]
}

const props = withDefaults(defineProps<{
  mode?: ShellMode
}>(), {
  mode: 'app'
})

const { t, locale, locales, setLocale } = useI18n()
const localePath = useLocalePath()
const route = useRoute()
const { user, isStaff, hasPermission, logout } = useAuth()
const siteApi = useSite()

const mobileSidebarOpen = ref(false)
const sidebarCollapsed = ref(false)
const loggingOut = ref(false)
const pageRefreshKey = ref(0)
const unreadMessageCount = ref(0)
const sidebarStorageKey = 'nextpt_sidebar_collapsed'

onMounted(() => {
  sidebarCollapsed.value = localStorage.getItem(sidebarStorageKey) === '1'
  void loadUnreadMessageCount()
  if (import.meta.client) {
    window.addEventListener('nextpt:messages-read', handleMessagesRead)
  }
})

onBeforeUnmount(() => {
  if (import.meta.client) {
    window.removeEventListener('nextpt:messages-read', handleMessagesRead)
  }
})

watch(sidebarCollapsed, (value) => {
  if (!import.meta.client) return
  localStorage.setItem(sidebarStorageKey, value ? '1' : '0')
})

watch(() => route.fullPath, () => {
  if (props.mode === 'app') {
    void loadUnreadMessageCount()
  }
})

function isActive(path: string) {
  const localized = localePath(path)
  if (localized === '/') return route.path === localized
  return route.path === localized || route.path.startsWith(`${localized}/`)
}

function isExactActive(path: string) {
  return route.path === localePath(path)
}

function isForumHomeActive() {
  const forumPath = localePath('/forum')
  const topicsPath = localePath('/forum/topics')
  const nodesPath = localePath('/forum/nodes')
  return route.path === forumPath || route.path.startsWith(`${topicsPath}/`) || route.path.startsWith(`${nodesPath}/`)
}

function isCatalogTorrentsActive() {
  const torrentsPath = localePath('/catalog/torrents')
  return route.path === torrentsPath || route.path.startsWith(`${torrentsPath}/`)
}

function refreshCurrentPage() {
  pageRefreshKey.value += 1
}

const currentAppPath = computed(() => {
  const codes = (locales.value as any[]).map((item) => item.code).filter(Boolean)
  for (const code of codes) {
    const prefix = `/${code}`
    if (route.path === prefix) return '/'
    if (route.path.startsWith(`${prefix}/`)) return route.path.slice(prefix.length)
  }
  return route.path
})

const appNavSections = computed(() => {
  const sections = [
    {
      key: 'overview',
      label: '',
      items: [
        { label: t('nav.home'), to: '/', icon: 'i-lucide-house', active: route.path === localePath('/') }
      ]
    },
    {
      key: 'catalog',
      label: t('nav.catalog'),
      items: [
        { label: t('catalog.torrents.title'), to: '/catalog/torrents', icon: 'i-lucide-library', active: isCatalogTorrentsActive() },
        { label: t('catalog.requests.navLabel'), to: '/catalog/requests', icon: 'i-lucide-hand-helping', active: isActive('/catalog/requests') },
        { label: t('nav.subtitles'), to: '/catalog/subtitles', icon: 'i-lucide-captions', active: isActive('/catalog/subtitles') },
        { label: t('catalog.bookmarks.title'), to: '/catalog/bookmarks', icon: 'i-lucide-bookmark', active: isActive('/catalog/bookmarks') }
      ]
    },
    {
      key: 'forum',
      label: t('nav.forum'),
      items: [
        { label: t('nav.forum'), to: '/forum', icon: 'i-lucide-messages-square', active: isForumHomeActive() },
        { label: t('forum.sidebar.bookmarks'), to: '/forum/bookmarks', icon: 'i-lucide-book-marked', active: isActive('/forum/bookmarks') }
      ]
    },
    {
      key: 'account',
      label: t('nav.my'),
      items: [
        { label: t('user.nav.overview'), to: '/iam/users/me', icon: 'i-lucide-user-round', active: isExactActive('/iam/users/me') },
        { label: t('user.nav.activity'), to: '/iam/users/me/activity', icon: 'i-lucide-chart-line', active: isActive('/iam/users/me/activity') },
        { label: t('user.nav.bonus'), to: '/iam/users/me/bonus', icon: 'i-lucide-coins', active: isActive('/iam/users/me/bonus') },
        { label: t('user.shop.title'), to: '/economy/shop', icon: 'i-lucide-store', active: isActive('/economy/shop') },
        { label: t('user.nav.settings'), to: '/iam/users/me/settings', icon: 'i-lucide-settings', active: isActive('/iam/users/me/settings') }
      ]
    }
  ]

  if (isStaff.value) {
    sections.push({
      key: 'admin',
      label: t('nav.admin'),
      items: [
        { label: t('nav.admin'), to: '/admin', icon: 'i-lucide-shield-check', active: isActive('/admin') }
      ]
    })
  }

  return sections
})

const adminNavSections = computed<ShellNavSection[]>(() => compactNavSections([
  {
    key: 'overview',
    label: '',
    items: [
      { label: t('admin.nav.overview'), to: '/admin', icon: 'i-lucide-layout-dashboard', active: route.path === localePath('/admin') }
    ]
  },
  {
    key: 'catalog',
    label: t('admin.nav.catalog'),
    items: compactNavItems([
      permissionNavItem(Permission.AdminCatalogCategoryManage, { label: t('admin.catalog.categories.title'), to: '/admin/catalog/categories', icon: 'i-lucide-tags', active: isActive('/admin/catalog/categories') })
    ])
  },
  {
    key: 'forum',
    label: t('admin.nav.forum'),
    items: compactNavItems([
      permissionNavItem(Permission.AdminForumCategoryManage, { label: t('admin.forum.categories.title'), to: '/admin/forum/categories', icon: 'i-lucide-folder-tree', active: isActive('/admin/forum/categories') }),
      permissionNavItem(Permission.AdminForumNodeManage, { label: t('admin.forum.nodes.title'), to: '/admin/forum/nodes', icon: 'i-lucide-panels-top-left', active: isActive('/admin/forum/nodes') })
    ])
  },
  {
    key: 'iam',
    label: t('admin.nav.iam'),
    items: compactNavItems([
      permissionNavItem(Permission.AdminIamUserManage, { label: t('admin.iam.users.title'), to: '/admin/iam/users', icon: 'i-lucide-users', active: isActive('/admin/iam/users') }),
      permissionNavItem(Permission.AdminIamUserManage, { label: t('admin.iam.loginLogs.title'), to: '/admin/iam/login-logs', icon: 'i-lucide-history', active: isActive('/admin/iam/login-logs') }),
      permissionNavItem(Permission.AdminIamRoleManage, { label: t('admin.iam.roles.title'), to: '/admin/iam/roles', icon: 'i-lucide-shield-check', active: isActive('/admin/iam/roles') }),
      permissionNavItem(Permission.AdminIamInviteManage, { label: t('admin.iam.invites.title'), to: '/admin/iam/invites', icon: 'i-lucide-ticket-plus', active: isActive('/admin/iam/invites') })
    ])
  },
  {
    key: 'moderation',
    label: t('admin.nav.moderation'),
    items: compactNavItems([
      permissionNavItem(Permission.AdminModReportManage, { label: t('admin.mod.reports.title'), to: '/admin/mod/reports', icon: 'i-lucide-flag', active: isActive('/admin/mod/reports') }),
      permissionNavItem(Permission.AdminModCheaterManage, { label: t('admin.mod.cheaters.title'), to: '/admin/mod/cheaters', icon: 'i-lucide-radar', active: isActive('/admin/mod/cheaters') })
    ])
  },
  {
    key: 'site',
    label: t('admin.nav.site'),
    items: compactNavItems([
      permissionNavItem(Permission.AdminSiteConfig, { label: t('admin.site.configs.title'), to: '/admin/site/configs', icon: 'i-lucide-settings-2', active: isActive('/admin/site/configs') }),
      permissionNavItem(Permission.AdminSiteAnnouncement, { label: t('admin.site.announcements.title'), to: '/admin/site/announcements', icon: 'i-lucide-megaphone', active: isActive('/admin/site/announcements') }),
      permissionNavItem(Permission.AdminSiteMessage, { label: t('admin.site.messages.title'), to: '/admin/site/messages', icon: 'i-lucide-bell', active: isActive('/admin/site/messages') }),
      permissionNavItem(Permission.AdminSiteAudit, { label: t('admin.site.audits.title'), to: '/admin/site/audits', icon: 'i-lucide-scroll-text', active: isActive('/admin/site/audits') })
    ])
  },
  {
    key: 'system',
    label: t('admin.nav.system'),
    items: compactNavItems([
      permissionNavItem(Permission.AdminSysCronManage, { label: t('admin.sys.crons.title'), to: '/admin/sys/crons', icon: 'i-lucide-clock-3', active: isActive('/admin/sys/crons') })
    ])
  }
]))

const navSections = computed(() => props.mode === 'admin' ? adminNavSections.value : appNavSections.value)
const shellUser = computed(() => user.value
  ? {
      id: user.value.user.id,
      username: user.value.user.username,
      email: user.value.user.email,
      roleName: user.value.role.name,
      avatar: user.value.profile.avatar
    }
  : null
)
const adminReturnAction = computed(() => props.mode === 'admin'
  ? { label: t('nav.backToUser'), to: '/', icon: 'i-lucide-arrow-left' }
  : null
)
const activeItem = computed(() => navSections.value.flatMap((section) => section.items).find((item) => item.active))
const routeSpecificLabel = computed(() => {
  const path = currentAppPath.value
  const torrentDetailMatch = path.match(/^\/catalog\/torrents\/([^/]+)$/)
  if (torrentDetailMatch?.[1]) return t('catalog.torrents.detail.titleFallback', { id: torrentDetailMatch[1] })
  if (path === '/catalog/torrents/upload') return t('catalog.torrents.upload.title')
  if (/^\/catalog\/torrents\/[^/]+\/edit$/.test(path)) return t('catalog.torrents.edit.title')
  if (path === '/catalog/requests/create') return t('catalog.requests.create.title')
  if (/^\/catalog\/requests\/[^/]+$/.test(path)) return t('catalog.requests.detail.title')
  if (path === '/catalog/requests') return t('catalog.requests.title')
  if (path === '/catalog/bookmarks') return t('catalog.bookmarks.title')
  if (path === '/catalog/subtitles') return t('catalog.subtitles.title')
  const forumTopicMatch = path.match(/^\/forum\/topics\/([^/]+)$/)
  if (forumTopicMatch?.[1]) return t('forum.detail.titleFallback', { id: forumTopicMatch[1] })
  if (path === '/forum/topics/create') return t('forum.create.title')
  if (path === '/forum/bookmarks') return t('forum.bookmarks.title')
  if (path === '/site/messages') return t('site.messages.title')
  if (path === '/iam/users/me/activity') return t('user.nav.activity')
  if (path === '/iam/users/me/bonus') return t('user.nav.bonus')
  if (path === '/economy/shop') return t('user.shop.title')
  if (path === '/iam/users/me/settings') return t('user.nav.settings')
  return ''
})
const activeItemLabel = computed(() => routeSpecificLabel.value || activeItem.value?.label || (props.mode === 'admin' ? t('nav.admin') : t('common.brand')))

const languageItems = computed<DropdownMenuItem[][]>(() => [
  (locales.value as any[]).map((item) => ({
    label: item.name,
    onSelect: () => {
      void handleLanguageSwitch(item.code)
    }
  }))
])

const userMenuItems = computed<DropdownMenuItem[][]>(() => {
  const firstGroup: DropdownMenuItem[] = [
    {
      label: t('nav.user'),
      icon: 'i-lucide-user-round',
      onSelect: () => navigateFromUserMenu('/iam/users/me')
    },
    {
      label: t('nav.catalog'),
      icon: 'i-lucide-library',
      onSelect: () => navigateFromUserMenu('/catalog/torrents')
    }
  ]

  if (isStaff.value) {
    firstGroup.push({
      label: t('nav.admin'),
      icon: 'i-lucide-shield-check',
      onSelect: () => navigateFromUserMenu('/admin')
    })
  }

  return [
    firstGroup,
    [
      {
        label: t('auth.logout'),
        icon: 'i-lucide-log-out',
        onSelect: () => {
          void handleLogout()
        }
      }
    ]
  ]
})

async function handleLanguageSwitch(code: string) {
  const supportedCodes = (locales.value as any[]).map((item) => item.code).filter(Boolean)
  const normalized = normalizeLocaleCode(code, supportedCodes, supportedCodes[0] || locale.value)
  if (locale.value === normalized) return
  await setLocale(normalized)
}

function toggleSidebar() {
  sidebarCollapsed.value = !sidebarCollapsed.value
}

function navigateFromUserMenu(path: string) {
  mobileSidebarOpen.value = false
  void navigateTo(localePath(path))
}

function permissionNavItem(permission: string, item: ShellNavItem) {
  return hasPermission(permission) ? item : null
}

function compactNavItems(items: Array<ShellNavItem | null>) {
  return items.filter((item): item is ShellNavItem => Boolean(item))
}

function compactNavSections(sections: ShellNavSection[]) {
  return sections.filter((section) => section.items.length > 0)
}

async function handleLogout() {
  mobileSidebarOpen.value = false
  loggingOut.value = true
  try {
    await logout()
    await navigateTo(localePath('/'))
  } finally {
    loggingOut.value = false
  }
}

async function loadUnreadMessageCount() {
  if (props.mode !== 'app' || !hasPermission(Permission.SiteMessageRead)) {
    unreadMessageCount.value = 0
    return
  }
  try {
    unreadMessageCount.value = await siteApi.unreadMessageCount()
  } catch {
    unreadMessageCount.value = 0
  }
}

function handleMessagesRead() {
  void loadUnreadMessageCount()
}
</script>
