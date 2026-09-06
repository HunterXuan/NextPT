<template>
  <div class="space-y-3">
    <p class="text-xs leading-5 text-slate-500 dark:text-slate-400">{{ $t('admin.site.configs.advertisements.hint') }}</p>

    <section v-for="advertisement in advertisements" :key="advertisement.placement" class="overflow-hidden rounded-md border border-slate-200 dark:border-slate-800">
      <header class="flex flex-col gap-2 border-b border-slate-200 bg-slate-50 px-3 py-2.5 sm:flex-row sm:items-center sm:justify-between dark:border-slate-800 dark:bg-slate-950/60">
        <div class="min-w-0">
          <h3 class="text-sm font-semibold text-slate-950 dark:text-white">{{ placementLabel(advertisement.placement) }}</h3>
          <p class="mt-0.5 text-xs text-slate-500 dark:text-slate-400">{{ placementHint(advertisement.placement) }}</p>
        </div>
        <USwitch v-model="advertisement.enabled" class="shrink-0" :disabled="disabled" />
      </header>

      <div class="space-y-3 p-3">
        <div class="space-y-3">
          <UFormField :label="$t('admin.site.configs.advertisements.title')">
            <UInput v-model="advertisement.title" class="w-full" size="lg" :disabled="disabled" />
          </UFormField>
          <UFormField :label="$t('admin.site.configs.advertisements.image')">
            <UInput v-model="advertisement.image" class="w-full" size="lg" type="url" :disabled="disabled" @update:model-value="clearImageError(advertisement.placement)" />
          </UFormField>
          <UFormField :label="$t('admin.site.configs.advertisements.aspectRatio')" :description="$t('admin.site.configs.advertisements.aspectRatioHint', { ratio: defaultAspectRatio(advertisement.placement) })">
            <UInput v-model="advertisement.aspectRatio" class="w-full" size="lg" :placeholder="defaultAspectRatio(advertisement.placement)" :disabled="disabled" />
          </UFormField>
          <UFormField :label="$t('admin.site.configs.advertisements.url')">
            <UInput v-model="advertisement.url" class="w-full" size="lg" type="url" :disabled="disabled" />
          </UFormField>
        </div>

        <div>
          <p class="mb-1.5 text-xs font-medium text-slate-500 dark:text-slate-400">{{ $t('admin.site.configs.advertisements.preview') }}</p>
          <div class="relative overflow-hidden rounded-md border border-slate-200 bg-slate-100 dark:border-slate-800 dark:bg-slate-950" :style="previewStyle(advertisement.aspectRatio, advertisement.placement)">
            <img
              v-if="advertisement.image && !imageErrors[advertisement.placement]"
              :src="advertisement.image"
              :alt="advertisement.title"
              class="h-full w-full object-cover"
              @error="imageErrors[advertisement.placement] = true"
            >
            <div v-else class="flex h-full items-center justify-center gap-2 px-3 text-center text-xs text-slate-400 dark:text-slate-500">
              <UIcon name="i-lucide-image" class="size-4 shrink-0" />
              <span>{{ $t('admin.site.configs.advertisements.imagePlaceholder') }}</span>
            </div>
          </div>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
type AdvertisementPlacement = 'home' | 'catalog_list' | 'forum_list'

interface AdvertisementForm {
  placement: AdvertisementPlacement
  enabled: boolean
  title: string
  image: string
  aspectRatio: string
  url: string
}

interface ValidationResult {
  valid: boolean
  value?: unknown
  message?: string
}

const advertisementPlacements: AdvertisementPlacement[] = ['home', 'catalog_list', 'forum_list']
const defaultAspectRatios: Record<AdvertisementPlacement, string> = {
  home: '8:1',
  catalog_list: '10:1',
  forum_list: '10:1'
}

const props = withDefaults(defineProps<{
  modelValue: unknown
  disabled?: boolean
}>(), {
  disabled: false
})

const emit = defineEmits<{
  'update:modelValue': [value: unknown]
}>()

const { t } = useI18n()
const advertisements = ref<AdvertisementForm[]>([])
const imageErrors = reactive<Record<string, boolean>>({})
let resetting = false
let lastEmittedSnapshot = ''

watch(() => props.modelValue, resetFromModelValue, { deep: true, immediate: true })
watch(advertisements, emitValue, { deep: true })

function resetFromModelValue() {
  const snapshot = JSON.stringify(props.modelValue)
  if (snapshot === lastEmittedSnapshot) return

  const value = props.modelValue && typeof props.modelValue === 'object' && !Array.isArray(props.modelValue)
    ? props.modelValue as Record<string, unknown>
    : {}
  resetting = true
  advertisements.value = advertisementPlacements.map(placement => {
    const item = value[placement] && typeof value[placement] === 'object' && !Array.isArray(value[placement])
      ? value[placement] as Record<string, unknown>
      : {}
    return {
      placement,
      enabled: Boolean(item.enabled),
      title: String(item.title || '').trim(),
      image: String(item.image || '').trim(),
      aspectRatio: normalizedAspectRatio(item.aspectRatio, placement),
      url: String(item.url || '').trim()
    }
  })
  Object.keys(imageErrors).forEach(key => { delete imageErrors[key] })
  nextTick(() => { resetting = false })
}

function emitValue() {
  if (resetting) return
  const value = buildValue()
  lastEmittedSnapshot = JSON.stringify(value)
  emit('update:modelValue', value)
}

function buildValue() {
  return Object.fromEntries(advertisements.value.map(advertisement => [advertisement.placement, {
    enabled: advertisement.enabled,
    title: advertisement.title.trim(),
    image: advertisement.image.trim(),
    aspectRatio: normalizedAspectRatio(advertisement.aspectRatio, advertisement.placement),
    url: advertisement.url.trim()
  }]))
}

function validate(): ValidationResult {
  const value = buildValue()
  for (const advertisement of advertisements.value) {
    const placement = placementLabel(advertisement.placement)
    if (!isValidAspectRatio(normalizedAspectRatio(advertisement.aspectRatio, advertisement.placement))) {
      return { valid: false, message: t('admin.site.configs.advertisements.errors.aspectRatio', { placement }) }
    }
    if (!advertisement.enabled) continue
    if (!advertisement.title.trim() || advertisement.title.trim().length > 120) {
      return { valid: false, message: t('admin.site.configs.advertisements.errors.title', { placement }) }
    }
    if (!isAllowedUrl(advertisement.image)) {
      return { valid: false, message: t('admin.site.configs.advertisements.errors.image', { placement }) }
    }
    if (!isAllowedUrl(advertisement.url)) {
      return { valid: false, message: t('admin.site.configs.advertisements.errors.url', { placement }) }
    }
  }
  return { valid: true, value }
}

function isAllowedUrl(value: string) {
  if (value.startsWith('/') && !value.startsWith('//')) return true
  try {
    const parsed = new URL(value)
    return parsed.protocol === 'http:' || parsed.protocol === 'https:'
  } catch {
    return false
  }
}

function placementLabel(placement: AdvertisementPlacement) {
  return t(`admin.site.configs.advertisements.placements.${placement}.label`)
}

function placementHint(placement: AdvertisementPlacement) {
  return t(`admin.site.configs.advertisements.placements.${placement}.hint`)
}

function defaultAspectRatio(placement: AdvertisementPlacement) {
  return defaultAspectRatios[placement]
}

function normalizedAspectRatio(value: unknown, placement: AdvertisementPlacement) {
  const matched = String(value || '').match(/^\s*(\d{1,3})\s*:\s*(\d{1,3})\s*$/)
  if (!matched) return String(value || '').trim() || defaultAspectRatio(placement)
  return `${matched[1]}:${matched[2]}`
}

function isValidAspectRatio(value: string) {
  const matched = value.match(/^(\d{1,3}):(\d{1,3})$/)
  return Boolean(matched && Number(matched[1]) > 0 && Number(matched[2]) > 0)
}

function previewStyle(value: string, placement: AdvertisementPlacement) {
  const ratio = normalizedAspectRatio(value, placement).replace(':', ' / ')
  return { aspectRatio: ratio }
}

function clearImageError(placement: AdvertisementPlacement) {
  imageErrors[placement] = false
}

defineExpose({ validate })
</script>
