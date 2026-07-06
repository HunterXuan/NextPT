<template>
  <div class="min-h-[calc(100vh-4rem)] bg-slate-50 py-6 dark:bg-slate-950">
    <div class="w-full px-3 sm:px-4 lg:px-5">
      <div class="grid gap-4 xl:grid-cols-[minmax(680px,860px)_minmax(480px,1fr)] 2xl:grid-cols-[minmax(740px,920px)_minmax(520px,1fr)]">
        <section class="overflow-hidden rounded-lg border border-slate-200 bg-white dark:border-slate-800 dark:bg-slate-900">
          <div class="border-b border-slate-200 px-4 py-3 dark:border-slate-800">
            <div class="flex flex-col gap-3 lg:flex-row lg:items-center lg:justify-between">
              <div class="flex items-center gap-2">
                <UIcon name="i-lucide-settings-2" class="size-5 text-indigo-600 dark:text-indigo-300" />
                <h2 class="text-sm font-semibold text-slate-950 dark:text-white">{{ $t('admin.site.configs.list') }}</h2>
                <span class="inline-flex h-6 items-center rounded-md bg-slate-100 px-2 text-xs font-medium text-slate-600 dark:bg-slate-800 dark:text-slate-300">
                  {{ $t('admin.site.configs.stats.total') }} {{ numberFormatter.format(configs.length) }}
                </span>
              </div>
              <div class="flex flex-wrap items-center gap-2">
                <UButton
                  v-for="group in groups"
                  :key="group.value"
                  size="sm"
                  :color="selectedGroup === group.value ? 'primary' : 'neutral'"
                  :variant="selectedGroup === group.value ? 'soft' : 'ghost'"
                  @click="selectGroup(group.value)"
                >
                  {{ group.label }}
                </UButton>
              </div>
            </div>
          </div>

          <div v-if="pending" class="overflow-x-auto">
            <table class="min-w-[780px] w-full table-fixed border-collapse text-left">
              <thead class="bg-slate-50 text-xs font-medium uppercase text-slate-500 dark:bg-slate-950/70 dark:text-slate-400">
                <tr>
                  <th class="w-[24%] border-b border-slate-200 px-3 py-2.5 dark:border-slate-800">{{ $t('admin.site.configs.table.key') }}</th>
                  <th class="w-[28%] border-b border-slate-200 px-3 py-2.5 dark:border-slate-800">{{ $t('admin.site.configs.table.value') }}</th>
                  <th class="w-[34%] border-b border-slate-200 px-3 py-2.5 dark:border-slate-800">{{ $t('admin.site.configs.table.description') }}</th>
                  <th class="w-[14%] border-b border-slate-200 px-3 py-2.5 text-right dark:border-slate-800">{{ $t('admin.site.configs.table.updatedAt') }}</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="index in 6" :key="index" class="border-b border-slate-200 last:border-b-0 dark:border-slate-800">
                  <td class="px-3 py-3"><div class="h-4 w-40 animate-pulse rounded bg-slate-200 dark:bg-slate-800" /></td>
                  <td class="px-3 py-3"><div class="h-5 w-52 animate-pulse rounded bg-slate-100 dark:bg-slate-800/70" /></td>
                  <td class="px-3 py-3"><div class="h-4 w-64 animate-pulse rounded bg-slate-100 dark:bg-slate-800/70" /></td>
                  <td class="px-3 py-3"><div class="ml-auto h-4 w-24 animate-pulse rounded bg-slate-100 dark:bg-slate-800/70" /></td>
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
            <table class="min-w-[780px] w-full table-fixed border-collapse text-left">
              <thead class="bg-slate-50 text-xs font-medium uppercase text-slate-500 dark:bg-slate-950/70 dark:text-slate-400">
                <tr>
                  <th class="w-[24%] border-b border-slate-200 px-3 py-2.5 dark:border-slate-800">{{ $t('admin.site.configs.table.key') }}</th>
                  <th class="w-[28%] border-b border-slate-200 px-3 py-2.5 dark:border-slate-800">{{ $t('admin.site.configs.table.value') }}</th>
                  <th class="w-[34%] border-b border-slate-200 px-3 py-2.5 dark:border-slate-800">{{ $t('admin.site.configs.table.description') }}</th>
                  <th class="w-[14%] border-b border-slate-200 px-3 py-2.5 text-right dark:border-slate-800">{{ $t('admin.site.configs.table.updatedAt') }}</th>
                </tr>
              </thead>
              <tbody>
                <tr
                  v-for="config in configs"
                  :key="config.id || `${config.group}.${config.key}`"
                  class="cursor-pointer border-b border-slate-200 transition-colors last:border-b-0 dark:border-slate-800"
                  :class="isSelectedConfig(config) ? 'bg-indigo-50/70 dark:bg-indigo-950/30' : 'hover:bg-slate-50 dark:hover:bg-slate-950/70'"
                  role="button"
                  tabindex="0"
                  @click="selectConfig(config)"
                  @keydown.enter.prevent="selectConfig(config)"
                  @keydown.space.prevent="selectConfig(config)"
                >
                  <td class="px-3 py-2.5 align-middle">
                    <span class="block truncate text-sm font-semibold text-slate-950 dark:text-white">
                      {{ displayConfigLabel(config) }}
                    </span>
                    <code class="mt-1 block truncate text-xs text-slate-500 dark:text-slate-400">{{ config.group }}.{{ config.key }}</code>
                  </td>
                  <td class="px-3 py-2.5 align-middle">
                    <code class="block truncate rounded-md bg-slate-100 px-2 py-1 text-xs text-slate-700 dark:bg-slate-800 dark:text-slate-200">
                      {{ previewConfigValue(config) }}
                    </code>
                  </td>
                  <td class="px-3 py-2.5 align-middle">
                    <span class="block truncate text-sm text-slate-600 dark:text-slate-300">{{ displayConfigDescription(config) }}</span>
                  </td>
                  <td class="px-3 py-2.5 text-right align-middle text-xs text-slate-600 dark:text-slate-300">
                    {{ formatDateTime(config.updatedAt || config.createdAt, locale) }}
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </section>

        <section class="overflow-hidden rounded-lg border border-slate-200 bg-white dark:border-slate-800 dark:bg-slate-900">
          <div class="border-b border-slate-200 px-4 py-3 dark:border-slate-800">
            <div class="flex items-start justify-between gap-3">
              <div class="min-w-0">
                <h2 class="truncate text-sm font-semibold text-slate-950 dark:text-white">
                  {{ selectedConfig ? selectedConfigLabel : $t('admin.site.configs.form.empty') }}
                </h2>
                <p v-if="selectedConfig" class="mt-1 truncate text-xs text-slate-500 dark:text-slate-400">
                  {{ selectedConfigPath }}
                </p>
              </div>
              <UBadge v-if="selectedConfig" color="primary" variant="soft">
                {{ selectedKindLabel }}
              </UBadge>
            </div>
          </div>

          <div v-if="!selectedConfig" class="flex min-h-80 flex-col items-center justify-center px-4 py-16 text-center">
            <UIcon name="i-lucide-mouse-pointer-2" class="size-9 text-slate-400" />
            <p class="mt-3 text-sm font-medium text-slate-950 dark:text-white">{{ $t('admin.site.configs.form.empty') }}</p>
            <p class="mt-1 max-w-xs text-sm text-slate-500 dark:text-slate-400">{{ $t('admin.site.configs.form.emptyHint') }}</p>
          </div>

          <form v-else class="space-y-4 p-4" @submit.prevent="saveConfig">
            <div class="rounded-md border border-slate-200 bg-slate-50 px-3 py-2 text-sm text-slate-600 dark:border-slate-800 dark:bg-slate-950/70 dark:text-slate-300">
              {{ selectedDescription }}
            </div>

            <dl class="grid gap-2 text-xs sm:grid-cols-2">
              <div class="rounded-md bg-slate-50 px-3 py-2 dark:bg-slate-950/70">
                <dt class="text-slate-400 dark:text-slate-500">{{ $t('admin.site.configs.form.group') }}</dt>
                <dd class="mt-1 font-medium text-slate-700 dark:text-slate-200">{{ currentGroupLabel }}</dd>
              </div>
              <div class="rounded-md bg-slate-50 px-3 py-2 dark:bg-slate-950/70">
                <dt class="text-slate-400 dark:text-slate-500">{{ $t('admin.site.configs.form.updatedAt') }}</dt>
                <dd class="mt-1 font-medium text-slate-700 dark:text-slate-200">{{ selectedUpdatedAtLabel }}</dd>
              </div>
            </dl>

            <label v-if="selectedKind === 'boolean'" class="flex items-center justify-between gap-3 rounded-md border border-slate-200 px-3 py-3 dark:border-slate-800">
              <span>
                <span class="block text-sm font-medium text-slate-950 dark:text-white">{{ $t('admin.site.configs.form.booleanValue') }}</span>
                <span class="mt-1 block text-xs text-slate-500 dark:text-slate-400">{{ previewBooleanValue(form.booleanValue) }}</span>
              </span>
              <button
                type="button"
                role="switch"
                :aria-checked="form.booleanValue"
                class="relative h-6 w-11 rounded-full transition disabled:cursor-not-allowed disabled:opacity-60"
                :class="form.booleanValue ? 'bg-indigo-600' : 'bg-slate-300 dark:bg-slate-700'"
                :disabled="saving"
                @click="form.booleanValue = !form.booleanValue"
              >
                <span
                  class="absolute top-0.5 size-5 rounded-full bg-white shadow-sm transition"
                  :class="form.booleanValue ? 'left-5' : 'left-0.5'"
                />
              </button>
            </label>

            <label v-else-if="selectedUsesRoleSelect" class="block">
              <span class="text-sm font-medium text-slate-700 dark:text-slate-200">{{ $t('admin.site.configs.form.value') }}</span>
              <select
                v-model.number="selectedRoleId"
                class="mt-1 h-10 w-full rounded-md border border-slate-200 bg-white px-3 text-sm text-slate-950 outline-none transition focus:border-indigo-400 focus:ring-2 focus:ring-indigo-100 disabled:cursor-not-allowed disabled:bg-slate-50 disabled:text-slate-400 dark:border-slate-700 dark:bg-slate-950 dark:text-white dark:focus:border-indigo-500 dark:focus:ring-indigo-950 dark:disabled:bg-slate-900 dark:disabled:text-slate-500"
                :disabled="saving || rolesPending || roleOptions.length === 0"
              >
                <option v-if="!hasSelectedRole" :value="selectedRoleId">{{ $t('admin.site.configs.form.unknownRole', { id: selectedRoleId }) }}</option>
                <option v-for="role in roleOptions" :key="role.id" :value="role.id">{{ roleNameWithLevel(role) }}</option>
              </select>
              <p v-if="rolesError" class="mt-1 text-xs text-red-600 dark:text-red-300">{{ rolesError }}</p>
            </label>

            <label v-else-if="selectedKind === 'int' || selectedKind === 'float'" class="block">
              <span class="text-sm font-medium text-slate-700 dark:text-slate-200">{{ $t('admin.site.configs.form.value') }}</span>
              <input
                v-model="form.textValue"
                type="number"
                :step="selectedKind === 'int' ? 1 : 'any'"
                class="mt-1 h-10 w-full rounded-md border border-slate-200 bg-white px-3 text-sm text-slate-950 outline-none transition focus:border-indigo-400 focus:ring-2 focus:ring-indigo-100 dark:border-slate-700 dark:bg-slate-950 dark:text-white dark:focus:border-indigo-500 dark:focus:ring-indigo-950"
                :disabled="saving"
              >
            </label>

            <label v-else-if="selectedKind === 'json'" class="block">
              <span class="text-sm font-medium text-slate-700 dark:text-slate-200">{{ $t('admin.site.configs.form.value') }}</span>
              <textarea
                v-model="form.textValue"
                rows="12"
                class="mt-1 w-full resize-y rounded-md border border-slate-200 bg-white px-3 py-2 font-mono text-sm text-slate-950 outline-none transition focus:border-indigo-400 focus:ring-2 focus:ring-indigo-100 dark:border-slate-700 dark:bg-slate-950 dark:text-white dark:focus:border-indigo-500 dark:focus:ring-indigo-950"
                :disabled="saving"
              />
            </label>

            <label v-else class="block">
              <span class="text-sm font-medium text-slate-700 dark:text-slate-200">{{ $t('admin.site.configs.form.value') }}</span>
              <input
                v-model="form.textValue"
                class="mt-1 h-10 w-full rounded-md border border-slate-200 bg-white px-3 text-sm text-slate-950 outline-none transition focus:border-indigo-400 focus:ring-2 focus:ring-indigo-100 dark:border-slate-700 dark:bg-slate-950 dark:text-white dark:focus:border-indigo-500 dark:focus:ring-indigo-950"
                :disabled="saving"
              >
            </label>

            <p v-if="formError" class="rounded-md border border-red-200 bg-red-50 px-3 py-2 text-sm text-red-700 dark:border-red-900 dark:bg-red-950/40 dark:text-red-200">
              {{ formError }}
            </p>

            <div class="flex justify-end gap-2 border-t border-slate-200 pt-4 dark:border-slate-800">
              <UButton type="button" color="neutral" variant="outline" icon="i-lucide-rotate-ccw" :disabled="saving || !isFormDirty" @click="resetFormFromSelected">
                {{ $t('admin.actions.reset') }}
              </UButton>
              <UButton type="submit" color="primary" icon="i-lucide-save" :loading="saving" :disabled="!canSave">
                {{ $t('common.save') }}
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
import type { AdminIamRole, AdminSiteConfig } from '~/composables/useAdmin'
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
const originalFormSnapshot = ref('')
const roles = ref<AdminIamRole[]>([])
const rolesPending = ref(false)
const rolesError = ref('')

const form = reactive({
  textValue: '',
  booleanValue: false
})

const numberFormatter = computed(() => new Intl.NumberFormat(locale.value))
const currentGroupLabel = computed(() => groups.value.find((item) => item.value === selectedGroup.value)?.label || selectedGroup.value)
const selectedRawValue = computed(() => selectedConfig.value ? selectedConfig.value.value : null)
const selectedKind = computed<ConfigValueKind>(() => selectedConfig.value ? normalizeValueKind(selectedConfig.value.valueType, selectedRawValue.value) : 'string')
const selectedKindLabel = computed(() => t(`admin.site.configs.types.${selectedKind.value}`))
const selectedConfigLabel = computed(() => selectedConfig.value ? displayConfigLabel(selectedConfig.value) : '')
const selectedConfigPath = computed(() => selectedConfig.value ? `${selectedConfig.value.group}.${selectedConfig.value.key}` : '')
const selectedDescription = computed(() => selectedConfig.value ? displayConfigDescription(selectedConfig.value) : '')
const selectedUpdatedAtLabel = computed(() => {
  if (!selectedConfig.value) return '-'
  return formatDateTime(selectedConfig.value.updatedAt || selectedConfig.value.createdAt, locale.value)
})
const isFormDirty = computed(() => Boolean(selectedConfig.value && formSnapshot() !== originalFormSnapshot.value))
const canSave = computed(() => Boolean(selectedConfig.value && isFormDirty.value && !saving.value))
const selectedUsesRoleSelect = computed(() => selectedConfigPath.value === 'iam.default_register_role')
const { roleOptions, roleNameWithLevel } = useAdminIamRoleLevels(roles)
const selectedRoleId = computed({
  get: () => Number.parseInt(getTextFormValue(), 10) || 0,
  set: (value) => {
    form.textValue = String(Number(value || 0))
  }
})
const hasSelectedRole = computed(() => roleOptions.value.some((role) => role.id === selectedRoleId.value))

onMounted(loadConfigs)

async function loadConfigs() {
  pending.value = true
  errorMessage.value = ''
  try {
    const previous = selectedConfig.value
    const data = await adminApi.listSiteConfigs(selectedGroup.value)
    configs.value = (data.configs || []).sort((a, b) => a.key.localeCompare(b.key))
    if (previous) {
      const next = configs.value.find((item) => item.group === previous.group && item.key === previous.key)
      selectConfig(next || configs.value[0] || null)
      return
    }
    selectConfig(configs.value[0] || null)
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
  if (group === 'iam') {
    loadRoles()
  }
}

function selectConfig(config: AdminSiteConfig | null) {
  selectedConfig.value = config
  formError.value = ''
  resetFormFromSelected(true)
  if (selectedUsesRoleSelect.value) {
    loadRoles()
  }
}

async function loadRoles() {
  if (rolesPending.value || roles.value.length > 0) return
  rolesPending.value = true
  rolesError.value = ''
  try {
    const data = await adminApi.listIamRoles()
    roles.value = data.roles || []
  } catch (error: unknown) {
    rolesError.value = error instanceof ApiError ? error.message : error instanceof Error ? error.message : t('common.requestFailed')
  } finally {
    rolesPending.value = false
  }
}

function resetFormFromSelected(updateSnapshot = false) {
  const value = selectedRawValue.value
  if (selectedKind.value === 'boolean') {
    form.booleanValue = Boolean(value)
    form.textValue = ''
  } else {
    form.textValue = formatEditableValue(value)
    form.booleanValue = false
  }
  formError.value = ''
  if (updateSnapshot) {
    originalFormSnapshot.value = selectedConfig.value ? formSnapshot() : ''
  }
}

function isSelectedConfig(config: AdminSiteConfig) {
  return selectedConfig.value?.group === config.group && selectedConfig.value?.key === config.key
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

function previewBooleanValue(value: boolean) {
  return value ? t('admin.site.configs.boolean.true') : t('admin.site.configs.boolean.false')
}

function previewConfigValue(config: AdminSiteConfig) {
  const value = config.value
  if (value === null || value === undefined || value === '') return '-'
  if (typeof value === 'string') return value
  if (typeof value === 'boolean') return previewBooleanValue(value)
  return JSON.stringify(value)
}

function displayConfigLabel(config: AdminSiteConfig) {
  return t(`admin.site.configs.labels.${config.group}.${config.key}`)
}

function displayConfigDescription(config: AdminSiteConfig) {
  return t(`admin.site.configs.descriptions.${config.group}.${config.key}`)
}

function getTextFormValue() {
  return String(form.textValue ?? '').trim()
}

function formSnapshot() {
  return JSON.stringify({
    kind: selectedKind.value,
    textValue: selectedKind.value === 'boolean' ? '' : form.textValue,
    booleanValue: selectedKind.value === 'boolean' ? form.booleanValue : false
  })
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
  if (!selectedConfig.value || !canSave.value) return
  formError.value = ''
  if (!validateFormValue()) return

  const group = selectedConfig.value.group
  const key = selectedConfig.value.key
  const value = buildSubmitValue()

  saving.value = true
  try {
    await adminApi.updateSiteConfig(group, key, value)
    toast.add({ title: t('admin.site.configs.saved') })
    await loadConfigs()
  } catch (error: unknown) {
    formError.value = error instanceof ApiError ? error.message : error instanceof Error ? error.message : t('common.requestFailed')
  } finally {
    saving.value = false
  }
}
</script>
