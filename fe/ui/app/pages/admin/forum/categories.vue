<template>
  <div class="min-h-[calc(100vh-4rem)] bg-slate-50 py-6 dark:bg-slate-950">
    <div class="w-full px-3 sm:px-4 lg:px-5">
      <div class="grid gap-4 lg:grid-cols-[minmax(0,1fr)_380px]">
        <section class="overflow-hidden rounded-lg border border-slate-200 bg-white dark:border-slate-800 dark:bg-slate-900">
          <div class="flex items-center justify-between gap-3 border-b border-slate-200 px-4 py-3 dark:border-slate-800">
            <h2 class="text-sm font-semibold text-slate-950 dark:text-white">{{ $t('admin.forum.categories.list') }}</h2>
          </div>

          <div v-if="pending" class="overflow-x-auto">
            <table class="min-w-[760px] w-full table-fixed border-collapse text-left">
              <thead class="bg-slate-50 text-xs font-medium uppercase text-slate-500 dark:bg-slate-950/70 dark:text-slate-400">
                <tr>
                  <th class="w-[34%] border-b border-slate-200 px-4 py-3 dark:border-slate-800">{{ $t('admin.forum.categories.table.name') }}</th>
                  <th class="w-[26%] border-b border-slate-200 px-4 py-3 dark:border-slate-800">{{ $t('admin.forum.categories.table.description') }}</th>
                  <th class="w-[12%] border-b border-slate-200 px-4 py-3 text-right dark:border-slate-800">{{ $t('admin.forum.categories.table.sort') }}</th>
                  <th class="w-[12%] border-b border-slate-200 px-4 py-3 text-right dark:border-slate-800">{{ $t('admin.forum.categories.table.minRoleView') }}</th>
                  <th class="w-[16%] border-b border-slate-200 px-4 py-3 text-right dark:border-slate-800">{{ $t('admin.forum.categories.table.actions') }}</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="index in 4" :key="index" class="border-b border-slate-200 last:border-b-0 dark:border-slate-800">
                  <td class="px-4 py-4"><div class="h-4 w-44 animate-pulse rounded bg-slate-200 dark:bg-slate-800" /></td>
                  <td class="px-4 py-4"><div class="h-4 w-36 animate-pulse rounded bg-slate-100 dark:bg-slate-800/70" /></td>
                  <td class="px-4 py-4"><div class="ml-auto h-4 w-10 animate-pulse rounded bg-slate-100 dark:bg-slate-800/70" /></td>
                  <td class="px-4 py-4"><div class="ml-auto h-4 w-10 animate-pulse rounded bg-slate-100 dark:bg-slate-800/70" /></td>
                  <td class="px-4 py-4"><div class="ml-auto h-4 w-24 animate-pulse rounded bg-slate-100 dark:bg-slate-800/70" /></td>
                </tr>
              </tbody>
            </table>
          </div>

          <div v-else-if="errorMessage" class="flex flex-col items-center justify-center px-4 py-16 text-center">
            <UIcon name="i-lucide-circle-alert" class="size-9 text-red-500" />
            <p class="mt-3 text-sm font-medium text-slate-950 dark:text-white">{{ errorMessage }}</p>
          </div>

          <div v-else-if="categories.length === 0" class="flex flex-col items-center justify-center px-4 py-16 text-center">
            <UIcon name="i-lucide-inbox" class="size-9 text-slate-400" />
            <p class="mt-3 text-sm font-medium text-slate-950 dark:text-white">{{ $t('admin.forum.categories.empty') }}</p>
          </div>

          <div v-else class="overflow-x-auto">
            <table class="min-w-[760px] w-full table-fixed border-collapse text-left">
              <thead class="bg-slate-50 text-xs font-medium uppercase text-slate-500 dark:bg-slate-950/70 dark:text-slate-400">
                <tr>
                  <th class="w-[34%] border-b border-slate-200 px-4 py-3 dark:border-slate-800">{{ $t('admin.forum.categories.table.name') }}</th>
                  <th class="w-[26%] border-b border-slate-200 px-4 py-3 dark:border-slate-800">{{ $t('admin.forum.categories.table.description') }}</th>
                  <th class="w-[12%] border-b border-slate-200 px-4 py-3 text-right dark:border-slate-800">{{ $t('admin.forum.categories.table.sort') }}</th>
                  <th class="w-[12%] border-b border-slate-200 px-4 py-3 text-right dark:border-slate-800">{{ $t('admin.forum.categories.table.minRoleView') }}</th>
                  <th class="w-[16%] border-b border-slate-200 px-4 py-3 text-right dark:border-slate-800">{{ $t('admin.forum.categories.table.actions') }}</th>
                </tr>
              </thead>
              <tbody>
                <tr
                  v-for="category in categories"
                  :key="category.id"
                  class="border-b border-slate-200 transition-colors last:border-b-0 dark:border-slate-800"
                  :class="selectedId === category.id ? 'bg-cyan-50/70 dark:bg-cyan-950/30' : 'hover:bg-slate-50 dark:hover:bg-slate-950/70'"
                >
                  <td class="px-4 py-3 align-middle">
                    <button type="button" class="block max-w-full text-left" @click="selectCategory(category)">
                      <span class="block truncate text-sm font-semibold text-slate-950 hover:text-cyan-700 dark:text-white dark:hover:text-cyan-300">
                        {{ categoryName(category) }}
                      </span>
                      <span class="mt-1 block text-xs text-slate-500 dark:text-slate-400">
                        {{ formatDateTime(category.updatedAt || category.createdAt, locale) }}
                      </span>
                    </button>
                  </td>
                  <td class="px-4 py-3 align-middle text-sm text-slate-600 dark:text-slate-300">
                    <span class="block truncate">{{ categoryDescription(category) || '-' }}</span>
                  </td>
                  <td class="px-4 py-3 text-right align-middle text-sm text-slate-600 dark:text-slate-300">{{ numberFormatter.format(category.sortOrder) }}</td>
                  <td class="px-4 py-3 text-right align-middle text-sm text-slate-600 dark:text-slate-300">{{ numberFormatter.format(category.minRoleView) }}</td>
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
                {{ selectedId ? $t('admin.forum.categories.form.edit') : $t('admin.forum.categories.form.create') }}
              </h2>
              <UButton v-if="selectedId" color="neutral" variant="ghost" size="sm" icon="i-lucide-plus" @click="resetForm">
                {{ $t('admin.actions.new') }}
              </UButton>
            </div>
          </div>

          <form class="space-y-4 p-4" @submit.prevent="saveCategory">
            <label v-for="field in nameFields" :key="field.key" class="block">
              <span class="text-sm font-medium text-slate-700 dark:text-slate-200">{{ field.label }}</span>
              <input v-model="form[field.key]" class="mt-1 h-10 w-full rounded-md border border-slate-200 bg-white px-3 text-sm text-slate-950 outline-none transition focus:border-cyan-400 focus:ring-2 focus:ring-cyan-100 dark:border-slate-700 dark:bg-slate-950 dark:text-white dark:focus:border-cyan-500 dark:focus:ring-cyan-950" :disabled="saving">
            </label>

            <label v-for="field in descFields" :key="field.key" class="block">
              <span class="text-sm font-medium text-slate-700 dark:text-slate-200">{{ field.label }}</span>
              <textarea v-model="form[field.key]" rows="2" class="mt-1 w-full resize-none rounded-md border border-slate-200 bg-white px-3 py-2 text-sm text-slate-950 outline-none transition focus:border-cyan-400 focus:ring-2 focus:ring-cyan-100 dark:border-slate-700 dark:bg-slate-950 dark:text-white dark:focus:border-cyan-500 dark:focus:ring-cyan-950" :disabled="saving" />
            </label>

            <div class="grid grid-cols-2 gap-3">
              <label class="block">
                <span class="text-sm font-medium text-slate-700 dark:text-slate-200">{{ $t('admin.forum.categories.form.sortOrder') }}</span>
                <input v-model.number="form.sortOrder" type="number" class="mt-1 h-10 w-full rounded-md border border-slate-200 bg-white px-3 text-sm text-slate-950 outline-none transition focus:border-cyan-400 focus:ring-2 focus:ring-cyan-100 dark:border-slate-700 dark:bg-slate-950 dark:text-white dark:focus:border-cyan-500 dark:focus:ring-cyan-950" :disabled="saving">
              </label>
              <label class="block">
                <span class="text-sm font-medium text-slate-700 dark:text-slate-200">{{ $t('admin.forum.categories.form.minRoleView') }}</span>
                <input v-model.number="form.minRoleView" type="number" min="0" class="mt-1 h-10 w-full rounded-md border border-slate-200 bg-white px-3 text-sm text-slate-950 outline-none transition focus:border-cyan-400 focus:ring-2 focus:ring-cyan-100 dark:border-slate-700 dark:bg-slate-950 dark:text-white dark:focus:border-cyan-500 dark:focus:ring-cyan-950" :disabled="saving">
              </label>
            </div>

            <p v-if="formError" class="rounded-md border border-red-200 bg-red-50 px-3 py-2 text-sm text-red-700 dark:border-red-900 dark:bg-red-950/40 dark:text-red-200">{{ formError }}</p>

            <div class="flex flex-col gap-2 sm:flex-row">
              <UButton type="submit" color="primary" icon="i-lucide-save" :loading="saving" :disabled="!canSubmit">{{ $t('common.save') }}</UButton>
              <UButton type="button" color="neutral" variant="outline" icon="i-lucide-rotate-ccw" :disabled="saving" @click="resetForm">{{ $t('admin.actions.reset') }}</UButton>
            </div>
          </form>
        </section>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ApiError } from '~/composables/useApi'
import type { AdminForumCategory, AdminForumCategoryInput } from '~/composables/useAdmin'
import { formatDateTime, localizeI18nName } from '~/utils/format'

definePageMeta({ layout: 'admin', middleware: 'admin' })

const { t, locale } = useI18n()
const toast = useToast()
const adminApi = useAdmin()

useHead({ title: t('admin.forum.categories.title') })

const categories = ref<AdminForumCategory[]>([])
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
  descZhCN: '',
  descZhTW: '',
  descEnUS: '',
  sortOrder: 0,
  minRoleView: 0
})

const numberFormatter = computed(() => new Intl.NumberFormat(locale.value))
const canSubmit = computed(() => Boolean(form.nameZhCN.trim() && !saving.value))
const nameFields = computed(() => [
  { key: 'nameZhCN' as const, label: t('admin.forum.categories.form.nameZhCN') },
  { key: 'nameZhTW' as const, label: t('admin.forum.categories.form.nameZhTW') },
  { key: 'nameEnUS' as const, label: t('admin.forum.categories.form.nameEnUS') }
])
const descFields = computed(() => [
  { key: 'descZhCN' as const, label: t('admin.forum.categories.form.descZhCN') },
  { key: 'descZhTW' as const, label: t('admin.forum.categories.form.descZhTW') },
  { key: 'descEnUS' as const, label: t('admin.forum.categories.form.descEnUS') }
])

onMounted(loadCategories)

async function loadCategories() {
  pending.value = true
  errorMessage.value = ''
  try {
    const data = await adminApi.listForumCategories()
    categories.value = (data.categories || []).sort((a, b) => a.sortOrder - b.sortOrder || a.id - b.id)
  } catch (error: unknown) {
    errorMessage.value = error instanceof ApiError ? error.message : t('common.requestFailed')
  } finally {
    pending.value = false
  }
}

function categoryName(category: AdminForumCategory) {
  return localizeI18nName(category.nameI18N, locale.value, `#${category.id}`)
}

function categoryDescription(category: AdminForumCategory) {
  return localizeI18nName(category.descI18N, locale.value, '')
}

function selectCategory(category: AdminForumCategory) {
  selectedId.value = category.id
  form.nameZhCN = String(category.nameI18N?.['zh-CN'] || '')
  form.nameZhTW = String(category.nameI18N?.['zh-TW'] || '')
  form.nameEnUS = String(category.nameI18N?.['en-US'] || '')
  form.descZhCN = String(category.descI18N?.['zh-CN'] || '')
  form.descZhTW = String(category.descI18N?.['zh-TW'] || '')
  form.descEnUS = String(category.descI18N?.['en-US'] || '')
  form.sortOrder = category.sortOrder
  form.minRoleView = category.minRoleView
  formError.value = ''
}

function resetForm() {
  selectedId.value = null
  form.nameZhCN = ''
  form.nameZhTW = ''
  form.nameEnUS = ''
  form.descZhCN = ''
  form.descZhTW = ''
  form.descEnUS = ''
  form.sortOrder = 0
  form.minRoleView = 0
  formError.value = ''
}

function buildInput(): AdminForumCategoryInput {
  const zhCN = form.nameZhCN.trim()
  const descZhCN = form.descZhCN.trim()
  return {
    nameI18N: {
      'zh-CN': zhCN,
      'zh-TW': form.nameZhTW.trim() || zhCN,
      'en-US': form.nameEnUS.trim() || zhCN
    },
    descI18N: {
      'zh-CN': descZhCN,
      'zh-TW': form.descZhTW.trim() || descZhCN,
      'en-US': form.descEnUS.trim() || descZhCN
    },
    sortOrder: Number(form.sortOrder || 0),
    minRoleView: Number(form.minRoleView || 0)
  }
}

async function saveCategory() {
  if (!canSubmit.value) return
  saving.value = true
  formError.value = ''
  try {
    const input = buildInput()
    if (selectedId.value) {
      await adminApi.updateForumCategory(selectedId.value, input)
    } else {
      await adminApi.createForumCategory(input)
    }
    toast.add({ title: t('admin.forum.categories.saved') })
    resetForm()
    await loadCategories()
  } catch (error: unknown) {
    formError.value = error instanceof ApiError ? error.message : t('common.requestFailed')
  } finally {
    saving.value = false
  }
}

async function deleteCategory(category: AdminForumCategory) {
  if (!window.confirm(t('admin.forum.categories.confirmDelete', { name: categoryName(category) }))) return
  deletingId.value = category.id
  try {
    await adminApi.deleteForumCategory(category.id)
    if (selectedId.value === category.id) resetForm()
    toast.add({ title: t('admin.forum.categories.deleted') })
    await loadCategories()
  } catch (error: unknown) {
    toast.add({ color: 'error', title: error instanceof ApiError ? error.message : t('common.requestFailed') })
  } finally {
    deletingId.value = null
  }
}
</script>
