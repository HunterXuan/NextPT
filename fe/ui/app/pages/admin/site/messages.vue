<template>
  <div class="min-h-[calc(100vh-4rem)] bg-slate-50 py-6 dark:bg-slate-950">
    <div class="w-full px-3 sm:px-4 lg:px-5">
      <div class="grid gap-4 xl:grid-cols-[minmax(440px,560px)_1fr]">
        <section class="overflow-hidden rounded-lg border border-slate-200 bg-white dark:border-slate-800 dark:bg-slate-900">
          <div class="border-b border-slate-200 px-4 py-3 dark:border-slate-800">
            <h2 class="text-sm font-semibold text-slate-950 dark:text-white">{{ $t('admin.site.messages.send') }}</h2>
            <p class="mt-1 text-xs text-slate-500 dark:text-slate-400">{{ $t('admin.site.messages.sendHint') }}</p>
          </div>
          <form class="space-y-4 p-4" @submit.prevent="sendMessage">
            <label class="block">
              <span class="text-sm font-medium text-slate-700 dark:text-slate-200">{{ $t('admin.site.messages.fields.receiverIds') }}</span>
              <textarea v-model="form.receiverIds" rows="4" class="mt-1 w-full resize-y rounded-md border border-slate-200 bg-white px-3 py-2 font-mono text-sm text-slate-950 outline-none transition focus:border-sky-400 focus:ring-2 focus:ring-sky-100 dark:border-slate-700 dark:bg-slate-950 dark:text-white dark:focus:border-sky-500 dark:focus:ring-sky-950" :placeholder="$t('admin.site.messages.receiverPlaceholder')" />
            </label>
            <label class="block">
              <span class="text-sm font-medium text-slate-700 dark:text-slate-200">{{ $t('admin.site.messages.fields.title') }}</span>
              <input v-model.trim="form.title" class="mt-1 h-10 w-full rounded-md border border-slate-200 bg-white px-3 text-sm text-slate-950 outline-none transition focus:border-sky-400 focus:ring-2 focus:ring-sky-100 dark:border-slate-700 dark:bg-slate-950 dark:text-white dark:focus:border-sky-500 dark:focus:ring-sky-950">
            </label>
            <label class="block">
              <span class="text-sm font-medium text-slate-700 dark:text-slate-200">{{ $t('admin.site.messages.fields.content') }}</span>
              <textarea v-model="form.content" rows="8" class="mt-1 w-full resize-y rounded-md border border-slate-200 bg-white px-3 py-2 text-sm leading-6 text-slate-950 outline-none transition focus:border-sky-400 focus:ring-2 focus:ring-sky-100 dark:border-slate-700 dark:bg-slate-950 dark:text-white dark:focus:border-sky-500 dark:focus:ring-sky-950" />
            </label>
            <div class="grid gap-3 sm:grid-cols-2">
              <label class="block">
                <span class="text-sm font-medium text-slate-700 dark:text-slate-200">{{ $t('admin.site.messages.fields.targetType') }}</span>
                <select v-model="form.targetType" class="mt-1 h-10 w-full rounded-md border border-slate-200 bg-white px-3 text-sm text-slate-950 outline-none transition focus:border-sky-400 focus:ring-2 focus:ring-sky-100 dark:border-slate-700 dark:bg-slate-950 dark:text-white dark:focus:border-sky-500 dark:focus:ring-sky-950">
                  <option value="">{{ $t('admin.site.messages.targetTypes.none') }}</option>
                  <option value="catalog_torrent">{{ $t('admin.site.messages.targetTypes.catalogTorrent') }}</option>
                  <option value="forum_topic">{{ $t('admin.site.messages.targetTypes.forumTopic') }}</option>
                </select>
              </label>
              <label class="block">
                <span class="text-sm font-medium text-slate-700 dark:text-slate-200">{{ $t('admin.site.messages.fields.targetId') }}</span>
                <input v-model.number="form.targetId" type="number" min="0" class="mt-1 h-10 w-full rounded-md border border-slate-200 bg-white px-3 text-sm text-slate-950 outline-none transition focus:border-sky-400 focus:ring-2 focus:ring-sky-100 dark:border-slate-700 dark:bg-slate-950 dark:text-white dark:focus:border-sky-500 dark:focus:ring-sky-950">
              </label>
            </div>
            <p v-if="formError" class="rounded-md border border-red-200 bg-red-50 px-3 py-2 text-sm text-red-700 dark:border-red-900 dark:bg-red-950/40 dark:text-red-200">{{ formError }}</p>
            <div class="flex justify-end border-t border-slate-200 pt-4 dark:border-slate-800">
              <UButton type="submit" color="primary" icon="i-lucide-send" :loading="sending">{{ $t('admin.site.messages.submit') }}</UButton>
            </div>
          </form>
        </section>

        <section class="overflow-hidden rounded-lg border border-slate-200 bg-white dark:border-slate-800 dark:bg-slate-900">
          <div class="flex flex-col gap-3 border-b border-slate-200 px-4 py-3 sm:flex-row sm:items-center sm:justify-between dark:border-slate-800">
            <h2 class="text-sm font-semibold text-slate-950 dark:text-white">{{ $t('admin.site.messages.list') }}</h2>
            <div class="flex items-center gap-2">
              <input v-model.trim="receiverFilter" inputmode="numeric" class="h-9 w-32 rounded-md border border-slate-200 bg-white px-3 text-sm outline-none transition focus:border-sky-300 dark:border-slate-700 dark:bg-slate-950 dark:focus:border-sky-700" :placeholder="$t('admin.site.messages.filters.receiverId')">
              <select v-model="readFilter" class="h-9 w-32 rounded-md border border-slate-200 bg-white px-3 text-sm outline-none transition focus:border-sky-300 dark:border-slate-700 dark:bg-slate-950 dark:focus:border-sky-700">
                <option value="">{{ $t('admin.site.messages.filters.all') }}</option>
                <option value="false">{{ $t('admin.site.messages.filters.unread') }}</option>
                <option value="true">{{ $t('admin.site.messages.filters.read') }}</option>
              </select>
              <UButton size="sm" color="neutral" variant="outline" @click="reloadFromFirstPage">{{ $t('admin.site.messages.filters.apply') }}</UButton>
            </div>
          </div>

          <div v-if="pending" class="space-y-2 p-4">
            <div v-for="item in 8" :key="item" class="h-20 animate-pulse rounded-md bg-slate-100 dark:bg-slate-800" />
          </div>
          <div v-else-if="errorMessage" class="flex flex-col items-center justify-center px-4 py-16 text-center">
            <UIcon name="i-lucide-circle-alert" class="size-9 text-red-500" />
            <p class="mt-3 text-sm font-medium text-slate-950 dark:text-white">{{ errorMessage }}</p>
          </div>
          <div v-else-if="messages.length === 0" class="flex flex-col items-center justify-center px-4 py-16 text-center">
            <UIcon name="i-lucide-inbox" class="size-9 text-slate-400" />
            <p class="mt-3 text-sm font-medium text-slate-950 dark:text-white">{{ $t('admin.site.messages.empty') }}</p>
          </div>
          <div v-else class="divide-y divide-slate-200 dark:divide-slate-800">
            <article v-for="item in messages" :key="item.id" class="p-4">
              <div class="flex items-start justify-between gap-3">
                <div class="min-w-0">
                  <p class="truncate text-sm font-semibold text-slate-950 dark:text-white">{{ item.title }}</p>
                  <p class="mt-1 text-xs text-slate-500 dark:text-slate-400">
                    {{ $t('admin.site.messages.toUser', { id: item.receiverId, name: item.receiver?.username || '-' }) }}
                    <span class="mx-2 text-slate-300 dark:text-slate-700">/</span>
                    {{ formatDateTime(item.createdAt, locale) }}
                  </p>
                </div>
                <UBadge :color="item.isRead ? 'neutral' : 'warning'" variant="soft">{{ item.isRead ? $t('admin.site.messages.readState.read') : $t('admin.site.messages.readState.unread') }}</UBadge>
              </div>
              <p class="mt-2 whitespace-pre-line text-sm leading-6 text-slate-600 dark:text-slate-300">{{ item.content }}</p>
            </article>
          </div>

          <AppPager
            class="border-t border-slate-200 px-4 py-3 dark:border-slate-800"
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
  </div>
</template>

<script setup lang="ts">
import { ApiError } from '~/composables/useApi'
import type { AdminSiteMessage } from '~/composables/useAdmin'
import { formatDateTime } from '~/utils/format'

definePageMeta({ layout: 'admin', middleware: 'admin' })

const { t, locale } = useI18n()
const toast = useToast()
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
const form = reactive({ receiverIds: '', title: '', content: '', targetType: '', targetId: 0 })

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
  sending.value = true
  try {
    const data = await adminApi.createSiteMessage({
      receiverIds,
      title: form.title,
      content: form.content,
      targetType: form.targetType,
      targetId: Number(form.targetId) || 0
    })
    toast.add({ color: 'success', title: t('admin.site.messages.sent', { count: data.count || 0 }), icon: 'i-lucide-check' })
    form.receiverIds = ''
    form.title = ''
    form.content = ''
    form.targetType = ''
    form.targetId = 0
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

function reloadFromFirstPage() {
  page.value = 1
  loadMessages()
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
</script>
