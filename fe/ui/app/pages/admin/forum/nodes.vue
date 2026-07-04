<template>
  <div class="min-h-[calc(100vh-4rem)] bg-slate-50 py-6 dark:bg-slate-950">
    <div class="w-full px-3 sm:px-4 lg:px-5">
      <div class="grid gap-4 xl:grid-cols-[minmax(760px,940px)_minmax(560px,1fr)] 2xl:grid-cols-[minmax(820px,1020px)_minmax(600px,1fr)]">
        <section class="overflow-hidden rounded-lg border border-slate-200 bg-white dark:border-slate-800 dark:bg-slate-900">
          <div class="border-b border-slate-200 px-4 py-3 dark:border-slate-800">
            <div class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
              <div class="flex items-center gap-2">
                <UIcon name="i-lucide-git-branch" class="size-5 text-cyan-600 dark:text-cyan-300" />
                <h2 class="text-sm font-semibold text-slate-950 dark:text-white">{{ $t('admin.forum.nodes.list') }}</h2>
              </div>
              <div class="flex flex-wrap items-center gap-2 text-xs">
                <span class="inline-flex h-7 items-center gap-1.5 rounded-md bg-slate-100 px-2.5 font-medium text-slate-600 dark:bg-slate-800 dark:text-slate-300">
                  <span>{{ $t('admin.forum.nodes.stats.total') }}</span>
                  <span class="font-semibold text-slate-950 dark:text-white">{{ numberFormatter.format(nodes.length) }}</span>
                </span>
                <span class="inline-flex h-7 items-center gap-1.5 rounded-md bg-cyan-50 px-2.5 font-medium text-cyan-700 dark:bg-cyan-950/40 dark:text-cyan-300">
                  <span>{{ $t('admin.forum.nodes.stats.categories') }}</span>
                  <span class="font-semibold">{{ numberFormatter.format(categories.length) }}</span>
                </span>
              </div>
            </div>
          </div>

          <div v-if="pending" class="overflow-x-auto">
            <table class="min-w-[1080px] w-full table-fixed border-collapse text-left">
              <thead class="bg-slate-50 text-xs font-medium uppercase text-slate-500 dark:bg-slate-950/70 dark:text-slate-400">
                <tr>
                  <th class="w-[23%] border-b border-slate-200 px-3 py-2.5 dark:border-slate-800">{{ $t('admin.forum.nodes.table.name') }}</th>
                  <th class="w-[15%] border-b border-slate-200 px-3 py-2.5 dark:border-slate-800">{{ $t('admin.forum.nodes.table.category') }}</th>
                  <th class="w-[14%] border-b border-slate-200 px-3 py-2.5 dark:border-slate-800">{{ $t('admin.forum.nodes.table.slug') }}</th>
                  <th class="w-[26%] border-b border-slate-200 px-3 py-2.5 dark:border-slate-800">{{ $t('admin.forum.nodes.table.permissions') }}</th>
                  <th class="w-[12%] border-b border-slate-200 px-3 py-2.5 text-right dark:border-slate-800">{{ $t('admin.forum.nodes.table.counts') }}</th>
                  <th class="w-[10%] border-b border-slate-200 px-3 py-2.5 text-right dark:border-slate-800">{{ $t('admin.forum.nodes.table.actions') }}</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="index in 5" :key="index" class="border-b border-slate-200 last:border-b-0 dark:border-slate-800">
                  <td class="px-3 py-3">
                    <div class="h-4 w-44 animate-pulse rounded bg-slate-200 dark:bg-slate-800" />
                    <div class="mt-2 h-3 w-32 animate-pulse rounded bg-slate-100 dark:bg-slate-800/70" />
                  </td>
                  <td class="px-3 py-3"><div class="h-4 w-28 animate-pulse rounded bg-slate-100 dark:bg-slate-800/70" /></td>
                  <td class="px-3 py-3"><div class="h-4 w-24 animate-pulse rounded bg-slate-100 dark:bg-slate-800/70" /></td>
                  <td class="px-3 py-3"><div class="ml-auto h-4 w-20 animate-pulse rounded bg-slate-100 dark:bg-slate-800/70" /></td>
                  <td class="px-3 py-3"><div class="ml-auto h-4 w-16 animate-pulse rounded bg-slate-100 dark:bg-slate-800/70" /></td>
                  <td class="px-3 py-3"><div class="ml-auto h-4 w-20 animate-pulse rounded bg-slate-100 dark:bg-slate-800/70" /></td>
                </tr>
              </tbody>
            </table>
          </div>

          <div v-else-if="errorMessage" class="flex flex-col items-center justify-center px-4 py-16 text-center">
            <UIcon name="i-lucide-circle-alert" class="size-9 text-red-500" />
            <p class="mt-3 text-sm font-medium text-slate-950 dark:text-white">{{ errorMessage }}</p>
          </div>

          <div v-else-if="nodes.length === 0" class="flex flex-col items-center justify-center px-4 py-16 text-center">
            <UIcon name="i-lucide-inbox" class="size-9 text-slate-400" />
            <p class="mt-3 text-sm font-medium text-slate-950 dark:text-white">{{ $t('admin.forum.nodes.empty') }}</p>
          </div>

          <div v-else class="overflow-x-auto">
            <table class="min-w-[1080px] w-full table-fixed border-collapse text-left">
              <thead class="bg-slate-50 text-xs font-medium uppercase text-slate-500 dark:bg-slate-950/70 dark:text-slate-400">
                <tr>
                  <th class="w-[23%] border-b border-slate-200 px-3 py-2.5 dark:border-slate-800">{{ $t('admin.forum.nodes.table.name') }}</th>
                  <th class="w-[15%] border-b border-slate-200 px-3 py-2.5 dark:border-slate-800">{{ $t('admin.forum.nodes.table.category') }}</th>
                  <th class="w-[14%] border-b border-slate-200 px-3 py-2.5 dark:border-slate-800">{{ $t('admin.forum.nodes.table.slug') }}</th>
                  <th class="w-[26%] border-b border-slate-200 px-3 py-2.5 dark:border-slate-800">{{ $t('admin.forum.nodes.table.permissions') }}</th>
                  <th class="w-[12%] border-b border-slate-200 px-3 py-2.5 text-right dark:border-slate-800">{{ $t('admin.forum.nodes.table.counts') }}</th>
                  <th class="w-[10%] border-b border-slate-200 px-3 py-2.5 text-right dark:border-slate-800">{{ $t('admin.forum.nodes.table.actions') }}</th>
                </tr>
              </thead>
              <tbody>
                <tr
                  v-for="node in nodes"
                  :key="node.id"
                  class="border-b border-slate-200 transition-colors last:border-b-0 dark:border-slate-800"
                  :class="selectedId === node.id ? 'bg-cyan-50/70 dark:bg-cyan-950/30' : 'hover:bg-slate-50 dark:hover:bg-slate-950/70'"
                >
                  <td class="px-3 py-2.5 align-middle">
                    <button type="button" class="block max-w-full text-left" @click="selectNode(node)">
                      <span class="block truncate text-sm font-semibold text-slate-950 hover:text-cyan-700 dark:text-white dark:hover:text-cyan-300">
                        {{ nodeName(node) }}
                      </span>
                      <span class="mt-1 block truncate text-xs text-slate-500 dark:text-slate-400">
                        {{ nodeDescription(node) || '-' }}
                      </span>
                    </button>
                  </td>
                  <td class="px-3 py-2.5 align-middle text-sm text-slate-600 dark:text-slate-300">
                    <span class="block truncate">{{ categoryName(node.categoryId) }}</span>
                  </td>
                  <td class="px-3 py-2.5 align-middle">
                    <code class="block truncate text-xs text-slate-600 dark:text-slate-300">{{ node.slug }}</code>
                  </td>
                  <td class="px-3 py-2.5 align-middle text-xs text-slate-600 dark:text-slate-300">
                    <div class="flex flex-wrap items-center gap-1.5">
                      <span
                        v-for="item in nodePermissionSummary(node)"
                        :key="item.key"
                        class="inline-flex shrink-0 items-center gap-1.5"
                      >
                        <span class="shrink-0 text-slate-400 dark:text-slate-500">{{ item.label }}</span>
                        <span class="whitespace-nowrap rounded-md bg-slate-100 px-1.5 py-0.5 font-medium text-slate-700 dark:bg-slate-800 dark:text-slate-200">
                          {{ item.role }}
                        </span>
                      </span>
                    </div>
                  </td>
                  <td class="px-3 py-2.5 text-right align-middle text-xs text-slate-600 dark:text-slate-300">
                    {{ numberFormatter.format(node.topicCount) }} / {{ numberFormatter.format(node.replyCount) }}
                  </td>
                  <td class="px-3 py-2.5 align-middle">
                    <div class="flex items-center justify-end gap-2">
                      <UTooltip
                        :text="$t('admin.actions.edit')"
                        :content="{ side: 'top', sideOffset: 8 }"
                        :delay-duration="120"
                      >
                        <UButton
                          color="neutral"
                          variant="ghost"
                          size="sm"
                          icon="i-lucide-pencil"
                          :aria-label="$t('admin.actions.edit')"
                          @click="selectNode(node)"
                        />
                      </UTooltip>
                      <UPopover
                        :content="{ side: 'top', align: 'end', sideOffset: 8 }"
                        :ui="{ content: 'w-72 p-3' }"
                      >
                        <UButton
                          color="error"
                          variant="ghost"
                          size="sm"
                          icon="i-lucide-trash-2"
                          :loading="deletingId === node.id"
                          :disabled="deletingId !== null"
                          :aria-label="$t('admin.actions.delete')"
                        />

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
                              <UButton color="error" size="xs" type="button" icon="i-lucide-trash-2" :loading="deletingId === node.id" :disabled="deletingId !== null" @click="deleteNode(node, close)">
                                {{ $t('admin.actions.delete') }}
                              </UButton>
                            </div>
                          </div>
                        </template>
                      </UPopover>
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
                {{ selectedId ? $t('admin.forum.nodes.form.edit') : $t('admin.forum.nodes.form.create') }}
              </h2>
              <div class="flex h-8 min-w-20 items-center justify-end">
                <UButton v-if="selectedId" color="primary" variant="soft" size="sm" icon="i-lucide-plus" @click="startCreate">
                  {{ $t('admin.actions.new') }}
                </UButton>
              </div>
            </div>
          </div>

          <form class="space-y-5 p-4" @submit.prevent="saveNode">
            <div>
              <div class="border-b border-slate-200 pb-2 dark:border-slate-800">
                <span class="text-sm font-semibold text-slate-950 dark:text-white">{{ $t('admin.forum.nodes.form.basic') }}</span>
              </div>
              <div class="mt-3 grid gap-3 sm:grid-cols-[minmax(0,1fr)_minmax(0,1fr)_112px]">
                <label class="block">
                  <span class="text-sm font-medium text-slate-700 dark:text-slate-200">{{ $t('admin.forum.nodes.form.category') }}</span>
                  <select
                    v-model.number="form.categoryId"
                    class="mt-1 h-10 w-full rounded-md border border-slate-200 bg-white px-3 text-sm text-slate-950 outline-none transition focus:border-cyan-400 focus:ring-2 focus:ring-cyan-100 dark:border-slate-700 dark:bg-slate-950 dark:text-white dark:focus:border-cyan-500 dark:focus:ring-cyan-950"
                    :disabled="saving || categories.length === 0"
                  >
                    <option :value="0">{{ $t('admin.forum.nodes.form.selectCategory') }}</option>
                    <option v-for="category in categories" :key="category.id" :value="category.id">{{ forumCategoryName(category) }}</option>
                  </select>
                </label>

                <label class="block">
                  <span class="text-sm font-medium text-slate-700 dark:text-slate-200">{{ $t('admin.forum.nodes.form.slug') }}</span>
                  <input
                    v-model="form.slug"
                    class="mt-1 h-10 w-full rounded-md border border-slate-200 bg-white px-3 text-sm text-slate-950 outline-none transition focus:border-cyan-400 focus:ring-2 focus:ring-cyan-100 dark:border-slate-700 dark:bg-slate-950 dark:text-white dark:focus:border-cyan-500 dark:focus:ring-cyan-950"
                    :disabled="saving"
                  >
                </label>

                <label class="block">
                  <span class="text-sm font-medium text-slate-700 dark:text-slate-200">{{ $t('admin.forum.nodes.form.sortOrder') }}</span>
                  <input
                    v-model.number="form.sortOrder"
                    type="number"
                    class="mt-1 h-10 w-full rounded-md border border-slate-200 bg-white px-3 text-sm text-slate-950 outline-none transition focus:border-cyan-400 focus:ring-2 focus:ring-cyan-100 dark:border-slate-700 dark:bg-slate-950 dark:text-white dark:focus:border-cyan-500 dark:focus:ring-cyan-950"
                    :disabled="saving"
                  >
                </label>
              </div>
            </div>

            <div>
              <div class="flex items-center justify-between gap-3 border-b border-slate-200 pb-2 dark:border-slate-800">
                <span class="text-sm font-semibold text-slate-950 dark:text-white">{{ $t('admin.forum.nodes.form.name') }}</span>
                <span class="text-xs text-slate-500 dark:text-slate-400">
                  {{ $t('admin.forum.nodes.form.primaryLocale', { locale: primaryLocaleLabel }) }}
                </span>
              </div>
              <div class="mt-3 grid gap-3 sm:grid-cols-2 2xl:grid-cols-3">
                <label v-for="item in localeOptions" :key="item.code" class="block">
                  <span class="flex items-center justify-between gap-2 text-sm font-medium text-slate-700 dark:text-slate-200">
                    <span>{{ item.name }}</span>
                    <span class="text-xs font-normal text-slate-400 dark:text-slate-500">{{ item.code }}</span>
                  </span>
                  <input
                    v-model="form.names[item.code]"
                    class="mt-1 h-10 w-full rounded-md border border-slate-200 bg-white px-3 text-sm text-slate-950 outline-none transition focus:border-cyan-400 focus:ring-2 focus:ring-cyan-100 dark:border-slate-700 dark:bg-slate-950 dark:text-white dark:focus:border-cyan-500 dark:focus:ring-cyan-950"
                    :disabled="saving"
                  >
                </label>
              </div>
            </div>

            <div>
              <div class="border-b border-slate-200 pb-2 dark:border-slate-800">
                <span class="text-sm font-semibold text-slate-950 dark:text-white">{{ $t('admin.forum.nodes.form.description') }}</span>
              </div>
              <div class="mt-3 grid gap-3 sm:grid-cols-2 2xl:grid-cols-3">
                <label v-for="item in localeOptions" :key="item.code" class="block">
                  <span class="flex items-center justify-between gap-2 text-sm font-medium text-slate-700 dark:text-slate-200">
                    <span>{{ item.name }}</span>
                    <span class="text-xs font-normal text-slate-400 dark:text-slate-500">{{ item.code }}</span>
                  </span>
                  <textarea
                    v-model="form.descs[item.code]"
                    rows="3"
                    class="mt-1 w-full resize-none rounded-md border border-slate-200 bg-white px-3 py-2 text-sm text-slate-950 outline-none transition focus:border-cyan-400 focus:ring-2 focus:ring-cyan-100 dark:border-slate-700 dark:bg-slate-950 dark:text-white dark:focus:border-cyan-500 dark:focus:ring-cyan-950"
                    :disabled="saving"
                  />
                </label>
              </div>
            </div>

            <div>
              <div class="border-b border-slate-200 pb-2 dark:border-slate-800">
                <span class="text-sm font-semibold text-slate-950 dark:text-white">{{ $t('admin.forum.nodes.form.permissions') }}</span>
              </div>
              <div class="mt-3 grid gap-3 sm:grid-cols-3">
                <label v-for="field in permissionFields" :key="field.key" class="block">
                  <span class="text-sm font-medium text-slate-700 dark:text-slate-200">{{ field.label }}</span>
                  <AdminIamRoleLevelSelect
                    v-model.number="form[field.key]"
                    :roles="roles"
                    :disabled="saving"
                  />
                </label>
              </div>
            </div>

            <label class="block">
              <span class="text-sm font-medium text-slate-700 dark:text-slate-200">{{ $t('admin.forum.nodes.form.moderators') }}</span>
              <input
                v-model="form.moderators"
                class="mt-1 h-10 w-full rounded-md border border-slate-200 bg-white px-3 text-sm text-slate-950 outline-none transition focus:border-cyan-400 focus:ring-2 focus:ring-cyan-100 dark:border-slate-700 dark:bg-slate-950 dark:text-white dark:focus:border-cyan-500 dark:focus:ring-cyan-950"
                :disabled="saving"
                placeholder="1,2,3"
              >
            </label>

            <p v-if="formError" class="rounded-md border border-red-200 bg-red-50 px-3 py-2 text-sm text-red-700 dark:border-red-900 dark:bg-red-950/40 dark:text-red-200">
              {{ formError }}
            </p>

            <div class="flex flex-col gap-2 border-t border-slate-200 pt-4 sm:flex-row sm:justify-end dark:border-slate-800">
              <UButton type="button" color="neutral" variant="outline" icon="i-lucide-rotate-ccw" :disabled="saving || !isFormDirty" @click="resetFormChanges">
                {{ selectedId ? $t('admin.actions.discardChanges') : $t('admin.actions.reset') }}
              </UButton>
              <UButton type="submit" color="primary" icon="i-lucide-save" :loading="saving" :disabled="!canSubmit">
                {{ $t('common.save') }}
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
import type { AdminForumCategory, AdminForumNode, AdminForumNodeInput, AdminIamRole } from '~/composables/useAdmin'
import { localizeI18nName } from '~/utils/format'

interface LocaleOption {
  code: string
  name: string
}

definePageMeta({ layout: 'admin', middleware: 'admin' })

const { t, locale, locales } = useI18n()
const toast = useToast()
const adminApi = useAdmin()

useHead({ title: t('admin.forum.nodes.title') })

const categories = ref<AdminForumCategory[]>([])
const nodes = ref<AdminForumNode[]>([])
const roles = ref<AdminIamRole[]>([])
const pending = ref(false)
const saving = ref(false)
const deletingId = ref<number | null>(null)
const selectedId = ref<number | null>(null)
const errorMessage = ref('')
const formError = ref('')
const originalFormSnapshot = ref('')

const form = reactive({
  categoryId: 0,
  slug: '',
  names: {} as Record<string, string>,
  descs: {} as Record<string, string>,
  sortOrder: 0,
  minRoleRead: 0,
  minRoleWrite: 0,
  minRoleCreate: 0,
  moderators: ''
})

const localeOptions = computed<LocaleOption[]>(() => {
  return (locales.value as Array<string | { code?: string, language?: string, name?: string }>)
    .map((item) => {
      if (typeof item === 'string') return { code: item, name: item }
      const code = String(item.code || item.language || '')
      return { code, name: String(item.name || code) }
    })
    .filter((item) => item.code)
})
const primaryLocaleCode = computed(() => localeOptions.value[0]?.code || locale.value)
const primaryLocaleLabel = computed(() => {
  const primary = localeOptions.value.find((item) => item.code === primaryLocaleCode.value)
  return primary ? `${primary.name} (${primary.code})` : primaryLocaleCode.value
})
const numberFormatter = computed(() => new Intl.NumberFormat(locale.value))
const categoryMap = computed(() => new Map(categories.value.map((item) => [item.id, item])))
const { roleNameByLevel } = useAdminIamRoleLevels(roles)
const selectedNode = computed(() => nodes.value.find((item) => item.id === selectedId.value) || null)
const isFormDirty = computed(() => formSnapshot() !== originalFormSnapshot.value)
const canSubmit = computed(() => Boolean(form.categoryId > 0 && form.slug.trim() && nameValue(primaryLocaleCode.value).trim() && isFormDirty.value && !saving.value))
const permissionFields = computed(() => [
  { key: 'minRoleRead' as const, label: t('admin.forum.nodes.form.minRoleRead') },
  { key: 'minRoleWrite' as const, label: t('admin.forum.nodes.form.minRoleWrite') },
  { key: 'minRoleCreate' as const, label: t('admin.forum.nodes.form.minRoleCreate') }
])

watch(localeOptions, () => {
  ensureFormLocales()
}, { immediate: true })

onMounted(() => {
  originalFormSnapshot.value = formSnapshot()
  loadAll()
})

async function loadAll() {
  pending.value = true
  errorMessage.value = ''
  try {
    const [categoryData, nodeData, roleData] = await Promise.all([
      adminApi.listForumCategories(),
      adminApi.listForumNodes(),
      adminApi.listIamRoles()
    ])
    categories.value = (categoryData.categories || []).sort((a, b) => a.sortOrder - b.sortOrder || a.id - b.id)
    nodes.value = (nodeData.nodes || []).sort((a, b) => a.categoryId - b.categoryId || a.sortOrder - b.sortOrder || a.id - b.id)
    roles.value = roleData.roles || []
    if (!selectedId.value && !form.categoryId && categories.value.length > 0) {
      form.categoryId = categories.value[0].id
      originalFormSnapshot.value = formSnapshot()
    }
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

function nodePermissionSummary(node: AdminForumNode) {
  return [
    { key: 'read', label: t('admin.forum.nodes.permissionShort.read'), role: roleNameByLevel(node.minRoleRead) },
    { key: 'write', label: t('admin.forum.nodes.permissionShort.write'), role: roleNameByLevel(node.minRoleWrite) },
    { key: 'create', label: t('admin.forum.nodes.permissionShort.create'), role: roleNameByLevel(node.minRoleCreate) }
  ]
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
  setI18nValues(form.names, node.nameI18N)
  setI18nValues(form.descs, node.descI18N)
  form.sortOrder = node.sortOrder
  form.minRoleRead = node.minRoleRead
  form.minRoleWrite = node.minRoleWrite
  form.minRoleCreate = node.minRoleCreate
  form.moderators = normalizeModerators(node.moderators)
  formError.value = ''
  originalFormSnapshot.value = formSnapshot()
}

function startCreate() {
  selectedId.value = null
  form.categoryId = categories.value[0]?.id || 0
  form.slug = ''
  setI18nValues(form.names)
  setI18nValues(form.descs)
  form.sortOrder = 0
  form.minRoleRead = 0
  form.minRoleWrite = 0
  form.minRoleCreate = 0
  form.moderators = ''
  formError.value = ''
  originalFormSnapshot.value = formSnapshot()
}

function resetFormChanges() {
  const node = selectedNode.value
  if (node) {
    selectNode(node)
    return
  }
  startCreate()
}

function formSnapshot() {
  return JSON.stringify({
    categoryId: Number(form.categoryId || 0),
    slug: form.slug,
    names: localeOptions.value.reduce<Record<string, string>>((names, item) => {
      names[item.code] = nameValue(item.code)
      return names
    }, {}),
    descs: localeOptions.value.reduce<Record<string, string>>((descs, item) => {
      descs[item.code] = descValue(item.code)
      return descs
    }, {}),
    sortOrder: Number(form.sortOrder || 0),
    minRoleRead: Number(form.minRoleRead || 0),
    minRoleWrite: Number(form.minRoleWrite || 0),
    minRoleCreate: Number(form.minRoleCreate || 0),
    moderators: form.moderators
  })
}

function buildInput(): AdminForumNodeInput {
  const primaryName = nameValue(primaryLocaleCode.value).trim()
  const fallbackDesc = descValue(primaryLocaleCode.value).trim() || firstFilledDesc()
  return {
    categoryId: Number(form.categoryId),
    slug: form.slug.trim(),
    nameI18N: localeOptions.value.reduce<Record<string, string>>((names, item) => {
      names[item.code] = nameValue(item.code).trim() || primaryName
      return names
    }, {}),
    descI18N: localeOptions.value.reduce<Record<string, string>>((descs, item) => {
      descs[item.code] = descValue(item.code).trim() || fallbackDesc
      return descs
    }, {}),
    sortOrder: Number(form.sortOrder || 0),
    minRoleRead: Number(form.minRoleRead || 0),
    minRoleWrite: Number(form.minRoleWrite || 0),
    minRoleCreate: Number(form.minRoleCreate || 0),
    moderators: parseModerators(form.moderators)
  }
}

function ensureFormLocales() {
  ensureI18nLocales(form.names)
  ensureI18nLocales(form.descs)
}

function ensureI18nLocales(target: Record<string, string>) {
  for (const item of localeOptions.value) {
    if (!(item.code in target)) {
      target[item.code] = ''
    }
  }
}

function setI18nValues(target: Record<string, string>, values?: Record<string, unknown> | null) {
  const supportedCodes = new Set(localeOptions.value.map((item) => item.code))
  for (const code of Object.keys(target)) {
    if (!supportedCodes.has(code)) {
      delete target[code]
    }
  }
  for (const item of localeOptions.value) {
    const value = values?.[item.code]
    target[item.code] = typeof value === 'string' ? value : ''
  }
}

function nameValue(code: string) {
  return form.names[code] || ''
}

function descValue(code: string) {
  return form.descs[code] || ''
}

function firstFilledDesc() {
  for (const item of localeOptions.value) {
    const value = descValue(item.code).trim()
    if (value) return value
  }
  return ''
}

async function saveNode() {
  if (!canSubmit.value) return
  saving.value = true
  formError.value = ''
  try {
    const input = buildInput()
    const editingId = selectedId.value
    if (editingId) {
      await adminApi.updateForumNode(editingId, input)
    } else {
      await adminApi.createForumNode(input)
    }
    toast.add({ title: t('admin.forum.nodes.saved') })
    await loadAll()
    if (editingId) {
      const updatedNode = nodes.value.find((item) => item.id === editingId)
      if (updatedNode) {
        selectNode(updatedNode)
      } else {
        startCreate()
      }
    } else {
      startCreate()
    }
  } catch (error: unknown) {
    formError.value = error instanceof ApiError || error instanceof Error ? error.message : t('common.requestFailed')
  } finally {
    saving.value = false
  }
}

async function deleteNode(node: AdminForumNode, close?: () => void) {
  deletingId.value = node.id
  try {
    await adminApi.deleteForumNode(node.id)
    if (selectedId.value === node.id) {
      startCreate()
    }
    toast.add({ title: t('admin.forum.nodes.deleted') })
    close?.()
    await loadAll()
  } catch (error: unknown) {
    toast.add({ color: 'error', title: error instanceof ApiError ? error.message : t('common.requestFailed') })
  } finally {
    deletingId.value = null
  }
}
</script>
