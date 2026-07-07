<template>
  <div class="min-h-[calc(100vh-4rem)] bg-slate-50 py-6 dark:bg-slate-950">
    <div class="w-full px-3 sm:px-4 lg:px-5">
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
      </div>

      <div v-else-if="topic && topicEditOpen" class="mx-auto grid w-full max-w-6xl grid-cols-1 gap-4 xl:grid-cols-[minmax(0,1fr)_300px] xl:items-start">
        <main class="min-w-0">
          <ForumTopicEditForm
            v-model:selected-category-id="editSelectedCategoryId"
            v-model:node-id="editForm.nodeId"
            v-model:subject="editForm.subject"
            v-model:content="editForm.content"
            v-model:editor-mode="topicEditEditorMode"
            :categories="editCategories"
            :nodes-pending="editNodesPending"
            :nodes-error="editNodesError"
            :action-pending="topicActionPending"
            :can-submit="canSubmitTopicEdit"
            :is-dirty="isTopicEditDirty"
            :submit-disabled-reason="topicEditSubmitDisabledReason"
            @submit="handleUpdateTopic"
            @cancel="closeTopicEdit"
          />
        </main>

        <aside class="app-sticky-offset space-y-3 xl:sticky">
          <section class="rounded-lg border border-slate-200 bg-white p-3 dark:border-slate-800 dark:bg-slate-900">
            <h2 class="text-sm font-semibold text-slate-950 dark:text-white">{{ $t('forum.detail.edit.guide.title') }}</h2>
            <ul class="mt-3 space-y-2 text-sm leading-6 text-slate-600 dark:text-slate-300">
              <li class="flex gap-2">
                <span class="mt-2 size-1.5 shrink-0 rounded-full bg-slate-300 dark:bg-slate-600" />
                <span>{{ $t('forum.detail.edit.guide.window') }}</span>
              </li>
              <li class="flex gap-2">
                <span class="mt-2 size-1.5 shrink-0 rounded-full bg-slate-300 dark:bg-slate-600" />
                <span>{{ $t('forum.detail.edit.guide.markdown') }}</span>
              </li>
              <li class="flex gap-2">
                <span class="mt-2 size-1.5 shrink-0 rounded-full bg-slate-300 dark:bg-slate-600" />
                <span>{{ $t('forum.detail.edit.guide.node') }}</span>
              </li>
              <li class="flex gap-2 rounded-md bg-amber-50 px-2 py-1.5 text-amber-800 dark:bg-amber-950/40 dark:text-amber-200">
                <UIcon name="i-lucide-triangle-alert" class="mt-1 size-4 shrink-0" />
                <span>{{ $t('forum.detail.edit.guide.discard') }}</span>
              </li>
            </ul>
          </section>

          <section class="rounded-lg border border-slate-200 bg-white p-3 dark:border-slate-800 dark:bg-slate-900">
            <dl class="space-y-2 text-sm">
              <div class="flex items-center justify-between gap-3">
                <dt class="text-slate-500 dark:text-slate-400">{{ $t('forum.detail.edit.windowStatus') }}</dt>
                <dd>
                  <UBadge :color="topicEditStatusColor" variant="soft">
                    {{ topicEditStatusText }}
                  </UBadge>
                </dd>
              </div>
              <div class="flex items-center justify-between gap-3">
                <dt class="text-slate-500 dark:text-slate-400">{{ $t('forum.create.summary.subjectLength') }}</dt>
                <dd class="font-medium text-slate-950 dark:text-white">{{ numberFormatter.format(editForm.subject.trim().length) }}</dd>
              </div>
              <div class="flex items-center justify-between gap-3">
                <dt class="text-slate-500 dark:text-slate-400">{{ $t('forum.create.summary.content') }}</dt>
                <dd class="font-medium text-slate-950 dark:text-white">{{ numberFormatter.format(editForm.content.trim().length) }}</dd>
              </div>
            </dl>
          </section>
        </aside>
      </div>

      <div v-else-if="topic" class="grid gap-6 lg:grid-cols-[minmax(0,1fr)_340px] lg:items-start">
        <main class="space-y-6">
          <ForumTopicCard
            v-model:append-content="topicAppendContent"
            v-model:append-editor-mode="topicAppendEditorMode"
            v-model:report-reason="topicReportReason"
            :topic="topic"
            :appends="topicAppends"
            :active-panel="activeTopicPanel"
            :action-pending="topicActionPending"
            :can-edit="canEditTopic"
            :can-append="canAppendTopic"
            :can-use-owner-actions="canUseTopicOwnerActions"
            @toggle-like="handleToggleTopicLike"
            @toggle-bookmark="handleToggleTopicBookmark"
            @edit="openTopicEdit"
            @open-panel="openTopicPanel"
            @close-panel="closeTopicPanel"
            @append="handleAppendTopic"
            @report="handleReportTopic"
          />

          <section class="overflow-hidden rounded-lg border border-slate-200 bg-white dark:border-slate-800 dark:bg-slate-900">
            <div class="flex items-center justify-between gap-3 border-b border-slate-200 px-4 py-3 dark:border-slate-800">
              <h2 class="text-base font-semibold text-slate-950 dark:text-white">{{ $t('forum.detail.replies.title') }}</h2>
              <UBadge color="neutral" variant="soft">
                {{ $t('forum.detail.replies.summary', { count: numberFormatter.format(replyTotal) }) }}
              </UBadge>
            </div>

            <RichTextComposer
              id="forum-reply-composer"
              v-model="replyContent"
              v-model:mode="replyEditorMode"
              class="border-b border-slate-200 bg-slate-50/60 p-4 dark:border-slate-800 dark:bg-slate-950/40"
              :placeholder="$t('forum.detail.replyForm.placeholder')"
              :disabled="replyCreatePending || topic.isLocked || !canCreateForumReply"
              :pending="replyCreatePending"
              :submit-disabled="!canCreateReply"
              :submit-disabled-text="canCreateForumReply ? '' : $t('common.noPermission')"
              :submit-label="$t('forum.detail.replyForm.submit')"
              :write-label="$t('common.editor.edit')"
              :preview-label="$t('common.editor.preview')"
              :preview-empty="$t('common.editor.previewEmpty')"
              :locked-text="!canCreateForumReply ? $t('common.noPermission') : topic.isLocked ? $t('forum.detail.replyForm.locked') : ''"
              @submit="handleCreateReply"
            />

            <ForumReplyList
              v-model:active-report-id="activeReplyReportId"
              v-model:report-reason="replyReportReason"
              :replies="replies"
              :total="replyTotal"
              :page="replyPage"
              :page-size="replySize"
              :pending="repliesPending"
              :error="repliesError"
              :action-pending="replyActionPending"
              :submit-reward="submitReplyReward"
              @page-change="goToReplyPage"
              @toggle-like="handleToggleReplyLike"
              @reward-success="handleReplyRewardSuccess"
              @quote="insertReplyQuote"
              @report="handleReportReply"
            />
          </section>
        </main>

        <ForumTopicSidebar
          v-model:move-node-id="moveNodeId"
          :topic="topic"
          :can-manage-topic="canManageForumTopic"
          :admin-nodes="adminNodes"
          :admin-categories="adminCategories"
          :admin-nodes-pending="adminNodesPending"
          :admin-action-pending="adminActionPending"
          :admin-error="adminError"
          :load-rewards="loadTopicRewards"
          :submit-reward="submitTopicReward"
          @admin-action="handleAdminTopicAction"
          @move="handleMoveTopic"
          @delete="handleAdminDeleteTopic"
        />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ApiError } from '~/composables/useApi'
import type { AdminForumCategory, AdminForumNode } from '~/composables/useAdmin'
import { useForum, type ForumNodeCategory, type ForumReplyItem, type ForumTopicAppend, type ForumTopicDetail } from '~/composables/useForum'

definePageMeta({
  middleware: 'auth'
})

type TopicPanel = 'append' | 'report'
type EditorMode = 'write' | 'preview'

interface TopicEditSnapshot {
  nodeId: string
  subject: string
  content: string
}

const { t, locale } = useI18n()
const route = useRoute()
const router = useRouter()
const localePath = useLocalePath()
const toast = useToast()
const forum = useForum()
const adminApi = useAdmin()
const workspaceTabs = useWorkspaceTabs('app')
const { hasPermission, user } = useAuth()

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
const adminCategories = ref<AdminForumCategory[]>([])
const moveNodeId = ref(0)
const topicEditOpen = ref(false)
const currentTimeMs = ref(Date.now())
const editCategories = ref<ForumNodeCategory[]>([])
const editNodesPending = ref(false)
const editNodesError = ref('')
const editSelectedCategoryId = ref(0)
const editForm = reactive({
  nodeId: '0',
  subject: '',
  content: ''
})
const editInitialForm = ref<TopicEditSnapshot | null>(null)

const replyPage = ref(readPositiveIntQuery('page', 1))
const replySize = 50
const activeTopicPanel = ref<TopicPanel | null>(null)
const activeReplyReportId = ref(0)
const topicAppendContent = ref('')
const topicEditEditorMode = ref<EditorMode>('write')
const topicAppendEditorMode = ref<EditorMode>('write')
const topicReportReason = ref('')
const replyReportReason = ref('')
const replyContent = ref('')
const replyEditorMode = ref<EditorMode>('write')
let currentTimeTimer: ReturnType<typeof setInterval> | null = null

const numberFormatter = computed(() => new Intl.NumberFormat(locale.value))
const replyTotalPages = computed(() => Math.max(1, Math.ceil(replyTotal.value / replySize)))
const topicAppends = computed<ForumTopicAppend[]>(() => {
  return Array.isArray(topic.value?.appends) ? topic.value.appends.filter((append) => append?.content) : []
})
watch(
  () => topic.value?.subject,
  (subject) => {
    if (subject) {
      workspaceTabs.updateTabTitle(route.fullPath, subject)
    }
  }
)
const topicEditWindowMs = 5 * 60 * 1000
const topicEditExpiresAtMs = computed(() => {
  const createdAt = Date.parse(topic.value?.createdAt || '')
  return Number.isFinite(createdAt) ? createdAt + topicEditWindowMs : 0
})
const topicEditRemainingMs = computed(() => {
  if (!topicEditExpiresAtMs.value) return 0
  return Math.max(0, topicEditExpiresAtMs.value - currentTimeMs.value)
})
const topicEditRemainingSeconds = computed(() => Math.ceil(topicEditRemainingMs.value / 1000))
const isTopicEditExpired = computed(() => Boolean(topicEditOpen.value && topicEditRemainingMs.value <= 0))
const topicEditStatusColor = computed<'warning' | 'error'>(() => isTopicEditExpired.value || topicEditRemainingSeconds.value <= 60 ? 'error' : 'warning')
const topicEditStatusText = computed(() => {
  if (isTopicEditExpired.value) return t('forum.detail.edit.expired')
  return t('forum.detail.edit.remaining', { time: formatTopicEditRemaining(topicEditRemainingSeconds.value) })
})
const canEditTopic = computed(() => {
  if (!topic.value || topic.value.isLocked || !user.value?.user.id || topic.value.author?.id !== user.value.user.id) return false
  const createdAt = Date.parse(topic.value.createdAt || '')
  return Number.isFinite(createdAt) && currentTimeMs.value - createdAt <= topicEditWindowMs
})
const canAppendTopic = computed(() => {
  return Boolean(topic.value && !topic.value.isLocked && user.value?.user.id && topic.value.author?.id === user.value.user.id && topicAppends.value.length < 3)
})
const canUseTopicOwnerActions = computed(() => canEditTopic.value || canAppendTopic.value || activeTopicPanel.value === 'append')
const isTopicEditDirty = computed(() => {
  const initial = editInitialForm.value
  if (!initial) return false
  return editForm.nodeId !== initial.nodeId
    || editForm.subject.trim() !== initial.subject
    || editForm.content.trim() !== initial.content
})
const canSubmitTopicEdit = computed(() => {
  return Boolean(
    topic.value
    && isTopicEditDirty.value
    && !isTopicEditExpired.value
    && Number(editForm.nodeId) > 0
    && editForm.subject.trim().length >= 2
    && editForm.content.trim().length >= 2
    && topicActionPending.value !== 'update'
  )
})
const topicEditSubmitDisabledReason = computed(() => {
  if (topicActionPending.value === 'update') return ''
  if (isTopicEditExpired.value) return t('forum.detail.edit.errors.expired')
  if (Number(editForm.nodeId) <= 0) return t('forum.detail.edit.errors.nodeRequired')
  if (editForm.subject.trim().length < 2) return t('forum.detail.edit.errors.subjectTooShort')
  if (editForm.content.trim().length < 2) return t('forum.detail.edit.errors.contentTooShort')
  if (!isTopicEditDirty.value) return t('forum.detail.edit.errors.noChanges')
  return ''
})
const canManageForumTopic = computed(() => hasPermission(Permission.AdminForumTopicManage))
const canCreateForumReply = computed(() => hasPermission(Permission.ForumReplyCreate))
const canCreateReply = computed(() => {
  return Boolean(canCreateForumReply.value && topic.value && !topic.value.isLocked && replyContent.value.trim().length >= 2 && !replyCreatePending.value)
})

useHead(() => ({
  title: topic.value?.subject || t('forum.detail.metaTitle')
}))

onMounted(() => {
  window.addEventListener('beforeunload', handleTopicEditBeforeUnload)
  currentTimeTimer = setInterval(() => {
    currentTimeMs.value = Date.now()
  }, 1000)
  loadPage()
  if (canManageForumTopic.value) loadAdminNodes()
})

onBeforeUnmount(() => {
  if (currentTimeTimer) {
    clearInterval(currentTimeTimer)
    currentTimeTimer = null
  }
  window.removeEventListener('beforeunload', handleTopicEditBeforeUnload)
})

onBeforeRouteLeave(() => {
  if (!shouldWarnTopicEditLeave()) return true
  return window.confirm(t('forum.detail.edit.leaveConfirm'))
})

watch(() => route.hash, () => {
  void nextTick(scrollToReplyHash)
})

watch(canManageForumTopic, (canManage) => {
  if (canManage && adminNodes.value.length === 0 && !adminNodesPending.value) {
    loadAdminNodes()
  }
})

function readFirstQueryValue(key: string) {
  const value = route.query[key]
  return Array.isArray(value) ? value[0] : value
}

function readPositiveIntQuery(key: string, fallback: number) {
  const parsed = Number(readFirstQueryValue(key))
  return Number.isInteger(parsed) && parsed > 0 ? parsed : fallback
}

function formatTopicEditRemaining(seconds: number) {
  const safeSeconds = Math.max(0, seconds)
  const minutes = Math.floor(safeSeconds / 60)
  const restSeconds = safeSeconds % 60
  return `${minutes}:${String(restSeconds).padStart(2, '0')}`
}

function shouldWarnTopicEditLeave() {
  return Boolean(topicEditOpen.value && isTopicEditDirty.value)
}

function handleTopicEditBeforeUnload(event: BeforeUnloadEvent) {
  if (!shouldWarnTopicEditLeave()) return
  event.preventDefault()
  event.returnValue = ''
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

async function openTopicEdit() {
  if (!topic.value || !canEditTopic.value) return

  closeTopicPanel()
  const snapshot: TopicEditSnapshot = {
    nodeId: String(topic.value.nodeId),
    subject: topic.value.subject.trim(),
    content: topic.value.content.trim()
  }
  editForm.nodeId = snapshot.nodeId
  editForm.subject = snapshot.subject
  editForm.content = snapshot.content
  editInitialForm.value = snapshot
  topicEditEditorMode.value = 'write'
  topicEditOpen.value = true
  await loadEditNodes()
}

function closeTopicEdit() {
  topicEditOpen.value = false
  editNodesError.value = ''
  editInitialForm.value = null
  topicEditEditorMode.value = 'write'
}

async function loadEditNodes() {
  editNodesPending.value = true
  editNodesError.value = ''
  try {
    const data = await forum.listNodes({ scope: 'create' })
    editCategories.value = (data.list || []).filter((category) => category.nodes.length > 0)
    resolveEditNodeSelection()
  } catch (error) {
    editCategories.value = []
    editSelectedCategoryId.value = 0
    editForm.nodeId = '0'
    editNodesError.value = error instanceof ApiError ? error.message : t('common.requestFailed')
  } finally {
    editNodesPending.value = false
  }
}

function resolveEditNodeSelection() {
  const flatItems = editCategories.value.flatMap((category) => category.nodes.map((node) => ({ category, node })))
  const matched = flatItems.find((item) => item.node.id === Number(editForm.nodeId))
  if (matched) {
    editSelectedCategoryId.value = matched.category.id
    return
  }

  const firstCategory = editCategories.value[0]
  editSelectedCategoryId.value = firstCategory?.id || 0
  editForm.nodeId = String(firstCategory?.nodes[0]?.id || 0)
}

async function handleUpdateTopic() {
  if (!topic.value || !canSubmitTopicEdit.value) return

  topicActionPending.value = 'update'
  try {
    await forum.updateTopic(topic.value.id, {
      nodeId: Number(editForm.nodeId),
      subject: editForm.subject.trim(),
      content: editForm.content.trim()
    })
    toast.add({ title: t('forum.detail.edit.success'), color: 'success', icon: 'i-lucide-check-circle' })
    closeTopicEdit()
    await loadTopic()
  } catch (error) {
    showErrorToast(error)
  } finally {
    topicActionPending.value = ''
  }
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
      return
    }

    await nextTick(scrollToReplyHash)
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
  void router.replace({
    path: route.path,
    query: nextReplyQuery(replyPage.value)
  })
  loadReplies()
}

function nextReplyQuery(page: number) {
  return {
    ...route.query,
    page: String(page)
  }
}

function scrollToReplyHash() {
  const id = route.hash.replace(/^#/, '')
  if (!id.startsWith('reply-')) return
  document.getElementById(id)?.scrollIntoView({ behavior: 'smooth', block: 'start' })
}

async function handleToggleTopicLike() {
  if (!topic.value) return
  topicActionPending.value = 'like'
  try {
    const nextLiked = !topic.value.isLiked
    await forum.toggleTopicLike(topic.value.id)
    topic.value.isLiked = nextLiked
    topic.value.likeCount = Math.max(0, Number(topic.value.likeCount || 0) + (nextLiked ? 1 : -1))
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
  if (topicEditOpen.value) closeTopicEdit()
  activeTopicPanel.value = activeTopicPanel.value === panel ? null : panel
  if (activeTopicPanel.value === 'append') {
    topicAppendEditorMode.value = 'write'
  }
}

function closeTopicPanel() {
  activeTopicPanel.value = null
  topicAppendContent.value = ''
  topicAppendEditorMode.value = 'write'
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

async function submitTopicReward(amount: number) {
  if (!topic.value) return
  await forum.rewardTopic(topic.value.id, amount)
}

async function loadTopicRewards(page: number, size: number) {
  if (!topic.value) return { list: [], total: 0 }
  return await forum.listTopicRewards(topic.value.id, page, size)
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
    replyEditorMode.value = 'write'
    replyPage.value = Math.max(1, Math.ceil((replyTotal.value + 1) / replySize))
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
    const nextLiked = !reply.isLiked
    await forum.toggleReplyLike(reply.id)
    reply.isLiked = nextLiked
    reply.likeCount = Math.max(0, Number(reply.likeCount || 0) + (nextLiked ? 1 : -1))
  } catch (error) {
    showErrorToast(error)
  } finally {
    replyActionPending.value = ''
  }
}

async function submitReplyReward(reply: ForumReplyItem, amount: number) {
  await forum.rewardReply(reply.id, amount)
}

function handleReplyRewardSuccess(reply: ForumReplyItem) {
  reply.rewardCount = Number(reply.rewardCount || 0) + 1
}

async function handleReportReply(replyId: number) {
  const reason = replyReportReason.value.trim()
  if (reason.length < 5) return

  replyActionPending.value = `report:${replyId}`
  try {
    await forum.reportReply(replyId, reason)
    toast.add({ title: t('forum.detail.report.success'), color: 'success', icon: 'i-lucide-check-circle' })
    activeReplyReportId.value = 0
    replyReportReason.value = ''
  } catch (error) {
    showErrorToast(error)
  } finally {
    replyActionPending.value = ''
  }
}

function insertReplyQuote(quote: string) {
  if (!canCreateForumReply.value) return
  replyContent.value = replyContent.value.trim() ? `${replyContent.value.trim()}\n\n${quote}` : quote
  replyEditorMode.value = 'write'
  void nextTick(() => {
    document.getElementById('forum-reply-composer')?.scrollIntoView({ behavior: 'smooth', block: 'start' })
  })
}

function showErrorToast(error: unknown) {
  toast.add({
    title: error instanceof ApiError ? error.message : t('common.requestFailed'),
    color: 'error',
    icon: 'i-lucide-circle-alert'
  })
}

async function loadAdminNodes() {
  if (!canManageForumTopic.value) return
  adminNodesPending.value = true
  try {
    const [nodesData, categoriesData] = await Promise.all([
      adminApi.listForumNodes(),
      adminApi.listForumCategories()
    ])
    adminNodes.value = nodesData.nodes || []
    adminCategories.value = categoriesData.categories || []
  } catch {
    adminNodes.value = []
    adminCategories.value = []
  } finally {
    adminNodesPending.value = false
  }
}

async function handleAdminTopicAction(action: 'lock' | 'unlock' | 'pin' | 'unpin') {
  if (!canManageForumTopic.value || !topic.value) return
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
  if (!canManageForumTopic.value || !topic.value || !moveNodeId.value || moveNodeId.value === topic.value.nodeId) return
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

async function handleAdminDeleteTopic() {
  if (!canManageForumTopic.value || !topic.value) return
  adminActionPending.value = 'delete'
  adminError.value = ''
  try {
    await adminApi.deleteForumTopic(topic.value.id)
    toast.add({ title: t('forum.detail.admin.deleted'), color: 'success', icon: 'i-lucide-check-circle' })
    await navigateTo(localePath('/forum'))
  } catch (error) {
    adminError.value = error instanceof ApiError ? error.message : t('common.requestFailed')
  } finally {
    adminActionPending.value = ''
  }
}

</script>
