<template>
  <div class="min-h-[calc(100vh-4rem)] bg-slate-50 py-6 dark:bg-slate-950">
    <div class="w-full px-3 sm:px-4 lg:px-5">
      <div class="mb-4 flex flex-wrap items-center justify-end gap-2">
        <UButton
          v-for="group in groups"
          :key="group.value"
          :color="selectedGroup === group.value ? 'primary' : 'neutral'"
          :variant="selectedGroup === group.value ? 'soft' : 'outline'"
          icon="i-lucide-folder-cog"
          @click="selectGroup(group.value)"
        >
          {{ group.label }}
        </UButton>
      </div>

      <div class="mb-4 grid grid-cols-2 gap-2 sm:flex sm:items-center">
        <div class="rounded-lg border border-slate-200 bg-white px-3 py-2 dark:border-slate-800 dark:bg-slate-900">
          <p class="text-xs text-slate-500 dark:text-slate-400">{{ $t('admin.site.configs.stats.group') }}</p>
          <p class="mt-1 text-lg font-semibold text-slate-950 dark:text-white">{{ currentGroupLabel }}</p>
        </div>
        <div class="rounded-lg border border-slate-200 bg-white px-3 py-2 dark:border-slate-800 dark:bg-slate-900">
          <p class="text-xs text-slate-500 dark:text-slate-400">{{ $t('admin.site.configs.stats.total') }}</p>
          <p class="mt-1 text-lg font-semibold text-indigo-700 dark:text-indigo-300">{{ numberFormatter.format(configs.length) }}</p>
        </div>
      </div>

      <div class="grid gap-4 xl:grid-cols-[minmax(0,1fr)_420px]">
        <section class="overflow-hidden rounded-lg border border-slate-200 bg-white dark:border-slate-800 dark:bg-slate-900">
          <div class="border-b border-slate-200 px-4 py-3 dark:border-slate-800">
            <div class="flex items-center gap-2">
              <UIcon name="i-lucide-settings-2" class="size-5 text-indigo-600 dark:text-indigo-300" />
              <h2 class="text-sm font-semibold text-slate-950 dark:text-white">{{ $t('admin.site.configs.list') }}</h2>
            </div>
          </div>

          <div v-if="pending" class="overflow-x-auto">
            <table class="min-w-[920px] w-full table-fixed border-collapse text-left">
              <thead class="bg-slate-50 text-xs font-medium uppercase text-slate-500 dark:bg-slate-950/70 dark:text-slate-400">
                <tr>
                  <th class="w-[22%] border-b border-slate-200 px-4 py-3 dark:border-slate-800">{{ $t('admin.site.configs.table.key') }}</th>
                  <th class="w-[26%] border-b border-slate-200 px-4 py-3 dark:border-slate-800">{{ $t('admin.site.configs.table.value') }}</th>
                  <th class="w-[28%] border-b border-slate-200 px-4 py-3 dark:border-slate-800">{{ $t('admin.site.configs.table.description') }}</th>
                  <th class="w-[14%] border-b border-slate-200 px-4 py-3 text-right dark:border-slate-800">{{ $t('admin.site.configs.table.updatedAt') }}</th>
                  <th class="w-[10%] border-b border-slate-200 px-4 py-3 text-right dark:border-slate-800">{{ $t('admin.site.configs.table.actions') }}</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="index in 6" :key="index" class="border-b border-slate-200 last:border-b-0 dark:border-slate-800">
                  <td class="px-4 py-4"><div class="h-4 w-40 animate-pulse rounded bg-slate-200 dark:bg-slate-800" /></td>
                  <td class="px-4 py-4"><div class="h-4 w-52 animate-pulse rounded bg-slate-100 dark:bg-slate-800/70" /></td>
                  <td class="px-4 py-4"><div class="h-4 w-56 animate-pulse rounded bg-slate-100 dark:bg-slate-800/70" /></td>
                  <td class="px-4 py-4"><div class="ml-auto h-4 w-28 animate-pulse rounded bg-slate-100 dark:bg-slate-800/70" /></td>
                  <td class="px-4 py-4"><div class="ml-auto h-4 w-16 animate-pulse rounded bg-slate-100 dark:bg-slate-800/70" /></td>
                </tr>
              </tbody>
            </table>
          </div>

          <div v-else-if="errorMessage" class="flex flex-col items-center justify-center px-4 py-16 text-center">
            <UIcon name="i-lucide-circle-alert" class="size-9 text-red-500" />
            <p class="mt-3 text-sm font-medium text-slate-950 dark:text-white">{{ errorMessage }}</p>
          </div>

          <div v-else-if="configs.length === 0" class="flex flex-col items-center justify-center px-4 py-16 text-center">
            <UIcon name="i-lucide-inbox" class="size-9 text-slate-400" />
            <p class="mt-3 text-sm font-medium text-slate-950 dark:text-white">{{ $t('admin.site.configs.empty') }}</p>
          </div>

          <div v-else class="overflow-x-auto">
            <table class="min-w-[920px] w-full table-fixed border-collapse text-left">
              <thead class="bg-slate-50 text-xs font-medium uppercase text-slate-500 dark:bg-slate-950/70 dark:text-slate-400">
                <tr>
                  <th class="w-[22%] border-b border-slate-200 px-4 py-3 dark:border-slate-800">{{ $t('admin.site.configs.table.key') }}</th>
                  <th class="w-[26%] border-b border-slate-200 px-4 py-3 dark:border-slate-800">{{ $t('admin.site.configs.table.value') }}</th>
                  <th class="w-[28%] border-b border-slate-200 px-4 py-3 dark:border-slate-800">{{ $t('admin.site.configs.table.description') }}</th>
                  <th class="w-[14%] border-b border-slate-200 px-4 py-3 text-right dark:border-slate-800">{{ $t('admin.site.configs.table.updatedAt') }}</th>
                  <th class="w-[10%] border-b border-slate-200 px-4 py-3 text-right dark:border-slate-800">{{ $t('admin.site.configs.table.actions') }}</th>
                </tr>
              </thead>
              <tbody>
                <tr
                  v-for="config in configs"
                  :key="config.id || `${config.group}.${config.key}`"
                  class="cursor-pointer border-b border-slate-200 transition-colors last:border-b-0 dark:border-slate-800"
                  :class="selectedConfig?.group === config.group && selectedConfig?.key === config.key ? 'bg-indigo-50/70 dark:bg-indigo-950/30' : 'hover:bg-slate-50 dark:hover:bg-slate-950/70'"
                  @click="selectConfig(config)"
                >
                  <td class="px-4 py-3 align-middle">
                    <code class="block truncate text-sm font-semibold text-slate-950 dark:text-white">
                      {{ config.key }}
                    </code>
                    <span class="mt-1 block truncate text-xs text-slate-500 dark:text-slate-400">{{ config.group }}</span>
                  </td>
                  <td class="px-4 py-3 align-middle">
                    <code class="block truncate rounded bg-slate-100 px-2 py-1 text-xs text-slate-700 dark:bg-slate-800 dark:text-slate-200">
                      {{ previewConfigValue(config) }}
                    </code>
                  </td>
                  <td class="px-4 py-3 align-middle">
                    <span class="block truncate text-sm text-slate-600 dark:text-slate-300">{{ config.description || '-' }}</span>
                  </td>
                  <td class="px-4 py-3 text-right align-middle text-sm text-slate-600 dark:text-slate-300">
                    {{ formatDateTime(config.updatedAt || config.createdAt, locale) }}
                  </td>
                  <td class="px-4 py-3 align-middle">
                    <div class="flex justify-end">
                      <UButton color="neutral" variant="outline" size="sm" icon="i-lucide-pencil" @click.stop="selectConfig(config)">
                        {{ $t('admin.actions.edit') }}
                      </UButton>
                    </div>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </section>

        <section class="rounded-lg border border-slate-200 bg-white dark:border-slate-800 dark:bg-slate-900">
          <div class="border-b border-slate-200 px-4 py-3 dark:border-slate-800">
            <h2 class="text-sm font-semibold text-slate-950 dark:text-white">
              {{ selectedConfig ? selectedConfig.key : $t('admin.site.configs.form.empty') }}
            </h2>
          </div>

          <form class="space-y-4 p-4" @submit.prevent="saveConfig">
            <div class="grid grid-cols-2 gap-3">
              <label class="block">
                <span class="text-sm font-medium text-slate-700 dark:text-slate-200">{{ $t('admin.site.configs.form.group') }}</span>
                <input :value="selectedConfig?.group || selectedGroup" class="mt-1 h-10 w-full rounded-md border border-slate-200 bg-slate-50 px-3 text-sm text-slate-600 outline-none dark:border-slate-700 dark:bg-slate-950 dark:text-slate-300" disabled>
              </label>

              <label class="block">
                <span class="text-sm font-medium text-slate-700 dark:text-slate-200">{{ $t('admin.site.configs.form.type') }}</span>
                <input :value="selectedKindLabel" class="mt-1 h-10 w-full rounded-md border border-slate-200 bg-slate-50 px-3 text-sm text-slate-600 outline-none dark:border-slate-700 dark:bg-slate-950 dark:text-slate-300" disabled>
              </label>
            </div>

            <label class="block">
              <span class="text-sm font-medium text-slate-700 dark:text-slate-200">{{ $t('admin.site.configs.form.key') }}</span>
              <input :value="selectedConfig?.key || ''" class="mt-1 h-10 w-full rounded-md border border-slate-200 bg-slate-50 px-3 font-mono text-sm text-slate-600 outline-none dark:border-slate-700 dark:bg-slate-950 dark:text-slate-300" disabled>
            </label>

            <label v-if="selectedKind === 'boolean'" class="flex items-center gap-2 text-sm font-medium text-slate-700 dark:text-slate-200">
              <input v-model="form.booleanValue" type="checkbox" class="size-4 rounded border-slate-300 text-indigo-600 focus:ring-indigo-500 dark:border-slate-600" :disabled="!selectedConfig || saving">
              <span>{{ $t('admin.site.configs.form.booleanValue') }}</span>
            </label>

            <label v-else-if="selectedKind === 'int' || selectedKind === 'float'" class="block">
              <span class="text-sm font-medium text-slate-700 dark:text-slate-200">{{ $t('admin.site.configs.form.value') }}</span>
              <input v-model="form.textValue" type="number" :step="selectedKind === 'int' ? 1 : 'any'" class="mt-1 h-10 w-full rounded-md border border-slate-200 bg-white px-3 text-sm text-slate-950 outline-none transition focus:border-indigo-400 focus:ring-2 focus:ring-indigo-100 dark:border-slate-700 dark:bg-slate-950 dark:text-white dark:focus:border-indigo-500 dark:focus:ring-indigo-950" :disabled="!selectedConfig || saving">
            </label>

            <label v-else class="block">
              <span class="text-sm font-medium text-slate-700 dark:text-slate-200">{{ $t('admin.site.configs.form.value') }}</span>
              <textarea v-model="form.textValue" :rows="selectedKind === 'json' ? 8 : 4" class="mt-1 w-full resize-none rounded-md border border-slate-200 bg-white px-3 py-2 font-mono text-sm text-slate-950 outline-none transition focus:border-indigo-400 focus:ring-2 focus:ring-indigo-100 dark:border-slate-700 dark:bg-slate-950 dark:text-white dark:focus:border-indigo-500 dark:focus:ring-indigo-950" :disabled="!selectedConfig || saving" />
            </label>

            <div v-if="selectedConfig?.description" class="rounded-md border border-slate-200 bg-slate-50 px-3 py-2 text-sm text-slate-600 dark:border-slate-800 dark:bg-slate-950/70 dark:text-slate-300">
              {{ selectedConfig.description }}
            </div>

            <p v-if="formError" class="rounded-md border border-red-200 bg-red-50 px-3 py-2 text-sm text-red-700 dark:border-red-900 dark:bg-red-950/40 dark:text-red-200">
              {{ formError }}
            </p>

            <div class="flex flex-col gap-2 sm:flex-row">
              <button
                type="button"
                class="inline-flex h-10 items-center justify-center gap-2 rounded-md bg-primary px-4 text-sm font-medium text-white transition hover:bg-primary/90 disabled:cursor-not-allowed disabled:opacity-60"
                :disabled="!selectedConfig || saving"
                @click="saveConfig"
              >
                <UIcon :name="saving ? 'i-lucide-loader-circle' : 'i-lucide-save'" class="size-4" :class="saving ? 'animate-spin' : ''" />
                <span>{{ $t('common.save') }}</span>
              </button>
              <UButton type="button" color="neutral" variant="outline" icon="i-lucide-rotate-ccw" :disabled="!selectedConfig || saving" @click="resetFormFromSelected">
                {{ $t('admin.actions.reset') }}
              </UButton>
            </div>
          </form>
        </section>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ApiError } from '~/composables/useApi'
import type { AdminSiteConfig } from '~/composables/useAdmin'
import { formatDateTime } from '~/utils/format'

type ConfigValueKind = 'string' | 'int' | 'float' | 'boolean' | 'json'

definePageMeta({ layout: 'admin', middleware: 'admin' })

const { t, locale } = useI18n()
const toast = useToast()
const adminApi = useAdmin()

useHead({ title: t('admin.site.configs.title') })

const groups = computed(() => [
  { value: 'tracker', label: t('admin.site.configs.groups.tracker') },
  { value: 'iam', label: t('admin.site.configs.groups.iam') },
  { value: 'catalog', label: t('admin.site.configs.groups.catalog') }
])

const configs = ref<AdminSiteConfig[]>([])
const selectedGroup = ref('tracker')
const selectedConfig = ref<AdminSiteConfig | null>(null)
const pending = ref(false)
const saving = ref(false)
const errorMessage = ref('')
const formError = ref('')

const form = reactive({
  textValue: '',
  booleanValue: false
})

const numberFormatter = computed(() => new Intl.NumberFormat(locale.value))
const currentGroupLabel = computed(() => groups.value.find((item) => item.value === selectedGroup.value)?.label || selectedGroup.value)
const selectedRawValue = computed(() => selectedConfig.value ? selectedConfig.value.value : null)
const selectedKind = computed<ConfigValueKind>(() => selectedConfig.value ? normalizeValueKind(selectedConfig.value.valueType, selectedRawValue.value) : 'string')
const selectedKindLabel = computed(() => t(`admin.site.configs.types.${selectedKind.value}`))

onMounted(loadConfigs)

async function loadConfigs() {
  pending.value = true
  errorMessage.value = ''
  try {
    const data = await adminApi.listSiteConfigs(selectedGroup.value)
    configs.value = (data.configs || []).sort((a, b) => a.key.localeCompare(b.key))
    if (selectedConfig.value) {
      const next = configs.value.find((item) => item.group === selectedConfig.value?.group && item.key === selectedConfig.value?.key)
      if (next) {
        selectConfig(next)
      } else {
        selectConfig(configs.value[0] || null)
      }
    } else {
      selectConfig(configs.value[0] || null)
    }
  } catch (error: unknown) {
    errorMessage.value = error instanceof ApiError ? error.message : t('common.requestFailed')
  } finally {
    pending.value = false
  }
}

function selectGroup(group: string) {
  if (group === selectedGroup.value) return
  selectedGroup.value = group
  selectedConfig.value = null
  formError.value = ''
  loadConfigs()
}

function selectConfig(config: AdminSiteConfig | null) {
  selectedConfig.value = config
  formError.value = ''
  resetFormFromSelected()
}

function resetFormFromSelected() {
  const value = selectedRawValue.value
  if (selectedKind.value === 'boolean') {
    form.booleanValue = Boolean(value)
    form.textValue = ''
  } else {
    form.textValue = formatEditableValue(value)
    form.booleanValue = false
  }
  formError.value = ''
}

function inferValueKind(value: unknown): ConfigValueKind {
  if (typeof value === 'boolean') return 'boolean'
  if (typeof value === 'number') return Number.isInteger(value) ? 'int' : 'float'
  if (value !== null && typeof value === 'object') return 'json'
  return 'string'
}

function normalizeValueKind(valueType: unknown, value: unknown): ConfigValueKind {
  if (valueType === 'string' || valueType === 'int' || valueType === 'float' || valueType === 'boolean' || valueType === 'json') {
    return valueType
  }
  return inferValueKind(value)
}

function formatEditableValue(value: unknown) {
  if (value === null || value === undefined) return ''
  if (typeof value === 'string') return value
  if (typeof value === 'number') return String(value)
  return JSON.stringify(value, null, 2)
}

function previewConfigValue(config: AdminSiteConfig) {
  const value = config.value
  if (value === null || value === undefined || value === '') return '-'
  if (typeof value === 'string') return value
  if (typeof value === 'boolean') return value ? t('admin.site.configs.boolean.true') : t('admin.site.configs.boolean.false')
  return JSON.stringify(value)
}

function getTextFormValue() {
  return String(form.textValue ?? '').trim()
}

function buildSubmitValue() {
  if (!selectedConfig.value) return ''
  if (selectedKind.value === 'boolean') return form.booleanValue
  if (selectedKind.value === 'int') return Number.parseInt(getTextFormValue(), 10)
  if (selectedKind.value === 'float') return Number(getTextFormValue())
  if (selectedKind.value === 'json') return JSON.parse(getTextFormValue() || 'null')
  return getTextFormValue()
}

function validateFormValue() {
  if (selectedKind.value === 'int') {
    const value = getTextFormValue()
    if (!value || !Number.isInteger(Number(value))) {
      formError.value = t('admin.site.configs.form.intInvalid')
      return false
    }
  }
  if (selectedKind.value === 'float') {
    const value = getTextFormValue()
    if (!value || !Number.isFinite(Number(value))) {
      formError.value = t('admin.site.configs.form.numberInvalid')
      return false
    }
  }
  if (selectedKind.value !== 'json') return true
  try {
    JSON.parse(getTextFormValue() || 'null')
    return true
  } catch {
    formError.value = t('admin.site.configs.form.jsonInvalid')
    return false
  }
}

async function saveConfig() {
  if (!selectedConfig.value) return
  if (!validateFormValue()) return

  const group = selectedConfig.value.group
  const key = selectedConfig.value.key
  const value = buildSubmitValue()
  const requestPath = `/api/admin/site/configs/${encodeURIComponent(group)}/${encodeURIComponent(key)}`

  saving.value = true
  formError.value = ''
  try {
    await fetchApi(requestPath, {
      method: 'PUT',
      body: { value }
    })
    toast.add({ title: t('admin.site.configs.saved') })
    await loadConfigs()
  } catch (error: unknown) {
    formError.value = error instanceof ApiError ? error.message : error instanceof Error ? error.message : t('common.requestFailed')
  } finally {
    saving.value = false
  }
}
</script>
