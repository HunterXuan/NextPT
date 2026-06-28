<template>
  <div class="min-h-[calc(100vh-4rem)] bg-slate-50 py-6 dark:bg-slate-950">
    <div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
      <div class="mb-6 flex flex-col gap-4 lg:flex-row lg:items-end lg:justify-between">
        <div>
          <p class="text-sm font-medium text-slate-500 dark:text-slate-400">{{ $t('admin.forum.eyebrow') }}</p>
          <h1 class="mt-1 text-2xl font-semibold text-slate-950 dark:text-white">{{ $t('admin.forum.nodes.title') }}</h1>
        </div>

        <div class="flex flex-wrap items-center gap-2">
          <UButton color="neutral" variant="outline" icon="i-lucide-folder-tree" :to="localePath('/admin/forum/categories')">
            {{ $t('admin.forum.categories.title') }}
          </UButton>
          <UButton color="primary" variant="soft" icon="i-lucide-panels-top-left" :to="localePath('/admin/forum/nodes')">
            {{ $t('admin.forum.nodes.title') }}
          </UButton>
        </div>
      </div>

      <div class="grid gap-4 xl:grid-cols-[minmax(0,1fr)_400px]">
        <section class="overflow-hidden rounded-lg border border-slate-200 bg-white dark:border-slate-800 dark:bg-slate-900">
          <div class="flex items-center justify-between gap-3 border-b border-slate-200 px-4 py-3 dark:border-slate-800">
            <h2 class="text-sm font-semibold text-slate-950 dark:text-white">{{ $t('admin.forum.nodes.list') }}</h2>
            <UButton color="neutral" variant="outline" icon="i-lucide-refresh-cw" :loading="pending" @click="loadAll">
              {{ $t('common.refresh') }}
            </UButton>
          </div>

          <div v-if="pending" class="overflow-x-auto">
            <table class="min-w-[960px] w-full table-fixed border-collapse text-left">
              <thead class="bg-slate-50 text-xs font-medium uppercase text-slate-500 dark:bg-slate-950/70 dark:text-slate-400">
                <tr>
                  <th class="w-[24%] border-b border-slate-200 px-4 py-3 dark:border-slate-800">{{ $t('admin.forum.nodes.table.name') }}</th>
                  <th class="w-[16%] border-b border-slate-200 px-4 py-3 dark:border-slate-800">{{ $t('admin.forum.nodes.table.category') }}</th>
                  <th class="w-[16%] border-b border-slate-200 px-4 py-3 dark:border-slate-800">{{ $t('admin.forum.nodes.table.slug') }}</th>
                  <th class="w-[12%] border-b border-slate-200 px-4 py-3 text-right dark:border-slate-800">{{ $t('admin.forum.nodes.table.permissions') }}</th>
                  <th class="w-[12%] border-b border-slate-200 px-4 py-3 text-right dark:border-slate-800">{{ $t('admin.forum.nodes.table.counts') }}</th>
                  <th class="w-[20%] border-b border-slate-200 px-4 py-3 text-right dark:border-slate-800">{{ $t('admin.forum.nodes.table.actions') }}</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="index in 5" :key="index" class="border-b border-slate-200 last:border-b-0 dark:border-slate-800">
                  <td class="px-4 py-4"><div class="h-4 w-44 animate-pulse rounded bg-slate-200 dark:bg-slate-800" /></td>
                  <td class="px-4 py-4"><div class="h-4 w-24 animate-pulse rounded bg-slate-100 dark:bg-slate-800/70" /></td>
                  <td class="px-4 py-4"><div class="h-4 w-24 animate-pulse rounded bg-slate-100 dark:bg-slate-800/70" /></td>
                  <td class="px-4 py-4"><div class="ml-auto h-4 w-16 animate-pulse rounded bg-slate-100 dark:bg-slate-800/70" /></td>
                  <td class="px-4 py-4"><div class="ml-auto h-4 w-16 animate-pulse rounded bg-slate-100 dark:bg-slate-800/70" /></td>
                  <td class="px-4 py-4"><div class="ml-auto h-4 w-24 animate-pulse rounded bg-slate-100 dark:bg-slate-800/70" /></td>
                </tr>
              </tbody>
            </table>
          </div>

          <div v-else-if="errorMessage" class="flex flex-col items-center justify-center px-4 py-16 text-center">
            <UIcon name="i-lucide-circle-alert" class="size-9 text-red-500" />
            <p class="mt-3 text-sm font-medium text-slate-950 dark:text-white">{{ errorMessage }}</p>
            <UButton class="mt-5" color="neutral" variant="outline" icon="i-lucide-refresh-cw" @click="loadAll">
              {{ $t('common.retry') }}
            </UButton>
          </div>

          <div v-else-if="nodes.length === 0" class="flex flex-col items-center justify-center px-4 py-16 text-center">
            <UIcon name="i-lucide-inbox" class="size-9 text-slate-400" />
            <p class="mt-3 text-sm font-medium text-slate-950 dark:text-white">{{ $t('admin.forum.nodes.empty') }}</p>
          </div>

          <div v-else class="overflow-x-auto">
            <table class="min-w-[960px] w-full table-fixed border-collapse text-left">
              <thead class="bg-slate-50 text-xs font-medium uppercase text-slate-500 dark:bg-slate-950/70 dark:text-slate-400">
                <tr>
                  <th class="w-[24%] border-b border-slate-200 px-4 py-3 dark:border-slate-800">{{ $t('admin.forum.nodes.table.name') }}</th>
                  <th class="w-[16%] border-b border-slate-200 px-4 py-3 dark:border-slate-800">{{ $t('admin.forum.nodes.table.category') }}</th>
                  <th class="w-[16%] border-b border-slate-200 px-4 py-3 dark:border-slate-800">{{ $t('admin.forum.nodes.table.slug') }}</th>
                  <th class="w-[12%] border-b border-slate-200 px-4 py-3 text-right dark:border-slate-800">{{ $t('admin.forum.nodes.table.permissions') }}</th>
                  <th class="w-[12%] border-b border-slate-200 px-4 py-3 text-right dark:border-slate-800">{{ $t('admin.forum.nodes.table.counts') }}</th>
                  <th class="w-[20%] border-b border-slate-200 px-4 py-3 text-right dark:border-slate-800">{{ $t('admin.forum.nodes.table.actions') }}</th>
                </tr>
              </thead>
              <tbody>
                <tr
                  v-for="node in nodes"
                  :key="node.id"
                  class="border-b border-slate-200 transition-colors last:border-b-0 dark:border-slate-800"
                  :class="selectedId === node.id ? 'bg-cyan-50/70 dark:bg-cyan-950/30' : 'hover:bg-slate-50 dark:hover:bg-slate-950/70'"
                >
                  <td class="px-4 py-3 align-middle">
                    <button type="button" class="block max-w-full text-left" @click="selectNode(node)">
                      <span class="block truncate text-sm font-semibold text-slate-950 hover:text-cyan-700 dark:text-white dark:hover:text-cyan-300">{{ nodeName(node) }}</span>
                      <span class="mt-1 block truncate text-xs text-slate-500 dark:text-slate-400">{{ nodeDescription(node) || '-' }}</span>
                    </button>
                  </td>
                  <td class="px-4 py-3 align-middle text-sm text-slate-600 dark:text-slate-300">
                    <span class="block truncate">{{ categoryName(node.categoryId) }}</span>
                  </td>
                  <td class="px-4 py-3 align-middle">
                    <code class="block truncate rounded bg-slate-100 px-2 py-1 text-xs text-slate-700 dark:bg-slate-800 dark:text-slate-200">{{ node.slug }}</code>
                  </td>
                  <td class="px-4 py-3 text-right align-middle text-xs text-slate-600 dark:text-slate-300">
                    {{ node.minRoleRead }} / {{ node.minRoleWrite }} / {{ node.minRoleCreate }}
                  </td>
                  <td class="px-4 py-3 text-right align-middle text-xs text-slate-600 dark:text-slate-300">
                    {{ numberFormatter.format(node.topicCount) }} / {{ numberFormatter.format(node.replyCount) }}
                  </td>
                  <td class="px-4 py-3 align-middle">
                    <div class="flex items-center justify-end gap-2">
                      <UButton color="neutral" variant="outline" size="sm" icon="i-lucide-pencil" @click="selectNode(node)">{{ $t('admin.actions.edit') }}</UButton>
                      <UButton color="error" variant="soft" size="sm" icon="i-lucide-trash-2" :loading="deletingId === node.id" @click="deleteNode(node)">{{ $t('admin.actions.delete') }}</UButton>
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
              <h2 class="text-sm font-semibold text-slate-950 dark:text-white">{{ selectedId ? $t('admin.forum.nodes.form.edit') : $t('admin.forum.nodes.form.create') }}</h2>
              <UButton v-if="selectedId" color="neutral" variant="ghost" size="sm" icon="i-lucide-plus" @click="resetForm">{{ $t('admin.actions.new') }}</UButton>
            </div>
          </div>

          <form class="space-y-4 p-4" @submit.prevent="saveNode">
            <label class="block">
              <span class="text-sm font-medium text-slate-700 dark:text-slate-200">{{ $t('admin.forum.nodes.form.category') }}</span>
              <select v-model.number="form.categoryId" class="mt-1 h-10 w-full rounded-md border border-slate-200 bg-white px-3 text-sm text-slate-950 outline-none transition focus:border-cyan-400 focus:ring-2 focus:ring-cyan-100 dark:border-slate-700 dark:bg-slate-950 dark:text-white dark:focus:border-cyan-500 dark:focus:ring-cyan-950" :disabled="saving || categories.length === 0">
                <option :value="0">{{ $t('admin.forum.nodes.form.selectCategory') }}</option>
                <option v-for="category in categories" :key="category.id" :value="category.id">{{ forumCategoryName(category) }}</option>
              </select>
            </label>

            <label class="block">
              <span class="text-sm font-medium text-slate-700 dark:text-slate-200">{{ $t('admin.forum.nodes.form.slug') }}</span>
              <input v-model="form.slug" class="mt-1 h-10 w-full rounded-md border border-slate-200 bg-white px-3 text-sm text-slate-950 outline-none transition focus:border-cyan-400 focus:ring-2 focus:ring-cyan-100 dark:border-slate-700 dark:bg-slate-950 dark:text-white dark:focus:border-cyan-500 dark:focus:ring-cyan-950" :disabled="saving">
            </label>

            <label v-for="field in nameFields" :key="field.key" class="block">
              <span class="text-sm font-medium text-slate-700 dark:text-slate-200">{{ field.label }}</span>
              <input v-model="form[field.key]" class="mt-1 h-10 w-full rounded-md border border-slate-200 bg-white px-3 text-sm text-slate-950 outline-none transition focus:border-cyan-400 focus:ring-2 focus:ring-cyan-100 dark:border-slate-700 dark:bg-slate-950 dark:text-white dark:focus:border-cyan-500 dark:focus:ring-cyan-950" :disabled="saving">
            </label>

            <label v-for="field in descFields" :key="field.key" class="block">
              <span class="text-sm font-medium text-slate-700 dark:text-slate-200">{{ field.label }}</span>
              <textarea v-model="form[field.key]" rows="2" class="mt-1 w-full resize-none rounded-md border border-slate-200 bg-white px-3 py-2 text-sm text-slate-950 outline-none transition focus:border-cyan-400 focus:ring-2 focus:ring-cyan-100 dark:border-slate-700 dark:bg-slate-950 dark:text-white dark:focus:border-cyan-500 dark:focus:ring-cyan-950" :disabled="saving" />
            </label>

            <div class="grid grid-cols-2 gap-3">
              <label v-for="field in numberFields" :key="field.key" class="block">
                <span class="text-sm font-medium text-slate-700 dark:text-slate-200">{{ field.label }}</span>
                <input v-model.number="form[field.key]" type="number" min="0" class="mt-1 h-10 w-full rounded-md border border-slate-200 bg-white px-3 text-sm text-slate-950 outline-none transition focus:border-cyan-400 focus:ring-2 focus:ring-cyan-100 dark:border-slate-700 dark:bg-slate-950 dark:text-white dark:focus:border-cyan-500 dark:focus:ring-cyan-950" :disabled="saving">
              </label>
            </div>

            <label class="block">
              <span class="text-sm font-medium text-slate-700 dark:text-slate-200">{{ $t('admin.forum.nodes.form.moderators') }}</span>
              <input v-model="form.moderators" class="mt-1 h-10 w-full rounded-md border border-slate-200 bg-white px-3 text-sm text-slate-950 outline-none transition focus:border-cyan-400 focus:ring-2 focus:ring-cyan-100 dark:border-slate-700 dark:bg-slate-950 dark:text-white dark:focus:border-cyan-500 dark:focus:ring-cyan-950" :disabled="saving" placeholder="1,2,3">
            </label>

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
import type { AdminForumCategory, AdminForumNode, AdminForumNodeInput } from '~/composables/useAdmin'
import { localizeI18nName } from '~/utils/format'

definePageMeta({ layout: 'admin', middleware: 'admin' })

const { t, locale } = useI18n()
const localePath = useLocalePath()
const toast = useToast()
const adminApi = useAdmin()

useHead({ title: `${t('admin.forum.nodes.title')} - NextPT` })

const categories = ref<AdminForumCategory[]>([])
const nodes = ref<AdminForumNode[]>([])
const pending = ref(false)
const saving = ref(false)
const deletingId = ref<number | null>(null)
const selectedId = ref<number | null>(null)
const errorMessage = ref('')
const formError = ref('')

const form = reactive({
  categoryId: 0,
  slug: '',
  nameZhCN: '',
  nameZhTW: '',
  nameEnUS: '',
  descZhCN: '',
  descZhTW: '',
  descEnUS: '',
  sortOrder: 0,
  minRoleRead: 0,
  minRoleWrite: 0,
  minRoleCreate: 0,
  moderators: ''
})

const numberFormatter = computed(() => new Intl.NumberFormat(locale.value))
const categoryMap = computed(() => new Map(categories.value.map((item) => [item.id, item])))
const canSubmit = computed(() => Boolean(form.categoryId > 0 && form.slug.trim() && form.nameZhCN.trim() && !saving.value))
const nameFields = computed(() => [
  { key: 'nameZhCN' as const, label: t('admin.forum.nodes.form.nameZhCN') },
  { key: 'nameZhTW' as const, label: t('admin.forum.nodes.form.nameZhTW') },
  { key: 'nameEnUS' as const, label: t('admin.forum.nodes.form.nameEnUS') }
])
const descFields = computed(() => [
  { key: 'descZhCN' as const, label: t('admin.forum.nodes.form.descZhCN') },
  { key: 'descZhTW' as const, label: t('admin.forum.nodes.form.descZhTW') },
  { key: 'descEnUS' as const, label: t('admin.forum.nodes.form.descEnUS') }
])
const numberFields = computed(() => [
  { key: 'sortOrder' as const, label: t('admin.forum.nodes.form.sortOrder') },
  { key: 'minRoleRead' as const, label: t('admin.forum.nodes.form.minRoleRead') },
  { key: 'minRoleWrite' as const, label: t('admin.forum.nodes.form.minRoleWrite') },
  { key: 'minRoleCreate' as const, label: t('admin.forum.nodes.form.minRoleCreate') }
])

onMounted(loadAll)

async function loadAll() {
  pending.value = true
  errorMessage.value = ''
  try {
    const [categoryData, nodeData] = await Promise.all([
      adminApi.listForumCategories(),
      adminApi.listForumNodes()
    ])
    categories.value = (categoryData.categories || []).sort((a, b) => a.sortOrder - b.sortOrder || a.id - b.id)
    nodes.value = (nodeData.nodes || []).sort((a, b) => a.categoryId - b.categoryId || a.sortOrder - b.sortOrder || a.id - b.id)
  } catch (error: unknown) {
    errorMessage.value = error instanceof ApiError ? error.message : t('common.requestFailed')
  } finally {
    pending.value = false
  }
}

function forumCategoryName(category: AdminForumCategory) {
  return localizeI18nName(category.nameI18N, locale.value, `#${category.id}`)
}

function categoryName(id: number) {
  const category = categoryMap.value.get(id)
  return category ? forumCategoryName(category) : `#${id}`
}

function nodeName(node: AdminForumNode) {
  return localizeI18nName(node.nameI18N, locale.value, node.slug || `#${node.id}`)
}

function nodeDescription(node: AdminForumNode) {
  return localizeI18nName(node.descI18N, locale.value, '')
}

function normalizeModerators(value: unknown) {
  if (!Array.isArray(value)) return ''
  return value.map((item) => Number(item)).filter((item) => Number.isInteger(item) && item > 0).join(',')
}

function parseModerators(value: string) {
  return value
    .split(',')
    .map((item) => Number(item.trim()))
    .filter((item, index, list) => Number.isInteger(item) && item > 0 && list.indexOf(item) === index)
}

function selectNode(node: AdminForumNode) {
  selectedId.value = node.id
  form.categoryId = node.categoryId
  form.slug = node.slug
  form.nameZhCN = String(node.nameI18N?.['zh-CN'] || '')
  form.nameZhTW = String(node.nameI18N?.['zh-TW'] || '')
  form.nameEnUS = String(node.nameI18N?.['en-US'] || '')
  form.descZhCN = String(node.descI18N?.['zh-CN'] || '')
  form.descZhTW = String(node.descI18N?.['zh-TW'] || '')
  form.descEnUS = String(node.descI18N?.['en-US'] || '')
  form.sortOrder = node.sortOrder
  form.minRoleRead = node.minRoleRead
  form.minRoleWrite = node.minRoleWrite
  form.minRoleCreate = node.minRoleCreate
  form.moderators = normalizeModerators(node.moderators)
  formError.value = ''
}

function resetForm() {
  selectedId.value = null
  form.categoryId = categories.value[0]?.id || 0
  form.slug = ''
  form.nameZhCN = ''
  form.nameZhTW = ''
  form.nameEnUS = ''
  form.descZhCN = ''
  form.descZhTW = ''
  form.descEnUS = ''
  form.sortOrder = 0
  form.minRoleRead = 0
  form.minRoleWrite = 0
  form.minRoleCreate = 0
  form.moderators = ''
  formError.value = ''
}

function buildInput(): AdminForumNodeInput {
  const zhCN = form.nameZhCN.trim()
  const descZhCN = form.descZhCN.trim()
  return {
    categoryId: Number(form.categoryId),
    slug: form.slug.trim(),
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
    minRoleRead: Number(form.minRoleRead || 0),
    minRoleWrite: Number(form.minRoleWrite || 0),
    minRoleCreate: Number(form.minRoleCreate || 0),
    moderators: parseModerators(form.moderators)
  }
}

async function saveNode() {
  if (!canSubmit.value) return
  saving.value = true
  formError.value = ''
  try {
    const input = buildInput()
    if (selectedId.value) {
      await adminApi.updateForumNode(selectedId.value, input)
    } else {
      await adminApi.createForumNode(input)
    }
    toast.add({ title: t('admin.forum.nodes.saved') })
    resetForm()
    await loadAll()
  } catch (error: unknown) {
    formError.value = error instanceof ApiError ? error.message : t('common.requestFailed')
  } finally {
    saving.value = false
  }
}

async function deleteNode(node: AdminForumNode) {
  if (!window.confirm(t('admin.forum.nodes.confirmDelete', { name: nodeName(node) }))) return
  deletingId.value = node.id
  try {
    await adminApi.deleteForumNode(node.id)
    if (selectedId.value === node.id) resetForm()
    toast.add({ title: t('admin.forum.nodes.deleted') })
    await loadAll()
  } catch (error: unknown) {
    toast.add({ color: 'error', title: error instanceof ApiError ? error.message : t('common.requestFailed') })
  } finally {
    deletingId.value = null
  }
}
</script>
