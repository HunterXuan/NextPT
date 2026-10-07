<template>
  <section class="overflow-hidden rounded-lg border border-neutral-200 bg-white dark:border-neutral-800 dark:bg-neutral-900">
    <div v-if="pending" class="divide-y divide-neutral-200 dark:divide-neutral-800">
      <div v-for="index in 8" :key="index" class="grid gap-3 px-3 py-3 md:grid-cols-[minmax(0,1fr)_auto] md:items-center">
        <div class="flex gap-3">
          <div class="size-10 animate-pulse rounded-md bg-neutral-100 dark:bg-neutral-800/70" />
          <div class="min-w-0 flex-1 space-y-2">
            <div class="h-4 w-4/5 animate-pulse rounded bg-neutral-200 dark:bg-neutral-800" />
            <div class="h-3 w-3/5 animate-pulse rounded bg-neutral-100 dark:bg-neutral-800/70" />
          </div>
        </div>
        <div class="h-7 w-12 animate-pulse rounded-full bg-neutral-100 md:justify-self-center dark:bg-neutral-800/70" />
      </div>
    </div>

    <div v-else-if="errorMessage" class="flex flex-col items-center justify-center px-4 py-16 text-center">
      <UIcon name="i-lucide-circle-alert" class="size-9 text-red-500" />
      <p class="mt-3 text-sm font-medium text-neutral-950 dark:text-white">{{ errorMessage }}</p>
    </div>

    <div v-else-if="topics.length === 0" class="flex flex-col items-center justify-center px-4 py-16 text-center">
      <UIcon :name="emptyIcon" class="size-9 text-neutral-400" />
      <p class="mt-3 text-sm font-medium text-neutral-950 dark:text-white">{{ emptyTitle }}</p>
      <p class="mt-1 max-w-md text-sm text-neutral-500 dark:text-neutral-400">{{ emptyDescription }}</p>
      <AppPermissionButton
        v-if="emptyActionTo && emptyActionLabel"
        class="mt-5"
        color="primary"
        :icon="emptyActionIcon"
        :to="emptyActionTo"
        :allowed="!emptyActionDisabled"
        :disabled-tooltip="emptyActionDisabledText"
        :tooltip="emptyActionLabel"
      >
        {{ emptyActionLabel }}
      </AppPermissionButton>
    </div>

    <div v-else class="divide-y divide-neutral-200 dark:divide-neutral-800">
      <article
        v-for="topic in topics"
        :key="topic.id"
        class="grid gap-3 px-3 py-3 transition-colors hover:bg-neutral-50 md:grid-cols-[minmax(0,1fr)_auto] md:items-center dark:hover:bg-neutral-950/70"
      >
        <div class="flex min-w-0 gap-3">
          <IamUserAvatar
            :user="topic.author"
            :id="topic.author?.id || topic.id"
            :username="topic.author?.username || topicAuthor(topic)"
            size="md"
            :title="topicAuthor(topic)"
          />

          <div class="min-w-0 flex-1">
            <h2 class="flex min-w-0 items-center gap-1.5 text-sm font-medium leading-5">
              <span
                v-if="topic.isSticky"
                class="inline-flex h-5 shrink-0 items-center rounded bg-neutral-900 px-1.5 text-[10px] font-semibold leading-none text-white dark:bg-white dark:text-neutral-950"
              >
                {{ $t('forum.topicList.badges.sticky') }}
              </span>
              <span
                v-if="topic.isLocked"
                class="inline-flex h-5 shrink-0 items-center rounded bg-neutral-100 px-1.5 text-[10px] font-semibold leading-none text-neutral-600 dark:bg-neutral-800 dark:text-neutral-300"
              >
                {{ $t('forum.topicList.badges.locked') }}
              </span>
              <NuxtLink
                :to="localePath(`/forum/topics/${topic.id}`)"
                class="min-w-0 truncate text-neutral-950 hover:text-primary-700 dark:text-white dark:hover:text-primary-300"
              >
                {{ topic.subject || `#${topic.id}` }}
              </NuxtLink>
            </h2>

            <div class="mt-1 flex flex-wrap items-center gap-x-2 gap-y-1 text-xs text-neutral-500 dark:text-neutral-400">
              <IamUserPopover :user="topic.author" :fallback="topicAuthor(topic)" class="font-medium text-neutral-600 dark:text-neutral-300" />
              <span class="text-neutral-300 dark:text-neutral-700">/</span>
              <span :title="formatDateTime(lastActivityAt(topic), locale)">
                {{ relativeDateTime(lastActivityAt(topic)) }}
              </span>
              <span v-if="hasReplies(topic)" class="text-neutral-300 dark:text-neutral-700">/</span>
              <span v-if="hasReplies(topic)">{{ $t('forum.topicList.meta.lastReplyBy') }}</span>
              <IamUserPopover v-if="hasReplies(topic)" :user="topic.lastReplyUser" :fallback="lastReplyUser(topic)" class="font-medium text-neutral-600 dark:text-neutral-300" />
              <span v-if="showViews" class="text-neutral-300 dark:text-neutral-700">/</span>
              <span v-if="showViews">{{ $t('forum.topicList.meta.views', { count: numberFormatter.format(topic.views) }) }}</span>
            </div>
          </div>
        </div>

        <div class="flex items-center gap-2 md:justify-self-end">
          <NuxtLink
            :to="localePath(`/forum/topics/${topic.id}`)"
            class="inline-flex h-7 min-w-10 items-center justify-center rounded-full bg-neutral-100 px-3 text-sm font-semibold text-neutral-600 transition-colors hover:bg-primary-100 hover:text-primary-800 dark:bg-neutral-800 dark:text-neutral-300 dark:hover:bg-primary-950 dark:hover:text-primary-200"
            :title="$t('forum.topicList.table.replies')"
          >
            {{ numberFormatter.format(topic.replyCount) }}
          </NuxtLink>
          <slot name="topic-actions" :topic="topic" />
        </div>
      </article>
    </div>
  </section>
</template>

<script setup lang="ts">
import type { ForumTopicListItem } from '~/composables/useForum'
import type { UserSummary } from '~/types/iam'
import { formatDateTime } from '~/utils/format'

withDefaults(defineProps<{
  topics: ForumTopicListItem[]
  pending?: boolean
  errorMessage?: string
  emptyTitle: string
  emptyDescription: string
  emptyIcon?: string
  emptyActionIcon?: string
  emptyActionLabel?: string
  emptyActionTo?: string
  emptyActionDisabled?: boolean
  emptyActionDisabledText?: string
  showViews?: boolean
}>(), {
  pending: false,
  errorMessage: '',
  emptyIcon: 'i-lucide-message-square-off',
  emptyActionIcon: 'i-lucide-square-pen',
  emptyActionLabel: '',
  emptyActionTo: '',
  emptyActionDisabled: false,
  emptyActionDisabledText: '',
  showViews: true
})

const { locale } = useI18n()
const localePath = useLocalePath()

const numberFormatter = computed(() => new Intl.NumberFormat(locale.value))
const relativeTimeFormatter = computed(() => new Intl.RelativeTimeFormat(locale.value, { numeric: 'auto' }))

function topicAuthor(topic: ForumTopicListItem) {
  return userDisplayName(topic.author)
}

function userDisplayName(user?: UserSummary | null) {
  if (!user) return '-'
  return user.username || (user.id > 0 ? `#${user.id}` : '-')
}

function hasReplies(topic: ForumTopicListItem) {
  return Number(topic.replyCount || 0) > 0
}

function lastActivityAt(topic: ForumTopicListItem) {
  return hasReplies(topic) && topic.lastReplyAt ? topic.lastReplyAt : topic.createdAt
}

function lastReplyUser(topic: ForumTopicListItem) {
  return userDisplayName(topic.lastReplyUser)
}

function relativeDateTime(value?: string | null) {
  const date = parseDateTime(value)
  if (!date) return '-'

  const diffSeconds = Math.round((date.getTime() - Date.now()) / 1000)
  const absSeconds = Math.abs(diffSeconds)

  if (absSeconds < 45) return relativeTimeFormatter.value.format(0, 'second')
  if (absSeconds < 45 * 60) return relativeTimeFormatter.value.format(Math.round(diffSeconds / 60), 'minute')
  if (absSeconds < 22 * 60 * 60) return relativeTimeFormatter.value.format(Math.round(diffSeconds / 60 / 60), 'hour')
  if (absSeconds < 26 * 60 * 60) return relativeTimeFormatter.value.format(Math.round(diffSeconds / 60 / 60 / 24), 'day')
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
</script>
