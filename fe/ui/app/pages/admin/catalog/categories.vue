<template>
  <div class="min-h-[calc(100vh-4rem)] bg-slate-50 py-6 dark:bg-slate-950">
    <div class="w-full px-3 sm:px-4 lg:px-5">
      <div class="grid gap-4 xl:grid-cols-[minmax(660px,800px)_minmax(560px,1fr)] 2xl:grid-cols-[minmax(700px,840px)_minmax(600px,1fr)]">
        <section class="overflow-hidden rounded-lg border border-slate-200 bg-white dark:border-slate-800 dark:bg-slate-900">
          <div class="border-b border-slate-200 px-4 py-3 dark:border-slate-800">
            <div class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
              <div class="flex items-center gap-2">
                <UIcon name="i-lucide-tags" class="size-5 text-sky-600 dark:text-sky-300" />
                <h2 class="text-sm font-semibold text-slate-950 dark:text-white">{{ $t('admin.catalog.categories.list') }}</h2>
              </div>
              <div class="flex flex-wrap items-center gap-2 text-xs">
                <span class="inline-flex h-7 items-center gap-1.5 rounded-md bg-slate-100 px-2.5 font-medium text-slate-600 dark:bg-slate-800 dark:text-slate-300">
                  <span>{{ $t('admin.catalog.categories.stats.total') }}</span>
                  <span class="font-semibold text-slate-950 dark:text-white">{{ numberFormatter.format(categories.length) }}</span>
                </span>
                <span class="inline-flex h-7 items-center gap-1.5 rounded-md bg-emerald-50 px-2.5 font-medium text-emerald-700 dark:bg-emerald-950/50 dark:text-emerald-300">
                  <span>{{ $t('admin.catalog.categories.stats.enabled') }}</span>
                  <span class="font-semibold">{{ numberFormatter.format(enabledCount) }}</span>
                </span>
                <span class="inline-flex h-7 items-center gap-1.5 rounded-md bg-amber-50 px-2.5 font-medium text-amber-700 dark:bg-amber-950/40 dark:text-amber-300">
                  <span>{{ $t('admin.catalog.categories.stats.disabled') }}</span>
                  <span class="font-semibold">{{ numberFormatter.format(disabledCount) }}</span>
                </span>
              </div>
            </div>
          </div>

          <div v-if="pending" class="overflow-x-auto">
            <table class="min-w-[760px] w-full table-fixed border-collapse text-left">
              <thead class="bg-slate-50 text-xs font-medium uppercase text-slate-500 dark:bg-slate-950/70 dark:text-slate-400">
                <tr>
                  <th class="w-[27%] border-b border-slate-200 px-3 py-2.5 dark:border-slate-800">{{ $t('admin.catalog.categories.table.name') }}</th>
                  <th class="w-[18%] border-b border-slate-200 px-3 py-2.5 dark:border-slate-800">{{ $t('admin.catalog.categories.table.slug') }}</th>
                  <th class="w-[18%] border-b border-slate-200 px-3 py-2.5 dark:border-slate-800">{{ $t('admin.catalog.categories.table.rule') }}</th>
                  <th class="w-[10%] border-b border-slate-200 px-3 py-2.5 text-right dark:border-slate-800">{{ $t('admin.catalog.categories.table.sort') }}</th>
                  <th class="w-[12%] border-b border-slate-200 px-3 py-2.5 dark:border-slate-800">{{ $t('admin.catalog.categories.table.status') }}</th>
                  <th class="w-[15%] border-b border-slate-200 px-3 py-2.5 text-right dark:border-slate-800">{{ $t('admin.catalog.categories.table.actions') }}</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="index in 5" :key="index" class="border-b border-slate-200 last:border-b-0 dark:border-slate-800">
                  <td class="px-3 py-3">
                    <div class="h-4 w-48 animate-pulse rounded bg-slate-200 dark:bg-slate-800" />
                    <div class="mt-2 h-3 w-28 animate-pulse rounded bg-slate-100 dark:bg-slate-800/70" />
                  </td>
                  <td class="px-3 py-3">
                    <div class="h-4 w-24 animate-pulse rounded bg-slate-100 dark:bg-slate-800/70" />
                  </td>
                  <td class="px-3 py-3">
                    <div class="h-5 w-24 animate-pulse rounded bg-slate-100 dark:bg-slate-800/70" />
                  </td>
                  <td class="px-3 py-3">
                    <div class="ml-auto h-4 w-12 animate-pulse rounded bg-slate-100 dark:bg-slate-800/70" />
                  </td>
                  <td class="px-3 py-3">
                    <div class="h-4 w-16 animate-pulse rounded bg-slate-100 dark:bg-slate-800/70" />
                  </td>
                  <td class="px-3 py-3">
                    <div class="ml-auto h-4 w-28 animate-pulse rounded bg-slate-100 dark:bg-slate-800/70" />
                  </td>
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
            <p class="mt-3 text-sm font-medium text-slate-950 dark:text-white">{{ $t('admin.catalog.categories.empty') }}</p>
          </div>

          <div v-else class="overflow-x-auto">
            <table class="min-w-[760px] w-full table-fixed border-collapse text-left">
              <thead class="bg-slate-50 text-xs font-medium uppercase text-slate-500 dark:bg-slate-950/70 dark:text-slate-400">
                <tr>
                  <th class="w-[27%] border-b border-slate-200 px-3 py-2.5 dark:border-slate-800">{{ $t('admin.catalog.categories.table.name') }}</th>
                  <th class="w-[18%] border-b border-slate-200 px-3 py-2.5 dark:border-slate-800">{{ $t('admin.catalog.categories.table.slug') }}</th>
                  <th class="w-[18%] border-b border-slate-200 px-3 py-2.5 dark:border-slate-800">{{ $t('admin.catalog.categories.table.rule') }}</th>
                  <th class="w-[10%] border-b border-slate-200 px-3 py-2.5 text-right dark:border-slate-800">{{ $t('admin.catalog.categories.table.sort') }}</th>
                  <th class="w-[12%] border-b border-slate-200 px-3 py-2.5 dark:border-slate-800">{{ $t('admin.catalog.categories.table.status') }}</th>
                  <th class="w-[15%] border-b border-slate-200 px-3 py-2.5 text-right dark:border-slate-800">{{ $t('admin.catalog.categories.table.actions') }}</th>
                </tr>
              </thead>
              <tbody>
                <tr
                  v-for="category in categories"
                  :key="category.id"
                  class="border-b border-slate-200 transition-colors last:border-b-0 dark:border-slate-800"
                  :class="selectedId === category.id ? 'bg-sky-50/70 dark:bg-sky-950/30' : 'hover:bg-slate-50 dark:hover:bg-slate-950/70'"
                >
                  <td class="px-3 py-2.5 align-middle">
                    <button class="block max-w-full text-left" type="button" @click="selectCategory(category)">
                      <span class="block truncate text-sm font-semibold text-slate-950 hover:text-sky-700 dark:text-white dark:hover:text-sky-300">
                        {{ categoryDisplayName(category) }}
                      </span>
                      <span class="mt-1 block truncate text-xs text-slate-500 dark:text-slate-400">
                        {{ formatDateTime(category.updatedAt || category.createdAt, locale) }}
                      </span>
                    </button>
                  </td>
                  <td class="px-3 py-2.5 align-middle">
                    <code class="block truncate text-xs text-slate-600 dark:text-slate-300">
                      {{ category.slug }}
                    </code>
                  </td>
                  <td class="px-3 py-2.5 align-middle">
                    <UBadge :color="category.uploadConfig ? 'primary' : 'neutral'" variant="soft">
                      {{ uploadRuleLabel(category) }}
                    </UBadge>
                  </td>
                  <td class="px-3 py-2.5 text-right align-middle text-sm text-slate-600 dark:text-slate-300">
                    {{ numberFormatter.format(category.sortOrder) }}
                  </td>
                  <td class="px-3 py-2.5 align-middle">
                    <UBadge :color="category.enabled ? 'success' : 'warning'" variant="soft">
                      {{ category.enabled ? $t('admin.status.enabled') : $t('admin.status.disabled') }}
                    </UBadge>
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
                              {{ $t('admin.catalog.categories.confirmDeleteTitle') }}
                            </p>
                            <p class="truncate text-xs text-slate-500 dark:text-slate-400" :title="categoryDisplayName(category)">
                              {{ categoryDisplayName(category) }}
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
                {{ selectedId ? $t('admin.catalog.categories.form.edit') : $t('admin.catalog.categories.form.create') }}
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
                <span class="text-sm font-semibold text-slate-950 dark:text-white">{{ $t('admin.catalog.categories.form.name') }}</span>
                <span class="text-xs text-slate-500 dark:text-slate-400">
                  {{ $t('admin.catalog.categories.form.primaryLocale', { locale: primaryLocaleLabel }) }}
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
                    class="mt-1 h-10 w-full rounded-md border border-slate-200 bg-white px-3 text-sm text-slate-950 outline-none transition focus:border-sky-400 focus:ring-2 focus:ring-sky-100 dark:border-slate-700 dark:bg-slate-950 dark:text-white dark:focus:border-sky-500 dark:focus:ring-sky-950"
                    :disabled="saving"
                  >
                </label>
              </div>
            </div>

            <div class="grid gap-3 sm:grid-cols-[minmax(0,1fr)_120px_112px] sm:items-end">
              <label class="block">
                <span class="text-sm font-medium text-slate-700 dark:text-slate-200">{{ $t('admin.catalog.categories.form.slug') }}</span>
                <input
                  v-model="form.slug"
                  class="mt-1 h-10 w-full rounded-md border border-slate-200 bg-white px-3 text-sm text-slate-950 outline-none transition focus:border-sky-400 focus:ring-2 focus:ring-sky-100 dark:border-slate-700 dark:bg-slate-950 dark:text-white dark:focus:border-sky-500 dark:focus:ring-sky-950"
                  :disabled="saving"
                >
              </label>
              <label class="block">
                <span class="text-sm font-medium text-slate-700 dark:text-slate-200">{{ $t('admin.catalog.categories.form.sortOrder') }}</span>
                <input
                  v-model.number="form.sortOrder"
                  type="number"
                  class="mt-1 h-10 w-full rounded-md border border-slate-200 bg-white px-3 text-sm text-slate-950 outline-none transition focus:border-sky-400 focus:ring-2 focus:ring-sky-100 dark:border-slate-700 dark:bg-slate-950 dark:text-white dark:focus:border-sky-500 dark:focus:ring-sky-950"
                  :disabled="saving"
                >
              </label>

              <label class="flex h-10 items-center justify-between gap-3 rounded-md border border-slate-200 bg-white px-3 text-sm font-medium text-slate-700 dark:border-slate-700 dark:bg-slate-950 dark:text-slate-200">
                <span>{{ $t('admin.catalog.categories.form.enabled') }}</span>
                <input
                  v-model="form.enabled"
                  type="checkbox"
                  class="peer sr-only"
                  :disabled="saving"
                >
                <span class="relative h-5 w-9 rounded-full bg-slate-200 transition-colors after:absolute after:left-0.5 after:top-0.5 after:size-4 after:rounded-full after:bg-white after:shadow-sm after:transition-transform peer-checked:bg-sky-500 peer-checked:after:translate-x-4 peer-disabled:opacity-60 dark:bg-slate-700" />
              </label>
            </div>

            <div class="border-t border-slate-200 pt-4 dark:border-slate-800">
              <div class="flex flex-col gap-3 sm:flex-row sm:items-start sm:justify-between">
                <div>
                  <h3 class="text-sm font-semibold text-slate-950 dark:text-white">{{ $t('admin.catalog.categories.form.uploadRule') }}</h3>
                  <p class="mt-1 text-xs text-slate-500 dark:text-slate-400">{{ $t('admin.catalog.categories.form.uploadRuleDescription') }}</p>
                </div>
                <div class="flex shrink-0 flex-wrap items-center gap-2">
                  <UButton type="button" color="neutral" variant="soft" size="xs" icon="i-lucide-eraser" :disabled="saving || !form.uploadConfigText.trim()" @click="clearUploadConfig">
                    {{ $t('admin.catalog.categories.form.uploadConfigClear') }}
                  </UButton>
                  <UButton type="button" color="neutral" variant="outline" size="xs" icon="i-lucide-align-left" :disabled="saving || !uploadConfigState.valid || !uploadConfigState.config" @click="formatUploadConfig">
                    {{ $t('admin.catalog.categories.form.uploadConfigFormat') }}
                  </UButton>
                </div>
              </div>

              <div class="mt-3 grid gap-2 sm:grid-cols-3">
                <div
                  v-for="item in uploadConfigSummaryItems"
                  :key="item.key"
                  class="rounded-md border border-slate-200 bg-slate-50 px-3 py-2 dark:border-slate-800 dark:bg-slate-950/50"
                >
                  <p class="text-xs text-slate-500 dark:text-slate-400">{{ item.label }}</p>
                  <p class="mt-1 truncate text-sm font-semibold text-slate-950 dark:text-white">{{ item.value }}</p>
                </div>
              </div>

              <div class="mt-4">
                <AdminJsonEditor
                  v-model="form.uploadConfigText"
                  :valid="uploadConfigState.valid"
                  :error-message="uploadConfigState.message"
                  :title="$t('admin.catalog.categories.form.uploadConfigJson')"
                  :placeholder="$t('admin.catalog.categories.form.uploadConfigPlaceholder')"
                  :hint="$t('admin.catalog.categories.form.uploadConfigEditorHint')"
                  :valid-text="$t('admin.catalog.categories.form.uploadConfigValid')"
                  :invalid-text="$t('admin.catalog.categories.form.uploadConfigInvalidShort')"
                  :fullscreen-text="$t('admin.catalog.categories.form.uploadConfigFullscreen')"
                  :close-text="$t('common.close')"
                  :disabled="saving"
                />
              </div>
              <p v-if="!uploadConfigState.valid" class="mt-2 text-xs text-red-600 dark:text-red-300">
                {{ uploadConfigState.message }}
              </p>
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
import type { AdminCatalogCategory, AdminCatalogCategoryInput } from '~/composables/useAdmin'
import type { UploadConfig } from '~/composables/useCatalogTorrents'
import { formatDateTime, localizeI18nName } from '~/utils/format'

interface LocaleOption {
  code: string
  name: string
}

interface UploadConfigParseState {
  valid: boolean
  config: UploadConfig | null
  message: string
}

type UploadConfigRecord = Record<string, unknown>

definePageMeta({
  layout: 'admin',
  middleware: 'admin'
})

const { t, locale, locales } = useI18n()
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
const originalFormSnapshot = ref('')

const form = reactive({
  names: {} as Record<string, string>,
  slug: '',
  sortOrder: 0,
  enabled: true,
  uploadConfigText: ''
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
const enabledCount = computed(() => categories.value.filter((item) => item.enabled).length)
const disabledCount = computed(() => categories.value.length - enabledCount.value)
const selectedCategory = computed(() => categories.value.find((item) => item.id === selectedId.value) || null)
const isFormDirty = computed(() => formSnapshot() !== originalFormSnapshot.value)
const uploadConfigState = computed(() => parseUploadConfigJson(form.uploadConfigText))
const uploadConfigSummaryItems = computed(() => {
  const config = uploadConfigState.value.config
  const fields = Array.isArray(config?.fields) ? config.fields : []
  const parts = Array.isArray(config?.title?.parts) ? config.title.parts : []

  return [
    {
      key: 'type',
      label: t('admin.catalog.categories.form.uploadRuleSummaryType'),
      value: uploadRuleTypeLabel(config, fields.length)
    },
    {
      key: 'title',
      label: t('admin.catalog.categories.form.uploadRuleSummaryTitle'),
      value: uploadRuleTitleLabel(config, parts.length)
    },
    {
      key: 'fields',
      label: t('admin.catalog.categories.form.uploadRuleSummaryFields'),
      value: fields.length > 0
        ? t('admin.catalog.categories.rules.fieldCount', { count: numberFormatter.value.format(fields.length) })
        : t('admin.catalog.categories.form.uploadRuleNoFields')
    }
  ]
})
const canSubmit = computed(() => Boolean(nameValue(primaryLocaleCode.value).trim() && form.slug.trim() && isFormDirty.value && !saving.value && uploadConfigState.value.valid))

watch(localeOptions, () => {
  ensureFormNameLocales()
}, { immediate: true })

onMounted(() => {
  originalFormSnapshot.value = formSnapshot()
  loadCategories()
})

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

function uploadRuleLabel(category: AdminCatalogCategory) {
  const count = uploadRuleFieldCount(category)
  if (count > 0) return t('admin.catalog.categories.rules.fieldCount', { count: numberFormatter.value.format(count) })
  return category.uploadConfig ? t('admin.catalog.categories.rules.custom') : t('admin.catalog.categories.rules.plain')
}

function uploadRuleFieldCount(category: AdminCatalogCategory) {
  return Array.isArray(category.uploadConfig?.fields) ? category.uploadConfig.fields.length : 0
}

function selectCategory(category: AdminCatalogCategory) {
  selectedId.value = category.id
  setFormNames(category.nameI18N)
  form.slug = category.slug
  form.sortOrder = category.sortOrder
  form.enabled = category.enabled
  form.uploadConfigText = category.uploadConfig ? JSON.stringify(category.uploadConfig, null, 2) : ''
  formError.value = ''
  originalFormSnapshot.value = formSnapshot()
}

function startCreate() {
  selectedId.value = null
  setFormNames()
  form.slug = ''
  form.sortOrder = 0
  form.enabled = true
  form.uploadConfigText = ''
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
    slug: form.slug,
    sortOrder: Number(form.sortOrder || 0),
    enabled: form.enabled,
    uploadConfigText: form.uploadConfigText
  })
}

function buildInput(): AdminCatalogCategoryInput {
  const primaryName = nameValue(primaryLocaleCode.value).trim()
  const nameI18N = localeOptions.value.reduce<Record<string, string>>((names, item) => {
    names[item.code] = nameValue(item.code).trim() || primaryName
    return names
  }, {})

  const uploadConfig = parseUploadConfigText()

  return {
    nameI18N,
    slug: form.slug.trim(),
    sortOrder: Number(form.sortOrder || 0),
    enabled: form.enabled,
    uploadConfig
  }
}

function ensureFormNameLocales() {
  for (const item of localeOptions.value) {
    if (!(item.code in form.names)) {
      form.names[item.code] = ''
    }
  }
}

function setFormNames(values?: Record<string, unknown> | null) {
  const supportedCodes = new Set(localeOptions.value.map((item) => item.code))
  for (const code of Object.keys(form.names)) {
    if (!supportedCodes.has(code)) {
      delete form.names[code]
    }
  }
  for (const item of localeOptions.value) {
    const value = values?.[item.code]
    form.names[item.code] = typeof value === 'string' ? value : ''
  }
}

function nameValue(code: string) {
  return form.names[code] || ''
}

function parseUploadConfigText() {
  const state = uploadConfigState.value
  if (!state.valid) {
    throw new Error(state.message || t('admin.catalog.categories.form.uploadConfigInvalid'))
  }
  return state.config
}

function parseUploadConfigJson(raw: string): UploadConfigParseState {
  const text = raw.trim()
  if (!text) return { valid: true, config: null, message: '' }

  try {
    const parsed = JSON.parse(text) as unknown
    if (parsed === null) return { valid: true, config: null, message: '' }
    if (typeof parsed !== 'object' || Array.isArray(parsed)) {
      return {
        valid: false,
        config: null,
        message: t('admin.catalog.categories.form.uploadConfigObjectInvalid')
      }
    }
    const config = parsed as UploadConfig
    const message = validateUploadConfig(config)
    if (message) {
      return { valid: false, config, message }
    }
    return { valid: true, config, message: '' }
  } catch (error: unknown) {
    return {
      valid: false,
      config: null,
      message: error instanceof Error ? error.message : t('admin.catalog.categories.form.uploadConfigInvalid')
    }
  }
}

function validateUploadConfig(config: UploadConfig): string {
  const records = config as UploadConfigRecord
  const fieldKeys = new Set<string>()

  if (records.fields !== undefined) {
    if (!Array.isArray(records.fields)) {
      return t('admin.catalog.categories.form.uploadConfigFieldsInvalid')
    }

    for (const [index, field] of records.fields.entries()) {
      if (!isPlainObject(field)) {
        return t('admin.catalog.categories.form.uploadConfigFieldInvalid', { index: index + 1 })
      }

      const key = typeof field.key === 'string' ? field.key.trim() : ''
      if (!key) {
        return t('admin.catalog.categories.form.uploadConfigFieldKeyRequired', { index: index + 1 })
      }
      if (fieldKeys.has(key)) {
        return t('admin.catalog.categories.form.uploadConfigFieldKeyDuplicate', { key })
      }
      fieldKeys.add(key)

      const type = typeof field.type === 'string' ? field.type : ''
      if (type === 'select' || type === 'multiSelect') {
        const message = validateUploadFieldOptions(key, field)
        if (message) return message
      }
    }
  }

  if (records.title !== undefined) {
    if (!isPlainObject(records.title)) {
      return t('admin.catalog.categories.form.uploadConfigTitleInvalid')
    }

    const mode = typeof records.title.mode === 'string' ? records.title.mode : ''
    if (records.title.parts !== undefined) {
      if (!Array.isArray(records.title.parts)) {
        return t('admin.catalog.categories.form.uploadConfigTitlePartsInvalid')
      }
      if (mode === 'generated' && records.title.parts.length === 0) {
        return t('admin.catalog.categories.form.uploadConfigGeneratedPartsRequired')
      }

      for (const [index, part] of records.title.parts.entries()) {
        if (!isPlainObject(part)) {
          return t('admin.catalog.categories.form.uploadConfigTitlePartInvalid', { index: index + 1 })
        }
        const field = typeof part.field === 'string' ? part.field.trim() : ''
        if (!field) {
          return t('admin.catalog.categories.form.uploadConfigTitlePartFieldRequired', { index: index + 1 })
        }
        if (!fieldKeys.has(field)) {
          return t('admin.catalog.categories.form.uploadConfigTitlePartFieldUnknown', { field })
        }
      }
    } else if (mode === 'generated') {
      return t('admin.catalog.categories.form.uploadConfigGeneratedPartsRequired')
    }
  }

  return ''
}

function validateUploadFieldOptions(fieldKey: string, field: UploadConfigRecord): string {
  if (!isPlainObject(field.options)) {
    return t('admin.catalog.categories.form.uploadConfigOptionsRequired', { field: fieldKey })
  }

  const source = typeof field.options.source === 'string' ? field.options.source : ''
  if (source === 'static') {
    if (!Array.isArray(field.options.items) || field.options.items.length === 0) {
      return t('admin.catalog.categories.form.uploadConfigStaticOptionsRequired', { field: fieldKey })
    }
    return ''
  }

  if (source === 'tagGroup') {
    const slug = typeof field.options.slug === 'string' ? field.options.slug.trim() : ''
    if (!slug) {
      return t('admin.catalog.categories.form.uploadConfigTagGroupSlugRequired', { field: fieldKey })
    }
    return ''
  }

  return t('admin.catalog.categories.form.uploadConfigUnsupportedOptionSource', {
    field: fieldKey,
    source: source || '-'
  })
}

function isPlainObject(value: unknown): value is UploadConfigRecord {
  return typeof value === 'object' && value !== null && !Array.isArray(value)
}

function uploadRuleTypeLabel(config: UploadConfig | null, fieldCount: number) {
  if (!config) return t('admin.catalog.categories.form.uploadRulePlain')
  if (fieldCount > 0) return t('admin.catalog.categories.form.uploadRuleStructured')
  return t('admin.catalog.categories.form.uploadRuleCustom')
}

function uploadRuleTitleLabel(config: UploadConfig | null, partCount: number) {
  if (!config) return t('admin.catalog.categories.form.uploadTitlePlain')
  if (config.title?.mode === 'manual') return t('admin.catalog.categories.form.uploadTitleManual')
  if (config.title?.mode === 'generated' || partCount > 0) {
    return t('admin.catalog.categories.form.uploadTitleGenerated', { count: numberFormatter.value.format(partCount) })
  }
  return t('admin.catalog.categories.form.uploadTitleCustom')
}

function clearUploadConfig() {
  form.uploadConfigText = ''
}

function formatUploadConfig() {
  if (!uploadConfigState.value.config) return
  form.uploadConfigText = JSON.stringify(uploadConfigState.value.config, null, 2)
}

async function saveCategory() {
  if (!canSubmit.value) return

  saving.value = true
  formError.value = ''
  try {
    const input = buildInput()
    const editingId = selectedId.value
    if (editingId) {
      await adminApi.updateCatalogCategory(editingId, input)
    } else {
      await adminApi.createCatalogCategory(input)
    }
    toast.add({ title: t('admin.catalog.categories.saved') })
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

async function deleteCategory(category: AdminCatalogCategory, close?: () => void) {
  deletingId.value = category.id
  try {
    await adminApi.deleteCatalogCategory(category.id)
    if (selectedId.value === category.id) {
      startCreate()
    }
    toast.add({ title: t('admin.catalog.categories.deleted') })
    close?.()
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
