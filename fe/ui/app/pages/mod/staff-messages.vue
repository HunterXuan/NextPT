<template>
  <div class="min-h-[calc(100vh-4rem)] bg-slate-50 py-6 dark:bg-slate-950">
    <div class="w-full px-3 sm:px-4 lg:px-5">
      <section class="overflow-hidden rounded-lg border border-slate-200 bg-white dark:border-slate-800 dark:bg-slate-900">
        <div class="flex flex-col gap-3 border-b border-slate-200 px-4 py-3 sm:flex-row sm:items-center sm:justify-between dark:border-slate-800">
          <div class="flex min-w-0 items-center gap-3">
            <span class="flex size-9 shrink-0 items-center justify-center rounded-md bg-sky-50 text-sky-600 dark:bg-sky-950/40 dark:text-sky-300">
              <UIcon name="i-lucide-mail-question" class="size-4" />
            </span>
            <div class="min-w-0">
              <h1 class="text-base font-semibold text-slate-950 dark:text-white">{{ $t('mod.staffMessages.title') }}</h1>
              <p class="mt-0.5 text-xs text-slate-500 dark:text-slate-400">{{ $t('mod.staffMessages.description') }}</p>
            </div>
          </div>
          <AppPermissionButton
            :permission="Permission.ModStaffMessageCreate"
            color="primary"
            size="sm"
            icon="i-lucide-pen-line"
            :disabled="pending"
            :tooltip="$t('mod.staffMessages.actions.create')"
            @click="composeOpen = true"
          >
            {{ $t('mod.staffMessages.actions.create') }}
          </AppPermissionButton>
        </div>

        <div class="flex gap-1 overflow-x-auto border-b border-slate-200 px-3 py-2 dark:border-slate-800">
          <button
            v-for="option in filterOptions"
            :key="option.key"
            type="button"
            class="shrink-0 rounded-md px-3 py-1.5 text-sm font-medium transition"
            :class="filter === option.key ? 'bg-sky-50 text-sky-700 dark:bg-sky-950/40 dark:text-sky-300' : 'text-slate-500 hover:bg-slate-50 hover:text-slate-950 dark:text-slate-400 dark:hover:bg-slate-800 dark:hover:text-white'"
            @click="changeFilter(option.key)"
          >
            {{ option.label }}
          </button>
        </div>

        <div v-if="pending" class="grid min-h-[520px] gap-0 lg:grid-cols-[380px_minmax(0,1fr)]">
          <div class="space-y-2 border-r border-slate-200 p-3 dark:border-slate-800">
            <div v-for="item in 7" :key="item" class="h-20 animate-pulse rounded-md bg-slate-100 dark:bg-slate-800" />
          </div>
          <div class="p-5"><div class="h-5 w-1/3 animate-pulse rounded bg-slate-100 dark:bg-slate-800" /></div>
        </div>

        <div v-else-if="errorMessage" class="flex flex-col items-center justify-center px-4 py-16 text-center">
          <UIcon name="i-lucide-circle-alert" class="size-9 text-red-500" />
          <p class="mt-3 text-sm font-medium text-slate-950 dark:text-white">{{ errorMessage }}</p>
        </div>

        <div v-else-if="messages.length === 0" class="flex flex-col items-center justify-center px-4 py-16 text-center">
          <UIcon name="i-lucide-mail-open" class="size-9 text-slate-400" />
          <p class="mt-3 text-sm font-medium text-slate-950 dark:text-white">{{ $t('mod.staffMessages.empty') }}</p>
        </div>

        <div v-else class="grid min-h-[520px] gap-0 lg:grid-cols-[380px_minmax(0,1fr)]">
          <div class="border-r border-slate-200 dark:border-slate-800">
            <div class="divide-y divide-slate-100 dark:divide-slate-800">
              <button
                v-for="item in messages"
                :key="item.id"
                type="button"
                class="flex w-full gap-3 px-4 py-3 text-left transition"
                :class="selectedMessageId === item.id ? 'bg-sky-50/80 dark:bg-sky-950/30' : 'hover:bg-slate-50 dark:hover:bg-slate-950/60'"
                @click="selectedMessageId = item.id"
              >
                <span class="mt-0.5 flex size-9 shrink-0 items-center justify-center rounded-md" :class="statusTone(item.status)">
                  <UIcon :name="statusIcon(item.status)" class="size-4" />
                </span>
                <span class="min-w-0 flex-1">
                  <span class="flex min-w-0 items-center justify-between gap-3">
                    <span class="truncate text-sm font-semibold text-slate-950 dark:text-white">{{ item.subject }}</span>
                    <span class="shrink-0 text-xs text-slate-500 dark:text-slate-400">{{ formatRelativeDateTime(item.createdAt, locale) }}</span>
                  </span>
                  <span class="mt-1 flex items-center gap-2 text-xs text-slate-500 dark:text-slate-400">
                    <span>{{ statusLabel(item.status) }}</span>
                    <span class="text-slate-300 dark:text-slate-700">/</span>
                    <span class="truncate">{{ messageExcerpt(item.content) }}</span>
                  </span>
                </span>
              </button>
            </div>
          </div>

          <article v-if="selectedMessage" class="min-w-0 p-5">
            <div class="flex flex-col gap-3 border-b border-slate-100 pb-4 sm:flex-row sm:items-start sm:justify-between dark:border-slate-800">
              <div class="min-w-0">
                <div class="flex items-center gap-2">
                  <h2 class="break-words text-base font-semibold text-slate-950 dark:text-white">{{ selectedMessage.subject }}</h2>
                  <UBadge :color="statusColor(selectedMessage.status)" variant="soft" class="shrink-0">{{ statusLabel(selectedMessage.status) }}</UBadge>
                </div>
                <p class="mt-1 text-xs text-slate-500 dark:text-slate-400">{{ formatDateTime(selectedMessage.createdAt, locale) }}</p>
              </div>
            </div>

            <div class="mt-5 space-y-5">
              <section>
                <p class="mb-2 text-xs font-medium text-slate-500 dark:text-slate-400">{{ $t('mod.staffMessages.detail.message') }}</p>
                <div class="rich-text rounded-md border border-slate-200 bg-slate-50 px-4 py-3 dark:border-slate-800 dark:bg-slate-950" v-html="renderUserMarkdown(selectedMessage.content)" />
              </section>
              <section v-if="selectedMessage.answer" class="border-l-2 border-emerald-400 pl-4 dark:border-emerald-500">
                <div class="flex items-center justify-between gap-3">
                  <p class="text-xs font-medium text-emerald-700 dark:text-emerald-300">{{ $t('mod.staffMessages.detail.reply') }}</p>
                  <span class="text-xs text-slate-500 dark:text-slate-400">{{ formatDateTime(selectedMessage.answeredAt, locale) }}</span>
                </div>
                <div class="rich-text mt-2 text-sm" v-html="renderUserMarkdown(selectedMessage.answer)" />
              </section>
              <p v-else class="text-sm text-slate-500 dark:text-slate-400">{{ selectedMessage.status === ModStaffMessageStatus.Processed ? $t('mod.staffMessages.detail.processedNoReply') : $t('mod.staffMessages.detail.waiting') }}</p>
            </div>
          </article>
        </div>

        <AppPager
          v-if="!pending && messages.length > 0"
          class="border-t border-slate-200 px-4 py-3 dark:border-slate-800"
          :page="page"
          :total="total"
          :page-size="size"
          :page-size-options="pageSizes"
          :disabled="pending"
          @page-change="changePage"
          @page-size-change="changeSize"
        />
      </section>
    </div>

    <UModal
      :open="composeOpen"
      :title="$t('mod.staffMessages.compose.title')"
      :description="$t('mod.staffMessages.compose.description')"
      :ui="{ content: 'sm:max-w-2xl', header: 'px-5 py-4', body: 'p-0 sm:p-0' }"
      @update:open="composeOpen = $event"
    >
      <template #body>
        <form class="space-y-4 px-5 pb-5 pt-5" @submit.prevent="submitMessage">
          <UFormField :label="$t('mod.staffMessages.compose.subject')" required>
            <UInput v-model="form.subject" class="w-full" :disabled="submitting" />
          </UFormField>
          <UFormField :label="$t('mod.staffMessages.compose.content')" required>
            <UTextarea v-model="form.content" class="w-full" :rows="8" :disabled="submitting" />
          </UFormField>
          <p v-if="formError" class="rounded-md border border-red-200 bg-red-50 px-3 py-2 text-sm text-red-700 dark:border-red-900 dark:bg-red-950/40 dark:text-red-200">{{ formError }}</p>
          <div class="flex justify-end gap-2">
            <UButton type="button" color="neutral" variant="ghost" :disabled="submitting" @click="composeOpen = false">{{ $t('common.cancel') }}</UButton>
            <UButton type="submit" color="primary" icon="i-lucide-send" :loading="submitting">{{ $t('mod.staffMessages.compose.submit') }}</UButton>
          </div>
        </form>
      </template>
    </UModal>
  </div>
</template>

<script setup lang="ts">
import { ApiError } from '~/composables/useApi'
import { ModStaffMessageStatus, type StaffMessage } from '~/composables/useMod'
import { Permission } from '~/composables/useAuth'
import { formatDateTime, formatRelativeDateTime } from '~/utils/format'
import { renderUserMarkdown } from '~/utils/richText'

definePageMeta({ middleware: 'auth' })

type MessageFilter = 'all' | 'pending' | 'processed'

const { t, locale } = useI18n()
const toast = useToast()
const modApi = useMod()

const messages = ref<StaffMessage[]>([])
const selectedMessageId = ref(0)
const total = ref(0)
const page = ref(1)
const size = ref(20)
const pageSizes = [20, 50, 100]
const filter = ref<MessageFilter>('all')
const pending = ref(false)
const errorMessage = ref('')
const composeOpen = ref(false)
const submitting = ref(false)
const formError = ref('')
const form = reactive({ subject: '', content: '' })

const selectedMessage = computed(() => messages.value.find(item => item.id === selectedMessageId.value) || null)
const filterOptions = computed(() => [
  { key: 'all' as const, label: t('mod.staffMessages.filters.all') },
  { key: 'pending' as const, label: t('mod.staffMessages.status.pending') },
  { key: 'processed' as const, label: t('mod.staffMessages.status.processed') }
])

useHead(() => ({ title: t('mod.staffMessages.title') }))
onMounted(loadMessages)

async function loadMessages() {
  pending.value = true
  errorMessage.value = ''
  try {
    const status = filter.value === 'all' ? undefined : statusForFilter(filter.value)
    const data = await modApi.listStaffMessages({ page: page.value, size: size.value, status })
    messages.value = data.list || []
    total.value = data.total || 0
    if (!messages.value.some(item => item.id === selectedMessageId.value)) selectedMessageId.value = messages.value[0]?.id || 0
  } catch (error) {
    messages.value = []
    total.value = 0
    selectedMessageId.value = 0
    errorMessage.value = error instanceof ApiError ? error.message : t('common.requestFailed')
  } finally {
    pending.value = false
  }
}

function statusForFilter(value: Exclude<MessageFilter, 'all'>) {
  return value === 'pending' ? ModStaffMessageStatus.Pending : ModStaffMessageStatus.Processed
}

function changeFilter(value: MessageFilter) {
  filter.value = value
  page.value = 1
  void loadMessages()
}

function changePage(value: number) {
  page.value = value
  void loadMessages()
}

function changeSize(value: number) {
  size.value = value
  page.value = 1
  void loadMessages()
}

async function submitMessage() {
  formError.value = ''
  if (!form.subject.trim() || !form.content.trim()) {
    formError.value = t('mod.staffMessages.compose.required')
    return
  }
  submitting.value = true
  try {
    await modApi.createStaffMessage({ subject: form.subject.trim(), content: form.content.trim() })
    toast.add({ title: t('mod.staffMessages.compose.sent'), color: 'success', icon: 'i-lucide-check' })
    form.subject = ''
    form.content = ''
    composeOpen.value = false
    filter.value = 'all'
    page.value = 1
    await loadMessages()
  } catch (error) {
    formError.value = error instanceof ApiError ? error.message : t('common.requestFailed')
  } finally {
    submitting.value = false
  }
}

function statusLabel(status: number) {
  if (status === ModStaffMessageStatus.Processed) return t('mod.staffMessages.status.processed')
  return t('mod.staffMessages.status.pending')
}

function statusColor(status: number) {
  if (status === ModStaffMessageStatus.Processed) return 'success' as const
  return 'warning' as const
}

function statusTone(status: number) {
  if (status === ModStaffMessageStatus.Processed) return 'bg-emerald-50 text-emerald-600 dark:bg-emerald-950/40 dark:text-emerald-300'
  return 'bg-amber-50 text-amber-600 dark:bg-amber-950/40 dark:text-amber-300'
}

function statusIcon(status: number) {
  return status === ModStaffMessageStatus.Processed ? 'i-lucide-mail-check' : 'i-lucide-mail'
}

function messageExcerpt(value: string) {
  const text = value.replace(/\s+/g, ' ').trim()
  return text.length > 64 ? `${text.slice(0, 64)}...` : text
}
</script>
