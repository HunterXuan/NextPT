<template>
  <InteractionCommentList
    v-model:active-report-id="activeReportIdValue"
    v-model:report-reason="reportReasonValue"
    :items="replies"
    :total="total"
    :page="page"
    :page-size="pageSize"
    query-page-key="page"
    anchor-prefix="reply"
    :empty-text="t('forum.detail.replies.empty')"
    :like-text="t('forum.detail.actions.like')"
    :reward-text="t('forum.detail.actions.reward')"
    :quote-text="t('forum.detail.replies.quote')"
    :report-text="t('forum.detail.actions.report')"
    :report-reason-placeholder="t('forum.detail.report.reason')"
    :report-submit-text="t('forum.detail.report.submit')"
    :pending="pending"
    :error="error"
    :like-pending-id="likePendingId"
    :report-pending-id="reportPendingId"
    :submit-reward="submitReward"
    @retry="$emit('retry')"
    @page-change="$emit('pageChange', $event)"
    @toggle-like="emitToggleLike"
    @reward-success="emitRewardSuccess"
    @quote="$emit('quote', $event)"
    @report="$emit('report', $event)"
  />
</template>

<script setup lang="ts">
import type { ForumReplyItem } from '~/composables/useForum'
import type { InteractionCommentItem } from '~/types/interaction'

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

const { t } = useI18n()

const activeReportIdValue = computed({
  get: () => props.activeReportId,
  set: (value) => emit('update:activeReportId', value)
})

const reportReasonValue = computed({
  get: () => props.reportReason,
  set: (value) => emit('update:reportReason', value)
})

const likePendingId = computed(() => parsePendingId('like'))
const reportPendingId = computed(() => parsePendingId('report'))

function submitReward(item: InteractionCommentItem, amount: number) {
  return props.submitReward(item as ForumReplyItem, amount)
}

function emitToggleLike(item: InteractionCommentItem) {
  emit('toggleLike', item as ForumReplyItem)
}

function emitRewardSuccess(item: InteractionCommentItem) {
  emit('rewardSuccess', item as ForumReplyItem)
}

function parsePendingId(action: string) {
  const prefix = `${action}:`
  if (!props.actionPending.startsWith(prefix)) return 0
  const id = Number(props.actionPending.slice(prefix.length))
  return Number.isInteger(id) && id > 0 ? id : 0
}
</script>
