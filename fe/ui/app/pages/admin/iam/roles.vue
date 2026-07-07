<template>
  <div class="min-h-[calc(100vh-4rem)] bg-slate-50 py-6 dark:bg-slate-950">
    <div class="w-full px-3 sm:px-4 lg:px-5">
      <div class="grid gap-4 xl:grid-cols-[minmax(620px,820px)_minmax(520px,1fr)] 2xl:grid-cols-[minmax(660px,880px)_minmax(560px,1fr)]">
        <section class="overflow-hidden rounded-lg border border-slate-200 bg-white dark:border-slate-800 dark:bg-slate-900">
          <div class="border-b border-slate-200 px-4 py-3 dark:border-slate-800">
            <div class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
              <div class="flex items-center gap-2">
                <UIcon name="i-lucide-shield-check" class="size-5 text-emerald-600 dark:text-emerald-300" />
                <h2 class="text-sm font-semibold text-slate-950 dark:text-white">{{ $t('admin.iam.roles.list') }}</h2>
              </div>
              <div class="flex flex-wrap items-center gap-2 text-xs">
                <span class="inline-flex h-7 items-center gap-1.5 rounded-md bg-slate-100 px-2.5 font-medium text-slate-600 dark:bg-slate-800 dark:text-slate-300">
                  <span>{{ $t('admin.iam.roles.stats.total') }}</span>
                  <span class="font-semibold text-slate-950 dark:text-white">{{ numberFormatter.format(roles.length) }}</span>
                </span>
                <span class="inline-flex h-7 items-center gap-1.5 rounded-md bg-emerald-50 px-2.5 font-medium text-emerald-700 dark:bg-emerald-950/40 dark:text-emerald-300">
                  <span>{{ $t('admin.iam.roles.stats.staff') }}</span>
                  <span class="font-semibold">{{ numberFormatter.format(staffCount) }}</span>
                </span>
              </div>
            </div>
          </div>

          <div v-if="pending" class="overflow-x-auto">
            <table class="min-w-[820px] w-full table-fixed border-collapse text-left">
              <thead class="bg-slate-50 text-xs font-medium uppercase text-slate-500 dark:bg-slate-950/70 dark:text-slate-400">
                <tr>
                  <th class="w-[34%] border-b border-slate-200 px-4 py-3 dark:border-slate-800">{{ $t('admin.iam.roles.table.name') }}</th>
                  <th class="w-[14%] border-b border-slate-200 px-4 py-3 text-right dark:border-slate-800">{{ $t('admin.iam.roles.table.level') }}</th>
                  <th class="w-[14%] border-b border-slate-200 px-4 py-3 dark:border-slate-800">{{ $t('admin.iam.roles.table.staff') }}</th>
                  <th class="w-[18%] border-b border-slate-200 px-4 py-3 text-right dark:border-slate-800">{{ $t('admin.iam.roles.table.permissions') }}</th>
                  <th class="w-[20%] border-b border-slate-200 px-4 py-3 dark:border-slate-800">{{ $t('admin.iam.roles.table.updatedAt') }}</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="index in 5" :key="index" class="border-b border-slate-200 last:border-b-0 dark:border-slate-800">
                  <td class="px-4 py-4"><div class="h-4 w-44 animate-pulse rounded bg-slate-200 dark:bg-slate-800" /></td>
                  <td class="px-4 py-4"><div class="ml-auto h-4 w-12 animate-pulse rounded bg-slate-100 dark:bg-slate-800/70" /></td>
                  <td class="px-4 py-4"><div class="h-4 w-16 animate-pulse rounded bg-slate-100 dark:bg-slate-800/70" /></td>
                  <td class="px-4 py-4"><div class="ml-auto h-4 w-16 animate-pulse rounded bg-slate-100 dark:bg-slate-800/70" /></td>
                  <td class="px-4 py-4"><div class="h-4 w-28 animate-pulse rounded bg-slate-100 dark:bg-slate-800/70" /></td>
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
            <table class="min-w-[820px] w-full table-fixed border-collapse text-left">
              <thead class="bg-slate-50 text-xs font-medium uppercase text-slate-500 dark:bg-slate-950/70 dark:text-slate-400">
                <tr>
                  <th class="w-[34%] border-b border-slate-200 px-4 py-3 dark:border-slate-800">{{ $t('admin.iam.roles.table.name') }}</th>
                  <th class="w-[14%] border-b border-slate-200 px-4 py-3 text-right dark:border-slate-800">{{ $t('admin.iam.roles.table.level') }}</th>
                  <th class="w-[14%] border-b border-slate-200 px-4 py-3 dark:border-slate-800">{{ $t('admin.iam.roles.table.staff') }}</th>
                  <th class="w-[18%] border-b border-slate-200 px-4 py-3 text-right dark:border-slate-800">{{ $t('admin.iam.roles.table.permissions') }}</th>
                  <th class="w-[20%] border-b border-slate-200 px-4 py-3 dark:border-slate-800">{{ $t('admin.iam.roles.table.updatedAt') }}</th>
                </tr>
              </thead>
              <tbody>
                <tr
                  v-for="role in roles"
                  :key="role.id"
                  class="cursor-pointer border-b border-slate-200 transition-colors last:border-b-0 dark:border-slate-800"
                  :class="selectedId === role.id ? 'bg-emerald-50/70 dark:bg-emerald-950/30' : 'hover:bg-slate-50 dark:hover:bg-slate-950/70'"
                  role="button"
                  tabindex="0"
                  @click="selectRole(role)"
                  @keydown.enter.prevent="selectRole(role)"
                  @keydown.space.prevent="selectRole(role)"
                >
                  <td class="px-4 py-3 align-middle">
                    <div class="block max-w-full text-left">
                      <span class="block truncate text-sm font-semibold text-slate-950 hover:text-emerald-700 dark:text-white dark:hover:text-emerald-300">
                        {{ roleName(role) }}
                      </span>
                      <span class="mt-1 block truncate text-xs text-slate-500 dark:text-slate-400">#{{ role.id }}</span>
                    </div>
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
              <div class="flex h-8 min-w-20 items-center justify-end">
                <UButton v-if="selectedId" color="primary" variant="soft" size="sm" icon="i-lucide-plus" @click="startCreate">
                  {{ $t('admin.actions.new') }}
                </UButton>
              </div>
            </div>
          </div>

          <form class="space-y-5 p-4" @submit.prevent="saveRole">
            <div>
              <div class="flex items-center justify-between gap-3 border-b border-slate-200 pb-2 dark:border-slate-800">
                <span class="text-sm font-semibold text-slate-950 dark:text-white">{{ $t('admin.iam.roles.form.name') }}</span>
                <span class="text-xs text-slate-500 dark:text-slate-400">
                  {{ $t('admin.iam.roles.form.primaryLocale', { locale: primaryLocaleLabel }) }}
                </span>
              </div>
              <div class="mt-3 grid gap-3 sm:grid-cols-2 2xl:grid-cols-3">
                <label v-for="item in localeOptions" :key="item.code" class="block">
                  <span class="flex items-center justify-between gap-2 text-sm font-medium text-slate-700 dark:text-slate-200">
                    <span>{{ item.name }}</span>
                    <span class="text-xs font-normal text-slate-400 dark:text-slate-500">{{ item.code }}</span>
                  </span>
                  <input
                    v-model="form.names[item.code]"
                    class="mt-1 h-10 w-full rounded-md border border-slate-200 bg-white px-3 text-sm text-slate-950 outline-none transition focus:border-emerald-400 focus:ring-2 focus:ring-emerald-100 dark:border-slate-700 dark:bg-slate-950 dark:text-white dark:focus:border-emerald-500 dark:focus:ring-emerald-950"
                    :disabled="saving"
                  >
                </label>
              </div>
            </div>

            <div class="grid gap-3 sm:grid-cols-[minmax(0,1fr)_180px]">
              <label class="block">
                <span class="text-sm font-medium text-slate-700 dark:text-slate-200">{{ $t('admin.iam.roles.form.level') }}</span>
                <input v-model.number="form.level" type="number" class="mt-1 h-10 w-full rounded-md border border-slate-200 bg-white px-3 text-sm text-slate-950 outline-none transition focus:border-emerald-400 focus:ring-2 focus:ring-emerald-100 dark:border-slate-700 dark:bg-slate-950 dark:text-white dark:focus:border-emerald-500 dark:focus:ring-emerald-950" :disabled="saving">
              </label>

              <label class="block">
                <span class="text-sm font-medium text-slate-700 dark:text-slate-200">{{ $t('admin.iam.roles.form.isStaff') }}</span>
                <span
                  class="mt-1 flex h-10 items-center justify-between gap-3 rounded-md border px-3 text-sm font-medium transition"
                  :class="form.isStaff ? 'border-emerald-200 bg-emerald-50 text-emerald-700 dark:border-emerald-900/70 dark:bg-emerald-950/30 dark:text-emerald-300' : 'border-slate-200 bg-white text-slate-600 dark:border-slate-700 dark:bg-slate-950 dark:text-slate-300'"
                >
                  <span>{{ form.isStaff ? $t('admin.iam.roles.staff.yes') : $t('admin.iam.roles.staff.no') }}</span>
                  <input
                    v-model="form.isStaff"
                    type="checkbox"
                    class="peer sr-only"
                    :disabled="saving"
                  >
                  <span class="relative h-5 w-9 rounded-full bg-slate-200 transition-colors after:absolute after:left-0.5 after:top-0.5 after:size-4 after:rounded-full after:bg-white after:shadow-sm after:transition-transform peer-checked:bg-emerald-500 peer-checked:after:translate-x-4 peer-disabled:opacity-60 dark:bg-slate-700" />
                </span>
              </label>
            </div>

            <div>
              <div class="mb-2 flex items-center justify-between gap-3">
                <span class="text-sm font-medium text-slate-700 dark:text-slate-200">{{ $t('admin.iam.roles.form.permissions') }}</span>
                <span class="text-xs text-slate-500 dark:text-slate-400">{{ numberFormatter.format(form.permissions.length) }}</span>
              </div>
              <div v-if="permissions.length" class="max-h-72 overflow-y-auto rounded-md border border-slate-200 dark:border-slate-700">
                <section
                  v-for="group in permissionGroups"
                  :key="group.key"
                  class="border-b border-slate-200 last:border-b-0 dark:border-slate-800"
                >
                  <div class="sticky top-0 z-10 flex items-center justify-between gap-3 border-b border-slate-200 bg-slate-50 px-3 py-2 dark:border-slate-800 dark:bg-slate-900">
                    <label class="flex min-w-0 items-center gap-2">
                      <input
                        type="checkbox"
                        class="size-4 rounded border-slate-300 text-emerald-600 focus:ring-emerald-500 dark:border-slate-600"
                        :checked="isPermissionGroupSelected(group)"
                        :indeterminate="isPermissionGroupPartial(group)"
                        :disabled="saving"
                        @change="togglePermissionGroup(group)"
                      >
                      <span class="truncate text-xs font-semibold uppercase tracking-wide text-slate-600 dark:text-slate-300">{{ group.label }}</span>
                    </label>
                    <span class="shrink-0 rounded bg-white px-1.5 py-0.5 text-xs font-medium text-slate-500 dark:bg-slate-950 dark:text-slate-400">
                      {{ permissionGroupSelectedCount(group) }}/{{ group.permissions.length }}
                    </span>
                  </div>
                  <label
                    v-for="permission in group.permissions"
                    :key="permission"
                    class="flex items-center gap-2 border-b border-slate-100 px-3 py-2 text-xs last:border-b-0 hover:bg-slate-50 dark:border-slate-800 dark:hover:bg-slate-950/70"
                  >
                    <input
                      type="checkbox"
                      class="size-4 rounded border-slate-300 text-emerald-600 focus:ring-emerald-500 dark:border-slate-600"
                      :checked="isPermissionSelected(permission)"
                      :disabled="saving"
                      @change="togglePermission(permission)"
                    >
                    <span class="w-20 shrink-0 rounded bg-slate-100 px-1.5 py-0.5 text-center font-mono text-[11px] font-medium text-slate-600 dark:bg-slate-800 dark:text-slate-300">
                      {{ permissionAction(permission) }}
                    </span>
                    <code class="min-w-0 truncate text-slate-700 dark:text-slate-200">{{ permission }}</code>
                  </label>
                </section>
              </div>
              <div v-else class="rounded-md border border-dashed border-slate-300 px-3 py-4 text-sm text-slate-500 dark:border-slate-700 dark:text-slate-400">
                {{ $t('admin.iam.roles.form.noPermissions') }}
              </div>
            </div>

            <div>
              <div class="mb-2 flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
                <span class="text-sm font-medium text-slate-700 dark:text-slate-200">{{ $t('admin.iam.roles.form.rules') }}</span>
                <UButton
                  type="button"
                  color="neutral"
                  variant="outline"
                  size="xs"
                  icon="i-lucide-align-left"
                  :disabled="saving || !rulesJsonState.valid"
                  @click="formatRulesJson"
                >
                  {{ $t('admin.iam.roles.form.rulesFormat') }}
                </UButton>
              </div>
              <AdminJsonEditor
                v-model="form.rulesJson"
                :valid="rulesJsonState.valid"
                :error-message="rulesJsonState.message"
                :title="$t('admin.iam.roles.form.rules')"
                :placeholder="$t('admin.iam.roles.form.rulesPlaceholder')"
                :hint="$t('admin.iam.roles.form.rulesEditorHint')"
                :valid-text="$t('admin.iam.roles.form.rulesValid')"
                :invalid-text="$t('admin.iam.roles.form.rulesInvalidShort')"
                :fullscreen-text="$t('admin.iam.roles.form.rulesFullscreen')"
                :close-text="$t('common.close')"
                :disabled="saving"
              />
              <p v-if="!rulesJsonState.valid" class="mt-2 text-xs text-red-600 dark:text-red-300">
                {{ rulesJsonState.message }}
              </p>
            </div>

            <p v-if="formError" class="rounded-md border border-red-200 bg-red-50 px-3 py-2 text-sm text-red-700 dark:border-red-900 dark:bg-red-950/40 dark:text-red-200">
              {{ formError }}
            </p>

            <div class="flex flex-col gap-3 border-t border-slate-200 pt-4 sm:flex-row sm:items-center sm:justify-between dark:border-slate-800">
              <div class="flex justify-start">
                <UPopover
                  v-if="selectedRole"
                  :content="{ side: 'top', align: 'start', sideOffset: 8 }"
                  :ui="{ content: 'w-72 p-3' }"
                >
                  <UButton
                    color="error"
                    variant="soft"
                    icon="i-lucide-trash-2"
                    :loading="deletingId === selectedRole.id"
                    :disabled="deletingId !== null || saving"
                  >
                    {{ $t('admin.actions.delete') }}
                  </UButton>

                  <template #content="{ close }">
                    <div class="space-y-3">
                      <p class="text-sm font-medium text-slate-950 dark:text-white">
                        {{ $t('admin.actions.confirmDeleteTitle') }}
                      </p>
                      <p class="text-xs text-slate-500 dark:text-slate-400">
                        {{ $t('admin.actions.deleteIrreversible') }}
                      </p>
                      <div class="flex justify-end gap-2">
                        <UButton color="neutral" variant="ghost" size="xs" type="button" @click="close()">
                          {{ $t('common.cancel') }}
                        </UButton>
                        <UButton color="error" size="xs" type="button" icon="i-lucide-trash-2" :loading="deletingId === selectedRole.id" :disabled="deletingId !== null" @click="deleteRole(selectedRole, close)">
                          {{ $t('admin.actions.delete') }}
                        </UButton>
                      </div>
                    </div>
                  </template>
                </UPopover>
              </div>
              <div class="flex flex-col gap-2 sm:flex-row sm:justify-end">
                <UButton type="button" color="neutral" variant="outline" icon="i-lucide-rotate-ccw" :disabled="saving || !isFormDirty" @click="resetFormChanges">
                  {{ selectedId ? $t('admin.actions.discardChanges') : $t('admin.actions.reset') }}
                </UButton>
                <UButton type="submit" color="primary" icon="i-lucide-save" :loading="saving" :disabled="!canSubmit">
                  {{ $t('common.save') }}
                </UButton>
              </div>
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

interface LocaleOption {
  code: string
  name: string
}

interface RulesJsonState {
  valid: boolean
  rules: Record<string, unknown>
  message: string
}

interface PermissionGroup {
  key: string
  label: string
  permissions: string[]
}

definePageMeta({ layout: 'admin', middleware: 'admin' })

const { t, locale, locales } = useI18n()
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
const originalFormSnapshot = ref('')

const form = reactive({
  names: {} as Record<string, string>,
  level: 10,
  isStaff: false,
  permissions: [] as string[],
  rulesJson: '{}'
})

const localeOptions = computed<LocaleOption[]>(() => {
  return (locales.value as Array<string | { code?: string, language?: string, name?: string }>)
    .map((item) => {
      if (typeof item === 'string') return { code: item, name: item }
      const code = String(item.code || item.language || '')
      return { code, name: String(item.name || code) }
    })
    .filter((item) => item.code)
})
const primaryLocaleCode = computed(() => localeOptions.value[0]?.code || locale.value)
const primaryLocaleLabel = computed(() => {
  const primary = localeOptions.value.find((item) => item.code === primaryLocaleCode.value)
  return primary ? `${primary.name} (${primary.code})` : primaryLocaleCode.value
})
const numberFormatter = computed(() => new Intl.NumberFormat(locale.value))
const staffCount = computed(() => roles.value.filter((item) => item.isStaff).length)
const selectedRole = computed(() => roles.value.find((item) => item.id === selectedId.value) || null)
const isFormDirty = computed(() => formSnapshot() !== originalFormSnapshot.value)
const rulesJsonState = computed(() => parseRulesJsonText(form.rulesJson))
const canSubmit = computed(() => Boolean(nameValue(primaryLocaleCode.value).trim() && isFormDirty.value && !saving.value && rulesJsonState.value.valid))
const permissionGroups = computed<PermissionGroup[]>(() => {
  const groups = new Map<string, PermissionGroup>()
  for (const permission of permissions.value) {
    const groupMeta = permissionGroupMeta(permission)
    const group = groups.get(groupMeta.key)
    if (group) {
      group.permissions.push(permission)
    } else {
      groups.set(groupMeta.key, {
        key: groupMeta.key,
        label: groupMeta.label,
        permissions: [permission]
      })
    }
  }

  return Array.from(groups.values())
    .map((group) => ({
      ...group,
      permissions: group.permissions.slice().sort(comparePermissions)
    }))
    .sort((a, b) => a.label.localeCompare(b.label))
})

watch(localeOptions, () => {
  ensureFormNameLocales()
}, { immediate: true })

onMounted(() => {
  originalFormSnapshot.value = formSnapshot()
  loadAll()
})

async function loadAll() {
  pending.value = true
  errorMessage.value = ''
  try {
    const data = await adminApi.listIamRoles()
    roles.value = data.roles || []
    const roleId = roles.value[0]?.id || 1
    const permissionData = await adminApi.listIamPermissions(roleId)
    permissions.value = (permissionData.permissions || []).slice().sort()
    if (selectedId.value) {
      const next = roles.value.find((item) => item.id === selectedId.value)
      if (next) {
        selectRole(next)
      } else {
        startCreate()
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
    return {}
  }
  const record = value as Record<string, unknown>
  return record
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

function permissionAction(permission: string) {
  const separatorIndex = permission.indexOf(':')
  return separatorIndex > 0 ? permission.slice(0, separatorIndex) : permission
}

function permissionGroupMeta(permission: string) {
  const parts = permission.split(':')
  const target = parts.length >= 2 ? parts[1] : permission
  const [domain, ...resourceParts] = target.split('/')
  const resource = resourceParts.join('/')
  if (!domain) {
    return { key: permission, label: permission }
  }
  if (!resource) {
    return { key: domain, label: domain }
  }
  return { key: `${domain}/${resource}`, label: `${domain} / ${resource}` }
}

function comparePermissions(a: string, b: string) {
  const groupA = permissionGroupMeta(a).label
  const groupB = permissionGroupMeta(b).label
  if (groupA !== groupB) return groupA.localeCompare(groupB)
  const actionA = permissionAction(a)
  const actionB = permissionAction(b)
  if (actionA !== actionB) return actionA.localeCompare(actionB)
  return a.localeCompare(b)
}

function sortSelectedPermissions(values: string[]) {
  const permissionOrder = new Map(permissions.value.map((permission, index) => [permission, index]))
  return Array.from(new Set(values)).sort((a, b) => {
    const orderA = permissionOrder.get(a)
    const orderB = permissionOrder.get(b)
    if (orderA !== undefined && orderB !== undefined) return orderA - orderB
    if (orderA !== undefined) return -1
    if (orderB !== undefined) return 1
    return comparePermissions(a, b)
  })
}

function isPermissionSelected(permission: string) {
  return form.permissions.includes(permission)
}

function togglePermission(permission: string) {
  if (saving.value) return
  const selected = new Set(form.permissions)
  if (selected.has(permission)) {
    selected.delete(permission)
  } else {
    selected.add(permission)
  }
  form.permissions = sortSelectedPermissions(Array.from(selected))
}

function permissionGroupSelectedCount(group: PermissionGroup) {
  return group.permissions.filter((permission) => form.permissions.includes(permission)).length
}

function isPermissionGroupSelected(group: PermissionGroup) {
  return group.permissions.length > 0 && permissionGroupSelectedCount(group) === group.permissions.length
}

function isPermissionGroupPartial(group: PermissionGroup) {
  const count = permissionGroupSelectedCount(group)
  return count > 0 && count < group.permissions.length
}

function togglePermissionGroup(group: PermissionGroup) {
  if (saving.value) return
  const selected = new Set(form.permissions)
  const shouldSelect = !isPermissionGroupSelected(group)
  for (const permission of group.permissions) {
    if (shouldSelect) {
      selected.add(permission)
    } else {
      selected.delete(permission)
    }
  }
  form.permissions = sortSelectedPermissions(Array.from(selected))
}

function selectRole(role: AdminIamRole) {
  selectedId.value = role.id
  setFormNames(role.nameI18N)
  form.level = role.level
  form.isStaff = role.isStaff
  form.permissions = rolePermissions(role)
  form.rulesJson = JSON.stringify(normalizeRules(role.rules), null, 2)
  formError.value = ''
  originalFormSnapshot.value = formSnapshot()
}

function startCreate() {
  selectedId.value = null
  setFormNames()
  form.level = 10
  form.isStaff = false
  form.permissions = []
  form.rulesJson = '{}'
  formError.value = ''
  originalFormSnapshot.value = formSnapshot()
}

function resetFormChanges() {
  const role = selectedRole.value
  if (role) {
    selectRole(role)
    return
  }
  startCreate()
}

function parseRulesJson() {
  const state = rulesJsonState.value
  if (!state.valid) {
    throw new Error(state.message || t('admin.iam.roles.form.rulesJsonInvalid'))
  }
  return state.rules
}

function parseRulesJsonText(raw: string): RulesJsonState {
  const text = raw.trim()
  if (!text) return { valid: true, rules: {}, message: '' }

  try {
    const parsed = JSON.parse(text) as unknown
    if (!parsed || typeof parsed !== 'object' || Array.isArray(parsed)) {
      return { valid: false, rules: {}, message: t('admin.iam.roles.form.rulesInvalid') }
    }
    return { valid: true, rules: parsed as Record<string, unknown>, message: '' }
  } catch {
    return { valid: false, rules: {}, message: t('admin.iam.roles.form.rulesJsonInvalid') }
  }
}

function formatRulesJson() {
  if (!rulesJsonState.value.valid) return
  form.rulesJson = JSON.stringify(rulesJsonState.value.rules, null, 2)
}

function ensureFormNameLocales() {
  const supportedCodes = new Set(localeOptions.value.map((item) => item.code))
  for (const code of Object.keys(form.names)) {
    if (!supportedCodes.has(code)) {
      delete form.names[code]
    }
  }
  for (const item of localeOptions.value) {
    if (!(item.code in form.names)) {
      form.names[item.code] = ''
    }
  }
}

function setFormNames(values?: Record<string, unknown> | null) {
  const supportedCodes = new Set(localeOptions.value.map((item) => item.code))
  for (const code of Object.keys(form.names)) {
    if (!supportedCodes.has(code)) {
      delete form.names[code]
    }
  }
  for (const item of localeOptions.value) {
    const value = values?.[item.code]
    form.names[item.code] = typeof value === 'string' ? value : ''
  }
}

function nameValue(code: string) {
  return form.names[code] || ''
}

function formSnapshot() {
  return JSON.stringify({
    names: localeOptions.value.reduce<Record<string, string>>((names, item) => {
      names[item.code] = nameValue(item.code).trim()
      return names
    }, {}),
    level: Number(form.level || 0),
    isStaff: Boolean(form.isStaff),
    permissions: form.permissions.slice().sort(),
    rulesJson: form.rulesJson.trim()
  })
}

function buildInput(): AdminIamRoleInput {
  const primaryName = nameValue(primaryLocaleCode.value).trim()
  return {
    level: Number(form.level || 0),
    nameI18N: localeOptions.value.reduce<Record<string, string>>((names, item) => {
      names[item.code] = nameValue(item.code).trim() || primaryName
      return names
    }, {}),
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
    const editingId = selectedId.value
    if (selectedId.value) {
      await adminApi.updateIamRole(selectedId.value, input)
    } else {
      await adminApi.createIamRole(input)
    }
    toast.add({ title: t('admin.iam.roles.saved') })
    await loadAll()
    if (editingId) {
      const updatedRole = roles.value.find((item) => item.id === editingId)
      if (updatedRole) {
        selectRole(updatedRole)
      } else {
        startCreate()
      }
    } else {
      startCreate()
    }
  } catch (error: unknown) {
    formError.value = error instanceof ApiError ? error.message : error instanceof Error ? error.message : t('common.requestFailed')
  } finally {
    saving.value = false
  }
}

async function deleteRole(role: AdminIamRole, close?: () => void) {
  deletingId.value = role.id
  try {
    await adminApi.deleteIamRole(role.id)
    if (selectedId.value === role.id) {
      startCreate()
    }
    toast.add({ title: t('admin.iam.roles.deleted') })
    close?.()
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
