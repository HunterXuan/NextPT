<template>
  <div class="min-h-[calc(100vh-4rem)] bg-slate-50 py-6 dark:bg-slate-950">
    <div class="w-full px-3 sm:px-4 lg:px-5">
      <div class="mb-4 grid grid-cols-2 gap-3 sm:flex sm:items-center sm:justify-end">
        <div class="rounded-lg border border-slate-200 bg-white px-4 py-3 dark:border-slate-800 dark:bg-slate-900">
          <p class="text-xs text-slate-500 dark:text-slate-400">{{ $t('forum.summary.nodes') }}</p>
          <p class="mt-1 text-lg font-semibold text-slate-950 dark:text-white">{{ numberFormatter.format(nodeCount) }}</p>
        </div>
        <div class="rounded-lg border border-slate-200 bg-white px-4 py-3 dark:border-slate-800 dark:bg-slate-900">
          <p class="text-xs text-slate-500 dark:text-slate-400">{{ $t('forum.summary.topics') }}</p>
          <p class="mt-1 text-lg font-semibold text-slate-950 dark:text-white">{{ numberFormatter.format(topicCount) }}</p>
        </div>
        <UButton class="col-span-2 sm:col-span-1" color="primary" icon="i-lucide-square-pen" :to="localePath('/forum/topics/create')">
          {{ $t('forum.actions.createTopic') }}
        </UButton>
      </div>

      <div v-if="pending" class="space-y-4">
        <div v-for="index in 3" :key="index" class="h-36 animate-pulse rounded-lg border border-slate-200 bg-white dark:border-slate-800 dark:bg-slate-900" />
      </div>

      <div v-else-if="errorMessage" class="flex flex-col items-center justify-center rounded-lg border border-slate-200 bg-white px-4 py-16 text-center dark:border-slate-800 dark:bg-slate-900">
        <UIcon name="i-lucide-circle-alert" class="size-9 text-red-500" />
        <p class="mt-3 text-sm font-medium text-slate-950 dark:text-white">{{ errorMessage }}</p>
        <UButton class="mt-5" color="neutral" variant="outline" icon="i-lucide-refresh-cw" @click="loadNodes">
          {{ $t('common.retry') }}
        </UButton>
      </div>

      <div v-else-if="categories.length === 0" class="flex flex-col items-center justify-center rounded-lg border border-slate-200 bg-white px-4 py-16 text-center dark:border-slate-800 dark:bg-slate-900">
        <UIcon name="i-lucide-inbox" class="size-9 text-slate-400" />
        <p class="mt-3 text-sm font-medium text-slate-950 dark:text-white">{{ $t('forum.empty.title') }}</p>
        <p class="mt-1 max-w-md text-sm text-slate-500 dark:text-slate-400">{{ $t('forum.empty.description') }}</p>
      </div>

      <div v-else class="grid gap-6 lg:grid-cols-[minmax(0,1fr)_300px] lg:items-start">
        <main class="space-y-5">
        <section
          v-for="category in categories"
          :key="category.id"
          class="overflow-hidden rounded-lg border border-slate-200 bg-white dark:border-slate-800 dark:bg-slate-900"
        >
          <div class="border-b border-slate-200 px-4 py-4 dark:border-slate-800">
            <h2 class="text-base font-semibold text-slate-950 dark:text-white">{{ categoryDisplayName(category) }}</h2>
            <p v-if="categoryDisplayDesc(category)" class="mt-1 text-sm text-slate-500 dark:text-slate-400">
              {{ categoryDisplayDesc(category) }}
            </p>
          </div>

          <div v-if="category.nodes.length === 0" class="px-4 py-8 text-center text-sm text-slate-500 dark:text-slate-400">
            {{ $t('forum.empty.nodes') }}
          </div>

          <div v-else class="grid gap-3 px-4 py-4 sm:grid-cols-2 xl:grid-cols-3">
            <NuxtLink
              v-for="node in category.nodes"
              :key="node.id"
              :to="localePath(`/forum/nodes/${node.slug}`)"
              class="rounded-md border border-slate-200 px-3 py-3 transition-colors hover:border-sky-200 hover:bg-sky-50/60 dark:border-slate-800 dark:hover:border-sky-800 dark:hover:bg-sky-950/30"
            >
              <div class="min-w-0">
                <div class="flex items-center gap-2">
                  <UIcon name="i-lucide-messages-square" class="size-4 shrink-0 text-sky-500" />
                  <h3 class="truncate text-sm font-semibold text-slate-950 dark:text-white">{{ nodeDisplayName(node) }}</h3>
                </div>
                <p v-if="nodeDisplayDesc(node)" class="mt-1 line-clamp-2 text-sm text-slate-500 dark:text-slate-400">
                  {{ nodeDisplayDesc(node) }}
                </p>
              </div>

              <div class="mt-3 flex items-center gap-3 text-xs text-slate-500 dark:text-slate-400">
                <span>{{ $t('forum.nodeStats.topics', { count: numberFormatter.format(node.topicCount) }) }}</span>
                <span>{{ $t('forum.nodeStats.replies', { count: numberFormatter.format(node.replyCount) }) }}</span>
              </div>
            </NuxtLink>
          </div>
        </section>
        </main>

        <aside class="space-y-4">
          <section class="rounded-lg border border-slate-200 bg-white p-4 dark:border-slate-800 dark:bg-slate-900">
            <h2 class="text-base font-semibold text-slate-950 dark:text-white">{{ $t('forum.sidebar.title') }}</h2>
            <div class="mt-4 grid gap-2">
              <UButton color="primary" icon="i-lucide-square-pen" block :to="localePath('/forum/topics/create')">
                {{ $t('forum.actions.createTopic') }}
              </UButton>
              <UButton color="neutral" variant="outline" icon="i-lucide-bookmark" block :to="localePath('/forum/bookmarks')">
                {{ $t('forum.sidebar.bookmarks') }}
              </UButton>
            </div>
          </section>

          <section class="rounded-lg border border-slate-200 bg-white p-4 dark:border-slate-800 dark:bg-slate-900">
            <h2 class="text-base font-semibold text-slate-950 dark:text-white">{{ $t('forum.sidebar.quickNodes') }}</h2>
            <div class="mt-3 flex flex-wrap gap-2">
              <NuxtLink
                v-for="node in quickNodes"
                :key="node.id"
                :to="localePath(`/forum/nodes/${node.slug}`)"
                class="rounded-md bg-slate-100 px-2.5 py-1 text-xs font-medium text-slate-700 transition-colors hover:bg-sky-100 hover:text-sky-800 dark:bg-slate-800 dark:text-slate-200 dark:hover:bg-sky-950 dark:hover:text-sky-200"
              >
                {{ nodeDisplayName(node) }}
              </NuxtLink>
            </div>
          </section>
        </aside>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ApiError } from '~/composables/useApi'
import { useForum, type ForumNode, type ForumNodeCategory } from '~/composables/useForum'
import { localizeI18nName } from '~/utils/format'

definePageMeta({
  middleware: 'auth'
})

const { t, locale } = useI18n()
const localePath = useLocalePath()
const forum = useForum()

const categories = ref<ForumNodeCategory[]>([])
const pending = ref(true)
const errorMessage = ref('')
const numberFormatter = computed(() => new Intl.NumberFormat(locale.value))

const nodeCount = computed(() => categories.value.reduce((sum, category) => sum + category.nodes.length, 0))
const topicCount = computed(() => categories.value.reduce((sum, category) => {
  return sum + category.nodes.reduce((nodeSum, node) => nodeSum + Number(node.topicCount || 0), 0)
}, 0))
const quickNodes = computed(() => {
  return categories.value
    .flatMap((category) => category.nodes)
    .sort((a, b) => Number(b.topicCount || 0) - Number(a.topicCount || 0))
    .slice(0, 12)
})

useHead(() => ({
  title: `${t('forum.metaTitle')} - NextPT`
}))

onMounted(loadNodes)

async function loadNodes() {
  pending.value = true
  errorMessage.value = ''
  try {
    const data = await forum.listNodes()
    categories.value = data.list || []
  } catch (error) {
    errorMessage.value = error instanceof ApiError ? error.message : t('common.requestFailed')
  } finally {
    pending.value = false
  }
}

function categoryDisplayName(category: ForumNodeCategory) {
  return localizeI18nName(category.nameI18n, locale.value, t('forum.fallback.category'))
}

function categoryDisplayDesc(category: ForumNodeCategory) {
  return localizeI18nName(category.descI18n, locale.value, '')
}

function nodeDisplayName(node: ForumNode) {
  return localizeI18nName(node.nameI18n, locale.value, t('forum.fallback.node'))
}

function nodeDisplayDesc(node: ForumNode) {
  return localizeI18nName(node.descI18n, locale.value, '')
}
</script>
