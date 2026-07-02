<template>
  <div class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
    <div class="flex flex-wrap items-center gap-3 text-sm text-slate-500 dark:text-slate-400">
      <span v-if="showTotal">{{ t('common.pagination.total', { total: numberFormatter.format(safeTotal) }) }}</span>
      <label v-if="normalizedPageSizeOptions.length > 0" class="flex items-center gap-2">
        <span>{{ t('common.pagination.pageSize') }}</span>
        <select
          :value="safePageSize"
          class="h-9 rounded-md border border-slate-200 bg-white px-2 text-sm text-slate-950 outline-none transition focus:border-sky-400 focus:ring-2 focus:ring-sky-100 disabled:cursor-not-allowed disabled:opacity-60 dark:border-slate-700 dark:bg-slate-900 dark:text-white dark:focus:border-sky-500 dark:focus:ring-sky-950"
          :disabled="disabled"
          @change="handlePageSizeChange"
        >
          <option v-for="option in normalizedPageSizeOptions" :key="option" :value="option">
            {{ option }}
          </option>
        </select>
      </label>
    </div>

    <nav v-if="totalPages > 0" class="flex flex-wrap items-center gap-1" :aria-label="t('common.pagination.label')">
      <UButton
        color="neutral"
        variant="outline"
        :size="size"
        icon="i-lucide-chevrons-left"
        class="h-9 min-w-9 justify-center px-2"
        :aria-label="t('common.pagination.first')"
        :disabled="safePage <= 1 || disabled"
        @click="changePage(1)"
      />
      <UButton
        color="neutral"
        variant="outline"
        :size="size"
        icon="i-lucide-chevron-left"
        class="h-9 min-w-9 justify-center px-2"
        :aria-label="t('common.previous')"
        :disabled="safePage <= 1 || disabled"
        @click="changePage(safePage - 1)"
      />

      <template v-for="item in pageItems" :key="item.key">
        <span
          v-if="item.type === 'ellipsis'"
          class="flex h-9 min-w-9 items-center justify-center px-2 text-sm text-slate-400 dark:text-slate-500"
        >
          ...
        </span>
        <UButton
          v-else
          :color="item.value === safePage ? 'primary' : 'neutral'"
          :variant="item.value === safePage ? 'soft' : 'outline'"
          :size="size"
          class="h-9 min-w-9 justify-center px-2"
          :aria-current="item.value === safePage ? 'page' : undefined"
          :aria-label="t('common.pagination.page', { page: item.value })"
          :disabled="disabled"
          @click="changePage(item.value)"
        >
          {{ item.value }}
        </UButton>
      </template>

      <UButton
        color="neutral"
        variant="outline"
        :size="size"
        icon="i-lucide-chevron-right"
        class="h-9 min-w-9 justify-center px-2"
        :aria-label="t('common.next')"
        :disabled="safePage >= totalPages || disabled"
        @click="changePage(safePage + 1)"
      />
      <UButton
        color="neutral"
        variant="outline"
        :size="size"
        icon="i-lucide-chevrons-right"
        class="h-9 min-w-9 justify-center px-2"
        :aria-label="t('common.pagination.last')"
        :disabled="safePage >= totalPages || disabled"
        @click="changePage(totalPages)"
      />
    </nav>
  </div>
</template>

<script setup lang="ts">
type PagerSize = 'xs' | 'sm' | 'md'
type PageItem = { type: 'page', key: string, value: number } | { type: 'ellipsis', key: string }

const props = withDefaults(defineProps<{
  page: number
  total: number
  pageSize: number
  pageSizeOptions?: number[]
  disabled?: boolean
  showTotal?: boolean
  size?: PagerSize
}>(), {
  pageSizeOptions: () => [],
  disabled: false,
  showTotal: true,
  size: 'sm'
})

const emit = defineEmits<{
  'page-change': [page: number]
  'page-size-change': [pageSize: number]
}>()

const { t, locale } = useI18n()

const numberFormatter = computed(() => new Intl.NumberFormat(locale.value))
const safeTotal = computed(() => Math.max(0, Number(props.total) || 0))
const safePageSize = computed(() => Math.max(1, Number(props.pageSize) || 1))
const totalPages = computed(() => safeTotal.value > 0 ? Math.ceil(safeTotal.value / safePageSize.value) : 0)
const safePage = computed(() => clampPage(props.page))

const normalizedPageSizeOptions = computed(() => {
  return Array.from(new Set(props.pageSizeOptions.map((value) => Number(value)).filter((value) => Number.isInteger(value) && value > 0)))
})

const pageItems = computed<PageItem[]>(() => {
  const pages = buildVisiblePages(safePage.value, totalPages.value)
  const items: PageItem[] = []

  pages.forEach((page, index) => {
    const previous = pages[index - 1]
    if (previous && page - previous > 1) {
      items.push({ type: 'ellipsis', key: `ellipsis-${previous}-${page}` })
    }
    items.push({ type: 'page', key: `page-${page}`, value: page })
  })

  return items
})

function buildVisiblePages(page: number, pages: number) {
  if (pages <= 0) return []

  if (pages <= 7) {
    return Array.from({ length: pages }, (_, index) => index + 1)
  }

  const visible = new Set([1, pages, page, page - 1, page + 1])
  if (page <= 4) {
    visible.add(2)
    visible.add(3)
    visible.add(4)
    visible.add(5)
  }
  if (page >= pages - 3) {
    visible.add(pages - 4)
    visible.add(pages - 3)
    visible.add(pages - 2)
    visible.add(pages - 1)
  }

  return Array.from(visible)
    .filter((value) => value >= 1 && value <= pages)
    .sort((left, right) => left - right)
}

function clampPage(page: number) {
  if (totalPages.value <= 0) return 1

  const parsed = Number(page)
  if (!Number.isFinite(parsed)) return 1
  return Math.min(Math.max(1, Math.trunc(parsed)), totalPages.value)
}

function changePage(page: number) {
  if (props.disabled) return

  const nextPage = clampPage(page)
  if (nextPage === safePage.value) return
  emit('page-change', nextPage)
}

function handlePageSizeChange(event: Event) {
  if (props.disabled) return

  const target = event.target as HTMLSelectElement | null
  const nextPageSize = Number(target?.value)
  if (!Number.isInteger(nextPageSize) || nextPageSize <= 0 || nextPageSize === safePageSize.value) return
  emit('page-size-change', nextPageSize)
}
</script>
