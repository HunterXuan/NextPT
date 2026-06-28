<template>
  <div class="min-h-[calc(100vh-4rem)] bg-slate-50 py-6 dark:bg-slate-950">
    <div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
      <div class="mb-6 flex flex-col gap-4 lg:flex-row lg:items-end lg:justify-between">
        <div>
          <UButton color="neutral" variant="ghost" icon="i-lucide-arrow-left" :to="localePath('/forum')">
            {{ $t('forum.detail.back') }}
          </UButton>
          <p class="mt-4 text-sm font-medium text-slate-500 dark:text-slate-400">{{ $t('forum.eyebrow') }}</p>
          <h1 class="mt-1 text-2xl font-semibold text-slate-950 dark:text-white">{{ $t('forum.bookmarks.title') }}</h1>
        </div>

        <UButton color="neutral" variant="outline" icon="i-lucide-refresh-cw" :loading="pending" @click="loadBookmarks">
          {{ $t('common.refresh') }}
        </UButton>
      </div>

      <div class="overflow-hidden rounded-lg border border-slate-200 bg-white dark:border-slate-800 dark:bg-slate-900">
        <div v-if="pending" class="divide-y divide-slate-200 dark:divide-slate-800">
          <div v-for="index in 6" :key="index" class="grid gap-4 px-4 py-4 md:grid-cols-[minmax(0,1fr)_64px]">
            <div class="space-y-2">
              <div class="h-4 w-3/4 animate-pulse rounded bg-slate-200 dark:bg-slate-800" />
              <div class="h-3 w-1/2 animate-pulse rounded bg-slate-100 dark:bg-slate-800/70" />
            </div>
            <div class="h-7 w-12 animate-pulse rounded-full bg-slate-100 dark:bg-slate-800/70" />
          </div>
        </div>

        <div v-else-if="errorMessage" class="flex flex-col items-center justify-center px-4 py-16 text-center">
          <UIcon name="i-lucide-circle-alert" class="size-9 text-red-500" />
          <p class="mt-3 text-sm font-medium text-slate-950 dark:text-white">{{ errorMessage }}</p>
          <UButton class="mt-5" color="neutral" variant="outline" icon="i-lucide-refresh-cw" @click="loadBookmarks">
            {{ $t('common.retry') }}
          </UButton>
        </div>

        <div v-else-if="topics.length === 0" class="flex flex-col items-center justify-center px-4 py-16 text-center">
          <UIcon name="i-lucide-bookmark-x" class="size-9 text-slate-400" />
          <p class="mt-3 text-sm font-medium text-slate-950 dark:text-white">{{ $t('forum.bookmarks.empty.title') }}</p>
          <p class="mt-1 max-w-md text-sm text-slate-500 dark:text-slate-400">{{ $t('forum.bookmarks.empty.description') }}</p>
        </div>

        <div v-else class="divide-y divide-slate-200 dark:divide-slate-800">
          <article
            v-for="topic in topics"
            :key="topic.id"
            class="grid gap-3 px-4 py-3 transition-colors hover:bg-slate-50 md:grid-cols-[minmax(0,1fr)_64px] md:items-center dark:hover:bg-slate-950/70"
          >
            <div class="min-w-0">
              <div class="flex flex-wrap items-center gap-2">
                <UBadge v-if="topic.isSticky" color="primary" variant="soft">{{ $t('forum.topicList.badges.sticky') }}</UBadge>
                <UBadge v-if="topic.isLocked" color="neutral" variant="outline">{{ $t('forum.topicList.badges.locked') }}</UBadge>
              </div>
              <h2 class="mt-2 truncate text-sm font-semibold">
                <NuxtLink
                  :to="localePath(`/forum/topics/${topic.id}`)"
                  class="text-slate-950 hover:text-sky-700 dark:text-white dark:hover:text-sky-300"
                >
                  {{ topic.subject || `#${topic.id}` }}
                </NuxtLink>
              </h2>
              <p class="mt-2 text-xs text-slate-500 dark:text-slate-400">
                {{ topic.username || `#${topic.userId}` }}
                <span class="mx-1 text-slate-300 dark:text-slate-700">/</span>
                {{ $t('forum.topicList.meta.created') }} {{ formatDateTime(topic.createdAt, locale) }}
                <span class="mx-1 text-slate-300 dark:text-slate-700">/</span>
                {{ $t('forum.topicList.meta.lastReply') }} {{ formatDateTime(topic.lastReplyAt || topic.createdAt, locale) }}
              </p>
            </div>

            <NuxtLink
              :to="localePath(`/forum/topics/${topic.id}`)"
              class="inline-flex h-7 min-w-10 items-center justify-center rounded-full bg-slate-100 px-3 text-sm font-semibold text-slate-600 transition-colors hover:bg-sky-100 hover:text-sky-800 md:justify-self-end dark:bg-slate-800 dark:text-slate-300 dark:hover:bg-sky-950 dark:hover:text-sky-200"
            >
              {{ numberFormatter.format(topic.replyCount) }}
            </NuxtLink>
          </article>
        </div>
      </div>

      <div class="mt-4 flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
        <p class="text-sm text-slate-500 dark:text-slate-400">
          {{ $t('forum.pagination.summary', { page: page, pages: totalPages }) }}
        </p>
        <div class="flex flex-wrap items-center gap-2">
          <UButton color="neutral" variant="outline" icon="i-lucide-chevron-left" :disabled="page <= 1 || pending" @click="goToPage(page - 1)">
            {{ $t('common.previous') }}
          </UButton>
          <UButton color="neutral" variant="outline" trailing-icon="i-lucide-chevron-right" :disabled="page >= totalPages || pending" @click="goToPage(page + 1)">
            {{ $t('common.next') }}
          </UButton>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ApiError } from '~/composables/useApi'
import { useForum, type ForumTopicListItem } from '~/composables/useForum'
import { formatDateTime } from '~/utils/format'

definePageMeta({
  middleware: 'auth'
})

const { t, locale } = useI18n()
const localePath = useLocalePath()
const route = useRoute()
const router = useRouter()
const forum = useForum()

const topics = ref<ForumTopicListItem[]>([])
const total = ref(0)
const pending = ref(false)
const errorMessage = ref('')
const page = ref(readPositiveIntQuery('page', 1))
const size = 20

const numberFormatter = computed(() => new Intl.NumberFormat(locale.value))
const totalPages = computed(() => Math.max(1, Math.ceil(total.value / size)))

useHead(() => ({
  title: `${t('forum.bookmarks.metaTitle')} - NextPT`
}))

onMounted(loadBookmarks)

function readFirstQueryValue(key: string) {
  const value = route.query[key]
  return Array.isArray(value) ? value[0] : value
}

function readPositiveIntQuery(key: string, fallback: number) {
  const parsed = Number(readFirstQueryValue(key))
  return Number.isInteger(parsed) && parsed > 0 ? parsed : fallback
}

async function loadBookmarks() {
  pending.value = true
  errorMessage.value = ''
  syncQuery()

  try {
    const data = await forum.listBookmarks(page.value, size)
    topics.value = data.list || []
    total.value = data.total || 0

    if (page.value > totalPages.value) {
      page.value = totalPages.value
      await loadBookmarks()
    }
  } catch (error) {
    topics.value = []
    total.value = 0
    errorMessage.value = error instanceof ApiError ? error.message : t('common.requestFailed')
  } finally {
    pending.value = false
  }
}

function goToPage(nextPage: number) {
  page.value = Math.min(Math.max(1, nextPage), totalPages.value)
  loadBookmarks()
}

function syncQuery() {
  router.replace({
    query: {
      ...route.query,
      page: page.value > 1 ? String(page.value) : undefined
    }
  })
}
</script>
