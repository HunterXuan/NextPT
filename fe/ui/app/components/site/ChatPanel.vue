<template>
  <UModal
    :open="open"
    :title="$t('site.chat.title')"
    :description="$t('site.chat.description')"
    :ui="modalUi"
    @update:open="emit('update:open', $event)"
  >
    <template #body>
      <div class="flex h-[min(68vh,32rem)] min-h-0 flex-col">
        <div ref="messageList" class="min-h-0 flex-1 overflow-y-auto px-4 py-3 sm:px-4">
          <div v-if="loading && messages.length === 0" class="space-y-3">
            <div v-for="index in 5" :key="index" class="flex gap-2.5">
              <div class="size-7 shrink-0 animate-pulse rounded-md bg-neutral-200 dark:bg-neutral-800" />
              <div class="min-w-0 flex-1 space-y-1.5 pt-0.5">
                <div class="h-3 w-24 animate-pulse rounded bg-neutral-200 dark:bg-neutral-800" />
                <div class="h-3 w-4/5 animate-pulse rounded bg-neutral-100 dark:bg-neutral-800/70" />
              </div>
            </div>
          </div>

          <p v-else-if="errorMessage" class="rounded-md border border-red-200 bg-red-50 px-3 py-2 text-sm text-red-700 dark:border-red-900 dark:bg-red-950/40 dark:text-red-200">
            {{ errorMessage }}
          </p>

          <p v-else-if="messages.length === 0" class="py-12 text-center text-sm text-neutral-500 dark:text-neutral-400">
            {{ $t('site.chat.empty') }}
          </p>

          <ol v-else class="space-y-3">
            <li v-for="message in messages" :key="message.id" class="flex gap-2.5" :class="isOwnMessage(message) ? 'justify-end' : ''">
              <IamUserAvatar v-if="!isOwnMessage(message)" :user="message.user" size="xs" />
              <div class="inline-flex min-w-0 max-w-[82%] flex-col" :class="isOwnMessage(message) ? 'items-end text-right' : 'items-start'">
                <div class="inline-flex min-w-0 items-baseline gap-2" :class="isOwnMessage(message) ? 'justify-end' : ''">
                  <IamUserPopover :user="message.user" class="min-w-0 text-sm font-semibold text-neutral-800 dark:text-neutral-100" />
                  <time class="shrink-0 text-[11px] text-neutral-400 dark:text-neutral-500" :title="formatDateTime(message.createdAt, locale)">
                    {{ formatRelativeDateTime(message.createdAt, locale) }}
                  </time>
                </div>
                <p
                  class="mt-0.5 inline-block max-w-full whitespace-pre-wrap break-words rounded-md px-2.5 py-1.5 text-left text-sm leading-5"
                  :class="isOwnMessage(message) ? 'bg-primary-50 text-neutral-800 dark:bg-primary-950/50 dark:text-neutral-100' : 'bg-neutral-50 text-neutral-700 dark:bg-neutral-900/70 dark:text-neutral-300'"
                >{{ message.content }}</p>
              </div>
              <IamUserAvatar v-if="isOwnMessage(message)" :user="message.user" size="xs" />
            </li>
          </ol>
        </div>

        <form class="border-t border-neutral-200 p-3 dark:border-neutral-800" @submit.prevent="sendMessage">
          <UTextarea
            v-model="draft"
            class="w-full"
            :rows="2"
            :placeholder="canSend ? $t('site.chat.placeholder') : $t('site.chat.readOnly')"
            :disabled="!canSend || sending"
            :maxlength="1000"
            autoresize
            @keydown.enter.exact.prevent="sendMessage"
          />
          <div class="mt-2 flex items-center justify-between gap-3">
            <p class="text-xs text-neutral-400 dark:text-neutral-500">{{ draftLength }}/1000</p>
            <UTooltip :text="canSend ? $t('site.chat.send') : $t('common.noPermission')">
              <UButton
                type="submit"
                color="primary"
                size="sm"
                icon="i-lucide-send"
                :aria-label="$t('site.chat.send')"
                :loading="sending"
                :disabled="!canSend || !draft.trim() || sending"
              />
            </UTooltip>
          </div>
        </form>
      </div>
    </template>
  </UModal>
</template>

<script setup lang="ts">
import { ApiError } from '~/composables/useApi'
import type { SiteChatMessage } from '~/composables/useSite'
import { formatDateTime, formatRelativeDateTime } from '~/utils/format'

const props = withDefaults(defineProps<{
  open: boolean
  canSend?: boolean
}>(), {
  canSend: false
})

const emit = defineEmits<{
  'update:open': [value: boolean]
}>()

const { t, locale } = useI18n()
const siteApi = useSite()
const { user } = useAuth()
const messages = ref<SiteChatMessage[]>([])
const draft = ref('')
const loading = ref(false)
const sending = ref(false)
const errorMessage = ref('')
const messageList = ref<HTMLElement | null>(null)
let pollTimer: ReturnType<typeof setInterval> | undefined

const modalUi = {
  content: 'w-[calc(100vw-2rem)] max-w-[34rem] sm:max-w-[34rem]',
  header: 'min-h-0 px-4 py-3 sm:px-4 sm:py-3',
  title: 'text-sm font-semibold',
  description: 'mt-0.5 text-xs',
  close: 'top-2 end-2',
  body: 'p-0 sm:p-0'
}
const draftLength = computed(() => Array.from(draft.value).length)

watch(() => props.open, (value) => {
  if (value) {
    void loadInitialMessages()
    startPolling()
    return
  }
  stopPolling()
}, { immediate: true })

onBeforeUnmount(stopPolling)

async function loadInitialMessages() {
  if (loading.value) return
  loading.value = true
  errorMessage.value = ''
  try {
    const data = await siteApi.listChatMessages()
    messages.value = data.list || []
    await scrollToLatest()
  } catch (error) {
    errorMessage.value = error instanceof ApiError ? error.message : t('common.requestFailed')
  } finally {
    loading.value = false
  }
}

async function pollMessages() {
  if (!props.open || loading.value || sending.value) return
  const afterId = messages.value.at(-1)?.id || 0
  try {
    const data = await siteApi.listChatMessages(afterId)
    appendMessages(data.list || [])
  } catch {
    // Keep the current conversation visible and retry on the next interval.
  }
}

async function sendMessage() {
  const content = draft.value.trim()
  if (!props.canSend || !content || sending.value) return
  sending.value = true
  errorMessage.value = ''
  try {
    const message = await siteApi.createChatMessage(content)
    draft.value = ''
    appendMessages([message])
  } catch (error) {
    errorMessage.value = error instanceof ApiError ? error.message : t('common.requestFailed')
  } finally {
    sending.value = false
  }
}

function appendMessages(items: SiteChatMessage[]) {
  if (items.length === 0) return
  const ids = new Set(messages.value.map(item => item.id))
  const appended = items.filter(item => !ids.has(item.id))
  if (appended.length === 0) return
  messages.value = [...messages.value, ...appended].slice(-200)
  void scrollToLatest()
}

function isOwnMessage(message: SiteChatMessage) {
  return Number(message.user.id) === Number(user.value?.user.id || 0)
}

function startPolling() {
  stopPolling()
  pollTimer = setInterval(() => { void pollMessages() }, 8000)
}

function stopPolling() {
  if (!pollTimer) return
  clearInterval(pollTimer)
  pollTimer = undefined
}

async function scrollToLatest() {
  await nextTick()
  if (messageList.value) messageList.value.scrollTop = messageList.value.scrollHeight
}
</script>
