<template>
  <div class="min-h-[calc(100vh-4rem)] bg-slate-50 py-6 dark:bg-slate-950">
    <div class="w-full px-3 sm:px-4 lg:px-5">
      <section class="overflow-hidden rounded-lg border border-slate-200 bg-white dark:border-slate-800 dark:bg-slate-900">
        <div class="flex flex-col gap-3 border-b border-slate-200 px-4 py-3 lg:flex-row lg:items-center lg:justify-between dark:border-slate-800">
          <div class="min-w-0">
            <h1 class="text-base font-semibold text-slate-950 dark:text-white">{{ $t('site.messages.title') }}</h1>
          </div>
          <div class="flex flex-wrap items-center gap-2">
            <div class="inline-flex h-9 items-center rounded-md border border-slate-200 bg-slate-50 p-0.5 dark:border-slate-800 dark:bg-slate-950/70">
              <button
                v-for="option in filterOptions"
                :key="option.value"
                type="button"
                class="h-8 rounded px-3 text-sm font-medium transition"
                :class="filter === option.value ? 'bg-white text-slate-950 shadow-sm dark:bg-slate-800 dark:text-white' : 'text-slate-500 hover:text-slate-950 dark:text-slate-400 dark:hover:text-white'"
                @click="setFilter(option.value)"
              >
                {{ option.label }}
              </button>
            </div>
            <UButton class="h-9" color="neutral" variant="outline" size="sm" icon="i-lucide-check-check" :disabled="pending || messages.length === 0" @click="markAllRead">
              {{ $t('site.messages.markAllRead') }}
            </UButton>
          </div>
        </div>

        <div v-if="pending" class="grid min-h-[520px] gap-0 lg:grid-cols-[380px_minmax(0,1fr)]">
          <div class="space-y-2 border-r border-slate-200 p-3 dark:border-slate-800">
            <div v-for="item in 8" :key="item" class="h-20 animate-pulse rounded-md bg-slate-100 dark:bg-slate-800" />
          </div>
          <div class="p-5">
            <div class="h-5 w-1/3 animate-pulse rounded bg-slate-100 dark:bg-slate-800" />
            <div class="mt-5 space-y-2">
              <div class="h-4 w-full animate-pulse rounded bg-slate-100 dark:bg-slate-800" />
              <div class="h-4 w-4/5 animate-pulse rounded bg-slate-100 dark:bg-slate-800" />
              <div class="h-4 w-2/3 animate-pulse rounded bg-slate-100 dark:bg-slate-800" />
            </div>
          </div>
        </div>

        <div v-else-if="errorMessage" class="flex flex-col items-center justify-center px-4 py-16 text-center">
          <UIcon name="i-lucide-circle-alert" class="size-9 text-red-500" />
          <p class="mt-3 text-sm font-medium text-slate-950 dark:text-white">{{ errorMessage }}</p>
        </div>

        <div v-else-if="messages.length === 0" class="flex flex-col items-center justify-center px-4 py-16 text-center">
          <UIcon name="i-lucide-bell-off" class="size-9 text-slate-400" />
          <p class="mt-3 text-sm font-medium text-slate-950 dark:text-white">{{ $t('site.messages.empty') }}</p>
        </div>

        <div v-else class="grid min-h-[520px] gap-0 lg:grid-cols-[380px_minmax(0,1fr)]">
          <div class="border-r border-slate-200 dark:border-slate-800">
            <div class="divide-y divide-slate-100 dark:divide-slate-800">
              <button
                v-for="item in messages"
                :key="item.id"
                type="button"
                class="flex w-full gap-3 px-4 py-3 text-left transition"
                :class="selectedMessageId === item.id ? 'bg-sky-50/80 dark:bg-sky-950/30' : 'hover:bg-slate-50 dark:hover:bg-slate-950/60'"
                @click="selectMessage(item)"
              >
                <span class="relative mt-1 flex size-9 shrink-0 items-center justify-center rounded-md" :class="item.isRead ? 'bg-slate-100 text-slate-400 dark:bg-slate-800 dark:text-slate-500' : 'bg-sky-50 text-sky-600 dark:bg-sky-950/40 dark:text-sky-300'">
                  <UIcon :name="item.isRead ? 'i-lucide-mail-open' : 'i-lucide-mail'" class="size-4" />
                  <span v-if="!item.isRead" class="absolute -right-0.5 -top-0.5 size-2.5 rounded-full bg-sky-500 ring-2 ring-white dark:ring-slate-900" />
                </span>
                <span class="min-w-0 flex-1">
                  <span class="flex min-w-0 items-center justify-between gap-3">
                    <UTooltip class="min-w-0" :text="item.title" :content="{ side: 'top', sideOffset: 8 }" :delay-duration="600">
                      <span class="block truncate text-sm text-slate-950 dark:text-white" :class="item.isRead ? 'font-medium' : 'font-semibold'">{{ item.title }}</span>
                    </UTooltip>
                    <span class="shrink-0 text-xs text-slate-500 dark:text-slate-400">{{ formatRelativeDateTime(item.createdAt, locale) }}</span>
                  </span>
                  <span class="mt-1 block truncate text-sm leading-5 text-slate-500 dark:text-slate-400">{{ messageExcerpt(item) }}</span>
                </span>
              </button>
            </div>
          </div>

          <article v-if="selectedMessage" class="min-w-0 p-5">
            <div class="flex flex-col gap-3 border-b border-slate-100 pb-4 sm:flex-row sm:items-start sm:justify-between dark:border-slate-800">
              <div class="min-w-0">
                <div class="flex items-center gap-2">
                  <span class="flex size-9 shrink-0 items-center justify-center rounded-md bg-slate-100 text-slate-500 dark:bg-slate-800 dark:text-slate-300">
                    <UIcon name="i-lucide-mail-open" class="size-4" />
                  </span>
                  <div class="min-w-0">
                    <h2 class="break-words text-base font-semibold text-slate-950 dark:text-white">{{ selectedMessage.title }}</h2>
                    <div class="mt-1 flex flex-wrap items-center gap-1.5 text-xs text-slate-500 dark:text-slate-400">
                      <IamUserPopover
                        v-if="selectedMessage.senderId"
                        :id="selectedMessage.senderId"
                        :user="selectedMessage.sender"
                        :fallback="messageSender(selectedMessage)"
                      />
                      <span v-else>{{ messageSender(selectedMessage) }}</span>
                      <span>/</span>
                      <span>{{ formatDateTime(selectedMessage.createdAt, locale) }}</span>
                    </div>
                  </div>
                </div>
              </div>
              <UTooltip v-if="messageTargetPath(selectedMessage)" :text="$t('site.messages.openTarget')" :content="{ side: 'top', sideOffset: 8 }" :delay-duration="300">
                <UButton color="neutral" variant="outline" size="sm" icon="i-lucide-arrow-up-right" :to="localePath(messageTargetPath(selectedMessage) || '/')" :aria-label="$t('site.messages.openTarget')" />
              </UTooltip>
            </div>
            <div class="mt-5">
              <div v-if="selectedMessageContent" class="rich-text" v-html="selectedMessageContent" />
              <p v-else class="text-sm text-slate-500 dark:text-slate-400">{{ $t('site.messages.noContent') }}</p>
            </div>
          </article>

          <div v-else class="flex items-center justify-center p-8 text-center">
            <div>
              <UIcon name="i-lucide-mail-open" class="mx-auto size-9 text-slate-400" />
              <p class="mt-3 text-sm text-slate-500 dark:text-slate-400">{{ $t('site.messages.selectHint') }}</p>
            </div>
          </div>
        </div>

        <AppPager
          class="border-t border-slate-200 px-4 py-3 dark:border-slate-800"
          :page="page"
          :total="total"
          :page-size="size"
          :page-size-options="pageSizes"
          :disabled="pending"
          @page-change="goToPage"
          @page-size-change="changeSize"
        />
      </section>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ApiError } from '~/composables/useApi'
import type { SiteMessage } from '~/composables/useSite'
import { formatDateTime, formatRelativeDateTime } from '~/utils/format'
import { renderUserMarkdown } from '~/utils/richText'

definePageMeta({ middleware: 'auth' })

type MessageFilter = 'all' | 'unread'

const { t, locale } = useI18n()
const localePath = useLocalePath()
const siteApi = useSite()

const messages = ref<SiteMessage[]>([])
const total = ref(0)
const page = ref(1)
const size = ref(20)
const pageSizes = [20, 50, 100]
const filter = ref<MessageFilter>('all')
const pending = ref(false)
const errorMessage = ref('')
const selectedMessageId = ref(0)

const selectedMessage = computed(() => messages.value.find((item) => item.id === selectedMessageId.value) || null)
const selectedMessageContent = computed(() => renderUserMarkdown(selectedMessage.value?.content || '').trim())
const filterOptions = computed(() => [
  { value: 'all' as const, label: t('site.messages.filters.all') },
  { value: 'unread' as const, label: t('site.messages.filters.unread') }
])

useHead(() => ({ title: t('site.messages.title') }))
onMounted(loadMessages)

async function loadMessages() {
  pending.value = true
  errorMessage.value = ''
  try {
    const data = await siteApi.listMessages({
      page: page.value,
      size: size.value,
      isRead: filter.value === 'unread' ? false : undefined
    })
    messages.value = data.list || []
    total.value = data.total || 0
    if (!messages.value.some((item) => item.id === selectedMessageId.value)) {
      selectedMessageId.value = 0
    }
  } catch (error) {
    messages.value = []
    total.value = 0
    selectedMessageId.value = 0
    errorMessage.value = error instanceof ApiError ? error.message : t('common.requestFailed')
  } finally {
    pending.value = false
  }
}

function setFilter(nextFilter: MessageFilter) {
  filter.value = nextFilter
  page.value = 1
  void loadMessages()
}

async function selectMessage(item: SiteMessage) {
  selectedMessageId.value = item.id
  if (item.isRead) return

  await siteApi.markMessageRead(item.id)
  item.isRead = true
  notifyUnreadCountChanged()
}

async function markAllRead() {
  await siteApi.markAllMessagesRead()
  notifyUnreadCountChanged()
  await loadMessages()
}

function goToPage(nextPage: number) {
  page.value = nextPage
  void loadMessages()
}

function changeSize(nextSize: number) {
  size.value = nextSize
  page.value = 1
  void loadMessages()
}

function messageExcerpt(item: SiteMessage) {
  return String(item.content || '').replace(/\s+/g, ' ').trim() || t('site.messages.noContent')
}

function messageSender(item: SiteMessage) {
  if (!item.senderId) return t('site.messages.systemSender')
  return item.sender?.username || `#${item.senderId}`
}

function messageTargetPath(item: SiteMessage) {
  if (!item.targetType || !item.targetId) return ''
  switch (item.targetType) {
    case 'catalog_torrent':
      return `/catalog/torrents/${item.targetId}`
    case 'catalog_request':
      return `/catalog/requests/${item.targetId}`
    case 'forum_topic':
      return `/forum/topics/${item.targetId}`
    default:
      return ''
  }
}

function notifyUnreadCountChanged() {
  if (!import.meta.client) return
  window.dispatchEvent(new CustomEvent('nextpt:messages-read'))
}
</script>
