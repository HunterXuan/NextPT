<template>
  <div class="min-h-[calc(100vh-4rem)] bg-slate-50 py-6 dark:bg-slate-950">
    <div class="w-full px-3 sm:px-4 lg:px-5">
      <div v-if="pending" class="h-96 animate-pulse rounded-lg border border-slate-200 bg-white dark:border-slate-800 dark:bg-slate-900" />

      <div v-else-if="errorMessage" class="flex flex-col items-center justify-center rounded-lg border border-slate-200 bg-white px-4 py-16 text-center dark:border-slate-800 dark:bg-slate-900">
        <UIcon name="i-lucide-circle-alert" class="size-9 text-red-500" />
        <p class="mt-3 text-sm font-medium text-slate-950 dark:text-white">{{ errorMessage }}</p>
      </div>

      <form v-else class="grid grid-cols-1 gap-6 lg:grid-cols-[minmax(0,1fr)_340px]" @submit.prevent="handleSubmit">
        <div class="space-y-6">
          <UCard class="rounded-lg">
            <template #header>
              <h2 class="text-base font-semibold text-slate-950 dark:text-white">{{ $t('catalog.torrents.edit.sections.info') }}</h2>
            </template>

            <div class="grid grid-cols-1 gap-4">
              <UFormField id="edit-category" :label="$t('catalog.torrents.upload.fields.category')" required>
                <select
                  v-model="form.categoryId"
                  class="h-10 w-full rounded-md border border-slate-200 bg-white px-3 text-sm text-slate-950 outline-none transition focus:border-sky-400 focus:ring-2 focus:ring-sky-100 dark:border-slate-700 dark:bg-slate-950 dark:text-white dark:focus:border-sky-500 dark:focus:ring-sky-950"
                  :disabled="savePending || categories.length === 0"
                >
                  <option value="0">{{ $t('catalog.torrents.upload.fields.categoryPlaceholder') }}</option>
                  <option v-for="category in categories" :key="category.id" :value="String(category.id)">
                    {{ categoryDisplayName(category) }}
                  </option>
                </select>
              </UFormField>

              <CatalogTorrentReleaseFieldsForm
                v-model="releaseFields"
                :category="selectedCategory"
                :tag-groups="tagGroups"
                :disabled="savePending"
                :options-pending="tagGroupsPending"
                :show-errors="submitAttempted"
                @state-change="handleReleaseStateChange"
              />

              <UFormField v-if="showManualTitleInput" :label="manualTitleLabel">
                <UInput v-model="form.name" class="w-full" :disabled="savePending" />
              </UFormField>

              <div id="edit-title" :class="finalTitleCardClass">
                <div class="flex items-center justify-between gap-3">
                  <p class="text-xs font-medium text-slate-500 dark:text-slate-400">{{ $t('catalog.torrents.upload.finalTitle') }}</p>
                  <UBadge v-if="isGeneratedTitleMode && generatedTitle && generatedTitle === finalTitle" color="neutral" variant="soft">{{ $t('catalog.torrents.upload.generatedTitle') }}</UBadge>
                </div>
                <p class="mt-1 break-words text-sm font-semibold text-slate-950 dark:text-white">{{ finalTitle || $t('catalog.torrents.upload.finalTitleEmpty') }}</p>
                <p v-if="isGeneratedTitleMode && generatedTitle && generatedTitle !== finalTitle" class="mt-2 break-words text-xs text-slate-500 dark:text-slate-400">
                  {{ $t('catalog.torrents.upload.generatedTitle') }}: {{ generatedTitle }}
                </p>
                <p v-if="submitAttempted && titleError" class="mt-2 text-xs text-red-600 dark:text-red-300">{{ titleError }}</p>
              </div>

              <UFormField :label="$t('catalog.torrents.upload.fields.subTitle')">
                <UInput v-model="form.subTitle" class="w-full" :disabled="savePending" />
              </UFormField>

              <CatalogTorrentDescriptionEditor v-model="form.description" :disabled="savePending" :rows="12" />
            </div>
          </UCard>
        </div>

        <aside class="space-y-6">
          <UCard class="rounded-lg">
            <template #header>
              <h2 class="text-base font-semibold text-slate-950 dark:text-white">{{ $t('catalog.torrents.edit.sections.publish') }}</h2>
            </template>

            <div class="space-y-4">
              <label class="flex items-center justify-between gap-4 rounded-md border border-slate-200 px-3 py-2 dark:border-slate-800">
                <span class="text-sm font-medium text-slate-700 dark:text-slate-200">{{ $t('catalog.torrents.upload.fields.anonymous') }}</span>
                <input
                  v-model="form.anonymous"
                  type="checkbox"
                  class="size-4 rounded border-slate-300 text-sky-600 focus:ring-sky-500 dark:border-slate-600"
                  :disabled="savePending"
                >
              </label>

              <dl class="space-y-3 text-sm">
                <div class="flex items-center justify-between gap-3">
                  <dt class="text-slate-500 dark:text-slate-400">{{ $t('catalog.torrents.edit.summary.torrentId') }}</dt>
                  <dd class="font-medium text-slate-950 dark:text-white">#{{ torrentId }}</dd>
                </div>
                <div class="flex items-center justify-between gap-3">
                  <dt class="text-slate-500 dark:text-slate-400">{{ $t('catalog.torrents.edit.summary.category') }}</dt>
                  <dd class="min-w-0 truncate font-medium text-slate-950 dark:text-white">{{ selectedCategoryName }}</dd>
                </div>
              </dl>

              <div v-if="submitAttempted && submitError" class="rounded-md border border-red-200 bg-red-50 px-3 py-2 text-sm text-red-700 dark:border-red-900 dark:bg-red-950 dark:text-red-200">
                {{ submitError }}
              </div>

              <div v-if="initialEditSnapshot && !hasChanges" class="rounded-md border border-slate-200 bg-slate-50 px-3 py-2 text-sm text-slate-600 dark:border-slate-800 dark:bg-slate-950 dark:text-slate-300">
                {{ $t('catalog.torrents.edit.noChanges') }}
              </div>

              <UButton type="submit" color="primary" icon="i-lucide-save" block :loading="savePending" :disabled="savePending || !torrent || !canEditTorrent || !hasChanges">
                {{ $t('catalog.torrents.edit.submit') }}
              </UButton>
            </div>
          </UCard>
        </aside>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ApiError } from '~/composables/useApi'
import type { CatalogCategory, CatalogTagGroup, ReleaseFieldsState, ReleaseFieldValue, TorrentDetail } from '~/composables/useCatalogTorrents'
import { localizeI18nName } from '~/utils/format'

definePageMeta({
  middleware: 'auth'
})

const { t, locale } = useI18n()
const localePath = useLocalePath()
const route = useRoute()
const toast = useToast()
const catalogTorrents = useCatalogTorrents()
const { user, hasPermission, fetchUser } = useAuth()

const torrent = ref<TorrentDetail | null>(null)
const categories = ref<CatalogCategory[]>([])
const tagGroups = ref<CatalogTagGroup[]>([])
const pending = ref(true)
const savePending = ref(false)
const tagGroupsPending = ref(false)
const errorMessage = ref('')
const submitAttempted = ref(false)

const form = reactive({
  categoryId: '0',
  name: '',
  subTitle: '',
  description: '',
  anonymous: false
})
const releaseFields = ref<Record<string, ReleaseFieldValue>>({})
const releaseState = ref<ReleaseFieldsState>({
  generatedTitle: '',
  valid: true,
  firstError: '',
  missingLabels: [],
  output: {}
})
const initialEditSnapshot = ref('')
const initialSnapshotPending = ref(false)
let initialSnapshotTimer: ReturnType<typeof setTimeout> | null = null
const initialSnapshotSettleMs = 30

const torrentId = computed(() => readRouteId())
const canEditTorrent = computed(() => Boolean(torrent.value && (hasPermission(Permission.AdminCatalogTorrentManage) || user.value?.user.id === torrent.value.owner?.id)))
const selectedCategory = computed(() => categories.value.find((item) => item.id === Number(form.categoryId)) || null)
const selectedCategoryName = computed(() => selectedCategory.value ? categoryDisplayName(selectedCategory.value) : '-')
const selectedUploadConfig = computed(() => selectedCategory.value?.uploadConfig || null)
const isGeneratedTitleMode = computed(() => selectedUploadConfig.value?.title?.mode === 'generated')
const showManualTitleInput = computed(() => !isGeneratedTitleMode.value || Boolean(selectedUploadConfig.value?.title?.allowManualOverride))
const manualTitleLabel = computed(() => isGeneratedTitleMode.value ? t('catalog.torrents.upload.fields.titleOverride') : t('catalog.torrents.upload.fields.name'))
const generatedTitle = computed(() => releaseState.value.generatedTitle)
const releaseFieldsValid = computed(() => releaseState.value.valid)
const effectiveTitle = computed(() => {
  const manualTitle = form.name.trim()
  if (showManualTitleInput.value && manualTitle) return manualTitle
  if (isGeneratedTitleMode.value) {
    if (generatedTitle.value) return generatedTitle.value
  }
  return manualTitle
})
const finalTitle = computed(() => effectiveTitle.value)
const titleError = computed(() => finalTitle.value ? '' : t('catalog.torrents.upload.errors.titleRequired'))
const submitError = computed(() => {
  if (Number(form.categoryId) <= 0) return t('catalog.torrents.upload.errors.categoryRequired')
  if (tagGroupsPending.value) return t('catalog.torrents.upload.errors.optionsLoading')
  if (!releaseFieldsValid.value) return releaseState.value.firstError
  if (titleError.value) return titleError.value
  return ''
})
const currentEditSnapshot = computed(() => editSnapshot())
const hasChanges = computed(() => Boolean(!initialSnapshotPending.value && initialEditSnapshot.value && currentEditSnapshot.value !== initialEditSnapshot.value))
const canSubmit = computed(() => Boolean(hasChanges.value && !submitError.value && torrent.value && canEditTorrent.value && !savePending.value))
const finalTitleCardClass = computed(() => [
  'rounded-md border bg-white px-3 py-2.5 dark:bg-slate-950',
  submitAttempted.value && titleError.value
    ? 'border-red-200 dark:border-red-900'
    : 'border-slate-200 dark:border-slate-800'
].join(' '))

useHead(() => ({
  title: t('catalog.torrents.edit.metaTitle')
}))

onMounted(loadPage)

onBeforeUnmount(clearInitialSnapshotTimer)

watch(currentEditSnapshot, () => {
  if (initialSnapshotPending.value) {
    scheduleInitialSnapshotCapture()
  }
})

watch(selectedCategory, (category) => {
  if (categoryNeedsTagGroups(category)) {
    void loadTagGroups()
  }
  if (isGeneratedTitleMode.value && !selectedUploadConfig.value?.title?.allowManualOverride) {
    form.name = ''
  }
})

function readRouteId() {
  const raw = Array.isArray(route.params.id) ? route.params.id[0] : route.params.id
  const id = Number(raw)
  return Number.isInteger(id) && id > 0 ? id : 0
}

async function loadPage() {
  if (torrentId.value <= 0) {
    errorMessage.value = t('catalog.torrents.detail.invalidId')
    pending.value = false
    return
  }

  pending.value = true
  errorMessage.value = ''
  initialEditSnapshot.value = ''
  initialSnapshotPending.value = true
  try {
    if (!user.value) {
      await fetchUser()
    }

    const [detail, categoryList] = await Promise.all([
      catalogTorrents.getTorrent(torrentId.value),
      catalogTorrents.listCategories().catch(() => ({ list: [] }))
    ])

    torrent.value = detail
    categories.value = categoryList.list || []
    form.categoryId = String(detail.categoryId || 0)
    form.name = detail.name || ''
    form.subTitle = detail.subTitle || ''
    form.description = detail.description || ''
    form.anonymous = Boolean(detail.anonymous)
    releaseFields.value = normalizeStoredReleaseFields(detail.releaseFields)
    if (categoryNeedsTagGroups(selectedCategory.value)) {
      await loadTagGroups()
    }

    if (!canEditTorrent.value) {
      errorMessage.value = t('catalog.torrents.edit.forbidden')
    }
    initialSnapshotPending.value = true
    scheduleInitialSnapshotCapture()
  } catch (error) {
    torrent.value = null
    categories.value = []
    initialSnapshotPending.value = false
    clearInitialSnapshotTimer()
    errorMessage.value = error instanceof ApiError ? error.message : t('common.requestFailed')
  } finally {
    pending.value = false
  }
}

async function handleSubmit() {
  submitAttempted.value = true
  if (!hasChanges.value) {
    toast.add({
      title: t('catalog.torrents.edit.noChanges'),
      color: 'neutral',
      icon: 'i-lucide-circle-check'
    })
    return
  }
  if (!canSubmit.value) {
    toast.add({
      title: submitError.value || t('common.requestFailed'),
      color: 'error',
      icon: 'i-lucide-circle-alert'
    })
    scrollToEditIssue()
    return
  }

  savePending.value = true
  try {
    await catalogTorrents.updateTorrent(torrentId.value, {
      categoryId: Number(form.categoryId),
      name: submitTitleValue(),
      subTitle: form.subTitle.trim(),
      description: form.description.trim(),
      releaseFields: releaseState.value.output,
      anonymous: form.anonymous
    })
    toast.add({
      title: t('catalog.torrents.edit.saved'),
      color: 'success',
      icon: 'i-lucide-check-circle'
    })
    await navigateTo(localePath(`/catalog/torrents/${torrentId.value}`))
  } catch (error) {
    toast.add({
      title: error instanceof ApiError ? error.message : t('common.requestFailed'),
      color: 'error',
      icon: 'i-lucide-circle-alert'
    })
  } finally {
    savePending.value = false
  }
}

function categoryDisplayName(category: CatalogCategory) {
  return localizeI18nName(category.name, locale.value, category.slug || `#${category.id}`)
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

function categoryNeedsTagGroups(category: CatalogCategory | null) {
  return Boolean(category?.uploadConfig?.fields?.some((field) => field.options?.source === 'tagGroup'))
}

function submitTitleValue() {
  if (showManualTitleInput.value && form.name.trim()) {
    return form.name.trim()
  }
  if (isGeneratedTitleMode.value) {
    return ''
  }
  return form.name.trim()
}

function normalizeStoredReleaseFields(source?: Record<string, unknown> | null) {
  const fields: Record<string, ReleaseFieldValue> = {}
  if (!source) return fields
  for (const [key, value] of Object.entries(source)) {
    if (Array.isArray(value)) {
      fields[key] = value.map((item) => String(item)).filter(Boolean)
      continue
    }
    fields[key] = String(value || '')
  }
  return fields
}

function handleReleaseStateChange(state: ReleaseFieldsState) {
  releaseState.value = state
  if (initialSnapshotPending.value) {
    scheduleInitialSnapshotCapture()
  }
}

function editSnapshot() {
  return JSON.stringify({
    categoryId: Number(form.categoryId),
    name: submitTitleValue(),
    subTitle: form.subTitle.trim(),
    description: form.description.trim(),
    anonymous: form.anonymous,
    releaseFields: releaseState.value.output
  })
}

function scheduleInitialSnapshotCapture() {
  clearInitialSnapshotTimer()
  initialSnapshotTimer = setTimeout(() => {
    initialSnapshotTimer = null
    if (!initialSnapshotPending.value) return
    initialEditSnapshot.value = currentEditSnapshot.value
    initialSnapshotPending.value = false
    submitAttempted.value = false
  }, initialSnapshotSettleMs)
}

function clearInitialSnapshotTimer() {
  if (!initialSnapshotTimer) return
  clearTimeout(initialSnapshotTimer)
  initialSnapshotTimer = null
}

function scrollToEditIssue() {
  let target = '#edit-title'
  if (Number(form.categoryId) <= 0) {
    target = '#edit-category'
  } else if (!releaseFieldsValid.value) {
    target = '#upload-release-fields'
  }
  document.querySelector(target)?.scrollIntoView({ behavior: 'smooth', block: 'center' })
}
</script>
