<template>
  <InteractionCommentList
    v-model:active-report-id="activeReportIdValue"
    v-model:report-reason="reportReasonValue"
    :items="comments"
    :total="total"
    :page="page"
    :page-size="pageSize"
    query-page-key="commentPage"
    anchor-prefix="comment"
    :empty-text="t('catalog.torrents.detail.comments.empty')"
    :like-text="t('catalog.torrents.detail.actions.like')"
    :reward-text="t('catalog.torrents.detail.reward.submit')"
    :quote-text="t('catalog.torrents.detail.comments.quote')"
    :report-text="t('catalog.torrents.detail.actions.report')"
    :report-reason-placeholder="t('catalog.torrents.detail.report.reason')"
    :report-submit-text="t('catalog.torrents.detail.report.submit')"
    :pending="pending"
    :error="error"
    :like-pending-id="likePendingId"
    :report-pending="reportPending"
    :submit-reward="submitReward"
    @page-change="$emit('pageChange', $event)"
    @toggle-like="emitToggleLike"
    @reward-success="emitRewardSuccess"
    @quote="emitQuote"
    @report="$emit('report', $event)"
  />
</template>

<script setup lang="ts">
import type { CommentItem } from '~/composables/useCatalogTorrents'
import type { InteractionCommentItem } from '~/types/interaction'

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
  pageChange: [page: number]
  toggleLike: [comment: CommentItem]
  rewardSuccess: [comment: CommentItem]
  quote: [quote: string]
  report: [commentId: number]
  'update:activeReportId': [commentId: number]
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

function submitReward(item: InteractionCommentItem, amount: number) {
  return props.submitReward(item as CommentItem, amount)
}

function emitToggleLike(item: InteractionCommentItem) {
  emit('toggleLike', item as CommentItem)
}

function emitRewardSuccess(item: InteractionCommentItem) {
  emit('rewardSuccess', item as CommentItem)
}

function emitQuote(quote: string) {
  emit('quote', quote)
}
</script>
