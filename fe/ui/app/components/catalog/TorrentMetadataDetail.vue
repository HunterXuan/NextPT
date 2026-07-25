<template>
  <section id="torrent-metadata" class="scroll-mt-24 overflow-hidden rounded-lg border border-slate-200 bg-white dark:border-slate-800 dark:bg-slate-900">
    <div class="border-b border-slate-200 px-4 py-3 dark:border-slate-800">
      <h2 class="text-base font-semibold text-slate-950 dark:text-white">{{ $t('catalog.torrents.metadata.title') }}</h2>
    </div>
    <div class="p-4">
      <div v-if="metadata.data" class="grid gap-4 sm:grid-cols-[128px_minmax(0,1fr)]">
        <img
          v-if="metadata.data.posterUrl"
          :src="metadata.data.posterUrl"
          :alt="metadata.data.title"
          class="aspect-[2/3] w-32 rounded-md object-cover ring-1 ring-slate-200 dark:ring-slate-700"
          referrerpolicy="no-referrer"
        >
        <div v-else class="flex aspect-[2/3] w-32 items-center justify-center rounded-md bg-slate-100 text-slate-400 ring-1 ring-slate-200 dark:bg-slate-800 dark:ring-slate-700">
          <UIcon name="i-lucide-image-off" class="size-6" />
        </div>

        <div class="min-w-0">
          <div class="flex flex-wrap items-baseline gap-x-2 gap-y-1">
            <h3 class="break-words text-lg font-semibold text-slate-950 dark:text-white">{{ metadata.data.title }}</h3>
            <span v-if="metadata.data.year" class="text-sm text-slate-500 dark:text-slate-400">{{ metadata.data.year }}</span>
          </div>
          <p v-if="metadata.data.originalTitle && metadata.data.originalTitle !== metadata.data.title" class="mt-1 text-sm text-slate-500 dark:text-slate-400">
            {{ metadata.data.originalTitle }}
          </p>

          <div v-if="sources.length" class="mt-3 flex max-w-full items-center overflow-x-auto py-0.5">
            <UTooltip v-for="source in sources" :key="source.provider" :text="source.label">
              <NuxtLink
                :to="source.url"
                target="_blank"
                rel="noopener noreferrer"
                class="group flex h-7 shrink-0 items-center gap-1.5 border-l border-slate-200 px-3 text-sm first:border-l-0 first:pl-0 dark:border-slate-700"
              >
                <span class="font-medium text-slate-500 transition-colors group-hover:text-slate-950 dark:text-slate-400 dark:group-hover:text-white">{{ source.label }}</span>
                <span v-if="source.rating" class="inline-flex items-center gap-1 font-semibold tabular-nums text-slate-950 dark:text-white">
                  <UIcon name="i-lucide-star" class="size-3.5 fill-amber-400 text-amber-500" />
                  {{ source.rating.toFixed(1) }}
                </span>
                <UIcon name="i-lucide-arrow-up-right" class="size-3 text-slate-400 transition-colors group-hover:text-slate-700 dark:group-hover:text-slate-200" />
              </NuxtLink>
            </UTooltip>
          </div>

          <div v-if="metadata.data.genres?.length" class="mt-3 flex flex-wrap gap-1.5">
            <UBadge v-for="genre in metadata.data.genres" :key="genre" color="neutral" variant="soft">{{ genre }}</UBadge>
          </div>
          <p v-if="metadata.data.overview" class="mt-3 text-sm leading-6 text-slate-600 dark:text-slate-300">{{ metadata.data.overview }}</p>
          <p v-else class="mt-3 text-sm text-slate-500 dark:text-slate-400">{{ $t('catalog.torrents.metadata.noOverview') }}</p>
        </div>
      </div>

    </div>
  </section>
</template>

<script setup lang="ts">
import type { TorrentMetadataOut } from '~/composables/useCatalogTorrents'

const props = defineProps<{
  metadata: TorrentMetadataOut
}>()

const { t } = useI18n()

const sources = computed(() => {
  const binding = props.metadata.binding
  const ratings = new Map((props.metadata.sources || []).map((source) => [source.provider, source.rating]))
  const items: Array<{ provider: string, label: string, url: string, rating: number }> = []
  if (binding.tmdbId && binding.tmdbType) items.push({ provider: 'tmdb', label: providerLabel('tmdb'), url: `https://www.themoviedb.org/${binding.tmdbType}/${binding.tmdbId}`, rating: ratings.get('tmdb') || 0 })
  if (binding.imdbId) items.push({ provider: 'imdb', label: providerLabel('imdb'), url: `https://www.imdb.com/title/${binding.imdbId}/`, rating: ratings.get('imdb') || 0 })
  if (binding.doubanId) items.push({ provider: 'douban', label: providerLabel('douban'), url: `https://movie.douban.com/subject/${binding.doubanId}/`, rating: ratings.get('douban') || 0 })
  if (binding.bangumiId) items.push({ provider: 'bangumi', label: providerLabel('bangumi'), url: `https://bgm.tv/subject/${binding.bangumiId}`, rating: ratings.get('bangumi') || 0 })
  return items
})

function providerLabel(provider: string) {
  if (provider === 'imdb') return 'IMDb'
  if (provider === 'tmdb' || provider === 'douban' || provider === 'bangumi') {
    return t(`catalog.torrents.metadata.providers.${provider}`)
  }
  return provider
}
</script>
