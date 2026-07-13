<template>
  <div class="min-h-[calc(100vh-4rem)] bg-slate-50 py-6 dark:bg-slate-950">
    <div class="w-full px-3 sm:px-4 lg:px-5">
      <div v-if="pending && !request" class="grid gap-4 lg:grid-cols-[minmax(0,1fr)_300px]">
        <div class="h-80 animate-pulse rounded-lg bg-white dark:bg-slate-900" />
        <div class="h-64 animate-pulse rounded-lg bg-white dark:bg-slate-900" />
      </div>

      <section v-else-if="errorMessage && !request" class="rounded-lg border border-red-200 bg-red-50 p-5 text-sm text-red-700 dark:border-red-900 dark:bg-red-950 dark:text-red-200">
        <p>{{ errorMessage }}</p>
        <UButton class="mt-3" color="error" variant="soft" icon="i-lucide-rotate-cw" @click="loadRequest">{{ $t('common.retry') }}</UButton>
      </section>

      <div v-else-if="request" class="grid gap-4 lg:grid-cols-[minmax(0,1fr)_300px] lg:items-start">
        <main class="min-w-0 space-y-3">
          <section class="overflow-hidden rounded-lg border border-slate-200 bg-white dark:border-slate-800 dark:bg-slate-900">
            <div class="p-4">
              <div class="flex flex-wrap items-center gap-2">
                <UBadge :color="request.requestType === CatalogRequestType.Reseed ? 'warning' : 'primary'" variant="soft">{{ requestTypeLabel }}</UBadge>
                <UBadge :color="requestStatusColor(request.status)" variant="soft">{{ requestStatusLabel(request.status) }}</UBadge>
                <UBadge v-if="categoryName" color="neutral" variant="soft">{{ categoryName }}</UBadge>
              </div>
              <h1 class="mt-3 break-words text-xl font-semibold leading-8 text-slate-950 dark:text-white">{{ request.title }}</h1>
              <div class="mt-2 flex flex-wrap items-center gap-x-2 text-xs text-slate-500 dark:text-slate-400">
                <span>{{ request.requester.username || `#${request.requester.id}` }}</span>
                <span>/</span>
                <UTooltip :text="formatDateTime(request.createdAt, locale)" :delay-duration="600"><span>{{ formatRelativeDateTime(request.createdAt, locale) }}</span></UTooltip>
              </div>
            </div>
            <div class="border-t border-slate-200 px-4 py-4 dark:border-slate-800">
              <h2 class="text-sm font-semibold text-slate-950 dark:text-white">{{ $t('catalog.requests.fields.description') }}</h2>
              <div v-if="renderedDescription" class="rich-text mt-3" v-html="renderedDescription" />
            </div>
          </section>

          <section v-if="hasRelatedTorrent" class="overflow-hidden rounded-lg border border-slate-200 bg-white dark:border-slate-800 dark:bg-slate-900">
            <div class="border-b border-slate-200 px-4 py-3 dark:border-slate-800">
              <h2 class="text-base font-semibold text-slate-950 dark:text-white">{{ $t('catalog.requests.detail.relatedTorrent') }}</h2>
            </div>
            <div class="grid gap-2 p-4" :class="hasBothRelatedTorrents ? 'xl:grid-cols-2' : 'grid-cols-1'">
              <NuxtLink v-if="request.targetTorrent?.exist" :to="localePath(`/catalog/torrents/${request.targetTorrent.id}`)" class="group flex min-w-0 items-center gap-3 rounded-md border border-slate-200 px-3 py-2.5 transition hover:border-slate-300 hover:bg-slate-50 dark:border-slate-800 dark:hover:border-slate-700 dark:hover:bg-slate-950">
                <span class="flex size-9 shrink-0 items-center justify-center rounded-md bg-slate-100 text-slate-500 dark:bg-slate-800 dark:text-slate-300">
                  <UIcon name="i-lucide-library" class="size-4" />
                </span>
                <span class="min-w-0 flex-1">
                  <span class="block text-xs text-slate-500 dark:text-slate-400">{{ $t('catalog.requests.fields.targetTorrent') }}</span>
                  <span class="mt-0.5 block truncate text-sm font-medium text-slate-950 dark:text-white">{{ request.targetTorrent.name }}</span>
                </span>
                <UIcon name="i-lucide-arrow-up-right" class="size-4 shrink-0 text-slate-400 transition group-hover:text-sky-600 dark:group-hover:text-sky-300" />
              </NuxtLink>
              <NuxtLink v-if="request.resultTorrent?.exist" :to="localePath(`/catalog/torrents/${request.resultTorrent.id}`)" class="group flex min-w-0 items-center gap-3 rounded-md border border-emerald-200 bg-emerald-50/40 px-3 py-2.5 transition hover:border-emerald-300 hover:bg-emerald-50 dark:border-emerald-900 dark:bg-emerald-950/20 dark:hover:border-emerald-800 dark:hover:bg-emerald-950/40">
                <span class="flex size-9 shrink-0 items-center justify-center rounded-md bg-emerald-100 text-emerald-600 dark:bg-emerald-950 dark:text-emerald-300">
                  <UIcon name="i-lucide-circle-check" class="size-4" />
                </span>
                <span class="min-w-0 flex-1">
                  <span class="block text-xs text-slate-500 dark:text-slate-400">{{ $t('catalog.requests.fields.resultTorrent') }}</span>
                  <span class="mt-0.5 block truncate text-sm font-medium text-slate-950 dark:text-white">{{ request.resultTorrent.name }}</span>
                </span>
                <UIcon name="i-lucide-arrow-up-right" class="size-4 shrink-0 text-slate-400 transition group-hover:text-emerald-600 dark:group-hover:text-emerald-300" />
              </NuxtLink>
            </div>
          </section>
        </main>

        <aside class="app-sticky-offset space-y-3 lg:sticky lg:col-start-2 lg:row-span-2 lg:row-start-1">
          <section class="overflow-hidden rounded-lg border border-slate-200 bg-white dark:border-slate-800 dark:bg-slate-900">
            <div class="border-b border-slate-200 px-4 py-3 dark:border-slate-800">
              <h2 class="text-base font-semibold text-slate-950 dark:text-white">{{ $t('catalog.requests.detail.progress') }}</h2>
            </div>
            <div class="p-4">
              <dl class="divide-y divide-slate-100 text-sm dark:divide-slate-800">
                <div class="flex items-center justify-between gap-3 pb-2.5"><dt class="text-slate-500 dark:text-slate-400">{{ $t('catalog.requests.fields.reward') }}</dt><dd class="inline-flex items-center gap-1 font-semibold text-amber-600 dark:text-amber-400"><UIcon name="i-lucide-coins" class="size-4" />{{ numberFormatter.format(request.rewardAmount) }}</dd></div>
                <div class="grid grid-cols-[80px_minmax(0,1fr)] gap-3 py-2.5"><dt class="text-slate-500 dark:text-slate-400">{{ $t('catalog.requests.fields.claimer') }}</dt><dd class="truncate text-right font-medium text-slate-950 dark:text-white">{{ request.claimer?.username || '-' }}</dd></div>
                <div v-if="request.claimExpiresAt" class="grid grid-cols-[80px_minmax(0,1fr)] gap-3 pt-2.5"><dt class="text-slate-500 dark:text-slate-400">{{ $t('catalog.requests.fields.claimExpiresAt') }}</dt><dd class="text-right font-medium text-slate-950 dark:text-white">{{ formatDateTime(request.claimExpiresAt, locale) }}</dd></div>
              </dl>

              <div v-if="request.status === CatalogRequestStatus.Cancelled && request.cancelReason" class="mt-3 border-t border-red-100 pt-3 text-sm dark:border-red-950">
                <p class="font-medium text-red-700 dark:text-red-300">{{ $t('catalog.requests.fields.cancelReason') }}</p>
                <p class="mt-1 break-words leading-5 text-slate-600 dark:text-slate-300">{{ request.cancelReason }}</p>
              </div>
            </div>
          </section>

          <section v-if="hasActions" class="rounded-lg border border-slate-200 bg-white p-3 dark:border-slate-800 dark:bg-slate-900">
            <div class="grid gap-2">
              <UButton v-if="request.actions.canClaim" color="primary" icon="i-lucide-hand" block :loading="actionPending === 'claim'" @click="handleClaim">{{ $t('catalog.requests.actions.claim') }}</UButton>
              <UButton v-if="request.actions.canAbandon" color="neutral" variant="outline" icon="i-lucide-log-out" block :loading="actionPending === 'abandon'" @click="handleAbandon">{{ $t('catalog.requests.actions.abandon') }}</UButton>

              <div v-if="request.actions.canSubmit && request.requestType === CatalogRequestType.Torrent" class="rounded-md border border-slate-200 p-3 dark:border-slate-800">
                <UFormField :label="$t('catalog.requests.fields.resultTorrent')">
                  <UInput v-model.number="resultTorrentId" type="number" min="1" :placeholder="$t('catalog.requests.detail.resultTorrentPlaceholder')" />
                </UFormField>
                <UButton class="mt-2" color="primary" variant="soft" icon="i-lucide-send" block :loading="actionPending === 'submit'" :disabled="resultTorrentId <= 0" @click="handleSubmitResult">{{ $t('catalog.requests.actions.submit') }}</UButton>
              </div>
              <UButton v-else-if="request.actions.canSubmit" color="primary" variant="soft" icon="i-lucide-send" block :loading="actionPending === 'submit'" @click="handleSubmitResult">{{ $t('catalog.requests.actions.submitReseed') }}</UButton>

              <UPopover v-if="request.actions.canComplete" :content="{ side: 'bottom', align: 'end', sideOffset: 8 }" :ui="{ content: 'w-64 p-3' }">
                <UButton color="success" icon="i-lucide-circle-check" block :loading="actionPending === 'complete'">{{ $t('catalog.requests.actions.complete') }}</UButton>
                <template #content="{ close }">
                  <p class="text-sm font-medium text-slate-950 dark:text-white">{{ $t('catalog.requests.actions.completeConfirm') }}</p>
                  <p class="mt-1 text-xs leading-5 text-slate-500 dark:text-slate-400">{{ $t('catalog.requests.actions.completeHint') }}</p>
                  <div class="mt-3 flex justify-end gap-2"><UButton color="neutral" variant="ghost" size="xs" @click="close">{{ $t('common.cancel') }}</UButton><UButton color="success" size="xs" :loading="actionPending === 'complete'" @click="handleComplete(close)">{{ $t('common.confirm') }}</UButton></div>
                </template>
              </UPopover>

              <UPopover v-if="request.actions.canCancel" :content="{ side: 'bottom', align: 'end', sideOffset: 8 }" :ui="{ content: 'w-72 p-3' }">
                <UButton color="error" variant="soft" icon="i-lucide-x-circle" block :loading="actionPending === 'cancel'">{{ $t('catalog.requests.actions.cancel') }}</UButton>
                <template #content="{ close }">
                  <p class="text-sm font-medium text-slate-950 dark:text-white">{{ $t('catalog.requests.actions.cancelConfirm') }}</p>
                  <p class="mt-1 text-xs leading-5 text-slate-500 dark:text-slate-400">{{ $t('catalog.requests.actions.cancelHint') }}</p>
                  <UTextarea v-model="cancelReason" class="mt-3 w-full" :rows="2" :placeholder="$t('catalog.requests.actions.cancelReason')" />
                  <div class="mt-3 flex justify-end gap-2"><UButton color="neutral" variant="ghost" size="xs" @click="close">{{ $t('common.cancel') }}</UButton><UButton color="error" size="xs" :loading="actionPending === 'cancel'" @click="handleCancel(close)">{{ $t('common.confirm') }}</UButton></div>
                </template>
              </UPopover>
            </div>
          </section>

          <section v-if="errorMessage" class="rounded-lg border border-red-200 bg-red-50 p-3 text-sm text-red-700 dark:border-red-900 dark:bg-red-950 dark:text-red-200">{{ errorMessage }}</section>
        </aside>

        <CatalogCommentSection
          v-model:content="commentContent"
          v-model:editor-mode="commentEditorMode"
          v-model:active-report-id="activeCommentReportId"
          v-model:report-reason="commentReportReason"
          class="lg:col-start-1"
          section-id="request-comments"
          composer-id="request-comment-composer"
          :comments="comments"
          :total="commentTotal"
          :page="commentPage"
          :page-size="commentPageSize"
          :can-create="canCreateCatalogComment"
          :submit-pending="commentSubmitPending"
          :pending="commentsPending"
          :error="commentsError"
          :like-pending-id="commentLikePendingId"
          :report-pending="commentReportPending"
          :submit-reward="submitCommentReward"
          @submit="handleCommentSubmit"
          @page-change="goToCommentPage"
          @toggle-like="handleCommentLike"
          @reward-success="loadComments"
          @report="handleCommentReport"
        />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ApiError } from '~/composables/useApi'
import { CatalogRequestStatus, CatalogRequestType, type CatalogRequestDetail } from '~/composables/useCatalogRequests'
import type { CommentItem, CatalogCategory } from '~/composables/useCatalogTorrents'
import { formatDateTime, formatRelativeDateTime, localizeI18nName } from '~/utils/format'
import { renderUserMarkdown } from '~/utils/richText'

type CommentEditorMode = 'write' | 'preview'

definePageMeta({ middleware: 'auth' })

const { t, locale } = useI18n()
const localePath = useLocalePath()
const route = useRoute()
const toast = useToast()
const requestApi = useCatalogRequests()
const catalogApi = useCatalogTorrents()
const { fetchUser, hasPermission } = useAuth()

const request = ref<CatalogRequestDetail | null>(null)
const categories = ref<CatalogCategory[]>([])
const pending = ref(true)
const errorMessage = ref('')
const actionPending = ref('')
const resultTorrentId = ref(0)
const cancelReason = ref('')
const comments = ref<CommentItem[]>([])
const commentTotal = ref(0)
const commentPage = ref(1)
const commentPageSize = 20
const commentsPending = ref(false)
const commentsError = ref('')
const commentContent = ref('')
const commentEditorMode = ref<CommentEditorMode>('write')
const commentSubmitPending = ref(false)
const commentLikePendingId = ref(0)
const activeCommentReportId = ref(0)
const commentReportReason = ref('')
const commentReportPending = ref(false)
const numberFormatter = computed(() => new Intl.NumberFormat(locale.value, { maximumFractionDigits: 1 }))

const requestId = computed(() => Number(route.params.id))
const renderedDescription = computed(() => renderUserMarkdown(request.value?.description || '').trim())
const canCreateCatalogComment = computed(() => hasPermission(Permission.CatalogCommentCreate))
const requestTypeLabel = computed(() => request.value?.requestType === CatalogRequestType.Reseed ? t('catalog.requests.types.reseed') : t('catalog.requests.types.torrent'))
const hasRelatedTorrent = computed(() => Boolean(request.value?.targetTorrent?.exist || request.value?.resultTorrent?.exist))
const hasBothRelatedTorrents = computed(() => Boolean(request.value?.targetTorrent?.exist && request.value?.resultTorrent?.exist))
const categoryName = computed(() => {
  const category = categories.value.find(item => item.id === request.value?.categoryId)
  return category ? localizeI18nName(category.name, locale.value) : ''
})
const hasActions = computed(() => request.value ? Object.values(request.value.actions).some(Boolean) : false)

useHead(() => ({ title: request.value?.title || t('catalog.requests.detail.metaTitle') }))

onMounted(async () => {
  await Promise.all([loadCategories(), loadRequest()])
  await loadComments()
})

async function loadCategories() {
  try { categories.value = (await catalogApi.listCategories()).list || [] } catch { categories.value = [] }
}

async function loadRequest() {
  if (!Number.isInteger(requestId.value) || requestId.value <= 0) {
    errorMessage.value = t('catalog.requests.detail.invalidId')
    pending.value = false
    return
  }
  pending.value = true
  errorMessage.value = ''
  try {
    request.value = await requestApi.getRequest(requestId.value)
  } catch (error) {
    request.value = null
    errorMessage.value = apiErrorMessage(error)
  } finally {
    pending.value = false
  }
}

async function runAction(name: string, action: () => Promise<void>, successKey: string, refreshUser = false) {
  if (actionPending.value) return
  actionPending.value = name
  errorMessage.value = ''
  try {
    await action()
    if (refreshUser) await fetchUser()
    toast.add({ title: t(successKey), color: 'success', icon: 'i-lucide-circle-check' })
    await loadRequest()
  } catch (error) {
    errorMessage.value = apiErrorMessage(error)
  } finally {
    actionPending.value = ''
  }
}

function handleClaim() { return runAction('claim', () => requestApi.claimRequest(requestId.value), 'catalog.requests.actions.claimed') }
function handleAbandon() { return runAction('abandon', () => requestApi.abandonRequest(requestId.value), 'catalog.requests.actions.abandoned') }
function handleSubmitResult() { return runAction('submit', () => requestApi.submitRequest(requestId.value, resultTorrentId.value), 'catalog.requests.actions.submitted') }
async function handleComplete(close: () => void) {
  const action = hasPermission(Permission.AdminCatalogRequestManage)
    ? () => requestApi.completeRequestByAdmin(requestId.value)
    : () => requestApi.completeRequest(requestId.value)
  await runAction('complete', action, 'catalog.requests.actions.completed', true)
  close()
}

async function handleCancel(close: () => void) {
  const action = hasPermission(Permission.AdminCatalogRequestManage)
    ? () => requestApi.cancelRequestByAdmin(requestId.value, cancelReason.value)
    : () => requestApi.cancelRequest(requestId.value, cancelReason.value)
  await runAction('cancel', action, 'catalog.requests.actions.cancelled', true)
  close()
}

async function loadComments() {
  if (!request.value) return
  commentsPending.value = true
  commentsError.value = ''
  try {
    const data = await requestApi.listComments(requestId.value, commentPage.value, commentPageSize)
    comments.value = data.list || []
    commentTotal.value = data.total || 0
  } catch (error) {
    commentsError.value = apiErrorMessage(error)
  } finally {
    commentsPending.value = false
  }
}

async function handleCommentSubmit() {
  const content = commentContent.value.trim()
  if (content.length < 3 || commentSubmitPending.value) return
  commentSubmitPending.value = true
  try {
    await requestApi.createComment(requestId.value, content)
    commentContent.value = ''
    commentEditorMode.value = 'write'
    commentPage.value = 1
    toast.add({ title: t('catalog.torrents.detail.comments.created'), color: 'success', icon: 'i-lucide-circle-check' })
    await loadComments()
  } catch (error) {
    commentsError.value = apiErrorMessage(error)
  } finally {
    commentSubmitPending.value = false
  }
}

async function handleCommentLike(comment: CommentItem) {
  if (commentLikePendingId.value) return
  commentLikePendingId.value = comment.id
  try {
    const data = await requestApi.toggleCommentLike(requestId.value, comment.id)
    comment.isLiked = data.isLiked
    comment.likeCount += data.isLiked ? 1 : -1
  } catch (error) {
    commentsError.value = apiErrorMessage(error)
  } finally {
    commentLikePendingId.value = 0
  }
}

async function submitCommentReward(comment: CommentItem, amount: number) {
  await requestApi.rewardComment(requestId.value, comment.id, amount)
  await fetchUser()
}

async function handleCommentReport(commentId: number) {
  const reason = commentReportReason.value.trim()
  if (reason.length < 5 || commentReportPending.value) return
  commentReportPending.value = true
  try {
    await requestApi.reportComment(requestId.value, commentId, reason)
    toast.add({ title: t('catalog.torrents.detail.report.success'), color: 'success', icon: 'i-lucide-circle-check' })
    activeCommentReportId.value = 0
    commentReportReason.value = ''
  } catch (error) {
    commentsError.value = apiErrorMessage(error)
  } finally {
    commentReportPending.value = false
  }
}

function goToCommentPage(value: number) {
  commentPage.value = value
  loadComments()
}

function requestStatusLabel(value: number) { return t(`catalog.requests.status.${value}`) }
function requestStatusColor(value: number): 'neutral' | 'primary' | 'warning' | 'success' | 'error' {
  if (value === CatalogRequestStatus.Open) return 'primary'
  if (value === CatalogRequestStatus.Claimed || value === CatalogRequestStatus.Submitted) return 'warning'
  if (value === CatalogRequestStatus.Completed) return 'success'
  if (value === CatalogRequestStatus.Cancelled) return 'error'
  return 'neutral'
}

function apiErrorMessage(error: unknown) { return error instanceof ApiError ? error.message : t('common.requestFailed') }
</script>
