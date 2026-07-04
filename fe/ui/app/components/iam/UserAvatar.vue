<template>
  <img
    v-if="displayedAvatarSrc"
    v-bind="rootAttrs"
    :src="displayedAvatarSrc"
    :alt="displayName"
    :class="imageClass"
    loading="lazy"
    @load="handleImageLoad"
    @error="handleImageError"
  >
  <div
    v-else
    v-bind="rootAttrs"
    :class="fallbackClass"
    :aria-label="displayName"
  >
    {{ initial }}
  </div>
</template>

<script setup lang="ts">
defineOptions({
  inheritAttrs: false
})

interface AvatarUser {
  id?: number | string
  username?: string
  avatar?: string | null
}

const props = withDefaults(defineProps<{
  user?: AvatarUser | null
  id?: number | string
  username?: string
  avatar?: string | null
  alt?: string
  size?: 'xs' | 'sm' | 'md' | 'lg' | 'xl' | '2xl' | '3xl'
}>(), {
  user: null,
  id: '',
  username: '',
  avatar: '',
  alt: '',
  size: 'md'
})

const emit = defineEmits<{
  loadError: []
  loadSuccess: []
}>()

const attrs = useAttrs()
const imageFailed = ref(false)

const sizeClasses = {
  xs: 'size-6 text-xs',
  sm: 'size-8 text-xs',
  md: 'size-10 text-sm',
  lg: 'size-12 text-lg',
  xl: 'size-14 text-xl',
  '2xl': 'size-16 text-2xl',
  '3xl': 'size-18 text-2xl'
}

const paletteClasses = [
  'bg-sky-100 text-sky-700 border-sky-200 dark:bg-sky-950 dark:text-sky-200 dark:border-sky-900',
  'bg-emerald-100 text-emerald-700 border-emerald-200 dark:bg-emerald-950 dark:text-emerald-200 dark:border-emerald-900',
  'bg-amber-100 text-amber-700 border-amber-200 dark:bg-amber-950 dark:text-amber-200 dark:border-amber-900',
  'bg-rose-100 text-rose-700 border-rose-200 dark:bg-rose-950 dark:text-rose-200 dark:border-rose-900',
  'bg-violet-100 text-violet-700 border-violet-200 dark:bg-violet-950 dark:text-violet-200 dark:border-violet-900',
  'bg-cyan-100 text-cyan-700 border-cyan-200 dark:bg-cyan-950 dark:text-cyan-200 dark:border-cyan-900'
]

const avatarSrc = computed(() => (props.avatar || props.user?.avatar || '').trim())
const displayedAvatarSrc = computed(() => imageFailed.value ? '' : avatarSrc.value)
const userId = computed(() => props.id || props.user?.id || '')
const username = computed(() => (props.username || props.user?.username || '').trim())
const displayName = computed(() => props.alt || username.value || (userId.value ? `#${userId.value}` : '?'))
const rootAttrs = computed(() => {
  const { class: _class, ...rest } = attrs
  return rest
})
const baseClass = computed(() => [sizeClasses[props.size], 'shrink-0 rounded-md border object-cover', attrs.class])
const imageClass = computed(() => [baseClass.value, 'border-slate-200 dark:border-slate-800'])
const fallbackClass = computed(() => [baseClass.value, 'flex items-center justify-center font-semibold leading-none', paletteClass.value])
const initial = computed(() => {
  if (username.value) return Array.from(username.value)[0].toUpperCase()
  if (userId.value) return String(userId.value).slice(-2)
  return '?'
})
const paletteClass = computed(() => paletteClasses[stableIndex(userId.value || username.value || displayName.value)])

watch(avatarSrc, () => {
  imageFailed.value = false
})

function handleImageError() {
  imageFailed.value = true
  emit('loadError')
}

function handleImageLoad() {
  emit('loadSuccess')
}

function stableIndex(value: number | string) {
  const numberValue = Number(value)
  if (Number.isFinite(numberValue) && String(value).trim() !== '') {
    return Math.abs(numberValue) % paletteClasses.length
  }

  const text = String(value || '')
  let total = 0
  for (const char of text) {
    total += char.charCodeAt(0)
  }
  return total % paletteClasses.length
}
</script>
