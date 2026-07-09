<template>
  <div class="min-h-[calc(100vh-4rem)] bg-slate-50 py-6 dark:bg-slate-950">
    <div class="w-full px-3 sm:px-4 lg:px-5">
      <div class="grid gap-4 xl:grid-cols-[minmax(620px,760px)_minmax(460px,1fr)]">
        <section class="overflow-hidden rounded-lg border border-slate-200 bg-white dark:border-slate-800 dark:bg-slate-900">
          <div class="flex items-center justify-between gap-3 border-b border-slate-200 px-4 py-3 dark:border-slate-800">
            <div class="flex items-center gap-2">
              <UIcon name="i-lucide-megaphone" class="size-5 text-sky-600 dark:text-sky-300" />
              <h2 class="text-sm font-semibold text-slate-950 dark:text-white">{{ $t('admin.site.announcements.list') }}</h2>
            </div>
            <UButton size="sm" icon="i-lucide-plus" @click="startCreate">{{ $t('admin.site.announcements.create') }}</UButton>
          </div>

          <div class="border-b border-slate-200 px-4 py-3 dark:border-slate-800">
            <select v-model.number="statusFilter" class="h-9 w-40 rounded-md border border-slate-200 bg-white px-3 text-sm outline-none transition focus:border-sky-300 dark:border-slate-700 dark:bg-slate-950 dark:focus:border-sky-700" @change="reloadFromFirstPage">
              <option :value="-1">{{ $t('admin.site.announcements.status.all') }}</option>
              <option :value="0">{{ $t('admin.site.announcements.status.draft') }}</option>
              <option :value="1">{{ $t('admin.site.announcements.status.published') }}</option>
              <option :value="2">{{ $t('admin.site.announcements.status.archived') }}</option>
            </select>
          </div>

          <div v-if="pending" class="space-y-2 p-4">
            <div v-for="item in 6" :key="item" class="h-16 animate-pulse rounded-md bg-slate-100 dark:bg-slate-800" />
          </div>
          <div v-else-if="errorMessage" class="flex flex-col items-center justify-center px-4 py-16 text-center">
            <UIcon name="i-lucide-circle-alert" class="size-9 text-red-500" />
            <p class="mt-3 text-sm font-medium text-slate-950 dark:text-white">{{ errorMessage }}</p>
          </div>
          <div v-else-if="announcements.length === 0" class="flex flex-col items-center justify-center px-4 py-16 text-center">
            <UIcon name="i-lucide-inbox" class="size-9 text-slate-400" />
            <p class="mt-3 text-sm font-medium text-slate-950 dark:text-white">{{ $t('admin.site.announcements.empty') }}</p>
          </div>
          <div v-else class="divide-y divide-slate-200 dark:divide-slate-800">
            <button
              v-for="item in announcements"
              :key="item.id"
              type="button"
              class="flex w-full items-center justify-between gap-3 px-4 py-3 text-left transition"
              :class="selected?.id === item.id ? 'bg-sky-50/70 dark:bg-sky-950/30' : 'hover:bg-slate-50 dark:hover:bg-slate-950/70'"
              @click="selectAnnouncement(item)"
            >
              <div class="min-w-0">
                <p class="truncate text-sm font-semibold text-slate-950 dark:text-white">{{ item.title }}</p>
                <p class="mt-1 text-xs text-slate-500 dark:text-slate-400">{{ formatDateTime(item.publishedAt || item.createdAt, locale) }}</p>
              </div>
              <UBadge :color="statusColor(item.status)" variant="soft">{{ statusLabel(item.status) }}</UBadge>
            </button>
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

        <section class="overflow-hidden rounded-lg border border-slate-200 bg-white dark:border-slate-800 dark:bg-slate-900">
          <div class="flex items-start justify-between gap-3 border-b border-slate-200 px-4 py-3 dark:border-slate-800">
            <div class="min-w-0">
              <h2 class="truncate text-sm font-semibold text-slate-950 dark:text-white">{{ form.id ? $t('admin.site.announcements.edit') : $t('admin.site.announcements.create') }}</h2>
              <p class="mt-1 text-xs text-slate-500 dark:text-slate-400">{{ $t('admin.site.announcements.formHint') }}</p>
            </div>
            <UBadge v-if="form.id" color="neutral" variant="soft">#{{ form.id }}</UBadge>
          </div>

          <form class="space-y-4 p-4" @submit.prevent="saveAnnouncement">
            <label class="block">
              <span class="text-sm font-medium text-slate-700 dark:text-slate-200">{{ $t('admin.site.announcements.fields.title') }}</span>
              <input v-model.trim="form.title" class="mt-1 h-10 w-full rounded-md border border-slate-200 bg-white px-3 text-sm text-slate-950 outline-none transition focus:border-sky-400 focus:ring-2 focus:ring-sky-100 dark:border-slate-700 dark:bg-slate-950 dark:text-white dark:focus:border-sky-500 dark:focus:ring-sky-950">
            </label>
            <div class="grid gap-3 sm:grid-cols-2">
              <label class="block">
                <span class="text-sm font-medium text-slate-700 dark:text-slate-200">{{ $t('admin.site.announcements.fields.status') }}</span>
                <select v-model.number="form.status" class="mt-1 h-10 w-full rounded-md border border-slate-200 bg-white px-3 text-sm text-slate-950 outline-none transition focus:border-sky-400 focus:ring-2 focus:ring-sky-100 dark:border-slate-700 dark:bg-slate-950 dark:text-white dark:focus:border-sky-500 dark:focus:ring-sky-950">
                  <option :value="0">{{ $t('admin.site.announcements.status.draft') }}</option>
                  <option :value="1">{{ $t('admin.site.announcements.status.published') }}</option>
                  <option :value="2">{{ $t('admin.site.announcements.status.archived') }}</option>
                </select>
              </label>
              <label class="block">
                <span class="text-sm font-medium text-slate-700 dark:text-slate-200">{{ $t('admin.site.announcements.fields.publishedAt') }}</span>
                <input v-model="form.publishedAt" type="datetime-local" class="mt-1 h-10 w-full rounded-md border border-slate-200 bg-white px-3 text-sm text-slate-950 outline-none transition focus:border-sky-400 focus:ring-2 focus:ring-sky-100 dark:border-slate-700 dark:bg-slate-950 dark:text-white dark:focus:border-sky-500 dark:focus:ring-sky-950">
              </label>
            </div>
            <div>
              <span class="text-sm font-medium text-slate-700 dark:text-slate-200">{{ $t('admin.site.announcements.fields.content') }}</span>
              <RichTextComposer
                v-model="form.content"
                v-model:mode="editorMode"
                class="mt-1"
                :as-form="false"
                :rows="14"
                :disabled="saving"
                :submit-label="$t('common.save')"
                :write-label="$t('common.editor.edit')"
                :preview-label="$t('common.editor.preview')"
                :preview-empty="$t('common.editor.previewEmpty')"
                :show-actions="false"
              />
            </div>
            <p v-if="formError" class="rounded-md border border-red-200 bg-red-50 px-3 py-2 text-sm text-red-700 dark:border-red-900 dark:bg-red-950/40 dark:text-red-200">{{ formError }}</p>
            <div class="flex items-center justify-between gap-3 border-t border-slate-200 pt-4 dark:border-slate-800">
              <UButton v-if="form.id" color="error" variant="ghost" icon="i-lucide-trash-2" :disabled="saving" @click="deleteAnnouncement">{{ $t('common.delete') }}</UButton>
              <span v-else />
              <div class="flex gap-2">
                <UButton type="button" color="neutral" variant="outline" :disabled="saving" @click="resetForm">{{ $t('admin.actions.reset') }}</UButton>
                <UButton type="submit" color="primary" icon="i-lucide-save" :loading="saving">{{ $t('common.save') }}</UButton>
              </div>
            </div>
          </form>
        </section>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ApiError } from '~/composables/useApi'
import type { AdminSiteAnnouncement } from '~/composables/useAdmin'
import { formatDateTime } from '~/utils/format'

definePageMeta({ layout: 'admin', middleware: 'admin' })

type EditorMode = 'write' | 'preview'

const { t, locale } = useI18n()
const toast = useToast()
const adminApi = useAdmin()

const announcements = ref<AdminSiteAnnouncement[]>([])
const selected = ref<AdminSiteAnnouncement | null>(null)
const total = ref(0)
const page = ref(1)
const size = ref(20)
const pageSizes = [20, 50, 100]
const statusFilter = ref(-1)
const pending = ref(false)
const saving = ref(false)
const errorMessage = ref('')
const formError = ref('')
const editorMode = ref<EditorMode>('write')
const form = reactive({ id: 0, title: '', content: '', status: 1, publishedAt: '' })

useHead(() => ({ title: t('admin.site.announcements.title') }))
onMounted(loadAnnouncements)

async function loadAnnouncements() {
  pending.value = true
  errorMessage.value = ''
  try {
    const data = await adminApi.listSiteAnnouncements({
      page: page.value,
      size: size.value,
      status: statusFilter.value >= 0 ? statusFilter.value : undefined
    })
    announcements.value = data.list || []
    total.value = data.total || 0
  } catch (error) {
    errorMessage.value = error instanceof ApiError ? error.message : t('common.requestFailed')
  } finally {
    pending.value = false
  }
}

function startCreate() {
  selected.value = null
  resetForm()
}

function selectAnnouncement(item: AdminSiteAnnouncement) {
  selected.value = item
  form.id = item.id
  form.title = item.title
  form.content = item.content
  form.status = item.status
  form.publishedAt = toDatetimeLocal(item.publishedAt)
  editorMode.value = 'write'
  formError.value = ''
}

function resetForm() {
  if (selected.value) {
    selectAnnouncement(selected.value)
    return
  }
  form.id = 0
  form.title = ''
  form.content = ''
  form.status = 1
  form.publishedAt = ''
  editorMode.value = 'write'
  formError.value = ''
}

async function saveAnnouncement() {
  formError.value = ''
  if (!form.title || !form.content) {
    formError.value = t('admin.site.announcements.required')
    return
  }
  saving.value = true
  try {
    const input = {
      title: form.title,
      content: form.content,
      status: form.status,
      publishedAt: fromDatetimeLocal(form.publishedAt)
    }
    if (form.id) {
      await adminApi.updateSiteAnnouncement(form.id, input)
    } else {
      await adminApi.createSiteAnnouncement(input)
    }
    toast.add({ color: 'success', title: t('admin.site.announcements.saved'), icon: 'i-lucide-check' })
    await loadAnnouncements()
  } catch (error) {
    formError.value = error instanceof ApiError ? error.message : t('common.requestFailed')
  } finally {
    saving.value = false
  }
}

async function deleteAnnouncement() {
  if (!form.id) return
  saving.value = true
  try {
    await adminApi.deleteSiteAnnouncement(form.id)
    toast.add({ color: 'success', title: t('admin.site.announcements.deleted'), icon: 'i-lucide-check' })
    startCreate()
    await loadAnnouncements()
  } catch (error) {
    formError.value = error instanceof ApiError ? error.message : t('common.requestFailed')
  } finally {
    saving.value = false
  }
}

function reloadFromFirstPage() {
  page.value = 1
  loadAnnouncements()
}

function goToPage(nextPage: number) {
  page.value = nextPage
  loadAnnouncements()
}

function changeSize(nextSize: number) {
  size.value = nextSize
  page.value = 1
  loadAnnouncements()
}

function statusLabel(status: number) {
  if (status === 1) return t('admin.site.announcements.status.published')
  if (status === 2) return t('admin.site.announcements.status.archived')
  return t('admin.site.announcements.status.draft')
}

function statusColor(status: number) {
  if (status === 1) return 'success'
  if (status === 2) return 'neutral'
  return 'warning'
}

function toDatetimeLocal(value?: string | null) {
  if (!value) return ''
  return String(value).replace(' ', 'T').slice(0, 16)
}

function fromDatetimeLocal(value: string) {
  return value ? `${value.replace('T', ' ')}:00` : null
}
</script>
