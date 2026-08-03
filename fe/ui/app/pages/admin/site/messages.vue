<template>
  <div class="min-h-[calc(100vh-4rem)] bg-slate-50 py-6 dark:bg-slate-950">
    <div class="w-full px-3 sm:px-4 lg:px-5">
      <div class="grid gap-4 xl:grid-cols-[minmax(360px,430px)_minmax(0,1fr)]">
        <section class="overflow-hidden rounded-lg border border-slate-200 bg-white dark:border-slate-800 dark:bg-slate-900">
          <div class="border-b border-slate-200 px-4 py-3 dark:border-slate-800">
            <h2 class="text-sm font-semibold text-slate-950 dark:text-white">{{ $t('admin.site.messages.send') }}</h2>
          </div>

          <form class="space-y-3 p-4" @submit.prevent="sendMessage">
            <label class="block">
              <span class="flex items-center justify-between gap-3">
                <span class="text-sm font-medium text-slate-700 dark:text-slate-200">{{ $t('admin.site.messages.fields.receiverIds') }}</span>
                <UBadge v-if="receiverCount > 0" color="neutral" variant="soft">{{ $t('admin.site.messages.receiverCount', { count: receiverCount }) }}</UBadge>
              </span>
              <textarea v-model="form.receiverIds" rows="2" class="mt-1 w-full resize-y rounded-md border border-slate-200 bg-white px-3 py-2 font-mono text-sm text-slate-950 outline-none transition focus:border-sky-400 focus:ring-2 focus:ring-sky-100 dark:border-slate-700 dark:bg-slate-950 dark:text-white dark:focus:border-sky-500 dark:focus:ring-sky-950" :placeholder="$t('admin.site.messages.receiverPlaceholder')" />
            </label>

            <label class="block">
              <span class="text-sm font-medium text-slate-700 dark:text-slate-200">{{ $t('admin.site.messages.fields.title') }}</span>
              <input v-model.trim="form.title" class="mt-1 h-10 w-full rounded-md border border-slate-200 bg-white px-3 text-sm text-slate-950 outline-none transition focus:border-sky-400 focus:ring-2 focus:ring-sky-100 dark:border-slate-700 dark:bg-slate-950 dark:text-white dark:focus:border-sky-500 dark:focus:ring-sky-950">
            </label>

            <div>
              <span class="text-sm font-medium text-slate-700 dark:text-slate-200">{{ $t('admin.site.messages.fields.content') }}</span>
              <RichTextComposer
                v-model="form.content"
                v-model:mode="editorMode"
                class="mt-1"
                :as-form="false"
                :rows="6"
                :disabled="sending"
                :placeholder="$t('admin.site.messages.contentPlaceholder')"
                :submit-label="$t('admin.site.messages.submit')"
                :write-label="$t('common.editor.edit')"
                :preview-label="$t('common.editor.preview')"
                :preview-empty="$t('common.editor.previewEmpty')"
                :show-actions="false"
              />
            </div>

            <div class="grid gap-3 sm:grid-cols-[minmax(0,1fr)_130px]">
              <UFormField :label="$t('admin.site.messages.fields.targetType')">
                <USelect v-model="form.targetType" class="w-full" size="lg" :ui="{ base: 'h-10 w-full' }" :items="messageTargetTypeOptions" value-key="value" />
              </UFormField>
              <label class="block">
                <span class="text-sm font-medium text-slate-700 dark:text-slate-200">{{ $t('admin.site.messages.fields.targetId') }}</span>
                <input v-model.number="form.targetId" type="number" min="0" class="mt-1 h-10 w-full rounded-md border border-slate-200 bg-white px-3 text-sm text-slate-950 outline-none transition disabled:cursor-not-allowed disabled:bg-slate-50 disabled:text-slate-400 focus:border-sky-400 focus:ring-2 focus:ring-sky-100 dark:border-slate-700 dark:bg-slate-950 dark:text-white dark:disabled:bg-slate-900 dark:disabled:text-slate-600 dark:focus:border-sky-500 dark:focus:ring-sky-950" :disabled="form.targetType === noMessageTargetValue">
              </label>
            </div>

            <p v-if="formError" class="rounded-md border border-red-200 bg-red-50 px-3 py-2 text-sm text-red-700 dark:border-red-900 dark:bg-red-950/40 dark:text-red-200">{{ formError }}</p>

            <div class="flex justify-end gap-2 border-t border-slate-200 pt-4 dark:border-slate-800">
              <UButton type="button" color="neutral" variant="outline" :disabled="sending || !isFormFilled" @click="resetForm">{{ $t('admin.actions.reset') }}</UButton>
              <UButton type="submit" color="primary" icon="i-lucide-send" :loading="sending">{{ $t('admin.site.messages.submit') }}</UButton>
            </div>
          </form>
        </section>

        <section class="overflow-hidden rounded-lg border border-slate-200 bg-white dark:border-slate-800 dark:bg-slate-900">
          <div class="border-b border-slate-200 px-4 py-3 dark:border-slate-800">
            <div class="flex flex-col gap-3 lg:flex-row lg:items-center lg:justify-between">
              <h2 class="text-sm font-semibold text-slate-950 dark:text-white">{{ $t('admin.site.messages.list') }}</h2>
              <form class="flex flex-wrap items-center gap-2" @submit.prevent="applyFilters">
                <input v-model.trim="receiverFilter" inputmode="numeric" class="h-8 w-28 rounded-md border border-slate-200 bg-white px-2.5 text-sm outline-none transition focus:border-sky-300 dark:border-slate-700 dark:bg-slate-950 dark:focus:border-sky-700" :placeholder="$t('admin.site.messages.filters.receiverId')">
                <UTooltip :text="$t('admin.site.messages.filters.search')" :content="{ side: 'top', sideOffset: 8 }" :delay-duration="300">
                  <UButton type="submit" class="h-8" size="sm" color="neutral" variant="outline" icon="i-lucide-search" :loading="pending" :aria-label="$t('admin.site.messages.filters.search')" />
                </UTooltip>
                <div class="inline-flex h-8 items-center rounded-md border border-slate-200 bg-slate-50 p-0.5 dark:border-slate-800 dark:bg-slate-950/70" :aria-label="$t('admin.site.messages.filters.readState')" role="group">
                  <button
                    v-for="option in readFilterOptions"
                    :key="option.value"
                    type="button"
                    class="h-7 rounded px-2.5 text-xs font-medium transition disabled:cursor-not-allowed disabled:opacity-60"
                    :class="readFilter === option.value ? 'bg-white text-slate-950 shadow-sm dark:bg-slate-800 dark:text-white' : 'text-slate-500 hover:text-slate-950 dark:text-slate-400 dark:hover:text-white'"
                    :disabled="pending"
                    @click="setReadFilter(option.value)"
                  >
                    {{ option.label }}
                  </button>
                </div>
              </form>
            </div>
          </div>

          <div v-if="pending" class="space-y-2 p-4">
            <div v-for="item in 8" :key="item" class="h-16 animate-pulse rounded-md bg-slate-100 dark:bg-slate-800" />
          </div>
          <div v-else-if="errorMessage" class="flex flex-col items-center justify-center px-4 py-16 text-center">
            <UIcon name="i-lucide-circle-alert" class="size-9 text-red-500" />
            <p class="mt-3 text-sm font-medium text-slate-950 dark:text-white">{{ errorMessage }}</p>
          </div>
          <div v-else-if="messages.length === 0" class="flex flex-col items-center justify-center px-4 py-16 text-center">
            <UIcon name="i-lucide-inbox" class="size-9 text-slate-400" />
            <p class="mt-3 text-sm font-medium text-slate-950 dark:text-white">{{ $t('admin.site.messages.empty') }}</p>
          </div>
          <div v-else class="overflow-x-auto">
            <table class="min-w-[820px] w-full table-fixed border-collapse text-left">
              <thead class="bg-slate-50 text-xs font-medium uppercase text-slate-500 dark:bg-slate-950/70 dark:text-slate-400">
                <tr>
                  <th class="w-[48%] border-b border-slate-200 px-4 py-3 dark:border-slate-800">{{ $t('admin.site.messages.table.message') }}</th>
                  <th class="w-[22%] border-b border-slate-200 px-4 py-3 dark:border-slate-800">{{ $t('admin.site.messages.table.receiver') }}</th>
                  <th class="w-[10%] border-b border-slate-200 px-4 py-3 dark:border-slate-800">{{ $t('admin.site.messages.table.state') }}</th>
                  <th class="w-[20%] border-b border-slate-200 px-4 py-3 text-right dark:border-slate-800">{{ $t('admin.site.messages.table.time') }}</th>
                </tr>
              </thead>
              <tbody>
                <tr
                  v-for="item in messages"
                  :key="item.id"
                  class="group cursor-pointer border-b border-slate-200 transition-colors last:border-b-0 hover:bg-slate-50 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-inset focus-visible:ring-sky-300 dark:border-slate-800 dark:hover:bg-slate-950/70 dark:focus-visible:ring-sky-700"
                  role="button"
                  tabindex="0"
                  @click="openMessage(item)"
                  @keydown.enter.prevent="openMessage(item)"
                  @keydown.space.prevent="openMessage(item)"
                >
                  <td class="px-4 py-3">
                    <div class="min-w-0">
                      <div class="flex min-w-0 items-center gap-2">
                        <p class="truncate text-sm font-semibold text-slate-950 transition group-hover:text-sky-700 dark:text-white dark:group-hover:text-sky-300">{{ item.title }}</p>
                        <UBadge v-if="messageTargetLabel(item)" color="neutral" variant="soft" class="shrink-0">{{ messageTargetLabel(item) }}</UBadge>
                      </div>
                      <p class="mt-1 line-clamp-2 text-sm leading-5 text-slate-500 dark:text-slate-400">{{ messageExcerpt(item) }}</p>
                    </div>
                  </td>
                  <td class="px-4 py-3">
                    <div class="flex min-w-0 items-center gap-2.5">
                      <IamUserAvatar :user="item.receiver" size="sm" />
                      <div class="min-w-0">
                        <IamUserPopover :id="item.receiverId" :user="item.receiver" class="truncate text-sm font-medium text-slate-950 dark:text-white" />
                        <p class="mt-0.5 truncate text-xs text-slate-500 dark:text-slate-400">#{{ item.receiverId }}</p>
                      </div>
                    </div>
                  </td>
                  <td class="px-4 py-3">
                    <UBadge :color="item.isRead ? 'neutral' : 'warning'" variant="soft">{{ item.isRead ? $t('admin.site.messages.readState.read') : $t('admin.site.messages.readState.unread') }}</UBadge>
                  </td>
                  <td class="px-4 py-3 text-right text-sm text-slate-600 dark:text-slate-300">
                    <UTooltip :text="formatDateTime(item.createdAt, locale)" :content="{ side: 'top', sideOffset: 8 }" :delay-duration="600">
                      <span>{{ formatDateTime(item.createdAt, locale) }}</span>
                    </UTooltip>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>

          <AppPager
            class="border-t border-slate-200 px-4 py-3 dark:border-slate-800"
            size="sm"
            :page="page"
            :total="total"
            :page-size="size"
            :page-size-options="pageSizes"
            :disabled="pending"
            @page-change="goToPage"
            @page-size-change="changeSize"
          />
        </section>
      </div>
    </div>

    <UModal
      :open="detailModalOpen"
      :title="selectedMessage?.title || $t('admin.site.messages.detailModal.title')"
      :description="selectedMessageDescription"
      :ui="{
        content: 'sm:max-w-3xl overflow-hidden',
        header: 'min-h-0 px-5 py-4 sm:px-5',
        body: 'p-0 sm:p-0',
        title: 'text-base font-semibold text-slate-950 dark:text-white',
        description: 'mt-1 text-sm text-slate-500 dark:text-slate-400',
        close: 'top-4 end-4'
      }"
      @update:open="setDetailModalOpen"
    >
      <template #body>
        <div>
          <div class="flex flex-wrap items-center justify-between gap-2 border-b border-slate-200 bg-slate-50 px-5 py-2.5 dark:border-slate-800 dark:bg-slate-950/70">
            <div class="flex min-w-0 flex-wrap items-center gap-2">
              <UBadge v-if="selectedMessage" :color="selectedMessage.isRead ? 'neutral' : 'warning'" variant="soft">
                {{ selectedMessage.isRead ? $t('admin.site.messages.readState.read') : $t('admin.site.messages.readState.unread') }}
              </UBadge>
              <UBadge v-if="selectedMessage && messageTargetLabel(selectedMessage)" color="neutral" variant="soft">
                {{ messageTargetLabel(selectedMessage) }}
              </UBadge>
            </div>
            <div class="flex items-center gap-2">
              <UButton v-if="selectedMessageTargetPath" color="neutral" variant="soft" size="xs" icon="i-lucide-arrow-up-right" :to="localePath(selectedMessageTargetPath)">
                {{ $t('site.messages.openTarget') }}
              </UButton>
              <UButton color="neutral" variant="soft" size="xs" icon="i-lucide-copy" :disabled="!selectedMessage?.content" @click="copyMessageContent">
                {{ $t('common.copy') }}
              </UButton>
            </div>
          </div>
          <div class="max-h-[62vh] overflow-auto px-5 py-4">
            <div v-if="selectedMessageContent" class="rich-text" v-html="selectedMessageContent" />
            <p v-else class="text-sm text-slate-500 dark:text-slate-400">{{ $t('admin.site.messages.noContent') }}</p>
          </div>
        </div>
      </template>
    </UModal>
  </div>
</template>

<script setup lang="ts">
import { ApiError } from '~/composables/useApi'
import type { AdminSiteMessage } from '~/composables/useAdmin'
import { formatDateTime } from '~/utils/format'
import { renderUserMarkdown } from '~/utils/richText'

definePageMeta({ layout: 'admin', middleware: 'admin' })

type EditorMode = 'write' | 'preview'

const { t, locale } = useI18n()
const toast = useToast()
const localePath = useLocalePath()
const adminApi = useAdmin()

const messages = ref<AdminSiteMessage[]>([])
const total = ref(0)
const page = ref(1)
const size = ref(20)
const pageSizes = [20, 50, 100]
const pending = ref(false)
const sending = ref(false)
const errorMessage = ref('')
const formError = ref('')
const receiverFilter = ref('')
const readFilter = ref('')
const editorMode = ref<EditorMode>('write')
const detailModalOpen = ref(false)
const selectedMessage = ref<AdminSiteMessage | null>(null)
const noMessageTargetValue = '__none__'
const form = reactive({ receiverIds: '', title: '', content: '', targetType: noMessageTargetValue, targetId: 0 })

const readFilterOptions = computed(() => [
  { value: '', label: t('admin.site.messages.filters.all') },
  { value: 'false', label: t('admin.site.messages.filters.unread') },
  { value: 'true', label: t('admin.site.messages.filters.read') }
])
const messageTargetTypeOptions = computed(() => [
  { value: noMessageTargetValue, label: t('admin.site.messages.targetTypes.none') },
  { value: 'catalog_torrent', label: t('admin.site.messages.targetTypes.catalogTorrent') },
  { value: 'forum_topic', label: t('admin.site.messages.targetTypes.forumTopic') }
])
const receiverCount = computed(() => parseReceiverIds(form.receiverIds).length)
const hasMessageTarget = computed(() => form.targetType !== noMessageTargetValue)
const isFormFilled = computed(() => Boolean(form.receiverIds || form.title || form.content || hasMessageTarget.value || form.targetId))
const selectedMessageContent = computed(() => renderUserMarkdown(selectedMessage.value?.content || '').trim())
const selectedMessageTargetPath = computed(() => selectedMessage.value ? messageTargetPath(selectedMessage.value) : '')
const selectedMessageDescription = computed(() => {
  if (!selectedMessage.value) return ''
  return `${t('admin.site.messages.toUser', { id: selectedMessage.value.receiverId, name: selectedMessage.value.receiver?.username || '-' })} / ${formatDateTime(selectedMessage.value.createdAt, locale.value)}`
})

watch(() => form.targetType, (targetType) => {
  if (targetType === noMessageTargetValue) form.targetId = 0
})

useHead(() => ({ title: t('admin.site.messages.title') }))
onMounted(loadMessages)

async function loadMessages() {
  pending.value = true
  errorMessage.value = ''
  try {
    const data = await adminApi.listSiteMessages({
      page: page.value,
      size: size.value,
      receiverId: receiverFilter.value ? Number(receiverFilter.value) : undefined,
      isRead: readFilter.value === '' ? undefined : readFilter.value === 'true'
    })
    messages.value = data.list || []
    total.value = data.total || 0
  } catch (error) {
    errorMessage.value = error instanceof ApiError ? error.message : t('common.requestFailed')
  } finally {
    pending.value = false
  }
}

async function sendMessage() {
  formError.value = ''
  const receiverIds = parseReceiverIds(form.receiverIds)
  if (receiverIds.length === 0 || !form.title || !form.content) {
    formError.value = t('admin.site.messages.required')
    return
  }
  if (hasMessageTarget.value && Number(form.targetId || 0) <= 0) {
    formError.value = t('admin.site.messages.targetRequired')
    return
  }
  sending.value = true
  try {
    const data = await adminApi.createSiteMessage({
      receiverIds,
      title: form.title,
      content: form.content,
      targetType: hasMessageTarget.value ? form.targetType : '',
      targetId: hasMessageTarget.value ? Number(form.targetId) || 0 : 0
    })
    toast.add({ color: 'success', title: t('admin.site.messages.sent', { count: data.count || 0 }), icon: 'i-lucide-check' })
    resetForm()
    await loadMessages()
  } catch (error) {
    formError.value = error instanceof ApiError ? error.message : t('common.requestFailed')
  } finally {
    sending.value = false
  }
}

function parseReceiverIds(value: string) {
  return Array.from(new Set(String(value || '')
    .split(/[\s,，]+/)
    .map(item => Number(item.trim()))
    .filter(item => Number.isInteger(item) && item > 0)))
}

function resetForm() {
  form.receiverIds = ''
  form.title = ''
  form.content = ''
  form.targetType = noMessageTargetValue
  form.targetId = 0
  editorMode.value = 'write'
  formError.value = ''
}

function applyFilters() {
  page.value = 1
  loadMessages()
}

function setReadFilter(value: string) {
  if (readFilter.value === value) return
  readFilter.value = value
  applyFilters()
}

function goToPage(nextPage: number) {
  page.value = nextPage
  loadMessages()
}

function changeSize(nextSize: number) {
  size.value = nextSize
  page.value = 1
  loadMessages()
}

function openMessage(item: AdminSiteMessage) {
  selectedMessage.value = item
  detailModalOpen.value = true
}

function setDetailModalOpen(value: boolean) {
  detailModalOpen.value = value
  if (!value) selectedMessage.value = null
}

function messageExcerpt(item: AdminSiteMessage) {
  return String(item.content || '')
    .replace(/[`*_>#\[\]()]/g, '')
    .replace(/\s+/g, ' ')
    .trim() || t('admin.site.messages.noContent')
}

function messageTargetLabel(item: AdminSiteMessage) {
  if (!item.targetType || !item.targetId) return ''
  return `${targetTypeLabel(item.targetType)} #${item.targetId}`
}

function targetTypeLabel(type: string) {
  if (type === 'catalog_torrent') return t('admin.site.messages.targetTypes.catalogTorrent')
  if (type === 'forum_topic') return t('admin.site.messages.targetTypes.forumTopic')
  return type
}

function messageTargetPath(item: AdminSiteMessage) {
  if (!item.targetType || !item.targetId) return ''
  if (item.targetType === 'catalog_torrent') return `/catalog/torrents/${item.targetId}`
  if (item.targetType === 'catalog_request') return `/catalog/requests/${item.targetId}`
  if (item.targetType === 'forum_topic') return `/forum/topics/${item.targetId}`
  return ''
}

async function copyMessageContent() {
  if (!selectedMessage.value?.content || typeof navigator === 'undefined' || !navigator.clipboard) return
  await navigator.clipboard.writeText(selectedMessage.value.content)
}
</script>
