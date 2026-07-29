<template>
  <div class="min-h-[calc(100vh-4rem)] bg-slate-50 py-6 dark:bg-slate-950">
    <div class="w-full px-3 sm:px-4 lg:px-5">
      <div class="grid gap-4 xl:grid-cols-[minmax(420px,620px)_minmax(620px,1fr)]">
        <section class="overflow-hidden rounded-lg border border-slate-200 bg-white dark:border-slate-800 dark:bg-slate-900">
          <div class="flex h-12 items-center justify-between border-b border-slate-200 px-4 dark:border-slate-800">
            <div class="flex items-center gap-2">
              <UIcon name="i-lucide-tags" class="size-4 text-sky-600 dark:text-sky-300" />
              <h2 class="text-sm font-semibold text-slate-950 dark:text-white">{{ $t('admin.catalog.tags.list') }}</h2>
            </div>
            <span class="text-xs tabular-nums text-slate-500 dark:text-slate-400">{{ numberFormatter.format(groups.length) }}</span>
          </div>

          <div v-if="pending" class="divide-y divide-slate-200 dark:divide-slate-800">
            <div v-for="index in 5" :key="index" class="space-y-2 px-4 py-3">
              <div class="h-4 w-36 animate-pulse rounded bg-slate-200 dark:bg-slate-800" />
              <div class="h-3 w-52 animate-pulse rounded bg-slate-100 dark:bg-slate-800/70" />
            </div>
          </div>
          <div v-else-if="errorMessage" class="px-4 py-16 text-center text-sm text-red-600 dark:text-red-300">{{ errorMessage }}</div>
          <div v-else-if="groups.length === 0" class="px-4 py-16 text-center text-sm text-slate-500 dark:text-slate-400">{{ $t('admin.catalog.tags.empty') }}</div>
          <div v-else class="divide-y divide-slate-200 dark:divide-slate-800">
            <button
              v-for="group in groups"
              :key="group.id"
              type="button"
              class="flex w-full items-center justify-between gap-4 px-4 py-3 text-left transition-colors"
              :class="selectedGroupId === group.id ? 'bg-sky-50/70 dark:bg-sky-950/30' : 'hover:bg-slate-50 dark:hover:bg-slate-950/70'"
              @click="selectGroup(group)"
            >
              <span class="min-w-0">
                <span class="block truncate text-sm font-semibold text-slate-950 dark:text-white">{{ groupName(group) }}</span>
                <span class="mt-1 flex min-w-0 items-center gap-2 text-xs text-slate-500 dark:text-slate-400">
                  <code class="truncate">{{ group.slug }}</code>
                  <span>·</span>
                  <span class="shrink-0">{{ $t('admin.catalog.tags.tagCount', { count: numberFormatter.format(group.tags.length) }) }}</span>
                </span>
              </span>
              <UIcon name="i-lucide-chevron-right" class="size-4 shrink-0 text-slate-400" />
            </button>
          </div>
        </section>

        <section class="overflow-hidden rounded-lg border border-slate-200 bg-white dark:border-slate-800 dark:bg-slate-900">
          <div class="flex h-12 items-center justify-between border-b border-slate-200 px-4 dark:border-slate-800">
            <h2 class="text-sm font-semibold text-slate-950 dark:text-white">
              {{ selectedGroupId ? $t('admin.catalog.tags.form.editGroup') : $t('admin.catalog.tags.form.createGroup') }}
            </h2>
            <UButton v-if="selectedGroupId" color="primary" variant="soft" size="sm" icon="i-lucide-plus" @click="startCreateGroup">
              {{ $t('admin.actions.new') }}
            </UButton>
          </div>

          <form class="space-y-5 p-4" @submit.prevent="saveGroup">
            <div class="grid gap-3 sm:grid-cols-2 2xl:grid-cols-3">
              <UFormField v-for="item in localeOptions" :key="item.code" :label="`${item.name} (${item.code})`">
                <UInput v-model="groupForm.names[item.code]" class="w-full" :ui="{ base: 'h-10 w-full' }" :disabled="savingGroup" />
              </UFormField>
            </div>

            <div class="grid gap-3 sm:grid-cols-[minmax(0,1fr)_120px]">
              <UFormField :label="$t('admin.catalog.tags.form.slug')">
                <UInput v-model="groupForm.slug" class="w-full" :ui="{ base: 'h-10 w-full font-mono' }" :disabled="savingGroup || Boolean(selectedGroupId)" />
              </UFormField>
              <UFormField :label="$t('admin.catalog.tags.form.sortOrder')">
                <UInput v-model.number="groupForm.sortOrder" type="number" class="w-full" :ui="{ base: 'h-10 w-full' }" :disabled="savingGroup" />
              </UFormField>
            </div>

            <div>
              <p class="text-sm font-medium text-slate-700 dark:text-slate-200">{{ $t('admin.catalog.tags.form.categories') }}</p>
              <div class="mt-2 flex flex-wrap gap-1.5">
                <UButton
                  v-for="category in categories"
                  :key="category.id"
                  type="button"
                  color="neutral"
                  size="xs"
                  :variant="groupForm.categoryIds.includes(category.id) ? 'soft' : 'outline'"
                  :class="groupForm.categoryIds.includes(category.id) ? 'ring-1 ring-sky-300 dark:ring-sky-700' : ''"
                  :disabled="savingGroup"
                  @click="toggleCategory(category.id)"
                >
                  {{ categoryName(category) }}
                </UButton>
                <span v-if="categories.length === 0" class="text-sm text-slate-500 dark:text-slate-400">{{ $t('admin.catalog.tags.form.allCategories') }}</span>
              </div>
              <p v-if="categories.length && groupForm.categoryIds.length === 0" class="mt-2 text-xs text-slate-500 dark:text-slate-400">{{ $t('admin.catalog.tags.form.allCategories') }}</p>
            </div>

            <p v-if="formError" class="rounded-md border border-red-200 bg-red-50 px-3 py-2 text-sm text-red-700 dark:border-red-900 dark:bg-red-950/40 dark:text-red-200">{{ formError }}</p>

            <div class="flex items-center justify-between border-t border-slate-200 pt-4 dark:border-slate-800">
              <UPopover v-if="selectedGroup" :content="{ side: 'top', align: 'start', sideOffset: 8 }" :ui="{ content: 'w-72 p-3' }">
                <UButton type="button" color="error" variant="ghost" size="sm" icon="i-lucide-trash-2" :disabled="busy" />
                <template #content="{ close }">
                  <DeleteConfirmContent :pending="deleting" @cancel="close()" @confirm="deleteGroup(close)" />
                </template>
              </UPopover>
              <span v-else />
              <div class="flex gap-2">
                <UButton type="button" color="neutral" variant="ghost" :disabled="busy || !groupDirty" @click="resetGroupForm">{{ $t('admin.actions.reset') }}</UButton>
                <UButton type="submit" color="primary" icon="i-lucide-save" :loading="savingGroup" :disabled="busy || !canSaveGroup">{{ $t('common.save') }}</UButton>
              </div>
            </div>
          </form>

          <div v-if="selectedGroup" class="border-t border-slate-200 dark:border-slate-800">
            <div class="flex h-12 items-center justify-between px-4">
              <h3 class="text-sm font-semibold text-slate-950 dark:text-white">{{ $t('admin.catalog.tags.tags') }}</h3>
              <UButton color="neutral" variant="soft" size="sm" icon="i-lucide-plus" @click="startCreateTag">{{ $t('admin.actions.new') }}</UButton>
            </div>
            <div class="border-t border-slate-200 dark:border-slate-800">
              <button
                v-for="tag in selectedGroup.tags"
                :key="tag.id"
                type="button"
                class="grid w-full grid-cols-[minmax(0,1fr)_minmax(100px,180px)_64px] items-center gap-3 border-b border-slate-100 px-4 py-2.5 text-left last:border-b-0 hover:bg-slate-50 dark:border-slate-800 dark:hover:bg-slate-950/70"
                :class="selectedTagId === tag.id ? 'bg-sky-50/60 dark:bg-sky-950/20' : ''"
                @click="selectTag(tag)"
              >
                <span class="truncate text-sm font-medium text-slate-950 dark:text-white">{{ tagName(tag) }}</span>
                <code class="truncate text-xs text-slate-500 dark:text-slate-400">{{ tag.value }}</code>
                <span class="text-right text-xs tabular-nums text-slate-500 dark:text-slate-400">{{ tag.sortOrder }}</span>
              </button>
              <p v-if="selectedGroup.tags.length === 0" class="px-4 py-8 text-center text-sm text-slate-500 dark:text-slate-400">{{ $t('admin.catalog.tags.noTags') }}</p>
            </div>

            <form v-if="tagEditorOpen" class="space-y-4 border-t border-slate-200 bg-slate-50/60 p-4 dark:border-slate-800 dark:bg-slate-950/30" @submit.prevent="saveTag">
              <div class="grid gap-3 sm:grid-cols-2 2xl:grid-cols-3">
                <UFormField v-for="item in localeOptions" :key="item.code" :label="`${item.name} (${item.code})`">
                  <UInput v-model="tagForm.names[item.code]" class="w-full" :ui="{ base: 'h-10 w-full' }" :disabled="savingTag" />
                </UFormField>
              </div>
              <div class="grid gap-3 sm:grid-cols-[minmax(0,1fr)_120px]">
                <UFormField :label="$t('admin.catalog.tags.form.value')">
                  <UInput v-model="tagForm.value" class="w-full" :ui="{ base: 'h-10 w-full font-mono' }" :disabled="savingTag || Boolean(selectedTagId)" />
                </UFormField>
                <UFormField :label="$t('admin.catalog.tags.form.sortOrder')">
                  <UInput v-model.number="tagForm.sortOrder" type="number" class="w-full" :ui="{ base: 'h-10 w-full' }" :disabled="savingTag" />
                </UFormField>
              </div>
              <div class="flex items-center justify-between">
                <UPopover v-if="selectedTag" :content="{ side: 'top', align: 'start', sideOffset: 8 }" :ui="{ content: 'w-72 p-3' }">
                  <UButton type="button" color="error" variant="ghost" size="sm" icon="i-lucide-trash-2" :disabled="busy" />
                  <template #content="{ close }">
                    <DeleteConfirmContent :pending="deleting" @cancel="close()" @confirm="deleteTag(close)" />
                  </template>
                </UPopover>
                <span v-else />
                <div class="flex gap-2">
                  <UButton type="button" color="neutral" variant="ghost" @click="tagEditorOpen = false">{{ $t('common.cancel') }}</UButton>
                  <UButton type="submit" color="primary" icon="i-lucide-save" :loading="savingTag" :disabled="busy || !canSaveTag">{{ $t('common.save') }}</UButton>
                </div>
              </div>
            </form>
          </div>
        </section>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ApiError } from '~/composables/useApi'
import type { AdminCatalogTag, AdminCatalogTagGroup, AdminCatalogTagGroupInput, AdminCatalogTagInput } from '~/composables/useAdmin'
import type { CatalogCategory } from '~/composables/useCatalogTorrents'

interface LocaleOption { code: string, name: string }

definePageMeta({ layout: 'admin', middleware: 'admin' })

const { t, locale, locales } = useI18n()
const adminApi = useAdmin()
const catalogApi = useCatalogTorrents()
const toast = useToast()

const groups = ref<AdminCatalogTagGroup[]>([])
const categories = ref<CatalogCategory[]>([])
const selectedGroupId = ref<number | null>(null)
const selectedTagId = ref<number | null>(null)
const pending = ref(true)
const savingGroup = ref(false)
const savingTag = ref(false)
const deleting = ref(false)
const errorMessage = ref('')
const formError = ref('')
const groupSnapshot = ref('')
const tagEditorOpen = ref(false)

const groupForm = reactive({ names: {} as Record<string, string>, slug: '', categoryIds: [] as number[], sortOrder: 0 })
const tagForm = reactive({ names: {} as Record<string, string>, value: '', sortOrder: 0 })

const localeOptions = computed<LocaleOption[]>(() => (locales.value as Array<string | { code?: string, language?: string, name?: string }>).map((item) => {
  if (typeof item === 'string') return { code: item, name: item }
  const code = String(item.code || item.language || '')
  return { code, name: String(item.name || code) }
}).filter(item => item.code))
const selectedGroup = computed(() => groups.value.find(group => group.id === selectedGroupId.value) || null)
const selectedTag = computed(() => selectedGroup.value?.tags.find(tag => tag.id === selectedTagId.value) || null)
const numberFormatter = computed(() => new Intl.NumberFormat(locale.value))
const busy = computed(() => savingGroup.value || savingTag.value || deleting.value)
const groupDirty = computed(() => groupFormSnapshot() !== groupSnapshot.value)
const canSaveGroup = computed(() => Boolean(primaryName(groupForm.names) && groupForm.slug.trim() && groupDirty.value))
const canSaveTag = computed(() => Boolean(primaryName(tagForm.names) && tagForm.value.trim()))

watch(localeOptions, ensureLocaleFields, { immediate: true })
onMounted(loadPage)

async function loadPage() {
  pending.value = true
  errorMessage.value = ''
  try {
    const [groupOut, categoryOut] = await Promise.all([adminApi.listCatalogTagGroups(), catalogApi.listCategories()])
    groups.value = groupOut.groups || []
    categories.value = categoryOut.list || []
    if (selectedGroupId.value) {
      const current = groups.value.find(group => group.id === selectedGroupId.value)
      if (current) selectGroup(current)
    }
  } catch (error) {
    errorMessage.value = error instanceof ApiError ? error.message : t('common.requestFailed')
  } finally {
    pending.value = false
  }
}

function selectGroup(group: AdminCatalogTagGroup) {
  selectedGroupId.value = group.id
  setNames(groupForm.names, group.nameI18N)
  groupForm.slug = group.slug
  groupForm.categoryIds = [...(group.categoryIds || [])]
  groupForm.sortOrder = group.sortOrder
  selectedTagId.value = null
  tagEditorOpen.value = false
  formError.value = ''
  groupSnapshot.value = groupFormSnapshot()
}

function startCreateGroup() {
  selectedGroupId.value = null
  setNames(groupForm.names)
  groupForm.slug = ''
  groupForm.categoryIds = []
  groupForm.sortOrder = 0
  selectedTagId.value = null
  tagEditorOpen.value = false
  formError.value = ''
  groupSnapshot.value = groupFormSnapshot()
}

function resetGroupForm() {
  if (selectedGroup.value) selectGroup(selectedGroup.value)
  else startCreateGroup()
}

async function saveGroup() {
  if (!canSaveGroup.value) return
  savingGroup.value = true
  formError.value = ''
  const input: AdminCatalogTagGroupInput = {
    nameI18N: buildNames(groupForm.names),
    slug: groupForm.slug.trim(),
    categoryIds: groupForm.categoryIds,
    sortOrder: Number(groupForm.sortOrder || 0)
  }
  try {
    const currentId = selectedGroupId.value
    if (currentId) await adminApi.updateCatalogTagGroup(currentId, input)
    else await adminApi.createCatalogTagGroup(input)
    toast.add({ title: t('admin.catalog.tags.saved') })
    await loadPage()
    const target = currentId ? groups.value.find(group => group.id === currentId) : groups.value.find(group => group.slug === input.slug)
    if (target) selectGroup(target)
  } catch (error) {
    formError.value = error instanceof ApiError ? error.message : t('common.requestFailed')
  } finally {
    savingGroup.value = false
  }
}

async function deleteGroup(close: () => void) {
  if (!selectedGroup.value) return
  deleting.value = true
  try {
    await adminApi.deleteCatalogTagGroup(selectedGroup.value.id)
    toast.add({ title: t('admin.catalog.tags.deleted') })
    close()
    startCreateGroup()
    await loadPage()
  } catch (error) {
    toast.add({ color: 'error', title: error instanceof ApiError ? error.message : t('common.requestFailed') })
  } finally {
    deleting.value = false
  }
}

function startCreateTag() {
  selectedTagId.value = null
  setNames(tagForm.names)
  tagForm.value = ''
  tagForm.sortOrder = 0
  tagEditorOpen.value = true
}

function selectTag(tag: AdminCatalogTag) {
  selectedTagId.value = tag.id
  setNames(tagForm.names, tag.nameI18N)
  tagForm.value = tag.value
  tagForm.sortOrder = tag.sortOrder
  tagEditorOpen.value = true
}

async function saveTag() {
  if (!selectedGroup.value || !canSaveTag.value) return
  savingTag.value = true
  const input: AdminCatalogTagInput = { nameI18N: buildNames(tagForm.names), value: tagForm.value.trim(), sortOrder: Number(tagForm.sortOrder || 0) }
  try {
    if (selectedTagId.value) await adminApi.updateCatalogTag(selectedTagId.value, input)
    else await adminApi.createCatalogTag(selectedGroup.value.id, input)
    toast.add({ title: t('admin.catalog.tags.tagSaved') })
    tagEditorOpen.value = false
    await loadPage()
  } catch (error) {
    toast.add({ color: 'error', title: error instanceof ApiError ? error.message : t('common.requestFailed') })
  } finally {
    savingTag.value = false
  }
}

async function deleteTag(close: () => void) {
  if (!selectedTag.value) return
  deleting.value = true
  try {
    await adminApi.deleteCatalogTag(selectedTag.value.id)
    toast.add({ title: t('admin.catalog.tags.tagDeleted') })
    close()
    tagEditorOpen.value = false
    selectedTagId.value = null
    await loadPage()
  } catch (error) {
    toast.add({ color: 'error', title: error instanceof ApiError ? error.message : t('common.requestFailed') })
  } finally {
    deleting.value = false
  }
}

function toggleCategory(id: number) {
  groupForm.categoryIds = groupForm.categoryIds.includes(id) ? groupForm.categoryIds.filter(item => item !== id) : [...groupForm.categoryIds, id]
}

function ensureLocaleFields() {
  setNames(groupForm.names, groupForm.names)
  setNames(tagForm.names, tagForm.names)
}

function setNames(target: Record<string, string>, source?: Record<string, unknown> | null) {
  const values = { ...(source || {}) }
  for (const key of Object.keys(target)) delete target[key]
  for (const item of localeOptions.value) target[item.code] = typeof values[item.code] === 'string' ? String(values[item.code]) : ''
}

function buildNames(source: Record<string, string>) {
  const fallback = primaryName(source)
  return Object.fromEntries(localeOptions.value.map(item => [item.code, source[item.code]?.trim() || fallback]))
}

function primaryName(source: Record<string, string>) {
  return source[localeOptions.value[0]?.code || locale.value]?.trim() || ''
}

function groupFormSnapshot() {
  return JSON.stringify({ names: groupForm.names, slug: groupForm.slug, categoryIds: [...groupForm.categoryIds].sort((a, b) => a - b), sortOrder: Number(groupForm.sortOrder || 0) })
}

function groupName(group: AdminCatalogTagGroup) { return localizeI18nName(group.nameI18N, locale.value, group.slug) }
function tagName(tag: AdminCatalogTag) { return localizeI18nName(tag.nameI18N, locale.value, tag.value) }
function categoryName(category: CatalogCategory) { return localizeI18nName(category.name, locale.value, category.slug) }

useHead(() => ({ title: t('admin.catalog.tags.title') }))
</script>
