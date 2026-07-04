<template>
  <div class="min-h-[calc(100vh-4rem)] bg-slate-50 py-6 dark:bg-slate-950">
    <div class="w-full px-3 sm:px-4 lg:px-5">
      <div v-if="nodesPending" class="space-y-3">
        <div class="h-24 animate-pulse rounded-lg border border-slate-200 bg-white dark:border-slate-800 dark:bg-slate-900" />
        <div class="h-96 animate-pulse rounded-lg border border-slate-200 bg-white dark:border-slate-800 dark:bg-slate-900" />
      </div>

      <div v-else-if="nodesError" class="flex flex-col items-center justify-center rounded-lg border border-slate-200 bg-white px-4 py-16 text-center dark:border-slate-800 dark:bg-slate-900">
        <UIcon name="i-lucide-circle-alert" class="size-9 text-red-500" />
        <p class="mt-3 text-sm font-medium text-slate-950 dark:text-white">{{ nodesError }}</p>
        <UButton class="mt-5" color="neutral" variant="outline" icon="i-lucide-refresh-cw" @click="loadForumHome">
          {{ $t('common.retry') }}
        </UButton>
      </div>

      <div v-else-if="categories.length === 0" class="flex flex-col items-center justify-center rounded-lg border border-slate-200 bg-white px-4 py-16 text-center dark:border-slate-800 dark:bg-slate-900">
        <UIcon name="i-lucide-inbox" class="size-9 text-slate-400" />
        <p class="mt-3 text-sm font-medium text-slate-950 dark:text-white">{{ $t('forum.empty.title') }}</p>
        <p class="mt-1 max-w-md text-sm text-slate-500 dark:text-slate-400">{{ $t('forum.empty.description') }}</p>
      </div>

      <div v-else class="grid gap-4 xl:grid-cols-[minmax(0,1fr)_300px] xl:items-start">
        <main class="min-w-0 space-y-3">
          <section class="rounded-lg border border-slate-200 bg-white px-3 py-2.5 dark:border-slate-800 dark:bg-slate-900">
            <div class="-mx-1 flex min-w-0 items-center gap-1 overflow-x-auto px-1">
              <button
                v-for="category in categories"
                :key="category.id"
                type="button"
                class="inline-flex h-8 shrink-0 items-center rounded-md px-2.5 text-sm font-medium transition-colors"
                :class="selectedCategoryId === category.id
                  ? 'bg-slate-100 text-slate-950 ring-1 ring-slate-200 dark:bg-slate-800 dark:text-white dark:ring-slate-700'
                  : 'text-slate-500 hover:bg-slate-50 hover:text-slate-950 dark:text-slate-400 dark:hover:bg-slate-800/70 dark:hover:text-white'"
                @click="selectCategory(category)"
              >
                {{ categoryDisplayName(category) }}
              </button>
            </div>

            <div v-if="selectedCategoryNodes.length > 0" class="mt-2 flex items-center gap-2 border-t border-slate-100 pt-2 dark:border-slate-800">
              <div class="-mx-1 flex min-w-0 flex-1 items-center gap-1 overflow-x-auto px-1">
                <button
                  v-for="node in selectedCategoryNodes"
                  :key="node.id"
                  type="button"
                  class="inline-flex h-7 shrink-0 items-center rounded-md px-2.5 text-sm transition-colors"
                  :class="selectedNodeSlug === node.slug
                    ? 'bg-sky-50 font-medium text-sky-700 ring-1 ring-sky-100 dark:bg-sky-950/50 dark:text-sky-200 dark:ring-sky-900'
                    : 'text-slate-500 hover:bg-slate-50 hover:text-slate-950 dark:text-slate-400 dark:hover:bg-slate-800/70 dark:hover:text-white'"
                  @click="selectNode(node)"
                >
                  {{ nodeDisplayName(node) }}
                </button>
              </div>
            </div>
            <div v-else class="mt-2 flex items-center gap-2 border-t border-slate-100 pt-3 text-sm text-slate-500 dark:border-slate-800 dark:text-slate-400">
              <UIcon name="i-lucide-circle-slash" class="size-4" />
              <span>{{ $t('forum.empty.nodes') }}</span>
            </div>
          </section>

          <ForumTopicList
            v-if="selectedNodeSlug"
            :topics="topics"
            :pending="topicsPending"
            :error-message="topicsError"
            :empty-title="$t('forum.topicList.empty.title')"
            :empty-description="$t('forum.topicList.empty.description')"
            :empty-action-label="$t('forum.actions.createTopic')"
            :empty-action-to="createTopicPath"
            @retry="loadTopics"
          />

          <AppPager
            v-if="selectedNodeSlug"
            class="mt-4"
            :page="page"
            :total="total"
            :page-size="Number(selectedSize)"
            :page-size-options="pageSizes"
            :disabled="topicsPending"
            @page-change="goToPage"
            @page-size-change="handlePageSizeChange"
          />
        </main>

        <aside class="app-sticky-offset space-y-3 xl:sticky">
          <section class="rounded-lg border border-slate-200 bg-white p-3 dark:border-slate-800 dark:bg-slate-900">
            <UButton color="primary" icon="i-lucide-square-pen" block :to="createTopicPath">
              {{ $t('forum.actions.createTopic') }}
            </UButton>
          </section>

          <section v-if="selectedNode" class="rounded-lg border border-slate-200 bg-white p-3 dark:border-slate-800 dark:bg-slate-900">
            <h2 class="text-sm font-semibold text-slate-950 dark:text-white">{{ $t('forum.sidebar.currentNode') }}</h2>
            <p class="mt-1 truncate text-sm text-slate-500 dark:text-slate-400" :title="nodeDisplayName(selectedNode)">
              {{ nodeDisplayName(selectedNode) }}
            </p>
            <dl class="mt-3 grid grid-cols-2 gap-2">
              <div class="rounded-md bg-slate-50 px-3 py-2 dark:bg-slate-950/60">
                <dt class="text-xs text-slate-500 dark:text-slate-400">{{ $t('forum.fields.topics') }}</dt>
                <dd class="mt-1 text-base font-semibold text-slate-950 dark:text-white">{{ numberFormatter.format(selectedNode.topicCount) }}</dd>
              </div>
              <div class="rounded-md bg-slate-50 px-3 py-2 dark:bg-slate-950/60">
                <dt class="text-xs text-slate-500 dark:text-slate-400">{{ $t('forum.fields.replies') }}</dt>
                <dd class="mt-1 text-base font-semibold text-slate-950 dark:text-white">{{ numberFormatter.format(selectedNode.replyCount) }}</dd>
              </div>
            </dl>
          </section>

          <section class="rounded-lg border border-slate-200 bg-white p-3 dark:border-slate-800 dark:bg-slate-900">
            <h2 class="text-sm font-semibold text-slate-950 dark:text-white">{{ $t('forum.sidebar.siteStats') }}</h2>
            <dl class="mt-3 space-y-2 text-sm">
              <div class="flex items-center justify-between gap-3">
                <dt class="text-slate-500 dark:text-slate-400">{{ $t('forum.summary.nodes') }}</dt>
                <dd class="font-medium text-slate-950 dark:text-white">{{ numberFormatter.format(forumStats.nodes) }}</dd>
              </div>
              <div class="flex items-center justify-between gap-3">
                <dt class="text-slate-500 dark:text-slate-400">{{ $t('forum.summary.topics') }}</dt>
                <dd class="font-medium text-slate-950 dark:text-white">{{ numberFormatter.format(forumStats.topics) }}</dd>
              </div>
              <div class="flex items-center justify-between gap-3">
                <dt class="text-slate-500 dark:text-slate-400">{{ $t('forum.fields.replies') }}</dt>
                <dd class="font-medium text-slate-950 dark:text-white">{{ numberFormatter.format(forumStats.replies) }}</dd>
              </div>
            </dl>
          </section>

          <section v-if="quickNodes.length > 0" class="rounded-lg border border-slate-200 bg-white p-3 dark:border-slate-800 dark:bg-slate-900">
            <h2 class="text-sm font-semibold text-slate-950 dark:text-white">{{ $t('forum.sidebar.quickNodes') }}</h2>
            <div class="mt-3 flex flex-wrap gap-1.5">
              <button
                v-for="node in quickNodes"
                :key="node.id"
                type="button"
                class="rounded-md bg-slate-50 px-2 py-1 text-xs font-medium text-slate-600 transition-colors hover:bg-sky-50 hover:text-sky-700 dark:bg-slate-950/70 dark:text-slate-300 dark:hover:bg-sky-950/50 dark:hover:text-sky-200"
                :class="selectedNodeSlug === node.slug ? 'bg-sky-50 text-sky-700 dark:bg-sky-950/50 dark:text-sky-200' : ''"
                @click="selectNode(node)"
              >
                {{ nodeDisplayName(node) }}
              </button>
            </div>
          </section>
        </aside>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ApiError } from '~/composables/useApi'
import { useForum, type ForumNode, type ForumNodeCategory, type ForumTopicListItem } from '~/composables/useForum'
import { localizeI18nName } from '~/utils/format'

definePageMeta({
  middleware: 'auth'
})

interface ForumHomeSelection {
  categoryId: number
  nodeSlug: string
}

interface ForumNodeWithCategory {
  category: ForumNodeCategory
  node: ForumNode
}

const forumHomeSelectionStorageKey = 'nextpt:forum:home-selection'

const { t, locale } = useI18n()
const localePath = useLocalePath()
const route = useRoute()
const router = useRouter()
const forum = useForum()

const categories = ref<ForumNodeCategory[]>([])
const topics = ref<ForumTopicListItem[]>([])
const total = ref(0)
const nodesPending = ref(true)
const topicsPending = ref(false)
const nodesError = ref('')
const topicsError = ref('')
const selectedCategoryId = ref(0)
const selectedNodeSlug = ref('')
const pageSizes = [20, 50, 100]

const page = ref(readPositiveIntQuery('page', 1))
const selectedSize = ref(String(readPageSizeQuery()))

const selectedCategory = computed(() => {
  return categories.value.find((category) => category.id === selectedCategoryId.value) || categories.value[0] || null
})
const selectedCategoryNodes = computed(() => selectedCategory.value?.nodes || [])
const selectedNode = computed(() => {
  return selectedCategoryNodes.value.find((node) => node.slug === selectedNodeSlug.value) || selectedCategoryNodes.value[0] || null
})
const totalPages = computed(() => Math.max(1, Math.ceil(total.value / Number(selectedSize.value || 20))))
const numberFormatter = computed(() => new Intl.NumberFormat(locale.value))
const flatNodeItems = computed<ForumNodeWithCategory[]>(() => {
  return categories.value.flatMap((category) => category.nodes.map((node) => ({ category, node })))
})
const forumStats = computed(() => {
  return flatNodeItems.value.reduce((stats, item) => {
    stats.nodes += 1
    stats.topics += Number(item.node.topicCount || 0)
    stats.replies += Number(item.node.replyCount || 0)
    return stats
  }, { nodes: 0, topics: 0, replies: 0 })
})
const quickNodes = computed(() => {
  return [...selectedCategoryNodes.value]
    .sort((left, right) => {
      const rightScore = Number(right.topicCount || 0) * 3 + Number(right.replyCount || 0)
      const leftScore = Number(left.topicCount || 0) * 3 + Number(left.replyCount || 0)
      return rightScore - leftScore
    })
    .slice(0, 8)
})
const createTopicPath = computed(() => {
  if (!selectedNode.value?.id) return localePath('/forum/topics/create')
  return localePath(`/forum/topics/create?nodeId=${selectedNode.value.id}`)
})

useHead(() => ({
  title: t('forum.metaTitle')
}))

onMounted(loadForumHome)

function readFirstQueryValue(key: string) {
  const value = route.query[key]
  return Array.isArray(value) ? value[0] : value
}

function readStringQuery(key: string) {
  return String(readFirstQueryValue(key) || '').trim()
}

function readPositiveIntQuery(key: string, fallback: number) {
  const parsed = Number(readFirstQueryValue(key))
  return Number.isInteger(parsed) && parsed > 0 ? parsed : fallback
}

function readPageSizeQuery() {
  const parsed = readPositiveIntQuery('size', 20)
  return [20, 50, 100].includes(parsed) ? parsed : 20
}

async function loadForumHome() {
  nodesPending.value = true
  nodesError.value = ''

  try {
    const data = await forum.listNodes()
    categories.value = data.list || []
    resolveSelectionFromQuery()
  } catch (error) {
    categories.value = []
    topics.value = []
    total.value = 0
    selectedCategoryId.value = 0
    selectedNodeSlug.value = ''
    nodesError.value = error instanceof ApiError ? error.message : t('common.requestFailed')
  } finally {
    nodesPending.value = false
  }

  if (selectedNodeSlug.value) {
    await loadTopics()
  }
}

function resolveSelectionFromQuery() {
  const queryNodeSlug = readStringQuery('node')
  const queryCategoryId = readPositiveIntQuery('category', 0)
  const hasQuerySelection = Boolean(queryNodeSlug || queryCategoryId)
  const rememberedSelection = hasQuerySelection ? null : readRememberedSelection()
  const flatNodes = categories.value.flatMap((category) => category.nodes.map((node) => ({ category, node })))
  const matchedNode = flatNodes.find((item) => item.node.slug === queryNodeSlug)

  if (matchedNode) {
    selectedCategoryId.value = matchedNode.category.id
    selectedNodeSlug.value = matchedNode.node.slug
    return
  }

  const matchedCategory = categories.value.find((category) => category.id === queryCategoryId && category.nodes.length > 0)
  const rememberedNode = rememberedSelection?.nodeSlug
    ? flatNodes.find((item) => item.node.slug === rememberedSelection.nodeSlug)
    : null
  if (rememberedNode) {
    selectedCategoryId.value = rememberedNode.category.id
    selectedNodeSlug.value = rememberedNode.node.slug
    return
  }

  const rememberedCategory = rememberedSelection?.categoryId
    ? categories.value.find((category) => category.id === rememberedSelection.categoryId && category.nodes.length > 0)
    : null
  const firstCategoryWithNodes = categories.value.find((category) => category.nodes.length > 0)
  const nextCategory = matchedCategory || rememberedCategory || firstCategoryWithNodes || categories.value[0]
  selectedCategoryId.value = nextCategory?.id || 0
  selectedNodeSlug.value = nextCategory?.nodes[0]?.slug || ''
}

function readRememberedSelection(): ForumHomeSelection | null {
  if (!import.meta.client) return null

  try {
    const rawValue = localStorage.getItem(forumHomeSelectionStorageKey)
    if (!rawValue) return null

    const parsed = JSON.parse(rawValue) as Partial<ForumHomeSelection>
    const categoryId = Number(parsed.categoryId || 0)
    const nodeSlug = typeof parsed.nodeSlug === 'string' ? parsed.nodeSlug : ''
    if (!categoryId && !nodeSlug) return null

    return { categoryId, nodeSlug }
  } catch {
    return null
  }
}

function rememberSelection() {
  if (!import.meta.client || !selectedCategoryId.value || !selectedNodeSlug.value) return

  try {
    localStorage.setItem(forumHomeSelectionStorageKey, JSON.stringify({
      categoryId: selectedCategoryId.value,
      nodeSlug: selectedNodeSlug.value
    }))
  } catch {
    // Ignore storage failures, browsing should continue normally.
  }
}

async function loadTopics() {
  if (!selectedNodeSlug.value) {
    topics.value = []
    total.value = 0
    return
  }

  topicsPending.value = true
  topicsError.value = ''
  rememberSelection()
  syncQuery()

  try {
    const data = await forum.listTopics(selectedNodeSlug.value, {
      page: page.value,
      size: Number(selectedSize.value)
    })
    topics.value = data.list || []
    total.value = data.total || 0

    if (page.value > totalPages.value) {
      page.value = totalPages.value
      await loadTopics()
    }
  } catch (error) {
    topics.value = []
    total.value = 0
    topicsError.value = error instanceof ApiError ? error.message : t('common.requestFailed')
  } finally {
    topicsPending.value = false
  }
}

function selectCategory(category: ForumNodeCategory) {
  if (selectedCategoryId.value === category.id) return

  selectedCategoryId.value = category.id
  selectedNodeSlug.value = category.nodes[0]?.slug || ''
  page.value = 1
  loadTopics()
}

function selectNode(node: ForumNode) {
  if (selectedNodeSlug.value === node.slug) return

  selectedNodeSlug.value = node.slug
  page.value = 1
  loadTopics()
}

function handlePageSizeChange(nextSize: number) {
  selectedSize.value = String(nextSize)
  page.value = 1
  loadTopics()
}

function goToPage(nextPage: number) {
  page.value = Math.min(Math.max(1, nextPage), totalPages.value)
  loadTopics()
}

function syncQuery() {
  router.replace({
    query: {
      category: selectedCategoryId.value > 0 ? String(selectedCategoryId.value) : undefined,
      node: selectedNodeSlug.value || undefined,
      page: page.value > 1 ? String(page.value) : undefined,
      size: selectedSize.value !== '20' ? selectedSize.value : undefined
    }
  })
}

function categoryDisplayName(category: ForumNodeCategory) {
  return localizeI18nName(category.nameI18n, locale.value, t('forum.fallback.category'))
}

function nodeDisplayName(node: ForumNode) {
  return localizeI18nName(node.nameI18n, locale.value, t('forum.fallback.node'))
}
</script>
