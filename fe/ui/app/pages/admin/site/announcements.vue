<template>
  <div class="min-h-[calc(100vh-4rem)] bg-slate-50 py-6 dark:bg-slate-950">
    <div class="w-full px-3 sm:px-4 lg:px-5">
      <div class="grid gap-4 xl:grid-cols-[minmax(360px,480px)_minmax(0,1fr)]">
        <section class="overflow-hidden rounded-lg border border-slate-200 bg-white dark:border-slate-800 dark:bg-slate-900">
          <div class="border-b border-slate-200 px-4 py-3 dark:border-slate-800">
            <div class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
              <h2 class="text-sm font-semibold text-slate-950 dark:text-white">{{ $t('admin.site.announcements.list') }}</h2>
              <div class="inline-flex h-8 max-w-full items-center overflow-x-auto rounded-md border border-slate-200 bg-slate-50 p-0.5 dark:border-slate-800 dark:bg-slate-950/70" :aria-label="$t('admin.site.announcements.fields.status')" role="group">
                <button
                  v-for="option in statusOptions"
                  :key="option.value"
                  type="button"
                  class="h-7 whitespace-nowrap rounded px-2.5 text-xs font-medium transition disabled:cursor-not-allowed disabled:opacity-60"
                  :class="statusFilter === option.value ? 'bg-white text-slate-950 shadow-sm dark:bg-slate-800 dark:text-white' : 'text-slate-500 hover:text-slate-950 dark:text-slate-400 dark:hover:text-white'"
                  :disabled="pending"
                  @click="setStatusFilter(option.value)"
                >
                  {{ option.label }}
                </button>
              </div>
            </div>
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
          <div v-else class="space-y-1 p-2">
            <button
              v-for="item in announcements"
              :key="item.id"
              type="button"
              class="group flex w-full cursor-pointer items-start gap-3 rounded-md px-3 py-2.5 text-left transition focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-sky-300 dark:focus-visible:ring-sky-700"
              :class="selected?.id === item.id ? 'bg-sky-50 ring-1 ring-sky-200 dark:bg-sky-950/30 dark:ring-sky-900' : 'hover:bg-slate-50 dark:hover:bg-slate-950/70'"
              @click="selectAnnouncement(item)"
            >
              <span class="mt-1 size-2.5 shrink-0 rounded-full" :class="statusDotClass(item.status)" />
              <div class="min-w-0 flex-1">
                <p class="truncate text-sm font-semibold text-slate-950 transition group-hover:text-sky-700 dark:text-white dark:group-hover:text-sky-300">{{ item.title }}</p>
                <p class="mt-1 truncate text-xs text-slate-500 dark:text-slate-400">
                  #{{ item.id }}
                  <span class="mx-1.5 text-slate-300 dark:text-slate-700">/</span>
                  {{ formatDateTime(item.publishedAt || item.createdAt, locale) }}
                </p>
              </div>
              <UBadge :color="statusColor(item.status)" variant="soft" class="shrink-0">{{ statusLabel(item.status) }}</UBadge>
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
          <div class="border-b border-slate-200 px-4 py-3 dark:border-slate-800">
            <div class="flex items-center justify-between gap-3">
              <h2 class="truncate text-sm font-semibold text-slate-950 dark:text-white">{{ form.id ? $t('admin.site.announcements.edit') : $t('admin.site.announcements.create') }}</h2>
              <div class="flex h-8 min-w-20 items-center justify-end">
                <UButton v-if="form.id" color="primary" variant="soft" size="sm" icon="i-lucide-plus" @click="startCreate">
                  {{ $t('admin.actions.new') }}
                </UButton>
              </div>
            </div>
          </div>

          <form class="space-y-4 p-4" @submit.prevent="saveAnnouncement">
            <label class="block">
              <span class="text-sm font-medium text-slate-700 dark:text-slate-200">{{ $t('admin.site.announcements.fields.title') }}</span>
              <input v-model.trim="form.title" class="mt-1 h-10 w-full rounded-md border border-slate-200 bg-white px-3 text-sm text-slate-950 outline-none transition focus:border-sky-400 focus:ring-2 focus:ring-sky-100 dark:border-slate-700 dark:bg-slate-950 dark:text-white dark:focus:border-sky-500 dark:focus:ring-sky-950">
            </label>
            <div class="grid gap-3 sm:grid-cols-2">
              <UFormField :label="$t('admin.site.announcements.fields.status')">
                <USelect v-model="form.status" class="w-full" size="lg" :ui="{ base: 'h-10 w-full' }" :items="formStatusOptions" value-key="value" />
              </UFormField>
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
              <UPopover
                v-if="form.id"
                :content="{ side: 'top', align: 'start', sideOffset: 8 }"
                :ui="{ content: 'w-72 p-3' }"
              >
                <UButton color="error" variant="soft" icon="i-lucide-trash-2" :disabled="saving">
                  {{ $t('common.delete') }}
                </UButton>

                <template #content="{ close }">
                  <div class="space-y-3">
                    <p class="text-sm font-medium text-slate-950 dark:text-white">
                      {{ $t('admin.actions.confirmDeleteTitle') }}
                    </p>
                    <p class="text-xs text-slate-500 dark:text-slate-400">
                      {{ $t('admin.actions.deleteIrreversible') }}
                    </p>
                    <div class="flex justify-end gap-2">
                      <UButton color="neutral" variant="ghost" size="xs" type="button" @click="close()">
                        {{ $t('common.cancel') }}
                      </UButton>
                      <UButton color="error" size="xs" type="button" icon="i-lucide-trash-2" :loading="saving" :disabled="saving" @click="deleteAnnouncement(close)">
                        {{ $t('common.delete') }}
                      </UButton>
                    </div>
                  </div>
                </template>
              </UPopover>
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
const statusOptions = computed(() => [
  { value: -1, label: t('admin.site.announcements.filters.all') },
  { value: 1, label: t('admin.site.announcements.status.published') },
  { value: 0, label: t('admin.site.announcements.status.draft') },
  { value: 2, label: t('admin.site.announcements.status.archived') }
])
const formStatusOptions = computed(() => statusOptions.value.filter(option => option.value >= 0))

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

async function deleteAnnouncement(close?: () => void) {
  if (!form.id) return
  saving.value = true
  try {
    await adminApi.deleteSiteAnnouncement(form.id)
    toast.add({ color: 'success', title: t('admin.site.announcements.deleted'), icon: 'i-lucide-check' })
    close?.()
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

function setStatusFilter(status: number) {
  if (statusFilter.value === status) return
  statusFilter.value = status
  reloadFromFirstPage()
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

function statusDotClass(status: number) {
  if (status === 1) return 'bg-emerald-500'
  if (status === 2) return 'bg-slate-300 dark:bg-slate-700'
  return 'bg-amber-500'
}

function toDatetimeLocal(value?: string | null) {
  if (!value) return ''
  return String(value).replace(' ', 'T').slice(0, 16)
}

function fromDatetimeLocal(value: string) {
  return value ? `${value.replace('T', ' ')}:00` : null
}
</script>
