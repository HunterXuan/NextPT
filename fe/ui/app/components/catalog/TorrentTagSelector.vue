<template>
  <section v-if="visibleGroups.length" id="torrent-tags" class="rounded-md border border-slate-200 bg-slate-50/70 px-3 py-3 dark:border-slate-800 dark:bg-slate-950/50">
    <div class="space-y-3">
      <div v-for="group in visibleGroups" :key="group.id" class="grid gap-2 sm:grid-cols-[96px_minmax(0,1fr)] sm:items-start">
        <span class="pt-1 text-sm font-medium text-slate-700 dark:text-slate-200">{{ groupName(group) }}</span>
        <div class="flex min-w-0 flex-wrap gap-1.5">
          <UButton
            v-for="tag in group.tags"
            :key="tag.id"
            type="button"
            color="neutral"
            size="xs"
            :variant="modelValue.includes(tag.id) ? 'soft' : 'outline'"
            :class="modelValue.includes(tag.id) ? 'ring-1 ring-sky-300 dark:ring-sky-700' : ''"
            :disabled="disabled"
            @click="toggleTag(tag.id)"
          >
            {{ tagName(tag) }}
          </UButton>
        </div>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
import type { CatalogCategory, CatalogTagGroup, CatalogTagItem } from '~/composables/useCatalogTorrents'
import { localizeI18nName } from '~/utils/format'

const props = withDefaults(defineProps<{
  modelValue: number[]
  category: CatalogCategory | null
  groups: CatalogTagGroup[]
  disabled?: boolean
}>(), {
  disabled: false
})

const emit = defineEmits<{
  'update:modelValue': [value: number[]]
}>()

const { locale } = useI18n()

const releaseGroupSlugs = computed(() => new Set(
  (props.category?.uploadConfig?.fields || [])
    .filter(field => field.options?.source === 'tagGroup' && field.options.slug)
    .map(field => String(field.options?.slug))
))

const visibleGroups = computed(() => props.groups.filter((group) => {
  if (releaseGroupSlugs.value.has(group.slug)) return false
  return !group.categories?.length || group.categories.includes(props.category?.id || 0)
}))

watch(() => [props.category?.id, props.groups] as const, () => {
  if (!props.groups.length) return
  const available = new Set(visibleGroups.value.flatMap(group => group.tags.map(tag => tag.id)))
  const next = props.modelValue.filter(id => available.has(id))
  if (next.length !== props.modelValue.length) emit('update:modelValue', next)
}, { deep: true })

function toggleTag(tagId: number) {
  const next = props.modelValue.includes(tagId)
    ? props.modelValue.filter(id => id !== tagId)
    : [...props.modelValue, tagId]
  emit('update:modelValue', next)
}

function groupName(group: CatalogTagGroup) {
  return localizeI18nName(group.name, locale.value, group.slug)
}

function tagName(tag: CatalogTagItem) {
  return localizeI18nName(tag.name, locale.value, tag.value)
}
</script>
