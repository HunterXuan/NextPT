<template>
  <div>
    <template v-if="!isLoggedIn">
      <section class="relative isolate overflow-hidden bg-slate-950 text-white">
        <img
          src="/images/home-hero.png"
          alt=""
          class="absolute inset-0 -z-10 h-full w-full object-cover"
        >
        <div class="absolute inset-0 -z-10 bg-slate-950/62" />

        <div class="mx-auto flex min-h-[calc(100svh-9rem)] max-w-7xl items-center px-4 py-16 sm:px-6 lg:min-h-[calc(100svh-4rem)] lg:px-8 lg:py-20">
          <div class="max-w-3xl">
            <p class="inline-flex items-center gap-2 rounded-md border border-white/15 bg-white/10 px-3 py-1.5 text-xs font-medium text-sky-100 backdrop-blur">
              <UIcon name="i-lucide-sparkles" class="size-4 text-amber-300" />
              {{ $t('home.guest.eyebrow') }}
            </p>
            <h1 class="mt-5 text-4xl font-semibold leading-tight sm:text-5xl lg:text-6xl">
              {{ $t('home.guest.title') }}
            </h1>
            <p class="mt-5 max-w-2xl text-base leading-7 text-slate-200 sm:text-lg">
              {{ $t('home.guest.description') }}
            </p>

            <div class="mt-8 flex flex-col gap-3 sm:flex-row">
              <UButton color="primary" size="lg" icon="i-lucide-user-plus" :to="localePath('/register')">
                {{ $t('home.guest.primaryAction') }}
              </UButton>
              <UButton color="neutral" variant="outline" size="lg" icon="i-lucide-log-in" :to="localePath('/login')">
                {{ $t('home.guest.secondaryAction') }}
              </UButton>
            </div>

            <dl class="mt-10 grid max-w-2xl grid-cols-1 gap-3 sm:grid-cols-3">
              <div v-for="item in guestStats" :key="item.label" class="rounded-lg border border-white/12 bg-white/10 px-4 py-3 backdrop-blur">
                <dt class="text-xs font-medium text-slate-300">{{ item.label }}</dt>
                <dd class="mt-1 text-lg font-semibold text-white">{{ item.value }}</dd>
              </div>
            </dl>
          </div>
        </div>
      </section>

      <section class="bg-white py-14 dark:bg-slate-950">
        <div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
          <div class="max-w-2xl">
            <p class="text-sm font-medium text-sky-600 dark:text-sky-400">{{ $t('home.guest.features.eyebrow') }}</p>
            <h2 class="mt-2 text-2xl font-semibold text-slate-950 dark:text-white">{{ $t('home.guest.features.title') }}</h2>
            <p class="mt-3 text-sm leading-6 text-slate-600 dark:text-slate-300">{{ $t('home.guest.features.description') }}</p>
          </div>

          <div class="mt-8 grid gap-4 md:grid-cols-2 xl:grid-cols-3">
            <article v-for="feature in guestFeatures" :key="feature.title" class="rounded-lg border border-slate-200 bg-white p-5 dark:border-slate-800 dark:bg-slate-900">
              <span class="flex size-10 items-center justify-center rounded-md" :class="feature.iconClass">
                <UIcon :name="feature.icon" class="size-5" />
              </span>
              <h3 class="mt-4 text-base font-semibold text-slate-950 dark:text-white">{{ feature.title }}</h3>
              <p class="mt-2 text-sm leading-6 text-slate-600 dark:text-slate-300">{{ feature.description }}</p>
            </article>
          </div>
        </div>
      </section>

      <section class="border-y border-slate-200 bg-slate-50 py-14 dark:border-slate-800 dark:bg-slate-900/40">
        <div class="mx-auto grid max-w-7xl gap-10 px-4 sm:px-6 lg:grid-cols-[360px_minmax(0,1fr)] lg:px-8">
          <div>
            <p class="text-sm font-medium text-sky-600 dark:text-sky-400">{{ $t('home.guest.flow.eyebrow') }}</p>
            <h2 class="mt-2 text-2xl font-semibold text-slate-950 dark:text-white">{{ $t('home.guest.flow.title') }}</h2>
            <p class="mt-3 text-sm leading-6 text-slate-600 dark:text-slate-300">{{ $t('home.guest.flow.description') }}</p>
          </div>

          <div class="grid gap-3 md:grid-cols-2">
            <article v-for="(step, index) in guestFlow" :key="step.title" class="rounded-lg border border-slate-200 bg-white p-5 dark:border-slate-800 dark:bg-slate-950">
              <p class="text-xs font-semibold text-sky-600 dark:text-sky-400">{{ numberFormatter.format(index + 1) }}</p>
              <h3 class="mt-2 text-base font-semibold text-slate-950 dark:text-white">{{ step.title }}</h3>
              <p class="mt-2 text-sm leading-6 text-slate-600 dark:text-slate-300">{{ step.description }}</p>
            </article>
          </div>
        </div>
      </section>

      <section class="bg-white py-14 dark:bg-slate-950">
        <div class="mx-auto grid max-w-7xl gap-10 px-4 sm:px-6 lg:grid-cols-[minmax(0,1fr)_420px] lg:px-8">
          <div>
            <p class="text-sm font-medium text-sky-600 dark:text-sky-400">{{ $t('home.guest.faq.eyebrow') }}</p>
            <h2 class="mt-2 text-2xl font-semibold text-slate-950 dark:text-white">{{ $t('home.guest.faq.title') }}</h2>
            <div class="mt-6 divide-y divide-slate-200 rounded-lg border border-slate-200 dark:divide-slate-800 dark:border-slate-800">
              <details v-for="item in guestFaqs" :key="item.question" class="group px-4 py-4">
                <summary class="flex cursor-pointer list-none items-center justify-between gap-4 text-sm font-semibold text-slate-950 dark:text-white">
                  {{ item.question }}
                  <UIcon name="i-lucide-chevron-down" class="size-4 shrink-0 text-slate-400 transition group-open:rotate-180" />
                </summary>
                <p class="mt-3 text-sm leading-6 text-slate-600 dark:text-slate-300">{{ item.answer }}</p>
              </details>
            </div>
          </div>

          <aside class="rounded-lg border border-slate-200 bg-slate-50 p-6 dark:border-slate-800 dark:bg-slate-900">
            <UIcon name="i-lucide-ticket-plus" class="size-8 text-sky-500" />
            <h2 class="mt-5 text-xl font-semibold text-slate-950 dark:text-white">{{ $t('home.guest.cta.title') }}</h2>
            <p class="mt-3 text-sm leading-6 text-slate-600 dark:text-slate-300">{{ $t('home.guest.cta.description') }}</p>
            <div class="mt-6 grid gap-3">
              <UButton color="primary" icon="i-lucide-user-plus" block :to="localePath('/register')">
                {{ $t('home.guest.primaryAction') }}
              </UButton>
              <UButton color="neutral" variant="outline" icon="i-lucide-log-in" block :to="localePath('/login')">
                {{ $t('home.guest.secondaryAction') }}
              </UButton>
            </div>
          </aside>
        </div>
      </section>
    </template>

    <template v-else>
      <section class="min-h-[calc(100vh-4rem)] bg-slate-50 py-6 dark:bg-slate-950">
        <div class="mx-auto max-w-7xl px-3 sm:px-4 lg:px-5">
          <div class="grid gap-4 xl:grid-cols-[minmax(0,1fr)_340px] xl:items-start">
            <main class="min-w-0 space-y-4">
              <section class="overflow-hidden rounded-lg border border-slate-200 bg-white dark:border-slate-800 dark:bg-slate-900">
                <div class="flex items-center justify-between gap-3 border-b border-slate-200 px-4 py-3 dark:border-slate-800">
                  <h2 class="text-sm font-semibold text-slate-950 dark:text-white">{{ $t('site.announcements.title') }}</h2>
                  <span v-if="unreadAnnouncementCount > 0" class="rounded-full bg-sky-50 px-2 py-0.5 text-xs font-medium text-sky-700 dark:bg-sky-950 dark:text-sky-300">
                    {{ numberFormatter.format(unreadAnnouncementCount) }}
                  </span>
                </div>

                <div v-if="announcementsPending" class="divide-y divide-slate-100 dark:divide-slate-800">
                  <div v-for="item in 3" :key="item" class="px-4 py-3">
                    <div class="h-4 w-2/5 animate-pulse rounded bg-slate-100 dark:bg-slate-800" />
                    <div class="mt-2 h-3 w-4/5 animate-pulse rounded bg-slate-100 dark:bg-slate-800" />
                  </div>
                </div>
                <div v-else-if="announcementError" class="px-4 py-4 text-sm text-red-600 dark:text-red-300">
                  {{ announcementError }}
                </div>
                <div v-else-if="announcements.length === 0" class="px-4 py-8 text-center text-sm text-slate-500 dark:text-slate-400">
                  {{ $t('site.announcements.empty') }}
                </div>
                <div v-else class="divide-y divide-slate-100 dark:divide-slate-800">
                  <article
                    v-for="item in announcements"
                    :key="item.id"
                    class="flex items-start gap-2 px-4 py-3 transition hover:bg-slate-50 dark:hover:bg-slate-950/60"
                  >
                    <div
                      role="button"
                      tabindex="0"
                      class="group flex min-w-0 flex-1 cursor-pointer gap-3 text-left focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-sky-300 dark:focus-visible:ring-sky-700"
                      @click="openAnnouncement(item)"
                      @keydown.enter.prevent="openAnnouncement(item)"
                      @keydown.space.prevent="openAnnouncement(item)"
                    >
                      <span class="mt-2 size-2 shrink-0 rounded-full transition" :class="item.isRead ? 'bg-slate-300 dark:bg-slate-700' : 'bg-sky-500'" />
                      <div class="min-w-0 flex-1">
                        <div class="flex min-w-0 flex-col gap-1 sm:flex-row sm:items-start sm:justify-between sm:gap-3">
                          <h3 class="truncate text-sm font-semibold text-slate-950 transition group-hover:text-sky-700 dark:text-white dark:group-hover:text-sky-300">{{ item.title }}</h3>
                          <p class="shrink-0 text-xs text-slate-500 dark:text-slate-400">{{ formatRelativeDateTime(item.publishedAt || item.createdAt, locale) }}</p>
                        </div>
                        <div
                          v-if="renderAnnouncementContent(item.content)"
                          class="rich-text rich-text-compact mt-2 max-h-20 overflow-hidden text-sm leading-6 text-slate-600 dark:text-slate-300"
                          v-html="renderAnnouncementContent(item.content)"
                        />
                      </div>
                    </div>
                    <UTooltip v-if="!item.isRead" :text="$t('site.announcements.markRead')" :content="{ side: 'top', sideOffset: 8 }" :delay-duration="300">
                      <UButton class="shrink-0" color="neutral" variant="ghost" size="xs" icon="i-lucide-check" :loading="readingAnnouncementId === item.id" :aria-label="$t('site.announcements.markRead')" @click.stop="markAnnouncementRead(item)" />
                    </UTooltip>
                  </article>
                </div>
              </section>

              <section class="overflow-hidden rounded-lg border border-slate-200 bg-white dark:border-slate-800 dark:bg-slate-900">
                <div class="flex flex-col gap-3 border-b border-slate-200 px-4 py-3 sm:flex-row sm:items-center sm:justify-between dark:border-slate-800">
                  <div>
                    <h2 class="text-sm font-semibold text-slate-950 dark:text-white">{{ $t('home.member.torrents.title') }}</h2>
                  </div>
                  <UButton color="neutral" variant="ghost" size="xs" :to="localePath('/catalog/torrents')">
                    {{ $t('home.member.viewAll') }}
                  </UButton>
                </div>
                <div v-if="torrentsPending" class="divide-y divide-slate-100 dark:divide-slate-800">
                  <div v-for="item in 5" :key="item" class="px-4 py-3">
                    <div class="h-4 w-3/4 animate-pulse rounded bg-slate-100 dark:bg-slate-800" />
                    <div class="mt-2 h-3 w-1/2 animate-pulse rounded bg-slate-100 dark:bg-slate-800" />
                  </div>
                </div>
                <div v-else-if="torrentsError" class="px-4 py-4 text-sm text-red-600 dark:text-red-300">{{ torrentsError }}</div>
                <div v-else-if="torrents.length === 0" class="px-4 py-8 text-center text-sm text-slate-500 dark:text-slate-400">{{ $t('home.member.empty') }}</div>
                <div v-else>
                  <NuxtLink v-for="torrent in torrents" :key="torrent.id" :to="localePath(`/catalog/torrents/${torrent.id}`)" class="grid gap-3 border-t border-slate-100 px-4 py-3 outline-none transition-colors first:border-t-0 hover:bg-slate-50 focus-visible:ring-2 focus-visible:ring-inset focus-visible:ring-sky-300 sm:grid-cols-[minmax(0,1fr)_auto] sm:items-center dark:border-slate-800 dark:hover:bg-slate-950/60 dark:focus-visible:ring-sky-700">
                    <div class="min-w-0">
                      <div class="flex min-w-0 items-center gap-2">
                        <span class="inline-flex h-5 max-w-24 shrink-0 items-center rounded border border-slate-200 bg-slate-50 px-1.5 text-[10px] font-medium text-slate-500 dark:border-slate-700 dark:bg-slate-950 dark:text-slate-400">{{ torrentCategoryName(torrent) }}</span>
                        <p class="truncate text-sm font-medium text-slate-950 dark:text-white">{{ torrent.name || `#${torrent.id}` }}</p>
                      </div>
                      <p class="mt-1 truncate text-xs text-slate-500 dark:text-slate-400">{{ formatRelativeDateTime(torrent.createdAt, locale) }} / {{ formatBytes(torrent.size) }}</p>
                    </div>
                    <div class="flex items-center gap-3 text-xs sm:justify-end">
                      <UTooltip v-for="metric in torrentMetrics(torrent)" :key="metric.key" :text="metric.label" :content="{ side: 'top', sideOffset: 8 }" :delay-duration="300">
                        <span class="inline-flex items-center gap-1 font-medium" :class="metric.class">
                          <UIcon :name="metric.icon" class="size-3.5" />
                          <span class="tabular-nums">{{ numberFormatter.format(metric.value) }}</span>
                        </span>
                      </UTooltip>
                    </div>
                  </NuxtLink>
                </div>
              </section>

              <section class="overflow-hidden rounded-lg border border-slate-200 bg-white dark:border-slate-800 dark:bg-slate-900">
                <div class="flex flex-col gap-3 border-b border-slate-200 px-4 py-3 sm:flex-row sm:items-center sm:justify-between dark:border-slate-800">
                  <div>
                    <h2 class="text-sm font-semibold text-slate-950 dark:text-white">{{ $t('home.member.forum.title') }}</h2>
                  </div>
                  <UButton color="neutral" variant="ghost" size="xs" :to="localePath('/forum')">
                    {{ $t('home.member.viewAll') }}
                  </UButton>
                </div>
                <div v-if="topicsPending" class="divide-y divide-slate-100 dark:divide-slate-800">
                  <div v-for="item in 5" :key="item" class="px-4 py-3">
                    <div class="h-4 w-3/4 animate-pulse rounded bg-slate-100 dark:bg-slate-800" />
                    <div class="mt-2 h-3 w-1/2 animate-pulse rounded bg-slate-100 dark:bg-slate-800" />
                  </div>
                </div>
                <div v-else-if="topicsError" class="px-4 py-4 text-sm text-red-600 dark:text-red-300">{{ topicsError }}</div>
                <div v-else-if="dashboardTopics.length === 0" class="px-4 py-8 text-center text-sm text-slate-500 dark:text-slate-400">{{ $t('home.member.empty') }}</div>
                <div v-else>
                  <NuxtLink v-for="topic in dashboardTopics" :key="topic.id" :to="localePath(`/forum/topics/${topic.id}`)" class="grid gap-3 border-t border-slate-100 px-4 py-3 outline-none transition-colors first:border-t-0 hover:bg-slate-50 focus-visible:ring-2 focus-visible:ring-inset focus-visible:ring-sky-300 sm:grid-cols-[minmax(0,1fr)_auto] sm:items-center dark:border-slate-800 dark:hover:bg-slate-950/60 dark:focus-visible:ring-sky-700">
                    <div class="min-w-0">
                      <p class="truncate text-sm font-medium text-slate-950 dark:text-white">{{ topic.subject || `#${topic.id}` }}</p>
                      <p class="mt-1 truncate text-xs text-slate-500 dark:text-slate-400">{{ topicNodeName(topic) }} / {{ topic.author?.username || '-' }} / {{ formatRelativeDateTime(topicLastActivityAt(topic), locale) }}</p>
                    </div>
                    <div class="flex items-center gap-3 text-xs sm:justify-end">
                      <UTooltip :text="$t('home.member.metrics.views')" :content="{ side: 'top', sideOffset: 8 }" :delay-duration="300">
                        <span class="inline-flex items-center gap-1 font-medium text-slate-500 dark:text-slate-400">
                          <UIcon name="i-lucide-eye" class="size-3.5" />
                          <span class="tabular-nums">{{ numberFormatter.format(topic.views) }}</span>
                        </span>
                      </UTooltip>
                      <UTooltip :text="$t('home.member.metrics.replies')" :content="{ side: 'top', sideOffset: 8 }" :delay-duration="300">
                        <span class="inline-flex items-center gap-1 font-medium text-slate-700 dark:text-slate-200">
                          <UIcon name="i-lucide-message-square" class="size-3.5" />
                          <span class="tabular-nums">{{ numberFormatter.format(topic.replyCount) }}</span>
                        </span>
                      </UTooltip>
                    </div>
                  </NuxtLink>
                </div>
              </section>
            </main>

            <aside>
              <section class="rounded-lg border border-slate-200 bg-white p-4 dark:border-slate-800 dark:bg-slate-900">
                <h2 class="text-sm font-semibold text-slate-950 dark:text-white">{{ $t('home.member.mySite') }}</h2>
                <div class="mt-4 flex items-center gap-3">
                  <IamUserAvatar :id="user?.user.id" :username="user?.user.username" :avatar="user?.profile.avatar" size="lg" />
                  <div class="min-w-0">
                    <p class="truncate text-sm font-semibold text-slate-950 dark:text-white">{{ user?.user.username || '-' }}</p>
                    <p class="mt-0.5 truncate text-xs text-slate-500 dark:text-slate-400">{{ user?.role.name || '-' }} / #{{ user?.user.id || '-' }}</p>
                  </div>
                </div>

                <div class="mt-4 grid grid-cols-2 gap-2">
                  <div v-for="item in memberStatCards" :key="item.key" class="rounded-md border border-slate-200 px-3 py-2 dark:border-slate-800">
                    <div class="flex items-center justify-between gap-2">
                      <p class="text-xs text-slate-500 dark:text-slate-400">{{ item.label }}</p>
                      <UIcon :name="item.icon" class="size-3.5 text-slate-400" />
                    </div>
                    <p class="mt-1 truncate text-sm font-semibold tabular-nums text-slate-950 dark:text-white">{{ item.value }}</p>
                  </div>
                </div>

                <div class="mt-4 border-t border-slate-200 pt-4 dark:border-slate-800">
                  <h3 class="text-xs font-medium text-slate-500 dark:text-slate-400">{{ $t('home.member.quickActions') }}</h3>
                  <div class="mt-3 grid grid-cols-2 gap-2">
                    <NuxtLink v-for="action in memberActions" :key="action.to" :to="localePath(action.to)" class="rounded-md border border-slate-200 px-3 py-2 text-sm transition hover:border-sky-300 hover:bg-sky-50/70 dark:border-slate-800 dark:hover:border-sky-800 dark:hover:bg-sky-950/30">
                      <div class="flex items-center gap-2">
                        <span class="flex size-7 shrink-0 items-center justify-center rounded-md" :class="action.iconClass">
                          <UIcon :name="action.icon" class="size-4" />
                        </span>
                        <span class="min-w-0 flex-1 truncate font-medium text-slate-700 dark:text-slate-200">{{ action.title }}</span>
                      </div>
                    </NuxtLink>
                  </div>
                </div>
              </section>
            </aside>
          </div>
        </div>
      </section>
    </template>

    <UModal
      :open="announcementModalOpen"
      :title="selectedAnnouncement?.title || $t('site.announcements.title')"
      :description="selectedAnnouncementDescription"
      :ui="{
        content: 'sm:max-w-3xl overflow-hidden',
        header: 'min-h-0 px-5 py-4 sm:px-5',
        body: 'p-0 sm:p-0',
        title: 'text-base font-semibold text-slate-950 dark:text-white',
        description: 'mt-1 text-sm text-slate-500 dark:text-slate-400',
        close: 'top-4 end-4'
      }"
      @update:open="setAnnouncementModalOpen"
    >
      <template #body>
        <div class="max-h-[70vh] overflow-auto px-5 py-4">
          <div
            v-if="selectedAnnouncementContent"
            class="rich-text"
            v-html="selectedAnnouncementContent"
          />
        </div>
      </template>
    </UModal>
  </div>
</template>

<script setup lang="ts">
import { useCatalogTorrents, type TorrentHotItem } from '~/composables/useCatalogTorrents'
import { ApiError } from '~/composables/useApi'
import type { SiteAnnouncement } from '~/composables/useSite'
import { useForum, type ForumTopicHotItem } from '~/composables/useForum'
import { formatBytes, formatDateTime, formatRelativeDateTime, localizeI18nName } from '~/utils/format'
import { renderUserMarkdown } from '~/utils/richText'

const { t, locale } = useI18n()
const localePath = useLocalePath()
const { isLoggedIn, user } = useAuth()
const siteApi = useSite()
const catalogApi = useCatalogTorrents()
const forumApi = useForum()

const numberFormatter = computed(() => new Intl.NumberFormat(locale.value))
const announcements = ref<SiteAnnouncement[]>([])
const announcementsPending = ref(false)
const announcementError = ref('')
const readingAnnouncementId = ref(0)
const announcementModalOpen = ref(false)
const selectedAnnouncement = ref<SiteAnnouncement | null>(null)
const torrents = ref<TorrentHotItem[]>([])
const torrentsPending = ref(false)
const torrentsError = ref('')
const dashboardTopics = ref<ForumTopicHotItem[]>([])
const topicsPending = ref(false)
const topicsError = ref('')

onMounted(() => {
  if (isLoggedIn.value) {
    void loadMemberDashboard()
  }
})

watch(isLoggedIn, (value) => {
  if (value) {
    void loadMemberDashboard()
    return
  }
  announcements.value = []
  torrents.value = []
  dashboardTopics.value = []
  announcementError.value = ''
  torrentsError.value = ''
  topicsError.value = ''
})

const guestStats = computed(() => [
  { label: t('home.guest.stats.privacy'), value: t('home.guest.stats.privacyValue') },
  { label: t('home.guest.stats.ratio'), value: t('home.guest.stats.ratioValue') },
  { label: t('home.guest.stats.community'), value: t('home.guest.stats.communityValue') }
])

const guestFeatures = computed(() => [
  {
    title: t('home.guest.featureItems.catalog.title'),
    description: t('home.guest.featureItems.catalog.description'),
    icon: 'i-lucide-library',
    iconClass: 'bg-sky-100 text-sky-700 dark:bg-sky-950 dark:text-sky-300'
  },
  {
    title: t('home.guest.featureItems.tracker.title'),
    description: t('home.guest.featureItems.tracker.description'),
    icon: 'i-lucide-radio-tower',
    iconClass: 'bg-emerald-100 text-emerald-700 dark:bg-emerald-950 dark:text-emerald-300'
  },
  {
    title: t('home.guest.featureItems.ratio.title'),
    description: t('home.guest.featureItems.ratio.description'),
    icon: 'i-lucide-scale',
    iconClass: 'bg-amber-100 text-amber-700 dark:bg-amber-950 dark:text-amber-300'
  },
  {
    title: t('home.guest.featureItems.forum.title'),
    description: t('home.guest.featureItems.forum.description'),
    icon: 'i-lucide-messages-square',
    iconClass: 'bg-cyan-100 text-cyan-700 dark:bg-cyan-950 dark:text-cyan-300'
  },
  {
    title: t('home.guest.featureItems.subtitle.title'),
    description: t('home.guest.featureItems.subtitle.description'),
    icon: 'i-lucide-captions',
    iconClass: 'bg-indigo-100 text-indigo-700 dark:bg-indigo-950 dark:text-indigo-300'
  },
  {
    title: t('home.guest.featureItems.account.title'),
    description: t('home.guest.featureItems.account.description'),
    icon: 'i-lucide-user-check',
    iconClass: 'bg-rose-100 text-rose-700 dark:bg-rose-950 dark:text-rose-300'
  }
])

const guestFlow = computed(() => [
  {
    title: t('home.guest.flow.items.invite.title'),
    description: t('home.guest.flow.items.invite.description')
  },
  {
    title: t('home.guest.flow.items.discover.title'),
    description: t('home.guest.flow.items.discover.description')
  },
  {
    title: t('home.guest.flow.items.seed.title'),
    description: t('home.guest.flow.items.seed.description')
  },
  {
    title: t('home.guest.flow.items.discuss.title'),
    description: t('home.guest.flow.items.discuss.description')
  }
])

const guestFaqs = computed(() => [
  {
    question: t('home.guest.faq.items.invite.question'),
    answer: t('home.guest.faq.items.invite.answer')
  },
  {
    question: t('home.guest.faq.items.ratio.question'),
    answer: t('home.guest.faq.items.ratio.answer')
  },
  {
    question: t('home.guest.faq.items.passkey.question'),
    answer: t('home.guest.faq.items.passkey.answer')
  },
  {
    question: t('home.guest.faq.items.forum.question'),
    answer: t('home.guest.faq.items.forum.answer')
  }
])

const memberActions = computed(() => [
  {
    title: t('home.member.actions.catalog.title'),
    description: t('home.member.actions.catalog.description'),
    to: '/catalog/torrents',
    icon: 'i-lucide-library',
    iconClass: 'bg-sky-100 text-sky-700 dark:bg-sky-950 dark:text-sky-300'
  },
  {
    title: t('home.member.actions.upload.title'),
    description: t('home.member.actions.upload.description'),
    to: '/catalog/torrents/upload',
    icon: 'i-lucide-upload',
    iconClass: 'bg-emerald-100 text-emerald-700 dark:bg-emerald-950 dark:text-emerald-300'
  },
  {
    title: t('home.member.actions.forum.title'),
    description: t('home.member.actions.forum.description'),
    to: '/forum',
    icon: 'i-lucide-messages-square',
    iconClass: 'bg-cyan-100 text-cyan-700 dark:bg-cyan-950 dark:text-cyan-300'
  },
  {
    title: t('home.member.actions.profile.title'),
    description: t('home.member.actions.profile.description'),
    to: '/iam/users/me',
    icon: 'i-lucide-user-round',
    iconClass: 'bg-amber-100 text-amber-700 dark:bg-amber-950 dark:text-amber-300'
  }
])

const memberStatCards = computed(() => {
  const stats = user.value?.stat

  return [
    {
      key: 'uploaded',
      label: t('home.member.stats.uploaded'),
      value: formatBytes(stats?.uploaded),
      icon: 'i-lucide-upload',
      iconClass: 'bg-emerald-100 text-emerald-700 dark:bg-emerald-950 dark:text-emerald-300'
    },
    {
      key: 'downloaded',
      label: t('home.member.stats.downloaded'),
      value: formatBytes(stats?.downloaded),
      icon: 'i-lucide-download',
      iconClass: 'bg-sky-100 text-sky-700 dark:bg-sky-950 dark:text-sky-300'
    },
    {
      key: 'ratio',
      label: t('home.member.stats.ratio'),
      value: formatRatio(stats?.shareRatio),
      icon: 'i-lucide-scale',
      iconClass: 'bg-amber-100 text-amber-700 dark:bg-amber-950 dark:text-amber-300'
    },
    {
      key: 'bonus',
      label: t('home.member.stats.bonus'),
      value: numberFormatter.value.format(Number(stats?.bonus || 0)),
      icon: 'i-lucide-coins',
      iconClass: 'bg-rose-100 text-rose-700 dark:bg-rose-950 dark:text-rose-300'
    }
  ]
})

const selectedAnnouncementContent = computed(() => renderAnnouncementContent(selectedAnnouncement.value?.content || ''))
const selectedAnnouncementDescription = computed(() => {
  const item = selectedAnnouncement.value
  if (!item) return ''
  return formatDateTime(item.publishedAt || item.createdAt, locale.value)
})
const unreadAnnouncementCount = computed(() => announcements.value.filter(item => !item.isRead).length)

async function loadMemberDashboard() {
  await Promise.allSettled([
    loadAnnouncements(),
    loadTorrentDashboard(),
    loadForumDashboard()
  ])
}

async function loadAnnouncements() {
  announcementsPending.value = true
  announcementError.value = ''
  try {
    const data = await siteApi.listAnnouncements({ page: 1, size: 5 })
    announcements.value = data.list || []
  } catch (error) {
    announcements.value = []
    announcementError.value = error instanceof ApiError ? error.message : t('common.requestFailed')
  } finally {
    announcementsPending.value = false
  }
}

async function markAnnouncementRead(item: SiteAnnouncement) {
  if (item.isRead) return

  readingAnnouncementId.value = item.id
  try {
    await siteApi.markAnnouncementRead(item.id)
    item.isRead = true
  } finally {
    readingAnnouncementId.value = 0
  }
}

function openAnnouncement(item: SiteAnnouncement) {
  selectedAnnouncement.value = item
  announcementModalOpen.value = true
  if (!item.isRead) {
    void markAnnouncementRead(item)
  }
}

function setAnnouncementModalOpen(value: boolean) {
  announcementModalOpen.value = value
}

async function loadTorrentDashboard() {
  torrentsPending.value = true
  torrentsError.value = ''
  try {
    const torrentData = await catalogApi.listHotTorrents(5)
    torrents.value = torrentData.list || []
  } catch (error) {
    torrents.value = []
    torrentsError.value = error instanceof ApiError ? error.message : t('common.requestFailed')
  } finally {
    torrentsPending.value = false
  }
}

async function loadForumDashboard() {
  topicsPending.value = true
  topicsError.value = ''
  try {
    const data = await forumApi.listHotTopics(5)
    dashboardTopics.value = data.list || []
  } catch (error) {
    dashboardTopics.value = []
    topicsError.value = error instanceof ApiError ? error.message : t('common.requestFailed')
  } finally {
    topicsPending.value = false
  }
}

function renderAnnouncementContent(content: string) {
  return renderUserMarkdown(content || '').trim()
}

function torrentCategoryName(torrent: TorrentHotItem) {
  const category = torrent.category
  return category ? localizeI18nName(category.name, locale.value, category.slug || `#${category.id}`) : `#${torrent.categoryId}`
}

function topicNodeName(topic: ForumTopicHotItem) {
  return localizeI18nName(topic.node?.nameI18n, locale.value, topic.node?.slug || `#${topic.node?.id || ''}`)
}

function torrentMetrics(torrent: TorrentHotItem) {
  return [
    {
      key: 'snatched',
      label: t('home.member.metrics.snatched'),
      value: Number(torrent.snatched || 0),
      icon: 'i-lucide-check',
      class: 'text-slate-500 dark:text-slate-400'
    },
    {
      key: 'seeders',
      label: t('home.member.metrics.seeders'),
      value: Number(torrent.seeders || 0),
      icon: 'i-lucide-upload',
      class: 'text-emerald-600 dark:text-emerald-400'
    },
    {
      key: 'leechers',
      label: t('home.member.metrics.leechers'),
      value: Number(torrent.leechers || 0),
      icon: 'i-lucide-download',
      class: 'text-sky-600 dark:text-sky-400'
    }
  ]
}

function topicLastActivityAt(topic: ForumTopicHotItem) {
  return Number(topic.replyCount || 0) > 0 && topic.lastReplyAt ? topic.lastReplyAt : topic.createdAt
}

function formatRatio(value?: number | null) {
  const ratio = Number(value || 0)
  return Number.isFinite(ratio) ? ratio.toFixed(2) : '-'
}

useSeoMeta({
  title: t('home.metaTitle'),
  description: t('home.metaDescription')
})
</script>
