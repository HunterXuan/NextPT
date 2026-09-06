<template>
  <a
    v-if="item && !imageFailed"
    :href="item.url"
    :title="item.title"
    :target="isExternal ? '_blank' : undefined"
    :rel="isExternal ? 'noopener noreferrer' : undefined"
    class="block overflow-hidden rounded-lg border border-slate-200 bg-slate-100 outline-none transition hover:border-sky-300 focus-visible:ring-2 focus-visible:ring-sky-400 dark:border-slate-800 dark:bg-slate-900 dark:hover:border-sky-700 dark:focus-visible:ring-sky-600"
    :style="advertisementStyle"
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
const advertisementStyle = computed(() => ({ aspectRatio: normalizedAspectRatio(item.value?.aspectRatio) }))

onMounted(load)
watch(() => item.value?.image, () => { imageFailed.value = false })

function normalizedAspectRatio(value?: string) {
  const matched = String(value || '').match(/^\s*(\d{1,3})\s*:\s*(\d{1,3})\s*$/)
  return matched && Number(matched[1]) > 0 && Number(matched[2]) > 0 ? `${matched[1]} / ${matched[2]}` : '8 / 1'
}
</script>
