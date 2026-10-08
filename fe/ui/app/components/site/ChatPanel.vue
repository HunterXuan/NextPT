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
        <div class="relative min-h-0 flex-1">
          <div ref="messageList" class="h-full overflow-y-auto px-4 py-3 sm:px-4" @scroll="handleScroll">
            <div v-if="loading && messages.length === 0" class="space-y-3">
              <div v-for="index in 5" :key="index" class="flex gap-2.5">
                <div class="size-7 shrink-0 animate-pulse rounded-md bg-neutral-200 dark:bg-neutral-800" />
                <div class="min-w-0 flex-1 space-y-1.5 pt-0.5">
                  <div class="h-3 w-24 animate-pulse rounded bg-neutral-200 dark:bg-neutral-800" />
                  <div class="h-3 w-4/5 animate-pulse rounded bg-neutral-100 dark:bg-neutral-800/70" />
                </div>
              </div>
            </div>

            <p v-else-if="errorMessage && messages.length === 0" class="rounded-md border border-red-200 bg-red-50 px-3 py-2 text-sm text-red-700 dark:border-red-900 dark:bg-red-950/40 dark:text-red-200">
              {{ errorMessage }}
            </p>

            <p v-else-if="messages.length === 0" class="py-12 text-center text-sm text-neutral-500 dark:text-neutral-400">
              {{ $t('site.chat.empty') }}
            </p>

            <ol v-else class="space-y-3">
              <li v-for="message in displayMessages" :key="message.id" class="flex gap-2.5" :class="isOwnMessage(message) ? 'justify-end' : ''">
                <IamUserAvatar v-if="!isOwnMessage(message)" :user="message.user" size="xs" />
                <div class="inline-flex min-w-0 max-w-[82%] flex-col" :class="isOwnMessage(message) ? 'items-end text-right' : 'items-start'">
                  <div class="inline-flex min-w-0 items-baseline gap-2" :class="isOwnMessage(message) ? 'justify-end' : ''">
                    <IamUserPopover :user="message.user" class="min-w-0 text-sm font-semibold text-neutral-800 dark:text-neutral-100" />
                    <time class="shrink-0 text-[11px] text-neutral-400 dark:text-neutral-500" :title="formatDateTime(message.createdAt, locale)">
                      {{ formatRelativeDateTime(message.createdAt, locale) }}
                    </time>
                  </div>
                  <UContextMenu
                    :items="[{ label: $t('site.chat.reply'), icon: 'i-lucide-reply', disabled: sending || !message.body.trim(), onSelect: () => replyToMessage(message) }]"
                    :disabled="!canSend"
                    :modal="false"
                    class="mt-0.5 min-w-0 max-w-full"
                  >
                    <div
                      class="min-w-0 max-w-full break-words rounded-md px-2.5 py-1.5 text-left text-sm leading-5 [overflow-wrap:anywhere]"
                      :class="isOwnMessage(message) ? 'bg-primary-50 text-neutral-800 dark:bg-primary-950/50 dark:text-neutral-100' : 'bg-neutral-50 text-neutral-700 dark:bg-neutral-900/70 dark:text-neutral-300'"
                    >
                      <blockquote v-if="message.quote" class="mb-1.5 whitespace-pre-wrap border-l-2 border-primary/40 pl-2 text-xs leading-5 text-muted">{{ message.quote }}</blockquote>
                      <div v-if="message.body" class="chat-message-content rich-text rich-text-compact text-sm leading-5" v-html="renderUserMarkdown(message.body)" />
                    </div>
                  </UContextMenu>
                </div>
                <IamUserAvatar v-if="isOwnMessage(message)" :user="message.user" size="xs" />
              </li>
            </ol>
          </div>
          <UTooltip v-if="hasNewMessages" :text="$t('site.chat.newMessages')">
            <UButton class="absolute bottom-3 end-3 size-8 justify-center p-0 shadow-sm" color="primary" icon="i-lucide-arrow-down" :aria-label="$t('site.chat.newMessages')" @click="scrollToLatest" />
          </UTooltip>
        </div>

        <form class="border-t border-neutral-200 p-3 dark:border-neutral-800" @submit.prevent="sendMessage">
          <div v-if="quotedText" class="mb-2 flex min-w-0 items-center gap-2 border-l-2 border-primary/40 bg-elevated px-2 py-1.5">
            <p class="min-w-0 flex-1 line-clamp-2 break-words text-xs leading-5 text-muted [overflow-wrap:anywhere]">{{ quotedText }}</p>
            <UTooltip :text="$t('site.chat.cancelQuote')">
              <UButton color="neutral" variant="ghost" size="xs" icon="i-lucide-x" class="size-6 shrink-0 p-0" :aria-label="$t('site.chat.cancelQuote')" :disabled="sending" @click="clearQuote" />
            </UTooltip>
          </div>
          <UTextarea
            ref="composer"
            v-model="draft"
            class="w-full"
            :rows="2"
            :placeholder="canSend ? $t('site.chat.placeholder') : $t('site.chat.readOnly')"
            :disabled="!canSend || sending"
            :maxlength="draftMaxLength"
            autoresize
            @keydown.enter.exact.prevent="sendMessage"
          />
          <p v-if="errorMessage && messages.length > 0" role="alert" class="mt-2 text-xs text-error">{{ errorMessage }}</p>
          <div class="mt-2 flex items-center justify-between gap-3">
            <p class="text-xs" :class="draftLength > 1000 ? 'text-error' : 'text-neutral-400 dark:text-neutral-500'">{{ draftLength }}/1000</p>
            <div class="flex items-center gap-2">
              <UPopover :content="{ side: 'top', align: 'end' }" :ui="{ content: 'p-3' }">
                <UTooltip :text="$t('site.chat.notificationSettings')">
                  <UButton type="button" color="neutral" variant="soft" icon="i-lucide-settings-2" size="sm" class="size-8 justify-center rounded-md p-0" :ui="{ leadingIcon: 'size-4' }" :aria-label="$t('site.chat.notificationSettings')" />
                </UTooltip>
                <template #content>
                  <USwitch v-model="notificationsEnabled" :label="$t('site.chat.notifyNewMessages')" />
                </template>
              </UPopover>
              <UPopover v-if="canSend" v-model:open="imagePopoverOpen" :content="{ side: 'top', align: 'end' }" :ui="{ content: 'w-[min(20rem,calc(100vw-3rem))] p-3' }">
                <UTooltip :text="$t('site.chat.insertImage')" :disabled="imagePopoverOpen">
                  <UButton type="button" color="neutral" variant="soft" icon="i-lucide-image" size="sm" class="size-8 justify-center rounded-md p-0" :ui="{ leadingIcon: 'size-4' }" :aria-label="$t('site.chat.insertImage')" :disabled="sending" />
                </UTooltip>
                <template #content>
                  <div class="space-y-2">
                    <UFormField :label="$t('site.chat.imageUrl')" :error="imageError || undefined">
                      <UInput v-model="imageUrl" class="w-full" type="url" placeholder="https://" @keydown.enter.prevent="insertImage" />
                    </UFormField>
                    <div class="flex justify-end gap-2">
                      <UButton type="button" color="neutral" variant="soft" size="sm" @click="cancelImage">{{ $t('common.cancel') }}</UButton>
                      <UButton type="button" size="sm" :disabled="!imageUrl.trim() || sending" @click="insertImage">{{ $t('common.confirm') }}</UButton>
                    </div>
                  </div>
                </template>
              </UPopover>
              <UTooltip :text="canSend ? $t('site.chat.send') : $t('common.noPermission')">
                <UButton
                  type="submit"
                  color="primary"
                  size="sm"
                  icon="i-lucide-send"
                  class="size-8 justify-center rounded-md p-0"
                  :ui="{ leadingIcon: 'size-4' }"
                  :aria-label="$t('site.chat.send')"
                  :loading="sending"
                  :disabled="!canSend || !draft.trim() || draftLength > 1000 || sending"
                />
              </UTooltip>
            </div>
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
import { composeChatMessage, countNewChatMessages, createChatImage, createChatQuote, splitChatMessage } from '~/utils/chat'
import { renderUserMarkdown } from '~/utils/richText'

const props = withDefaults(defineProps<{
  open: boolean
  canSend?: boolean
}>(), {
  canSend: false
})

const emit = defineEmits<{
  'update:open': [value: boolean]
  'update:unreadCount': [value: number]
}>()

const { t, locale } = useI18n()
const siteApi = useSite()
const { user } = useAuth()
const messages = ref<SiteChatMessage[]>([])
const draft = ref('')
const quotedText = ref('')
const imagePopoverOpen = ref(false)
const imageUrl = ref('')
const imageError = ref('')
const composer = useTemplateRef('composer')
const hasNewMessages = ref(false)
const loading = ref(false)
const sending = ref(false)
const errorMessage = ref('')
const messageList = ref<HTMLElement | null>(null)
const notificationsEnabled = ref(true)
const unreadCount = ref(0)
const notificationStorageKey = 'nextpt_chat_notifications'
let lastMessageId: number | null = null
let pollTimer: ReturnType<typeof setInterval> | undefined

const modalUi = {
  content: 'w-[calc(100vw-2rem)] max-w-[34rem] sm:max-w-[34rem]',
  header: 'min-h-0 px-4 py-3 sm:px-4 sm:py-3',
  title: 'text-sm font-semibold',
  description: 'mt-0.5 text-xs',
  close: 'top-2 end-2',
  body: 'p-0 sm:p-0'
}
const displayMessages = computed(() => messages.value.map(message => ({ ...message, ...splitChatMessage(message.content) })))
const draftLength = computed(() => Array.from(composeChatMessage(draft.value, quotedText.value)).length)
const draftMaxLength = computed(() => Math.max(0, 1000 - (quotedText.value ? Array.from(`> ${quotedText.value}\n\n`).length : 0)))

async function replyToMessage(message: SiteChatMessage) {
  if (!props.canSend || sending.value) return
  quotedText.value = createChatQuote(message.user.username || `#${message.user.id}`, message.content, t('site.chat.imageSummary'))
  await nextTick()
  composer.value?.textareaRef?.focus()
}

async function cancelImage() {
  imagePopoverOpen.value = false
  imageUrl.value = ''
  imageError.value = ''
  await nextTick()
  composer.value?.textareaRef?.focus()
}

async function insertImage() {
  if (!props.canSend || sending.value) return
  imageError.value = ''
  const image = createChatImage(imageUrl.value)
  if (!image) {
    imageError.value = t('site.chat.invalidImageUrl')
    return
  }
  const content = [draft.value.trimEnd(), image].filter(Boolean).join('\n')
  if (Array.from(composeChatMessage(content, quotedText.value)).length > 1000) {
    imageError.value = t('site.chat.messageTooLong')
    return
  }
  draft.value = content
  imageUrl.value = ''
  imagePopoverOpen.value = false
  await nextTick()
  composer.value?.textareaRef?.focus()
}

async function clearQuote() {
  quotedText.value = ''
  await nextTick()
  composer.value?.textareaRef?.focus()
}

watch(() => props.open, (value) => {
  if (value) {
    void loadInitialMessages()
  }
  startPolling()
})

watch(unreadCount, value => emit('update:unreadCount', value), { immediate: true })

watch(notificationsEnabled, (value) => {
  unreadCount.value = 0
  lastMessageId = null
  try {
    localStorage.setItem(notificationStorageKey, value ? '1' : '0')
  } catch {
    // Keep the preference usable when browser storage is unavailable.
  }
  if (value) void pollMessages()
  startPolling()
})

onMounted(() => {
  try {
    notificationsEnabled.value = localStorage.getItem(notificationStorageKey) !== '0'
  } catch {
    // Default to reminders when browser storage is unavailable.
  }
  if (props.open) void loadInitialMessages()
  else void pollMessages()
  startPolling()
})

onBeforeUnmount(stopPolling)

async function loadInitialMessages() {
  if (loading.value) {
    await scrollToLatest()
    return
  }
  loading.value = true
  errorMessage.value = ''
  try {
    const data = await siteApi.listChatMessages()
    messages.value = data.list || []
    lastMessageId = Math.max(lastMessageId || 0, ...messages.value.map(message => message.id))
    unreadCount.value = 0
    await scrollToLatest()
  } catch (error) {
    errorMessage.value = error instanceof ApiError ? error.message : t('common.requestFailed')
  } finally {
    loading.value = false
  }
}

async function pollMessages() {
  if ((!props.open && !notificationsEnabled.value) || loading.value || sending.value) return
  const afterId = lastMessageId || 0
  loading.value = true
  try {
    const data = await siteApi.listChatMessages(afterId)
    const items = data.list || []
    if (props.open) unreadCount.value = 0
    else if (notificationsEnabled.value && lastMessageId !== null) {
      unreadCount.value += countNewChatMessages(items, lastMessageId, Number(user.value?.user.id || 0))
    }
    lastMessageId = Math.max(afterId, ...items.map(message => message.id))
    appendMessages(items)
  } catch {
    // Keep the current conversation visible and retry on the next interval.
  } finally {
    loading.value = false
  }
}

async function sendMessage() {
  if (!props.canSend || !draft.value.trim() || draftLength.value > 1000 || sending.value) return
  const content = composeChatMessage(draft.value, quotedText.value)
  sending.value = true
  errorMessage.value = ''
  try {
    const message = await siteApi.createChatMessage(content)
    draft.value = ''
    quotedText.value = ''
    appendMessages([message], true)
  } catch (error) {
    errorMessage.value = error instanceof ApiError ? error.message : t('common.requestFailed')
  } finally {
    sending.value = false
  }
}

function appendMessages(items: SiteChatMessage[], forceScroll = false) {
  if (items.length === 0) return
  const ids = new Set(messages.value.map(item => item.id))
  const appended = items.filter(item => !ids.has(item.id))
  if (appended.length === 0) {
    if (forceScroll) void scrollToLatest()
    return
  }
  const followLatest = forceScroll || isAtBottom()
  messages.value = [...messages.value, ...appended].slice(-200)
  if (followLatest) void scrollToLatest()
  else hasNewMessages.value = true
}

function isAtBottom() {
  const element = messageList.value
  return !element || element.scrollHeight - element.scrollTop - element.clientHeight < 48
}

function handleScroll() {
  if (isAtBottom()) hasNewMessages.value = false
}

function isOwnMessage(message: SiteChatMessage) {
  return Number(message.user.id) === Number(user.value?.user.id || 0)
}

function startPolling() {
  stopPolling()
  if (!props.open && !notificationsEnabled.value) return
  pollTimer = setInterval(() => { void pollMessages() }, props.open ? 8000 : 60000)
}

function stopPolling() {
  if (!pollTimer) return
  clearInterval(pollTimer)
  pollTimer = undefined
}

async function scrollToLatest() {
  await nextTick()
  if (messageList.value) messageList.value.scrollTop = messageList.value.scrollHeight
  hasNewMessages.value = false
}
</script>

<style scoped>
.chat-message-content :deep(img) {
  max-height: 12rem;
  max-width: 100%;
  object-fit: contain;
}
</style>
