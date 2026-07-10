<template>
  <UPopover
    :open="open"
    :content="{ side: 'bottom', align: 'end', sideOffset: 8 }"
    :ui="{ content: 'w-[min(24rem,calc(100vw-1.5rem))] p-4' }"
    @update:open="open = $event"
  >
    <UTooltip
      :text="triggerTooltip"
      :disabled="open"
      :content="{ side: 'bottom', sideOffset: 8 }"
      :delay-duration="300"
    >
      <span class="inline-flex">
        <UButton
          type="button"
          color="neutral"
          variant="soft"
          icon="i-lucide-rss"
          class="h-10 w-10 justify-center p-0"
          :aria-label="$t('catalog.torrents.rss.title')"
          :disabled="!canUseRss"
        />
      </span>
    </UTooltip>

    <template #content>
      <div class="flex items-center gap-2">
        <UIcon name="i-lucide-rss" class="size-4 text-orange-500" />
        <h2 class="text-sm font-semibold text-slate-950 dark:text-white">
          {{ $t('catalog.torrents.rss.title') }}
        </h2>
      </div>

      <div class="mt-4 grid grid-cols-[minmax(0,1fr)_116px] items-end gap-3">
        <div class="flex h-9 items-center justify-between rounded-md border border-slate-200 px-3 dark:border-slate-700">
          <span class="text-sm text-slate-700 dark:text-slate-200">{{ $t('catalog.torrents.rss.promotionOnly') }}</span>
          <USwitch v-model="promotionOnly" size="sm" />
        </div>
        <UFormField :label="$t('catalog.torrents.rss.size')">
          <USelect
            v-model="feedSize"
            class="w-full"
            :items="feedSizeOptions"
            value-key="value"
          />
        </UFormField>
      </div>

      <div class="mt-4">
        <p class="text-xs font-medium text-slate-500 dark:text-slate-400">
          {{ $t('catalog.torrents.rss.url') }}
        </p>
        <div class="mt-1.5 flex items-center gap-1.5">
          <input
            :value="rssUrl"
            readonly
            class="h-9 min-w-0 flex-1 truncate rounded-md border border-slate-200 bg-slate-50 px-2.5 font-mono text-xs text-slate-700 outline-none dark:border-slate-700 dark:bg-slate-950 dark:text-slate-300"
            @focus="selectInput"
          >
          <UTooltip :text="$t('common.copy')" :content="{ side: 'top', sideOffset: 8 }" :delay-duration="120">
            <UButton
              type="button"
              color="neutral"
              variant="outline"
              icon="i-lucide-copy"
              class="h-9 w-9 justify-center p-0"
              :aria-label="$t('common.copy')"
              @click="copyRssUrl"
            />
          </UTooltip>
        </div>
        <p class="mt-2 text-xs leading-5 text-amber-700 dark:text-amber-300">
          {{ $t('catalog.torrents.rss.passkeyWarning') }}
        </p>
      </div>
    </template>
  </UPopover>
</template>

<script setup lang="ts">
const props = withDefaults(defineProps<{
  keyword?: string
  categoryIds?: number[]
}>(), {
  keyword: '',
  categoryIds: () => []
})

const { t, locale } = useI18n()
const toast = useToast()
const requestUrl = useRequestURL()
const { user, hasPermission } = useAuth()

const open = ref(false)
const promotionOnly = ref(false)
const feedSize = ref(50)
const feedSizeOptions = computed(() => [20, 50, 100].map((value) => ({
  label: t('catalog.torrents.rss.sizeOption', { count: value }),
  value
})))

const passkey = computed(() => user.value?.user.passkey || '')
const canUseRss = computed(() => Boolean(passkey.value) && hasPermission(Permission.CatalogTorrentDownload))
const triggerTooltip = computed(() => {
  if (!hasPermission(Permission.CatalogTorrentDownload)) return t('common.noPermission')
  if (!passkey.value) return t('catalog.torrents.rss.passkeyRequired')
  return t('catalog.torrents.rss.title')
})

const rssUrl = computed(() => {
  const origin = import.meta.client ? window.location.origin : requestUrl.origin
  const url = new URL('/api/tracker/rss', origin)
  url.searchParams.set('passkey', passkey.value)
  url.searchParams.set('lang', locale.value)
  url.searchParams.set('size', String(feedSize.value))

  const keyword = props.keyword.trim()
  if (keyword) url.searchParams.set('keyword', keyword)
  for (const categoryId of props.categoryIds) {
    url.searchParams.append('categoryIds[]', String(categoryId))
  }
  if (promotionOnly.value) url.searchParams.set('promotionOnly', 'true')
  return url.toString()
})

function selectInput(event: FocusEvent) {
  (event.target as HTMLInputElement).select()
}

async function copyRssUrl() {
  if (!canUseRss.value || !navigator?.clipboard) return
  await navigator.clipboard.writeText(rssUrl.value)
  toast.add({
    title: t('catalog.torrents.rss.copied'),
    color: 'success',
    icon: 'i-lucide-check-circle'
  })
}
</script>
