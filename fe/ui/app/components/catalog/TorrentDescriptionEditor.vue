<template>
  <UFormField :label="t('catalog.torrents.upload.fields.description')">
    <UTextarea
      v-if="descriptionMode === 'write'"
      :model-value="modelValue"
      class="w-full"
      :rows="rows"
      :disabled="disabled"
      @update:model-value="emit('update:modelValue', String($event || ''))"
    />
    <div v-else class="min-h-72 rounded-md border border-slate-200 bg-slate-50 px-3 py-2.5 dark:border-slate-800 dark:bg-slate-950">
      <div v-if="renderedDescriptionPreview" class="rich-text" v-html="renderedDescriptionPreview" />
      <p v-else class="text-sm text-slate-500 dark:text-slate-400">{{ t('catalog.torrents.upload.preview.empty') }}</p>
    </div>

    <div class="mt-3 flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
      <div class="inline-flex w-fit rounded-md border border-slate-200 bg-slate-50 p-0.5 dark:border-slate-800 dark:bg-slate-950">
        <button
          type="button"
          :class="descriptionModeButtonClass('write')"
          @click="descriptionMode = 'write'"
        >
          {{ t('catalog.torrents.upload.preview.write') }}
        </button>
        <button
          type="button"
          :class="descriptionModeButtonClass('preview')"
          @click="descriptionMode = 'preview'"
        >
          {{ t('catalog.torrents.upload.preview.preview') }}
        </button>
      </div>
      <span class="text-xs text-slate-500 dark:text-slate-400">
        {{ t('catalog.torrents.upload.summary.descriptionLength', { count: numberFormatter.format(modelValue.trim().length) }) }}
      </span>
    </div>
  </UFormField>
</template>

<script setup lang="ts">
import { renderUserMarkdown } from '~/utils/richText'

type DescriptionMode = 'write' | 'preview'

const props = withDefaults(defineProps<{
  modelValue: string
  disabled?: boolean
  rows?: number
}>(), {
  disabled: false,
  rows: 12
})

const emit = defineEmits<{
  'update:modelValue': [value: string]
}>()

const { t, locale } = useI18n()
const descriptionMode = ref<DescriptionMode>('write')
const numberFormatter = computed(() => new Intl.NumberFormat(locale.value))
const renderedDescriptionPreview = computed(() => renderUserMarkdown(props.modelValue).trim())

function descriptionModeButtonClass(mode: DescriptionMode) {
  const active = descriptionMode.value === mode
  return [
    'h-7 rounded px-3 text-xs font-medium transition-colors',
    active
      ? 'bg-white text-slate-950 shadow-sm ring-1 ring-slate-200 dark:bg-slate-800 dark:text-white dark:ring-slate-700'
      : 'text-slate-500 hover:text-slate-800 dark:text-slate-400 dark:hover:text-slate-100'
  ].join(' ')
}
</script>
