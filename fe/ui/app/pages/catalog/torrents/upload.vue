<template>
  <div class="min-h-[calc(100vh-4rem)] bg-slate-50 py-6 dark:bg-slate-950">
    <div class="w-full px-3 sm:px-4 lg:px-5">
      <form class="grid min-w-0 grid-cols-1 gap-4 lg:grid-cols-[minmax(0,1fr)_340px] lg:items-start" @submit.prevent="handleSubmit">
        <main class="min-w-0 space-y-4">
          <section class="overflow-hidden rounded-lg border border-slate-200 bg-white dark:border-slate-800 dark:bg-slate-900">
            <div class="flex items-center justify-between gap-3 border-b border-slate-200 px-4 py-3 dark:border-slate-800">
              <h2 class="text-base font-semibold text-slate-950 dark:text-white">{{ $t('catalog.torrents.upload.sections.file') }}</h2>
              <UBadge :color="selectedFile ? 'success' : 'neutral'" variant="soft">
                {{ selectedFile ? $t('catalog.torrents.upload.file.selected') : $t('catalog.torrents.upload.file.required') }}
              </UBadge>
            </div>

            <div class="p-4">
              <div
                v-if="selectedFile"
                class="grid gap-3 rounded-md border border-sky-200 bg-sky-50 p-3 sm:grid-cols-[minmax(0,1fr)_auto] sm:items-center dark:border-sky-800 dark:bg-sky-950/40"
                @dragover.prevent
                @drop.prevent="handleDrop"
              >
                <div class="flex min-w-0 items-center gap-3">
                  <span class="flex size-11 shrink-0 items-center justify-center rounded-md bg-white text-sky-600 ring-1 ring-sky-200 dark:bg-slate-900 dark:text-sky-300 dark:ring-sky-800">
                    <UIcon name="i-lucide-file-check-2" class="size-5" />
                  </span>
                  <div class="min-w-0">
                    <p class="truncate text-sm font-semibold text-slate-950 dark:text-white">{{ selectedFile.name }}</p>
                    <p class="mt-1 text-xs text-slate-500 dark:text-slate-400">{{ formatBytes(selectedFile.size) }}</p>
                  </div>
                </div>

                <div class="grid grid-cols-2 gap-2 sm:flex sm:items-center sm:justify-end">
                  <label class="inline-flex h-8 cursor-pointer items-center justify-center rounded-md border border-slate-200 bg-white px-3 text-sm font-medium text-slate-700 transition hover:bg-slate-50 dark:border-slate-700 dark:bg-slate-900 dark:text-slate-200 dark:hover:bg-slate-800">
                    {{ $t('catalog.torrents.upload.file.replace') }}
                    <input
                      :key="fileInputKey"
                      class="sr-only"
                      type="file"
                      accept=".torrent,application/x-bittorrent"
                      :disabled="pending"
                      @change="handleFileChange"
                    >
                  </label>
                  <UButton type="button" color="neutral" variant="ghost" size="sm" icon="i-lucide-x" :disabled="pending" @click="clearFile">
                    {{ $t('catalog.torrents.upload.file.remove') }}
                  </UButton>
                </div>
              </div>

              <label
                v-else
                class="flex min-h-36 cursor-pointer flex-col items-center justify-center rounded-md border border-dashed border-slate-300 bg-slate-50 px-4 py-8 text-center transition hover:border-slate-400 dark:border-slate-700 dark:bg-slate-950 dark:hover:border-slate-600"
                @dragover.prevent
                @drop.prevent="handleDrop"
              >
                <input
                  :key="fileInputKey"
                  class="sr-only"
                  type="file"
                  accept=".torrent,application/x-bittorrent"
                  :disabled="pending"
                  @change="handleFileChange"
                >
                <UIcon name="i-lucide-file-up" class="size-9 text-slate-400" />
                <span class="mt-3 text-sm font-medium text-slate-950 dark:text-white">
                  {{ $t('catalog.torrents.upload.file.choose') }}
                </span>
              </label>
            </div>
          </section>

          <section class="overflow-hidden rounded-lg border border-slate-200 bg-white dark:border-slate-800 dark:bg-slate-900">
            <div class="border-b border-slate-200 px-4 py-3 dark:border-slate-800">
              <h2 class="text-base font-semibold text-slate-950 dark:text-white">{{ $t('catalog.torrents.upload.sections.info') }}</h2>
            </div>

            <div class="grid grid-cols-1 gap-4 p-4">
              <UFormField :label="$t('catalog.torrents.upload.fields.category')" required>
                <select
                  v-model="form.categoryId"
                  class="h-10 w-full rounded-md border border-slate-200 bg-white px-3 text-sm text-slate-950 outline-none transition focus:border-sky-400 focus:ring-2 focus:ring-sky-100 dark:border-slate-700 dark:bg-slate-950 dark:text-white dark:focus:border-sky-500 dark:focus:ring-sky-950"
                  :disabled="pending || categoriesPending || categories.length === 0"
                >
                  <option value="0">{{ $t('catalog.torrents.upload.fields.categoryPlaceholder') }}</option>
                  <option v-for="category in categories" :key="category.id" :value="String(category.id)">
                    {{ categoryDisplayName(category) }}
                  </option>
                </select>
              </UFormField>

              <div v-if="categoriesError" class="rounded-md border border-red-200 bg-red-50 px-3 py-2 text-sm text-red-700 dark:border-red-900 dark:bg-red-950 dark:text-red-200">
                {{ categoriesError }}
              </div>

              <div v-if="schemaFields.length > 0" class="rounded-md border border-slate-200 bg-slate-50/70 p-3 dark:border-slate-800 dark:bg-slate-950/40">
                <div class="mb-3 flex items-center justify-between gap-3">
                  <h3 class="text-sm font-semibold text-slate-950 dark:text-white">{{ $t('catalog.torrents.upload.sections.release') }}</h3>
                  <UBadge color="neutral" variant="soft">{{ selectedCategoryName }}</UBadge>
                </div>

                <div class="grid grid-cols-1 gap-4">
                  <UFormField
                    v-for="field in schemaFields"
                    :key="field.key"
                    :label="fieldLabel(field)"
                    :description="fieldDescription(field)"
                    :required="Boolean(field.required)"
                  >
                    <select
                      v-if="field.type === 'select'"
                      :value="String(releaseFields[field.key] || '')"
                      class="h-10 w-full rounded-md border border-slate-200 bg-white px-3 text-sm text-slate-950 outline-none transition focus:border-sky-400 focus:ring-2 focus:ring-sky-100 dark:border-slate-700 dark:bg-slate-950 dark:text-white dark:focus:border-sky-500 dark:focus:ring-sky-950"
                      :disabled="pending || tagGroupsPending"
                      @change="handleReleaseSelectChange(field.key, $event)"
                    >
                      <option value="">{{ $t('catalog.torrents.upload.fields.optionPlaceholder') }}</option>
                      <option v-for="option in fieldOptions(field)" :key="option.value" :value="option.value">
                        {{ optionLabel(option) }}
                      </option>
                    </select>

                    <div v-else-if="field.type === 'multiSelect'" class="grid gap-2 rounded-md border border-slate-200 bg-white p-2 sm:grid-cols-2 dark:border-slate-700 dark:bg-slate-950">
                      <label
                        v-for="option in fieldOptions(field)"
                        :key="option.value"
                        class="flex min-h-8 cursor-pointer items-center gap-2 rounded px-2 text-sm text-slate-700 transition hover:bg-slate-50 dark:text-slate-200 dark:hover:bg-slate-900"
                      >
                        <input
                          type="checkbox"
                          class="size-4 rounded border-slate-300 text-sky-600 focus:ring-sky-500 dark:border-slate-600"
                          :checked="releaseFieldListValue(field.key).includes(option.value)"
                          :disabled="pending || tagGroupsPending"
                          @change="toggleReleaseFieldOption(field.key, option.value)"
                        >
                        <span class="min-w-0 truncate">{{ optionLabel(option) }}</span>
                      </label>
                      <p v-if="fieldOptions(field).length === 0" class="px-2 py-1 text-sm text-slate-500 dark:text-slate-400">
                        {{ $t('catalog.torrents.upload.fields.noOptions') }}
                      </p>
                    </div>

                    <UTextarea
                      v-else-if="field.type === 'textarea'"
                      :model-value="String(releaseFields[field.key] || '')"
                      class="w-full"
                      :rows="3"
                      :disabled="pending"
                      :placeholder="fieldPlaceholder(field)"
                      @update:model-value="setReleaseField(field.key, String($event || ''))"
                    />

                    <UInput
                      v-else
                      :model-value="String(releaseFields[field.key] || '')"
                      class="w-full"
                      :disabled="pending"
                      :placeholder="fieldPlaceholder(field)"
                      @update:model-value="setReleaseField(field.key, String($event || ''))"
                    />
                  </UFormField>
                </div>
              </div>

              <div v-if="isGeneratedTitleMode" class="rounded-md border border-slate-200 bg-white px-3 py-2.5 dark:border-slate-800 dark:bg-slate-950">
                <p class="text-xs font-medium text-slate-500 dark:text-slate-400">{{ $t('catalog.torrents.upload.generatedTitle') }}</p>
                <p class="mt-1 break-words text-sm font-semibold text-slate-950 dark:text-white">{{ generatedTitle || '-' }}</p>
              </div>

              <UFormField v-if="showManualTitleInput" :label="manualTitleLabel">
                <UInput v-model="form.name" class="w-full" :disabled="pending" @update:model-value="titleManuallyEdited = true" />
              </UFormField>

              <UFormField :label="$t('catalog.torrents.upload.fields.subTitle')">
                <UInput v-model="form.subTitle" class="w-full" :disabled="pending" />
              </UFormField>

              <UFormField :label="$t('catalog.torrents.upload.fields.description')">
                <UTextarea
                  v-if="descriptionMode === 'write'"
                  v-model="form.description"
                  class="w-full"
                  :rows="12"
                  :disabled="pending"
                />
                <div v-else class="min-h-72 rounded-md border border-slate-200 bg-slate-50 px-3 py-2.5 dark:border-slate-800 dark:bg-slate-950">
                  <div v-if="renderedDescriptionPreview" class="rich-text" v-html="renderedDescriptionPreview" />
                  <p v-else class="text-sm text-slate-500 dark:text-slate-400">{{ $t('catalog.torrents.upload.preview.empty') }}</p>
                </div>

                <div class="mt-3 flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
                  <div class="inline-flex w-fit rounded-md border border-slate-200 bg-slate-50 p-0.5 dark:border-slate-800 dark:bg-slate-950">
                    <button
                      type="button"
                      :class="descriptionModeButtonClass('write')"
                      @click="descriptionMode = 'write'"
                    >
                      {{ $t('catalog.torrents.upload.preview.write') }}
                    </button>
                    <button
                      type="button"
                      :class="descriptionModeButtonClass('preview')"
                      @click="descriptionMode = 'preview'"
                    >
                      {{ $t('catalog.torrents.upload.preview.preview') }}
                    </button>
                  </div>
                  <span class="text-xs text-slate-500 dark:text-slate-400">
                    {{ $t('catalog.torrents.upload.summary.descriptionLength', { count: numberFormatter.format(form.description.trim().length) }) }}
                  </span>
                </div>
              </UFormField>
            </div>
          </section>
        </main>

        <aside class="min-w-0 space-y-4 lg:sticky lg:top-[5.5rem]">
          <section class="rounded-lg border border-slate-200 bg-white p-4 dark:border-slate-800 dark:bg-slate-900">
            <h2 class="text-sm font-semibold text-slate-950 dark:text-white">{{ $t('catalog.torrents.upload.sections.publish') }}</h2>

            <div class="mt-3 space-y-4">
              <label class="flex items-center justify-between gap-4 rounded-md border border-slate-200 px-3 py-2 dark:border-slate-800">
                <span class="text-sm font-medium text-slate-700 dark:text-slate-200">{{ $t('catalog.torrents.upload.fields.anonymous') }}</span>
                <input
                  v-model="form.anonymous"
                  type="checkbox"
                  class="size-4 rounded border-slate-300 text-sky-600 focus:ring-sky-500 dark:border-slate-600"
                  :disabled="pending"
                >
              </label>

              <div class="divide-y divide-slate-100 rounded-md border border-slate-200 dark:divide-slate-800 dark:border-slate-800">
                <div
                  v-for="check in publishChecks"
                  :key="check.key"
                  class="flex items-center justify-between gap-3 px-3 py-2.5 text-sm"
                >
                  <div class="flex min-w-0 items-center gap-2">
                    <UIcon :name="check.passed ? 'i-lucide-circle-check' : 'i-lucide-circle'" :class="check.passed ? 'text-emerald-500' : 'text-slate-300 dark:text-slate-600'" class="size-4 shrink-0" />
                    <span class="text-slate-600 dark:text-slate-300">{{ check.label }}</span>
                  </div>
                  <span class="min-w-0 truncate text-right text-xs font-medium text-slate-500 dark:text-slate-400">{{ check.value }}</span>
                </div>
              </div>

              <div class="grid grid-cols-2 gap-2">
                <UButton color="neutral" variant="outline" block :to="localePath('/catalog/torrents')">
                  {{ $t('common.cancel') }}
                </UButton>
                <UButton type="submit" color="primary" icon="i-lucide-upload" block :loading="pending" :disabled="!canSubmit">
                  {{ $t('catalog.torrents.upload.submit') }}
                </UButton>
              </div>
            </div>
          </section>
        </aside>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ApiError } from '~/composables/useApi'
import type { CatalogCategory, CatalogTagGroup, UploadFieldConfig, UploadOptionItem, UploadTitlePart } from '~/composables/useCatalogTorrents'
import { renderUserMarkdown } from '~/utils/richText'

type DescriptionMode = 'write' | 'preview'
type ReleaseFieldValue = string | string[]

definePageMeta({
  middleware: 'auth'
})

const { t, locale } = useI18n()
const localePath = useLocalePath()
const route = useRoute()
const toast = useToast()
const catalogTorrents = useCatalogTorrents()

const categories = ref<CatalogCategory[]>([])
const tagGroups = ref<CatalogTagGroup[]>([])
const categoriesPending = ref(true)
const tagGroupsPending = ref(false)
const categoriesError = ref('')
const selectedFile = ref<File | null>(null)
const fileInputKey = ref(0)
const descriptionMode = ref<DescriptionMode>('write')
const titleManuallyEdited = ref(false)
const pending = ref(false)
const numberFormatter = computed(() => new Intl.NumberFormat(locale.value))

const form = reactive({
  categoryId: String(readCategoryIdQuery()),
  name: '',
  subTitle: '',
  description: '',
  anonymous: false
})
const releaseFields = reactive<Record<string, ReleaseFieldValue>>({})

const selectedCategory = computed(() => {
  const categoryId = Number(form.categoryId)
  return categories.value.find((item) => item.id === categoryId) || null
})

const selectedCategoryName = computed(() => {
  return selectedCategory.value ? categoryDisplayName(selectedCategory.value) : '-'
})

const selectedUploadConfig = computed(() => selectedCategory.value?.uploadConfig || null)
const schemaFields = computed(() => selectedUploadConfig.value?.fields || [])
const isGeneratedTitleMode = computed(() => selectedUploadConfig.value?.title?.mode === 'generated')
const showManualTitleInput = computed(() => !isGeneratedTitleMode.value || Boolean(selectedUploadConfig.value?.title?.allowManualOverride))
const manualTitleLabel = computed(() => isGeneratedTitleMode.value ? t('catalog.torrents.upload.fields.titleOverride') : t('catalog.torrents.upload.fields.name'))
const generatedTitle = computed(() => buildGeneratedTitle(selectedUploadConfig.value?.title?.parts || []))
const releaseFieldsValid = computed(() => schemaFields.value.every((field) => {
  if (!field.required) return true
  const value = releaseFields[field.key]
  return Array.isArray(value) ? value.length > 0 : Boolean(String(value || '').trim())
}))

const effectiveTitle = computed(() => {
  const manualTitle = form.name.trim()
  if (isGeneratedTitleMode.value) {
    if (showManualTitleInput.value && titleManuallyEdited.value && manualTitle) return manualTitle
    if (generatedTitle.value) return generatedTitle.value
  }
  if (manualTitle) return manualTitle
  return selectedFile.value?.name.replace(/\.torrent$/i, '') || '-'
})

const renderedDescriptionPreview = computed(() => renderUserMarkdown(form.description).trim())
const canSubmit = computed(() => Boolean(selectedFile.value && selectedCategory.value && releaseFieldsValid.value && effectiveTitle.value !== '-' && !categoriesPending.value && !tagGroupsPending.value && !pending.value))
const publishChecks = computed(() => [
  {
    key: 'file',
    label: t('catalog.torrents.upload.checks.file'),
    value: selectedFile.value ? formatBytes(selectedFile.value.size) : '-',
    passed: Boolean(selectedFile.value)
  },
  {
    key: 'category',
    label: t('catalog.torrents.upload.checks.category'),
    value: selectedCategoryName.value,
    passed: Boolean(selectedCategory.value)
  },
  {
    key: 'title',
    label: t('catalog.torrents.upload.checks.title'),
    value: numberFormatter.value.format(effectiveTitle.value === '-' ? 0 : effectiveTitle.value.length),
    passed: effectiveTitle.value !== '-'
  },
  {
    key: 'description',
    label: t('catalog.torrents.upload.checks.description'),
    value: numberFormatter.value.format(form.description.trim().length),
    passed: form.description.trim().length > 0
  }
])

onMounted(loadCategories)

watch(selectedCategory, (category) => {
  initializeReleaseFields()
  if (categoryNeedsTagGroups(category)) {
    void loadTagGroups()
  }
  if (isGeneratedTitleMode.value && !titleManuallyEdited.value) {
    form.name = ''
  }
  if (!isGeneratedTitleMode.value && !titleManuallyEdited.value && selectedFile.value && !form.name.trim()) {
    form.name = selectedFile.value.name.replace(/\.torrent$/i, '')
  }
})

function readCategoryIdQuery() {
  const raw = Array.isArray(route.query.categoryId) ? route.query.categoryId[0] : route.query.categoryId
  const parsed = Number(raw)
  return Number.isInteger(parsed) && parsed > 0 ? parsed : 0
}

async function loadCategories() {
  categoriesPending.value = true
  categoriesError.value = ''
  try {
    const data = await catalogTorrents.listCategories()
    categories.value = data.list || []
    if (Number(form.categoryId) > 0 && !categories.value.some((item) => item.id === Number(form.categoryId))) {
      form.categoryId = '0'
    }
    if (categories.value.some(categoryNeedsTagGroups)) {
      await loadTagGroups()
    }
  } catch (error) {
    categories.value = []
    categoriesError.value = error instanceof ApiError ? error.message : t('common.requestFailed')
  } finally {
    categoriesPending.value = false
  }
}

async function loadTagGroups() {
  if (tagGroupsPending.value || tagGroups.value.length > 0) return
  tagGroupsPending.value = true
  try {
    const data = await catalogTorrents.listTagGroups()
    tagGroups.value = data.list || []
  } catch (error) {
    toast.add({
      title: error instanceof ApiError ? error.message : t('common.requestFailed'),
      color: 'error',
      icon: 'i-lucide-circle-alert'
    })
  } finally {
    tagGroupsPending.value = false
  }
}

function handleFileChange(event: Event) {
  const input = event.target as HTMLInputElement
  setFile(input.files?.[0] || null)
}

function handleDrop(event: DragEvent) {
  setFile(event.dataTransfer?.files?.[0] || null)
}

function clearFile() {
  selectedFile.value = null
  fileInputKey.value += 1
}

function setFile(file: File | null) {
  if (!file) {
    clearFile()
    return
  }

  if (!file.name.toLowerCase().endsWith('.torrent')) {
    toast.add({
      title: t('catalog.torrents.upload.file.invalid'),
      color: 'error',
      icon: 'i-lucide-circle-alert'
    })
    return
  }

  selectedFile.value = file
  fileInputKey.value += 1
  if (!isGeneratedTitleMode.value && !titleManuallyEdited.value && !form.name.trim()) {
    form.name = file.name.replace(/\.torrent$/i, '')
  }
}

async function handleSubmit() {
  if (!selectedFile.value || Number(form.categoryId) <= 0) return

  pending.value = true
  try {
    const out = await catalogTorrents.uploadTorrent({
      file: selectedFile.value,
      categoryId: Number(form.categoryId),
      name: submitTitleValue(),
      subTitle: form.subTitle,
      description: form.description,
      releaseFields: buildReleaseFieldsInput(),
      anonymous: form.anonymous
    })

    toast.add({
      title: t('catalog.torrents.upload.success', { id: out.torrentId }),
      color: 'success',
      icon: 'i-lucide-check-circle'
    })
    await navigateTo(localePath(`/catalog/torrents/${out.torrentId}`))
  } catch (error) {
    toast.add({
      title: error instanceof ApiError ? error.message : t('common.requestFailed'),
      color: 'error',
      icon: 'i-lucide-circle-alert'
    })
  } finally {
    pending.value = false
  }
}

function categoryDisplayName(category: CatalogCategory) {
  return localizeI18nName(category.name, locale.value, category.slug || `#${category.id}`)
}

function fieldLabel(field: UploadFieldConfig) {
  return localizeI18nName(field.label, locale.value, field.key)
}

function fieldDescription(field: UploadFieldConfig) {
  return localizeI18nName(field.description, locale.value, '')
}

function fieldPlaceholder(field: UploadFieldConfig) {
  return localizeI18nName(field.placeholder, locale.value, '')
}

function optionLabel(option: UploadOptionItem) {
  return localizeI18nName(option.label, locale.value, option.value)
}

function fieldOptions(field: UploadFieldConfig): UploadOptionItem[] {
  if (!field.options) return []
  if (field.options.source === 'static') return field.options.items || []
  if (field.options.source !== 'tagGroup' || !field.options.slug) return []

  const group = tagGroups.value.find((item) => item.slug === field.options?.slug && tagGroupAppliesToCategory(item))
  return (group?.tags || [])
    .filter((tag) => Boolean(tag.value))
    .map((tag) => ({ value: tag.value, label: tag.name }))
}

function tagGroupAppliesToCategory(group: CatalogTagGroup) {
  if (!selectedCategory.value || !group.categories?.length) return true
  return group.categories.includes(selectedCategory.value.id)
}

function categoryNeedsTagGroups(category: CatalogCategory | null) {
  return Boolean(category?.uploadConfig?.fields?.some((field) => field.options?.source === 'tagGroup'))
}

function initializeReleaseFields() {
  const nextKeys = new Set(schemaFields.value.map((field) => field.key).filter(Boolean))
  for (const key of Object.keys(releaseFields)) {
    if (!nextKeys.has(key)) delete releaseFields[key]
  }
  for (const field of schemaFields.value) {
    if (!field.key || releaseFields[field.key] !== undefined) continue
    releaseFields[field.key] = field.type === 'multiSelect' ? [] : ''
  }
}

function setReleaseField(key: string, value: string) {
  releaseFields[key] = value
}

function handleReleaseSelectChange(key: string, event: Event) {
  setReleaseField(key, (event.target as HTMLSelectElement).value)
}

function releaseFieldListValue(key: string) {
  const value = releaseFields[key]
  return Array.isArray(value) ? value : []
}

function toggleReleaseFieldOption(key: string, value: string) {
  const current = releaseFieldListValue(key)
  releaseFields[key] = current.includes(value)
    ? current.filter((item) => item !== value)
    : [...current, value]
}

function buildGeneratedTitle(parts: UploadTitlePart[]) {
  return parts
    .map((part) => {
      const value = titlePartValue(releaseFields[part.field], part.separator)
      return value ? `${part.prefix || ''}${value}${part.suffix || ''}` : ''
    })
    .join('')
    .trim()
}

function titlePartValue(value: ReleaseFieldValue | undefined, separator = '/') {
  if (Array.isArray(value)) return value.join(separator || '/')
  return String(value || '').trim()
}

function buildReleaseFieldsInput() {
  const input: Record<string, unknown> = {}
  for (const field of schemaFields.value) {
    const value = releaseFields[field.key]
    if (Array.isArray(value)) {
      const list = value.filter(Boolean)
      if (list.length > 0) input[field.key] = list
      continue
    }
    const text = String(value || '').trim()
    if (text) input[field.key] = text
  }
  return input
}

function submitTitleValue() {
  if (isGeneratedTitleMode.value && !(showManualTitleInput.value && titleManuallyEdited.value)) {
    return ''
  }
  return form.name
}

function descriptionModeButtonClass(mode: DescriptionMode) {
  const active = descriptionMode.value === mode
  return [
    'h-7 rounded px-3 text-xs font-medium transition-colors',
    active
      ? 'bg-white text-slate-950 shadow-sm ring-1 ring-slate-200 dark:bg-slate-800 dark:text-white dark:ring-slate-700'
      : 'text-slate-500 hover:text-slate-800 dark:text-slate-400 dark:hover:text-slate-100'
  ].join(' ')
}

useSeoMeta({
  title: t('catalog.torrents.upload.metaTitle'),
  robots: 'noindex, nofollow'
})
</script>
