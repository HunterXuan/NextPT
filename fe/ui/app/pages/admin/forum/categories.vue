<template>
  <div class="min-h-[calc(100vh-4rem)] bg-slate-50 py-6 dark:bg-slate-950">
    <div class="w-full px-3 sm:px-4 lg:px-5">
      <div class="grid gap-4 xl:grid-cols-[minmax(620px,780px)_minmax(520px,1fr)] 2xl:grid-cols-[minmax(660px,820px)_minmax(560px,1fr)]">
        <section class="overflow-hidden rounded-lg border border-slate-200 bg-white dark:border-slate-800 dark:bg-slate-900">
          <div class="border-b border-slate-200 px-4 py-3 dark:border-slate-800">
            <div class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
              <div class="flex items-center gap-2">
                <UIcon name="i-lucide-folder-tree" class="size-5 text-cyan-600 dark:text-cyan-300" />
                <h2 class="text-sm font-semibold text-slate-950 dark:text-white">{{ $t('admin.forum.categories.list') }}</h2>
              </div>
              <div class="flex flex-wrap items-center gap-2 text-xs">
                <span class="inline-flex h-7 items-center gap-1.5 rounded-md bg-slate-100 px-2.5 font-medium text-slate-600 dark:bg-slate-800 dark:text-slate-300">
                  <span>{{ $t('admin.forum.categories.stats.total') }}</span>
                  <span class="font-semibold text-slate-950 dark:text-white">{{ numberFormatter.format(categories.length) }}</span>
                </span>
              </div>
            </div>
          </div>

          <div v-if="pending" class="overflow-x-auto">
            <table class="min-w-[760px] w-full table-fixed border-collapse text-left">
              <thead class="bg-slate-50 text-xs font-medium uppercase text-slate-500 dark:bg-slate-950/70 dark:text-slate-400">
                <tr>
                  <th class="w-[30%] border-b border-slate-200 px-3 py-2.5 dark:border-slate-800">{{ $t('admin.forum.categories.table.name') }}</th>
                  <th class="w-[32%] border-b border-slate-200 px-3 py-2.5 dark:border-slate-800">{{ $t('admin.forum.categories.table.description') }}</th>
                  <th class="w-[12%] border-b border-slate-200 px-3 py-2.5 text-right dark:border-slate-800">{{ $t('admin.forum.categories.table.sort') }}</th>
                  <th class="w-[14%] border-b border-slate-200 px-3 py-2.5 text-right dark:border-slate-800">{{ $t('admin.forum.categories.table.minRoleView') }}</th>
                  <th class="w-[12%] border-b border-slate-200 px-3 py-2.5 text-right dark:border-slate-800">{{ $t('admin.forum.categories.table.actions') }}</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="index in 4" :key="index" class="border-b border-slate-200 last:border-b-0 dark:border-slate-800">
                  <td class="px-3 py-3">
                    <div class="h-4 w-44 animate-pulse rounded bg-slate-200 dark:bg-slate-800" />
                    <div class="mt-2 h-3 w-28 animate-pulse rounded bg-slate-100 dark:bg-slate-800/70" />
                  </td>
                  <td class="px-3 py-3"><div class="h-4 w-40 animate-pulse rounded bg-slate-100 dark:bg-slate-800/70" /></td>
                  <td class="px-3 py-3"><div class="ml-auto h-4 w-10 animate-pulse rounded bg-slate-100 dark:bg-slate-800/70" /></td>
                  <td class="px-3 py-3"><div class="ml-auto h-4 w-10 animate-pulse rounded bg-slate-100 dark:bg-slate-800/70" /></td>
                  <td class="px-3 py-3"><div class="ml-auto h-4 w-20 animate-pulse rounded bg-slate-100 dark:bg-slate-800/70" /></td>
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
                  <th class="w-[30%] border-b border-slate-200 px-3 py-2.5 dark:border-slate-800">{{ $t('admin.forum.categories.table.name') }}</th>
                  <th class="w-[32%] border-b border-slate-200 px-3 py-2.5 dark:border-slate-800">{{ $t('admin.forum.categories.table.description') }}</th>
                  <th class="w-[12%] border-b border-slate-200 px-3 py-2.5 text-right dark:border-slate-800">{{ $t('admin.forum.categories.table.sort') }}</th>
                  <th class="w-[14%] border-b border-slate-200 px-3 py-2.5 text-right dark:border-slate-800">{{ $t('admin.forum.categories.table.minRoleView') }}</th>
                  <th class="w-[12%] border-b border-slate-200 px-3 py-2.5 text-right dark:border-slate-800">{{ $t('admin.forum.categories.table.actions') }}</th>
                </tr>
              </thead>
              <tbody>
                <tr
                  v-for="category in categories"
                  :key="category.id"
                  class="border-b border-slate-200 transition-colors last:border-b-0 dark:border-slate-800"
                  :class="selectedId === category.id ? 'bg-cyan-50/70 dark:bg-cyan-950/30' : 'hover:bg-slate-50 dark:hover:bg-slate-950/70'"
                >
                  <td class="px-3 py-2.5 align-middle">
                    <button type="button" class="block max-w-full text-left" @click="selectCategory(category)">
                      <span class="block truncate text-sm font-semibold text-slate-950 hover:text-cyan-700 dark:text-white dark:hover:text-cyan-300">
                        {{ categoryName(category) }}
                      </span>
                      <span class="mt-1 block truncate text-xs text-slate-500 dark:text-slate-400">
                        {{ formatDateTime(category.updatedAt || category.createdAt, locale) }}
                      </span>
                    </button>
                  </td>
                  <td class="px-3 py-2.5 align-middle text-sm text-slate-600 dark:text-slate-300">
                    <span class="block truncate">{{ categoryDescription(category) || '-' }}</span>
                  </td>
                  <td class="px-3 py-2.5 text-right align-middle text-sm text-slate-600 dark:text-slate-300">
                    {{ numberFormatter.format(category.sortOrder) }}
                  </td>
                  <td class="px-3 py-2.5 text-right align-middle text-sm text-slate-600 dark:text-slate-300">
                    <span class="inline-flex max-w-full items-center justify-end">
                      <span class="truncate rounded-md bg-slate-100 px-2 py-1 text-xs font-medium text-slate-700 dark:bg-slate-800 dark:text-slate-200">
                        {{ roleNameByLevel(category.minRoleView) }}
                      </span>
                    </span>
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
                          @click="selectCategory(category)"
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
                          :loading="deletingId === category.id"
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
                              <UButton color="error" size="xs" type="button" icon="i-lucide-trash-2" :loading="deletingId === category.id" :disabled="deletingId !== null" @click="deleteCategory(category, close)">
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
                {{ selectedId ? $t('admin.forum.categories.form.edit') : $t('admin.forum.categories.form.create') }}
              </h2>
              <div class="flex h-8 min-w-20 items-center justify-end">
                <UButton v-if="selectedId" color="primary" variant="soft" size="sm" icon="i-lucide-plus" @click="startCreate">
                  {{ $t('admin.actions.new') }}
                </UButton>
              </div>
            </div>
          </div>

          <form class="space-y-5 p-4" @submit.prevent="saveCategory">
            <div>
              <div class="flex items-center justify-between gap-3 border-b border-slate-200 pb-2 dark:border-slate-800">
                <span class="text-sm font-semibold text-slate-950 dark:text-white">{{ $t('admin.forum.categories.form.name') }}</span>
                <span class="text-xs text-slate-500 dark:text-slate-400">
                  {{ $t('admin.forum.categories.form.primaryLocale', { locale: primaryLocaleLabel }) }}
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
                <span class="text-sm font-semibold text-slate-950 dark:text-white">{{ $t('admin.forum.categories.form.description') }}</span>
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

            <div class="grid gap-3 sm:grid-cols-2">
              <label class="block">
                <span class="text-sm font-medium text-slate-700 dark:text-slate-200">{{ $t('admin.forum.categories.form.sortOrder') }}</span>
                <input
                  v-model.number="form.sortOrder"
                  type="number"
                  class="mt-1 h-10 w-full rounded-md border border-slate-200 bg-white px-3 text-sm text-slate-950 outline-none transition focus:border-cyan-400 focus:ring-2 focus:ring-cyan-100 dark:border-slate-700 dark:bg-slate-950 dark:text-white dark:focus:border-cyan-500 dark:focus:ring-cyan-950"
                  :disabled="saving"
                >
              </label>
              <label class="block">
                <span class="text-sm font-medium text-slate-700 dark:text-slate-200">{{ $t('admin.forum.categories.form.minRoleView') }}</span>
                <AdminIamRoleLevelSelect
                  v-model.number="form.minRoleView"
                  :roles="roles"
                  :disabled="saving"
                />
              </label>
            </div>

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
import type { AdminForumCategory, AdminForumCategoryInput, AdminIamRole } from '~/composables/useAdmin'
import { formatDateTime, localizeI18nName } from '~/utils/format'

interface LocaleOption {
  code: string
  name: string
}

definePageMeta({ layout: 'admin', middleware: 'admin' })

const { t, locale, locales } = useI18n()
const toast = useToast()
const adminApi = useAdmin()

useHead({ title: t('admin.forum.categories.title') })

const categories = ref<AdminForumCategory[]>([])
const roles = ref<AdminIamRole[]>([])
const pending = ref(false)
const saving = ref(false)
const deletingId = ref<number | null>(null)
const selectedId = ref<number | null>(null)
const errorMessage = ref('')
const formError = ref('')
const originalFormSnapshot = ref('')

const form = reactive({
  names: {} as Record<string, string>,
  descs: {} as Record<string, string>,
  sortOrder: 0,
  minRoleView: 0
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
const { roleNameByLevel } = useAdminIamRoleLevels(roles)
const selectedCategory = computed(() => categories.value.find((item) => item.id === selectedId.value) || null)
const isFormDirty = computed(() => formSnapshot() !== originalFormSnapshot.value)
const canSubmit = computed(() => Boolean(nameValue(primaryLocaleCode.value).trim() && isFormDirty.value && !saving.value))

watch(localeOptions, () => {
  ensureFormLocales()
}, { immediate: true })

onMounted(() => {
  originalFormSnapshot.value = formSnapshot()
  loadCategories()
})

async function loadCategories() {
  pending.value = true
  errorMessage.value = ''
  try {
    const [categoryData, roleData] = await Promise.all([
      adminApi.listForumCategories(),
      adminApi.listIamRoles()
    ])
    categories.value = (categoryData.categories || []).sort((a, b) => a.sortOrder - b.sortOrder || a.id - b.id)
    roles.value = roleData.roles || []
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
  setI18nValues(form.names, category.nameI18N)
  setI18nValues(form.descs, category.descI18N)
  form.sortOrder = category.sortOrder
  form.minRoleView = category.minRoleView
  formError.value = ''
  originalFormSnapshot.value = formSnapshot()
}

function startCreate() {
  selectedId.value = null
  setI18nValues(form.names)
  setI18nValues(form.descs)
  form.sortOrder = 0
  form.minRoleView = 0
  formError.value = ''
  originalFormSnapshot.value = formSnapshot()
}

function resetFormChanges() {
  const category = selectedCategory.value
  if (category) {
    selectCategory(category)
    return
  }
  startCreate()
}

function formSnapshot() {
  return JSON.stringify({
    names: localeOptions.value.reduce<Record<string, string>>((names, item) => {
      names[item.code] = nameValue(item.code)
      return names
    }, {}),
    descs: localeOptions.value.reduce<Record<string, string>>((descs, item) => {
      descs[item.code] = descValue(item.code)
      return descs
    }, {}),
    sortOrder: Number(form.sortOrder || 0),
    minRoleView: Number(form.minRoleView || 0)
  })
}

function buildInput(): AdminForumCategoryInput {
  const primaryName = nameValue(primaryLocaleCode.value).trim()
  const fallbackDesc = descValue(primaryLocaleCode.value).trim() || firstFilledDesc()

  return {
    nameI18N: localeOptions.value.reduce<Record<string, string>>((names, item) => {
      names[item.code] = nameValue(item.code).trim() || primaryName
      return names
    }, {}),
    descI18N: localeOptions.value.reduce<Record<string, string>>((descs, item) => {
      descs[item.code] = descValue(item.code).trim() || fallbackDesc
      return descs
    }, {}),
    sortOrder: Number(form.sortOrder || 0),
    minRoleView: Number(form.minRoleView || 0)
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

async function saveCategory() {
  if (!canSubmit.value) return
  saving.value = true
  formError.value = ''
  try {
    const input = buildInput()
    const editingId = selectedId.value
    if (editingId) {
      await adminApi.updateForumCategory(editingId, input)
    } else {
      await adminApi.createForumCategory(input)
    }
    toast.add({ title: t('admin.forum.categories.saved') })
    await loadCategories()
    if (editingId) {
      const updatedCategory = categories.value.find((item) => item.id === editingId)
      if (updatedCategory) {
        selectCategory(updatedCategory)
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

async function deleteCategory(category: AdminForumCategory, close?: () => void) {
  deletingId.value = category.id
  try {
    await adminApi.deleteForumCategory(category.id)
    if (selectedId.value === category.id) {
      startCreate()
    }
    toast.add({ title: t('admin.forum.categories.deleted') })
    close?.()
    await loadCategories()
  } catch (error: unknown) {
    toast.add({ color: 'error', title: error instanceof ApiError ? error.message : t('common.requestFailed') })
  } finally {
    deletingId.value = null
  }
}
</script>
