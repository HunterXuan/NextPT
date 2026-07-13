<template>
  <section :id="sectionId" class="scroll-mt-24 overflow-hidden rounded-lg border border-slate-200 bg-white dark:border-slate-800 dark:bg-slate-900">
    <div class="flex items-center justify-between gap-3 border-b border-slate-200 px-4 py-3 dark:border-slate-800">
      <h2 class="text-base font-semibold text-slate-950 dark:text-white">{{ t('catalog.torrents.detail.comments.title') }}</h2>
      <UBadge color="neutral" variant="soft">
        {{ t('catalog.torrents.detail.comments.summary', { count: numberFormatter.format(total) }) }}
      </UBadge>
    </div>

    <form :id="composerId" class="border-b border-slate-200 bg-slate-50/60 p-4 dark:border-slate-800 dark:bg-slate-950/40" @submit.prevent="emit('submit')">
      <UTextarea
        v-if="editorModeValue === 'write'"
        ref="commentInput"
        v-model="contentValue"
        class="w-full"
        :rows="4"
        :placeholder="canCreate ? t('catalog.torrents.detail.comments.placeholder') : t('common.noPermission')"
        :disabled="!canCreate || submitPending"
      />
      <div v-else class="min-h-28 rounded-md border border-slate-200 bg-slate-50 px-3 py-2.5 dark:border-slate-800 dark:bg-slate-950">
        <div v-if="renderedPreview" class="rich-text rich-text-compact" v-html="renderedPreview" />
        <p v-else class="text-sm text-slate-500 dark:text-slate-400">{{ t('catalog.torrents.detail.comments.previewEmpty') }}</p>
      </div>

      <div class="mt-3 flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
        <div class="inline-flex w-fit rounded-md border border-slate-200 bg-slate-50 p-0.5 dark:border-slate-800 dark:bg-slate-950">
          <button type="button" :class="editorTabClass('write')" @click="editorModeValue = 'write'">
            {{ t('catalog.torrents.detail.comments.edit') }}
          </button>
          <button type="button" :class="editorTabClass('preview')" @click="editorModeValue = 'preview'">
            {{ t('catalog.torrents.detail.comments.preview') }}
          </button>
        </div>

        <AppPermissionButton
          :permission="Permission.CatalogCommentCreate"
          type="submit"
          color="primary"
          icon="i-lucide-send"
          :tooltip="t('catalog.torrents.detail.comments.submit')"
          :loading="submitPending"
          :disabled="!canCreate || contentValue.trim().length < 3"
        >
          {{ t('catalog.torrents.detail.comments.submit') }}
        </AppPermissionButton>
      </div>
    </form>

    <CatalogCommentList
      v-model:active-report-id="activeReportIdValue"
      v-model:report-reason="reportReasonValue"
      :comments="comments"
      :total="total"
      :page="page"
      :page-size="pageSize"
      :pending="pending"
      :error="error"
      :like-pending-id="likePendingId"
      :report-pending="reportPending"
      :submit-reward="submitReward"
      @page-change="emit('pageChange', $event)"
      @toggle-like="emit('toggleLike', $event)"
      @reward-success="emit('rewardSuccess', $event)"
      @quote="insertQuote"
      @report="emit('report', $event)"
    />
  </section>
</template>

<script setup lang="ts">
import type { CommentItem } from '~/composables/useCatalogTorrents'
import { renderUserMarkdown } from '~/utils/richText'

type EditorMode = 'write' | 'preview'

const props = withDefaults(defineProps<{
  content: string
  editorMode: EditorMode
  comments: CommentItem[]
  total: number
  page: number
  pageSize: number
  activeReportId: number
  reportReason: string
  submitReward: (comment: CommentItem, amount: number) => Promise<void>
  sectionId?: string
  composerId?: string
  canCreate?: boolean
  submitPending?: boolean
  pending?: boolean
  error?: string
  likePendingId?: number
  reportPending?: boolean
}>(), {
  sectionId: 'catalog-comments',
  composerId: 'catalog-comment-composer',
  canCreate: true,
  submitPending: false,
  pending: false,
  error: '',
  likePendingId: 0,
  reportPending: false
})

const emit = defineEmits<{
  submit: []
  pageChange: [page: number]
  toggleLike: [comment: CommentItem]
  rewardSuccess: [comment: CommentItem]
  report: [commentId: number]
  'update:content': [content: string]
  'update:editorMode': [mode: EditorMode]
  'update:activeReportId': [commentId: number]
  'update:reportReason': [reason: string]
}>()

const { t, locale } = useI18n()
const commentInput = ref<{ textareaRef?: HTMLTextAreaElement } | null>(null)
const numberFormatter = computed(() => new Intl.NumberFormat(locale.value))

const contentValue = computed({
  get: () => props.content,
  set: value => emit('update:content', value)
})

const editorModeValue = computed({
  get: () => props.editorMode,
  set: value => emit('update:editorMode', value)
})

const activeReportIdValue = computed({
  get: () => props.activeReportId,
  set: value => emit('update:activeReportId', value)
})

const reportReasonValue = computed({
  get: () => props.reportReason,
  set: value => emit('update:reportReason', value)
})

const renderedPreview = computed(() => renderUserMarkdown(contentValue.value).trim())

function editorTabClass(mode: EditorMode) {
  const base = 'h-8 rounded px-3 text-sm font-medium transition focus:outline-none focus-visible:ring-2 focus-visible:ring-sky-200 dark:focus-visible:ring-sky-900'
  if (editorModeValue.value === mode) {
    return `${base} bg-white text-slate-950 shadow-sm dark:bg-slate-800 dark:text-white`
  }
  return `${base} text-slate-500 hover:text-slate-950 dark:text-slate-400 dark:hover:text-white`
}

function insertQuote(quote: string) {
  if (!props.canCreate) return
  contentValue.value = contentValue.value.trim() ? `${contentValue.value.trim()}\n\n${quote}` : quote
  editorModeValue.value = 'write'
  void nextTick(() => {
    document.getElementById(props.composerId)?.scrollIntoView({ behavior: 'smooth', block: 'start' })
    commentInput.value?.textareaRef?.focus()
  })
}
</script>
