<template>
  <div>
    <div class="grid gap-3 sm:grid-cols-2 xl:grid-cols-4">
      <UFormField :label="$t('catalog.torrents.filters.promotion')">
        <USelect
          v-model="draft.promotion"
          class="w-full"
          size="lg"
          :ui="{ base: 'h-10 w-full' }"
          :items="promotionOptions"
          value-key="value"
          :disabled="pending"
        />
      </UFormField>

      <UFormField :label="$t('catalog.torrents.filters.seedStatus')">
        <USelect
          v-model="draft.seedStatus"
          class="w-full"
          size="lg"
          :ui="{ base: 'h-10 w-full' }"
          :items="seedStatusOptions"
          value-key="value"
          :disabled="pending"
        />
      </UFormField>

      <UFormField :label="$t('catalog.torrents.filters.publishedWithin')">
        <USelect
          v-model="draft.publishedWithin"
          class="w-full"
          size="lg"
          :ui="{ base: 'h-10 w-full' }"
          :items="publishedWithinOptions"
          value-key="value"
          :disabled="pending"
        />
      </UFormField>

      <UFormField :label="$t('catalog.torrents.filters.sort')">
        <USelect
          v-model="draft.sort"
          class="w-full"
          size="lg"
          :ui="{ base: 'h-10 w-full' }"
          :items="sortOptions"
          value-key="value"
          :disabled="pending"
        />
      </UFormField>

      <div class="sm:col-span-2 xl:col-span-2">
        <span class="text-sm font-medium text-slate-700 dark:text-slate-200">{{ $t('catalog.torrents.filters.sizeRange') }}</span>
        <div class="mt-1 grid grid-cols-[minmax(0,1fr)_auto_minmax(0,1fr)] items-center gap-2">
          <UInput
            v-model="minSizeGiB"
            type="number"
            min="0"
            step="0.1"
            :ui="{ base: 'h-10 w-full' }"
            :placeholder="$t('catalog.torrents.filters.minSize')"
            :disabled="pending"
          >
            <template #trailing><span class="text-xs text-slate-400">GiB</span></template>
          </UInput>
          <span class="text-sm text-slate-400">-</span>
          <UInput
            v-model="maxSizeGiB"
            type="number"
            min="0"
            step="0.1"
            :ui="{ base: 'h-10 w-full' }"
            :placeholder="$t('catalog.torrents.filters.maxSize')"
            :disabled="pending"
          >
            <template #trailing><span class="text-xs text-slate-400">GiB</span></template>
          </UInput>
        </div>
      </div>

      <div>
        <span class="text-sm font-medium text-slate-700 dark:text-slate-200">{{ $t('catalog.torrents.filters.featured') }}</span>
        <label class="mt-1 flex h-10 items-center justify-between rounded-md border border-slate-200 px-3 dark:border-slate-700">
          <span class="text-sm text-slate-600 dark:text-slate-300">{{ $t('catalog.torrents.filters.featuredOnly') }}</span>
          <USwitch v-model="draft.featuredOnly" size="sm" :disabled="pending" />
        </label>
      </div>
    </div>

    <p v-if="sizeRangeError" class="mt-2 text-sm text-red-600 dark:text-red-300">
      {{ sizeRangeError }}
    </p>

    <div class="mt-3 flex justify-end gap-2">
      <UButton
        type="button"
        color="neutral"
        variant="ghost"
        icon="i-lucide-rotate-ccw"
        :disabled="pending || (!hasDraftFilters && !hasAppliedFilters)"
        @click="resetFilters"
      >
        {{ $t('catalog.torrents.filters.reset') }}
      </UButton>
      <UButton
        type="button"
        color="primary"
        icon="i-lucide-check"
        :loading="pending"
        :disabled="Boolean(sizeRangeError)"
        @click="applyFilters"
      >
        {{ $t('catalog.torrents.filters.apply') }}
      </UButton>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { TorrentAdvancedFilters } from '~/composables/useCatalogTorrents'

const props = withDefaults(defineProps<{
  value: TorrentAdvancedFilters
  pending?: boolean
}>(), {
  pending: false
})

const emit = defineEmits<{
  apply: [value: TorrentAdvancedFilters]
}>()

const { t } = useI18n()
const bytesPerGiB = 1024 ** 3
const minSizeGiB = ref<string | number>('')
const maxSizeGiB = ref<string | number>('')
const draft = reactive<TorrentAdvancedFilters>(defaultFilters())

const promotionOptions = computed(() => [
  { value: 'all', label: t('catalog.torrents.filters.options.allPromotions') },
  { value: 'promoted', label: t('catalog.torrents.filters.options.promoted') },
  { value: 'normal', label: t('catalog.torrents.filters.options.normal') },
  { value: 'free', label: t('catalog.torrents.status.promotion.1') },
  { value: '2x', label: t('catalog.torrents.status.promotion.2') },
  { value: '2x_free', label: t('catalog.torrents.status.promotion.3') },
  { value: '50_percent', label: t('catalog.torrents.status.promotion.4') },
  { value: '2x_50_percent', label: t('catalog.torrents.status.promotion.5') },
  { value: '30_percent', label: t('catalog.torrents.status.promotion.6') }
])
const seedStatusOptions = computed(() => [
  { value: 'all', label: t('catalog.torrents.filters.options.allSeedStatus') },
  { value: 'seeded', label: t('catalog.torrents.filters.options.seeded') },
  { value: 'unseeded', label: t('catalog.torrents.filters.options.unseeded') }
])
const publishedWithinOptions = computed(() => [
  { value: 0, label: t('catalog.torrents.filters.options.anyTime') },
  { value: 1, label: t('catalog.torrents.filters.options.withinDays', { days: 1 }) },
  { value: 7, label: t('catalog.torrents.filters.options.withinDays', { days: 7 }) },
  { value: 30, label: t('catalog.torrents.filters.options.withinDays', { days: 30 }) },
  { value: 90, label: t('catalog.torrents.filters.options.withinDays', { days: 90 }) }
])
const sortOptions = computed(() => [
  { value: 'newest', label: t('catalog.torrents.filters.options.sortNewest') },
  { value: 'oldest', label: t('catalog.torrents.filters.options.sortOldest') },
  { value: 'seeders', label: t('catalog.torrents.filters.options.sortSeeders') },
  { value: 'leechers', label: t('catalog.torrents.filters.options.sortLeechers') },
  { value: 'completed', label: t('catalog.torrents.filters.options.sortCompleted') },
  { value: 'size_asc', label: t('catalog.torrents.filters.options.sortSizeAsc') },
  { value: 'size_desc', label: t('catalog.torrents.filters.options.sortSizeDesc') }
])
const sizeRangeError = computed(() => {
  const minSize = gibToBytes(minSizeGiB.value)
  const maxSize = gibToBytes(maxSizeGiB.value)
  return minSize > 0 && maxSize > 0 && minSize > maxSize
    ? t('catalog.torrents.filters.sizeRangeInvalid')
    : ''
})
const hasDraftFilters = computed(() => filterCount({
  ...draft,
  minSize: gibToBytes(minSizeGiB.value),
  maxSize: gibToBytes(maxSizeGiB.value)
}) > 0)
const hasAppliedFilters = computed(() => filterCount(props.value) > 0)

watch(() => props.value, resetDraftFromValue, { deep: true, immediate: true })

function applyFilters() {
  if (sizeRangeError.value) return
  emit('apply', {
    ...draft,
    minSize: gibToBytes(minSizeGiB.value),
    maxSize: gibToBytes(maxSizeGiB.value)
  })
}

function resetFilters() {
  const value = defaultFilters()
  Object.assign(draft, value)
  minSizeGiB.value = ''
  maxSizeGiB.value = ''
  emit('apply', value)
}

function resetDraftFromValue(value: TorrentAdvancedFilters) {
  Object.assign(draft, value)
  minSizeGiB.value = bytesToGiB(value.minSize)
  maxSizeGiB.value = bytesToGiB(value.maxSize)
}

function defaultFilters(): TorrentAdvancedFilters {
  return {
    promotion: 'all',
    seedStatus: 'all',
    featuredOnly: false,
    minSize: 0,
    maxSize: 0,
    publishedWithin: 0,
    sort: 'newest'
  }
}

function filterCount(value: TorrentAdvancedFilters) {
  return [
    value.promotion !== 'all',
    value.seedStatus !== 'all',
    value.featuredOnly,
    value.minSize > 0,
    value.maxSize > 0,
    value.publishedWithin > 0,
    value.sort !== 'newest'
  ].filter(Boolean).length
}

function bytesToGiB(bytes: number) {
  if (!Number.isFinite(bytes) || bytes <= 0) return ''
  return String(Number((bytes / bytesPerGiB).toFixed(2)))
}

function gibToBytes(value: string | number) {
  const gib = Number(value)
  if (!Number.isFinite(gib) || gib <= 0) return 0
  return Math.min(Number.MAX_SAFE_INTEGER, Math.round(gib * bytesPerGiB))
}
</script>
