<template>
  <section class="overflow-hidden rounded-lg border border-slate-200 bg-white dark:border-slate-800 dark:bg-slate-900">
    <div class="border-b border-slate-200 px-5 py-4 dark:border-slate-800">
      <div class="flex flex-col gap-4 md:flex-row md:items-start md:justify-between">
        <div class="min-w-0">
          <div v-if="topic.isSticky || topic.isLocked" class="mb-3 flex flex-wrap items-center gap-2">
            <UBadge v-if="topic.isSticky" color="primary" variant="soft">{{ t('forum.topicList.badges.sticky') }}</UBadge>
            <UBadge v-if="topic.isLocked" color="neutral" variant="outline">{{ t('forum.topicList.badges.locked') }}</UBadge>
          </div>
          <h1 class="break-words [overflow-wrap:anywhere] text-2xl font-semibold leading-8 text-slate-950 dark:text-white">
            {{ topic.subject || t('forum.detail.titleFallback', { id: topic.id }) }}
          </h1>
          <p class="mt-3 flex flex-wrap items-center gap-x-2 gap-y-1 text-sm text-slate-500 dark:text-slate-400">
            <IamUserPopover :user="topic.author" :fallback="userDisplayName(topic.author)" class="font-medium text-slate-600 dark:text-slate-300" />
            <span class="text-slate-300 dark:text-slate-700">/</span>
            <UTooltip
              :text="formatDateTime(topic.createdAt, locale)"
              :content="{ side: 'top', sideOffset: 8 }"
              :delay-duration="120"
            >
              <span>{{ relativeDateTime(topic.createdAt) }}</span>
            </UTooltip>
            <span class="text-slate-300 dark:text-slate-700">/</span>
            <span>{{ t('forum.topicList.meta.views', { count: numberFormatter.format(topic.views) }) }}</span>
          </p>
        </div>

        <div class="flex shrink-0 items-center gap-2 md:pt-1">
          <UTooltip
            :text="t('forum.detail.actions.like')"
            :content="{ side: 'top', sideOffset: 8 }"
            :delay-duration="120"
          >
            <UButton
              :color="topic.isLiked ? 'error' : 'neutral'"
              :variant="topic.isLiked ? 'soft' : 'outline'"
              size="sm"
              icon="i-lucide-heart"
              :loading="actionPending === 'like'"
              :aria-label="t('forum.detail.actions.like')"
              @click="$emit('toggleLike')"
            >
              {{ numberFormatter.format(topic.likeCount || 0) }}
            </UButton>
          </UTooltip>
          <UTooltip
            :text="topic.isBookmarked ? t('forum.detail.actions.bookmarked') : t('forum.detail.actions.bookmark')"
            :content="{ side: 'top', sideOffset: 8 }"
            :delay-duration="120"
          >
            <UButton
              color="neutral"
              :variant="topic.isBookmarked ? 'soft' : 'outline'"
              size="sm"
              :icon="topic.isBookmarked ? 'i-lucide-bookmark-check' : 'i-lucide-bookmark-plus'"
              :loading="actionPending === 'bookmark'"
              :aria-label="topic.isBookmarked ? t('forum.detail.actions.bookmarked') : t('forum.detail.actions.bookmark')"
              @click="$emit('toggleBookmark')"
            />
          </UTooltip>
        </div>
      </div>
    </div>

    <div class="px-5 py-5">
      <div class="rich-text text-sm leading-6 text-slate-700 dark:text-slate-200" v-html="renderRichText(topic.content)" />
    </div>

    <div v-if="appends.length > 0" class="space-y-3 border-t border-slate-200 px-5 py-4 dark:border-slate-800">
      <article v-for="(append, index) in appends" :key="`${append.created_at || append.createdAt || index}`" class="rounded-md bg-slate-50 px-3 py-3 dark:bg-slate-950">
        <div class="flex items-center justify-between gap-3">
          <p class="text-xs font-medium text-slate-500 dark:text-slate-400">
            {{ t('forum.detail.append.index', { index: index + 1 }) }}
          </p>
          <p class="text-xs text-slate-500 dark:text-slate-400">{{ formatDateTime(append.created_at || append.createdAt, locale) }}</p>
        </div>
        <div class="rich-text rich-text-compact mt-2 text-sm leading-6 text-slate-700 dark:text-slate-200" v-html="renderRichText(append.content || '')" />
      </article>
    </div>

    <div v-if="activePanel === 'append'" class="border-t border-slate-200 bg-slate-50/60 px-5 py-4 dark:border-slate-800 dark:bg-slate-950/40">
      <RichTextComposer
        v-model="appendContentValue"
        v-model:mode="appendEditorModeValue"
        :placeholder="t('forum.detail.append.placeholder')"
        :disabled="actionPending === 'append'"
        :pending="actionPending === 'append'"
        :submit-disabled="appendContentValue.trim().length < 2"
        :submit-label="t('forum.detail.append.submit')"
        submit-icon="i-lucide-plus-square"
        :write-label="t('common.editor.edit')"
        :preview-label="t('common.editor.preview')"
        :preview-empty="t('common.editor.previewEmpty')"
        @submit="$emit('append')"
      >
        <template #actions>
          <div class="flex justify-end gap-2">
            <UButton type="button" color="neutral" variant="ghost" size="sm" @click="$emit('closePanel')">
              {{ t('common.cancel') }}
            </UButton>
            <UButton type="submit" color="primary" size="sm" icon="i-lucide-plus-square" :loading="actionPending === 'append'" :disabled="appendContentValue.trim().length < 2">
              {{ t('forum.detail.append.submit') }}
            </UButton>
          </div>
        </template>
      </RichTextComposer>
    </div>

    <div class="flex flex-col gap-3 border-t border-slate-200 bg-slate-50/70 px-5 py-3 sm:flex-row sm:items-center sm:justify-between dark:border-slate-800 dark:bg-slate-950/50">
      <UButton color="neutral" :variant="activePanel === 'report' ? 'soft' : 'ghost'" size="sm" icon="i-lucide-flag" class="w-fit" @click="$emit('openPanel', 'report')">
        {{ t('forum.detail.actions.report') }}
      </UButton>
      <div v-if="canUseOwnerActions" class="flex flex-wrap items-center justify-end gap-2">
        <UButton v-if="canEdit" color="neutral" variant="ghost" size="sm" icon="i-lucide-pencil" @click="$emit('edit')">
          {{ t('forum.detail.actions.edit') }}
        </UButton>
        <UButton v-if="canAppend" color="neutral" :variant="activePanel === 'append' ? 'soft' : 'ghost'" size="sm" icon="i-lucide-plus-square" @click="$emit('openPanel', 'append')">
          {{ t('forum.detail.actions.append') }}
        </UButton>
      </div>
    </div>

    <div v-if="activePanel === 'report'" class="border-t border-slate-200 bg-slate-50/60 px-5 py-4 dark:border-slate-800 dark:bg-slate-950/40">
      <form class="space-y-3" @submit.prevent="$emit('report')">
        <UTextarea v-model="reportReasonValue" class="w-full" :rows="3" :placeholder="t('forum.detail.report.reason')" :disabled="actionPending === 'report'" />
        <div class="flex justify-end gap-2">
          <UButton type="button" color="neutral" variant="ghost" size="sm" @click="$emit('closePanel')">
            {{ t('common.cancel') }}
          </UButton>
          <UButton type="submit" color="error" variant="soft" size="sm" icon="i-lucide-flag" :loading="actionPending === 'report'" :disabled="reportReasonValue.trim().length < 5">
            {{ t('forum.detail.report.submit') }}
          </UButton>
        </div>
      </form>
    </div>
  </section>
</template>

<script setup lang="ts">
import type { ForumTopicAppend, ForumTopicDetail } from '~/composables/useForum'
import type { UserSummary } from '~/types/iam'
import { formatDateTime } from '~/utils/format'
import { renderUserMarkdown } from '~/utils/richText'

type TopicPanel = 'append' | 'report'
type EditorMode = 'write' | 'preview'

const props = withDefaults(defineProps<{
  topic: ForumTopicDetail
  appends: ForumTopicAppend[]
  activePanel?: TopicPanel | null
  actionPending?: string
  canEdit?: boolean
  canAppend?: boolean
  canUseOwnerActions?: boolean
  appendContent: string
  appendEditorMode: EditorMode
  reportReason: string
}>(), {
  activePanel: null,
  actionPending: '',
  canEdit: false,
  canAppend: false,
  canUseOwnerActions: false
})

const emit = defineEmits<{
  toggleLike: []
  toggleBookmark: []
  edit: []
  openPanel: [panel: TopicPanel]
  closePanel: []
  append: []
  report: []
  'update:appendContent': [value: string]
  'update:appendEditorMode': [value: EditorMode]
  'update:reportReason': [value: string]
}>()

const { t, locale } = useI18n()

const numberFormatter = computed(() => new Intl.NumberFormat(locale.value))
const relativeTimeFormatter = computed(() => new Intl.RelativeTimeFormat(locale.value, { numeric: 'auto' }))
const appendContentValue = computed({
  get: () => props.appendContent,
  set: (value) => emit('update:appendContent', value)
})
const appendEditorModeValue = computed({
  get: () => props.appendEditorMode,
  set: (value) => emit('update:appendEditorMode', value)
})
const reportReasonValue = computed({
  get: () => props.reportReason,
  set: (value) => emit('update:reportReason', value)
})

function relativeDateTime(value?: string | null) {
  const date = parseDateTime(value)
  if (!date) return '-'

  const diffSeconds = Math.round((date.getTime() - Date.now()) / 1000)
  const absSeconds = Math.abs(diffSeconds)

  if (absSeconds < 45) return relativeTimeFormatter.value.format(0, 'second')
  if (absSeconds < 45 * 60) return relativeTimeFormatter.value.format(Math.round(diffSeconds / 60), 'minute')
  if (absSeconds < 22 * 60 * 60) return relativeTimeFormatter.value.format(Math.round(diffSeconds / 60 / 60), 'hour')
  if (absSeconds < 30 * 24 * 60 * 60) return relativeTimeFormatter.value.format(Math.round(diffSeconds / 60 / 60 / 24), 'day')
  if (absSeconds < 12 * 30 * 24 * 60 * 60) return relativeTimeFormatter.value.format(Math.round(diffSeconds / 60 / 60 / 24 / 30), 'month')
  return relativeTimeFormatter.value.format(Math.round(diffSeconds / 60 / 60 / 24 / 365), 'year')
}

function parseDateTime(value?: string | null) {
  if (!value) return null

  const date = new Date(value)
  if (!Number.isNaN(date.getTime())) return date

  const normalizedDate = new Date(value.replace(' ', 'T'))
  return Number.isNaN(normalizedDate.getTime()) ? null : normalizedDate
}

function renderRichText(content?: string) {
  return renderUserMarkdown(content || '').trim()
}

function userDisplayName(user?: UserSummary | null) {
  if (!user) return '-'
  return user.username || (user.id > 0 ? `#${user.id}` : '-')
}
</script>
