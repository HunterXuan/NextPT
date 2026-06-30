<template>
  <div class="min-h-[calc(100vh-4rem)] bg-slate-50 py-6 dark:bg-slate-950">
    <div class="w-full px-3 sm:px-4 lg:px-5">
      <div v-if="pending" class="h-96 animate-pulse rounded-lg border border-slate-200 bg-white dark:border-slate-800 dark:bg-slate-900" />

      <div v-else-if="errorMessage" class="flex flex-col items-center justify-center rounded-lg border border-slate-200 bg-white px-4 py-16 text-center dark:border-slate-800 dark:bg-slate-900">
        <UIcon name="i-lucide-circle-alert" class="size-9 text-red-500" />
        <p class="mt-3 text-sm font-medium text-slate-950 dark:text-white">{{ errorMessage }}</p>
        <UButton class="mt-5" color="neutral" variant="outline" icon="i-lucide-refresh-cw" @click="loadPage">
          {{ $t('common.retry') }}
        </UButton>
      </div>

      <form v-else class="grid grid-cols-1 gap-6 lg:grid-cols-[minmax(0,1fr)_340px]" @submit.prevent="handleSubmit">
        <div class="space-y-6">
          <UCard class="rounded-lg">
            <template #header>
              <h2 class="text-base font-semibold text-slate-950 dark:text-white">{{ $t('catalog.torrents.edit.sections.info') }}</h2>
            </template>

            <div class="grid grid-cols-1 gap-4">
              <UFormField :label="$t('catalog.torrents.upload.fields.category')" required>
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
                      :disabled="savePending || tagGroupsPending"
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
                          :disabled="savePending || tagGroupsPending"
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
                      :disabled="savePending"
                      :placeholder="fieldPlaceholder(field)"
                      @update:model-value="setReleaseField(field.key, String($event || ''))"
                    />

                    <UInput
                      v-else
                      :model-value="String(releaseFields[field.key] || '')"
                      class="w-full"
                      :disabled="savePending"
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
                <UInput v-model="form.name" class="w-full" :disabled="savePending" @update:model-value="titleManuallyEdited = true" />
              </UFormField>

              <UFormField :label="$t('catalog.torrents.upload.fields.subTitle')">
                <UInput v-model="form.subTitle" class="w-full" :disabled="savePending" />
              </UFormField>

              <UFormField :label="$t('catalog.torrents.upload.fields.description')">
                <UTextarea v-model="form.description" class="w-full" :rows="10" :disabled="savePending" />
              </UFormField>
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

              <UButton type="submit" color="primary" icon="i-lucide-save" block :loading="savePending" :disabled="!canSubmit">
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
import type { CatalogCategory, CatalogTagGroup, TorrentDetail, UploadFieldConfig, UploadOptionItem, UploadTitlePart } from '~/composables/useCatalogTorrents'
import { localizeI18nName } from '~/utils/format'

type ReleaseFieldValue = string | string[]

definePageMeta({
  middleware: 'auth'
})

const { t, locale } = useI18n()
const localePath = useLocalePath()
const route = useRoute()
const toast = useToast()
const catalogTorrents = useCatalogTorrents()
const { user, isStaff, fetchUser } = useAuth()

const torrent = ref<TorrentDetail | null>(null)
const categories = ref<CatalogCategory[]>([])
const tagGroups = ref<CatalogTagGroup[]>([])
const pending = ref(true)
const savePending = ref(false)
const tagGroupsPending = ref(false)
const errorMessage = ref('')
const titleManuallyEdited = ref(false)

const form = reactive({
  categoryId: '0',
  name: '',
  subTitle: '',
  description: '',
  anonymous: false
})
const releaseFields = reactive<Record<string, ReleaseFieldValue>>({})

const torrentId = computed(() => readRouteId())
const canEditTorrent = computed(() => Boolean(torrent.value && (isStaff.value || user.value?.id === torrent.value.ownerId)))
const selectedCategory = computed(() => categories.value.find((item) => item.id === Number(form.categoryId)) || null)
const selectedCategoryName = computed(() => selectedCategory.value ? categoryDisplayName(selectedCategory.value) : '-')
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
const canSubmit = computed(() => Number(form.categoryId) > 0 && Boolean(torrent.value) && canEditTorrent.value && releaseFieldsValid.value && !tagGroupsPending.value && !savePending.value)

useHead(() => ({
  title: `${t('catalog.torrents.edit.metaTitle')} - NextPT`
}))

onMounted(loadPage)

watch(selectedCategory, (category) => {
  initializeReleaseFields()
  if (categoryNeedsTagGroups(category)) {
    void loadTagGroups()
  }
  if (isGeneratedTitleMode.value && !selectedUploadConfig.value?.title?.allowManualOverride) {
    form.name = ''
    titleManuallyEdited.value = false
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
    titleManuallyEdited.value = false
    initializeReleaseFields(detail.releaseFields)
    if (categoryNeedsTagGroups(selectedCategory.value)) {
      await loadTagGroups()
    }

    if (!canEditTorrent.value) {
      errorMessage.value = t('catalog.torrents.edit.forbidden')
    }
  } catch (error) {
    torrent.value = null
    categories.value = []
    errorMessage.value = error instanceof ApiError ? error.message : t('common.requestFailed')
  } finally {
    pending.value = false
  }
}

async function handleSubmit() {
  if (!canSubmit.value) return

  savePending.value = true
  try {
    await catalogTorrents.updateTorrent(torrentId.value, {
      categoryId: Number(form.categoryId),
      name: submitTitleValue(),
      subTitle: form.subTitle.trim(),
      description: form.description.trim(),
      releaseFields: buildReleaseFieldsInput(),
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

function initializeReleaseFields(source?: Record<string, unknown> | null) {
  const nextKeys = new Set(schemaFields.value.map((field) => field.key).filter(Boolean))
  for (const key of Object.keys(releaseFields)) {
    if (!nextKeys.has(key)) delete releaseFields[key]
  }
  for (const field of schemaFields.value) {
    if (!field.key) continue
    const sourceValue = source && Object.prototype.hasOwnProperty.call(source, field.key) ? source[field.key] : undefined
    const currentValue = releaseFields[field.key]
    releaseFields[field.key] = normalizeReleaseFieldValue(field, sourceValue ?? currentValue)
  }
}

function normalizeReleaseFieldValue(field: UploadFieldConfig, value: unknown): ReleaseFieldValue {
  if (field.type === 'multiSelect') {
    if (Array.isArray(value)) return value.map((item) => String(item)).filter(Boolean)
    const text = String(value || '').trim()
    return text ? [text] : []
  }
  if (Array.isArray(value)) return value.map((item) => String(item)).filter(Boolean).join('/')
  return String(value || '')
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
  return form.name.trim()
}
</script>
