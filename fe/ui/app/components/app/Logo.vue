<template>
  <span class="inline-flex shrink-0 items-center justify-center" :class="small ? 'size-8' : 'size-9'" aria-hidden="true">
    <template v-for="variant in variants" :key="variant.key">
      <img
        v-if="variant.src"
        :src="variant.src"
        alt=""
        class="h-full w-full object-contain"
        :class="variant.class"
        @error="failed[variant.src] = true"
      >
      <span
        v-else
        class="h-full w-full items-center justify-center rounded-md"
        :class="[variant.class, variant.key === 'dark' ? 'bg-white text-neutral-950' : 'bg-neutral-950 text-white']"
      >
        <UIcon name="i-lucide-radio-tower" :class="small ? 'size-4' : 'size-5'" />
      </span>
    </template>
  </span>
</template>

<script setup lang="ts">
const props = defineProps<{ small?: boolean; dark?: boolean }>()
const { siteLogo, siteLogoDark } = useSiteBranding()
const failed = reactive<Record<string, boolean>>({})
const lightSource = computed(() => failed[siteLogo.value] ? '' : siteLogo.value)
const darkSource = computed(() => siteLogoDark.value && !failed[siteLogoDark.value] ? siteLogoDark.value : lightSource.value)
const variants = computed(() => props.dark
  ? [{ key: 'dark', src: darkSource.value, class: 'flex' }]
  : [
      { key: 'light', src: lightSource.value, class: 'flex dark:hidden' },
      { key: 'dark', src: darkSource.value, class: 'hidden dark:flex' }
    ])
</script>
