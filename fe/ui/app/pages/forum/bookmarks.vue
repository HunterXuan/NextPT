<template>
  <div class="min-h-[calc(100vh-4rem)] bg-slate-50 py-6 dark:bg-slate-950">
    <div class="w-full px-3 sm:px-4 lg:px-5">
      <div class="space-y-3">
        <div class="grid grid-cols-2 gap-2 sm:flex sm:justify-end">
          <UButton class="w-full justify-center sm:w-auto" color="neutral" variant="outline" size="sm" icon="i-lucide-message-square" :to="localePath('/forum')">
            {{ $t('forum.bookmarks.actions.browse') }}
          </UButton>
          <UButton class="w-full justify-center sm:w-auto" color="primary" size="sm" icon="i-lucide-square-pen" :to="localePath('/forum/topics/create')">
            {{ $t('forum.actions.createTopic') }}
          </UButton>
        </div>

        <ForumTopicList
          :topics="topics"
          :pending="pending"
          :error-message="errorMessage"
          :empty-title="$t('forum.bookmarks.empty.title')"
          :empty-description="$t('forum.bookmarks.empty.description')"
          empty-icon="i-lucide-bookmark-x"
          empty-action-icon="i-lucide-search"
          :empty-action-label="$t('forum.bookmarks.empty.action')"
          :empty-action-to="localePath('/forum')"
          @retry="loadBookmarks"
        >
          <template #topic-actions="{ topic }">
            <UTooltip
              :text="$t('forum.bookmarks.actions.remove')"
              :content="{ side: 'top', sideOffset: 8 }"
              :delay-duration="120"
            >
              <UButton
                class="h-7 min-w-10 justify-center rounded-full px-0"
                color="warning"
                variant="soft"
                size="sm"
                icon="i-lucide-bookmark-x"
                :aria-label="$t('forum.bookmarks.actions.remove')"
                :loading="removingBookmarkId === topic.id"
                :disabled="removingBookmarkId > 0 || pending"
                @click="handleRemoveBookmark(topic)"
              />
            </UTooltip>
          </template>
        </ForumTopicList>

        <AppPager
          :page="page"
          :total="total"
          :page-size="Number(selectedSize)"
          :page-size-options="pageSizes"
          :disabled="pending || removingBookmarkId > 0"
          @page-change="goToPage"
          @page-size-change="handlePageSizeChange"
        />
      </div>
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
const toast = useToast()
const forum = useForum()

const topics = ref<ForumTopicListItem[]>([])
const total = ref(0)
const pending = ref(true)
const errorMessage = ref('')
const removingBookmarkId = ref(0)
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

async function handleRemoveBookmark(topic: ForumTopicListItem) {
  if (removingBookmarkId.value > 0 || pending.value) return

  removingBookmarkId.value = topic.id
  try {
    await forum.unbookmarkTopic(topic.id)
    toast.add({
      title: t('forum.bookmarks.actions.removed'),
      color: 'success',
      icon: 'i-lucide-bookmark-x'
    })

    if (topics.value.length === 1 && page.value > 1) {
      page.value -= 1
    }
    await loadBookmarks()
  } catch (error) {
    toast.add({
      title: error instanceof ApiError ? error.message : t('common.requestFailed'),
      color: 'error',
      icon: 'i-lucide-circle-alert'
    })
  } finally {
    removingBookmarkId.value = 0
  }
}
</script>
