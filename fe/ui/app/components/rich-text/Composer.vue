<template>
  <component :is="asForm ? 'form' : 'div'" class="space-y-3" @submit.prevent="handleSubmit">
    <UTextarea
      v-if="modeValue === 'write'"
      v-model="contentValue"
      class="w-full"
      :rows="rows"
      :placeholder="placeholder"
      :disabled="disabled"
    />
    <div v-else class="min-h-28 rounded-md border border-slate-200 bg-slate-50 px-3 py-2.5 dark:border-slate-800 dark:bg-slate-950">
      <div v-if="renderedPreview" class="rich-text rich-text-compact" v-html="renderedPreview" />
      <p v-else class="text-sm text-slate-500 dark:text-slate-400">{{ previewEmpty }}</p>
    </div>

    <div class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
      <div class="inline-flex w-fit rounded-md border border-slate-200 bg-slate-50 p-0.5 dark:border-slate-800 dark:bg-slate-950">
        <button type="button" :class="tabClass('write')" @click="modeValue = 'write'">
          {{ writeLabel }}
        </button>
        <button type="button" :class="tabClass('preview')" @click="modeValue = 'preview'">
          {{ previewLabel }}
        </button>
      </div>
      <div class="flex items-center justify-end gap-3">
        <span v-if="lockedText" class="text-xs text-slate-500 dark:text-slate-400">
          {{ lockedText }}
        </span>
        <slot name="actions">
          <UButton type="submit" color="primary" :icon="submitIcon" :loading="pending" :disabled="submitDisabled">
            {{ submitLabel }}
          </UButton>
        </slot>
      </div>
    </div>
  </component>
</template>

<script setup lang="ts">
import { renderUserMarkdown } from '~/utils/richText'

type EditorMode = 'write' | 'preview'

const props = withDefaults(defineProps<{
  modelValue: string
  mode: EditorMode
  rows?: number
  placeholder?: string
  disabled?: boolean
  pending?: boolean
  submitDisabled?: boolean
  asForm?: boolean
  submitLabel: string
  submitIcon?: string
  writeLabel: string
  previewLabel: string
  previewEmpty: string
  lockedText?: string
}>(), {
  rows: 4,
  placeholder: '',
  disabled: false,
  pending: false,
  submitDisabled: false,
  asForm: true,
  submitIcon: 'i-lucide-send',
  lockedText: ''
})

const emit = defineEmits<{
  'update:modelValue': [value: string]
  'update:mode': [value: EditorMode]
  submit: []
}>()

const contentValue = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value)
})

const modeValue = computed({
  get: () => props.mode,
  set: (value) => emit('update:mode', value)
})

const renderedPreview = computed(() => renderUserMarkdown(contentValue.value || '').trim())

function handleSubmit() {
  if (props.asForm) {
    emit('submit')
  }
}

function tabClass(mode: EditorMode) {
  const base = 'rounded px-3 py-1.5 text-xs font-medium transition-colors'
  if (modeValue.value === mode) {
    return `${base} bg-white text-slate-950 shadow-sm dark:bg-slate-800 dark:text-white`
  }
  return `${base} text-slate-500 hover:text-slate-950 dark:text-slate-400 dark:hover:text-white`
}
</script>
