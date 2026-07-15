<template>
  <section
    v-if="schemaFields.length > 0"
    id="upload-release-fields"
    class="rounded-md border border-slate-200 bg-slate-50/70 p-3 dark:border-slate-800 dark:bg-slate-950/40"
  >
    <div class="mb-3 flex items-center justify-between gap-3">
      <h3 class="text-sm font-semibold text-slate-950 dark:text-white">{{ t('catalog.torrents.upload.sections.release') }}</h3>
      <UBadge color="neutral" variant="soft">{{ selectedCategoryName }}</UBadge>
    </div>

    <div class="grid grid-cols-1 gap-4">
      <UFormField
        v-for="field in schemaFields"
        :key="field.key"
        :label="fieldLabel(field)"
        :description="fieldDescription(field)"
        :required="Boolean(field.required)"
        :error="fieldError(field)"
      >
        <USelect
          v-if="field.type === 'select'"
          :model-value="String(modelValue[field.key] || '')"
          class="w-full"
          size="lg"
          :ui="{ base: 'h-10 w-full' }"
          :items="releaseSelectOptions(field)"
          value-key="value"
          :disabled="disabled || optionsPending"
          @update:model-value="setReleaseField(field.key, String($event || ''))"
        />

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
              :disabled="disabled || optionsPending"
              @change="toggleReleaseFieldOption(field.key, option.value)"
            >
            <span class="min-w-0 truncate">{{ optionLabel(option) }}</span>
          </label>
          <p v-if="fieldOptions(field).length === 0" class="px-2 py-1 text-sm text-slate-500 dark:text-slate-400">
            {{ t('catalog.torrents.upload.fields.noOptions') }}
          </p>
        </div>

        <UTextarea
          v-else-if="field.type === 'textarea'"
          :model-value="String(modelValue[field.key] || '')"
          class="w-full"
          :rows="3"
          :disabled="disabled"
          :placeholder="fieldPlaceholder(field)"
          @update:model-value="setReleaseField(field.key, String($event || ''))"
        />

        <UInput
          v-else
          :model-value="String(modelValue[field.key] || '')"
          class="w-full"
          :disabled="disabled"
          :placeholder="fieldPlaceholder(field)"
          @update:model-value="setReleaseField(field.key, String($event || ''))"
        />
      </UFormField>
    </div>
  </section>
</template>

<script setup lang="ts">
import type { CatalogCategory, CatalogTagGroup, ReleaseFieldsState, ReleaseFieldValue, UploadFieldConfig, UploadOptionItem, UploadTitlePart } from '~/composables/useCatalogTorrents'
import { localizeI18nName } from '~/utils/format'

const props = withDefaults(defineProps<{
  category: CatalogCategory | null
  tagGroups: CatalogTagGroup[]
  modelValue: Record<string, ReleaseFieldValue>
  disabled?: boolean
  optionsPending?: boolean
  showErrors?: boolean
}>(), {
  disabled: false,
  optionsPending: false,
  showErrors: false
})

const emit = defineEmits<{
  'update:modelValue': [value: Record<string, ReleaseFieldValue>]
  'state-change': [state: ReleaseFieldsState]
}>()

const { t, locale } = useI18n()

const selectedCategoryName = computed(() => props.category ? localizeI18nName(props.category.name, locale.value, props.category.slug || `#${props.category.id}`) : '-')
const schemaFields = computed(() => props.category?.uploadConfig?.fields || [])
const generatedTitle = computed(() => buildGeneratedTitle(props.category?.uploadConfig?.title?.parts || [], props.modelValue))
const missingFields = computed(() => schemaFields.value.filter((field) => isFieldMissing(field, props.modelValue)))
const releaseState = computed<ReleaseFieldsState>(() => {
  const missingLabels = missingFields.value.map(fieldLabel)
  return {
    generatedTitle: generatedTitle.value,
    valid: missingLabels.length === 0,
    firstError: missingLabels.length > 0 ? t('catalog.torrents.upload.errors.releaseFieldRequired', { field: missingLabels[0] }) : '',
    missingLabels,
    output: buildReleaseFieldsInput(props.modelValue)
  }
})

watch(
  () => [schemaFields.value, props.modelValue] as const,
  () => {
    const normalized = normalizeReleaseFields(props.modelValue)
    if (!releaseFieldsEqual(normalized, props.modelValue)) {
      emit('update:modelValue', normalized)
    }
  },
  { immediate: true, deep: true }
)

watch(releaseState, (state) => emit('state-change', state), { immediate: true })

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

  const group = props.tagGroups.find((item) => item.slug === field.options?.slug && tagGroupAppliesToCategory(item))
  return (group?.tags || [])
    .filter((tag) => Boolean(tag.value))
    .map((tag) => ({ value: tag.value, label: tag.name }))
}

function releaseSelectOptions(field: UploadFieldConfig) {
  return [
    { value: '', label: t('catalog.torrents.upload.fields.optionPlaceholder') },
    ...fieldOptions(field).map(option => ({ value: option.value, label: optionLabel(option) }))
  ]
}

function tagGroupAppliesToCategory(group: CatalogTagGroup) {
  if (!props.category || !group.categories?.length) return true
  return group.categories.includes(props.category.id)
}

function fieldError(field: UploadFieldConfig) {
  if (!props.showErrors || !isFieldMissing(field, props.modelValue)) return undefined
  return t('catalog.torrents.upload.errors.releaseFieldRequired', { field: fieldLabel(field) })
}

function setReleaseField(key: string, value: string) {
  emitNormalizedFields({
    ...props.modelValue,
    [key]: value
  })
}

function releaseFieldListValue(key: string) {
  const value = props.modelValue[key]
  return Array.isArray(value) ? value : []
}

function toggleReleaseFieldOption(key: string, value: string) {
  const current = releaseFieldListValue(key)
  emitNormalizedFields({
    ...props.modelValue,
    [key]: current.includes(value)
      ? current.filter((item) => item !== value)
      : [...current, value]
  })
}

function emitNormalizedFields(value: Record<string, ReleaseFieldValue>) {
  const normalized = normalizeReleaseFields(value)
  emit('update:modelValue', normalized)
}

function normalizeReleaseFields(source: Record<string, unknown>) {
  const fields: Record<string, ReleaseFieldValue> = {}
  for (const field of schemaFields.value) {
    if (!field.key) continue
    fields[field.key] = normalizeReleaseFieldValue(field, source[field.key])
  }
  return fields
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

function isFieldMissing(field: UploadFieldConfig, source: Record<string, ReleaseFieldValue>) {
  if (!field.required) return false
  const value = source[field.key]
  return Array.isArray(value) ? value.length === 0 : !String(value || '').trim()
}

function buildGeneratedTitle(parts: UploadTitlePart[], source: Record<string, ReleaseFieldValue>) {
  return parts
    .map((part) => {
      const value = titlePartValue(source[part.field], part.separator)
      return value ? `${part.prefix || ''}${value}${part.suffix || ''}` : ''
    })
    .join('')
    .trim()
}

function titlePartValue(value: ReleaseFieldValue | undefined, separator = '/') {
  if (Array.isArray(value)) return value.join(separator || '/')
  return String(value || '').trim()
}

function buildReleaseFieldsInput(source: Record<string, ReleaseFieldValue>) {
  const input: Record<string, unknown> = {}
  for (const field of schemaFields.value) {
    const value = source[field.key]
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

function releaseFieldsEqual(left: Record<string, ReleaseFieldValue>, right: Record<string, ReleaseFieldValue>) {
  return JSON.stringify(left) === JSON.stringify(right || {})
}
</script>
