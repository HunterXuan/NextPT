<template>
  <div class="min-h-[calc(100vh-4rem)] bg-slate-50 py-6 dark:bg-slate-950">
    <div class="w-full px-3 sm:px-4 lg:px-5">
      <div class="mb-6 flex flex-col gap-4 lg:flex-row lg:items-start lg:justify-between">
        <div class="min-w-0">
          <h1 class="break-words text-2xl font-semibold text-slate-950 dark:text-white">
            {{ topic?.subject || $t('forum.detail.titleFallback', { id: topicId }) }}
          </h1>
          <p v-if="topic" class="mt-2 text-sm text-slate-500 dark:text-slate-400">
            {{ $t('forum.topicList.byline', { user: topic.username || `#${topic.userId}`, time: formatDateTime(topic.createdAt, locale) }) }}
          </p>
        </div>

        <div class="grid grid-cols-2 gap-2 sm:flex sm:items-center">
          <UButton
            color="neutral"
            :variant="topic?.isLiked ? 'soft' : 'outline'"
            :icon="topic?.isLiked ? 'i-lucide-heart' : 'i-lucide-heart-plus'"
            :loading="topicActionPending === 'like'"
            :disabled="!topic || pending"
            @click="handleToggleTopicLike"
          >
            {{ topic?.isLiked ? $t('forum.detail.actions.liked') : $t('forum.detail.actions.like') }}
          </UButton>
          <UButton
            color="neutral"
            :variant="topic?.isBookmarked ? 'soft' : 'outline'"
            :icon="topic?.isBookmarked ? 'i-lucide-bookmark-check' : 'i-lucide-bookmark-plus'"
            :loading="topicActionPending === 'bookmark'"
            :disabled="!topic || pending"
            @click="handleToggleTopicBookmark"
          >
            {{ topic?.isBookmarked ? $t('forum.detail.actions.bookmarked') : $t('forum.detail.actions.bookmark') }}
          </UButton>
        </div>
      </div>

      <div v-if="pending" class="grid gap-6 lg:grid-cols-[minmax(0,1fr)_340px]">
        <div class="space-y-4">
          <div class="h-64 animate-pulse rounded-lg border border-slate-200 bg-white dark:border-slate-800 dark:bg-slate-900" />
          <div class="h-80 animate-pulse rounded-lg border border-slate-200 bg-white dark:border-slate-800 dark:bg-slate-900" />
        </div>
        <div class="h-64 animate-pulse rounded-lg border border-slate-200 bg-white dark:border-slate-800 dark:bg-slate-900" />
      </div>

      <div v-else-if="errorMessage" class="flex flex-col items-center justify-center rounded-lg border border-slate-200 bg-white px-4 py-16 text-center dark:border-slate-800 dark:bg-slate-900">
        <UIcon name="i-lucide-circle-alert" class="size-9 text-red-500" />
        <p class="mt-3 text-sm font-medium text-slate-950 dark:text-white">{{ errorMessage }}</p>
        <UButton class="mt-5" color="neutral" variant="outline" icon="i-lucide-refresh-cw" @click="loadPage">
          {{ $t('common.retry') }}
        </UButton>
      </div>

      <div v-else-if="topic" class="grid gap-6 lg:grid-cols-[minmax(0,1fr)_340px] lg:items-start">
        <main class="space-y-6">
          <section class="rounded-lg border border-slate-200 bg-white p-4 dark:border-slate-800 dark:bg-slate-900">
            <div class="flex flex-wrap items-center gap-2">
              <UBadge v-if="topic.isSticky" color="primary" variant="soft">{{ $t('forum.topicList.badges.sticky') }}</UBadge>
              <UBadge v-if="topic.isLocked" color="neutral" variant="outline">{{ $t('forum.topicList.badges.locked') }}</UBadge>
              <UBadge color="neutral" variant="soft">{{ $t('forum.detail.topicId', { id: topic.id }) }}</UBadge>
            </div>

            <p class="mt-5 whitespace-pre-wrap break-words text-sm leading-6 text-slate-700 dark:text-slate-200">
              {{ topic.content }}
            </p>

            <div v-if="topicAppends.length > 0" class="mt-5 space-y-3 border-t border-slate-200 pt-4 dark:border-slate-800">
              <article v-for="(append, index) in topicAppends" :key="`${append.created_at || append.createdAt || index}`" class="rounded-md bg-slate-50 px-3 py-3 dark:bg-slate-950">
                <div class="flex items-center justify-between gap-3">
                  <p class="text-xs font-medium text-slate-500 dark:text-slate-400">
                    {{ $t('forum.detail.append.index', { index: index + 1 }) }}
                  </p>
                  <p class="text-xs text-slate-500 dark:text-slate-400">{{ formatDateTime(append.created_at || append.createdAt, locale) }}</p>
                </div>
                <p class="mt-2 whitespace-pre-wrap break-words text-sm leading-6 text-slate-700 dark:text-slate-200">{{ append.content }}</p>
              </article>
            </div>
          </section>

          <section class="overflow-hidden rounded-lg border border-slate-200 bg-white dark:border-slate-800 dark:bg-slate-900">
            <div class="flex flex-col gap-3 border-b border-slate-200 px-4 py-3 sm:flex-row sm:items-center sm:justify-between dark:border-slate-800">
              <div>
                <h2 class="text-base font-semibold text-slate-950 dark:text-white">{{ $t('forum.detail.replies.title') }}</h2>
                <p class="mt-1 text-sm text-slate-500 dark:text-slate-400">
                  {{ $t('forum.detail.replies.summary', { count: numberFormatter.format(replyTotal) }) }}
                </p>
              </div>
              <UButton color="neutral" variant="outline" icon="i-lucide-refresh-cw" :loading="repliesPending" @click="loadReplies">
                {{ $t('common.refresh') }}
              </UButton>
            </div>

            <div v-if="repliesError" class="flex flex-col items-center justify-center px-4 py-10 text-center">
              <UIcon name="i-lucide-circle-alert" class="size-8 text-red-500" />
              <p class="mt-3 text-sm font-medium text-slate-950 dark:text-white">{{ repliesError }}</p>
              <UButton class="mt-5" color="neutral" variant="outline" size="sm" icon="i-lucide-refresh-cw" @click="loadReplies">
                {{ $t('common.retry') }}
              </UButton>
            </div>

            <div v-else-if="repliesPending && replies.length === 0" class="space-y-2 px-4 py-4">
              <div v-for="index in 4" :key="index" class="h-24 animate-pulse rounded-md bg-slate-100 dark:bg-slate-800" />
            </div>

            <div v-else-if="replies.length === 0" class="px-4 py-10 text-center text-sm text-slate-500 dark:text-slate-400">
              {{ $t('forum.detail.replies.empty') }}
            </div>

            <div v-else class="divide-y divide-slate-200 dark:divide-slate-800">
              <article v-for="(reply, index) in replies" :key="reply.id" class="px-4 py-4">
                <div class="flex flex-col gap-3 sm:flex-row sm:items-start sm:justify-between">
                  <div class="min-w-0">
                    <p class="text-sm font-medium text-slate-950 dark:text-white">{{ reply.username || `#${reply.userId}` }}</p>
                    <p class="mt-1 text-xs text-slate-500 dark:text-slate-400">
                      {{ $t('forum.detail.replies.floor', { floor: replyFloor(index) }) }}
                      <span class="mx-1 text-slate-300 dark:text-slate-700">/</span>
                      {{ formatDateTime(reply.createdAt, locale) }}
                    </p>
                  </div>
                  <div class="flex shrink-0 flex-wrap items-center gap-1">
                    <UButton
                      color="neutral"
                      :variant="reply.isLiked ? 'soft' : 'ghost'"
                      size="xs"
                      :icon="reply.isLiked ? 'i-lucide-heart' : 'i-lucide-heart-plus'"
                      :loading="replyActionPending === `like:${reply.id}`"
                      @click="handleToggleReplyLike(reply)"
                    >
                      {{ reply.isLiked ? $t('forum.detail.actions.liked') : $t('forum.detail.actions.like') }}
                    </UButton>
                    <UButton color="neutral" variant="ghost" size="xs" icon="i-lucide-coins" @click="openReplyPanel(reply.id, 'reward')">
                      {{ $t('forum.detail.actions.reward') }}
                    </UButton>
                    <UButton color="neutral" variant="ghost" size="xs" icon="i-lucide-flag" @click="openReplyPanel(reply.id, 'report')">
                      {{ $t('forum.detail.actions.report') }}
                    </UButton>
                  </div>
                </div>

                <p class="mt-3 whitespace-pre-wrap break-words text-sm leading-6 text-slate-700 dark:text-slate-200">
                  {{ reply.content }}
                </p>

                <div v-if="activeReplyPanel?.id === reply.id" class="mt-4 rounded-md border border-slate-200 bg-slate-50 p-3 dark:border-slate-800 dark:bg-slate-950">
                  <form v-if="activeReplyPanel.type === 'reward'" class="grid gap-3 sm:grid-cols-[160px_auto]" @submit.prevent="handleRewardReply(reply.id)">
                    <UInput v-model="replyRewardAmount" type="number" min="1" step="1" :placeholder="$t('forum.detail.reward.amount')" :disabled="replyActionPending === `reward:${reply.id}`" />
                    <div class="flex items-center gap-2">
                      <UButton type="submit" color="primary" size="sm" icon="i-lucide-coins" :loading="replyActionPending === `reward:${reply.id}`">
                        {{ $t('forum.detail.reward.submit') }}
                      </UButton>
                      <UButton type="button" color="neutral" variant="ghost" size="sm" @click="closeReplyPanel">
                        {{ $t('common.cancel') }}
                      </UButton>
                    </div>
                  </form>
                  <form v-else class="grid gap-3" @submit.prevent="handleReportReply(reply.id)">
                    <UTextarea v-model="replyReportReason" :rows="3" :placeholder="$t('forum.detail.report.reason')" :disabled="replyActionPending === `report:${reply.id}`" />
                    <div class="flex items-center gap-2">
                      <UButton type="submit" color="primary" size="sm" icon="i-lucide-flag" :loading="replyActionPending === `report:${reply.id}`">
                        {{ $t('forum.detail.report.submit') }}
                      </UButton>
                      <UButton type="button" color="neutral" variant="ghost" size="sm" @click="closeReplyPanel">
                        {{ $t('common.cancel') }}
                      </UButton>
                    </div>
                  </form>
                </div>
              </article>
            </div>

            <AppPager
              class="border-t border-slate-200 px-4 py-3 dark:border-slate-800"
              size="sm"
              :page="replyPage"
              :total="replyTotal"
              :page-size="replySize"
              :disabled="repliesPending"
              @page-change="goToReplyPage"
            />
          </section>

          <section class="rounded-lg border border-slate-200 bg-white p-4 dark:border-slate-800 dark:bg-slate-900">
            <h2 class="text-base font-semibold text-slate-950 dark:text-white">{{ $t('forum.detail.replyForm.title') }}</h2>
            <form class="mt-4 space-y-3" @submit.prevent="handleCreateReply">
              <UTextarea v-model="replyContent" class="w-full" :rows="5" :placeholder="$t('forum.detail.replyForm.placeholder')" :disabled="replyCreatePending || topic.isLocked" />
              <div class="flex justify-end">
                <UButton type="submit" color="primary" icon="i-lucide-send" :loading="replyCreatePending" :disabled="!canCreateReply">
                  {{ $t('forum.detail.replyForm.submit') }}
                </UButton>
              </div>
            </form>
          </section>
        </main>

        <aside class="space-y-6">
          <UCard class="rounded-lg">
            <template #header>
              <h2 class="text-base font-semibold text-slate-950 dark:text-white">{{ $t('forum.detail.info.title') }}</h2>
            </template>

            <dl class="space-y-3 text-sm">
              <div class="flex items-center justify-between gap-3">
                <dt class="text-slate-500 dark:text-slate-400">{{ $t('forum.detail.info.author') }}</dt>
                <dd class="min-w-0 truncate font-medium text-slate-950 dark:text-white">{{ topic.username || `#${topic.userId}` }}</dd>
              </div>
              <div class="flex items-center justify-between gap-3">
                <dt class="text-slate-500 dark:text-slate-400">{{ $t('forum.detail.info.createdAt') }}</dt>
                <dd class="font-medium text-slate-950 dark:text-white">{{ formatDateTime(topic.createdAt, locale) }}</dd>
              </div>
              <div class="flex items-center justify-between gap-3">
                <dt class="text-slate-500 dark:text-slate-400">{{ $t('forum.detail.info.views') }}</dt>
                <dd class="font-medium text-slate-950 dark:text-white">{{ numberFormatter.format(topic.views) }}</dd>
              </div>
              <div class="flex items-center justify-between gap-3">
                <dt class="text-slate-500 dark:text-slate-400">{{ $t('forum.detail.info.replies') }}</dt>
                <dd class="font-medium text-slate-950 dark:text-white">{{ numberFormatter.format(topic.replyCount) }}</dd>
              </div>
            </dl>
          </UCard>

          <UCard v-if="isStaff" class="rounded-lg">
            <template #header>
              <h2 class="text-base font-semibold text-slate-950 dark:text-white">{{ $t('forum.detail.admin.title') }}</h2>
            </template>

            <div class="grid gap-2">
              <UButton
                color="neutral"
                variant="outline"
                :icon="topic.isLocked ? 'i-lucide-lock-open' : 'i-lucide-lock'"
                :loading="adminActionPending === 'lock'"
                @click="handleAdminTopicAction(topic.isLocked ? 'unlock' : 'lock')"
              >
                {{ topic.isLocked ? $t('forum.detail.admin.unlock') : $t('forum.detail.admin.lock') }}
              </UButton>
              <UButton
                color="neutral"
                variant="outline"
                :icon="topic.isSticky ? 'i-lucide-pin-off' : 'i-lucide-pin'"
                :loading="adminActionPending === 'pin'"
                @click="handleAdminTopicAction(topic.isSticky ? 'unpin' : 'pin')"
              >
                {{ topic.isSticky ? $t('forum.detail.admin.unpin') : $t('forum.detail.admin.pin') }}
              </UButton>
            </div>

            <form class="mt-4 grid gap-2 border-t border-slate-200 pt-4 dark:border-slate-800" @submit.prevent="handleMoveTopic">
              <label class="block">
                <span class="text-sm font-medium text-slate-700 dark:text-slate-200">{{ $t('forum.detail.admin.moveTo') }}</span>
                <select v-model.number="moveNodeId" class="mt-1 h-10 w-full rounded-md border border-slate-200 bg-white px-3 text-sm outline-none dark:border-slate-700 dark:bg-slate-950" :disabled="adminNodesPending || adminActionPending === 'move'">
                  <option v-for="node in adminNodes" :key="node.id" :value="node.id">{{ forumNodeName(node) }}</option>
                </select>
              </label>
              <UButton type="submit" color="primary" variant="soft" icon="i-lucide-move-right" :loading="adminActionPending === 'move'" :disabled="!moveNodeId || moveNodeId === topic.nodeId">
                {{ $t('forum.detail.admin.move') }}
              </UButton>
              <p v-if="adminError" class="text-sm text-red-600 dark:text-red-300">{{ adminError }}</p>
            </form>
          </UCard>

          <UCard class="rounded-lg">
            <template #header>
              <h2 class="text-base font-semibold text-slate-950 dark:text-white">{{ $t('forum.detail.actions.title') }}</h2>
            </template>

            <div class="grid grid-cols-1 gap-2">
              <UButton color="neutral" variant="outline" icon="i-lucide-plus-square" :disabled="topic.isLocked" @click="openTopicPanel('append')">
                {{ $t('forum.detail.actions.append') }}
              </UButton>
              <UButton color="neutral" variant="outline" icon="i-lucide-coins" @click="openTopicPanel('reward')">
                {{ $t('forum.detail.actions.reward') }}
              </UButton>
              <UButton color="neutral" variant="outline" icon="i-lucide-flag" @click="openTopicPanel('report')">
                {{ $t('forum.detail.actions.report') }}
              </UButton>
            </div>

            <div v-if="activeTopicPanel" class="mt-4 border-t border-slate-200 pt-4 dark:border-slate-800">
              <form v-if="activeTopicPanel === 'append'" class="space-y-3" @submit.prevent="handleAppendTopic">
                <UTextarea v-model="topicAppendContent" :rows="4" :placeholder="$t('forum.detail.append.placeholder')" :disabled="topicActionPending === 'append'" />
                <div class="flex items-center gap-2">
                  <UButton type="submit" color="primary" size="sm" icon="i-lucide-plus-square" :loading="topicActionPending === 'append'">
                    {{ $t('forum.detail.append.submit') }}
                  </UButton>
                  <UButton type="button" color="neutral" variant="ghost" size="sm" @click="closeTopicPanel">
                    {{ $t('common.cancel') }}
                  </UButton>
                </div>
              </form>

              <form v-else-if="activeTopicPanel === 'reward'" class="space-y-3" @submit.prevent="handleRewardTopic">
                <UInput v-model="topicRewardAmount" type="number" min="1" step="1" :placeholder="$t('forum.detail.reward.amount')" :disabled="topicActionPending === 'reward'" />
                <div class="flex items-center gap-2">
                  <UButton type="submit" color="primary" size="sm" icon="i-lucide-coins" :loading="topicActionPending === 'reward'">
                    {{ $t('forum.detail.reward.submit') }}
                  </UButton>
                  <UButton type="button" color="neutral" variant="ghost" size="sm" @click="closeTopicPanel">
                    {{ $t('common.cancel') }}
                  </UButton>
                </div>
              </form>

              <form v-else class="space-y-3" @submit.prevent="handleReportTopic">
                <UTextarea v-model="topicReportReason" :rows="4" :placeholder="$t('forum.detail.report.reason')" :disabled="topicActionPending === 'report'" />
                <div class="flex items-center gap-2">
                  <UButton type="submit" color="primary" size="sm" icon="i-lucide-flag" :loading="topicActionPending === 'report'">
                    {{ $t('forum.detail.report.submit') }}
                  </UButton>
                  <UButton type="button" color="neutral" variant="ghost" size="sm" @click="closeTopicPanel">
                    {{ $t('common.cancel') }}
                  </UButton>
                </div>
              </form>
            </div>
          </UCard>
        </aside>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ApiError } from '~/composables/useApi'
import type { AdminForumNode } from '~/composables/useAdmin'
import { useForum, type ForumReplyItem, type ForumTopicAppend, type ForumTopicDetail } from '~/composables/useForum'
import { formatDateTime, localizeI18nName } from '~/utils/format'

definePageMeta({
  middleware: 'auth'
})

type TopicPanel = 'append' | 'reward' | 'report'
type ReplyPanelType = 'reward' | 'report'

const { t, locale } = useI18n()
const route = useRoute()
const toast = useToast()
const forum = useForum()
const adminApi = useAdmin()
const { isStaff } = useAuth()

const topicId = computed(() => Number(route.params.id || 0))
const topic = ref<ForumTopicDetail | null>(null)
const replies = ref<ForumReplyItem[]>([])
const replyTotal = ref(0)
const pending = ref(true)
const repliesPending = ref(false)
const errorMessage = ref('')
const repliesError = ref('')
const topicActionPending = ref('')
const replyActionPending = ref('')
const replyCreatePending = ref(false)
const adminActionPending = ref('')
const adminNodesPending = ref(false)
const adminError = ref('')
const adminNodes = ref<AdminForumNode[]>([])
const moveNodeId = ref(0)

const replyPage = ref(readPositiveIntQuery('page', 1))
const replySize = 50
const activeTopicPanel = ref<TopicPanel | null>(null)
const activeReplyPanel = ref<{ id: number, type: ReplyPanelType } | null>(null)
const topicAppendContent = ref('')
const topicRewardAmount = ref('')
const topicReportReason = ref('')
const replyRewardAmount = ref('')
const replyReportReason = ref('')
const replyContent = ref('')

const numberFormatter = computed(() => new Intl.NumberFormat(locale.value))
const replyTotalPages = computed(() => Math.max(1, Math.ceil(replyTotal.value / replySize)))
const topicAppends = computed<ForumTopicAppend[]>(() => {
  return Array.isArray(topic.value?.appends) ? topic.value.appends.filter((append) => append?.content) : []
})
const canCreateReply = computed(() => {
  return Boolean(topic.value && !topic.value.isLocked && replyContent.value.trim().length >= 2 && !replyCreatePending.value)
})

useHead(() => ({
  title: `${topic.value?.subject || t('forum.detail.metaTitle')} - NextPT`
}))

onMounted(() => {
  loadPage()
  if (isStaff.value) loadAdminNodes()
})

function readFirstQueryValue(key: string) {
  const value = route.query[key]
  return Array.isArray(value) ? value[0] : value
}

function readPositiveIntQuery(key: string, fallback: number) {
  const parsed = Number(readFirstQueryValue(key))
  return Number.isInteger(parsed) && parsed > 0 ? parsed : fallback
}

async function loadPage() {
  if (!Number.isInteger(topicId.value) || topicId.value <= 0) {
    pending.value = false
    errorMessage.value = t('forum.detail.invalidId')
    return
  }

  pending.value = true
  errorMessage.value = ''
  try {
    await Promise.all([loadTopic(), loadReplies()])
  } catch (error) {
    errorMessage.value = error instanceof ApiError ? error.message : t('common.requestFailed')
  } finally {
    pending.value = false
  }
}

async function loadTopic() {
  topic.value = await forum.getTopic(topicId.value)
  moveNodeId.value = topic.value.nodeId
}

async function loadReplies() {
  repliesPending.value = true
  repliesError.value = ''
  try {
    const data = await forum.listReplies(topicId.value, replyPage.value, replySize)
    replies.value = data.list || []
    replyTotal.value = data.total || 0

    if (replyPage.value > replyTotalPages.value) {
      replyPage.value = replyTotalPages.value
      await loadReplies()
    }
  } catch (error) {
    replies.value = []
    replyTotal.value = 0
    repliesError.value = error instanceof ApiError ? error.message : t('common.requestFailed')
    if (pending.value) throw error
  } finally {
    repliesPending.value = false
  }
}

function goToReplyPage(nextPage: number) {
  replyPage.value = Math.min(Math.max(1, nextPage), replyTotalPages.value)
  loadReplies()
}

function replyFloor(index: number) {
  return (replyPage.value - 1) * replySize + index + 1
}

async function handleToggleTopicLike() {
  if (!topic.value) return
  topicActionPending.value = 'like'
  try {
    await forum.toggleTopicLike(topic.value.id)
    topic.value.isLiked = !topic.value.isLiked
  } catch (error) {
    showErrorToast(error)
  } finally {
    topicActionPending.value = ''
  }
}

async function handleToggleTopicBookmark() {
  if (!topic.value) return
  topicActionPending.value = 'bookmark'
  try {
    if (topic.value.isBookmarked) {
      await forum.unbookmarkTopic(topic.value.id)
    } else {
      await forum.bookmarkTopic(topic.value.id)
    }
    topic.value.isBookmarked = !topic.value.isBookmarked
  } catch (error) {
    showErrorToast(error)
  } finally {
    topicActionPending.value = ''
  }
}

function openTopicPanel(panel: TopicPanel) {
  activeTopicPanel.value = activeTopicPanel.value === panel ? null : panel
}

function closeTopicPanel() {
  activeTopicPanel.value = null
  topicAppendContent.value = ''
  topicRewardAmount.value = ''
  topicReportReason.value = ''
}

async function handleAppendTopic() {
  if (!topic.value || !topicAppendContent.value.trim()) return

  topicActionPending.value = 'append'
  try {
    await forum.appendTopic(topic.value.id, topicAppendContent.value.trim())
    toast.add({ title: t('forum.detail.append.success'), color: 'success', icon: 'i-lucide-check-circle' })
    closeTopicPanel()
    await loadTopic()
  } catch (error) {
    showErrorToast(error)
  } finally {
    topicActionPending.value = ''
  }
}

async function handleRewardTopic() {
  if (!topic.value) return
  const amount = Number(topicRewardAmount.value)
  if (!Number.isFinite(amount) || amount <= 0) return

  topicActionPending.value = 'reward'
  try {
    await forum.rewardTopic(topic.value.id, amount)
    toast.add({ title: t('forum.detail.reward.success'), color: 'success', icon: 'i-lucide-check-circle' })
    closeTopicPanel()
  } catch (error) {
    showErrorToast(error)
  } finally {
    topicActionPending.value = ''
  }
}

async function handleReportTopic() {
  if (!topic.value || !topicReportReason.value.trim()) return

  topicActionPending.value = 'report'
  try {
    await forum.reportTopic(topic.value.id, topicReportReason.value.trim())
    toast.add({ title: t('forum.detail.report.success'), color: 'success', icon: 'i-lucide-check-circle' })
    closeTopicPanel()
  } catch (error) {
    showErrorToast(error)
  } finally {
    topicActionPending.value = ''
  }
}

async function handleCreateReply() {
  if (!topic.value || !canCreateReply.value) return

  replyCreatePending.value = true
  try {
    await forum.createReply(topic.value.id, replyContent.value.trim())
    toast.add({ title: t('forum.detail.replyForm.success'), color: 'success', icon: 'i-lucide-check-circle' })
    replyContent.value = ''
    await Promise.all([loadTopic(), loadReplies()])
  } catch (error) {
    showErrorToast(error)
  } finally {
    replyCreatePending.value = false
  }
}

async function handleToggleReplyLike(reply: ForumReplyItem) {
  replyActionPending.value = `like:${reply.id}`
  try {
    await forum.toggleReplyLike(reply.id)
    reply.isLiked = !reply.isLiked
  } catch (error) {
    showErrorToast(error)
  } finally {
    replyActionPending.value = ''
  }
}

function openReplyPanel(replyId: number, type: ReplyPanelType) {
  const current = activeReplyPanel.value
  activeReplyPanel.value = current?.id === replyId && current.type === type ? null : { id: replyId, type }
  replyRewardAmount.value = ''
  replyReportReason.value = ''
}

function closeReplyPanel() {
  activeReplyPanel.value = null
  replyRewardAmount.value = ''
  replyReportReason.value = ''
}

async function handleRewardReply(replyId: number) {
  const amount = Number(replyRewardAmount.value)
  if (!Number.isFinite(amount) || amount <= 0) return

  replyActionPending.value = `reward:${replyId}`
  try {
    await forum.rewardReply(replyId, amount)
    toast.add({ title: t('forum.detail.reward.success'), color: 'success', icon: 'i-lucide-check-circle' })
    closeReplyPanel()
  } catch (error) {
    showErrorToast(error)
  } finally {
    replyActionPending.value = ''
  }
}

async function handleReportReply(replyId: number) {
  if (!replyReportReason.value.trim()) return

  replyActionPending.value = `report:${replyId}`
  try {
    await forum.reportReply(replyId, replyReportReason.value.trim())
    toast.add({ title: t('forum.detail.report.success'), color: 'success', icon: 'i-lucide-check-circle' })
    closeReplyPanel()
  } catch (error) {
    showErrorToast(error)
  } finally {
    replyActionPending.value = ''
  }
}

function showErrorToast(error: unknown) {
  toast.add({
    title: error instanceof ApiError ? error.message : t('common.requestFailed'),
    color: 'error',
    icon: 'i-lucide-circle-alert'
  })
}

async function loadAdminNodes() {
  adminNodesPending.value = true
  try {
    const data = await adminApi.listForumNodes()
    adminNodes.value = data.nodes || []
  } catch {
    adminNodes.value = []
  } finally {
    adminNodesPending.value = false
  }
}

async function handleAdminTopicAction(action: 'lock' | 'unlock' | 'pin' | 'unpin') {
  if (!topic.value) return
  adminActionPending.value = action === 'lock' || action === 'unlock' ? 'lock' : 'pin'
  adminError.value = ''
  try {
    if (action === 'lock') await adminApi.lockForumTopic(topic.value.id)
    if (action === 'unlock') await adminApi.unlockForumTopic(topic.value.id)
    if (action === 'pin') await adminApi.pinForumTopic(topic.value.id)
    if (action === 'unpin') await adminApi.unpinForumTopic(topic.value.id)
    toast.add({ title: t('forum.detail.admin.saved'), color: 'success', icon: 'i-lucide-check-circle' })
    await loadTopic()
  } catch (error) {
    adminError.value = error instanceof ApiError ? error.message : t('common.requestFailed')
  } finally {
    adminActionPending.value = ''
  }
}

async function handleMoveTopic() {
  if (!topic.value || !moveNodeId.value || moveNodeId.value === topic.value.nodeId) return
  adminActionPending.value = 'move'
  adminError.value = ''
  try {
    await adminApi.moveForumTopic(topic.value.id, moveNodeId.value)
    toast.add({ title: t('forum.detail.admin.saved'), color: 'success', icon: 'i-lucide-check-circle' })
    await loadTopic()
  } catch (error) {
    adminError.value = error instanceof ApiError ? error.message : t('common.requestFailed')
  } finally {
    adminActionPending.value = ''
  }
}

function forumNodeName(node: AdminForumNode) {
  return localizeI18nName(node.nameI18N as any, locale.value, node.slug || `#${node.id}`)
}
</script>
