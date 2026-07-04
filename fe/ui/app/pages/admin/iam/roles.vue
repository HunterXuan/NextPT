<template>
  <div class="min-h-[calc(100vh-4rem)] bg-slate-50 py-6 dark:bg-slate-950">
    <div class="w-full px-3 sm:px-4 lg:px-5">
      <div class="mb-4 grid grid-cols-2 gap-2 sm:flex sm:items-center">
        <div class="rounded-lg border border-slate-200 bg-white px-3 py-2 dark:border-slate-800 dark:bg-slate-900">
          <p class="text-xs text-slate-500 dark:text-slate-400">{{ $t('admin.iam.roles.stats.total') }}</p>
          <p class="mt-1 text-lg font-semibold text-slate-950 dark:text-white">{{ numberFormatter.format(roles.length) }}</p>
        </div>
        <div class="rounded-lg border border-slate-200 bg-white px-3 py-2 dark:border-slate-800 dark:bg-slate-900">
          <p class="text-xs text-slate-500 dark:text-slate-400">{{ $t('admin.iam.roles.stats.staff') }}</p>
          <p class="mt-1 text-lg font-semibold text-emerald-700 dark:text-emerald-300">{{ numberFormatter.format(staffCount) }}</p>
        </div>
      </div>

      <div class="grid gap-4 xl:grid-cols-[minmax(0,1fr)_440px]">
        <section class="overflow-hidden rounded-lg border border-slate-200 bg-white dark:border-slate-800 dark:bg-slate-900">
          <div class="flex items-center justify-between gap-3 border-b border-slate-200 px-4 py-3 dark:border-slate-800">
            <h2 class="text-sm font-semibold text-slate-950 dark:text-white">{{ $t('admin.iam.roles.list') }}</h2>
          </div>

          <div v-if="pending" class="overflow-x-auto">
            <table class="min-w-[920px] w-full table-fixed border-collapse text-left">
              <thead class="bg-slate-50 text-xs font-medium uppercase text-slate-500 dark:bg-slate-950/70 dark:text-slate-400">
                <tr>
                  <th class="w-[28%] border-b border-slate-200 px-4 py-3 dark:border-slate-800">{{ $t('admin.iam.roles.table.name') }}</th>
                  <th class="w-[12%] border-b border-slate-200 px-4 py-3 text-right dark:border-slate-800">{{ $t('admin.iam.roles.table.level') }}</th>
                  <th class="w-[12%] border-b border-slate-200 px-4 py-3 dark:border-slate-800">{{ $t('admin.iam.roles.table.staff') }}</th>
                  <th class="w-[18%] border-b border-slate-200 px-4 py-3 text-right dark:border-slate-800">{{ $t('admin.iam.roles.table.permissions') }}</th>
                  <th class="w-[14%] border-b border-slate-200 px-4 py-3 dark:border-slate-800">{{ $t('admin.iam.roles.table.updatedAt') }}</th>
                  <th class="w-[16%] border-b border-slate-200 px-4 py-3 text-right dark:border-slate-800">{{ $t('admin.iam.roles.table.actions') }}</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="index in 5" :key="index" class="border-b border-slate-200 last:border-b-0 dark:border-slate-800">
                  <td class="px-4 py-4"><div class="h-4 w-44 animate-pulse rounded bg-slate-200 dark:bg-slate-800" /></td>
                  <td class="px-4 py-4"><div class="ml-auto h-4 w-12 animate-pulse rounded bg-slate-100 dark:bg-slate-800/70" /></td>
                  <td class="px-4 py-4"><div class="h-4 w-16 animate-pulse rounded bg-slate-100 dark:bg-slate-800/70" /></td>
                  <td class="px-4 py-4"><div class="ml-auto h-4 w-16 animate-pulse rounded bg-slate-100 dark:bg-slate-800/70" /></td>
                  <td class="px-4 py-4"><div class="h-4 w-28 animate-pulse rounded bg-slate-100 dark:bg-slate-800/70" /></td>
                  <td class="px-4 py-4"><div class="ml-auto h-4 w-24 animate-pulse rounded bg-slate-100 dark:bg-slate-800/70" /></td>
                </tr>
              </tbody>
            </table>
          </div>

          <div v-else-if="errorMessage" class="flex flex-col items-center justify-center px-4 py-16 text-center">
            <UIcon name="i-lucide-circle-alert" class="size-9 text-red-500" />
            <p class="mt-3 text-sm font-medium text-slate-950 dark:text-white">{{ errorMessage }}</p>
          </div>

          <div v-else-if="roles.length === 0" class="flex flex-col items-center justify-center px-4 py-16 text-center">
            <UIcon name="i-lucide-inbox" class="size-9 text-slate-400" />
            <p class="mt-3 text-sm font-medium text-slate-950 dark:text-white">{{ $t('admin.iam.roles.empty') }}</p>
          </div>

          <div v-else class="overflow-x-auto">
            <table class="min-w-[920px] w-full table-fixed border-collapse text-left">
              <thead class="bg-slate-50 text-xs font-medium uppercase text-slate-500 dark:bg-slate-950/70 dark:text-slate-400">
                <tr>
                  <th class="w-[28%] border-b border-slate-200 px-4 py-3 dark:border-slate-800">{{ $t('admin.iam.roles.table.name') }}</th>
                  <th class="w-[12%] border-b border-slate-200 px-4 py-3 text-right dark:border-slate-800">{{ $t('admin.iam.roles.table.level') }}</th>
                  <th class="w-[12%] border-b border-slate-200 px-4 py-3 dark:border-slate-800">{{ $t('admin.iam.roles.table.staff') }}</th>
                  <th class="w-[18%] border-b border-slate-200 px-4 py-3 text-right dark:border-slate-800">{{ $t('admin.iam.roles.table.permissions') }}</th>
                  <th class="w-[14%] border-b border-slate-200 px-4 py-3 dark:border-slate-800">{{ $t('admin.iam.roles.table.updatedAt') }}</th>
                  <th class="w-[16%] border-b border-slate-200 px-4 py-3 text-right dark:border-slate-800">{{ $t('admin.iam.roles.table.actions') }}</th>
                </tr>
              </thead>
              <tbody>
                <tr
                  v-for="role in roles"
                  :key="role.id"
                  class="border-b border-slate-200 transition-colors last:border-b-0 dark:border-slate-800"
                  :class="selectedId === role.id ? 'bg-emerald-50/70 dark:bg-emerald-950/30' : 'hover:bg-slate-50 dark:hover:bg-slate-950/70'"
                >
                  <td class="px-4 py-3 align-middle">
                    <button type="button" class="block max-w-full text-left" @click="selectRole(role)">
                      <span class="block truncate text-sm font-semibold text-slate-950 hover:text-emerald-700 dark:text-white dark:hover:text-emerald-300">
                        {{ roleName(role) }}
                      </span>
                      <span class="mt-1 block truncate text-xs text-slate-500 dark:text-slate-400">#{{ role.id }}</span>
                    </button>
                  </td>
                  <td class="px-4 py-3 text-right align-middle text-sm text-slate-600 dark:text-slate-300">
                    {{ numberFormatter.format(role.level) }}
                  </td>
                  <td class="px-4 py-3 align-middle">
                    <UBadge :color="role.isStaff ? 'success' : 'neutral'" variant="soft">
                      {{ role.isStaff ? $t('admin.iam.roles.staff.yes') : $t('admin.iam.roles.staff.no') }}
                    </UBadge>
                  </td>
                  <td class="px-4 py-3 text-right align-middle text-sm text-slate-600 dark:text-slate-300">
                    {{ numberFormatter.format(rolePermissions(role).length) }}
                  </td>
                  <td class="px-4 py-3 align-middle text-sm text-slate-600 dark:text-slate-300">
                    {{ formatDateTime(role.updatedAt || role.createdAt, locale) }}
                  </td>
                  <td class="px-4 py-3 align-middle">
                    <div class="flex items-center justify-end gap-2">
                      <UButton color="neutral" variant="outline" size="sm" icon="i-lucide-pencil" @click="selectRole(role)">
                        {{ $t('admin.actions.edit') }}
                      </UButton>
                      <UButton color="error" variant="soft" size="sm" icon="i-lucide-trash-2" :loading="deletingId === role.id" @click="deleteRole(role)">
                        {{ $t('admin.actions.delete') }}
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
            <div class="flex items-center justify-between gap-3">
              <h2 class="text-sm font-semibold text-slate-950 dark:text-white">
                {{ selectedId ? $t('admin.iam.roles.form.edit') : $t('admin.iam.roles.form.create') }}
              </h2>
              <UButton v-if="selectedId" color="neutral" variant="ghost" size="sm" icon="i-lucide-plus" @click="resetForm">
                {{ $t('admin.actions.new') }}
              </UButton>
            </div>
          </div>

          <form class="space-y-4 p-4" @submit.prevent="saveRole">
            <label v-for="field in nameFields" :key="field.key" class="block">
              <span class="text-sm font-medium text-slate-700 dark:text-slate-200">{{ field.label }}</span>
              <input v-model="form[field.key]" class="mt-1 h-10 w-full rounded-md border border-slate-200 bg-white px-3 text-sm text-slate-950 outline-none transition focus:border-emerald-400 focus:ring-2 focus:ring-emerald-100 dark:border-slate-700 dark:bg-slate-950 dark:text-white dark:focus:border-emerald-500 dark:focus:ring-emerald-950" :disabled="saving">
            </label>

            <div class="grid grid-cols-[1fr_auto] gap-3">
              <label class="block">
                <span class="text-sm font-medium text-slate-700 dark:text-slate-200">{{ $t('admin.iam.roles.form.level') }}</span>
                <input v-model.number="form.level" type="number" class="mt-1 h-10 w-full rounded-md border border-slate-200 bg-white px-3 text-sm text-slate-950 outline-none transition focus:border-emerald-400 focus:ring-2 focus:ring-emerald-100 dark:border-slate-700 dark:bg-slate-950 dark:text-white dark:focus:border-emerald-500 dark:focus:ring-emerald-950" :disabled="saving">
              </label>

              <label class="flex items-end gap-2 pb-2 text-sm font-medium text-slate-700 dark:text-slate-200">
                <input v-model="form.isStaff" type="checkbox" class="size-4 rounded border-slate-300 text-emerald-600 focus:ring-emerald-500 dark:border-slate-600" :disabled="saving">
                <span>{{ $t('admin.iam.roles.form.isStaff') }}</span>
              </label>
            </div>

            <div>
              <div class="mb-2 flex items-center justify-between gap-3">
                <span class="text-sm font-medium text-slate-700 dark:text-slate-200">{{ $t('admin.iam.roles.form.permissions') }}</span>
                <span class="text-xs text-slate-500 dark:text-slate-400">{{ numberFormatter.format(form.permissions.length) }}</span>
              </div>
              <div v-if="permissions.length" class="max-h-60 overflow-y-auto rounded-md border border-slate-200 dark:border-slate-700">
                <label
                  v-for="permission in permissions"
                  :key="permission"
                  class="flex items-center gap-2 border-b border-slate-100 px-3 py-2 text-xs last:border-b-0 dark:border-slate-800"
                >
                  <input v-model="form.permissions" type="checkbox" class="size-4 rounded border-slate-300 text-emerald-600 focus:ring-emerald-500 dark:border-slate-600" :value="permission" :disabled="saving">
                  <code class="min-w-0 truncate text-slate-700 dark:text-slate-200">{{ permission }}</code>
                </label>
              </div>
              <div v-else class="rounded-md border border-dashed border-slate-300 px-3 py-4 text-sm text-slate-500 dark:border-slate-700 dark:text-slate-400">
                {{ $t('admin.iam.roles.form.noPermissions') }}
              </div>
            </div>

            <div class="flex gap-2">
              <input v-model="customPermission" class="h-10 min-w-0 flex-1 rounded-md border border-slate-200 bg-white px-3 font-mono text-sm text-slate-950 outline-none transition focus:border-emerald-400 focus:ring-2 focus:ring-emerald-100 dark:border-slate-700 dark:bg-slate-950 dark:text-white dark:focus:border-emerald-500 dark:focus:ring-emerald-950" :placeholder="$t('admin.iam.roles.form.customPermission')" :disabled="saving">
              <UButton type="button" color="neutral" variant="outline" icon="i-lucide-plus" :disabled="!customPermission.trim() || saving" @click="addCustomPermission">
                {{ $t('admin.iam.roles.form.addPermission') }}
              </UButton>
            </div>

            <label class="block">
              <span class="text-sm font-medium text-slate-700 dark:text-slate-200">{{ $t('admin.iam.roles.form.rules') }}</span>
              <textarea v-model="form.rulesJson" rows="6" class="mt-1 w-full resize-none rounded-md border border-slate-200 bg-white px-3 py-2 font-mono text-xs text-slate-950 outline-none transition focus:border-emerald-400 focus:ring-2 focus:ring-emerald-100 dark:border-slate-700 dark:bg-slate-950 dark:text-white dark:focus:border-emerald-500 dark:focus:ring-emerald-950" :disabled="saving" />
            </label>

            <p v-if="formError" class="rounded-md border border-red-200 bg-red-50 px-3 py-2 text-sm text-red-700 dark:border-red-900 dark:bg-red-950/40 dark:text-red-200">
              {{ formError }}
            </p>

            <div class="flex flex-col gap-2 sm:flex-row">
              <UButton type="submit" color="primary" icon="i-lucide-save" :loading="saving" :disabled="!canSubmit">
                {{ $t('common.save') }}
              </UButton>
              <UButton type="button" color="neutral" variant="outline" icon="i-lucide-rotate-ccw" :disabled="saving" @click="resetForm">
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
import type { AdminIamRole, AdminIamRoleInput } from '~/composables/useAdmin'
import { formatDateTime, localizeI18nName } from '~/utils/format'

definePageMeta({ layout: 'admin', middleware: 'admin' })

const { t, locale } = useI18n()
const toast = useToast()
const adminApi = useAdmin()

useHead({ title: t('admin.iam.roles.title') })

const roles = ref<AdminIamRole[]>([])
const permissions = ref<string[]>([])
const pending = ref(false)
const saving = ref(false)
const deletingId = ref<number | null>(null)
const selectedId = ref<number | null>(null)
const errorMessage = ref('')
const formError = ref('')
const customPermission = ref('')

const form = reactive({
  nameZhCN: '',
  nameZhTW: '',
  nameEnUS: '',
  level: 10,
  isStaff: false,
  permissions: [] as string[],
  rulesJson: '{}'
})

const numberFormatter = computed(() => new Intl.NumberFormat(locale.value))
const staffCount = computed(() => roles.value.filter((item) => item.isStaff).length)
const canSubmit = computed(() => Boolean(form.nameZhCN.trim() && !saving.value))
const nameFields = computed(() => [
  { key: 'nameZhCN' as const, label: t('admin.iam.roles.form.nameZhCN') },
  { key: 'nameZhTW' as const, label: t('admin.iam.roles.form.nameZhTW') },
  { key: 'nameEnUS' as const, label: t('admin.iam.roles.form.nameEnUS') }
])

onMounted(loadAll)

async function loadAll() {
  pending.value = true
  errorMessage.value = ''
  try {
    const data = await adminApi.listIamRoles()
    roles.value = (data.roles || []).sort((a, b) => b.level - a.level || a.id - b.id)
    const roleId = roles.value[0]?.id || 1
    const permissionData = await adminApi.listIamPermissions(roleId)
    permissions.value = (permissionData.permissions || []).slice().sort()
    if (selectedId.value) {
      const next = roles.value.find((item) => item.id === selectedId.value)
      if (next) {
        selectRole(next)
      } else {
        resetForm()
      }
    }
  } catch (error: unknown) {
    errorMessage.value = error instanceof ApiError ? error.message : t('common.requestFailed')
  } finally {
    pending.value = false
  }
}

function normalizeI18n(value: unknown) {
  if (!value || typeof value !== 'object' || Array.isArray(value)) {
    return { 'zh-CN': '', 'zh-TW': '', 'en-US': '' }
  }
  const record = value as Record<string, unknown>
  return {
    'zh-CN': String(record['zh-CN'] || ''),
    'zh-TW': String(record['zh-TW'] || ''),
    'en-US': String(record['en-US'] || '')
  }
}

function normalizeRules(value: unknown) {
  if (!value || typeof value !== 'object' || Array.isArray(value)) return {}
  return value as Record<string, unknown>
}

function normalizePermissions(value: unknown) {
  if (!Array.isArray(value)) return []
  return value
    .map((item) => String(item || '').trim())
    .filter((item, index, list) => item && list.indexOf(item) === index)
    .sort()
}

function roleName(role: AdminIamRole) {
  return localizeI18nName(normalizeI18n(role.nameI18N), locale.value, `#${role.id}`)
}

function rolePermissions(role: AdminIamRole) {
  return normalizePermissions(role.permissions)
}

function selectRole(role: AdminIamRole) {
  const name = normalizeI18n(role.nameI18N)
  form.nameZhCN = name['zh-CN']
  form.nameZhTW = name['zh-TW']
  form.nameEnUS = name['en-US']
  form.level = role.level
  form.isStaff = role.isStaff
  form.permissions = rolePermissions(role)
  form.rulesJson = JSON.stringify(normalizeRules(role.rules), null, 2)
  selectedId.value = role.id
  formError.value = ''
  customPermission.value = ''
}

function resetForm() {
  selectedId.value = null
  form.nameZhCN = ''
  form.nameZhTW = ''
  form.nameEnUS = ''
  form.level = 10
  form.isStaff = false
  form.permissions = []
  form.rulesJson = '{}'
  formError.value = ''
  customPermission.value = ''
}

function addCustomPermission() {
  const value = customPermission.value.trim()
  if (!value || form.permissions.includes(value)) return
  form.permissions = [...form.permissions, value].sort()
  customPermission.value = ''
}

function parseRulesJson() {
  const raw = form.rulesJson.trim()
  if (!raw) return {}
  const parsed = JSON.parse(raw) as unknown
  if (!parsed || typeof parsed !== 'object' || Array.isArray(parsed)) {
    throw new Error(t('admin.iam.roles.form.rulesInvalid'))
  }
  return parsed as Record<string, unknown>
}

function buildInput(): AdminIamRoleInput {
  const zhCN = form.nameZhCN.trim()
  const zhTW = form.nameZhTW.trim() || zhCN
  const enUS = form.nameEnUS.trim() || zhCN
  return {
    level: Number(form.level || 0),
    nameI18N: {
      'zh-CN': zhCN,
      'zh-TW': zhTW,
      'en-US': enUS
    },
    rules: parseRulesJson(),
    permissions: form.permissions.slice().sort(),
    isStaff: form.isStaff
  }
}

async function saveRole() {
  if (!canSubmit.value) return

  saving.value = true
  formError.value = ''
  try {
    const input = buildInput()
    if (selectedId.value) {
      await adminApi.updateIamRole(selectedId.value, input)
    } else {
      await adminApi.createIamRole(input)
    }
    toast.add({ title: t('admin.iam.roles.saved') })
    resetForm()
    await loadAll()
  } catch (error: unknown) {
    formError.value = error instanceof ApiError ? error.message : error instanceof Error ? error.message : t('common.requestFailed')
  } finally {
    saving.value = false
  }
}

async function deleteRole(role: AdminIamRole) {
  if (!window.confirm(t('admin.iam.roles.confirmDelete', { name: roleName(role) }))) return

  deletingId.value = role.id
  try {
    await adminApi.deleteIamRole(role.id)
    if (selectedId.value === role.id) {
      resetForm()
    }
    toast.add({ title: t('admin.iam.roles.deleted') })
    await loadAll()
  } catch (error: unknown) {
    toast.add({
      color: 'error',
      title: error instanceof ApiError ? error.message : t('common.requestFailed')
    })
  } finally {
    deletingId.value = null
  }
}
</script>
