<template>
  <section class="space-y-4 rounded-md border border-slate-200 bg-slate-50/70 p-3 dark:border-slate-800 dark:bg-slate-950/40">
    <div class="min-w-0">
      <h3 class="text-sm font-semibold text-slate-950 dark:text-white">{{ $t('catalog.torrents.metadata.title') }}</h3>
      <p class="mt-1 text-xs leading-5 text-slate-500 dark:text-slate-400">{{ $t('catalog.torrents.metadata.description') }}</p>
    </div>

    <div class="grid grid-cols-1 gap-3 sm:grid-cols-[minmax(0,1fr)_140px_40px] sm:items-end">
      <UFormField :label="$t('catalog.torrents.metadata.searchLabel')">
        <UInput
          v-model="query"
          class="w-full"
          size="lg"
          :ui="{ base: 'h-10 w-full' }"
          :placeholder="$t('catalog.torrents.metadata.searchPlaceholder')"
          :disabled="disabled || searchPending"
          @keyup.enter="search"
        />
      </UFormField>
      <UFormField :label="$t('catalog.torrents.metadata.typeLabel')">
        <USelect
          v-model="tmdbType"
          class="w-full"
          size="lg"
          :ui="{ base: 'h-10 w-full' }"
          :items="typeOptions"
          value-key="value"
          :disabled="disabled || searchPending"
        />
      </UFormField>
      <UTooltip :text="$t('catalog.torrents.metadata.search')">
        <span class="inline-flex size-10">
          <UButton
            type="button"
            color="neutral"
            variant="outline"
            icon="i-lucide-search"
            class="size-10 justify-center p-0"
            :aria-label="$t('catalog.torrents.metadata.search')"
            :loading="searchPending"
            :disabled="disabled || !query.trim()"
            @click="search"
          />
        </span>
      </UTooltip>
    </div>

    <div v-if="searchError" class="rounded-md border border-red-200 bg-red-50 px-3 py-2 text-sm text-red-700 dark:border-red-900 dark:bg-red-950 dark:text-red-200">
      {{ searchError }}
    </div>

    <div v-if="results.length > 0" class="divide-y divide-slate-200 overflow-hidden rounded-md border border-slate-200 bg-white dark:divide-slate-800 dark:border-slate-800 dark:bg-slate-900">
      <button
        v-for="item in results"
        :key="`${item.tmdbType}-${item.providerId}`"
        type="button"
        class="flex w-full min-w-0 gap-3 px-3 py-2.5 text-left transition hover:bg-slate-50 dark:hover:bg-slate-800/70"
        :disabled="disabled"
        @click="selectResult(item)"
      >
        <img v-if="item.posterUrl" :src="item.posterUrl" :alt="item.title" class="size-12 shrink-0 rounded object-cover" loading="lazy" referrerpolicy="no-referrer">
        <span v-else class="flex size-12 shrink-0 items-center justify-center rounded bg-slate-100 text-slate-400 dark:bg-slate-800">
          <UIcon name="i-lucide-image-off" class="size-4" />
        </span>
        <span class="min-w-0 flex-1">
          <span class="flex min-w-0 items-center gap-2">
            <span class="truncate text-sm font-semibold text-slate-950 dark:text-white">{{ item.title || item.originalTitle }}</span>
            <span v-if="item.year" class="shrink-0 text-xs text-slate-500 dark:text-slate-400">{{ item.year }}</span>
          </span>
          <span class="mt-1 flex items-center gap-2 text-xs text-slate-500 dark:text-slate-400">
            <span v-if="item.originalTitle" class="truncate">{{ item.originalTitle }}</span>
            <span v-if="item.rating > 0" class="shrink-0 text-amber-600 dark:text-amber-400">★ {{ item.rating.toFixed(1) }}</span>
          </span>
        </span>
        <UIcon name="i-lucide-arrow-up-right" class="mt-1 size-4 shrink-0 text-slate-400" />
      </button>
    </div>

    <div v-if="binding.tmdbId" class="flex w-full items-start gap-3 rounded-md border border-slate-200 bg-white p-3 dark:border-slate-700 dark:bg-slate-900 sm:gap-4">
      <img v-if="boundMetadata?.posterUrl" :src="boundMetadata.posterUrl" :alt="boundMetadata.title" class="h-20 w-14 shrink-0 rounded object-cover" loading="lazy" referrerpolicy="no-referrer">
      <span v-else class="flex h-20 w-14 shrink-0 items-center justify-center rounded bg-slate-100 text-slate-400 dark:bg-slate-800">
        <UIcon name="i-lucide-image-off" class="size-4" />
      </span>
      <div class="min-w-0 flex-1 py-0.5">
        <div class="flex min-w-0 flex-wrap items-baseline gap-x-2 gap-y-0.5">
          <p class="truncate text-sm font-semibold text-slate-950 dark:text-white">{{ boundMetadata?.title || `TMDB #${binding.tmdbId}` }}</p>
          <span v-if="boundMetadata?.year" class="shrink-0 text-xs text-slate-500 dark:text-slate-400">{{ boundMetadata.year }}</span>
        </div>
        <p v-if="boundMetadata?.originalTitle && boundMetadata.originalTitle !== boundMetadata.title" class="mt-1 truncate text-xs text-slate-500 dark:text-slate-400">
          {{ boundMetadata.originalTitle }}
        </p>
        <p v-if="boundMetadata?.overview" class="mt-1.5 line-clamp-1 text-xs leading-5 text-slate-600 dark:text-slate-300">
          {{ boundMetadata.overview }}
        </p>
        <div class="mt-2 flex flex-wrap items-center gap-x-3 gap-y-1 text-xs text-slate-500 dark:text-slate-400">
          <span class="font-medium text-slate-700 dark:text-slate-200">TMDB #{{ binding.tmdbId }}</span>
          <span class="inline-flex items-center gap-1"><UIcon name="i-lucide-film" class="size-3.5" />{{ binding.tmdbType === 'tv' ? $t('catalog.torrents.metadata.tv') : $t('catalog.torrents.metadata.movie') }}</span>
          <span v-if="boundMetadata?.rating" class="inline-flex items-center gap-1 text-amber-600 dark:text-amber-400"><UIcon name="i-lucide-star" class="size-3.5" />{{ boundMetadata.rating.toFixed(1) }}</span>
        </div>
      </div>
      <UTooltip :text="$t('catalog.torrents.metadata.removeBinding')">
        <UButton type="button" color="neutral" variant="ghost" size="xs" icon="i-lucide-x" class="size-8 shrink-0 justify-center p-0" :aria-label="$t('catalog.torrents.metadata.removeBinding')" :disabled="disabled" @click="clearTmdb" />
      </UTooltip>
    </div>

    <div class="grid grid-cols-1 gap-3 md:grid-cols-3">
      <UFormField :label="$t('catalog.torrents.metadata.imdb')">
        <UInput :model-value="binding.imdbId" class="w-full" size="lg" :ui="{ base: 'h-10 w-full' }" placeholder="tt..." :disabled="disabled" @update:model-value="updateField('imdbId', $event)" />
      </UFormField>
      <UFormField :label="$t('catalog.torrents.metadata.douban')">
        <UInput :model-value="binding.doubanId" class="w-full" size="lg" :ui="{ base: 'h-10 w-full' }" :disabled="disabled" @update:model-value="updateField('doubanId', $event)" />
      </UFormField>
      <UFormField :label="$t('catalog.torrents.metadata.bangumi')">
        <UInput :model-value="binding.bangumiId" class="w-full" size="lg" :ui="{ base: 'h-10 w-full' }" :disabled="disabled" @update:model-value="updateField('bangumiId', $event)" />
      </UFormField>
    </div>
  </section>
</template>

<script setup lang="ts">
import { ApiError } from '~/composables/useApi'
import type { TorrentMetadataBinding, TorrentMetadataItem } from '~/composables/useCatalogTorrents'

const props = withDefaults(defineProps<{
  modelValue?: TorrentMetadataBinding | null
  initialData?: TorrentMetadataItem | null
  disabled?: boolean
}>(), {
  modelValue: null,
  initialData: null,
  disabled: false
})

const emit = defineEmits<{
  'update:modelValue': [value: TorrentMetadataBinding]
}>()

const { t } = useI18n()
const catalogTorrents = useCatalogTorrents()
const query = ref('')
const tmdbType = ref<'movie' | 'tv'>('movie')
const results = ref<TorrentMetadataItem[]>([])
const selectedResult = ref<TorrentMetadataItem | null>(null)
const searchPending = ref(false)
const searchError = ref('')

const binding = computed<TorrentMetadataBinding>(() => ({
  imdbId: props.modelValue?.imdbId || '',
  doubanId: props.modelValue?.doubanId || '',
  bangumiId: props.modelValue?.bangumiId || '',
  tmdbId: props.modelValue?.tmdbId || '',
  tmdbType: props.modelValue?.tmdbType || '',
  imdbRating: props.modelValue?.imdbRating,
  doubanRating: props.modelValue?.doubanRating,
  bangumiRating: props.modelValue?.bangumiRating,
  tmdbRating: props.modelValue?.tmdbRating
}))
const boundMetadata = computed(() => {
  if (selectedResult.value?.providerId === binding.value.tmdbId && selectedResult.value.tmdbType === binding.value.tmdbType) return selectedResult.value
  if (props.initialData?.providerId === binding.value.tmdbId && props.initialData.tmdbType === binding.value.tmdbType) return props.initialData
  return null
})

const typeOptions = computed(() => [
  { value: 'movie', label: t('catalog.torrents.metadata.movie') },
  { value: 'tv', label: t('catalog.torrents.metadata.tv') }
])

async function search() {
  if (!query.value.trim() || searchPending.value) return
  searchPending.value = true
  searchError.value = ''
  try {
    const data = await catalogTorrents.searchMetadata(query.value.trim(), tmdbType.value)
    results.value = data.list || []
    if (results.value.length === 0) searchError.value = t('catalog.torrents.metadata.empty')
  } catch (error) {
    results.value = []
    searchError.value = error instanceof ApiError ? error.message : t('common.requestFailed')
  } finally {
    searchPending.value = false
  }
}

function selectResult(item: TorrentMetadataItem) {
  selectedResult.value = item
  emit('update:modelValue', {
    ...binding.value,
    tmdbId: item.providerId,
    tmdbType: item.tmdbType,
    imdbId: binding.value.imdbId || item.imdbId
  })
  results.value = []
}

function updateField(field: 'imdbId' | 'doubanId' | 'bangumiId', value: unknown) {
  emit('update:modelValue', { ...binding.value, [field]: String(value || '') })
}

function clearTmdb() {
  selectedResult.value = null
  emit('update:modelValue', { ...binding.value, tmdbId: '', tmdbType: '', tmdbRating: undefined })
}
</script>
