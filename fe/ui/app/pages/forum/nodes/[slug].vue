<template>
  <div class="min-h-[calc(100vh-4rem)] bg-slate-50 py-6 dark:bg-slate-950">
    <div class="w-full px-3 sm:px-4 lg:px-5">
      <div class="mb-6 flex flex-col gap-4 lg:flex-row lg:items-start lg:justify-between">
        <div class="min-w-0">
          <h1 class="break-words text-2xl font-semibold text-slate-950 dark:text-white">
            {{ node ? nodeDisplayName(node) : $t('forum.fallback.node') }}
          </h1>
          <p v-if="node && nodeDisplayDesc(node)" class="mt-2 max-w-3xl break-words text-sm text-slate-500 dark:text-slate-400">
            {{ nodeDisplayDesc(node) }}
          </p>
        </div>

        <div class="flex flex-wrap items-center gap-2">
          <UButton color="neutral" variant="outline" icon="i-lucide-refresh-cw" :loading="pending" @click="loadTopics">
            {{ $t('common.refresh') }}
          </UButton>
          <UButton color="primary" icon="i-lucide-square-pen" :to="createTopicPath">
            {{ $t('forum.actions.createTopic') }}
          </UButton>
        </div>
      </div>

      <div class="grid gap-6 lg:grid-cols-[minmax(0,1fr)_300px] lg:items-start">
        <main>
          <ForumTopicList
            :topics="topics"
            :pending="pending"
            :error-message="errorMessage"
            :empty-title="$t('forum.topicList.empty.title')"
            :empty-description="$t('forum.topicList.empty.description')"
            :empty-action-label="$t('forum.actions.createTopic')"
            :empty-action-to="createTopicPath"
            @retry="loadTopics"
          />
        </main>

        <aside class="space-y-4">
          <section class="rounded-lg border border-slate-200 bg-white p-4 dark:border-slate-800 dark:bg-slate-900">
            <h2 class="text-base font-semibold text-slate-950 dark:text-white">{{ $t('forum.sidebar.title') }}</h2>
            <div class="mt-4 grid gap-2">
              <UButton color="primary" icon="i-lucide-square-pen" block :to="createTopicPath">
                {{ $t('forum.actions.createTopic') }}
              </UButton>
              <UButton color="neutral" variant="outline" icon="i-lucide-bookmark" block :to="localePath('/forum/bookmarks')">
                {{ $t('forum.sidebar.bookmarks') }}
              </UButton>
            </div>
          </section>

          <section v-if="node" class="rounded-lg border border-slate-200 bg-white p-4 dark:border-slate-800 dark:bg-slate-900">
            <h2 class="text-base font-semibold text-slate-950 dark:text-white">{{ $t('forum.topicList.nodeStats') }}</h2>
            <dl class="mt-4 space-y-3 text-sm">
              <div class="flex items-center justify-between gap-3">
                <dt class="text-slate-500 dark:text-slate-400">{{ $t('forum.fields.topics') }}</dt>
                <dd class="font-medium text-slate-950 dark:text-white">{{ numberFormatter.format(node.topicCount) }}</dd>
              </div>
              <div class="flex items-center justify-between gap-3">
                <dt class="text-slate-500 dark:text-slate-400">{{ $t('forum.fields.replies') }}</dt>
                <dd class="font-medium text-slate-950 dark:text-white">{{ numberFormatter.format(node.replyCount) }}</dd>
              </div>
            </dl>
          </section>
        </aside>
      </div>

      <AppPager
        class="mt-4"
        :page="page"
        :total="total"
        :page-size="Number(selectedSize)"
        :page-size-options="pageSizes"
        :disabled="pending"
        @page-change="goToPage"
        @page-size-change="handlePageSizeChange"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ApiError } from '~/composables/useApi'
import { useForum, type ForumNode, type ForumTopicListItem } from '~/composables/useForum'
import { localizeI18nName } from '~/utils/format'

definePageMeta({
  middleware: 'auth'
})

const { t, locale } = useI18n()
const localePath = useLocalePath()
const route = useRoute()
const router = useRouter()
const forum = useForum()

const slug = computed(() => String(route.params.slug || ''))
const node = ref<ForumNode | null>(null)
const topics = ref<ForumTopicListItem[]>([])
const total = ref(0)
const pending = ref(false)
const errorMessage = ref('')
const pageSizes = [20, 50, 100]

const page = ref(readPositiveIntQuery('page', 1))
const selectedSize = ref(String(readPageSizeQuery()))

const numberFormatter = computed(() => new Intl.NumberFormat(locale.value))
const totalPages = computed(() => Math.max(1, Math.ceil(total.value / Number(selectedSize.value || 20))))
const createTopicPath = computed(() => {
  if (!node.value?.id) return localePath('/forum/topics/create')
  return localePath(`/forum/topics/create?nodeId=${node.value.id}`)
})

useHead(() => ({
  title: `${node.value ? nodeDisplayName(node.value) : t('forum.metaTitle')} - NextPT`
}))

onMounted(loadTopics)

function readFirstQueryValue(key: string) {
  const value = route.query[key]
  return Array.isArray(value) ? value[0] : value
}

function readPositiveIntQuery(key: string, fallback: number) {
  const parsed = Number(readFirstQueryValue(key))
  return Number.isInteger(parsed) && parsed > 0 ? parsed : fallback
}

function readPageSizeQuery() {
  const parsed = readPositiveIntQuery('size', 20)
  return [20, 50, 100].includes(parsed) ? parsed : 20
}

async function loadTopics() {
  if (!slug.value) return

  pending.value = true
  errorMessage.value = ''
  syncQuery()

  try {
    const data = await forum.listTopics(slug.value, {
      page: page.value,
      size: Number(selectedSize.value)
    })
    node.value = data.node
    topics.value = data.list || []
    total.value = data.total || 0

    if (page.value > totalPages.value) {
      page.value = totalPages.value
      await loadTopics()
    }
  } catch (error) {
    topics.value = []
    total.value = 0
    errorMessage.value = error instanceof ApiError ? error.message : t('common.requestFailed')
  } finally {
    pending.value = false
  }
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
      ...route.query,
      page: page.value > 1 ? String(page.value) : undefined,
      size: selectedSize.value !== '20' ? selectedSize.value : undefined
    }
  })
}

function nodeDisplayName(item: ForumNode) {
  return localizeI18nName(item.nameI18n, locale.value, item.slug || `#${item.id}`)
}

function nodeDisplayDesc(item: ForumNode) {
  return localizeI18nName(item.descI18n, locale.value, '')
}
</script>
