<template>
  <div class="min-h-[calc(100vh-4rem)] bg-slate-50 py-6 dark:bg-slate-950">
    <div class="w-full px-3 sm:px-4 lg:px-5">
      <div class="grid gap-4 xl:grid-cols-[260px_minmax(0,1fr)] xl:items-start">
        <aside class="xl:sticky xl:top-20">
          <nav class="flex gap-2 overflow-x-auto rounded-lg border border-slate-200 bg-white p-1 dark:border-slate-800 dark:bg-slate-900 xl:block xl:space-y-1 xl:overflow-visible">
            <button
              v-for="section in settingSections"
              :key="section.value"
              type="button"
              class="flex min-w-max items-center gap-2 rounded-md px-3 py-2 text-left text-sm transition xl:w-full xl:min-w-0"
              :class="activeSection === section.value
                ? 'bg-slate-100 text-slate-950 dark:bg-slate-800 dark:text-white'
                : 'text-slate-600 hover:bg-slate-50 hover:text-slate-950 dark:text-slate-300 dark:hover:bg-slate-950 dark:hover:text-white'"
              :aria-current="activeSection === section.value ? 'page' : undefined"
              @click="setActiveSection(section.value)"
            >
              <span
                class="flex size-8 shrink-0 items-center justify-center rounded-md"
                :class="activeSection === section.value
                  ? 'bg-white text-slate-950 dark:bg-slate-950 dark:text-white'
                  : 'bg-slate-100 text-slate-500 dark:bg-slate-800 dark:text-slate-400'"
              >
                <UIcon :name="section.icon" class="size-4" />
              </span>
              <span class="min-w-0">
                <span class="block truncate font-medium">{{ section.label }}</span>
                <span class="mt-0.5 hidden truncate text-xs text-slate-500 dark:text-slate-400 xl:block">{{ section.description }}</span>
              </span>
            </button>
          </nav>
        </aside>

        <main class="min-w-0">
          <div v-if="activeSection === 'profile'">
            <IamUserProfileForm />
          </div>

          <div v-else-if="activeSection === 'security'" class="space-y-4">
            <IamUserPasskeyCard />
            <IamUserSecurityForm />
          </div>

          <div v-else>
            <IamUserInvitesPanel />
          </div>
        </main>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
type SettingSection = 'profile' | 'security' | 'invites'

definePageMeta({
  middleware: 'auth'
})

const { t } = useI18n()
const route = useRoute()
const router = useRouter()

const settingSections = computed(() => [
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
    value: 'invites' as const,
    label: t('user.settings.sections.invites.label'),
    description: t('user.settings.sections.invites.description'),
    icon: 'i-lucide-ticket'
  }
])

const activeSection = ref<SettingSection>(normalizeSection(route.query.section))

watch(
  () => route.query.section,
  (section) => {
    activeSection.value = normalizeSection(section)
  }
)

function normalizeSection(section: unknown): SettingSection {
  return section === 'security' || section === 'invites' ? section : 'profile'
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

useSeoMeta({
  title: t('user.nav.settings'),
  robots: 'noindex, nofollow'
})
</script>
