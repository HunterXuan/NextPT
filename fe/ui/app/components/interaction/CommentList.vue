<template>
  <div v-if="error" class="flex flex-col items-center justify-center px-4 py-10 text-center">
    <UIcon name="i-lucide-circle-alert" class="size-8 text-red-500" />
    <p class="mt-3 text-sm font-medium text-slate-950 dark:text-white">{{ error }}</p>
  </div>

  <div v-else-if="pending && items.length === 0" class="space-y-2 px-4 py-4">
    <div v-for="index in 4" :key="index" class="h-24 animate-pulse rounded-md bg-slate-100 dark:bg-slate-800" />
  </div>

  <div v-else-if="items.length === 0" class="px-4 py-10 text-center text-sm text-slate-500 dark:text-slate-400">
    {{ emptyText }}
  </div>

  <div v-else class="divide-y divide-slate-200 dark:divide-slate-800">
    <article
      v-for="(item, index) in items"
      :id="anchorId(item)"
      :key="item.id"
      class="scroll-mt-24 px-4 py-4 transition-colors target:bg-sky-50/60 dark:target:bg-sky-950/30"
    >
      <div class="flex min-w-0 items-start gap-4">
        <IamUserAvatar :user="item.author" size="lg" />

        <div class="min-w-0 flex-1">
          <div class="flex min-w-0 items-start justify-between gap-3">
            <div class="min-w-0">
              <div class="flex flex-wrap items-center gap-x-2 gap-y-1">
                <p class="truncate text-sm font-semibold text-slate-950 dark:text-white">{{ displayName(item) }}</p>
                <span class="text-xs text-slate-500 dark:text-slate-400">{{ formatDateTime(item.createdAt, locale) }}</span>
              </div>
            </div>
            <a
              :href="anchorHref(item)"
              class="shrink-0 rounded px-1 text-xs font-medium tabular-nums text-slate-400 transition hover:bg-slate-100 hover:text-sky-600 dark:text-slate-500 dark:hover:bg-slate-800 dark:hover:text-sky-300"
              @click.prevent="handleAnchorClick(item)"
            >
              #{{ floor(index) }}
            </a>
          </div>

          <div class="rich-text rich-text-compact mt-2 text-sm leading-6 text-slate-700 dark:text-slate-200" v-html="renderRichText(item.content)" />

          <div class="mt-2 flex flex-wrap items-center justify-between gap-2">
            <div class="flex flex-wrap items-center gap-1">
              <UTooltip
                :text="likeText"
                :content="{ side: 'top', sideOffset: 8 }"
                :delay-duration="120"
              >
                <UButton
                  :color="item.isLiked ? 'error' : 'neutral'"
                  :variant="item.isLiked ? 'soft' : 'ghost'"
                  size="xs"
                  icon="i-lucide-heart"
                  :loading="likePendingId === item.id"
                  :aria-label="likeText"
                  @click="$emit('toggleLike', item)"
                >
                  {{ numberFormatter.format(item.likeCount || 0) }}
                </UButton>
              </UTooltip>
              <RewardQuickButton
                :count="item.rewardCount || 0"
                :aria-label="rewardText"
                :submit-reward="(amount) => submitReward(item, amount)"
                @success="$emit('rewardSuccess', item)"
              />
              <UTooltip
                :text="quoteText"
                :content="{ side: 'top', sideOffset: 8 }"
                :delay-duration="120"
              >
                <UButton
                  color="neutral"
                  variant="ghost"
                  size="xs"
                  icon="i-lucide-quote"
                  :aria-label="quoteText"
                  @click="$emit('quote', quoteTextFor(item, index), item)"
                />
              </UTooltip>
            </div>
            <UTooltip
              :text="reportText"
              :content="{ side: 'top', sideOffset: 8 }"
              :delay-duration="120"
            >
              <UButton
                color="neutral"
                variant="ghost"
                size="xs"
                icon="i-lucide-flag"
                :aria-label="reportText"
                @click="toggleReport(item.id)"
              />
            </UTooltip>
          </div>

          <form v-if="activeReportIdValue === item.id" class="mt-3 grid gap-2 rounded-md bg-slate-50 p-3 dark:bg-slate-950" @submit.prevent="$emit('report', item.id)">
            <UTextarea v-model="reportReasonValue" :rows="2" :placeholder="reportReasonPlaceholder" :disabled="isReportPending(item.id)" />
            <div class="flex justify-end gap-2">
              <UButton color="neutral" variant="ghost" size="xs" type="button" @click="activeReportIdValue = 0">{{ t('common.cancel') }}</UButton>
              <UButton color="error" variant="soft" size="xs" type="submit" :loading="isReportPending(item.id)" :disabled="reportReasonValue.trim().length < 5">
                {{ reportSubmitText }}
              </UButton>
            </div>
          </form>
        </div>
      </div>
    </article>
  </div>

  <AppPager
    v-if="totalPages > 1"
    class="border-t border-slate-200 px-4 py-3 dark:border-slate-800"
    size="sm"
    :page="page"
    :total="total"
    :page-size="pageSize"
    :disabled="pending"
    @page-change="$emit('pageChange', $event)"
  />
</template>

<script setup lang="ts">
import type { InteractionCommentItem } from '~/types/interaction'
import { formatDateTime } from '~/utils/format'
import { renderUserMarkdown } from '~/utils/richText'

const props = withDefaults(defineProps<{
  items: InteractionCommentItem[]
  total: number
  page: number
  pageSize: number
  queryPageKey: string
  anchorPrefix: string
  emptyText: string
  likeText: string
  rewardText: string
  quoteText: string
  reportText: string
  reportReasonPlaceholder: string
  reportSubmitText: string
  pending?: boolean
  error?: string
  likePendingId?: number
  reportPendingId?: number
  reportPending?: boolean
  activeReportId: number
  reportReason: string
  submitReward: (item: InteractionCommentItem, amount: number) => Promise<void>
}>(), {
  pending: false,
  error: '',
  likePendingId: 0,
  reportPendingId: 0,
  reportPending: false
})

const emit = defineEmits<{
  pageChange: [page: number]
  toggleLike: [item: InteractionCommentItem]
  rewardSuccess: [item: InteractionCommentItem]
  quote: [quote: string, item: InteractionCommentItem]
  report: [itemId: number]
  'update:activeReportId': [itemId: number]
  'update:reportReason': [reason: string]
}>()

const { t, locale } = useI18n()
const route = useRoute()

const numberFormatter = computed(() => new Intl.NumberFormat(locale.value))
const totalPages = computed(() => Math.max(1, Math.ceil(props.total / props.pageSize)))

const activeReportIdValue = computed({
  get: () => props.activeReportId,
  set: (value) => emit('update:activeReportId', value)
})

const reportReasonValue = computed({
  get: () => props.reportReason,
  set: (value) => emit('update:reportReason', value)
})

function submitReward(item: InteractionCommentItem, amount: number) {
  return props.submitReward(item, amount)
}

function floor(index: number) {
  return (props.page - 1) * props.pageSize + index + 1
}

function displayName(item: InteractionCommentItem) {
  return item.author?.username || (item.author?.id ? `#${item.author.id}` : '-')
}

function anchorId(item: InteractionCommentItem) {
  return `${props.anchorPrefix}-${item.id}`
}

function anchorHref(item: InteractionCommentItem) {
  const query = new URLSearchParams()
  for (const [key, value] of Object.entries(route.query)) {
    if (key === props.queryPageKey) continue
    if (Array.isArray(value)) {
      for (const entry of value) {
        if (entry != null) query.append(key, entry)
      }
      continue
    }
    if (value != null) query.set(key, value)
  }
  query.set(props.queryPageKey, String(props.page))
  const queryString = query.toString()
  return `${route.path}${queryString ? `?${queryString}` : ''}#${anchorId(item)}`
}

function handleAnchorClick(item: InteractionCommentItem) {
  const href = anchorHref(item)
  document.getElementById(anchorId(item))?.scrollIntoView({ behavior: 'smooth', block: 'start' })
  window.history.replaceState(null, '', href)
}

function quoteTextFor(item: InteractionCommentItem, index: number) {
  const author = item.author?.username ? `@${item.author.username}` : displayName(item)
  return `> #${floor(index)} ${author}\n\n`
}

function toggleReport(itemId: number) {
  activeReportIdValue.value = activeReportIdValue.value === itemId ? 0 : itemId
  reportReasonValue.value = ''
}

function isReportPending(itemId: number) {
  return props.reportPending || props.reportPendingId === itemId
}

function renderRichText(content?: string) {
  return renderUserMarkdown(content || '').trim()
}
</script>
