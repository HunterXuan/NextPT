<template>
  <div v-if="error" class="flex flex-col items-center justify-center px-4 py-10 text-center">
    <UIcon name="i-lucide-circle-alert" class="size-8 text-red-500" />
    <p class="mt-3 text-sm font-medium text-slate-950 dark:text-white">{{ error }}</p>
    <UButton class="mt-5" color="neutral" variant="outline" size="sm" icon="i-lucide-refresh-cw" @click="$emit('retry')">
      {{ t('common.retry') }}
    </UButton>
  </div>

  <div v-else-if="pending && replies.length === 0" class="space-y-2 px-4 py-4">
    <div v-for="index in 4" :key="index" class="h-24 animate-pulse rounded-md bg-slate-100 dark:bg-slate-800" />
  </div>

  <div v-else-if="replies.length === 0" class="px-4 py-10 text-center text-sm text-slate-500 dark:text-slate-400">
    {{ t('forum.detail.replies.empty') }}
  </div>

  <div v-else class="divide-y divide-slate-200 dark:divide-slate-800">
    <article
      v-for="(reply, index) in replies"
      :id="replyAnchorId(reply)"
      :key="reply.id"
      class="scroll-mt-24 px-4 py-4 transition-colors target:bg-sky-50/60 dark:target:bg-sky-950/30"
    >
      <div class="flex min-w-0 items-start gap-4">
        <img
          v-if="reply.author.avatar"
          class="size-12 shrink-0 rounded-md border border-slate-200 object-cover dark:border-slate-800"
          :src="reply.author.avatar"
          :alt="replyDisplayName(reply)"
          loading="lazy"
        >
        <div v-else :class="replyAvatarClass(reply)" :aria-label="replyDisplayName(reply)">
          {{ replyAvatarInitial(reply) }}
        </div>

        <div class="min-w-0 flex-1">
          <div class="flex min-w-0 items-start justify-between gap-3">
            <div class="min-w-0">
              <div class="flex flex-wrap items-center gap-x-2 gap-y-1">
                <p class="truncate text-sm font-semibold text-slate-950 dark:text-white">{{ replyDisplayName(reply) }}</p>
                <span class="text-xs text-slate-500 dark:text-slate-400">{{ formatDateTime(reply.createdAt, locale) }}</span>
              </div>
            </div>
            <a
              :href="replyAnchorHref(reply)"
              class="shrink-0 rounded px-1 text-xs font-medium tabular-nums text-slate-400 transition hover:bg-slate-100 hover:text-sky-600 dark:text-slate-500 dark:hover:bg-slate-800 dark:hover:text-sky-300"
              @click.prevent="handleReplyAnchorClick(reply)"
            >
              #{{ replyFloor(index) }}
            </a>
          </div>

          <div class="rich-text rich-text-compact mt-2 text-sm leading-6 text-slate-700 dark:text-slate-200" v-html="renderRichText(reply.content)" />

          <div class="mt-2 flex flex-wrap items-center justify-between gap-2">
            <div class="flex flex-wrap items-center gap-1">
              <UTooltip
                :text="t('forum.detail.actions.like')"
                :content="{ side: 'top', sideOffset: 8 }"
                :delay-duration="120"
              >
                <UButton
                  :color="reply.isLiked ? 'error' : 'neutral'"
                  :variant="reply.isLiked ? 'soft' : 'ghost'"
                  size="xs"
                  icon="i-lucide-heart"
                  :loading="actionPending === `like:${reply.id}`"
                  :aria-label="t('forum.detail.actions.like')"
                  @click="$emit('toggleLike', reply)"
                >
                  {{ numberFormatter.format(reply.likeCount || 0) }}
                </UButton>
              </UTooltip>
              <RewardQuickButton
                :count="reply.rewardCount || 0"
                :aria-label="t('forum.detail.actions.reward')"
                :submit-reward="(amount) => submitReward(reply, amount)"
                @success="$emit('rewardSuccess', reply)"
              />
              <UTooltip
                :text="t('forum.detail.replies.quote')"
                :content="{ side: 'top', sideOffset: 8 }"
                :delay-duration="120"
              >
                <UButton
                  color="neutral"
                  variant="ghost"
                  size="xs"
                  icon="i-lucide-quote"
                  :aria-label="t('forum.detail.replies.quote')"
                  @click="$emit('quote', replyQuoteText(reply, index))"
                />
              </UTooltip>
            </div>
            <UTooltip
              :text="t('forum.detail.actions.report')"
              :content="{ side: 'top', sideOffset: 8 }"
              :delay-duration="120"
            >
              <UButton
                color="neutral"
                variant="ghost"
                size="xs"
                icon="i-lucide-flag"
                :aria-label="t('forum.detail.actions.report')"
                @click="toggleReport(reply.id)"
              />
            </UTooltip>
          </div>

          <form v-if="activeReportIdValue === reply.id" class="mt-3 grid gap-2 rounded-md bg-slate-50 p-3 dark:bg-slate-950" @submit.prevent="$emit('report', reply.id)">
            <UTextarea v-model="reportReasonValue" :rows="2" :placeholder="t('forum.detail.report.reason')" :disabled="actionPending === `report:${reply.id}`" />
            <div class="flex justify-end gap-2">
              <UButton color="neutral" variant="ghost" size="xs" type="button" @click="activeReportIdValue = 0">{{ t('common.cancel') }}</UButton>
              <UButton color="error" variant="soft" size="xs" type="submit" :loading="actionPending === `report:${reply.id}`" :disabled="reportReasonValue.trim().length < 5">
                {{ t('forum.detail.report.submit') }}
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
import type { ForumReplyItem } from '~/composables/useForum'
import { formatDateTime } from '~/utils/format'
import { renderUserMarkdown } from '~/utils/richText'

const props = withDefaults(defineProps<{
  replies: ForumReplyItem[]
  total: number
  page: number
  pageSize: number
  pending?: boolean
  error?: string
  actionPending?: string
  activeReportId: number
  reportReason: string
  submitReward: (reply: ForumReplyItem, amount: number) => Promise<void>
}>(), {
  pending: false,
  error: '',
  actionPending: ''
})

const emit = defineEmits<{
  retry: []
  pageChange: [page: number]
  toggleLike: [reply: ForumReplyItem]
  rewardSuccess: [reply: ForumReplyItem]
  quote: [quote: string]
  report: [replyId: number]
  'update:activeReportId': [replyId: number]
  'update:reportReason': [reason: string]
}>()

const { t, locale } = useI18n()
const route = useRoute()

const numberFormatter = computed(() => new Intl.NumberFormat(locale.value))
const totalPages = computed(() => Math.max(1, Math.ceil(props.total / props.pageSize)))
const replyAvatarPalettes = [
  'bg-sky-100 text-sky-700 border-sky-200 dark:bg-sky-950 dark:text-sky-200 dark:border-sky-900',
  'bg-emerald-100 text-emerald-700 border-emerald-200 dark:bg-emerald-950 dark:text-emerald-200 dark:border-emerald-900',
  'bg-amber-100 text-amber-700 border-amber-200 dark:bg-amber-950 dark:text-amber-200 dark:border-amber-900',
  'bg-rose-100 text-rose-700 border-rose-200 dark:bg-rose-950 dark:text-rose-200 dark:border-rose-900',
  'bg-violet-100 text-violet-700 border-violet-200 dark:bg-violet-950 dark:text-violet-200 dark:border-violet-900',
  'bg-cyan-100 text-cyan-700 border-cyan-200 dark:bg-cyan-950 dark:text-cyan-200 dark:border-cyan-900'
]

const activeReportIdValue = computed({
  get: () => props.activeReportId,
  set: (value) => emit('update:activeReportId', value)
})

const reportReasonValue = computed({
  get: () => props.reportReason,
  set: (value) => emit('update:reportReason', value)
})

function submitReward(reply: ForumReplyItem, amount: number) {
  return props.submitReward(reply, amount)
}

function replyFloor(index: number) {
  return (props.page - 1) * props.pageSize + index + 1
}

function replyDisplayName(reply: ForumReplyItem) {
  return reply.author?.username || (reply.author?.id ? `#${reply.author.id}` : '-')
}

function replyAvatarInitial(reply: ForumReplyItem) {
  const value = replyDisplayName(reply).trim()
  return value ? Array.from(value)[0].toUpperCase() : '?'
}

function replyAvatarClass(reply: ForumReplyItem) {
  const base = 'flex size-12 shrink-0 items-center justify-center rounded-md border text-lg font-semibold'
  const paletteIndex = Math.abs(Number(reply.author?.id || reply.id || 0)) % replyAvatarPalettes.length
  return `${base} ${replyAvatarPalettes[paletteIndex]}`
}

function replyAnchorId(reply: ForumReplyItem) {
  return `reply-${reply.id}`
}

function replyAnchorHref(reply: ForumReplyItem) {
  const query = new URLSearchParams()
  for (const [key, value] of Object.entries(route.query)) {
    if (key === 'page') continue
    if (Array.isArray(value)) {
      for (const item of value) {
        if (item != null) query.append(key, item)
      }
      continue
    }
    if (value != null) query.set(key, value)
  }
  query.set('page', String(props.page))
  const queryString = query.toString()
  return `${route.path}${queryString ? `?${queryString}` : ''}#${replyAnchorId(reply)}`
}

function handleReplyAnchorClick(reply: ForumReplyItem) {
  const href = replyAnchorHref(reply)
  document.getElementById(replyAnchorId(reply))?.scrollIntoView({ behavior: 'smooth', block: 'start' })
  window.history.replaceState(null, '', href)
}

function replyQuoteText(reply: ForumReplyItem, index: number) {
  const author = reply.author?.username ? `@${reply.author.username}` : replyDisplayName(reply)
  return `> #${replyFloor(index)} ${author}\n\n`
}

function toggleReport(replyId: number) {
  activeReportIdValue.value = activeReportIdValue.value === replyId ? 0 : replyId
  reportReasonValue.value = ''
}

function renderRichText(content?: string) {
  return renderUserMarkdown(content || '').trim()
}
</script>
