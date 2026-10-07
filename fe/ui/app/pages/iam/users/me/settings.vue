<template>
  <div class="min-h-[calc(100vh-4rem)] bg-neutral-50 py-6 dark:bg-neutral-950">
    <div class="w-full px-3 sm:px-4 lg:px-5">
      <div class="grid gap-4 xl:grid-cols-[260px_minmax(0,1fr)] xl:items-start">
        <aside class="app-sticky-offset xl:sticky">
          <nav class="flex gap-2 overflow-x-auto rounded-lg border border-neutral-200 bg-white p-1 dark:border-neutral-800 dark:bg-neutral-900 xl:block xl:space-y-1 xl:overflow-visible">
            <button
              v-for="section in settingSections"
              :key="section.value"
              type="button"
              class="flex min-w-max items-center gap-2 rounded-md px-3 py-2 text-left text-sm transition xl:w-full xl:min-w-0"
              :class="activeSection === section.value
                ? 'app-selected'
                : 'text-neutral-600 hover:bg-neutral-50 hover:text-neutral-950 dark:text-neutral-300 dark:hover:bg-neutral-950 dark:hover:text-white'"
              :aria-current="activeSection === section.value ? 'page' : undefined"
              @click="setActiveSection(section.value)"
            >
              <span
                class="flex size-8 shrink-0 items-center justify-center rounded-md"
                :class="activeSection === section.value
                  ? 'bg-default text-primary'
                  : 'bg-neutral-100 text-neutral-500 dark:bg-neutral-800 dark:text-neutral-400'"
              >
                <UIcon :name="section.icon" class="size-4" />
              </span>
              <span class="min-w-0">
                <span class="block truncate font-medium">{{ section.label }}</span>
                <span class="mt-0.5 hidden truncate text-xs text-neutral-500 dark:text-neutral-400 xl:block">{{ section.description }}</span>
              </span>
            </button>
          </nav>
        </aside>

        <main class="min-w-0">
          <IamUserAppearance v-if="visitedSections.includes('appearance')" v-show="activeSection === 'appearance'" />
          <div v-if="visitedSections.includes('profile')" v-show="activeSection === 'profile'">
            <IamUserProfileForm />
          </div>

          <div
            v-if="visitedSections.includes('security')"
            v-show="activeSection === 'security'"
            class="grid gap-4"
          >
            <div class="grid gap-4 xl:grid-cols-[minmax(400px,0.85fr)_minmax(0,1.15fr)] xl:items-start">
              <div class="grid gap-4">
                <IamUserTwoStepCard />
                <IamUserPasskeyCard />
              </div>
              <IamUserSecurityForm />
            </div>
            <IamUserSessions />
            <IamLoginLogs />
          </div>

          <div v-if="canReadInvites && visitedSections.includes('invites')" v-show="activeSection === 'invites'">
            <IamUserInvitesPanel />
          </div>
        </main>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
type SettingSection = 'profile' | 'security' | 'invites' | 'appearance'

definePageMeta({
  middleware: 'auth'
})

const { t } = useI18n()
const route = useRoute()
const router = useRouter()
const { hasPermission } = useAuth()

const canReadInvites = computed(() => hasPermission(Permission.IamInviteRead))
const settingSections = computed(() => compactSettingSections([
  {
    value: 'profile' as const,
    label: t('user.settings.sections.profile.label'),
    description: t('user.settings.sections.profile.description'),
    icon: 'i-lucide-user'
  },
  {
    value: 'security' as const,
    label: t('user.settings.sections.security.label'),
    description: t('user.settings.sections.security.description'),
    icon: 'i-lucide-shield-check'
  },
  {
    value: 'appearance' as const,
    label: t('appearance.title'),
    description: t('appearance.description'),
    icon: 'i-lucide-palette'
  },
  canReadInvites.value ? {
    value: 'invites' as const,
    label: t('user.settings.sections.invites.label'),
    description: t('user.settings.sections.invites.description'),
    icon: 'i-lucide-ticket'
  } : null
]))

const activeSection = ref<SettingSection>(normalizeSection(route.query.section))
const visitedSections = ref<SettingSection[]>([activeSection.value])

watch(
  () => route.query.section,
  (section) => {
    setActiveSectionValue(normalizeSection(section))
  }
)

function normalizeSection(section: unknown): SettingSection {
  if (section === 'appearance') return 'appearance'
  if (section === 'security') return 'security'
  if (section === 'invites' && canReadInvites.value) return 'invites'
  return 'profile'
}

function setActiveSection(section: SettingSection) {
  if (activeSection.value === section) return

  const query = { ...route.query }
  if (section === 'profile') {
    delete query.section
  } else {
    query.section = section
  }
  router.replace({ query })
}

function setActiveSectionValue(section: SettingSection) {
  activeSection.value = section
  if (!visitedSections.value.includes(section)) {
    visitedSections.value.push(section)
  }
}

watch(canReadInvites, (canRead) => {
  if (!canRead && activeSection.value === 'invites') {
    setActiveSection('profile')
  }
})

function compactSettingSections(sections: Array<{
  value: SettingSection
  label: string
  description: string
  icon: string
} | null>) {
  return sections.filter((section): section is NonNullable<typeof section> => Boolean(section))
}

useSeoMeta({
  title: t('user.nav.settings'),
  robots: 'noindex, nofollow'
})
</script>
