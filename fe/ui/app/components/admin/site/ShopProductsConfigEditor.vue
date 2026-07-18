<template>
  <div class="space-y-4">
    <div class="flex flex-col gap-3 rounded-md border border-slate-200 bg-slate-50 px-3 py-3 sm:flex-row sm:items-center sm:justify-between dark:border-slate-800 dark:bg-slate-950/70">
      <div>
        <p class="text-sm font-medium text-slate-950 dark:text-white">{{ $t('admin.site.configs.shop.editMode') }}</p>
        <p class="mt-1 text-xs text-slate-500 dark:text-slate-400">{{ $t('admin.site.configs.shop.editModeHint') }}</p>
      </div>
      <div class="inline-flex h-9 rounded-md border border-slate-200 bg-white p-0.5 dark:border-slate-700 dark:bg-slate-950">
        <button
          v-for="mode in editorModes"
          :key="mode.value"
          type="button"
          class="inline-flex min-w-20 items-center justify-center gap-1.5 rounded px-3 text-sm font-medium transition"
          :class="editorMode === mode.value ? 'bg-indigo-600 text-white shadow-sm' : 'text-slate-600 hover:bg-slate-100 dark:text-slate-300 dark:hover:bg-slate-800'"
          :disabled="disabled"
          @click="setEditorMode(mode.value)"
        >
          <UIcon :name="mode.icon" class="size-4" />
          {{ mode.label }}
        </button>
      </div>
    </div>

    <div v-if="editorMode === 'form'" class="space-y-3">
      <div class="flex flex-col gap-2 sm:flex-row sm:items-center sm:justify-end">
        <USelect
          v-model="newProductType"
          class="w-full sm:w-44"
          size="lg"
          :items="productTypeOptions"
          value-key="value"
          :disabled="disabled"
          :aria-label="$t('admin.site.configs.shop.newProductType')"
        />
        <UButton type="button" color="neutral" variant="outline" size="lg" icon="i-lucide-plus" :disabled="disabled" @click="addProduct">
          {{ $t('admin.site.configs.shop.addProduct') }}
        </UButton>
      </div>

      <div v-if="products.length === 0" class="rounded-md border border-dashed border-slate-300 px-4 py-10 text-center text-sm text-slate-500 dark:border-slate-700 dark:text-slate-400">
        {{ $t('admin.site.configs.shop.empty') }}
      </div>

      <div v-for="product in products" :key="product.key" class="overflow-hidden rounded-md border border-slate-200 dark:border-slate-800">
        <div class="flex items-center justify-between gap-3 border-b border-slate-200 bg-slate-50 px-3 py-2.5 dark:border-slate-800 dark:bg-slate-950/60">
          <div class="min-w-0">
            <p class="truncate text-sm font-semibold text-slate-950 dark:text-white">{{ productTitle(product.type, product.key) }}</p>
            <p class="mt-0.5 truncate text-xs text-slate-500 dark:text-slate-400">{{ product.key }} / {{ product.type }}</p>
          </div>
          <USwitch v-model="product.enabled" :disabled="disabled" />
        </div>

        <div class="grid gap-3 p-3 sm:grid-cols-2">
          <UFormField :label="$t('admin.site.configs.shop.price')">
            <UInput v-model="product.price" type="number" min="0.01" step="0.01" class="w-full" size="lg" :disabled="disabled" />
          </UFormField>
          <UFormField :label="$t('admin.site.configs.shop.sortOrder')">
            <UInput v-model="product.sortOrder" type="number" step="1" class="w-full" size="lg" :disabled="disabled" />
          </UFormField>
          <UFormField v-if="product.type === 'invite'" :label="$t('admin.site.configs.shop.inviteAmount')" class="sm:col-span-2">
            <UInput v-model="product.amount" type="number" min="1" step="1" class="w-full" size="lg" :disabled="disabled" />
          </UFormField>
          <UFormField v-if="product.type === 'vip'" :label="$t('admin.site.configs.shop.vipDurationDays')" class="sm:col-span-2">
            <UInput v-model="product.durationDays" type="number" min="1" step="1" class="w-full" size="lg" :disabled="disabled" />
          </UFormField>
          <UFormField v-if="product.type === 'upload' || product.type === 'download'" :label="$t('admin.site.configs.shop.trafficAmountGiB')" class="sm:col-span-2">
            <UInput v-model="product.amountGiB" type="number" min="1" step="1" class="w-full" size="lg" :disabled="disabled" />
          </UFormField>
        </div>
      </div>

      <p class="text-xs leading-5 text-slate-500 dark:text-slate-400">{{ $t('admin.site.configs.shop.structureHint') }}</p>
    </div>

    <div v-else class="space-y-3">
      <div class="flex justify-end">
        <UButton type="button" color="neutral" variant="outline" size="xs" icon="i-lucide-align-left" :disabled="disabled || !jsonState.valid" @click="formatJson">
          {{ $t('admin.site.configs.shop.formatJson') }}
        </UButton>
      </div>
      <AdminJsonEditor
        v-model="jsonText"
        :valid="jsonState.valid"
        :error-message="jsonState.message"
        :title="$t('admin.site.configs.shop.jsonTitle')"
        :placeholder="$t('admin.site.configs.shop.jsonPlaceholder')"
        :hint="$t('admin.site.configs.shop.jsonHint')"
        :valid-text="$t('admin.site.configs.shop.jsonValid')"
        :invalid-text="$t('admin.site.configs.shop.jsonInvalidShort')"
        :fullscreen-text="$t('admin.site.configs.shop.jsonFullscreen')"
        :close-text="$t('common.close')"
        :disabled="disabled"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
type EditorMode = 'form' | 'json'

interface ProductForm {
  key: string
  type: string
  enabled: boolean
  price: string
  sortOrder: string
  amount: string
  durationDays: string
  amountGiB: string
  options: Record<string, unknown>
}

interface ValidationResult {
  valid: boolean
  value?: unknown
  message?: string
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
const productTypes = ['invite', 'vip', 'upload', 'download']
const productTypeSet = new Set(productTypes)
const editorMode = ref<EditorMode>('form')
const products = ref<ProductForm[]>([])
const jsonText = ref('')
const newProductType = ref('invite')
let resetting = false
let lastEmittedSnapshot = ''

const editorModes = computed(() => [
  { value: 'form' as const, label: t('admin.site.configs.shop.modeForm'), icon: 'i-lucide-sliders-horizontal' },
  { value: 'json' as const, label: t('admin.site.configs.shop.modeJson'), icon: 'i-lucide-braces' }
])
const productTypeOptions = computed(() => productTypes.map(type => ({
  value: type,
  label: productTitle(type, type)
})))
const jsonState = computed(() => {
  if (editorMode.value !== 'json') return { valid: true, message: '' }
  try {
    validateProducts(JSON.parse(jsonText.value || 'null'))
    return { valid: true, message: '' }
  } catch (error: unknown) {
    return { valid: false, message: errorMessage(error) }
  }
})

watch(() => props.modelValue, resetFromModelValue, { deep: true, immediate: true })
watch(products, emitFormValue, { deep: true })
watch(jsonText, emitJsonValue)

function resetFromModelValue() {
  const snapshot = JSON.stringify(props.modelValue)
  if (snapshot === lastEmittedSnapshot) return
  resetting = true
  jsonText.value = JSON.stringify(props.modelValue ?? [], null, 2)
  try {
    products.value = toProductForms(validateProducts(props.modelValue))
  } catch {
    products.value = []
    editorMode.value = 'json'
  }
  nextTick(() => { resetting = false })
}

function setEditorMode(mode: EditorMode) {
  if (mode === editorMode.value) return
  if (mode === 'json') {
    jsonText.value = JSON.stringify(buildProductsValue(), null, 2)
    editorMode.value = mode
    return
  }
  try {
    products.value = toProductForms(validateProducts(JSON.parse(jsonText.value || 'null')))
    editorMode.value = mode
  } catch {
    // Keep the JSON editor visible so the invalid value can be corrected.
  }
}

function emitFormValue() {
  if (resetting || editorMode.value !== 'form') return
  emitValue(buildProductsValue())
}

function emitJsonValue() {
  if (resetting || editorMode.value !== 'json') return
  try {
    const value = validateProducts(JSON.parse(jsonText.value || 'null'))
    emitValue(value)
  } catch {
    // Invalid JSON is surfaced by jsonState and is not emitted.
  }
}

function emitValue(value: unknown) {
  lastEmittedSnapshot = JSON.stringify(value)
  emit('update:modelValue', value)
}

function addProduct() {
  const type = newProductType.value
  const baseKey = {
    invite: 'invite',
    vip: 'vip_30d',
    upload: 'upload_100_gib',
    download: 'download_100_gib'
  }[type] || type
  const keys = new Set(products.value.map(product => product.key))
  let key = baseKey
  let suffix = 2
  while (keys.has(key)) {
    key = `${type}_${suffix}`
    suffix++
  }

  products.value.push({
    key,
    type,
    enabled: true,
    price: '1000',
    sortOrder: String(Math.max(0, ...products.value.map(product => Number(product.sortOrder) || 0)) + 10),
    amount: '1',
    durationDays: '30',
    amountGiB: '100',
    options: {}
  })
}

function buildProductsValue() {
  return products.value.map(product => ({
    key: product.key,
    type: product.type,
    enabled: product.enabled,
    price: Number(product.price),
    sortOrder: Number(product.sortOrder),
    options: buildProductOptions(product)
  }))
}

function buildProductOptions(product: ProductForm) {
  const options = { ...product.options }
  if (product.type === 'invite') options.amount = Number(product.amount)
  if (product.type === 'vip') options.durationDays = Number(product.durationDays)
  if (product.type === 'upload' || product.type === 'download') options.amountGiB = Number(product.amountGiB)
  return options
}

function validateProducts(value: unknown) {
  if (!Array.isArray(value)) throw new Error(t('admin.site.configs.shop.errors.array'))
  const keys = new Set<string>()
  return value.map((item, index) => {
    if (!item || typeof item !== 'object' || Array.isArray(item)) throw new Error(t('admin.site.configs.shop.errors.object', { index: index + 1 }))
    const product = item as Record<string, unknown>
    const key = String(product.key || '').trim()
    const type = String(product.type || '').trim()
    const price = Number(product.price)
    const sortOrder = Number(product.sortOrder ?? 0)
    const options = product.options && typeof product.options === 'object' && !Array.isArray(product.options)
      ? product.options as Record<string, unknown>
      : {}

    if (!key) throw new Error(t('admin.site.configs.shop.errors.key', { index: index + 1 }))
    if (keys.has(key)) throw new Error(t('admin.site.configs.shop.errors.duplicateKey', { key }))
    keys.add(key)
    if (!productTypeSet.has(type)) throw new Error(t('admin.site.configs.shop.errors.type', { index: index + 1 }))
    if (!Number.isFinite(price) || price <= 0) throw new Error(t('admin.site.configs.shop.errors.price', { index: index + 1 }))
    if (!Number.isInteger(sortOrder)) throw new Error(t('admin.site.configs.shop.errors.sortOrder', { index: index + 1 }))
    if (type === 'invite' && (!Number.isInteger(Number(options.amount)) || Number(options.amount) <= 0)) {
      throw new Error(t('admin.site.configs.shop.errors.inviteAmount', { index: index + 1 }))
    }
    if (type === 'vip' && (!Number.isInteger(Number(options.durationDays)) || Number(options.durationDays) <= 0)) {
      throw new Error(t('admin.site.configs.shop.errors.vipDurationDays', { index: index + 1 }))
    }
    if ((type === 'upload' || type === 'download') && (!Number.isInteger(Number(options.amountGiB)) || Number(options.amountGiB) <= 0)) {
      throw new Error(t('admin.site.configs.shop.errors.trafficAmountGiB', { index: index + 1 }))
    }
    return { ...product, key, type, enabled: Boolean(product.enabled), price, sortOrder, options }
  })
}

function toProductForms(value: Array<Record<string, unknown>>) {
  return value.map(product => {
    const options = product.options as Record<string, unknown>
    return {
      key: String(product.key),
      type: String(product.type),
      enabled: Boolean(product.enabled),
      price: String(product.price),
      sortOrder: String(product.sortOrder),
      amount: String(options.amount ?? 1),
      durationDays: String(options.durationDays ?? 30),
      amountGiB: String(options.amountGiB ?? 100),
      options: { ...options }
    }
  })
}

function formatJson() {
  if (!jsonState.value.valid) return
  jsonText.value = JSON.stringify(JSON.parse(jsonText.value), null, 2)
}

function productTitle(type: string, fallback: string) {
  const key = `user.shop.products.${type}.title`
  const translated = t(key)
  return translated === key ? fallback : translated
}

function errorMessage(error: unknown) {
  return error instanceof Error ? error.message : t('admin.site.configs.form.jsonInvalid')
}

function validate(): ValidationResult {
  try {
    const value = editorMode.value === 'json'
      ? validateProducts(JSON.parse(jsonText.value || 'null'))
      : validateProducts(buildProductsValue())
    return { valid: true, value }
  } catch (error: unknown) {
    return { valid: false, message: errorMessage(error) }
  }
}

defineExpose({ validate })
</script>
