<template>
  <section class="overflow-hidden rounded-lg border border-slate-200 bg-white dark:border-slate-800 dark:bg-slate-900">
    <form class="space-y-4 p-4" @submit.prevent="$emit('submit')">
      <div class="flex flex-col gap-2 sm:flex-row sm:items-center sm:justify-between">
        <div>
          <h2 class="text-base font-semibold text-slate-950 dark:text-white">{{ t('forum.detail.edit.title') }}</h2>
          <p class="mt-1 text-sm text-slate-500 dark:text-slate-400">{{ t('forum.detail.edit.hint') }}</p>
        </div>
        <UButton type="button" color="neutral" variant="ghost" size="sm" @click="$emit('cancel')">
          {{ t('common.cancel') }}
        </UButton>
      </div>

      <div v-if="nodesPending" class="grid gap-3 md:grid-cols-[minmax(0,240px)_minmax(0,1fr)]">
        <div class="h-10 animate-pulse rounded-md bg-slate-100 dark:bg-slate-800" />
        <div class="h-10 animate-pulse rounded-md bg-slate-100 dark:bg-slate-800" />
      </div>
      <div v-else class="grid gap-3 md:grid-cols-[minmax(0,240px)_minmax(0,1fr)]">
        <UFormField :label="t('forum.create.fields.category')">
          <select
            v-model.number="selectedCategoryIdValue"
            class="h-10 w-full rounded-md border border-slate-200 bg-white px-3 text-sm text-slate-950 outline-none transition focus:border-sky-400 focus:ring-2 focus:ring-sky-100 dark:border-slate-700 dark:bg-slate-950 dark:text-white dark:focus:border-sky-500 dark:focus:ring-sky-950"
            :disabled="isUpdating || categories.length === 0"
            @change="handleCategoryChange"
          >
            <option v-for="category in categories" :key="category.id" :value="category.id">
              {{ categoryDisplayName(category) }}
            </option>
          </select>
        </UFormField>

        <UFormField :label="t('forum.create.fields.node')" required>
          <select
            v-model="nodeIdValue"
            class="h-10 w-full rounded-md border border-slate-200 bg-white px-3 text-sm text-slate-950 outline-none transition focus:border-sky-400 focus:ring-2 focus:ring-sky-100 dark:border-slate-700 dark:bg-slate-950 dark:text-white dark:focus:border-sky-500 dark:focus:ring-sky-950"
            :disabled="isUpdating || selectedCategoryNodes.length === 0"
          >
            <option value="0">{{ t('forum.create.fields.nodePlaceholder') }}</option>
            <option v-for="node in selectedCategoryNodes" :key="node.id" :value="String(node.id)">
              {{ nodeDisplayName(node) }}
            </option>
          </select>
        </UFormField>
      </div>

      <p v-if="nodesError" class="rounded-md border border-red-200 bg-red-50 px-3 py-2 text-sm text-red-700 dark:border-red-900 dark:bg-red-950 dark:text-red-200">
        {{ nodesError }}
      </p>

      <UFormField :label="t('forum.create.fields.subject')" required>
        <UInput v-model="subjectValue" class="w-full" :disabled="isUpdating" />
      </UFormField>

      <UFormField :label="t('forum.create.fields.content')" required>
        <RichTextComposer
          v-model="contentValue"
          v-model:mode="editorModeValue"
          :as-form="false"
          :rows="10"
          :disabled="isUpdating"
          :submit-label="t('common.save')"
          :write-label="t('common.editor.edit')"
          :preview-label="t('common.editor.preview')"
          :preview-empty="t('common.editor.previewEmpty')"
        >
          <template #actions>
            <span />
          </template>
        </RichTextComposer>
      </UFormField>

      <div class="flex justify-end gap-2 border-t border-slate-200 pt-4 dark:border-slate-800">
        <UButton type="button" color="neutral" variant="ghost" @click="$emit('cancel')">
          {{ t('common.cancel') }}
        </UButton>
        <UButton type="submit" color="primary" icon="i-lucide-save" :loading="isUpdating" :disabled="!canSubmit">
          {{ t('common.save') }}
        </UButton>
      </div>
    </form>
  </section>
</template>

<script setup lang="ts">
import type { ForumNode, ForumNodeCategory } from '~/composables/useForum'
import { localizeI18nName } from '~/utils/format'

type EditorMode = 'write' | 'preview'

const props = withDefaults(defineProps<{
  categories: ForumNodeCategory[]
  nodesPending?: boolean
  nodesError?: string
  actionPending?: string
  canSubmit?: boolean
  selectedCategoryId: number
  nodeId: string
  subject: string
  content: string
  editorMode: EditorMode
}>(), {
  nodesPending: false,
  nodesError: '',
  actionPending: '',
  canSubmit: false
})

const emit = defineEmits<{
  submit: []
  cancel: []
  'update:selectedCategoryId': [value: number]
  'update:nodeId': [value: string]
  'update:subject': [value: string]
  'update:content': [value: string]
  'update:editorMode': [value: EditorMode]
}>()

const { t, locale } = useI18n()

const isUpdating = computed(() => props.actionPending === 'update')
const selectedCategoryIdValue = computed({
  get: () => props.selectedCategoryId,
  set: (value) => emit('update:selectedCategoryId', value)
})
const nodeIdValue = computed({
  get: () => props.nodeId,
  set: (value) => emit('update:nodeId', value)
})
const subjectValue = computed({
  get: () => props.subject,
  set: (value) => emit('update:subject', value)
})
const contentValue = computed({
  get: () => props.content,
  set: (value) => emit('update:content', value)
})
const editorModeValue = computed({
  get: () => props.editorMode,
  set: (value) => emit('update:editorMode', value)
})
const selectedCategory = computed(() => {
  return props.categories.find((category) => category.id === selectedCategoryIdValue.value) || props.categories[0] || null
})
const selectedCategoryNodes = computed(() => selectedCategory.value?.nodes || [])

function handleCategoryChange() {
  const category = props.categories.find((item) => item.id === selectedCategoryIdValue.value)
  if (!category) {
    nodeIdValue.value = '0'
    return
  }

  const currentNodeInCategory = category.nodes.some((node) => node.id === Number(nodeIdValue.value))
  if (!currentNodeInCategory) {
    nodeIdValue.value = String(category.nodes[0]?.id || 0)
  }
}

function categoryDisplayName(category: ForumNodeCategory) {
  return localizeI18nName(category.nameI18n, locale.value, t('forum.fallback.category'))
}

function nodeDisplayName(node: ForumNode) {
  return localizeI18nName(node.nameI18n, locale.value, node.slug || `#${node.id}`)
}
</script>
