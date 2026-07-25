<template>
  <div class="min-h-[calc(100vh-4rem)] bg-slate-50 py-6 dark:bg-slate-950">
    <div class="w-full px-3 sm:px-4 lg:px-5">
      <form class="grid min-w-0 grid-cols-1 gap-4 lg:grid-cols-[minmax(0,1fr)_340px] lg:items-start" @submit.prevent="handleSubmit">
        <main class="min-w-0 space-y-4">
          <section id="upload-file" class="overflow-hidden rounded-lg border border-slate-200 bg-white dark:border-slate-800 dark:bg-slate-900">
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
              <UFormField id="upload-category" :label="$t('catalog.torrents.upload.fields.category')" required>
                <USelect
                  v-model="form.categoryId"
                  class="w-full"
                  size="lg"
                  :ui="{ base: 'h-10 w-full' }"
                  :items="categoryOptions"
                  value-key="value"
                  :disabled="pending || categoriesPending || categories.length === 0"
                />
              </UFormField>

              <div v-if="categoriesError" class="rounded-md border border-red-200 bg-red-50 px-3 py-2 text-sm text-red-700 dark:border-red-900 dark:bg-red-950 dark:text-red-200">
                {{ categoriesError }}
              </div>

              <CatalogTorrentReleaseFieldsForm
                v-model="releaseFields"
                :category="selectedCategory"
                :tag-groups="tagGroups"
                :disabled="pending"
                :options-pending="tagGroupsPending"
                :show-errors="submitAttempted"
                @state-change="releaseState = $event"
              />

              <UFormField v-if="showManualTitleInput" :label="manualTitleLabel">
                <UInput v-model="form.name" class="w-full" :disabled="pending" @update:model-value="titleManuallyEdited = true" />
              </UFormField>

              <div id="upload-title" :class="finalTitleCardClass">
                <div class="flex items-center justify-between gap-3">
                  <p class="text-xs font-medium text-slate-500 dark:text-slate-400">{{ $t('catalog.torrents.upload.finalTitle') }}</p>
                  <UBadge v-if="isGeneratedTitleMode" color="neutral" variant="soft">{{ $t('catalog.torrents.upload.generatedTitle') }}</UBadge>
                </div>
                <p class="mt-1 break-words text-sm font-semibold text-slate-950 dark:text-white">{{ finalTitle || $t('catalog.torrents.upload.finalTitleEmpty') }}</p>
                <p v-if="submitAttempted && titleError" class="mt-2 text-xs text-red-600 dark:text-red-300">{{ titleError }}</p>
              </div>

              <UFormField :label="$t('catalog.torrents.upload.fields.subTitle')">
                <UInput v-model="form.subTitle" class="w-full" :disabled="pending" />
              </UFormField>

              <CatalogTorrentMetadataForm v-model="metadataBinding" :disabled="pending" />

              <CatalogTorrentDescriptionEditor v-model="form.description" :disabled="pending" :rows="12" />
            </div>
          </section>
        </main>

        <aside class="app-sticky-offset min-w-0 space-y-4 lg:sticky">
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

              <div v-if="submitAttempted && submitError" class="rounded-md border border-red-200 bg-red-50 px-3 py-2 text-sm text-red-700 dark:border-red-900 dark:bg-red-950 dark:text-red-200">
                {{ submitError }}
              </div>

              <div class="grid grid-cols-2 gap-2">
                <UButton color="neutral" variant="outline" block :to="localePath('/catalog/torrents')">
                  {{ $t('common.cancel') }}
                </UButton>
                <AppPermissionButton
                  :permission="Permission.CatalogTorrentCreate"
                  type="submit"
                  color="primary"
                  icon="i-lucide-upload"
                  block
                  :tooltip="$t('catalog.torrents.upload.submit')"
                  :loading="pending"
                  :disabled="!canSubmit"
                >
                  {{ $t('catalog.torrents.upload.submit') }}
                </AppPermissionButton>
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
import type { CatalogCategory, CatalogTagGroup, ReleaseFieldsState, ReleaseFieldValue, TorrentMetadataBinding } from '~/composables/useCatalogTorrents'
import { localizeI18nName } from '~/utils/format'

definePageMeta({
  middleware: 'auth'
})

const { t, locale } = useI18n()
const localePath = useLocalePath()
const route = useRoute()
const toast = useToast()
const catalogTorrents = useCatalogTorrents()
const { hasPermission } = useAuth()

const categories = ref<CatalogCategory[]>([])
const tagGroups = ref<CatalogTagGroup[]>([])
const categoriesPending = ref(true)
const tagGroupsPending = ref(false)
const categoriesError = ref('')
const selectedFile = ref<File | null>(null)
const fileInputKey = ref(0)
const titleManuallyEdited = ref(false)
const submitAttempted = ref(false)
const pending = ref(false)
const numberFormatter = computed(() => new Intl.NumberFormat(locale.value))

const form = reactive({
  categoryId: String(readCategoryIdQuery()),
  name: '',
  subTitle: '',
  description: '',
  anonymous: false
})
const releaseFields = ref<Record<string, ReleaseFieldValue>>({})
const metadataBinding = ref<TorrentMetadataBinding>({
  imdbId: '',
  doubanId: '',
  bangumiId: '',
  tmdbId: '',
  tmdbType: ''
})
const releaseState = ref<ReleaseFieldsState>({
  generatedTitle: '',
  valid: true,
  firstError: '',
  missingLabels: [],
  output: {}
})

const selectedCategory = computed(() => {
  const categoryId = Number(form.categoryId)
  return categories.value.find((item) => item.id === categoryId) || null
})

const selectedCategoryName = computed(() => {
  return selectedCategory.value ? categoryDisplayName(selectedCategory.value) : '-'
})
const categoryOptions = computed(() => [
  { value: '0', label: t('catalog.torrents.upload.fields.categoryPlaceholder') },
  ...categories.value.map(category => ({ value: String(category.id), label: categoryDisplayName(category) }))
])

const selectedUploadConfig = computed(() => selectedCategory.value?.uploadConfig || null)
const isGeneratedTitleMode = computed(() => selectedUploadConfig.value?.title?.mode === 'generated')
const showManualTitleInput = computed(() => !isGeneratedTitleMode.value || Boolean(selectedUploadConfig.value?.title?.allowManualOverride))
const manualTitleLabel = computed(() => isGeneratedTitleMode.value ? t('catalog.torrents.upload.fields.titleOverride') : t('catalog.torrents.upload.fields.name'))
const generatedTitle = computed(() => releaseState.value.generatedTitle)
const releaseFieldsValid = computed(() => releaseState.value.valid)

const effectiveTitle = computed(() => {
  const manualTitle = form.name.trim()
  if (isGeneratedTitleMode.value) {
    if (showManualTitleInput.value && titleManuallyEdited.value && manualTitle) return manualTitle
    if (generatedTitle.value) return generatedTitle.value
  }
  if (manualTitle) return manualTitle
  return selectedFile.value?.name.replace(/\.torrent$/i, '') || '-'
})

const finalTitle = computed(() => effectiveTitle.value === '-' ? '' : effectiveTitle.value)
const titleError = computed(() => finalTitle.value ? '' : t('catalog.torrents.upload.errors.titleRequired'))
const canCreateTorrent = computed(() => hasPermission(Permission.CatalogTorrentCreate))
const submitError = computed(() => {
  if (!canCreateTorrent.value) return t('common.noPermission')
  if (!selectedFile.value) return t('catalog.torrents.upload.errors.fileRequired')
  if (!selectedCategory.value) return t('catalog.torrents.upload.errors.categoryRequired')
  if (tagGroupsPending.value) return t('catalog.torrents.upload.errors.optionsLoading')
  if (!releaseFieldsValid.value) return releaseState.value.firstError
  if (titleError.value) return titleError.value
  return ''
})
const canSubmit = computed(() => Boolean(!submitError.value && !categoriesPending.value && !pending.value))
const finalTitleCardClass = computed(() => [
  'rounded-md border bg-white px-3 py-2.5 dark:bg-slate-950',
  submitAttempted.value && titleError.value
    ? 'border-red-200 dark:border-red-900'
    : 'border-slate-200 dark:border-slate-800'
].join(' '))
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
  submitAttempted.value = true
  if (!canSubmit.value || !selectedFile.value || Number(form.categoryId) <= 0) {
    toast.add({
      title: submitError.value || t('common.requestFailed'),
      color: 'error',
      icon: 'i-lucide-circle-alert'
    })
    scrollToUploadIssue()
    return
  }

  pending.value = true
  try {
    const out = await catalogTorrents.uploadTorrent({
      file: selectedFile.value,
      categoryId: Number(form.categoryId),
      name: submitTitleValue(),
      subTitle: form.subTitle,
      description: form.description,
      releaseFields: releaseState.value.output,
      metadata: metadataBinding.value,
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

function categoryNeedsTagGroups(category: CatalogCategory | null) {
  return Boolean(category?.uploadConfig?.fields?.some((field) => field.options?.source === 'tagGroup'))
}

function submitTitleValue() {
  if (isGeneratedTitleMode.value && !(showManualTitleInput.value && titleManuallyEdited.value)) {
    return ''
  }
  return form.name
}

function scrollToUploadIssue() {
  let target = '#upload-title'
  if (!selectedFile.value) {
    target = '#upload-file'
  } else if (!selectedCategory.value) {
    target = '#upload-category'
  } else if (!releaseFieldsValid.value) {
    target = '#upload-release-fields'
  }
  document.querySelector(target)?.scrollIntoView({ behavior: 'smooth', block: 'center' })
}

useSeoMeta({
  title: t('catalog.torrents.upload.metaTitle'),
  robots: 'noindex, nofollow'
})
</script>
