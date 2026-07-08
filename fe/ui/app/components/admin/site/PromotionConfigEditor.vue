<template>
  <div class="space-y-4">
    <div class="flex flex-col gap-3 rounded-md border border-slate-200 bg-slate-50 px-3 py-3 sm:flex-row sm:items-center sm:justify-between dark:border-slate-800 dark:bg-slate-950/70">
      <div>
        <p class="text-sm font-medium text-slate-950 dark:text-white">{{ $t('admin.site.configs.promotion.editMode') }}</p>
        <p class="mt-1 text-xs text-slate-500 dark:text-slate-400">{{ $t('admin.site.configs.promotion.editModeHint') }}</p>
      </div>
      <div class="inline-flex h-9 rounded-md border border-slate-200 bg-white p-0.5 dark:border-slate-700 dark:bg-slate-950">
        <button
          type="button"
          class="inline-flex min-w-20 items-center justify-center gap-1.5 rounded px-3 text-sm font-medium transition"
          :class="editorMode === 'form' ? 'bg-indigo-600 text-white shadow-sm' : 'text-slate-600 hover:bg-slate-100 dark:text-slate-300 dark:hover:bg-slate-800'"
          :disabled="disabled"
          @click="setEditorMode('form')"
        >
          <UIcon name="i-lucide-sliders-horizontal" class="size-4" />
          {{ $t('admin.site.configs.promotion.modeForm') }}
        </button>
        <button
          type="button"
          class="inline-flex min-w-20 items-center justify-center gap-1.5 rounded px-3 text-sm font-medium transition"
          :class="editorMode === 'json' ? 'bg-indigo-600 text-white shadow-sm' : 'text-slate-600 hover:bg-slate-100 dark:text-slate-300 dark:hover:bg-slate-800'"
          :disabled="disabled"
          @click="setEditorMode('json')"
        >
          <UIcon name="i-lucide-braces" class="size-4" />
          {{ $t('admin.site.configs.promotion.modeJson') }}
        </button>
      </div>
    </div>

    <div v-if="editorMode === 'form'" class="space-y-4">
      <div v-if="isGlobalPromotion" class="space-y-4">
        <label class="flex items-center justify-between gap-3 rounded-md border border-slate-200 px-3 py-3 dark:border-slate-800">
          <span>
            <span class="block text-sm font-medium text-slate-950 dark:text-white">{{ $t('admin.site.configs.promotion.enabled') }}</span>
            <span class="mt-1 block text-xs text-slate-500 dark:text-slate-400">{{ $t('admin.site.configs.promotion.globalEnabledHint') }}</span>
          </span>
          <button
            type="button"
            role="switch"
            :aria-checked="globalForm.enabled"
            class="relative h-6 w-11 rounded-full transition disabled:cursor-not-allowed disabled:opacity-60"
            :class="globalForm.enabled ? 'bg-indigo-600' : 'bg-slate-300 dark:bg-slate-700'"
            :disabled="disabled"
            @click="globalForm.enabled = !globalForm.enabled"
          >
            <span
              class="absolute top-0.5 size-5 rounded-full bg-white shadow-sm transition"
              :class="globalForm.enabled ? 'left-5' : 'left-0.5'"
            />
          </button>
        </label>

        <div class="grid gap-3 sm:grid-cols-2">
          <label class="block">
            <span class="text-sm font-medium text-slate-700 dark:text-slate-200">{{ $t('admin.site.configs.promotion.state') }}</span>
            <select
              v-model="globalForm.state"
              class="mt-1 h-10 w-full rounded-md border border-slate-200 bg-white px-3 text-sm text-slate-950 outline-none transition focus:border-indigo-400 focus:ring-2 focus:ring-indigo-100 disabled:cursor-not-allowed disabled:bg-slate-50 disabled:text-slate-400 dark:border-slate-700 dark:bg-slate-950 dark:text-white dark:focus:border-indigo-500 dark:focus:ring-indigo-950 dark:disabled:bg-slate-900 dark:disabled:text-slate-500"
              :disabled="disabled || !globalForm.enabled"
            >
              <option v-for="option in promotionStateOptions" :key="option.value" :value="option.value">{{ option.label }}</option>
            </select>
          </label>
          <label class="block">
            <span class="text-sm font-medium text-slate-700 dark:text-slate-200">{{ $t('admin.site.configs.promotion.expireAt') }}</span>
            <input
              v-model="globalForm.expireAt"
              type="datetime-local"
              step="1"
              class="mt-1 h-10 w-full rounded-md border border-slate-200 bg-white px-3 text-sm text-slate-950 outline-none transition focus:border-indigo-400 focus:ring-2 focus:ring-indigo-100 disabled:cursor-not-allowed disabled:bg-slate-50 disabled:text-slate-400 dark:border-slate-700 dark:bg-slate-950 dark:text-white dark:focus:border-indigo-500 dark:focus:ring-indigo-950 dark:disabled:bg-slate-900 dark:disabled:text-slate-500"
              :disabled="disabled || !globalForm.enabled"
            >
            <span class="mt-1 block text-xs text-slate-500 dark:text-slate-400">{{ $t('admin.site.configs.promotion.expireAtHint') }}</span>
          </label>
        </div>
      </div>

      <div v-else class="space-y-4">
        <label class="flex items-center justify-between gap-3 rounded-md border border-slate-200 px-3 py-3 dark:border-slate-800">
          <span>
            <span class="block text-sm font-medium text-slate-950 dark:text-white">{{ $t('admin.site.configs.promotion.enabled') }}</span>
            <span class="mt-1 block text-xs text-slate-500 dark:text-slate-400">{{ $t('admin.site.configs.promotion.newEnabledHint') }}</span>
          </span>
          <button
            type="button"
            role="switch"
            :aria-checked="newForm.enabled"
            class="relative h-6 w-11 rounded-full transition disabled:cursor-not-allowed disabled:opacity-60"
            :class="newForm.enabled ? 'bg-indigo-600' : 'bg-slate-300 dark:bg-slate-700'"
            :disabled="disabled"
            @click="newForm.enabled = !newForm.enabled"
          >
            <span
              class="absolute top-0.5 size-5 rounded-full bg-white shadow-sm transition"
              :class="newForm.enabled ? 'left-5' : 'left-0.5'"
            />
          </button>
        </label>

        <div class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
          <div>
            <h3 class="text-sm font-semibold text-slate-950 dark:text-white">{{ $t('admin.site.configs.promotion.rules') }}</h3>
            <p class="mt-1 text-xs text-slate-500 dark:text-slate-400">{{ $t('admin.site.configs.promotion.rulesHint') }}</p>
          </div>
          <UTooltip :text="$t('admin.site.configs.promotion.addRule')" :content="{ side: 'top', sideOffset: 8 }" :delay-duration="120">
            <UButton
              type="button"
              color="neutral"
              variant="outline"
              size="xs"
              icon="i-lucide-plus"
              :aria-label="$t('admin.site.configs.promotion.addRule')"
              :disabled="disabled"
              @click="addRule"
            />
          </UTooltip>
        </div>

        <div class="space-y-3">
          <div
            v-for="(rule, ruleIndex) in newForm.rules"
            :key="rule.uid"
            class="rounded-md border border-slate-200 bg-white dark:border-slate-800 dark:bg-slate-950/40"
          >
            <div class="flex flex-col gap-3 border-b border-slate-200 px-3 py-3 sm:flex-row sm:items-start sm:justify-between dark:border-slate-800">
              <div class="min-w-0">
                <p class="text-sm font-semibold text-slate-950 dark:text-white">{{ $t('admin.site.configs.promotion.ruleTitle', { index: ruleIndex + 1 }) }}</p>
                <p class="mt-1 text-xs text-slate-500 dark:text-slate-400">{{ ruleSummary(rule) }}</p>
              </div>
              <UButton
                type="button"
                color="error"
                variant="ghost"
                size="xs"
                icon="i-lucide-trash-2"
                :disabled="disabled || newForm.rules.length <= 1"
                @click="removeRule(ruleIndex)"
              >
                {{ $t('admin.actions.delete') }}
              </UButton>
            </div>

            <div class="space-y-3 px-3 py-3">
              <div class="grid gap-3 sm:grid-cols-2">
                <label class="block">
                  <span class="text-sm font-medium text-slate-700 dark:text-slate-200">{{ $t('admin.site.configs.promotion.minGiB') }}</span>
                  <input
                    v-model="rule.minGiB"
                    type="number"
                    min="0"
                    step="any"
                    class="mt-1 h-10 w-full rounded-md border border-slate-200 bg-white px-3 text-sm text-slate-950 outline-none transition focus:border-indigo-400 focus:ring-2 focus:ring-indigo-100 disabled:cursor-not-allowed disabled:bg-slate-50 disabled:text-slate-400 dark:border-slate-700 dark:bg-slate-950 dark:text-white dark:focus:border-indigo-500 dark:focus:ring-indigo-950 dark:disabled:bg-slate-900 dark:disabled:text-slate-500"
                    :disabled="disabled"
                  >
                </label>
                <label class="block">
                  <span class="text-sm font-medium text-slate-700 dark:text-slate-200">{{ $t('admin.site.configs.promotion.durationHours') }}</span>
                  <input
                    v-model="rule.durationHours"
                    type="number"
                    min="0"
                    step="1"
                    class="mt-1 h-10 w-full rounded-md border border-slate-200 bg-white px-3 text-sm text-slate-950 outline-none transition focus:border-indigo-400 focus:ring-2 focus:ring-indigo-100 disabled:cursor-not-allowed disabled:bg-slate-50 disabled:text-slate-400 dark:border-slate-700 dark:bg-slate-950 dark:text-white dark:focus:border-indigo-500 dark:focus:ring-indigo-950 dark:disabled:bg-slate-900 dark:disabled:text-slate-500"
                    :disabled="disabled"
                  >
                  <span class="mt-1 block text-xs text-slate-500 dark:text-slate-400">{{ $t('admin.site.configs.promotion.durationHint') }}</span>
                </label>
              </div>

              <div class="rounded-md border border-slate-200 dark:border-slate-800">
                <div class="flex items-center justify-between gap-3 border-b border-slate-200 bg-slate-50 px-3 py-2 dark:border-slate-800 dark:bg-slate-900/60">
                  <div>
                    <p class="text-sm font-medium text-slate-950 dark:text-white">{{ $t('admin.site.configs.promotion.options') }}</p>
                    <p class="mt-0.5 text-xs text-slate-500 dark:text-slate-400">{{ $t('admin.site.configs.promotion.optionsHint') }}</p>
                  </div>
                  <UBadge color="neutral" variant="soft">{{ $t('admin.site.configs.promotion.totalWeight', { total: ruleTotalWeight(rule) }) }}</UBadge>
                </div>

                <div class="divide-y divide-slate-200 dark:divide-slate-800">
                  <div
                    v-for="(option, optionIndex) in rule.options"
                    :key="option.uid"
                    class="grid gap-2 px-3 py-2 sm:grid-cols-[minmax(0,1fr)_7rem_auto] sm:items-center"
                  >
                    <select
                      v-model="option.state"
                      class="h-9 w-full rounded-md border border-slate-200 bg-white px-2 text-sm text-slate-950 outline-none transition focus:border-indigo-400 focus:ring-2 focus:ring-indigo-100 disabled:cursor-not-allowed disabled:bg-slate-50 disabled:text-slate-400 dark:border-slate-700 dark:bg-slate-950 dark:text-white dark:focus:border-indigo-500 dark:focus:ring-indigo-950 dark:disabled:bg-slate-900 dark:disabled:text-slate-500"
                      :disabled="disabled"
                    >
                      <option v-for="state in promotionStateOptions" :key="state.value" :value="state.value">{{ state.label }}</option>
                    </select>
                    <input
                      v-model="option.weight"
                      type="number"
                      min="1"
                      step="1"
                      class="h-9 w-full rounded-md border border-slate-200 bg-white px-2 text-sm text-slate-950 outline-none transition focus:border-indigo-400 focus:ring-2 focus:ring-indigo-100 disabled:cursor-not-allowed disabled:bg-slate-50 disabled:text-slate-400 dark:border-slate-700 dark:bg-slate-950 dark:text-white dark:focus:border-indigo-500 dark:focus:ring-indigo-950 dark:disabled:bg-slate-900 dark:disabled:text-slate-500"
                      :disabled="disabled"
                      :aria-label="$t('admin.site.configs.promotion.weight')"
                    >
                    <UButton
                      type="button"
                      color="error"
                      variant="ghost"
                      size="xs"
                      icon="i-lucide-minus"
                      :disabled="disabled || rule.options.length <= 1"
                      :aria-label="$t('admin.site.configs.promotion.removeOption')"
                      @click="removeOption(rule, optionIndex)"
                    />
                  </div>
                </div>

                <div class="border-t border-slate-200 px-3 py-2 dark:border-slate-800">
                  <UButton type="button" color="neutral" variant="soft" size="xs" icon="i-lucide-plus" :disabled="disabled" @click="addOption(rule)">
                    {{ $t('admin.site.configs.promotion.addOption') }}
                  </UButton>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div v-else class="space-y-3">
      <div class="flex justify-end">
        <UButton type="button" color="neutral" variant="outline" size="xs" icon="i-lucide-align-left" :disabled="disabled || !jsonState.valid" @click="formatJson">
          {{ $t('admin.site.configs.promotion.formatJson') }}
        </UButton>
      </div>
      <AdminJsonEditor
        v-model="jsonText"
        :valid="jsonState.valid"
        :error-message="jsonState.message"
        :title="$t('admin.site.configs.promotion.jsonTitle')"
        :placeholder="$t('admin.site.configs.promotion.jsonPlaceholder')"
        :hint="$t('admin.site.configs.promotion.jsonHint')"
        :valid-text="$t('admin.site.configs.promotion.jsonValid')"
        :invalid-text="$t('admin.site.configs.promotion.jsonInvalidShort')"
        :fullscreen-text="$t('admin.site.configs.promotion.jsonFullscreen')"
        :close-text="$t('common.close')"
        :disabled="disabled"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
type PromotionEditorMode = 'form' | 'json'
type PromotionState = 'normal' | 'free' | '2x' | '2x_free' | '50_percent' | '2x_50_percent' | '30_percent'

interface GlobalPromotionForm {
  enabled: boolean
  state: PromotionState
  expireAt: string
}

interface PromotionOptionForm {
  uid: number
  state: PromotionState
  weight: string
}

interface PromotionRuleForm {
  uid: number
  minGiB: string
  durationHours: string
  options: PromotionOptionForm[]
}

interface NewPromotionForm {
  enabled: boolean
  rules: PromotionRuleForm[]
}

interface ValidationResult {
  valid: boolean
  value?: unknown
  message?: string
}

const props = withDefaults(defineProps<{
  modelValue: unknown
  configPath: string
  disabled?: boolean
}>(), {
  disabled: false
})

const emit = defineEmits<{
  'update:modelValue': [value: unknown]
}>()

const { t } = useI18n()

const promotionStates: PromotionState[] = ['normal', 'free', '2x', '2x_free', '50_percent', '2x_50_percent', '30_percent']
const promotionStateSet = new Set<string>(promotionStates)

const editorMode = ref<PromotionEditorMode>('form')
const jsonText = ref('')
const globalForm = reactive<GlobalPromotionForm>({
  enabled: false,
  state: 'free',
  expireAt: ''
})
const newForm = reactive<NewPromotionForm>({
  enabled: true,
  rules: []
})

let uid = 0
let resetting = false
let lastEmittedSnapshot = ''

const isGlobalPromotion = computed(() => props.configPath === 'catalog.global_promotion')
const promotionStateOptions = computed(() => promotionStates.map((value) => ({
  value,
  label: t(`admin.site.configs.promotion.states.${value}`)
})))
const jsonState = computed(() => {
  if (editorMode.value !== 'json') {
    return { valid: true, message: '' }
  }
  try {
    validateConfigValue(JSON.parse(getJsonText() || 'null'))
    return { valid: true, message: '' }
  } catch (error: unknown) {
    return { valid: false, message: errorMessage(error) }
  }
})

watch(
  () => [props.configPath, props.modelValue] as const,
  () => {
    const nextSnapshot = modelSnapshot(props.configPath, props.modelValue)
    if (nextSnapshot === lastEmittedSnapshot) return
    resetFromModelValue()
  },
  { deep: true, immediate: true }
)

watch(
  [editorMode, jsonText, globalForm, newForm],
  () => {
    if (resetting) return
    emitCurrentValue()
  },
  { deep: true }
)

function setEditorMode(mode: PromotionEditorMode) {
  if (mode === editorMode.value) return
  if (mode === 'json') {
    jsonText.value = JSON.stringify(snapshotValue(), null, 2)
    editorMode.value = 'json'
    return
  }

  try {
    const value = JSON.parse(getJsonText() || 'null')
    validateConfigValue(value)
    resetFormsFromValue(value)
    editorMode.value = 'form'
  } catch {
    // Keep JSON mode so the visible validation message can guide the user.
  }
}

function formatJson() {
  if (!jsonState.value.valid) return
  jsonText.value = JSON.stringify(JSON.parse(getJsonText() || 'null'), null, 2)
}

function addRule() {
  newForm.rules.push(createRule({
    minGiB: 0,
    durationHours: 72,
    options: [{ state: 'normal', weight: 100 }]
  }))
}

function removeRule(index: number) {
  if (newForm.rules.length <= 1) return
  newForm.rules.splice(index, 1)
}

function addOption(rule: PromotionRuleForm) {
  rule.options.push(createOption({ state: 'free', weight: 10 }))
}

function removeOption(rule: PromotionRuleForm, index: number) {
  if (rule.options.length <= 1) return
  rule.options.splice(index, 1)
}

function ruleTotalWeight(rule: PromotionRuleForm) {
  return rule.options.reduce((total, option) => {
    const weight = Number(option.weight)
    return Number.isFinite(weight) && weight > 0 ? total + weight : total
  }, 0)
}

function ruleSummary(rule: PromotionRuleForm) {
  const minGiB = displayNumber(rule.minGiB)
  const duration = displayNumber(rule.durationHours)
  const durationLabel = Number(duration) === 0
    ? t('admin.site.configs.promotion.noExpire')
    : t('admin.site.configs.promotion.hours', { hours: duration })
  return t('admin.site.configs.promotion.ruleSummary', {
    minGiB,
    duration: durationLabel,
    total: ruleTotalWeight(rule)
  })
}

function validate(): ValidationResult {
  try {
    return { valid: true, value: buildSubmitValue() }
  } catch (error: unknown) {
    return { valid: false, message: errorMessage(error) }
  }
}

defineExpose({ validate })

function resetFromModelValue() {
  resetting = true
  editorMode.value = 'form'
  jsonText.value = formatEditableValue(props.modelValue)
  resetFormsFromValue(props.modelValue)
  resetting = false
}

function resetFormsFromValue(value: unknown) {
  if (isGlobalPromotion.value) {
    resetGlobalForm(value)
    return
  }
  resetNewForm(value)
}

function resetGlobalForm(value: unknown) {
  const record = asRecord(value)
  globalForm.enabled = Boolean(record.enabled)
  globalForm.state = normalizeState(record.state, 'free')
  globalForm.expireAt = apiDateToLocalInput(record.expireAt)
}

function resetNewForm(value: unknown) {
  const record = asRecord(value)
  newForm.enabled = record.enabled === undefined ? true : Boolean(record.enabled)
  const rules = Array.isArray(record.rules) ? record.rules : []
  newForm.rules = rules.length > 0
    ? rules.map((rule) => createRule(rule))
    : defaultRules()
}

function defaultRules(): PromotionRuleForm[] {
  return [
    createRule({
      minGiB: 0,
      durationHours: 72,
      options: [
        { state: 'normal', weight: 70 },
        { state: 'free', weight: 10 },
        { state: '2x', weight: 10 },
        { state: '50_percent', weight: 10 }
      ]
    }),
    createRule({
      minGiB: 10,
      durationHours: 96,
      options: [
        { state: 'normal', weight: 50 },
        { state: 'free', weight: 20 },
        { state: '2x', weight: 20 },
        { state: '2x_free', weight: 10 }
      ]
    }),
    createRule({
      minGiB: 50,
      durationHours: 168,
      options: [
        { state: 'free', weight: 40 },
        { state: '2x_free', weight: 30 },
        { state: '2x_50_percent', weight: 30 }
      ]
    })
  ]
}

function createRule(value: unknown = {}): PromotionRuleForm {
  const record = asRecord(value)
  const options = Array.isArray(record.options) ? record.options : []
  return {
    uid: nextUid(),
    minGiB: formatNumberInput(record.minGiB, '0'),
    durationHours: formatNumberInput(record.durationHours, '72'),
    options: options.length > 0 ? options.map((option) => createOption(option)) : [createOption({ state: 'normal', weight: 100 })]
  }
}

function createOption(value: unknown = {}): PromotionOptionForm {
  const record = asRecord(value)
  return {
    uid: nextUid(),
    state: normalizeState(record.state, 'normal'),
    weight: formatNumberInput(record.weight, '100')
  }
}

function emitCurrentValue() {
  const value = snapshotValue()
  lastEmittedSnapshot = modelSnapshot(props.configPath, value)
  emit('update:modelValue', value)
}

function snapshotValue() {
  if (editorMode.value === 'json') {
    try {
      return JSON.parse(getJsonText() || 'null')
    } catch {
      return { invalidJson: getJsonText() }
    }
  }
  return isGlobalPromotion.value ? buildGlobalValue(false) : buildNewValue(false)
}

function buildSubmitValue() {
  if (editorMode.value === 'json') {
    const value = JSON.parse(getJsonText() || 'null')
    validateConfigValue(value)
    return value
  }
  return isGlobalPromotion.value ? buildGlobalValue(true) : buildNewValue(true)
}

function buildGlobalValue(strict: boolean) {
  const state = normalizeState(globalForm.state, 'free')
  if (strict && globalForm.enabled && !promotionStateSet.has(state)) {
    throw new Error(t('admin.site.configs.promotion.errors.state'))
  }
  return {
    enabled: globalForm.enabled,
    state,
    expireAt: localInputToApiDate(globalForm.expireAt)
  }
}

function buildNewValue(strict: boolean) {
  const rules = newForm.rules.map((rule, ruleIndex) => ({
    minGiB: parseNumber(rule.minGiB, {
      strict,
      min: 0,
      field: t('admin.site.configs.promotion.minGiB'),
      index: ruleIndex + 1
    }),
    durationHours: parseNumber(rule.durationHours, {
      strict,
      min: 0,
      integer: true,
      field: t('admin.site.configs.promotion.durationHours'),
      index: ruleIndex + 1
    }),
    options: rule.options.map((option, optionIndex) => ({
      state: normalizeState(option.state, 'normal'),
      weight: parseNumber(option.weight, {
        strict,
        min: 1,
        integer: true,
        field: t('admin.site.configs.promotion.weight'),
        index: optionIndex + 1
      })
    }))
  }))

  if (strict && newForm.enabled && rules.length === 0) {
    throw new Error(t('admin.site.configs.promotion.errors.rules'))
  }
  if (strict) {
    for (const [ruleIndex, rule] of rules.entries()) {
      if (rule.options.length === 0) {
        throw new Error(t('admin.site.configs.promotion.errors.options', { index: ruleIndex + 1 }))
      }
      for (const option of rule.options) {
        if (!promotionStateSet.has(option.state)) {
          throw new Error(t('admin.site.configs.promotion.errors.state'))
        }
      }
    }
  }

  return {
    enabled: newForm.enabled,
    rules
  }
}

function validateConfigValue(value: unknown) {
  if (isGlobalPromotion.value) {
    validateGlobalValue(value)
    return
  }
  validateNewValue(value)
}

function validateGlobalValue(value: unknown) {
  const record = requireRecord(value)
  if (typeof record.enabled !== 'boolean') {
    throw new Error(t('admin.site.configs.promotion.errors.enabled'))
  }
  if (typeof record.state !== 'string' || !promotionStateSet.has(record.state)) {
    throw new Error(t('admin.site.configs.promotion.errors.state'))
  }
  if (record.expireAt !== undefined && record.expireAt !== null) {
    if (typeof record.expireAt !== 'string') {
      throw new Error(t('admin.site.configs.promotion.errors.expireAt'))
    }
    const expireAt = record.expireAt.trim()
    if (expireAt && !isApiDateTime(expireAt)) {
      throw new Error(t('admin.site.configs.promotion.errors.expireAt'))
    }
  }
}

function validateNewValue(value: unknown) {
  const record = requireRecord(value)
  if (typeof record.enabled !== 'boolean') {
    throw new Error(t('admin.site.configs.promotion.errors.enabled'))
  }
  if (!Array.isArray(record.rules)) {
    throw new Error(t('admin.site.configs.promotion.errors.rules'))
  }
  if (record.enabled && record.rules.length === 0) {
    throw new Error(t('admin.site.configs.promotion.errors.rules'))
  }

  record.rules.forEach((ruleValue, ruleIndex) => {
    const rule = requireRecord(ruleValue)
    assertFiniteNumber(rule.minGiB, { min: 0, field: t('admin.site.configs.promotion.minGiB'), index: ruleIndex + 1 })
    assertFiniteNumber(rule.durationHours, { min: 0, integer: true, field: t('admin.site.configs.promotion.durationHours'), index: ruleIndex + 1 })
    if (!Array.isArray(rule.options) || rule.options.length === 0) {
      throw new Error(t('admin.site.configs.promotion.errors.options', { index: ruleIndex + 1 }))
    }
    rule.options.forEach((optionValue, optionIndex) => {
      const option = requireRecord(optionValue)
      if (typeof option.state !== 'string' || !promotionStateSet.has(option.state)) {
        throw new Error(t('admin.site.configs.promotion.errors.state'))
      }
      assertFiniteNumber(option.weight, {
        min: 1,
        integer: true,
        field: t('admin.site.configs.promotion.weight'),
        index: optionIndex + 1
      })
    })
  })
}

function parseNumber(value: string, options: { strict: boolean, min: number, integer?: boolean, field: string, index: number }) {
  const text = String(value ?? '').trim()
  if (!text) {
    if (!options.strict) return ''
    throw new Error(numberError(options))
  }
  const parsed = Number(text)
  if (!options.strict) {
    return Number.isFinite(parsed) && (!options.integer || Number.isInteger(parsed)) ? parsed : text
  }
  assertFiniteNumber(parsed, options)
  return parsed
}

function assertFiniteNumber(value: unknown, options: { min: number, integer?: boolean, field: string, index: number }) {
  const parsed = typeof value === 'number' ? value : typeof value === 'string' && value.trim() !== '' ? Number(value) : Number.NaN
  if (!Number.isFinite(parsed) || parsed < options.min || (options.integer && !Number.isInteger(parsed))) {
    throw new Error(numberError(options))
  }
}

function numberError(options: { min: number, field: string, index: number }) {
  return t('admin.site.configs.promotion.errors.number', {
    field: options.field,
    index: options.index,
    min: options.min
  })
}

function asRecord(value: unknown): Record<string, unknown> {
  return value && typeof value === 'object' && !Array.isArray(value) ? value as Record<string, unknown> : {}
}

function requireRecord(value: unknown): Record<string, unknown> {
  if (!value || typeof value !== 'object' || Array.isArray(value)) {
    throw new Error(t('admin.site.configs.promotion.errors.object'))
  }
  return value as Record<string, unknown>
}

function normalizeState(value: unknown, fallback: PromotionState): PromotionState {
  const text = String(value ?? '').trim()
  return promotionStateSet.has(text) ? text as PromotionState : fallback
}

function formatNumberInput(value: unknown, fallback: string) {
  if (typeof value === 'number' && Number.isFinite(value)) return String(value)
  if (typeof value === 'string' && value.trim() !== '') return value.trim()
  return fallback
}

function displayNumber(value: string) {
  const parsed = Number(value)
  if (!Number.isFinite(parsed)) return value || '-'
  return Number.isInteger(parsed) ? String(parsed) : parsed.toFixed(2).replace(/\.?0+$/, '')
}

function apiDateToLocalInput(value: unknown) {
  const text = String(value ?? '').trim()
  if (!text) return ''
  return text.replace(' ', 'T').slice(0, 19)
}

function localInputToApiDate(value: string) {
  const text = String(value ?? '').trim()
  if (!text) return ''
  const normalized = text.replace('T', ' ')
  return normalized.length === 16 ? `${normalized}:00` : normalized
}

function isApiDateTime(value: string) {
  return /^\d{4}-\d{2}-\d{2} \d{2}:\d{2}(:\d{2})?$/.test(value)
}

function formatEditableValue(value: unknown) {
  if (value === null || value === undefined) return ''
  if (typeof value === 'string') return value
  if (typeof value === 'number') return String(value)
  return JSON.stringify(value, null, 2)
}

function getJsonText() {
  return String(jsonText.value ?? '').trim()
}

function stableStringify(value: unknown) {
  return JSON.stringify(value)
}

function modelSnapshot(configPath: string, value: unknown) {
  return `${configPath}:${stableStringify(value)}`
}

function errorMessage(error: unknown) {
  return error instanceof Error ? error.message : t('admin.site.configs.form.jsonInvalid')
}

function nextUid() {
  uid += 1
  return uid
}
</script>
