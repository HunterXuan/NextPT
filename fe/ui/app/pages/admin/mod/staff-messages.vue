<template>
  <div class="min-h-[calc(100vh-4rem)] bg-slate-50 py-6 dark:bg-slate-950">
    <div class="w-full px-3 sm:px-4 lg:px-5">
      <div class="grid gap-4 xl:grid-cols-[minmax(0,1fr)_380px]">
        <section class="min-w-0 overflow-hidden rounded-lg border border-slate-200 bg-white dark:border-slate-800 dark:bg-slate-900">
          <div class="border-b border-slate-200 px-4 py-3 dark:border-slate-800">
            <div class="flex flex-col gap-3 xl:flex-row xl:items-center xl:justify-between">
              <div class="flex items-center gap-2">
                <UIcon name="i-lucide-mail-question" class="size-5 text-sky-600 dark:text-sky-300" />
                <h1 class="text-sm font-semibold text-slate-950 dark:text-white">{{ $t('admin.mod.staffMessages.title') }}</h1>
              </div>
              <form class="flex min-w-0 items-center gap-2" @submit.prevent="search">
                <UInput v-model="senderIdInput" class="w-32" size="sm" :placeholder="$t('admin.mod.staffMessages.filters.senderId')" />
                <UButton type="submit" color="neutral" variant="soft" size="sm" icon="i-lucide-search" :loading="pending" :aria-label="$t('common.search')" />
                <USelect v-model="statusFilter" class="w-32" size="sm" :items="statusOptions" value-key="value" @update:model-value="search" />
              </form>
            </div>
          </div>

          <div v-if="pending" class="space-y-2 p-4">
            <div v-for="item in 7" :key="item" class="h-16 animate-pulse rounded-md bg-slate-100 dark:bg-slate-800" />
          </div>
          <div v-else-if="errorMessage" class="flex flex-col items-center justify-center px-4 py-16 text-center">
            <UIcon name="i-lucide-circle-alert" class="size-9 text-red-500" />
            <p class="mt-3 text-sm font-medium text-slate-950 dark:text-white">{{ errorMessage }}</p>
          </div>
          <div v-else-if="messages.length === 0" class="flex flex-col items-center justify-center px-4 py-16 text-center">
            <UIcon name="i-lucide-inbox" class="size-9 text-slate-400" />
            <p class="mt-3 text-sm font-medium text-slate-950 dark:text-white">{{ $t('admin.mod.staffMessages.empty') }}</p>
          </div>
          <div v-else class="overflow-x-auto">
            <table class="w-full min-w-[720px] table-fixed border-collapse text-left">
              <thead class="bg-slate-50 text-xs font-medium text-slate-500 dark:bg-slate-950/70 dark:text-slate-400">
                <tr>
                  <th class="w-[43%] border-b border-slate-200 px-4 py-3 dark:border-slate-800">{{ $t('admin.mod.staffMessages.table.message') }}</th>
                  <th class="w-[20%] border-b border-slate-200 px-4 py-3 dark:border-slate-800">{{ $t('admin.mod.staffMessages.table.sender') }}</th>
                  <th class="w-[17%] border-b border-slate-200 px-4 py-3 text-center dark:border-slate-800">{{ $t('admin.mod.staffMessages.table.status') }}</th>
                  <th class="w-[20%] border-b border-slate-200 px-4 py-3 text-right dark:border-slate-800">{{ $t('admin.mod.staffMessages.table.time') }}</th>
                </tr>
              </thead>
              <tbody>
                <tr
                  v-for="item in messages"
                  :key="item.id"
                  class="cursor-pointer border-b border-slate-200 transition last:border-b-0 hover:bg-slate-50 dark:border-slate-800 dark:hover:bg-slate-950/70"
                  :class="selectedMessageId === item.id ? 'bg-sky-50/70 dark:bg-sky-950/30' : ''"
                  @click="selectedMessageId = item.id"
                >
                  <td class="px-4 py-3">
                    <p class="truncate text-sm font-semibold text-slate-950 dark:text-white">{{ item.subject }}</p>
                    <p class="mt-1 truncate text-xs text-slate-500 dark:text-slate-400">{{ messageExcerpt(item.content) }}</p>
                  </td>
                  <td class="px-4 py-3">
                    <IamUserPopover :user="item.sender" :fallback="userName(item.sender, item.senderId)" class="truncate text-sm text-slate-700 dark:text-slate-200" />
                  </td>
                  <td class="px-4 py-3 text-center"><UBadge :color="statusColor(item.status)" variant="soft" class="whitespace-nowrap">{{ statusLabel(item.status) }}</UBadge></td>
                  <td class="px-4 py-3 text-right text-xs text-slate-500 dark:text-slate-400">{{ formatRelativeDateTime(item.createdAt, locale) }}</td>
                </tr>
              </tbody>
            </table>
          </div>

          <AppPager
            v-if="!pending && messages.length > 0"
            class="border-t border-slate-200 px-4 py-3 dark:border-slate-800"
            size="sm"
            :page="page"
            :total="total"
            :page-size="size"
            :disabled="pending"
            @page-change="changePage"
          />
        </section>

        <section class="overflow-hidden rounded-lg border border-slate-200 bg-white dark:border-slate-800 dark:bg-slate-900 xl:sticky xl:top-20 xl:self-start">
          <div v-if="!selectedMessage" class="flex min-h-[360px] flex-col items-center justify-center px-6 py-12 text-center">
            <span class="flex size-12 items-center justify-center rounded-lg bg-slate-100 text-slate-400 dark:bg-slate-800 dark:text-slate-500"><UIcon name="i-lucide-mouse-pointer-2" class="size-5" /></span>
            <h2 class="mt-4 text-sm font-semibold text-slate-950 dark:text-white">{{ $t('admin.mod.staffMessages.select') }}</h2>
            <p class="mt-2 max-w-xs text-sm leading-6 text-slate-500 dark:text-slate-400">{{ $t('admin.mod.staffMessages.selectHint') }}</p>
          </div>
          <template v-else>
            <div class="flex items-center justify-between gap-3 border-b border-slate-200 px-4 py-3 dark:border-slate-800">
              <div class="min-w-0">
                <p class="text-xs text-slate-500 dark:text-slate-400">{{ $t('admin.mod.staffMessages.detail.current') }}</p>
                <h2 class="mt-1 truncate text-sm font-semibold text-slate-950 dark:text-white">{{ selectedMessage.subject }}</h2>
              </div>
              <UBadge :color="statusColor(selectedMessage.status)" variant="soft" class="shrink-0">{{ statusLabel(selectedMessage.status) }}</UBadge>
            </div>
            <div class="space-y-5 p-4">
              <div class="flex items-center gap-2 border-b border-slate-100 pb-4 dark:border-slate-800">
                <IamUserPopover :user="selectedMessage.sender" :fallback="userName(selectedMessage.sender, selectedMessage.senderId)" show-avatar avatar-size="sm" class="text-sm font-medium text-slate-950 dark:text-white" />
                <span class="text-xs text-slate-500 dark:text-slate-400">#{{ selectedMessage.senderId }}</span>
              </div>
              <section>
                <p class="mb-2 text-xs font-medium text-slate-500 dark:text-slate-400">{{ $t('admin.mod.staffMessages.detail.message') }}</p>
                <div class="rich-text rounded-md border border-slate-200 bg-slate-50 px-3 py-2.5 dark:border-slate-800 dark:bg-slate-950" v-html="renderUserMarkdown(selectedMessage.content)" />
              </section>
              <section v-if="selectedMessage.answer" class="border-l-2 border-emerald-400 pl-3 dark:border-emerald-500">
                <p class="text-xs font-medium text-emerald-700 dark:text-emerald-300">{{ $t('admin.mod.staffMessages.detail.existingReply') }}</p>
                <div class="rich-text mt-2 text-sm" v-html="renderUserMarkdown(selectedMessage.answer)" />
              </section>
              <form v-if="selectedMessage.status === ModStaffMessageStatus.Pending" class="space-y-3 border-t border-slate-200 pt-4 dark:border-slate-800" @submit.prevent="submitUpdate">
                <UFormField :label="$t('admin.mod.staffMessages.detail.replyOptional')">
                  <UTextarea v-model="answer" class="w-full" :rows="5" :disabled="processing" :placeholder="$t('admin.mod.staffMessages.detail.replyPlaceholder')" />
                </UFormField>
                <p v-if="formError" class="rounded-md border border-red-200 bg-red-50 px-3 py-2 text-sm text-red-700 dark:border-red-900 dark:bg-red-950/40 dark:text-red-200">{{ formError }}</p>
                <div class="flex flex-wrap justify-end gap-2">
                  <UButton type="submit" color="primary" icon="i-lucide-check" :loading="processing">{{ $t('admin.mod.staffMessages.actions.process') }}</UButton>
                </div>
              </form>
              <p v-else class="text-sm text-slate-500 dark:text-slate-400">{{ $t('admin.mod.staffMessages.detail.processed') }}</p>
            </div>
          </template>
        </section>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ApiError } from '~/composables/useApi'
import { ModStaffMessageStatus } from '~/composables/useMod'
import type { AdminModStaffMessage, AdminUserSummary } from '~/composables/useAdmin'
import { formatRelativeDateTime } from '~/utils/format'
import { renderUserMarkdown } from '~/utils/richText'

definePageMeta({ layout: 'admin', middleware: 'admin' })

const { t, locale } = useI18n()
const toast = useToast()
const adminApi = useAdmin()

const messages = ref<AdminModStaffMessage[]>([])
const selectedMessageId = ref(0)
const total = ref(0)
const page = ref(1)
const size = ref(30)
const statusFilter = ref(-1)
const senderIdInput = ref('')
const senderId = ref(0)
const pending = ref(false)
const processing = ref(false)
const errorMessage = ref('')
const formError = ref('')
const answer = ref('')

const selectedMessage = computed(() => messages.value.find(item => item.id === selectedMessageId.value) || null)
const statusOptions = computed(() => [
  { value: -1, label: t('admin.mod.staffMessages.status.all') },
  { value: ModStaffMessageStatus.Pending, label: t('admin.mod.staffMessages.status.pending') },
  { value: ModStaffMessageStatus.Processed, label: t('admin.mod.staffMessages.status.processed') }
])

watch(selectedMessage, (message) => {
  answer.value = message?.answer || ''
}, { immediate: true })

useHead(() => ({ title: t('admin.mod.staffMessages.title') }))
onMounted(loadMessages)

async function loadMessages() {
  pending.value = true
  errorMessage.value = ''
  try {
    const data = await adminApi.listModStaffMessages({
      page: page.value,
      size: size.value,
      status: statusFilter.value >= 0 ? statusFilter.value : undefined,
      senderId: senderId.value || undefined
    })
    messages.value = data.list || []
    total.value = data.total || 0
    if (!messages.value.some(item => item.id === selectedMessageId.value)) {
      selectedMessageId.value = messages.value[0]?.id || 0
      answer.value = selectedMessage.value?.answer || ''
    }
  } catch (error) {
    messages.value = []
    total.value = 0
    selectedMessageId.value = 0
    errorMessage.value = error instanceof ApiError ? error.message : t('common.requestFailed')
  } finally {
    pending.value = false
  }
}

function search() {
  senderId.value = Number(senderIdInput.value) > 0 ? Number(senderIdInput.value) : 0
  page.value = 1
  void loadMessages()
}

function changePage(value: number) {
  page.value = value
  void loadMessages()
}

async function submitUpdate() {
  if (!selectedMessage.value) return
  formError.value = ''
  processing.value = true
  try {
    await adminApi.updateModStaffMessage(selectedMessage.value.id, { answer: answer.value.trim() })
    toast.add({ title: t('admin.mod.staffMessages.updated'), color: 'success', icon: 'i-lucide-check' })
    await loadMessages()
    answer.value = selectedMessage.value?.answer || ''
  } catch (error) {
    formError.value = error instanceof ApiError ? error.message : t('common.requestFailed')
  } finally {
    processing.value = false
  }
}

function statusLabel(status: number) {
  if (status === ModStaffMessageStatus.Processed) return t('admin.mod.staffMessages.status.processed')
  return t('admin.mod.staffMessages.status.pending')
}

function statusColor(status: number) {
  if (status === ModStaffMessageStatus.Processed) return 'success' as const
  return 'warning' as const
}

function userName(user: AdminUserSummary, id: number) {
  return user?.username || `#${id}`
}

function messageExcerpt(value: string) {
  const text = value.replace(/\s+/g, ' ').trim()
  return text.length > 70 ? `${text.slice(0, 70)}...` : text
}
</script>
