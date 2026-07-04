<template>
  <div
    class="overflow-hidden rounded-md border bg-white transition focus-within:ring-2 dark:bg-slate-950"
    :class="valid
      ? 'border-slate-200 focus-within:border-sky-400 focus-within:ring-sky-100 dark:border-slate-700 dark:focus-within:border-sky-500 dark:focus-within:ring-sky-950'
      : 'border-red-300 focus-within:border-red-400 focus-within:ring-red-100 dark:border-red-800 dark:focus-within:border-red-600 dark:focus-within:ring-red-950'"
  >
    <div class="flex items-center justify-between gap-3 border-b border-slate-200 bg-slate-50 px-3 py-2 dark:border-slate-800 dark:bg-slate-900/70">
      <span class="inline-flex items-center gap-1.5 text-xs font-semibold uppercase tracking-wide text-slate-500 dark:text-slate-400">
        <UIcon name="i-lucide-braces" class="size-3.5" />
        {{ title }}
      </span>
      <div class="flex items-center gap-2">
        <span
          class="inline-flex items-center gap-1 rounded px-1.5 py-0.5 text-xs font-medium"
          :class="valid ? 'bg-emerald-50 text-emerald-700 dark:bg-emerald-950/40 dark:text-emerald-300' : 'bg-red-50 text-red-700 dark:bg-red-950/40 dark:text-red-300'"
        >
          <UIcon :name="valid ? 'i-lucide-check' : 'i-lucide-circle-alert'" class="size-3.5" />
          {{ valid ? validText : invalidText }}
        </span>
        <UTooltip :text="fullscreenText" :content="{ side: 'top', sideOffset: 8 }" :delay-duration="120">
          <UButton
            type="button"
            color="neutral"
            variant="ghost"
            size="xs"
            icon="i-lucide-maximize-2"
            :aria-label="fullscreenText"
            @click="openFullscreen"
          />
        </UTooltip>
      </div>
    </div>

    <div class="grid max-h-80 min-h-44 grid-cols-[3rem_minmax(0,1fr)] items-stretch overflow-auto">
      <pre class="select-none border-r border-slate-200 bg-slate-50 px-2 py-3 text-right font-mono text-xs leading-5 text-slate-400 dark:border-slate-800 dark:bg-slate-900/70 dark:text-slate-500">{{ lineNumbers }}</pre>
      <textarea
        :value="modelValue"
        class="h-full min-h-44 w-full resize-none bg-transparent px-3 py-3 font-mono text-xs leading-5 text-slate-950 outline-none dark:text-white"
        :placeholder="placeholder"
        :disabled="disabled"
        spellcheck="false"
        @input="emitValue"
        @keydown.tab.prevent="handleTab"
      />
    </div>

    <div class="flex flex-col gap-1 border-t border-slate-200 bg-slate-50 px-3 py-2 text-xs text-slate-500 sm:flex-row sm:items-center sm:justify-between dark:border-slate-800 dark:bg-slate-900/70 dark:text-slate-400">
      <span>{{ hint }}</span>
      <span class="font-mono">{{ stats }}</span>
    </div>
  </div>

  <Teleport to="body">
    <div v-if="fullscreenOpen" class="fixed inset-0 z-[100] bg-slate-950/70 p-3 backdrop-blur-sm sm:p-5">
      <div class="flex h-full min-h-0 flex-col overflow-hidden rounded-lg border border-slate-200 bg-white shadow-2xl dark:border-slate-800 dark:bg-slate-950">
        <div class="flex items-center justify-between gap-3 border-b border-slate-200 px-4 py-3 dark:border-slate-800">
          <div class="min-w-0">
            <h2 class="truncate text-sm font-semibold text-slate-950 dark:text-white">{{ fullscreenText }}</h2>
            <p class="mt-0.5 text-xs text-slate-500 dark:text-slate-400">{{ title }}</p>
          </div>
          <div class="flex items-center gap-2">
            <span
              class="inline-flex items-center gap-1 rounded px-1.5 py-0.5 text-xs font-medium"
              :class="valid ? 'bg-emerald-50 text-emerald-700 dark:bg-emerald-950/40 dark:text-emerald-300' : 'bg-red-50 text-red-700 dark:bg-red-950/40 dark:text-red-300'"
            >
              <UIcon :name="valid ? 'i-lucide-check' : 'i-lucide-circle-alert'" class="size-3.5" />
              {{ valid ? validText : invalidText }}
            </span>
            <UButton type="button" color="neutral" variant="ghost" size="sm" icon="i-lucide-x" :aria-label="closeText" @click="closeFullscreen" />
          </div>
        </div>

        <div class="grid min-h-0 flex-1 grid-cols-[3.5rem_minmax(0,1fr)] overflow-auto">
          <pre class="select-none border-r border-slate-200 bg-slate-50 px-2 py-3 text-right font-mono text-xs leading-5 text-slate-400 dark:border-slate-800 dark:bg-slate-900/70 dark:text-slate-500">{{ lineNumbers }}</pre>
          <textarea
            ref="fullscreenTextarea"
            :value="modelValue"
            class="h-full min-h-[60vh] w-full resize-none bg-transparent px-3 py-3 font-mono text-xs leading-5 text-slate-950 outline-none dark:text-white"
            :placeholder="placeholder"
            :disabled="disabled"
            spellcheck="false"
            @input="emitValue"
            @keydown.tab.prevent="handleTab"
          />
        </div>

        <div class="flex flex-col gap-2 border-t border-slate-200 bg-slate-50 px-4 py-3 text-xs text-slate-500 sm:flex-row sm:items-center sm:justify-between dark:border-slate-800 dark:bg-slate-900/70 dark:text-slate-400">
          <span>{{ valid ? hint : errorMessage }}</span>
          <span class="font-mono">{{ stats }}</span>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<script setup lang="ts">
const props = withDefaults(defineProps<{
  modelValue: string
  valid: boolean
  errorMessage?: string
  title: string
  placeholder?: string
  hint: string
  validText: string
  invalidText: string
  fullscreenText: string
  closeText: string
  disabled?: boolean
}>(), {
  errorMessage: '',
  placeholder: '',
  disabled: false
})

const emit = defineEmits<{
  'update:modelValue': [value: string]
}>()

const fullscreenOpen = ref(false)
const fullscreenTextarea = ref<HTMLTextAreaElement | null>(null)

const lineCount = computed(() => Math.max(1, props.modelValue.split('\n').length))
const lineNumbers = computed(() => Array.from({ length: lineCount.value }, (_, index) => index + 1).join('\n'))
const stats = computed(() => `${lineCount.value}L / ${props.modelValue.length}C`)

function emitValue(event: Event) {
  emit('update:modelValue', (event.target as HTMLTextAreaElement).value)
}

function handleTab(event: KeyboardEvent) {
  const textarea = event.target as HTMLTextAreaElement
  const start = textarea.selectionStart
  const end = textarea.selectionEnd
  const nextValue = `${props.modelValue.slice(0, start)}  ${props.modelValue.slice(end)}`
  emit('update:modelValue', nextValue)
  nextTick(() => {
    textarea.selectionStart = start + 2
    textarea.selectionEnd = start + 2
  })
}

function openFullscreen() {
  fullscreenOpen.value = true
  nextTick(() => fullscreenTextarea.value?.focus())
}

function closeFullscreen() {
  fullscreenOpen.value = false
}
</script>
