<template>
  <div class="p-4">
    <div v-if="error" class="flex flex-col gap-3 rounded-md border border-red-200 bg-red-50 px-3 py-3 text-sm text-red-700 sm:flex-row sm:items-center sm:justify-between dark:border-red-900 dark:bg-red-950 dark:text-red-200">
      <span>{{ error }}</span>
      <UButton color="error" variant="soft" size="xs" icon="i-lucide-refresh-cw" :loading="pending" @click="$emit('retry')">
        {{ t('common.retry') }}
      </UButton>
    </div>

    <div v-else-if="pending && comments.length === 0" class="mt-4 space-y-3">
      <div v-for="item in 3" :key="item" class="h-28 animate-pulse rounded-md bg-slate-100 dark:bg-slate-800" />
    </div>

    <div v-else-if="comments.length === 0" class="mt-4 rounded-md border border-dashed border-slate-200 px-4 py-10 text-center text-sm text-slate-500 dark:border-slate-800 dark:text-slate-400">
      {{ t('catalog.torrents.detail.comments.empty') }}
    </div>

    <div v-else class="divide-y divide-slate-100 dark:divide-slate-800">
      <article
        v-for="(comment, index) in comments"
        :id="commentAnchorId(comment)"
        :key="comment.id"
        class="scroll-mt-24 py-4 transition-colors target:bg-sky-50/60 dark:target:bg-sky-950/30"
      >
        <div class="flex min-w-0 items-start gap-4">
          <img
            v-if="comment.author.avatar"
            class="size-12 shrink-0 rounded-md border border-slate-200 object-cover dark:border-slate-800"
            :src="comment.author.avatar"
            :alt="commentDisplayName(comment)"
            loading="lazy"
          >
          <div v-else :class="commentAvatarClass(comment)" :aria-label="commentDisplayName(comment)">
            {{ commentAvatarInitial(comment) }}
          </div>

          <div class="min-w-0 flex-1">
            <div class="flex min-w-0 items-start justify-between gap-3">
              <div class="min-w-0">
                <div class="flex flex-wrap items-center gap-x-2 gap-y-1">
                  <p class="truncate text-sm font-semibold text-slate-950 dark:text-white">{{ commentDisplayName(comment) }}</p>
                  <span class="text-xs text-slate-500 dark:text-slate-400">{{ formatDateTime(comment.createdAt, locale) }}</span>
                </div>
              </div>
              <a
                :href="commentAnchorHref(comment)"
                class="shrink-0 rounded px-1 text-xs font-medium tabular-nums text-slate-400 transition hover:bg-slate-100 hover:text-sky-600 dark:text-slate-500 dark:hover:bg-slate-800 dark:hover:text-sky-300"
                @click.prevent="handleCommentAnchorClick(comment)"
              >
                #{{ commentFloor(index) }}
              </a>
            </div>

            <div class="rich-text rich-text-compact mt-2 text-sm leading-6 text-slate-700 dark:text-slate-200" v-html="renderRichText(comment.content)" />

            <div class="mt-2 flex flex-wrap items-center justify-between gap-2">
              <div class="flex flex-wrap items-center gap-1">
                <UTooltip
                  :text="t('catalog.torrents.detail.actions.like')"
                  :content="{ side: 'top', sideOffset: 8 }"
                  :delay-duration="120"
                >
                  <UButton
                    :color="comment.isLiked ? 'error' : 'neutral'"
                    :variant="comment.isLiked ? 'soft' : 'ghost'"
                    size="xs"
                    icon="i-lucide-heart"
                    :loading="likePendingId === comment.id"
                    :aria-label="t('catalog.torrents.detail.actions.like')"
                    @click="$emit('toggleLike', comment)"
                  >
                    {{ numberFormatter.format(comment.likeCount) }}
                  </UButton>
                </UTooltip>
                <RewardQuickButton
                  :count="comment.rewardCount"
                  :aria-label="t('catalog.torrents.detail.reward.submit')"
                  :submit-reward="(amount) => submitReward(comment, amount)"
                  @success="$emit('rewardSuccess', comment)"
                />
                <UTooltip
                  :text="t('catalog.torrents.detail.comments.quote')"
                  :content="{ side: 'top', sideOffset: 8 }"
                  :delay-duration="120"
                >
                  <UButton
                    color="neutral"
                    variant="ghost"
                    size="xs"
                    icon="i-lucide-quote"
                    :aria-label="t('catalog.torrents.detail.comments.quote')"
                    @click="$emit('quote', commentQuoteText(comment, index))"
                  />
                </UTooltip>
              </div>
              <UTooltip
                :text="t('catalog.torrents.detail.actions.report')"
                :content="{ side: 'top', sideOffset: 8 }"
                :delay-duration="120"
              >
                <UButton
                  color="neutral"
                  variant="ghost"
                  size="xs"
                  icon="i-lucide-flag"
                  :aria-label="t('catalog.torrents.detail.actions.report')"
                  @click="toggleReport(comment.id)"
                />
              </UTooltip>
            </div>

            <form v-if="activeReportIdValue === comment.id" class="mt-3 grid gap-2 rounded-md bg-slate-50 p-3 dark:bg-slate-950" @submit.prevent="$emit('report', comment.id)">
              <UTextarea v-model="reportReasonValue" :rows="2" :placeholder="t('catalog.torrents.detail.report.reason')" :disabled="reportPending" />
              <div class="flex justify-end gap-2">
                <UButton color="neutral" variant="ghost" size="xs" type="button" @click="activeReportIdValue = 0">{{ t('common.cancel') }}</UButton>
                <UButton color="error" variant="soft" size="xs" type="submit" :loading="reportPending" :disabled="reportReasonValue.trim().length < 5">
                  {{ t('catalog.torrents.detail.report.submit') }}
                </UButton>
              </div>
            </form>
          </div>
        </div>
      </article>
    </div>

    <AppPager
      v-if="totalPages > 1"
      class="mt-4 border-t border-slate-200 pt-4 dark:border-slate-800"
      size="sm"
      :page="page"
      :total="total"
      :page-size="pageSize"
      :disabled="pending"
      @page-change="$emit('pageChange', $event)"
    />
  </div>
</template>

<script setup lang="ts">
import type { CommentItem } from '~/composables/useCatalogTorrents'
import { formatDateTime } from '~/utils/format'
import { renderUserMarkdown } from '~/utils/richText'

const props = withDefaults(defineProps<{
  comments: CommentItem[]
  total: number
  page: number
  pageSize: number
  pending?: boolean
  error?: string
  likePendingId?: number
  activeReportId: number
  reportReason: string
  reportPending?: boolean
  submitReward: (comment: CommentItem, amount: number) => Promise<void>
}>(), {
  pending: false,
  error: '',
  likePendingId: 0,
  reportPending: false
})

const emit = defineEmits<{
  retry: []
  pageChange: [page: number]
  toggleLike: [comment: CommentItem]
  rewardSuccess: [comment: CommentItem]
  quote: [quote: string]
  report: [commentId: number]
  'update:activeReportId': [commentId: number]
  'update:reportReason': [reason: string]
}>()

const { t, locale } = useI18n()
const route = useRoute()

const numberFormatter = computed(() => new Intl.NumberFormat(locale.value))
const totalPages = computed(() => Math.max(1, Math.ceil(props.total / props.pageSize)))
const commentAvatarPalettes = [
  'border-sky-200 bg-sky-50 text-sky-700 dark:border-sky-900 dark:bg-sky-950 dark:text-sky-300',
  'border-emerald-200 bg-emerald-50 text-emerald-700 dark:border-emerald-900 dark:bg-emerald-950 dark:text-emerald-300',
  'border-amber-200 bg-amber-50 text-amber-700 dark:border-amber-900 dark:bg-amber-950 dark:text-amber-300',
  'border-rose-200 bg-rose-50 text-rose-700 dark:border-rose-900 dark:bg-rose-950 dark:text-rose-300',
  'border-indigo-200 bg-indigo-50 text-indigo-700 dark:border-indigo-900 dark:bg-indigo-950 dark:text-indigo-300',
  'border-teal-200 bg-teal-50 text-teal-700 dark:border-teal-900 dark:bg-teal-950 dark:text-teal-300'
]

const activeReportIdValue = computed({
  get: () => props.activeReportId,
  set: (value) => emit('update:activeReportId', value)
})

const reportReasonValue = computed({
  get: () => props.reportReason,
  set: (value) => emit('update:reportReason', value)
})

function submitReward(comment: CommentItem, amount: number) {
  return props.submitReward(comment, amount)
}

function commentFloor(index: number) {
  return (props.page - 1) * props.pageSize + index + 1
}

function commentDisplayName(comment: CommentItem) {
  return comment.author?.username || (comment.author?.id ? `#${comment.author.id}` : '-')
}

function commentAvatarInitial(comment: CommentItem) {
  const value = commentDisplayName(comment).trim()
  return value ? Array.from(value)[0].toUpperCase() : '?'
}

function commentAvatarClass(comment: CommentItem) {
  const base = 'flex size-12 shrink-0 items-center justify-center rounded-md border text-lg font-semibold'
  const paletteIndex = Math.abs(Number(comment.author?.id || comment.id || 0)) % commentAvatarPalettes.length
  return `${base} ${commentAvatarPalettes[paletteIndex]}`
}

function commentAnchorId(comment: CommentItem) {
  return `comment-${comment.id}`
}

function commentAnchorHref(comment: CommentItem) {
  const query = new URLSearchParams()
  for (const [key, value] of Object.entries(route.query)) {
    if (key === 'commentPage') continue
    if (Array.isArray(value)) {
      for (const item of value) {
        if (item != null) query.append(key, item)
      }
      continue
    }
    if (value != null) query.set(key, value)
  }
  query.set('commentPage', String(props.page))
  const queryString = query.toString()
  return `${route.path}${queryString ? `?${queryString}` : ''}#${commentAnchorId(comment)}`
}

function handleCommentAnchorClick(comment: CommentItem) {
  const href = commentAnchorHref(comment)
  document.getElementById(commentAnchorId(comment))?.scrollIntoView({ behavior: 'smooth', block: 'start' })
  window.history.replaceState(null, '', href)
}

function commentQuoteText(comment: CommentItem, index: number) {
  const author = comment.author?.username ? `@${comment.author.username}` : commentDisplayName(comment)
  return `> #${commentFloor(index)} ${author}\n\n`
}

function toggleReport(commentId: number) {
  activeReportIdValue.value = activeReportIdValue.value === commentId ? 0 : commentId
  reportReasonValue.value = ''
}

function renderRichText(content?: string) {
  return renderUserMarkdown(content || '').trim()
}
</script>
