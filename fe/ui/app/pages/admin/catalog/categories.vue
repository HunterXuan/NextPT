<template>
  <div class="min-h-[calc(100vh-4rem)] bg-slate-50 py-6 dark:bg-slate-950">
    <div class="w-full px-3 sm:px-4 lg:px-5">
      <div class="mb-4 grid grid-cols-3 gap-2 sm:flex sm:items-center sm:justify-end">
        <div class="rounded-lg border border-slate-200 bg-white px-3 py-2 dark:border-slate-800 dark:bg-slate-900">
          <p class="text-xs text-slate-500 dark:text-slate-400">{{ $t('admin.catalog.categories.stats.total') }}</p>
          <p class="mt-1 text-lg font-semibold text-slate-950 dark:text-white">{{ numberFormatter.format(categories.length) }}</p>
        </div>
        <div class="rounded-lg border border-slate-200 bg-white px-3 py-2 dark:border-slate-800 dark:bg-slate-900">
          <p class="text-xs text-slate-500 dark:text-slate-400">{{ $t('admin.catalog.categories.stats.enabled') }}</p>
          <p class="mt-1 text-lg font-semibold text-emerald-700 dark:text-emerald-300">{{ numberFormatter.format(enabledCount) }}</p>
        </div>
        <div class="rounded-lg border border-slate-200 bg-white px-3 py-2 dark:border-slate-800 dark:bg-slate-900">
          <p class="text-xs text-slate-500 dark:text-slate-400">{{ $t('admin.catalog.categories.stats.disabled') }}</p>
          <p class="mt-1 text-lg font-semibold text-slate-700 dark:text-slate-200">{{ numberFormatter.format(disabledCount) }}</p>
        </div>
      </div>

      <div class="grid gap-4 lg:grid-cols-[minmax(0,1fr)_360px]">
        <section class="overflow-hidden rounded-lg border border-slate-200 bg-white dark:border-slate-800 dark:bg-slate-900">
          <div class="flex items-center justify-between gap-3 border-b border-slate-200 px-4 py-3 dark:border-slate-800">
            <div class="flex items-center gap-2">
              <UIcon name="i-lucide-tags" class="size-5 text-sky-600 dark:text-sky-300" />
              <h2 class="text-sm font-semibold text-slate-950 dark:text-white">{{ $t('admin.catalog.categories.list') }}</h2>
            </div>
            <UButton color="neutral" variant="outline" icon="i-lucide-refresh-cw" :loading="pending" @click="loadCategories">
              {{ $t('common.refresh') }}
            </UButton>
          </div>

          <div v-if="pending" class="overflow-x-auto">
            <table class="min-w-[820px] w-full table-fixed border-collapse text-left">
              <thead class="bg-slate-50 text-xs font-medium uppercase text-slate-500 dark:bg-slate-950/70 dark:text-slate-400">
                <tr>
                  <th class="w-[32%] border-b border-slate-200 px-4 py-3 dark:border-slate-800">{{ $t('admin.catalog.categories.table.name') }}</th>
                  <th class="w-[18%] border-b border-slate-200 px-4 py-3 dark:border-slate-800">{{ $t('admin.catalog.categories.table.slug') }}</th>
                  <th class="w-[12%] border-b border-slate-200 px-4 py-3 text-right dark:border-slate-800">{{ $t('admin.catalog.categories.table.sort') }}</th>
                  <th class="w-[14%] border-b border-slate-200 px-4 py-3 dark:border-slate-800">{{ $t('admin.catalog.categories.table.status') }}</th>
                  <th class="w-[24%] border-b border-slate-200 px-4 py-3 text-right dark:border-slate-800">{{ $t('admin.catalog.categories.table.actions') }}</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="index in 5" :key="index" class="border-b border-slate-200 last:border-b-0 dark:border-slate-800">
                  <td class="px-4 py-4">
                    <div class="h-4 w-48 animate-pulse rounded bg-slate-200 dark:bg-slate-800" />
                    <div class="mt-2 h-3 w-28 animate-pulse rounded bg-slate-100 dark:bg-slate-800/70" />
                  </td>
                  <td class="px-4 py-4">
                    <div class="h-4 w-24 animate-pulse rounded bg-slate-100 dark:bg-slate-800/70" />
                  </td>
                  <td class="px-4 py-4">
                    <div class="ml-auto h-4 w-12 animate-pulse rounded bg-slate-100 dark:bg-slate-800/70" />
                  </td>
                  <td class="px-4 py-4">
                    <div class="h-4 w-16 animate-pulse rounded bg-slate-100 dark:bg-slate-800/70" />
                  </td>
                  <td class="px-4 py-4">
                    <div class="ml-auto h-4 w-28 animate-pulse rounded bg-slate-100 dark:bg-slate-800/70" />
                  </td>
                </tr>
              </tbody>
            </table>
          </div>

          <div v-else-if="errorMessage" class="flex flex-col items-center justify-center px-4 py-16 text-center">
            <UIcon name="i-lucide-circle-alert" class="size-9 text-red-500" />
            <p class="mt-3 text-sm font-medium text-slate-950 dark:text-white">{{ errorMessage }}</p>
            <UButton class="mt-5" color="neutral" variant="outline" icon="i-lucide-refresh-cw" @click="loadCategories">
              {{ $t('common.retry') }}
            </UButton>
          </div>

          <div v-else-if="categories.length === 0" class="flex flex-col items-center justify-center px-4 py-16 text-center">
            <UIcon name="i-lucide-inbox" class="size-9 text-slate-400" />
            <p class="mt-3 text-sm font-medium text-slate-950 dark:text-white">{{ $t('admin.catalog.categories.empty') }}</p>
          </div>

          <div v-else class="overflow-x-auto">
            <table class="min-w-[820px] w-full table-fixed border-collapse text-left">
              <thead class="bg-slate-50 text-xs font-medium uppercase text-slate-500 dark:bg-slate-950/70 dark:text-slate-400">
                <tr>
                  <th class="w-[32%] border-b border-slate-200 px-4 py-3 dark:border-slate-800">{{ $t('admin.catalog.categories.table.name') }}</th>
                  <th class="w-[18%] border-b border-slate-200 px-4 py-3 dark:border-slate-800">{{ $t('admin.catalog.categories.table.slug') }}</th>
                  <th class="w-[12%] border-b border-slate-200 px-4 py-3 text-right dark:border-slate-800">{{ $t('admin.catalog.categories.table.sort') }}</th>
                  <th class="w-[14%] border-b border-slate-200 px-4 py-3 dark:border-slate-800">{{ $t('admin.catalog.categories.table.status') }}</th>
                  <th class="w-[24%] border-b border-slate-200 px-4 py-3 text-right dark:border-slate-800">{{ $t('admin.catalog.categories.table.actions') }}</th>
                </tr>
              </thead>
              <tbody>
                <tr
                  v-for="category in categories"
                  :key="category.id"
                  class="border-b border-slate-200 transition-colors last:border-b-0 dark:border-slate-800"
                  :class="selectedId === category.id ? 'bg-sky-50/70 dark:bg-sky-950/30' : 'hover:bg-slate-50 dark:hover:bg-slate-950/70'"
                >
                  <td class="px-4 py-3 align-middle">
                    <button class="block max-w-full text-left" type="button" @click="selectCategory(category)">
                      <span class="block truncate text-sm font-semibold text-slate-950 hover:text-sky-700 dark:text-white dark:hover:text-sky-300">
                        {{ categoryDisplayName(category) }}
                      </span>
                      <span class="mt-1 block truncate text-xs text-slate-500 dark:text-slate-400">
                        {{ formatDateTime(category.updatedAt || category.createdAt, locale) }}
                      </span>
                    </button>
                  </td>
                  <td class="px-4 py-3 align-middle">
                    <code class="block truncate rounded bg-slate-100 px-2 py-1 text-xs text-slate-700 dark:bg-slate-800 dark:text-slate-200">
                      {{ category.slug }}
                    </code>
                  </td>
                  <td class="px-4 py-3 text-right align-middle text-sm text-slate-600 dark:text-slate-300">
                    {{ numberFormatter.format(category.sortOrder) }}
                  </td>
                  <td class="px-4 py-3 align-middle">
                    <UBadge :color="category.enabled ? 'success' : 'neutral'" variant="soft">
                      {{ category.enabled ? $t('admin.status.enabled') : $t('admin.status.disabled') }}
                    </UBadge>
                  </td>
                  <td class="px-4 py-3 align-middle">
                    <div class="flex items-center justify-end gap-2">
                      <UButton color="neutral" variant="outline" size="sm" icon="i-lucide-pencil" @click="selectCategory(category)">
                        {{ $t('admin.actions.edit') }}
                      </UButton>
                      <UButton color="error" variant="soft" size="sm" icon="i-lucide-trash-2" :loading="deletingId === category.id" @click="deleteCategory(category)">
                        {{ $t('admin.actions.delete') }}
                      </UButton>
                    </div>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </section>

        <section class="rounded-lg border border-slate-200 bg-white dark:border-slate-800 dark:bg-slate-900">
          <div class="border-b border-slate-200 px-4 py-3 dark:border-slate-800">
            <div class="flex items-center justify-between gap-3">
              <h2 class="text-sm font-semibold text-slate-950 dark:text-white">
                {{ selectedId ? $t('admin.catalog.categories.form.edit') : $t('admin.catalog.categories.form.create') }}
              </h2>
              <UButton v-if="selectedId" color="neutral" variant="ghost" size="sm" icon="i-lucide-plus" @click="resetForm">
                {{ $t('admin.actions.new') }}
              </UButton>
            </div>
          </div>

          <form class="space-y-4 p-4" @submit.prevent="saveCategory">
            <label class="block">
              <span class="text-sm font-medium text-slate-700 dark:text-slate-200">{{ $t('admin.catalog.categories.form.nameZhCN') }}</span>
              <input
                v-model="form.nameZhCN"
                class="mt-1 h-10 w-full rounded-md border border-slate-200 bg-white px-3 text-sm text-slate-950 outline-none transition focus:border-sky-400 focus:ring-2 focus:ring-sky-100 dark:border-slate-700 dark:bg-slate-950 dark:text-white dark:focus:border-sky-500 dark:focus:ring-sky-950"
                :disabled="saving"
              >
            </label>

            <label class="block">
              <span class="text-sm font-medium text-slate-700 dark:text-slate-200">{{ $t('admin.catalog.categories.form.nameZhTW') }}</span>
              <input
                v-model="form.nameZhTW"
                class="mt-1 h-10 w-full rounded-md border border-slate-200 bg-white px-3 text-sm text-slate-950 outline-none transition focus:border-sky-400 focus:ring-2 focus:ring-sky-100 dark:border-slate-700 dark:bg-slate-950 dark:text-white dark:focus:border-sky-500 dark:focus:ring-sky-950"
                :disabled="saving"
              >
            </label>

            <label class="block">
              <span class="text-sm font-medium text-slate-700 dark:text-slate-200">{{ $t('admin.catalog.categories.form.nameEnUS') }}</span>
              <input
                v-model="form.nameEnUS"
                class="mt-1 h-10 w-full rounded-md border border-slate-200 bg-white px-3 text-sm text-slate-950 outline-none transition focus:border-sky-400 focus:ring-2 focus:ring-sky-100 dark:border-slate-700 dark:bg-slate-950 dark:text-white dark:focus:border-sky-500 dark:focus:ring-sky-950"
                :disabled="saving"
              >
            </label>

            <label class="block">
              <span class="text-sm font-medium text-slate-700 dark:text-slate-200">{{ $t('admin.catalog.categories.form.slug') }}</span>
              <input
                v-model="form.slug"
                class="mt-1 h-10 w-full rounded-md border border-slate-200 bg-white px-3 text-sm text-slate-950 outline-none transition focus:border-sky-400 focus:ring-2 focus:ring-sky-100 dark:border-slate-700 dark:bg-slate-950 dark:text-white dark:focus:border-sky-500 dark:focus:ring-sky-950"
                :disabled="saving"
              >
            </label>

            <div class="grid grid-cols-[1fr_auto] gap-3">
              <label class="block">
                <span class="text-sm font-medium text-slate-700 dark:text-slate-200">{{ $t('admin.catalog.categories.form.sortOrder') }}</span>
                <input
                  v-model.number="form.sortOrder"
                  type="number"
                  class="mt-1 h-10 w-full rounded-md border border-slate-200 bg-white px-3 text-sm text-slate-950 outline-none transition focus:border-sky-400 focus:ring-2 focus:ring-sky-100 dark:border-slate-700 dark:bg-slate-950 dark:text-white dark:focus:border-sky-500 dark:focus:ring-sky-950"
                  :disabled="saving"
                >
              </label>

              <label class="flex items-end gap-2 pb-2 text-sm font-medium text-slate-700 dark:text-slate-200">
                <input
                  v-model="form.enabled"
                  type="checkbox"
                  class="size-4 rounded border-slate-300 text-sky-600 focus:ring-sky-500 dark:border-slate-600"
                  :disabled="saving"
                >
                <span>{{ $t('admin.catalog.categories.form.enabled') }}</span>
              </label>
            </div>

            <label class="block">
              <span class="text-sm font-medium text-slate-700 dark:text-slate-200">{{ $t('admin.catalog.categories.form.uploadConfig') }}</span>
              <textarea
                v-model="form.uploadConfigText"
                class="mt-1 min-h-52 w-full rounded-md border border-slate-200 bg-white px-3 py-2 font-mono text-xs leading-5 text-slate-950 outline-none transition focus:border-sky-400 focus:ring-2 focus:ring-sky-100 dark:border-slate-700 dark:bg-slate-950 dark:text-white dark:focus:border-sky-500 dark:focus:ring-sky-950"
                :placeholder="$t('admin.catalog.categories.form.uploadConfigPlaceholder')"
                :disabled="saving"
              />
            </label>

            <p v-if="formError" class="rounded-md border border-red-200 bg-red-50 px-3 py-2 text-sm text-red-700 dark:border-red-900 dark:bg-red-950/40 dark:text-red-200">
              {{ formError }}
            </p>

            <div class="flex flex-col gap-2 sm:flex-row">
              <UButton type="submit" color="primary" icon="i-lucide-save" :loading="saving" :disabled="!canSubmit">
                {{ $t('common.save') }}
              </UButton>
              <UButton type="button" color="neutral" variant="outline" icon="i-lucide-rotate-ccw" :disabled="saving" @click="resetForm">
                {{ $t('admin.actions.reset') }}
              </UButton>
            </div>
          </form>
        </section>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ApiError } from '~/composables/useApi'
import type { AdminCatalogCategory, AdminCatalogCategoryInput } from '~/composables/useAdmin'
import { formatDateTime, localizeI18nName } from '~/utils/format'

definePageMeta({
  layout: 'admin',
  middleware: 'admin'
})

const { t, locale } = useI18n()
const toast = useToast()
const adminApi = useAdmin()

useHead({
  title: t('admin.catalog.categories.title')
})

const categories = ref<AdminCatalogCategory[]>([])
const pending = ref(false)
const saving = ref(false)
const deletingId = ref<number | null>(null)
const selectedId = ref<number | null>(null)
const errorMessage = ref('')
const formError = ref('')

const form = reactive({
  nameZhCN: '',
  nameZhTW: '',
  nameEnUS: '',
  slug: '',
  sortOrder: 0,
  enabled: true,
  uploadConfigText: ''
})

const numberFormatter = computed(() => new Intl.NumberFormat(locale.value))
const enabledCount = computed(() => categories.value.filter((item) => item.enabled).length)
const disabledCount = computed(() => categories.value.length - enabledCount.value)
const canSubmit = computed(() => Boolean(form.nameZhCN.trim() && form.slug.trim() && !saving.value))

onMounted(loadCategories)

async function loadCategories() {
  pending.value = true
  errorMessage.value = ''
  try {
    const data = await adminApi.listCatalogCategories()
    categories.value = (data.categories || []).sort((a, b) => a.sortOrder - b.sortOrder || a.id - b.id)
  } catch (error: unknown) {
    errorMessage.value = error instanceof ApiError ? error.message : t('common.requestFailed')
  } finally {
    pending.value = false
  }
}

function categoryDisplayName(category: AdminCatalogCategory) {
  return localizeI18nName(category.nameI18N, locale.value, category.slug || `#${category.id}`)
}

function selectCategory(category: AdminCatalogCategory) {
  selectedId.value = category.id
  form.nameZhCN = String(category.nameI18N?.['zh-CN'] || '')
  form.nameZhTW = String(category.nameI18N?.['zh-TW'] || '')
  form.nameEnUS = String(category.nameI18N?.['en-US'] || '')
  form.slug = category.slug
  form.sortOrder = category.sortOrder
  form.enabled = category.enabled
  form.uploadConfigText = category.uploadConfig ? JSON.stringify(category.uploadConfig, null, 2) : ''
  formError.value = ''
}

function resetForm() {
  selectedId.value = null
  form.nameZhCN = ''
  form.nameZhTW = ''
  form.nameEnUS = ''
  form.slug = ''
  form.sortOrder = 0
  form.enabled = true
  form.uploadConfigText = ''
  formError.value = ''
}

function buildInput(): AdminCatalogCategoryInput {
  const zhCN = form.nameZhCN.trim()
  const zhTW = form.nameZhTW.trim() || zhCN
  const enUS = form.nameEnUS.trim() || form.slug.trim()

  const uploadConfig = parseUploadConfigText()

  return {
    nameI18N: {
      'zh-CN': zhCN,
      'zh-TW': zhTW,
      'en-US': enUS
    },
    slug: form.slug.trim(),
    sortOrder: Number(form.sortOrder || 0),
    enabled: form.enabled,
    uploadConfig
  }
}

function parseUploadConfigText() {
  const text = form.uploadConfigText.trim()
  if (!text) return null
  try {
    return JSON.parse(text)
  } catch {
    throw new Error(t('admin.catalog.categories.form.uploadConfigInvalid'))
  }
}

async function saveCategory() {
  if (!canSubmit.value) return

  saving.value = true
  formError.value = ''
  try {
    const input = buildInput()
    if (selectedId.value) {
      await adminApi.updateCatalogCategory(selectedId.value, input)
    } else {
      await adminApi.createCatalogCategory(input)
    }
    toast.add({ title: t('admin.catalog.categories.saved') })
    resetForm()
    await loadCategories()
  } catch (error: unknown) {
    formError.value = error instanceof ApiError || error instanceof Error ? error.message : t('common.requestFailed')
  } finally {
    saving.value = false
  }
}

async function deleteCategory(category: AdminCatalogCategory) {
  if (!window.confirm(t('admin.catalog.categories.confirmDelete', { name: categoryDisplayName(category) }))) return

  deletingId.value = category.id
  try {
    await adminApi.deleteCatalogCategory(category.id)
    if (selectedId.value === category.id) {
      resetForm()
    }
    toast.add({ title: t('admin.catalog.categories.deleted') })
    await loadCategories()
  } catch (error: unknown) {
    toast.add({
      color: 'error',
      title: error instanceof ApiError ? error.message : t('common.requestFailed')
    })
  } finally {
    deletingId.value = null
  }
}
</script>
