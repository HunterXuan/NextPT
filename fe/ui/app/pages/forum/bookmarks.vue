<template>
  <div class="min-h-[calc(100vh-4rem)] bg-slate-50 py-6 dark:bg-slate-950">
    <div class="w-full px-3 sm:px-4 lg:px-5">
      <ForumTopicList
        :topics="topics"
        :pending="pending"
        :error-message="errorMessage"
        :empty-title="$t('forum.bookmarks.empty.title')"
        :empty-description="$t('forum.bookmarks.empty.description')"
        empty-icon="i-lucide-bookmark-x"
        :empty-action-label="$t('forum.bookmarks.empty.action')"
        :empty-action-to="localePath('/forum')"
        @retry="loadBookmarks"
      />

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
import { useForum, type ForumTopicListItem } from '~/composables/useForum'

definePageMeta({
  middleware: 'auth'
})

const { t } = useI18n()
const localePath = useLocalePath()
const route = useRoute()
const router = useRouter()
const forum = useForum()

const topics = ref<ForumTopicListItem[]>([])
const total = ref(0)
const pending = ref(false)
const errorMessage = ref('')
const pageSizes = [20, 50, 100]

const page = ref(readPositiveIntQuery('page', 1))
const selectedSize = ref(String(readPageSizeQuery()))

const totalPages = computed(() => Math.max(1, Math.ceil(total.value / Number(selectedSize.value || 20))))

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

function readPageSizeQuery() {
  const parsed = readPositiveIntQuery('size', 20)
  return [20, 50, 100].includes(parsed) ? parsed : 20
}

async function loadBookmarks() {
  pending.value = true
  errorMessage.value = ''
  syncQuery()

  try {
    const data = await forum.listBookmarks(page.value, Number(selectedSize.value))
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

function handlePageSizeChange(nextSize: number) {
  selectedSize.value = String(nextSize)
  page.value = 1
  loadBookmarks()
}

function goToPage(nextPage: number) {
  page.value = Math.min(Math.max(1, nextPage), totalPages.value)
  loadBookmarks()
}

function syncQuery() {
  router.replace({
    query: {
      page: page.value > 1 ? String(page.value) : undefined,
      size: selectedSize.value !== '20' ? selectedSize.value : undefined
    }
  })
}
</script>
