<template>
  <a
    v-if="item && !imageFailed"
    :href="item.url"
    :title="item.title"
    :target="isExternal ? '_blank' : undefined"
    :rel="isExternal ? 'noopener noreferrer' : undefined"
    class="block overflow-hidden rounded-lg border border-slate-200 bg-slate-100 outline-none transition hover:border-sky-300 focus-visible:ring-2 focus-visible:ring-sky-400 dark:border-slate-800 dark:bg-slate-900 dark:hover:border-sky-700 dark:focus-visible:ring-sky-600"
    :class="placementClass"
  >
    <img :src="item.image" :alt="item.title" class="h-full w-full object-cover" loading="lazy" @error="imageFailed = true">
  </a>
</template>

<script setup lang="ts">
import type { SiteAdvertisementPlacement } from '~/composables/useSite'

const props = defineProps<{
  placement: SiteAdvertisementPlacement
}>()

const { advertisement, load } = useSiteAdvertisements()
const item = advertisement(props.placement)
const imageFailed = ref(false)
const isExternal = computed(() => /^https?:\/\//i.test(item.value?.url || ''))
const placementClass = computed(() => ({
  home: 'aspect-[8/1]',
  catalog_list: 'h-24 sm:h-28 xl:h-32',
  forum_list: 'h-24 sm:h-28 xl:h-32'
})[props.placement])

onMounted(load)
watch(() => item.value?.image, () => { imageFailed.value = false })
</script>
